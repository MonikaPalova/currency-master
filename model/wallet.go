package model

// To get assets for user you get this class from db
type UserAsset struct {
	username string
	assetId  string
	quantity uint8
}

type Acquisition struct {
	username string
	assetId  string
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
