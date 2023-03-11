package main

import (
	"context"
	"errors"
	"github.com/clambin/solaredge"
	"github.com/clambin/solaredge-exporter/collector"
	"github.com/clambin/solaredge-exporter/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
)

var (
	cmd = cobra.Command{
		Use:     "solaredge-exporter",
		Short:   "exports SolarEdge metrics to Prometheus",
		Run:     Main,
		Version: version.BuildVersion,
	}
)

func main() {
	if err := cmd.Execute(); err != nil {
		slog.Error("failed to start", "err", err)
		return
	}
}

func Main(_ *cobra.Command, _ []string) {
	if viper.GetBool("debug") {
		opts := slog.HandlerOptions{Level: slog.LevelDebug}
		slog.SetDefault(slog.New(opts.NewTextHandler(os.Stderr)))
	}

	sites, err := getSites()
	if err != nil {
		slog.Error("failed to get SolarEdge sites", "err", err)
		return
	}

	slog.Info("solaredge-exporter started", "version", version.BuildVersion)

	coll := collector.Collector{Sites: sites}
	if err := prometheus.Register(&coll); err != nil {
		slog.Error("failed register Prometheus metrics", "err", err)
		return
	}

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(viper.GetString("addr"), nil); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start Prometheus metrics server", "err", err)
	}

	slog.Info("solaredge-exporter stopped")
	return
}

func getSites() ([]collector.Site, error) {
	c := solaredge.Client{
		Token:      viper.GetString("apikey"),
		HTTPClient: http.DefaultClient,
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
