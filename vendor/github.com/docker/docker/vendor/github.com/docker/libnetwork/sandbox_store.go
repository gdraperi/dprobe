package libnetwork

import (
	"container/heap"
	"encoding/json"
	"sync"

	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/osl"
	"github.com/sirupsen/logrus"
)

const (
	sandboxPrefix = "sandbox"
)

type epState struct ***REMOVED***
	Eid string
	Nid string
***REMOVED***

type sbState struct ***REMOVED***
	ID         string
	Cid        string
	c          *controller
	dbIndex    uint64
	dbExists   bool
	Eps        []epState
	EpPriority map[string]int
	// external servers have to be persisted so that on restart of a live-restore
	// enabled daemon we get the external servers for the running containers.
	// We have two versions of ExtDNS to support upgrade & downgrade of the daemon
	// between >=1.14 and <1.14 versions.
	ExtDNS  []string
	ExtDNS2 []extDNSEntry
***REMOVED***

func (sbs *sbState) Key() []string ***REMOVED***
	return []string***REMOVED***sandboxPrefix, sbs.ID***REMOVED***
***REMOVED***

func (sbs *sbState) KeyPrefix() []string ***REMOVED***
	return []string***REMOVED***sandboxPrefix***REMOVED***
***REMOVED***

func (sbs *sbState) Value() []byte ***REMOVED***
	b, err := json.Marshal(sbs)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	return b
***REMOVED***

func (sbs *sbState) SetValue(value []byte) error ***REMOVED***
	return json.Unmarshal(value, sbs)
***REMOVED***

func (sbs *sbState) Index() uint64 ***REMOVED***
	sbi, err := sbs.c.SandboxByID(sbs.ID)
	if err != nil ***REMOVED***
		return sbs.dbIndex
	***REMOVED***

	sb := sbi.(*sandbox)
	maxIndex := sb.dbIndex
	if sbs.dbIndex > maxIndex ***REMOVED***
		maxIndex = sbs.dbIndex
	***REMOVED***

	return maxIndex
***REMOVED***

func (sbs *sbState) SetIndex(index uint64) ***REMOVED***
	sbs.dbIndex = index
	sbs.dbExists = true

	sbi, err := sbs.c.SandboxByID(sbs.ID)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	sb := sbi.(*sandbox)
	sb.dbIndex = index
	sb.dbExists = true
***REMOVED***

func (sbs *sbState) Exists() bool ***REMOVED***
	if sbs.dbExists ***REMOVED***
		return sbs.dbExists
	***REMOVED***

	sbi, err := sbs.c.SandboxByID(sbs.ID)
	if err != nil ***REMOVED***
		return false
	***REMOVED***

	sb := sbi.(*sandbox)
	return sb.dbExists
***REMOVED***

func (sbs *sbState) Skip() bool ***REMOVED***
	return false
***REMOVED***

func (sbs *sbState) New() datastore.KVObject ***REMOVED***
	return &sbState***REMOVED***c: sbs.c***REMOVED***
***REMOVED***

func (sbs *sbState) CopyTo(o datastore.KVObject) error ***REMOVED***
	dstSbs := o.(*sbState)
	dstSbs.c = sbs.c
	dstSbs.ID = sbs.ID
	dstSbs.Cid = sbs.Cid
	dstSbs.dbIndex = sbs.dbIndex
	dstSbs.dbExists = sbs.dbExists
	dstSbs.EpPriority = sbs.EpPriority

	dstSbs.Eps = append(dstSbs.Eps, sbs.Eps...)

	if len(sbs.ExtDNS2) > 0 ***REMOVED***
		for _, dns := range sbs.ExtDNS2 ***REMOVED***
			dstSbs.ExtDNS2 = append(dstSbs.ExtDNS2, dns)
			dstSbs.ExtDNS = append(dstSbs.ExtDNS, dns.IPStr)
		***REMOVED***
		return nil
	***REMOVED***
	for _, dns := range sbs.ExtDNS ***REMOVED***
		dstSbs.ExtDNS = append(dstSbs.ExtDNS, dns)
		dstSbs.ExtDNS2 = append(dstSbs.ExtDNS2, extDNSEntry***REMOVED***IPStr: dns***REMOVED***)
	***REMOVED***

	return nil
***REMOVED***

func (sbs *sbState) DataScope() string ***REMOVED***
	return datastore.LocalScope
***REMOVED***

func (sb *sandbox) storeUpdate() error ***REMOVED***
	sbs := &sbState***REMOVED***
		c:          sb.controller,
		ID:         sb.id,
		Cid:        sb.containerID,
		EpPriority: sb.epPriority,
		ExtDNS2:    sb.extDNS,
	***REMOVED***

	for _, ext := range sb.extDNS ***REMOVED***
		sbs.ExtDNS = append(sbs.ExtDNS, ext.IPStr)
	***REMOVED***

retry:
	sbs.Eps = nil
	for _, ep := range sb.getConnectedEndpoints() ***REMOVED***
		// If the endpoint is not persisted then do not add it to
		// the sandbox checkpoint
		if ep.Skip() ***REMOVED***
			continue
		***REMOVED***

		eps := epState***REMOVED***
			Nid: ep.getNetwork().ID(),
			Eid: ep.ID(),
		***REMOVED***

		sbs.Eps = append(sbs.Eps, eps)
	***REMOVED***

	err := sb.controller.updateToStore(sbs)
	if err == datastore.ErrKeyModified ***REMOVED***
		// When we get ErrKeyModified it is sufficient to just
		// go back and retry.  No need to get the object from
		// the store because we always regenerate the store
		// state from in memory sandbox state
		goto retry
	***REMOVED***

	return err
***REMOVED***

func (sb *sandbox) storeDelete() error ***REMOVED***
	sbs := &sbState***REMOVED***
		c:        sb.controller,
		ID:       sb.id,
		Cid:      sb.containerID,
		dbIndex:  sb.dbIndex,
		dbExists: sb.dbExists,
	***REMOVED***

	return sb.controller.deleteFromStore(sbs)
***REMOVED***

func (c *controller) sandboxCleanup(activeSandboxes map[string]interface***REMOVED******REMOVED***) ***REMOVED***
	store := c.getStore(datastore.LocalScope)
	if store == nil ***REMOVED***
		logrus.Error("Could not find local scope store while trying to cleanup sandboxes")
		return
	***REMOVED***

	kvol, err := store.List(datastore.Key(sandboxPrefix), &sbState***REMOVED***c: c***REMOVED***)
	if err != nil && err != datastore.ErrKeyNotFound ***REMOVED***
		logrus.Errorf("failed to get sandboxes for scope %s: %v", store.Scope(), err)
		return
	***REMOVED***

	// It's normal for no sandboxes to be found. Just bail out.
	if err == datastore.ErrKeyNotFound ***REMOVED***
		return
	***REMOVED***

	for _, kvo := range kvol ***REMOVED***
		sbs := kvo.(*sbState)

		sb := &sandbox***REMOVED***
			id:                 sbs.ID,
			controller:         sbs.c,
			containerID:        sbs.Cid,
			endpoints:          epHeap***REMOVED******REMOVED***,
			populatedEndpoints: map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***,
			dbIndex:            sbs.dbIndex,
			isStub:             true,
			dbExists:           true,
		***REMOVED***
		// If we are restoring from a older version extDNSEntry won't have the
		// HostLoopback field
		if len(sbs.ExtDNS2) > 0 ***REMOVED***
			sb.extDNS = sbs.ExtDNS2
		***REMOVED*** else ***REMOVED***
			for _, dns := range sbs.ExtDNS ***REMOVED***
				sb.extDNS = append(sb.extDNS, extDNSEntry***REMOVED***IPStr: dns***REMOVED***)
			***REMOVED***
		***REMOVED***

		msg := " for cleanup"
		create := true
		isRestore := false
		if val, ok := activeSandboxes[sb.ID()]; ok ***REMOVED***
			msg = ""
			sb.isStub = false
			isRestore = true
			opts := val.([]SandboxOption)
			sb.processOptions(opts...)
			sb.restorePath()
			create = !sb.config.useDefaultSandBox
			heap.Init(&sb.endpoints)
		***REMOVED***
		sb.osSbox, err = osl.NewSandbox(sb.Key(), create, isRestore)
		if err != nil ***REMOVED***
			logrus.Errorf("failed to create osl sandbox while trying to restore sandbox %s%s: %v", sb.ID()[0:7], msg, err)
			continue
		***REMOVED***

		c.Lock()
		c.sandboxes[sb.id] = sb
		c.Unlock()

		for _, eps := range sbs.Eps ***REMOVED***
			n, err := c.getNetworkFromStore(eps.Nid)
			var ep *endpoint
			if err != nil ***REMOVED***
				logrus.Errorf("getNetworkFromStore for nid %s failed while trying to build sandbox for cleanup: %v", eps.Nid, err)
				n = &network***REMOVED***id: eps.Nid, ctrlr: c, drvOnce: &sync.Once***REMOVED******REMOVED***, persist: true***REMOVED***
				ep = &endpoint***REMOVED***id: eps.Eid, network: n, sandboxID: sbs.ID***REMOVED***
			***REMOVED*** else ***REMOVED***
				ep, err = n.getEndpointFromStore(eps.Eid)
				if err != nil ***REMOVED***
					logrus.Errorf("getEndpointFromStore for eid %s failed while trying to build sandbox for cleanup: %v", eps.Eid, err)
					ep = &endpoint***REMOVED***id: eps.Eid, network: n, sandboxID: sbs.ID***REMOVED***
				***REMOVED***
			***REMOVED***
			if _, ok := activeSandboxes[sb.ID()]; ok && err != nil ***REMOVED***
				logrus.Errorf("failed to restore endpoint %s in %s for container %s due to %v", eps.Eid, eps.Nid, sb.ContainerID(), err)
				continue
			***REMOVED***
			heap.Push(&sb.endpoints, ep)
		***REMOVED***

		if _, ok := activeSandboxes[sb.ID()]; !ok ***REMOVED***
			logrus.Infof("Removing stale sandbox %s (%s)", sb.id, sb.containerID)
			if err := sb.delete(true); err != nil ***REMOVED***
				logrus.Errorf("Failed to delete sandbox %s while trying to cleanup: %v", sb.id, err)
			***REMOVED***
			continue
		***REMOVED***

		// reconstruct osl sandbox field
		if !sb.config.useDefaultSandBox ***REMOVED***
			if err := sb.restoreOslSandbox(); err != nil ***REMOVED***
				logrus.Errorf("failed to populate fields for osl sandbox %s", sb.ID())
				continue
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			c.sboxOnce.Do(func() ***REMOVED***
				c.defOsSbox = sb.osSbox
			***REMOVED***)
		***REMOVED***

		for _, ep := range sb.endpoints ***REMOVED***
			// Watch for service records
			if !c.isAgent() ***REMOVED***
				c.watchSvcRecord(ep)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
