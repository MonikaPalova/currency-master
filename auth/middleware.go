package auth

import (
	"context"
	"log"
	"net/http"

	"github.com/MonikaPalova/currency-master/config"
	"github.com/MonikaPalova/currency-master/svc"
	"github.com/MonikaPalova/currency-master/utils"
)

type sessionUserCtxKey string

const callerCtxKey sessionUserCtxKey = "caller"

type SessionAuth struct {
	Config *config.Session
	Svc    *svc.Sessions
}

func (s SessionAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie(s.Config.SessionCookieName)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, err, "This action requires an active session cookie")
			return
		}
		session, err := s.Svc.GetByID(sessionCookie.Value)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, err, "Session doesn't exist or is already expired")
			return
		}

		ctx := context.WithValue(r.Context(), callerCtxKey, session.Username)
		log.Printf("User %s successfully authenticated", session.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUser(r *http.Request) string {
	return r.Context().Value("caller").(string)
}
