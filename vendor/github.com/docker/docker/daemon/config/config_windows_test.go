// +build windows

package config

import (
	"io/ioutil"
	"testing"

	"github.com/docker/docker/opts"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDaemonConfigurationMerge(t *testing.T) ***REMOVED***
	f, err := ioutil.TempFile("", "docker-config-")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	configFile := f.Name()

	f.Write([]byte(`
		***REMOVED***
			"debug": true,
			"log-opts": ***REMOVED***
				"tag": "test_tag"
			***REMOVED***
		***REMOVED***`))

	f.Close()

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
	flags.Var(opts.NewNamedMapOpts("log-opts", nil, nil), "log-opt", "")

	cc, err := MergeDaemonConfigurations(c, flags, configFile)
	require.NoError(t, err)

	assert.True(t, cc.Debug)
	assert.True(t, cc.AutoRestart)

	expectedLogConfig := LogConfig***REMOVED***
		Type:   "syslog",
		Config: map[string]string***REMOVED***"tag": "test_tag"***REMOVED***,
	***REMOVED***

	assert.Equal(t, expectedLogConfig, cc.LogConfig)
***REMOVED***
