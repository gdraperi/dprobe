// +build remote_daemon

package libcontainerd

import (
	"os"
)

const (
	grpcPipeName  = `\\.\pipe\docker-containerd-containerd`
	debugPipeName = `\\.\pipe\docker-containerd-debug`
)

func (r *remote) setDefaults() ***REMOVED***
	if r.GRPC.Address == "" ***REMOVED***
		r.GRPC.Address = grpcPipeName
	***REMOVED***
	if r.Debug.Address == "" ***REMOVED***
		r.Debug.Address = debugPipeName
	***REMOVED***
	if r.Debug.Level == "" ***REMOVED***
		r.Debug.Level = "info"
	***REMOVED***
	if r.snapshotter == "" ***REMOVED***
		r.snapshotter = "naive" // TODO(mlaventure): switch to "windows" once implemented
	***REMOVED***
***REMOVED***

func (r *remote) stopDaemon() ***REMOVED***
	p, err := os.FindProcess(r.daemonPid)
	if err != nil ***REMOVED***
		r.logger.WithField("pid", r.daemonPid).Warn("could not find daemon process")
		return
	***REMOVED***

	if err = p.Kill(); err != nil ***REMOVED***
		r.logger.WithError(err).WithField("pid", r.daemonPid).Warn("could not kill daemon process")
		return
	***REMOVED***

	_, err = p.Wait()
	if err != nil ***REMOVED***
		r.logger.WithError(err).WithField("pid", r.daemonPid).Warn("wait for daemon process")
		return
	***REMOVED***
***REMOVED***

func (r *remote) platformCleanup() ***REMOVED***
	// Nothing to do
***REMOVED***
