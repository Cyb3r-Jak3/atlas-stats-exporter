package atlas

import (
	"context"
	"net/http"
	"testing"
)

func TestAPI_GetCredits(t *testing.T) {
	setup()
	defer teardown()

	// Mock the API response for GetCredits
	mux.HandleFunc("/credits", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{
			"current_balance": 1000,
			"credit_checked": true,
			"max_daily_credits": 5000,
			"estimated_daily_income": 200,
			"estimated_daily_expenditure": 150,
			"estimated_daily_balance": 1050,
			"calculation_time": "2023-10-01T00:00:00Z"
		}`))
		if err != nil {
			return
		}
	})
	apiResponse, err := client.GetCredits(context.Background())
	if err != nil {
		t.Fatalf("GetCredits failed: %v", err)
	}
	if apiResponse.CurrentBalance != 1000 {
		t.Errorf("Expected current balance 1000, got %d", apiResponse.CurrentBalance)
	}
}
