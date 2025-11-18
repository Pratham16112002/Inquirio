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
		CheckUsernameExists(ctx context.Context, username string) bool
		CheckEmailExists(ctx context.Context, email string) bool
		RegisterUser(ctx context.Context, user *models.User, token string) error
		ActivateUser(ctx context.Context, token string) error
		GetUserByEmail(ctx context.Context, email string) (*models.User, error)
		AuthenticatePassword(ctx context.Context, user *models.User, pass *models.PasswordType) error
	}
	MentorServices interface {
		CheckMentorUsernameExists(ctx context.Context, username string) bool
		CheckMentorEmailExists(ctx context.Context, email string) bool
		GetMentorByEmail(ctx context.Context, email string) (*models.Mentor, error)
		AuthenticateMentorPassword(ctx context.Context, mentor *models.Mentor, pass *models.PasswordType) error
		ActivateMentor(ctx context.Context, token string) error
		RegisterMentor(ctx context.Context, mentor *models.Mentor, token string) error
	}
}

func NewService(repo repositories.Storage, logger *zap.SugaredLogger, mailer mailer.Client) Service {
	return Service{
		UserServices: UserServices{
			repo:   repo,
			logger: logger,
		},
		MentorServices: MentorServices{
			repo:   repo,
			logger: logger,
		},
	}
}
