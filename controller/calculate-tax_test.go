package controller

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateTax(t *testing.T) {
	testCases := []struct {
		name        string
		totalIncome float64
		allowances  []Allowance
		expectedTax float64
		wht         float64
		taxLevels   []TaxLevel
	}{
		{
			name:        "Case 1: Total income 500,000 with no allowances",
			totalIncome: 500000,
			wht:         0,
			allowances:  nil,
			expectedTax: 29000,
		},
		{
			name:        "Case 2: Total income 500,000 with WHT 25000",
			totalIncome: 500000,
			wht:         25000,
			allowances:  nil,
			expectedTax: 4000,
		},
		{
			name:        "Case 3: Total income 500,000 with donation allowance",
			totalIncome: 500000,
			allowances: []Allowance{
				{
					AllowanceType: "donation",
					Amount:        200000,
				},
			},
			expectedTax: 19000,
		},
		{
			name:        "Case 4: Total income 500,000 with donation allowance output taxLevel",
			totalIncome: 500000,
			wht:         0,
			allowances: []Allowance{
				{
					AllowanceType: "donation",
					Amount:        200000,
				},
			},
			expectedTax: 19000,
			taxLevels: []TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 19000.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			},
		},
		{
			name:        "Case 5: Total income 500,000 with donation&&k-receipt output taxLevel",
			totalIncome: 500000,
			wht:         0,
			allowances: []Allowance{
				{
					AllowanceType: "k-receipt",
					Amount:        200000,
				},
				{
					AllowanceType: "donation",
					Amount:        100000,
				},
			},
			expectedTax: 14000,
			taxLevels: []TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 14000.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CalculateTax(TaxCalculationInput{
				TotalIncome: tc.totalIncome,
				WHT:         tc.wht,
				Allowances:  tc.allowances,
			})

			assert.Equal(t, tc.expectedTax, result.Tax)

			// Assert tax levels
			if len(tc.taxLevels) > 0 {
				assert.Equal(t, tc.taxLevels, result.TaxLevels)
			}
			fmt.Println(tc.totalIncome)
		})
	}
}
