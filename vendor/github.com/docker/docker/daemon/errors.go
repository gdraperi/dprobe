package daemon

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/docker/docker/errdefs"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func errNotRunning(id string) error ***REMOVED***
	return errdefs.Conflict(errors.Errorf("Container %s is not running", id))
***REMOVED***

func containerNotFound(id string) error ***REMOVED***
	return objNotFoundError***REMOVED***"container", id***REMOVED***
***REMOVED***

func volumeNotFound(id string) error ***REMOVED***
	return objNotFoundError***REMOVED***"volume", id***REMOVED***
***REMOVED***

type objNotFoundError struct ***REMOVED***
	object string
	id     string
***REMOVED***

func (e objNotFoundError) Error() string ***REMOVED***
	return "No such " + e.object + ": " + e.id
***REMOVED***

func (e objNotFoundError) NotFound() ***REMOVED******REMOVED***

func errContainerIsRestarting(containerID string) error ***REMOVED***
	cause := errors.Errorf("Container %s is restarting, wait until the container is running", containerID)
	return errdefs.Conflict(cause)
***REMOVED***

func errExecNotFound(id string) error ***REMOVED***
	return objNotFoundError***REMOVED***"exec instance", id***REMOVED***
***REMOVED***

func errExecPaused(id string) error ***REMOVED***
	cause := errors.Errorf("Container %s is paused, unpause the container before exec", id)
	return errdefs.Conflict(cause)
***REMOVED***

func errNotPaused(id string) error ***REMOVED***
	cause := errors.Errorf("Container %s is already paused", id)
	return errdefs.Conflict(cause)
***REMOVED***

type nameConflictError struct ***REMOVED***
	id   string
	name string
***REMOVED***

func (e nameConflictError) Error() string ***REMOVED***
	return fmt.Sprintf("Conflict. The container name %q is already in use by container %q. You have to remove (or rename) that container to be able to reuse that name.", e.name, e.id)
***REMOVED***

func (nameConflictError) Conflict() ***REMOVED******REMOVED***

type containerNotModifiedError struct ***REMOVED***
	running bool
***REMOVED***

func (e containerNotModifiedError) Error() string ***REMOVED***
	if e.running ***REMOVED***
		return "Container is already started"
	***REMOVED***
	return "Container is already stopped"
***REMOVED***

func (e containerNotModifiedError) NotModified() ***REMOVED******REMOVED***

type invalidIdentifier string

func (e invalidIdentifier) Error() string ***REMOVED***
	return fmt.Sprintf("invalid name or ID supplied: %q", string(e))
***REMOVED***

func (invalidIdentifier) InvalidParameter() ***REMOVED******REMOVED***

type duplicateMountPointError string

func (e duplicateMountPointError) Error() string ***REMOVED***
	return "Duplicate mount point: " + string(e)
***REMOVED***
func (duplicateMountPointError) InvalidParameter() ***REMOVED******REMOVED***

type containerFileNotFound struct ***REMOVED***
	file      string
	container string
***REMOVED***

func (e containerFileNotFound) Error() string ***REMOVED***
	return "Could not find the file " + e.file + " in container " + e.container
***REMOVED***

func (containerFileNotFound) NotFound() ***REMOVED******REMOVED***

type invalidFilter struct ***REMOVED***
	filter string
	value  interface***REMOVED******REMOVED***
***REMOVED***

func (e invalidFilter) Error() string ***REMOVED***
	msg := "Invalid filter '" + e.filter
	if e.value != nil ***REMOVED***
		msg += fmt.Sprintf("=%s", e.value)
	***REMOVED***
	return msg + "'"
***REMOVED***

func (e invalidFilter) InvalidParameter() ***REMOVED******REMOVED***

type startInvalidConfigError string

func (e startInvalidConfigError) Error() string ***REMOVED***
	return string(e)
***REMOVED***

func (e startInvalidConfigError) InvalidParameter() ***REMOVED******REMOVED*** // Is this right???

func translateContainerdStartErr(cmd string, setExitCode func(int), err error) error ***REMOVED***
	errDesc := grpc.ErrorDesc(err)
	contains := func(s1, s2 string) bool ***REMOVED***
		return strings.Contains(strings.ToLower(s1), s2)
	***REMOVED***
	var retErr = errdefs.Unknown(errors.New(errDesc))
	// if we receive an internal error from the initial start of a container then lets
	// return it instead of entering the restart loop
	// set to 127 for container cmd not found/does not exist)
	if contains(errDesc, cmd) &&
		(contains(errDesc, "executable file not found") ||
			contains(errDesc, "no such file or directory") ||
			contains(errDesc, "system cannot find the file specified")) ***REMOVED***
		setExitCode(127)
		retErr = startInvalidConfigError(errDesc)
	***REMOVED***
	// set to 126 for container cmd can't be invoked errors
	if contains(errDesc, syscall.EACCES.Error()) ***REMOVED***
		setExitCode(126)
		retErr = startInvalidConfigError(errDesc)
	***REMOVED***

	// attempted to mount a file onto a directory, or a directory onto a file, maybe from user specified bind mounts
	if contains(errDesc, syscall.ENOTDIR.Error()) ***REMOVED***
		errDesc += ": Are you trying to mount a directory onto a file (or vice-versa)? Check if the specified host path exists and is the expected type"
		setExitCode(127)
		retErr = startInvalidConfigError(errDesc)
	***REMOVED***

	// TODO: it would be nice to get some better errors from containerd so we can return better errors here
	return retErr
***REMOVED***
