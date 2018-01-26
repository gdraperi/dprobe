package genericresource

import (
	"github.com/docker/swarmkit/api"
)

// NewSet creates a set object
func NewSet(key string, vals ...string) []*api.GenericResource ***REMOVED***
	rs := make([]*api.GenericResource, 0, len(vals))

	for _, v := range vals ***REMOVED***
		rs = append(rs, NewString(key, v))
	***REMOVED***

	return rs
***REMOVED***

// NewString creates a String resource
func NewString(key, val string) *api.GenericResource ***REMOVED***
	return &api.GenericResource***REMOVED***
		Resource: &api.GenericResource_NamedResourceSpec***REMOVED***
			NamedResourceSpec: &api.NamedGenericResource***REMOVED***
				Kind:  key,
				Value: val,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// NewDiscrete creates a Discrete resource
func NewDiscrete(key string, val int64) *api.GenericResource ***REMOVED***
	return &api.GenericResource***REMOVED***
		Resource: &api.GenericResource_DiscreteResourceSpec***REMOVED***
			DiscreteResourceSpec: &api.DiscreteGenericResource***REMOVED***
				Kind:  key,
				Value: val,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// GetResource returns resources from the "resources" parameter matching the kind key
func GetResource(kind string, resources []*api.GenericResource) []*api.GenericResource ***REMOVED***
	var res []*api.GenericResource

	for _, r := range resources ***REMOVED***
		if Kind(r) != kind ***REMOVED***
			continue
		***REMOVED***

		res = append(res, r)
	***REMOVED***

	return res
***REMOVED***

// ConsumeNodeResources removes "res" from nodeAvailableResources
func ConsumeNodeResources(nodeAvailableResources *[]*api.GenericResource, res []*api.GenericResource) ***REMOVED***
	if nodeAvailableResources == nil ***REMOVED***
		return
	***REMOVED***

	w := 0

loop:
	for _, na := range *nodeAvailableResources ***REMOVED***
		for _, r := range res ***REMOVED***
			if Kind(na) != Kind(r) ***REMOVED***
				continue
			***REMOVED***

			if remove(na, r) ***REMOVED***
				continue loop
			***REMOVED***
			// If this wasn't the right element then
			// we need to continue
		***REMOVED***

		(*nodeAvailableResources)[w] = na
		w++
	***REMOVED***

	*nodeAvailableResources = (*nodeAvailableResources)[:w]
***REMOVED***

// Returns true if the element is to be removed from the list
func remove(na, r *api.GenericResource) bool ***REMOVED***
	switch tr := r.Resource.(type) ***REMOVED***
	case *api.GenericResource_DiscreteResourceSpec:
		if na.GetDiscreteResourceSpec() == nil ***REMOVED***
			return false // Type change, ignore
		***REMOVED***

		na.GetDiscreteResourceSpec().Value -= tr.DiscreteResourceSpec.Value
		if na.GetDiscreteResourceSpec().Value <= 0 ***REMOVED***
			return true
		***REMOVED***
	case *api.GenericResource_NamedResourceSpec:
		if na.GetNamedResourceSpec() == nil ***REMOVED***
			return false // Type change, ignore
		***REMOVED***

		if tr.NamedResourceSpec.Value != na.GetNamedResourceSpec().Value ***REMOVED***
			return false // not the right item, ignore
		***REMOVED***

		return true
	***REMOVED***

	return false
***REMOVED***
