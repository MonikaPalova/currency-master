package application

import (
	"log"
	"net/http"

	. "github.com/MonikaPalova/currency-master/db"
	"github.com/MonikaPalova/currency-master/handlers"
	"github.com/MonikaPalova/currency-master/svc"
	"github.com/gorilla/mux"
)

const (
	usersApiV1        = "/api/v1/users"
	assetsApiV1       = "/api/v1/assets"
	userAssetsApiV1   = "/api/v1/users/{username}/assets"
	acquisitionsApiV1 = "/api/v1/acquisitions"
)

type Application struct {
	db     *Database
	router *mux.Router
	ASvc   *svc.Assets
	USvc   *svc.Users
}

func New() Application {
	var a Application
	a.initDB()
	a.ASvc = svc.NewAssets()
	a.USvc = &svc.Users{ASvc: a.ASvc, UDB: a.db.UsersDBHandler, UaDB: a.db.UserAssetsDBHandler}
	a.setupHTTP()
	// setup app
	return a
}

func (a *Application) initDB() {
	var err error
	a.db, err = NewDB()
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Successful database connection!")
}

func (a Application) Start() error {
	return http.ListenAndServe(":7777", a.router)
}

func (a *Application) setupHTTP() {
	a.router = mux.NewRouter()
	a.setupAssetsHandler()
	a.setupUsersHandler()
	a.setupUserAssetsHandler()
	a.setupAcquisitionsHandler()
}

func (a *Application) setupAssetsHandler() {
	assetsHandler := handlers.AssetsHandler{Svc: a.ASvc}
	a.router.Path(assetsApiV1).Methods(http.MethodGet).HandlerFunc(assetsHandler.GetAll)
	a.router.Path(assetsApiV1 + "/{id}").Methods(http.MethodGet).HandlerFunc(assetsHandler.GetById)
}

func (a *Application) setupUsersHandler() {
	usersHandler := handlers.UsersHandler{Svc: a.USvc}
	a.router.Path(usersApiV1).Methods(http.MethodGet).HandlerFunc(usersHandler.GetAll)
	a.router.Path(usersApiV1 + "/{username}").Methods(http.MethodGet).HandlerFunc(usersHandler.GetByUsername)
	a.router.Path(usersApiV1).Methods(http.MethodPost).HandlerFunc(usersHandler.Post)
}

func (a *Application) setupUserAssetsHandler() {
	userAssetsHandler := handlers.UserAssetsHandler{ASvc: a.ASvc, USvc: a.USvc}
	a.router.Path(userAssetsApiV1).Methods(http.MethodGet).HandlerFunc(userAssetsHandler.GetAll)
	a.router.Path(userAssetsApiV1 + "/{id}").Methods(http.MethodGet).HandlerFunc(userAssetsHandler.GetByID)
	a.router.Path(userAssetsApiV1 + "/{id}/buy").Methods(http.MethodPost).HandlerFunc(userAssetsHandler.Buy)
	a.router.Path(userAssetsApiV1 + "/{id}/sell").Methods(http.MethodPost).HandlerFunc(userAssetsHandler.Sell)
}

func (a *Application) setupAcquisitionsHandler() {
	acquisitionsHandler := handlers.AcquisitionsHandler{DB: a.db.AcquisitionsDBHandler}
	a.router.Path(acquisitionsApiV1).Methods(http.MethodGet).HandlerFunc(acquisitionsHandler.GetAll)
}
