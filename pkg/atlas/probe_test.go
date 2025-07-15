package atlas

import (
	"context"
	"net/http"
	"testing"

	"github.com/Cyb3r-Jak3/common/v5"
)

func TestAPI_GetMyProbes(t *testing.T) {
	setup()
	defer teardown()

	// Mock the API response for GetMyProbes
	mux.HandleFunc("/probes/my", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{
   			"count": 1,
			"next": null,
			"previous": null,
			"results": [
				{
            "address_v4": "45.138.229.91",
            "address_v6": "2a10:3781:e22:1:220:4aff:fec8:23d7",
            "asn_v4": 206238,
            "asn_v6": 206238,
            "country_code": "NL",
            "description": "Robert #1 100/10 Freedom.nl",
            "firmware_version": 4790,
            "first_connected": 1288367583,
            "geometry": {
                "type": "Point",
                "coordinates": [
                    4.9275,
                    52.3475
                ]
            },
            "id": 1,
            "is_anchor": false,
            "is_public": true,
            "last_connected": 1752371567,
            "prefix_v4": "45.138.228.0/22",
            "prefix_v6": "2a10:3780::/29",
            "status": {
                "id": 1,
                "name": "Connected",
                "since": "2025-07-11T14:37:28Z"
            },
            "status_since": 1752244648,
            "tags": [
                {
                    "name": "Home",
                    "slug": "home"
                },
                {
                    "name": "NAT",
                    "slug": "nat"
                },
                {
                    "name": "Native IPv6",
                    "slug": "native-ipv6"
                },
                {
                    "name": "IPv6",
                    "slug": "ipv6"
                },
                {
                    "name": "system: V1",
                    "slug": "system-v1"
                },
                {
                    "name": "system: IPv6 Capable",
                    "slug": "system-ipv6-capable"
                },
                {
                    "name": "system: IPv4 RFC1918",
                    "slug": "system-ipv4-rfc1918"
                },
                {
                    "name": "system: IPv4 Capable",
                    "slug": "system-ipv4-capable"
                },
                {
                    "name": "freedom nl",
                    "slug": "freedom-nl"
                },
                {
                    "name": "system: IPv4 Works",
                    "slug": "system-ipv4-works"
                },
                {
                    "name": "system: Resolves A Correctly",
                    "slug": "system-resolves-a-correctly"
                },
                {
                    "name": "system: IPv6 Works",
                    "slug": "system-ipv6-works"
                },
                {
                    "name": "system: Resolves AAAA Correctly",
                    "slug": "system-resolves-aaaa-correctly"
                },
                {
                    "name": "system: IPv4 Stable 1d",
                    "slug": "system-ipv4-stable-1d"
                }
            ],
            "total_uptime": 449427153,
            "type": "Probe"
        }
			]
		}`))
		if err != nil {
			return
		}
	})
	apiResponse, err := client.GetMyProbes(context.Background())
	if err != nil {
		t.Fatalf("GetMyProbes failed: %v", err)
	}
	if len(apiResponse) != 1 {
		t.Errorf("Expected 1 probe, got %d", len(apiResponse))
	}
	if apiResponse[0].ID != 1 {
		t.Errorf("Expected probe ID 1, got %d", apiResponse[0].ID)
	}
	if apiResponse[0].CountryCode != "NL" {
		t.Errorf("Expected country code 'NL', got '%s'", apiResponse[0].CountryCode)
	}
	if apiResponse[0].Description != "Robert #1 100/10 Freedom.nl" {
		t.Errorf("Expected description 'Robert #1 100/10 Freedom.nl', got '%s'", apiResponse[0].Description)
	}
}

func TestAPI_GetProbeMeasurements(t *testing.T) {
	setup()
	defer teardown()

	// Mock the API response for GetMyProbes
	mux.HandleFunc("/probes/my", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{
   			"count": 1,
			"next": null,
			"previous": null,
			"results": [
				{
            "address_v4": "45.138.229.91",
            "address_v6": "2a10:3781:e22:1:220:4aff:fec8:23d7",
            "asn_v4": 206238,
            "asn_v6": 206238,
            "country_code": "NL",
            "description": "Robert #1 100/10 Freedom.nl",
            "firmware_version": 4790,
            "first_connected": 1288367583,
            "geometry": {
                "type": "Point",
                "coordinates": [
                    4.9275,
                    52.3475
                ]
            },
            "id": 1,
            "is_anchor": false,
            "is_public": true,
            "last_connected": 1752371567,
            "prefix_v4": "45.138.228.0/22",
            "prefix_v6": "2a10:3780::/29",
            "status": {
                "id": 1,
                "name": "Connected",
                "since": "2025-07-11T14:37:28Z"
            },
            "status_since": 1752244648,
            "tags": [
                {
                    "name": "Home",
                    "slug": "home"
                },
                {
                    "name": "NAT",
                    "slug": "nat"
                },
                {
                    "name": "Native IPv6",
                    "slug": "native-ipv6"
                },
                {
                    "name": "IPv6",
                    "slug": "ipv6"
                },
                {
                    "name": "system: V1",
                    "slug": "system-v1"
                },
                {
                    "name": "system: IPv6 Capable",
                    "slug": "system-ipv6-capable"
                },
                {
                    "name": "system: IPv4 RFC1918",
                    "slug": "system-ipv4-rfc1918"
                },
                {
                    "name": "system: IPv4 Capable",
                    "slug": "system-ipv4-capable"
                },
                {
                    "name": "freedom nl",
                    "slug": "freedom-nl"
                },
                {
                    "name": "system: IPv4 Works",
                    "slug": "system-ipv4-works"
                },
                {
                    "name": "system: Resolves A Correctly",
                    "slug": "system-resolves-a-correctly"
                },
                {
                    "name": "system: IPv6 Works",
                    "slug": "system-ipv6-works"
                },
                {
                    "name": "system: Resolves AAAA Correctly",
                    "slug": "system-resolves-aaaa-correctly"
                },
                {
                    "name": "system: IPv4 Stable 1d",
                    "slug": "system-ipv4-stable-1d"
                }
            ],
            "total_uptime": 449427153,
            "type": "Probe"
        }
			]
		}`))
		if err != nil {
			return
		}
	})
	mux.HandleFunc("/probes/1/measurements", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{
   			"count": 1,
			"next": null,
			"previous": null,
			"results": [
				{
					"id": "12345",
					"type": "ping",
					"description": "Ping measurement",
					"status": "completed",
					"start_time": "2025-07-11T14:37:28Z",
					"stop_time": "2025-07-11T14:37:30Z",
					"target": "example.com"
				}
			]
		}`))
		if err != nil {
			return
		}
	},
	)
	apiResponse, err := client.GetMyProbesMeasurements(context.Background())
	if err != nil {
		t.Fatalf("GetMyProbesMeasurements failed: %v", err)
	}
	if len(apiResponse) != 1 {
		t.Errorf("Expected 1 measurement, got %d", len(apiResponse))
	}
	StartTime, err := common.ParseResilientTime("2025-07-11T14:37:28Z")
	if err != nil {
		t.Fatalf("Failed to parse start time: %v", err)
	}
	StopTime, err := common.ParseResilientTime("2025-07-11T14:37:30Z")
	if err != nil {
		t.Fatalf("Failed to parse stop time: %v", err)
	}
	expectedMeasurement := ProbeInfoMeasurement{
		ProbeID:       1,
		MeasurementID: "12345",
		Type:          "ping",
		Description:   "Ping measurement",
		Status:        "completed",
		StartTime:     StartTime,
		StopTime:      StopTime,
		Target:        "example.com",
	}
	if apiResponse[0] != expectedMeasurement {
		t.Errorf("Expected measurement %+v, got %+v", expectedMeasurement, apiResponse[0])
	}
}
