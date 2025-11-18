package controller

import (
	"Inquiro/config"
	"Inquiro/services"
	"net/http"
)

type Controller struct {
	User interface {
		UserSignUp(w http.ResponseWriter, r *http.Request)
		UserLogin(w http.ResponseWriter, r *http.Request)
		UserActivation(w http.ResponseWriter, r *http.Request)
	}
	Resume interface {
		ProcessResume(w http.ResponseWriter, r *http.Request)
	}
	Mentor interface {
		MentorSignUp(w http.ResponseWriter, r *http.Request)
		MentorLogin(w http.ResponseWriter, r *http.Request)
		MentorActivation(w http.ResponseWriter, r *http.Request)
	}
}

func NewController(service services.Service, cfg config.Application) Controller {
	return Controller{
		User: User{
			srv: service,
			cfg: cfg,
		},
		Resume: Resume{
			srv: service,
			cfg: cfg,
		},
		Mentor: Mentor{
			srv: service,
			cfg: cfg,
		},
	}
}
