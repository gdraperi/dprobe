package daemon

import (
	apitypes "github.com/docker/docker/api/types"
	lncluster "github.com/docker/libnetwork/cluster"
)

// Cluster is the interface for github.com/docker/docker/daemon/cluster.(*Cluster).
type Cluster interface ***REMOVED***
	ClusterStatus
	NetworkManager
	SendClusterEvent(event lncluster.ConfigEventType)
***REMOVED***

// ClusterStatus interface provides information about the Swarm status of the Cluster
type ClusterStatus interface ***REMOVED***
	IsAgent() bool
	IsManager() bool
***REMOVED***

// NetworkManager provides methods to manage networks
type NetworkManager interface ***REMOVED***
	GetNetwork(input string) (apitypes.NetworkResource, error)
	GetNetworks() ([]apitypes.NetworkResource, error)
	RemoveNetwork(input string) error
***REMOVED***
