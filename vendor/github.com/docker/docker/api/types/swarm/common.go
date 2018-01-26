package swarm

import "time"

// Version represents the internal object version.
type Version struct ***REMOVED***
	Index uint64 `json:",omitempty"`
***REMOVED***

// Meta is a base object inherited by most of the other once.
type Meta struct ***REMOVED***
	Version   Version   `json:",omitempty"`
	CreatedAt time.Time `json:",omitempty"`
	UpdatedAt time.Time `json:",omitempty"`
***REMOVED***

// Annotations represents how to describe an object.
type Annotations struct ***REMOVED***
	Name   string            `json:",omitempty"`
	Labels map[string]string `json:"Labels"`
***REMOVED***

// Driver represents a driver (network, logging, secrets backend).
type Driver struct ***REMOVED***
	Name    string            `json:",omitempty"`
	Options map[string]string `json:",omitempty"`
***REMOVED***

// TLSInfo represents the TLS information about what CA certificate is trusted,
// and who the issuer for a TLS certificate is
type TLSInfo struct ***REMOVED***
	// TrustRoot is the trusted CA root certificate in PEM format
	TrustRoot string `json:",omitempty"`

	// CertIssuer is the raw subject bytes of the issuer
	CertIssuerSubject []byte `json:",omitempty"`

	// CertIssuerPublicKey is the raw public key bytes of the issuer
	CertIssuerPublicKey []byte `json:",omitempty"`
***REMOVED***
