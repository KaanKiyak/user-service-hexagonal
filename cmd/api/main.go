package main

import (
	"log"
	"user-service-hexagonal/internal/app"
)

func main() {

	app := app.SetupApp()
	//Server başlat

	log.Fatal(app.Listen(":8080"))

}
