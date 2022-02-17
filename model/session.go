// Package model contains structs that represent object created by the application
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
