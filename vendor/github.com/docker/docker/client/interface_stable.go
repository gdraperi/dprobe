package client

// APIClient is an interface that clients that talk with a docker server must implement.
type APIClient interface ***REMOVED***
	CommonAPIClient
	apiClientExperimental
***REMOVED***

// Ensure that Client always implements APIClient.
var _ APIClient = &Client***REMOVED******REMOVED***
