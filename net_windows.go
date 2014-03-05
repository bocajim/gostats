package gostats

import (
	//"github.com/bocajim/helpers/log"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

func sysSocket(f, t, p int) (syscall.Handle, error) {
	// See ../syscall/exec_unix.go for description of ForkLock.
	syscall.ForkLock.RLock()
	s, err := syscall.Socket(f, t, p)
	if err == nil {
		syscall.CloseOnExec(s)
	}
	syscall.ForkLock.RUnlock()
	return s, err
}

func bytePtrToString(p *uint8) string {
	a := (*[10000]uint8)(unsafe.Pointer(p))
	i := 0
	for a[i] != 0 {
		i++
	}
	return string(a[:i])
}

func getAdapterList() (*syscall.IpAdapterInfo, error) {
	b := make([]byte, 1000)
	l := uint32(len(b))
	a := (*syscall.IpAdapterInfo)(unsafe.Pointer(&b[0]))
	// TODO(mikio): GetAdaptersInfo returns IP_ADAPTER_INFO that
	// contains IPv4 address list only. We should use another API
	// for fetching IPv6 stuff from the kernel.
	err := syscall.GetAdaptersInfo(a, &l)
	if err == syscall.ERROR_BUFFER_OVERFLOW {
		b = make([]byte, l)
		a = (*syscall.IpAdapterInfo)(unsafe.Pointer(&b[0]))
		err = syscall.GetAdaptersInfo(a, &l)
	}
	if err != nil {
		return nil, os.NewSyscallError("GetAdaptersInfo", err)
	}
	return a, nil
}

func getInterfaceList() ([]syscall.InterfaceInfo, error) {
	s, err := sysSocket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		return nil, os.NewSyscallError("Socket", err)
	}
	defer syscall.Closesocket(s)

	ii := [20]syscall.InterfaceInfo{}
	ret := uint32(0)
	size := uint32(unsafe.Sizeof(ii))
	err = syscall.WSAIoctl(s, syscall.SIO_GET_INTERFACE_LIST, nil, 0, (*byte)(unsafe.Pointer(&ii[0])), size, &ret, nil, 0)
	if err != nil {
		return nil, os.NewSyscallError("WSAIoctl", err)
	}
	c := ret / uint32(unsafe.Sizeof(ii[0]))
	return ii[:c-1], nil
}

// If the ifindex is zero, interfaceTable returns mappings of all
// network interfaces.  Otherwise it returns a mapping of a specific
// interface.
func interfaces(ifindex int) (map[string]Interface, error) {
	ai, err := getAdapterList()
	if err != nil {
		return nil, err
	}

	ifm := make(map[string]Interface)
	for ; ai != nil; ai = ai.Next {
		index := ai.Index
		if ifindex == 0 || ifindex == int(index) {

			row := syscall.MibIfRow{Index: index}
			e := syscall.GetIfEntry(&row)
			if e != nil {
				return nil, os.NewSyscallError("GetIfEntry", e)
			}

			name := bytePtrToString(&ai.Description[0])
			name = strings.Replace(name,"(R)","",-1)
			name = strings.Replace(name,".","_",-1)
			name = strings.Replace(name,"  "," ",-1)
			name = strings.Replace(name," ","_",-1)

			ifi := Interface{
				Name:         name,
				BytesIn:      int64(row.InOctets),
				BytesOut:     int64(row.OutOctets),
				PacketsIn:    int64(row.InUcastPkts),
				PacketsOut:   int64(row.OutUcastPkts),
			}
			ifm[ifi.Name] = ifi
		}
	}
	return ifm, nil
}
