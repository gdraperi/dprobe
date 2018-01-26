package container

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/docker/docker/daemon"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/swarmkit/api"
)

func newTestControllerWithMount(m api.Mount) (*controller, error) ***REMOVED***
	return newController(&daemon.Daemon***REMOVED******REMOVED***, &api.Task***REMOVED***
		ID:        stringid.GenerateRandomID(),
		ServiceID: stringid.GenerateRandomID(),
		Spec: api.TaskSpec***REMOVED***
			Runtime: &api.TaskSpec_Container***REMOVED***
				Container: &api.ContainerSpec***REMOVED***
					Image: "image_name",
					Labels: map[string]string***REMOVED***
						"com.docker.swarm.task.id": "id",
					***REMOVED***,
					Mounts: []api.Mount***REMOVED***m***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***, nil,
		nil)
***REMOVED***

func TestControllerValidateMountBind(t *testing.T) ***REMOVED***
	// with improper source
	if _, err := newTestControllerWithMount(api.Mount***REMOVED***
		Type:   api.MountTypeBind,
		Source: "foo",
		Target: testAbsPath,
	***REMOVED***); err == nil || !strings.Contains(err.Error(), "invalid bind mount source") ***REMOVED***
		t.Fatalf("expected  error, got: %v", err)
	***REMOVED***

	// with non-existing source
	if _, err := newTestControllerWithMount(api.Mount***REMOVED***
		Type:   api.MountTypeBind,
		Source: testAbsNonExistent,
		Target: testAbsPath,
	***REMOVED***); err != nil ***REMOVED***
		t.Fatalf("controller should not error at creation: %v", err)
	***REMOVED***

	// with proper source
	tmpdir, err := ioutil.TempDir("", "TestControllerValidateMountBind")
	if err != nil ***REMOVED***
		t.Fatalf("failed to create temp dir: %v", err)
	***REMOVED***
	defer os.Remove(tmpdir)

	if _, err := newTestControllerWithMount(api.Mount***REMOVED***
		Type:   api.MountTypeBind,
		Source: tmpdir,
		Target: testAbsPath,
	***REMOVED***); err != nil ***REMOVED***
		t.Fatalf("expected  error, got: %v", err)
	***REMOVED***
***REMOVED***

func TestControllerValidateMountVolume(t *testing.T) ***REMOVED***
	// with improper source
	if _, err := newTestControllerWithMount(api.Mount***REMOVED***
		Type:   api.MountTypeVolume,
		Source: testAbsPath,
		Target: testAbsPath,
	***REMOVED***); err == nil || !strings.Contains(err.Error(), "invalid volume mount source") ***REMOVED***
		t.Fatalf("expected error, got: %v", err)
	***REMOVED***

	// with proper source
	if _, err := newTestControllerWithMount(api.Mount***REMOVED***
		Type:   api.MountTypeVolume,
		Source: "foo",
		Target: testAbsPath,
	***REMOVED***); err != nil ***REMOVED***
		t.Fatalf("expected error, got: %v", err)
	***REMOVED***
***REMOVED***

func TestControllerValidateMountTarget(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "TestControllerValidateMountTarget")
	if err != nil ***REMOVED***
		t.Fatalf("failed to create temp dir: %v", err)
	***REMOVED***
	defer os.Remove(tmpdir)

	// with improper target
	if _, err := newTestControllerWithMount(api.Mount***REMOVED***
		Type:   api.MountTypeBind,
		Source: testAbsPath,
		Target: "foo",
	***REMOVED***); err == nil || !strings.Contains(err.Error(), "invalid mount target") ***REMOVED***
		t.Fatalf("expected error, got: %v", err)
	***REMOVED***

	// with proper target
	if _, err := newTestControllerWithMount(api.Mount***REMOVED***
		Type:   api.MountTypeBind,
		Source: tmpdir,
		Target: testAbsPath,
	***REMOVED***); err != nil ***REMOVED***
		t.Fatalf("expected no error, got: %v", err)
	***REMOVED***
***REMOVED***

func TestControllerValidateMountTmpfs(t *testing.T) ***REMOVED***
	// with improper target
	if _, err := newTestControllerWithMount(api.Mount***REMOVED***
		Type:   api.MountTypeTmpfs,
		Source: "foo",
		Target: testAbsPath,
	***REMOVED***); err == nil || !strings.Contains(err.Error(), "invalid tmpfs source") ***REMOVED***
		t.Fatalf("expected error, got: %v", err)
	***REMOVED***

	// with proper target
	if _, err := newTestControllerWithMount(api.Mount***REMOVED***
		Type:   api.MountTypeTmpfs,
		Target: testAbsPath,
	***REMOVED***); err != nil ***REMOVED***
		t.Fatalf("expected no error, got: %v", err)
	***REMOVED***
***REMOVED***

func TestControllerValidateMountInvalidType(t *testing.T) ***REMOVED***
	// with improper target
	if _, err := newTestControllerWithMount(api.Mount***REMOVED***
		Type:   api.Mount_MountType(9999),
		Source: "foo",
		Target: testAbsPath,
	***REMOVED***); err == nil || !strings.Contains(err.Error(), "invalid mount type") ***REMOVED***
		t.Fatalf("expected error, got: %v", err)
	***REMOVED***
***REMOVED***
