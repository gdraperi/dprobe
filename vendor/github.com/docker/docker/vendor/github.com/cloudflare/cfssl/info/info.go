// Package info contains the definitions for the info endpoint
package info

// Req is the request struct for an info API request.
type Req struct ***REMOVED***
	Label   string `json:"label"`
	Profile string `json:"profile"`
***REMOVED***

// Resp is the response for an Info API request.
type Resp struct ***REMOVED***
	Certificate  string   `json:"certificate"`
	Usage        []string `json:"usages"`
	ExpiryString string   `json:"expiry"`
***REMOVED***
