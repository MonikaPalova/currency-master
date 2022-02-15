package svc

import "github.com/MonikaPalova/currency-master/db"

// object that contains all services used in project
type Service struct {
	ASvc *Assets
	USvc *Users
}

// cosntructor
func NewSvc(db *db.Database) *Service {
	aSvc := NewAssets()
	uSvc := &Users{ASvc: aSvc, UDB: db.UsersDBHandler, UaDB: db.UserAssetsDBHandler}

	return &Service{ASvc: aSvc, USvc: uSvc}
}
