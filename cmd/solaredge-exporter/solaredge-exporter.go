package main

import (
	"context"
	"errors"
	"github.com/clambin/go-common/http/metrics"
	"github.com/clambin/go-common/http/roundtripper"
	"github.com/clambin/solaredge"
	"github.com/clambin/solaredge-exporter/internal/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var (
	version = "change-me"
	cmd     = cobra.Command{
		Use:     "solaredge-exporter",
		Short:   "exports SolarEdge metrics to Prometheus",
		Run:     Main,
		Version: version,
	}
)

func main() {
	if err := cmd.Execute(); err != nil {
		slog.Error("failed to start", "err", err)
		return
	}
}

func Main(cmd *cobra.Command, _ []string) {
	if viper.GetBool("debug") {
		opts := slog.HandlerOptions{Level: slog.LevelDebug}
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &opts)))
	}

	m := metrics.NewRequestSummaryMetrics("solaredge", "exporter", nil)
	prometheus.MustRegister(m)

	httpClient := http.Client{
		Transport: roundtripper.New(
			roundtripper.WithCache(roundtripper.DefaultCacheTable, 5*time.Minute, 0),
			roundtripper.WithRequestMetrics(m),
		),
		Timeout: 10 * time.Second,
	}

	sites, err := getSites(&httpClient)
	if err != nil {
		slog.Error("failed to get SolarEdge sites", "err", err)
		return
	}

	inverters, err := getInverters(sites)
	if err != nil {
		slog.Error("failed to get SolarEdge inverters", "err", err)
		return
	}

	slog.Info("solaredge-exporter started", "version", cmd.Version)

	coll := collector.Collector{Sites: sites, Inverters: inverters}
	if err := prometheus.Register(&coll); err != nil {
		slog.Error("failed register Prometheus metrics", "err", err)
		return
	}

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(viper.GetString("addr"), nil); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start Prometheus metrics server", "err", err)
	}

	slog.Info("solaredge-exporter stopped")
}

func getSites(httpClient *http.Client) ([]collector.Site, error) {
	c := solaredge.Client{
		Token:      viper.GetString("apikey"),
		HTTPClient: httpClient,
	}

	sites, err := c.GetSites(context.Background())
	if err != nil {
		return nil, err
	}

	var result []collector.Site
	for _, site := range sites {
		result = append(result, &site)
	}

	return result, nil
}

func getInverters(sites []collector.Site) (map[int][]collector.Inverter, error) {
	inverters := make(map[int][]collector.Inverter)
	for _, site := range sites {
		result, err := site.GetInverters(context.Background())
		if err != nil {
			return nil, err
		}
		invertersAtSite := make([]collector.Inverter, len(result))
		for index := range result {
			invertersAtSite[index] = collector.Inverter(&result[index])
		}
		inverters[site.GetID()] = invertersAtSite
	}
	return inverters, nil
}

func init() {
	cobra.OnInitialize(initConfig)
	cmd.Flags().Bool("debug", false, "Log debug messages")
	_ = viper.BindPFlag("debug", cmd.Flags().Lookup("debug"))
	cmd.Flags().String("addr", ":9090", "Listener address for Prometheus metrics")
	_ = viper.BindPFlag("addr", cmd.Flags().Lookup("addr"))
	cmd.Flags().String("apikey", "", "SolarEdge API key")
	_ = viper.BindPFlag("apikey", cmd.Flags().Lookup("apikey"))
}

func initConfig() {
	viper.AddConfigPath("/etc/solaredge-exporter/")
	viper.AddConfigPath("$HOME/.solaredge-exporter")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	viper.SetDefault("debug", false)
	viper.SetDefault("addr", ":9090")

	viper.SetEnvPrefix("SOLAREDGE")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()
}
