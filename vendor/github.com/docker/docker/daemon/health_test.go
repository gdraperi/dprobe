package daemon

import (
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	eventtypes "github.com/docker/docker/api/types/events"
	"github.com/docker/docker/container"
	"github.com/docker/docker/daemon/events"
)

func reset(c *container.Container) ***REMOVED***
	c.State = &container.State***REMOVED******REMOVED***
	c.State.Health = &container.Health***REMOVED******REMOVED***
	c.State.Health.SetStatus(types.Starting)
***REMOVED***

func TestNoneHealthcheck(t *testing.T) ***REMOVED***
	c := &container.Container***REMOVED***
		ID:   "container_id",
		Name: "container_name",
		Config: &containertypes.Config***REMOVED***
			Image: "image_name",
			Healthcheck: &containertypes.HealthConfig***REMOVED***
				Test: []string***REMOVED***"NONE"***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		State: &container.State***REMOVED******REMOVED***,
	***REMOVED***
	store, err := container.NewViewDB()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	daemon := &Daemon***REMOVED***
		containersReplica: store,
	***REMOVED***

	daemon.initHealthMonitor(c)
	if c.State.Health != nil ***REMOVED***
		t.Error("Expecting Health to be nil, but was not")
	***REMOVED***
***REMOVED***

// FIXME(vdemeester) This takes around 3s… This is *way* too long
func TestHealthStates(t *testing.T) ***REMOVED***
	e := events.New()
	_, l, _ := e.Subscribe()
	defer e.Evict(l)

	expect := func(expected string) ***REMOVED***
		select ***REMOVED***
		case event := <-l:
			ev := event.(eventtypes.Message)
			if ev.Status != expected ***REMOVED***
				t.Errorf("Expecting event %#v, but got %#v\n", expected, ev.Status)
			***REMOVED***
		case <-time.After(1 * time.Second):
			t.Errorf("Expecting event %#v, but got nothing\n", expected)
		***REMOVED***
	***REMOVED***

	c := &container.Container***REMOVED***
		ID:   "container_id",
		Name: "container_name",
		Config: &containertypes.Config***REMOVED***
			Image: "image_name",
		***REMOVED***,
	***REMOVED***

	store, err := container.NewViewDB()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	daemon := &Daemon***REMOVED***
		EventsService:     e,
		containersReplica: store,
	***REMOVED***

	c.Config.Healthcheck = &containertypes.HealthConfig***REMOVED***
		Retries: 1,
	***REMOVED***

	reset(c)

	handleResult := func(startTime time.Time, exitCode int) ***REMOVED***
		handleProbeResult(daemon, c, &types.HealthcheckResult***REMOVED***
			Start:    startTime,
			End:      startTime,
			ExitCode: exitCode,
		***REMOVED***, nil)
	***REMOVED***

	// starting -> failed -> success -> failed

	handleResult(c.State.StartedAt.Add(1*time.Second), 1)
	expect("health_status: unhealthy")

	handleResult(c.State.StartedAt.Add(2*time.Second), 0)
	expect("health_status: healthy")

	handleResult(c.State.StartedAt.Add(3*time.Second), 1)
	expect("health_status: unhealthy")

	// Test retries

	reset(c)
	c.Config.Healthcheck.Retries = 3

	handleResult(c.State.StartedAt.Add(20*time.Second), 1)
	handleResult(c.State.StartedAt.Add(40*time.Second), 1)
	if status := c.State.Health.Status(); status != types.Starting ***REMOVED***
		t.Errorf("Expecting starting, but got %#v\n", status)
	***REMOVED***
	if c.State.Health.FailingStreak != 2 ***REMOVED***
		t.Errorf("Expecting FailingStreak=2, but got %d\n", c.State.Health.FailingStreak)
	***REMOVED***
	handleResult(c.State.StartedAt.Add(60*time.Second), 1)
	expect("health_status: unhealthy")

	handleResult(c.State.StartedAt.Add(80*time.Second), 0)
	expect("health_status: healthy")
	if c.State.Health.FailingStreak != 0 ***REMOVED***
		t.Errorf("Expecting FailingStreak=0, but got %d\n", c.State.Health.FailingStreak)
	***REMOVED***

	// Test start period

	reset(c)
	c.Config.Healthcheck.Retries = 2
	c.Config.Healthcheck.StartPeriod = 30 * time.Second

	handleResult(c.State.StartedAt.Add(20*time.Second), 1)
	if status := c.State.Health.Status(); status != types.Starting ***REMOVED***
		t.Errorf("Expecting starting, but got %#v\n", status)
	***REMOVED***
	if c.State.Health.FailingStreak != 0 ***REMOVED***
		t.Errorf("Expecting FailingStreak=0, but got %d\n", c.State.Health.FailingStreak)
	***REMOVED***
	handleResult(c.State.StartedAt.Add(50*time.Second), 1)
	if status := c.State.Health.Status(); status != types.Starting ***REMOVED***
		t.Errorf("Expecting starting, but got %#v\n", status)
	***REMOVED***
	if c.State.Health.FailingStreak != 1 ***REMOVED***
		t.Errorf("Expecting FailingStreak=1, but got %d\n", c.State.Health.FailingStreak)
	***REMOVED***
	handleResult(c.State.StartedAt.Add(80*time.Second), 0)
	expect("health_status: healthy")
	if c.State.Health.FailingStreak != 0 ***REMOVED***
		t.Errorf("Expecting FailingStreak=0, but got %d\n", c.State.Health.FailingStreak)
	***REMOVED***
***REMOVED***
