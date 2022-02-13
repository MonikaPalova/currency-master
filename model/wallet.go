package model

import "time"

// To get assets for user you get this class from db
type UserAsset struct {
	Username  string  `json:"username,omitempty"`
	AssetId   string  `json:"assetId"`
	Name      string  `json:"name"`
	Quantity  float64 `json:"quantity"`
	Valuation float64 `json:"valuation"`
	// TODO : add usdSpent, usdEarned
}

type Acquisition struct {
	Username string    `json:"username"`
	AssetId  string    `json:"assetId"`
	Quantity float64   `json:"quantity"`
	PriceUSD float64   `json:"priceUSD"`
	Created  time.Time `json:"purchaseDate"`
}
