package system

import "github.com/prometheus/client_golang/prometheus"

var (
	CpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_usage_percent",
		Help: "Total CPU usage percentage.",
	})
)
