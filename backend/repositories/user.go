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
	"go.uber.org/zap"
)

type UserRepository struct {
	DB     *sql.DB
	logger *zap.SugaredLogger
}

var (
	QueryTimeOut         = 10 * time.Second
	ErrUserNotFound      = errors.New("user not found")
	ErrDuplicateEmail    = errors.New("email already registered")
	ErrDuplicateUsername = errors.New("duplicate username")
	InvitationExpiryTime = 50 * time.Minute
)

func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, first_name, last_name, provider, provider_id, password, email,is_active, is_verified FROM users WHERE email = $1`
	err := u.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Provider, &user.ProviderID, &user.Password.Hash, &user.Email, &user.IsActive, &user.IsVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			u.logger.Warnw("user does not exist", "error :", err.Error())
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
			u.logger.Warnw("user does not exist", "error :", err.Error())
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
			u.logger.Warnw("duplicate email", "error :", errString)
			return ErrDuplicateEmail
		case strings.Contains(errString, `"users_username_key"`):
			u.logger.Warnw("duplicate username", "error :", errString)
			return ErrDuplicateUsername
		case errors.Is(err, sql.ErrNoRows):
			u.logger.Errorw("Failed to insert the user", errString)
			return fmt.Errorf("UserRepository.Create failed: INSERT did not return a row: %w", err)
		default:
			u.logger.Errorw("Failed to insert the user", errString)
			return fmt.Errorf("UserRepository.Create failed: %w", err)
		}
	}
	return nil
}

func (u *UserRepository) createInvitation(tx *sql.Tx, ctx context.Context, userId uuid.UUID, token string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	row := tx.QueryRowContext(ctx, "INSERT INTO user_invitation (id,user_id, token,expiry) VALUES ($1, $2,$3,$4) RETURNING id,created_at", uuid.New(), userId, token, time.Now().Add(InvitationExpiryTime))
	if row.Err() != nil {
		u.logger.Errorw("insertion to user_invitation failed", "error :", row.Err().Error())
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
	query := `SELECT u.id , u.username , u.email , u.created_at , u.is_active FROM users u JOIN user_invitation
	ui ON u.id = ui.user_id WHERE ui.token = $1 AND ui.expiry > $2`
	user := &models.User{}
	hash := sha256.Sum256([]byte(token))
	hashedToken := hex.EncodeToString(hash[:])
	err := tx.QueryRowContext(ctx, query, hashedToken, time.Now()).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.IsActive)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			u.logger.Errorw("user does not exist", "error :", err.Error())
			return nil, ErrUserNotFound
		default:
			u.logger.Errorw("user extraction failed", "error :", err.Error())
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

func (u *UserRepository) update(tx *sql.Tx, ctx context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	query := `UPDATE users SET username = $1 , email = $2 , is_verified = $3 WHERE id = $4`

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsVerified, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	row := u.DB.QueryRowContext(ctx, "SELECT id, username, first_name, last_name, is_active , is_verified, password, email FROM user WHERE email = $1", email)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.IsActive, &user.IsVerified, &user.Password, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) deleteInvitation(tx *sql.Tx, ctx context.Context, userId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	query := `DELETE FROM user_invitation WHERE user_id = $1`

	_, err := tx.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	row := u.DB.QueryRowContext(ctx, "SELECT id, username, first_name, last_name, is_active , is_verified, email, role_id FROM user WHERE id = $1", id)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.IsActive, &user.IsVerified, &user.Password, &user.Email, &user.RoleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
