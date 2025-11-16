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

type MentorRepository struct {
	DB *sql.DB
}

func (u *MentorRepository) FindByEmail(ctx context.Context, email string) (*models.Mentor, error) {
	user := &models.Mentor{}
	query := `SELECT id, username, first_name, last_name, provider, provider_id, password, email,is_active, is_verified FROM Mentor WHERE email = $1`
	err := u.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Provider, &user.ProviderID, &user.Password.Hash, &user.Email, &user.IsActive, &user.IsVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (u *MentorRepository) FindByUsername(ctx context.Context, userName string) (*models.Mentor, error) {
	row := u.DB.QueryRowContext(ctx, "SELECT * FROM user WHERE username = $1", userName)
	user := &models.Mentor{}
	err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Provider, &user.ProviderID, &user.Password, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (u *MentorRepository) create(tx *sql.Tx, ctx context.Context, user *models.Mentor) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	row := tx.QueryRowContext(ctx, "INSERT INTO user (id,username,first_name,last_name,provider,provider_id,password,email) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id,created_at , updated_at", user.ID, user.Username, user.FirstName, user.LastName, user.Provider, user.ProviderID, user.Password.Hash, user.Email)
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
			return fmt.Errorf("MentorRepository.Create failed: INSERT did not return a row: %w", err)
		default:
			// For all other errors (e.g., connection issue, SQL syntax), return a wrapped error
			return fmt.Errorf("MentorRepository.Create failed: %w", err)
		}
	}

	// Check if the context was cancelled after the query executed but before processing
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return nil
}

func (u *MentorRepository) createInvitation(tx *sql.Tx, ctx context.Context, userId uuid.UUID, token string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	row := tx.QueryRowContext(ctx, "INSERT INTO user_invitation (id,user_id, token,expiry) VALUES ($1, $2,$3,$4) RETURNING id,created_at", uuid.New(), userId, token, time.Now().Add(InvitationExpiryTime))
	if row.Err() != nil {
		return row.Err()
	}
	return nil
}

func (u *MentorRepository) CreateAndInvite(ctx context.Context, token string, user *models.Mentor) error {
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

func (u *MentorRepository) getUserFromToken(tx *sql.Tx, ctx context.Context, token string) (*models.Mentor, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	query := `SELECT u.id , u.username , u.email , u.created_at , u.is_active FROM user u JOIN user_invitation
	ui ON u.id = ui.user_id WHERE ui.token = $1 AND ui.expiry > $2`
	user := &models.Mentor{}
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

func (u *MentorRepository) Activate(ctx context.Context, token string) error {
	return WithTx(u.DB, ctx, func(tx *sql.Tx) error {
		user, err := u.getUserFromToken(tx, ctx, token)
		if err != nil {
			return err
		}
		user.IsVerified = true
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

func (u *MentorRepository) update(tx *sql.Tx, ctx context.Context, user *models.Mentor) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	query := `UPDATE user SET username = $1 , email = $2 , is_verified = $3 WHERE id = $4`

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsVerified, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *MentorRepository) GetByEmail(ctx context.Context, email string) (*models.Mentor, error) {
	row := u.DB.QueryRowContext(ctx, "SELECT id, username, first_name, last_name, is_active , is_verified, password, email FROM user WHERE email = $1", email)
	user := &models.Mentor{}
	err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.IsActive, &user.IsVerified, &user.Password, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (u *MentorRepository) deleteInvitation(tx *sql.Tx, ctx context.Context, userId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	query := `DELETE FROM user_invitation WHERE user_id = $1`

	_, err := tx.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}
	return nil
}

func (u *MentorRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Mentor, error) {
	row := u.DB.QueryRowContext(ctx, "SELECT id, username, first_name, last_name, is_active , is_verified, email, role_id FROM user WHERE id = $1", id)
	user := &models.Mentor{}
	err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.IsActive, &user.IsVerified, &user.Password, &user.Email, &user.RoleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
