// Package svc contains services that prepare data for the handlers
package svc

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/MonikaPalova/currency-master/coinapi"
	"github.com/MonikaPalova/currency-master/model"
)

type stubClient struct {
	assets []coinapi.Asset
	err    error
}

func (s stubClient) GetAssets() ([]coinapi.Asset, error) {
	return s.assets, s.err
}

func TestAssets_GetAssetPage(t *testing.T) {
	type fields struct {
		client coinAPIClient
	}
	type args struct {
		page int
		size int
	}
	a := coinapi.Asset{ID: "id1"}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *coinapi.AssetPage
		wantErr bool
	}{
		{"no assets", fields{stubClient{[]coinapi.Asset{}, nil}}, args{1, 1}, &coinapi.AssetPage{Assets: []coinapi.Asset{}, Page: 1, Size: 1, Total: 0}, false},
		{"valid asset page", fields{stubClient{[]coinapi.Asset{a}, nil}}, args{1, 1}, &coinapi.AssetPage{Assets: []coinapi.Asset{a}, Page: 1, Size: 1, Total: 1}, false},
		{"cache update error", fields{stubClient{[]coinapi.Asset{}, fmt.Errorf("")}}, args{1, 1}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAssets(tt.fields.client)
			got, err := a.GetAssetPage(tt.args.page, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("Assets.GetAssetPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Assets.GetAssetPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssets_GetAssetById(t *testing.T) {
	type fields struct {
		client coinAPIClient
	}
	type args struct {
		id string
	}
	a := coinapi.Asset{ID: "id1"}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *coinapi.Asset
		wantErr bool
	}{
		{"no such asset", fields{stubClient{[]coinapi.Asset{a}, nil}}, args{"id3"}, nil, false},
		{"valid asset", fields{stubClient{[]coinapi.Asset{a}, nil}}, args{"id1"}, &a, false},
		{"cache update error", fields{stubClient{[]coinapi.Asset{}, fmt.Errorf("")}}, args{"id1"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAssets(tt.fields.client)
			got, err := a.GetAssetById(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Assets.GetAssetById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Assets.GetAssetById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssets_Valuate(t *testing.T) {
	type fields struct {
		client coinAPIClient
	}
	type args struct {
		ua model.UserAsset
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{"no such asset", fields{stubClient{[]coinapi.Asset{{ID: "id1", PriceUSD: 0.01}}, nil}}, args{model.UserAsset{AssetId: "id3", Quantity: 2}}, -1, true},
		{"valid asset", fields{stubClient{[]coinapi.Asset{{ID: "id1", PriceUSD: 0.01}}, nil}}, args{model.UserAsset{AssetId: "id1", Quantity: 2}}, 0.02, false},
		{"cache update error", fields{stubClient{[]coinapi.Asset{}, fmt.Errorf("")}}, args{model.UserAsset{AssetId: "id1", Quantity: 2}}, -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAssets(tt.fields.client)
			got, err := a.Valuate(tt.args.ua)
			if (err != nil) != tt.wantErr {
				t.Errorf("Assets.Valuate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Assets.Valuate() = %v, want %v", got, tt.want)
			}
		})
	}
}
