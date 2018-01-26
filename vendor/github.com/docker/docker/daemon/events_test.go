package daemon

import (
	"testing"
	"time"

	containertypes "github.com/docker/docker/api/types/container"
	eventtypes "github.com/docker/docker/api/types/events"
	"github.com/docker/docker/container"
	"github.com/docker/docker/daemon/events"
)

func TestLogContainerEventCopyLabels(t *testing.T) ***REMOVED***
	e := events.New()
	_, l, _ := e.Subscribe()
	defer e.Evict(l)

	container := &container.Container***REMOVED***
		ID:   "container_id",
		Name: "container_name",
		Config: &containertypes.Config***REMOVED***
			Image: "image_name",
			Labels: map[string]string***REMOVED***
				"node": "1",
				"os":   "alpine",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	daemon := &Daemon***REMOVED***
		EventsService: e,
	***REMOVED***
	daemon.LogContainerEvent(container, "create")

	if _, mutated := container.Config.Labels["image"]; mutated ***REMOVED***
		t.Fatalf("Expected to not mutate the container labels, got %q", container.Config.Labels)
	***REMOVED***

	validateTestAttributes(t, l, map[string]string***REMOVED***
		"node": "1",
		"os":   "alpine",
	***REMOVED***)
***REMOVED***

func TestLogContainerEventWithAttributes(t *testing.T) ***REMOVED***
	e := events.New()
	_, l, _ := e.Subscribe()
	defer e.Evict(l)

	container := &container.Container***REMOVED***
		ID:   "container_id",
		Name: "container_name",
		Config: &containertypes.Config***REMOVED***
			Labels: map[string]string***REMOVED***
				"node": "1",
				"os":   "alpine",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	daemon := &Daemon***REMOVED***
		EventsService: e,
	***REMOVED***
	attributes := map[string]string***REMOVED***
		"node": "2",
		"foo":  "bar",
	***REMOVED***
	daemon.LogContainerEventWithAttributes(container, "create", attributes)

	validateTestAttributes(t, l, map[string]string***REMOVED***
		"node": "1",
		"foo":  "bar",
	***REMOVED***)
***REMOVED***

func validateTestAttributes(t *testing.T, l chan interface***REMOVED******REMOVED***, expectedAttributesToTest map[string]string) ***REMOVED***
	select ***REMOVED***
	case ev := <-l:
		event, ok := ev.(eventtypes.Message)
		if !ok ***REMOVED***
			t.Fatalf("Unexpected event message: %q", ev)
		***REMOVED***
		for key, expected := range expectedAttributesToTest ***REMOVED***
			actual, ok := event.Actor.Attributes[key]
			if !ok || actual != expected ***REMOVED***
				t.Fatalf("Expected value for key %s to be %s, but was %s (event:%v)", key, expected, actual, event)
			***REMOVED***
		***REMOVED***
	case <-time.After(10 * time.Second):
		t.Fatal("LogEvent test timed out")
	***REMOVED***
***REMOVED***
