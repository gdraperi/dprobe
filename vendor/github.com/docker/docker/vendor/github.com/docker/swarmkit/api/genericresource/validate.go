package genericresource

import (
	"fmt"
	"github.com/docker/swarmkit/api"
)

// ValidateTask validates that the task only uses integers
// for generic resources
func ValidateTask(resources *api.Resources) error ***REMOVED***
	for _, v := range resources.Generic ***REMOVED***
		if v.GetDiscreteResourceSpec() != nil ***REMOVED***
			continue
		***REMOVED***

		return fmt.Errorf("invalid argument for resource %s", Kind(v))
	***REMOVED***

	return nil
***REMOVED***

// HasEnough returns true if node can satisfy the task's GenericResource request
func HasEnough(nodeRes []*api.GenericResource, taskRes *api.GenericResource) (bool, error) ***REMOVED***
	t := taskRes.GetDiscreteResourceSpec()
	if t == nil ***REMOVED***
		return false, fmt.Errorf("task should only hold Discrete type")
	***REMOVED***

	if nodeRes == nil ***REMOVED***
		return false, nil
	***REMOVED***

	nrs := GetResource(t.Kind, nodeRes)
	if len(nrs) == 0 ***REMOVED***
		return false, nil
	***REMOVED***

	switch nr := nrs[0].Resource.(type) ***REMOVED***
	case *api.GenericResource_DiscreteResourceSpec:
		if t.Value > nr.DiscreteResourceSpec.Value ***REMOVED***
			return false, nil
		***REMOVED***
	case *api.GenericResource_NamedResourceSpec:
		if t.Value > int64(len(nrs)) ***REMOVED***
			return false, nil
		***REMOVED***
	***REMOVED***

	return true, nil
***REMOVED***

// HasResource checks if there is enough "res" in the "resources" argument
func HasResource(res *api.GenericResource, resources []*api.GenericResource) bool ***REMOVED***
	for _, r := range resources ***REMOVED***
		if Kind(res) != Kind(r) ***REMOVED***
			continue
		***REMOVED***

		switch rtype := r.Resource.(type) ***REMOVED***
		case *api.GenericResource_DiscreteResourceSpec:
			if res.GetDiscreteResourceSpec() == nil ***REMOVED***
				return false
			***REMOVED***

			if res.GetDiscreteResourceSpec().Value < rtype.DiscreteResourceSpec.Value ***REMOVED***
				return false
			***REMOVED***

			return true
		case *api.GenericResource_NamedResourceSpec:
			if res.GetNamedResourceSpec() == nil ***REMOVED***
				return false
			***REMOVED***

			if res.GetNamedResourceSpec().Value != rtype.NamedResourceSpec.Value ***REMOVED***
				continue
			***REMOVED***

			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***
