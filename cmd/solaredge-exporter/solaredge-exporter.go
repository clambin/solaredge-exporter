package main

import (
	"errors"
	"fmt"
	"github.com/clambin/solaredge-exporter/collector"
	"github.com/clambin/solaredge-exporter/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	Port     int
	Debug    bool
	Interval time.Duration
	APIKey   string
)

func main() {
	if err := Main(); err != nil {
		slog.Error("failed to start", err)
		return
	}
}

func Main() error {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout)))
	parseOptions()

	slog.Info("solaredge-exporter started", "version", version.BuildVersion)

	coll := collector.New(APIKey)
	if err := prometheus.Register(coll); err != nil {
		return fmt.Errorf("failed register Prometheus metrics: %w", err)
	}

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(fmt.Sprintf(":%d", Port), nil); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start Prometheus metrics server: %w", err)
	}

	slog.Info("solaredge-exporter stopped")
	return nil
}

func parseOptions() {
	a := kingpin.New(filepath.Base(os.Args[0]), "solaredge-exporter")

	a.Version(version.BuildVersion)
	a.HelpFlag.Short('h')
	a.VersionFlag.Short('v')
	a.Flag("debug", "Log debug messages").Short('d').BoolVar(&Debug)
	a.Flag("port", "Prometheus listener port").Short('p').Default("8080").IntVar(&Port)
	a.Flag("interval", "Measurement interval").Short('i').Default("15m").DurationVar(&Interval)
	a.Flag("apikey", "SolarEdge API key").Short('a').Required().StringVar(&APIKey)

	if _, err := a.Parse(os.Args[1:]); err != nil {
		a.Usage(os.Args[1:])
		os.Exit(1)
	}
}
