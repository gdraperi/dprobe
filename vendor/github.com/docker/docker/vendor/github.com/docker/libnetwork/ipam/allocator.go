package ipam

import (
	"fmt"
	"net"
	"sort"
	"sync"

	"github.com/docker/libnetwork/bitseq"
	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/discoverapi"
	"github.com/docker/libnetwork/ipamapi"
	"github.com/docker/libnetwork/ipamutils"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

const (
	localAddressSpace  = "LocalDefault"
	globalAddressSpace = "GlobalDefault"
	// The biggest configurable host subnets
	minNetSize   = 8
	minNetSizeV6 = 64
	// datastore keyes for ipam objects
	dsConfigKey = "ipam/" + ipamapi.DefaultIPAM + "/config"
	dsDataKey   = "ipam/" + ipamapi.DefaultIPAM + "/data"
)

// Allocator provides per address space ipv4/ipv6 book keeping
type Allocator struct ***REMOVED***
	// Predefined pools for default address spaces
	predefined map[string][]*net.IPNet
	addrSpaces map[string]*addrSpace
	// stores        []datastore.Datastore
	// Allocated addresses in each address space's subnet
	addresses map[SubnetKey]*bitseq.Handle
	sync.Mutex
***REMOVED***

// NewAllocator returns an instance of libnetwork ipam
func NewAllocator(lcDs, glDs datastore.DataStore) (*Allocator, error) ***REMOVED***
	a := &Allocator***REMOVED******REMOVED***

	// Load predefined subnet pools
	a.predefined = map[string][]*net.IPNet***REMOVED***
		localAddressSpace:  ipamutils.PredefinedBroadNetworks,
		globalAddressSpace: ipamutils.PredefinedGranularNetworks,
	***REMOVED***

	// Initialize bitseq map
	a.addresses = make(map[SubnetKey]*bitseq.Handle)

	// Initialize address spaces
	a.addrSpaces = make(map[string]*addrSpace)
	for _, aspc := range []struct ***REMOVED***
		as string
		ds datastore.DataStore
	***REMOVED******REMOVED***
		***REMOVED***localAddressSpace, lcDs***REMOVED***,
		***REMOVED***globalAddressSpace, glDs***REMOVED***,
	***REMOVED*** ***REMOVED***
		a.initializeAddressSpace(aspc.as, aspc.ds)
	***REMOVED***

	return a, nil
***REMOVED***

func (a *Allocator) refresh(as string) error ***REMOVED***
	aSpace, err := a.getAddressSpaceFromStore(as)
	if err != nil ***REMOVED***
		return types.InternalErrorf("error getting pools config from store: %v", err)
	***REMOVED***

	if aSpace == nil ***REMOVED***
		return nil
	***REMOVED***

	a.Lock()
	a.addrSpaces[as] = aSpace
	a.Unlock()

	return nil
***REMOVED***

func (a *Allocator) updateBitMasks(aSpace *addrSpace) error ***REMOVED***
	var inserterList []func() error

	aSpace.Lock()
	for k, v := range aSpace.subnets ***REMOVED***
		if v.Range == nil ***REMOVED***
			kk := k
			vv := v
			inserterList = append(inserterList, func() error ***REMOVED*** return a.insertBitMask(kk, vv.Pool) ***REMOVED***)
		***REMOVED***
	***REMOVED***
	aSpace.Unlock()

	// Add the bitmasks (data could come from datastore)
	if inserterList != nil ***REMOVED***
		for _, f := range inserterList ***REMOVED***
			if err := f(); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// Checks for and fixes damaged bitmask.
func (a *Allocator) checkConsistency(as string) ***REMOVED***
	var sKeyList []SubnetKey

	// Retrieve this address space's configuration and bitmasks from the datastore
	a.refresh(as)
	a.Lock()
	aSpace, ok := a.addrSpaces[as]
	a.Unlock()
	if !ok ***REMOVED***
		return
	***REMOVED***
	a.updateBitMasks(aSpace)

	aSpace.Lock()
	for sk, pd := range aSpace.subnets ***REMOVED***
		if pd.Range != nil ***REMOVED***
			continue
		***REMOVED***
		sKeyList = append(sKeyList, sk)
	***REMOVED***
	aSpace.Unlock()

	for _, sk := range sKeyList ***REMOVED***
		a.Lock()
		bm := a.addresses[sk]
		a.Unlock()
		if err := bm.CheckConsistency(); err != nil ***REMOVED***
			logrus.Warnf("Error while running consistency check for %s: %v", sk, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (a *Allocator) initializeAddressSpace(as string, ds datastore.DataStore) error ***REMOVED***
	scope := ""
	if ds != nil ***REMOVED***
		scope = ds.Scope()
	***REMOVED***

	a.Lock()
	if currAS, ok := a.addrSpaces[as]; ok ***REMOVED***
		if currAS.ds != nil ***REMOVED***
			a.Unlock()
			return types.ForbiddenErrorf("a datastore is already configured for the address space %s", as)
		***REMOVED***
	***REMOVED***
	a.addrSpaces[as] = &addrSpace***REMOVED***
		subnets: map[SubnetKey]*PoolData***REMOVED******REMOVED***,
		id:      dsConfigKey + "/" + as,
		scope:   scope,
		ds:      ds,
		alloc:   a,
	***REMOVED***
	a.Unlock()

	a.checkConsistency(as)

	return nil
***REMOVED***

// DiscoverNew informs the allocator about a new global scope datastore
func (a *Allocator) DiscoverNew(dType discoverapi.DiscoveryType, data interface***REMOVED******REMOVED***) error ***REMOVED***
	if dType != discoverapi.DatastoreConfig ***REMOVED***
		return nil
	***REMOVED***

	dsc, ok := data.(discoverapi.DatastoreConfigData)
	if !ok ***REMOVED***
		return types.InternalErrorf("incorrect data in datastore update notification: %v", data)
	***REMOVED***

	ds, err := datastore.NewDataStoreFromConfig(dsc)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return a.initializeAddressSpace(globalAddressSpace, ds)
***REMOVED***

// DiscoverDelete is a notification of no interest for the allocator
func (a *Allocator) DiscoverDelete(dType discoverapi.DiscoveryType, data interface***REMOVED******REMOVED***) error ***REMOVED***
	return nil
***REMOVED***

// GetDefaultAddressSpaces returns the local and global default address spaces
func (a *Allocator) GetDefaultAddressSpaces() (string, string, error) ***REMOVED***
	return localAddressSpace, globalAddressSpace, nil
***REMOVED***

// RequestPool returns an address pool along with its unique id.
func (a *Allocator) RequestPool(addressSpace, pool, subPool string, options map[string]string, v6 bool) (string, *net.IPNet, map[string]string, error) ***REMOVED***
	logrus.Debugf("RequestPool(%s, %s, %s, %v, %t)", addressSpace, pool, subPool, options, v6)

	k, nw, ipr, err := a.parsePoolRequest(addressSpace, pool, subPool, v6)
	if err != nil ***REMOVED***
		return "", nil, nil, types.InternalErrorf("failed to parse pool request for address space %q pool %q subpool %q: %v", addressSpace, pool, subPool, err)
	***REMOVED***

	pdf := k == nil

retry:
	if pdf ***REMOVED***
		if nw, err = a.getPredefinedPool(addressSpace, v6); err != nil ***REMOVED***
			return "", nil, nil, err
		***REMOVED***
		k = &SubnetKey***REMOVED***AddressSpace: addressSpace, Subnet: nw.String()***REMOVED***
	***REMOVED***

	if err := a.refresh(addressSpace); err != nil ***REMOVED***
		return "", nil, nil, err
	***REMOVED***

	aSpace, err := a.getAddrSpace(addressSpace)
	if err != nil ***REMOVED***
		return "", nil, nil, err
	***REMOVED***

	insert, err := aSpace.updatePoolDBOnAdd(*k, nw, ipr, pdf)
	if err != nil ***REMOVED***
		if _, ok := err.(types.MaskableError); ok ***REMOVED***
			logrus.Debugf("Retrying predefined pool search: %v", err)
			goto retry
		***REMOVED***
		return "", nil, nil, err
	***REMOVED***

	if err := a.writeToStore(aSpace); err != nil ***REMOVED***
		if _, ok := err.(types.RetryError); !ok ***REMOVED***
			return "", nil, nil, types.InternalErrorf("pool configuration failed because of %s", err.Error())
		***REMOVED***

		goto retry
	***REMOVED***

	return k.String(), nw, nil, insert()
***REMOVED***

// ReleasePool releases the address pool identified by the passed id
func (a *Allocator) ReleasePool(poolID string) error ***REMOVED***
	logrus.Debugf("ReleasePool(%s)", poolID)
	k := SubnetKey***REMOVED******REMOVED***
	if err := k.FromString(poolID); err != nil ***REMOVED***
		return types.BadRequestErrorf("invalid pool id: %s", poolID)
	***REMOVED***

retry:
	if err := a.refresh(k.AddressSpace); err != nil ***REMOVED***
		return err
	***REMOVED***

	aSpace, err := a.getAddrSpace(k.AddressSpace)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	remove, err := aSpace.updatePoolDBOnRemoval(k)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err = a.writeToStore(aSpace); err != nil ***REMOVED***
		if _, ok := err.(types.RetryError); !ok ***REMOVED***
			return types.InternalErrorf("pool (%s) removal failed because of %v", poolID, err)
		***REMOVED***
		goto retry
	***REMOVED***

	return remove()
***REMOVED***

// Given the address space, returns the local or global PoolConfig based on the
// address space is local or global. AddressSpace locality is being registered with IPAM out of band.
func (a *Allocator) getAddrSpace(as string) (*addrSpace, error) ***REMOVED***
	a.Lock()
	defer a.Unlock()
	aSpace, ok := a.addrSpaces[as]
	if !ok ***REMOVED***
		return nil, types.BadRequestErrorf("cannot find address space %s (most likely the backing datastore is not configured)", as)
	***REMOVED***
	return aSpace, nil
***REMOVED***

func (a *Allocator) parsePoolRequest(addressSpace, pool, subPool string, v6 bool) (*SubnetKey, *net.IPNet, *AddressRange, error) ***REMOVED***
	var (
		nw  *net.IPNet
		ipr *AddressRange
		err error
	)

	if addressSpace == "" ***REMOVED***
		return nil, nil, nil, ipamapi.ErrInvalidAddressSpace
	***REMOVED***

	if pool == "" && subPool != "" ***REMOVED***
		return nil, nil, nil, ipamapi.ErrInvalidSubPool
	***REMOVED***

	if pool == "" ***REMOVED***
		return nil, nil, nil, nil
	***REMOVED***

	if _, nw, err = net.ParseCIDR(pool); err != nil ***REMOVED***
		return nil, nil, nil, ipamapi.ErrInvalidPool
	***REMOVED***

	if subPool != "" ***REMOVED***
		if ipr, err = getAddressRange(subPool, nw); err != nil ***REMOVED***
			return nil, nil, nil, err
		***REMOVED***
	***REMOVED***

	return &SubnetKey***REMOVED***AddressSpace: addressSpace, Subnet: nw.String(), ChildSubnet: subPool***REMOVED***, nw, ipr, nil
***REMOVED***

func (a *Allocator) insertBitMask(key SubnetKey, pool *net.IPNet) error ***REMOVED***
	//logrus.Debugf("Inserting bitmask (%s, %s)", key.String(), pool.String())

	store := a.getStore(key.AddressSpace)
	ipVer := getAddressVersion(pool.IP)
	ones, bits := pool.Mask.Size()
	numAddresses := uint64(1 << uint(bits-ones))

	// Allow /64 subnet
	if ipVer == v6 && numAddresses == 0 ***REMOVED***
		numAddresses--
	***REMOVED***

	// Generate the new address masks. AddressMask content may come from datastore
	h, err := bitseq.NewHandle(dsDataKey, store, key.String(), numAddresses)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Do not let network identifier address be reserved
	// Do the same for IPv6 so that bridge ip starts with XXXX...::1
	h.Set(0)

	// Do not let broadcast address be reserved
	if ipVer == v4 ***REMOVED***
		h.Set(numAddresses - 1)
	***REMOVED***

	a.Lock()
	a.addresses[key] = h
	a.Unlock()
	return nil
***REMOVED***

func (a *Allocator) retrieveBitmask(k SubnetKey, n *net.IPNet) (*bitseq.Handle, error) ***REMOVED***
	a.Lock()
	bm, ok := a.addresses[k]
	a.Unlock()
	if !ok ***REMOVED***
		logrus.Debugf("Retrieving bitmask (%s, %s)", k.String(), n.String())
		if err := a.insertBitMask(k, n); err != nil ***REMOVED***
			return nil, types.InternalErrorf("could not find bitmask in datastore for %s", k.String())
		***REMOVED***
		a.Lock()
		bm = a.addresses[k]
		a.Unlock()
	***REMOVED***
	return bm, nil
***REMOVED***

func (a *Allocator) getPredefineds(as string) []*net.IPNet ***REMOVED***
	a.Lock()
	defer a.Unlock()
	l := make([]*net.IPNet, 0, len(a.predefined[as]))
	for _, pool := range a.predefined[as] ***REMOVED***
		l = append(l, pool)
	***REMOVED***
	return l
***REMOVED***

func (a *Allocator) getPredefinedPool(as string, ipV6 bool) (*net.IPNet, error) ***REMOVED***
	var v ipVersion
	v = v4
	if ipV6 ***REMOVED***
		v = v6
	***REMOVED***

	if as != localAddressSpace && as != globalAddressSpace ***REMOVED***
		return nil, types.NotImplementedErrorf("no default pool availbale for non-default addresss spaces")
	***REMOVED***

	aSpace, err := a.getAddrSpace(as)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	for _, nw := range a.getPredefineds(as) ***REMOVED***
		if v != getAddressVersion(nw.IP) ***REMOVED***
			continue
		***REMOVED***
		aSpace.Lock()
		_, ok := aSpace.subnets[SubnetKey***REMOVED***AddressSpace: as, Subnet: nw.String()***REMOVED***]
		aSpace.Unlock()
		if ok ***REMOVED***
			continue
		***REMOVED***

		if !aSpace.contains(as, nw) ***REMOVED***
			return nw, nil
		***REMOVED***
	***REMOVED***

	return nil, types.NotFoundErrorf("could not find an available, non-overlapping IPv%d address pool among the defaults to assign to the network", v)
***REMOVED***

// RequestAddress returns an address from the specified pool ID
func (a *Allocator) RequestAddress(poolID string, prefAddress net.IP, opts map[string]string) (*net.IPNet, map[string]string, error) ***REMOVED***
	logrus.Debugf("RequestAddress(%s, %v, %v)", poolID, prefAddress, opts)
	k := SubnetKey***REMOVED******REMOVED***
	if err := k.FromString(poolID); err != nil ***REMOVED***
		return nil, nil, types.BadRequestErrorf("invalid pool id: %s", poolID)
	***REMOVED***

	if err := a.refresh(k.AddressSpace); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	aSpace, err := a.getAddrSpace(k.AddressSpace)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	aSpace.Lock()
	p, ok := aSpace.subnets[k]
	if !ok ***REMOVED***
		aSpace.Unlock()
		return nil, nil, types.NotFoundErrorf("cannot find address pool for poolID:%s", poolID)
	***REMOVED***

	if prefAddress != nil && !p.Pool.Contains(prefAddress) ***REMOVED***
		aSpace.Unlock()
		return nil, nil, ipamapi.ErrIPOutOfRange
	***REMOVED***

	c := p
	for c.Range != nil ***REMOVED***
		k = c.ParentKey
		c = aSpace.subnets[k]
	***REMOVED***
	aSpace.Unlock()

	bm, err := a.retrieveBitmask(k, c.Pool)
	if err != nil ***REMOVED***
		return nil, nil, types.InternalErrorf("could not find bitmask in datastore for %s on address %v request from pool %s: %v",
			k.String(), prefAddress, poolID, err)
	***REMOVED***
	// In order to request for a serial ip address allocation, callers can pass in the option to request
	// IP allocation serially or first available IP in the subnet
	var serial bool
	if opts != nil ***REMOVED***
		if val, ok := opts[ipamapi.AllocSerialPrefix]; ok ***REMOVED***
			serial = (val == "true")
		***REMOVED***
	***REMOVED***
	ip, err := a.getAddress(p.Pool, bm, prefAddress, p.Range, serial)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	return &net.IPNet***REMOVED***IP: ip, Mask: p.Pool.Mask***REMOVED***, nil, nil
***REMOVED***

// ReleaseAddress releases the address from the specified pool ID
func (a *Allocator) ReleaseAddress(poolID string, address net.IP) error ***REMOVED***
	logrus.Debugf("ReleaseAddress(%s, %v)", poolID, address)
	k := SubnetKey***REMOVED******REMOVED***
	if err := k.FromString(poolID); err != nil ***REMOVED***
		return types.BadRequestErrorf("invalid pool id: %s", poolID)
	***REMOVED***

	if err := a.refresh(k.AddressSpace); err != nil ***REMOVED***
		return err
	***REMOVED***

	aSpace, err := a.getAddrSpace(k.AddressSpace)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	aSpace.Lock()
	p, ok := aSpace.subnets[k]
	if !ok ***REMOVED***
		aSpace.Unlock()
		return types.NotFoundErrorf("cannot find address pool for poolID:%s", poolID)
	***REMOVED***

	if address == nil ***REMOVED***
		aSpace.Unlock()
		return types.BadRequestErrorf("invalid address: nil")
	***REMOVED***

	if !p.Pool.Contains(address) ***REMOVED***
		aSpace.Unlock()
		return ipamapi.ErrIPOutOfRange
	***REMOVED***

	c := p
	for c.Range != nil ***REMOVED***
		k = c.ParentKey
		c = aSpace.subnets[k]
	***REMOVED***
	aSpace.Unlock()

	mask := p.Pool.Mask

	h, err := types.GetHostPartIP(address, mask)
	if err != nil ***REMOVED***
		return types.InternalErrorf("failed to release address %s: %v", address.String(), err)
	***REMOVED***

	bm, err := a.retrieveBitmask(k, c.Pool)
	if err != nil ***REMOVED***
		return types.InternalErrorf("could not find bitmask in datastore for %s on address %v release from pool %s: %v",
			k.String(), address, poolID, err)
	***REMOVED***

	return bm.Unset(ipToUint64(h))
***REMOVED***

func (a *Allocator) getAddress(nw *net.IPNet, bitmask *bitseq.Handle, prefAddress net.IP, ipr *AddressRange, serial bool) (net.IP, error) ***REMOVED***
	var (
		ordinal uint64
		err     error
		base    *net.IPNet
	)

	base = types.GetIPNetCopy(nw)

	if bitmask.Unselected() <= 0 ***REMOVED***
		return nil, ipamapi.ErrNoAvailableIPs
	***REMOVED***
	if ipr == nil && prefAddress == nil ***REMOVED***
		ordinal, err = bitmask.SetAny(serial)
	***REMOVED*** else if prefAddress != nil ***REMOVED***
		hostPart, e := types.GetHostPartIP(prefAddress, base.Mask)
		if e != nil ***REMOVED***
			return nil, types.InternalErrorf("failed to allocate requested address %s: %v", prefAddress.String(), e)
		***REMOVED***
		ordinal = ipToUint64(types.GetMinimalIP(hostPart))
		err = bitmask.Set(ordinal)
	***REMOVED*** else ***REMOVED***
		ordinal, err = bitmask.SetAnyInRange(ipr.Start, ipr.End, serial)
	***REMOVED***

	switch err ***REMOVED***
	case nil:
		// Convert IP ordinal for this subnet into IP address
		return generateAddress(ordinal, base), nil
	case bitseq.ErrBitAllocated:
		return nil, ipamapi.ErrIPAlreadyAllocated
	case bitseq.ErrNoBitAvailable:
		return nil, ipamapi.ErrNoAvailableIPs
	default:
		return nil, err
	***REMOVED***
***REMOVED***

// DumpDatabase dumps the internal info
func (a *Allocator) DumpDatabase() string ***REMOVED***
	a.Lock()
	aspaces := make(map[string]*addrSpace, len(a.addrSpaces))
	orderedAS := make([]string, 0, len(a.addrSpaces))
	for as, aSpace := range a.addrSpaces ***REMOVED***
		orderedAS = append(orderedAS, as)
		aspaces[as] = aSpace
	***REMOVED***
	a.Unlock()

	sort.Strings(orderedAS)

	var s string
	for _, as := range orderedAS ***REMOVED***
		aSpace := aspaces[as]
		s = fmt.Sprintf("\n\n%s Config", as)
		aSpace.Lock()
		for k, config := range aSpace.subnets ***REMOVED***
			s += fmt.Sprintf("\n%v: %v", k, config)
			if config.Range == nil ***REMOVED***
				a.retrieveBitmask(k, config.Pool)
			***REMOVED***
		***REMOVED***
		aSpace.Unlock()
	***REMOVED***

	s = fmt.Sprintf("%s\n\nBitmasks", s)
	for k, bm := range a.addresses ***REMOVED***
		s += fmt.Sprintf("\n%s: %s", k, bm)
	***REMOVED***

	return s
***REMOVED***

// IsBuiltIn returns true for builtin drivers
func (a *Allocator) IsBuiltIn() bool ***REMOVED***
	return true
***REMOVED***
