package main

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"os"
	"strings"
	"testing"
)

func TestBuildInfoCollector(t *testing.T) {
	collector := BuildInfoCollector()
	if err := testutil.CollectAndCompare(collector, strings.NewReader(`
# HELP atlas_exporter_build_info Build information about the Atlas exporter
# TYPE atlas_exporter_build_info gauge
atlas_exporter_build_info{commit="unknown",date="unknown",go_version="go1.24.4",version="unknown"} 1
`)); err != nil {
		t.Errorf("BuildInfoCollector failed: %v", err)
	}
}

func TestCreditsCollector(t *testing.T) {
	if os.Getenv("ATLAS_EXPORTER_API_TOKEN") == "" {
		t.Skip("Skipping TestCreditsCollector because ATLAS_EXPORTER_API_TOKEN is not set")
	}
	ctx := context.Background()
	collector := CreditsCollector(ctx, 10)
	if err := testutil.CollectAndCompare(collector, strings.NewReader(`
# HELP atlas_exporter_credits Current number of credits available in the Atlas account
# TYPE atlas_exporter_credits gauge
atlas_exporter_credits 1000
`)); err != nil {
		t.Errorf("CreditsCollector failed: %v", err)
	}
}

func TestProbeLastConnectedCollector(t *testing.T) {
	if os.Getenv("ATLAS_EXPORTER_API_TOKEN") == "" {
		t.Skip("Skipping TestProbeLastConnectedCollector because ATLAS_EXPORTER_API_TOKEN is not set")
	}
	ctx := context.Background()
	collector := ProbeLastConnectedCollector{
		ctx:     ctx,
		timeout: 10,
	}
	if err := testutil.CollectAndCompare(&collector, strings.NewReader(`
# HELP atlas_exporter_probe_last_connected Last time the probe was connected
# TYPE atlas_exporter_probe_last_connected gauge
atlas_exporter_probe_last_connected{probe_id="12345"} 1700000000
}`)); err != nil {
		t.Errorf("ProbeLastConnectedCollector failed: %v", err)
	}
}
