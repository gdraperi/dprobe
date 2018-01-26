package libnetwork

import (
	"fmt"
)

// ErrNoSuchNetwork is returned when a network query finds no result
type ErrNoSuchNetwork string

func (nsn ErrNoSuchNetwork) Error() string ***REMOVED***
	return fmt.Sprintf("network %s not found", string(nsn))
***REMOVED***

// NotFound denotes the type of this error
func (nsn ErrNoSuchNetwork) NotFound() ***REMOVED******REMOVED***

// ErrNoSuchEndpoint is returned when an endpoint query finds no result
type ErrNoSuchEndpoint string

func (nse ErrNoSuchEndpoint) Error() string ***REMOVED***
	return fmt.Sprintf("endpoint %s not found", string(nse))
***REMOVED***

// NotFound denotes the type of this error
func (nse ErrNoSuchEndpoint) NotFound() ***REMOVED******REMOVED***

// ErrInvalidNetworkDriver is returned if an invalid driver
// name is passed.
type ErrInvalidNetworkDriver string

func (ind ErrInvalidNetworkDriver) Error() string ***REMOVED***
	return fmt.Sprintf("invalid driver bound to network: %s", string(ind))
***REMOVED***

// BadRequest denotes the type of this error
func (ind ErrInvalidNetworkDriver) BadRequest() ***REMOVED******REMOVED***

// ErrInvalidJoin is returned if a join is attempted on an endpoint
// which already has a container joined.
type ErrInvalidJoin struct***REMOVED******REMOVED***

func (ij ErrInvalidJoin) Error() string ***REMOVED***
	return "a container has already joined the endpoint"
***REMOVED***

// BadRequest denotes the type of this error
func (ij ErrInvalidJoin) BadRequest() ***REMOVED******REMOVED***

// ErrNoContainer is returned when the endpoint has no container
// attached to it.
type ErrNoContainer struct***REMOVED******REMOVED***

func (nc ErrNoContainer) Error() string ***REMOVED***
	return "no container is attached to the endpoint"
***REMOVED***

// Maskable denotes the type of this error
func (nc ErrNoContainer) Maskable() ***REMOVED******REMOVED***

// ErrInvalidID is returned when a query-by-id method is being invoked
// with an empty id parameter
type ErrInvalidID string

func (ii ErrInvalidID) Error() string ***REMOVED***
	return fmt.Sprintf("invalid id: %s", string(ii))
***REMOVED***

// BadRequest denotes the type of this error
func (ii ErrInvalidID) BadRequest() ***REMOVED******REMOVED***

// ErrInvalidName is returned when a query-by-name or resource create method is
// invoked with an empty name parameter
type ErrInvalidName string

func (in ErrInvalidName) Error() string ***REMOVED***
	return fmt.Sprintf("invalid name: %s", string(in))
***REMOVED***

// BadRequest denotes the type of this error
func (in ErrInvalidName) BadRequest() ***REMOVED******REMOVED***

// ErrInvalidConfigFile type is returned when an invalid LibNetwork config file is detected
type ErrInvalidConfigFile string

func (cf ErrInvalidConfigFile) Error() string ***REMOVED***
	return fmt.Sprintf("Invalid Config file %q", string(cf))
***REMOVED***

// NetworkTypeError type is returned when the network type string is not
// known to libnetwork.
type NetworkTypeError string

func (nt NetworkTypeError) Error() string ***REMOVED***
	return fmt.Sprintf("unknown driver %q", string(nt))
***REMOVED***

// NotFound denotes the type of this error
func (nt NetworkTypeError) NotFound() ***REMOVED******REMOVED***

// NetworkNameError is returned when a network with the same name already exists.
type NetworkNameError string

func (nnr NetworkNameError) Error() string ***REMOVED***
	return fmt.Sprintf("network with name %s already exists", string(nnr))
***REMOVED***

// Forbidden denotes the type of this error
func (nnr NetworkNameError) Forbidden() ***REMOVED******REMOVED***

// UnknownNetworkError is returned when libnetwork could not find in its database
// a network with the same name and id.
type UnknownNetworkError struct ***REMOVED***
	name string
	id   string
***REMOVED***

func (une *UnknownNetworkError) Error() string ***REMOVED***
	return fmt.Sprintf("unknown network %s id %s", une.name, une.id)
***REMOVED***

// NotFound denotes the type of this error
func (une *UnknownNetworkError) NotFound() ***REMOVED******REMOVED***

// ActiveEndpointsError is returned when a network is deleted which has active
// endpoints in it.
type ActiveEndpointsError struct ***REMOVED***
	name string
	id   string
***REMOVED***

func (aee *ActiveEndpointsError) Error() string ***REMOVED***
	return fmt.Sprintf("network %s id %s has active endpoints", aee.name, aee.id)
***REMOVED***

// Forbidden denotes the type of this error
func (aee *ActiveEndpointsError) Forbidden() ***REMOVED******REMOVED***

// UnknownEndpointError is returned when libnetwork could not find in its database
// an endpoint with the same name and id.
type UnknownEndpointError struct ***REMOVED***
	name string
	id   string
***REMOVED***

func (uee *UnknownEndpointError) Error() string ***REMOVED***
	return fmt.Sprintf("unknown endpoint %s id %s", uee.name, uee.id)
***REMOVED***

// NotFound denotes the type of this error
func (uee *UnknownEndpointError) NotFound() ***REMOVED******REMOVED***

// ActiveContainerError is returned when an endpoint is deleted which has active
// containers attached to it.
type ActiveContainerError struct ***REMOVED***
	name string
	id   string
***REMOVED***

func (ace *ActiveContainerError) Error() string ***REMOVED***
	return fmt.Sprintf("endpoint with name %s id %s has active containers", ace.name, ace.id)
***REMOVED***

// Forbidden denotes the type of this error
func (ace *ActiveContainerError) Forbidden() ***REMOVED******REMOVED***

// InvalidContainerIDError is returned when an invalid container id is passed
// in Join/Leave
type InvalidContainerIDError string

func (id InvalidContainerIDError) Error() string ***REMOVED***
	return fmt.Sprintf("invalid container id %s", string(id))
***REMOVED***

// BadRequest denotes the type of this error
func (id InvalidContainerIDError) BadRequest() ***REMOVED******REMOVED***

// ManagerRedirectError is returned when the request should be redirected to Manager
type ManagerRedirectError string

func (mr ManagerRedirectError) Error() string ***REMOVED***
	return "Redirect the request to the manager"
***REMOVED***

// Maskable denotes the type of this error
func (mr ManagerRedirectError) Maskable() ***REMOVED******REMOVED***

// ErrDataStoreNotInitialized is returned if an invalid data scope is passed
// for getting data store
type ErrDataStoreNotInitialized string

func (dsni ErrDataStoreNotInitialized) Error() string ***REMOVED***
	return fmt.Sprintf("datastore for scope %q is not initialized", string(dsni))
***REMOVED***
