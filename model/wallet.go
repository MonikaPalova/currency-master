package model

import "time"

// To get assets for user you get this class from db
type UserAsset struct {
	Username string  `json:"username,omitempty"`
	AssetId  string  `json:"assetId"`
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
}

type Acquisition struct {
	Username string    `json:"username"`
	AssetId  string    `json:"assetId"`
	Quantity float64   `json:"quantity"`
	PriceUSD float64   `json:"priceUSD"`
	Created  time.Time `json:"purchaseDate"`
}

// func (w *Wallet) GetValuation() float64 {
// 	// TODO
// }

// func (w *Wallet) Buy(asset Asset, quantity float64, maxPrice float64) (Acquisition, error) {
// 	// TODO
// }

// func (w *Wallet) Sell(asset Asset, quantity float64, minPrice float64) (float64, error) {
// 	// TODO
// }

// func (w *Wallet) GetAcquisitions() ([]Acquisition, error) {
// 	//TODO
// }
