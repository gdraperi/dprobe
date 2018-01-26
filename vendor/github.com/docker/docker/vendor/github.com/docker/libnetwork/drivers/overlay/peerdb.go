package overlay

import (
	"context"
	"fmt"
	"net"
	"sync"
	"syscall"

	"github.com/docker/libnetwork/common"
	"github.com/docker/libnetwork/osl"
	"github.com/sirupsen/logrus"
)

const ovPeerTable = "overlay_peer_table"

type peerKey struct ***REMOVED***
	peerIP  net.IP
	peerMac net.HardwareAddr
***REMOVED***

type peerEntry struct ***REMOVED***
	eid        string
	vtep       net.IP
	peerIPMask net.IPMask
	isLocal    bool
***REMOVED***

func (p *peerEntry) MarshalDB() peerEntryDB ***REMOVED***
	ones, bits := p.peerIPMask.Size()
	return peerEntryDB***REMOVED***
		eid:            p.eid,
		vtep:           p.vtep.String(),
		peerIPMaskOnes: ones,
		peerIPMaskBits: bits,
		isLocal:        p.isLocal,
	***REMOVED***
***REMOVED***

// This the structure saved into the set (SetMatrix), due to the implementation of it
// the value inserted in the set has to be Hashable so the []byte had to be converted into
// strings
type peerEntryDB struct ***REMOVED***
	eid            string
	vtep           string
	peerIPMaskOnes int
	peerIPMaskBits int
	isLocal        bool
***REMOVED***

func (p *peerEntryDB) UnMarshalDB() peerEntry ***REMOVED***
	return peerEntry***REMOVED***
		eid:        p.eid,
		vtep:       net.ParseIP(p.vtep),
		peerIPMask: net.CIDRMask(p.peerIPMaskOnes, p.peerIPMaskBits),
		isLocal:    p.isLocal,
	***REMOVED***
***REMOVED***

type peerMap struct ***REMOVED***
	// set of peerEntry, note they have to be objects and not pointers to maintain the proper equality checks
	mp common.SetMatrix
	sync.Mutex
***REMOVED***

type peerNetworkMap struct ***REMOVED***
	// map with key peerKey
	mp map[string]*peerMap
	sync.Mutex
***REMOVED***

func (pKey peerKey) String() string ***REMOVED***
	return fmt.Sprintf("%s %s", pKey.peerIP, pKey.peerMac)
***REMOVED***

func (pKey *peerKey) Scan(state fmt.ScanState, verb rune) error ***REMOVED***
	ipB, err := state.Token(true, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	pKey.peerIP = net.ParseIP(string(ipB))

	macB, err := state.Token(true, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	pKey.peerMac, err = net.ParseMAC(string(macB))
	return err
***REMOVED***

func (d *driver) peerDbWalk(f func(string, *peerKey, *peerEntry) bool) error ***REMOVED***
	d.peerDb.Lock()
	nids := []string***REMOVED******REMOVED***
	for nid := range d.peerDb.mp ***REMOVED***
		nids = append(nids, nid)
	***REMOVED***
	d.peerDb.Unlock()

	for _, nid := range nids ***REMOVED***
		d.peerDbNetworkWalk(nid, func(pKey *peerKey, pEntry *peerEntry) bool ***REMOVED***
			return f(nid, pKey, pEntry)
		***REMOVED***)
	***REMOVED***
	return nil
***REMOVED***

func (d *driver) peerDbNetworkWalk(nid string, f func(*peerKey, *peerEntry) bool) error ***REMOVED***
	d.peerDb.Lock()
	pMap, ok := d.peerDb.mp[nid]
	d.peerDb.Unlock()

	if !ok ***REMOVED***
		return nil
	***REMOVED***

	mp := map[string]peerEntry***REMOVED******REMOVED***
	pMap.Lock()
	for _, pKeyStr := range pMap.mp.Keys() ***REMOVED***
		entryDBList, ok := pMap.mp.Get(pKeyStr)
		if ok ***REMOVED***
			peerEntryDB := entryDBList[0].(peerEntryDB)
			mp[pKeyStr] = peerEntryDB.UnMarshalDB()
		***REMOVED***
	***REMOVED***
	pMap.Unlock()

	for pKeyStr, pEntry := range mp ***REMOVED***
		var pKey peerKey
		if _, err := fmt.Sscan(pKeyStr, &pKey); err != nil ***REMOVED***
			logrus.Warnf("Peer key scan on network %s failed: %v", nid, err)
		***REMOVED***
		if f(&pKey, &pEntry) ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) peerDbSearch(nid string, peerIP net.IP) (*peerKey, *peerEntry, error) ***REMOVED***
	var pKeyMatched *peerKey
	var pEntryMatched *peerEntry
	err := d.peerDbNetworkWalk(nid, func(pKey *peerKey, pEntry *peerEntry) bool ***REMOVED***
		if pKey.peerIP.Equal(peerIP) ***REMOVED***
			pKeyMatched = pKey
			pEntryMatched = pEntry
			return true
		***REMOVED***

		return false
	***REMOVED***)

	if err != nil ***REMOVED***
		return nil, nil, fmt.Errorf("peerdb search for peer ip %q failed: %v", peerIP, err)
	***REMOVED***

	if pKeyMatched == nil || pEntryMatched == nil ***REMOVED***
		return nil, nil, fmt.Errorf("peer ip %q not found in peerdb", peerIP)
	***REMOVED***

	return pKeyMatched, pEntryMatched, nil
***REMOVED***

func (d *driver) peerDbAdd(nid, eid string, peerIP net.IP, peerIPMask net.IPMask,
	peerMac net.HardwareAddr, vtep net.IP, isLocal bool) (bool, int) ***REMOVED***

	d.peerDb.Lock()
	pMap, ok := d.peerDb.mp[nid]
	if !ok ***REMOVED***
		d.peerDb.mp[nid] = &peerMap***REMOVED***
			mp: common.NewSetMatrix(),
		***REMOVED***

		pMap = d.peerDb.mp[nid]
	***REMOVED***
	d.peerDb.Unlock()

	pKey := peerKey***REMOVED***
		peerIP:  peerIP,
		peerMac: peerMac,
	***REMOVED***

	pEntry := peerEntry***REMOVED***
		eid:        eid,
		vtep:       vtep,
		peerIPMask: peerIPMask,
		isLocal:    isLocal,
	***REMOVED***

	pMap.Lock()
	defer pMap.Unlock()
	b, i := pMap.mp.Insert(pKey.String(), pEntry.MarshalDB())
	if i != 1 ***REMOVED***
		// Transient case, there is more than one endpoint that is using the same IP,MAC pair
		s, _ := pMap.mp.String(pKey.String())
		logrus.Warnf("peerDbAdd transient condition - Key:%s cardinality:%d db state:%s", pKey.String(), i, s)
	***REMOVED***
	return b, i
***REMOVED***

func (d *driver) peerDbDelete(nid, eid string, peerIP net.IP, peerIPMask net.IPMask,
	peerMac net.HardwareAddr, vtep net.IP, isLocal bool) (bool, int) ***REMOVED***

	d.peerDb.Lock()
	pMap, ok := d.peerDb.mp[nid]
	if !ok ***REMOVED***
		d.peerDb.Unlock()
		return false, 0
	***REMOVED***
	d.peerDb.Unlock()

	pKey := peerKey***REMOVED***
		peerIP:  peerIP,
		peerMac: peerMac,
	***REMOVED***

	pEntry := peerEntry***REMOVED***
		eid:        eid,
		vtep:       vtep,
		peerIPMask: peerIPMask,
		isLocal:    isLocal,
	***REMOVED***

	pMap.Lock()
	defer pMap.Unlock()
	b, i := pMap.mp.Remove(pKey.String(), pEntry.MarshalDB())
	if i != 0 ***REMOVED***
		// Transient case, there is more than one endpoint that is using the same IP,MAC pair
		s, _ := pMap.mp.String(pKey.String())
		logrus.Warnf("peerDbDelete transient condition - Key:%s cardinality:%d db state:%s", pKey.String(), i, s)
	***REMOVED***
	return b, i
***REMOVED***

// The overlay uses a lazy initialization approach, this means that when a network is created
// and the driver registered the overlay does not allocate resources till the moment that a
// sandbox is actually created.
// At the moment of this call, that happens when a sandbox is initialized, is possible that
// networkDB has already delivered some events of peers already available on remote nodes,
// these peers are saved into the peerDB and this function is used to properly configure
// the network sandbox with all those peers that got previously notified.
// Note also that this method sends a single message on the channel and the go routine on the
// other side, will atomically loop on the whole table of peers and will program their state
// in one single atomic operation. This is fundamental to guarantee consistency, and avoid that
// new peerAdd or peerDelete gets reordered during the sandbox init.
func (d *driver) initSandboxPeerDB(nid string) ***REMOVED***
	d.peerInit(nid)
***REMOVED***

type peerOperationType int32

const (
	peerOperationINIT peerOperationType = iota
	peerOperationADD
	peerOperationDELETE
	peerOperationFLUSH
)

type peerOperation struct ***REMOVED***
	opType     peerOperationType
	networkID  string
	endpointID string
	peerIP     net.IP
	peerIPMask net.IPMask
	peerMac    net.HardwareAddr
	vtepIP     net.IP
	l2Miss     bool
	l3Miss     bool
	localPeer  bool
	callerName string
***REMOVED***

func (d *driver) peerOpRoutine(ctx context.Context, ch chan *peerOperation) ***REMOVED***
	var err error
	for ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return
		case op := <-ch:
			switch op.opType ***REMOVED***
			case peerOperationINIT:
				err = d.peerInitOp(op.networkID)
			case peerOperationADD:
				err = d.peerAddOp(op.networkID, op.endpointID, op.peerIP, op.peerIPMask, op.peerMac, op.vtepIP, op.l2Miss, op.l3Miss, true, op.localPeer)
			case peerOperationDELETE:
				err = d.peerDeleteOp(op.networkID, op.endpointID, op.peerIP, op.peerIPMask, op.peerMac, op.vtepIP, op.localPeer)
			case peerOperationFLUSH:
				err = d.peerFlushOp(op.networkID)
			***REMOVED***
			if err != nil ***REMOVED***
				logrus.Warnf("Peer operation failed:%s op:%v", err, op)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (d *driver) peerInit(nid string) ***REMOVED***
	callerName := common.CallerName(1)
	d.peerOpCh <- &peerOperation***REMOVED***
		opType:     peerOperationINIT,
		networkID:  nid,
		callerName: callerName,
	***REMOVED***
***REMOVED***

func (d *driver) peerInitOp(nid string) error ***REMOVED***
	return d.peerDbNetworkWalk(nid, func(pKey *peerKey, pEntry *peerEntry) bool ***REMOVED***
		// Local entries do not need to be added
		if pEntry.isLocal ***REMOVED***
			return false
		***REMOVED***

		d.peerAddOp(nid, pEntry.eid, pKey.peerIP, pEntry.peerIPMask, pKey.peerMac, pEntry.vtep, false, false, false, pEntry.isLocal)
		// return false to loop on all entries
		return false
	***REMOVED***)
***REMOVED***

func (d *driver) peerAdd(nid, eid string, peerIP net.IP, peerIPMask net.IPMask,
	peerMac net.HardwareAddr, vtep net.IP, l2Miss, l3Miss, localPeer bool) ***REMOVED***
	d.peerOpCh <- &peerOperation***REMOVED***
		opType:     peerOperationADD,
		networkID:  nid,
		endpointID: eid,
		peerIP:     peerIP,
		peerIPMask: peerIPMask,
		peerMac:    peerMac,
		vtepIP:     vtep,
		l2Miss:     l2Miss,
		l3Miss:     l3Miss,
		localPeer:  localPeer,
		callerName: common.CallerName(1),
	***REMOVED***
***REMOVED***

func (d *driver) peerAddOp(nid, eid string, peerIP net.IP, peerIPMask net.IPMask,
	peerMac net.HardwareAddr, vtep net.IP, l2Miss, l3Miss, updateDB, localPeer bool) error ***REMOVED***

	if err := validateID(nid, eid); err != nil ***REMOVED***
		return err
	***REMOVED***

	var dbEntries int
	var inserted bool
	if updateDB ***REMOVED***
		inserted, dbEntries = d.peerDbAdd(nid, eid, peerIP, peerIPMask, peerMac, vtep, localPeer)
		if !inserted ***REMOVED***
			logrus.Warnf("Entry already present in db: nid:%s eid:%s peerIP:%v peerMac:%v isLocal:%t vtep:%v",
				nid, eid, peerIP, peerMac, localPeer, vtep)
		***REMOVED***
	***REMOVED***

	// Local peers do not need any further configuration
	if localPeer ***REMOVED***
		return nil
	***REMOVED***

	n := d.network(nid)
	if n == nil ***REMOVED***
		return nil
	***REMOVED***

	sbox := n.sandbox()
	if sbox == nil ***REMOVED***
		// We are hitting this case for all the events that are arriving before that the sandbox
		// is being created. The peer got already added into the database and the sanbox init will
		// call the peerDbUpdateSandbox that will configure all these peers from the database
		return nil
	***REMOVED***

	IP := &net.IPNet***REMOVED***
		IP:   peerIP,
		Mask: peerIPMask,
	***REMOVED***

	s := n.getSubnetforIP(IP)
	if s == nil ***REMOVED***
		return fmt.Errorf("couldn't find the subnet %q in network %q", IP.String(), n.id)
	***REMOVED***

	if err := n.obtainVxlanID(s); err != nil ***REMOVED***
		return fmt.Errorf("couldn't get vxlan id for %q: %v", s.subnetIP.String(), err)
	***REMOVED***

	if err := n.joinSubnetSandbox(s, false); err != nil ***REMOVED***
		return fmt.Errorf("subnet sandbox join failed for %q: %v", s.subnetIP.String(), err)
	***REMOVED***

	if err := d.checkEncryption(nid, vtep, n.vxlanID(s), false, true); err != nil ***REMOVED***
		logrus.Warn(err)
	***REMOVED***

	// Add neighbor entry for the peer IP
	if err := sbox.AddNeighbor(peerIP, peerMac, l3Miss, sbox.NeighborOptions().LinkName(s.vxlanName)); err != nil ***REMOVED***
		if _, ok := err.(osl.NeighborSearchError); ok && dbEntries > 1 ***REMOVED***
			// We are in the transient case so only the first configuration is programmed into the kernel
			// Upon deletion if the active configuration is deleted the next one from the database will be restored
			// Note we are skipping also the next configuration
			return nil
		***REMOVED***
		return fmt.Errorf("could not add neighbor entry for nid:%s eid:%s into the sandbox:%v", nid, eid, err)
	***REMOVED***

	// Add fdb entry to the bridge for the peer mac
	if err := sbox.AddNeighbor(vtep, peerMac, l2Miss, sbox.NeighborOptions().LinkName(s.vxlanName),
		sbox.NeighborOptions().Family(syscall.AF_BRIDGE)); err != nil ***REMOVED***
		return fmt.Errorf("could not add fdb entry for nid:%s eid:%s into the sandbox:%v", nid, eid, err)
	***REMOVED***

	return nil
***REMOVED***

func (d *driver) peerDelete(nid, eid string, peerIP net.IP, peerIPMask net.IPMask,
	peerMac net.HardwareAddr, vtep net.IP, localPeer bool) ***REMOVED***
	d.peerOpCh <- &peerOperation***REMOVED***
		opType:     peerOperationDELETE,
		networkID:  nid,
		endpointID: eid,
		peerIP:     peerIP,
		peerIPMask: peerIPMask,
		peerMac:    peerMac,
		vtepIP:     vtep,
		callerName: common.CallerName(1),
		localPeer:  localPeer,
	***REMOVED***
***REMOVED***

func (d *driver) peerDeleteOp(nid, eid string, peerIP net.IP, peerIPMask net.IPMask,
	peerMac net.HardwareAddr, vtep net.IP, localPeer bool) error ***REMOVED***

	if err := validateID(nid, eid); err != nil ***REMOVED***
		return err
	***REMOVED***

	deleted, dbEntries := d.peerDbDelete(nid, eid, peerIP, peerIPMask, peerMac, vtep, localPeer)
	if !deleted ***REMOVED***
		logrus.Warnf("Entry was not in db: nid:%s eid:%s peerIP:%v peerMac:%v isLocal:%t vtep:%v",
			nid, eid, peerIP, peerMac, localPeer, vtep)
	***REMOVED***

	n := d.network(nid)
	if n == nil ***REMOVED***
		return nil
	***REMOVED***

	sbox := n.sandbox()
	if sbox == nil ***REMOVED***
		return nil
	***REMOVED***

	if err := d.checkEncryption(nid, vtep, 0, localPeer, false); err != nil ***REMOVED***
		logrus.Warn(err)
	***REMOVED***

	// Local peers do not have any local configuration to delete
	if !localPeer ***REMOVED***
		// Remove fdb entry to the bridge for the peer mac
		if err := sbox.DeleteNeighbor(vtep, peerMac, true); err != nil ***REMOVED***
			if _, ok := err.(osl.NeighborSearchError); ok && dbEntries > 0 ***REMOVED***
				// We fall in here if there is a transient state and if the neighbor that is being deleted
				// was never been configured into the kernel (we allow only 1 configuration at the time per <ip,mac> mapping)
				return nil
			***REMOVED***
			return fmt.Errorf("could not delete fdb entry for nid:%s eid:%s into the sandbox:%v", nid, eid, err)
		***REMOVED***

		// Delete neighbor entry for the peer IP
		if err := sbox.DeleteNeighbor(peerIP, peerMac, true); err != nil ***REMOVED***
			return fmt.Errorf("could not delete neighbor entry for nid:%s eid:%s into the sandbox:%v", nid, eid, err)
		***REMOVED***
	***REMOVED***

	if dbEntries == 0 ***REMOVED***
		return nil
	***REMOVED***

	// If there is still an entry into the database and the deletion went through without errors means that there is now no
	// configuration active in the kernel.
	// Restore one configuration for the <ip,mac> directly from the database, note that is guaranteed that there is one
	peerKey, peerEntry, err := d.peerDbSearch(nid, peerIP)
	if err != nil ***REMOVED***
		logrus.Errorf("peerDeleteOp unable to restore a configuration for nid:%s ip:%v mac:%v err:%s", nid, peerIP, peerMac, err)
		return err
	***REMOVED***
	return d.peerAddOp(nid, peerEntry.eid, peerIP, peerEntry.peerIPMask, peerKey.peerMac, peerEntry.vtep, false, false, false, peerEntry.isLocal)
***REMOVED***

func (d *driver) peerFlush(nid string) ***REMOVED***
	d.peerOpCh <- &peerOperation***REMOVED***
		opType:     peerOperationFLUSH,
		networkID:  nid,
		callerName: common.CallerName(1),
	***REMOVED***
***REMOVED***

func (d *driver) peerFlushOp(nid string) error ***REMOVED***
	d.peerDb.Lock()
	defer d.peerDb.Unlock()
	_, ok := d.peerDb.mp[nid]
	if !ok ***REMOVED***
		return fmt.Errorf("Unable to find the peerDB for nid:%s", nid)
	***REMOVED***
	delete(d.peerDb.mp, nid)
	return nil
***REMOVED***

func (d *driver) pushLocalDb() ***REMOVED***
	d.peerDbWalk(func(nid string, pKey *peerKey, pEntry *peerEntry) bool ***REMOVED***
		if pEntry.isLocal ***REMOVED***
			d.pushLocalEndpointEvent("join", nid, pEntry.eid)
		***REMOVED***
		return false
	***REMOVED***)
***REMOVED***

func (d *driver) peerDBUpdateSelf() ***REMOVED***
	d.peerDbWalk(func(nid string, pkey *peerKey, pEntry *peerEntry) bool ***REMOVED***
		if pEntry.isLocal ***REMOVED***
			pEntry.vtep = net.ParseIP(d.advertiseAddress)
		***REMOVED***
		return false
	***REMOVED***)
***REMOVED***
