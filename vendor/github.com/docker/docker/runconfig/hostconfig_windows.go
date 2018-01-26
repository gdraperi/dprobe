package runconfig

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/sysinfo"
)

// DefaultDaemonNetworkMode returns the default network stack the daemon should
// use.
func DefaultDaemonNetworkMode() container.NetworkMode ***REMOVED***
	return container.NetworkMode("nat")
***REMOVED***

// IsPreDefinedNetwork indicates if a network is predefined by the daemon
func IsPreDefinedNetwork(network string) bool ***REMOVED***
	return !container.NetworkMode(network).IsUserDefined()
***REMOVED***

// validateNetMode ensures that the various combinations of requested
// network settings are valid.
func validateNetMode(c *container.Config, hc *container.HostConfig) error ***REMOVED***
	if hc == nil ***REMOVED***
		return nil
	***REMOVED***

	err := validateNetContainerMode(c, hc)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if hc.NetworkMode.IsContainer() && hc.Isolation.IsHyperV() ***REMOVED***
		return fmt.Errorf("Using the network stack of another container is not supported while using Hyper-V Containers")
	***REMOVED***

	return nil
***REMOVED***

// validateIsolation performs platform specific validation of the
// isolation in the hostconfig structure. Windows supports 'default' (or
// blank), 'process', or 'hyperv'.
func validateIsolation(hc *container.HostConfig) error ***REMOVED***
	// We may not be passed a host config, such as in the case of docker commit
	if hc == nil ***REMOVED***
		return nil
	***REMOVED***
	if !hc.Isolation.IsValid() ***REMOVED***
		return fmt.Errorf("Invalid isolation: %q. Windows supports 'default', 'process', or 'hyperv'", hc.Isolation)
	***REMOVED***
	return nil
***REMOVED***

// validateQoS performs platform specific validation of the Qos settings
func validateQoS(hc *container.HostConfig) error ***REMOVED***
	return nil
***REMOVED***

// validateResources performs platform specific validation of the resource settings
func validateResources(hc *container.HostConfig, si *sysinfo.SysInfo) error ***REMOVED***
	// We may not be passed a host config, such as in the case of docker commit
	if hc == nil ***REMOVED***
		return nil
	***REMOVED***
	if hc.Resources.CPURealtimePeriod != 0 ***REMOVED***
		return fmt.Errorf("Windows does not support CPU real-time period")
	***REMOVED***
	if hc.Resources.CPURealtimeRuntime != 0 ***REMOVED***
		return fmt.Errorf("Windows does not support CPU real-time runtime")
	***REMOVED***
	return nil
***REMOVED***

// validatePrivileged performs platform specific validation of the Privileged setting
func validatePrivileged(hc *container.HostConfig) error ***REMOVED***
	// We may not be passed a host config, such as in the case of docker commit
	if hc == nil ***REMOVED***
		return nil
	***REMOVED***
	if hc.Privileged ***REMOVED***
		return fmt.Errorf("Windows does not support privileged mode")
	***REMOVED***
	return nil
***REMOVED***

// validateReadonlyRootfs performs platform specific validation of the ReadonlyRootfs setting
func validateReadonlyRootfs(hc *container.HostConfig) error ***REMOVED***
	// We may not be passed a host config, such as in the case of docker commit
	if hc == nil ***REMOVED***
		return nil
	***REMOVED***
	if hc.ReadonlyRootfs ***REMOVED***
		return fmt.Errorf("Windows does not support root filesystem in read-only mode")
	***REMOVED***
	return nil
***REMOVED***
