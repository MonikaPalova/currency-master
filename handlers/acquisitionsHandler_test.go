package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MonikaPalova/currency-master/config"
	"github.com/MonikaPalova/currency-master/model"
	"github.com/stretchr/testify/mock"
)

var testAppConfig config.App = *config.NewApp()

type mockAcqDB struct {
	mock.Mock
}

func (m mockAcqDB) GetAll() ([]model.Acquisition, error) {
	args := m.Called()
	return args.Get(0).([]model.Acquisition), args.Error(1)
}

func (m mockAcqDB) GetByUsername(username string) ([]model.Acquisition, error) {
	args := m.Called(username)
	return args.Get(0).([]model.Acquisition), args.Error(1)
}

func (m mockAcqDB) Create(acq model.Acquisition) (*model.Acquisition, error) {
	args := m.Called(acq)
	return args.Get(0).(*model.Acquisition), args.Error(1)
}

func TestAcquisitionsHandler_GetAll(t *testing.T) {
	type fields struct {
		acqs []model.Acquisition
		err  error
	}
	type args struct {
		w     *httptest.ResponseRecorder
		query string
	}
	acq := model.Acquisition{Username: "p1", AssetId: "id1", Quantity: 3, PriceUSD: 0.1, TotalUSD: 0.3, Created: time.Now().UTC()}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantStatusCode int
		wantResponse   []model.Acquisition
	}{
		{"db error", fields{acqs: nil, err: fmt.Errorf("")}, args{httptest.NewRecorder(), ""}, http.StatusInternalServerError, nil},
		{"with username, db error", fields{acqs: nil, err: fmt.Errorf("")}, args{httptest.NewRecorder(), "?username=p1"}, http.StatusInternalServerError, nil},
		{"ok", fields{acqs: []model.Acquisition{acq}, err: nil}, args{httptest.NewRecorder(), ""}, http.StatusOK, []model.Acquisition{acq}},
		{"with username, ok", fields{acqs: []model.Acquisition{acq}, err: nil}, args{httptest.NewRecorder(), "?username=p1"}, http.StatusOK, []model.Acquisition{acq}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", testAppConfig.AcquisitionsApiV1+tt.args.query, nil)

			mockAcqDB := new(mockAcqDB)
			if tt.args.query != "" {
				mockAcqDB.On("GetByUsername", "p1").Return(tt.fields.acqs, tt.fields.err)

			} else {
				mockAcqDB.On("GetAll").Return(tt.fields.acqs, tt.fields.err)
			}

			a := AcquisitionsHandler{DB: mockAcqDB}
			a.GetAll(tt.args.w, r)

			if tt.args.w.Code != tt.wantStatusCode {
				t.Fatalf("unexpected status code: got %v want %v", tt.args.w.Code, tt.wantStatusCode)
			}
			mockAcqDB.AssertExpectations(t)
		})
	}
}
