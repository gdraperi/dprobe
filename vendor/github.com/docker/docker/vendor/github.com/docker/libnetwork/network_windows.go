// +build windows

package libnetwork

import (
	"runtime"
	"time"

	"github.com/Microsoft/hcsshim"
	"github.com/docker/libnetwork/drivers/windows"
	"github.com/docker/libnetwork/ipamapi"
	"github.com/docker/libnetwork/ipams/windowsipam"
	"github.com/sirupsen/logrus"
)

func executeInCompartment(compartmentID uint32, x func()) ***REMOVED***
	runtime.LockOSThread()

	if err := hcsshim.SetCurrentThreadCompartmentId(compartmentID); err != nil ***REMOVED***
		logrus.Error(err)
	***REMOVED***
	defer func() ***REMOVED***
		hcsshim.SetCurrentThreadCompartmentId(0)
		runtime.UnlockOSThread()
	***REMOVED***()

	x()
***REMOVED***

func (n *network) startResolver() ***REMOVED***
	if n.networkType == "ics" ***REMOVED***
		return
	***REMOVED***
	n.resolverOnce.Do(func() ***REMOVED***
		logrus.Debugf("Launching DNS server for network %q", n.Name())
		options := n.Info().DriverOptions()
		hnsid := options[windows.HNSID]

		if hnsid == "" ***REMOVED***
			return
		***REMOVED***

		hnsresponse, err := hcsshim.HNSNetworkRequest("GET", hnsid, "")
		if err != nil ***REMOVED***
			logrus.Errorf("Resolver Setup/Start failed for container %s, %q", n.Name(), err)
			return
		***REMOVED***

		for _, subnet := range hnsresponse.Subnets ***REMOVED***
			if subnet.GatewayAddress != "" ***REMOVED***
				for i := 0; i < 3; i++ ***REMOVED***
					resolver := NewResolver(subnet.GatewayAddress, false, "", n)
					logrus.Debugf("Binding a resolver on network %s gateway %s", n.Name(), subnet.GatewayAddress)
					executeInCompartment(hnsresponse.DNSServerCompartment, resolver.SetupFunc(53))

					if err = resolver.Start(); err != nil ***REMOVED***
						logrus.Errorf("Resolver Setup/Start failed for container %s, %q", n.Name(), err)
						time.Sleep(1 * time.Second)
					***REMOVED*** else ***REMOVED***
						logrus.Debugf("Resolver bound successfully for network %s", n.Name())
						n.resolver = append(n.resolver, resolver)
						break
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***

func defaultIpamForNetworkType(networkType string) string ***REMOVED***
	if windows.IsBuiltinLocalDriver(networkType) ***REMOVED***
		return windowsipam.DefaultIPAM
	***REMOVED***
	return ipamapi.DefaultIPAM
***REMOVED***
