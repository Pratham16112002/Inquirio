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

func (mr MentorRoutes) RegisterMentorRoutes(chi_router *chi.Mux) {
	chi_router.Route("/mentor", func(r chi.Router) {
		r.Put("/activate/{token}", func(w http.ResponseWriter, r *http.Request) {
			mr.controller.Mentor.MentorActivation(w, r)
		})
		r.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
			mr.controller.Mentor.MentorSignUp(w, r)
		})
		r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
			mr.controller.Mentor.MentorLogin(w, r)
		})
	})
}
