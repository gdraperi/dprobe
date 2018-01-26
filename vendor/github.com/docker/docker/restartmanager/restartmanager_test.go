package restartmanager

import (
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
)

func TestRestartManagerTimeout(t *testing.T) ***REMOVED***
	rm := New(container.RestartPolicy***REMOVED***Name: "always"***REMOVED***, 0).(*restartManager)
	var duration = time.Duration(1 * time.Second)
	should, _, err := rm.ShouldRestart(0, false, duration)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !should ***REMOVED***
		t.Fatal("container should be restarted")
	***REMOVED***
	if rm.timeout != defaultTimeout ***REMOVED***
		t.Fatalf("restart manager should have a timeout of 100 ms but has %s", rm.timeout)
	***REMOVED***
***REMOVED***

func TestRestartManagerTimeoutReset(t *testing.T) ***REMOVED***
	rm := New(container.RestartPolicy***REMOVED***Name: "always"***REMOVED***, 0).(*restartManager)
	rm.timeout = 5 * time.Second
	var duration = time.Duration(10 * time.Second)
	_, _, err := rm.ShouldRestart(0, false, duration)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if rm.timeout != defaultTimeout ***REMOVED***
		t.Fatalf("restart manager should have a timeout of 100 ms but has %s", rm.timeout)
	***REMOVED***
***REMOVED***
