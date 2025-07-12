package atlas

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type ProbeAPIResponse struct {
	Pagination Pagination  `json:"pagination"`
	Results    []ProbeInfo `json:"results"`
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

// LastConnectedTime Helper to get LastConnected as time.Time.
func (p *ProbeInfo) LastConnectedTime() time.Time {
	return time.Unix(int64(p.LastConnected), 0)
}

func (api *API) GetMyProbes(ctx context.Context) ([]ProbeInfo, error) {
	if api.APIToken == "" {
		return nil, ErrMissingToken
	}
	nextPage := true
	var probes []ProbeInfo
	var P Pagination
	for nextPage {
		resp, err := api.request(ctx, "GET", buildURI("/probes/my", P), nil, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get probes: %w", err)
		}
		var probeResponse ProbeAPIResponse
		if err := json.Unmarshal(resp.Body, &probeResponse); err != nil {
			return nil, fmt.Errorf("failed to unmarshal probes response: %w", err)
		}
		probes = append(probes, probeResponse.Results...)
		if probeResponse.Pagination.Done() {
			nextPage = false
		} else {
			P = probeResponse.Pagination.Next()
		}
	}
	return probes, nil
}
