package metrics

import (
	"github.com/mitchellh/mapstructure"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"time"
)

var (
	Version       = "dev"
	Build         = "none"
	BuildDate     = "unknown"
	PromNamespace = "bserver"
	promAlive     prometheus.Gauge
)

func init() {}

var labels prometheus.Labels
var aliveTimer *time.Ticker

func GetGlobalLabels() prometheus.Labels {
	if labels != nil {
		return labels
	}
	if viper.IsSet("metrics.labels") {
		configLabels := viper.GetStringMapString("metrics.labels")
		if len(configLabels) > 0 {
			mapstructure.Decode(configLabels, &labels)
		}
		labels["version"] = Version
		labels["build"] = Build
		labels["build_date"] = BuildDate
	}
	intiAlive()
	return labels
}

func intiAlive() {

	promAlive = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   PromNamespace,
		Name:        "info",
		Help:        "The server alive ",
		ConstLabels: labels,
	})

	prometheus.MustRegister(promAlive)
	aliveCheckInterval, _ := time.ParseDuration("1m")

	if viper.IsSet("metrics.check-alive") {
		aliveCheckInterval = viper.GetDuration("metrics.check-alive")
	}

	go func() {
		aliveTimer = time.NewTicker(aliveCheckInterval)

		for range aliveTimer.C {
			promAlive.Set(1)
		}
	}()

	promAlive.Set(1)
}
func StopAlive() {

	if aliveTimer != nil {
		aliveTimer.Stop()
	}
}
