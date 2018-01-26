package containerd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	goruntime "runtime"
	"strings"
	"syscall"
	"time"

	"github.com/containerd/containerd/api/services/tasks/v1"
	"github.com/containerd/containerd/api/types"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/diff"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/mount"
	"github.com/containerd/containerd/plugin"
	"github.com/containerd/containerd/rootfs"
	"github.com/containerd/typeurl"
	google_protobuf "github.com/gogo/protobuf/types"
	digest "github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/specs-go/v1"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
)

// UnknownExitStatus is returned when containerd is unable to
// determine the exit status of a process. This can happen if the process never starts
// or if an error was encountered when obtaining the exit status, it is set to 255.
const UnknownExitStatus = 255

const (
	checkpointDateFormat = "01-02-2006-15:04:05"
	checkpointNameFormat = "containerd.io/checkpoint/%s:%s"
)

// Status returns process status and exit information
type Status struct ***REMOVED***
	// Status of the process
	Status ProcessStatus
	// ExitStatus returned by the process
	ExitStatus uint32
	// ExitedTime is the time at which the process died
	ExitTime time.Time
***REMOVED***

// ProcessInfo provides platform specific process information
type ProcessInfo struct ***REMOVED***
	// Pid is the process ID
	Pid uint32
	// Info includes additional process information
	// Info varies by platform
	Info *google_protobuf.Any
***REMOVED***

// ProcessStatus returns a human readable status for the Process representing its current status
type ProcessStatus string

const (
	// Running indicates the process is currently executing
	Running ProcessStatus = "running"
	// Created indicates the process has been created within containerd but the
	// user's defined process has not started
	Created ProcessStatus = "created"
	// Stopped indicates that the process has ran and exited
	Stopped ProcessStatus = "stopped"
	// Paused indicates that the process is currently paused
	Paused ProcessStatus = "paused"
	// Pausing indicates that the process is currently switching from a
	// running state into a paused state
	Pausing ProcessStatus = "pausing"
	// Unknown indicates that we could not determine the status from the runtime
	Unknown ProcessStatus = "unknown"
)

// IOCloseInfo allows specific io pipes to be closed on a process
type IOCloseInfo struct ***REMOVED***
	Stdin bool
***REMOVED***

// IOCloserOpts allows the caller to set specific pipes as closed on a process
type IOCloserOpts func(*IOCloseInfo)

// WithStdinCloser closes the stdin of a process
func WithStdinCloser(r *IOCloseInfo) ***REMOVED***
	r.Stdin = true
***REMOVED***

// CheckpointTaskInfo allows specific checkpoint information to be set for the task
type CheckpointTaskInfo struct ***REMOVED***
	Name string
	// ParentCheckpoint is the digest of a parent checkpoint
	ParentCheckpoint digest.Digest
	// Options hold runtime specific settings for checkpointing a task
	Options interface***REMOVED******REMOVED***
***REMOVED***

// CheckpointTaskOpts allows the caller to set checkpoint options
type CheckpointTaskOpts func(*CheckpointTaskInfo) error

// TaskInfo sets options for task creation
type TaskInfo struct ***REMOVED***
	// Checkpoint is the Descriptor for an existing checkpoint that can be used
	// to restore a task's runtime and memory state
	Checkpoint *types.Descriptor
	// RootFS is a list of mounts to use as the task's root filesystem
	RootFS []mount.Mount
	// Options hold runtime specific settings for task creation
	Options interface***REMOVED******REMOVED***
***REMOVED***

// Task is the executable object within containerd
type Task interface ***REMOVED***
	Process

	// Pause suspends the execution of the task
	Pause(context.Context) error
	// Resume the execution of the task
	Resume(context.Context) error
	// Exec creates a new process inside the task
	Exec(context.Context, string, *specs.Process, cio.Creator) (Process, error)
	// Pids returns a list of system specific process ids inside the task
	Pids(context.Context) ([]ProcessInfo, error)
	// Checkpoint serializes the runtime and memory information of a task into an
	// OCI Index that can be push and pulled from a remote resource.
	//
	// Additional software like CRIU maybe required to checkpoint and restore tasks
	Checkpoint(context.Context, ...CheckpointTaskOpts) (Image, error)
	// Update modifies executing tasks with updated settings
	Update(context.Context, ...UpdateTaskOpts) error
	// LoadProcess loads a previously created exec'd process
	LoadProcess(context.Context, string, cio.Attach) (Process, error)
	// Metrics returns task metrics for runtime specific metrics
	//
	// The metric types are generic to containerd and change depending on the runtime
	// For the built in Linux runtime, github.com/containerd/cgroups.Metrics
	// are returned in protobuf format
	Metrics(context.Context) (*types.Metric, error)
***REMOVED***

var _ = (Task)(&task***REMOVED******REMOVED***)

type task struct ***REMOVED***
	client *Client

	io  cio.IO
	id  string
	pid uint32
***REMOVED***

// Pid returns the pid or process id for the task
func (t *task) Pid() uint32 ***REMOVED***
	return t.pid
***REMOVED***

func (t *task) Start(ctx context.Context) error ***REMOVED***
	r, err := t.client.TaskService().Start(ctx, &tasks.StartRequest***REMOVED***
		ContainerID: t.id,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.io.Cancel()
		t.io.Close()
		return errdefs.FromGRPC(err)
	***REMOVED***
	t.pid = r.Pid
	return nil
***REMOVED***

func (t *task) Kill(ctx context.Context, s syscall.Signal, opts ...KillOpts) error ***REMOVED***
	var i KillInfo
	for _, o := range opts ***REMOVED***
		if err := o(ctx, &i); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	_, err := t.client.TaskService().Kill(ctx, &tasks.KillRequest***REMOVED***
		Signal:      uint32(s),
		ContainerID: t.id,
		ExecID:      i.ExecID,
		All:         i.All,
	***REMOVED***)
	if err != nil ***REMOVED***
		return errdefs.FromGRPC(err)
	***REMOVED***
	return nil
***REMOVED***

func (t *task) Pause(ctx context.Context) error ***REMOVED***
	_, err := t.client.TaskService().Pause(ctx, &tasks.PauseTaskRequest***REMOVED***
		ContainerID: t.id,
	***REMOVED***)
	return errdefs.FromGRPC(err)
***REMOVED***

func (t *task) Resume(ctx context.Context) error ***REMOVED***
	_, err := t.client.TaskService().Resume(ctx, &tasks.ResumeTaskRequest***REMOVED***
		ContainerID: t.id,
	***REMOVED***)
	return errdefs.FromGRPC(err)
***REMOVED***

func (t *task) Status(ctx context.Context) (Status, error) ***REMOVED***
	r, err := t.client.TaskService().Get(ctx, &tasks.GetRequest***REMOVED***
		ContainerID: t.id,
	***REMOVED***)
	if err != nil ***REMOVED***
		return Status***REMOVED******REMOVED***, errdefs.FromGRPC(err)
	***REMOVED***
	return Status***REMOVED***
		Status:     ProcessStatus(strings.ToLower(r.Process.Status.String())),
		ExitStatus: r.Process.ExitStatus,
		ExitTime:   r.Process.ExitedAt,
	***REMOVED***, nil
***REMOVED***

func (t *task) Wait(ctx context.Context) (<-chan ExitStatus, error) ***REMOVED***
	c := make(chan ExitStatus, 1)
	go func() ***REMOVED***
		defer close(c)
		r, err := t.client.TaskService().Wait(ctx, &tasks.WaitRequest***REMOVED***
			ContainerID: t.id,
		***REMOVED***)
		if err != nil ***REMOVED***
			c <- ExitStatus***REMOVED***
				code: UnknownExitStatus,
				err:  err,
			***REMOVED***
			return
		***REMOVED***
		c <- ExitStatus***REMOVED***
			code:     r.ExitStatus,
			exitedAt: r.ExitedAt,
		***REMOVED***
	***REMOVED***()
	return c, nil
***REMOVED***

// Delete deletes the task and its runtime state
// it returns the exit status of the task and any errors that were encountered
// during cleanup
func (t *task) Delete(ctx context.Context, opts ...ProcessDeleteOpts) (*ExitStatus, error) ***REMOVED***
	for _, o := range opts ***REMOVED***
		if err := o(ctx, t); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	status, err := t.Status(ctx)
	if err != nil && errdefs.IsNotFound(err) ***REMOVED***
		return nil, err
	***REMOVED***
	switch status.Status ***REMOVED***
	case Stopped, Unknown, "":
	case Created:
		if t.client.runtime == fmt.Sprintf("%s.%s", plugin.RuntimePlugin, "windows") ***REMOVED***
			// On windows Created is akin to Stopped
			break
		***REMOVED***
		fallthrough
	default:
		return nil, errors.Wrapf(errdefs.ErrFailedPrecondition, "task must be stopped before deletion: %s", status.Status)
	***REMOVED***
	if t.io != nil ***REMOVED***
		t.io.Cancel()
		t.io.Wait()
		t.io.Close()
	***REMOVED***
	r, err := t.client.TaskService().Delete(ctx, &tasks.DeleteTaskRequest***REMOVED***
		ContainerID: t.id,
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, errdefs.FromGRPC(err)
	***REMOVED***
	return &ExitStatus***REMOVED***code: r.ExitStatus, exitedAt: r.ExitedAt***REMOVED***, nil
***REMOVED***

func (t *task) Exec(ctx context.Context, id string, spec *specs.Process, ioCreate cio.Creator) (_ Process, err error) ***REMOVED***
	if id == "" ***REMOVED***
		return nil, errors.Wrapf(errdefs.ErrInvalidArgument, "exec id must not be empty")
	***REMOVED***
	i, err := ioCreate(id)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil && i != nil ***REMOVED***
			i.Cancel()
			i.Close()
		***REMOVED***
	***REMOVED***()
	any, err := typeurl.MarshalAny(spec)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	cfg := i.Config()
	request := &tasks.ExecProcessRequest***REMOVED***
		ContainerID: t.id,
		ExecID:      id,
		Terminal:    cfg.Terminal,
		Stdin:       cfg.Stdin,
		Stdout:      cfg.Stdout,
		Stderr:      cfg.Stderr,
		Spec:        any,
	***REMOVED***
	if _, err := t.client.TaskService().Exec(ctx, request); err != nil ***REMOVED***
		i.Cancel()
		i.Wait()
		i.Close()
		return nil, errdefs.FromGRPC(err)
	***REMOVED***
	return &process***REMOVED***
		id:   id,
		task: t,
		io:   i,
	***REMOVED***, nil
***REMOVED***

func (t *task) Pids(ctx context.Context) ([]ProcessInfo, error) ***REMOVED***
	response, err := t.client.TaskService().ListPids(ctx, &tasks.ListPidsRequest***REMOVED***
		ContainerID: t.id,
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, errdefs.FromGRPC(err)
	***REMOVED***
	var processList []ProcessInfo
	for _, p := range response.Processes ***REMOVED***
		processList = append(processList, ProcessInfo***REMOVED***
			Pid:  p.Pid,
			Info: p.Info,
		***REMOVED***)
	***REMOVED***
	return processList, nil
***REMOVED***

func (t *task) CloseIO(ctx context.Context, opts ...IOCloserOpts) error ***REMOVED***
	r := &tasks.CloseIORequest***REMOVED***
		ContainerID: t.id,
	***REMOVED***
	var i IOCloseInfo
	for _, o := range opts ***REMOVED***
		o(&i)
	***REMOVED***
	r.Stdin = i.Stdin
	_, err := t.client.TaskService().CloseIO(ctx, r)
	return errdefs.FromGRPC(err)
***REMOVED***

func (t *task) IO() cio.IO ***REMOVED***
	return t.io
***REMOVED***

func (t *task) Resize(ctx context.Context, w, h uint32) error ***REMOVED***
	_, err := t.client.TaskService().ResizePty(ctx, &tasks.ResizePtyRequest***REMOVED***
		ContainerID: t.id,
		Width:       w,
		Height:      h,
	***REMOVED***)
	return errdefs.FromGRPC(err)
***REMOVED***

func (t *task) Checkpoint(ctx context.Context, opts ...CheckpointTaskOpts) (Image, error) ***REMOVED***
	ctx, done, err := t.client.WithLease(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer done()

	request := &tasks.CheckpointTaskRequest***REMOVED***
		ContainerID: t.id,
	***REMOVED***
	var i CheckpointTaskInfo
	for _, o := range opts ***REMOVED***
		if err := o(&i); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	// set a default name
	if i.Name == "" ***REMOVED***
		i.Name = fmt.Sprintf(checkpointNameFormat, t.id, time.Now().Format(checkpointDateFormat))
	***REMOVED***
	request.ParentCheckpoint = i.ParentCheckpoint
	if i.Options != nil ***REMOVED***
		any, err := typeurl.MarshalAny(i.Options)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		request.Options = any
	***REMOVED***
	// make sure we pause it and resume after all other filesystem operations are completed
	if err := t.Pause(ctx); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer t.Resume(ctx)
	cr, err := t.client.ContainerService().Get(ctx, t.id)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	index := v1.Index***REMOVED***
		Annotations: make(map[string]string),
	***REMOVED***
	if err := t.checkpointTask(ctx, &index, request); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if cr.Image != "" ***REMOVED***
		if err := t.checkpointImage(ctx, &index, cr.Image); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		index.Annotations["image.name"] = cr.Image
	***REMOVED***
	if cr.SnapshotKey != "" ***REMOVED***
		if err := t.checkpointRWSnapshot(ctx, &index, cr.Snapshotter, cr.SnapshotKey); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	desc, err := t.writeIndex(ctx, &index)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	im := images.Image***REMOVED***
		Name:   i.Name,
		Target: desc,
		Labels: map[string]string***REMOVED***
			"containerd.io/checkpoint": "true",
		***REMOVED***,
	***REMOVED***
	if im, err = t.client.ImageService().Create(ctx, im); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &image***REMOVED***
		client: t.client,
		i:      im,
	***REMOVED***, nil
***REMOVED***

// UpdateTaskInfo allows updated specific settings to be changed on a task
type UpdateTaskInfo struct ***REMOVED***
	// Resources updates a tasks resource constraints
	Resources interface***REMOVED******REMOVED***
***REMOVED***

// UpdateTaskOpts allows a caller to update task settings
type UpdateTaskOpts func(context.Context, *Client, *UpdateTaskInfo) error

func (t *task) Update(ctx context.Context, opts ...UpdateTaskOpts) error ***REMOVED***
	request := &tasks.UpdateTaskRequest***REMOVED***
		ContainerID: t.id,
	***REMOVED***
	var i UpdateTaskInfo
	for _, o := range opts ***REMOVED***
		if err := o(ctx, t.client, &i); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if i.Resources != nil ***REMOVED***
		any, err := typeurl.MarshalAny(i.Resources)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		request.Resources = any
	***REMOVED***
	_, err := t.client.TaskService().Update(ctx, request)
	return errdefs.FromGRPC(err)
***REMOVED***

func (t *task) LoadProcess(ctx context.Context, id string, ioAttach cio.Attach) (Process, error) ***REMOVED***
	response, err := t.client.TaskService().Get(ctx, &tasks.GetRequest***REMOVED***
		ContainerID: t.id,
		ExecID:      id,
	***REMOVED***)
	if err != nil ***REMOVED***
		err = errdefs.FromGRPC(err)
		if errdefs.IsNotFound(err) ***REMOVED***
			return nil, errors.Wrapf(err, "no running process found")
		***REMOVED***
		return nil, err
	***REMOVED***
	var i cio.IO
	if ioAttach != nil ***REMOVED***
		if i, err = attachExistingIO(response, ioAttach); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return &process***REMOVED***
		id:   id,
		task: t,
		io:   i,
	***REMOVED***, nil
***REMOVED***

func (t *task) Metrics(ctx context.Context) (*types.Metric, error) ***REMOVED***
	response, err := t.client.TaskService().Metrics(ctx, &tasks.MetricsRequest***REMOVED***
		Filters: []string***REMOVED***
			"id==" + t.id,
		***REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, errdefs.FromGRPC(err)
	***REMOVED***

	if response.Metrics == nil ***REMOVED***
		_, err := t.Status(ctx)
		if err != nil && errdefs.IsNotFound(err) ***REMOVED***
			return nil, err
		***REMOVED***
		return nil, errors.New("no metrics received")
	***REMOVED***

	return response.Metrics[0], nil
***REMOVED***

func (t *task) checkpointTask(ctx context.Context, index *v1.Index, request *tasks.CheckpointTaskRequest) error ***REMOVED***
	response, err := t.client.TaskService().Checkpoint(ctx, request)
	if err != nil ***REMOVED***
		return errdefs.FromGRPC(err)
	***REMOVED***
	// add the checkpoint descriptors to the index
	for _, d := range response.Descriptors ***REMOVED***
		index.Manifests = append(index.Manifests, v1.Descriptor***REMOVED***
			MediaType: d.MediaType,
			Size:      d.Size_,
			Digest:    d.Digest,
			Platform: &v1.Platform***REMOVED***
				OS:           goruntime.GOOS,
				Architecture: goruntime.GOARCH,
			***REMOVED***,
		***REMOVED***)
	***REMOVED***
	return nil
***REMOVED***

func (t *task) checkpointRWSnapshot(ctx context.Context, index *v1.Index, snapshotterName string, id string) error ***REMOVED***
	opts := []diff.Opt***REMOVED***
		diff.WithReference(fmt.Sprintf("checkpoint-rw-%s", id)),
	***REMOVED***
	rw, err := rootfs.Diff(ctx, id, t.client.SnapshotService(snapshotterName), t.client.DiffService(), opts...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	rw.Platform = &v1.Platform***REMOVED***
		OS:           goruntime.GOOS,
		Architecture: goruntime.GOARCH,
	***REMOVED***
	index.Manifests = append(index.Manifests, rw)
	return nil
***REMOVED***

func (t *task) checkpointImage(ctx context.Context, index *v1.Index, image string) error ***REMOVED***
	if image == "" ***REMOVED***
		return fmt.Errorf("cannot checkpoint image with empty name")
	***REMOVED***
	ir, err := t.client.ImageService().Get(ctx, image)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	index.Manifests = append(index.Manifests, ir.Target)
	return nil
***REMOVED***

func (t *task) writeIndex(ctx context.Context, index *v1.Index) (d v1.Descriptor, err error) ***REMOVED***
	labels := map[string]string***REMOVED******REMOVED***
	for i, m := range index.Manifests ***REMOVED***
		labels[fmt.Sprintf("containerd.io/gc.ref.content.%d", i)] = m.Digest.String()
	***REMOVED***
	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(index); err != nil ***REMOVED***
		return v1.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***
	return writeContent(ctx, t.client.ContentStore(), v1.MediaTypeImageIndex, t.id, buf, content.WithLabels(labels))
***REMOVED***

func writeContent(ctx context.Context, store content.Store, mediaType, ref string, r io.Reader, opts ...content.Opt) (d v1.Descriptor, err error) ***REMOVED***
	writer, err := store.Writer(ctx, ref, 0, "")
	if err != nil ***REMOVED***
		return d, err
	***REMOVED***
	defer writer.Close()
	size, err := io.Copy(writer, r)
	if err != nil ***REMOVED***
		return d, err
	***REMOVED***
	if err := writer.Commit(ctx, size, "", opts...); err != nil ***REMOVED***
		return d, err
	***REMOVED***
	return v1.Descriptor***REMOVED***
		MediaType: mediaType,
		Digest:    writer.Digest(),
		Size:      size,
	***REMOVED***, nil
***REMOVED***
