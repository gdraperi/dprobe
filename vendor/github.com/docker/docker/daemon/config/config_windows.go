package config

import (
	"github.com/docker/docker/api/types"
)

// BridgeConfig stores all the bridge driver specific
// configuration.
type BridgeConfig struct ***REMOVED***
	commonBridgeConfig
***REMOVED***

// Config defines the configuration of a docker daemon.
// These are the configuration settings that you pass
// to the docker daemon when you launch it with say: `dockerd -e windows`
type Config struct ***REMOVED***
	CommonConfig

	// Fields below here are platform specific. (There are none presently
	// for the Windows daemon.)
***REMOVED***

// GetRuntime returns the runtime path and arguments for a given
// runtime name
func (conf *Config) GetRuntime(name string) *types.Runtime ***REMOVED***
	return nil
***REMOVED***

// GetInitPath returns the configure docker-init path
func (conf *Config) GetInitPath() string ***REMOVED***
	return ""
***REMOVED***

// GetDefaultRuntimeName returns the current default runtime
func (conf *Config) GetDefaultRuntimeName() string ***REMOVED***
	return StockRuntimeName
***REMOVED***

// GetAllRuntimes returns a copy of the runtimes map
func (conf *Config) GetAllRuntimes() map[string]types.Runtime ***REMOVED***
	return map[string]types.Runtime***REMOVED******REMOVED***
***REMOVED***

// GetExecRoot returns the user configured Exec-root
func (conf *Config) GetExecRoot() string ***REMOVED***
	return ""
***REMOVED***

// IsSwarmCompatible defines if swarm mode can be enabled in this config
func (conf *Config) IsSwarmCompatible() error ***REMOVED***
	return nil
***REMOVED***

// ValidatePlatformConfig checks if any platform-specific configuration settings are invalid.
func (conf *Config) ValidatePlatformConfig() error ***REMOVED***
	return nil
***REMOVED***
