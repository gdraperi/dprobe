package hostdiscovery

import (
	"net"
	"sync"

	"github.com/sirupsen/logrus"

	mapset "github.com/deckarep/golang-set"
	"github.com/docker/docker/pkg/discovery"
	// Including KV
	_ "github.com/docker/docker/pkg/discovery/kv"
	"github.com/docker/libkv/store/consul"
	"github.com/docker/libkv/store/etcd"
	"github.com/docker/libkv/store/zookeeper"
	"github.com/docker/libnetwork/types"
)

type hostDiscovery struct ***REMOVED***
	watcher  discovery.Watcher
	nodes    mapset.Set
	stopChan chan struct***REMOVED******REMOVED***
	sync.Mutex
***REMOVED***

func init() ***REMOVED***
	consul.Register()
	etcd.Register()
	zookeeper.Register()
***REMOVED***

// NewHostDiscovery function creates a host discovery object
func NewHostDiscovery(watcher discovery.Watcher) HostDiscovery ***REMOVED***
	return &hostDiscovery***REMOVED***watcher: watcher, nodes: mapset.NewSet(), stopChan: make(chan struct***REMOVED******REMOVED***)***REMOVED***
***REMOVED***

func (h *hostDiscovery) Watch(activeCallback ActiveCallback, joinCallback JoinCallback, leaveCallback LeaveCallback) error ***REMOVED***
	h.Lock()
	d := h.watcher
	h.Unlock()
	if d == nil ***REMOVED***
		return types.BadRequestErrorf("invalid discovery watcher")
	***REMOVED***
	discoveryCh, errCh := d.Watch(h.stopChan)
	go h.monitorDiscovery(discoveryCh, errCh, activeCallback, joinCallback, leaveCallback)
	return nil
***REMOVED***

func (h *hostDiscovery) monitorDiscovery(ch <-chan discovery.Entries, errCh <-chan error,
	activeCallback ActiveCallback, joinCallback JoinCallback, leaveCallback LeaveCallback) ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case entries := <-ch:
			h.processCallback(entries, activeCallback, joinCallback, leaveCallback)
		case err := <-errCh:
			if err != nil ***REMOVED***
				logrus.Errorf("discovery error: %v", err)
			***REMOVED***
		case <-h.stopChan:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (h *hostDiscovery) StopDiscovery() error ***REMOVED***
	h.Lock()
	stopChan := h.stopChan
	h.watcher = nil
	h.Unlock()

	close(stopChan)
	return nil
***REMOVED***

func (h *hostDiscovery) processCallback(entries discovery.Entries,
	activeCallback ActiveCallback, joinCallback JoinCallback, leaveCallback LeaveCallback) ***REMOVED***
	updated := hosts(entries)
	h.Lock()
	existing := h.nodes
	added, removed := diff(existing, updated)
	h.nodes = updated
	h.Unlock()

	activeCallback()
	if len(added) > 0 ***REMOVED***
		joinCallback(added)
	***REMOVED***
	if len(removed) > 0 ***REMOVED***
		leaveCallback(removed)
	***REMOVED***
***REMOVED***

func diff(existing mapset.Set, updated mapset.Set) (added []net.IP, removed []net.IP) ***REMOVED***
	addSlice := updated.Difference(existing).ToSlice()
	removeSlice := existing.Difference(updated).ToSlice()
	for _, ip := range addSlice ***REMOVED***
		added = append(added, net.ParseIP(ip.(string)))
	***REMOVED***
	for _, ip := range removeSlice ***REMOVED***
		removed = append(removed, net.ParseIP(ip.(string)))
	***REMOVED***
	return
***REMOVED***

func (h *hostDiscovery) Fetch() []net.IP ***REMOVED***
	h.Lock()
	defer h.Unlock()
	ips := []net.IP***REMOVED******REMOVED***
	for _, ipstr := range h.nodes.ToSlice() ***REMOVED***
		ips = append(ips, net.ParseIP(ipstr.(string)))
	***REMOVED***
	return ips
***REMOVED***

func hosts(entries discovery.Entries) mapset.Set ***REMOVED***
	hosts := mapset.NewSet()
	for _, entry := range entries ***REMOVED***
		hosts.Add(entry.Host)
	***REMOVED***
	return hosts
***REMOVED***
