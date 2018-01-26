package bridge

import "github.com/docker/libnetwork/iptables"

func (n *bridgeNetwork) setupFirewalld(config *networkConfiguration, i *bridgeInterface) error ***REMOVED***
	d := n.driver
	d.Lock()
	driverConfig := d.config
	d.Unlock()

	// Sanity check.
	if !driverConfig.EnableIPTables ***REMOVED***
		return IPTableCfgError(config.BridgeName)
	***REMOVED***

	iptables.OnReloaded(func() ***REMOVED*** n.setupIPTables(config, i) ***REMOVED***)
	iptables.OnReloaded(n.portMapper.ReMapAll)

	return nil
***REMOVED***
