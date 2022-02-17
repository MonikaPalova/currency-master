package svc

import (
	"fmt"
	"log"

	"github.com/MonikaPalova/currency-master/coinapi"
	"github.com/MonikaPalova/currency-master/model"
)

// Assets service which handles assets retrieval from cache and external api
type Assets struct {
	cache  *coinapi.Cache
	client coinAPIClient
}

type coinAPIClient interface {
	GetAssets() ([]coinapi.Asset, error)
}

// Constructor
func NewAssets(client coinAPIClient) *Assets {
	return &Assets{cache: coinapi.NewCache(), client: client}
}

// Gets a specific asset page
func (a Assets) GetAssetPage(page, size int) (*coinapi.AssetPage, error) {
	if err := a.updateCacheIfNeeded(); err != nil {
		return nil, err
	}

	assetsPage := a.cache.GetPage(page, size)
	return &assetsPage, nil
}

// Get specific asset by id
func (a Assets) GetAssetById(id string) (*coinapi.Asset, error) {
	if err := a.updateCacheIfNeeded(); err != nil {
		return nil, err
	}

	return a.cache.GetAsset(id), nil
}

func (a Assets) updateCacheIfNeeded() error {
	if a.cache.IsExpired() {
		assets, err := a.client.GetAssets()
		if err != nil {
			return fmt.Errorf("error retrieving assets from external api: %v", err)
		}
		a.cache.Fill(assets)
		log.Println("Updated cache")
	}

	return nil
}

// Calculates the gain if all quantity is sold now
func (a Assets) Valuate(ua model.UserAsset) (float64, error) {
	asset, err := a.GetAssetById(ua.AssetId)
	if err != nil {
		return -1, err
	}
	if asset == nil {
		return -1, fmt.Errorf("there is no asset with id %s", ua.AssetId)
	}

	return asset.PriceUSD * ua.Quantity, nil
}
