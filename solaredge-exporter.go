package main

import (
	"fmt"
	"github.com/clambin/solaredge-exporter/internal/exporter"
	"github.com/clambin/solaredge-exporter/internal/version"
	"github.com/clambin/solaredge-exporter/pkg/solaredge"
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
	var err error

	a := kingpin.New(filepath.Base(os.Args[0]), "solaredge-exporter")

	a.Version(version.BuildVersion)
	a.HelpFlag.Short('h')
	a.VersionFlag.Short('v')
	a.Flag("debug", "Log debug messages").Short('d').BoolVar(&Debug)
	a.Flag("port", "Prometheus listener port").Short('p').Default("8080").IntVar(&Port)
	a.Flag("interval", "Measurement interval").Short('i').Default("15m").DurationVar(&Interval)
	a.Flag("apikey", "SolarEdge API key").Short('a').Required().StringVar(&APIKey)

	_, err = a.Parse(os.Args[1:])
	if err != nil {
		a.Usage(os.Args[1:])
		os.Exit(2)
	}

	if Debug {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	parseOptions()

	log.WithField("version", version.BuildVersion).Info("solaredge-exporter started")

	go func() {
		client := solaredge.NewClient(APIKey, nil)
		err := exporter.Run(client, Interval)
		if err != nil {
			log.WithError(err).Fatal("exporter failed. aborting ...")
		}
	}()

	// Run initialized & runs the metrics
	listenAddress := fmt.Sprintf(":%d", Port)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(listenAddress, nil)
	log.WithError(err).Fatal("Failed to start Prometheus http handler")
}
