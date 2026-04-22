package atlas

import (
	"context"
	"net/http"
	"testing"
)

func TestAPI_GetRootData(t *testing.T) {
	setup()
	defer teardown()

	// Mock the API response for GetCredits
	mux.HandleFunc("/root/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{
			"measurements":"https://atlas.ripe.net/api/v2/measurements/?format=json",
			"probes":["https://atlas.ripe.net/api/v2/probes/?format=json",{"tags":"https://atlas.ripe.net/api/v2/probes/tags/?format=json"}],
			"anchors":"https://atlas.ripe.net/api/v2/anchors/?format=json",
			"anchor-measurements":"https://atlas.ripe.net/api/v2/anchor-measurements/?format=json"
		}`))
		if err != nil {
			return
		}
	})
	apiResponse, err := client.GetRootData(context.Background())
	if err != nil {
		t.Fatalf("GetCredits failed: %v", err)
	}
	if apiResponse.Measurements != "https://atlas.ripe.net/api/v2/measurements/?format=json" {
		t.Errorf("Expected measurements URL 'https://atlas.ripe.net/api/v2/measurements/?format=json', got '%s'", apiResponse.Measurements)
	}
	if apiResponse.Anchors != "https://atlas.ripe.net/api/v2/anchors/?format=json" {
		t.Errorf("Expected anchors URL 'https://atlas.ripe.net/api/v2/anchors/?format=json', got '%s'", apiResponse.Anchors)
	}
	if apiResponse.AnchorMeasurements != "https://atlas.ripe.net/api/v2/anchor-measurements/?format=json" {
		t.Errorf("Expected anchor measurements URL 'https://atlas.ripe.net/api/v2/anchor-measurements/?format=json', got '%s'", apiResponse.AnchorMeasurements)
	}
}
