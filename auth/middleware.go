// Package auth handles authentication middleware
package auth

import (
	"context"
	"log"
	"net/http"

	"github.com/MonikaPalova/currency-master/config"
	"github.com/MonikaPalova/currency-master/svc"
	"github.com/MonikaPalova/currency-master/httputils"
)

type SessionUserCtxKey string

const CallerCtxKey SessionUserCtxKey = "caller"

type SessionAuth struct {
	Config *config.Session
	Svc    *svc.Sessions
}

// Prodvides Middleware function for authentication
func (s SessionAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie(s.Config.SessionCookieName)
		if err != nil {
			httputils.RespondWithError(w, http.StatusUnauthorized, err, "This action requires an active session cookie")
			return
		}
		session, err := s.Svc.GetByID(sessionCookie.Value)
		if err != nil {
			httputils.RespondWithError(w, http.StatusUnauthorized, err, "Invalid session")
			return
		}

		ctx := context.WithValue(r.Context(), CallerCtxKey, session.Username)
		log.Printf("User %s successfully authenticated", session.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Retrieves user from request context
func GetUser(r *http.Request) string {
	return r.Context().Value(CallerCtxKey).(string)
}
