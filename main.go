package main

import (
	"log"

	"github.com/kiitii00/assessment-tax/routes"
)

func main() {

	echo := routes.SetupRoutes()
	err := echo.Start((":8080"))
	if err != nil {
		log.Fatal(err)
	}
}
