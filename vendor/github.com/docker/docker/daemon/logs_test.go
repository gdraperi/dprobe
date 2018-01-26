package daemon

import (
	"testing"

	containertypes "github.com/docker/docker/api/types/container"
)

func TestMergeAndVerifyLogConfigNilConfig(t *testing.T) ***REMOVED***
	d := &Daemon***REMOVED***defaultLogConfig: containertypes.LogConfig***REMOVED***Type: "json-file", Config: map[string]string***REMOVED***"max-file": "1"***REMOVED******REMOVED******REMOVED***
	cfg := containertypes.LogConfig***REMOVED***Type: d.defaultLogConfig.Type***REMOVED***
	if err := d.mergeAndVerifyLogConfig(&cfg); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
