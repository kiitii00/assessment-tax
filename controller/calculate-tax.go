package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}
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
	Tax       float64    `json:"tax"`
	TaxLevels []TaxLevel `json:"taxLevels"`
	TaxRefund float64    `json:"taxRefund"`
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
	kReceiptFound := false
	// Add up all allowances
	// for _, allowance := range input.Allowances {
	// 	totalDeductions += allowance.Amount
	// fmt.Println("ลดหย่อน:",totalDeductions)
	// }
	for _, allowance := range input.Allowances {
		if allowance.AllowanceType == "k-receipt" {
			kReceiptFound = true
			if allowance.Amount > 50000 {
				totalDeductions += 50000
			} else {
				totalDeductions += allowance.Amount
			}
		} else if allowance.AllowanceType == "donation" && allowance.Amount > 100000 {
			totalDeductions += 100000
		} else {
			totalDeductions += allowance.Amount
		}

		fmt.Println("ลดหย่อน:", totalDeductions, "$")
	}
	if !kReceiptFound {
		totalDeductions += 0.0
	}
	// Ensure minimum tax reduction for personal deduction
	if totalDeductions < 10000 {
		totalDeductions = 10000
	}

	// Calculate taxable income
	taxableIncome := input.TotalIncome - totalDeductions
	fmt.Println(taxableIncome)

	taxLevels := []TaxLevel{
		{Level: "0-150,000", Tax: 0.0},
		{Level: "150,001-500,000", Tax: 0.0},
		{Level: "500,001-1,000,000", Tax: 0.0},
		{Level: "1,000,001-2,000,000", Tax: 0.0},
		{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
	}
	// Calculate tax
	var tax float64

	switch {
	case taxableIncome <= 150000:
		taxLevels[0].Tax = 0
		tax = 0
	case taxableIncome <= 500000:
		taxLevels[1].Tax = (taxableIncome - 150000) * 0.10
		tax = taxLevels[1].Tax
	case taxableIncome <= 1000000:
		taxLevels[2].Tax = (taxableIncome - 500000) * 0.15
		tax = taxLevels[2].Tax
	case taxableIncome <= 2000000:
		taxLevels[3].Tax = (taxableIncome - 1000000) * 0.20
		tax = taxLevels[3].Tax
	default:
		taxLevels[4].Tax = (taxableIncome - 2000000) * 0.35
		tax = taxLevels[4].Tax
	}

	tax -= input.WHT

	return TaxCalculationResult{
		Tax:       tax,
		TaxLevels: taxLevels,
	}
}
