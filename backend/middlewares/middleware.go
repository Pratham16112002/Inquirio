package middlewares

import (
	"Inquiro/config"
	"net/http"
)

type Middleware struct {
	Auth interface {
		LoadUser() func(http.Handler) http.Handler
	}
}

func NewMiddleware(cfg config.Application) Middleware {
	return Middleware{
		Auth: &Auth{
			cfg: cfg,
		},
	}
}
