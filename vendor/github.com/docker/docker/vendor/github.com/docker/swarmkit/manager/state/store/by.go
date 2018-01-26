package store

import "github.com/docker/swarmkit/api"

// By is an interface type passed to Find methods. Implementations must be
// defined in this package.
type By interface ***REMOVED***
	// isBy allows this interface to only be satisfied by certain internal
	// types.
	isBy()
***REMOVED***

type byAll struct***REMOVED******REMOVED***

func (a byAll) isBy() ***REMOVED***
***REMOVED***

// All is an argument that can be passed to find to list all items in the
// set.
var All byAll

type byNamePrefix string

func (b byNamePrefix) isBy() ***REMOVED***
***REMOVED***

// ByNamePrefix creates an object to pass to Find to select by query.
func ByNamePrefix(namePrefix string) By ***REMOVED***
	return byNamePrefix(namePrefix)
***REMOVED***

type byIDPrefix string

func (b byIDPrefix) isBy() ***REMOVED***
***REMOVED***

// ByIDPrefix creates an object to pass to Find to select by query.
func ByIDPrefix(idPrefix string) By ***REMOVED***
	return byIDPrefix(idPrefix)
***REMOVED***

type byName string

func (b byName) isBy() ***REMOVED***
***REMOVED***

// ByName creates an object to pass to Find to select by name.
func ByName(name string) By ***REMOVED***
	return byName(name)
***REMOVED***

type byService string

func (b byService) isBy() ***REMOVED***
***REMOVED***

type byRuntime string

func (b byRuntime) isBy() ***REMOVED***
***REMOVED***

// ByRuntime creates an object to pass to Find to select by runtime.
func ByRuntime(runtime string) By ***REMOVED***
	return byRuntime(runtime)
***REMOVED***

// ByServiceID creates an object to pass to Find to select by service.
func ByServiceID(serviceID string) By ***REMOVED***
	return byService(serviceID)
***REMOVED***

type byNode string

func (b byNode) isBy() ***REMOVED***
***REMOVED***

// ByNodeID creates an object to pass to Find to select by node.
func ByNodeID(nodeID string) By ***REMOVED***
	return byNode(nodeID)
***REMOVED***

type bySlot struct ***REMOVED***
	serviceID string
	slot      uint64
***REMOVED***

func (b bySlot) isBy() ***REMOVED***
***REMOVED***

// BySlot creates an object to pass to Find to select by slot.
func BySlot(serviceID string, slot uint64) By ***REMOVED***
	return bySlot***REMOVED***serviceID: serviceID, slot: slot***REMOVED***
***REMOVED***

type byDesiredState api.TaskState

func (b byDesiredState) isBy() ***REMOVED***
***REMOVED***

// ByDesiredState creates an object to pass to Find to select by desired state.
func ByDesiredState(state api.TaskState) By ***REMOVED***
	return byDesiredState(state)
***REMOVED***

type byTaskState api.TaskState

func (b byTaskState) isBy() ***REMOVED***
***REMOVED***

// ByTaskState creates an object to pass to Find to select by task state.
func ByTaskState(state api.TaskState) By ***REMOVED***
	return byTaskState(state)
***REMOVED***

type byRole api.NodeRole

func (b byRole) isBy() ***REMOVED***
***REMOVED***

// ByRole creates an object to pass to Find to select by role.
func ByRole(role api.NodeRole) By ***REMOVED***
	return byRole(role)
***REMOVED***

type byMembership api.NodeSpec_Membership

func (b byMembership) isBy() ***REMOVED***
***REMOVED***

// ByMembership creates an object to pass to Find to select by Membership.
func ByMembership(membership api.NodeSpec_Membership) By ***REMOVED***
	return byMembership(membership)
***REMOVED***

type byReferencedNetworkID string

func (b byReferencedNetworkID) isBy() ***REMOVED***
***REMOVED***

// ByReferencedNetworkID creates an object to pass to Find to search for a
// service or task that references a network with the given ID.
func ByReferencedNetworkID(networkID string) By ***REMOVED***
	return byReferencedNetworkID(networkID)
***REMOVED***

type byReferencedSecretID string

func (b byReferencedSecretID) isBy() ***REMOVED***
***REMOVED***

// ByReferencedSecretID creates an object to pass to Find to search for a
// service or task that references a secret with the given ID.
func ByReferencedSecretID(secretID string) By ***REMOVED***
	return byReferencedSecretID(secretID)
***REMOVED***

type byReferencedConfigID string

func (b byReferencedConfigID) isBy() ***REMOVED***
***REMOVED***

// ByReferencedConfigID creates an object to pass to Find to search for a
// service or task that references a config with the given ID.
func ByReferencedConfigID(configID string) By ***REMOVED***
	return byReferencedConfigID(configID)
***REMOVED***

type byKind string

func (b byKind) isBy() ***REMOVED***
***REMOVED***

// ByKind creates an object to pass to Find to search for a Resource of a
// particular kind.
func ByKind(kind string) By ***REMOVED***
	return byKind(kind)
***REMOVED***

type byCustom struct ***REMOVED***
	objType string
	index   string
	value   string
***REMOVED***

func (b byCustom) isBy() ***REMOVED***
***REMOVED***

// ByCustom creates an object to pass to Find to search a custom index.
func ByCustom(objType, index, value string) By ***REMOVED***
	return byCustom***REMOVED***
		objType: objType,
		index:   index,
		value:   value,
	***REMOVED***
***REMOVED***

type byCustomPrefix struct ***REMOVED***
	objType string
	index   string
	value   string
***REMOVED***

func (b byCustomPrefix) isBy() ***REMOVED***
***REMOVED***

// ByCustomPrefix creates an object to pass to Find to search a custom index by
// a value prefix.
func ByCustomPrefix(objType, index, value string) By ***REMOVED***
	return byCustomPrefix***REMOVED***
		objType: objType,
		index:   index,
		value:   value,
	***REMOVED***
***REMOVED***
