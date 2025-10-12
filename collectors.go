package main

import (
	"context"
	"fmt"
	"runtime/debug"
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

func CreditsCollector(ctx context.Context, timeout int) prometheus.Collector {
	return prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "atlas_exporter_credits",
			Help: "Current number of credits available in the Atlas account",
		},
		func() float64 {
			cancelCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
			defer cancel()
			resp, err := AtlasAPIClient.GetCredits(cancelCtx)
			if err != nil {
				logger.WithError(err).Error("Failed to get credits")
				return 0
			}
			if resp == nil {
				logger.Error("Received nil response from GetCredits")
				return 0
			}
			return float64(resp.CurrentBalance)
		},
	)
}

type ProbeLastConnectedCollector struct {
	ctx     context.Context
	timeout int
}

func (c *ProbeLastConnectedCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc(
		"atlas_exporter_probe_last_connected",
		"Last connected time (Unix timestamp) for each probe",
		[]string{"probe_id", "country_code", "description"},
		nil,
	)
}

func (c *ProbeLastConnectedCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(c.ctx, time.Duration(c.timeout)*time.Second)
	defer cancel()
	logger.Debug("Collecting last connected time (Unix timestamp) for each probe")
	resp, err := AtlasAPIClient.GetMyProbes(ctx)
	if err != nil {
		logger.WithError(err).Error("Failed to get probe last connected")
		return
	}
	desc := prometheus.NewDesc(
		"atlas_exporter_probe_last_connected",
		"Last connected time (Unix timestamp) for each probe",
		[]string{"probe_id", "country_code", "description"},
		nil,
	)
	for _, probe := range resp {
		ch <- prometheus.MustNewConstMetric(
			desc,
			prometheus.GaugeValue,
			float64(probe.LastConnected),
			fmt.Sprintf("%d", probe.ID),
			probe.CountryCode,
			probe.Description,
		)
	}
}

func ProbeLastConnectedCollectorFactory(ctx context.Context, timeout int) prometheus.Collector {
	return &ProbeLastConnectedCollector{timeout: timeout, ctx: ctx}
}

type ProbeMeasurementsCollector struct {
	ctx     context.Context
	timeout int
}

func (c *ProbeMeasurementsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc(
		"atlas_exporter_probe_measurements",
		"Measurements for each probe",
		[]string{"probe_id", "type", "status"},
		nil,
	)
}

func (c *ProbeMeasurementsCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(c.ctx, time.Duration(c.timeout)*time.Second)
	defer cancel()
	logger.Debug("Collecting measurements for each probe")
	resp, err := AtlasAPIClient.GetMyProbesMeasurements(ctx)
	if err != nil {
		logger.WithError(err).Error("Failed to get probe measurements")
		return
	}
	desc := prometheus.NewDesc(
		"atlas_exporter_probe_measurements",
		"Measurements for each probe",
		[]string{"probe_id", "type", "status"},
		nil,
	)
	matrix := make(map[int]map[string]map[string]int)
	for _, measurement := range resp {
		probeID := measurement.ProbeID
		typ := measurement.Type
		status := measurement.Status

		if _, ok := matrix[probeID]; !ok {
			matrix[probeID] = make(map[string]map[string]int)
		}
		if _, ok := matrix[probeID][typ]; !ok {
			matrix[probeID][typ] = make(map[string]int)
		}
		matrix[probeID][typ][status]++
	}
	for probeID, probeMeasurements := range matrix {
		for typ, statuses := range probeMeasurements {
			for status, count := range statuses {
				ch <- prometheus.MustNewConstMetric(
					desc,
					prometheus.GaugeValue,
					float64(count),
					fmt.Sprintf("%d", probeID),
					typ,
					status,
				)
			}
		}
	}
}

func ProbeMeasurementsCollectorFactory(ctx context.Context, timeout int) prometheus.Collector {
	return &ProbeMeasurementsCollector{timeout: timeout, ctx: ctx}
}
