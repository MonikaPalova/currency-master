package application

import (
	"log"
	"net/http"

	. "github.com/MonikaPalova/currency-master/db"
	"github.com/MonikaPalova/currency-master/handlers"
	"github.com/gorilla/mux"
)

const (
	user     = "root"
	password = ""
	dbname   = "currency-master"

	USERS_API_V1       = "/api/v1/users"
	ASSETS_API_V1      = "/api/v1/assets"
	USER_ASSETS_API_V1 = "/api/v1/users/{username}/assets"
)

type Application struct {
	db     *Database
	router *mux.Router
}

func New() Application {
	var a Application
	a.initDB()
	a.setupHTTP()
	// setup app
	return a
}

func (a *Application) initDB() {
	var err error
	a.db, err = NewDB(user, password)
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
	assetsHandler := handlers.NewAssetsHandler()
	a.router.Path(ASSETS_API_V1).Methods(http.MethodGet).HandlerFunc(assetsHandler.GetAll)
	a.router.Path(ASSETS_API_V1 + "/{id}").Methods(http.MethodGet).HandlerFunc(assetsHandler.GetById)
}

func (a *Application) setupUsersHandler() {
	usersHandler := handlers.UsersHandler{DB: a.db.UsersDBHandler}
	a.router.Path(USERS_API_V1).Methods(http.MethodGet).HandlerFunc(usersHandler.GetAll)
	a.router.Path(USERS_API_V1 + "/{username}").Methods(http.MethodGet).HandlerFunc(usersHandler.GetByUsername)
	a.router.Path(USERS_API_V1).Methods(http.MethodPost).HandlerFunc(usersHandler.Post)
}

func (a *Application) setupUserAssetsHandler() {
	userAssetsHandler := handlers.UserAssetsHanlder(a.db.UserAssetsDBHandler)

	a.router.Path(USER_ASSETS_API_V1).Methods(http.MethodGet).HandlerFunc(userAssetsHandler.GetAll)
	a.router.Path(USER_ASSETS_API_V1 + "/{id}").Methods(http.MethodGet).HandlerFunc(userAssetsHandler.GetByID)
	a.router.Path(USER_ASSETS_API_V1 + "/{id}/buy").Methods(http.MethodPost).HandlerFunc(userAssetsHandler.Buy)
	a.router.Path(USER_ASSETS_API_V1 + "/{id}/sell").Methods(http.MethodPost).HandlerFunc(userAssetsHandler.Sell)
}
