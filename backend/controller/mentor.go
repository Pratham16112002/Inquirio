package controller

import (
	"Inquiro/config"
	"Inquiro/models"
	"Inquiro/repositories"
	"Inquiro/services"
	"Inquiro/utils/json"
	"Inquiro/utils/response"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Mentor struct {
	srv services.Service
	cfg config.Application
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email,max=50"`
	Password string `json:"password" validate:"required,min=3,max=88"`
}

func (u Mentor) MentorLogin(w http.ResponseWriter, r *http.Request) {
	var payload LoginPayload
	err := json.Read(w, r, &payload)
	if err != nil {
		u.cfg.Logger.Warnw("Bad request", "error : ", err.Error())
		response.Error(w, r, "Bad request", err.Error(), 400, http.StatusBadRequest)
		return
	}
	if err := json.Validate.Struct(payload); err != nil {
		u.cfg.Logger.Warnw("Bad request", "error : ", err.Error())
		response.Error(w, r, "Bad request", err.Error(), 400, http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	mentor, err := u.srv.MentorServices.GetMentorByEmail(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			response.Error(w, r, "Login failed", "Mentor does not exist", 404, http.StatusNotFound)
			return
		}
		response.Error(w, r, "Login failed", "Something went wrong", 500, http.StatusInternalServerError)
		return
	}
	if mentor.IsActive == false {
		u.cfg.Logger.Warnw("Mentor not active", "error : ", "mentor has been deactivated")
		response.Error(w, r, "Login Failed", "Mentor does not exist", 404, http.StatusNotFound)
		return
	}

	if mentor.IsVerified == false {
		u.cfg.Logger.Warnw("Mentor not verified", "error : ", "mentor's email has not been verified")
		response.Error(w, r, "Login Failed", "Please verify your email", 404, http.StatusNotFound)
		return
	}

	pass := &models.PasswordType{}
	pass.Set(payload.Password)
	err = u.srv.MentorServices.AuthenticateMentorPassword(ctx, mentor, pass)
	if err != nil {
		response.Error(w, r, "Login failed", "Incorrect credentials", 404, http.StatusNotFound)
		return
	}
	u.cfg.Session.Put(ctx, "mentorId", mentor.ID.String())
	u.cfg.Session.Put(ctx, "userName", mentor.Username)
	u.cfg.Session.Put(ctx, "mentorEmail", mentor.Email)
	response.Success(w, r, "Login successfull", nil, http.StatusOK)
}

type signUpPayload struct {
	Username        string `json:"username" validate:"required,max=50"`
	FirstName       string `json:"first_name" validate:"required,max=100"`
	LastName        string `json:"last_name" validate:"max=100"`
	ExperienceYear  int    `json:"experience_year" validate:"required"`
	ExperienceMonth int    `json:"experience_month" validate:"required"`
	Bio             string `json:"bio"`
	Email           string `json:"email" validate:"required,email,max=50"`
	Password        string `json:"password" validate:"required,min=3,max=88"`
}

func (m Mentor) MentorSignUp(w http.ResponseWriter, r *http.Request) {
	var payload signUpPayload

	err := json.Read(w, r, &payload)
	if err != nil {
		response.Error(w, r, "Bad request", err.Error(), 400, http.StatusBadRequest)
		return
	}
	if err := json.Validate.Struct(payload); err != nil {
		response.Error(w, r, "Bad request", err.Error(), 400, http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	if found := m.srv.MentorServices.CheckMentorUsernameExists(ctx, payload.Username); found == true {
		response.Error(w, r, "Invalid username", "Mentor already taken", 409, http.StatusConflict)
		return
	}
	if found := m.srv.MentorServices.CheckMentorEmailExists(ctx, payload.Email); found == true {
		response.Error(w, r, "SignUp failed", "Email already taken", 409, http.StatusConflict)
		return
	}
	pass := models.PasswordType{}
	pass.Set(payload.Password)
	mentor := &models.Mentor{
		Username:        payload.Username,
		FirstName:       payload.FirstName,
		LastName:        payload.LastName,
		ExperienceYears: convertExperienceToYears(payload.ExperienceYear, payload.ExperienceMonth),
		Bio:             payload.Bio,
		Email:           payload.Email,
		Password:        pass,
		RoleID:          1,
	}
	token := uuid.New().String()
	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])
	if err := m.srv.MentorServices.RegisterMentor(ctx, mentor, hashToken); err != nil {
		response.Error(w, r, "Signup failed", "Account could not be created", 500, http.StatusInternalServerError)
		return
	}
	activationURL := fmt.Sprintf("%s/activate/%s", "http://localhost:3000/mentor", token)
	err = m.cfg.Mail.Send("user_invitation.tmpl", payload.Username, []string{payload.Email}, map[string]string{"ActivationURL": activationURL})
	if err != nil {
		response.Error(w, r, "Signup failed", "Verification email not sent", 500, http.StatusInternalServerError)
		return
	}
	response.Success(w, r, "Signup successful", nil, http.StatusCreated)
}

func (m Mentor) MentorActivation(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	ctx := r.Context()
	err := m.srv.MentorServices.ActivateMentor(ctx, token)
	if err != nil {
		response.Error(w, r, "Activation failed", "Could not activate account", 500, http.StatusInternalServerError)
		return
	}
	response.Success(w, r, "Activation successful", nil, http.StatusOK)
}

func convertExperienceToYears(ExperienceMonth, ExperienceYear int) float32 {
	return float32(ExperienceMonth) + float32(ExperienceYear)
}
