// +build !windows

package runconfig

import (
	"fmt"
	"runtime"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/sysinfo"
)

// DefaultDaemonNetworkMode returns the default network stack the daemon should
// use.
func DefaultDaemonNetworkMode() container.NetworkMode ***REMOVED***
	return container.NetworkMode("bridge")
***REMOVED***

// IsPreDefinedNetwork indicates if a network is predefined by the daemon
func IsPreDefinedNetwork(network string) bool ***REMOVED***
	n := container.NetworkMode(network)
	return n.IsBridge() || n.IsHost() || n.IsNone() || n.IsDefault()
***REMOVED***

// validateNetMode ensures that the various combinations of requested
// network settings are valid.
func validateNetMode(c *container.Config, hc *container.HostConfig) error ***REMOVED***
	// We may not be passed a host config, such as in the case of docker commit
	if hc == nil ***REMOVED***
		return nil
	***REMOVED***

	err := validateNetContainerMode(c, hc)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if hc.UTSMode.IsHost() && c.Hostname != "" ***REMOVED***
		return ErrConflictUTSHostname
	***REMOVED***

	if hc.NetworkMode.IsHost() && len(hc.Links) > 0 ***REMOVED***
		return ErrConflictHostNetworkAndLinks
	***REMOVED***

	return nil
***REMOVED***

// validateIsolation performs platform specific validation of
// isolation in the hostconfig structure. Linux only supports "default"
// which is LXC container isolation
func validateIsolation(hc *container.HostConfig) error ***REMOVED***
	// We may not be passed a host config, such as in the case of docker commit
	if hc == nil ***REMOVED***
		return nil
	***REMOVED***
	if !hc.Isolation.IsValid() ***REMOVED***
		return fmt.Errorf("Invalid isolation: %q - %s only supports 'default'", hc.Isolation, runtime.GOOS)
	***REMOVED***
	return nil
***REMOVED***

// validateQoS performs platform specific validation of the QoS settings
func validateQoS(hc *container.HostConfig) error ***REMOVED***
	// We may not be passed a host config, such as in the case of docker commit
	if hc == nil ***REMOVED***
		return nil
	***REMOVED***

	if hc.IOMaximumBandwidth != 0 ***REMOVED***
		return fmt.Errorf("Invalid QoS settings: %s does not support configuration of maximum bandwidth", runtime.GOOS)
	***REMOVED***

	if hc.IOMaximumIOps != 0 ***REMOVED***
		return fmt.Errorf("Invalid QoS settings: %s does not support configuration of maximum IOPs", runtime.GOOS)
	***REMOVED***
	return nil
***REMOVED***

// validateResources performs platform specific validation of the resource settings
// cpu-rt-runtime and cpu-rt-period can not be greater than their parent, cpu-rt-runtime requires sys_nice
func validateResources(hc *container.HostConfig, si *sysinfo.SysInfo) error ***REMOVED***
	// We may not be passed a host config, such as in the case of docker commit
	if hc == nil ***REMOVED***
		return nil
	***REMOVED***

	if hc.Resources.CPURealtimePeriod > 0 && !si.CPURealtimePeriod ***REMOVED***
		return fmt.Errorf("Your kernel does not support cgroup cpu real-time period")
	***REMOVED***

	if hc.Resources.CPURealtimeRuntime > 0 && !si.CPURealtimeRuntime ***REMOVED***
		return fmt.Errorf("Your kernel does not support cgroup cpu real-time runtime")
	***REMOVED***

	if hc.Resources.CPURealtimePeriod != 0 && hc.Resources.CPURealtimeRuntime != 0 && hc.Resources.CPURealtimeRuntime > hc.Resources.CPURealtimePeriod ***REMOVED***
		return fmt.Errorf("cpu real-time runtime cannot be higher than cpu real-time period")
	***REMOVED***
	return nil
***REMOVED***

// validatePrivileged performs platform specific validation of the Privileged setting
func validatePrivileged(hc *container.HostConfig) error ***REMOVED***
	return nil
***REMOVED***

// validateReadonlyRootfs performs platform specific validation of the ReadonlyRootfs setting
func validateReadonlyRootfs(hc *container.HostConfig) error ***REMOVED***
	return nil
***REMOVED***
