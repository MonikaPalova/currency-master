package svc

import (
	"github.com/MonikaPalova/currency-master/coinapi"
	"github.com/MonikaPalova/currency-master/config"
	"github.com/MonikaPalova/currency-master/db"
	"github.com/MonikaPalova/currency-master/model"
)

// object that contains all services used in project
type Service struct {
	ASvc  *Assets
	USvc  *Users
	UaSvc *UserAssets
	SSvc  *Sessions
}

// cosntructor
func NewSvc(db *db.Database) *Service {
	aSvc := NewAssets(coinapi.NewClient())
	uSvc := &Users{UDB: db.UsersDBHandler, v: valuator{svc: aSvc}}
	uaSvc := &UserAssets{UaDB: db.UserAssetsDBHandler, v: valuator{svc: aSvc}}
	sSvc := &Sessions{sessions: map[string]model.Session{}, config: config.NewSession()}

	return &Service{ASvc: aSvc, USvc: uSvc, UaSvc: uaSvc, SSvc: sSvc}
}
