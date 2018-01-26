package container

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/integration/util/request"
	"github.com/docker/docker/pkg/stdcopy"
)

func TestUpdateCPUQUota(t *testing.T) ***REMOVED***
	t.Parallel()

	client := request.NewAPIClient(t)
	ctx := context.Background()

	c, err := client.ContainerCreate(ctx, &container.Config***REMOVED***
		Image: "busybox",
		Cmd:   []string***REMOVED***"top"***REMOVED***,
	***REMOVED***, nil, nil, "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer func() ***REMOVED***
		if err := client.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions***REMOVED***Force: true***REMOVED***); err != nil ***REMOVED***
			panic(fmt.Sprintf("failed to clean up after test: %v", err))
		***REMOVED***
	***REMOVED***()

	if err := client.ContainerStart(ctx, c.ID, types.ContainerStartOptions***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	for _, test := range []struct ***REMOVED***
		desc   string
		update int64
	***REMOVED******REMOVED***
		***REMOVED***desc: "some random value", update: 15000***REMOVED***,
		***REMOVED***desc: "a higher value", update: 20000***REMOVED***,
		***REMOVED***desc: "a lower value", update: 10000***REMOVED***,
		***REMOVED***desc: "unset value", update: -1***REMOVED***,
	***REMOVED*** ***REMOVED***
		if _, err := client.ContainerUpdate(ctx, c.ID, container.UpdateConfig***REMOVED***
			Resources: container.Resources***REMOVED***
				CPUQuota: test.update,
			***REMOVED***,
		***REMOVED***); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		inspect, err := client.ContainerInspect(ctx, c.ID)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		if inspect.HostConfig.CPUQuota != test.update ***REMOVED***
			t.Fatalf("quota not updated in the API, expected %d, got: %d", test.update, inspect.HostConfig.CPUQuota)
		***REMOVED***

		execCreate, err := client.ContainerExecCreate(ctx, c.ID, types.ExecConfig***REMOVED***
			Cmd:          []string***REMOVED***"/bin/cat", "/sys/fs/cgroup/cpu/cpu.cfs_quota_us"***REMOVED***,
			AttachStdout: true,
			AttachStderr: true,
		***REMOVED***)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		attach, err := client.ContainerExecAttach(ctx, execCreate.ID, types.ExecStartCheck***REMOVED******REMOVED***)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		if err := client.ContainerExecStart(ctx, execCreate.ID, types.ExecStartCheck***REMOVED******REMOVED***); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		buf := bytes.NewBuffer(nil)
		ready := make(chan error)

		go func() ***REMOVED***
			_, err := stdcopy.StdCopy(buf, buf, attach.Reader)
			ready <- err
		***REMOVED***()

		select ***REMOVED***
		case <-time.After(60 * time.Second):
			t.Fatal("timeout waiting for exec to complete")
		case err := <-ready:
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
		***REMOVED***

		actual := strings.TrimSpace(buf.String())
		if actual != strconv.Itoa(int(test.update)) ***REMOVED***
			t.Fatalf("expected cgroup value %d, got: %s", test.update, actual)
		***REMOVED***
	***REMOVED***

***REMOVED***
