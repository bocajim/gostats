package gostats

import (
	//"github.com/bocajim/helpers/log"
	"io/ioutil"
	"strconv"
	"strings"
)

var lastUserTime int64
var lastNiceTime int64
var lastSystemTime int64
var lastIdleTime int64

func Cpu() int {
	b, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return -1
	}

	lines := strings.Split(string(b), "\n")

	for _, line := range lines {
		ss := strings.Split(line, " ")

		if ss[0] == "cpu" && len(ss) >= 5 {
			//user
			user, _ := strconv.ParseInt(ss[2], 10, 64)

			//nice
			nice, _ := strconv.ParseInt(ss[3], 10, 64)

			//system
			system, _ := strconv.ParseInt(ss[4], 10, 64)

			//idle
			idle, _ := strconv.ParseInt(ss[5], 10, 64)

			if lastIdleTime == 0 {
				lastUserTime = user
				lastNiceTime = nice
				lastSystemTime = system
				lastIdleTime = idle
				return 0
			} else {
				deltaUserTime := user - lastUserTime
				deltaNiceTime := nice - lastNiceTime
				deltaSystemTime := system - lastSystemTime
				deltaIdleTime := idle - lastIdleTime

				lastUserTime = user
				lastNiceTime = nice
				lastSystemTime = system
				lastIdleTime = idle

				return int(((deltaUserTime + deltaNiceTime + deltaSystemTime) * 100) / (deltaUserTime + deltaNiceTime + deltaSystemTime + deltaIdleTime))
			}

		}
	}

	return 0
}

func getValue(line string) uint64 {
	ss := strings.Split(line, " ")
	for _, s := range ss {
		if v, err := strconv.ParseUint(s, 10, 64); err == nil {
			return v
		}
	}
	return 0
}

func MemoryPhysical() (int, uint64, uint64) {
	b, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return -1, 0, 0
	}

	lines := strings.Split(string(b), "\n")

	totalKb := uint64(0)
	freeKb := uint64(0)
	cachedKb := uint64(0)

	for _, line := range lines {

		if strings.HasPrefix(line, "MemTotal") {
			totalKb = getValue(line)
		} else if strings.HasPrefix(line, "MemFree") {
			freeKb = getValue(line)
		} else if strings.HasPrefix(line, "Cached") {
			cachedKb = getValue(line)
		}
	}

	if totalKb < 0 || freeKb < 0 || cachedKb < 0 {
		return 0, 0, 0
	}

	usage := int((totalKb - (freeKb + cachedKb)) * 100 / (totalKb))
	used := (totalKb - (freeKb + cachedKb)) * 1024
	total := (totalKb) * 1024

	return usage, used, total
}

func MemoryVirtual() (int, uint64, uint64) {
	b, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return -1, 0, 0
	}

	lines := strings.Split(string(b), "\n")

	totalKb := uint64(0)
	freeKb := uint64(0)
	cachedKb := uint64(0)
	swapTotalKb := uint64(0)
	swapFreeKb := uint64(0)

	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal") {
			totalKb = getValue(line)
		} else if strings.HasPrefix(line, "MemFree") {
			freeKb = getValue(line)
		} else if strings.HasPrefix(line, "Cached") {
			cachedKb = getValue(line)
		} else if strings.HasPrefix(line, "SwapTotal") {
			swapTotalKb = getValue(line)
		} else if strings.HasPrefix(line, "SwapFree") {
			swapFreeKb = getValue(line)
		}
	}

	if totalKb < 0 || freeKb < 0 || cachedKb < 0 || swapTotalKb < 0 || swapFreeKb < 0 {
		return 0, 0, 0
	}

	usage := int((totalKb + swapTotalKb - (freeKb + cachedKb + swapFreeKb)) * 100 / (totalKb + swapTotalKb))
	used := (totalKb + swapTotalKb - (freeKb + cachedKb + swapFreeKb)) * 1024
	total := (totalKb + swapTotalKb) * 1024

	return usage, used, total
}
