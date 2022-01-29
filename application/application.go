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

func (a *Application) Start() error {
	return http.ListenAndServe(":7777", a.router)
}

func (a *Application) setupHTTP() {
	a.router = mux.NewRouter()

	a.setupAssetsHandler()
	a.setupUsersHandler()
}

func (a *Application) setupAssetsHandler() {
	assetsHandler := handlers.NewAssetsHandler()
	a.router.Path("/assets").Methods(http.MethodGet).HandlerFunc(assetsHandler.GetAll)
	a.router.Path("/assets/{id}").Methods(http.MethodGet).HandlerFunc(assetsHandler.GetById)
}

func (a *Application) setupUsersHandler() {
	usersHandler := handlers.UsersHandler{DB: a.db.UsersDBHandler}
	a.router.Path("/users").Methods(http.MethodGet).HandlerFunc(usersHandler.GetAll)
	a.router.Path("/users/{username}").Methods(http.MethodGet).HandlerFunc(usersHandler.GetByUsername)
	a.router.Path("/users").Methods(http.MethodPost).HandlerFunc(usersHandler.Post)
}
