package oci

import (
	"context"

	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func createDefaultSpec(ctx context.Context, id string) (*specs.Spec, error) ***REMOVED***
	return &specs.Spec***REMOVED***
		Version: specs.Version,
		Root:    &specs.Root***REMOVED******REMOVED***,
		Process: &specs.Process***REMOVED***
			Cwd: `C:\`,
			ConsoleSize: &specs.Box***REMOVED***
				Width:  80,
				Height: 20,
			***REMOVED***,
		***REMOVED***,
		Windows: &specs.Windows***REMOVED***
			IgnoreFlushesDuringBoot: true,
			Network: &specs.WindowsNetwork***REMOVED***
				AllowUnqualifiedDNSQuery: true,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***, nil
***REMOVED***
