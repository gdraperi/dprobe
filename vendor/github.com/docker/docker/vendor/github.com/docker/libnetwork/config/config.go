package config

import (
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/docker/docker/pkg/discovery"
	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/docker/libkv/store"
	"github.com/docker/libnetwork/cluster"
	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/netlabel"
	"github.com/docker/libnetwork/osl"
	"github.com/sirupsen/logrus"
)

const (
	warningThNetworkControlPlaneMTU = 1500
	minimumNetworkControlPlaneMTU   = 500
)

// Config encapsulates configurations of various Libnetwork components
type Config struct ***REMOVED***
	Daemon          DaemonCfg
	Cluster         ClusterCfg
	Scopes          map[string]*datastore.ScopeCfg
	ActiveSandboxes map[string]interface***REMOVED******REMOVED***
	PluginGetter    plugingetter.PluginGetter
***REMOVED***

// DaemonCfg represents libnetwork core configuration
type DaemonCfg struct ***REMOVED***
	Debug                  bool
	Experimental           bool
	DataDir                string
	DefaultNetwork         string
	DefaultDriver          string
	Labels                 []string
	DriverCfg              map[string]interface***REMOVED******REMOVED***
	ClusterProvider        cluster.Provider
	NetworkControlPlaneMTU int
***REMOVED***

// ClusterCfg represents cluster configuration
type ClusterCfg struct ***REMOVED***
	Watcher   discovery.Watcher
	Address   string
	Discovery string
	Heartbeat uint64
***REMOVED***

// LoadDefaultScopes loads default scope configs for scopes which
// doesn't have explicit user specified configs.
func (c *Config) LoadDefaultScopes(dataDir string) ***REMOVED***
	for k, v := range datastore.DefaultScopes(dataDir) ***REMOVED***
		if _, ok := c.Scopes[k]; !ok ***REMOVED***
			c.Scopes[k] = v
		***REMOVED***
	***REMOVED***
***REMOVED***

// ParseConfig parses the libnetwork configuration file
func ParseConfig(tomlCfgFile string) (*Config, error) ***REMOVED***
	cfg := &Config***REMOVED***
		Scopes: map[string]*datastore.ScopeCfg***REMOVED******REMOVED***,
	***REMOVED***

	if _, err := toml.DecodeFile(tomlCfgFile, cfg); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	cfg.LoadDefaultScopes(cfg.Daemon.DataDir)
	return cfg, nil
***REMOVED***

// ParseConfigOptions parses the configuration options and returns
// a reference to the corresponding Config structure
func ParseConfigOptions(cfgOptions ...Option) *Config ***REMOVED***
	cfg := &Config***REMOVED***
		Daemon: DaemonCfg***REMOVED***
			DriverCfg: make(map[string]interface***REMOVED******REMOVED***),
		***REMOVED***,
		Scopes: make(map[string]*datastore.ScopeCfg),
	***REMOVED***

	cfg.ProcessOptions(cfgOptions...)
	cfg.LoadDefaultScopes(cfg.Daemon.DataDir)

	return cfg
***REMOVED***

// Option is an option setter function type used to pass various configurations
// to the controller
type Option func(c *Config)

// OptionDefaultNetwork function returns an option setter for a default network
func OptionDefaultNetwork(dn string) Option ***REMOVED***
	return func(c *Config) ***REMOVED***
		logrus.Debugf("Option DefaultNetwork: %s", dn)
		c.Daemon.DefaultNetwork = strings.TrimSpace(dn)
	***REMOVED***
***REMOVED***

// OptionDefaultDriver function returns an option setter for default driver
func OptionDefaultDriver(dd string) Option ***REMOVED***
	return func(c *Config) ***REMOVED***
		logrus.Debugf("Option DefaultDriver: %s", dd)
		c.Daemon.DefaultDriver = strings.TrimSpace(dd)
	***REMOVED***
***REMOVED***

// OptionDriverConfig returns an option setter for driver configuration.
func OptionDriverConfig(networkType string, config map[string]interface***REMOVED******REMOVED***) Option ***REMOVED***
	return func(c *Config) ***REMOVED***
		c.Daemon.DriverCfg[networkType] = config
	***REMOVED***
***REMOVED***

// OptionLabels function returns an option setter for labels
func OptionLabels(labels []string) Option ***REMOVED***
	return func(c *Config) ***REMOVED***
		for _, label := range labels ***REMOVED***
			if strings.HasPrefix(label, netlabel.Prefix) ***REMOVED***
				c.Daemon.Labels = append(c.Daemon.Labels, label)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// OptionKVProvider function returns an option setter for kvstore provider
func OptionKVProvider(provider string) Option ***REMOVED***
	return func(c *Config) ***REMOVED***
		logrus.Debugf("Option OptionKVProvider: %s", provider)
		if _, ok := c.Scopes[datastore.GlobalScope]; !ok ***REMOVED***
			c.Scopes[datastore.GlobalScope] = &datastore.ScopeCfg***REMOVED******REMOVED***
		***REMOVED***
		c.Scopes[datastore.GlobalScope].Client.Provider = strings.TrimSpace(provider)
	***REMOVED***
***REMOVED***

// OptionKVProviderURL function returns an option setter for kvstore url
func OptionKVProviderURL(url string) Option ***REMOVED***
	return func(c *Config) ***REMOVED***
		logrus.Debugf("Option OptionKVProviderURL: %s", url)
		if _, ok := c.Scopes[datastore.GlobalScope]; !ok ***REMOVED***
			c.Scopes[datastore.GlobalScope] = &datastore.ScopeCfg***REMOVED******REMOVED***
		***REMOVED***
		c.Scopes[datastore.GlobalScope].Client.Address = strings.TrimSpace(url)
	***REMOVED***
***REMOVED***

// OptionKVOpts function returns an option setter for kvstore options
func OptionKVOpts(opts map[string]string) Option ***REMOVED***
	return func(c *Config) ***REMOVED***
		if opts["kv.cacertfile"] != "" && opts["kv.certfile"] != "" && opts["kv.keyfile"] != "" ***REMOVED***
			logrus.Info("Option Initializing KV with TLS")
			tlsConfig, err := tlsconfig.Client(tlsconfig.Options***REMOVED***
				CAFile:   opts["kv.cacertfile"],
				CertFile: opts["kv.certfile"],
				KeyFile:  opts["kv.keyfile"],
			***REMOVED***)
			if err != nil ***REMOVED***
				logrus.Errorf("Unable to set up TLS: %s", err)
				return
			***REMOVED***
			if _, ok := c.Scopes[datastore.GlobalScope]; !ok ***REMOVED***
				c.Scopes[datastore.GlobalScope] = &datastore.ScopeCfg***REMOVED******REMOVED***
			***REMOVED***
			if c.Scopes[datastore.GlobalScope].Client.Config == nil ***REMOVED***
				c.Scopes[datastore.GlobalScope].Client.Config = &store.Config***REMOVED***TLS: tlsConfig***REMOVED***
			***REMOVED*** else ***REMOVED***
				c.Scopes[datastore.GlobalScope].Client.Config.TLS = tlsConfig
			***REMOVED***
			// Workaround libkv/etcd bug for https
			c.Scopes[datastore.GlobalScope].Client.Config.ClientTLS = &store.ClientTLSConfig***REMOVED***
				CACertFile: opts["kv.cacertfile"],
				CertFile:   opts["kv.certfile"],
				KeyFile:    opts["kv.keyfile"],
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			logrus.Info("Option Initializing KV without TLS")
		***REMOVED***
	***REMOVED***
***REMOVED***

// OptionDiscoveryWatcher function returns an option setter for discovery watcher
func OptionDiscoveryWatcher(watcher discovery.Watcher) Option ***REMOVED***
	return func(c *Config) ***REMOVED***
		c.Cluster.Watcher = watcher
	***REMOVED***
***REMOVED***

// OptionDiscoveryAddress function returns an option setter for self discovery address
func OptionDiscoveryAddress(address string) Option ***REMOVED***
	return func(c *Config) ***REMOVED***
		c.Cluster.Address = address
	***REMOVED***
***REMOVED***

// OptionDataDir function returns an option setter for data folder
func OptionDataDir(dataDir string) Option ***REMOVED***
	return func(c *Config) ***REMOVED***
		c.Daemon.DataDir = dataDir
	***REMOVED***
***REMOVED***

// OptionExecRoot function returns an option setter for exec root folder
func OptionExecRoot(execRoot string) Option ***REMOVED***
	return func(c *Config) ***REMOVED***
		osl.SetBasePath(execRoot)
	***REMOVED***
***REMOVED***

// OptionPluginGetter returns a plugingetter for remote drivers.
func OptionPluginGetter(pg plugingetter.PluginGetter) Option ***REMOVED***
	return func(c *Config) ***REMOVED***
		c.PluginGetter = pg
	***REMOVED***
***REMOVED***

// OptionExperimental function returns an option setter for experimental daemon
func OptionExperimental(exp bool) Option ***REMOVED***
	return func(c *Config) ***REMOVED***
		logrus.Debugf("Option Experimental: %v", exp)
		c.Daemon.Experimental = exp
	***REMOVED***
***REMOVED***

// OptionNetworkControlPlaneMTU function returns an option setter for control plane MTU
func OptionNetworkControlPlaneMTU(exp int) Option ***REMOVED***
	return func(c *Config) ***REMOVED***
		logrus.Debugf("Network Control Plane MTU: %d", exp)
		if exp < warningThNetworkControlPlaneMTU ***REMOVED***
			logrus.Warnf("Received a MTU of %d, this value is very low, the network control plane can misbehave,"+
				" defaulting to minimum value (%d)", exp, minimumNetworkControlPlaneMTU)
			if exp < minimumNetworkControlPlaneMTU ***REMOVED***
				exp = minimumNetworkControlPlaneMTU
			***REMOVED***
		***REMOVED***
		c.Daemon.NetworkControlPlaneMTU = exp
	***REMOVED***
***REMOVED***

// ProcessOptions processes options and stores it in config
func (c *Config) ProcessOptions(options ...Option) ***REMOVED***
	for _, opt := range options ***REMOVED***
		if opt != nil ***REMOVED***
			opt(c)
		***REMOVED***
	***REMOVED***
***REMOVED***

// IsValidName validates configuration objects supported by libnetwork
func IsValidName(name string) bool ***REMOVED***
	return strings.TrimSpace(name) != ""
***REMOVED***

// OptionLocalKVProvider function returns an option setter for kvstore provider
func OptionLocalKVProvider(provider string) Option ***REMOVED***
	return func(c *Config) ***REMOVED***
		logrus.Debugf("Option OptionLocalKVProvider: %s", provider)
		if _, ok := c.Scopes[datastore.LocalScope]; !ok ***REMOVED***
			c.Scopes[datastore.LocalScope] = &datastore.ScopeCfg***REMOVED******REMOVED***
		***REMOVED***
		c.Scopes[datastore.LocalScope].Client.Provider = strings.TrimSpace(provider)
	***REMOVED***
***REMOVED***

// OptionLocalKVProviderURL function returns an option setter for kvstore url
func OptionLocalKVProviderURL(url string) Option ***REMOVED***
	return func(c *Config) ***REMOVED***
		logrus.Debugf("Option OptionLocalKVProviderURL: %s", url)
		if _, ok := c.Scopes[datastore.LocalScope]; !ok ***REMOVED***
			c.Scopes[datastore.LocalScope] = &datastore.ScopeCfg***REMOVED******REMOVED***
		***REMOVED***
		c.Scopes[datastore.LocalScope].Client.Address = strings.TrimSpace(url)
	***REMOVED***
***REMOVED***

// OptionLocalKVProviderConfig function returns an option setter for kvstore config
func OptionLocalKVProviderConfig(config *store.Config) Option ***REMOVED***
	return func(c *Config) ***REMOVED***
		logrus.Debugf("Option OptionLocalKVProviderConfig: %v", config)
		if _, ok := c.Scopes[datastore.LocalScope]; !ok ***REMOVED***
			c.Scopes[datastore.LocalScope] = &datastore.ScopeCfg***REMOVED******REMOVED***
		***REMOVED***
		c.Scopes[datastore.LocalScope].Client.Config = config
	***REMOVED***
***REMOVED***

// OptionActiveSandboxes function returns an option setter for passing the sandboxes
// which were active during previous daemon life
func OptionActiveSandboxes(sandboxes map[string]interface***REMOVED******REMOVED***) Option ***REMOVED***
	return func(c *Config) ***REMOVED***
		c.ActiveSandboxes = sandboxes
	***REMOVED***
***REMOVED***
