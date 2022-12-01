package main

import (
	"fmt"
	"github.com/clambin/httpserver"
	"github.com/clambin/solaredge-exporter/collector"
	"github.com/clambin/solaredge-exporter/version"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
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

func parseOptions() error {
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
		return err
	}

	if Debug {
		log.SetLevel(log.DebugLevel)
	}
	return nil
}

func main() {
	if err := Main(); err != nil {
		log.Fatal(err)
	}
}
func Main() (err error) {
	if err = parseOptions(); err != nil {
		return err
	}

	log.WithField("version", version.BuildVersion).Info("solaredge-exporter started")

	coll := collector.New(APIKey)
	prometheus.MustRegister(coll)

	// Run initialized & runs the metrics
	var server *httpserver.Server
	if server, err = httpserver.New(httpserver.WithPort{Port: Port}, httpserver.WithPrometheus{}); err == nil {
		err = server.Run()
	}
	if err != nil {
		return fmt.Errorf("failed to start Prometheus metrics server: %w", err)
	}

	log.WithField("version", version.BuildVersion).Info("solaredge-exporter stopped")
	return nil
}
