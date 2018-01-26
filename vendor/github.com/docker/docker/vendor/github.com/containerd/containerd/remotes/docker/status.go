package docker

import (
	"sync"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/errdefs"
	"github.com/pkg/errors"
)

// Status of a content operation
type Status struct ***REMOVED***
	content.Status

	// UploadUUID is used by the Docker registry to reference blob uploads
	UploadUUID string
***REMOVED***

// StatusTracker to track status of operations
type StatusTracker interface ***REMOVED***
	GetStatus(string) (Status, error)
	SetStatus(string, Status)
***REMOVED***

type memoryStatusTracker struct ***REMOVED***
	statuses map[string]Status
	m        sync.Mutex
***REMOVED***

// NewInMemoryTracker returns a StatusTracker that tracks content status in-memory
func NewInMemoryTracker() StatusTracker ***REMOVED***
	return &memoryStatusTracker***REMOVED***
		statuses: map[string]Status***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

func (t *memoryStatusTracker) GetStatus(ref string) (Status, error) ***REMOVED***
	t.m.Lock()
	defer t.m.Unlock()
	status, ok := t.statuses[ref]
	if !ok ***REMOVED***
		return Status***REMOVED******REMOVED***, errors.Wrapf(errdefs.ErrNotFound, "status for ref %v", ref)
	***REMOVED***
	return status, nil
***REMOVED***

func (t *memoryStatusTracker) SetStatus(ref string, status Status) ***REMOVED***
	t.m.Lock()
	t.statuses[ref] = status
	t.m.Unlock()
***REMOVED***
