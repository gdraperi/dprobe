package driverapi

import (
	"fmt"
)

// ErrNoNetwork is returned if no network with the specified id exists
type ErrNoNetwork string

func (enn ErrNoNetwork) Error() string ***REMOVED***
	return fmt.Sprintf("No network (%s) exists", string(enn))
***REMOVED***

// NotFound denotes the type of this error
func (enn ErrNoNetwork) NotFound() ***REMOVED******REMOVED***

// ErrEndpointExists is returned if more than one endpoint is added to the network
type ErrEndpointExists string

func (ee ErrEndpointExists) Error() string ***REMOVED***
	return fmt.Sprintf("Endpoint (%s) already exists (Only one endpoint allowed)", string(ee))
***REMOVED***

// Forbidden denotes the type of this error
func (ee ErrEndpointExists) Forbidden() ***REMOVED******REMOVED***

// ErrNotImplemented is returned when a Driver has not implemented an API yet
type ErrNotImplemented struct***REMOVED******REMOVED***

func (eni *ErrNotImplemented) Error() string ***REMOVED***
	return "The API is not implemented yet"
***REMOVED***

// NotImplemented denotes the type of this error
func (eni *ErrNotImplemented) NotImplemented() ***REMOVED******REMOVED***

// ErrNoEndpoint is returned if no endpoint with the specified id exists
type ErrNoEndpoint string

func (ene ErrNoEndpoint) Error() string ***REMOVED***
	return fmt.Sprintf("No endpoint (%s) exists", string(ene))
***REMOVED***

// NotFound denotes the type of this error
func (ene ErrNoEndpoint) NotFound() ***REMOVED******REMOVED***

// ErrActiveRegistration represents an error when a driver is registered to a networkType that is previously registered
type ErrActiveRegistration string

// Error interface for ErrActiveRegistration
func (ar ErrActiveRegistration) Error() string ***REMOVED***
	return fmt.Sprintf("Driver already registered for type %q", string(ar))
***REMOVED***

// Forbidden denotes the type of this error
func (ar ErrActiveRegistration) Forbidden() ***REMOVED******REMOVED***
