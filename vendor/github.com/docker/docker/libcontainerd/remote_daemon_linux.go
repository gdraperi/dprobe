package libcontainerd

import (
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/docker/docker/pkg/system"
)

const (
	sockFile      = "docker-containerd.sock"
	debugSockFile = "docker-containerd-debug.sock"
)

func (r *remote) setDefaults() ***REMOVED***
	if r.GRPC.Address == "" ***REMOVED***
		r.GRPC.Address = filepath.Join(r.stateDir, sockFile)
	***REMOVED***
	if r.Debug.Address == "" ***REMOVED***
		r.Debug.Address = filepath.Join(r.stateDir, debugSockFile)
	***REMOVED***
	if r.Debug.Level == "" ***REMOVED***
		r.Debug.Level = "info"
	***REMOVED***
	if r.OOMScore == 0 ***REMOVED***
		r.OOMScore = -999
	***REMOVED***
	if r.snapshotter == "" ***REMOVED***
		r.snapshotter = "overlay"
	***REMOVED***
***REMOVED***

func (r *remote) stopDaemon() ***REMOVED***
	// Ask the daemon to quit
	syscall.Kill(r.daemonPid, syscall.SIGTERM)
	// Wait up to 15secs for it to stop
	for i := time.Duration(0); i < shutdownTimeout; i += time.Second ***REMOVED***
		if !system.IsProcessAlive(r.daemonPid) ***REMOVED***
			break
		***REMOVED***
		time.Sleep(time.Second)
	***REMOVED***

	if system.IsProcessAlive(r.daemonPid) ***REMOVED***
		r.logger.WithField("pid", r.daemonPid).Warn("daemon didn't stop within 15 secs, killing it")
		syscall.Kill(r.daemonPid, syscall.SIGKILL)
	***REMOVED***
***REMOVED***

func (r *remote) platformCleanup() ***REMOVED***
	os.Remove(filepath.Join(r.stateDir, sockFile))
***REMOVED***
