package containerd

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/containerd/containerd/api/services/tasks/v1"
	"github.com/containerd/containerd/api/types"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/typeurl"
	prototypes "github.com/gogo/protobuf/types"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
)

// Container is a metadata object for container resources and task creation
type Container interface ***REMOVED***
	// ID identifies the container
	ID() string
	// Info returns the underlying container record type
	Info(context.Context) (containers.Container, error)
	// Delete removes the container
	Delete(context.Context, ...DeleteOpts) error
	// NewTask creates a new task based on the container metadata
	NewTask(context.Context, cio.Creator, ...NewTaskOpts) (Task, error)
	// Spec returns the OCI runtime specification
	Spec(context.Context) (*specs.Spec, error)
	// Task returns the current task for the container
	//
	// If cio.Attach options are passed the client will reattach to the IO for the running
	// task. If no task exists for the container a NotFound error is returned
	//
	// Clients must make sure that only one reader is attached to the task and consuming
	// the output from the task's fifos
	Task(context.Context, cio.Attach) (Task, error)
	// Image returns the image that the container is based on
	Image(context.Context) (Image, error)
	// Labels returns the labels set on the container
	Labels(context.Context) (map[string]string, error)
	// SetLabels sets the provided labels for the container and returns the final label set
	SetLabels(context.Context, map[string]string) (map[string]string, error)
	// Extensions returns the extensions set on the container
	Extensions(context.Context) (map[string]prototypes.Any, error)
	// Update a container
	Update(context.Context, ...UpdateContainerOpts) error
***REMOVED***

func containerFromRecord(client *Client, c containers.Container) *container ***REMOVED***
	return &container***REMOVED***
		client: client,
		id:     c.ID,
	***REMOVED***
***REMOVED***

var _ = (Container)(&container***REMOVED******REMOVED***)

type container struct ***REMOVED***
	client *Client
	id     string
***REMOVED***

// ID returns the container's unique id
func (c *container) ID() string ***REMOVED***
	return c.id
***REMOVED***

func (c *container) Info(ctx context.Context) (containers.Container, error) ***REMOVED***
	return c.get(ctx)
***REMOVED***

func (c *container) Extensions(ctx context.Context) (map[string]prototypes.Any, error) ***REMOVED***
	r, err := c.get(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return r.Extensions, nil
***REMOVED***

func (c *container) Labels(ctx context.Context) (map[string]string, error) ***REMOVED***
	r, err := c.get(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return r.Labels, nil
***REMOVED***

func (c *container) SetLabels(ctx context.Context, labels map[string]string) (map[string]string, error) ***REMOVED***
	container := containers.Container***REMOVED***
		ID:     c.id,
		Labels: labels,
	***REMOVED***

	var paths []string
	// mask off paths so we only muck with the labels encountered in labels.
	// Labels not in the passed in argument will be left alone.
	for k := range labels ***REMOVED***
		paths = append(paths, strings.Join([]string***REMOVED***"labels", k***REMOVED***, "."))
	***REMOVED***

	r, err := c.client.ContainerService().Update(ctx, container, paths...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return r.Labels, nil
***REMOVED***

// Spec returns the current OCI specification for the container
func (c *container) Spec(ctx context.Context) (*specs.Spec, error) ***REMOVED***
	r, err := c.get(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var s specs.Spec
	if err := json.Unmarshal(r.Spec.Value, &s); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &s, nil
***REMOVED***

// Delete deletes an existing container
// an error is returned if the container has running tasks
func (c *container) Delete(ctx context.Context, opts ...DeleteOpts) error ***REMOVED***
	if _, err := c.loadTask(ctx, nil); err == nil ***REMOVED***
		return errors.Wrapf(errdefs.ErrFailedPrecondition, "cannot delete running task %v", c.id)
	***REMOVED***
	r, err := c.get(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, o := range opts ***REMOVED***
		if err := o(ctx, c.client, r); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return c.client.ContainerService().Delete(ctx, c.id)
***REMOVED***

func (c *container) Task(ctx context.Context, attach cio.Attach) (Task, error) ***REMOVED***
	return c.loadTask(ctx, attach)
***REMOVED***

// Image returns the image that the container is based on
func (c *container) Image(ctx context.Context) (Image, error) ***REMOVED***
	r, err := c.get(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if r.Image == "" ***REMOVED***
		return nil, errors.Wrap(errdefs.ErrNotFound, "container not created from an image")
	***REMOVED***
	i, err := c.client.ImageService().Get(ctx, r.Image)
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "failed to get image %s for container", r.Image)
	***REMOVED***
	return &image***REMOVED***
		client: c.client,
		i:      i,
	***REMOVED***, nil
***REMOVED***

func (c *container) NewTask(ctx context.Context, ioCreate cio.Creator, opts ...NewTaskOpts) (_ Task, err error) ***REMOVED***
	i, err := ioCreate(c.id)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil && i != nil ***REMOVED***
			i.Cancel()
			i.Close()
		***REMOVED***
	***REMOVED***()
	cfg := i.Config()
	request := &tasks.CreateTaskRequest***REMOVED***
		ContainerID: c.id,
		Terminal:    cfg.Terminal,
		Stdin:       cfg.Stdin,
		Stdout:      cfg.Stdout,
		Stderr:      cfg.Stderr,
	***REMOVED***
	r, err := c.get(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if r.SnapshotKey != "" ***REMOVED***
		if r.Snapshotter == "" ***REMOVED***
			return nil, errors.Wrapf(errdefs.ErrInvalidArgument, "unable to resolve rootfs mounts without snapshotter on container")
		***REMOVED***

		// get the rootfs from the snapshotter and add it to the request
		mounts, err := c.client.SnapshotService(r.Snapshotter).Mounts(ctx, r.SnapshotKey)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		for _, m := range mounts ***REMOVED***
			request.Rootfs = append(request.Rootfs, &types.Mount***REMOVED***
				Type:    m.Type,
				Source:  m.Source,
				Options: m.Options,
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	var info TaskInfo
	for _, o := range opts ***REMOVED***
		if err := o(ctx, c.client, &info); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	if info.RootFS != nil ***REMOVED***
		for _, m := range info.RootFS ***REMOVED***
			request.Rootfs = append(request.Rootfs, &types.Mount***REMOVED***
				Type:    m.Type,
				Source:  m.Source,
				Options: m.Options,
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	if info.Options != nil ***REMOVED***
		any, err := typeurl.MarshalAny(info.Options)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		request.Options = any
	***REMOVED***
	t := &task***REMOVED***
		client: c.client,
		io:     i,
		id:     c.id,
	***REMOVED***
	if info.Checkpoint != nil ***REMOVED***
		request.Checkpoint = info.Checkpoint
	***REMOVED***
	response, err := c.client.TaskService().Create(ctx, request)
	if err != nil ***REMOVED***
		return nil, errdefs.FromGRPC(err)
	***REMOVED***
	t.pid = response.Pid
	return t, nil
***REMOVED***

func (c *container) Update(ctx context.Context, opts ...UpdateContainerOpts) error ***REMOVED***
	// fetch the current container config before updating it
	r, err := c.get(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, o := range opts ***REMOVED***
		if err := o(ctx, c.client, &r); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if _, err := c.client.ContainerService().Update(ctx, r); err != nil ***REMOVED***
		return errdefs.FromGRPC(err)
	***REMOVED***
	return nil
***REMOVED***

func (c *container) loadTask(ctx context.Context, ioAttach cio.Attach) (Task, error) ***REMOVED***
	response, err := c.client.TaskService().Get(ctx, &tasks.GetRequest***REMOVED***
		ContainerID: c.id,
	***REMOVED***)
	if err != nil ***REMOVED***
		err = errdefs.FromGRPC(err)
		if errdefs.IsNotFound(err) ***REMOVED***
			return nil, errors.Wrapf(err, "no running task found")
		***REMOVED***
		return nil, err
	***REMOVED***
	var i cio.IO
	if ioAttach != nil ***REMOVED***
		if i, err = attachExistingIO(response, ioAttach); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	t := &task***REMOVED***
		client: c.client,
		io:     i,
		id:     response.Process.ID,
		pid:    response.Process.Pid,
	***REMOVED***
	return t, nil
***REMOVED***

func (c *container) get(ctx context.Context) (containers.Container, error) ***REMOVED***
	return c.client.ContainerService().Get(ctx, c.id)
***REMOVED***

// get the existing fifo paths from the task information stored by the daemon
func attachExistingIO(response *tasks.GetResponse, ioAttach cio.Attach) (cio.IO, error) ***REMOVED***
	path := getFifoDir([]string***REMOVED***
		response.Process.Stdin,
		response.Process.Stdout,
		response.Process.Stderr,
	***REMOVED***)
	closer := func() error ***REMOVED***
		return os.RemoveAll(path)
	***REMOVED***
	fifoSet := cio.NewFIFOSet(cio.Config***REMOVED***
		Stdin:    response.Process.Stdin,
		Stdout:   response.Process.Stdout,
		Stderr:   response.Process.Stderr,
		Terminal: response.Process.Terminal,
	***REMOVED***, closer)
	return ioAttach(fifoSet)
***REMOVED***

// getFifoDir looks for any non-empty path for a stdio fifo
// and returns the dir for where it is located
func getFifoDir(paths []string) string ***REMOVED***
	for _, p := range paths ***REMOVED***
		if p != "" ***REMOVED***
			return filepath.Dir(p)
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***
