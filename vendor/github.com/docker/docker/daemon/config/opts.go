package config

import (
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/daemon/cluster/convert"
	"github.com/docker/swarmkit/api/genericresource"
)

// ParseGenericResources parses and validates the specified string as a list of GenericResource
func ParseGenericResources(value []string) ([]swarm.GenericResource, error) ***REMOVED***
	if len(value) == 0 ***REMOVED***
		return nil, nil
	***REMOVED***

	resources, err := genericresource.Parse(value)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	obj := convert.GenericResourcesFromGRPC(resources)
	return obj, nil
***REMOVED***
