// +build !windows

package config

import (
	"testing"

	"github.com/docker/docker/api/types"
)

func TestCommonUnixValidateConfigurationErrors(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		config *Config
	***REMOVED******REMOVED***
		// Can't override the stock runtime
		***REMOVED***
			config: &Config***REMOVED***
				CommonUnixConfig: CommonUnixConfig***REMOVED***
					Runtimes: map[string]types.Runtime***REMOVED***
						StockRuntimeName: ***REMOVED******REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		// Default runtime should be present in runtimes
		***REMOVED***
			config: &Config***REMOVED***
				CommonUnixConfig: CommonUnixConfig***REMOVED***
					Runtimes: map[string]types.Runtime***REMOVED***
						"foo": ***REMOVED******REMOVED***,
					***REMOVED***,
					DefaultRuntime: "bar",
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		err := Validate(tc.config)
		if err == nil ***REMOVED***
			t.Fatalf("expected error, got nil for config %v", tc.config)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestCommonUnixGetInitPath(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		config           *Config
		expectedInitPath string
	***REMOVED******REMOVED***
		***REMOVED***
			config: &Config***REMOVED***
				InitPath: "some-init-path",
			***REMOVED***,
			expectedInitPath: "some-init-path",
		***REMOVED***,
		***REMOVED***
			config: &Config***REMOVED***
				CommonUnixConfig: CommonUnixConfig***REMOVED***
					DefaultInitBinary: "foo-init-bin",
				***REMOVED***,
			***REMOVED***,
			expectedInitPath: "foo-init-bin",
		***REMOVED***,
		***REMOVED***
			config: &Config***REMOVED***
				InitPath: "init-path-A",
				CommonUnixConfig: CommonUnixConfig***REMOVED***
					DefaultInitBinary: "init-path-B",
				***REMOVED***,
			***REMOVED***,
			expectedInitPath: "init-path-A",
		***REMOVED***,
		***REMOVED***
			config:           &Config***REMOVED******REMOVED***,
			expectedInitPath: "docker-init",
		***REMOVED***,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		initPath := tc.config.GetInitPath()
		if initPath != tc.expectedInitPath ***REMOVED***
			t.Fatalf("expected initPath to be %v, got %v", tc.expectedInitPath, initPath)
		***REMOVED***
	***REMOVED***
***REMOVED***
