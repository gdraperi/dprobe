package auth

import (
	"net/http"
	"strings"
)

// APIVersion represents a version of an API including its
// type and version number.
type APIVersion struct ***REMOVED***
	// Type refers to the name of a specific API specification
	// such as "registry"
	Type string

	// Version is the version of the API specification implemented,
	// This may omit the revision number and only include
	// the major and minor version, such as "2.0"
	Version string
***REMOVED***

// String returns the string formatted API Version
func (v APIVersion) String() string ***REMOVED***
	return v.Type + "/" + v.Version
***REMOVED***

// APIVersions gets the API versions out of an HTTP response using the provided
// version header as the key for the HTTP header.
func APIVersions(resp *http.Response, versionHeader string) []APIVersion ***REMOVED***
	versions := []APIVersion***REMOVED******REMOVED***
	if versionHeader != "" ***REMOVED***
		for _, supportedVersions := range resp.Header[http.CanonicalHeaderKey(versionHeader)] ***REMOVED***
			for _, version := range strings.Fields(supportedVersions) ***REMOVED***
				versions = append(versions, ParseAPIVersion(version))
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return versions
***REMOVED***

// ParseAPIVersion parses an API version string into an APIVersion
// Format (Expected, not enforced):
// API version string = <API type> '/' <API version>
// API type = [a-z][a-z0-9]*
// API version = [0-9]+(\.[0-9]+)?
// TODO(dmcgowan): Enforce format, add error condition, remove unknown type
func ParseAPIVersion(versionStr string) APIVersion ***REMOVED***
	idx := strings.IndexRune(versionStr, '/')
	if idx == -1 ***REMOVED***
		return APIVersion***REMOVED***
			Type:    "unknown",
			Version: versionStr,
		***REMOVED***
	***REMOVED***
	return APIVersion***REMOVED***
		Type:    strings.ToLower(versionStr[:idx]),
		Version: versionStr[idx+1:],
	***REMOVED***
***REMOVED***
