// +build !windows

package libnetwork

import "github.com/docker/libnetwork/ipamapi"

// Stub implementations for DNS related functions

func (n *network) startResolver() ***REMOVED***
***REMOVED***

func defaultIpamForNetworkType(networkType string) string ***REMOVED***
	return ipamapi.DefaultIPAM
***REMOVED***
