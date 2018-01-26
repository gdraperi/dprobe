package genericresource

import (
	"fmt"
	"github.com/docker/swarmkit/api"
)

// Claim assigns GenericResources to a task by taking them from the
// node's GenericResource list and storing them in the task's available list
func Claim(nodeAvailableResources, taskAssigned *[]*api.GenericResource,
	taskReservations []*api.GenericResource) error ***REMOVED***
	var resSelected []*api.GenericResource

	for _, res := range taskReservations ***REMOVED***
		tr := res.GetDiscreteResourceSpec()
		if tr == nil ***REMOVED***
			return fmt.Errorf("task should only hold Discrete type")
		***REMOVED***

		// Select the resources
		nrs, err := selectNodeResources(*nodeAvailableResources, tr)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		resSelected = append(resSelected, nrs...)
	***REMOVED***

	ClaimResources(nodeAvailableResources, taskAssigned, resSelected)
	return nil
***REMOVED***

// ClaimResources adds the specified resources to the task's list
// and removes them from the node's generic resource list
func ClaimResources(nodeAvailableResources, taskAssigned *[]*api.GenericResource,
	resSelected []*api.GenericResource) ***REMOVED***
	*taskAssigned = append(*taskAssigned, resSelected...)
	ConsumeNodeResources(nodeAvailableResources, resSelected)
***REMOVED***

func selectNodeResources(nodeRes []*api.GenericResource,
	tr *api.DiscreteGenericResource) ([]*api.GenericResource, error) ***REMOVED***
	var nrs []*api.GenericResource

	for _, res := range nodeRes ***REMOVED***
		if Kind(res) != tr.Kind ***REMOVED***
			continue
		***REMOVED***

		switch nr := res.Resource.(type) ***REMOVED***
		case *api.GenericResource_DiscreteResourceSpec:
			if nr.DiscreteResourceSpec.Value >= tr.Value && tr.Value != 0 ***REMOVED***
				nrs = append(nrs, NewDiscrete(tr.Kind, tr.Value))
			***REMOVED***

			return nrs, nil
		case *api.GenericResource_NamedResourceSpec:
			nrs = append(nrs, res.Copy())

			if int64(len(nrs)) == tr.Value ***REMOVED***
				return nrs, nil
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if len(nrs) == 0 ***REMOVED***
		return nil, fmt.Errorf("not enough resources available for task reservations: %+v", tr)
	***REMOVED***

	return nrs, nil
***REMOVED***

// Reclaim adds the resources taken by the task to the node's store
func Reclaim(nodeAvailableResources *[]*api.GenericResource, taskAssigned, nodeRes []*api.GenericResource) error ***REMOVED***
	err := reclaimResources(nodeAvailableResources, taskAssigned)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	sanitize(nodeRes, nodeAvailableResources)

	return nil
***REMOVED***

func reclaimResources(nodeAvailableResources *[]*api.GenericResource, taskAssigned []*api.GenericResource) error ***REMOVED***
	// The node could have been updated
	if nodeAvailableResources == nil ***REMOVED***
		return fmt.Errorf("node no longer has any resources")
	***REMOVED***

	for _, res := range taskAssigned ***REMOVED***
		switch tr := res.Resource.(type) ***REMOVED***
		case *api.GenericResource_DiscreteResourceSpec:
			nrs := GetResource(tr.DiscreteResourceSpec.Kind, *nodeAvailableResources)

			// If the resource went down to 0 it's no longer in the
			// available list
			if len(nrs) == 0 ***REMOVED***
				*nodeAvailableResources = append(*nodeAvailableResources, res.Copy())
			***REMOVED***

			if len(nrs) != 1 ***REMOVED***
				continue // Type change
			***REMOVED***

			nr := nrs[0].GetDiscreteResourceSpec()
			if nr == nil ***REMOVED***
				continue // Type change
			***REMOVED***

			nr.Value += tr.DiscreteResourceSpec.Value
		case *api.GenericResource_NamedResourceSpec:
			*nodeAvailableResources = append(*nodeAvailableResources, res.Copy())
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// sanitize checks that nodeAvailableResources does not add resources unknown
// to the nodeSpec (nodeRes) or goes over the integer bound specified
// by the spec.
// Note this is because the user is able to update a node's resources
func sanitize(nodeRes []*api.GenericResource, nodeAvailableResources *[]*api.GenericResource) ***REMOVED***
	// - We add the sanitized resources at the end, after
	// having removed the elements from the list

	// - When a set changes to a Discrete we also need
	// to make sure that we don't add the Discrete multiple
	// time hence, the need of a map to remember that
	var sanitized []*api.GenericResource
	kindSanitized := make(map[string]struct***REMOVED******REMOVED***)
	w := 0

	for _, na := range *nodeAvailableResources ***REMOVED***
		ok, nrs := sanitizeResource(nodeRes, na)
		if !ok ***REMOVED***
			if _, ok = kindSanitized[Kind(na)]; ok ***REMOVED***
				continue
			***REMOVED***

			kindSanitized[Kind(na)] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			sanitized = append(sanitized, nrs...)

			continue
		***REMOVED***

		(*nodeAvailableResources)[w] = na
		w++
	***REMOVED***

	*nodeAvailableResources = (*nodeAvailableResources)[:w]
	*nodeAvailableResources = append(*nodeAvailableResources, sanitized...)
***REMOVED***

// Returns true if the element is in nodeRes and "sane"
// Returns false if the element isn't in nodeRes and "sane" and the element(s) that should be replacing it
func sanitizeResource(nodeRes []*api.GenericResource, res *api.GenericResource) (ok bool, nrs []*api.GenericResource) ***REMOVED***
	switch na := res.Resource.(type) ***REMOVED***
	case *api.GenericResource_DiscreteResourceSpec:
		nrs := GetResource(na.DiscreteResourceSpec.Kind, nodeRes)

		// Type change or removed: reset
		if len(nrs) != 1 ***REMOVED***
			return false, nrs
		***REMOVED***

		// Type change: reset
		nr := nrs[0].GetDiscreteResourceSpec()
		if nr == nil ***REMOVED***
			return false, nrs
		***REMOVED***

		// Amount change: reset
		if na.DiscreteResourceSpec.Value > nr.Value ***REMOVED***
			return false, nrs
		***REMOVED***
	case *api.GenericResource_NamedResourceSpec:
		nrs := GetResource(na.NamedResourceSpec.Kind, nodeRes)

		// Type change
		if len(nrs) == 0 ***REMOVED***
			return false, nrs
		***REMOVED***

		for _, nr := range nrs ***REMOVED***
			// Type change: reset
			if nr.GetDiscreteResourceSpec() != nil ***REMOVED***
				return false, nrs
			***REMOVED***

			if na.NamedResourceSpec.Value == nr.GetNamedResourceSpec().Value ***REMOVED***
				return true, nil
			***REMOVED***
		***REMOVED***

		// Removed
		return false, nil
	***REMOVED***

	return true, nil
***REMOVED***
