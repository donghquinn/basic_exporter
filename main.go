package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	psnet "github.com/shirou/gopsutil/v4/net"
)

// 누적 네트워크 바이트 메트릭 (카운터)
var (
	netRecv = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "network_receive_bytes_total",
		Help: "Total number of bytes received on all network interfaces",
	})
	netSent = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "network_transmit_bytes_total",
		Help: "Total number of bytes sent on all network interfaces",
	})
)

// 이전 값들을 저장해서 증분값을 계산할 경우 필요하면 전역변수 사용
var prevRecv, prevSent uint64

func updateMetrics() {
	// IOCounters(false): 모든 인터페이스를 합산한 결과 반환
	counters, err := psnet.IOCounters(false)
	if err != nil || len(counters) == 0 {
		log.Printf("Failed to get network counters: %v", err)
		return
	}

	currentRecv := counters[0].BytesRecv
	currentSent := counters[0].BytesSent

	// Prometheus의 Counter는 직접 감소시키거나 설정할 수 없으므로,
	// 최초 값만 세팅하는 경우엔 아래처럼 초기값을 기록하고,
	// 이후에는 누적 증가분만 Add() 해주면 됩니다.
	// (혹은, 메트릭을 Gauge로 등록하고 직접 설정할 수도 있습니다.)
	if prevRecv == 0 && prevSent == 0 {
		// 초기값 설정 (이미 이전 누적 값이 있는 경우 Prometheus에서는 rate()로 계산 가능)
		netRecv.Add(float64(currentRecv))
		netSent.Add(float64(currentSent))
	} else {
		// 누적 값의 증가량을 더함
		if currentRecv >= prevRecv {
			netRecv.Add(float64(currentRecv - prevRecv))
		}
		if currentSent >= prevSent {
			netSent.Add(float64(currentSent - prevSent))
		}
	}
	prevRecv = currentRecv
	prevSent = currentSent
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
	// Prometheus에 메트릭 등록
	prometheus.MustRegister(netRecv)
	prometheus.MustRegister(netSent)

	// 메트릭 업데이트 시작
	recordMetrics()

	// /metrics endpoint 노출
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Exporter listening on :8080/metrics")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
