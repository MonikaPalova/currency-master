package svc

import "github.com/MonikaPalova/currency-master/db"

type Service struct {
	ASvc *Assets
	USvc *Users
}

func NewSvc(db *db.Database) *Service {
	aSvc := NewAssets()
	uSvc := &Users{ASvc: aSvc, UDB: db.UsersDBHandler, UaDB: db.UserAssetsDBHandler}

	return &Service{ASvc: aSvc, USvc: uSvc}
}
