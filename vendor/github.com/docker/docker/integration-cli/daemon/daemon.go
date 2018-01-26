package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/request"
	"github.com/docker/docker/opts"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/go-connections/sockets"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

type testingT interface ***REMOVED***
	require.TestingT
	logT
	Fatalf(string, ...interface***REMOVED******REMOVED***)
***REMOVED***

type logT interface ***REMOVED***
	Logf(string, ...interface***REMOVED******REMOVED***)
***REMOVED***

// SockRoot holds the path of the default docker integration daemon socket
var SockRoot = filepath.Join(os.TempDir(), "docker-integration")

var errDaemonNotStarted = errors.New("daemon not started")

// Daemon represents a Docker daemon for the testing framework.
type Daemon struct ***REMOVED***
	GlobalFlags       []string
	Root              string
	Folder            string
	Wait              chan error
	UseDefaultHost    bool
	UseDefaultTLSHost bool

	id             string
	logFile        *os.File
	stdin          io.WriteCloser
	stdout, stderr io.ReadCloser
	cmd            *exec.Cmd
	storageDriver  string
	userlandProxy  bool
	execRoot       string
	experimental   bool
	dockerBinary   string
	dockerdBinary  string
	log            logT
***REMOVED***

// Config holds docker daemon integration configuration
type Config struct ***REMOVED***
	Experimental bool
***REMOVED***

type clientConfig struct ***REMOVED***
	transport *http.Transport
	scheme    string
	addr      string
***REMOVED***

// New returns a Daemon instance to be used for testing.
// This will create a directory such as d123456789 in the folder specified by $DOCKER_INTEGRATION_DAEMON_DEST or $DEST.
// The daemon will not automatically start.
func New(t testingT, dockerBinary string, dockerdBinary string, config Config) *Daemon ***REMOVED***
	dest := os.Getenv("DOCKER_INTEGRATION_DAEMON_DEST")
	if dest == "" ***REMOVED***
		dest = os.Getenv("DEST")
	***REMOVED***
	if dest == "" ***REMOVED***
		t.Fatalf("Please set the DOCKER_INTEGRATION_DAEMON_DEST or the DEST environment variable")
	***REMOVED***

	if err := os.MkdirAll(SockRoot, 0700); err != nil ***REMOVED***
		t.Fatalf("could not create daemon socket root")
	***REMOVED***

	id := fmt.Sprintf("d%s", stringid.TruncateID(stringid.GenerateRandomID()))
	dir := filepath.Join(dest, id)
	daemonFolder, err := filepath.Abs(dir)
	if err != nil ***REMOVED***
		t.Fatalf("Could not make %q an absolute path", dir)
	***REMOVED***
	daemonRoot := filepath.Join(daemonFolder, "root")

	if err := os.MkdirAll(daemonRoot, 0755); err != nil ***REMOVED***
		t.Fatalf("Could not create daemon root %q", dir)
	***REMOVED***

	userlandProxy := true
	if env := os.Getenv("DOCKER_USERLANDPROXY"); env != "" ***REMOVED***
		if val, err := strconv.ParseBool(env); err != nil ***REMOVED***
			userlandProxy = val
		***REMOVED***
	***REMOVED***

	return &Daemon***REMOVED***
		id:            id,
		Folder:        daemonFolder,
		Root:          daemonRoot,
		storageDriver: os.Getenv("DOCKER_GRAPHDRIVER"),
		userlandProxy: userlandProxy,
		execRoot:      filepath.Join(os.TempDir(), "docker-execroot", id),
		dockerBinary:  dockerBinary,
		dockerdBinary: dockerdBinary,
		experimental:  config.Experimental,
		log:           t,
	***REMOVED***
***REMOVED***

// RootDir returns the root directory of the daemon.
func (d *Daemon) RootDir() string ***REMOVED***
	return d.Root
***REMOVED***

// ID returns the generated id of the daemon
func (d *Daemon) ID() string ***REMOVED***
	return d.id
***REMOVED***

// StorageDriver returns the configured storage driver of the daemon
func (d *Daemon) StorageDriver() string ***REMOVED***
	return d.storageDriver
***REMOVED***

// CleanupExecRoot cleans the daemon exec root (network namespaces, ...)
func (d *Daemon) CleanupExecRoot(c *check.C) ***REMOVED***
	cleanupExecRoot(c, d.execRoot)
***REMOVED***

func (d *Daemon) getClientConfig() (*clientConfig, error) ***REMOVED***
	var (
		transport *http.Transport
		scheme    string
		addr      string
		proto     string
	)
	if d.UseDefaultTLSHost ***REMOVED***
		option := &tlsconfig.Options***REMOVED***
			CAFile:   "fixtures/https/ca.pem",
			CertFile: "fixtures/https/client-cert.pem",
			KeyFile:  "fixtures/https/client-key.pem",
		***REMOVED***
		tlsConfig, err := tlsconfig.Client(*option)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		transport = &http.Transport***REMOVED***
			TLSClientConfig: tlsConfig,
		***REMOVED***
		addr = fmt.Sprintf("%s:%d", opts.DefaultHTTPHost, opts.DefaultTLSHTTPPort)
		scheme = "https"
		proto = "tcp"
	***REMOVED*** else if d.UseDefaultHost ***REMOVED***
		addr = opts.DefaultUnixSocket
		proto = "unix"
		scheme = "http"
		transport = &http.Transport***REMOVED******REMOVED***
	***REMOVED*** else ***REMOVED***
		addr = d.sockPath()
		proto = "unix"
		scheme = "http"
		transport = &http.Transport***REMOVED******REMOVED***
	***REMOVED***

	if err := sockets.ConfigureTransport(transport, proto, addr); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	transport.DisableKeepAlives = true

	return &clientConfig***REMOVED***
		transport: transport,
		scheme:    scheme,
		addr:      addr,
	***REMOVED***, nil
***REMOVED***

// Start starts the daemon and return once it is ready to receive requests.
func (d *Daemon) Start(t testingT, args ...string) ***REMOVED***
	if err := d.StartWithError(args...); err != nil ***REMOVED***
		t.Fatalf("Error starting daemon with arguments: %v", args)
	***REMOVED***
***REMOVED***

// StartWithError starts the daemon and return once it is ready to receive requests.
// It returns an error in case it couldn't start.
func (d *Daemon) StartWithError(args ...string) error ***REMOVED***
	logFile, err := os.OpenFile(filepath.Join(d.Folder, "docker.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "[%s] Could not create %s/docker.log", d.id, d.Folder)
	***REMOVED***

	return d.StartWithLogFile(logFile, args...)
***REMOVED***

// StartWithLogFile will start the daemon and attach its streams to a given file.
func (d *Daemon) StartWithLogFile(out *os.File, providedArgs ...string) error ***REMOVED***
	dockerdBinary, err := exec.LookPath(d.dockerdBinary)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "[%s] could not find docker binary in $PATH", d.id)
	***REMOVED***
	args := append(d.GlobalFlags,
		"--containerd", "/var/run/docker/containerd/docker-containerd.sock",
		"--data-root", d.Root,
		"--exec-root", d.execRoot,
		"--pidfile", fmt.Sprintf("%s/docker.pid", d.Folder),
		fmt.Sprintf("--userland-proxy=%t", d.userlandProxy),
	)
	if d.experimental ***REMOVED***
		args = append(args, "--experimental", "--init")
	***REMOVED***
	if !(d.UseDefaultHost || d.UseDefaultTLSHost) ***REMOVED***
		args = append(args, []string***REMOVED***"--host", d.Sock()***REMOVED***...)
	***REMOVED***
	if root := os.Getenv("DOCKER_REMAP_ROOT"); root != "" ***REMOVED***
		args = append(args, []string***REMOVED***"--userns-remap", root***REMOVED***...)
	***REMOVED***

	// If we don't explicitly set the log-level or debug flag(-D) then
	// turn on debug mode
	foundLog := false
	foundSd := false
	for _, a := range providedArgs ***REMOVED***
		if strings.Contains(a, "--log-level") || strings.Contains(a, "-D") || strings.Contains(a, "--debug") ***REMOVED***
			foundLog = true
		***REMOVED***
		if strings.Contains(a, "--storage-driver") ***REMOVED***
			foundSd = true
		***REMOVED***
	***REMOVED***
	if !foundLog ***REMOVED***
		args = append(args, "--debug")
	***REMOVED***
	if d.storageDriver != "" && !foundSd ***REMOVED***
		args = append(args, "--storage-driver", d.storageDriver)
	***REMOVED***

	args = append(args, providedArgs...)
	d.cmd = exec.Command(dockerdBinary, args...)
	d.cmd.Env = append(os.Environ(), "DOCKER_SERVICE_PREFER_OFFLINE_IMAGE=1")
	d.cmd.Stdout = out
	d.cmd.Stderr = out
	d.logFile = out

	if err := d.cmd.Start(); err != nil ***REMOVED***
		return errors.Errorf("[%s] could not start daemon container: %v", d.id, err)
	***REMOVED***

	wait := make(chan error)

	go func() ***REMOVED***
		wait <- d.cmd.Wait()
		d.log.Logf("[%s] exiting daemon", d.id)
		close(wait)
	***REMOVED***()

	d.Wait = wait

	tick := time.Tick(500 * time.Millisecond)
	// make sure daemon is ready to receive requests
	startTime := time.Now().Unix()
	for ***REMOVED***
		d.log.Logf("[%s] waiting for daemon to start", d.id)
		if time.Now().Unix()-startTime > 5 ***REMOVED***
			// After 5 seconds, give up
			return errors.Errorf("[%s] Daemon exited and never started", d.id)
		***REMOVED***
		select ***REMOVED***
		case <-time.After(2 * time.Second):
			return errors.Errorf("[%s] timeout: daemon does not respond", d.id)
		case <-tick:
			clientConfig, err := d.getClientConfig()
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			client := &http.Client***REMOVED***
				Transport: clientConfig.transport,
			***REMOVED***

			req, err := http.NewRequest("GET", "/_ping", nil)
			if err != nil ***REMOVED***
				return errors.Wrapf(err, "[%s] could not create new request", d.id)
			***REMOVED***
			req.URL.Host = clientConfig.addr
			req.URL.Scheme = clientConfig.scheme
			resp, err := client.Do(req)
			if err != nil ***REMOVED***
				continue
			***REMOVED***
			resp.Body.Close()
			if resp.StatusCode != http.StatusOK ***REMOVED***
				d.log.Logf("[%s] received status != 200 OK: %s\n", d.id, resp.Status)
			***REMOVED***
			d.log.Logf("[%s] daemon started\n", d.id)
			d.Root, err = d.queryRootDir()
			if err != nil ***REMOVED***
				return errors.Errorf("[%s] error querying daemon for root directory: %v", d.id, err)
			***REMOVED***
			return nil
		case <-d.Wait:
			return errors.Errorf("[%s] Daemon exited during startup", d.id)
		***REMOVED***
	***REMOVED***
***REMOVED***

// StartWithBusybox will first start the daemon with Daemon.Start()
// then save the busybox image from the main daemon and load it into this Daemon instance.
func (d *Daemon) StartWithBusybox(t testingT, arg ...string) ***REMOVED***
	d.Start(t, arg...)
	d.LoadBusybox(t)
***REMOVED***

// Kill will send a SIGKILL to the daemon
func (d *Daemon) Kill() error ***REMOVED***
	if d.cmd == nil || d.Wait == nil ***REMOVED***
		return errDaemonNotStarted
	***REMOVED***

	defer func() ***REMOVED***
		d.logFile.Close()
		d.cmd = nil
	***REMOVED***()

	if err := d.cmd.Process.Kill(); err != nil ***REMOVED***
		return err
	***REMOVED***

	return os.Remove(fmt.Sprintf("%s/docker.pid", d.Folder))
***REMOVED***

// Pid returns the pid of the daemon
func (d *Daemon) Pid() int ***REMOVED***
	return d.cmd.Process.Pid
***REMOVED***

// Interrupt stops the daemon by sending it an Interrupt signal
func (d *Daemon) Interrupt() error ***REMOVED***
	return d.Signal(os.Interrupt)
***REMOVED***

// Signal sends the specified signal to the daemon if running
func (d *Daemon) Signal(signal os.Signal) error ***REMOVED***
	if d.cmd == nil || d.Wait == nil ***REMOVED***
		return errDaemonNotStarted
	***REMOVED***
	return d.cmd.Process.Signal(signal)
***REMOVED***

// DumpStackAndQuit sends SIGQUIT to the daemon, which triggers it to dump its
// stack to its log file and exit
// This is used primarily for gathering debug information on test timeout
func (d *Daemon) DumpStackAndQuit() ***REMOVED***
	if d.cmd == nil || d.cmd.Process == nil ***REMOVED***
		return
	***REMOVED***
	SignalDaemonDump(d.cmd.Process.Pid)
***REMOVED***

// Stop will send a SIGINT every second and wait for the daemon to stop.
// If it times out, a SIGKILL is sent.
// Stop will not delete the daemon directory. If a purged daemon is needed,
// instantiate a new one with NewDaemon.
// If an error occurs while starting the daemon, the test will fail.
func (d *Daemon) Stop(t testingT) ***REMOVED***
	err := d.StopWithError()
	if err != nil ***REMOVED***
		if err != errDaemonNotStarted ***REMOVED***
			t.Fatalf("Error while stopping the daemon %s : %v", d.id, err)
		***REMOVED*** else ***REMOVED***
			t.Logf("Daemon %s is not started", d.id)
		***REMOVED***
	***REMOVED***
***REMOVED***

// StopWithError will send a SIGINT every second and wait for the daemon to stop.
// If it timeouts, a SIGKILL is sent.
// Stop will not delete the daemon directory. If a purged daemon is needed,
// instantiate a new one with NewDaemon.
func (d *Daemon) StopWithError() error ***REMOVED***
	if d.cmd == nil || d.Wait == nil ***REMOVED***
		return errDaemonNotStarted
	***REMOVED***

	defer func() ***REMOVED***
		d.logFile.Close()
		d.cmd = nil
	***REMOVED***()

	i := 1
	tick := time.Tick(time.Second)

	if err := d.cmd.Process.Signal(os.Interrupt); err != nil ***REMOVED***
		if strings.Contains(err.Error(), "os: process already finished") ***REMOVED***
			return errDaemonNotStarted
		***REMOVED***
		return errors.Errorf("could not send signal: %v", err)
	***REMOVED***
out1:
	for ***REMOVED***
		select ***REMOVED***
		case err := <-d.Wait:
			return err
		case <-time.After(20 * time.Second):
			// time for stopping jobs and run onShutdown hooks
			d.log.Logf("[%s] daemon started", d.id)
			break out1
		***REMOVED***
	***REMOVED***

out2:
	for ***REMOVED***
		select ***REMOVED***
		case err := <-d.Wait:
			return err
		case <-tick:
			i++
			if i > 5 ***REMOVED***
				d.log.Logf("tried to interrupt daemon for %d times, now try to kill it", i)
				break out2
			***REMOVED***
			d.log.Logf("Attempt #%d: daemon is still running with pid %d", i, d.cmd.Process.Pid)
			if err := d.cmd.Process.Signal(os.Interrupt); err != nil ***REMOVED***
				return errors.Errorf("could not send signal: %v", err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if err := d.cmd.Process.Kill(); err != nil ***REMOVED***
		d.log.Logf("Could not kill daemon: %v", err)
		return err
	***REMOVED***

	d.cmd.Wait()

	return os.Remove(fmt.Sprintf("%s/docker.pid", d.Folder))
***REMOVED***

// Restart will restart the daemon by first stopping it and the starting it.
// If an error occurs while starting the daemon, the test will fail.
func (d *Daemon) Restart(t testingT, args ...string) ***REMOVED***
	d.Stop(t)
	d.handleUserns()
	d.Start(t, args...)
***REMOVED***

// RestartWithError will restart the daemon by first stopping it and then starting it.
func (d *Daemon) RestartWithError(arg ...string) error ***REMOVED***
	if err := d.StopWithError(); err != nil ***REMOVED***
		return err
	***REMOVED***
	d.handleUserns()
	return d.StartWithError(arg...)
***REMOVED***

func (d *Daemon) handleUserns() ***REMOVED***
	// in the case of tests running a user namespace-enabled daemon, we have resolved
	// d.Root to be the actual final path of the graph dir after the "uid.gid" of
	// remapped root is added--we need to subtract it from the path before calling
	// start or else we will continue making subdirectories rather than truly restarting
	// with the same location/root:
	if root := os.Getenv("DOCKER_REMAP_ROOT"); root != "" ***REMOVED***
		d.Root = filepath.Dir(d.Root)
	***REMOVED***
***REMOVED***

// LoadBusybox image into the daemon
func (d *Daemon) LoadBusybox(t testingT) ***REMOVED***
	clientHost, err := client.NewEnvClient()
	require.NoError(t, err, "failed to create client")
	defer clientHost.Close()

	ctx := context.Background()
	reader, err := clientHost.ImageSave(ctx, []string***REMOVED***"busybox:latest"***REMOVED***)
	require.NoError(t, err, "failed to download busybox")
	defer reader.Close()

	client, err := d.NewClient()
	require.NoError(t, err, "failed to create client")
	defer client.Close()

	resp, err := client.ImageLoad(ctx, reader, true)
	require.NoError(t, err, "failed to load busybox")
	defer resp.Body.Close()
***REMOVED***

func (d *Daemon) queryRootDir() (string, error) ***REMOVED***
	// update daemon root by asking /info endpoint (to support user
	// namespaced daemon with root remapped uid.gid directory)
	clientConfig, err := d.getClientConfig()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	client := &http.Client***REMOVED***
		Transport: clientConfig.transport,
	***REMOVED***

	req, err := http.NewRequest("GET", "/info", nil)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	req.Header.Set("Content-Type", "application/json")
	req.URL.Host = clientConfig.addr
	req.URL.Scheme = clientConfig.scheme

	resp, err := client.Do(req)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	body := ioutils.NewReadCloserWrapper(resp.Body, func() error ***REMOVED***
		return resp.Body.Close()
	***REMOVED***)

	type Info struct ***REMOVED***
		DockerRootDir string
	***REMOVED***
	var b []byte
	var i Info
	b, err = request.ReadBody(body)
	if err == nil && resp.StatusCode == http.StatusOK ***REMOVED***
		// read the docker root dir
		if err = json.Unmarshal(b, &i); err == nil ***REMOVED***
			return i.DockerRootDir, nil
		***REMOVED***
	***REMOVED***
	return "", err
***REMOVED***

// Sock returns the socket path of the daemon
func (d *Daemon) Sock() string ***REMOVED***
	return fmt.Sprintf("unix://" + d.sockPath())
***REMOVED***

func (d *Daemon) sockPath() string ***REMOVED***
	return filepath.Join(SockRoot, d.id+".sock")
***REMOVED***

// WaitRun waits for a container to be running for 10s
func (d *Daemon) WaitRun(contID string) error ***REMOVED***
	args := []string***REMOVED***"--host", d.Sock()***REMOVED***
	return WaitInspectWithArgs(d.dockerBinary, contID, "***REMOVED******REMOVED***.State.Running***REMOVED******REMOVED***", "true", 10*time.Second, args...)
***REMOVED***

// Info returns the info struct for this daemon
func (d *Daemon) Info(t require.TestingT) types.Info ***REMOVED***
	apiclient, err := request.NewClientForHost(d.Sock())
	require.NoError(t, err)
	info, err := apiclient.Info(context.Background())
	require.NoError(t, err)
	return info
***REMOVED***

// Cmd executes a docker CLI command against this daemon.
// Example: d.Cmd("version") will run docker -H unix://path/to/unix.sock version
func (d *Daemon) Cmd(args ...string) (string, error) ***REMOVED***
	result := icmd.RunCmd(d.Command(args...))
	return result.Combined(), result.Error
***REMOVED***

// Command creates a docker CLI command against this daemon, to be executed later.
// Example: d.Command("version") creates a command to run "docker -H unix://path/to/unix.sock version"
func (d *Daemon) Command(args ...string) icmd.Cmd ***REMOVED***
	return icmd.Command(d.dockerBinary, d.PrependHostArg(args)...)
***REMOVED***

// PrependHostArg prepend the specified arguments by the daemon host flags
func (d *Daemon) PrependHostArg(args []string) []string ***REMOVED***
	for _, arg := range args ***REMOVED***
		if arg == "--host" || arg == "-H" ***REMOVED***
			return args
		***REMOVED***
	***REMOVED***
	return append([]string***REMOVED***"--host", d.Sock()***REMOVED***, args...)
***REMOVED***

// SockRequest executes a socket request on a daemon and returns statuscode and output.
func (d *Daemon) SockRequest(method, endpoint string, data interface***REMOVED******REMOVED***) (int, []byte, error) ***REMOVED***
	jsonData := bytes.NewBuffer(nil)
	if err := json.NewEncoder(jsonData).Encode(data); err != nil ***REMOVED***
		return -1, nil, err
	***REMOVED***

	res, body, err := d.SockRequestRaw(method, endpoint, jsonData, "application/json")
	if err != nil ***REMOVED***
		return -1, nil, err
	***REMOVED***
	b, err := request.ReadBody(body)
	return res.StatusCode, b, err
***REMOVED***

// SockRequestRaw executes a socket request on a daemon and returns an http
// response and a reader for the output data.
// Deprecated: use request package instead
func (d *Daemon) SockRequestRaw(method, endpoint string, data io.Reader, ct string) (*http.Response, io.ReadCloser, error) ***REMOVED***
	return request.SockRequestRaw(method, endpoint, data, ct, d.Sock())
***REMOVED***

// LogFileName returns the path the daemon's log file
func (d *Daemon) LogFileName() string ***REMOVED***
	return d.logFile.Name()
***REMOVED***

// GetIDByName returns the ID of an object (container, volume, …) given its name
func (d *Daemon) GetIDByName(name string) (string, error) ***REMOVED***
	return d.inspectFieldWithError(name, "Id")
***REMOVED***

// ActiveContainers returns the list of ids of the currently running containers
func (d *Daemon) ActiveContainers() (ids []string) ***REMOVED***
	// FIXME(vdemeester) shouldn't ignore the error
	out, _ := d.Cmd("ps", "-q")
	for _, id := range strings.Split(out, "\n") ***REMOVED***
		if id = strings.TrimSpace(id); id != "" ***REMOVED***
			ids = append(ids, id)
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// ReadLogFile returns the content of the daemon log file
func (d *Daemon) ReadLogFile() ([]byte, error) ***REMOVED***
	return ioutil.ReadFile(d.logFile.Name())
***REMOVED***

// InspectField returns the field filter by 'filter'
func (d *Daemon) InspectField(name, filter string) (string, error) ***REMOVED***
	return d.inspectFilter(name, filter)
***REMOVED***

func (d *Daemon) inspectFilter(name, filter string) (string, error) ***REMOVED***
	format := fmt.Sprintf("***REMOVED******REMOVED***%s***REMOVED******REMOVED***", filter)
	out, err := d.Cmd("inspect", "-f", format, name)
	if err != nil ***REMOVED***
		return "", errors.Errorf("failed to inspect %s: %s", name, out)
	***REMOVED***
	return strings.TrimSpace(out), nil
***REMOVED***

func (d *Daemon) inspectFieldWithError(name, field string) (string, error) ***REMOVED***
	return d.inspectFilter(name, fmt.Sprintf(".%s", field))
***REMOVED***

// FindContainerIP returns the ip of the specified container
func (d *Daemon) FindContainerIP(id string) (string, error) ***REMOVED***
	out, err := d.Cmd("inspect", "--format='***REMOVED******REMOVED*** .NetworkSettings.Networks.bridge.IPAddress ***REMOVED******REMOVED***'", id)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return strings.Trim(out, " \r\n'"), nil
***REMOVED***

// BuildImageWithOut builds an image with the specified dockerfile and options and returns the output
func (d *Daemon) BuildImageWithOut(name, dockerfile string, useCache bool, buildFlags ...string) (string, int, error) ***REMOVED***
	buildCmd := BuildImageCmdWithHost(d.dockerBinary, name, dockerfile, d.Sock(), useCache, buildFlags...)
	result := icmd.RunCmd(icmd.Cmd***REMOVED***
		Command: buildCmd.Args,
		Env:     buildCmd.Env,
		Dir:     buildCmd.Dir,
		Stdin:   buildCmd.Stdin,
		Stdout:  buildCmd.Stdout,
	***REMOVED***)
	return result.Combined(), result.ExitCode, result.Error
***REMOVED***

// CheckActiveContainerCount returns the number of active containers
// FIXME(vdemeester) should re-use ActivateContainers in some way
func (d *Daemon) CheckActiveContainerCount(c *check.C) (interface***REMOVED******REMOVED***, check.CommentInterface) ***REMOVED***
	out, err := d.Cmd("ps", "-q")
	c.Assert(err, checker.IsNil)
	if len(strings.TrimSpace(out)) == 0 ***REMOVED***
		return 0, nil
	***REMOVED***
	return len(strings.Split(strings.TrimSpace(out), "\n")), check.Commentf("output: %q", string(out))
***REMOVED***

// ReloadConfig asks the daemon to reload its configuration
func (d *Daemon) ReloadConfig() error ***REMOVED***
	if d.cmd == nil || d.cmd.Process == nil ***REMOVED***
		return errors.New("daemon is not running")
	***REMOVED***

	errCh := make(chan error)
	started := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		_, body, err := request.DoOnHost(d.Sock(), "/events", request.Method(http.MethodGet))
		close(started)
		if err != nil ***REMOVED***
			errCh <- err
		***REMOVED***
		defer body.Close()
		dec := json.NewDecoder(body)
		for ***REMOVED***
			var e events.Message
			if err := dec.Decode(&e); err != nil ***REMOVED***
				errCh <- err
				return
			***REMOVED***
			if e.Type != events.DaemonEventType ***REMOVED***
				continue
			***REMOVED***
			if e.Action != "reload" ***REMOVED***
				continue
			***REMOVED***
			close(errCh) // notify that we are done
			return
		***REMOVED***
	***REMOVED***()

	<-started
	if err := signalDaemonReload(d.cmd.Process.Pid); err != nil ***REMOVED***
		return errors.Errorf("error signaling daemon reload: %v", err)
	***REMOVED***
	select ***REMOVED***
	case err := <-errCh:
		if err != nil ***REMOVED***
			return errors.Errorf("error waiting for daemon reload event: %v", err)
		***REMOVED***
	case <-time.After(30 * time.Second):
		return errors.New("timeout waiting for daemon reload event")
	***REMOVED***
	return nil
***REMOVED***

// NewClient creates new client based on daemon's socket path
func (d *Daemon) NewClient() (*client.Client, error) ***REMOVED***
	httpClient, err := request.NewHTTPClient(d.Sock())
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return client.NewClient(d.Sock(), api.DefaultVersion, httpClient, nil)
***REMOVED***

// WaitInspectWithArgs waits for the specified expression to be equals to the specified expected string in the given time.
// Deprecated: use cli.WaitCmd instead
func WaitInspectWithArgs(dockerBinary, name, expr, expected string, timeout time.Duration, arg ...string) error ***REMOVED***
	after := time.After(timeout)

	args := append(arg, "inspect", "-f", expr, name)
	for ***REMOVED***
		result := icmd.RunCommand(dockerBinary, args...)
		if result.Error != nil ***REMOVED***
			if !strings.Contains(strings.ToLower(result.Stderr()), "no such") ***REMOVED***
				return errors.Errorf("error executing docker inspect: %v\n%s",
					result.Stderr(), result.Stdout())
			***REMOVED***
			select ***REMOVED***
			case <-after:
				return result.Error
			default:
				time.Sleep(10 * time.Millisecond)
				continue
			***REMOVED***
		***REMOVED***

		out := strings.TrimSpace(result.Stdout())
		if out == expected ***REMOVED***
			break
		***REMOVED***

		select ***REMOVED***
		case <-after:
			return errors.Errorf("condition \"%q == %q\" not true in time (%v)", out, expected, timeout)
		default:
		***REMOVED***

		time.Sleep(100 * time.Millisecond)
	***REMOVED***
	return nil
***REMOVED***

// BuildImageCmdWithHost create a build command with the specified arguments.
// Deprecated
// FIXME(vdemeester) move this away
func BuildImageCmdWithHost(dockerBinary, name, dockerfile, host string, useCache bool, buildFlags ...string) *exec.Cmd ***REMOVED***
	args := []string***REMOVED******REMOVED***
	if host != "" ***REMOVED***
		args = append(args, "--host", host)
	***REMOVED***
	args = append(args, "build", "-t", name)
	if !useCache ***REMOVED***
		args = append(args, "--no-cache")
	***REMOVED***
	args = append(args, buildFlags...)
	args = append(args, "-")
	buildCmd := exec.Command(dockerBinary, args...)
	buildCmd.Stdin = strings.NewReader(dockerfile)
	return buildCmd
***REMOVED***
