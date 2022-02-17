package model

import "time"

// information about asset owned by user
type UserAsset struct {
	//username
	Username string `json:"username,omitempty"`

	// asset id
	AssetId string `json:"assetId"`

	//asset name
	Name string `json:"name"`

	//quantity of asset owned by user
	Quantity float64 `json:"quantity"`

	// the usd value of the quantity if sold now
	Valuation float64 `json:"valuation"`
}

// information about a specific asset purchase - receipt
type Acquisition struct {
	Username string    `json:"username"`
	AssetId  string    `json:"assetId"`
	Quantity float64   `json:"quantity"`
	PriceUSD float64   `json:"priceUSD"`
	TotalUSD float64   `json:"totalUSD"`
	Created  time.Time `json:"purchaseDate"`
}
