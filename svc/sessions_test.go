package svc

import (
	"reflect"
	"testing"
	"time"

	"github.com/MonikaPalova/currency-master/model"
)

func TestSessions_GetByID(t *testing.T) {
	type fields struct {
		session model.Session
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"ok", fields{model.Session{ID: "sid1", Username: "u1", Expiration: time.Now().Add(time.Hour)}}, args{"sid1"}, false},
		{"not exist", fields{model.Session{ID: "sid2"}}, args{"sid1"}, true},
		{"expired", fields{model.Session{ID: "sid1", Username: "u1", Expiration: time.Now().Add(-time.Hour)}}, args{"sid1"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Sessions{sessions: map[string]model.Session{tt.fields.session.ID: tt.fields.session}}
			got, err := s.GetByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sessions.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				want := &tt.fields.session
				if !reflect.DeepEqual(got, want) {
					t.Errorf("Sessions.GetByID() = %v, want %v", got, want)
				}
			}
		})
	}
}

func TestSessions_Delete(t *testing.T) {
	sid1 := model.Session{ID: "sid1", Username: "u1", Expiration: time.Now().Add(time.Hour)}
	s := Sessions{sessions: map[string]model.Session{sid1.ID: sid1}}

	s.Delete(sid1.ID)

	_, err := s.GetByID(sid1.ID)
	if err == nil {
		t.Fatalf("session should be invalidated")
	}
}

func TestSessions_ClearExpired(t *testing.T) {
	sid1 := model.Session{ID: "sid1", Username: "u1", Expiration: time.Now().Add(time.Hour)}
	sid2 := model.Session{ID: "sid2", Username: "u2", Expiration: time.Now().Add(-time.Hour)}
	s := Sessions{sessions: map[string]model.Session{sid1.ID: sid1, sid2.ID: sid2}}

	s.ClearExpired()

	want := map[string]model.Session{sid1.ID: sid1}
	if !reflect.DeepEqual(s.sessions, want) {
		t.Fatalf("expired sessions should be deleted. Want: %v, got: %v", want, s.sessions)
	}
}
