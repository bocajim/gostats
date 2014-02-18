package gostats

import (
	//"github.com/bocajim/helpers/log"
	"io/ioutil"
	"strconv"
	"strings"
)

func interfaces(ifindex int) (map[string]Interface, error) {

	b, err := ioutil.ReadFile("/proc/net/dev")
	if err != nil {
		return nil,err
	}

	lines := strings.Split(string(b), "\n")

	ifm := make(map[string]Interface)
	for _, line := range lines {
	
		if strings.HasPrefix(line," face") {
			continue
		}
	
		ss := strings.Fields(line)
		
		if len(ss)<11 {
			continue
		}
		
		name := strings.TrimSpace(ss[0])
		name = strings.TrimRight(name,":")
		
		//R bytes
		bytesIn, _ := strconv.ParseInt(ss[1],10,64)
		//R pkts
		pktsIn, _ := strconv.ParseInt(ss[2],10,64)
		
		//T bytes
		bytesOut, _ := strconv.ParseInt(ss[9],10,64)
		//T pkts
		pktsOut, _ := strconv.ParseInt(ss[10],10,64)
		
		ifi := Interface{
				Name:         name,
				BytesIn:      bytesIn,
				BytesOut:     bytesOut,
				PacketsIn:    pktsIn,
				PacketsOut:   pktsOut,
			}
			ifm[ifi.Name] = ifi
		
	}
	
	return ifm,nil
}
