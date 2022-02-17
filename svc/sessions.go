package svc

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MonikaPalova/currency-master/config"
	"github.com/MonikaPalova/currency-master/model"
	"github.com/google/uuid"
)

type Sessions struct {
	sessions map[string]model.Session
	config   *config.Session
}

func (s Sessions) GetByID(id string) (*model.Session, error) {
	session, ok := s.sessions[id]
	if !ok {
		return nil, fmt.Errorf("session with id %s doesn't exist", id)
	}
	if session.IsExpired() {
		return nil, fmt.Errorf("session with id %s is expired", id)
	}
	return &session, nil
}

func (s Sessions) CreateCookie(username string) *http.Cookie {
	session := model.Session{
		ID:         uuid.New().String(),
		Username:   username,
		Expiration: time.Now().Add(s.config.SessionDuration),
	}
	s.sessions[session.ID] = session

	sessionCookie := http.Cookie{
		Name:    s.config.SessionCookieName,
		Value:   session.ID,
		Expires: session.Expiration,
	}

	return &sessionCookie
}

func (s Sessions) ClearExpired() {
	validSessions := map[string]model.Session{}
	for _, session := range s.sessions {
		if !session.IsExpired() {
			validSessions[session.ID] = session
		}
	}
	s.sessions = validSessions
}