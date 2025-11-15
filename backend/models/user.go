package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         uuid.UUID    `json:"id"`
	Username   string       `json:"username"`
	FirstName  string       `json:"first_name"`
	LastName   string       `json:"last_name"`
	Provider   string       `json:"provider"`
	ProviderID string       `json:"provider_id"`
	Password   PasswordType `json:"-"`
	Email      string       `json:"email"`
	IsActive   bool         `json:"is_active"`
	IsVerified bool         `json:"is_verified"`
	Role       Role         `json:"role"`
	RoleID     int          `json:"role_id"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

type PasswordType struct {
	Text *string
	Hash []byte
}

// Set sets the password to the hash of the password_txt
func (p *PasswordType) Set(password_txt string) error {
	var err error
	p.Hash, err = bcrypt.GenerateFromPassword([]byte(password_txt), bcrypt.MinCost)
	p.Text = &password_txt
	if err != nil {
		return err
	}
	return nil
}

func (p *PasswordType) Compare(pass string) error {
	if err := bcrypt.CompareHashAndPassword(p.Hash, []byte(pass)); err != nil {
		return err
	}
	return nil
}
