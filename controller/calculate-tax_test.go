package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	_"strings"
	"github.com/labstack/echo/v4"
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

func TestUploadTaxFile(t *testing.T) {
	// Prepare a sample CSV file content
	csvContent := `totalIncome,wht,donation
	500000,0,0
	600000,40000,20000
	750000,50000,15000`

	// Create a request body with the sample CSV content
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("taxFile", "taxes.csv")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte(csvContent))
	writer.Close()

	// Create a mock HTTP request
	req := httptest.NewRequest(http.MethodPost, "/upload-csv", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())

	// Create a mock HTTP response recorder
	rec := httptest.NewRecorder()

	// Create a new Echo instance
	e := echo.New()

	// Create a context for the request
	c := e.NewContext(req, rec)
	fmt.Println("c",c)
	// Call the UploadTaxFile handler
	if err := UploadTaxFile(c); err != nil {
		t.Fatalf("failed to upload tax file: %v", err)
	}
	
	// Check the response status code
	assert.Equal(t, http.StatusOK, rec.Code)

	// Define the expected response body
	expectedResponse := []Tax{
		{TotalIncome: 500000, Tax: 29000},
		{TotalIncome: 600000, Tax: 33000},
		{TotalIncome: 750000, Tax: 53750},
	}

	// Unmarshal the response body JSON
	var response struct {
		Taxes []Tax `json:"taxes"`
	}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}
	// fmt.Println("t",t)
	// fmt.Println("expectedResponse",expectedResponse)
	// fmt.Println("response.Taxes",response.Taxes)
	// Check the response body against the expected response
	assert.ElementsMatch(t, expectedResponse, response.Taxes)
	
}

