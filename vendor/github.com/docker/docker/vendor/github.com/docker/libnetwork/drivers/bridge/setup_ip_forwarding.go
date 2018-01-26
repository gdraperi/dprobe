package bridge

import (
	"fmt"
	"io/ioutil"

	"github.com/docker/libnetwork/iptables"
	"github.com/sirupsen/logrus"
)

const (
	ipv4ForwardConf     = "/proc/sys/net/ipv4/ip_forward"
	ipv4ForwardConfPerm = 0644
)

func configureIPForwarding(enable bool) error ***REMOVED***
	var val byte
	if enable ***REMOVED***
		val = '1'
	***REMOVED***
	return ioutil.WriteFile(ipv4ForwardConf, []byte***REMOVED***val, '\n'***REMOVED***, ipv4ForwardConfPerm)
***REMOVED***

func setupIPForwarding(enableIPTables bool) error ***REMOVED***
	// Get current IPv4 forward setup
	ipv4ForwardData, err := ioutil.ReadFile(ipv4ForwardConf)
	if err != nil ***REMOVED***
		return fmt.Errorf("Cannot read IP forwarding setup: %v", err)
	***REMOVED***

	// Enable IPv4 forwarding only if it is not already enabled
	if ipv4ForwardData[0] != '1' ***REMOVED***
		// Enable IPv4 forwarding
		if err := configureIPForwarding(true); err != nil ***REMOVED***
			return fmt.Errorf("Enabling IP forwarding failed: %v", err)
		***REMOVED***
		// When enabling ip_forward set the default policy on forward chain to
		// drop only if the daemon option iptables is not set to false.
		if !enableIPTables ***REMOVED***
			return nil
		***REMOVED***
		if err := iptables.SetDefaultPolicy(iptables.Filter, "FORWARD", iptables.Drop); err != nil ***REMOVED***
			if err := configureIPForwarding(false); err != nil ***REMOVED***
				logrus.Errorf("Disabling IP forwarding failed, %v", err)
			***REMOVED***
			return err
		***REMOVED***
		iptables.OnReloaded(func() ***REMOVED***
			logrus.Debug("Setting the default DROP policy on firewall reload")
			if err := iptables.SetDefaultPolicy(iptables.Filter, "FORWARD", iptables.Drop); err != nil ***REMOVED***
				logrus.Warnf("Settig the default DROP policy on firewall reload failed, %v", err)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
	return nil
***REMOVED***
