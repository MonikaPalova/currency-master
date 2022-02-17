package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MonikaPalova/currency-master/model"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

type mockUsersSvc struct {
	mock.Mock
}

func (m mockUsersSvc) Create(user model.User) (*model.User, error) {
	args := m.Called(user)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m mockUsersSvc) GetAll() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (m mockUsersSvc) GetByUsername(username string, valuation bool) (user *model.User, err error) {
	args := m.Called(username, valuation)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m mockUsersSvc) AddUSD(username string, usd float64) (float64, error) {
	args := m.Called(username, usd)
	return args.Get(0).(float64), args.Error(1)
}

func (m mockUsersSvc) DeductUSD(username string, usd float64) (float64, error) {
	args := m.Called(username, usd)
	return args.Get(0).(float64), args.Error(1)
}

func TestUsersHandler_Post_InvalidData(t *testing.T) {
	type fields struct {
		user *model.User
		err  error
	}
	type args struct {
		w *httptest.ResponseRecorder
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		body           string
		wantStatusCode int
	}{
		{"non-json body", fields{}, args{httptest.NewRecorder()}, ``, http.StatusBadRequest},
		{"blank username", fields{}, args{httptest.NewRecorder()}, `{"username":"","password":"p1","email":"e1"}`, http.StatusBadRequest},
		{"blank password", fields{}, args{httptest.NewRecorder()}, `{"username":"u1","password":"","email":"e1"}`, http.StatusBadRequest},
		{"blank email", fields{}, args{httptest.NewRecorder()}, `{"username":"u1","password":"p1","email":""}`, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", testAppConfig.UsersApiV1, strings.NewReader(tt.body))

			u := UsersHandler{}
			u.Post(tt.args.w, r)

			if tt.args.w.Code != tt.wantStatusCode {
				t.Fatalf("unexpected status code: got %v want %v", tt.args.w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestUsersHandler_Post(t *testing.T) {
	type fields struct {
		user        model.User
		createdUser *model.User
		err         error
	}
	type args struct {
		w *httptest.ResponseRecorder
	}
	u := model.User{Username: "u1", Password: "p1", Email: "e1"}
	tests := []struct {
		name           string
		fields         fields
		args           args
		body           string
		wantStatusCode int
	}{
		{"users svc error", fields{user: u, createdUser: &u, err: fmt.Errorf("")}, args{httptest.NewRecorder()}, `{"username":"u1","password":"p1","email":"e1"}`, http.StatusInternalServerError},
		{"user already exists", fields{user: u, createdUser: nil, err: nil}, args{httptest.NewRecorder()}, `{"username":"u1","password":"p1","email":"e1"}`, http.StatusConflict},
		{"ok", fields{user: u, createdUser: &u, err: nil}, args{httptest.NewRecorder()}, `{"username":"u1","password":"p1","email":"e1"}`, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", testAppConfig.UsersApiV1, strings.NewReader(tt.body))

			mockUsersSvc := new(mockUsersSvc)
			mockUsersSvc.On("Create", tt.fields.user).Return(tt.fields.createdUser, tt.fields.err)

			u := UsersHandler{Svc: mockUsersSvc}
			u.Post(tt.args.w, r)

			if tt.args.w.Code != tt.wantStatusCode {
				t.Fatalf("unexpected status code: got %v want %v", tt.args.w.Code, tt.wantStatusCode)
			}
			mockUsersSvc.AssertExpectations(t)
		})
	}
}

func TestUsersHandler_GetAll(t *testing.T) {
	type fields struct {
		users []model.User
		err   error
	}
	type args struct {
		w *httptest.ResponseRecorder
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantStatusCode int
	}{
		{"ok", fields{users: []model.User{{Username: "u1", Password: "p1", Email: "e1"}}, err: nil}, args{httptest.NewRecorder()}, http.StatusOK},
		{"ok no users", fields{users: []model.User{}}, args{httptest.NewRecorder()}, http.StatusOK},
		{"users svc error", fields{err: fmt.Errorf("")}, args{httptest.NewRecorder()}, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", testAppConfig.UsersApiV1, nil)

			mockUsersSvc := new(mockUsersSvc)
			mockUsersSvc.On("GetAll").Return(tt.fields.users, tt.fields.err)

			u := UsersHandler{Svc: mockUsersSvc}
			u.GetAll(tt.args.w, r)

			if tt.args.w.Code != tt.wantStatusCode {
				t.Fatalf("unexpected status code: got %v want %v", tt.args.w.Code, tt.wantStatusCode)
			}
			mockUsersSvc.AssertExpectations(t)
		})
	}
}

func TestUsersHandler_GetByUsername(t *testing.T) {
	type fields struct {
		user *model.User
		err  error
	}
	type args struct {
		w        *httptest.ResponseRecorder
		username string
	}
	u := model.User{Username: "u1", Password: "p1", Email: "e1"}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantStatusCode int
	}{
		{"ok", fields{user: &u}, args{httptest.NewRecorder(), u.Username}, http.StatusOK},
		{"no such user", fields{}, args{httptest.NewRecorder(), u.Username}, http.StatusNotFound},
		{"users svc error", fields{err: fmt.Errorf("")}, args{httptest.NewRecorder(), u.Username}, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", testAppConfig.UsersApiV1+"/"+tt.args.username, nil)
			r = mux.SetURLVars(r, map[string]string{
				"username": tt.args.username,
			})

			mockUsersSvc := new(mockUsersSvc)
			mockUsersSvc.On("GetByUsername", tt.args.username, true).Return(tt.fields.user, tt.fields.err)

			u := UsersHandler{Svc: mockUsersSvc}
			u.GetByUsername(tt.args.w, r)

			if tt.args.w.Code != tt.wantStatusCode {
				t.Fatalf("unexpected status code: got %v want %v", tt.args.w.Code, tt.wantStatusCode)
			}
			mockUsersSvc.AssertExpectations(t)
		})
	}
}
