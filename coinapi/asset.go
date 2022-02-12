package coinapi

import (
	"encoding/json"
	"fmt"
)

type Asset struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	IsCrypto bool    `json:"isCrypto"`
	PriceUSD float64 `json:"priceUSD"`
}

type AssetPage struct {
	Assets []Asset `json:"assets"`
	Page   int     `json:"page"`
	Size   int     `json:"size"`
	Total  int     `json:"totalResults"`
}

func (a *Asset) UnmarshalJSON(bytes []byte) (err error) {
	var asset struct {
		ID       string  `json:"asset_id"`
		Name     string  `json:"name"`
		IsCrypto float64 `json:"type_is_crypto"`
		PriceUSD float64 `json:"price_usd"`
	}
	if err = json.Unmarshal(bytes, &asset); err != nil {
		fmt.Println(string(bytes))
		fmt.Println(err.Error())
		return err
	}

	a.ID = asset.ID
	a.Name = asset.Name
	a.IsCrypto = asset.IsCrypto != 0
	a.PriceUSD = asset.PriceUSD

	return err
}

func (a Asset) String() string {
	bytes, _ := json.Marshal(a)
	return string(bytes)
}
