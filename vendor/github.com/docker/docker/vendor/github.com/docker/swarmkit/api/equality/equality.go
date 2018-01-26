package equality

import (
	"crypto/subtle"
	"reflect"

	"github.com/docker/swarmkit/api"
)

// TasksEqualStable returns true if the tasks are functionally equal, ignoring status,
// version and other superfluous fields.
//
// This used to decide whether or not to propagate a task update to a controller.
func TasksEqualStable(a, b *api.Task) bool ***REMOVED***
	// shallow copy
	copyA, copyB := *a, *b

	copyA.Status, copyB.Status = api.TaskStatus***REMOVED******REMOVED***, api.TaskStatus***REMOVED******REMOVED***
	copyA.Meta, copyB.Meta = api.Meta***REMOVED******REMOVED***, api.Meta***REMOVED******REMOVED***

	return reflect.DeepEqual(&copyA, &copyB)
***REMOVED***

// TaskStatusesEqualStable compares the task status excluding timestamp fields.
func TaskStatusesEqualStable(a, b *api.TaskStatus) bool ***REMOVED***
	copyA, copyB := *a, *b

	copyA.Timestamp, copyB.Timestamp = nil, nil
	return reflect.DeepEqual(&copyA, &copyB)
***REMOVED***

// RootCAEqualStable compares RootCAs, excluding join tokens, which are randomly generated
func RootCAEqualStable(a, b *api.RootCA) bool ***REMOVED***
	if a == nil && b == nil ***REMOVED***
		return true
	***REMOVED***
	if a == nil || b == nil ***REMOVED***
		return false
	***REMOVED***

	var aRotationKey, bRotationKey []byte
	if a.RootRotation != nil ***REMOVED***
		aRotationKey = a.RootRotation.CAKey
	***REMOVED***
	if b.RootRotation != nil ***REMOVED***
		bRotationKey = b.RootRotation.CAKey
	***REMOVED***
	if subtle.ConstantTimeCompare(a.CAKey, b.CAKey) != 1 || subtle.ConstantTimeCompare(aRotationKey, bRotationKey) != 1 ***REMOVED***
		return false
	***REMOVED***

	copyA, copyB := *a, *b
	copyA.JoinTokens, copyB.JoinTokens = api.JoinTokens***REMOVED******REMOVED***, api.JoinTokens***REMOVED******REMOVED***
	return reflect.DeepEqual(copyA, copyB)
***REMOVED***

// ExternalCAsEqualStable compares lists of external CAs and determines whether they are equal.
func ExternalCAsEqualStable(a, b []*api.ExternalCA) bool ***REMOVED***
	// because DeepEqual will treat an empty list and a nil list differently, we want to manually check this first
	if len(a) == 0 && len(b) == 0 ***REMOVED***
		return true
	***REMOVED***
	// The assumption is that each individual api.ExternalCA within both lists are created from deserializing from a
	// protobuf, so no special affordances are made to treat a nil map and empty map in the Options field of an
	// api.ExternalCA as equivalent.
	return reflect.DeepEqual(a, b)
***REMOVED***
