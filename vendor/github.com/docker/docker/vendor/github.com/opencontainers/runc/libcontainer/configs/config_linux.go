package configs

import "fmt"

// HostUID gets the translated uid for the process on host which could be
// different when user namespaces are enabled.
func (c Config) HostUID(containerId int) (int, error) ***REMOVED***
	if c.Namespaces.Contains(NEWUSER) ***REMOVED***
		if c.UidMappings == nil ***REMOVED***
			return -1, fmt.Errorf("User namespaces enabled, but no uid mappings found.")
		***REMOVED***
		id, found := c.hostIDFromMapping(containerId, c.UidMappings)
		if !found ***REMOVED***
			return -1, fmt.Errorf("User namespaces enabled, but no user mapping found.")
		***REMOVED***
		return id, nil
	***REMOVED***
	// Return unchanged id.
	return containerId, nil
***REMOVED***

// HostRootUID gets the root uid for the process on host which could be non-zero
// when user namespaces are enabled.
func (c Config) HostRootUID() (int, error) ***REMOVED***
	return c.HostUID(0)
***REMOVED***

// HostGID gets the translated gid for the process on host which could be
// different when user namespaces are enabled.
func (c Config) HostGID(containerId int) (int, error) ***REMOVED***
	if c.Namespaces.Contains(NEWUSER) ***REMOVED***
		if c.GidMappings == nil ***REMOVED***
			return -1, fmt.Errorf("User namespaces enabled, but no gid mappings found.")
		***REMOVED***
		id, found := c.hostIDFromMapping(containerId, c.GidMappings)
		if !found ***REMOVED***
			return -1, fmt.Errorf("User namespaces enabled, but no group mapping found.")
		***REMOVED***
		return id, nil
	***REMOVED***
	// Return unchanged id.
	return containerId, nil
***REMOVED***

// HostRootGID gets the root gid for the process on host which could be non-zero
// when user namespaces are enabled.
func (c Config) HostRootGID() (int, error) ***REMOVED***
	return c.HostGID(0)
***REMOVED***

// Utility function that gets a host ID for a container ID from user namespace map
// if that ID is present in the map.
func (c Config) hostIDFromMapping(containerID int, uMap []IDMap) (int, bool) ***REMOVED***
	for _, m := range uMap ***REMOVED***
		if (containerID >= m.ContainerID) && (containerID <= (m.ContainerID + m.Size - 1)) ***REMOVED***
			hostID := m.HostID + (containerID - m.ContainerID)
			return hostID, true
		***REMOVED***
	***REMOVED***
	return -1, false
***REMOVED***
