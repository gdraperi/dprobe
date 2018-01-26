// Package ipamutils provides utility functions for ipam management
package ipamutils

import (
	"net"
	"sync"
)

var (
	// PredefinedBroadNetworks contains a list of 31 IPv4 private networks with host size 16 and 12
	// (172.17-31.x.x/16, 192.168.x.x/20) which do not overlap with the networks in `PredefinedGranularNetworks`
	PredefinedBroadNetworks []*net.IPNet
	// PredefinedGranularNetworks contains a list of 64K IPv4 private networks with host size 8
	// (10.x.x.x/24) which do not overlap with the networks in `PredefinedBroadNetworks`
	PredefinedGranularNetworks []*net.IPNet

	initNetworksOnce sync.Once
)

// InitNetworks initializes the pre-defined networks used by the built-in IP allocator
func InitNetworks() ***REMOVED***
	initNetworksOnce.Do(func() ***REMOVED***
		PredefinedBroadNetworks = initBroadPredefinedNetworks()
		PredefinedGranularNetworks = initGranularPredefinedNetworks()
	***REMOVED***)
***REMOVED***

func initBroadPredefinedNetworks() []*net.IPNet ***REMOVED***
	pl := make([]*net.IPNet, 0, 31)
	mask := []byte***REMOVED***255, 255, 0, 0***REMOVED***
	for i := 17; i < 32; i++ ***REMOVED***
		pl = append(pl, &net.IPNet***REMOVED***IP: []byte***REMOVED***172, byte(i), 0, 0***REMOVED***, Mask: mask***REMOVED***)
	***REMOVED***
	mask20 := []byte***REMOVED***255, 255, 240, 0***REMOVED***
	for i := 0; i < 16; i++ ***REMOVED***
		pl = append(pl, &net.IPNet***REMOVED***IP: []byte***REMOVED***192, 168, byte(i << 4), 0***REMOVED***, Mask: mask20***REMOVED***)
	***REMOVED***
	return pl
***REMOVED***

func initGranularPredefinedNetworks() []*net.IPNet ***REMOVED***
	pl := make([]*net.IPNet, 0, 256*256)
	mask := []byte***REMOVED***255, 255, 255, 0***REMOVED***
	for i := 0; i < 256; i++ ***REMOVED***
		for j := 0; j < 256; j++ ***REMOVED***
			pl = append(pl, &net.IPNet***REMOVED***IP: []byte***REMOVED***10, byte(i), byte(j), 0***REMOVED***, Mask: mask***REMOVED***)
		***REMOVED***
	***REMOVED***
	return pl
***REMOVED***
