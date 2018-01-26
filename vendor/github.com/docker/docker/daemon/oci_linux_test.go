package daemon

import (
	"testing"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/container"
	"github.com/docker/docker/daemon/config"
	"github.com/docker/docker/oci"
	"github.com/docker/docker/pkg/idtools"

	"github.com/stretchr/testify/assert"
)

// TestTmpfsDevShmNoDupMount checks that a user-specified /dev/shm tmpfs
// mount (as in "docker run --tmpfs /dev/shm:rw,size=NNN") does not result
// in "Duplicate mount point" error from the engine.
// https://github.com/moby/moby/issues/35455
func TestTmpfsDevShmNoDupMount(t *testing.T) ***REMOVED***
	d := Daemon***REMOVED***
		// some empty structs to avoid getting a panic
		// caused by a null pointer dereference
		idMappings:  &idtools.IDMappings***REMOVED******REMOVED***,
		configStore: &config.Config***REMOVED******REMOVED***,
	***REMOVED***
	c := &container.Container***REMOVED***
		ShmPath: "foobar", // non-empty, for c.IpcMounts() to work
		HostConfig: &containertypes.HostConfig***REMOVED***
			IpcMode: containertypes.IpcMode("shareable"), // default mode
			// --tmpfs /dev/shm:rw,exec,size=NNN
			Tmpfs: map[string]string***REMOVED***
				"/dev/shm": "rw,exec,size=1g",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	// Mimick the code flow of daemon.createSpec(), enough to reproduce the issue
	ms, err := d.setupMounts(c)
	assert.NoError(t, err)

	ms = append(ms, c.IpcMounts()...)

	tmpfsMounts, err := c.TmpfsMounts()
	assert.NoError(t, err)
	ms = append(ms, tmpfsMounts...)

	s := oci.DefaultSpec()
	err = setMounts(&d, &s, c, ms)
	assert.NoError(t, err)
***REMOVED***
