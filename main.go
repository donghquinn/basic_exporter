package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	psnet "github.com/shirou/gopsutil/v4/net"
	"org.donghyuns.com/exporter/network/network"
	"org.donghyuns.com/exporter/network/system"
)

// 네트워크 증분 계산을 위한 이전 값 저장
var prevNetRecv, prevNetSent uint64

// 시스템 메트릭 업데이트 함수
func updateMetrics() {
	// 메모리 정보 수집
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Failed to get memory info: %v", err)
	} else {
		system.MemoryTotal.Set(float64(vmStat.Total))
		system.MemoryUsed.Set(float64(vmStat.Used))
		system.MemoryAvailable.Set(float64(vmStat.Available))
		system.MemoryUsagePercent.Set(vmStat.UsedPercent)
	}

	// CPU 사용률 (0초 간격은 즉시 반환하므로, 참고용으로 사용)
	cpuPercents, err := cpu.Percent(0, false)
	if err != nil {
		log.Printf("Failed to get CPU percent: %v", err)
	} else if len(cpuPercents) > 0 {
		system.CpuUsage.Set(cpuPercents[0])
	}

	// Load Average 수집
	loadStat, err := load.Avg()
	if err != nil {
		log.Printf("Failed to get load average: %v", err)
	} else {
		system.Load1.Set(loadStat.Load1)
		system.Load5.Set(loadStat.Load5)
		system.Load15.Set(loadStat.Load15)
	}

	// 네트워크 정보 수집 (모든 인터페이스의 합산)
	netCounters, err := psnet.IOCounters(false)
	if err != nil || len(netCounters) == 0 {
		log.Printf("Failed to get network counters: %v", err)
	} else {
		currentRecv := netCounters[0].BytesRecv
		currentSent := netCounters[0].BytesSent

		// 최초 호출인 경우 이전값이 없으므로 바로 설정
		if prevNetRecv == 0 && prevNetSent == 0 {
			network.NetRecv.Add(float64(currentRecv))
			network.NetSent.Add(float64(currentSent))
		} else {
			// 누적 값의 증가분만 추가 (Prometheus Counter는 감소 불가)
			if currentRecv >= prevNetRecv {
				network.NetRecv.Add(float64(currentRecv - prevNetRecv))
			}
			if currentSent >= prevNetSent {
				network.NetSent.Add(float64(currentSent - prevNetSent))
			}
		}
		prevNetRecv = currentRecv
		prevNetSent = currentSent
	}
}

func recordMetrics() {
	// 주기적으로 메트릭 업데이트
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			updateMetrics()
		}
	}()
}

func main() {
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

	// 메트릭 업데이트 시작
	recordMetrics()

	// /metrics endpoint 노출
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Exporter listening on :8080/metrics")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
