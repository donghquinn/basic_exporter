package network

import "github.com/prometheus/client_golang/prometheus"

// 네트워크 대역폭 메트릭 (실시간 Gauge: bytes per second)
var (
	NetworkRecvBps = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "network_receive_bytes_per_second",
		Help: "Current network receive bandwidth in bytes per second.",
	})
	NetworkSentBps = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "network_transmit_bytes_per_second",
		Help: "Current network transmit bandwidth in bytes per second.",
	})
)
