// +build linux freebsd

package daemon

import (
	"bytes"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/daemon/config"
)

// reloadPlatform updates configuration with platform specific options
// and updates the passed attributes
func (daemon *Daemon) reloadPlatform(conf *config.Config, attributes map[string]string) error ***REMOVED***
	if err := conf.ValidatePlatformConfig(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if conf.IsValueSet("runtimes") ***REMOVED***
		// Always set the default one
		conf.Runtimes[config.StockRuntimeName] = types.Runtime***REMOVED***Path: DefaultRuntimeBinary***REMOVED***
		if err := daemon.initRuntimes(conf.Runtimes); err != nil ***REMOVED***
			return err
		***REMOVED***
		daemon.configStore.Runtimes = conf.Runtimes
	***REMOVED***

	if conf.DefaultRuntime != "" ***REMOVED***
		daemon.configStore.DefaultRuntime = conf.DefaultRuntime
	***REMOVED***

	if conf.IsValueSet("default-shm-size") ***REMOVED***
		daemon.configStore.ShmSize = conf.ShmSize
	***REMOVED***

	if conf.IpcMode != "" ***REMOVED***
		daemon.configStore.IpcMode = conf.IpcMode
	***REMOVED***

	// Update attributes
	var runtimeList bytes.Buffer
	for name, rt := range daemon.configStore.Runtimes ***REMOVED***
		if runtimeList.Len() > 0 ***REMOVED***
			runtimeList.WriteRune(' ')
		***REMOVED***
		runtimeList.WriteString(fmt.Sprintf("%s:%s", name, rt))
	***REMOVED***

	attributes["runtimes"] = runtimeList.String()
	attributes["default-runtime"] = daemon.configStore.DefaultRuntime
	attributes["default-shm-size"] = fmt.Sprintf("%d", daemon.configStore.ShmSize)
	attributes["default-ipc-mode"] = daemon.configStore.IpcMode

	return nil
***REMOVED***
