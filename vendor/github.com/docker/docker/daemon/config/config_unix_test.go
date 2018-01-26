// +build !windows

package config

import (
	"testing"

	"github.com/docker/docker/opts"
	units "github.com/docker/go-units"
	"github.com/gotestyourself/gotestyourself/fs"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConflictFreeConfiguration(t *testing.T) ***REMOVED***
	configFileData := `
		***REMOVED***
			"debug": true,
			"default-ulimits": ***REMOVED***
				"nofile": ***REMOVED***
					"Name": "nofile",
					"Hard": 2048,
					"Soft": 1024
				***REMOVED***
			***REMOVED***,
			"log-opts": ***REMOVED***
				"tag": "test_tag"
			***REMOVED***
		***REMOVED***`

	file := fs.NewFile(t, "docker-config", fs.WithContent(configFileData))
	defer file.Remove()

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	var debug bool
	flags.BoolVarP(&debug, "debug", "D", false, "")
	flags.Var(opts.NewNamedUlimitOpt("default-ulimits", nil), "default-ulimit", "")
	flags.Var(opts.NewNamedMapOpts("log-opts", nil, nil), "log-opt", "")

	cc, err := getConflictFreeConfiguration(file.Path(), flags)
	require.NoError(t, err)

	assert.True(t, cc.Debug)

	expectedUlimits := map[string]*units.Ulimit***REMOVED***
		"nofile": ***REMOVED***
			Name: "nofile",
			Hard: 2048,
			Soft: 1024,
		***REMOVED***,
	***REMOVED***

	assert.Equal(t, expectedUlimits, cc.Ulimits)
***REMOVED***

func TestDaemonConfigurationMerge(t *testing.T) ***REMOVED***
	configFileData := `
		***REMOVED***
			"debug": true,
			"default-ulimits": ***REMOVED***
				"nofile": ***REMOVED***
					"Name": "nofile",
					"Hard": 2048,
					"Soft": 1024
				***REMOVED***
			***REMOVED***,
			"log-opts": ***REMOVED***
				"tag": "test_tag"
			***REMOVED***
		***REMOVED***`

	file := fs.NewFile(t, "docker-config", fs.WithContent(configFileData))
	defer file.Remove()

	c := &Config***REMOVED***
		CommonConfig: CommonConfig***REMOVED***
			AutoRestart: true,
			LogConfig: LogConfig***REMOVED***
				Type:   "syslog",
				Config: map[string]string***REMOVED***"tag": "test"***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)

	var debug bool
	flags.BoolVarP(&debug, "debug", "D", false, "")
	flags.Var(opts.NewNamedUlimitOpt("default-ulimits", nil), "default-ulimit", "")
	flags.Var(opts.NewNamedMapOpts("log-opts", nil, nil), "log-opt", "")

	cc, err := MergeDaemonConfigurations(c, flags, file.Path())
	require.NoError(t, err)

	assert.True(t, cc.Debug)
	assert.True(t, cc.AutoRestart)

	expectedLogConfig := LogConfig***REMOVED***
		Type:   "syslog",
		Config: map[string]string***REMOVED***"tag": "test_tag"***REMOVED***,
	***REMOVED***

	assert.Equal(t, expectedLogConfig, cc.LogConfig)

	expectedUlimits := map[string]*units.Ulimit***REMOVED***
		"nofile": ***REMOVED***
			Name: "nofile",
			Hard: 2048,
			Soft: 1024,
		***REMOVED***,
	***REMOVED***

	assert.Equal(t, expectedUlimits, cc.Ulimits)
***REMOVED***

func TestDaemonConfigurationMergeShmSize(t *testing.T) ***REMOVED***
	data := `***REMOVED***"default-shm-size": "1g"***REMOVED***`

	file := fs.NewFile(t, "docker-config", fs.WithContent(data))
	defer file.Remove()

	c := &Config***REMOVED******REMOVED***

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	shmSize := opts.MemBytes(DefaultShmSize)
	flags.Var(&shmSize, "default-shm-size", "")

	cc, err := MergeDaemonConfigurations(c, flags, file.Path())
	require.NoError(t, err)

	expectedValue := 1 * 1024 * 1024 * 1024
	assert.Equal(t, int64(expectedValue), cc.ShmSize.Value())
***REMOVED***
