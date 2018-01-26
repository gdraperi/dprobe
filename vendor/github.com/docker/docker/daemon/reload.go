package daemon

import (
	"encoding/json"
	"fmt"

	"github.com/docker/docker/daemon/config"
	"github.com/docker/docker/daemon/discovery"
	"github.com/sirupsen/logrus"
)

// Reload reads configuration changes and modifies the
// daemon according to those changes.
// These are the settings that Reload changes:
// - Platform runtime
// - Daemon debug log level
// - Daemon max concurrent downloads
// - Daemon max concurrent uploads
// - Daemon shutdown timeout (in seconds)
// - Cluster discovery (reconfigure and restart)
// - Daemon labels
// - Insecure registries
// - Registry mirrors
// - Daemon live restore
func (daemon *Daemon) Reload(conf *config.Config) (err error) ***REMOVED***
	daemon.configStore.Lock()
	attributes := map[string]string***REMOVED******REMOVED***

	defer func() ***REMOVED***
		jsonString, _ := json.Marshal(daemon.configStore)

		// we're unlocking here, because
		// LogDaemonEventWithAttributes() -> SystemInfo() -> GetAllRuntimes()
		// holds that lock too.
		daemon.configStore.Unlock()
		if err == nil ***REMOVED***
			logrus.Infof("Reloaded configuration: %s", jsonString)
			daemon.LogDaemonEventWithAttributes("reload", attributes)
		***REMOVED***
	***REMOVED***()

	if err := daemon.reloadPlatform(conf, attributes); err != nil ***REMOVED***
		return err
	***REMOVED***
	daemon.reloadDebug(conf, attributes)
	daemon.reloadMaxConcurrentDownloadsAndUploads(conf, attributes)
	daemon.reloadShutdownTimeout(conf, attributes)

	if err := daemon.reloadClusterDiscovery(conf, attributes); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := daemon.reloadLabels(conf, attributes); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := daemon.reloadAllowNondistributableArtifacts(conf, attributes); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := daemon.reloadInsecureRegistries(conf, attributes); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := daemon.reloadRegistryMirrors(conf, attributes); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := daemon.reloadLiveRestore(conf, attributes); err != nil ***REMOVED***
		return err
	***REMOVED***
	return daemon.reloadNetworkDiagnosticPort(conf, attributes)
***REMOVED***

// reloadDebug updates configuration with Debug option
// and updates the passed attributes
func (daemon *Daemon) reloadDebug(conf *config.Config, attributes map[string]string) ***REMOVED***
	// update corresponding configuration
	if conf.IsValueSet("debug") ***REMOVED***
		daemon.configStore.Debug = conf.Debug
	***REMOVED***
	// prepare reload event attributes with updatable configurations
	attributes["debug"] = fmt.Sprintf("%t", daemon.configStore.Debug)
***REMOVED***

// reloadMaxConcurrentDownloadsAndUploads updates configuration with max concurrent
// download and upload options and updates the passed attributes
func (daemon *Daemon) reloadMaxConcurrentDownloadsAndUploads(conf *config.Config, attributes map[string]string) ***REMOVED***
	// If no value is set for max-concurrent-downloads we assume it is the default value
	// We always "reset" as the cost is lightweight and easy to maintain.
	if conf.IsValueSet("max-concurrent-downloads") && conf.MaxConcurrentDownloads != nil ***REMOVED***
		*daemon.configStore.MaxConcurrentDownloads = *conf.MaxConcurrentDownloads
	***REMOVED*** else ***REMOVED***
		maxConcurrentDownloads := config.DefaultMaxConcurrentDownloads
		daemon.configStore.MaxConcurrentDownloads = &maxConcurrentDownloads
	***REMOVED***
	logrus.Debugf("Reset Max Concurrent Downloads: %d", *daemon.configStore.MaxConcurrentDownloads)
	if daemon.downloadManager != nil ***REMOVED***
		daemon.downloadManager.SetConcurrency(*daemon.configStore.MaxConcurrentDownloads)
	***REMOVED***

	// prepare reload event attributes with updatable configurations
	attributes["max-concurrent-downloads"] = fmt.Sprintf("%d", *daemon.configStore.MaxConcurrentDownloads)

	// If no value is set for max-concurrent-upload we assume it is the default value
	// We always "reset" as the cost is lightweight and easy to maintain.
	if conf.IsValueSet("max-concurrent-uploads") && conf.MaxConcurrentUploads != nil ***REMOVED***
		*daemon.configStore.MaxConcurrentUploads = *conf.MaxConcurrentUploads
	***REMOVED*** else ***REMOVED***
		maxConcurrentUploads := config.DefaultMaxConcurrentUploads
		daemon.configStore.MaxConcurrentUploads = &maxConcurrentUploads
	***REMOVED***
	logrus.Debugf("Reset Max Concurrent Uploads: %d", *daemon.configStore.MaxConcurrentUploads)
	if daemon.uploadManager != nil ***REMOVED***
		daemon.uploadManager.SetConcurrency(*daemon.configStore.MaxConcurrentUploads)
	***REMOVED***

	// prepare reload event attributes with updatable configurations
	attributes["max-concurrent-uploads"] = fmt.Sprintf("%d", *daemon.configStore.MaxConcurrentUploads)
***REMOVED***

// reloadShutdownTimeout updates configuration with daemon shutdown timeout option
// and updates the passed attributes
func (daemon *Daemon) reloadShutdownTimeout(conf *config.Config, attributes map[string]string) ***REMOVED***
	// update corresponding configuration
	if conf.IsValueSet("shutdown-timeout") ***REMOVED***
		daemon.configStore.ShutdownTimeout = conf.ShutdownTimeout
		logrus.Debugf("Reset Shutdown Timeout: %d", daemon.configStore.ShutdownTimeout)
	***REMOVED***

	// prepare reload event attributes with updatable configurations
	attributes["shutdown-timeout"] = fmt.Sprintf("%d", daemon.configStore.ShutdownTimeout)
***REMOVED***

// reloadClusterDiscovery updates configuration with cluster discovery options
// and updates the passed attributes
func (daemon *Daemon) reloadClusterDiscovery(conf *config.Config, attributes map[string]string) (err error) ***REMOVED***
	defer func() ***REMOVED***
		// prepare reload event attributes with updatable configurations
		attributes["cluster-store"] = conf.ClusterStore
		attributes["cluster-advertise"] = conf.ClusterAdvertise

		attributes["cluster-store-opts"] = "***REMOVED******REMOVED***"
		if daemon.configStore.ClusterOpts != nil ***REMOVED***
			opts, err2 := json.Marshal(conf.ClusterOpts)
			if err != nil ***REMOVED***
				err = err2
			***REMOVED***
			attributes["cluster-store-opts"] = string(opts)
		***REMOVED***
	***REMOVED***()

	newAdvertise := conf.ClusterAdvertise
	newClusterStore := daemon.configStore.ClusterStore
	if conf.IsValueSet("cluster-advertise") ***REMOVED***
		if conf.IsValueSet("cluster-store") ***REMOVED***
			newClusterStore = conf.ClusterStore
		***REMOVED***
		newAdvertise, err = config.ParseClusterAdvertiseSettings(newClusterStore, conf.ClusterAdvertise)
		if err != nil && err != discovery.ErrDiscoveryDisabled ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if daemon.clusterProvider != nil ***REMOVED***
		if err := conf.IsSwarmCompatible(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// check discovery modifications
	if !config.ModifiedDiscoverySettings(daemon.configStore, newClusterStore, newAdvertise, conf.ClusterOpts) ***REMOVED***
		return nil
	***REMOVED***

	// enable discovery for the first time if it was not previously enabled
	if daemon.discoveryWatcher == nil ***REMOVED***
		discoveryWatcher, err := discovery.Init(newClusterStore, newAdvertise, conf.ClusterOpts)
		if err != nil ***REMOVED***
			return fmt.Errorf("failed to initialize discovery: %v", err)
		***REMOVED***
		daemon.discoveryWatcher = discoveryWatcher
	***REMOVED*** else if err == discovery.ErrDiscoveryDisabled ***REMOVED***
		// disable discovery if it was previously enabled and it's disabled now
		daemon.discoveryWatcher.Stop()
	***REMOVED*** else if err = daemon.discoveryWatcher.Reload(conf.ClusterStore, newAdvertise, conf.ClusterOpts); err != nil ***REMOVED***
		// reload discovery
		return err
	***REMOVED***

	daemon.configStore.ClusterStore = newClusterStore
	daemon.configStore.ClusterOpts = conf.ClusterOpts
	daemon.configStore.ClusterAdvertise = newAdvertise

	if daemon.netController == nil ***REMOVED***
		return nil
	***REMOVED***
	netOptions, err := daemon.networkOptions(daemon.configStore, daemon.PluginStore, nil)
	if err != nil ***REMOVED***
		logrus.WithError(err).Warnf("failed to get options with network controller")
		return nil
	***REMOVED***
	err = daemon.netController.ReloadConfiguration(netOptions...)
	if err != nil ***REMOVED***
		logrus.Warnf("Failed to reload configuration with network controller: %v", err)
	***REMOVED***
	return nil
***REMOVED***

// reloadLabels updates configuration with engine labels
// and updates the passed attributes
func (daemon *Daemon) reloadLabels(conf *config.Config, attributes map[string]string) error ***REMOVED***
	// update corresponding configuration
	if conf.IsValueSet("labels") ***REMOVED***
		daemon.configStore.Labels = conf.Labels
	***REMOVED***

	// prepare reload event attributes with updatable configurations
	if daemon.configStore.Labels != nil ***REMOVED***
		labels, err := json.Marshal(daemon.configStore.Labels)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		attributes["labels"] = string(labels)
	***REMOVED*** else ***REMOVED***
		attributes["labels"] = "[]"
	***REMOVED***

	return nil
***REMOVED***

// reloadAllowNondistributableArtifacts updates the configuration with allow-nondistributable-artifacts options
// and updates the passed attributes.
func (daemon *Daemon) reloadAllowNondistributableArtifacts(conf *config.Config, attributes map[string]string) error ***REMOVED***
	// Update corresponding configuration.
	if conf.IsValueSet("allow-nondistributable-artifacts") ***REMOVED***
		daemon.configStore.AllowNondistributableArtifacts = conf.AllowNondistributableArtifacts
		if err := daemon.RegistryService.LoadAllowNondistributableArtifacts(conf.AllowNondistributableArtifacts); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Prepare reload event attributes with updatable configurations.
	if daemon.configStore.AllowNondistributableArtifacts != nil ***REMOVED***
		v, err := json.Marshal(daemon.configStore.AllowNondistributableArtifacts)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		attributes["allow-nondistributable-artifacts"] = string(v)
	***REMOVED*** else ***REMOVED***
		attributes["allow-nondistributable-artifacts"] = "[]"
	***REMOVED***

	return nil
***REMOVED***

// reloadInsecureRegistries updates configuration with insecure registry option
// and updates the passed attributes
func (daemon *Daemon) reloadInsecureRegistries(conf *config.Config, attributes map[string]string) error ***REMOVED***
	// update corresponding configuration
	if conf.IsValueSet("insecure-registries") ***REMOVED***
		daemon.configStore.InsecureRegistries = conf.InsecureRegistries
		if err := daemon.RegistryService.LoadInsecureRegistries(conf.InsecureRegistries); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// prepare reload event attributes with updatable configurations
	if daemon.configStore.InsecureRegistries != nil ***REMOVED***
		insecureRegistries, err := json.Marshal(daemon.configStore.InsecureRegistries)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		attributes["insecure-registries"] = string(insecureRegistries)
	***REMOVED*** else ***REMOVED***
		attributes["insecure-registries"] = "[]"
	***REMOVED***

	return nil
***REMOVED***

// reloadRegistryMirrors updates configuration with registry mirror options
// and updates the passed attributes
func (daemon *Daemon) reloadRegistryMirrors(conf *config.Config, attributes map[string]string) error ***REMOVED***
	// update corresponding configuration
	if conf.IsValueSet("registry-mirrors") ***REMOVED***
		daemon.configStore.Mirrors = conf.Mirrors
		if err := daemon.RegistryService.LoadMirrors(conf.Mirrors); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// prepare reload event attributes with updatable configurations
	if daemon.configStore.Mirrors != nil ***REMOVED***
		mirrors, err := json.Marshal(daemon.configStore.Mirrors)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		attributes["registry-mirrors"] = string(mirrors)
	***REMOVED*** else ***REMOVED***
		attributes["registry-mirrors"] = "[]"
	***REMOVED***

	return nil
***REMOVED***

// reloadLiveRestore updates configuration with live retore option
// and updates the passed attributes
func (daemon *Daemon) reloadLiveRestore(conf *config.Config, attributes map[string]string) error ***REMOVED***
	// update corresponding configuration
	if conf.IsValueSet("live-restore") ***REMOVED***
		daemon.configStore.LiveRestoreEnabled = conf.LiveRestoreEnabled
	***REMOVED***

	// prepare reload event attributes with updatable configurations
	attributes["live-restore"] = fmt.Sprintf("%t", daemon.configStore.LiveRestoreEnabled)
	return nil
***REMOVED***

// reloadNetworkDiagnosticPort updates the network controller starting the diagnose mode if the config is valid
func (daemon *Daemon) reloadNetworkDiagnosticPort(conf *config.Config, attributes map[string]string) error ***REMOVED***
	if conf == nil || daemon.netController == nil ***REMOVED***
		return nil
	***REMOVED***
	// Enable the network diagnose if the flag is set with a valid port withing the range
	if conf.IsValueSet("network-diagnostic-port") && conf.NetworkDiagnosticPort > 0 && conf.NetworkDiagnosticPort < 65536 ***REMOVED***
		logrus.Warnf("Calling the diagnostic start with %d", conf.NetworkDiagnosticPort)
		daemon.netController.StartDiagnose(conf.NetworkDiagnosticPort)
	***REMOVED*** else ***REMOVED***
		daemon.netController.StopDiagnose()
	***REMOVED***
	return nil
***REMOVED***
