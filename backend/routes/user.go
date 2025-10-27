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
			ur.controller.Users.Activate(w, r)
		})
		r.Route("/{provider}", func(r chi.Router) {
			r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
				ur.controller.Users.Login(w, r)
			})
			r.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
				ur.controller.Users.SignUp(w, r)
			})
		})

	})
}
