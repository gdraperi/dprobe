package daemon

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/container"
	"github.com/docker/docker/internal/testutil"
	"github.com/stretchr/testify/require"
)

func newDaemonWithTmpRoot(t *testing.T) (*Daemon, func()) ***REMOVED***
	tmp, err := ioutil.TempDir("", "docker-daemon-unix-test-")
	require.NoError(t, err)
	d := &Daemon***REMOVED***
		repository: tmp,
		root:       tmp,
	***REMOVED***
	d.containers = container.NewMemoryStore()
	return d, func() ***REMOVED*** os.RemoveAll(tmp) ***REMOVED***
***REMOVED***

func newContainerWithState(state *container.State) *container.Container ***REMOVED***
	return &container.Container***REMOVED***
		ID:     "test",
		State:  state,
		Config: &containertypes.Config***REMOVED******REMOVED***,
	***REMOVED***

***REMOVED***

// TestContainerDelete tests that a useful error message and instructions is
// given when attempting to remove a container (#30842)
func TestContainerDelete(t *testing.T) ***REMOVED***
	tt := []struct ***REMOVED***
		errMsg        string
		fixMsg        string
		initContainer func() *container.Container
	***REMOVED******REMOVED***
		// a paused container
		***REMOVED***
			errMsg: "cannot remove a paused container",
			fixMsg: "Unpause and then stop the container before attempting removal or force remove",
			initContainer: func() *container.Container ***REMOVED***
				return newContainerWithState(&container.State***REMOVED***Paused: true, Running: true***REMOVED***)
			***REMOVED******REMOVED***,
		// a restarting container
		***REMOVED***
			errMsg: "cannot remove a restarting container",
			fixMsg: "Stop the container before attempting removal or force remove",
			initContainer: func() *container.Container ***REMOVED***
				c := newContainerWithState(container.NewState())
				c.SetRunning(0, true)
				c.SetRestarting(&container.ExitStatus***REMOVED******REMOVED***)
				return c
			***REMOVED******REMOVED***,
		// a running container
		***REMOVED***
			errMsg: "cannot remove a running container",
			fixMsg: "Stop the container before attempting removal or force remove",
			initContainer: func() *container.Container ***REMOVED***
				return newContainerWithState(&container.State***REMOVED***Running: true***REMOVED***)
			***REMOVED******REMOVED***,
	***REMOVED***

	for _, te := range tt ***REMOVED***
		c := te.initContainer()
		d, cleanup := newDaemonWithTmpRoot(t)
		defer cleanup()
		d.containers.Add(c.ID, c)

		err := d.ContainerRm(c.ID, &types.ContainerRmConfig***REMOVED***ForceRemove: false***REMOVED***)
		testutil.ErrorContains(t, err, te.errMsg)
		testutil.ErrorContains(t, err, te.fixMsg)
	***REMOVED***
***REMOVED***

func TestContainerDoubleDelete(t *testing.T) ***REMOVED***
	c := newContainerWithState(container.NewState())

	// Mark the container as having a delete in progress
	c.SetRemovalInProgress()

	d, cleanup := newDaemonWithTmpRoot(t)
	defer cleanup()
	d.containers.Add(c.ID, c)

	// Try to remove the container when its state is removalInProgress.
	// It should return an error indicating it is under removal progress.
	err := d.ContainerRm(c.ID, &types.ContainerRmConfig***REMOVED***ForceRemove: true***REMOVED***)
	testutil.ErrorContains(t, err, fmt.Sprintf("removal of container %s is already in progress", c.ID))
***REMOVED***
