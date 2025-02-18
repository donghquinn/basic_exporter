package system

import (
	"log"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
)

func UpdateSystemMetrics() {
	UpdateMemory()
	UpdateCpu()
	UpdateLoad()
}

// 메모리 정보 수집
func UpdateMemory() {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Failed to get memory info: %v", err)
	} else {
		MemoryTotal.Set(float64(vmStat.Total))
		MemoryUsed.Set(float64(vmStat.Used))
		MemoryAvailable.Set(float64(vmStat.Available))
		MemoryUsagePercent.Set(vmStat.UsedPercent)
	}

}

// CPU 사용률 (0초 간격은 즉시 반환하므로, 참고용으로 사용)
func UpdateCpu() {
	cpuPercents, err := cpu.Percent(0, false)
	if err != nil {
		log.Printf("Failed to get CPU percent: %v", err)
	} else if len(cpuPercents) > 0 {
		CpuUsage.Set(cpuPercents[0])
	}
}

// Load Average 수집
func UpdateLoad() {
	loadStat, err := load.Avg()
	if err != nil {
		log.Printf("Failed to get load average: %v", err)
	} else {
		Load1.Set(loadStat.Load1)
		Load5.Set(loadStat.Load5)
		Load15.Set(loadStat.Load15)
	}
}
