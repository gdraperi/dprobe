package container

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/integration-cli/daemon"
	"github.com/docker/docker/pkg/stdcopy"
)

func TestContainerShmNoLeak(t *testing.T) ***REMOVED***
	t.Parallel()
	d := daemon.New(t, "docker", "dockerd", daemon.Config***REMOVED******REMOVED***)
	client, err := d.NewClient()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	d.StartWithBusybox(t)
	defer d.Stop(t)

	ctx := context.Background()
	cfg := container.Config***REMOVED***
		Image: "busybox",
		Cmd:   []string***REMOVED***"top"***REMOVED***,
	***REMOVED***

	ctr, err := client.ContainerCreate(ctx, &cfg, nil, nil, "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer client.ContainerRemove(ctx, ctr.ID, types.ContainerRemoveOptions***REMOVED***Force: true***REMOVED***)

	if err := client.ContainerStart(ctx, ctr.ID, types.ContainerStartOptions***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// this should recursively bind mount everything in the test daemons root
	// except of course we are hoping that the previous containers /dev/shm mount did not leak into this new container
	hc := container.HostConfig***REMOVED***
		Mounts: []mount.Mount***REMOVED***
			***REMOVED***
				Type:        mount.TypeBind,
				Source:      d.Root,
				Target:      "/testdaemonroot",
				BindOptions: &mount.BindOptions***REMOVED***Propagation: mount.PropagationRPrivate***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***
	cfg.Cmd = []string***REMOVED***"/bin/sh", "-c", fmt.Sprintf("mount | grep testdaemonroot | grep containers | grep %s", ctr.ID)***REMOVED***
	cfg.AttachStdout = true
	cfg.AttachStderr = true
	ctrLeak, err := client.ContainerCreate(ctx, &cfg, &hc, nil, "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	attach, err := client.ContainerAttach(ctx, ctrLeak.ID, types.ContainerAttachOptions***REMOVED***
		Stream: true,
		Stdout: true,
		Stderr: true,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := client.ContainerStart(ctx, ctrLeak.ID, types.ContainerStartOptions***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	buf := bytes.NewBuffer(nil)

	if _, err := stdcopy.StdCopy(buf, buf, attach.Reader); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	out := bytes.TrimSpace(buf.Bytes())
	if !bytes.Equal(out, []byte***REMOVED******REMOVED***) ***REMOVED***
		t.Fatalf("mount leaked: %s", string(out))
	***REMOVED***
***REMOVED***
