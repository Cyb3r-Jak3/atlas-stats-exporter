package main

import (
	"context"
	"fmt"
	"net/http"
	"net/mail"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/Cyb3r-Jak3/atlas-stats-exporter/pkg/atlas"
	"github.com/Cyb3r-Jak3/atlas-stats-exporter/pkg/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

var (
	logger         = logrus.New()
	AtlasAPIClient *atlas.API
	versionString  = version.String()
)

func buildApp() *cli.Command {
	app := &cli.Command{
		Name:    "atlas_exporter",
		Usage:   "A Prometheus exporter for RIPE Atlas credits and probe statistics",
		Version: versionString,
		Suggest: true,
		Authors: []any{
			&mail.Address{
				Name:    "Cyber-Jak3",
				Address: "git@cyberjake.xyz",
			},
		},
		Action: Run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "listen_address",
				Aliases: []string{"l"},
				Usage:   "Address to listen on for HTTP requests",
				Value:   ":8080",
				Sources: cli.EnvVars("ATLAS_EXPORTER_LISTEN_ADDRESS"),
			},
			&cli.StringFlag{
				Name:    "metrics_path",
				Aliases: []string{"m"},
				Usage:   "Path to expose metrics",
				Value:   "/metrics",
				Sources: cli.EnvVars("ATLAS_EXPORTER_METRICS_PATH"),
			},
			&cli.StringFlag{
				Name:    "api_token",
				Aliases: []string{"t"},
				Usage:   "API token for authentication with the Atlas API. Can also be set via the ATLAS_API_TOKEN environment variable",
				Value:   "",
				Sources: cli.EnvVars("ATLAS_EXPORTER_API_TOKEN"),
			},
			&cli.IntFlag{
				Name:    "timeout",
				Usage:   "Timeout for API requests in seconds",
				Value:   60,
				Sources: cli.EnvVars("ATLAS_EXPORTER_TIMEOUT"),
			},
			&cli.BoolFlag{
				Name:    "tls_enabled",
				Aliases: []string{"tls"},
				Usage:   "Enable TLS for the HTTP server",
				Value:   false,
				Sources: cli.EnvVars("ATLAS_EXPORTER_TLS_ENABLED"),
			},
			&cli.StringFlag{
				Name:    "tls_cert_chain_path",
				Aliases: []string{"tls_cert"},
				Usage:   "Path to the TLS certificate chain file (PEM format)",
				Value:   "cert.pem",
				Sources: cli.EnvVars("ATLAS_EXPORTER_TLS_CERT_CHAIN_PATH"),
			},
			&cli.StringFlag{
				Name:    "tls_key_path",
				Aliases: []string{"tls_key"},
				Usage:   "Path to the TLS private key file (PEM format)",
				Value:   "key.pem",
				Sources: cli.EnvVars("ATLAS_EXPORTER_TLS_KEY_PATH"),
			},
			&cli.StringFlag{
				Name:    "log_level",
				Aliases: []string{"ll"},
				Usage:   "Set the logging level (debug, info, warn, error, fatal, panic)",
				Value:   "info",
				Sources: cli.EnvVars("ATLAS_EXPORTER_LOG_LEVEL"),
			},
			&cli.StringFlag{
				Name:    "base_url",
				Usage:   "Base URL for the Atlas API. Useful for testing or custom deployments.",
				Value:   "https://atlas.ripe.net/api/v2",
				Sources: cli.EnvVars("ATLAS_EXPORTER_BASE_URL"),
				Hidden:  true,
			},
		},
		EnableShellCompletion: true,
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	return app
}

func Run(ctx context.Context, c *cli.Command) error {
	err := SetLogLevel(c)
	if err != nil {
		logger.Fatalf("Failed to set log level: %v", err)
	}

	apiToken := strings.TrimSpace(c.String("api_token"))
	if apiToken == "" {
		logger.Fatal("API token is required. Please set the API token using the --api_token flag or the ATLAS_API_TOKEN environment variable.")
	}
	apiClientOptions := []atlas.Option{
		atlas.WithUserAgent("go-atlas-stats-exporter/" + version.Version),
		atlas.WithAPIToken(apiToken),
		atlas.WithBaseURL(c.String("base_url")),
	}

	if logger.Level >= logrus.DebugLevel {
		apiClientOptions = append(apiClientOptions, atlas.WithDebug(true))
	}
	AtlasAPIClient, err = atlas.New(apiClientOptions...)
	if err != nil {
		logger.Fatalf("Failed to initialize Atlas API client: %v", err)
	}
	scrapeTimeout := c.Int("timeout")
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		BuildInfoCollector(),
		CreditsCollector(ctx, scrapeTimeout),
		ProbeLastConnectedCollectorFactory(ctx, scrapeTimeout),
		ProbeMeasurementsCollectorFactory(ctx, scrapeTimeout),
	)
	logger.Infof("Starting atlas exporter (Version: %s)", version.Version)
	listenAddress := c.String("listen_address")
	metricsPath := c.String("metrics_path")
	tlsEnabled := c.Bool("tls_enabled")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, handleErr := w.Write([]byte(`<html>
					<head><title>Cyb3r-Jak3 RIPE Atlas Exporter (Version ` + version.Version + `)</title></head>
					<body>
					<h1> Cyb3r-Jak3 RIPE Atlas Exporter</h1>
					<h2>Example</h2>
					<p>Metrics for measurement configured in configuration file:</p>
					<p><a href="` + metricsPath + `">` + r.Host + metricsPath + `</a></p>
					<h2>More information</h2>
					<p><a href="https://github.com/Cyb3r-Jak3/atlas_exporter">github.com/Cyb3r-Jak3/atlas_exporter</a></p>
					</body>
					<footer> Commit: ` + version.Commit + `, Date: ` + version.Date + `, Version: ` + version.Version + `</footer>
					</html>`))
		if handleErr != nil {
			logger.Errorf("Failed to write response: %v", handleErr)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
	http.HandleFunc("/version", func(w http.ResponseWriter, _ *http.Request) {
		_, handleErr := w.Write([]byte(versionString))
		if handleErr != nil {
			logger.Errorf("Failed to write response: %v", handleErr)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
	http.Handle(metricsPath, promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	logger.Infof("Listening for %s on %s (TLS: %v)", metricsPath, listenAddress, tlsEnabled)
	httpServer := &http.Server{
		Addr:              listenAddress,
		ReadHeaderTimeout: 5 * time.Second,
	}
	// Channel to listen for errors from ListenAndServe
	serverErr := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		if tlsEnabled {
			serverErr <- httpServer.ListenAndServeTLS(
				c.String("tls_cert_chain_path"),
				c.String("tls_key_path"),
			)
		} else {
			serverErr <- httpServer.ListenAndServe()
		}
	}()

	// Listen for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	select {
	case <-quit:
		logger.Info("Shutting down server...")
		ShutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return httpServer.Shutdown(ShutdownContext)
	case ServerErr := <-serverErr:
		return ServerErr
	}
}

func main() {
	app := buildApp()
	err := app.Run(context.Background(), os.Args)
	if err != nil {
		fmt.Printf("Error running app: %s\n", err)
		os.Exit(1)
	}
}

func SetLogLevel(c *cli.Command) error {
	logLevel := c.String("log_level")
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("invalid log level: %s, error: %w", logLevel, err)
	}
	logger.SetLevel(level)

	logger.Debugf("Log Level set to %v", logger.Level)
	return nil
}
