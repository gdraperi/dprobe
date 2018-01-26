package runconfig

import (
	"encoding/json"
	"io"

	"github.com/docker/docker/api/types/container"
	networktypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/sysinfo"
)

// ContainerDecoder implements httputils.ContainerDecoder
// calling DecodeContainerConfig.
type ContainerDecoder struct***REMOVED******REMOVED***

// DecodeConfig makes ContainerDecoder to implement httputils.ContainerDecoder
func (r ContainerDecoder) DecodeConfig(src io.Reader) (*container.Config, *container.HostConfig, *networktypes.NetworkingConfig, error) ***REMOVED***
	return decodeContainerConfig(src)
***REMOVED***

// DecodeHostConfig makes ContainerDecoder to implement httputils.ContainerDecoder
func (r ContainerDecoder) DecodeHostConfig(src io.Reader) (*container.HostConfig, error) ***REMOVED***
	return decodeHostConfig(src)
***REMOVED***

// decodeContainerConfig decodes a json encoded config into a ContainerConfigWrapper
// struct and returns both a Config and a HostConfig struct
// Be aware this function is not checking whether the resulted structs are nil,
// it's your business to do so
func decodeContainerConfig(src io.Reader) (*container.Config, *container.HostConfig, *networktypes.NetworkingConfig, error) ***REMOVED***
	var w ContainerConfigWrapper

	decoder := json.NewDecoder(src)
	if err := decoder.Decode(&w); err != nil ***REMOVED***
		return nil, nil, nil, err
	***REMOVED***

	hc := w.getHostConfig()

	// Perform platform-specific processing of Volumes and Binds.
	if w.Config != nil && hc != nil ***REMOVED***

		// Initialize the volumes map if currently nil
		if w.Config.Volumes == nil ***REMOVED***
			w.Config.Volumes = make(map[string]struct***REMOVED******REMOVED***)
		***REMOVED***
	***REMOVED***

	// Certain parameters need daemon-side validation that cannot be done
	// on the client, as only the daemon knows what is valid for the platform.
	if err := validateNetMode(w.Config, hc); err != nil ***REMOVED***
		return nil, nil, nil, err
	***REMOVED***

	// Validate isolation
	if err := validateIsolation(hc); err != nil ***REMOVED***
		return nil, nil, nil, err
	***REMOVED***

	// Validate QoS
	if err := validateQoS(hc); err != nil ***REMOVED***
		return nil, nil, nil, err
	***REMOVED***

	// Validate Resources
	if err := validateResources(hc, sysinfo.New(true)); err != nil ***REMOVED***
		return nil, nil, nil, err
	***REMOVED***

	// Validate Privileged
	if err := validatePrivileged(hc); err != nil ***REMOVED***
		return nil, nil, nil, err
	***REMOVED***

	// Validate ReadonlyRootfs
	if err := validateReadonlyRootfs(hc); err != nil ***REMOVED***
		return nil, nil, nil, err
	***REMOVED***

	return w.Config, hc, w.NetworkingConfig, nil
***REMOVED***
