package application

import (
	"net/http"

	"github.com/MonikaPalova/currency-master/handlers"
	"github.com/gorilla/mux"
)

type Application struct {
	router *mux.Router
}

func New() Application {
	var a Application
	a.setupHTTP()
	// setup app
	return a
}

func (a *Application) Start() error {
	return http.ListenAndServe(":7777", a.router)
}

func (a *Application) setupHTTP() {
	a.router = mux.NewRouter()

	a.setupAssetsHandler()
}

func (a *Application) setupAssetsHandler() {
	assetsHandler := handlers.NewAssetsHandler()
	a.router.Path("/assets").Methods(http.MethodGet).HandlerFunc(assetsHandler.Get)
	a.router.Path("/assets/{id}").Methods(http.MethodGet).HandlerFunc(assetsHandler.GetById)
}
