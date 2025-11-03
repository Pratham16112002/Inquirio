package services

import (
	"Inquiro/models"
	"Inquiro/repositories"
	"Inquiro/utils/mailer"
	"context"

	"go.uber.org/zap"
)

type Service struct {
	UserServices interface {
		CheckUsernameExists(ctx context.Context, username string) (bool, error)
		CheckEmailExists(ctx context.Context, email string) (bool, error)
		RegisterUser(ctx context.Context, user *models.User, token string) error
		ActivateUser(ctx context.Context, token string) error
		GetUserByEmail(ctx context.Context, email string) (*models.User, error)
		AuthenticatePassword(ctx context.Context, user *models.User, pass *models.PasswordType) error
	}
}

func NewService(repo repositories.Storage, logger *zap.SugaredLogger, mailer mailer.Client) Service {
	return Service{
		UserServices: UserServices{
			repo:   repo,
			logger: logger,
		},
	}
}
