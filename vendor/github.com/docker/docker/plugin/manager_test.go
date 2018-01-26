package plugin

import (
	"testing"

	"github.com/docker/docker/api/types"
)

func TestValidatePrivileges(t *testing.T) ***REMOVED***
	testData := map[string]struct ***REMOVED***
		requiredPrivileges types.PluginPrivileges
		privileges         types.PluginPrivileges
		result             bool
	***REMOVED******REMOVED***
		"diff-len": ***REMOVED***
			requiredPrivileges: []types.PluginPrivilege***REMOVED***
				***REMOVED***Name: "Privilege1", Description: "Description", Value: []string***REMOVED***"abc", "def", "ghi"***REMOVED******REMOVED***,
			***REMOVED***,
			privileges: []types.PluginPrivilege***REMOVED***
				***REMOVED***Name: "Privilege1", Description: "Description", Value: []string***REMOVED***"abc", "def", "ghi"***REMOVED******REMOVED***,
				***REMOVED***Name: "Privilege2", Description: "Description", Value: []string***REMOVED***"123", "456", "789"***REMOVED******REMOVED***,
			***REMOVED***,
			result: false,
		***REMOVED***,
		"diff-value": ***REMOVED***
			requiredPrivileges: []types.PluginPrivilege***REMOVED***
				***REMOVED***Name: "Privilege1", Description: "Description", Value: []string***REMOVED***"abc", "def", "GHI"***REMOVED******REMOVED***,
				***REMOVED***Name: "Privilege2", Description: "Description", Value: []string***REMOVED***"123", "456", "***"***REMOVED******REMOVED***,
			***REMOVED***,
			privileges: []types.PluginPrivilege***REMOVED***
				***REMOVED***Name: "Privilege1", Description: "Description", Value: []string***REMOVED***"abc", "def", "ghi"***REMOVED******REMOVED***,
				***REMOVED***Name: "Privilege2", Description: "Description", Value: []string***REMOVED***"123", "456", "789"***REMOVED******REMOVED***,
			***REMOVED***,
			result: false,
		***REMOVED***,
		"diff-order-but-same-value": ***REMOVED***
			requiredPrivileges: []types.PluginPrivilege***REMOVED***
				***REMOVED***Name: "Privilege1", Description: "Description", Value: []string***REMOVED***"abc", "def", "GHI"***REMOVED******REMOVED***,
				***REMOVED***Name: "Privilege2", Description: "Description", Value: []string***REMOVED***"123", "456", "789"***REMOVED******REMOVED***,
			***REMOVED***,
			privileges: []types.PluginPrivilege***REMOVED***
				***REMOVED***Name: "Privilege2", Description: "Description", Value: []string***REMOVED***"123", "456", "789"***REMOVED******REMOVED***,
				***REMOVED***Name: "Privilege1", Description: "Description", Value: []string***REMOVED***"GHI", "abc", "def"***REMOVED******REMOVED***,
			***REMOVED***,
			result: true,
		***REMOVED***,
	***REMOVED***

	for key, data := range testData ***REMOVED***
		err := validatePrivileges(data.requiredPrivileges, data.privileges)
		if (err == nil) != data.result ***REMOVED***
			t.Fatalf("Test item %s expected result to be %t, got %t", key, data.result, (err == nil))
		***REMOVED***
	***REMOVED***
***REMOVED***
