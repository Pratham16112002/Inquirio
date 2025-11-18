package services

import (
	"Inquiro/models"
	"Inquiro/repositories"
	"context"
	"errors"

	"go.uber.org/zap"
)

type MentorServices struct {
	repo   repositories.Storage
	logger *zap.SugaredLogger
}

func (m MentorServices) CheckMentorUsernameExists(ctx context.Context, username string) bool {
	if _, err := m.repo.Mentor.FindByUsername(ctx, username); err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return false
		}
	}
	return true
}

func (m MentorServices) CheckMentorEmailExists(ctx context.Context, email string) bool {
	if _, err := m.repo.Mentor.FindByEmail(ctx, email); err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return false
		}
	}
	return true
}

func (m MentorServices) GetMentorByEmail(ctx context.Context, email string) (*models.Mentor, error) {
	return m.repo.Mentor.FindByEmail(ctx, email)
}

func (m MentorServices) AuthenticateMentorPassword(ctx context.Context, mentor *models.Mentor, pass *models.PasswordType) error {
	if err := mentor.Password.Compare(*pass.Text); err != nil {
		m.logger.Warnw("Incorrect credentials", "error : ", err.Error())
		return err
	}
	return nil
}

func (m MentorServices) ActivateMentor(ctx context.Context, token string) error {
	return m.repo.Mentor.Activate(ctx, token)
}

func (m MentorServices) RegisterMentor(ctx context.Context, mentor *models.Mentor, token string) error {
	return m.repo.Mentor.CreateAndInvite(ctx, token, mentor)
}
