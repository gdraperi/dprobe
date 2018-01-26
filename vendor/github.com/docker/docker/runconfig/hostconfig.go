package runconfig

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/docker/docker/api/types/container"
)

// DecodeHostConfig creates a HostConfig based on the specified Reader.
// It assumes the content of the reader will be JSON, and decodes it.
func decodeHostConfig(src io.Reader) (*container.HostConfig, error) ***REMOVED***
	decoder := json.NewDecoder(src)

	var w ContainerConfigWrapper
	if err := decoder.Decode(&w); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	hc := w.getHostConfig()
	return hc, nil
***REMOVED***

// SetDefaultNetModeIfBlank changes the NetworkMode in a HostConfig structure
// to default if it is not populated. This ensures backwards compatibility after
// the validation of the network mode was moved from the docker CLI to the
// docker daemon.
func SetDefaultNetModeIfBlank(hc *container.HostConfig) ***REMOVED***
	if hc != nil ***REMOVED***
		if hc.NetworkMode == container.NetworkMode("") ***REMOVED***
			hc.NetworkMode = container.NetworkMode("default")
		***REMOVED***
	***REMOVED***
***REMOVED***

// validateNetContainerMode ensures that the various combinations of requested
// network settings wrt container mode are valid.
func validateNetContainerMode(c *container.Config, hc *container.HostConfig) error ***REMOVED***
	// We may not be passed a host config, such as in the case of docker commit
	if hc == nil ***REMOVED***
		return nil
	***REMOVED***
	parts := strings.Split(string(hc.NetworkMode), ":")
	if parts[0] == "container" ***REMOVED***
		if len(parts) < 2 || parts[1] == "" ***REMOVED***
			return validationError("Invalid network mode: invalid container format container:<name|id>")
		***REMOVED***
	***REMOVED***

	if hc.NetworkMode.IsContainer() && c.Hostname != "" ***REMOVED***
		return ErrConflictNetworkHostname
	***REMOVED***

	if hc.NetworkMode.IsContainer() && len(hc.Links) > 0 ***REMOVED***
		return ErrConflictContainerNetworkAndLinks
	***REMOVED***

	if hc.NetworkMode.IsContainer() && len(hc.DNS) > 0 ***REMOVED***
		return ErrConflictNetworkAndDNS
	***REMOVED***

	if hc.NetworkMode.IsContainer() && len(hc.ExtraHosts) > 0 ***REMOVED***
		return ErrConflictNetworkHosts
	***REMOVED***

	if (hc.NetworkMode.IsContainer() || hc.NetworkMode.IsHost()) && c.MacAddress != "" ***REMOVED***
		return ErrConflictContainerNetworkAndMac
	***REMOVED***

	if hc.NetworkMode.IsContainer() && (len(hc.PortBindings) > 0 || hc.PublishAllPorts) ***REMOVED***
		return ErrConflictNetworkPublishPorts
	***REMOVED***

	if hc.NetworkMode.IsContainer() && len(c.ExposedPorts) > 0 ***REMOVED***
		return ErrConflictNetworkExposePorts
	***REMOVED***
	return nil
***REMOVED***
