package controller

import (
	"Inquiro/config"
	"Inquiro/services"
	"net/http"
)

type Controller struct {
	Users interface {
		SignUp(w http.ResponseWriter, r *http.Request)
		Login(w http.ResponseWriter, r *http.Request)
		Activate(w http.ResponseWriter, r *http.Request)
	}
	Resume interface {
		ProcessResume(w http.ResponseWriter, r *http.Request)
	}
	Mentor interface {
		SignUp(w http.ResponseWriter, r *http.Request)
		Login(w http.ResponseWriter, r *http.Request)
	}
}

func NewController(service services.Service, cfg config.Application) Controller {
	return Controller{
		Users: User{
			srv: service,
			cfg: cfg,
		},
		Resume: Resume{
			srv: service,
			cfg: cfg,
		},
	}
}
