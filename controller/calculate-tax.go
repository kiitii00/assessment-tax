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

// HandleTaxCalculations handles the HTTP POST request for tax calculations
func HandleTaxCalculations(c echo.Context) error {
	var input TaxCalculationInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	result := CalculateTax(input)

	return c.JSON(http.StatusOK, result)
}

// CalculateTax calculates tax based on the input data
func CalculateTax(input TaxCalculationInput) TaxCalculationResult {
	totalDeductions := 60000.0

	// Add up all allowances
	// for _, allowance := range input.Allowances {
	// 	totalDeductions += allowance.Amount
	// }
	for _, allowance := range input.Allowances {
		if allowance.AllowanceType == "donation" && allowance.Amount > 100000 {
			totalDeductions += 100000
		} else {
			totalDeductions += allowance.Amount
		}
		fmt.Println(totalDeductions)
	}

	// Ensure minimum tax reduction for personal deduction
	if totalDeductions < 10000 {
		totalDeductions = 10000
	}

	// Calculate taxable income
	taxableIncome := input.TotalIncome - totalDeductions
	fmt.Println(taxableIncome)
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
	// for _, allowance := range input.Allowances {
	// 	if allowance.AllowanceType == "donation" && allowance.Amount > 100000 {
	// 		tax -= 100000
	// 	} else {
	// 		tax -= allowance.Amount
	// 	}
	// 	//fmt.Println(tax)
	// }

	// // Ensure minimum tax reduction for personal deduction
	// if totalDeductions < 10000 {
	// 	totalDeductions = 10000
	// }

	// Calculate tax refund
	taxRefund := 0.0
	if input.TotalIncome-totalDeductions-input.WHT < 0 {
		taxRefund = -(input.TotalIncome - totalDeductions - input.WHT - tax)
	}

	// tax -= input.WHT
	// if tax < 0 {
	// 	tax = 0
	// }

	return TaxCalculationResult{
		Tax:       tax,
		TaxRefund: taxRefund,
	}
}
