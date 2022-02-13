package svc

import (
	"fmt"

	"github.com/MonikaPalova/currency-master/coinapi"
	"github.com/MonikaPalova/currency-master/model"
)

type Assets struct {
	cache  *coinapi.Cache
	client *coinapi.Client
}

func NewAssets() *Assets {
	return &Assets{cache: coinapi.NewCache(), client: coinapi.NewClient()}
}

func (a Assets) GetAssetPage(page, size int) (*coinapi.AssetPage, error) {
	if err := a.updateCacheIfNeeded(); err != nil {
		return nil, err
	}

	assetsPage := a.cache.GetPage(page, size)
	return &assetsPage, nil
}

func (a Assets) GetAssetById(id string) (*coinapi.Asset, error) {
	if err := a.updateCacheIfNeeded(); err != nil {
		return nil, err
	}

	return a.cache.GetAsset(id), nil
}

func (a Assets) updateCacheIfNeeded() error {
	if a.cache.IsExpired() {
		fmt.Println("UPDATING CACHE")
		assets, err := a.client.GetAssets()
		if err != nil {
			return fmt.Errorf("error retrieving assets from external api: %v", err)
		}
		a.cache.Fill(assets)
	}

	return nil
}

func (a Assets) Valuate(ua model.UserAsset) (float64, error) {
	asset, err := a.GetAssetById(ua.AssetId)
	if err != nil {
		return -1, err
	}

	return asset.PriceUSD * ua.Quantity, nil
}
