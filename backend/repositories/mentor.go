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

type MentorRepository struct {
	DB     *sql.DB
	logger *zap.SugaredLogger
}

func (u *MentorRepository) FindByEmail(ctx context.Context, email string) (*models.Mentor, error) {
	mentor := &models.Mentor{}
	query := `SELECT id, username, first_name, last_name, email, password , is_active, is_verified , experience_years, bio FROM mentor WHERE email = $1`
	err := u.DB.QueryRowContext(ctx, query, email).Scan(&mentor.ID, &mentor.Username, &mentor.FirstName, &mentor.LastName, &mentor.Email, &mentor.Password.Hash, &mentor.IsActive, &mentor.IsVerified, &mentor.ExperienceYears)
	if err != nil {
		if err == sql.ErrNoRows {
			u.logger.Warnw("user does not exist with this email", "error :", err.Error())
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return mentor, nil
}

func (u *MentorRepository) FindByUsername(ctx context.Context, username string) (*models.Mentor, error) {
	row := u.DB.QueryRowContext(ctx, "SELECT  id, username, first_name, last_name, email, password , is_active, is_verified , experience_years, bio FROM mentor WHERE username = $1", username)
	mentor := &models.Mentor{}
	err := row.Scan(&mentor.ID, &mentor.Username, &mentor.FirstName, &mentor.LastName, &mentor.Email, &mentor.Password.Hash, &mentor.IsActive, &mentor.IsVerified, &mentor.ExperienceYears)
	if err != nil {
		if err == sql.ErrNoRows {
			u.logger.Warnw("user does not exist with this username", "error :", err.Error())
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return mentor, nil
}

func (u *MentorRepository) create(tx *sql.Tx, ctx context.Context, mentor *models.Mentor) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	query := "INSERT INTO mentor (username,first_name,last_name,provider,provider_id,password,email,experience_years,bio) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id,created_at , updated_at"
	row := tx.QueryRowContext(ctx, query, mentor.Username, mentor.FirstName, mentor.LastName, mentor.Provider, mentor.ProviderID, mentor.Password.Hash, mentor.Email, mentor.ExperienceYears, mentor.Bio)
	err := row.Scan(&mentor.ID, &mentor.CreatedAt, &mentor.UpdatedAt)

	if err != nil {
		// If the query failed, check the error for specific database violation messages
		errString := err.Error()

		switch {
		// Check for specific PostgreSQL duplicate key constraints
		case strings.Contains(errString, `"mentor_email_key"`):
			u.logger.Warnw("duplicate email", "error :", errString)
			return ErrDuplicateEmail
		case strings.Contains(errString, `"mentor_username_key"`):
			u.logger.Warnw("duplicate username", "error :", errString)
			return ErrDuplicateUsername
		case errors.Is(err, sql.ErrNoRows):
			// This case is unlikely for a RETURNING clause but is good practice.
			u.logger.Errorw("Failed to insert the user", errString)
			return fmt.Errorf("MentorRepository.Create failed: INSERT did not return a row: %w", err)
		default:
			// For all other errors (e.g., connection issue, SQL syntax), return a wrapped error
			u.logger.Errorw("Failed to insert the user", errString)
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

func (u *MentorRepository) CreateAndInvite(ctx context.Context, token string, mentor *models.Mentor) error {
	return WithTx(u.DB, ctx, func(tx *sql.Tx) error {
		// create mentor
		if err := u.create(tx, ctx, mentor); err != nil {
			return err
		}
		// create mentor invitation
		if err := u.createInvitation(tx, ctx, mentor.ID, token); err != nil {
			return err
		}
		return nil
	})
}

func (u *MentorRepository) getMentorFromToken(tx *sql.Tx, ctx context.Context, token string) (*models.Mentor, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()
	query := `SELECT u.id , u.username , u.email , u.created_at , u.is_active FROM mentor u JOIN user_invitation
	ui ON u.id = ui.user_id WHERE ui.token = $1 AND ui.expiry > $2`
	user := &models.Mentor{}
	hash := sha256.Sum256([]byte(token))
	hashedToken := hex.EncodeToString(hash[:])
	err := tx.QueryRowContext(ctx, query, hashedToken, time.Now()).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.IsActive)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			u.logger.Errorw("mentor does not exist", "error :", err.Error())
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (u *MentorRepository) Activate(ctx context.Context, token string) error {
	return WithTx(u.DB, ctx, func(tx *sql.Tx) error {
		user, err := u.getMentorFromToken(tx, ctx, token)
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

func (u *MentorRepository) update(tx *sql.Tx, ctx context.Context, mentor *models.Mentor) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	query := `UPDATE mentor SET username = $1 , email = $2 , is_verified = $3 WHERE id = $4`

	_, err := tx.ExecContext(ctx, query, mentor.Username, mentor.Email, mentor.IsVerified, mentor.ID)
	if err != nil {
		return err
	}
	return nil
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
