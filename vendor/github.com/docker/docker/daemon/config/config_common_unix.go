// +build linux freebsd

package config

import (
	"net"

	"github.com/docker/docker/api/types"
)

// CommonUnixConfig defines configuration of a docker daemon that is
// common across Unix platforms.
type CommonUnixConfig struct ***REMOVED***
	Runtimes          map[string]types.Runtime `json:"runtimes,omitempty"`
	DefaultRuntime    string                   `json:"default-runtime,omitempty"`
	DefaultInitBinary string                   `json:"default-init,omitempty"`
***REMOVED***

type commonUnixBridgeConfig struct ***REMOVED***
	DefaultIP                   net.IP `json:"ip,omitempty"`
	IP                          string `json:"bip,omitempty"`
	DefaultGatewayIPv4          net.IP `json:"default-gateway,omitempty"`
	DefaultGatewayIPv6          net.IP `json:"default-gateway-v6,omitempty"`
	InterContainerCommunication bool   `json:"icc,omitempty"`
***REMOVED***

// GetRuntime returns the runtime path and arguments for a given
// runtime name
func (conf *Config) GetRuntime(name string) *types.Runtime ***REMOVED***
	conf.Lock()
	defer conf.Unlock()
	if rt, ok := conf.Runtimes[name]; ok ***REMOVED***
		return &rt
	***REMOVED***
	return nil
***REMOVED***

// GetDefaultRuntimeName returns the current default runtime
func (conf *Config) GetDefaultRuntimeName() string ***REMOVED***
	conf.Lock()
	rt := conf.DefaultRuntime
	conf.Unlock()

	return rt
***REMOVED***

// GetAllRuntimes returns a copy of the runtimes map
func (conf *Config) GetAllRuntimes() map[string]types.Runtime ***REMOVED***
	conf.Lock()
	rts := conf.Runtimes
	conf.Unlock()
	return rts
***REMOVED***

// GetExecRoot returns the user configured Exec-root
func (conf *Config) GetExecRoot() string ***REMOVED***
	return conf.ExecRoot
***REMOVED***

// GetInitPath returns the configured docker-init path
func (conf *Config) GetInitPath() string ***REMOVED***
	conf.Lock()
	defer conf.Unlock()
	if conf.InitPath != "" ***REMOVED***
		return conf.InitPath
	***REMOVED***
	if conf.DefaultInitBinary != "" ***REMOVED***
		return conf.DefaultInitBinary
	***REMOVED***
	return DefaultInitBinary
***REMOVED***
