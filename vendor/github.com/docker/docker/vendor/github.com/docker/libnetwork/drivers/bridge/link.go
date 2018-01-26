package bridge

import (
	"fmt"
	"net"

	"github.com/docker/libnetwork/iptables"
	"github.com/docker/libnetwork/types"
	"github.com/sirupsen/logrus"
)

type link struct ***REMOVED***
	parentIP string
	childIP  string
	ports    []types.TransportPort
	bridge   string
***REMOVED***

func (l *link) String() string ***REMOVED***
	return fmt.Sprintf("%s <-> %s [%v] on %s", l.parentIP, l.childIP, l.ports, l.bridge)
***REMOVED***

func newLink(parentIP, childIP string, ports []types.TransportPort, bridge string) *link ***REMOVED***
	return &link***REMOVED***
		childIP:  childIP,
		parentIP: parentIP,
		ports:    ports,
		bridge:   bridge,
	***REMOVED***

***REMOVED***

func (l *link) Enable() error ***REMOVED***
	// -A == iptables append flag
	linkFunction := func() error ***REMOVED***
		return linkContainers("-A", l.parentIP, l.childIP, l.ports, l.bridge, false)
	***REMOVED***

	iptables.OnReloaded(func() ***REMOVED*** linkFunction() ***REMOVED***)
	return linkFunction()
***REMOVED***

func (l *link) Disable() ***REMOVED***
	// -D == iptables delete flag
	err := linkContainers("-D", l.parentIP, l.childIP, l.ports, l.bridge, true)
	if err != nil ***REMOVED***
		logrus.Errorf("Error removing IPTables rules for a link %s due to %s", l.String(), err.Error())
	***REMOVED***
	// Return proper error once we move to use a proper iptables package
	// that returns typed errors
***REMOVED***

func linkContainers(action, parentIP, childIP string, ports []types.TransportPort, bridge string,
	ignoreErrors bool) error ***REMOVED***
	var nfAction iptables.Action

	switch action ***REMOVED***
	case "-A":
		nfAction = iptables.Append
	case "-I":
		nfAction = iptables.Insert
	case "-D":
		nfAction = iptables.Delete
	default:
		return InvalidIPTablesCfgError(action)
	***REMOVED***

	ip1 := net.ParseIP(parentIP)
	if ip1 == nil ***REMOVED***
		return InvalidLinkIPAddrError(parentIP)
	***REMOVED***
	ip2 := net.ParseIP(childIP)
	if ip2 == nil ***REMOVED***
		return InvalidLinkIPAddrError(childIP)
	***REMOVED***

	chain := iptables.ChainInfo***REMOVED***Name: DockerChain***REMOVED***
	for _, port := range ports ***REMOVED***
		err := chain.Link(nfAction, ip1, ip2, int(port.Port), port.Proto.String(), bridge)
		if !ignoreErrors && err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
