package libnetwork

import (
	"fmt"
	"strings"

	"github.com/docker/libkv/store/boltdb"
	"github.com/docker/libkv/store/consul"
	"github.com/docker/libkv/store/etcd"
	"github.com/docker/libkv/store/zookeeper"
	"github.com/docker/libnetwork/datastore"
	"github.com/sirupsen/logrus"
)

func registerKVStores() ***REMOVED***
	consul.Register()
	zookeeper.Register()
	etcd.Register()
	boltdb.Register()
***REMOVED***

func (c *controller) initScopedStore(scope string, scfg *datastore.ScopeCfg) error ***REMOVED***
	store, err := datastore.NewDataStore(scope, scfg)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	c.Lock()
	c.stores = append(c.stores, store)
	c.Unlock()

	return nil
***REMOVED***

func (c *controller) initStores() error ***REMOVED***
	registerKVStores()

	c.Lock()
	if c.cfg == nil ***REMOVED***
		c.Unlock()
		return nil
	***REMOVED***
	scopeConfigs := c.cfg.Scopes
	c.stores = nil
	c.Unlock()

	for scope, scfg := range scopeConfigs ***REMOVED***
		if err := c.initScopedStore(scope, scfg); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	c.startWatch()
	return nil
***REMOVED***

func (c *controller) closeStores() ***REMOVED***
	for _, store := range c.getStores() ***REMOVED***
		store.Close()
	***REMOVED***
***REMOVED***

func (c *controller) getStore(scope string) datastore.DataStore ***REMOVED***
	c.Lock()
	defer c.Unlock()

	for _, store := range c.stores ***REMOVED***
		if store.Scope() == scope ***REMOVED***
			return store
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (c *controller) getStores() []datastore.DataStore ***REMOVED***
	c.Lock()
	defer c.Unlock()

	return c.stores
***REMOVED***

func (c *controller) getNetworkFromStore(nid string) (*network, error) ***REMOVED***
	for _, store := range c.getStores() ***REMOVED***
		n := &network***REMOVED***id: nid, ctrlr: c***REMOVED***
		err := store.GetObject(datastore.Key(n.Key()...), n)
		// Continue searching in the next store if the key is not found in this store
		if err != nil ***REMOVED***
			if err != datastore.ErrKeyNotFound ***REMOVED***
				logrus.Debugf("could not find network %s: %v", nid, err)
			***REMOVED***
			continue
		***REMOVED***

		ec := &endpointCnt***REMOVED***n: n***REMOVED***
		err = store.GetObject(datastore.Key(ec.Key()...), ec)
		if err != nil && !n.inDelete ***REMOVED***
			return nil, fmt.Errorf("could not find endpoint count for network %s: %v", n.Name(), err)
		***REMOVED***

		n.epCnt = ec
		if n.scope == "" ***REMOVED***
			n.scope = store.Scope()
		***REMOVED***
		return n, nil
	***REMOVED***

	return nil, fmt.Errorf("network %s not found", nid)
***REMOVED***

func (c *controller) getNetworksForScope(scope string) ([]*network, error) ***REMOVED***
	var nl []*network

	store := c.getStore(scope)
	if store == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	kvol, err := store.List(datastore.Key(datastore.NetworkKeyPrefix),
		&network***REMOVED***ctrlr: c***REMOVED***)
	if err != nil && err != datastore.ErrKeyNotFound ***REMOVED***
		return nil, fmt.Errorf("failed to get networks for scope %s: %v",
			scope, err)
	***REMOVED***

	for _, kvo := range kvol ***REMOVED***
		n := kvo.(*network)
		n.ctrlr = c

		ec := &endpointCnt***REMOVED***n: n***REMOVED***
		err = store.GetObject(datastore.Key(ec.Key()...), ec)
		if err != nil && !n.inDelete ***REMOVED***
			logrus.Warnf("Could not find endpoint count key %s for network %s while listing: %v", datastore.Key(ec.Key()...), n.Name(), err)
			continue
		***REMOVED***

		n.epCnt = ec
		if n.scope == "" ***REMOVED***
			n.scope = scope
		***REMOVED***
		nl = append(nl, n)
	***REMOVED***

	return nl, nil
***REMOVED***

func (c *controller) getNetworksFromStore() ([]*network, error) ***REMOVED***
	var nl []*network

	for _, store := range c.getStores() ***REMOVED***
		kvol, err := store.List(datastore.Key(datastore.NetworkKeyPrefix),
			&network***REMOVED***ctrlr: c***REMOVED***)
		// Continue searching in the next store if no keys found in this store
		if err != nil ***REMOVED***
			if err != datastore.ErrKeyNotFound ***REMOVED***
				logrus.Debugf("failed to get networks for scope %s: %v", store.Scope(), err)
			***REMOVED***
			continue
		***REMOVED***

		kvep, err := store.Map(datastore.Key(epCntKeyPrefix), &endpointCnt***REMOVED******REMOVED***)
		if err != nil ***REMOVED***
			if err != datastore.ErrKeyNotFound ***REMOVED***
				logrus.Warnf("failed to get endpoint_count map for scope %s: %v", store.Scope(), err)
			***REMOVED***
		***REMOVED***

		for _, kvo := range kvol ***REMOVED***
			n := kvo.(*network)
			n.Lock()
			n.ctrlr = c
			ec := &endpointCnt***REMOVED***n: n***REMOVED***
			// Trim the leading & trailing "/" to make it consistent across all stores
			if val, ok := kvep[strings.Trim(datastore.Key(ec.Key()...), "/")]; ok ***REMOVED***
				ec = val.(*endpointCnt)
				ec.n = n
				n.epCnt = ec
			***REMOVED***
			if n.scope == "" ***REMOVED***
				n.scope = store.Scope()
			***REMOVED***
			n.Unlock()
			nl = append(nl, n)
		***REMOVED***
	***REMOVED***

	return nl, nil
***REMOVED***

func (n *network) getEndpointFromStore(eid string) (*endpoint, error) ***REMOVED***
	var errors []string
	for _, store := range n.ctrlr.getStores() ***REMOVED***
		ep := &endpoint***REMOVED***id: eid, network: n***REMOVED***
		err := store.GetObject(datastore.Key(ep.Key()...), ep)
		// Continue searching in the next store if the key is not found in this store
		if err != nil ***REMOVED***
			if err != datastore.ErrKeyNotFound ***REMOVED***
				errors = append(errors, fmt.Sprintf("***REMOVED***%s:%v***REMOVED***, ", store.Scope(), err))
				logrus.Debugf("could not find endpoint %s in %s: %v", eid, store.Scope(), err)
			***REMOVED***
			continue
		***REMOVED***
		return ep, nil
	***REMOVED***
	return nil, fmt.Errorf("could not find endpoint %s: %v", eid, errors)
***REMOVED***

func (n *network) getEndpointsFromStore() ([]*endpoint, error) ***REMOVED***
	var epl []*endpoint

	tmp := endpoint***REMOVED***network: n***REMOVED***
	for _, store := range n.getController().getStores() ***REMOVED***
		kvol, err := store.List(datastore.Key(tmp.KeyPrefix()...), &endpoint***REMOVED***network: n***REMOVED***)
		// Continue searching in the next store if no keys found in this store
		if err != nil ***REMOVED***
			if err != datastore.ErrKeyNotFound ***REMOVED***
				logrus.Debugf("failed to get endpoints for network %s scope %s: %v",
					n.Name(), store.Scope(), err)
			***REMOVED***
			continue
		***REMOVED***

		for _, kvo := range kvol ***REMOVED***
			ep := kvo.(*endpoint)
			epl = append(epl, ep)
		***REMOVED***
	***REMOVED***

	return epl, nil
***REMOVED***

func (c *controller) updateToStore(kvObject datastore.KVObject) error ***REMOVED***
	cs := c.getStore(kvObject.DataScope())
	if cs == nil ***REMOVED***
		return ErrDataStoreNotInitialized(kvObject.DataScope())
	***REMOVED***

	if err := cs.PutObjectAtomic(kvObject); err != nil ***REMOVED***
		if err == datastore.ErrKeyModified ***REMOVED***
			return err
		***REMOVED***
		return fmt.Errorf("failed to update store for object type %T: %v", kvObject, err)
	***REMOVED***

	return nil
***REMOVED***

func (c *controller) deleteFromStore(kvObject datastore.KVObject) error ***REMOVED***
	cs := c.getStore(kvObject.DataScope())
	if cs == nil ***REMOVED***
		return ErrDataStoreNotInitialized(kvObject.DataScope())
	***REMOVED***

retry:
	if err := cs.DeleteObjectAtomic(kvObject); err != nil ***REMOVED***
		if err == datastore.ErrKeyModified ***REMOVED***
			if err := cs.GetObject(datastore.Key(kvObject.Key()...), kvObject); err != nil ***REMOVED***
				return fmt.Errorf("could not update the kvobject to latest when trying to delete: %v", err)
			***REMOVED***
			logrus.Warnf("Error (%v) deleting object %v, retrying....", err, kvObject.Key())
			goto retry
		***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

type netWatch struct ***REMOVED***
	localEps  map[string]*endpoint
	remoteEps map[string]*endpoint
	stopCh    chan struct***REMOVED******REMOVED***
***REMOVED***

func (c *controller) getLocalEps(nw *netWatch) []*endpoint ***REMOVED***
	c.Lock()
	defer c.Unlock()

	var epl []*endpoint
	for _, ep := range nw.localEps ***REMOVED***
		epl = append(epl, ep)
	***REMOVED***

	return epl
***REMOVED***

func (c *controller) watchSvcRecord(ep *endpoint) ***REMOVED***
	c.watchCh <- ep
***REMOVED***

func (c *controller) unWatchSvcRecord(ep *endpoint) ***REMOVED***
	c.unWatchCh <- ep
***REMOVED***

func (c *controller) networkWatchLoop(nw *netWatch, ep *endpoint, ecCh <-chan datastore.KVObject) ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case <-nw.stopCh:
			return
		case o := <-ecCh:
			ec := o.(*endpointCnt)

			epl, err := ec.n.getEndpointsFromStore()
			if err != nil ***REMOVED***
				break
			***REMOVED***

			c.Lock()
			var addEp []*endpoint

			delEpMap := make(map[string]*endpoint)
			renameEpMap := make(map[string]bool)
			for k, v := range nw.remoteEps ***REMOVED***
				delEpMap[k] = v
			***REMOVED***

			for _, lEp := range epl ***REMOVED***
				if _, ok := nw.localEps[lEp.ID()]; ok ***REMOVED***
					continue
				***REMOVED***

				if ep, ok := nw.remoteEps[lEp.ID()]; ok ***REMOVED***
					// On a container rename EP ID will remain
					// the same but the name will change. service
					// records should reflect the change.
					// Keep old EP entry in the delEpMap and add
					// EP from the store (which has the new name)
					// into the new list
					if lEp.name == ep.name ***REMOVED***
						delete(delEpMap, lEp.ID())
						continue
					***REMOVED***
					renameEpMap[lEp.ID()] = true
				***REMOVED***
				nw.remoteEps[lEp.ID()] = lEp
				addEp = append(addEp, lEp)
			***REMOVED***

			// EPs whose name are to be deleted from the svc records
			// should also be removed from nw's remote EP list, except
			// the ones that are getting renamed.
			for _, lEp := range delEpMap ***REMOVED***
				if !renameEpMap[lEp.ID()] ***REMOVED***
					delete(nw.remoteEps, lEp.ID())
				***REMOVED***
			***REMOVED***
			c.Unlock()

			for _, lEp := range delEpMap ***REMOVED***
				ep.getNetwork().updateSvcRecord(lEp, c.getLocalEps(nw), false)

			***REMOVED***
			for _, lEp := range addEp ***REMOVED***
				ep.getNetwork().updateSvcRecord(lEp, c.getLocalEps(nw), true)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *controller) processEndpointCreate(nmap map[string]*netWatch, ep *endpoint) ***REMOVED***
	n := ep.getNetwork()
	if !c.isDistributedControl() && n.Scope() == datastore.SwarmScope && n.driverIsMultihost() ***REMOVED***
		return
	***REMOVED***

	c.Lock()
	nw, ok := nmap[n.ID()]
	c.Unlock()

	if ok ***REMOVED***
		// Update the svc db for the local endpoint join right away
		n.updateSvcRecord(ep, c.getLocalEps(nw), true)

		c.Lock()
		nw.localEps[ep.ID()] = ep

		// If we had learned that from the kv store remove it
		// from remote ep list now that we know that this is
		// indeed a local endpoint
		delete(nw.remoteEps, ep.ID())
		c.Unlock()
		return
	***REMOVED***

	nw = &netWatch***REMOVED***
		localEps:  make(map[string]*endpoint),
		remoteEps: make(map[string]*endpoint),
	***REMOVED***

	// Update the svc db for the local endpoint join right away
	// Do this before adding this ep to localEps so that we don't
	// try to update this ep's container's svc records
	n.updateSvcRecord(ep, c.getLocalEps(nw), true)

	c.Lock()
	nw.localEps[ep.ID()] = ep
	nmap[n.ID()] = nw
	nw.stopCh = make(chan struct***REMOVED******REMOVED***)
	c.Unlock()

	store := c.getStore(n.DataScope())
	if store == nil ***REMOVED***
		return
	***REMOVED***

	if !store.Watchable() ***REMOVED***
		return
	***REMOVED***

	ch, err := store.Watch(n.getEpCnt(), nw.stopCh)
	if err != nil ***REMOVED***
		logrus.Warnf("Error creating watch for network: %v", err)
		return
	***REMOVED***

	go c.networkWatchLoop(nw, ep, ch)
***REMOVED***

func (c *controller) processEndpointDelete(nmap map[string]*netWatch, ep *endpoint) ***REMOVED***
	n := ep.getNetwork()
	if !c.isDistributedControl() && n.Scope() == datastore.SwarmScope && n.driverIsMultihost() ***REMOVED***
		return
	***REMOVED***

	c.Lock()
	nw, ok := nmap[n.ID()]

	if ok ***REMOVED***
		delete(nw.localEps, ep.ID())
		c.Unlock()

		// Update the svc db about local endpoint leave right away
		// Do this after we remove this ep from localEps so that we
		// don't try to remove this svc record from this ep's container.
		n.updateSvcRecord(ep, c.getLocalEps(nw), false)

		c.Lock()
		if len(nw.localEps) == 0 ***REMOVED***
			close(nw.stopCh)

			// This is the last container going away for the network. Destroy
			// this network's svc db entry
			delete(c.svcRecords, n.ID())

			delete(nmap, n.ID())
		***REMOVED***
	***REMOVED***
	c.Unlock()
***REMOVED***

func (c *controller) watchLoop() ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case ep := <-c.watchCh:
			c.processEndpointCreate(c.nmap, ep)
		case ep := <-c.unWatchCh:
			c.processEndpointDelete(c.nmap, ep)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *controller) startWatch() ***REMOVED***
	if c.watchCh != nil ***REMOVED***
		return
	***REMOVED***
	c.watchCh = make(chan *endpoint)
	c.unWatchCh = make(chan *endpoint)
	c.nmap = make(map[string]*netWatch)

	go c.watchLoop()
***REMOVED***

func (c *controller) networkCleanup() ***REMOVED***
	networks, err := c.getNetworksFromStore()
	if err != nil ***REMOVED***
		logrus.Warnf("Could not retrieve networks from store(s) during network cleanup: %v", err)
		return
	***REMOVED***

	for _, n := range networks ***REMOVED***
		if n.inDelete ***REMOVED***
			logrus.Infof("Removing stale network %s (%s)", n.Name(), n.ID())
			if err := n.delete(true); err != nil ***REMOVED***
				logrus.Debugf("Error while removing stale network: %v", err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

var populateSpecial NetworkWalker = func(nw Network) bool ***REMOVED***
	if n := nw.(*network); n.hasSpecialDriver() && !n.ConfigOnly() ***REMOVED***
		if err := n.getController().addNetwork(n); err != nil ***REMOVED***
			logrus.Warnf("Failed to populate network %q with driver %q", nw.Name(), nw.Type())
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
