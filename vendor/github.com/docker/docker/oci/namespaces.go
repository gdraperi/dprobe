package oci

import specs "github.com/opencontainers/runtime-spec/specs-go"

// RemoveNamespace removes the `nsType` namespace from OCI spec `s`
func RemoveNamespace(s *specs.Spec, nsType specs.LinuxNamespaceType) ***REMOVED***
	for i, n := range s.Linux.Namespaces ***REMOVED***
		if n.Type == nsType ***REMOVED***
			s.Linux.Namespaces = append(s.Linux.Namespaces[:i], s.Linux.Namespaces[i+1:]...)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***
