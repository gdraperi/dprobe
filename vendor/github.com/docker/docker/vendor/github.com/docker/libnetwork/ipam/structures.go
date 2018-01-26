package ipam

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/ipamapi"
	"github.com/docker/libnetwork/types"
)

// SubnetKey is the pointer to the configured pools in each address space
type SubnetKey struct ***REMOVED***
	AddressSpace string
	Subnet       string
	ChildSubnet  string
***REMOVED***

// PoolData contains the configured pool data
type PoolData struct ***REMOVED***
	ParentKey SubnetKey
	Pool      *net.IPNet
	Range     *AddressRange `json:",omitempty"`
	RefCount  int
***REMOVED***

// addrSpace contains the pool configurations for the address space
type addrSpace struct ***REMOVED***
	subnets  map[SubnetKey]*PoolData
	dbIndex  uint64
	dbExists bool
	id       string
	scope    string
	ds       datastore.DataStore
	alloc    *Allocator
	sync.Mutex
***REMOVED***

// AddressRange specifies first and last ip ordinal which
// identifies a range in a pool of addresses
type AddressRange struct ***REMOVED***
	Sub        *net.IPNet
	Start, End uint64
***REMOVED***

// String returns the string form of the AddressRange object
func (r *AddressRange) String() string ***REMOVED***
	return fmt.Sprintf("Sub: %s, range [%d, %d]", r.Sub, r.Start, r.End)
***REMOVED***

// MarshalJSON returns the JSON encoding of the Range object
func (r *AddressRange) MarshalJSON() ([]byte, error) ***REMOVED***
	m := map[string]interface***REMOVED******REMOVED******REMOVED***
		"Sub":   r.Sub.String(),
		"Start": r.Start,
		"End":   r.End,
	***REMOVED***
	return json.Marshal(m)
***REMOVED***

// UnmarshalJSON decodes data into the Range object
func (r *AddressRange) UnmarshalJSON(data []byte) error ***REMOVED***
	m := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	err := json.Unmarshal(data, &m)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if r.Sub, err = types.ParseCIDR(m["Sub"].(string)); err != nil ***REMOVED***
		return err
	***REMOVED***
	r.Start = uint64(m["Start"].(float64))
	r.End = uint64(m["End"].(float64))
	return nil
***REMOVED***

// String returns the string form of the SubnetKey object
func (s *SubnetKey) String() string ***REMOVED***
	k := fmt.Sprintf("%s/%s", s.AddressSpace, s.Subnet)
	if s.ChildSubnet != "" ***REMOVED***
		k = fmt.Sprintf("%s/%s", k, s.ChildSubnet)
	***REMOVED***
	return k
***REMOVED***

// FromString populates the SubnetKey object reading it from string
func (s *SubnetKey) FromString(str string) error ***REMOVED***
	if str == "" || !strings.Contains(str, "/") ***REMOVED***
		return types.BadRequestErrorf("invalid string form for subnetkey: %s", str)
	***REMOVED***

	p := strings.Split(str, "/")
	if len(p) != 3 && len(p) != 5 ***REMOVED***
		return types.BadRequestErrorf("invalid string form for subnetkey: %s", str)
	***REMOVED***
	s.AddressSpace = p[0]
	s.Subnet = fmt.Sprintf("%s/%s", p[1], p[2])
	if len(p) == 5 ***REMOVED***
		s.ChildSubnet = fmt.Sprintf("%s/%s", p[3], p[4])
	***REMOVED***

	return nil
***REMOVED***

// String returns the string form of the PoolData object
func (p *PoolData) String() string ***REMOVED***
	return fmt.Sprintf("ParentKey: %s, Pool: %s, Range: %s, RefCount: %d",
		p.ParentKey.String(), p.Pool.String(), p.Range, p.RefCount)
***REMOVED***

// MarshalJSON returns the JSON encoding of the PoolData object
func (p *PoolData) MarshalJSON() ([]byte, error) ***REMOVED***
	m := map[string]interface***REMOVED******REMOVED******REMOVED***
		"ParentKey": p.ParentKey,
		"RefCount":  p.RefCount,
	***REMOVED***
	if p.Pool != nil ***REMOVED***
		m["Pool"] = p.Pool.String()
	***REMOVED***
	if p.Range != nil ***REMOVED***
		m["Range"] = p.Range
	***REMOVED***
	return json.Marshal(m)
***REMOVED***

// UnmarshalJSON decodes data into the PoolData object
func (p *PoolData) UnmarshalJSON(data []byte) error ***REMOVED***
	var (
		err error
		t   struct ***REMOVED***
			ParentKey SubnetKey
			Pool      string
			Range     *AddressRange `json:",omitempty"`
			RefCount  int
		***REMOVED***
	)

	if err = json.Unmarshal(data, &t); err != nil ***REMOVED***
		return err
	***REMOVED***

	p.ParentKey = t.ParentKey
	p.Range = t.Range
	p.RefCount = t.RefCount
	if t.Pool != "" ***REMOVED***
		if p.Pool, err = types.ParseCIDR(t.Pool); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// MarshalJSON returns the JSON encoding of the addrSpace object
func (aSpace *addrSpace) MarshalJSON() ([]byte, error) ***REMOVED***
	aSpace.Lock()
	defer aSpace.Unlock()

	m := map[string]interface***REMOVED******REMOVED******REMOVED***
		"Scope": string(aSpace.scope),
	***REMOVED***

	if aSpace.subnets != nil ***REMOVED***
		s := map[string]*PoolData***REMOVED******REMOVED***
		for k, v := range aSpace.subnets ***REMOVED***
			s[k.String()] = v
		***REMOVED***
		m["Subnets"] = s
	***REMOVED***

	return json.Marshal(m)
***REMOVED***

// UnmarshalJSON decodes data into the addrSpace object
func (aSpace *addrSpace) UnmarshalJSON(data []byte) error ***REMOVED***
	aSpace.Lock()
	defer aSpace.Unlock()

	m := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	err := json.Unmarshal(data, &m)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	aSpace.scope = datastore.LocalScope
	s := m["Scope"].(string)
	if s == string(datastore.GlobalScope) ***REMOVED***
		aSpace.scope = datastore.GlobalScope
	***REMOVED***

	if v, ok := m["Subnets"]; ok ***REMOVED***
		sb, _ := json.Marshal(v)
		var s map[string]*PoolData
		err := json.Unmarshal(sb, &s)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		for ks, v := range s ***REMOVED***
			k := SubnetKey***REMOVED******REMOVED***
			k.FromString(ks)
			aSpace.subnets[k] = v
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// CopyTo deep copies the pool data to the destination pooldata
func (p *PoolData) CopyTo(dstP *PoolData) error ***REMOVED***
	dstP.ParentKey = p.ParentKey
	dstP.Pool = types.GetIPNetCopy(p.Pool)

	if p.Range != nil ***REMOVED***
		dstP.Range = &AddressRange***REMOVED******REMOVED***
		dstP.Range.Sub = types.GetIPNetCopy(p.Range.Sub)
		dstP.Range.Start = p.Range.Start
		dstP.Range.End = p.Range.End
	***REMOVED***

	dstP.RefCount = p.RefCount
	return nil
***REMOVED***

func (aSpace *addrSpace) CopyTo(o datastore.KVObject) error ***REMOVED***
	aSpace.Lock()
	defer aSpace.Unlock()

	dstAspace := o.(*addrSpace)

	dstAspace.id = aSpace.id
	dstAspace.ds = aSpace.ds
	dstAspace.alloc = aSpace.alloc
	dstAspace.scope = aSpace.scope
	dstAspace.dbIndex = aSpace.dbIndex
	dstAspace.dbExists = aSpace.dbExists

	dstAspace.subnets = make(map[SubnetKey]*PoolData)
	for k, v := range aSpace.subnets ***REMOVED***
		dstAspace.subnets[k] = &PoolData***REMOVED******REMOVED***
		v.CopyTo(dstAspace.subnets[k])
	***REMOVED***

	return nil
***REMOVED***

func (aSpace *addrSpace) New() datastore.KVObject ***REMOVED***
	aSpace.Lock()
	defer aSpace.Unlock()

	return &addrSpace***REMOVED***
		id:    aSpace.id,
		ds:    aSpace.ds,
		alloc: aSpace.alloc,
		scope: aSpace.scope,
	***REMOVED***
***REMOVED***

func (aSpace *addrSpace) updatePoolDBOnAdd(k SubnetKey, nw *net.IPNet, ipr *AddressRange, pdf bool) (func() error, error) ***REMOVED***
	aSpace.Lock()
	defer aSpace.Unlock()

	// Check if already allocated
	if p, ok := aSpace.subnets[k]; ok ***REMOVED***
		if pdf ***REMOVED***
			return nil, types.InternalMaskableErrorf("predefined pool %s is already reserved", nw)
		***REMOVED***
		aSpace.incRefCount(p, 1)
		return func() error ***REMOVED*** return nil ***REMOVED***, nil
	***REMOVED***

	// If master pool, check for overlap
	if ipr == nil ***REMOVED***
		if aSpace.contains(k.AddressSpace, nw) ***REMOVED***
			return nil, ipamapi.ErrPoolOverlap
		***REMOVED***
		// This is a new master pool, add it along with corresponding bitmask
		aSpace.subnets[k] = &PoolData***REMOVED***Pool: nw, RefCount: 1***REMOVED***
		return func() error ***REMOVED*** return aSpace.alloc.insertBitMask(k, nw) ***REMOVED***, nil
	***REMOVED***

	// This is a new non-master pool
	p := &PoolData***REMOVED***
		ParentKey: SubnetKey***REMOVED***AddressSpace: k.AddressSpace, Subnet: k.Subnet***REMOVED***,
		Pool:      nw,
		Range:     ipr,
		RefCount:  1,
	***REMOVED***
	aSpace.subnets[k] = p

	// Look for parent pool
	pp, ok := aSpace.subnets[p.ParentKey]
	if ok ***REMOVED***
		aSpace.incRefCount(pp, 1)
		return func() error ***REMOVED*** return nil ***REMOVED***, nil
	***REMOVED***

	// Parent pool does not exist, add it along with corresponding bitmask
	aSpace.subnets[p.ParentKey] = &PoolData***REMOVED***Pool: nw, RefCount: 1***REMOVED***
	return func() error ***REMOVED*** return aSpace.alloc.insertBitMask(p.ParentKey, nw) ***REMOVED***, nil
***REMOVED***

func (aSpace *addrSpace) updatePoolDBOnRemoval(k SubnetKey) (func() error, error) ***REMOVED***
	aSpace.Lock()
	defer aSpace.Unlock()

	p, ok := aSpace.subnets[k]
	if !ok ***REMOVED***
		return nil, ipamapi.ErrBadPool
	***REMOVED***

	aSpace.incRefCount(p, -1)

	c := p
	for ok ***REMOVED***
		if c.RefCount == 0 ***REMOVED***
			delete(aSpace.subnets, k)
			if c.Range == nil ***REMOVED***
				return func() error ***REMOVED***
					bm, err := aSpace.alloc.retrieveBitmask(k, c.Pool)
					if err != nil ***REMOVED***
						return types.InternalErrorf("could not find bitmask in datastore for pool %s removal: %v", k.String(), err)
					***REMOVED***
					return bm.Destroy()
				***REMOVED***, nil
			***REMOVED***
		***REMOVED***
		k = c.ParentKey
		c, ok = aSpace.subnets[k]
	***REMOVED***

	return func() error ***REMOVED*** return nil ***REMOVED***, nil
***REMOVED***

func (aSpace *addrSpace) incRefCount(p *PoolData, delta int) ***REMOVED***
	c := p
	ok := true
	for ok ***REMOVED***
		c.RefCount += delta
		c, ok = aSpace.subnets[c.ParentKey]
	***REMOVED***
***REMOVED***

// Checks whether the passed subnet is a superset or subset of any of the subset in this config db
func (aSpace *addrSpace) contains(space string, nw *net.IPNet) bool ***REMOVED***
	for k, v := range aSpace.subnets ***REMOVED***
		if space == k.AddressSpace && k.ChildSubnet == "" ***REMOVED***
			if nw.Contains(v.Pool.IP) || v.Pool.Contains(nw.IP) ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (aSpace *addrSpace) store() datastore.DataStore ***REMOVED***
	aSpace.Lock()
	defer aSpace.Unlock()

	return aSpace.ds
***REMOVED***
