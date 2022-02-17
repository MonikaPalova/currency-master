package model

import "time"

type Session struct {
	ID         string
	Username   string
	Expiration time.Time
}

func (s Session) IsExpired() bool {
	return s.Expiration.Before(time.Now())
}
