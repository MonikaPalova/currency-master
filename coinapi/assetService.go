package coinapi

import "fmt"

type AssetService struct {
	cache  *cache
	client *client
}

func NewAssetService() *AssetService {
	return &AssetService{cache: newCache(), client: newClient()}
}

func (a AssetService) GetAssetPage(page, size int) (*AssetPage, error) {
	if err := a.updateCacheIfNeeded(); err != nil {
		return nil, err
	}

	assetsPage := a.cache.getPage(page, size)
	return &assetsPage, nil
}

func (a AssetService) GetAssetById(id string) (*Asset, error) {
	if err := a.updateCacheIfNeeded(); err != nil {
		return nil, err
	}

	return a.cache.getAsset(id), nil
}

func (a AssetService) updateCacheIfNeeded() error {
	if a.cache.isExpired() {
		fmt.Println("UPDATING CACHE")
		assets, err := a.client.getAssets()
		if err != nil {
			return fmt.Errorf("error retrieving assets from external api: %v", err)
		}
		a.cache.fill(assets)
	}

	return nil
}
