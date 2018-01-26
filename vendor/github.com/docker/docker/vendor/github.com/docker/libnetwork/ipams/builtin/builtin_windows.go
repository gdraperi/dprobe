// +build windows

package builtin

import (
	"errors"

	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/ipam"
	"github.com/docker/libnetwork/ipamapi"
	"github.com/docker/libnetwork/ipamutils"

	windowsipam "github.com/docker/libnetwork/ipams/windowsipam"
)

// InitDockerDefault registers the built-in ipam service with libnetwork
func InitDockerDefault(ic ipamapi.Callback, l, g interface***REMOVED******REMOVED***) error ***REMOVED***
	var (
		ok                bool
		localDs, globalDs datastore.DataStore
	)

	if l != nil ***REMOVED***
		if localDs, ok = l.(datastore.DataStore); !ok ***REMOVED***
			return errors.New("incorrect local datastore passed to built-in ipam init")
		***REMOVED***
	***REMOVED***

	if g != nil ***REMOVED***
		if globalDs, ok = g.(datastore.DataStore); !ok ***REMOVED***
			return errors.New("incorrect global datastore passed to built-in ipam init")
		***REMOVED***
	***REMOVED***

	ipamutils.InitNetworks()

	a, err := ipam.NewAllocator(localDs, globalDs)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	cps := &ipamapi.Capability***REMOVED***RequiresRequestReplay: true***REMOVED***

	return ic.RegisterIpamDriverWithCapabilities(ipamapi.DefaultIPAM, a, cps)
***REMOVED***

// Init registers the built-in ipam service with libnetwork
func Init(ic ipamapi.Callback, l, g interface***REMOVED******REMOVED***) error ***REMOVED***
	initFunc := windowsipam.GetInit(windowsipam.DefaultIPAM)

	err := InitDockerDefault(ic, l, g)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return initFunc(ic, l, g)
***REMOVED***
