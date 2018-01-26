// Package config in the ocsp directory provides configuration data for an OCSP
// signer.
package config

import "time"

// Config contains configuration information required to set up an OCSP signer.
type Config struct ***REMOVED***
	CACertFile        string
	ResponderCertFile string
	KeyFile           string
	Interval          time.Duration
***REMOVED***
