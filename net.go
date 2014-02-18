package gostats

import (
	"math"
	"time"
)

type Interface struct {
	Name       string // e.g., "en0", "lo0", "eth0.100"
	BytesIn    int64
	BytesOut   int64
	PacketsIn  int64
	PacketsOut int64
}

type Usage struct {
	MegaBitsPerSecondIn  float64 `json:"mbpsIn"   bson:"mbpsIn"`
	MegaBitsPerSecondOut float64 `json:"mbpsOut"  bson:"mbpsOut"`
	PacketsPerSecondIn   float64 `json:"ppsIn"    bson:"ppsIn"`
	PacketsPerSecondOut  float64 `json:"ppsOut"   bson:"ppsOut"`
}

var prevTime time.Time
var prevInterfaces map[string]Interface
var currentUsage map[string]Usage

func networkDiff() {

	rm := make(map[string]Usage)

	newt := time.Now()
	difft := newt.Sub(prevTime)

	ifm, err := interfaces(0)
	if err != nil {
		return
	}

	if prevInterfaces == nil {
		prevInterfaces = ifm
		prevTime = newt
	}

	for k, v := range ifm {
		p, f := prevInterfaces[k]
		if !f {
			continue
		}
		var u Usage
		if p.BytesIn > v.BytesIn {
			u.MegaBitsPerSecondIn = float64(math.MaxUint32-p.BytesIn+v.BytesIn) / difft.Seconds() / 131072.0
		} else {
			u.MegaBitsPerSecondIn = float64(v.BytesIn-p.BytesIn) / difft.Seconds() / 131072.0
		}
		if p.BytesOut > v.BytesOut {
			u.MegaBitsPerSecondOut = float64(math.MaxUint32-p.BytesOut+v.BytesOut) / difft.Seconds() / 131072.0
		} else {
			u.MegaBitsPerSecondOut = float64(v.BytesOut-p.BytesOut) / difft.Seconds() / 131072.0
		}
		if p.PacketsIn > v.PacketsIn {
			u.PacketsPerSecondIn = float64(math.MaxUint32-p.PacketsIn+v.PacketsIn) / difft.Seconds()
		} else {
			u.PacketsPerSecondIn = float64(v.PacketsIn-p.PacketsIn) / difft.Seconds()
		}
		if p.PacketsOut > v.PacketsOut {
			u.PacketsPerSecondOut = float64(math.MaxUint32-p.PacketsOut+v.PacketsOut) / difft.Seconds()
		} else {
			u.PacketsPerSecondOut = float64(v.PacketsOut-p.PacketsOut) / difft.Seconds()
		}
		
		rm[k] = u
	}
	currentUsage = rm
	prevTime = newt
	prevInterfaces = ifm
}

func NetworkInterfaceUsage() map[string]Usage {
	return currentUsage
}
