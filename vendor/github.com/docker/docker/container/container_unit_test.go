package container

import (
	"path/filepath"
	"testing"

	"github.com/docker/docker/api/types/container"
	swarmtypes "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/pkg/signal"
)

func TestContainerStopSignal(t *testing.T) ***REMOVED***
	c := &Container***REMOVED***
		Config: &container.Config***REMOVED******REMOVED***,
	***REMOVED***

	def, err := signal.ParseSignal(signal.DefaultStopSignal)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	s := c.StopSignal()
	if s != int(def) ***REMOVED***
		t.Fatalf("Expected %v, got %v", def, s)
	***REMOVED***

	c = &Container***REMOVED***
		Config: &container.Config***REMOVED***StopSignal: "SIGKILL"***REMOVED***,
	***REMOVED***
	s = c.StopSignal()
	if s != 9 ***REMOVED***
		t.Fatalf("Expected 9, got %v", s)
	***REMOVED***
***REMOVED***

func TestContainerStopTimeout(t *testing.T) ***REMOVED***
	c := &Container***REMOVED***
		Config: &container.Config***REMOVED******REMOVED***,
	***REMOVED***

	s := c.StopTimeout()
	if s != DefaultStopTimeout ***REMOVED***
		t.Fatalf("Expected %v, got %v", DefaultStopTimeout, s)
	***REMOVED***

	stopTimeout := 15
	c = &Container***REMOVED***
		Config: &container.Config***REMOVED***StopTimeout: &stopTimeout***REMOVED***,
	***REMOVED***
	s = c.StopSignal()
	if s != 15 ***REMOVED***
		t.Fatalf("Expected 15, got %v", s)
	***REMOVED***
***REMOVED***

func TestContainerSecretReferenceDestTarget(t *testing.T) ***REMOVED***
	ref := &swarmtypes.SecretReference***REMOVED***
		File: &swarmtypes.SecretReferenceFileTarget***REMOVED***
			Name: "app",
		***REMOVED***,
	***REMOVED***

	d := getSecretTargetPath(ref)
	expected := filepath.Join(containerSecretMountPath, "app")
	if d != expected ***REMOVED***
		t.Fatalf("expected secret dest %q; received %q", expected, d)
	***REMOVED***
***REMOVED***
