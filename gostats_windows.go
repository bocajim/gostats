package gostats

import (
	//"github.com/bocajim/helpers/log"
	"syscall"
	"unsafe"
)

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")
	//modpsapi                   = syscall.NewLazyDLL("psapi.dll")
	procGetSystemTimes       = modkernel32.NewProc("GetSystemTimes")
	procGlobalMemoryStatusEx = modkernel32.NewProc("GlobalMemoryStatusEx")
)

var lastIdleTime int64
var lastKernelTime int64
var lastUserTime int64

func Cpu() int {
	var idleTime int64
	var kernelTime int64
	var userTime int64
	syscall.Syscall(procGetSystemTimes.Addr(), 3, uintptr(unsafe.Pointer(&idleTime)), uintptr(unsafe.Pointer(&kernelTime)), uintptr(unsafe.Pointer(&userTime)))

	if lastIdleTime == 0 {
		lastIdleTime = idleTime
		lastKernelTime = kernelTime
		lastUserTime = userTime
		return 0
	} else {
		deltaIdleTime := idleTime - lastIdleTime
		deltaKernelTime := kernelTime - lastKernelTime
		deltaUserTime := userTime - lastUserTime

		lastIdleTime = idleTime
		lastKernelTime = kernelTime
		lastUserTime = userTime
		
		total := deltaKernelTime + deltaUserTime
		if total == 0 {
			return 0
		}

		return int(((deltaKernelTime + deltaUserTime - deltaIdleTime) * 100) / total)
	}
}

type MEMORYSTATUSEX struct {
	dwLength                uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

func MemoryPhysical() (int, uint64, uint64) {

	var mex MEMORYSTATUSEX
	mex.dwLength = 64

	syscall.Syscall(procGlobalMemoryStatusEx.Addr(), 1, uintptr(unsafe.Pointer(&mex)), 0, 0)

	return int((mex.ullTotalPhys - mex.ullAvailPhys) * 100 / mex.ullTotalPhys), (mex.ullTotalPhys - mex.ullAvailPhys), mex.ullTotalPhys
}

func MemoryVirtual() (int, uint64, uint64) {

	var mex MEMORYSTATUSEX
	mex.dwLength = 64

	syscall.Syscall(procGlobalMemoryStatusEx.Addr(), 1, uintptr(unsafe.Pointer(&mex)), 0, 0)

	return int((mex.ullTotalPageFile - mex.ullAvailPageFile) * 100 / mex.ullTotalPageFile), (mex.ullTotalPageFile - mex.ullAvailPageFile), mex.ullTotalPageFile
}
