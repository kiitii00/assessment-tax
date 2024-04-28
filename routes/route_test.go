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
		amount         float64
		expectedStatus int
		expectedBody   interface{} // Use interface{} to handle any JSON response
	}{
		{
			name:           "Case 1: Valid amount",
			amount:         70000.0,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]float64{
				"personalDeduction": 70000.0,
			},
		},
	}

	// Create a new echo instance
	e := echo.New()

	// Set up request and recorder
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", strings.NewReader(`{"amount": 70000.0}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Call the handler
			HandlePersonalDeduction(c)

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
