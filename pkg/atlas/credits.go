package atlas

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Cyb3r-Jak3/common/v5"
)

type CreditAPIResponse struct {
	CurrentBalance            int                  `json:"current_balance"`
	CreditChecked             bool                 `json:"credit_checked"`
	MaxDailyCredits           int                  `json:"max_daily_credits"`
	EstimatedDailyIncome      int                  `json:"estimated_daily_income"`
	EstimatedDailyExpenditure int                  `json:"estimated_daily_expenditure"`
	EstimatedDailyBalance     int                  `json:"estimated_daily_balance"`
	CalculationTime           common.ResilientTime `json:"calculation_time"`
	EstimatedRunOutSeconds    int                  `json:"estimated_run_out_seconds,omitempty"`
	PastDayMeasurementResults int                  `json:"past_day_measurement_results,omitempty"`
	PastDayCreditsSpending    int                  `json:"past_day_credits_spending,omitempty"`
	LastDateDebited           common.ResilientTime `json:"last_date_debited,omitempty"`
	LastDateCredited          common.ResilientTime `json:"last_date_credited,omitempty"`
}

func (api *API) GetCredits(ctx context.Context) (*CreditAPIResponse, error) {
	if api.APIToken == "" {
		return nil, ErrMissingToken
	}
	resp, err := api.request(ctx, "GET", "/credits", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get credits: %w", err)
	}
	var creditResponse CreditAPIResponse
	if err = json.Unmarshal(resp.Body, &creditResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal credits response: %w", err)
	}
	return &creditResponse, nil
}
