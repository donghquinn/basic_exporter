package system

import "github.com/prometheus/client_golang/prometheus"

var (
	Load1 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "load1",
		Help: "1 minute load average.",
	})
	Load5 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "load5",
		Help: "5 minutes load average.",
	})
	Load15 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "load15",
		Help: "15 minutes load average.",
	})
)
