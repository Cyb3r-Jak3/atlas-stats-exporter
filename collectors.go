package main

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/Cyb3r-Jak3/atlas-stats-exporter/pkg/version"
	"github.com/prometheus/client_golang/prometheus"
)

func BuildInfoCollector() prometheus.Collector {
	goVersion := "unknown"
	if buildInfo, available := debug.ReadBuildInfo(); available {
		goVersion = buildInfo.GoVersion
	}
	return prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "atlas_exporter_build_info",
			Help: "Build information about the Atlas exporter",
			ConstLabels: prometheus.Labels{
				"version":    version.Version,
				"commit":     version.Commit,
				"date":       version.Date,
				"go_version": goVersion,
			},
		},
		func() float64 {
			return 1
		},
	)
}

// RefreshTracker records the last successful refresh of each collector.
type RefreshTracker struct {
	mu   sync.Mutex
	last map[string]time.Time
	desc *prometheus.Desc
}

func NewRefreshTracker() *RefreshTracker {
	return &RefreshTracker{
		last: make(map[string]time.Time),
		desc: prometheus.NewDesc(
			"atlas_exporter_last_successful_refresh_timestamp",
			"Unix timestamp of the last successful refresh of data from the Atlas API, per collector",
			[]string{"collector"},
			nil,
		),
	}
}

func (t *RefreshTracker) record(name string, at time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.last[name] = at
}

func (t *RefreshTracker) Describe(ch chan<- *prometheus.Desc) {
	ch <- t.desc
}

// Collect only emits collectors that have refreshed successfully at least once.
func (t *RefreshTracker) Collect(ch chan<- prometheus.Metric) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for name, at := range t.last {
		ch <- prometheus.MustNewConstMetric(t.desc, prometheus.GaugeValue, float64(at.Unix()), name)
	}
}

// cachedCollector fetches fresh data from the Atlas API on every scrape. If the fetch fails,
// the metrics from the last successful fetch are served instead.
type cachedCollector struct {
	ctx       context.Context
	name      string
	descs     []*prometheus.Desc
	fetch     func(ctx context.Context) ([]prometheus.Metric, error)
	refreshes *RefreshTracker

	// mu is held for the whole refresh so concurrent scrapes don't duplicate upstream requests.
	mu      sync.Mutex
	metrics []prometheus.Metric
}

func (c *cachedCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, d := range c.descs {
		ch <- d
	}
}

func (c *cachedCollector) Collect(ch chan<- prometheus.Metric) {
	c.mu.Lock()
	defer c.mu.Unlock()
	logger.Debugf("Refreshing %s", c.name)
	metrics, err := c.fetch(c.ctx)
	if err != nil {
		logger.WithError(err).Errorf("Failed to refresh %s, serving previous metrics", c.name)
	} else {
		c.metrics = metrics
		c.refreshes.record(c.name, time.Now())
	}
	for _, m := range c.metrics {
		ch <- m
	}
}

func CreditsCollector(ctx context.Context, refreshes *RefreshTracker) prometheus.Collector {
	desc := prometheus.NewDesc(
		"atlas_exporter_credits",
		"Current number of credits available in the Atlas account",
		nil,
		nil,
	)
	return &cachedCollector{
		ctx:       ctx,
		name:      "credits",
		descs:     []*prometheus.Desc{desc},
		refreshes: refreshes,
		fetch: func(ctx context.Context) ([]prometheus.Metric, error) {
			resp, err := AtlasAPIClient.GetCredits(ctx)
			if err != nil {
				return nil, err
			}
			if resp == nil {
				return nil, fmt.Errorf("received nil response from GetCredits")
			}
			return []prometheus.Metric{
				prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(resp.CurrentBalance)),
			}, nil
		},
	}
}

func ProbeLastConnectedCollectorFactory(ctx context.Context, refreshes *RefreshTracker) prometheus.Collector {
	desc := prometheus.NewDesc(
		"atlas_exporter_probe_last_connected",
		"Last connected time (Unix timestamp) for each probe",
		[]string{"probe_id", "country_code", "description"},
		nil,
	)
	return &cachedCollector{
		ctx:       ctx,
		name:      "probe_last_connected",
		descs:     []*prometheus.Desc{desc},
		refreshes: refreshes,
		fetch: func(ctx context.Context) ([]prometheus.Metric, error) {
			resp, err := AtlasAPIClient.GetMyProbes(ctx)
			if err != nil {
				return nil, err
			}
			metrics := make([]prometheus.Metric, 0, len(resp))
			for _, probe := range resp {
				metrics = append(metrics, prometheus.MustNewConstMetric(
					desc,
					prometheus.GaugeValue,
					float64(probe.LastConnected),
					fmt.Sprintf("%d", probe.ID),
					probe.CountryCode,
					probe.Description,
				))
			}
			return metrics, nil
		},
	}
}

func ProbeMeasurementsCollectorFactory(ctx context.Context, refreshes *RefreshTracker) prometheus.Collector {
	desc := prometheus.NewDesc(
		"atlas_exporter_probe_measurements",
		"Measurements for each probe",
		[]string{"probe_id", "type", "status"},
		nil,
	)
	return &cachedCollector{
		ctx:       ctx,
		name:      "probe_measurements",
		descs:     []*prometheus.Desc{desc},
		refreshes: refreshes,
		fetch: func(ctx context.Context) ([]prometheus.Metric, error) {
			resp, err := AtlasAPIClient.GetMyProbesMeasurements(ctx)
			if err != nil {
				return nil, err
			}
			matrix := make(map[int]map[string]map[string]int)
			for _, measurement := range resp {
				if _, ok := matrix[measurement.ProbeID]; !ok {
					matrix[measurement.ProbeID] = make(map[string]map[string]int)
				}
				if _, ok := matrix[measurement.ProbeID][measurement.Type]; !ok {
					matrix[measurement.ProbeID][measurement.Type] = make(map[string]int)
				}
				matrix[measurement.ProbeID][measurement.Type][measurement.Status]++
			}
			var metrics []prometheus.Metric
			for probeID, probeMeasurements := range matrix {
				for typ, statuses := range probeMeasurements {
					for status, count := range statuses {
						metrics = append(metrics, prometheus.MustNewConstMetric(
							desc,
							prometheus.GaugeValue,
							float64(count),
							fmt.Sprintf("%d", probeID),
							typ,
							status,
						))
					}
				}
			}
			return metrics, nil
		},
	}
}
