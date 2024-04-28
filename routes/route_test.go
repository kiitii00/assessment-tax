// In routes package

package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPersonalDeductionHandler(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		path           string // Add path to distinguish between personal and k-receipt deductions
		amount         float64
		expectedStatus int
		expectedBody   interface{} // Use interface{} to handle any JSON response
	}{
		{
			name:           "Case 1: Valid amount for personal deduction",
			path:           "/admin/deductions/personal",
			amount:         70000.0,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]float64{
				"personalDeduction": 70000.0,
			},
		},
		{
			name:           "Case 2: Valid amount for k-receipt deduction",
			path:           "/admin/deductions/k-receipt",
			amount:         70000.0,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]float64{
				"kReceipt": 70000.0,
			},
		},
	}

	// Create a new echo instance
	e := echo.New()

	// Set up request and recorder
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tc.path, strings.NewReader(`{"amount": 70000.0}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Call the handler based on the path
			switch tc.path {
			case "/admin/deductions/personal":
				HandlePersonalDeduction(c)
			case "/admin/deductions/k-receipt":
				HandlekReceipt(c)
			}

			// Check the status code
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Marshal the expected body
			expectedJSON, err := json.Marshal(tc.expectedBody)
			if err != nil {
				t.Errorf("error marshalling expected body: %v", err)
				return
			}

			// Check the response body
			assert.JSONEq(t, string(expectedJSON), rec.Body.String())
		})
	}
}
