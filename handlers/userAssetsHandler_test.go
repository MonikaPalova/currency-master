package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/MonikaPalova/currency-master/auth"
	"github.com/MonikaPalova/currency-master/coinapi"
	"github.com/MonikaPalova/currency-master/model"
	"github.com/stretchr/testify/mock"
)

type mockUserAssetsSvc struct {
	mock.Mock
}

func (m mockUserAssetsSvc) GetByUsername(username string) ([]model.UserAsset, error) {
	args := m.Called(username)
	return args.Get(0).([]model.UserAsset), args.Error(1)
}

func (m mockUserAssetsSvc) GetByUsernameAndId(username, id string) (*model.UserAsset, error) {
	args := m.Called(username, id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.UserAsset), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m mockUserAssetsSvc) Create(asset model.UserAsset) (*model.UserAsset, error) {
	args := m.Called(asset)
	return args.Get(0).(*model.UserAsset), args.Error(1)
}

func (m mockUserAssetsSvc) Update(asset model.UserAsset) (*model.UserAsset, error) {
	args := m.Called(asset)
	return args.Get(0).(*model.UserAsset), args.Error(1)
}

func (m mockUserAssetsSvc) Delete(asset model.UserAsset) error {
	args := m.Called(asset)
	return args.Error(0)
}

type testCtx struct {
	username string
	id       string
}

func (m testCtx) Value(key interface{}) interface{} {
	if reflect.TypeOf(key) == reflect.TypeOf(auth.CallerCtxKey) {
		return m.username // wants caller
	} else { // wants vars map
		vars := map[string]string{}
		vars["username"] = m.username
		vars["id"] = m.id
		return vars
	}
}

func (m testCtx) Err() error {
	return nil
}
func (m testCtx) Done() <-chan struct{} {
	return nil
}
func (m testCtx) Deadline() (deadline time.Time, ok bool) {
	return time.Now(), false
}

func TestUserAssetsHandler_GetAll(t *testing.T) {
	type fields struct {
		assets []model.UserAsset
		err    error
	}
	type args struct {
		w        *httptest.ResponseRecorder
		username string
	}
	ua := model.UserAsset{Username: "u1", AssetId: "id1", Quantity: 3, Valuation: 0.3}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantStatusCode int
	}{
		{"no such user", fields{assets: []model.UserAsset{ua}, err: nil}, args{httptest.NewRecorder(), ua.Username}, http.StatusOK},
		{"svc error", fields{assets: []model.UserAsset{ua}, err: fmt.Errorf("")}, args{httptest.NewRecorder(), ua.Username}, http.StatusInternalServerError},
		{"no assets", fields{assets: []model.UserAsset{}, err: nil}, args{httptest.NewRecorder(), ua.Username}, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", testAppConfig.AssetsApiV1+"/"+tt.args.username, nil)

			mockUserAssetsSvc := new(mockUserAssetsSvc)
			mockUserAssetsSvc.On("GetByUsername", tt.args.username).Return(tt.fields.assets, tt.fields.err)

			u := UserAssetsHandler{UaSvc: mockUserAssetsSvc}
			r = r.WithContext(testCtx{username: tt.args.username})
			u.GetAll(tt.args.w, r)

			if tt.args.w.Code != tt.wantStatusCode {
				t.Fatalf("unexpected status code: got %v want %v", tt.args.w.Code, tt.wantStatusCode)
			}
			mockUserAssetsSvc.AssertExpectations(t)
		})
	}
}

func TestUserAssetsHandler_GetByID(t *testing.T) {
	type fields struct {
		asset *model.UserAsset
		err   error
	}
	type args struct {
		w        *httptest.ResponseRecorder
		username string
		id       string
	}
	ua := model.UserAsset{Username: "u2", AssetId: "id1", Quantity: 3, Valuation: 0.3}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantStatusCode int
	}{
		{"no such user/asset", fields{asset: nil, err: nil}, args{httptest.NewRecorder(), ua.Username, ua.AssetId}, http.StatusNotFound},
		{"svc error", fields{asset: &ua, err: fmt.Errorf("")}, args{httptest.NewRecorder(), ua.Username, ua.AssetId}, http.StatusInternalServerError},
		{"ok", fields{asset: &ua, err: nil}, args{httptest.NewRecorder(), ua.Username, ua.AssetId}, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", testAppConfig.AssetsApiV1+"/"+tt.args.username+"/"+tt.args.id, nil)

			mockUserAssetsSvc := new(mockUserAssetsSvc)
			mockUserAssetsSvc.On("GetByUsernameAndId", tt.args.username, tt.args.id).Return(tt.fields.asset, tt.fields.err)

			u := UserAssetsHandler{UaSvc: mockUserAssetsSvc}
			r = r.WithContext(testCtx{username: tt.args.username, id: tt.args.id})
			u.GetByID(tt.args.w, r)

			if tt.args.w.Code != tt.wantStatusCode {
				t.Fatalf("unexpected status code: got %v want %v", tt.args.w.Code, tt.wantStatusCode)
			}
			mockUserAssetsSvc.AssertExpectations(t)
		})
	}
}

func TestUserAssetsHandler_Buy_BadRequest(t *testing.T) {
	type args struct {
		w        *httptest.ResponseRecorder
		username string
		id       string
		query    string
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{"no quantity", args{httptest.NewRecorder(), "u1", "id1", ""}, http.StatusBadRequest},
		{"invalid quantity", args{httptest.NewRecorder(), "u1", "id1", "?quantity=-1"}, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", testAppConfig.UsersApiV1+"/"+tt.args.username+"/assets/"+tt.args.id+"/buy"+tt.args.query, nil)

			u := UserAssetsHandler{}
			r = r.WithContext(testCtx{username: tt.args.username, id: tt.args.id})
			u.Buy(tt.args.w, r)

			if tt.args.w.Code != tt.wantStatusCode {
				t.Fatalf("unexpected status code: got %v want %v", tt.args.w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestUserAssetsHandler_Buy_NotFound(t *testing.T) {
	type fields struct {
		asset *coinapi.Asset
		user  *model.User
	}
	type args struct {
		w        *httptest.ResponseRecorder
		username string
		id       string
		query    string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantStatusCode int
	}{
		{"no such asset in external api", fields{asset: nil, user: &model.User{Username: "u1"}}, args{httptest.NewRecorder(), "u1", "id1", "?quantity=1"}, http.StatusNotFound},
		{"no such user", fields{asset: &coinapi.Asset{ID: "id1"}, user: nil}, args{httptest.NewRecorder(), "u1", "id1", "?quantity=1"}, http.StatusNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", testAppConfig.UsersApiV1+"/"+tt.args.username+"/assets/"+tt.args.id+"/buy"+tt.args.query, nil)

			mockUsersSvc := new(mockUsersSvc)
			mockUsersSvc.On("GetByUsername", tt.args.username, false).Return(tt.fields.user, nil)
			mockAssetsSvc := new(mockAssetsSvc)
			mockAssetsSvc.On("GetAssetById", tt.args.id).Return(tt.fields.asset, nil)

			u := UserAssetsHandler{USvc: mockUsersSvc, ASvc: mockAssetsSvc}
			r = r.WithContext(testCtx{username: tt.args.username, id: tt.args.id})
			u.Buy(tt.args.w, r)

			if tt.args.w.Code != tt.wantStatusCode {
				t.Fatalf("unexpected status code: got %v want %v", tt.args.w.Code, tt.wantStatusCode)
			}
			mockAssetsSvc.AssertExpectations(t)
		})
	}
}

func TestUserAssetsHandler_Buy_NoMoney(t *testing.T) {
	username := "u1"
	id := "id1"
	r := httptest.NewRequest("POST", testAppConfig.UsersApiV1+"/"+username+"/assets/"+id+"/buy?quantity=1", nil)

	mockUsersSvc := new(mockUsersSvc)
	mockUsersSvc.On("GetByUsername", username, false).Return(&model.User{Username: username, USD: 2}, nil)
	mockAssetsSvc := new(mockAssetsSvc)
	mockAssetsSvc.On("GetAssetById", id).Return(&coinapi.Asset{ID: id, PriceUSD: 100}, nil)

	u := UserAssetsHandler{USvc: mockUsersSvc, ASvc: mockAssetsSvc}
	w := httptest.NewRecorder()
	r = r.WithContext(testCtx{username: username, id: id})
	u.Buy(w, r)

	if w.Code != http.StatusConflict {
		t.Fatalf("unexpected status code: got %v want %v", w.Code, http.StatusConflict)
	}
	mockUsersSvc.AssertExpectations(t)
	mockAssetsSvc.AssertExpectations(t)
}

func TestUserAssetsHandler_Buy_NewAsset(t *testing.T) {
	username := "u1"
	id := "id1"
	a := coinapi.Asset{ID: id, PriceUSD: 2, Name: "n1"}
	q := 1.0
	ua := model.UserAsset{Username: username, AssetId: id, Name: a.Name, Quantity: q}
	r := httptest.NewRequest("POST", testAppConfig.UsersApiV1+"/"+username+"/assets/"+id+"/buy?quantity=1", nil)

	mockUsersSvc := new(mockUsersSvc)
	mockUsersSvc.On("GetByUsername", username, false).Return(&model.User{Username: username, USD: 2}, nil)
	mockUsersSvc.On("DeductUSD", username, q*a.PriceUSD).Return(2-q*a.PriceUSD, nil)
	mockAssetsSvc := new(mockAssetsSvc)
	mockAssetsSvc.On("GetAssetById", id).Return(&a, nil)
	mockUserAssetsSvc := new(mockUserAssetsSvc)
	mockUserAssetsSvc.On("GetByUsernameAndId", username, id).Return(nil, nil)
	mockUserAssetsSvc.On("Create", ua).Return(&ua, nil)
	mockAcqDB := new(mockAcqDB)
	mockAcqDB.On("Create", mock.Anything).Return(&model.Acquisition{}, nil)

	u := UserAssetsHandler{USvc: mockUsersSvc, ASvc: mockAssetsSvc, UaSvc: mockUserAssetsSvc, ADB: mockAcqDB}
	w := httptest.NewRecorder()
	r = r.WithContext(testCtx{username: username, id: id})
	u.Buy(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %v want %v", w.Code, http.StatusOK)
	}
	mockUsersSvc.AssertExpectations(t)
	mockAssetsSvc.AssertExpectations(t)
	mockUserAssetsSvc.AssertExpectations(t)
	mockAcqDB.AssertExpectations(t)
}

func TestUserAssetsHandler_Buy_UpdateAsset(t *testing.T) {
	username := "u1"
	id := "id1"
	a := coinapi.Asset{ID: id, PriceUSD: 2, Name: "n1"}
	q := 1.0
	ua := model.UserAsset{Username: username, AssetId: id, Name: a.Name, Quantity: 2}
	updatedUa := model.UserAsset{Username: username, AssetId: id, Name: a.Name, Quantity: 2 + q}
	r := httptest.NewRequest("POST", testAppConfig.UsersApiV1+"/"+username+"/assets/"+id+"/buy?quantity=1", nil)

	mockUsersSvc := new(mockUsersSvc)
	mockUsersSvc.On("GetByUsername", username, false).Return(&model.User{Username: username, USD: 2}, nil)
	mockUsersSvc.On("DeductUSD", username, q*a.PriceUSD).Return(2-q*a.PriceUSD, nil)
	mockAssetsSvc := new(mockAssetsSvc)
	mockAssetsSvc.On("GetAssetById", id).Return(&a, nil)
	mockUserAssetsSvc := new(mockUserAssetsSvc)
	mockUserAssetsSvc.On("GetByUsernameAndId", username, id).Return(&ua, nil)
	mockUserAssetsSvc.On("Update", updatedUa).Return(&updatedUa, nil)
	mockAcqDB := new(mockAcqDB)
	mockAcqDB.On("Create", mock.Anything).Return(&model.Acquisition{}, nil)

	u := UserAssetsHandler{USvc: mockUsersSvc, ASvc: mockAssetsSvc, UaSvc: mockUserAssetsSvc, ADB: mockAcqDB}
	w := httptest.NewRecorder()
	r = r.WithContext(testCtx{username: username, id: id})
	u.Buy(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %v want %v", w.Code, http.StatusOK)
	}
	mockUsersSvc.AssertExpectations(t)
	mockAssetsSvc.AssertExpectations(t)
	mockUserAssetsSvc.AssertExpectations(t)
	mockAcqDB.AssertExpectations(t)
}

func TestUserAssetsHandler_Sell_BadRequest(t *testing.T) {
	type args struct {
		w        *httptest.ResponseRecorder
		username string
		id       string
		query    string
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{"no quantity", args{httptest.NewRecorder(), "u1", "id1", ""}, http.StatusBadRequest},
		{"invalid quantity", args{httptest.NewRecorder(), "u1", "id1", "?quantity=-1"}, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", testAppConfig.UsersApiV1+"/"+tt.args.username+"/assets/"+tt.args.id+"/sell"+tt.args.query, nil)

			u := UserAssetsHandler{}
			r = r.WithContext(testCtx{username: tt.args.username, id: tt.args.id})
			u.Sell(tt.args.w, r)

			if tt.args.w.Code != tt.wantStatusCode {
				t.Fatalf("unexpected status code: got %v want %v", tt.args.w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestUserAssetsHandler_Sell_Gone(t *testing.T) {
	username := "u1"
	id := "id1"
	ua := model.UserAsset{Username: username, AssetId: id, Name: "name", Quantity: 2}
	r := httptest.NewRequest("POST", testAppConfig.UsersApiV1+"/"+username+"/assets/"+id+"/sell?quantity=1", nil)

	mockAssetsSvc := new(mockAssetsSvc)
	mockAssetsSvc.On("GetAssetById", id).Return(nil, nil)
	mockUserAssetsSvc := new(mockUserAssetsSvc)
	mockUserAssetsSvc.On("GetByUsernameAndId", username, id).Return(&ua, nil)

	u := UserAssetsHandler{ASvc: mockAssetsSvc, UaSvc: mockUserAssetsSvc}
	w := httptest.NewRecorder()
	r = r.WithContext(testCtx{username: username, id: id})
	u.Sell(w, r)

	if w.Code != http.StatusGone {
		t.Fatalf("unexpected status code: got %v want %v", w.Code, http.StatusGone)
	}
	mockAssetsSvc.AssertExpectations(t)
	mockUserAssetsSvc.AssertExpectations(t)
}

func TestUserAssetsHandler_Sell_NoQuantity(t *testing.T) {
	username := "u1"
	id := "id1"
	ua := model.UserAsset{Username: username, AssetId: id, Name: "name", Quantity: 2}
	r := httptest.NewRequest("POST", testAppConfig.UsersApiV1+"/"+username+"/assets/"+id+"/sell?quantity=5", nil)

	mockUserAssetsSvc := new(mockUserAssetsSvc)
	mockUserAssetsSvc.On("GetByUsernameAndId", username, id).Return(&ua, nil)

	u := UserAssetsHandler{UaSvc: mockUserAssetsSvc}
	w := httptest.NewRecorder()
	r = r.WithContext(testCtx{username: username, id: id})
	u.Sell(w, r)

	if w.Code != http.StatusConflict {
		t.Fatalf("unexpected status code: got %v want %v", w.Code, http.StatusConflict)
	}
	mockUserAssetsSvc.AssertExpectations(t)
}

func TestUserAssetsHandler_Sell_NoAsset(t *testing.T) {
	username := "u1"
	id := "id1"
	r := httptest.NewRequest("POST", testAppConfig.UsersApiV1+"/"+username+"/assets/"+id+"/sell?quantity=5", nil)

	mockUserAssetsSvc := new(mockUserAssetsSvc)
	mockUserAssetsSvc.On("GetByUsernameAndId", username, id).Return(nil, nil)

	u := UserAssetsHandler{UaSvc: mockUserAssetsSvc}
	w := httptest.NewRecorder()
	r = r.WithContext(testCtx{username: username, id: id})
	u.Sell(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("unexpected status code: got %v want %v", w.Code, http.StatusNotFound)
	}
	mockUserAssetsSvc.AssertExpectations(t)
}

func TestUserAssetsHandler_Sell_Ok(t *testing.T) {
	username := "u1"
	id := "id1"
	a := coinapi.Asset{ID: id, PriceUSD: 2, Name: "n1"}
	ua := model.UserAsset{Username: username, AssetId: id, Name: a.Name, Quantity: 5}
	updatedUa := model.UserAsset{Username: username, AssetId: id, Name: a.Name, Quantity: 5 - 1}
	r := httptest.NewRequest("POST", testAppConfig.UsersApiV1+"/"+username+"/assets/"+id+"/sell?quantity=1", nil)

	mockUserAssetsSvc := new(mockUserAssetsSvc)
	mockUserAssetsSvc.On("GetByUsernameAndId", username, id).Return(&ua, nil)
	mockUserAssetsSvc.On("Update", updatedUa).Return(&updatedUa, nil)
	mockAssetsSvc := new(mockAssetsSvc)
	mockAssetsSvc.On("GetAssetById", id).Return(&a, nil)
	mockUsersSvc := new(mockUsersSvc)
	mockUsersSvc.On("AddUSD", username, 1*a.PriceUSD).Return(3.0, nil)

	u := UserAssetsHandler{UaSvc: mockUserAssetsSvc, ASvc: mockAssetsSvc, USvc: mockUsersSvc}
	w := httptest.NewRecorder()
	r = r.WithContext(testCtx{username: username, id: id})
	u.Sell(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %v want %v", w.Code, http.StatusOK)
	}
	mockUserAssetsSvc.AssertExpectations(t)
	mockAssetsSvc.AssertExpectations(t)
	mockUsersSvc.AssertExpectations(t)
}
