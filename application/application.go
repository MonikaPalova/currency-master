// Package application setups the application
package application

import (
	"log"
	"net/http"

	"github.com/MonikaPalova/currency-master/auth"
	"github.com/MonikaPalova/currency-master/config"
	. "github.com/MonikaPalova/currency-master/db"
	"github.com/MonikaPalova/currency-master/handlers"
	"github.com/MonikaPalova/currency-master/svc"
	"github.com/gorilla/mux"
	"github.com/robfig/cron"
)

type Application struct {
	db     *Database
	router *mux.Router
	svc    *svc.Service
	config *config.App
	auth   *mux.Router
}

// Application construtor
func New() Application {
	var a Application
	a.initDB()
	a.config = config.NewApp()
	a.svc = svc.NewSvc(a.db)
	a.setupHTTP()
	a.triggerSessionsCleaner()

	return a
}

// Initializes application's Database field
func (a *Application) initDB() {
	var err error
	a.db, err = NewDB()
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Successful database connection!")
}

// Starts application
func (a Application) Start() error {
	log.Println("Starting server")
	return http.ListenAndServe(a.config.Host+":"+a.config.Port, a.router)
}

func (a *Application) setupHTTP() {
	a.router = mux.NewRouter()

	a.auth = a.router.NewRoute().Subrouter()
	sessionAuth := auth.SessionAuth{Svc: a.svc.SSvc, Config: config.NewSession()}
	a.auth.Use(sessionAuth.Middleware)

	a.setupAuthHandler()
	a.setupAssetsHandler()
	a.setupUsersHandler()
	a.setupUserAssetsHandler()
	a.setupAcquisitionsHandler()
}

func (a *Application) setupAuthHandler() {
	authHandler := handlers.AuthHandler{USvc: a.svc.USvc, SSvc: a.svc.SSvc}
	a.router.Path("/login").Methods(http.MethodPost).HandlerFunc(authHandler.Login)
	a.auth.Path("/logout").Methods(http.MethodPost).HandlerFunc(authHandler.Logout)
}

func (a *Application) setupAssetsHandler() {
	assetsHandler := handlers.AssetsHandler{Svc: a.svc.ASvc}
	a.router.Path(a.config.AssetsApiV1).Methods(http.MethodGet).HandlerFunc(assetsHandler.GetAll)
	a.router.Path(a.config.AssetsApiV1 + "/{id}").Methods(http.MethodGet).HandlerFunc(assetsHandler.GetById)
}

func (a *Application) setupUsersHandler() {
	usersHandler := handlers.UsersHandler{Svc: a.svc.USvc}
	a.router.Path(a.config.UsersApiV1).Methods(http.MethodGet).HandlerFunc(usersHandler.GetAll)
	a.router.Path(a.config.UsersApiV1 + "/{username}").Methods(http.MethodGet).HandlerFunc(usersHandler.GetByUsername)
	a.router.Path(a.config.UsersApiV1).Methods(http.MethodPost).HandlerFunc(usersHandler.Post)
}

func (a *Application) setupUserAssetsHandler() {
	userAssetsHandler := handlers.UserAssetsHandler{ASvc: a.svc.ASvc, USvc: a.svc.USvc, UaSvc: a.svc.UaSvc, ADB: a.db.AcquisitionsDBHandler}
	a.router.Path(a.config.UserAssetsApiV1).Methods(http.MethodGet).HandlerFunc(userAssetsHandler.GetAll)
	a.router.Path(a.config.UserAssetsApiV1 + "/{id}").Methods(http.MethodGet).HandlerFunc(userAssetsHandler.GetByID)
	a.auth.Path(a.config.UserAssetsApiV1 + "/{id}/buy").Methods(http.MethodPost).HandlerFunc(userAssetsHandler.Buy)
	a.auth.Path(a.config.UserAssetsApiV1 + "/{id}/sell").Methods(http.MethodPost).HandlerFunc(userAssetsHandler.Sell)
}

func (a *Application) setupAcquisitionsHandler() {
	acquisitionsHandler := handlers.AcquisitionsHandler{DB: a.db.AcquisitionsDBHandler}
	a.router.Path(a.config.AcquisitionsApiV1).Methods(http.MethodGet).HandlerFunc(acquisitionsHandler.GetAll)
}

func (a Application) triggerSessionsCleaner() {
	c := cron.New()
	c.AddFunc("@hourly", func() { a.svc.SSvc.ClearExpired() })
	c.Start()
}
