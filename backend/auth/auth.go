package auth

import (
	"Inquiro/models"
	"Inquiro/repositories"
	"context"

	"github.com/alexedwards/scs/v2"
)

type Auth struct {
	LocalAuth interface {
		EmailPasswordAuthenticate(ctx context.Context, email string, password string) (*models.User, error)
		LogOut(ctx context.Context) error
	}
}

func NewAuth(repo repositories.Storage, sessions *scs.SessionManager) Auth {
	return Auth{
		LocalAuth: &LocalAuth{store: repo, sessions: sessions},
	}
}
