// Package v1p20 provides specific API types for the API version 1, patch 20.
package v1p20

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

// ContainerJSON is a backcompatibility struct for the API 1.20
type ContainerJSON struct ***REMOVED***
	*types.ContainerJSONBase
	Mounts          []types.MountPoint
	Config          *ContainerConfig
	NetworkSettings *NetworkSettings
***REMOVED***

// ContainerConfig is a backcompatibility struct used in ContainerJSON for the API 1.20
type ContainerConfig struct ***REMOVED***
	*container.Config

	MacAddress      string
	NetworkDisabled bool
	ExposedPorts    map[nat.Port]struct***REMOVED******REMOVED***

	// backward compatibility, they now live in HostConfig
	VolumeDriver string
***REMOVED***

// StatsJSON is a backcompatibility struct used in Stats for APIs prior to 1.21
type StatsJSON struct ***REMOVED***
	types.Stats
	Network types.NetworkStats `json:"network,omitempty"`
***REMOVED***

// NetworkSettings is a backward compatible struct for APIs prior to 1.21
type NetworkSettings struct ***REMOVED***
	types.NetworkSettingsBase
	types.DefaultNetworkSettings
***REMOVED***
