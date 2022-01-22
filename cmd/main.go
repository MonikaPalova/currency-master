package main

import (
	"log"

	"github.com/MonikaPalova/currency-master/application"
)

func main() {
	app := application.New()
	if err := app.Start(); err != nil {
		log.Fatalf("An error occured when starting the aplication: %s", err.Error())
	}
}
