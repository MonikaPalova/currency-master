package svc

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/MonikaPalova/currency-master/coinapi"
	"github.com/MonikaPalova/currency-master/model"
)

type stubUaDB struct {
	assets []model.UserAsset
	asset  *model.UserAsset
	err    error
}

func (s stubUaDB) GetByUsername(username string) ([]model.UserAsset, error) {
	return s.assets, s.err
}

func (s stubUaDB) GetByUsernameAndId(username, id string) (*model.UserAsset, error) {
	return s.asset, s.err
}

func (s stubUaDB) Create(asset model.UserAsset) (*model.UserAsset, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &asset, nil
}

func (s stubUaDB) Update(asset model.UserAsset) (*model.UserAsset, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &asset, nil
}

func (s stubUaDB) Delete(asset model.UserAsset) error {
	return s.err
}

func TestUsers_GetAssetsByUsername(t *testing.T) {
	type fields struct {
		UaDB   userAssetsDB
		client coinAPIClient
	}
	type args struct {
		username string
	}
	a1 := coinapi.Asset{ID: "id1", PriceUSD: 0.1}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.UserAsset
		wantErr bool
	}{
		{"no assets", fields{stubUaDB{assets: []model.UserAsset{}, err: nil}, stubClient{assets: []coinapi.Asset{}, err: nil}}, args{"u1"}, []model.UserAsset{}, false},
		{"ok", fields{stubUaDB{assets: []model.UserAsset{{AssetId: "id1", Quantity: 3}}, err: nil}, stubClient{assets: []coinapi.Asset{a1}, err: nil}}, args{"u1"}, []model.UserAsset{{AssetId: "id1", Quantity: 3, Valuation: 3 * a1.PriceUSD}}, false},
		{"db error", fields{stubUaDB{err: fmt.Errorf("")}, nil}, args{"u1"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UserAssets{
				v:    valuator{svc: NewAssets(tt.fields.client)},
				UaDB: tt.fields.UaDB,
			}
			got, err := u.GetByUsername(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.GetAssetsByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Users.GetAssetsByUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsers_GetAssetByUsernameAndId(t *testing.T) {
	type fields struct {
		UaDB   userAssetsDB
		client coinAPIClient
	}
	type args struct {
		username string
		id       string
	}
	a := coinapi.Asset{ID: "id1", PriceUSD: 0.1}
	ua := model.UserAsset{AssetId: "id1", Quantity: 3, Valuation: 3 * a.PriceUSD}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.UserAsset
		wantErr bool
	}{
		{"db error", fields{stubUaDB{err: fmt.Errorf("")}, nil}, args{"u1", "id1"}, nil, true},
		{"no such user", fields{stubUaDB{asset: nil, err: nil}, nil}, args{"u1", "id1"}, nil, false},
		{"no such asset", fields{stubUaDB{asset: nil, err: nil}, nil}, args{"u1", "id1"}, nil, false},
		{"ok", fields{stubUaDB{asset: &model.UserAsset{AssetId: "id1", Quantity: 3}, err: nil}, stubClient{assets: []coinapi.Asset{a}, err: nil}}, args{"u1", "id1"}, &ua, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UserAssets{
				v:    valuator{svc: NewAssets(tt.fields.client)},
				UaDB: tt.fields.UaDB,
			}
			got, err := u.GetByUsernameAndId(tt.args.username, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.GetAssetByUsernameAndId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Users.GetAssetByUsernameAndId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsers_CreateAsset(t *testing.T) {
	type fields struct {
		UaDB userAssetsDB
	}
	type args struct {
		asset model.UserAsset
	}
	ua := model.UserAsset{Username: "u1", AssetId: "a1"}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.UserAsset
		wantErr bool
	}{
		{"valid", fields{stubUaDB{}}, args{ua}, &ua, false},
		{"error saving in db", fields{stubUaDB{err: fmt.Errorf("")}}, args{ua}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UserAssets{UaDB: tt.fields.UaDB}
			got, err := u.Create(tt.args.asset)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.CreateAsset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Users.CreateAsset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsers_DeleteAsset(t *testing.T) {
	type fields struct {
		UaDB userAssetsDB
	}
	type args struct {
		asset model.UserAsset
	}
	ua := model.UserAsset{Username: "u1", AssetId: "a1"}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"valid", fields{stubUaDB{}}, args{ua}, false},
		{"error deleting from db", fields{stubUaDB{err: fmt.Errorf("")}}, args{ua}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UserAssets{UaDB: tt.fields.UaDB}
			if err := u.Delete(tt.args.asset); (err != nil) != tt.wantErr {
				t.Errorf("Users.DeleteAsset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsers_UpdateAsset(t *testing.T) {
	type fields struct {
		UaDB userAssetsDB
	}
	type args struct {
		asset model.UserAsset
	}
	ua := model.UserAsset{Username: "u1", AssetId: "a1"}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.UserAsset
		wantErr bool
	}{
		{"valid", fields{stubUaDB{}}, args{ua}, &ua, false},
		{"error updating db", fields{stubUaDB{err: fmt.Errorf("")}}, args{ua}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UserAssets{UaDB: tt.fields.UaDB}
			got, err := u.Update(tt.args.asset)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.UpdateAsset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Users.UpdateAsset() = %v, want %v", got, tt.want)
			}
		})
	}
}
