package coinapi

import (
	"encoding/json"
)

// Asset object received from external api
type Asset struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	IsCrypto bool    `json:"isCrypto"`
	PriceUSD float64 `json:"priceUSD"`
}

// Page with assets
type AssetPage struct {
	Assets []Asset `json:"assets"`
	Page   int     `json:"page"`
	Size   int     `json:"size"`
	Total  int     `json:"totalResults"`
}

// formats asset object received from external api to Asset object
func (a *Asset) UnmarshalJSON(bytes []byte) (err error) {
	var asset struct {
		ID       string  `json:"asset_id"`
		Name     string  `json:"name"`
		IsCrypto float64 `json:"type_is_crypto"`
		PriceUSD float64 `json:"price_usd"`
	}
	if err = json.Unmarshal(bytes, &asset); err != nil {
		return err
	}

	a.ID = asset.ID
	a.Name = asset.Name
	a.IsCrypto = asset.IsCrypto != 0
	a.PriceUSD = asset.PriceUSD

	return err
}
