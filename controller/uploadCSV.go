package controller

import (
	"encoding/csv"
	"fmt"
	_ "fmt"
	"io"
	"net/http"
	_ "os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type Tax struct {
	TotalIncome float64 `json:"totalIncome"`
	Tax         float64 `json:"tax"`
}

func UploadTaxFile(c echo.Context) error {
	// Get the uploaded file
	file, err := c.FormFile("taxFile")
	if err != nil {
		return err
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create a CSV reader
	reader := csv.NewReader(src)

	// Skip the header row
	_, err = reader.Read()
	if err != nil {
		return err
	}

	// Initialize tax results slice
	var taxes []Tax

	// Read records from the CSV file
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Parse CSV fields
		totalIncome, err := strconv.ParseFloat(strings.TrimSpace(record[0]), 64)
		//fmt.Println(totalIncome)
		fmt.Println("record", strings.TrimSpace(record[0]))
		fmt.Println("record", strings.TrimSpace(record[1]))
		fmt.Println("record", strings.TrimSpace(record[2]))
		if err != nil {
			return err
		}
		wht, err := strconv.ParseFloat(strings.TrimSpace(record[1]), 64)
		if err != nil {
			return err
		}
		donation, err := strconv.ParseFloat(strings.TrimSpace(record[2]), 64)
		if err != nil {
			return err
		}

		// Calculate tax
		tax := calculateTax(totalIncome, wht, donation)

		// Append tax result to the taxes slice
		taxes = append(taxes, Tax{
			TotalIncome: totalIncome,
			Tax:         tax,
		})
	}

	// Return tax results as JSON
	return c.JSON(http.StatusOK, echo.Map{"taxes": taxes})
}

func calculateTax(totalIncome, wht, donation float64) float64 {
	fmt.Println("calculateTax",totalIncome,wht,donation)
	// Calculate taxable income (total income minus WHT and donation)
	totalDeductions := 60000.0
	taxableIncome := totalIncome - wht - donation - totalDeductions

	// Tax rates for different income levels
	var taxRate float64

	switch {
	case taxableIncome <= 150000:
		taxRate = 0.0
	case taxableIncome <= 500000:
		taxRate = (taxableIncome - 150000) * 0.10
	case taxableIncome <= 1000000:
		taxRate = (500000-150000)*0.10 + (taxableIncome-500000)*0.15
	case taxableIncome <= 2000000:
		taxRate = (1000000-500000)*0.15 + (taxableIncome-1000000)*0.20
	default:
		taxRate = (2000000-1000000)*0.20 + (taxableIncome-2000000)*0.35
	}

	// Calculate tax based on taxable income and tax rate
	tax := taxRate

	return tax
}
