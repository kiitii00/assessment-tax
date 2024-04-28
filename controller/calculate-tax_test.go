package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateTax(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		totalIncome    float64
		expectedTax    float64
	}{
		{
			name:           "Case 1: Total income 500,000",
			totalIncome:    500000,
			expectedTax:    29000,
		},
		{
			name:           "Case 2: Total income 50,000",
			totalIncome:    50000,
			expectedTax:    0.0,
		},
		
	}

	// Iterate over test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the CalculateTax function with the test case input
			result := CalculateTax(TaxCalculationInput{
				TotalIncome: tc.totalIncome,
				WHT:         0, // Assuming no withholding tax for simplicity
				Allowances:  nil, // Assuming no other allowances for simplicity
			})

			// Assert that the calculated tax matches the expected tax
			assert.Equal(t, tc.expectedTax, result.Tax)
		})
	}
}
