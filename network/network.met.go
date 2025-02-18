package network

import "github.com/prometheus/client_golang/prometheus"

// 누적 네트워크 바이트 메트릭 (카운터)
var (
	NetRecv = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "network_receive_bytes_total",
		Help: "Total number of bytes received on all network interfaces",
	})
	NetSent = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "network_transmit_bytes_total",
		Help: "Total number of bytes sent on all network interfaces",
	})
)
