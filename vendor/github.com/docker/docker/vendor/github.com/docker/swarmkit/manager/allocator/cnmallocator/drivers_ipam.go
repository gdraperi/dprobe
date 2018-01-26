package cnmallocator

import (
	"github.com/docker/libnetwork/drvregistry"
	"github.com/docker/libnetwork/ipamapi"
	builtinIpam "github.com/docker/libnetwork/ipams/builtin"
	nullIpam "github.com/docker/libnetwork/ipams/null"
	remoteIpam "github.com/docker/libnetwork/ipams/remote"
)

func initIPAMDrivers(r *drvregistry.DrvRegistry) error ***REMOVED***
	for _, fn := range [](func(ipamapi.Callback, interface***REMOVED******REMOVED***, interface***REMOVED******REMOVED***) error)***REMOVED***
		builtinIpam.Init,
		remoteIpam.Init,
		nullIpam.Init,
	***REMOVED*** ***REMOVED***
		if err := fn(r, nil, nil); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***
