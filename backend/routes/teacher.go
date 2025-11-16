package routes

import (
	"Inquiro/controller"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type MentorRoutes struct {
	controller controller.Controller
}

func NewMentorRoutes(controller controller.Controller) MentorRoutes {
	return MentorRoutes{
		controller: controller,
	}
}

func (tr MentorRoutes) RegisterMentorRoutes(chi_router *chi.Mux) {
	chi_router.Route("/Mentor", func(r chi.Router) {
		r.Post("/signup", func(w http.ResponseWriter, r *http.Request) {

		})
		r.Post("/login", func(w http.ResponseWriter, r *http.Request) {

		})
	})
}
