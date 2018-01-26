package container

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/integration-cli/daemon"
)

func TestDaemonRestartKillContainers(t *testing.T) ***REMOVED***
	type testCase struct ***REMOVED***
		desc       string
		config     *container.Config
		hostConfig *container.HostConfig

		xRunning            bool
		xRunningLiveRestore bool
	***REMOVED***

	for _, c := range []testCase***REMOVED***
		***REMOVED***
			desc:                "container without restart policy",
			config:              &container.Config***REMOVED***Image: "busybox", Cmd: []string***REMOVED***"top"***REMOVED******REMOVED***,
			xRunningLiveRestore: true,
		***REMOVED***,
		***REMOVED***
			desc:                "container with restart=always",
			config:              &container.Config***REMOVED***Image: "busybox", Cmd: []string***REMOVED***"top"***REMOVED******REMOVED***,
			hostConfig:          &container.HostConfig***REMOVED***RestartPolicy: container.RestartPolicy***REMOVED***Name: "always"***REMOVED******REMOVED***,
			xRunning:            true,
			xRunningLiveRestore: true,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		for _, liveRestoreEnabled := range []bool***REMOVED***false, true***REMOVED*** ***REMOVED***
			for fnName, stopDaemon := range map[string]func(*testing.T, *daemon.Daemon)***REMOVED***
				"kill-daemon": func(t *testing.T, d *daemon.Daemon) ***REMOVED***
					if err := d.Kill(); err != nil ***REMOVED***
						t.Fatal(err)
					***REMOVED***
				***REMOVED***,
				"stop-daemon": func(t *testing.T, d *daemon.Daemon) ***REMOVED***
					d.Stop(t)
				***REMOVED***,
			***REMOVED*** ***REMOVED***
				t.Run(fmt.Sprintf("live-restore=%v/%s/%s", liveRestoreEnabled, c.desc, fnName), func(t *testing.T) ***REMOVED***
					c := c
					liveRestoreEnabled := liveRestoreEnabled
					stopDaemon := stopDaemon

					t.Parallel()

					d := daemon.New(t, "", "dockerd", daemon.Config***REMOVED******REMOVED***)
					client, err := d.NewClient()
					if err != nil ***REMOVED***
						t.Fatal(err)
					***REMOVED***

					args := []string***REMOVED***"--iptables=false"***REMOVED***
					if liveRestoreEnabled ***REMOVED***
						args = append(args, "--live-restore")
					***REMOVED***

					d.StartWithBusybox(t, args...)
					defer d.Stop(t)
					ctx := context.Background()

					resp, err := client.ContainerCreate(ctx, c.config, c.hostConfig, nil, "")
					if err != nil ***REMOVED***
						t.Fatal(err)
					***REMOVED***
					defer client.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions***REMOVED***Force: true***REMOVED***)

					if err := client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions***REMOVED******REMOVED***); err != nil ***REMOVED***
						t.Fatal(err)
					***REMOVED***

					stopDaemon(t, d)
					d.Start(t, args...)

					expected := c.xRunning
					if liveRestoreEnabled ***REMOVED***
						expected = c.xRunningLiveRestore
					***REMOVED***

					var running bool
					for i := 0; i < 30; i++ ***REMOVED***
						inspect, err := client.ContainerInspect(ctx, resp.ID)
						if err != nil ***REMOVED***
							t.Fatal(err)
						***REMOVED***

						running = inspect.State.Running
						if running == expected ***REMOVED***
							break
						***REMOVED***
						time.Sleep(2 * time.Second)

					***REMOVED***

					if running != expected ***REMOVED***
						t.Fatalf("got unexpected running state, expected %v, got: %v", expected, running)
					***REMOVED***
					// TODO(cpuguy83): test pause states... this seems to be rather undefined currently
				***REMOVED***)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
