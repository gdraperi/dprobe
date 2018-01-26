package drvregistry

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/ipamapi"
	"github.com/docker/libnetwork/types"
)

type driverData struct ***REMOVED***
	driver     driverapi.Driver
	capability driverapi.Capability
***REMOVED***

type ipamData struct ***REMOVED***
	driver     ipamapi.Ipam
	capability *ipamapi.Capability
	// default address spaces are provided by ipam driver at registration time
	defaultLocalAddressSpace, defaultGlobalAddressSpace string
***REMOVED***

type driverTable map[string]*driverData
type ipamTable map[string]*ipamData

// DrvRegistry holds the registry of all network drivers and IPAM drivers that it knows about.
type DrvRegistry struct ***REMOVED***
	sync.Mutex
	drivers      driverTable
	ipamDrivers  ipamTable
	dfn          DriverNotifyFunc
	ifn          IPAMNotifyFunc
	pluginGetter plugingetter.PluginGetter
***REMOVED***

// Functors definition

// InitFunc defines the driver initialization function signature.
type InitFunc func(driverapi.DriverCallback, map[string]interface***REMOVED******REMOVED***) error

// IPAMWalkFunc defines the IPAM driver table walker function signature.
type IPAMWalkFunc func(name string, driver ipamapi.Ipam, cap *ipamapi.Capability) bool

// DriverWalkFunc defines the network driver table walker function signature.
type DriverWalkFunc func(name string, driver driverapi.Driver, capability driverapi.Capability) bool

// IPAMNotifyFunc defines the notify function signature when a new IPAM driver gets registered.
type IPAMNotifyFunc func(name string, driver ipamapi.Ipam, cap *ipamapi.Capability) error

// DriverNotifyFunc defines the notify function signature when a new network driver gets registered.
type DriverNotifyFunc func(name string, driver driverapi.Driver, capability driverapi.Capability) error

// New retruns a new driver registry handle.
func New(lDs, gDs interface***REMOVED******REMOVED***, dfn DriverNotifyFunc, ifn IPAMNotifyFunc, pg plugingetter.PluginGetter) (*DrvRegistry, error) ***REMOVED***
	r := &DrvRegistry***REMOVED***
		drivers:      make(driverTable),
		ipamDrivers:  make(ipamTable),
		dfn:          dfn,
		ifn:          ifn,
		pluginGetter: pg,
	***REMOVED***

	return r, nil
***REMOVED***

// AddDriver adds a network driver to the registry.
func (r *DrvRegistry) AddDriver(ntype string, fn InitFunc, config map[string]interface***REMOVED******REMOVED***) error ***REMOVED***
	return fn(r, config)
***REMOVED***

// WalkIPAMs walks the IPAM drivers registered in the registry and invokes the passed walk function and each one of them.
func (r *DrvRegistry) WalkIPAMs(ifn IPAMWalkFunc) ***REMOVED***
	type ipamVal struct ***REMOVED***
		name string
		data *ipamData
	***REMOVED***

	r.Lock()
	ivl := make([]ipamVal, 0, len(r.ipamDrivers))
	for k, v := range r.ipamDrivers ***REMOVED***
		ivl = append(ivl, ipamVal***REMOVED***name: k, data: v***REMOVED***)
	***REMOVED***
	r.Unlock()

	for _, iv := range ivl ***REMOVED***
		if ifn(iv.name, iv.data.driver, iv.data.capability) ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

// WalkDrivers walks the network drivers registered in the registry and invokes the passed walk function and each one of them.
func (r *DrvRegistry) WalkDrivers(dfn DriverWalkFunc) ***REMOVED***
	type driverVal struct ***REMOVED***
		name string
		data *driverData
	***REMOVED***

	r.Lock()
	dvl := make([]driverVal, 0, len(r.drivers))
	for k, v := range r.drivers ***REMOVED***
		dvl = append(dvl, driverVal***REMOVED***name: k, data: v***REMOVED***)
	***REMOVED***
	r.Unlock()

	for _, dv := range dvl ***REMOVED***
		if dfn(dv.name, dv.data.driver, dv.data.capability) ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

// Driver returns the actual network driver instance and its capability  which registered with the passed name.
func (r *DrvRegistry) Driver(name string) (driverapi.Driver, *driverapi.Capability) ***REMOVED***
	r.Lock()
	defer r.Unlock()

	d, ok := r.drivers[name]
	if !ok ***REMOVED***
		return nil, nil
	***REMOVED***

	return d.driver, &d.capability
***REMOVED***

// IPAM returns the actual IPAM driver instance and its capability which registered with the passed name.
func (r *DrvRegistry) IPAM(name string) (ipamapi.Ipam, *ipamapi.Capability) ***REMOVED***
	r.Lock()
	defer r.Unlock()

	i, ok := r.ipamDrivers[name]
	if !ok ***REMOVED***
		return nil, nil
	***REMOVED***

	return i.driver, i.capability
***REMOVED***

// IPAMDefaultAddressSpaces returns the default address space strings for the passed IPAM driver name.
func (r *DrvRegistry) IPAMDefaultAddressSpaces(name string) (string, string, error) ***REMOVED***
	r.Lock()
	defer r.Unlock()

	i, ok := r.ipamDrivers[name]
	if !ok ***REMOVED***
		return "", "", fmt.Errorf("ipam %s not found", name)
	***REMOVED***

	return i.defaultLocalAddressSpace, i.defaultGlobalAddressSpace, nil
***REMOVED***

// GetPluginGetter returns the plugingetter
func (r *DrvRegistry) GetPluginGetter() plugingetter.PluginGetter ***REMOVED***
	return r.pluginGetter
***REMOVED***

// RegisterDriver registers the network driver when it gets discovered.
func (r *DrvRegistry) RegisterDriver(ntype string, driver driverapi.Driver, capability driverapi.Capability) error ***REMOVED***
	if strings.TrimSpace(ntype) == "" ***REMOVED***
		return errors.New("network type string cannot be empty")
	***REMOVED***

	r.Lock()
	dd, ok := r.drivers[ntype]
	r.Unlock()

	if ok && dd.driver.IsBuiltIn() ***REMOVED***
		return driverapi.ErrActiveRegistration(ntype)
	***REMOVED***

	if r.dfn != nil ***REMOVED***
		if err := r.dfn(ntype, driver, capability); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	dData := &driverData***REMOVED***driver, capability***REMOVED***

	r.Lock()
	r.drivers[ntype] = dData
	r.Unlock()

	return nil
***REMOVED***

func (r *DrvRegistry) registerIpamDriver(name string, driver ipamapi.Ipam, caps *ipamapi.Capability) error ***REMOVED***
	if strings.TrimSpace(name) == "" ***REMOVED***
		return errors.New("ipam driver name string cannot be empty")
	***REMOVED***

	r.Lock()
	dd, ok := r.ipamDrivers[name]
	r.Unlock()
	if ok && dd.driver.IsBuiltIn() ***REMOVED***
		return types.ForbiddenErrorf("ipam driver %q already registered", name)
	***REMOVED***

	locAS, glbAS, err := driver.GetDefaultAddressSpaces()
	if err != nil ***REMOVED***
		return types.InternalErrorf("ipam driver %q failed to return default address spaces: %v", name, err)
	***REMOVED***

	if r.ifn != nil ***REMOVED***
		if err := r.ifn(name, driver, caps); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	r.Lock()
	r.ipamDrivers[name] = &ipamData***REMOVED***driver: driver, defaultLocalAddressSpace: locAS, defaultGlobalAddressSpace: glbAS, capability: caps***REMOVED***
	r.Unlock()

	return nil
***REMOVED***

// RegisterIpamDriver registers the IPAM driver discovered with default capabilities.
func (r *DrvRegistry) RegisterIpamDriver(name string, driver ipamapi.Ipam) error ***REMOVED***
	return r.registerIpamDriver(name, driver, &ipamapi.Capability***REMOVED******REMOVED***)
***REMOVED***

// RegisterIpamDriverWithCapabilities registers the IPAM driver discovered with specified capabilities.
func (r *DrvRegistry) RegisterIpamDriverWithCapabilities(name string, driver ipamapi.Ipam, caps *ipamapi.Capability) error ***REMOVED***
	return r.registerIpamDriver(name, driver, caps)
***REMOVED***
