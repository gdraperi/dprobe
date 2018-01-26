package system

import (
	"fmt"
	"runtime"
	"strings"

	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// ValidatePlatform determines if a platform structure is valid.
// TODO This is a temporary function - can be replaced by parsing from
// https://github.com/containerd/containerd/pull/1403/files at a later date.
// @jhowardmsft
func ValidatePlatform(platform *specs.Platform) error ***REMOVED***
	platform.Architecture = strings.ToLower(platform.Architecture)
	platform.OS = strings.ToLower(platform.OS)
	// Based on https://github.com/moby/moby/pull/34642#issuecomment-330375350, do
	// not support anything except operating system.
	if platform.Architecture != "" ***REMOVED***
		return fmt.Errorf("invalid platform architecture %q", platform.Architecture)
	***REMOVED***
	if platform.OS != "" ***REMOVED***
		if !(platform.OS == runtime.GOOS || (LCOWSupported() && platform.OS == "linux")) ***REMOVED***
			return fmt.Errorf("invalid platform os %q", platform.OS)
		***REMOVED***
	***REMOVED***
	if len(platform.OSFeatures) != 0 ***REMOVED***
		return fmt.Errorf("invalid platform osfeatures %q", platform.OSFeatures)
	***REMOVED***
	if platform.OSVersion != "" ***REMOVED***
		return fmt.Errorf("invalid platform osversion %q", platform.OSVersion)
	***REMOVED***
	if platform.Variant != "" ***REMOVED***
		return fmt.Errorf("invalid platform variant %q", platform.Variant)
	***REMOVED***
	return nil
***REMOVED***

// ParsePlatform parses a platform string in the format os[/arch[/variant]
// into an OCI image-spec platform structure.
// TODO This is a temporary function - can be replaced by parsing from
// https://github.com/containerd/containerd/pull/1403/files at a later date.
// @jhowardmsft
func ParsePlatform(in string) *specs.Platform ***REMOVED***
	p := &specs.Platform***REMOVED******REMOVED***
	elements := strings.SplitN(strings.ToLower(in), "/", 3)
	if len(elements) == 3 ***REMOVED***
		p.Variant = elements[2]
	***REMOVED***
	if len(elements) >= 2 ***REMOVED***
		p.Architecture = elements[1]
	***REMOVED***
	if len(elements) >= 1 ***REMOVED***
		p.OS = elements[0]
	***REMOVED***
	return p
***REMOVED***

// IsOSSupported determines if an operating system is supported by the host
func IsOSSupported(os string) bool ***REMOVED***
	if runtime.GOOS == os ***REMOVED***
		return true
	***REMOVED***
	if LCOWSupported() && os == "linux" ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***
