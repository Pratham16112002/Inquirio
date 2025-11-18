package routes

import (
	"Inquiro/controller"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserRoutes struct {
	controller controller.Controller
}

func NewUserRoutes(controller controller.Controller) UserRoutes {
	return UserRoutes{
		controller: controller,
	}
}

func (ur UserRoutes) RegisterUserRoutes(chi_router *chi.Mux) {
	chi_router.Route("/user", func(r chi.Router) {
		r.Put("/activate/{token}", func(w http.ResponseWriter, r *http.Request) {
			ur.controller.User.UserActivation(w, r)
		})
		r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
			ur.controller.User.UserLogin(w, r)
		})
		r.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
			ur.controller.User.UserSignUp(w, r)
		})
	})
}
