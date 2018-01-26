package libcontainerd

import (
	"testing"
)

func TestEnvironmentParsing(t *testing.T) ***REMOVED***
	env := []string***REMOVED***"foo=bar", "car=hat", "a=b=c"***REMOVED***
	result := setupEnvironmentVariables(env)
	if len(result) != 3 || result["foo"] != "bar" || result["car"] != "hat" || result["a"] != "b=c" ***REMOVED***
		t.Fatalf("Expected map[foo:bar car:hat a:b=c], got %v", result)
	***REMOVED***
***REMOVED***
