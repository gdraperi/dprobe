// +build !windows

package network

import (
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

func TestFilterNetworks(t *testing.T) ***REMOVED***
	networks := []types.NetworkResource***REMOVED***
		***REMOVED***
			Name:   "host",
			Driver: "host",
			Scope:  "local",
		***REMOVED***,
		***REMOVED***
			Name:   "bridge",
			Driver: "bridge",
			Scope:  "local",
		***REMOVED***,
		***REMOVED***
			Name:   "none",
			Driver: "null",
			Scope:  "local",
		***REMOVED***,
		***REMOVED***
			Name:   "myoverlay",
			Driver: "overlay",
			Scope:  "swarm",
		***REMOVED***,
		***REMOVED***
			Name:   "mydrivernet",
			Driver: "mydriver",
			Scope:  "local",
		***REMOVED***,
		***REMOVED***
			Name:   "mykvnet",
			Driver: "mykvdriver",
			Scope:  "global",
		***REMOVED***,
	***REMOVED***

	bridgeDriverFilters := filters.NewArgs()
	bridgeDriverFilters.Add("driver", "bridge")

	overlayDriverFilters := filters.NewArgs()
	overlayDriverFilters.Add("driver", "overlay")

	nonameDriverFilters := filters.NewArgs()
	nonameDriverFilters.Add("driver", "noname")

	customDriverFilters := filters.NewArgs()
	customDriverFilters.Add("type", "custom")

	builtinDriverFilters := filters.NewArgs()
	builtinDriverFilters.Add("type", "builtin")

	invalidDriverFilters := filters.NewArgs()
	invalidDriverFilters.Add("type", "invalid")

	localScopeFilters := filters.NewArgs()
	localScopeFilters.Add("scope", "local")

	swarmScopeFilters := filters.NewArgs()
	swarmScopeFilters.Add("scope", "swarm")

	globalScopeFilters := filters.NewArgs()
	globalScopeFilters.Add("scope", "global")

	testCases := []struct ***REMOVED***
		filter      filters.Args
		resultCount int
		err         string
	***REMOVED******REMOVED***
		***REMOVED***
			filter:      bridgeDriverFilters,
			resultCount: 1,
			err:         "",
		***REMOVED***,
		***REMOVED***
			filter:      overlayDriverFilters,
			resultCount: 1,
			err:         "",
		***REMOVED***,
		***REMOVED***
			filter:      nonameDriverFilters,
			resultCount: 0,
			err:         "",
		***REMOVED***,
		***REMOVED***
			filter:      customDriverFilters,
			resultCount: 3,
			err:         "",
		***REMOVED***,
		***REMOVED***
			filter:      builtinDriverFilters,
			resultCount: 3,
			err:         "",
		***REMOVED***,
		***REMOVED***
			filter:      invalidDriverFilters,
			resultCount: 0,
			err:         "Invalid filter: 'type'='invalid'",
		***REMOVED***,
		***REMOVED***
			filter:      localScopeFilters,
			resultCount: 4,
			err:         "",
		***REMOVED***,
		***REMOVED***
			filter:      swarmScopeFilters,
			resultCount: 1,
			err:         "",
		***REMOVED***,
		***REMOVED***
			filter:      globalScopeFilters,
			resultCount: 1,
			err:         "",
		***REMOVED***,
	***REMOVED***

	for _, testCase := range testCases ***REMOVED***
		result, err := filterNetworks(networks, testCase.filter)
		if testCase.err != "" ***REMOVED***
			if err == nil ***REMOVED***
				t.Fatalf("expect error '%s', got no error", testCase.err)

			***REMOVED*** else if !strings.Contains(err.Error(), testCase.err) ***REMOVED***
				t.Fatalf("expect error '%s', got '%s'", testCase.err, err)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if err != nil ***REMOVED***
				t.Fatalf("expect no error, got error '%s'", err)
			***REMOVED***
			// Make sure result is not nil
			if result == nil ***REMOVED***
				t.Fatal("filterNetworks should not return nil")
			***REMOVED***

			if len(result) != testCase.resultCount ***REMOVED***
				t.Fatalf("expect '%d' networks, got '%d' networks", testCase.resultCount, len(result))
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
