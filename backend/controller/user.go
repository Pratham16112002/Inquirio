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

type User struct {
	srv services.Service
	cfg config.Application
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email,max=50"`
	Password string `json:"password" validate:"required,min=3,max=88"`
}

func (u User) Login(w http.ResponseWriter, r *http.Request) {
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
	user, err := u.srv.UserServices.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			response.Error(w, r, "Login Failed", "User does not exist", 404, http.StatusNotFound)
			return
		}
		response.Error(w, r, "Something is wrong on our side", err.Error(), 500, http.StatusInternalServerError)
		return
	}
	if user.IsActive == false {
		u.cfg.Logger.Warnw("User not active", "error : ", "user has been deactivated")
		response.Error(w, r, "Login Failed", "User does not exist", 404, http.StatusNotFound)
		return
	}

	if user.IsVerified == false {
		u.cfg.Logger.Warnw("User not verified", "error : ", "user's email has not been verified")
		response.Error(w, r, "Login Failed", "Please verify your email", 404, http.StatusNotFound)
		return
	}

	pass := &models.PasswordType{}
	pass.Set(payload.Password)
	err = u.srv.UserServices.AuthenticatePassword(ctx, user, pass)
	if err != nil {
		response.Error(w, r, "Login Failed", "Incorrect credentials", 404, http.StatusNotFound)
		return
	}
	u.cfg.Session.Put(ctx, "userId", user.ID)
	u.cfg.Session.Put(ctx, "userName", user.Username)
	u.cfg.Session.Put(ctx, "userEmail", user.Email)
	response.Success(w, r, "Login Successfull", nil, http.StatusOK)
}

type SingUpPayload struct {
	Username   string `json:"username" validate:"required,max=50"`
	FirstName  string `json:"first_name" validate:"required,max=100"`
	LastName   string `json:"last_name" validate:"max=100"`
	Provider   string `json:"provider"`
	ProviderID string `json:"provider_id"`
	Email      string `json:"email" validate:"required,email,max=50"`
	Password   string `json:"password" validate:"required,min=3,max=88"`
}

func (u User) SignUp(w http.ResponseWriter, r *http.Request) {
	var payload SingUpPayload

	provider := chi.URLParam(r, "provider")

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
	if found, err := u.srv.UserServices.CheckUsernameExists(ctx, payload.Username); found == true {
		response.Error(w, r, "Invalid username", err.Error(), 409, http.StatusConflict)
		return
	}
	if found, err := u.srv.UserServices.CheckEmailExists(ctx, payload.Email); found == true {
		response.Error(w, r, "Invalid email", err.Error(), 409, http.StatusConflict)
		return
	}
	pass := &models.PasswordType{}
	pass.Set(payload.Password)
	user := &models.User{
		ID:         uuid.New(),
		Username:   payload.Username,
		FirstName:  payload.FirstName,
		LastName:   payload.LastName,
		Email:      payload.Email,
		Provider:   provider,
		ProviderID: "",
		Password:   pass,
	}
	token := uuid.New().String()
	hash := sha256.Sum256([]byte(token))
	hashedToken := hex.EncodeToString(hash[:])
	err = u.srv.UserServices.RegisterUser(ctx, user, hashedToken)
	if err != nil {
		response.Error(w, r, "Signup failed", err.Error(), 500, http.StatusInternalServerError)
		return
	}
	activationURL := fmt.Sprintf("%s/activate/%s", "http://localhost:3000/user", token)
	err = u.cfg.Mail.Send("user_invitation.tmpl", payload.Username, []string{payload.Email}, map[string]string{"ActivationURL": activationURL})
	if err != nil {
		response.Error(w, r, "Verfication email was not sent", err.Error(), 500, http.StatusInternalServerError)
		return
	}
	response.Success(w, r, "Signup Successful", nil, http.StatusCreated)
}

func (u User) Activate(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	ctx := r.Context()
	err := u.srv.UserServices.ActivateUser(ctx, token)
	if err != nil {
		u.cfg.Logger.Infow("Activation not completed", "error : ", err.Error())
		response.Error(w, r, "Activation failed", "Could not activate account", 500, http.StatusInternalServerError)
		return
	}
	response.Success(w, r, "Activation Successful", nil, http.StatusOK)
}
