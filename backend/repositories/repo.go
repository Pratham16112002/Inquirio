package repositories

import (
	"Inquiro/models"
	"context"
	"database/sql"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Storage struct {
	Users interface {
		FindByEmail(ctx context.Context, email string) (*models.User, error)
		FindByUsername(ctx context.Context, username string) (*models.User, error)
		CreateAndInvite(ctx context.Context, token string, user *models.User) error
		create(tx *sql.Tx, ctx context.Context, user *models.User) error
		createInvitation(tx *sql.Tx, ctx context.Context, userId uuid.UUID, token string) error
		Activate(ctx context.Context, token string) error
		getUserFromToken(tx *sql.Tx, ctx context.Context, token string) (*models.User, error)
		update(tx *sql.Tx, ctx context.Context, user *models.User) error
		deleteInvitation(tx *sql.Tx, ctx context.Context, userId uuid.UUID) error
		GetByEmail(ctx context.Context, email string) (*models.User, error)
		GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	}
	Mentor interface {
		FindByEmail(ctx context.Context, email string) (*models.Mentor, error)
		FindByUsername(ctx context.Context, username string) (*models.Mentor, error)
		CreateAndInvite(ctx context.Context, token string, user *models.Mentor) error
		create(tx *sql.Tx, ctx context.Context, mentor *models.Mentor) error
		createInvitation(tx *sql.Tx, ctx context.Context, userId uuid.UUID, token string) error
		Activate(ctx context.Context, token string) error
		getMentorFromToken(tx *sql.Tx, ctx context.Context, token string) (*models.Mentor, error)
		update(tx *sql.Tx, ctx context.Context, mentor *models.Mentor) error
		deleteInvitation(tx *sql.Tx, ctx context.Context, userId uuid.UUID) error
		GetByID(ctx context.Context, id uuid.UUID) (*models.Mentor, error)
	}
	Role interface {
		GetRoleByID(ctx context.Context, id int) (models.Role, error)
	}
}

func NewStorage(db *sql.DB, logger *zap.SugaredLogger) Storage {
	return Storage{
		Users: &UserRepository{DB: db,
			logger: logger},
		Mentor: &MentorRepository{DB: db,
			logger: logger},
		Role: &RoleRepository{DB: db,
			logger: logger},
	}
}

func WithTx(db *sql.DB, ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
