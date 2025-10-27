package auth

import (
	"Inquiro/models"
	"Inquiro/repositories"
	"context"

	"github.com/alexedwards/scs/v2"
)

type LocalAuth struct {
	store    repositories.Storage
	sessions *scs.SessionManager
}

func NewLocalAuth(repo repositories.Storage, sessions *scs.SessionManager) *LocalAuth {
	return &LocalAuth{
		store:    repo,
		sessions: sessions,
	}
}

func (l *LocalAuth) EmailPasswordAuthenticate(ctx context.Context, email string, password string) (*models.User, error) {
	user, err := l.store.Users.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	l.sessions.Put(ctx, "user", user.ID)
	return user, nil
}

func (l *LocalAuth) LogOut(ctx context.Context) error {
	err := l.sessions.Clear(ctx)
	return err
}
