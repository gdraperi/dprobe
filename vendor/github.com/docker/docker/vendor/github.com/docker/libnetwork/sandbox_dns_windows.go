// +build windows

package libnetwork

import (
	"github.com/docker/libnetwork/etchosts"
)

// Stub implementations for DNS related functions

func (sb *sandbox) startResolver(bool) ***REMOVED***
***REMOVED***

func (sb *sandbox) setupResolutionFiles() error ***REMOVED***
	return nil
***REMOVED***

func (sb *sandbox) restorePath() ***REMOVED***
***REMOVED***

func (sb *sandbox) updateHostsFile(ifaceIP string) error ***REMOVED***
	return nil
***REMOVED***

func (sb *sandbox) addHostsEntries(recs []etchosts.Record) ***REMOVED***

***REMOVED***

func (sb *sandbox) deleteHostsEntries(recs []etchosts.Record) ***REMOVED***

***REMOVED***

func (sb *sandbox) updateDNS(ipv6Enabled bool) error ***REMOVED***
	return nil
***REMOVED***
