package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/docker/distribution/uuid"
	"github.com/docker/docker/api"
	apiserver "github.com/docker/docker/api/server"
	buildbackend "github.com/docker/docker/api/server/backend/build"
	"github.com/docker/docker/api/server/middleware"
	"github.com/docker/docker/api/server/router"
	"github.com/docker/docker/api/server/router/build"
	checkpointrouter "github.com/docker/docker/api/server/router/checkpoint"
	"github.com/docker/docker/api/server/router/container"
	distributionrouter "github.com/docker/docker/api/server/router/distribution"
	"github.com/docker/docker/api/server/router/image"
	"github.com/docker/docker/api/server/router/network"
	pluginrouter "github.com/docker/docker/api/server/router/plugin"
	sessionrouter "github.com/docker/docker/api/server/router/session"
	swarmrouter "github.com/docker/docker/api/server/router/swarm"
	systemrouter "github.com/docker/docker/api/server/router/system"
	"github.com/docker/docker/api/server/router/volume"
	"github.com/docker/docker/builder/dockerfile"
	"github.com/docker/docker/builder/fscache"
	"github.com/docker/docker/cli/debug"
	"github.com/docker/docker/daemon"
	"github.com/docker/docker/daemon/cluster"
	"github.com/docker/docker/daemon/config"
	"github.com/docker/docker/daemon/listeners"
	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/dockerversion"
	"github.com/docker/docker/libcontainerd"
	dopts "github.com/docker/docker/opts"
	"github.com/docker/docker/pkg/authorization"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/pidfile"
	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/docker/pkg/signal"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/docker/plugin"
	"github.com/docker/docker/registry"
	"github.com/docker/docker/runconfig"
	"github.com/docker/go-connections/tlsconfig"
	swarmapi "github.com/docker/swarmkit/api"
	"github.com/moby/buildkit/session"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// DaemonCli represents the daemon CLI.
type DaemonCli struct ***REMOVED***
	*config.Config
	configFile *string
	flags      *pflag.FlagSet

	api             *apiserver.Server
	d               *daemon.Daemon
	authzMiddleware *authorization.Middleware // authzMiddleware enables to dynamically reload the authorization plugins
***REMOVED***

// NewDaemonCli returns a daemon CLI
func NewDaemonCli() *DaemonCli ***REMOVED***
	return &DaemonCli***REMOVED******REMOVED***
***REMOVED***

func (cli *DaemonCli) start(opts *daemonOptions) (err error) ***REMOVED***
	stopc := make(chan bool)
	defer close(stopc)

	// warn from uuid package when running the daemon
	uuid.Loggerf = logrus.Warnf

	opts.SetDefaultOptions(opts.flags)

	if cli.Config, err = loadDaemonCliConfig(opts); err != nil ***REMOVED***
		return err
	***REMOVED***
	cli.configFile = &opts.configFile
	cli.flags = opts.flags

	if cli.Config.Debug ***REMOVED***
		debug.Enable()
	***REMOVED***

	if cli.Config.Experimental ***REMOVED***
		logrus.Warn("Running experimental build")
	***REMOVED***

	logrus.SetFormatter(&logrus.TextFormatter***REMOVED***
		TimestampFormat: jsonmessage.RFC3339NanoFixed,
		DisableColors:   cli.Config.RawLogs,
		FullTimestamp:   true,
	***REMOVED***)

	system.InitLCOW(cli.Config.Experimental)

	if err := setDefaultUmask(); err != nil ***REMOVED***
		return fmt.Errorf("Failed to set umask: %v", err)
	***REMOVED***

	if len(cli.LogConfig.Config) > 0 ***REMOVED***
		if err := logger.ValidateLogOpts(cli.LogConfig.Type, cli.LogConfig.Config); err != nil ***REMOVED***
			return fmt.Errorf("Failed to set log opts: %v", err)
		***REMOVED***
	***REMOVED***

	// Create the daemon root before we create ANY other files (PID, or migrate keys)
	// to ensure the appropriate ACL is set (particularly relevant on Windows)
	if err := daemon.CreateDaemonRoot(cli.Config); err != nil ***REMOVED***
		return err
	***REMOVED***

	if cli.Pidfile != "" ***REMOVED***
		pf, err := pidfile.New(cli.Pidfile)
		if err != nil ***REMOVED***
			return fmt.Errorf("Error starting daemon: %v", err)
		***REMOVED***
		defer func() ***REMOVED***
			if err := pf.Remove(); err != nil ***REMOVED***
				logrus.Error(err)
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	// TODO: extract to newApiServerConfig()
	serverConfig := &apiserver.Config***REMOVED***
		Logging:     true,
		SocketGroup: cli.Config.SocketGroup,
		Version:     dockerversion.Version,
		CorsHeaders: cli.Config.CorsHeaders,
	***REMOVED***

	if cli.Config.TLS ***REMOVED***
		tlsOptions := tlsconfig.Options***REMOVED***
			CAFile:             cli.Config.CommonTLSOptions.CAFile,
			CertFile:           cli.Config.CommonTLSOptions.CertFile,
			KeyFile:            cli.Config.CommonTLSOptions.KeyFile,
			ExclusiveRootPools: true,
		***REMOVED***

		if cli.Config.TLSVerify ***REMOVED***
			// server requires and verifies client's certificate
			tlsOptions.ClientAuth = tls.RequireAndVerifyClientCert
		***REMOVED***
		tlsConfig, err := tlsconfig.Server(tlsOptions)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		serverConfig.TLSConfig = tlsConfig
	***REMOVED***

	if len(cli.Config.Hosts) == 0 ***REMOVED***
		cli.Config.Hosts = make([]string, 1)
	***REMOVED***

	cli.api = apiserver.New(serverConfig)

	var hosts []string

	for i := 0; i < len(cli.Config.Hosts); i++ ***REMOVED***
		var err error
		if cli.Config.Hosts[i], err = dopts.ParseHost(cli.Config.TLS, cli.Config.Hosts[i]); err != nil ***REMOVED***
			return fmt.Errorf("error parsing -H %s : %v", cli.Config.Hosts[i], err)
		***REMOVED***

		protoAddr := cli.Config.Hosts[i]
		protoAddrParts := strings.SplitN(protoAddr, "://", 2)
		if len(protoAddrParts) != 2 ***REMOVED***
			return fmt.Errorf("bad format %s, expected PROTO://ADDR", protoAddr)
		***REMOVED***

		proto := protoAddrParts[0]
		addr := protoAddrParts[1]

		// It's a bad idea to bind to TCP without tlsverify.
		if proto == "tcp" && (serverConfig.TLSConfig == nil || serverConfig.TLSConfig.ClientAuth != tls.RequireAndVerifyClientCert) ***REMOVED***
			logrus.Warn("[!] DON'T BIND ON ANY IP ADDRESS WITHOUT setting --tlsverify IF YOU DON'T KNOW WHAT YOU'RE DOING [!]")
		***REMOVED***
		ls, err := listeners.Init(proto, addr, serverConfig.SocketGroup, serverConfig.TLSConfig)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		ls = wrapListeners(proto, ls)
		// If we're binding to a TCP port, make sure that a container doesn't try to use it.
		if proto == "tcp" ***REMOVED***
			if err := allocateDaemonPort(addr); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		logrus.Debugf("Listener created for HTTP on %s (%s)", proto, addr)
		hosts = append(hosts, protoAddrParts[1])
		cli.api.Accept(addr, ls...)
	***REMOVED***

	registryService, err := registry.NewService(cli.Config.ServiceOptions)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	rOpts, err := cli.getRemoteOptions()
	if err != nil ***REMOVED***
		return fmt.Errorf("Failed to generate containerd options: %s", err)
	***REMOVED***
	containerdRemote, err := libcontainerd.New(filepath.Join(cli.Config.Root, "containerd"), filepath.Join(cli.Config.ExecRoot, "containerd"), rOpts...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	signal.Trap(func() ***REMOVED***
		cli.stop()
		<-stopc // wait for daemonCli.start() to return
	***REMOVED***, logrus.StandardLogger())

	// Notify that the API is active, but before daemon is set up.
	preNotifySystem()

	pluginStore := plugin.NewStore()

	if err := cli.initMiddlewares(cli.api, serverConfig, pluginStore); err != nil ***REMOVED***
		logrus.Fatalf("Error creating middlewares: %v", err)
	***REMOVED***

	d, err := daemon.NewDaemon(cli.Config, registryService, containerdRemote, pluginStore)
	if err != nil ***REMOVED***
		return fmt.Errorf("Error starting daemon: %v", err)
	***REMOVED***

	d.StoreHosts(hosts)

	// validate after NewDaemon has restored enabled plugins. Dont change order.
	if err := validateAuthzPlugins(cli.Config.AuthorizationPlugins, pluginStore); err != nil ***REMOVED***
		return fmt.Errorf("Error validating authorization plugin: %v", err)
	***REMOVED***

	// TODO: move into startMetricsServer()
	if cli.Config.MetricsAddress != "" ***REMOVED***
		if !d.HasExperimental() ***REMOVED***
			return fmt.Errorf("metrics-addr is only supported when experimental is enabled")
		***REMOVED***
		if err := startMetricsServer(cli.Config.MetricsAddress); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// TODO: createAndStartCluster()
	name, _ := os.Hostname()

	// Use a buffered channel to pass changes from store watch API to daemon
	// A buffer allows store watch API and daemon processing to not wait for each other
	watchStream := make(chan *swarmapi.WatchMessage, 32)

	c, err := cluster.New(cluster.Config***REMOVED***
		Root:                   cli.Config.Root,
		Name:                   name,
		Backend:                d,
		PluginBackend:          d.PluginManager(),
		NetworkSubnetsProvider: d,
		DefaultAdvertiseAddr:   cli.Config.SwarmDefaultAdvertiseAddr,
		RuntimeRoot:            cli.getSwarmRunRoot(),
		WatchStream:            watchStream,
	***REMOVED***)
	if err != nil ***REMOVED***
		logrus.Fatalf("Error creating cluster component: %v", err)
	***REMOVED***
	d.SetCluster(c)
	err = c.Start()
	if err != nil ***REMOVED***
		logrus.Fatalf("Error starting cluster component: %v", err)
	***REMOVED***

	// Restart all autostart containers which has a swarm endpoint
	// and is not yet running now that we have successfully
	// initialized the cluster.
	d.RestartSwarmContainers()

	logrus.Info("Daemon has completed initialization")

	cli.d = d

	routerOptions, err := newRouterOptions(cli.Config, d)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	routerOptions.api = cli.api
	routerOptions.cluster = c

	initRouter(routerOptions)

	// process cluster change notifications
	watchCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go d.ProcessClusterNotifications(watchCtx, watchStream)

	cli.setupConfigReloadTrap()

	// The serve API routine never exits unless an error occurs
	// We need to start it as a goroutine and wait on it so
	// daemon doesn't exit
	serveAPIWait := make(chan error)
	go cli.api.Wait(serveAPIWait)

	// after the daemon is done setting up we can notify systemd api
	notifySystem()

	// Daemon is fully initialized and handling API traffic
	// Wait for serve API to complete
	errAPI := <-serveAPIWait
	c.Cleanup()
	shutdownDaemon(d)
	containerdRemote.Cleanup()
	if errAPI != nil ***REMOVED***
		return fmt.Errorf("Shutting down due to ServeAPI error: %v", errAPI)
	***REMOVED***

	return nil
***REMOVED***

type routerOptions struct ***REMOVED***
	sessionManager *session.Manager
	buildBackend   *buildbackend.Backend
	buildCache     *fscache.FSCache
	daemon         *daemon.Daemon
	api            *apiserver.Server
	cluster        *cluster.Cluster
***REMOVED***

func newRouterOptions(config *config.Config, daemon *daemon.Daemon) (routerOptions, error) ***REMOVED***
	opts := routerOptions***REMOVED******REMOVED***
	sm, err := session.NewManager()
	if err != nil ***REMOVED***
		return opts, errors.Wrap(err, "failed to create sessionmanager")
	***REMOVED***

	builderStateDir := filepath.Join(config.Root, "builder")

	buildCache, err := fscache.NewFSCache(fscache.Opt***REMOVED***
		Backend: fscache.NewNaiveCacheBackend(builderStateDir),
		Root:    builderStateDir,
		GCPolicy: fscache.GCPolicy***REMOVED*** // TODO: expose this in config
			MaxSize:         1024 * 1024 * 512,  // 512MB
			MaxKeepDuration: 7 * 24 * time.Hour, // 1 week
		***REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		return opts, errors.Wrap(err, "failed to create fscache")
	***REMOVED***

	manager, err := dockerfile.NewBuildManager(daemon, sm, buildCache, daemon.IDMappings())
	if err != nil ***REMOVED***
		return opts, err
	***REMOVED***

	bb, err := buildbackend.NewBackend(daemon, manager, buildCache)
	if err != nil ***REMOVED***
		return opts, errors.Wrap(err, "failed to create buildmanager")
	***REMOVED***

	return routerOptions***REMOVED***
		sessionManager: sm,
		buildBackend:   bb,
		buildCache:     buildCache,
		daemon:         daemon,
	***REMOVED***, nil
***REMOVED***

func (cli *DaemonCli) reloadConfig() ***REMOVED***
	reload := func(config *config.Config) ***REMOVED***

		// Revalidate and reload the authorization plugins
		if err := validateAuthzPlugins(config.AuthorizationPlugins, cli.d.PluginStore); err != nil ***REMOVED***
			logrus.Fatalf("Error validating authorization plugin: %v", err)
			return
		***REMOVED***
		cli.authzMiddleware.SetPlugins(config.AuthorizationPlugins)

		if err := cli.d.Reload(config); err != nil ***REMOVED***
			logrus.Errorf("Error reconfiguring the daemon: %v", err)
			return
		***REMOVED***

		if config.IsValueSet("debug") ***REMOVED***
			debugEnabled := debug.IsEnabled()
			switch ***REMOVED***
			case debugEnabled && !config.Debug: // disable debug
				debug.Disable()
			case config.Debug && !debugEnabled: // enable debug
				debug.Enable()
			***REMOVED***

		***REMOVED***
	***REMOVED***

	if err := config.Reload(*cli.configFile, cli.flags, reload); err != nil ***REMOVED***
		logrus.Error(err)
	***REMOVED***
***REMOVED***

func (cli *DaemonCli) stop() ***REMOVED***
	cli.api.Close()
***REMOVED***

// shutdownDaemon just wraps daemon.Shutdown() to handle a timeout in case
// d.Shutdown() is waiting too long to kill container or worst it's
// blocked there
func shutdownDaemon(d *daemon.Daemon) ***REMOVED***
	shutdownTimeout := d.ShutdownTimeout()
	ch := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		d.Shutdown()
		close(ch)
	***REMOVED***()
	if shutdownTimeout < 0 ***REMOVED***
		<-ch
		logrus.Debug("Clean shutdown succeeded")
		return
	***REMOVED***
	select ***REMOVED***
	case <-ch:
		logrus.Debug("Clean shutdown succeeded")
	case <-time.After(time.Duration(shutdownTimeout) * time.Second):
		logrus.Error("Force shutdown daemon")
	***REMOVED***
***REMOVED***

func loadDaemonCliConfig(opts *daemonOptions) (*config.Config, error) ***REMOVED***
	conf := opts.daemonConfig
	flags := opts.flags
	conf.Debug = opts.Debug
	conf.Hosts = opts.Hosts
	conf.LogLevel = opts.LogLevel
	conf.TLS = opts.TLS
	conf.TLSVerify = opts.TLSVerify
	conf.CommonTLSOptions = config.CommonTLSOptions***REMOVED******REMOVED***

	if opts.TLSOptions != nil ***REMOVED***
		conf.CommonTLSOptions.CAFile = opts.TLSOptions.CAFile
		conf.CommonTLSOptions.CertFile = opts.TLSOptions.CertFile
		conf.CommonTLSOptions.KeyFile = opts.TLSOptions.KeyFile
	***REMOVED***

	if conf.TrustKeyPath == "" ***REMOVED***
		conf.TrustKeyPath = filepath.Join(
			getDaemonConfDir(conf.Root),
			defaultTrustKeyFile)
	***REMOVED***

	if flags.Changed("graph") && flags.Changed("data-root") ***REMOVED***
		return nil, fmt.Errorf(`cannot specify both "--graph" and "--data-root" option`)
	***REMOVED***

	if opts.configFile != "" ***REMOVED***
		c, err := config.MergeDaemonConfigurations(conf, flags, opts.configFile)
		if err != nil ***REMOVED***
			if flags.Changed("config-file") || !os.IsNotExist(err) ***REMOVED***
				return nil, fmt.Errorf("unable to configure the Docker daemon with file %s: %v", opts.configFile, err)
			***REMOVED***
		***REMOVED***
		// the merged configuration can be nil if the config file didn't exist.
		// leave the current configuration as it is if when that happens.
		if c != nil ***REMOVED***
			conf = c
		***REMOVED***
	***REMOVED***

	if err := config.Validate(conf); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if runtime.GOOS != "windows" ***REMOVED***
		if flags.Changed("disable-legacy-registry") ***REMOVED***
			// TODO: Remove this error after 3 release cycles (18.03)
			return nil, errors.New("ERROR: The '--disable-legacy-registry' flag has been removed. Interacting with legacy (v1) registries is no longer supported")
		***REMOVED***
		if !conf.V2Only ***REMOVED***
			// TODO: Remove this error after 3 release cycles (18.03)
			return nil, errors.New("ERROR: The 'disable-legacy-registry' configuration option has been removed. Interacting with legacy (v1) registries is no longer supported")
		***REMOVED***
	***REMOVED***

	if flags.Changed("graph") ***REMOVED***
		logrus.Warnf(`The "-g / --graph" flag is deprecated. Please use "--data-root" instead`)
	***REMOVED***

	// Check if duplicate label-keys with different values are found
	newLabels, err := config.GetConflictFreeLabels(conf.Labels)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	conf.Labels = newLabels

	// Regardless of whether the user sets it to true or false, if they
	// specify TLSVerify at all then we need to turn on TLS
	if conf.IsValueSet(FlagTLSVerify) ***REMOVED***
		conf.TLS = true
	***REMOVED***

	// ensure that the log level is the one set after merging configurations
	setLogLevel(conf.LogLevel)

	return conf, nil
***REMOVED***

func initRouter(opts routerOptions) ***REMOVED***
	decoder := runconfig.ContainerDecoder***REMOVED******REMOVED***

	routers := []router.Router***REMOVED***
		// we need to add the checkpoint router before the container router or the DELETE gets masked
		checkpointrouter.NewRouter(opts.daemon, decoder),
		container.NewRouter(opts.daemon, decoder),
		image.NewRouter(opts.daemon, decoder),
		systemrouter.NewRouter(opts.daemon, opts.cluster, opts.buildCache),
		volume.NewRouter(opts.daemon),
		build.NewRouter(opts.buildBackend, opts.daemon),
		sessionrouter.NewRouter(opts.sessionManager),
		swarmrouter.NewRouter(opts.cluster),
		pluginrouter.NewRouter(opts.daemon.PluginManager()),
		distributionrouter.NewRouter(opts.daemon),
	***REMOVED***

	if opts.daemon.NetworkControllerEnabled() ***REMOVED***
		routers = append(routers, network.NewRouter(opts.daemon, opts.cluster))
	***REMOVED***

	if opts.daemon.HasExperimental() ***REMOVED***
		for _, r := range routers ***REMOVED***
			for _, route := range r.Routes() ***REMOVED***
				if experimental, ok := route.(router.ExperimentalRoute); ok ***REMOVED***
					experimental.Enable()
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	opts.api.InitRouter(routers...)
***REMOVED***

// TODO: remove this from cli and return the authzMiddleware
func (cli *DaemonCli) initMiddlewares(s *apiserver.Server, cfg *apiserver.Config, pluginStore plugingetter.PluginGetter) error ***REMOVED***
	v := cfg.Version

	exp := middleware.NewExperimentalMiddleware(cli.Config.Experimental)
	s.UseMiddleware(exp)

	vm := middleware.NewVersionMiddleware(v, api.DefaultVersion, api.MinVersion)
	s.UseMiddleware(vm)

	if cfg.CorsHeaders != "" ***REMOVED***
		c := middleware.NewCORSMiddleware(cfg.CorsHeaders)
		s.UseMiddleware(c)
	***REMOVED***

	cli.authzMiddleware = authorization.NewMiddleware(cli.Config.AuthorizationPlugins, pluginStore)
	cli.Config.AuthzMiddleware = cli.authzMiddleware
	s.UseMiddleware(cli.authzMiddleware)
	return nil
***REMOVED***

func (cli *DaemonCli) getRemoteOptions() ([]libcontainerd.RemoteOption, error) ***REMOVED***
	opts := []libcontainerd.RemoteOption***REMOVED******REMOVED***

	pOpts, err := cli.getPlatformRemoteOptions()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	opts = append(opts, pOpts...)
	return opts, nil
***REMOVED***

// validates that the plugins requested with the --authorization-plugin flag are valid AuthzDriver
// plugins present on the host and available to the daemon
func validateAuthzPlugins(requestedPlugins []string, pg plugingetter.PluginGetter) error ***REMOVED***
	for _, reqPlugin := range requestedPlugins ***REMOVED***
		if _, err := pg.Get(reqPlugin, authorization.AuthZApiImplements, plugingetter.Lookup); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
