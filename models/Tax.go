package models

// TaxCalculationInput represents the data sent in the tax calculation request
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