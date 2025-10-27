package repositories

import (
	"Inquiro/models"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UserRepository struct {
	DB *sql.DB
}

var (
	QueryTimeOut         = 10 * time.Second
	ErrUserNotFound      = errors.New("user not found")
	ErrDuplicateEmail    = errors.New("email already registered")
	ErrDuplicateUsername = errors.New("duplicate username")
)

func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	row := u.DB.QueryRowContext(ctx, "SELECT * FROM users WHERE email = $1", email)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Provider, &user.ProviderID, &user.Password, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (u *UserRepository) FindByUsername(ctx context.Context, userName string) (*models.User, error) {
	row := u.DB.QueryRowContext(ctx, "SELECT * FROM users WHERE username = $1", userName)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Provider, &user.ProviderID, &user.Password, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) create(tx *sql.Tx, ctx context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	row := tx.QueryRowContext(ctx, "INSERT INTO users (id,username,first_name,last_name,provider,provider_id,password,email) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id,created_at , updated_at", user.ID, user.Username, user.FirstName, user.LastName, user.Provider, user.ProviderID, user.Password.Hash, user.Email)
	err := row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		// If the query failed, check the error for specific database violation messages
		errString := err.Error()

		switch {
		// Check for specific PostgreSQL duplicate key constraints
		case strings.Contains(errString, `"users_email_key"`):
			return ErrDuplicateEmail
		case strings.Contains(errString, `"users_username_key"`):
			return ErrDuplicateUsername
		case errors.Is(err, sql.ErrNoRows):
			// This case is unlikely for a RETURNING clause but is good practice.
			return fmt.Errorf("UserRepository.Create failed: INSERT did not return a row: %w", err)
		default:
			// For all other errors (e.g., connection issue, SQL syntax), return a wrapped error
			return fmt.Errorf("UserRepository.Create failed: %w", err)
		}
	}

	// Check if the context was cancelled after the query executed but before processing
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return nil
}

func (u *UserRepository) createInvitation(tx *sql.Tx, ctx context.Context, userId uuid.UUID, token string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	row := tx.QueryRowContext(ctx, "INSERT INTO user_invitations (id,user_id, token) VALUES ($1, $2,$3) RETURNING id,created_at", uuid.New(), userId, token)
	if row.Err() != nil {
		return row.Err()
	}
	return nil
}

func (u *UserRepository) CreateAndInvite(ctx context.Context, token string, user *models.User) error {
	return WithTx(u.DB, ctx, func(tx *sql.Tx) error {
		// create user
		if err := u.create(tx, ctx, user); err != nil {
			return err
		}
		// create user invitation
		if err := u.createInvitation(tx, ctx, user.ID, token); err != nil {
			return err
		}
		return nil
	})
}

func (u *UserRepository) getUserFromToken(tx *sql.Tx, ctx context.Context, token string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	query := `SELECT u.id , u.username , u.email , u.created_at , u.is_active FROM users u JOIN user_invitations
	ui ON u.id = ui.user_id WHERE ui.token = $1 AND ui.expiry > $2`
	user := &models.User{}
	hash := sha256.Sum256([]byte(token))
	hashedToken := hex.EncodeToString(hash[:])
	err := tx.QueryRowContext(ctx, query, hashedToken, time.Now()).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.IsActive)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (u *UserRepository) Activate(ctx context.Context, token string) error {
	return WithTx(u.DB, ctx, func(tx *sql.Tx) error {
		user, err := u.getUserFromToken(tx, ctx, token)
		if err != nil {
			return err
		}
		user.IsActive = true
		err = u.update(tx, ctx, user)
		if err != nil {
			return err
		}
		if err := u.deleteInvitation(tx, ctx, user.ID); err != nil {
			return err
		}
		return nil

	})
}

func (u *UserRepository) update(tx *sql.Tx, ctx context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	query := `UPDATE users SET username = $1 , email = $2 , is_active = $3 WHERE id = $4`

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) deleteInvitation(tx *sql.Tx, ctx context.Context, userId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	query := `DELETE FROM user_invitations WHERE user_id = $1`

	_, err := tx.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}
	return nil
}
