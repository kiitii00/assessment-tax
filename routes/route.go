package routes

import (
    "net/http"

    "github.com/kiitii00/assessment-tax/controller"
    "github.com/labstack/echo/v4"
)

func SetupRoutes() *echo.Echo {
    e := echo.New()
    e.POST("/tax/calculations", controller.HandleTaxCalculations)
    e.POST("/admin/deductions/personal", HandlePersonalDeduction)
    e.POST("/admin/deductions/k-receipt", HandlekReceipt)  
    return e
}

func HandlePersonalDeduction(c echo.Context) error {
    type RequestBody struct {
        Amount float64 `json:"amount"`
    }

    var reqBody RequestBody
    if err := c.Bind(&reqBody); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    // Process deduction logic
    personalDeduction := reqBody.Amount

    return c.JSON(http.StatusOK, map[string]float64{"personalDeduction": personalDeduction})
}

func HandlekReceipt(c echo.Context) error {
    type RequestBody struct {
        Amount float64 `json:"amount"`
    }

    var reqBody RequestBody
    if err := c.Bind(&reqBody); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    // Process deduction logic
    kReceiptDeduction := reqBody.Amount

    return c.JSON(http.StatusOK, map[string]float64{"kReceipt": kReceiptDeduction})
}