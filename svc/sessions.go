package svc

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MonikaPalova/currency-master/config"
	"github.com/MonikaPalova/currency-master/model"
	"github.com/google/uuid"
)

// Sessions service to handle sessions
type Sessions struct {
	sessions map[string]model.Session
	Config   *config.Session
}

// Gets session by id if exists.
// Returns error if sessions doesn't exist or is expired
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

// Creates new session for username and a session cookie
func (s Sessions) CreateCookie(username string) *http.Cookie {
	session := model.Session{
		ID:         uuid.New().String(),
		Username:   username,
		Expiration: time.Now().Add(s.Config.SessionDuration),
	}
	s.sessions[session.ID] = session

	sessionCookie := http.Cookie{
		Name:    s.Config.SessionCookieName,
		Value:   session.ID,
		Expires: session.Expiration,
	}

	return &sessionCookie
}

// Invalidates the cookie but keeps it in the array, to be deleted by the cron job
func (s Sessions) Delete(id string) {
	session, ok := s.sessions[id]
	if ok {
		session.Expiration = time.Now().Add(-time.Hour)
		s.sessions[id] = session
	}
	log.Printf("Invalidated session with id %s", id)
}

// Clears expired sessions
func (s *Sessions) ClearExpired() {
	old := len(s.sessions)
	validSessions := map[string]model.Session{}
	for _, session := range s.sessions {
		if !session.IsExpired() {
			fmt.Println("adding session " + session.ID)
			validSessions[session.ID] = session
		}
	}
	s.sessions = validSessions
	log.Printf("Cleared expired sessions. Deleted: %d", old-len(s.sessions))
}
