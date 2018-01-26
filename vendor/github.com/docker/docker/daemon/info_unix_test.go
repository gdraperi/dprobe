// +build !windows

package daemon

import (
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/dockerversion"
	"github.com/stretchr/testify/assert"
)

func TestParseInitVersion(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		version string
		result  types.Commit
		invalid bool
	***REMOVED******REMOVED***
		***REMOVED***
			version: "tini version 0.13.0 - git.949e6fa",
			result:  types.Commit***REMOVED***ID: "949e6fa", Expected: dockerversion.InitCommitID[0:7]***REMOVED***,
		***REMOVED***, ***REMOVED***
			version: "tini version 0.13.0\n",
			result:  types.Commit***REMOVED***ID: "v0.13.0", Expected: dockerversion.InitCommitID***REMOVED***,
		***REMOVED***, ***REMOVED***
			version: "tini version 0.13.2",
			result:  types.Commit***REMOVED***ID: "v0.13.2", Expected: dockerversion.InitCommitID***REMOVED***,
		***REMOVED***, ***REMOVED***
			version: "tini version0.13.2",
			result:  types.Commit***REMOVED***ID: "N/A", Expected: dockerversion.InitCommitID***REMOVED***,
			invalid: true,
		***REMOVED***, ***REMOVED***
			version: "",
			result:  types.Commit***REMOVED***ID: "N/A", Expected: dockerversion.InitCommitID***REMOVED***,
			invalid: true,
		***REMOVED***, ***REMOVED***
			version: "hello world",
			result:  types.Commit***REMOVED***ID: "N/A", Expected: dockerversion.InitCommitID***REMOVED***,
			invalid: true,
		***REMOVED***,
	***REMOVED***

	for _, test := range tests ***REMOVED***
		ver, err := parseInitVersion(string(test.version))
		if test.invalid ***REMOVED***
			assert.Error(t, err)
		***REMOVED*** else ***REMOVED***
			assert.NoError(t, err)
		***REMOVED***
		assert.Equal(t, test.result, ver)
	***REMOVED***
***REMOVED***
