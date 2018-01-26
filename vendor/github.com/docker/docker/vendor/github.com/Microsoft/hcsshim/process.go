package hcsshim

import (
	"encoding/json"
	"io"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

// ContainerError is an error encountered in HCS
type process struct ***REMOVED***
	handleLock     sync.RWMutex
	handle         hcsProcess
	processID      int
	container      *container
	cachedPipes    *cachedPipes
	callbackNumber uintptr
***REMOVED***

type cachedPipes struct ***REMOVED***
	stdIn  syscall.Handle
	stdOut syscall.Handle
	stdErr syscall.Handle
***REMOVED***

type processModifyRequest struct ***REMOVED***
	Operation   string
	ConsoleSize *consoleSize `json:",omitempty"`
	CloseHandle *closeHandle `json:",omitempty"`
***REMOVED***

type consoleSize struct ***REMOVED***
	Height uint16
	Width  uint16
***REMOVED***

type closeHandle struct ***REMOVED***
	Handle string
***REMOVED***

type processStatus struct ***REMOVED***
	ProcessID      uint32
	Exited         bool
	ExitCode       uint32
	LastWaitResult int32
***REMOVED***

const (
	stdIn  string = "StdIn"
	stdOut string = "StdOut"
	stdErr string = "StdErr"
)

const (
	modifyConsoleSize string = "ConsoleSize"
	modifyCloseHandle string = "CloseHandle"
)

// Pid returns the process ID of the process within the container.
func (process *process) Pid() int ***REMOVED***
	return process.processID
***REMOVED***

// Kill signals the process to terminate but does not wait for it to finish terminating.
func (process *process) Kill() error ***REMOVED***
	process.handleLock.RLock()
	defer process.handleLock.RUnlock()
	operation := "Kill"
	title := "HCSShim::Process::" + operation
	logrus.Debugf(title+" processid=%d", process.processID)

	if process.handle == 0 ***REMOVED***
		return makeProcessError(process, operation, "", ErrAlreadyClosed)
	***REMOVED***

	var resultp *uint16
	err := hcsTerminateProcess(process.handle, &resultp)
	err = processHcsResult(err, resultp)
	if err != nil ***REMOVED***
		return makeProcessError(process, operation, "", err)
	***REMOVED***

	logrus.Debugf(title+" succeeded processid=%d", process.processID)
	return nil
***REMOVED***

// Wait waits for the process to exit.
func (process *process) Wait() error ***REMOVED***
	operation := "Wait"
	title := "HCSShim::Process::" + operation
	logrus.Debugf(title+" processid=%d", process.processID)

	err := waitForNotification(process.callbackNumber, hcsNotificationProcessExited, nil)
	if err != nil ***REMOVED***
		return makeProcessError(process, operation, "", err)
	***REMOVED***

	logrus.Debugf(title+" succeeded processid=%d", process.processID)
	return nil
***REMOVED***

// WaitTimeout waits for the process to exit or the duration to elapse. It returns
// false if timeout occurs.
func (process *process) WaitTimeout(timeout time.Duration) error ***REMOVED***
	operation := "WaitTimeout"
	title := "HCSShim::Process::" + operation
	logrus.Debugf(title+" processid=%d", process.processID)

	err := waitForNotification(process.callbackNumber, hcsNotificationProcessExited, &timeout)
	if err != nil ***REMOVED***
		return makeProcessError(process, operation, "", err)
	***REMOVED***

	logrus.Debugf(title+" succeeded processid=%d", process.processID)
	return nil
***REMOVED***

// ExitCode returns the exit code of the process. The process must have
// already terminated.
func (process *process) ExitCode() (int, error) ***REMOVED***
	process.handleLock.RLock()
	defer process.handleLock.RUnlock()
	operation := "ExitCode"
	title := "HCSShim::Process::" + operation
	logrus.Debugf(title+" processid=%d", process.processID)

	if process.handle == 0 ***REMOVED***
		return 0, makeProcessError(process, operation, "", ErrAlreadyClosed)
	***REMOVED***

	properties, err := process.properties()
	if err != nil ***REMOVED***
		return 0, makeProcessError(process, operation, "", err)
	***REMOVED***

	if properties.Exited == false ***REMOVED***
		return 0, makeProcessError(process, operation, "", ErrInvalidProcessState)
	***REMOVED***

	if properties.LastWaitResult != 0 ***REMOVED***
		return 0, makeProcessError(process, operation, "", syscall.Errno(properties.LastWaitResult))
	***REMOVED***

	logrus.Debugf(title+" succeeded processid=%d exitCode=%d", process.processID, properties.ExitCode)
	return int(properties.ExitCode), nil
***REMOVED***

// ResizeConsole resizes the console of the process.
func (process *process) ResizeConsole(width, height uint16) error ***REMOVED***
	process.handleLock.RLock()
	defer process.handleLock.RUnlock()
	operation := "ResizeConsole"
	title := "HCSShim::Process::" + operation
	logrus.Debugf(title+" processid=%d", process.processID)

	if process.handle == 0 ***REMOVED***
		return makeProcessError(process, operation, "", ErrAlreadyClosed)
	***REMOVED***

	modifyRequest := processModifyRequest***REMOVED***
		Operation: modifyConsoleSize,
		ConsoleSize: &consoleSize***REMOVED***
			Height: height,
			Width:  width,
		***REMOVED***,
	***REMOVED***

	modifyRequestb, err := json.Marshal(modifyRequest)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	modifyRequestStr := string(modifyRequestb)

	var resultp *uint16
	err = hcsModifyProcess(process.handle, modifyRequestStr, &resultp)
	err = processHcsResult(err, resultp)
	if err != nil ***REMOVED***
		return makeProcessError(process, operation, "", err)
	***REMOVED***

	logrus.Debugf(title+" succeeded processid=%d", process.processID)
	return nil
***REMOVED***

func (process *process) properties() (*processStatus, error) ***REMOVED***
	operation := "properties"
	title := "HCSShim::Process::" + operation
	logrus.Debugf(title+" processid=%d", process.processID)

	var (
		resultp     *uint16
		propertiesp *uint16
	)
	err := hcsGetProcessProperties(process.handle, &propertiesp, &resultp)
	err = processHcsResult(err, resultp)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if propertiesp == nil ***REMOVED***
		return nil, ErrUnexpectedValue
	***REMOVED***
	propertiesRaw := convertAndFreeCoTaskMemBytes(propertiesp)

	properties := &processStatus***REMOVED******REMOVED***
	if err := json.Unmarshal(propertiesRaw, properties); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	logrus.Debugf(title+" succeeded processid=%d, properties=%s", process.processID, propertiesRaw)
	return properties, nil
***REMOVED***

// Stdio returns the stdin, stdout, and stderr pipes, respectively. Closing
// these pipes does not close the underlying pipes; it should be possible to
// call this multiple times to get multiple interfaces.
func (process *process) Stdio() (io.WriteCloser, io.ReadCloser, io.ReadCloser, error) ***REMOVED***
	process.handleLock.RLock()
	defer process.handleLock.RUnlock()
	operation := "Stdio"
	title := "HCSShim::Process::" + operation
	logrus.Debugf(title+" processid=%d", process.processID)

	if process.handle == 0 ***REMOVED***
		return nil, nil, nil, makeProcessError(process, operation, "", ErrAlreadyClosed)
	***REMOVED***

	var stdIn, stdOut, stdErr syscall.Handle

	if process.cachedPipes == nil ***REMOVED***
		var (
			processInfo hcsProcessInformation
			resultp     *uint16
		)
		err := hcsGetProcessInfo(process.handle, &processInfo, &resultp)
		err = processHcsResult(err, resultp)
		if err != nil ***REMOVED***
			return nil, nil, nil, makeProcessError(process, operation, "", err)
		***REMOVED***

		stdIn, stdOut, stdErr = processInfo.StdInput, processInfo.StdOutput, processInfo.StdError
	***REMOVED*** else ***REMOVED***
		// Use cached pipes
		stdIn, stdOut, stdErr = process.cachedPipes.stdIn, process.cachedPipes.stdOut, process.cachedPipes.stdErr

		// Invalidate the cache
		process.cachedPipes = nil
	***REMOVED***

	pipes, err := makeOpenFiles([]syscall.Handle***REMOVED***stdIn, stdOut, stdErr***REMOVED***)
	if err != nil ***REMOVED***
		return nil, nil, nil, makeProcessError(process, operation, "", err)
	***REMOVED***

	logrus.Debugf(title+" succeeded processid=%d", process.processID)
	return pipes[0], pipes[1], pipes[2], nil
***REMOVED***

// CloseStdin closes the write side of the stdin pipe so that the process is
// notified on the read side that there is no more data in stdin.
func (process *process) CloseStdin() error ***REMOVED***
	process.handleLock.RLock()
	defer process.handleLock.RUnlock()
	operation := "CloseStdin"
	title := "HCSShim::Process::" + operation
	logrus.Debugf(title+" processid=%d", process.processID)

	if process.handle == 0 ***REMOVED***
		return makeProcessError(process, operation, "", ErrAlreadyClosed)
	***REMOVED***

	modifyRequest := processModifyRequest***REMOVED***
		Operation: modifyCloseHandle,
		CloseHandle: &closeHandle***REMOVED***
			Handle: stdIn,
		***REMOVED***,
	***REMOVED***

	modifyRequestb, err := json.Marshal(modifyRequest)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	modifyRequestStr := string(modifyRequestb)

	var resultp *uint16
	err = hcsModifyProcess(process.handle, modifyRequestStr, &resultp)
	err = processHcsResult(err, resultp)
	if err != nil ***REMOVED***
		return makeProcessError(process, operation, "", err)
	***REMOVED***

	logrus.Debugf(title+" succeeded processid=%d", process.processID)
	return nil
***REMOVED***

// Close cleans up any state associated with the process but does not kill
// or wait on it.
func (process *process) Close() error ***REMOVED***
	process.handleLock.Lock()
	defer process.handleLock.Unlock()
	operation := "Close"
	title := "HCSShim::Process::" + operation
	logrus.Debugf(title+" processid=%d", process.processID)

	// Don't double free this
	if process.handle == 0 ***REMOVED***
		return nil
	***REMOVED***

	if err := process.unregisterCallback(); err != nil ***REMOVED***
		return makeProcessError(process, operation, "", err)
	***REMOVED***

	if err := hcsCloseProcess(process.handle); err != nil ***REMOVED***
		return makeProcessError(process, operation, "", err)
	***REMOVED***

	process.handle = 0

	logrus.Debugf(title+" succeeded processid=%d", process.processID)
	return nil
***REMOVED***

func (process *process) registerCallback() error ***REMOVED***
	context := &notifcationWatcherContext***REMOVED***
		channels: newChannels(),
	***REMOVED***

	callbackMapLock.Lock()
	callbackNumber := nextCallback
	nextCallback++
	callbackMap[callbackNumber] = context
	callbackMapLock.Unlock()

	var callbackHandle hcsCallback
	err := hcsRegisterProcessCallback(process.handle, notificationWatcherCallback, callbackNumber, &callbackHandle)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	context.handle = callbackHandle
	process.callbackNumber = callbackNumber

	return nil
***REMOVED***

func (process *process) unregisterCallback() error ***REMOVED***
	callbackNumber := process.callbackNumber

	callbackMapLock.RLock()
	context := callbackMap[callbackNumber]
	callbackMapLock.RUnlock()

	if context == nil ***REMOVED***
		return nil
	***REMOVED***

	handle := context.handle

	if handle == 0 ***REMOVED***
		return nil
	***REMOVED***

	// hcsUnregisterProcessCallback has its own syncronization
	// to wait for all callbacks to complete. We must NOT hold the callbackMapLock.
	err := hcsUnregisterProcessCallback(handle)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	closeChannels(context.channels)

	callbackMapLock.Lock()
	callbackMap[callbackNumber] = nil
	callbackMapLock.Unlock()

	handle = 0

	return nil
***REMOVED***
