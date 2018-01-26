// +build linux

package linux

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	eventstypes "github.com/containerd/containerd/api/events"
	"github.com/containerd/containerd/api/types"
	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/events/exchange"
	"github.com/containerd/containerd/identifiers"
	"github.com/containerd/containerd/linux/proc"
	"github.com/containerd/containerd/linux/runctypes"
	shim "github.com/containerd/containerd/linux/shim/v1"
	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/metadata"
	"github.com/containerd/containerd/mount"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/platforms"
	"github.com/containerd/containerd/plugin"
	"github.com/containerd/containerd/reaper"
	"github.com/containerd/containerd/runtime"
	"github.com/containerd/containerd/sys"
	runc "github.com/containerd/go-runc"
	"github.com/containerd/typeurl"
	ptypes "github.com/gogo/protobuf/types"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

var (
	pluginID = fmt.Sprintf("%s.%s", plugin.RuntimePlugin, "linux")
	empty    = &ptypes.Empty***REMOVED******REMOVED***
)

const (
	configFilename = "config.json"
	defaultRuntime = "runc"
	defaultShim    = "containerd-shim"
)

func init() ***REMOVED***
	plugin.Register(&plugin.Registration***REMOVED***
		Type:   plugin.RuntimePlugin,
		ID:     "linux",
		InitFn: New,
		Requires: []plugin.Type***REMOVED***
			plugin.TaskMonitorPlugin,
			plugin.MetadataPlugin,
		***REMOVED***,
		Config: &Config***REMOVED***
			Shim:    defaultShim,
			Runtime: defaultRuntime,
		***REMOVED***,
	***REMOVED***)
***REMOVED***

var _ = (runtime.Runtime)(&Runtime***REMOVED******REMOVED***)

// Config options for the runtime
type Config struct ***REMOVED***
	// Shim is a path or name of binary implementing the Shim GRPC API
	Shim string `toml:"shim"`
	// Runtime is a path or name of an OCI runtime used by the shim
	Runtime string `toml:"runtime"`
	// RuntimeRoot is the path that shall be used by the OCI runtime for its data
	RuntimeRoot string `toml:"runtime_root"`
	// NoShim calls runc directly from within the pkg
	NoShim bool `toml:"no_shim"`
	// Debug enable debug on the shim
	ShimDebug bool `toml:"shim_debug"`
***REMOVED***

// New returns a configured runtime
func New(ic *plugin.InitContext) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	ic.Meta.Platforms = []ocispec.Platform***REMOVED***platforms.DefaultSpec()***REMOVED***

	if err := os.MkdirAll(ic.Root, 0711); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := os.MkdirAll(ic.State, 0711); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	monitor, err := ic.Get(plugin.TaskMonitorPlugin)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	m, err := ic.Get(plugin.MetadataPlugin)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	cfg := ic.Config.(*Config)
	r := &Runtime***REMOVED***
		root:    ic.Root,
		state:   ic.State,
		monitor: monitor.(runtime.TaskMonitor),
		tasks:   runtime.NewTaskList(),
		db:      m.(*metadata.DB),
		address: ic.Address,
		events:  ic.Events,
		config:  cfg,
	***REMOVED***
	tasks, err := r.restoreTasks(ic.Context)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// TODO: need to add the tasks to the monitor
	for _, t := range tasks ***REMOVED***
		if err := r.tasks.AddWithNamespace(t.namespace, t); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return r, nil
***REMOVED***

// Runtime for a linux based system
type Runtime struct ***REMOVED***
	root    string
	state   string
	address string

	monitor runtime.TaskMonitor
	tasks   *runtime.TaskList
	db      *metadata.DB
	events  *exchange.Exchange

	config *Config
***REMOVED***

// ID of the runtime
func (r *Runtime) ID() string ***REMOVED***
	return pluginID
***REMOVED***

// Create a new task
func (r *Runtime) Create(ctx context.Context, id string, opts runtime.CreateOpts) (_ runtime.Task, err error) ***REMOVED***
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := identifiers.Validate(id); err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "invalid task id")
	***REMOVED***

	ropts, err := r.getRuncOptions(ctx, id)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ec := reaper.Default.Subscribe()
	defer reaper.Default.Unsubscribe(ec)

	bundle, err := newBundle(id,
		filepath.Join(r.state, namespace),
		filepath.Join(r.root, namespace),
		opts.Spec.Value)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			bundle.Delete()
		***REMOVED***
	***REMOVED***()

	shimopt := ShimLocal(r.events)
	if !r.config.NoShim ***REMOVED***
		var cgroup string
		if opts.Options != nil ***REMOVED***
			v, err := typeurl.UnmarshalAny(opts.Options)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			cgroup = v.(*runctypes.CreateOptions).ShimCgroup
		***REMOVED***
		exitHandler := func() ***REMOVED***
			log.G(ctx).WithField("id", id).Info("shim reaped")
			t, err := r.tasks.Get(ctx, id)
			if err != nil ***REMOVED***
				// Task was never started or was already sucessfully deleted
				return
			***REMOVED***
			lc := t.(*Task)

			// Stop the monitor
			if err := r.monitor.Stop(lc); err != nil ***REMOVED***
				log.G(ctx).WithError(err).WithFields(logrus.Fields***REMOVED***
					"id":        id,
					"namespace": namespace,
				***REMOVED***).Warn("failed to stop monitor")
			***REMOVED***

			log.G(ctx).WithFields(logrus.Fields***REMOVED***
				"id":        id,
				"namespace": namespace,
			***REMOVED***).Warn("cleaning up after killed shim")
			err = r.cleanupAfterDeadShim(context.Background(), bundle, namespace, id, lc.pid, ec)
			if err == nil ***REMOVED***
				r.tasks.Delete(ctx, lc)
			***REMOVED*** else ***REMOVED***
				log.G(ctx).WithError(err).WithFields(logrus.Fields***REMOVED***
					"id":        id,
					"namespace": namespace,
				***REMOVED***).Warn("failed to clen up after killed shim")
			***REMOVED***
		***REMOVED***
		shimopt = ShimRemote(r.config.Shim, r.address, cgroup, r.config.ShimDebug, exitHandler)
	***REMOVED***

	s, err := bundle.NewShimClient(ctx, namespace, shimopt, ropts)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if kerr := s.KillShim(ctx); kerr != nil ***REMOVED***
				log.G(ctx).WithError(err).Error("failed to kill shim")
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	rt := r.config.Runtime
	if ropts != nil && ropts.Runtime != "" ***REMOVED***
		rt = ropts.Runtime
	***REMOVED***
	sopts := &shim.CreateTaskRequest***REMOVED***
		ID:         id,
		Bundle:     bundle.path,
		Runtime:    rt,
		Stdin:      opts.IO.Stdin,
		Stdout:     opts.IO.Stdout,
		Stderr:     opts.IO.Stderr,
		Terminal:   opts.IO.Terminal,
		Checkpoint: opts.Checkpoint,
		Options:    opts.Options,
	***REMOVED***
	for _, m := range opts.Rootfs ***REMOVED***
		sopts.Rootfs = append(sopts.Rootfs, &types.Mount***REMOVED***
			Type:    m.Type,
			Source:  m.Source,
			Options: m.Options,
		***REMOVED***)
	***REMOVED***
	cr, err := s.Create(ctx, sopts)
	if err != nil ***REMOVED***
		return nil, errdefs.FromGRPC(err)
	***REMOVED***
	t, err := newTask(id, namespace, int(cr.Pid), s, r.monitor, r.events,
		proc.NewRunc(ropts.RuntimeRoot, sopts.Bundle, namespace, rt, ropts.CriuPath, ropts.SystemdCgroup))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := r.tasks.Add(ctx, t); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// after the task is created, add it to the monitor if it has a cgroup
	// this can be different on a checkpoint/restore
	if t.cg != nil ***REMOVED***
		if err = r.monitor.Monitor(t); err != nil ***REMOVED***
			if _, err := r.Delete(ctx, t); err != nil ***REMOVED***
				log.G(ctx).WithError(err).Error("deleting task after failed monitor")
			***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	r.events.Publish(ctx, runtime.TaskCreateEventTopic, &eventstypes.TaskCreate***REMOVED***
		ContainerID: sopts.ID,
		Bundle:      sopts.Bundle,
		Rootfs:      sopts.Rootfs,
		IO: &eventstypes.TaskIO***REMOVED***
			Stdin:    sopts.Stdin,
			Stdout:   sopts.Stdout,
			Stderr:   sopts.Stderr,
			Terminal: sopts.Terminal,
		***REMOVED***,
		Checkpoint: sopts.Checkpoint,
		Pid:        uint32(t.pid),
	***REMOVED***)

	return t, nil
***REMOVED***

// Delete a task removing all on disk state
func (r *Runtime) Delete(ctx context.Context, c runtime.Task) (*runtime.Exit, error) ***REMOVED***
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	lc, ok := c.(*Task)
	if !ok ***REMOVED***
		return nil, fmt.Errorf("task cannot be cast as *linux.Task")
	***REMOVED***
	if err := r.monitor.Stop(lc); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	bundle := loadBundle(
		lc.id,
		filepath.Join(r.state, namespace, lc.id),
		filepath.Join(r.root, namespace, lc.id),
	)

	rsp, err := lc.shim.Delete(ctx, empty)
	if err != nil ***REMOVED***
		if cerr := r.cleanupAfterDeadShim(ctx, bundle, namespace, c.ID(), lc.pid, nil); cerr != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("unable to cleanup task")
		***REMOVED***
		return nil, errdefs.FromGRPC(err)
	***REMOVED***
	r.tasks.Delete(ctx, lc)
	if err := lc.shim.KillShim(ctx); err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("failed to kill shim")
	***REMOVED***

	if err := bundle.Delete(); err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("failed to delete bundle")
	***REMOVED***
	r.events.Publish(ctx, runtime.TaskDeleteEventTopic, &eventstypes.TaskDelete***REMOVED***
		ContainerID: lc.id,
		ExitStatus:  rsp.ExitStatus,
		ExitedAt:    rsp.ExitedAt,
		Pid:         rsp.Pid,
	***REMOVED***)
	return &runtime.Exit***REMOVED***
		Status:    rsp.ExitStatus,
		Timestamp: rsp.ExitedAt,
		Pid:       rsp.Pid,
	***REMOVED***, nil
***REMOVED***

// Tasks returns all tasks known to the runtime
func (r *Runtime) Tasks(ctx context.Context) ([]runtime.Task, error) ***REMOVED***
	return r.tasks.GetAll(ctx)
***REMOVED***

func (r *Runtime) restoreTasks(ctx context.Context) ([]*Task, error) ***REMOVED***
	dir, err := ioutil.ReadDir(r.state)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var o []*Task
	for _, namespace := range dir ***REMOVED***
		if !namespace.IsDir() ***REMOVED***
			continue
		***REMOVED***
		name := namespace.Name()
		log.G(ctx).WithField("namespace", name).Debug("loading tasks in namespace")
		tasks, err := r.loadTasks(ctx, name)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		o = append(o, tasks...)
	***REMOVED***
	return o, nil
***REMOVED***

// Get a specific task by task id
func (r *Runtime) Get(ctx context.Context, id string) (runtime.Task, error) ***REMOVED***
	return r.tasks.Get(ctx, id)
***REMOVED***

func (r *Runtime) loadTasks(ctx context.Context, ns string) ([]*Task, error) ***REMOVED***
	dir, err := ioutil.ReadDir(filepath.Join(r.state, ns))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var o []*Task
	for _, path := range dir ***REMOVED***
		if !path.IsDir() ***REMOVED***
			continue
		***REMOVED***
		id := path.Name()
		bundle := loadBundle(
			id,
			filepath.Join(r.state, ns, id),
			filepath.Join(r.root, ns, id),
		)
		ctx = namespaces.WithNamespace(ctx, ns)
		pid, _ := runc.ReadPidFile(filepath.Join(bundle.path, proc.InitPidFile))
		s, err := bundle.NewShimClient(ctx, ns, ShimConnect(), nil)
		if err != nil ***REMOVED***
			log.G(ctx).WithError(err).WithFields(logrus.Fields***REMOVED***
				"id":        id,
				"namespace": ns,
			***REMOVED***).Error("connecting to shim")
			err := r.cleanupAfterDeadShim(ctx, bundle, ns, id, pid, nil)
			if err != nil ***REMOVED***
				log.G(ctx).WithError(err).WithField("bundle", bundle.path).
					Error("cleaning up after dead shim")
			***REMOVED***
			continue
		***REMOVED***
		ropts, err := r.getRuncOptions(ctx, id)
		if err != nil ***REMOVED***
			log.G(ctx).WithError(err).WithField("id", id).
				Error("get runtime options")
			continue
		***REMOVED***

		t, err := newTask(id, ns, pid, s, r.monitor, r.events,
			proc.NewRunc(ropts.RuntimeRoot, bundle.path, ns, ropts.Runtime, ropts.CriuPath, ropts.SystemdCgroup))
		if err != nil ***REMOVED***
			log.G(ctx).WithError(err).Error("loading task type")
			continue
		***REMOVED***
		o = append(o, t)
	***REMOVED***
	return o, nil
***REMOVED***

func (r *Runtime) cleanupAfterDeadShim(ctx context.Context, bundle *bundle, ns, id string, pid int, ec chan runc.Exit) error ***REMOVED***
	ctx = namespaces.WithNamespace(ctx, ns)
	if err := r.terminate(ctx, bundle, ns, id); err != nil ***REMOVED***
		if r.config.ShimDebug ***REMOVED***
			return errors.Wrap(err, "failed to terminate task, leaving bundle for debugging")
		***REMOVED***
		log.G(ctx).WithError(err).Warn("failed to terminate task")
	***REMOVED***

	if ec != nil ***REMOVED***
		// if sub-reaper is set, reap our new child
		if v, err := sys.GetSubreaper(); err == nil && v == 1 ***REMOVED***
			for e := range ec ***REMOVED***
				if e.Pid == pid ***REMOVED***
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Notify Client
	exitedAt := time.Now().UTC()
	r.events.Publish(ctx, runtime.TaskExitEventTopic, &eventstypes.TaskExit***REMOVED***
		ContainerID: id,
		ID:          id,
		Pid:         uint32(pid),
		ExitStatus:  128 + uint32(unix.SIGKILL),
		ExitedAt:    exitedAt,
	***REMOVED***)

	if err := bundle.Delete(); err != nil ***REMOVED***
		log.G(ctx).WithError(err).Error("delete bundle")
	***REMOVED***

	r.events.Publish(ctx, runtime.TaskDeleteEventTopic, &eventstypes.TaskDelete***REMOVED***
		ContainerID: id,
		Pid:         uint32(pid),
		ExitStatus:  128 + uint32(unix.SIGKILL),
		ExitedAt:    exitedAt,
	***REMOVED***)

	return nil
***REMOVED***

func (r *Runtime) terminate(ctx context.Context, bundle *bundle, ns, id string) error ***REMOVED***
	ctx = namespaces.WithNamespace(ctx, ns)
	rt, err := r.getRuntime(ctx, ns, id)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := rt.Delete(ctx, id, &runc.DeleteOpts***REMOVED***
		Force: true,
	***REMOVED***); err != nil ***REMOVED***
		log.G(ctx).WithError(err).Warnf("delete runtime state %s", id)
	***REMOVED***
	if err := mount.Unmount(filepath.Join(bundle.path, "rootfs"), 0); err != nil ***REMOVED***
		log.G(ctx).WithError(err).WithFields(logrus.Fields***REMOVED***
			"path": bundle.path,
			"id":   id,
		***REMOVED***).Warnf("unmount task rootfs")
	***REMOVED***
	return nil
***REMOVED***

func (r *Runtime) getRuntime(ctx context.Context, ns, id string) (*runc.Runc, error) ***REMOVED***
	ropts, err := r.getRuncOptions(ctx, id)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var (
		cmd  = r.config.Runtime
		root = proc.RuncRoot
	)
	if ropts != nil ***REMOVED***
		if ropts.Runtime != "" ***REMOVED***
			cmd = ropts.Runtime
		***REMOVED***
		if ropts.RuntimeRoot != "" ***REMOVED***
			root = ropts.RuntimeRoot
		***REMOVED***
	***REMOVED***

	return &runc.Runc***REMOVED***
		Command:      cmd,
		LogFormat:    runc.JSON,
		PdeathSignal: unix.SIGKILL,
		Root:         filepath.Join(root, ns),
	***REMOVED***, nil
***REMOVED***

func (r *Runtime) getRuncOptions(ctx context.Context, id string) (*runctypes.RuncOptions, error) ***REMOVED***
	var container containers.Container

	if err := r.db.View(func(tx *bolt.Tx) error ***REMOVED***
		store := metadata.NewContainerStore(tx)
		var err error
		container, err = store.Get(ctx, id)
		return err
	***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if container.Runtime.Options != nil ***REMOVED***
		v, err := typeurl.UnmarshalAny(container.Runtime.Options)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		ropts, ok := v.(*runctypes.RuncOptions)
		if !ok ***REMOVED***
			return nil, errors.New("invalid runtime options format")
		***REMOVED***

		return ropts, nil
	***REMOVED***
	return &runctypes.RuncOptions***REMOVED******REMOVED***, nil
***REMOVED***
