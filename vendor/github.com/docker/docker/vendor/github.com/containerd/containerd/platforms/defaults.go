package platforms

import (
	"runtime"

	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// Default returns the default specifier for the platform.
func Default() string ***REMOVED***
	return Format(DefaultSpec())
***REMOVED***

// DefaultSpec returns the current platform's default platform specification.
func DefaultSpec() specs.Platform ***REMOVED***
	return specs.Platform***REMOVED***
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		// TODO(stevvooe): Need to resolve GOARM for arm hosts.
	***REMOVED***
***REMOVED***
