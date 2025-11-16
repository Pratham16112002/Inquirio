package models

import (
	"time"

	"github.com/google/uuid"
)

type Mentor struct {
	ID              uuid.UUID    `json:"id"`
	Username        string       `json:"username"`
	FirstName       string       `json:"first_name"`
	LastName        string       `json:"last_name"`
	Provider        string       `json:"provider"`
	ProviderID      string       `json:"provider_id"`
	Password        PasswordType `json:"-"`
	Email           string       `json:"email"`
	IsActive        bool         `json:"is_active"`
	IsVerified      bool         `json:"is_verified"`
	ExperienceYears float32      `json:"experience_years"`
	Bio             string       `json:"bio"`
	Role            Role         `json:"role"`
	RoleID          int          `json:"role_id"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}
