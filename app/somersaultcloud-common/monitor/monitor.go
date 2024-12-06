package monitor

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func GetSystemMetrics() (availableMem uint64, cpuIdleTime float64) {
	memStat, _ := mem.VirtualMemory()
	availableMem = memStat.Available
	cTimes, _ := cpu.Times(false)
	cpuIdleTime = cTimes[0].Idle
	return
}
