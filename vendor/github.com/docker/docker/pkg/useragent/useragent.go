// Package useragent provides helper functions to pack
// version information into a single User-Agent header.
package useragent

import (
	"strings"
)

// VersionInfo is used to model UserAgent versions.
type VersionInfo struct ***REMOVED***
	Name    string
	Version string
***REMOVED***

func (vi *VersionInfo) isValid() bool ***REMOVED***
	const stopChars = " \t\r\n/"
	name := vi.Name
	vers := vi.Version
	if len(name) == 0 || strings.ContainsAny(name, stopChars) ***REMOVED***
		return false
	***REMOVED***
	if len(vers) == 0 || strings.ContainsAny(vers, stopChars) ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

// AppendVersions converts versions to a string and appends the string to the string base.
//
// Each VersionInfo will be converted to a string in the format of
// "product/version", where the "product" is get from the name field, while
// version is get from the version field. Several pieces of version information
// will be concatenated and separated by space.
//
// Example:
// AppendVersions("base", VersionInfo***REMOVED***"foo", "1.0"***REMOVED***, VersionInfo***REMOVED***"bar", "2.0"***REMOVED***)
// results in "base foo/1.0 bar/2.0".
func AppendVersions(base string, versions ...VersionInfo) string ***REMOVED***
	if len(versions) == 0 ***REMOVED***
		return base
	***REMOVED***

	verstrs := make([]string, 0, 1+len(versions))
	if len(base) > 0 ***REMOVED***
		verstrs = append(verstrs, base)
	***REMOVED***

	for _, v := range versions ***REMOVED***
		if !v.isValid() ***REMOVED***
			continue
		***REMOVED***
		verstrs = append(verstrs, v.Name+"/"+v.Version)
	***REMOVED***
	return strings.Join(verstrs, " ")
***REMOVED***
