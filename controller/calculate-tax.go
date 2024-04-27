package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TaxCalculationInput struct {
	TotalIncome float64     `json:"totalIncome"`
	WHT         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

// Allowance represents an allowance type and amount
type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type TaxCalculationResult struct {
	Tax       float64 `json:"tax"`
	TaxRefund float64 `json:"taxRefund"`
}

// commit#1
// HandleTaxCalculations handles the HTTP POST request for tax calculations
func HandleTaxCalculations(c echo.Context) error {
	// Parse request body
	var input TaxCalculationInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	fmt.Println("commit#1")
	// Perform tax calculation
	result := CalculateTax(input)

	// Return JSON response
	return c.JSON(http.StatusOK, result)
}

// calculateTax calculates tax based on the input data
func CalculateTax(input TaxCalculationInput) TaxCalculationResult {
	totalDeductions := 60000.0 //input.WHT

	// Add up all allowances
	// for _, allowance := range input.Allowances {
	// 	totalDeductions += allowance.Amount
	// }

	// Calculate taxable income

	taxableIncome := input.TotalIncome - totalDeductions

	// Calculate tax
	var tax float64

	switch {
	case taxableIncome <= 150000:
		tax = 0
	case taxableIncome <= 500000:
		tax = (taxableIncome - 150000) * 0.10
	case taxableIncome <= 1000000:
		tax = (500000-150000)*0.10 + (taxableIncome-500000)*0.15
	case taxableIncome <= 2000000:
		tax = (500000-150000)*0.10 + (1000000-500000)*0.15 + (taxableIncome-1000000)*0.20
	default:
		tax = (500000-150000)*0.10 + (1000000-500000)*0.15 + (2000000-1000000)*0.20 + (taxableIncome-2000000)*0.35
	}

	// Apply maximum tax reduction for donation
	for _, allowance := range input.Allowances {
		if allowance.AllowanceType == "donation" && allowance.Amount > 100000 {
			tax -= 100000
		} else {
			tax -= allowance.Amount
		}
	}

	// Ensure minimum tax reduction for personal deduction
	if totalDeductions < 10000 {
		totalDeductions = 10000
	}

	// Calculate tax refund
	taxRefund := 0.0
	if input.TotalIncome-totalDeductions-input.WHT < 0 {
		taxRefund = -(input.TotalIncome - totalDeductions - input.WHT - tax)
	}

	return TaxCalculationResult{

		Tax:       tax,
		TaxRefund: taxRefund,
	}
}
