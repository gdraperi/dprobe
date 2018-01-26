// +build !windows

package libcontainerd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/api/events"
	eventsapi "github.com/containerd/containerd/api/services/events/v1"
	"github.com/containerd/containerd/api/types"
	"github.com/containerd/containerd/archive"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/content"
	containerderrors "github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/linux/runctypes"
	"github.com/containerd/typeurl"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/opencontainers/image-spec/specs-go/v1"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// InitProcessName is the name given to the first process of a
// container
const InitProcessName = "init"

type container struct ***REMOVED***
	mu sync.Mutex

	bundleDir string
	ctr       containerd.Container
	task      containerd.Task
	execs     map[string]containerd.Process
	oomKilled bool
***REMOVED***

func (c *container) setTask(t containerd.Task) ***REMOVED***
	c.mu.Lock()
	c.task = t
	c.mu.Unlock()
***REMOVED***

func (c *container) getTask() containerd.Task ***REMOVED***
	c.mu.Lock()
	t := c.task
	c.mu.Unlock()
	return t
***REMOVED***

func (c *container) addProcess(id string, p containerd.Process) ***REMOVED***
	c.mu.Lock()
	if c.execs == nil ***REMOVED***
		c.execs = make(map[string]containerd.Process)
	***REMOVED***
	c.execs[id] = p
	c.mu.Unlock()
***REMOVED***

func (c *container) deleteProcess(id string) ***REMOVED***
	c.mu.Lock()
	delete(c.execs, id)
	c.mu.Unlock()
***REMOVED***

func (c *container) getProcess(id string) containerd.Process ***REMOVED***
	c.mu.Lock()
	p := c.execs[id]
	c.mu.Unlock()
	return p
***REMOVED***

func (c *container) setOOMKilled(killed bool) ***REMOVED***
	c.mu.Lock()
	c.oomKilled = killed
	c.mu.Unlock()
***REMOVED***

func (c *container) getOOMKilled() bool ***REMOVED***
	c.mu.Lock()
	killed := c.oomKilled
	c.mu.Unlock()
	return killed
***REMOVED***

type client struct ***REMOVED***
	sync.RWMutex // protects containers map

	remote   *containerd.Client
	stateDir string
	logger   *logrus.Entry

	namespace  string
	backend    Backend
	eventQ     queue
	containers map[string]*container
***REMOVED***

func (c *client) Version(ctx context.Context) (containerd.Version, error) ***REMOVED***
	return c.remote.Version(ctx)
***REMOVED***

func (c *client) Restore(ctx context.Context, id string, attachStdio StdioCallback) (alive bool, pid int, err error) ***REMOVED***
	c.Lock()
	defer c.Unlock()

	var dio *cio.DirectIO
	defer func() ***REMOVED***
		if err != nil && dio != nil ***REMOVED***
			dio.Cancel()
			dio.Close()
		***REMOVED***
		err = wrapError(err)
	***REMOVED***()

	ctr, err := c.remote.LoadContainer(ctx, id)
	if err != nil ***REMOVED***
		return false, -1, errors.WithStack(err)
	***REMOVED***

	attachIO := func(fifos *cio.FIFOSet) (cio.IO, error) ***REMOVED***
		// dio must be assigned to the previously defined dio for the defer above
		// to handle cleanup
		dio, err = cio.NewDirectIO(ctx, fifos)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return attachStdio(dio)
	***REMOVED***
	t, err := ctr.Task(ctx, attachIO)
	if err != nil && !errdefs.IsNotFound(errors.Cause(err)) ***REMOVED***
		return false, -1, err
	***REMOVED***

	if t != nil ***REMOVED***
		s, err := t.Status(ctx)
		if err != nil ***REMOVED***
			return false, -1, err
		***REMOVED***

		alive = s.Status != containerd.Stopped
		pid = int(t.Pid())
	***REMOVED***
	c.containers[id] = &container***REMOVED***
		bundleDir: filepath.Join(c.stateDir, id),
		ctr:       ctr,
		task:      t,
		// TODO(mlaventure): load execs
	***REMOVED***

	c.logger.WithFields(logrus.Fields***REMOVED***
		"container": id,
		"alive":     alive,
		"pid":       pid,
	***REMOVED***).Debug("restored container")

	return alive, pid, nil
***REMOVED***

func (c *client) Create(ctx context.Context, id string, ociSpec *specs.Spec, runtimeOptions interface***REMOVED******REMOVED***) error ***REMOVED***
	if ctr := c.getContainer(id); ctr != nil ***REMOVED***
		return errors.WithStack(newConflictError("id already in use"))
	***REMOVED***

	bdir, err := prepareBundleDir(filepath.Join(c.stateDir, id), ociSpec)
	if err != nil ***REMOVED***
		return errdefs.System(errors.Wrap(err, "prepare bundle dir failed"))
	***REMOVED***

	c.logger.WithField("bundle", bdir).WithField("root", ociSpec.Root.Path).Debug("bundle dir created")

	cdCtr, err := c.remote.NewContainer(ctx, id,
		containerd.WithSpec(ociSpec),
		// TODO(mlaventure): when containerd support lcow, revisit runtime value
		containerd.WithRuntime(fmt.Sprintf("io.containerd.runtime.v1.%s", runtime.GOOS), runtimeOptions))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.Lock()
	c.containers[id] = &container***REMOVED***
		bundleDir: bdir,
		ctr:       cdCtr,
	***REMOVED***
	c.Unlock()

	return nil
***REMOVED***

// Start create and start a task for the specified containerd id
func (c *client) Start(ctx context.Context, id, checkpointDir string, withStdin bool, attachStdio StdioCallback) (int, error) ***REMOVED***
	ctr := c.getContainer(id)
	if ctr == nil ***REMOVED***
		return -1, errors.WithStack(newNotFoundError("no such container"))
	***REMOVED***
	if t := ctr.getTask(); t != nil ***REMOVED***
		return -1, errors.WithStack(newConflictError("container already started"))
	***REMOVED***

	var (
		cp             *types.Descriptor
		t              containerd.Task
		rio            cio.IO
		err            error
		stdinCloseSync = make(chan struct***REMOVED******REMOVED***)
	)

	if checkpointDir != "" ***REMOVED***
		// write checkpoint to the content store
		tar := archive.Diff(ctx, "", checkpointDir)
		cp, err = c.writeContent(ctx, images.MediaTypeContainerd1Checkpoint, checkpointDir, tar)
		// remove the checkpoint when we're done
		defer func() ***REMOVED***
			if cp != nil ***REMOVED***
				err := c.remote.ContentStore().Delete(context.Background(), cp.Digest)
				if err != nil ***REMOVED***
					c.logger.WithError(err).WithFields(logrus.Fields***REMOVED***
						"ref":    checkpointDir,
						"digest": cp.Digest,
					***REMOVED***).Warnf("failed to delete temporary checkpoint entry")
				***REMOVED***
			***REMOVED***
		***REMOVED***()
		if err := tar.Close(); err != nil ***REMOVED***
			return -1, errors.Wrap(err, "failed to close checkpoint tar stream")
		***REMOVED***
		if err != nil ***REMOVED***
			return -1, errors.Wrapf(err, "failed to upload checkpoint to containerd")
		***REMOVED***
	***REMOVED***

	spec, err := ctr.ctr.Spec(ctx)
	if err != nil ***REMOVED***
		return -1, errors.Wrap(err, "failed to retrieve spec")
	***REMOVED***
	uid, gid := getSpecUser(spec)
	t, err = ctr.ctr.NewTask(ctx,
		func(id string) (cio.IO, error) ***REMOVED***
			fifos := newFIFOSet(ctr.bundleDir, InitProcessName, withStdin, spec.Process.Terminal)
			rio, err = c.createIO(fifos, id, InitProcessName, stdinCloseSync, attachStdio)
			return rio, err
		***REMOVED***,
		func(_ context.Context, _ *containerd.Client, info *containerd.TaskInfo) error ***REMOVED***
			info.Checkpoint = cp
			info.Options = &runctypes.CreateOptions***REMOVED***
				IoUid:       uint32(uid),
				IoGid:       uint32(gid),
				NoPivotRoot: os.Getenv("DOCKER_RAMDISK") != "",
			***REMOVED***
			return nil
		***REMOVED***)
	if err != nil ***REMOVED***
		close(stdinCloseSync)
		if rio != nil ***REMOVED***
			rio.Cancel()
			rio.Close()
		***REMOVED***
		return -1, err
	***REMOVED***

	ctr.setTask(t)

	// Signal c.createIO that it can call CloseIO
	close(stdinCloseSync)

	if err := t.Start(ctx); err != nil ***REMOVED***
		if _, err := t.Delete(ctx); err != nil ***REMOVED***
			c.logger.WithError(err).WithField("container", id).
				Error("failed to delete task after fail start")
		***REMOVED***
		ctr.setTask(nil)
		return -1, err
	***REMOVED***

	return int(t.Pid()), nil
***REMOVED***

func (c *client) Exec(ctx context.Context, containerID, processID string, spec *specs.Process, withStdin bool, attachStdio StdioCallback) (int, error) ***REMOVED***
	ctr := c.getContainer(containerID)
	if ctr == nil ***REMOVED***
		return -1, errors.WithStack(newNotFoundError("no such container"))
	***REMOVED***
	t := ctr.getTask()
	if t == nil ***REMOVED***
		return -1, errors.WithStack(newInvalidParameterError("container is not running"))
	***REMOVED***

	if p := ctr.getProcess(processID); p != nil ***REMOVED***
		return -1, errors.WithStack(newConflictError("id already in use"))
	***REMOVED***

	var (
		p              containerd.Process
		rio            cio.IO
		err            error
		stdinCloseSync = make(chan struct***REMOVED******REMOVED***)
	)

	fifos := newFIFOSet(ctr.bundleDir, processID, withStdin, spec.Terminal)

	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if rio != nil ***REMOVED***
				rio.Cancel()
				rio.Close()
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	p, err = t.Exec(ctx, processID, spec, func(id string) (cio.IO, error) ***REMOVED***
		rio, err = c.createIO(fifos, containerID, processID, stdinCloseSync, attachStdio)
		return rio, err
	***REMOVED***)
	if err != nil ***REMOVED***
		close(stdinCloseSync)
		return -1, err
	***REMOVED***

	ctr.addProcess(processID, p)

	// Signal c.createIO that it can call CloseIO
	close(stdinCloseSync)

	if err = p.Start(ctx); err != nil ***REMOVED***
		p.Delete(context.Background())
		ctr.deleteProcess(processID)
		return -1, err
	***REMOVED***

	return int(p.Pid()), nil
***REMOVED***

func (c *client) SignalProcess(ctx context.Context, containerID, processID string, signal int) error ***REMOVED***
	p, err := c.getProcess(containerID, processID)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return wrapError(p.Kill(ctx, syscall.Signal(signal)))
***REMOVED***

func (c *client) ResizeTerminal(ctx context.Context, containerID, processID string, width, height int) error ***REMOVED***
	p, err := c.getProcess(containerID, processID)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return p.Resize(ctx, uint32(width), uint32(height))
***REMOVED***

func (c *client) CloseStdin(ctx context.Context, containerID, processID string) error ***REMOVED***
	p, err := c.getProcess(containerID, processID)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return p.CloseIO(ctx, containerd.WithStdinCloser)
***REMOVED***

func (c *client) Pause(ctx context.Context, containerID string) error ***REMOVED***
	p, err := c.getProcess(containerID, InitProcessName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return p.(containerd.Task).Pause(ctx)
***REMOVED***

func (c *client) Resume(ctx context.Context, containerID string) error ***REMOVED***
	p, err := c.getProcess(containerID, InitProcessName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return p.(containerd.Task).Resume(ctx)
***REMOVED***

func (c *client) Stats(ctx context.Context, containerID string) (*Stats, error) ***REMOVED***
	p, err := c.getProcess(containerID, InitProcessName)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	m, err := p.(containerd.Task).Metrics(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	v, err := typeurl.UnmarshalAny(m.Data)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return interfaceToStats(m.Timestamp, v), nil
***REMOVED***

func (c *client) ListPids(ctx context.Context, containerID string) ([]uint32, error) ***REMOVED***
	p, err := c.getProcess(containerID, InitProcessName)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	pis, err := p.(containerd.Task).Pids(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var pids []uint32
	for _, i := range pis ***REMOVED***
		pids = append(pids, i.Pid)
	***REMOVED***

	return pids, nil
***REMOVED***

func (c *client) Summary(ctx context.Context, containerID string) ([]Summary, error) ***REMOVED***
	p, err := c.getProcess(containerID, InitProcessName)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	pis, err := p.(containerd.Task).Pids(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var infos []Summary
	for _, pi := range pis ***REMOVED***
		i, err := typeurl.UnmarshalAny(pi.Info)
		if err != nil ***REMOVED***
			return nil, errors.Wrap(err, "unable to decode process details")
		***REMOVED***
		s, err := summaryFromInterface(i)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		infos = append(infos, *s)
	***REMOVED***

	return infos, nil
***REMOVED***

func (c *client) DeleteTask(ctx context.Context, containerID string) (uint32, time.Time, error) ***REMOVED***
	p, err := c.getProcess(containerID, InitProcessName)
	if err != nil ***REMOVED***
		return 255, time.Now(), nil
	***REMOVED***

	status, err := p.(containerd.Task).Delete(ctx)
	if err != nil ***REMOVED***
		return 255, time.Now(), nil
	***REMOVED***

	if ctr := c.getContainer(containerID); ctr != nil ***REMOVED***
		ctr.setTask(nil)
	***REMOVED***
	return status.ExitCode(), status.ExitTime(), nil
***REMOVED***

func (c *client) Delete(ctx context.Context, containerID string) error ***REMOVED***
	ctr := c.getContainer(containerID)
	if ctr == nil ***REMOVED***
		return errors.WithStack(newNotFoundError("no such container"))
	***REMOVED***

	if err := ctr.ctr.Delete(ctx); err != nil ***REMOVED***
		return err
	***REMOVED***

	if os.Getenv("LIBCONTAINERD_NOCLEAN") != "1" ***REMOVED***
		if err := os.RemoveAll(ctr.bundleDir); err != nil ***REMOVED***
			c.logger.WithError(err).WithFields(logrus.Fields***REMOVED***
				"container": containerID,
				"bundle":    ctr.bundleDir,
			***REMOVED***).Error("failed to remove state dir")
		***REMOVED***
	***REMOVED***

	c.removeContainer(containerID)

	return nil
***REMOVED***

func (c *client) Status(ctx context.Context, containerID string) (Status, error) ***REMOVED***
	ctr := c.getContainer(containerID)
	if ctr == nil ***REMOVED***
		return StatusUnknown, errors.WithStack(newNotFoundError("no such container"))
	***REMOVED***

	t := ctr.getTask()
	if t == nil ***REMOVED***
		return StatusUnknown, errors.WithStack(newNotFoundError("no such task"))
	***REMOVED***

	s, err := t.Status(ctx)
	if err != nil ***REMOVED***
		return StatusUnknown, err
	***REMOVED***

	return Status(s.Status), nil
***REMOVED***

func (c *client) CreateCheckpoint(ctx context.Context, containerID, checkpointDir string, exit bool) error ***REMOVED***
	p, err := c.getProcess(containerID, InitProcessName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	img, err := p.(containerd.Task).Checkpoint(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// Whatever happens, delete the checkpoint from containerd
	defer func() ***REMOVED***
		err := c.remote.ImageService().Delete(context.Background(), img.Name())
		if err != nil ***REMOVED***
			c.logger.WithError(err).WithField("digest", img.Target().Digest).
				Warnf("failed to delete checkpoint image")
		***REMOVED***
	***REMOVED***()

	b, err := content.ReadBlob(ctx, c.remote.ContentStore(), img.Target().Digest)
	if err != nil ***REMOVED***
		return errdefs.System(errors.Wrapf(err, "failed to retrieve checkpoint data"))
	***REMOVED***
	var index v1.Index
	if err := json.Unmarshal(b, &index); err != nil ***REMOVED***
		return errdefs.System(errors.Wrapf(err, "failed to decode checkpoint data"))
	***REMOVED***

	var cpDesc *v1.Descriptor
	for _, m := range index.Manifests ***REMOVED***
		if m.MediaType == images.MediaTypeContainerd1Checkpoint ***REMOVED***
			cpDesc = &m
			break
		***REMOVED***
	***REMOVED***
	if cpDesc == nil ***REMOVED***
		return errdefs.System(errors.Wrapf(err, "invalid checkpoint"))
	***REMOVED***

	rat, err := c.remote.ContentStore().ReaderAt(ctx, cpDesc.Digest)
	if err != nil ***REMOVED***
		return errdefs.System(errors.Wrapf(err, "failed to get checkpoint reader"))
	***REMOVED***
	defer rat.Close()
	_, err = archive.Apply(ctx, checkpointDir, content.NewReader(rat))
	if err != nil ***REMOVED***
		return errdefs.System(errors.Wrapf(err, "failed to read checkpoint reader"))
	***REMOVED***

	return err
***REMOVED***

func (c *client) getContainer(id string) *container ***REMOVED***
	c.RLock()
	ctr := c.containers[id]
	c.RUnlock()

	return ctr
***REMOVED***

func (c *client) removeContainer(id string) ***REMOVED***
	c.Lock()
	delete(c.containers, id)
	c.Unlock()
***REMOVED***

func (c *client) getProcess(containerID, processID string) (containerd.Process, error) ***REMOVED***
	ctr := c.getContainer(containerID)
	if ctr == nil ***REMOVED***
		return nil, errors.WithStack(newNotFoundError("no such container"))
	***REMOVED***

	t := ctr.getTask()
	if t == nil ***REMOVED***
		return nil, errors.WithStack(newNotFoundError("container is not running"))
	***REMOVED***
	if processID == InitProcessName ***REMOVED***
		return t, nil
	***REMOVED***

	p := ctr.getProcess(processID)
	if p == nil ***REMOVED***
		return nil, errors.WithStack(newNotFoundError("no such exec"))
	***REMOVED***
	return p, nil
***REMOVED***

// createIO creates the io to be used by a process
// This needs to get a pointer to interface as upon closure the process may not have yet been registered
func (c *client) createIO(fifos *cio.FIFOSet, containerID, processID string, stdinCloseSync chan struct***REMOVED******REMOVED***, attachStdio StdioCallback) (cio.IO, error) ***REMOVED***
	io, err := cio.NewDirectIO(context.Background(), fifos)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if io.Stdin != nil ***REMOVED***
		var (
			err       error
			stdinOnce sync.Once
		)
		pipe := io.Stdin
		io.Stdin = ioutils.NewWriteCloserWrapper(pipe, func() error ***REMOVED***
			stdinOnce.Do(func() ***REMOVED***
				err = pipe.Close()
				// Do the rest in a new routine to avoid a deadlock if the
				// Exec/Start call failed.
				go func() ***REMOVED***
					<-stdinCloseSync
					p, err := c.getProcess(containerID, processID)
					if err == nil ***REMOVED***
						err = p.CloseIO(context.Background(), containerd.WithStdinCloser)
						if err != nil && strings.Contains(err.Error(), "transport is closing") ***REMOVED***
							err = nil
						***REMOVED***
					***REMOVED***
				***REMOVED***()
			***REMOVED***)
			return err
		***REMOVED***)
	***REMOVED***

	rio, err := attachStdio(io)
	if err != nil ***REMOVED***
		io.Cancel()
		io.Close()
	***REMOVED***
	return rio, err
***REMOVED***

func (c *client) processEvent(ctr *container, et EventType, ei EventInfo) ***REMOVED***
	c.eventQ.append(ei.ContainerID, func() ***REMOVED***
		err := c.backend.ProcessEvent(ei.ContainerID, et, ei)
		if err != nil ***REMOVED***
			c.logger.WithError(err).WithFields(logrus.Fields***REMOVED***
				"container":  ei.ContainerID,
				"event":      et,
				"event-info": ei,
			***REMOVED***).Error("failed to process event")
		***REMOVED***

		if et == EventExit && ei.ProcessID != ei.ContainerID ***REMOVED***
			p := ctr.getProcess(ei.ProcessID)
			if p == nil ***REMOVED***
				c.logger.WithError(errors.New("no such process")).
					WithFields(logrus.Fields***REMOVED***
						"container": ei.ContainerID,
						"process":   ei.ProcessID,
					***REMOVED***).Error("exit event")
				return
			***REMOVED***
			_, err = p.Delete(context.Background())
			if err != nil ***REMOVED***
				c.logger.WithError(err).WithFields(logrus.Fields***REMOVED***
					"container": ei.ContainerID,
					"process":   ei.ProcessID,
				***REMOVED***).Warn("failed to delete process")
			***REMOVED***
			ctr.deleteProcess(ei.ProcessID)

			ctr := c.getContainer(ei.ContainerID)
			if ctr == nil ***REMOVED***
				c.logger.WithFields(logrus.Fields***REMOVED***
					"container": ei.ContainerID,
				***REMOVED***).Error("failed to find container")
			***REMOVED*** else ***REMOVED***
				newFIFOSet(ctr.bundleDir, ei.ProcessID, true, false).Close()
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***

func (c *client) processEventStream(ctx context.Context) ***REMOVED***
	var (
		err         error
		eventStream eventsapi.Events_SubscribeClient
		ev          *eventsapi.Envelope
		et          EventType
		ei          EventInfo
		ctr         *container
	)
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			select ***REMOVED***
			case <-ctx.Done():
				c.logger.WithError(ctx.Err()).
					Info("stopping event stream following graceful shutdown")
			default:
				go c.processEventStream(ctx)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	eventStream, err = c.remote.EventService().Subscribe(ctx, &eventsapi.SubscribeRequest***REMOVED***
		Filters: []string***REMOVED***
			// Filter on both namespace *and* topic. To create an "and" filter,
			// this must be a single, comma-separated string
			"namespace==" + c.namespace + ",topic~=|^/tasks/|",
		***REMOVED***,
	***REMOVED***, grpc.FailFast(false))
	if err != nil ***REMOVED***
		return
	***REMOVED***

	var oomKilled bool
	for ***REMOVED***
		ev, err = eventStream.Recv()
		if err != nil ***REMOVED***
			errStatus, ok := status.FromError(err)
			if !ok || errStatus.Code() != codes.Canceled ***REMOVED***
				c.logger.WithError(err).Error("failed to get event")
			***REMOVED***
			return
		***REMOVED***

		if ev.Event == nil ***REMOVED***
			c.logger.WithField("event", ev).Warn("invalid event")
			continue
		***REMOVED***

		v, err := typeurl.UnmarshalAny(ev.Event)
		if err != nil ***REMOVED***
			c.logger.WithError(err).WithField("event", ev).Warn("failed to unmarshal event")
			continue
		***REMOVED***

		c.logger.WithField("topic", ev.Topic).Debug("event")

		switch t := v.(type) ***REMOVED***
		case *events.TaskCreate:
			et = EventCreate
			ei = EventInfo***REMOVED***
				ContainerID: t.ContainerID,
				ProcessID:   t.ContainerID,
				Pid:         t.Pid,
			***REMOVED***
		case *events.TaskStart:
			et = EventStart
			ei = EventInfo***REMOVED***
				ContainerID: t.ContainerID,
				ProcessID:   t.ContainerID,
				Pid:         t.Pid,
			***REMOVED***
		case *events.TaskExit:
			et = EventExit
			ei = EventInfo***REMOVED***
				ContainerID: t.ContainerID,
				ProcessID:   t.ID,
				Pid:         t.Pid,
				ExitCode:    t.ExitStatus,
				ExitedAt:    t.ExitedAt,
			***REMOVED***
		case *events.TaskOOM:
			et = EventOOM
			ei = EventInfo***REMOVED***
				ContainerID: t.ContainerID,
				OOMKilled:   true,
			***REMOVED***
			oomKilled = true
		case *events.TaskExecAdded:
			et = EventExecAdded
			ei = EventInfo***REMOVED***
				ContainerID: t.ContainerID,
				ProcessID:   t.ExecID,
			***REMOVED***
		case *events.TaskExecStarted:
			et = EventExecStarted
			ei = EventInfo***REMOVED***
				ContainerID: t.ContainerID,
				ProcessID:   t.ExecID,
				Pid:         t.Pid,
			***REMOVED***
		case *events.TaskPaused:
			et = EventPaused
			ei = EventInfo***REMOVED***
				ContainerID: t.ContainerID,
			***REMOVED***
		case *events.TaskResumed:
			et = EventResumed
			ei = EventInfo***REMOVED***
				ContainerID: t.ContainerID,
			***REMOVED***
		default:
			c.logger.WithFields(logrus.Fields***REMOVED***
				"topic": ev.Topic,
				"type":  reflect.TypeOf(t)***REMOVED***,
			).Info("ignoring event")
			continue
		***REMOVED***

		ctr = c.getContainer(ei.ContainerID)
		if ctr == nil ***REMOVED***
			c.logger.WithField("container", ei.ContainerID).Warn("unknown container")
			continue
		***REMOVED***

		if oomKilled ***REMOVED***
			ctr.setOOMKilled(true)
			oomKilled = false
		***REMOVED***
		ei.OOMKilled = ctr.getOOMKilled()

		c.processEvent(ctr, et, ei)
	***REMOVED***
***REMOVED***

func (c *client) writeContent(ctx context.Context, mediaType, ref string, r io.Reader) (*types.Descriptor, error) ***REMOVED***
	writer, err := c.remote.ContentStore().Writer(ctx, ref, 0, "")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer writer.Close()
	size, err := io.Copy(writer, r)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	labels := map[string]string***REMOVED***
		"containerd.io/gc.root": time.Now().UTC().Format(time.RFC3339),
	***REMOVED***
	if err := writer.Commit(ctx, 0, "", content.WithLabels(labels)); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &types.Descriptor***REMOVED***
		MediaType: mediaType,
		Digest:    writer.Digest(),
		Size_:     size,
	***REMOVED***, nil
***REMOVED***

func wrapError(err error) error ***REMOVED***
	switch ***REMOVED***
	case err == nil:
		return nil
	case containerderrors.IsNotFound(err):
		return errdefs.NotFound(err)
	***REMOVED***

	msg := err.Error()
	for _, s := range []string***REMOVED***"container does not exist", "not found", "no such container"***REMOVED*** ***REMOVED***
		if strings.Contains(msg, s) ***REMOVED***
			return errdefs.NotFound(err)
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***
