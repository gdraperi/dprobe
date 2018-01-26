// +build !windows

package libcontainerd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/server"
	"github.com/docker/docker/pkg/system"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	maxConnectionRetryCount = 3
	healthCheckTimeout      = 3 * time.Second
	shutdownTimeout         = 15 * time.Second
	configFile              = "containerd.toml"
	binaryName              = "docker-containerd"
	pidFile                 = "docker-containerd.pid"
)

type pluginConfigs struct ***REMOVED***
	Plugins map[string]interface***REMOVED******REMOVED*** `toml:"plugins"`
***REMOVED***

type remote struct ***REMOVED***
	sync.RWMutex
	server.Config

	daemonPid int
	logger    *logrus.Entry

	daemonWaitCh    chan struct***REMOVED******REMOVED***
	clients         []*client
	shutdownContext context.Context
	shutdownCancel  context.CancelFunc
	shutdown        bool

	// Options
	startDaemon bool
	rootDir     string
	stateDir    string
	snapshotter string
	pluginConfs pluginConfigs
***REMOVED***

// New creates a fresh instance of libcontainerd remote.
func New(rootDir, stateDir string, options ...RemoteOption) (rem Remote, err error) ***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			err = errors.Wrap(err, "Failed to connect to containerd")
		***REMOVED***
	***REMOVED***()

	r := &remote***REMOVED***
		rootDir:  rootDir,
		stateDir: stateDir,
		Config: server.Config***REMOVED***
			Root:  filepath.Join(rootDir, "daemon"),
			State: filepath.Join(stateDir, "daemon"),
		***REMOVED***,
		pluginConfs: pluginConfigs***REMOVED***make(map[string]interface***REMOVED******REMOVED***)***REMOVED***,
		daemonPid:   -1,
		logger:      logrus.WithField("module", "libcontainerd"),
	***REMOVED***
	r.shutdownContext, r.shutdownCancel = context.WithCancel(context.Background())

	rem = r
	for _, option := range options ***REMOVED***
		if err = option.Apply(r); err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	r.setDefaults()

	if err = system.MkdirAll(stateDir, 0700, ""); err != nil ***REMOVED***
		return
	***REMOVED***

	if r.startDaemon ***REMOVED***
		os.Remove(r.GRPC.Address)
		if err = r.startContainerd(); err != nil ***REMOVED***
			return
		***REMOVED***
		defer func() ***REMOVED***
			if err != nil ***REMOVED***
				r.Cleanup()
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	// This connection is just used to monitor the connection
	client, err := containerd.New(r.GRPC.Address)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if _, err := client.Version(context.Background()); err != nil ***REMOVED***
		system.KillProcess(r.daemonPid)
		return nil, errors.Wrapf(err, "unable to get containerd version")
	***REMOVED***

	go r.monitorConnection(client)

	return r, nil
***REMOVED***

func (r *remote) NewClient(ns string, b Backend) (Client, error) ***REMOVED***
	c := &client***REMOVED***
		stateDir:   r.stateDir,
		logger:     r.logger.WithField("namespace", ns),
		namespace:  ns,
		backend:    b,
		containers: make(map[string]*container),
	***REMOVED***

	rclient, err := containerd.New(r.GRPC.Address, containerd.WithDefaultNamespace(ns))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	c.remote = rclient

	go c.processEventStream(r.shutdownContext)

	r.Lock()
	r.clients = append(r.clients, c)
	r.Unlock()
	return c, nil
***REMOVED***

func (r *remote) Cleanup() ***REMOVED***
	if r.daemonPid != -1 ***REMOVED***
		r.shutdownCancel()
		r.stopDaemon()
	***REMOVED***

	// cleanup some files
	os.Remove(filepath.Join(r.stateDir, pidFile))

	r.platformCleanup()
***REMOVED***

func (r *remote) getContainerdPid() (int, error) ***REMOVED***
	pidFile := filepath.Join(r.stateDir, pidFile)
	f, err := os.OpenFile(pidFile, os.O_RDWR, 0600)
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return -1, nil
		***REMOVED***
		return -1, err
	***REMOVED***
	defer f.Close()

	b := make([]byte, 8)
	n, err := f.Read(b)
	if err != nil && err != io.EOF ***REMOVED***
		return -1, err
	***REMOVED***

	if n > 0 ***REMOVED***
		pid, err := strconv.ParseUint(string(b[:n]), 10, 64)
		if err != nil ***REMOVED***
			return -1, err
		***REMOVED***
		if system.IsProcessAlive(int(pid)) ***REMOVED***
			return int(pid), nil
		***REMOVED***
	***REMOVED***

	return -1, nil
***REMOVED***

func (r *remote) getContainerdConfig() (string, error) ***REMOVED***
	path := filepath.Join(r.stateDir, configFile)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil ***REMOVED***
		return "", errors.Wrapf(err, "failed to open containerd config file at %s", path)
	***REMOVED***
	defer f.Close()

	enc := toml.NewEncoder(f)
	if err = enc.Encode(r.Config); err != nil ***REMOVED***
		return "", errors.Wrapf(err, "failed to encode general config")
	***REMOVED***
	if err = enc.Encode(r.pluginConfs); err != nil ***REMOVED***
		return "", errors.Wrapf(err, "failed to encode plugin configs")
	***REMOVED***

	return path, nil
***REMOVED***

func (r *remote) startContainerd() error ***REMOVED***
	pid, err := r.getContainerdPid()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if pid != -1 ***REMOVED***
		r.daemonPid = pid
		logrus.WithField("pid", pid).
			Infof("libcontainerd: %s is still running", binaryName)
		return nil
	***REMOVED***

	configFile, err := r.getContainerdConfig()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	args := []string***REMOVED***"--config", configFile***REMOVED***
	cmd := exec.Command(binaryName, args...)
	// redirect containerd logs to docker logs
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = containerdSysProcAttr()
	// clear the NOTIFY_SOCKET from the env when starting containerd
	cmd.Env = nil
	for _, e := range os.Environ() ***REMOVED***
		if !strings.HasPrefix(e, "NOTIFY_SOCKET") ***REMOVED***
			cmd.Env = append(cmd.Env, e)
		***REMOVED***
	***REMOVED***
	if err := cmd.Start(); err != nil ***REMOVED***
		return err
	***REMOVED***

	r.daemonWaitCh = make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		// Reap our child when needed
		if err := cmd.Wait(); err != nil ***REMOVED***
			r.logger.WithError(err).Errorf("containerd did not exit successfully")
		***REMOVED***
		close(r.daemonWaitCh)
	***REMOVED***()

	r.daemonPid = cmd.Process.Pid

	err = ioutil.WriteFile(filepath.Join(r.stateDir, pidFile), []byte(fmt.Sprintf("%d", r.daemonPid)), 0660)
	if err != nil ***REMOVED***
		system.KillProcess(r.daemonPid)
		return errors.Wrap(err, "libcontainerd: failed to save daemon pid to disk")
	***REMOVED***

	logrus.WithField("pid", r.daemonPid).
		Infof("libcontainerd: started new %s process", binaryName)

	return nil
***REMOVED***

func (r *remote) monitorConnection(client *containerd.Client) ***REMOVED***
	var transientFailureCount = 0

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for ***REMOVED***
		<-ticker.C
		ctx, cancel := context.WithTimeout(r.shutdownContext, healthCheckTimeout)
		_, err := client.IsServing(ctx)
		cancel()
		if err == nil ***REMOVED***
			transientFailureCount = 0
			continue
		***REMOVED***

		select ***REMOVED***
		case <-r.shutdownContext.Done():
			r.logger.Info("stopping healthcheck following graceful shutdown")
			client.Close()
			return
		default:
		***REMOVED***

		r.logger.WithError(err).WithField("binary", binaryName).Debug("daemon is not responding")

		if r.daemonPid != -1 ***REMOVED***
			transientFailureCount++
			if transientFailureCount >= maxConnectionRetryCount || !system.IsProcessAlive(r.daemonPid) ***REMOVED***
				transientFailureCount = 0
				if system.IsProcessAlive(r.daemonPid) ***REMOVED***
					r.logger.WithField("pid", r.daemonPid).Info("killing and restarting containerd")
					// Try to get a stack trace
					syscall.Kill(r.daemonPid, syscall.SIGUSR1)
					<-time.After(100 * time.Millisecond)
					system.KillProcess(r.daemonPid)
				***REMOVED***
				<-r.daemonWaitCh
				var err error
				client.Close()
				os.Remove(r.GRPC.Address)
				if err = r.startContainerd(); err != nil ***REMOVED***
					r.logger.WithError(err).Error("failed restarting containerd")
				***REMOVED*** else ***REMOVED***
					newClient, err := containerd.New(r.GRPC.Address)
					if err != nil ***REMOVED***
						r.logger.WithError(err).Error("failed connect to containerd")
					***REMOVED*** else ***REMOVED***
						client = newClient
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
