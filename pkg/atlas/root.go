package atlas

import (
	"context"
	"encoding/json"
	"fmt"
)

type RootAPIResponse struct {
	Measurements       string          `json:"measurements"`
	Anchors            string          `json:"anchors"`
	AnchorMeasurements string          `json:"anchor-measurements"`
	Probes             json.RawMessage `json:"probes"`
}

func (api *API) GetRootData(ctx context.Context) (*RootAPIResponse, error) {
	if api.APIToken == "" {
		return nil, ErrMissingToken
	}
	resp, err := api.request(ctx, "GET", "/root/?format=json", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get root data: %w", err)
	}
	var rootResponse RootAPIResponse
	if err = json.Unmarshal(resp.Body, &rootResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal root data response: %w", err)
	}
	return &rootResponse, nil
}
