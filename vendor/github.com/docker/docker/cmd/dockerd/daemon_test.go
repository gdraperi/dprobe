package main

import (
	"testing"

	"github.com/docker/docker/daemon/config"
	"github.com/docker/docker/internal/testutil"
	"github.com/gotestyourself/gotestyourself/fs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func defaultOptions(configFile string) *daemonOptions ***REMOVED***
	opts := newDaemonOptions(&config.Config***REMOVED******REMOVED***)
	opts.flags = &pflag.FlagSet***REMOVED******REMOVED***
	opts.InstallFlags(opts.flags)
	installConfigFlags(opts.daemonConfig, opts.flags)
	opts.flags.StringVar(&opts.configFile, "config-file", defaultDaemonConfigFile, "")
	opts.configFile = configFile
	return opts
***REMOVED***

func TestLoadDaemonCliConfigWithoutOverriding(t *testing.T) ***REMOVED***
	opts := defaultOptions("")
	opts.Debug = true

	loadedConfig, err := loadDaemonCliConfig(opts)
	require.NoError(t, err)
	require.NotNil(t, loadedConfig)
	if !loadedConfig.Debug ***REMOVED***
		t.Fatalf("expected debug to be copied from the common flags, got false")
	***REMOVED***
***REMOVED***

func TestLoadDaemonCliConfigWithTLS(t *testing.T) ***REMOVED***
	opts := defaultOptions("")
	opts.TLSOptions.CAFile = "/tmp/ca.pem"
	opts.TLS = true

	loadedConfig, err := loadDaemonCliConfig(opts)
	require.NoError(t, err)
	require.NotNil(t, loadedConfig)
	assert.Equal(t, "/tmp/ca.pem", loadedConfig.CommonTLSOptions.CAFile)
***REMOVED***

func TestLoadDaemonCliConfigWithConflicts(t *testing.T) ***REMOVED***
	tempFile := fs.NewFile(t, "config", fs.WithContent(`***REMOVED***"labels": ["l3=foo"]***REMOVED***`))
	defer tempFile.Remove()
	configFile := tempFile.Path()

	opts := defaultOptions(configFile)
	flags := opts.flags

	assert.NoError(t, flags.Set("config-file", configFile))
	assert.NoError(t, flags.Set("label", "l1=bar"))
	assert.NoError(t, flags.Set("label", "l2=baz"))

	_, err := loadDaemonCliConfig(opts)
	testutil.ErrorContains(t, err, "as a flag and in the configuration file: labels")
***REMOVED***

func TestLoadDaemonCliWithConflictingLabels(t *testing.T) ***REMOVED***
	opts := defaultOptions("")
	flags := opts.flags

	assert.NoError(t, flags.Set("label", "foo=bar"))
	assert.NoError(t, flags.Set("label", "foo=baz"))

	_, err := loadDaemonCliConfig(opts)
	assert.EqualError(t, err, "conflict labels for foo=baz and foo=bar")
***REMOVED***

func TestLoadDaemonCliWithDuplicateLabels(t *testing.T) ***REMOVED***
	opts := defaultOptions("")
	flags := opts.flags

	assert.NoError(t, flags.Set("label", "foo=the-same"))
	assert.NoError(t, flags.Set("label", "foo=the-same"))

	_, err := loadDaemonCliConfig(opts)
	assert.NoError(t, err)
***REMOVED***

func TestLoadDaemonCliConfigWithTLSVerify(t *testing.T) ***REMOVED***
	tempFile := fs.NewFile(t, "config", fs.WithContent(`***REMOVED***"tlsverify": true***REMOVED***`))
	defer tempFile.Remove()

	opts := defaultOptions(tempFile.Path())
	opts.TLSOptions.CAFile = "/tmp/ca.pem"

	loadedConfig, err := loadDaemonCliConfig(opts)
	require.NoError(t, err)
	require.NotNil(t, loadedConfig)
	assert.Equal(t, loadedConfig.TLS, true)
***REMOVED***

func TestLoadDaemonCliConfigWithExplicitTLSVerifyFalse(t *testing.T) ***REMOVED***
	tempFile := fs.NewFile(t, "config", fs.WithContent(`***REMOVED***"tlsverify": false***REMOVED***`))
	defer tempFile.Remove()

	opts := defaultOptions(tempFile.Path())
	opts.TLSOptions.CAFile = "/tmp/ca.pem"

	loadedConfig, err := loadDaemonCliConfig(opts)
	require.NoError(t, err)
	require.NotNil(t, loadedConfig)
	assert.True(t, loadedConfig.TLS)
***REMOVED***

func TestLoadDaemonCliConfigWithoutTLSVerify(t *testing.T) ***REMOVED***
	tempFile := fs.NewFile(t, "config", fs.WithContent(`***REMOVED******REMOVED***`))
	defer tempFile.Remove()

	opts := defaultOptions(tempFile.Path())
	opts.TLSOptions.CAFile = "/tmp/ca.pem"

	loadedConfig, err := loadDaemonCliConfig(opts)
	require.NoError(t, err)
	require.NotNil(t, loadedConfig)
	assert.False(t, loadedConfig.TLS)
***REMOVED***

func TestLoadDaemonCliConfigWithLogLevel(t *testing.T) ***REMOVED***
	tempFile := fs.NewFile(t, "config", fs.WithContent(`***REMOVED***"log-level": "warn"***REMOVED***`))
	defer tempFile.Remove()

	opts := defaultOptions(tempFile.Path())
	loadedConfig, err := loadDaemonCliConfig(opts)
	require.NoError(t, err)
	require.NotNil(t, loadedConfig)
	assert.Equal(t, "warn", loadedConfig.LogLevel)
	assert.Equal(t, logrus.WarnLevel, logrus.GetLevel())
***REMOVED***

func TestLoadDaemonConfigWithEmbeddedOptions(t *testing.T) ***REMOVED***
	content := `***REMOVED***"tlscacert": "/etc/certs/ca.pem", "log-driver": "syslog"***REMOVED***`
	tempFile := fs.NewFile(t, "config", fs.WithContent(content))
	defer tempFile.Remove()

	opts := defaultOptions(tempFile.Path())
	loadedConfig, err := loadDaemonCliConfig(opts)
	require.NoError(t, err)
	require.NotNil(t, loadedConfig)
	assert.Equal(t, "/etc/certs/ca.pem", loadedConfig.CommonTLSOptions.CAFile)
	assert.Equal(t, "syslog", loadedConfig.LogConfig.Type)
***REMOVED***

func TestLoadDaemonConfigWithRegistryOptions(t *testing.T) ***REMOVED***
	content := `***REMOVED***
		"allow-nondistributable-artifacts": ["allow-nondistributable-artifacts.com"],
		"registry-mirrors": ["https://mirrors.docker.com"],
		"insecure-registries": ["https://insecure.docker.com"]
	***REMOVED***`
	tempFile := fs.NewFile(t, "config", fs.WithContent(content))
	defer tempFile.Remove()

	opts := defaultOptions(tempFile.Path())
	loadedConfig, err := loadDaemonCliConfig(opts)
	require.NoError(t, err)
	require.NotNil(t, loadedConfig)

	assert.Len(t, loadedConfig.AllowNondistributableArtifacts, 1)
	assert.Len(t, loadedConfig.Mirrors, 1)
	assert.Len(t, loadedConfig.InsecureRegistries, 1)
***REMOVED***
