package svc

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/MonikaPalova/currency-master/coinapi"
	"github.com/MonikaPalova/currency-master/model"
)

type stubUDB struct {
	users  []model.User
	user   *model.User
	assets []model.UserAsset
	err    error
}

func (s stubUDB) Create(user model.User) (*model.User, error) {
	user.Password = ""
	user.Assets = []model.UserAsset{}

	if s.err != nil {
		return nil, s.err
	}
	return &user, nil
}

func (s stubUDB) GetAll() ([]model.User, error) {
	return s.users, s.err
}
func (s stubUDB) GetByUsername(username string) (*model.User, error) {
	return s.user, s.err
}
func (s stubUDB) GetByUsernameWithAssets(username string) (*model.User, error) {
	if s.err != nil || s.user == nil {
		return nil, s.err
	}
	s.user.Assets = s.assets
	return s.user, nil
}
func (s stubUDB) UpdateUSD(username string, money float64) error {
	s.user.USD = money
	return s.err
}

func (s stubUDB) Exists(username, password string) (bool, error) {
	return false, nil
}

func TestUsers_Create(t *testing.T) {
	type fields struct {
		UDB usersDB
	}
	type args struct {
		user model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		{"valid", fields{stubUDB{}}, args{model.User{Username: "u1", Password: "P1", Email: "e1"}}, &model.User{Username: "u1", Password: "", Email: "e1", USD: startUserUSD, Assets: []model.UserAsset{}, Valuation: 0}, false},
		{"error saving in db", fields{stubUDB{err: fmt.Errorf("")}}, args{model.User{Username: "u1", Password: "P1", Email: "e1"}}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Users{UDB: tt.fields.UDB}
			got, err := u.Create(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Users.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsers_GetAll(t *testing.T) {
	type fields struct {
		UDB    usersDB
		client coinAPIClient
	}
	a1 := coinapi.Asset{ID: "id1", PriceUSD: 0.1}
	a2 := coinapi.Asset{ID: "id2", PriceUSD: 0.2}
	dbU1 := model.User{Username: "u1", Assets: []model.UserAsset{{AssetId: "id1", Quantity: 3}}}
	dbU2 := model.User{Username: "u2", Assets: []model.UserAsset{{AssetId: "id1", Quantity: 4}, {AssetId: "id2", Quantity: 2}}}
	u1 := model.User{Username: "u1", Assets: []model.UserAsset{{AssetId: "id1", Quantity: 3, Valuation: 3 * a1.PriceUSD}}, Valuation: 3 * a1.PriceUSD}
	u2 := model.User{Username: "u2", Assets: []model.UserAsset{{AssetId: "id1", Quantity: 4, Valuation: 4 * a1.PriceUSD}, {AssetId: "id2", Quantity: 2, Valuation: 2 * a2.PriceUSD}}, Valuation: 0.8}

	tests := []struct {
		name    string
		fields  fields
		want    []model.User
		wantErr bool
	}{
		{"valid", fields{stubUDB{users: []model.User{dbU1, dbU2}}, stubClient{[]coinapi.Asset{a1, a2}, nil}}, []model.User{u1, u2}, false},
		{"asset does not exist - no price data", fields{stubUDB{users: []model.User{dbU1}}, stubClient{[]coinapi.Asset{}, fmt.Errorf((""))}}, nil, true},
		{"could not fetch users from db", fields{stubUDB{err: fmt.Errorf("")}, stubClient{[]coinapi.Asset{}, nil}}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Users{
				v:   valuator{svc: NewAssets(tt.fields.client)},
				UDB: tt.fields.UDB,
			}
			got, err := u.GetAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Users.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsers_GetByUsernameValuationTrue(t *testing.T) {
	type fields struct {
		UDB    usersDB
		client coinAPIClient
	}
	type args struct {
		username string
	}
	a := coinapi.Asset{ID: "id1", PriceUSD: 0.1}
	dbU := model.User{Username: "u1", Assets: []model.UserAsset{{AssetId: "id1", Quantity: 3}}}
	ua := []model.UserAsset{{AssetId: "id1", Quantity: 3, Valuation: 3 * a.PriceUSD}}
	u := model.User{Username: "u1", Assets: ua, Valuation: 3 * a.PriceUSD}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantUser *model.User
		wantErr  bool
	}{
		{"db error", fields{stubUDB{err: fmt.Errorf("")}, nil}, args{"u1"}, nil, true},
		{"no such user", fields{stubUDB{user: nil}, nil}, args{"u1"}, nil, false},
		{"ok", fields{stubUDB{user: &dbU, assets: ua}, stubClient{[]coinapi.Asset{a}, nil}}, args{"u1"}, &u, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Users{
				v:   valuator{svc: NewAssets(tt.fields.client)},
				UDB: tt.fields.UDB,
			}
			gotUser, err := u.GetByUsername(tt.args.username, true)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.GetByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("Users.GetByUsername() = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}

func TestUsers_GetByUsernameValuationFalse(t *testing.T) {
	type fields struct {
		UDB usersDB
	}
	type args struct {
		username string
	}
	dbU := model.User{Username: "u1", Assets: []model.UserAsset{}}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantUser *model.User
		wantErr  bool
	}{
		{"db error", fields{stubUDB{err: fmt.Errorf("")}}, args{"u1"}, nil, true},
		{"no such user", fields{stubUDB{user: nil}}, args{"u1"}, nil, false},
		{"ok", fields{stubUDB{user: &dbU}}, args{"u1"}, &dbU, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Users{UDB: tt.fields.UDB}
			gotUser, err := u.GetByUsername(tt.args.username, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.GetByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("Users.GetByUsername() = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}

func TestUsers_AddUSD(t *testing.T) {
	type fields struct {
		UDB usersDB
	}
	type args struct {
		username string
		usd      float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{"negative usd", fields{nil}, args{"u1", -5}, -1, true},
		{"err getting user form db", fields{stubUDB{err: fmt.Errorf("")}}, args{"u1", 5}, -1, true},
		{"no such user", fields{stubUDB{user: nil}}, args{"u1", 5}, -1, true},
		{"ok", fields{stubUDB{user: &model.User{Username: "u1", USD: 15}}}, args{"u1", 5}, 20, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Users{UDB: tt.fields.UDB}
			got, err := u.AddUSD(tt.args.username, tt.args.usd)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.AddUSD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Users.AddUSD() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsers_DeductUSD(t *testing.T) {
	type fields struct {
		UDB usersDB
	}
	type args struct {
		username string
		usd      float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{"negative usd", fields{nil}, args{"u1", -5}, -1, true},
		{"err getting user form db", fields{stubUDB{err: fmt.Errorf("")}}, args{"u1", 5}, -1, true},
		{"no such user", fields{stubUDB{user: nil}}, args{"u1", 5}, -1, true},
		{"ok", fields{stubUDB{user: &model.User{Username: "u1", USD: 15}}}, args{"u1", 5}, 10, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Users{UDB: tt.fields.UDB}
			got, err := u.DeductUSD(tt.args.username, tt.args.usd)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.DeductUSD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Users.DeductUSD() = %v, want %v", got, tt.want)
			}
		})
	}
}
