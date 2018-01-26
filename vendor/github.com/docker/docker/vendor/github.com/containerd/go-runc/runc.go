package runc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// Format is the type of log formatting options avaliable
type Format string

const (
	none Format = ""
	JSON Format = "json"
	Text Format = "text"
	// DefaultCommand is the default command for Runc
	DefaultCommand = "runc"
)

// Runc is the client to the runc cli
type Runc struct ***REMOVED***
	//If command is empty, DefaultCommand is used
	Command       string
	Root          string
	Debug         bool
	Log           string
	LogFormat     Format
	PdeathSignal  syscall.Signal
	Setpgid       bool
	Criu          string
	SystemdCgroup bool
***REMOVED***

// List returns all containers created inside the provided runc root directory
func (r *Runc) List(context context.Context) ([]*Container, error) ***REMOVED***
	data, err := cmdOutput(r.command(context, "list", "--format=json"), false)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var out []*Container
	if err := json.Unmarshal(data, &out); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return out, nil
***REMOVED***

// State returns the state for the container provided by id
func (r *Runc) State(context context.Context, id string) (*Container, error) ***REMOVED***
	data, err := cmdOutput(r.command(context, "state", id), true)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("%s: %s", err, data)
	***REMOVED***
	var c Container
	if err := json.Unmarshal(data, &c); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &c, nil
***REMOVED***

type ConsoleSocket interface ***REMOVED***
	Path() string
***REMOVED***

type CreateOpts struct ***REMOVED***
	IO
	// PidFile is a path to where a pid file should be created
	PidFile       string
	ConsoleSocket ConsoleSocket
	Detach        bool
	NoPivot       bool
	NoNewKeyring  bool
	ExtraFiles    []*os.File
***REMOVED***

func (o *CreateOpts) args() (out []string, err error) ***REMOVED***
	if o.PidFile != "" ***REMOVED***
		abs, err := filepath.Abs(o.PidFile)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		out = append(out, "--pid-file", abs)
	***REMOVED***
	if o.ConsoleSocket != nil ***REMOVED***
		out = append(out, "--console-socket", o.ConsoleSocket.Path())
	***REMOVED***
	if o.NoPivot ***REMOVED***
		out = append(out, "--no-pivot")
	***REMOVED***
	if o.NoNewKeyring ***REMOVED***
		out = append(out, "--no-new-keyring")
	***REMOVED***
	if o.Detach ***REMOVED***
		out = append(out, "--detach")
	***REMOVED***
	if o.ExtraFiles != nil ***REMOVED***
		out = append(out, "--preserve-fds", strconv.Itoa(len(o.ExtraFiles)))
	***REMOVED***
	return out, nil
***REMOVED***

// Create creates a new container and returns its pid if it was created successfully
func (r *Runc) Create(context context.Context, id, bundle string, opts *CreateOpts) error ***REMOVED***
	args := []string***REMOVED***"create", "--bundle", bundle***REMOVED***
	if opts != nil ***REMOVED***
		oargs, err := opts.args()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		args = append(args, oargs...)
	***REMOVED***
	cmd := r.command(context, append(args, id)...)
	if opts != nil && opts.IO != nil ***REMOVED***
		opts.Set(cmd)
	***REMOVED***
	cmd.ExtraFiles = opts.ExtraFiles

	if cmd.Stdout == nil && cmd.Stderr == nil ***REMOVED***
		data, err := cmdOutput(cmd, true)
		if err != nil ***REMOVED***
			return fmt.Errorf("%s: %s", err, data)
		***REMOVED***
		return nil
	***REMOVED***
	ec, err := Monitor.Start(cmd)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if opts != nil && opts.IO != nil ***REMOVED***
		if c, ok := opts.IO.(StartCloser); ok ***REMOVED***
			if err := c.CloseAfterStart(); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	status, err := Monitor.Wait(cmd, ec)
	if err == nil && status != 0 ***REMOVED***
		err = fmt.Errorf("%s did not terminate sucessfully", cmd.Args[0])
	***REMOVED***
	return err
***REMOVED***

// Start will start an already created container
func (r *Runc) Start(context context.Context, id string) error ***REMOVED***
	return r.runOrError(r.command(context, "start", id))
***REMOVED***

type ExecOpts struct ***REMOVED***
	IO
	PidFile       string
	ConsoleSocket ConsoleSocket
	Detach        bool
***REMOVED***

func (o *ExecOpts) args() (out []string, err error) ***REMOVED***
	if o.ConsoleSocket != nil ***REMOVED***
		out = append(out, "--console-socket", o.ConsoleSocket.Path())
	***REMOVED***
	if o.Detach ***REMOVED***
		out = append(out, "--detach")
	***REMOVED***
	if o.PidFile != "" ***REMOVED***
		abs, err := filepath.Abs(o.PidFile)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		out = append(out, "--pid-file", abs)
	***REMOVED***
	return out, nil
***REMOVED***

// Exec executres and additional process inside the container based on a full
// OCI Process specification
func (r *Runc) Exec(context context.Context, id string, spec specs.Process, opts *ExecOpts) error ***REMOVED***
	f, err := ioutil.TempFile("", "runc-process")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer os.Remove(f.Name())
	err = json.NewEncoder(f).Encode(spec)
	f.Close()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	args := []string***REMOVED***"exec", "--process", f.Name()***REMOVED***
	if opts != nil ***REMOVED***
		oargs, err := opts.args()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		args = append(args, oargs...)
	***REMOVED***
	cmd := r.command(context, append(args, id)...)
	if opts != nil && opts.IO != nil ***REMOVED***
		opts.Set(cmd)
	***REMOVED***
	if cmd.Stdout == nil && cmd.Stderr == nil ***REMOVED***
		data, err := cmdOutput(cmd, true)
		if err != nil ***REMOVED***
			return fmt.Errorf("%s: %s", err, data)
		***REMOVED***
		return nil
	***REMOVED***
	ec, err := Monitor.Start(cmd)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if opts != nil && opts.IO != nil ***REMOVED***
		if c, ok := opts.IO.(StartCloser); ok ***REMOVED***
			if err := c.CloseAfterStart(); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	status, err := Monitor.Wait(cmd, ec)
	if err == nil && status != 0 ***REMOVED***
		err = fmt.Errorf("%s did not terminate sucessfully", cmd.Args[0])
	***REMOVED***
	return err
***REMOVED***

// Run runs the create, start, delete lifecycle of the container
// and returns its exit status after it has exited
func (r *Runc) Run(context context.Context, id, bundle string, opts *CreateOpts) (int, error) ***REMOVED***
	args := []string***REMOVED***"run", "--bundle", bundle***REMOVED***
	if opts != nil ***REMOVED***
		oargs, err := opts.args()
		if err != nil ***REMOVED***
			return -1, err
		***REMOVED***
		args = append(args, oargs...)
	***REMOVED***
	cmd := r.command(context, append(args, id)...)
	if opts != nil && opts.IO != nil ***REMOVED***
		opts.Set(cmd)
	***REMOVED***
	ec, err := Monitor.Start(cmd)
	if err != nil ***REMOVED***
		return -1, err
	***REMOVED***
	return Monitor.Wait(cmd, ec)
***REMOVED***

type DeleteOpts struct ***REMOVED***
	Force bool
***REMOVED***

func (o *DeleteOpts) args() (out []string) ***REMOVED***
	if o.Force ***REMOVED***
		out = append(out, "--force")
	***REMOVED***
	return out
***REMOVED***

// Delete deletes the container
func (r *Runc) Delete(context context.Context, id string, opts *DeleteOpts) error ***REMOVED***
	args := []string***REMOVED***"delete"***REMOVED***
	if opts != nil ***REMOVED***
		args = append(args, opts.args()...)
	***REMOVED***
	return r.runOrError(r.command(context, append(args, id)...))
***REMOVED***

// KillOpts specifies options for killing a container and its processes
type KillOpts struct ***REMOVED***
	All bool
***REMOVED***

func (o *KillOpts) args() (out []string) ***REMOVED***
	if o.All ***REMOVED***
		out = append(out, "--all")
	***REMOVED***
	return out
***REMOVED***

// Kill sends the specified signal to the container
func (r *Runc) Kill(context context.Context, id string, sig int, opts *KillOpts) error ***REMOVED***
	args := []string***REMOVED***
		"kill",
	***REMOVED***
	if opts != nil ***REMOVED***
		args = append(args, opts.args()...)
	***REMOVED***
	return r.runOrError(r.command(context, append(args, id, strconv.Itoa(sig))...))
***REMOVED***

// Stats return the stats for a container like cpu, memory, and io
func (r *Runc) Stats(context context.Context, id string) (*Stats, error) ***REMOVED***
	cmd := r.command(context, "events", "--stats", id)
	rd, err := cmd.StdoutPipe()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ec, err := Monitor.Start(cmd)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer func() ***REMOVED***
		rd.Close()
		Monitor.Wait(cmd, ec)
	***REMOVED***()
	var e Event
	if err := json.NewDecoder(rd).Decode(&e); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return e.Stats, nil
***REMOVED***

// Events returns an event stream from runc for a container with stats and OOM notifications
func (r *Runc) Events(context context.Context, id string, interval time.Duration) (chan *Event, error) ***REMOVED***
	cmd := r.command(context, "events", fmt.Sprintf("--interval=%ds", int(interval.Seconds())), id)
	rd, err := cmd.StdoutPipe()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ec, err := Monitor.Start(cmd)
	if err != nil ***REMOVED***
		rd.Close()
		return nil, err
	***REMOVED***
	var (
		dec = json.NewDecoder(rd)
		c   = make(chan *Event, 128)
	)
	go func() ***REMOVED***
		defer func() ***REMOVED***
			close(c)
			rd.Close()
			Monitor.Wait(cmd, ec)
		***REMOVED***()
		for ***REMOVED***
			var e Event
			if err := dec.Decode(&e); err != nil ***REMOVED***
				if err == io.EOF ***REMOVED***
					return
				***REMOVED***
				e = Event***REMOVED***
					Type: "error",
					Err:  err,
				***REMOVED***
			***REMOVED***
			c <- &e
		***REMOVED***
	***REMOVED***()
	return c, nil
***REMOVED***

// Pause the container with the provided id
func (r *Runc) Pause(context context.Context, id string) error ***REMOVED***
	return r.runOrError(r.command(context, "pause", id))
***REMOVED***

// Resume the container with the provided id
func (r *Runc) Resume(context context.Context, id string) error ***REMOVED***
	return r.runOrError(r.command(context, "resume", id))
***REMOVED***

// Ps lists all the processes inside the container returning their pids
func (r *Runc) Ps(context context.Context, id string) ([]int, error) ***REMOVED***
	data, err := cmdOutput(r.command(context, "ps", "--format", "json", id), true)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("%s: %s", err, data)
	***REMOVED***
	var pids []int
	if err := json.Unmarshal(data, &pids); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return pids, nil
***REMOVED***

type CheckpointOpts struct ***REMOVED***
	// ImagePath is the path for saving the criu image file
	ImagePath string
	// WorkDir is the working directory for criu
	WorkDir string
	// ParentPath is the path for previous image files from a pre-dump
	ParentPath string
	// AllowOpenTCP allows open tcp connections to be checkpointed
	AllowOpenTCP bool
	// AllowExternalUnixSockets allows external unix sockets to be checkpointed
	AllowExternalUnixSockets bool
	// AllowTerminal allows the terminal(pty) to be checkpointed with a container
	AllowTerminal bool
	// CriuPageServer is the address:port for the criu page server
	CriuPageServer string
	// FileLocks handle file locks held by the container
	FileLocks bool
	// Cgroups is the cgroup mode for how to handle the checkpoint of a container's cgroups
	Cgroups CgroupMode
	// EmptyNamespaces creates a namespace for the container but does not save its properties
	// Provide the namespaces you wish to be checkpointed without their settings on restore
	EmptyNamespaces []string
***REMOVED***

type CgroupMode string

const (
	Soft   CgroupMode = "soft"
	Full   CgroupMode = "full"
	Strict CgroupMode = "strict"
)

func (o *CheckpointOpts) args() (out []string) ***REMOVED***
	if o.ImagePath != "" ***REMOVED***
		out = append(out, "--image-path", o.ImagePath)
	***REMOVED***
	if o.WorkDir != "" ***REMOVED***
		out = append(out, "--work-path", o.WorkDir)
	***REMOVED***
	if o.ParentPath != "" ***REMOVED***
		out = append(out, "--parent-path", o.ParentPath)
	***REMOVED***
	if o.AllowOpenTCP ***REMOVED***
		out = append(out, "--tcp-established")
	***REMOVED***
	if o.AllowExternalUnixSockets ***REMOVED***
		out = append(out, "--ext-unix-sk")
	***REMOVED***
	if o.AllowTerminal ***REMOVED***
		out = append(out, "--shell-job")
	***REMOVED***
	if o.CriuPageServer != "" ***REMOVED***
		out = append(out, "--page-server", o.CriuPageServer)
	***REMOVED***
	if o.FileLocks ***REMOVED***
		out = append(out, "--file-locks")
	***REMOVED***
	if string(o.Cgroups) != "" ***REMOVED***
		out = append(out, "--manage-cgroups-mode", string(o.Cgroups))
	***REMOVED***
	for _, ns := range o.EmptyNamespaces ***REMOVED***
		out = append(out, "--empty-ns", ns)
	***REMOVED***
	return out
***REMOVED***

type CheckpointAction func([]string) []string

// LeaveRunning keeps the container running after the checkpoint has been completed
func LeaveRunning(args []string) []string ***REMOVED***
	return append(args, "--leave-running")
***REMOVED***

// PreDump allows a pre-dump of the checkpoint to be made and completed later
func PreDump(args []string) []string ***REMOVED***
	return append(args, "--pre-dump")
***REMOVED***

// Checkpoint allows you to checkpoint a container using criu
func (r *Runc) Checkpoint(context context.Context, id string, opts *CheckpointOpts, actions ...CheckpointAction) error ***REMOVED***
	args := []string***REMOVED***"checkpoint"***REMOVED***
	if opts != nil ***REMOVED***
		args = append(args, opts.args()...)
	***REMOVED***
	for _, a := range actions ***REMOVED***
		args = a(args)
	***REMOVED***
	return r.runOrError(r.command(context, append(args, id)...))
***REMOVED***

type RestoreOpts struct ***REMOVED***
	CheckpointOpts
	IO

	Detach      bool
	PidFile     string
	NoSubreaper bool
	NoPivot     bool
***REMOVED***

func (o *RestoreOpts) args() ([]string, error) ***REMOVED***
	out := o.CheckpointOpts.args()
	if o.Detach ***REMOVED***
		out = append(out, "--detach")
	***REMOVED***
	if o.PidFile != "" ***REMOVED***
		abs, err := filepath.Abs(o.PidFile)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		out = append(out, "--pid-file", abs)
	***REMOVED***
	if o.NoPivot ***REMOVED***
		out = append(out, "--no-pivot")
	***REMOVED***
	if o.NoSubreaper ***REMOVED***
		out = append(out, "-no-subreaper")
	***REMOVED***
	return out, nil
***REMOVED***

// Restore restores a container with the provide id from an existing checkpoint
func (r *Runc) Restore(context context.Context, id, bundle string, opts *RestoreOpts) (int, error) ***REMOVED***
	args := []string***REMOVED***"restore"***REMOVED***
	if opts != nil ***REMOVED***
		oargs, err := opts.args()
		if err != nil ***REMOVED***
			return -1, err
		***REMOVED***
		args = append(args, oargs...)
	***REMOVED***
	args = append(args, "--bundle", bundle)
	cmd := r.command(context, append(args, id)...)
	if opts != nil && opts.IO != nil ***REMOVED***
		opts.Set(cmd)
	***REMOVED***
	ec, err := Monitor.Start(cmd)
	if err != nil ***REMOVED***
		return -1, err
	***REMOVED***
	if opts != nil && opts.IO != nil ***REMOVED***
		if c, ok := opts.IO.(StartCloser); ok ***REMOVED***
			if err := c.CloseAfterStart(); err != nil ***REMOVED***
				return -1, err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return Monitor.Wait(cmd, ec)
***REMOVED***

// Update updates the current container with the provided resource spec
func (r *Runc) Update(context context.Context, id string, resources *specs.LinuxResources) error ***REMOVED***
	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(resources); err != nil ***REMOVED***
		return err
	***REMOVED***
	args := []string***REMOVED***"update", "--resources", "-", id***REMOVED***
	cmd := r.command(context, args...)
	cmd.Stdin = buf
	return r.runOrError(cmd)
***REMOVED***

var ErrParseRuncVersion = errors.New("unable to parse runc version")

type Version struct ***REMOVED***
	Runc   string
	Commit string
	Spec   string
***REMOVED***

// Version returns the runc and runtime-spec versions
func (r *Runc) Version(context context.Context) (Version, error) ***REMOVED***
	data, err := cmdOutput(r.command(context, "--version"), false)
	if err != nil ***REMOVED***
		return Version***REMOVED******REMOVED***, err
	***REMOVED***
	return parseVersion(data)
***REMOVED***

func parseVersion(data []byte) (Version, error) ***REMOVED***
	var v Version
	parts := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(parts) != 3 ***REMOVED***
		return v, ErrParseRuncVersion
	***REMOVED***

	for i, p := range []struct ***REMOVED***
		dest  *string
		split string
	***REMOVED******REMOVED***
		***REMOVED***
			dest:  &v.Runc,
			split: "version ",
		***REMOVED***,
		***REMOVED***
			dest:  &v.Commit,
			split: ": ",
		***REMOVED***,
		***REMOVED***
			dest:  &v.Spec,
			split: ": ",
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		p2 := strings.Split(parts[i], p.split)
		if len(p2) != 2 ***REMOVED***
			return v, fmt.Errorf("unable to parse version line %q", parts[i])
		***REMOVED***
		*p.dest = p2[1]
	***REMOVED***
	return v, nil
***REMOVED***

func (r *Runc) args() (out []string) ***REMOVED***
	if r.Root != "" ***REMOVED***
		out = append(out, "--root", r.Root)
	***REMOVED***
	if r.Debug ***REMOVED***
		out = append(out, "--debug")
	***REMOVED***
	if r.Log != "" ***REMOVED***
		out = append(out, "--log", r.Log)
	***REMOVED***
	if r.LogFormat != none ***REMOVED***
		out = append(out, "--log-format", string(r.LogFormat))
	***REMOVED***
	if r.Criu != "" ***REMOVED***
		out = append(out, "--criu", r.Criu)
	***REMOVED***
	if r.SystemdCgroup ***REMOVED***
		out = append(out, "--systemd-cgroup")
	***REMOVED***
	return out
***REMOVED***

// runOrError will run the provided command.  If an error is
// encountered and neither Stdout or Stderr was set the error and the
// stderr of the command will be returned in the format of <error>:
// <stderr>
func (r *Runc) runOrError(cmd *exec.Cmd) error ***REMOVED***
	if cmd.Stdout != nil || cmd.Stderr != nil ***REMOVED***
		ec, err := Monitor.Start(cmd)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		status, err := Monitor.Wait(cmd, ec)
		if err == nil && status != 0 ***REMOVED***
			err = fmt.Errorf("%s did not terminate sucessfully", cmd.Args[0])
		***REMOVED***
		return err
	***REMOVED***
	data, err := cmdOutput(cmd, true)
	if err != nil ***REMOVED***
		return fmt.Errorf("%s: %s", err, data)
	***REMOVED***
	return nil
***REMOVED***

func cmdOutput(cmd *exec.Cmd, combined bool) ([]byte, error) ***REMOVED***
	var b bytes.Buffer

	cmd.Stdout = &b
	if combined ***REMOVED***
		cmd.Stderr = &b
	***REMOVED***
	ec, err := Monitor.Start(cmd)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	status, err := Monitor.Wait(cmd, ec)
	if err == nil && status != 0 ***REMOVED***
		err = fmt.Errorf("%s did not terminate sucessfully", cmd.Args[0])
	***REMOVED***

	return b.Bytes(), err
***REMOVED***
