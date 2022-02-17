package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MonikaPalova/currency-master/coinapi"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

type mockAssetsSvc struct {
	mock.Mock
}

func (m *mockAssetsSvc) GetAssetPage(page, size int) (*coinapi.AssetPage, error) {
	args := m.Called(page, size)
	return args.Get(0).(*coinapi.AssetPage), args.Error(1)
}

func (m *mockAssetsSvc) GetAssetById(id string) (*coinapi.Asset, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*coinapi.Asset), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestAssetsHandler_GetAll(t *testing.T) {
	type fields struct {
		page int
		size int
	}
	type args struct {
		w     *httptest.ResponseRecorder
		query string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantStatusCode int
	}{
		{"default page default size", fields{page: defaultPage, size: defaultSize}, args{httptest.NewRecorder(), ""}, http.StatusOK},
		{"custom page default size", fields{page: 2, size: defaultSize}, args{httptest.NewRecorder(), "?page=2"}, http.StatusOK},
		{"default page custom size", fields{page: defaultPage, size: 2}, args{httptest.NewRecorder(), "?size=2"}, http.StatusOK},
		{"custom page custom size", fields{page: 2, size: 3}, args{httptest.NewRecorder(), "?size=3&page=2"}, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", testAppConfig.AssetsApiV1+tt.args.query, nil)

			mockAssetsSvc := new(mockAssetsSvc)
			mockAssetsSvc.On("GetAssetPage", tt.fields.page, tt.fields.size).Return(&coinapi.AssetPage{Page: tt.fields.page, Size: tt.fields.size, Total: 1}, nil)

			a := AssetsHandler{mockAssetsSvc}
			a.GetAll(tt.args.w, r)

			if tt.args.w.Code != tt.wantStatusCode {
				t.Fatalf("unexpected status code: got %v want %v", tt.args.w.Code, tt.wantStatusCode)
			}
			mockAssetsSvc.AssertExpectations(t)
		})
	}
}

func TestAssetsHandler_GetAll_BadRequest(t *testing.T) {
	type fields struct {
		page int
		size int
		err  error
	}
	type args struct {
		w     *httptest.ResponseRecorder
		query string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantStatusCode int
	}{
		{"negative page", fields{page: -2, size: defaultSize}, args{httptest.NewRecorder(), "?page=-2"}, http.StatusBadRequest},
		{"negative size", fields{page: defaultPage, size: -2}, args{httptest.NewRecorder(), "?size=-2"}, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", testAppConfig.AssetsApiV1+tt.args.query, nil)

			a := AssetsHandler{}
			a.GetAll(tt.args.w, r)

			if tt.args.w.Code != http.StatusBadRequest {
				t.Fatalf("unexpected status code: got %v want %v", tt.args.w.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestAssetsHandler_GetById(t *testing.T) {
	type fields struct {
		asset *coinapi.Asset
		err   error
	}
	type args struct {
		w  *httptest.ResponseRecorder
		id string
	}
	a := coinapi.Asset{ID: "id1", Name: "name1", IsCrypto: true, PriceUSD: 0.1}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantStatusCode int
	}{
		{"ok", fields{asset: &a, err: nil}, args{httptest.NewRecorder(), "id1"}, http.StatusOK},
		{"no such asset", fields{asset: nil, err: nil}, args{httptest.NewRecorder(), "id1"}, http.StatusNotFound},
		{"asset svc error", fields{asset: &a, err: fmt.Errorf("")}, args{httptest.NewRecorder(), "id1"}, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", testAppConfig.AssetsApiV1+"/"+tt.args.id, nil)
			r = mux.SetURLVars(r, map[string]string{
				"id": tt.args.id,
			})

			mockAssetsSvc := new(mockAssetsSvc)
			mockAssetsSvc.On("GetAssetById", tt.args.id).Return(tt.fields.asset, tt.fields.err)

			a := AssetsHandler{mockAssetsSvc}
			a.GetById(tt.args.w, r)

			if tt.args.w.Code != tt.wantStatusCode {
				t.Fatalf("unexpected status code: got %v want %v", tt.args.w.Code, tt.wantStatusCode)
			}
			mockAssetsSvc.AssertExpectations(t)
		})
	}
}
