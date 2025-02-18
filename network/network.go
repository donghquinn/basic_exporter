package network

import (
	"log"
	"sync"

	psnet "github.com/shirou/gopsutil/v4/net"
)

// 이전 네트워크 누적값 저장 및 동시성 안전을 위한 뮤텍스
var (
	prevNetRecv uint64
	prevNetSent uint64
	mu          sync.Mutex
)

// UpdateNetwork는 intervalSeconds 간격으로 네트워크 IO 카운터를 수집
// 누적 카운터와 초당 대역폭(Gauge) 메트릭을 업데이트
// 기본적으로 모든 인터페이스의 합산값(netCounters[0])을 사용
// 필요에 따라 특정 인터페이스 필터링을 추가 가능
func UpdateNetwork(intervalSeconds int) {
	// 네트워크 정보 수집 (모든 인터페이스의 합산)
	netCounters, err := psnet.IOCounters(false)
	if err != nil || len(netCounters) == 0 {
		log.Printf("Failed to get network counters: %v", err)
		return
	}

	// netCounters[0]는 모든 인터페이스의 합산값
	currentRecv := netCounters[0].BytesRecv
	currentSent := netCounters[0].BytesSent

	mu.Lock()
	defer mu.Unlock()

	// 첫 실행 시, 누적 값만 업데이트 (대역폭 계산은 이후부터 가능)
	if prevNetRecv == 0 && prevNetSent == 0 {
		NetRecv.Add(float64(currentRecv))
		NetSent.Add(float64(currentSent))
	} else {
		// 누적 값 증가분 계산
		if currentRecv >= prevNetRecv {
			deltaRecv := currentRecv - prevNetRecv
			NetRecv.Add(float64(deltaRecv))
			// 간단히 측정 간격으로 나누어 초당 수신 바이트 수 계산
			NetworkRecvBps.Set(float64(deltaRecv) / float64(intervalSeconds))
		}
		if currentSent >= prevNetSent {
			deltaSent := currentSent - prevNetSent
			NetSent.Add(float64(deltaSent))
			NetworkSentBps.Set(float64(deltaSent) / float64(intervalSeconds))
		}
	}

	// 현재 값을 이전 값으로 업데이트
	prevNetRecv = currentRecv
	prevNetSent = currentSent
}
