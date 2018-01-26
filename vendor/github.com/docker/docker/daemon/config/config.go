package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"

	daemondiscovery "github.com/docker/docker/daemon/discovery"
	"github.com/docker/docker/opts"
	"github.com/docker/docker/pkg/authorization"
	"github.com/docker/docker/pkg/discovery"
	"github.com/docker/docker/registry"
	"github.com/imdario/mergo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

const (
	// DefaultMaxConcurrentDownloads is the default value for
	// maximum number of downloads that
	// may take place at a time for each pull.
	DefaultMaxConcurrentDownloads = 3
	// DefaultMaxConcurrentUploads is the default value for
	// maximum number of uploads that
	// may take place at a time for each push.
	DefaultMaxConcurrentUploads = 5
	// StockRuntimeName is the reserved name/alias used to represent the
	// OCI runtime being shipped with the docker daemon package.
	StockRuntimeName = "runc"
	// DefaultShmSize is the default value for container's shm size
	DefaultShmSize = int64(67108864)
	// DefaultNetworkMtu is the default value for network MTU
	DefaultNetworkMtu = 1500
	// DisableNetworkBridge is the default value of the option to disable network bridge
	DisableNetworkBridge = "none"
	// DefaultInitBinary is the name of the default init binary
	DefaultInitBinary = "docker-init"
)

// flatOptions contains configuration keys
// that MUST NOT be parsed as deep structures.
// Use this to differentiate these options
// with others like the ones in CommonTLSOptions.
var flatOptions = map[string]bool***REMOVED***
	"cluster-store-opts": true,
	"log-opts":           true,
	"runtimes":           true,
	"default-ulimits":    true,
***REMOVED***

// LogConfig represents the default log configuration.
// It includes json tags to deserialize configuration from a file
// using the same names that the flags in the command line use.
type LogConfig struct ***REMOVED***
	Type   string            `json:"log-driver,omitempty"`
	Config map[string]string `json:"log-opts,omitempty"`
***REMOVED***

// commonBridgeConfig stores all the platform-common bridge driver specific
// configuration.
type commonBridgeConfig struct ***REMOVED***
	Iface     string `json:"bridge,omitempty"`
	FixedCIDR string `json:"fixed-cidr,omitempty"`
***REMOVED***

// CommonTLSOptions defines TLS configuration for the daemon server.
// It includes json tags to deserialize configuration from a file
// using the same names that the flags in the command line use.
type CommonTLSOptions struct ***REMOVED***
	CAFile   string `json:"tlscacert,omitempty"`
	CertFile string `json:"tlscert,omitempty"`
	KeyFile  string `json:"tlskey,omitempty"`
***REMOVED***

// CommonConfig defines the configuration of a docker daemon which is
// common across platforms.
// It includes json tags to deserialize configuration from a file
// using the same names that the flags in the command line use.
type CommonConfig struct ***REMOVED***
	AuthzMiddleware       *authorization.Middleware `json:"-"`
	AuthorizationPlugins  []string                  `json:"authorization-plugins,omitempty"` // AuthorizationPlugins holds list of authorization plugins
	AutoRestart           bool                      `json:"-"`
	Context               map[string][]string       `json:"-"`
	DisableBridge         bool                      `json:"-"`
	DNS                   []string                  `json:"dns,omitempty"`
	DNSOptions            []string                  `json:"dns-opts,omitempty"`
	DNSSearch             []string                  `json:"dns-search,omitempty"`
	ExecOptions           []string                  `json:"exec-opts,omitempty"`
	GraphDriver           string                    `json:"storage-driver,omitempty"`
	GraphOptions          []string                  `json:"storage-opts,omitempty"`
	Labels                []string                  `json:"labels,omitempty"`
	Mtu                   int                       `json:"mtu,omitempty"`
	NetworkDiagnosticPort int                       `json:"network-diagnostic-port,omitempty"`
	Pidfile               string                    `json:"pidfile,omitempty"`
	RawLogs               bool                      `json:"raw-logs,omitempty"`
	RootDeprecated        string                    `json:"graph,omitempty"`
	Root                  string                    `json:"data-root,omitempty"`
	ExecRoot              string                    `json:"exec-root,omitempty"`
	SocketGroup           string                    `json:"group,omitempty"`
	CorsHeaders           string                    `json:"api-cors-header,omitempty"`

	// TrustKeyPath is used to generate the daemon ID and for signing schema 1 manifests
	// when pushing to a registry which does not support schema 2. This field is marked as
	// deprecated because schema 1 manifests are deprecated in favor of schema 2 and the
	// daemon ID will use a dedicated identifier not shared with exported signatures.
	TrustKeyPath string `json:"deprecated-key-path,omitempty"`

	// LiveRestoreEnabled determines whether we should keep containers
	// alive upon daemon shutdown/start
	LiveRestoreEnabled bool `json:"live-restore,omitempty"`

	// ClusterStore is the storage backend used for the cluster information. It is used by both
	// multihost networking (to store networks and endpoints information) and by the node discovery
	// mechanism.
	ClusterStore string `json:"cluster-store,omitempty"`

	// ClusterOpts is used to pass options to the discovery package for tuning libkv settings, such
	// as TLS configuration settings.
	ClusterOpts map[string]string `json:"cluster-store-opts,omitempty"`

	// ClusterAdvertise is the network endpoint that the Engine advertises for the purpose of node
	// discovery. This should be a 'host:port' combination on which that daemon instance is
	// reachable by other hosts.
	ClusterAdvertise string `json:"cluster-advertise,omitempty"`

	// MaxConcurrentDownloads is the maximum number of downloads that
	// may take place at a time for each pull.
	MaxConcurrentDownloads *int `json:"max-concurrent-downloads,omitempty"`

	// MaxConcurrentUploads is the maximum number of uploads that
	// may take place at a time for each push.
	MaxConcurrentUploads *int `json:"max-concurrent-uploads,omitempty"`

	// ShutdownTimeout is the timeout value (in seconds) the daemon will wait for the container
	// to stop when daemon is being shutdown
	ShutdownTimeout int `json:"shutdown-timeout,omitempty"`

	Debug     bool     `json:"debug,omitempty"`
	Hosts     []string `json:"hosts,omitempty"`
	LogLevel  string   `json:"log-level,omitempty"`
	TLS       bool     `json:"tls,omitempty"`
	TLSVerify bool     `json:"tlsverify,omitempty"`

	// Embedded structs that allow config
	// deserialization without the full struct.
	CommonTLSOptions

	// SwarmDefaultAdvertiseAddr is the default host/IP or network interface
	// to use if a wildcard address is specified in the ListenAddr value
	// given to the /swarm/init endpoint and no advertise address is
	// specified.
	SwarmDefaultAdvertiseAddr string `json:"swarm-default-advertise-addr"`
	MetricsAddress            string `json:"metrics-addr"`

	LogConfig
	BridgeConfig // bridgeConfig holds bridge network specific configuration.
	registry.ServiceOptions

	sync.Mutex
	// FIXME(vdemeester) This part is not that clear and is mainly dependent on cli flags
	// It should probably be handled outside this package.
	ValuesSet map[string]interface***REMOVED******REMOVED*** `json:"-"`

	Experimental bool `json:"experimental"` // Experimental indicates whether experimental features should be exposed or not

	// Exposed node Generic Resources
	// e.g: ["orange=red", "orange=green", "orange=blue", "apple=3"]
	NodeGenericResources []string `json:"node-generic-resources,omitempty"`
	// NetworkControlPlaneMTU allows to specify the control plane MTU, this will allow to optimize the network use in some components
	NetworkControlPlaneMTU int `json:"network-control-plane-mtu,omitempty"`

	// ContainerAddr is the address used to connect to containerd if we're
	// not starting it ourselves
	ContainerdAddr string `json:"containerd,omitempty"`
***REMOVED***

// IsValueSet returns true if a configuration value
// was explicitly set in the configuration file.
func (conf *Config) IsValueSet(name string) bool ***REMOVED***
	if conf.ValuesSet == nil ***REMOVED***
		return false
	***REMOVED***
	_, ok := conf.ValuesSet[name]
	return ok
***REMOVED***

// New returns a new fully initialized Config struct
func New() *Config ***REMOVED***
	config := Config***REMOVED******REMOVED***
	config.LogConfig.Config = make(map[string]string)
	config.ClusterOpts = make(map[string]string)

	if runtime.GOOS != "linux" ***REMOVED***
		config.V2Only = true
	***REMOVED***
	return &config
***REMOVED***

// ParseClusterAdvertiseSettings parses the specified advertise settings
func ParseClusterAdvertiseSettings(clusterStore, clusterAdvertise string) (string, error) ***REMOVED***
	if clusterAdvertise == "" ***REMOVED***
		return "", daemondiscovery.ErrDiscoveryDisabled
	***REMOVED***
	if clusterStore == "" ***REMOVED***
		return "", errors.New("invalid cluster configuration. --cluster-advertise must be accompanied by --cluster-store configuration")
	***REMOVED***

	advertise, err := discovery.ParseAdvertise(clusterAdvertise)
	if err != nil ***REMOVED***
		return "", fmt.Errorf("discovery advertise parsing failed (%v)", err)
	***REMOVED***
	return advertise, nil
***REMOVED***

// GetConflictFreeLabels validates Labels for conflict
// In swarm the duplicates for labels are removed
// so we only take same values here, no conflict values
// If the key-value is the same we will only take the last label
func GetConflictFreeLabels(labels []string) ([]string, error) ***REMOVED***
	labelMap := map[string]string***REMOVED******REMOVED***
	for _, label := range labels ***REMOVED***
		stringSlice := strings.SplitN(label, "=", 2)
		if len(stringSlice) > 1 ***REMOVED***
			// If there is a conflict we will return an error
			if v, ok := labelMap[stringSlice[0]]; ok && v != stringSlice[1] ***REMOVED***
				return nil, fmt.Errorf("conflict labels for %s=%s and %s=%s", stringSlice[0], stringSlice[1], stringSlice[0], v)
			***REMOVED***
			labelMap[stringSlice[0]] = stringSlice[1]
		***REMOVED***
	***REMOVED***

	newLabels := []string***REMOVED******REMOVED***
	for k, v := range labelMap ***REMOVED***
		newLabels = append(newLabels, fmt.Sprintf("%s=%s", k, v))
	***REMOVED***
	return newLabels, nil
***REMOVED***

// Reload reads the configuration in the host and reloads the daemon and server.
func Reload(configFile string, flags *pflag.FlagSet, reload func(*Config)) error ***REMOVED***
	logrus.Infof("Got signal to reload configuration, reloading from: %s", configFile)
	newConfig, err := getConflictFreeConfiguration(configFile, flags)
	if err != nil ***REMOVED***
		if flags.Changed("config-file") || !os.IsNotExist(err) ***REMOVED***
			return fmt.Errorf("unable to configure the Docker daemon with file %s: %v", configFile, err)
		***REMOVED***
		newConfig = New()
	***REMOVED***

	if err := Validate(newConfig); err != nil ***REMOVED***
		return fmt.Errorf("file configuration validation failed (%v)", err)
	***REMOVED***

	// Check if duplicate label-keys with different values are found
	newLabels, err := GetConflictFreeLabels(newConfig.Labels)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	newConfig.Labels = newLabels

	reload(newConfig)
	return nil
***REMOVED***

// boolValue is an interface that boolean value flags implement
// to tell the command line how to make -name equivalent to -name=true.
type boolValue interface ***REMOVED***
	IsBoolFlag() bool
***REMOVED***

// MergeDaemonConfigurations reads a configuration file,
// loads the file configuration in an isolated structure,
// and merges the configuration provided from flags on top
// if there are no conflicts.
func MergeDaemonConfigurations(flagsConfig *Config, flags *pflag.FlagSet, configFile string) (*Config, error) ***REMOVED***
	fileConfig, err := getConflictFreeConfiguration(configFile, flags)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := Validate(fileConfig); err != nil ***REMOVED***
		return nil, fmt.Errorf("configuration validation from file failed (%v)", err)
	***REMOVED***

	// merge flags configuration on top of the file configuration
	if err := mergo.Merge(fileConfig, flagsConfig); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// We need to validate again once both fileConfig and flagsConfig
	// have been merged
	if err := Validate(fileConfig); err != nil ***REMOVED***
		return nil, fmt.Errorf("merged configuration validation from file and command line flags failed (%v)", err)
	***REMOVED***

	return fileConfig, nil
***REMOVED***

// getConflictFreeConfiguration loads the configuration from a JSON file.
// It compares that configuration with the one provided by the flags,
// and returns an error if there are conflicts.
func getConflictFreeConfiguration(configFile string, flags *pflag.FlagSet) (*Config, error) ***REMOVED***
	b, err := ioutil.ReadFile(configFile)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var config Config
	var reader io.Reader
	if flags != nil ***REMOVED***
		var jsonConfig map[string]interface***REMOVED******REMOVED***
		reader = bytes.NewReader(b)
		if err := json.NewDecoder(reader).Decode(&jsonConfig); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		configSet := configValuesSet(jsonConfig)

		if err := findConfigurationConflicts(configSet, flags); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		// Override flag values to make sure the values set in the config file with nullable values, like `false`,
		// are not overridden by default truthy values from the flags that were not explicitly set.
		// See https://github.com/docker/docker/issues/20289 for an example.
		//
		// TODO: Rewrite configuration logic to avoid same issue with other nullable values, like numbers.
		namedOptions := make(map[string]interface***REMOVED******REMOVED***)
		for key, value := range configSet ***REMOVED***
			f := flags.Lookup(key)
			if f == nil ***REMOVED*** // ignore named flags that don't match
				namedOptions[key] = value
				continue
			***REMOVED***

			if _, ok := f.Value.(boolValue); ok ***REMOVED***
				f.Value.Set(fmt.Sprintf("%v", value))
			***REMOVED***
		***REMOVED***
		if len(namedOptions) > 0 ***REMOVED***
			// set also default for mergeVal flags that are boolValue at the same time.
			flags.VisitAll(func(f *pflag.Flag) ***REMOVED***
				if opt, named := f.Value.(opts.NamedOption); named ***REMOVED***
					v, set := namedOptions[opt.Name()]
					_, boolean := f.Value.(boolValue)
					if set && boolean ***REMOVED***
						f.Value.Set(fmt.Sprintf("%v", v))
					***REMOVED***
				***REMOVED***
			***REMOVED***)
		***REMOVED***

		config.ValuesSet = configSet
	***REMOVED***

	reader = bytes.NewReader(b)
	if err := json.NewDecoder(reader).Decode(&config); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if config.RootDeprecated != "" ***REMOVED***
		logrus.Warn(`The "graph" config file option is deprecated. Please use "data-root" instead.`)

		if config.Root != "" ***REMOVED***
			return nil, fmt.Errorf(`cannot specify both "graph" and "data-root" config file options`)
		***REMOVED***

		config.Root = config.RootDeprecated
	***REMOVED***

	return &config, nil
***REMOVED***

// configValuesSet returns the configuration values explicitly set in the file.
func configValuesSet(config map[string]interface***REMOVED******REMOVED***) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	flatten := make(map[string]interface***REMOVED******REMOVED***)
	for k, v := range config ***REMOVED***
		if m, isMap := v.(map[string]interface***REMOVED******REMOVED***); isMap && !flatOptions[k] ***REMOVED***
			for km, vm := range m ***REMOVED***
				flatten[km] = vm
			***REMOVED***
			continue
		***REMOVED***

		flatten[k] = v
	***REMOVED***
	return flatten
***REMOVED***

// findConfigurationConflicts iterates over the provided flags searching for
// duplicated configurations and unknown keys. It returns an error with all the conflicts if
// it finds any.
func findConfigurationConflicts(config map[string]interface***REMOVED******REMOVED***, flags *pflag.FlagSet) error ***REMOVED***
	// 1. Search keys from the file that we don't recognize as flags.
	unknownKeys := make(map[string]interface***REMOVED******REMOVED***)
	for key, value := range config ***REMOVED***
		if flag := flags.Lookup(key); flag == nil ***REMOVED***
			unknownKeys[key] = value
		***REMOVED***
	***REMOVED***

	// 2. Discard values that implement NamedOption.
	// Their configuration name differs from their flag name, like `labels` and `label`.
	if len(unknownKeys) > 0 ***REMOVED***
		unknownNamedConflicts := func(f *pflag.Flag) ***REMOVED***
			if namedOption, ok := f.Value.(opts.NamedOption); ok ***REMOVED***
				if _, valid := unknownKeys[namedOption.Name()]; valid ***REMOVED***
					delete(unknownKeys, namedOption.Name())
				***REMOVED***
			***REMOVED***
		***REMOVED***
		flags.VisitAll(unknownNamedConflicts)
	***REMOVED***

	if len(unknownKeys) > 0 ***REMOVED***
		var unknown []string
		for key := range unknownKeys ***REMOVED***
			unknown = append(unknown, key)
		***REMOVED***
		return fmt.Errorf("the following directives don't match any configuration option: %s", strings.Join(unknown, ", "))
	***REMOVED***

	var conflicts []string
	printConflict := func(name string, flagValue, fileValue interface***REMOVED******REMOVED***) string ***REMOVED***
		return fmt.Sprintf("%s: (from flag: %v, from file: %v)", name, flagValue, fileValue)
	***REMOVED***

	// 3. Search keys that are present as a flag and as a file option.
	duplicatedConflicts := func(f *pflag.Flag) ***REMOVED***
		// search option name in the json configuration payload if the value is a named option
		if namedOption, ok := f.Value.(opts.NamedOption); ok ***REMOVED***
			if optsValue, ok := config[namedOption.Name()]; ok ***REMOVED***
				conflicts = append(conflicts, printConflict(namedOption.Name(), f.Value.String(), optsValue))
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// search flag name in the json configuration payload
			for _, name := range []string***REMOVED***f.Name, f.Shorthand***REMOVED*** ***REMOVED***
				if value, ok := config[name]; ok ***REMOVED***
					conflicts = append(conflicts, printConflict(name, f.Value.String(), value))
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	flags.Visit(duplicatedConflicts)

	if len(conflicts) > 0 ***REMOVED***
		return fmt.Errorf("the following directives are specified both as a flag and in the configuration file: %s", strings.Join(conflicts, ", "))
	***REMOVED***
	return nil
***REMOVED***

// Validate validates some specific configs.
// such as config.DNS, config.Labels, config.DNSSearch,
// as well as config.MaxConcurrentDownloads, config.MaxConcurrentUploads.
func Validate(config *Config) error ***REMOVED***
	// validate DNS
	for _, dns := range config.DNS ***REMOVED***
		if _, err := opts.ValidateIPAddress(dns); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// validate DNSSearch
	for _, dnsSearch := range config.DNSSearch ***REMOVED***
		if _, err := opts.ValidateDNSSearch(dnsSearch); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// validate Labels
	for _, label := range config.Labels ***REMOVED***
		if _, err := opts.ValidateLabel(label); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	// validate MaxConcurrentDownloads
	if config.MaxConcurrentDownloads != nil && *config.MaxConcurrentDownloads < 0 ***REMOVED***
		return fmt.Errorf("invalid max concurrent downloads: %d", *config.MaxConcurrentDownloads)
	***REMOVED***
	// validate MaxConcurrentUploads
	if config.MaxConcurrentUploads != nil && *config.MaxConcurrentUploads < 0 ***REMOVED***
		return fmt.Errorf("invalid max concurrent uploads: %d", *config.MaxConcurrentUploads)
	***REMOVED***

	// validate that "default" runtime is not reset
	if runtimes := config.GetAllRuntimes(); len(runtimes) > 0 ***REMOVED***
		if _, ok := runtimes[StockRuntimeName]; ok ***REMOVED***
			return fmt.Errorf("runtime name '%s' is reserved", StockRuntimeName)
		***REMOVED***
	***REMOVED***

	if _, err := ParseGenericResources(config.NodeGenericResources); err != nil ***REMOVED***
		return err
	***REMOVED***

	if defaultRuntime := config.GetDefaultRuntimeName(); defaultRuntime != "" && defaultRuntime != StockRuntimeName ***REMOVED***
		runtimes := config.GetAllRuntimes()
		if _, ok := runtimes[defaultRuntime]; !ok ***REMOVED***
			return fmt.Errorf("specified default runtime '%s' does not exist", defaultRuntime)
		***REMOVED***
	***REMOVED***

	// validate platform-specific settings
	return config.ValidatePlatformConfig()
***REMOVED***

// ModifiedDiscoverySettings returns whether the discovery configuration has been modified or not.
func ModifiedDiscoverySettings(config *Config, backendType, advertise string, clusterOpts map[string]string) bool ***REMOVED***
	if config.ClusterStore != backendType || config.ClusterAdvertise != advertise ***REMOVED***
		return true
	***REMOVED***

	if (config.ClusterOpts == nil && clusterOpts == nil) ||
		(config.ClusterOpts == nil && len(clusterOpts) == 0) ||
		(len(config.ClusterOpts) == 0 && clusterOpts == nil) ***REMOVED***
		return false
	***REMOVED***

	return !reflect.DeepEqual(config.ClusterOpts, clusterOpts)
***REMOVED***
