package environment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/integration-cli/fixtures/load"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// Execution contains information about the current test execution and daemon
// under test
type Execution struct ***REMOVED***
	client            client.APIClient
	DaemonInfo        types.Info
	OSType            string
	PlatformDefaults  PlatformDefaults
	protectedElements protectedElements
***REMOVED***

// PlatformDefaults are defaults values for the platform of the daemon under test
type PlatformDefaults struct ***REMOVED***
	BaseImage            string
	VolumesConfigPath    string
	ContainerStoragePath string
***REMOVED***

// New creates a new Execution struct
func New() (*Execution, error) ***REMOVED***
	client, err := client.NewEnvClient()
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "failed to create client")
	***REMOVED***

	info, err := client.Info(context.Background())
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "failed to get info from daemon")
	***REMOVED***

	osType := getOSType(info)

	return &Execution***REMOVED***
		client:            client,
		DaemonInfo:        info,
		OSType:            osType,
		PlatformDefaults:  getPlatformDefaults(info, osType),
		protectedElements: newProtectedElements(),
	***REMOVED***, nil
***REMOVED***

func getOSType(info types.Info) string ***REMOVED***
	// Docker EE does not set the OSType so allow the user to override this value.
	userOsType := os.Getenv("TEST_OSTYPE")
	if userOsType != "" ***REMOVED***
		return userOsType
	***REMOVED***
	return info.OSType
***REMOVED***

func getPlatformDefaults(info types.Info, osType string) PlatformDefaults ***REMOVED***
	volumesPath := filepath.Join(info.DockerRootDir, "volumes")
	containersPath := filepath.Join(info.DockerRootDir, "containers")

	switch osType ***REMOVED***
	case "linux":
		return PlatformDefaults***REMOVED***
			BaseImage:            "scratch",
			VolumesConfigPath:    toSlash(volumesPath),
			ContainerStoragePath: toSlash(containersPath),
		***REMOVED***
	case "windows":
		baseImage := "microsoft/windowsservercore"
		if override := os.Getenv("WINDOWS_BASE_IMAGE"); override != "" ***REMOVED***
			baseImage = override
			fmt.Println("INFO: Windows Base image is ", baseImage)
		***REMOVED***
		return PlatformDefaults***REMOVED***
			BaseImage:            baseImage,
			VolumesConfigPath:    filepath.FromSlash(volumesPath),
			ContainerStoragePath: filepath.FromSlash(containersPath),
		***REMOVED***
	default:
		panic(fmt.Sprintf("unknown OSType for daemon: %s", osType))
	***REMOVED***
***REMOVED***

// Make sure in context of daemon, not the local platform. Note we can't
// use filepath.FromSlash or ToSlash here as they are a no-op on Unix.
func toSlash(path string) string ***REMOVED***
	return strings.Replace(path, `\`, `/`, -1)
***REMOVED***

// IsLocalDaemon is true if the daemon under test is on the same
// host as the CLI.
//
// Deterministically working out the environment in which CI is running
// to evaluate whether the daemon is local or remote is not possible through
// a build tag.
//
// For example Windows to Linux CI under Jenkins tests the 64-bit
// Windows binary build with the daemon build tag, but calls a remote
// Linux daemon.
//
// We can't just say if Windows then assume the daemon is local as at
// some point, we will be testing the Windows CLI against a Windows daemon.
//
// Similarly, it will be perfectly valid to also run CLI tests from
// a Linux CLI (built with the daemon tag) against a Windows daemon.
func (e *Execution) IsLocalDaemon() bool ***REMOVED***
	return os.Getenv("DOCKER_REMOTE_DAEMON") == ""
***REMOVED***

// Print the execution details to stdout
// TODO: print everything
func (e *Execution) Print() ***REMOVED***
	if e.IsLocalDaemon() ***REMOVED***
		fmt.Println("INFO: Testing against a local daemon")
	***REMOVED*** else ***REMOVED***
		fmt.Println("INFO: Testing against a remote daemon")
	***REMOVED***
***REMOVED***

// APIClient returns an APIClient connected to the daemon under test
func (e *Execution) APIClient() client.APIClient ***REMOVED***
	return e.client
***REMOVED***

// EnsureFrozenImagesLinux loads frozen test images into the daemon
// if they aren't already loaded
func EnsureFrozenImagesLinux(testEnv *Execution) error ***REMOVED***
	if testEnv.OSType == "linux" ***REMOVED***
		err := load.FrozenImagesLinux(testEnv.APIClient(), frozenImages...)
		if err != nil ***REMOVED***
			return errors.Wrap(err, "error loading frozen images")
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
