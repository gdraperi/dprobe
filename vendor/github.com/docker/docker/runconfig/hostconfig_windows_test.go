// +build windows

package runconfig

import (
	"testing"

	"github.com/docker/docker/api/types/container"
)

func TestValidatePrivileged(t *testing.T) ***REMOVED***
	expected := "Windows does not support privileged mode"
	err := validatePrivileged(&container.HostConfig***REMOVED***Privileged: true***REMOVED***)
	if err == nil || err.Error() != expected ***REMOVED***
		t.Fatalf("Expected %s", expected)
	***REMOVED***
***REMOVED***
