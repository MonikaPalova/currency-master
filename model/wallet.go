package model

import "github.com/MonikaPalova/currency-master/coinapi"

type Wallet struct {
	assets       *[]coinapi.Asset
	acquisitions *[]Acquisition

	assetsClient *coinapi.Client
}

type Acquisition struct {
	asset    *coinapi.Asset
	quantity float64
	priceUSD float64
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

func (a Acquisition) Total() float64 {
	return a.quantity * a.priceUSD
}
