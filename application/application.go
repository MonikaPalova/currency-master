package application

import (
	"log"
	"net/http"

	"github.com/MonikaPalova/currency-master/coinapi"
	. "github.com/MonikaPalova/currency-master/db"
	"github.com/MonikaPalova/currency-master/handlers"
	"github.com/gorilla/mux"
)

const (
	usersApiV1      = "/api/v1/users"
	assetsApiV1     = "/api/v1/assets"
	userAssetsApiV1 = "/api/v1/users/{username}/assets"
)

type Application struct {
	db            *Database
	coinapiClient *coinapi.Client
	router        *mux.Router
}

func New() Application {
	var a Application
	a.initDB()
	a.coinapiClient = coinapi.NewClient()
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
}

func (a *Application) setupAssetsHandler() {
	assetsHandler := handlers.AssetsHandler{Client: a.coinapiClient}
	a.router.Path(assetsApiV1).Methods(http.MethodGet).HandlerFunc(assetsHandler.GetAll)
	a.router.Path(assetsApiV1 + "/{id}").Methods(http.MethodGet).HandlerFunc(assetsHandler.GetById)
}

func (a *Application) setupUsersHandler() {
	usersHandler := handlers.UsersHandler{DB: a.db.UsersDBHandler}
	a.router.Path(usersApiV1).Methods(http.MethodGet).HandlerFunc(usersHandler.GetAll)
	a.router.Path(usersApiV1 + "/{username}").Methods(http.MethodGet).HandlerFunc(usersHandler.GetByUsername)
	a.router.Path(usersApiV1).Methods(http.MethodPost).HandlerFunc(usersHandler.Post)
}

func (a *Application) setupUserAssetsHandler() {
	userAssetsHandler := handlers.UserAssetsHandler{UaDB: a.db.UserAssetsDBHandler, UDB: a.db.UsersDBHandler, Client: a.coinapiClient}

	a.router.Path(userAssetsApiV1).Methods(http.MethodGet).HandlerFunc(userAssetsHandler.GetAll)
	a.router.Path(userAssetsApiV1 + "/{id}").Methods(http.MethodGet).HandlerFunc(userAssetsHandler.GetByID)
	a.router.Path(userAssetsApiV1 + "/{id}/buy").Methods(http.MethodPost).HandlerFunc(userAssetsHandler.Buy)
	a.router.Path(userAssetsApiV1 + "/{id}/sell").Methods(http.MethodPost).HandlerFunc(userAssetsHandler.Sell)
}
