package main

import (
	"log"

	"github.com/kiitii00/assessment-tax/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Basic Auth Middleware
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Check if the provided username and password are correct
		if username == "adminTax" && password == "admin!" {
			return true, nil
		}
		return false, nil
	}))
	echo := routes.SetupRoutes()
	err := echo.Start((":8080"))
	if err != nil {
		log.Fatal(err)
	}
}
