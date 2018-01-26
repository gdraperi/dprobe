package config

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/docker/docker/daemon/discovery"
	"github.com/docker/docker/internal/testutil"
	"github.com/docker/docker/opts"
	"github.com/gotestyourself/gotestyourself/fs"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestDaemonConfigurationNotFound(t *testing.T) ***REMOVED***
	_, err := MergeDaemonConfigurations(&Config***REMOVED******REMOVED***, nil, "/tmp/foo-bar-baz-docker")
	if err == nil || !os.IsNotExist(err) ***REMOVED***
		t.Fatalf("expected does not exist error, got %v", err)
	***REMOVED***
***REMOVED***

func TestDaemonBrokenConfiguration(t *testing.T) ***REMOVED***
	f, err := ioutil.TempFile("", "docker-config-")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	configFile := f.Name()
	f.Write([]byte(`***REMOVED***"Debug": tru`))
	f.Close()

	_, err = MergeDaemonConfigurations(&Config***REMOVED******REMOVED***, nil, configFile)
	if err == nil ***REMOVED***
		t.Fatalf("expected error, got %v", err)
	***REMOVED***
***REMOVED***

func TestParseClusterAdvertiseSettings(t *testing.T) ***REMOVED***
	_, err := ParseClusterAdvertiseSettings("something", "")
	if err != discovery.ErrDiscoveryDisabled ***REMOVED***
		t.Fatalf("expected discovery disabled error, got %v\n", err)
	***REMOVED***

	_, err = ParseClusterAdvertiseSettings("", "something")
	if err == nil ***REMOVED***
		t.Fatalf("expected discovery store error, got %v\n", err)
	***REMOVED***

	_, err = ParseClusterAdvertiseSettings("etcd", "127.0.0.1:8080")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestFindConfigurationConflicts(t *testing.T) ***REMOVED***
	config := map[string]interface***REMOVED******REMOVED******REMOVED***"authorization-plugins": "foobar"***REMOVED***
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)

	flags.String("authorization-plugins", "", "")
	assert.NoError(t, flags.Set("authorization-plugins", "asdf"))

	testutil.ErrorContains(t,
		findConfigurationConflicts(config, flags),
		"authorization-plugins: (from flag: asdf, from file: foobar)")
***REMOVED***

func TestFindConfigurationConflictsWithNamedOptions(t *testing.T) ***REMOVED***
	config := map[string]interface***REMOVED******REMOVED******REMOVED***"hosts": []string***REMOVED***"qwer"***REMOVED******REMOVED***
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)

	var hosts []string
	flags.VarP(opts.NewNamedListOptsRef("hosts", &hosts, opts.ValidateHost), "host", "H", "Daemon socket(s) to connect to")
	assert.NoError(t, flags.Set("host", "tcp://127.0.0.1:4444"))
	assert.NoError(t, flags.Set("host", "unix:///var/run/docker.sock"))

	testutil.ErrorContains(t, findConfigurationConflicts(config, flags), "hosts")
***REMOVED***

func TestDaemonConfigurationMergeConflicts(t *testing.T) ***REMOVED***
	f, err := ioutil.TempFile("", "docker-config-")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	configFile := f.Name()
	f.Write([]byte(`***REMOVED***"debug": true***REMOVED***`))
	f.Close()

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.Bool("debug", false, "")
	flags.Set("debug", "false")

	_, err = MergeDaemonConfigurations(&Config***REMOVED******REMOVED***, flags, configFile)
	if err == nil ***REMOVED***
		t.Fatal("expected error, got nil")
	***REMOVED***
	if !strings.Contains(err.Error(), "debug") ***REMOVED***
		t.Fatalf("expected debug conflict, got %v", err)
	***REMOVED***
***REMOVED***

func TestDaemonConfigurationMergeConcurrent(t *testing.T) ***REMOVED***
	f, err := ioutil.TempFile("", "docker-config-")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	configFile := f.Name()
	f.Write([]byte(`***REMOVED***"max-concurrent-downloads": 1***REMOVED***`))
	f.Close()

	_, err = MergeDaemonConfigurations(&Config***REMOVED******REMOVED***, nil, configFile)
	if err != nil ***REMOVED***
		t.Fatal("expected error, got nil")
	***REMOVED***
***REMOVED***

func TestDaemonConfigurationMergeConcurrentError(t *testing.T) ***REMOVED***
	f, err := ioutil.TempFile("", "docker-config-")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	configFile := f.Name()
	f.Write([]byte(`***REMOVED***"max-concurrent-downloads": -1***REMOVED***`))
	f.Close()

	_, err = MergeDaemonConfigurations(&Config***REMOVED******REMOVED***, nil, configFile)
	if err == nil ***REMOVED***
		t.Fatalf("expected no error, got error %v", err)
	***REMOVED***
***REMOVED***

func TestDaemonConfigurationMergeConflictsWithInnerStructs(t *testing.T) ***REMOVED***
	f, err := ioutil.TempFile("", "docker-config-")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	configFile := f.Name()
	f.Write([]byte(`***REMOVED***"tlscacert": "/etc/certificates/ca.pem"***REMOVED***`))
	f.Close()

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("tlscacert", "", "")
	flags.Set("tlscacert", "~/.docker/ca.pem")

	_, err = MergeDaemonConfigurations(&Config***REMOVED******REMOVED***, flags, configFile)
	if err == nil ***REMOVED***
		t.Fatal("expected error, got nil")
	***REMOVED***
	if !strings.Contains(err.Error(), "tlscacert") ***REMOVED***
		t.Fatalf("expected tlscacert conflict, got %v", err)
	***REMOVED***
***REMOVED***

func TestFindConfigurationConflictsWithUnknownKeys(t *testing.T) ***REMOVED***
	config := map[string]interface***REMOVED******REMOVED******REMOVED***"tls-verify": "true"***REMOVED***
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)

	flags.Bool("tlsverify", false, "")
	err := findConfigurationConflicts(config, flags)
	if err == nil ***REMOVED***
		t.Fatal("expected error, got nil")
	***REMOVED***
	if !strings.Contains(err.Error(), "the following directives don't match any configuration option: tls-verify") ***REMOVED***
		t.Fatalf("expected tls-verify conflict, got %v", err)
	***REMOVED***
***REMOVED***

func TestFindConfigurationConflictsWithMergedValues(t *testing.T) ***REMOVED***
	var hosts []string
	config := map[string]interface***REMOVED******REMOVED******REMOVED***"hosts": "tcp://127.0.0.1:2345"***REMOVED***
	flags := pflag.NewFlagSet("base", pflag.ContinueOnError)
	flags.VarP(opts.NewNamedListOptsRef("hosts", &hosts, nil), "host", "H", "")

	err := findConfigurationConflicts(config, flags)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	flags.Set("host", "unix:///var/run/docker.sock")
	err = findConfigurationConflicts(config, flags)
	if err == nil ***REMOVED***
		t.Fatal("expected error, got nil")
	***REMOVED***
	if !strings.Contains(err.Error(), "hosts: (from flag: [unix:///var/run/docker.sock], from file: tcp://127.0.0.1:2345)") ***REMOVED***
		t.Fatalf("expected hosts conflict, got %v", err)
	***REMOVED***
***REMOVED***

func TestValidateConfigurationErrors(t *testing.T) ***REMOVED***
	minusNumber := -10
	testCases := []struct ***REMOVED***
		config *Config
	***REMOVED******REMOVED***
		***REMOVED***
			config: &Config***REMOVED***
				CommonConfig: CommonConfig***REMOVED***
					Labels: []string***REMOVED***"one"***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			config: &Config***REMOVED***
				CommonConfig: CommonConfig***REMOVED***
					Labels: []string***REMOVED***"foo=bar", "one"***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			config: &Config***REMOVED***
				CommonConfig: CommonConfig***REMOVED***
					DNS: []string***REMOVED***"1.1.1.1o"***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			config: &Config***REMOVED***
				CommonConfig: CommonConfig***REMOVED***
					DNS: []string***REMOVED***"2.2.2.2", "1.1.1.1o"***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			config: &Config***REMOVED***
				CommonConfig: CommonConfig***REMOVED***
					DNSSearch: []string***REMOVED***"123456"***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			config: &Config***REMOVED***
				CommonConfig: CommonConfig***REMOVED***
					DNSSearch: []string***REMOVED***"a.b.c", "123456"***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			config: &Config***REMOVED***
				CommonConfig: CommonConfig***REMOVED***
					MaxConcurrentDownloads: &minusNumber,
					// This is weird...
					ValuesSet: map[string]interface***REMOVED******REMOVED******REMOVED***
						"max-concurrent-downloads": -1,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			config: &Config***REMOVED***
				CommonConfig: CommonConfig***REMOVED***
					MaxConcurrentUploads: &minusNumber,
					// This is weird...
					ValuesSet: map[string]interface***REMOVED******REMOVED******REMOVED***
						"max-concurrent-uploads": -1,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			config: &Config***REMOVED***
				CommonConfig: CommonConfig***REMOVED***
					NodeGenericResources: []string***REMOVED***"foo"***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			config: &Config***REMOVED***
				CommonConfig: CommonConfig***REMOVED***
					NodeGenericResources: []string***REMOVED***"foo=bar", "foo=1"***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		err := Validate(tc.config)
		if err == nil ***REMOVED***
			t.Fatalf("expected error, got nil for config %v", tc.config)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestValidateConfiguration(t *testing.T) ***REMOVED***
	minusNumber := 4
	testCases := []struct ***REMOVED***
		config *Config
	***REMOVED******REMOVED***
		***REMOVED***
			config: &Config***REMOVED***
				CommonConfig: CommonConfig***REMOVED***
					Labels: []string***REMOVED***"one=two"***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			config: &Config***REMOVED***
				CommonConfig: CommonConfig***REMOVED***
					DNS: []string***REMOVED***"1.1.1.1"***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			config: &Config***REMOVED***
				CommonConfig: CommonConfig***REMOVED***
					DNSSearch: []string***REMOVED***"a.b.c"***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			config: &Config***REMOVED***
				CommonConfig: CommonConfig***REMOVED***
					MaxConcurrentDownloads: &minusNumber,
					// This is weird...
					ValuesSet: map[string]interface***REMOVED******REMOVED******REMOVED***
						"max-concurrent-downloads": -1,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			config: &Config***REMOVED***
				CommonConfig: CommonConfig***REMOVED***
					MaxConcurrentUploads: &minusNumber,
					// This is weird...
					ValuesSet: map[string]interface***REMOVED******REMOVED******REMOVED***
						"max-concurrent-uploads": -1,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			config: &Config***REMOVED***
				CommonConfig: CommonConfig***REMOVED***
					NodeGenericResources: []string***REMOVED***"foo=bar", "foo=baz"***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			config: &Config***REMOVED***
				CommonConfig: CommonConfig***REMOVED***
					NodeGenericResources: []string***REMOVED***"foo=1"***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		err := Validate(tc.config)
		if err != nil ***REMOVED***
			t.Fatalf("expected no error, got error %v", err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestModifiedDiscoverySettings(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		current  *Config
		modified *Config
		expected bool
	***REMOVED******REMOVED***
		***REMOVED***
			current:  discoveryConfig("foo", "bar", map[string]string***REMOVED******REMOVED***),
			modified: discoveryConfig("foo", "bar", map[string]string***REMOVED******REMOVED***),
			expected: false,
		***REMOVED***,
		***REMOVED***
			current:  discoveryConfig("foo", "bar", map[string]string***REMOVED***"foo": "bar"***REMOVED***),
			modified: discoveryConfig("foo", "bar", map[string]string***REMOVED***"foo": "bar"***REMOVED***),
			expected: false,
		***REMOVED***,
		***REMOVED***
			current:  discoveryConfig("foo", "bar", map[string]string***REMOVED******REMOVED***),
			modified: discoveryConfig("foo", "bar", nil),
			expected: false,
		***REMOVED***,
		***REMOVED***
			current:  discoveryConfig("foo", "bar", nil),
			modified: discoveryConfig("foo", "bar", map[string]string***REMOVED******REMOVED***),
			expected: false,
		***REMOVED***,
		***REMOVED***
			current:  discoveryConfig("foo", "bar", nil),
			modified: discoveryConfig("baz", "bar", nil),
			expected: true,
		***REMOVED***,
		***REMOVED***
			current:  discoveryConfig("foo", "bar", nil),
			modified: discoveryConfig("foo", "baz", nil),
			expected: true,
		***REMOVED***,
		***REMOVED***
			current:  discoveryConfig("foo", "bar", nil),
			modified: discoveryConfig("foo", "bar", map[string]string***REMOVED***"foo": "bar"***REMOVED***),
			expected: true,
		***REMOVED***,
	***REMOVED***

	for _, c := range cases ***REMOVED***
		got := ModifiedDiscoverySettings(c.current, c.modified.ClusterStore, c.modified.ClusterAdvertise, c.modified.ClusterOpts)
		if c.expected != got ***REMOVED***
			t.Fatalf("expected %v, got %v: current config %v, new config %v", c.expected, got, c.current, c.modified)
		***REMOVED***
	***REMOVED***
***REMOVED***

func discoveryConfig(backendAddr, advertiseAddr string, opts map[string]string) *Config ***REMOVED***
	return &Config***REMOVED***
		CommonConfig: CommonConfig***REMOVED***
			ClusterStore:     backendAddr,
			ClusterAdvertise: advertiseAddr,
			ClusterOpts:      opts,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// TestReloadSetConfigFileNotExist tests that when `--config-file` is set
// and it doesn't exist the `Reload` function returns an error.
func TestReloadSetConfigFileNotExist(t *testing.T) ***REMOVED***
	configFile := "/tmp/blabla/not/exists/config.json"
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("config-file", "", "")
	flags.Set("config-file", configFile)

	err := Reload(configFile, flags, func(c *Config) ***REMOVED******REMOVED***)
	assert.Error(t, err)
	testutil.ErrorContains(t, err, "unable to configure the Docker daemon with file")
***REMOVED***

// TestReloadDefaultConfigNotExist tests that if the default configuration file
// doesn't exist the daemon still will be reloaded.
func TestReloadDefaultConfigNotExist(t *testing.T) ***REMOVED***
	reloaded := false
	configFile := "/etc/docker/daemon.json"
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("config-file", configFile, "")
	err := Reload(configFile, flags, func(c *Config) ***REMOVED***
		reloaded = true
	***REMOVED***)
	assert.Nil(t, err)
	assert.True(t, reloaded)
***REMOVED***

// TestReloadBadDefaultConfig tests that when `--config-file` is not set
// and the default configuration file exists and is bad return an error
func TestReloadBadDefaultConfig(t *testing.T) ***REMOVED***
	f, err := ioutil.TempFile("", "docker-config-")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	configFile := f.Name()
	f.Write([]byte(`***REMOVED***wrong: "configuration"***REMOVED***`))
	f.Close()

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("config-file", configFile, "")
	err = Reload(configFile, flags, func(c *Config) ***REMOVED******REMOVED***)
	assert.Error(t, err)
	testutil.ErrorContains(t, err, "unable to configure the Docker daemon with file")
***REMOVED***

func TestReloadWithConflictingLabels(t *testing.T) ***REMOVED***
	tempFile := fs.NewFile(t, "config", fs.WithContent(`***REMOVED***"labels":["foo=bar","foo=baz"]***REMOVED***`))
	defer tempFile.Remove()
	configFile := tempFile.Path()

	var lbls []string
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("config-file", configFile, "")
	flags.StringSlice("labels", lbls, "")
	err := Reload(configFile, flags, func(c *Config) ***REMOVED******REMOVED***)
	testutil.ErrorContains(t, err, "conflict labels for foo=baz and foo=bar")
***REMOVED***

func TestReloadWithDuplicateLabels(t *testing.T) ***REMOVED***
	tempFile := fs.NewFile(t, "config", fs.WithContent(`***REMOVED***"labels":["foo=the-same","foo=the-same"]***REMOVED***`))
	defer tempFile.Remove()
	configFile := tempFile.Path()

	var lbls []string
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("config-file", configFile, "")
	flags.StringSlice("labels", lbls, "")
	err := Reload(configFile, flags, func(c *Config) ***REMOVED******REMOVED***)
	assert.NoError(t, err)
***REMOVED***
