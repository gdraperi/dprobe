package hcsshim

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	defaultTimeout = time.Minute * 4
)

const (
	pendingUpdatesQuery    = `***REMOVED*** "PropertyTypes" : ["PendingUpdates"]***REMOVED***`
	statisticsQuery        = `***REMOVED*** "PropertyTypes" : ["Statistics"]***REMOVED***`
	processListQuery       = `***REMOVED*** "PropertyTypes" : ["ProcessList"]***REMOVED***`
	mappedVirtualDiskQuery = `***REMOVED*** "PropertyTypes" : ["MappedVirtualDisk"]***REMOVED***`
)

type container struct ***REMOVED***
	handleLock     sync.RWMutex
	handle         hcsSystem
	id             string
	callbackNumber uintptr
***REMOVED***

// ContainerProperties holds the properties for a container and the processes running in that container
type ContainerProperties struct ***REMOVED***
	ID                           string `json:"Id"`
	Name                         string
	SystemType                   string
	Owner                        string
	SiloGUID                     string                              `json:"SiloGuid,omitempty"`
	RuntimeID                    string                              `json:"RuntimeId,omitempty"`
	IsRuntimeTemplate            bool                                `json:",omitempty"`
	RuntimeImagePath             string                              `json:",omitempty"`
	Stopped                      bool                                `json:",omitempty"`
	ExitType                     string                              `json:",omitempty"`
	AreUpdatesPending            bool                                `json:",omitempty"`
	ObRoot                       string                              `json:",omitempty"`
	Statistics                   Statistics                          `json:",omitempty"`
	ProcessList                  []ProcessListItem                   `json:",omitempty"`
	MappedVirtualDiskControllers map[int]MappedVirtualDiskController `json:",omitempty"`
***REMOVED***

// MemoryStats holds the memory statistics for a container
type MemoryStats struct ***REMOVED***
	UsageCommitBytes            uint64 `json:"MemoryUsageCommitBytes,omitempty"`
	UsageCommitPeakBytes        uint64 `json:"MemoryUsageCommitPeakBytes,omitempty"`
	UsagePrivateWorkingSetBytes uint64 `json:"MemoryUsagePrivateWorkingSetBytes,omitempty"`
***REMOVED***

// ProcessorStats holds the processor statistics for a container
type ProcessorStats struct ***REMOVED***
	TotalRuntime100ns  uint64 `json:",omitempty"`
	RuntimeUser100ns   uint64 `json:",omitempty"`
	RuntimeKernel100ns uint64 `json:",omitempty"`
***REMOVED***

// StorageStats holds the storage statistics for a container
type StorageStats struct ***REMOVED***
	ReadCountNormalized  uint64 `json:",omitempty"`
	ReadSizeBytes        uint64 `json:",omitempty"`
	WriteCountNormalized uint64 `json:",omitempty"`
	WriteSizeBytes       uint64 `json:",omitempty"`
***REMOVED***

// NetworkStats holds the network statistics for a container
type NetworkStats struct ***REMOVED***
	BytesReceived          uint64 `json:",omitempty"`
	BytesSent              uint64 `json:",omitempty"`
	PacketsReceived        uint64 `json:",omitempty"`
	PacketsSent            uint64 `json:",omitempty"`
	DroppedPacketsIncoming uint64 `json:",omitempty"`
	DroppedPacketsOutgoing uint64 `json:",omitempty"`
	EndpointId             string `json:",omitempty"`
	InstanceId             string `json:",omitempty"`
***REMOVED***

// Statistics is the structure returned by a statistics call on a container
type Statistics struct ***REMOVED***
	Timestamp          time.Time      `json:",omitempty"`
	ContainerStartTime time.Time      `json:",omitempty"`
	Uptime100ns        uint64         `json:",omitempty"`
	Memory             MemoryStats    `json:",omitempty"`
	Processor          ProcessorStats `json:",omitempty"`
	Storage            StorageStats   `json:",omitempty"`
	Network            []NetworkStats `json:",omitempty"`
***REMOVED***

// ProcessList is the structure of an item returned by a ProcessList call on a container
type ProcessListItem struct ***REMOVED***
	CreateTimestamp              time.Time `json:",omitempty"`
	ImageName                    string    `json:",omitempty"`
	KernelTime100ns              uint64    `json:",omitempty"`
	MemoryCommitBytes            uint64    `json:",omitempty"`
	MemoryWorkingSetPrivateBytes uint64    `json:",omitempty"`
	MemoryWorkingSetSharedBytes  uint64    `json:",omitempty"`
	ProcessId                    uint32    `json:",omitempty"`
	UserTime100ns                uint64    `json:",omitempty"`
***REMOVED***

// MappedVirtualDiskController is the structure of an item returned by a MappedVirtualDiskList call on a container
type MappedVirtualDiskController struct ***REMOVED***
	MappedVirtualDisks map[int]MappedVirtualDisk `json:",omitempty"`
***REMOVED***

// Type of Request Support in ModifySystem
type RequestType string

// Type of Resource Support in ModifySystem
type ResourceType string

// RequestType const
const (
	Add     RequestType  = "Add"
	Remove  RequestType  = "Remove"
	Network ResourceType = "Network"
)

// ResourceModificationRequestResponse is the structure used to send request to the container to modify the system
// Supported resource types are Network and Request Types are Add/Remove
type ResourceModificationRequestResponse struct ***REMOVED***
	Resource ResourceType `json:"ResourceType"`
	Data     interface***REMOVED******REMOVED***  `json:"Settings"`
	Request  RequestType  `json:"RequestType,omitempty"`
***REMOVED***

// createContainerAdditionalJSON is read from the environment at initialisation
// time. It allows an environment variable to define additional JSON which
// is merged in the CreateContainer call to HCS.
var createContainerAdditionalJSON string

func init() ***REMOVED***
	createContainerAdditionalJSON = os.Getenv("HCSSHIM_CREATECONTAINER_ADDITIONALJSON")
***REMOVED***

// CreateContainer creates a new container with the given configuration but does not start it.
func CreateContainer(id string, c *ContainerConfig) (Container, error) ***REMOVED***
	return createContainerWithJSON(id, c, "")
***REMOVED***

// CreateContainerWithJSON creates a new container with the given configuration but does not start it.
// It is identical to CreateContainer except that optional additional JSON can be merged before passing to HCS.
func CreateContainerWithJSON(id string, c *ContainerConfig, additionalJSON string) (Container, error) ***REMOVED***
	return createContainerWithJSON(id, c, additionalJSON)
***REMOVED***

func createContainerWithJSON(id string, c *ContainerConfig, additionalJSON string) (Container, error) ***REMOVED***
	operation := "CreateContainer"
	title := "HCSShim::" + operation

	container := &container***REMOVED***
		id: id,
	***REMOVED***

	configurationb, err := json.Marshal(c)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	configuration := string(configurationb)
	logrus.Debugf(title+" id=%s config=%s", id, configuration)

	// Merge any additional JSON. Priority is given to what is passed in explicitly,
	// falling back to what's set in the environment.
	if additionalJSON == "" && createContainerAdditionalJSON != "" ***REMOVED***
		additionalJSON = createContainerAdditionalJSON
	***REMOVED***
	if additionalJSON != "" ***REMOVED***
		configurationMap := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
		if err := json.Unmarshal([]byte(configuration), &configurationMap); err != nil ***REMOVED***
			return nil, fmt.Errorf("failed to unmarshal %s: %s", configuration, err)
		***REMOVED***

		additionalMap := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
		if err := json.Unmarshal([]byte(additionalJSON), &additionalMap); err != nil ***REMOVED***
			return nil, fmt.Errorf("failed to unmarshal %s: %s", additionalJSON, err)
		***REMOVED***

		mergedMap := mergeMaps(additionalMap, configurationMap)
		mergedJSON, err := json.Marshal(mergedMap)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("failed to marshal merged configuration map %+v: %s", mergedMap, err)
		***REMOVED***

		configuration = string(mergedJSON)
		logrus.Debugf(title+" id=%s merged config=%s", id, configuration)
	***REMOVED***

	var (
		resultp  *uint16
		identity syscall.Handle
	)
	createError := hcsCreateComputeSystem(id, configuration, identity, &container.handle, &resultp)

	if createError == nil || IsPending(createError) ***REMOVED***
		if err := container.registerCallback(); err != nil ***REMOVED***
			// Terminate the container if it still exists. We're okay to ignore a failure here.
			container.Terminate()
			return nil, makeContainerError(container, operation, "", err)
		***REMOVED***
	***REMOVED***

	err = processAsyncHcsResult(createError, resultp, container.callbackNumber, hcsNotificationSystemCreateCompleted, &defaultTimeout)
	if err != nil ***REMOVED***
		if err == ErrTimeout ***REMOVED***
			// Terminate the container if it still exists. We're okay to ignore a failure here.
			container.Terminate()
		***REMOVED***
		return nil, makeContainerError(container, operation, configuration, err)
	***REMOVED***

	logrus.Debugf(title+" succeeded id=%s handle=%d", id, container.handle)
	return container, nil
***REMOVED***

// mergeMaps recursively merges map `fromMap` into map `ToMap`. Any pre-existing values
// in ToMap are overwritten. Values in fromMap are added to ToMap.
// From http://stackoverflow.com/questions/40491438/merging-two-json-strings-in-golang
func mergeMaps(fromMap, ToMap interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	switch fromMap := fromMap.(type) ***REMOVED***
	case map[string]interface***REMOVED******REMOVED***:
		ToMap, ok := ToMap.(map[string]interface***REMOVED******REMOVED***)
		if !ok ***REMOVED***
			return fromMap
		***REMOVED***
		for keyToMap, valueToMap := range ToMap ***REMOVED***
			if valueFromMap, ok := fromMap[keyToMap]; ok ***REMOVED***
				fromMap[keyToMap] = mergeMaps(valueFromMap, valueToMap)
			***REMOVED*** else ***REMOVED***
				fromMap[keyToMap] = valueToMap
			***REMOVED***
		***REMOVED***
	case nil:
		// merge(nil, map[string]interface***REMOVED***...***REMOVED***) -> map[string]interface***REMOVED***...***REMOVED***
		ToMap, ok := ToMap.(map[string]interface***REMOVED******REMOVED***)
		if ok ***REMOVED***
			return ToMap
		***REMOVED***
	***REMOVED***
	return fromMap
***REMOVED***

// OpenContainer opens an existing container by ID.
func OpenContainer(id string) (Container, error) ***REMOVED***
	operation := "OpenContainer"
	title := "HCSShim::" + operation
	logrus.Debugf(title+" id=%s", id)

	container := &container***REMOVED***
		id: id,
	***REMOVED***

	var (
		handle  hcsSystem
		resultp *uint16
	)
	err := hcsOpenComputeSystem(id, &handle, &resultp)
	err = processHcsResult(err, resultp)
	if err != nil ***REMOVED***
		return nil, makeContainerError(container, operation, "", err)
	***REMOVED***

	container.handle = handle

	if err := container.registerCallback(); err != nil ***REMOVED***
		return nil, makeContainerError(container, operation, "", err)
	***REMOVED***

	logrus.Debugf(title+" succeeded id=%s handle=%d", id, handle)
	return container, nil
***REMOVED***

// GetContainers gets a list of the containers on the system that match the query
func GetContainers(q ComputeSystemQuery) ([]ContainerProperties, error) ***REMOVED***
	operation := "GetContainers"
	title := "HCSShim::" + operation

	queryb, err := json.Marshal(q)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	query := string(queryb)
	logrus.Debugf(title+" query=%s", query)

	var (
		resultp         *uint16
		computeSystemsp *uint16
	)
	err = hcsEnumerateComputeSystems(query, &computeSystemsp, &resultp)
	err = processHcsResult(err, resultp)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if computeSystemsp == nil ***REMOVED***
		return nil, ErrUnexpectedValue
	***REMOVED***
	computeSystemsRaw := convertAndFreeCoTaskMemBytes(computeSystemsp)
	computeSystems := []ContainerProperties***REMOVED******REMOVED***
	if err := json.Unmarshal(computeSystemsRaw, &computeSystems); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	logrus.Debugf(title + " succeeded")
	return computeSystems, nil
***REMOVED***

// Start synchronously starts the container.
func (container *container) Start() error ***REMOVED***
	container.handleLock.RLock()
	defer container.handleLock.RUnlock()
	operation := "Start"
	title := "HCSShim::Container::" + operation
	logrus.Debugf(title+" id=%s", container.id)

	if container.handle == 0 ***REMOVED***
		return makeContainerError(container, operation, "", ErrAlreadyClosed)
	***REMOVED***

	var resultp *uint16
	err := hcsStartComputeSystem(container.handle, "", &resultp)
	err = processAsyncHcsResult(err, resultp, container.callbackNumber, hcsNotificationSystemStartCompleted, &defaultTimeout)
	if err != nil ***REMOVED***
		return makeContainerError(container, operation, "", err)
	***REMOVED***

	logrus.Debugf(title+" succeeded id=%s", container.id)
	return nil
***REMOVED***

// Shutdown requests a container shutdown, if IsPending() on the error returned is true,
// it may not actually be shut down until Wait() succeeds.
func (container *container) Shutdown() error ***REMOVED***
	container.handleLock.RLock()
	defer container.handleLock.RUnlock()
	operation := "Shutdown"
	title := "HCSShim::Container::" + operation
	logrus.Debugf(title+" id=%s", container.id)

	if container.handle == 0 ***REMOVED***
		return makeContainerError(container, operation, "", ErrAlreadyClosed)
	***REMOVED***

	var resultp *uint16
	err := hcsShutdownComputeSystem(container.handle, "", &resultp)
	err = processHcsResult(err, resultp)
	if err != nil ***REMOVED***
		return makeContainerError(container, operation, "", err)
	***REMOVED***

	logrus.Debugf(title+" succeeded id=%s", container.id)
	return nil
***REMOVED***

// Terminate requests a container terminate, if IsPending() on the error returned is true,
// it may not actually be shut down until Wait() succeeds.
func (container *container) Terminate() error ***REMOVED***
	container.handleLock.RLock()
	defer container.handleLock.RUnlock()
	operation := "Terminate"
	title := "HCSShim::Container::" + operation
	logrus.Debugf(title+" id=%s", container.id)

	if container.handle == 0 ***REMOVED***
		return makeContainerError(container, operation, "", ErrAlreadyClosed)
	***REMOVED***

	var resultp *uint16
	err := hcsTerminateComputeSystem(container.handle, "", &resultp)
	err = processHcsResult(err, resultp)
	if err != nil ***REMOVED***
		return makeContainerError(container, operation, "", err)
	***REMOVED***

	logrus.Debugf(title+" succeeded id=%s", container.id)
	return nil
***REMOVED***

// Wait synchronously waits for the container to shutdown or terminate.
func (container *container) Wait() error ***REMOVED***
	operation := "Wait"
	title := "HCSShim::Container::" + operation
	logrus.Debugf(title+" id=%s", container.id)

	err := waitForNotification(container.callbackNumber, hcsNotificationSystemExited, nil)
	if err != nil ***REMOVED***
		return makeContainerError(container, operation, "", err)
	***REMOVED***

	logrus.Debugf(title+" succeeded id=%s", container.id)
	return nil
***REMOVED***

// WaitTimeout synchronously waits for the container to terminate or the duration to elapse.
// If the timeout expires, IsTimeout(err) == true
func (container *container) WaitTimeout(timeout time.Duration) error ***REMOVED***
	operation := "WaitTimeout"
	title := "HCSShim::Container::" + operation
	logrus.Debugf(title+" id=%s", container.id)

	err := waitForNotification(container.callbackNumber, hcsNotificationSystemExited, &timeout)
	if err != nil ***REMOVED***
		return makeContainerError(container, operation, "", err)
	***REMOVED***

	logrus.Debugf(title+" succeeded id=%s", container.id)
	return nil
***REMOVED***

func (container *container) properties(query string) (*ContainerProperties, error) ***REMOVED***
	var (
		resultp     *uint16
		propertiesp *uint16
	)
	err := hcsGetComputeSystemProperties(container.handle, query, &propertiesp, &resultp)
	err = processHcsResult(err, resultp)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if propertiesp == nil ***REMOVED***
		return nil, ErrUnexpectedValue
	***REMOVED***
	propertiesRaw := convertAndFreeCoTaskMemBytes(propertiesp)
	properties := &ContainerProperties***REMOVED******REMOVED***
	if err := json.Unmarshal(propertiesRaw, properties); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return properties, nil
***REMOVED***

// HasPendingUpdates returns true if the container has updates pending to install
func (container *container) HasPendingUpdates() (bool, error) ***REMOVED***
	container.handleLock.RLock()
	defer container.handleLock.RUnlock()
	operation := "HasPendingUpdates"
	title := "HCSShim::Container::" + operation
	logrus.Debugf(title+" id=%s", container.id)

	if container.handle == 0 ***REMOVED***
		return false, makeContainerError(container, operation, "", ErrAlreadyClosed)
	***REMOVED***

	properties, err := container.properties(pendingUpdatesQuery)
	if err != nil ***REMOVED***
		return false, makeContainerError(container, operation, "", err)
	***REMOVED***

	logrus.Debugf(title+" succeeded id=%s", container.id)
	return properties.AreUpdatesPending, nil
***REMOVED***

// Statistics returns statistics for the container
func (container *container) Statistics() (Statistics, error) ***REMOVED***
	container.handleLock.RLock()
	defer container.handleLock.RUnlock()
	operation := "Statistics"
	title := "HCSShim::Container::" + operation
	logrus.Debugf(title+" id=%s", container.id)

	if container.handle == 0 ***REMOVED***
		return Statistics***REMOVED******REMOVED***, makeContainerError(container, operation, "", ErrAlreadyClosed)
	***REMOVED***

	properties, err := container.properties(statisticsQuery)
	if err != nil ***REMOVED***
		return Statistics***REMOVED******REMOVED***, makeContainerError(container, operation, "", err)
	***REMOVED***

	logrus.Debugf(title+" succeeded id=%s", container.id)
	return properties.Statistics, nil
***REMOVED***

// ProcessList returns an array of ProcessListItems for the container
func (container *container) ProcessList() ([]ProcessListItem, error) ***REMOVED***
	container.handleLock.RLock()
	defer container.handleLock.RUnlock()
	operation := "ProcessList"
	title := "HCSShim::Container::" + operation
	logrus.Debugf(title+" id=%s", container.id)

	if container.handle == 0 ***REMOVED***
		return nil, makeContainerError(container, operation, "", ErrAlreadyClosed)
	***REMOVED***

	properties, err := container.properties(processListQuery)
	if err != nil ***REMOVED***
		return nil, makeContainerError(container, operation, "", err)
	***REMOVED***

	logrus.Debugf(title+" succeeded id=%s", container.id)
	return properties.ProcessList, nil
***REMOVED***

// MappedVirtualDisks returns a map of the controllers and the disks mapped
// to a container.
//
// Example of JSON returned by the query.
//***REMOVED***
//   "Id":"1126e8d7d279c707a666972a15976371d365eaf622c02cea2c442b84f6f550a3_svm",
//   "SystemType":"Container",
//   "RuntimeOsType":"Linux",
//   "RuntimeId":"00000000-0000-0000-0000-000000000000",
//   "State":"Running",
//   "MappedVirtualDiskControllers":***REMOVED***
//      "0":***REMOVED***
//         "MappedVirtualDisks":***REMOVED***
//            "2":***REMOVED***
//               "HostPath":"C:\\lcow\\lcow\\scratch\\1126e8d7d279c707a666972a15976371d365eaf622c02cea2c442b84f6f550a3.vhdx",
//               "ContainerPath":"/mnt/gcs/LinuxServiceVM/scratch",
//               "Lun":2,
//               "CreateInUtilityVM":true
//        ***REMOVED***,
//            "3":***REMOVED***
//               "HostPath":"C:\\lcow\\lcow\\1126e8d7d279c707a666972a15976371d365eaf622c02cea2c442b84f6f550a3\\sandbox.vhdx",
//               "Lun":3,
//               "CreateInUtilityVM":true,
//               "AttachOnly":true
//        ***REMOVED***
//     ***REMOVED***
//  ***REMOVED***
//   ***REMOVED***
//***REMOVED***
func (container *container) MappedVirtualDisks() (map[int]MappedVirtualDiskController, error) ***REMOVED***
	container.handleLock.RLock()
	defer container.handleLock.RUnlock()
	operation := "MappedVirtualDiskList"
	title := "HCSShim::Container::" + operation
	logrus.Debugf(title+" id=%s", container.id)

	if container.handle == 0 ***REMOVED***
		return nil, makeContainerError(container, operation, "", ErrAlreadyClosed)
	***REMOVED***

	properties, err := container.properties(mappedVirtualDiskQuery)
	if err != nil ***REMOVED***
		return nil, makeContainerError(container, operation, "", err)
	***REMOVED***

	logrus.Debugf(title+" succeeded id=%s", container.id)
	return properties.MappedVirtualDiskControllers, nil
***REMOVED***

// Pause pauses the execution of the container. This feature is not enabled in TP5.
func (container *container) Pause() error ***REMOVED***
	container.handleLock.RLock()
	defer container.handleLock.RUnlock()
	operation := "Pause"
	title := "HCSShim::Container::" + operation
	logrus.Debugf(title+" id=%s", container.id)

	if container.handle == 0 ***REMOVED***
		return makeContainerError(container, operation, "", ErrAlreadyClosed)
	***REMOVED***

	var resultp *uint16
	err := hcsPauseComputeSystem(container.handle, "", &resultp)
	err = processAsyncHcsResult(err, resultp, container.callbackNumber, hcsNotificationSystemPauseCompleted, &defaultTimeout)
	if err != nil ***REMOVED***
		return makeContainerError(container, operation, "", err)
	***REMOVED***

	logrus.Debugf(title+" succeeded id=%s", container.id)
	return nil
***REMOVED***

// Resume resumes the execution of the container. This feature is not enabled in TP5.
func (container *container) Resume() error ***REMOVED***
	container.handleLock.RLock()
	defer container.handleLock.RUnlock()
	operation := "Resume"
	title := "HCSShim::Container::" + operation
	logrus.Debugf(title+" id=%s", container.id)

	if container.handle == 0 ***REMOVED***
		return makeContainerError(container, operation, "", ErrAlreadyClosed)
	***REMOVED***

	var resultp *uint16
	err := hcsResumeComputeSystem(container.handle, "", &resultp)
	err = processAsyncHcsResult(err, resultp, container.callbackNumber, hcsNotificationSystemResumeCompleted, &defaultTimeout)
	if err != nil ***REMOVED***
		return makeContainerError(container, operation, "", err)
	***REMOVED***

	logrus.Debugf(title+" succeeded id=%s", container.id)
	return nil
***REMOVED***

// CreateProcess launches a new process within the container.
func (container *container) CreateProcess(c *ProcessConfig) (Process, error) ***REMOVED***
	container.handleLock.RLock()
	defer container.handleLock.RUnlock()
	operation := "CreateProcess"
	title := "HCSShim::Container::" + operation
	var (
		processInfo   hcsProcessInformation
		processHandle hcsProcess
		resultp       *uint16
	)

	if container.handle == 0 ***REMOVED***
		return nil, makeContainerError(container, operation, "", ErrAlreadyClosed)
	***REMOVED***

	// If we are not emulating a console, ignore any console size passed to us
	if !c.EmulateConsole ***REMOVED***
		c.ConsoleSize[0] = 0
		c.ConsoleSize[1] = 0
	***REMOVED***

	configurationb, err := json.Marshal(c)
	if err != nil ***REMOVED***
		return nil, makeContainerError(container, operation, "", err)
	***REMOVED***

	configuration := string(configurationb)
	logrus.Debugf(title+" id=%s config=%s", container.id, configuration)

	err = hcsCreateProcess(container.handle, configuration, &processInfo, &processHandle, &resultp)
	err = processHcsResult(err, resultp)
	if err != nil ***REMOVED***
		return nil, makeContainerError(container, operation, configuration, err)
	***REMOVED***

	process := &process***REMOVED***
		handle:    processHandle,
		processID: int(processInfo.ProcessId),
		container: container,
		cachedPipes: &cachedPipes***REMOVED***
			stdIn:  processInfo.StdInput,
			stdOut: processInfo.StdOutput,
			stdErr: processInfo.StdError,
		***REMOVED***,
	***REMOVED***

	if err := process.registerCallback(); err != nil ***REMOVED***
		return nil, makeContainerError(container, operation, "", err)
	***REMOVED***

	logrus.Debugf(title+" succeeded id=%s processid=%d", container.id, process.processID)
	return process, nil
***REMOVED***

// OpenProcess gets an interface to an existing process within the container.
func (container *container) OpenProcess(pid int) (Process, error) ***REMOVED***
	container.handleLock.RLock()
	defer container.handleLock.RUnlock()
	operation := "OpenProcess"
	title := "HCSShim::Container::" + operation
	logrus.Debugf(title+" id=%s, processid=%d", container.id, pid)
	var (
		processHandle hcsProcess
		resultp       *uint16
	)

	if container.handle == 0 ***REMOVED***
		return nil, makeContainerError(container, operation, "", ErrAlreadyClosed)
	***REMOVED***

	err := hcsOpenProcess(container.handle, uint32(pid), &processHandle, &resultp)
	err = processHcsResult(err, resultp)
	if err != nil ***REMOVED***
		return nil, makeContainerError(container, operation, "", err)
	***REMOVED***

	process := &process***REMOVED***
		handle:    processHandle,
		processID: pid,
		container: container,
	***REMOVED***

	if err := process.registerCallback(); err != nil ***REMOVED***
		return nil, makeContainerError(container, operation, "", err)
	***REMOVED***

	logrus.Debugf(title+" succeeded id=%s processid=%s", container.id, process.processID)
	return process, nil
***REMOVED***

// Close cleans up any state associated with the container but does not terminate or wait for it.
func (container *container) Close() error ***REMOVED***
	container.handleLock.Lock()
	defer container.handleLock.Unlock()
	operation := "Close"
	title := "HCSShim::Container::" + operation
	logrus.Debugf(title+" id=%s", container.id)

	// Don't double free this
	if container.handle == 0 ***REMOVED***
		return nil
	***REMOVED***

	if err := container.unregisterCallback(); err != nil ***REMOVED***
		return makeContainerError(container, operation, "", err)
	***REMOVED***

	if err := hcsCloseComputeSystem(container.handle); err != nil ***REMOVED***
		return makeContainerError(container, operation, "", err)
	***REMOVED***

	container.handle = 0

	logrus.Debugf(title+" succeeded id=%s", container.id)
	return nil
***REMOVED***

func (container *container) registerCallback() error ***REMOVED***
	context := &notifcationWatcherContext***REMOVED***
		channels: newChannels(),
	***REMOVED***

	callbackMapLock.Lock()
	callbackNumber := nextCallback
	nextCallback++
	callbackMap[callbackNumber] = context
	callbackMapLock.Unlock()

	var callbackHandle hcsCallback
	err := hcsRegisterComputeSystemCallback(container.handle, notificationWatcherCallback, callbackNumber, &callbackHandle)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	context.handle = callbackHandle
	container.callbackNumber = callbackNumber

	return nil
***REMOVED***

func (container *container) unregisterCallback() error ***REMOVED***
	callbackNumber := container.callbackNumber

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

	// hcsUnregisterComputeSystemCallback has its own syncronization
	// to wait for all callbacks to complete. We must NOT hold the callbackMapLock.
	err := hcsUnregisterComputeSystemCallback(handle)
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

// Modifies the System by sending a request to HCS
func (container *container) Modify(config *ResourceModificationRequestResponse) error ***REMOVED***
	container.handleLock.RLock()
	defer container.handleLock.RUnlock()
	operation := "Modify"
	title := "HCSShim::Container::" + operation

	if container.handle == 0 ***REMOVED***
		return makeContainerError(container, operation, "", ErrAlreadyClosed)
	***REMOVED***

	requestJSON, err := json.Marshal(config)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	requestString := string(requestJSON)
	logrus.Debugf(title+" id=%s request=%s", container.id, requestString)

	var resultp *uint16
	err = hcsModifyComputeSystem(container.handle, requestString, &resultp)
	err = processHcsResult(err, resultp)
	if err != nil ***REMOVED***
		return makeContainerError(container, operation, "", err)
	***REMOVED***
	logrus.Debugf(title+" succeeded id=%s", container.id)
	return nil
***REMOVED***
