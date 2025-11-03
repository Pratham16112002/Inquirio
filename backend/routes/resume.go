package routes

import (
	"Inquiro/controller"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ResumeRoutes struct {
	controller controller.Controller
}

func NewResumeRoutes(controller controller.Controller) ResumeRoutes {
	return ResumeRoutes{
		controller: controller,
	}
}

func (rr ResumeRoutes) RegisterResumeRoutes(chi_router *chi.Mux) {
	chi_router.Route("/resume", func(r chi.Router) {
		r.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
			rr.controller.Resume.ProcessResume(w, r)
		})
	})
}
