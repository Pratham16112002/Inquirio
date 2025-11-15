package services

import (
	"Inquiro/models"
	"Inquiro/repositories"
	"context"
	"errors"

	"go.uber.org/zap"
)

type UserServices struct {
	repo   repositories.Storage
	logger *zap.SugaredLogger
}

var (
	ErrUserNameNotAvailable = errors.New("username not available")
)

func (u UserServices) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	if _, err := u.repo.Users.FindByUsername(ctx, username); err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return false, nil
		}
	}
	return true, ErrUserNameNotAvailable
}

func (u UserServices) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	if _, err := u.repo.Users.FindByEmail(ctx, email); err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return false, nil
		}
	}
	return true, ErrUserNameNotAvailable
}

func (u UserServices) ActivateUser(ctx context.Context, token string) error {
	return u.repo.Users.Activate(ctx, token)
}

func (u UserServices) RegisterUser(ctx context.Context, user *models.User, token string) error {
	return u.repo.Users.CreateAndInvite(ctx, token, user)
}

func (u UserServices) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := u.repo.Users.FindByEmail(ctx, email)
	if err != nil {
		u.logger.Warnw("User does not exist with this credentials", "error : ", err.Error())
		return nil, err
	}
	return user, nil
}

func (u UserServices) AuthenticatePassword(ctx context.Context, user *models.User, pass *models.PasswordType) error {
	if err := user.Password.Compare(*pass.Text); err != nil {
		u.logger.Warnw("Incorrect credentials", "error : ", err.Error())
		return err
	}
	return nil
}
