package repositories

import (
	"Inquiro/models"
	"context"
	"database/sql"

	"github.com/google/uuid"
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
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users: &UserRepository{DB: db},
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
