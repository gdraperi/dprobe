// +build !windows

package daemon

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/container"
	"github.com/pkg/errors"
)

// Resolve Network SandboxID in case the container reuse another container's network stack
func (daemon *Daemon) getNetworkSandboxID(c *container.Container) (string, error) ***REMOVED***
	curr := c
	for curr.HostConfig.NetworkMode.IsContainer() ***REMOVED***
		containerID := curr.HostConfig.NetworkMode.ConnectedContainer()
		connected, err := daemon.GetContainer(containerID)
		if err != nil ***REMOVED***
			return "", errors.Wrapf(err, "Could not get container for %s", containerID)
		***REMOVED***
		curr = connected
	***REMOVED***
	return curr.NetworkSettings.SandboxID, nil
***REMOVED***

func (daemon *Daemon) getNetworkStats(c *container.Container) (map[string]types.NetworkStats, error) ***REMOVED***
	sandboxID, err := daemon.getNetworkSandboxID(c)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	sb, err := daemon.netController.SandboxByID(sandboxID)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	lnstats, err := sb.Statistics()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	stats := make(map[string]types.NetworkStats)
	// Convert libnetwork nw stats into api stats
	for ifName, ifStats := range lnstats ***REMOVED***
		stats[ifName] = types.NetworkStats***REMOVED***
			RxBytes:   ifStats.RxBytes,
			RxPackets: ifStats.RxPackets,
			RxErrors:  ifStats.RxErrors,
			RxDropped: ifStats.RxDropped,
			TxBytes:   ifStats.TxBytes,
			TxPackets: ifStats.TxPackets,
			TxErrors:  ifStats.TxErrors,
			TxDropped: ifStats.TxDropped,
		***REMOVED***
	***REMOVED***

	return stats, nil
***REMOVED***
