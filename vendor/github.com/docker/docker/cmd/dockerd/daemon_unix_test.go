// +build !windows

package main

import (
	"testing"

	"github.com/docker/docker/daemon/config"
	"github.com/gotestyourself/gotestyourself/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadDaemonCliConfigWithDaemonFlags(t *testing.T) ***REMOVED***
	content := `***REMOVED***"log-opts": ***REMOVED***"max-size": "1k"***REMOVED******REMOVED***`
	tempFile := fs.NewFile(t, "config", fs.WithContent(content))
	defer tempFile.Remove()

	opts := defaultOptions(tempFile.Path())
	opts.Debug = true
	opts.LogLevel = "info"
	assert.NoError(t, opts.flags.Set("selinux-enabled", "true"))

	loadedConfig, err := loadDaemonCliConfig(opts)
	require.NoError(t, err)
	require.NotNil(t, loadedConfig)

	assert.True(t, loadedConfig.Debug)
	assert.Equal(t, "info", loadedConfig.LogLevel)
	assert.True(t, loadedConfig.EnableSelinuxSupport)
	assert.Equal(t, "json-file", loadedConfig.LogConfig.Type)
	assert.Equal(t, "1k", loadedConfig.LogConfig.Config["max-size"])
***REMOVED***

func TestLoadDaemonConfigWithNetwork(t *testing.T) ***REMOVED***
	content := `***REMOVED***"bip": "127.0.0.2", "ip": "127.0.0.1"***REMOVED***`
	tempFile := fs.NewFile(t, "config", fs.WithContent(content))
	defer tempFile.Remove()

	opts := defaultOptions(tempFile.Path())
	loadedConfig, err := loadDaemonCliConfig(opts)
	require.NoError(t, err)
	require.NotNil(t, loadedConfig)

	assert.Equal(t, "127.0.0.2", loadedConfig.IP)
	assert.Equal(t, "127.0.0.1", loadedConfig.DefaultIP.String())
***REMOVED***

func TestLoadDaemonConfigWithMapOptions(t *testing.T) ***REMOVED***
	content := `***REMOVED***
		"cluster-store-opts": ***REMOVED***"kv.cacertfile": "/var/lib/docker/discovery_certs/ca.pem"***REMOVED***,
		"log-opts": ***REMOVED***"tag": "test"***REMOVED***
***REMOVED***`
	tempFile := fs.NewFile(t, "config", fs.WithContent(content))
	defer tempFile.Remove()

	opts := defaultOptions(tempFile.Path())
	loadedConfig, err := loadDaemonCliConfig(opts)
	require.NoError(t, err)
	require.NotNil(t, loadedConfig)
	assert.NotNil(t, loadedConfig.ClusterOpts)

	expectedPath := "/var/lib/docker/discovery_certs/ca.pem"
	assert.Equal(t, expectedPath, loadedConfig.ClusterOpts["kv.cacertfile"])
	assert.NotNil(t, loadedConfig.LogConfig.Config)
	assert.Equal(t, "test", loadedConfig.LogConfig.Config["tag"])
***REMOVED***

func TestLoadDaemonConfigWithTrueDefaultValues(t *testing.T) ***REMOVED***
	content := `***REMOVED*** "userland-proxy": false ***REMOVED***`
	tempFile := fs.NewFile(t, "config", fs.WithContent(content))
	defer tempFile.Remove()

	opts := defaultOptions(tempFile.Path())
	loadedConfig, err := loadDaemonCliConfig(opts)
	require.NoError(t, err)
	require.NotNil(t, loadedConfig)

	assert.False(t, loadedConfig.EnableUserlandProxy)

	// make sure reloading doesn't generate configuration
	// conflicts after normalizing boolean values.
	reload := func(reloadedConfig *config.Config) ***REMOVED***
		assert.False(t, reloadedConfig.EnableUserlandProxy)
	***REMOVED***
	assert.NoError(t, config.Reload(opts.configFile, opts.flags, reload))
***REMOVED***

func TestLoadDaemonConfigWithTrueDefaultValuesLeaveDefaults(t *testing.T) ***REMOVED***
	tempFile := fs.NewFile(t, "config", fs.WithContent(`***REMOVED******REMOVED***`))
	defer tempFile.Remove()

	opts := defaultOptions(tempFile.Path())
	loadedConfig, err := loadDaemonCliConfig(opts)
	require.NoError(t, err)
	require.NotNil(t, loadedConfig)

	assert.True(t, loadedConfig.EnableUserlandProxy)
***REMOVED***
