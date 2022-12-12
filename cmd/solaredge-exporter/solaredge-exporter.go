package main

import (
	"errors"
	"fmt"
	"github.com/clambin/solaredge-exporter/collector"
	"github.com/clambin/solaredge-exporter/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
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
	}

	if Debug {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	if err := Main(); err != nil {
		log.Fatal(err)
	}
}
func Main() error {
	parseOptions()

	log.WithField("version", version.BuildVersion).Info("solaredge-exporter started")

	coll := collector.New(APIKey)
	if err := prometheus.Register(coll); err != nil {
		return fmt.Errorf("failed registered Prometheus metrics: %w", err)
	}

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(fmt.Sprintf(":%d", Port), nil); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start Prometheus metrics server: %w", err)
	}

	log.WithField("version", version.BuildVersion).Info("solaredge-exporter stopped")
	return nil
}
