package api

import (
	"errors"
	"fmt"
	"strings"

	"github.com/docker/go-events"
)

var (
	errUnknownStoreAction = errors.New("unrecognized action type")
	errConflictingFilters = errors.New("conflicting filters specified")
	errNoKindSpecified    = errors.New("no kind of object specified")
	errUnrecognizedAction = errors.New("unrecognized action")
)

// StoreObject is an abstract object that can be handled by the store.
type StoreObject interface ***REMOVED***
	GetID() string                           // Get ID
	GetMeta() Meta                           // Retrieve metadata
	SetMeta(Meta)                            // Set metadata
	CopyStoreObject() StoreObject            // Return a copy of this object
	EventCreate() Event                      // Return a creation event
	EventUpdate(oldObject StoreObject) Event // Return an update event
	EventDelete() Event                      // Return a deletion event
***REMOVED***

// Event is the type used for events passed over watcher channels, and also
// the type used to specify filtering in calls to Watch.
type Event interface ***REMOVED***
	// TODO(stevvooe): Consider whether it makes sense to squish both the
	// matcher type and the primary type into the same type. It might be better
	// to build a matcher from an event prototype.

	// Matches checks if this item in a watch queue Matches the event
	// description.
	Matches(events.Event) bool
***REMOVED***

func customIndexer(kind string, annotations *Annotations) (bool, [][]byte, error) ***REMOVED***
	var converted [][]byte

	for _, entry := range annotations.Indices ***REMOVED***
		index := make([]byte, 0, len(kind)+1+len(entry.Key)+1+len(entry.Val)+1)
		if kind != "" ***REMOVED***
			index = append(index, []byte(kind)...)
			index = append(index, '|')
		***REMOVED***
		index = append(index, []byte(entry.Key)...)
		index = append(index, '|')
		index = append(index, []byte(entry.Val)...)
		index = append(index, '\x00')
		converted = append(converted, index)
	***REMOVED***

	// Add the null character as a terminator
	return len(converted) != 0, converted, nil
***REMOVED***

func fromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	if len(args) != 1 ***REMOVED***
		return nil, fmt.Errorf("must provide only a single argument")
	***REMOVED***
	arg, ok := args[0].(string)
	if !ok ***REMOVED***
		return nil, fmt.Errorf("argument must be a string: %#v", args[0])
	***REMOVED***
	// Add the null character as a terminator
	arg += "\x00"
	return []byte(arg), nil
***REMOVED***

func prefixFromArgs(args ...interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	val, err := fromArgs(args...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Strip the null terminator, the rest is a prefix
	n := len(val)
	if n > 0 ***REMOVED***
		return val[:n-1], nil
	***REMOVED***
	return val, nil
***REMOVED***

func checkCustom(a1, a2 Annotations) bool ***REMOVED***
	if len(a1.Indices) == 1 ***REMOVED***
		for _, ind := range a2.Indices ***REMOVED***
			if ind.Key == a1.Indices[0].Key && ind.Val == a1.Indices[0].Val ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func checkCustomPrefix(a1, a2 Annotations) bool ***REMOVED***
	if len(a1.Indices) == 1 ***REMOVED***
		for _, ind := range a2.Indices ***REMOVED***
			if ind.Key == a1.Indices[0].Key && strings.HasPrefix(ind.Val, a1.Indices[0].Val) ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
