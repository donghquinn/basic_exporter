package system

import "github.com/prometheus/client_golang/prometheus"

// 메모리 메트릭
var (
	MemoryTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_total_bytes",
		Help: "Total physical memory in bytes.",
	})
	MemoryUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_used_bytes",
		Help: "Used memory in bytes.",
	})
	MemoryAvailable = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_available_bytes",
		Help: "Available memory in bytes.",
	})
	MemoryUsagePercent = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_usage_percent",
		Help: "Memory usage percentage.",
	})
)
