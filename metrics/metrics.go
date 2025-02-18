package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"org.donghyuns.com/exporter/basic/network"
	"org.donghyuns.com/exporter/basic/system"
)

// 시스템 메트릭 업데이트 함수
func updateMetrics(intervalSeconds int) {
	system.UpdateSystemMetrics()
	network.UpdateNetwork(intervalSeconds)
}

func recordMetrics(interval time.Duration) {
	// 주기적으로 메트릭 업데이트
	ticker := time.NewTicker(interval * time.Second)
	go func() {
		for range ticker.C {
			updateMetrics(int(interval))
		}
	}()
}

func MetricsScheduler(interval time.Duration) {
	// Memory Metrics
	prometheus.MustRegister(system.MemoryTotal)
	prometheus.MustRegister(system.MemoryUsed)
	prometheus.MustRegister(system.MemoryAvailable)
	prometheus.MustRegister(system.MemoryUsagePercent)

	// CPU Metrics
	prometheus.MustRegister(system.CpuUsage)

	// Load Metrics
	prometheus.MustRegister(system.Load1)
	prometheus.MustRegister(system.Load5)
	prometheus.MustRegister(system.Load15)

	// Network Metrics
	prometheus.MustRegister(network.NetRecv)
	prometheus.MustRegister(network.NetSent)
	prometheus.MustRegister(network.NetworkRecvBps)
	prometheus.MustRegister(network.NetworkSentBps)

	// 메트릭 업데이트 시작
	recordMetrics(interval)
}
