package gostats

import (
	"container/ring"
	"time"
)

var cpuRing *ring.Ring
var cpuAverage int

func CalcAverages() {

	cpuRing = ring.New(12)
	cpuAverage = Cpu()
	for i := 0; i < 12; i++ {
		cpuRing.Value = cpuAverage
		cpuRing = cpuRing.Next()
	}

	go func() {
		for {
			select {
			case <-time.After(time.Second * 5):
				//calc CPU
				cpuRing.Value = Cpu()
				cpuRing = cpuRing.Next()
				tempAvgSum = 0
				cpuRing.Do(ringAvg)
				cpuAverage = tempAvgSum / 12
			}
		}
	}()
}

var tempAvgSum int

func ringAvg(i interface{}) {
	tempAvgSum = tempAvgSum + i.(int)
}

func CpuAverage() int {
	return cpuAverage
}

func MemoryPhysicalUsage() int {
	val, _, _ := MemoryPhysical()
	return val
}

func MemoryVirtualUsage() int {
	val, _, _ := MemoryVirtual()
	return val
}
