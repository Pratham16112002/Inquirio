package middlewares

import (
	"Inquiro/config"
	"Inquiro/utils/response"
	"context"
	"net/http"

	"github.com/google/uuid"
)

type Auth struct {
	cfg config.Application
}

func (a Auth) LoadUser() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if !a.cfg.Session.Exists(ctx, "userId") {
				a.cfg.Logger.Errorw("user not logged in", "error :", "user not logged in")
				response.Error(w, r, "Failed", "Not authorized", 401, http.StatusUnauthorized)
				return
			}
			userId := a.cfg.Session.GetString(ctx, "userId")
			uuid, err := uuid.Parse(userId)
			if err != nil {
				// Invalid data in session (rare)
				a.cfg.Session.Clear(r.Context())
				a.cfg.Logger.Errorw("invalid session data in request", "error :", err.Error())
				response.Error(w, r, "Failed", "Not authorized", 401, http.StatusUnauthorized)
				return
			}
			user, err := a.cfg.Store.Users.GetByID(ctx, uuid)
			if err != nil {
				a.cfg.Session.Clear(ctx)
				a.cfg.Logger.Errorw("no user found with this id", "error :", err.Error())
				response.Error(w, r, "Failed", "No user found", 401, http.StatusUnauthorized)
				return
			}
			role, err := a.cfg.Store.Role.GetRoleByID(ctx, user.RoleID)
			if err != nil {
				a.cfg.Logger.Errorw("no role found with this id", "error :", err.Error())
				response.Error(w, r, "Failed", "Invalid user", 401, http.StatusUnauthorized)
				return
			}
			user.Role = role
			ctxWithUser := context.WithValue(ctx, "sessionUser", user)
			next.ServeHTTP(w, r.WithContext(ctxWithUser))
		})
	}
}
