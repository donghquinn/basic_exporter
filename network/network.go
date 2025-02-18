package network

import (
	"log"

	psnet "github.com/shirou/gopsutil/v4/net"
)

// 네트워크 증분 계산을 위한 이전 값 저장
var prevNetRecv, prevNetSent uint64

func UpdateNetwork(intervalSeconds int) {
	// 네트워크 정보 수집 (모든 인터페이스의 합산)
	netCounters, err := psnet.IOCounters(false)
	if err != nil || len(netCounters) == 0 {
		log.Printf("Failed to get network counters: %v", err)
	} else {
		currentRecv := netCounters[0].BytesRecv
		currentSent := netCounters[0].BytesSent

		// 첫 실행 시, 누적 값만 등록 (대역폭 계산은 이후부터)
		if prevNetRecv == 0 && prevNetSent == 0 {
			NetRecv.Add(float64(currentRecv))
			NetSent.Add(float64(currentSent))
		} else {
			// 누적 값 증가분 계산
			if currentRecv >= prevNetRecv {
				deltaRecv := currentRecv - prevNetRecv
				NetRecv.Add(float64(deltaRecv))
				// 초당 수신 바이트 수 (간단히 interval로 나눔)
				NetworkRecvBps.Set(float64(deltaRecv) / float64(intervalSeconds))
			}
			if currentSent >= prevNetSent {
				deltaSent := currentSent - prevNetSent
				NetSent.Add(float64(deltaSent))
				// 초당 전송 바이트 수
				NetworkSentBps.Set(float64(deltaSent) / float64(intervalSeconds))
			}
		}
		prevNetRecv = currentRecv
		prevNetSent = currentSent
	}
}
