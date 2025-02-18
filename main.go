package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	psnet "github.com/shirou/gopsutil/v4/net"
	"org.donghyuns.com/exporter/basic/network"
	"org.donghyuns.com/exporter/basic/system"
)

// 네트워크 증분 계산을 위한 이전 값 저장
var prevNetRecv, prevNetSent uint64

// 시스템 메트릭 업데이트 함수
func updateMetrics(intervalSeconds int) {
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

		// 첫 실행 시, 누적 값만 등록 (대역폭 계산은 이후부터)
		if prevNetRecv == 0 && prevNetSent == 0 {
			network.NetRecv.Add(float64(currentRecv))
			network.NetSent.Add(float64(currentSent))
		} else {
			// 누적 값 증가분 계산
			if currentRecv >= prevNetRecv {
				deltaRecv := currentRecv - prevNetRecv
				network.NetRecv.Add(float64(deltaRecv))
				// 초당 수신 바이트 수 (간단히 interval로 나눔)
				network.NetworkRecvBps.Set(float64(deltaRecv) / float64(intervalSeconds))
			}
			if currentSent >= prevNetSent {
				deltaSent := currentSent - prevNetSent
				network.NetSent.Add(float64(deltaSent))
				// 초당 전송 바이트 수
				network.NetworkSentBps.Set(float64(deltaSent) / float64(intervalSeconds))
			}
		}
		prevNetRecv = currentRecv
		prevNetSent = currentSent
	}
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

func main() {
	if loadErr := godotenv.Load(".env"); loadErr != nil {
		log.Fatalf("Error loading .env file: %v", loadErr)
	}

	convInt, convErr := strconv.Atoi(os.Getenv("METRICS_INTERVAL"))
	if convErr != nil {
		log.Fatalf("Error converting METRICS_INTERVAL to integer: %v", convErr)
	}

	interval := time.Duration(convInt)

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

	// /metrics endpoint 노출
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Exporter listening on :9468/metrics")
	log.Printf("Metrics update interval: %d seconds", interval)
	log.Fatal(http.ListenAndServe(":9468", nil))
}
