// +build !windows

package runconfig

import (
	"github.com/docker/docker/api/types/container"
	networktypes "github.com/docker/docker/api/types/network"
)

// ContainerConfigWrapper is a Config wrapper that holds the container Config (portable)
// and the corresponding HostConfig (non-portable).
type ContainerConfigWrapper struct ***REMOVED***
	*container.Config
	InnerHostConfig       *container.HostConfig          `json:"HostConfig,omitempty"`
	Cpuset                string                         `json:",omitempty"` // Deprecated. Exported for backwards compatibility.
	NetworkingConfig      *networktypes.NetworkingConfig `json:"NetworkingConfig,omitempty"`
	*container.HostConfig                                // Deprecated. Exported to read attributes from json that are not in the inner host config structure.
***REMOVED***

// getHostConfig gets the HostConfig of the Config.
// It's mostly there to handle Deprecated fields of the ContainerConfigWrapper
func (w *ContainerConfigWrapper) getHostConfig() *container.HostConfig ***REMOVED***
	hc := w.HostConfig

	if hc == nil && w.InnerHostConfig != nil ***REMOVED***
		hc = w.InnerHostConfig
	***REMOVED*** else if w.InnerHostConfig != nil ***REMOVED***
		if hc.Memory != 0 && w.InnerHostConfig.Memory == 0 ***REMOVED***
			w.InnerHostConfig.Memory = hc.Memory
		***REMOVED***
		if hc.MemorySwap != 0 && w.InnerHostConfig.MemorySwap == 0 ***REMOVED***
			w.InnerHostConfig.MemorySwap = hc.MemorySwap
		***REMOVED***
		if hc.CPUShares != 0 && w.InnerHostConfig.CPUShares == 0 ***REMOVED***
			w.InnerHostConfig.CPUShares = hc.CPUShares
		***REMOVED***
		if hc.CpusetCpus != "" && w.InnerHostConfig.CpusetCpus == "" ***REMOVED***
			w.InnerHostConfig.CpusetCpus = hc.CpusetCpus
		***REMOVED***

		if hc.VolumeDriver != "" && w.InnerHostConfig.VolumeDriver == "" ***REMOVED***
			w.InnerHostConfig.VolumeDriver = hc.VolumeDriver
		***REMOVED***

		hc = w.InnerHostConfig
	***REMOVED***

	if hc != nil ***REMOVED***
		if w.Cpuset != "" && hc.CpusetCpus == "" ***REMOVED***
			hc.CpusetCpus = w.Cpuset
		***REMOVED***
	***REMOVED***

	// Make sure NetworkMode has an acceptable value. We do this to ensure
	// backwards compatible API behavior.
	SetDefaultNetModeIfBlank(hc)

	return hc
***REMOVED***
