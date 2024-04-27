package routes

import (
	_ "encoding/json"

	"github.com/kiitii00/assessment-tax/controller"
	_ "github.com/kiitii00/assessment-tax/models"
	"github.com/labstack/echo/v4"
)

func SetupRoutes() *echo.Echo {
	e := echo.New()
	e.POST("/tax/calculations", controller.HandleTaxCalculations)
	return e
}
