package atlas

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Cyb3r-Jak3/common/v5"
)

type ProbeAPIResponse struct {
	Count   int         `json:"count"`
	Next    string      `json:"next"`
	Results []ProbeInfo `json:"results"`
}

type ProbeInfoGeometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type ProbeInfoStatus struct {
	ID    int       `json:"id"`
	Name  string    `json:"name"`
	Since time.Time `json:"since"`
}

type ProbeInfoTags struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type ProbeInfo struct {
	AddressV4       string            `json:"address_v4"`
	AddressV6       string            `json:"address_v6"`
	ASNv4           int               `json:"asn_v4"`
	ASNv6           int               `json:"asn_v6"`
	CountryCode     string            `json:"country_code"`
	Description     string            `json:"description"`
	FirmwareVersion int               `json:"firmware_version"`
	FirstConnected  int               `json:"first_connected"`
	Geometry        ProbeInfoGeometry `json:"geometry"`
	ID              int               `json:"id"`
	Anchor          bool              `json:"is_anchor"`
	Public          bool              `json:"is_public"`
	LastConnected   int               `json:"last_connected"`
	PrefixV4        string            `json:"prefix_v4"`
	PrefixV6        string            `json:"prefix_v6"`
	Status          ProbeInfoStatus   `json:"status"`
	StatusSince     int               `json:"status_since"`
	TotalUptime     int               `json:"total_uptime"`
	Type            string            `json:"type"`
}

type ProbeInfoMeasurement struct {
	ProbeID       int
	MeasurementID string               `json:"id"`
	Type          string               `json:"type"`
	Description   string               `json:"description"`
	Status        string               `json:"status"`
	StartTime     common.ResilientTime `json:"start_time"`
	StopTime      common.ResilientTime `json:"stop_time"`
	Target        string               `json:"target"`
}

type ProbeMeasurementResponse struct {
	Count   int                    `json:"count"`
	Next    string                 `json:"next"`
	Results []ProbeInfoMeasurement `json:"results"`
}

// LastConnectedTime Helper to get LastConnected as time.Time.
func (p *ProbeInfo) LastConnectedTime() time.Time {
	return time.Unix(int64(p.LastConnected), 0)
}

func (api *API) GetMyProbes(ctx context.Context) ([]ProbeInfo, error) {
	if api.APIToken == "" {
		return nil, ErrMissingToken
	}
	var probes []ProbeInfo
	page := 1
	for {
		resp, err := api.request(ctx, "GET", fmt.Sprintf("/probes/my?page=%d", page), nil, nil)
		if err != nil {
			return nil, fmt.Errorf("API request failed: %w", err)
		}
		var probeResponse ProbeAPIResponse
		if unmarshalErr := json.Unmarshal(resp.Body, &probeResponse); unmarshalErr != nil {
			return nil, fmt.Errorf("failed to unmarshal probes response: %w", unmarshalErr)
		}
		probes = append(probes, probeResponse.Results...)
		if probeResponse.Count == 0 || probeResponse.Next == "" {
			break
		}
		page++
	}
	return probes, nil
}

func (api *API) GetMyProbesMeasurements(ctx context.Context) ([]ProbeInfoMeasurement, error) {
	if api.APIToken == "" {
		return nil, ErrMissingToken
	}
	myProbes, err := api.GetMyProbes(ctx)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	var probeMeasurements []ProbeInfoMeasurement
	for _, probe := range myProbes {
		page := 1
		for {
			resp, respErr := api.request(ctx, "GET", fmt.Sprintf("/probes/%d/measurements?page=%d", probe.ID, page), nil, nil)
			if respErr != nil {
				return nil, fmt.Errorf("API request failed to get probe measurements, %d: %w", probe.ID, respErr)
			}
			var probeMeasurementResponse ProbeMeasurementResponse
			if unmarshalErr := json.Unmarshal(resp.Body, &probeMeasurementResponse); unmarshalErr != nil {
				return nil, fmt.Errorf("failed to unmarshal probe measurement response: %w", unmarshalErr)
			}
			for i := range probeMeasurementResponse.Results {
				probeMeasurementResponse.Results[i].ProbeID = probe.ID
			}
			probeMeasurements = append(probeMeasurements, probeMeasurementResponse.Results...)
			if probeMeasurementResponse.Count == 0 || probeMeasurementResponse.Next == "" {
				break
			}
			page++
		}
	}
	return probeMeasurements, nil
}
