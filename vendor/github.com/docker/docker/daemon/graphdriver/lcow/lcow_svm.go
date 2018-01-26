// +build windows

package lcow

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/Microsoft/hcsshim"
	"github.com/Microsoft/opengcs/client"
	"github.com/sirupsen/logrus"
)

// Code for all the service VM management for the LCOW graphdriver

var errVMisTerminating = errors.New("service VM is shutting down")
var errVMUnknown = errors.New("service vm id is unknown")
var errVMStillHasReference = errors.New("Attemping to delete a VM that is still being used")

// serviceVMMap is the struct representing the id -> service VM mapping.
type serviceVMMap struct ***REMOVED***
	sync.Mutex
	svms map[string]*serviceVMMapItem
***REMOVED***

// serviceVMMapItem is our internal structure representing an item in our
// map of service VMs we are maintaining.
type serviceVMMapItem struct ***REMOVED***
	svm      *serviceVM // actual service vm object
	refCount int        // refcount for VM
***REMOVED***

type serviceVM struct ***REMOVED***
	sync.Mutex                     // Serialises operations being performed in this service VM.
	scratchAttached bool           // Has a scratch been attached?
	config          *client.Config // Represents the service VM item.

	// Indicates that the vm is started
	startStatus chan interface***REMOVED******REMOVED***
	startError  error

	// Indicates that the vm is stopped
	stopStatus chan interface***REMOVED******REMOVED***
	stopError  error

	attachedVHDs map[string]int // Map ref counting all the VHDS we've hot-added/hot-removed.
	unionMounts  map[string]int // Map ref counting all the union filesystems we mounted.
***REMOVED***

// add will add an id to the service vm map. There are three cases:
// 	- entry doesn't exist:
// 		- add id to map and return a new vm that the caller can manually configure+start
//	- entry does exist
//  	- return vm in map and increment ref count
//  - entry does exist but the ref count is 0
//		- return the svm and errVMisTerminating. Caller can call svm.getStopError() to wait for stop
func (svmMap *serviceVMMap) add(id string) (svm *serviceVM, alreadyExists bool, err error) ***REMOVED***
	svmMap.Lock()
	defer svmMap.Unlock()
	if svm, ok := svmMap.svms[id]; ok ***REMOVED***
		if svm.refCount == 0 ***REMOVED***
			return svm.svm, true, errVMisTerminating
		***REMOVED***
		svm.refCount++
		return svm.svm, true, nil
	***REMOVED***

	// Doesn't exist, so create an empty svm to put into map and return
	newSVM := &serviceVM***REMOVED***
		startStatus:  make(chan interface***REMOVED******REMOVED***),
		stopStatus:   make(chan interface***REMOVED******REMOVED***),
		attachedVHDs: make(map[string]int),
		unionMounts:  make(map[string]int),
		config:       &client.Config***REMOVED******REMOVED***,
	***REMOVED***
	svmMap.svms[id] = &serviceVMMapItem***REMOVED***
		svm:      newSVM,
		refCount: 1,
	***REMOVED***
	return newSVM, false, nil
***REMOVED***

// get will get the service vm from the map. There are three cases:
// 	- entry doesn't exist:
// 		- return errVMUnknown
//	- entry does exist
//  	- return vm with no error
//  - entry does exist but the ref count is 0
//		- return the svm and errVMisTerminating. Caller can call svm.getStopError() to wait for stop
func (svmMap *serviceVMMap) get(id string) (*serviceVM, error) ***REMOVED***
	svmMap.Lock()
	defer svmMap.Unlock()
	svm, ok := svmMap.svms[id]
	if !ok ***REMOVED***
		return nil, errVMUnknown
	***REMOVED***
	if svm.refCount == 0 ***REMOVED***
		return svm.svm, errVMisTerminating
	***REMOVED***
	return svm.svm, nil
***REMOVED***

// decrementRefCount decrements the ref count of the given ID from the map. There are four cases:
// 	- entry doesn't exist:
// 		- return errVMUnknown
//  - entry does exist but the ref count is 0
//		- return the svm and errVMisTerminating. Caller can call svm.getStopError() to wait for stop
//	- entry does exist but ref count is 1
//  	- return vm and set lastRef to true. The caller can then stop the vm, delete the id from this map
//      - and execute svm.signalStopFinished to signal the threads that the svm has been terminated.
//	- entry does exist and ref count > 1
//		- just reduce ref count and return svm
func (svmMap *serviceVMMap) decrementRefCount(id string) (_ *serviceVM, lastRef bool, _ error) ***REMOVED***
	svmMap.Lock()
	defer svmMap.Unlock()

	svm, ok := svmMap.svms[id]
	if !ok ***REMOVED***
		return nil, false, errVMUnknown
	***REMOVED***
	if svm.refCount == 0 ***REMOVED***
		return svm.svm, false, errVMisTerminating
	***REMOVED***
	svm.refCount--
	return svm.svm, svm.refCount == 0, nil
***REMOVED***

// setRefCountZero works the same way as decrementRefCount, but sets ref count to 0 instead of decrementing it.
func (svmMap *serviceVMMap) setRefCountZero(id string) (*serviceVM, error) ***REMOVED***
	svmMap.Lock()
	defer svmMap.Unlock()

	svm, ok := svmMap.svms[id]
	if !ok ***REMOVED***
		return nil, errVMUnknown
	***REMOVED***
	if svm.refCount == 0 ***REMOVED***
		return svm.svm, errVMisTerminating
	***REMOVED***
	svm.refCount = 0
	return svm.svm, nil
***REMOVED***

// deleteID deletes the given ID from the map. If the refcount is not 0 or the
// VM does not exist, then this function returns an error.
func (svmMap *serviceVMMap) deleteID(id string) error ***REMOVED***
	svmMap.Lock()
	defer svmMap.Unlock()
	svm, ok := svmMap.svms[id]
	if !ok ***REMOVED***
		return errVMUnknown
	***REMOVED***
	if svm.refCount != 0 ***REMOVED***
		return errVMStillHasReference
	***REMOVED***
	delete(svmMap.svms, id)
	return nil
***REMOVED***

func (svm *serviceVM) signalStartFinished(err error) ***REMOVED***
	svm.Lock()
	svm.startError = err
	svm.Unlock()
	close(svm.startStatus)
***REMOVED***

func (svm *serviceVM) getStartError() error ***REMOVED***
	<-svm.startStatus
	svm.Lock()
	defer svm.Unlock()
	return svm.startError
***REMOVED***

func (svm *serviceVM) signalStopFinished(err error) ***REMOVED***
	svm.Lock()
	svm.stopError = err
	svm.Unlock()
	close(svm.stopStatus)
***REMOVED***

func (svm *serviceVM) getStopError() error ***REMOVED***
	<-svm.stopStatus
	svm.Lock()
	defer svm.Unlock()
	return svm.stopError
***REMOVED***

// hotAddVHDs waits for the service vm to start and then attaches the vhds.
func (svm *serviceVM) hotAddVHDs(mvds ...hcsshim.MappedVirtualDisk) error ***REMOVED***
	if err := svm.getStartError(); err != nil ***REMOVED***
		return err
	***REMOVED***
	return svm.hotAddVHDsAtStart(mvds...)
***REMOVED***

// hotAddVHDsAtStart works the same way as hotAddVHDs but does not wait for the VM to start.
func (svm *serviceVM) hotAddVHDsAtStart(mvds ...hcsshim.MappedVirtualDisk) error ***REMOVED***
	svm.Lock()
	defer svm.Unlock()
	for i, mvd := range mvds ***REMOVED***
		if _, ok := svm.attachedVHDs[mvd.HostPath]; ok ***REMOVED***
			svm.attachedVHDs[mvd.HostPath]++
			continue
		***REMOVED***

		if err := svm.config.HotAddVhd(mvd.HostPath, mvd.ContainerPath, mvd.ReadOnly, !mvd.AttachOnly); err != nil ***REMOVED***
			svm.hotRemoveVHDsAtStart(mvds[:i]...)
			return err
		***REMOVED***
		svm.attachedVHDs[mvd.HostPath] = 1
	***REMOVED***
	return nil
***REMOVED***

// hotRemoveVHDs waits for the service vm to start and then removes the vhds.
func (svm *serviceVM) hotRemoveVHDs(mvds ...hcsshim.MappedVirtualDisk) error ***REMOVED***
	if err := svm.getStartError(); err != nil ***REMOVED***
		return err
	***REMOVED***
	return svm.hotRemoveVHDsAtStart(mvds...)
***REMOVED***

// hotRemoveVHDsAtStart works the same way as hotRemoveVHDs but does not wait for the VM to start.
func (svm *serviceVM) hotRemoveVHDsAtStart(mvds ...hcsshim.MappedVirtualDisk) error ***REMOVED***
	svm.Lock()
	defer svm.Unlock()
	var retErr error
	for _, mvd := range mvds ***REMOVED***
		if _, ok := svm.attachedVHDs[mvd.HostPath]; !ok ***REMOVED***
			// We continue instead of returning an error if we try to hot remove a non-existent VHD.
			// This is because one of the callers of the function is graphdriver.Put(). Since graphdriver.Get()
			// defers the VM start to the first operation, it's possible that nothing have been hot-added
			// when Put() is called. To avoid Put returning an error in that case, we simply continue if we
			// don't find the vhd attached.
			continue
		***REMOVED***

		if svm.attachedVHDs[mvd.HostPath] > 1 ***REMOVED***
			svm.attachedVHDs[mvd.HostPath]--
			continue
		***REMOVED***

		// last VHD, so remove from VM and map
		if err := svm.config.HotRemoveVhd(mvd.HostPath); err == nil ***REMOVED***
			delete(svm.attachedVHDs, mvd.HostPath)
		***REMOVED*** else ***REMOVED***
			// Take note of the error, but still continue to remove the other VHDs
			logrus.Warnf("Failed to hot remove %s: %s", mvd.HostPath, err)
			if retErr == nil ***REMOVED***
				retErr = err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return retErr
***REMOVED***

func (svm *serviceVM) createExt4VHDX(destFile string, sizeGB uint32, cacheFile string) error ***REMOVED***
	if err := svm.getStartError(); err != nil ***REMOVED***
		return err
	***REMOVED***

	svm.Lock()
	defer svm.Unlock()
	return svm.config.CreateExt4Vhdx(destFile, sizeGB, cacheFile)
***REMOVED***

func (svm *serviceVM) createUnionMount(mountName string, mvds ...hcsshim.MappedVirtualDisk) (err error) ***REMOVED***
	if len(mvds) == 0 ***REMOVED***
		return fmt.Errorf("createUnionMount: error must have at least 1 layer")
	***REMOVED***

	if err = svm.getStartError(); err != nil ***REMOVED***
		return err
	***REMOVED***

	svm.Lock()
	defer svm.Unlock()
	if _, ok := svm.unionMounts[mountName]; ok ***REMOVED***
		svm.unionMounts[mountName]++
		return nil
	***REMOVED***

	var lowerLayers []string
	if mvds[0].ReadOnly ***REMOVED***
		lowerLayers = append(lowerLayers, mvds[0].ContainerPath)
	***REMOVED***

	for i := 1; i < len(mvds); i++ ***REMOVED***
		lowerLayers = append(lowerLayers, mvds[i].ContainerPath)
	***REMOVED***

	logrus.Debugf("Doing the overlay mount with union directory=%s", mountName)
	if err = svm.runProcess(fmt.Sprintf("mkdir -p %s", mountName), nil, nil, nil); err != nil ***REMOVED***
		return err
	***REMOVED***

	var cmd string
	if len(mvds) == 1 ***REMOVED***
		// `FROM SCRATCH` case and the only layer. No overlay required.
		cmd = fmt.Sprintf("mount %s %s", mvds[0].ContainerPath, mountName)
	***REMOVED*** else if mvds[0].ReadOnly ***REMOVED***
		// Readonly overlay
		cmd = fmt.Sprintf("mount -t overlay overlay -olowerdir=%s %s",
			strings.Join(lowerLayers, ","),
			mountName)
	***REMOVED*** else ***REMOVED***
		upper := fmt.Sprintf("%s/upper", mvds[0].ContainerPath)
		work := fmt.Sprintf("%s/work", mvds[0].ContainerPath)

		if err = svm.runProcess(fmt.Sprintf("mkdir -p %s %s", upper, work), nil, nil, nil); err != nil ***REMOVED***
			return err
		***REMOVED***

		cmd = fmt.Sprintf("mount -t overlay overlay -olowerdir=%s,upperdir=%s,workdir=%s %s",
			strings.Join(lowerLayers, ":"),
			upper,
			work,
			mountName)
	***REMOVED***

	logrus.Debugf("createUnionMount: Executing mount=%s", cmd)
	if err = svm.runProcess(cmd, nil, nil, nil); err != nil ***REMOVED***
		return err
	***REMOVED***

	svm.unionMounts[mountName] = 1
	return nil
***REMOVED***

func (svm *serviceVM) deleteUnionMount(mountName string, disks ...hcsshim.MappedVirtualDisk) error ***REMOVED***
	if err := svm.getStartError(); err != nil ***REMOVED***
		return err
	***REMOVED***

	svm.Lock()
	defer svm.Unlock()
	if _, ok := svm.unionMounts[mountName]; !ok ***REMOVED***
		return nil
	***REMOVED***

	if svm.unionMounts[mountName] > 1 ***REMOVED***
		svm.unionMounts[mountName]--
		return nil
	***REMOVED***

	logrus.Debugf("Removing union mount %s", mountName)
	if err := svm.runProcess(fmt.Sprintf("umount %s", mountName), nil, nil, nil); err != nil ***REMOVED***
		return err
	***REMOVED***

	delete(svm.unionMounts, mountName)
	return nil
***REMOVED***

func (svm *serviceVM) runProcess(command string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error ***REMOVED***
	process, err := svm.config.RunProcess(command, stdin, stdout, stderr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer process.Close()

	process.WaitTimeout(time.Duration(int(time.Second) * svm.config.UvmTimeoutSeconds))
	exitCode, err := process.ExitCode()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if exitCode != 0 ***REMOVED***
		return fmt.Errorf("svm.runProcess: command %s failed with exit code %d", command, exitCode)
	***REMOVED***
	return nil
***REMOVED***
