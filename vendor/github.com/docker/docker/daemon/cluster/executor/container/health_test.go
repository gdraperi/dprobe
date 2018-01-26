// +build !windows

package container

import (
	"testing"
	"time"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/container"
	"github.com/docker/docker/daemon"
	"github.com/docker/docker/daemon/events"
	"github.com/docker/swarmkit/api"
	"golang.org/x/net/context"
)

func TestHealthStates(t *testing.T) ***REMOVED***

	// set up environment: events, task, container ....
	e := events.New()
	_, l, _ := e.Subscribe()
	defer e.Evict(l)

	task := &api.Task***REMOVED***
		ID:        "id",
		ServiceID: "sid",
		Spec: api.TaskSpec***REMOVED***
			Runtime: &api.TaskSpec_Container***REMOVED***
				Container: &api.ContainerSpec***REMOVED***
					Image: "image_name",
					Labels: map[string]string***REMOVED***
						"com.docker.swarm.task.id": "id",
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		Annotations: api.Annotations***REMOVED***Name: "name"***REMOVED***,
	***REMOVED***

	c := &container.Container***REMOVED***
		ID:   "id",
		Name: "name",
		Config: &containertypes.Config***REMOVED***
			Image: "image_name",
			Labels: map[string]string***REMOVED***
				"com.docker.swarm.task.id": "id",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	daemon := &daemon.Daemon***REMOVED***
		EventsService: e,
	***REMOVED***

	controller, err := newController(daemon, task, nil, nil)
	if err != nil ***REMOVED***
		t.Fatalf("create controller fail %v", err)
	***REMOVED***

	errChan := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// fire checkHealth
	go func() ***REMOVED***
		err := controller.checkHealth(ctx)
		select ***REMOVED***
		case errChan <- err:
		case <-ctx.Done():
		***REMOVED***
	***REMOVED***()

	// send an event and expect to get expectedErr
	// if expectedErr is nil, shouldn't get any error
	logAndExpect := func(msg string, expectedErr error) ***REMOVED***
		daemon.LogContainerEvent(c, msg)

		timer := time.NewTimer(1 * time.Second)
		defer timer.Stop()

		select ***REMOVED***
		case err := <-errChan:
			if err != expectedErr ***REMOVED***
				t.Fatalf("expect error %v, but get %v", expectedErr, err)
			***REMOVED***
		case <-timer.C:
			if expectedErr != nil ***REMOVED***
				t.Fatal("time limit exceeded, didn't get expected error")
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// events that are ignored by checkHealth
	logAndExpect("health_status: running", nil)
	logAndExpect("health_status: healthy", nil)
	logAndExpect("die", nil)

	// unhealthy event will be caught by checkHealth
	logAndExpect("health_status: unhealthy", ErrContainerUnhealthy)
***REMOVED***
