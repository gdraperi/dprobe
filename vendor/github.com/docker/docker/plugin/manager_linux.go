package plugin

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/daemon/initlayer"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/containerfs"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/pkg/plugins"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/plugin/v2"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

func (pm *Manager) enable(p *v2.Plugin, c *controller, force bool) (err error) ***REMOVED***
	p.Rootfs = filepath.Join(pm.config.Root, p.PluginObj.ID, "rootfs")
	if p.IsEnabled() && !force ***REMOVED***
		return errors.Wrap(enabledError(p.Name()), "plugin already enabled")
	***REMOVED***
	spec, err := p.InitSpec(pm.config.ExecRoot)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c.restart = true
	c.exitChan = make(chan bool)

	pm.mu.Lock()
	pm.cMap[p] = c
	pm.mu.Unlock()

	var propRoot string
	if p.PropagatedMount != "" ***REMOVED***
		propRoot = filepath.Join(filepath.Dir(p.Rootfs), "propagated-mount")

		if err = os.MkdirAll(propRoot, 0755); err != nil ***REMOVED***
			logrus.Errorf("failed to create PropagatedMount directory at %s: %v", propRoot, err)
		***REMOVED***

		if err = mount.MakeRShared(propRoot); err != nil ***REMOVED***
			return errors.Wrap(err, "error setting up propagated mount dir")
		***REMOVED***

		if err = mount.Mount(propRoot, p.PropagatedMount, "none", "rbind"); err != nil ***REMOVED***
			return errors.Wrap(err, "error creating mount for propagated mount")
		***REMOVED***
	***REMOVED***

	rootFS := containerfs.NewLocalContainerFS(filepath.Join(pm.config.Root, p.PluginObj.ID, rootFSFileName))
	if err := initlayer.Setup(rootFS, idtools.IDPair***REMOVED***0, 0***REMOVED***); err != nil ***REMOVED***
		return errors.WithStack(err)
	***REMOVED***

	stdout, stderr := makeLoggerStreams(p.GetID())
	if err := pm.executor.Create(p.GetID(), *spec, stdout, stderr); err != nil ***REMOVED***
		if p.PropagatedMount != "" ***REMOVED***
			if err := mount.Unmount(p.PropagatedMount); err != nil ***REMOVED***
				logrus.Warnf("Could not unmount %s: %v", p.PropagatedMount, err)
			***REMOVED***
			if err := mount.Unmount(propRoot); err != nil ***REMOVED***
				logrus.Warnf("Could not unmount %s: %v", propRoot, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return pm.pluginPostStart(p, c)
***REMOVED***

func (pm *Manager) pluginPostStart(p *v2.Plugin, c *controller) error ***REMOVED***
	sockAddr := filepath.Join(pm.config.ExecRoot, p.GetID(), p.GetSocket())
	client, err := plugins.NewClientWithTimeout("unix://"+sockAddr, nil, time.Duration(c.timeoutInSecs)*time.Second)
	if err != nil ***REMOVED***
		c.restart = false
		shutdownPlugin(p, c, pm.executor)
		return errors.WithStack(err)
	***REMOVED***

	p.SetPClient(client)

	// Initial sleep before net Dial to allow plugin to listen on socket.
	time.Sleep(500 * time.Millisecond)
	maxRetries := 3
	var retries int
	for ***REMOVED***
		// net dial into the unix socket to see if someone's listening.
		conn, err := net.Dial("unix", sockAddr)
		if err == nil ***REMOVED***
			conn.Close()
			break
		***REMOVED***

		time.Sleep(3 * time.Second)
		retries++

		if retries > maxRetries ***REMOVED***
			logrus.Debugf("error net dialing plugin: %v", err)
			c.restart = false
			// While restoring plugins, we need to explicitly set the state to disabled
			pm.config.Store.SetState(p, false)
			shutdownPlugin(p, c, pm.executor)
			return err
		***REMOVED***

	***REMOVED***
	pm.config.Store.SetState(p, true)
	pm.config.Store.CallHandler(p)

	return pm.save(p)
***REMOVED***

func (pm *Manager) restore(p *v2.Plugin) error ***REMOVED***
	stdout, stderr := makeLoggerStreams(p.GetID())
	if err := pm.executor.Restore(p.GetID(), stdout, stderr); err != nil ***REMOVED***
		return err
	***REMOVED***

	if pm.config.LiveRestoreEnabled ***REMOVED***
		c := &controller***REMOVED******REMOVED***
		if isRunning, _ := pm.executor.IsRunning(p.GetID()); !isRunning ***REMOVED***
			// plugin is not running, so follow normal startup procedure
			return pm.enable(p, c, true)
		***REMOVED***

		c.exitChan = make(chan bool)
		c.restart = true
		pm.mu.Lock()
		pm.cMap[p] = c
		pm.mu.Unlock()
		return pm.pluginPostStart(p, c)
	***REMOVED***

	return nil
***REMOVED***

func shutdownPlugin(p *v2.Plugin, c *controller, executor Executor) ***REMOVED***
	pluginID := p.GetID()

	err := executor.Signal(pluginID, int(unix.SIGTERM))
	if err != nil ***REMOVED***
		logrus.Errorf("Sending SIGTERM to plugin failed with error: %v", err)
	***REMOVED*** else ***REMOVED***
		select ***REMOVED***
		case <-c.exitChan:
			logrus.Debug("Clean shutdown of plugin")
		case <-time.After(time.Second * 10):
			logrus.Debug("Force shutdown plugin")
			if err := executor.Signal(pluginID, int(unix.SIGKILL)); err != nil ***REMOVED***
				logrus.Errorf("Sending SIGKILL to plugin failed with error: %v", err)
			***REMOVED***
			select ***REMOVED***
			case <-c.exitChan:
				logrus.Debug("SIGKILL plugin shutdown")
			case <-time.After(time.Second * 10):
				logrus.Debug("Force shutdown plugin FAILED")
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func setupRoot(root string) error ***REMOVED***
	if err := mount.MakePrivate(root); err != nil ***REMOVED***
		return errors.Wrap(err, "error setting plugin manager root to private")
	***REMOVED***
	return nil
***REMOVED***

func (pm *Manager) disable(p *v2.Plugin, c *controller) error ***REMOVED***
	if !p.IsEnabled() ***REMOVED***
		return errors.Wrap(errDisabled(p.Name()), "plugin is already disabled")
	***REMOVED***

	c.restart = false
	shutdownPlugin(p, c, pm.executor)
	pm.config.Store.SetState(p, false)
	return pm.save(p)
***REMOVED***

// Shutdown stops all plugins and called during daemon shutdown.
func (pm *Manager) Shutdown() ***REMOVED***
	plugins := pm.config.Store.GetAll()
	for _, p := range plugins ***REMOVED***
		pm.mu.RLock()
		c := pm.cMap[p]
		pm.mu.RUnlock()

		if pm.config.LiveRestoreEnabled && p.IsEnabled() ***REMOVED***
			logrus.Debug("Plugin active when liveRestore is set, skipping shutdown")
			continue
		***REMOVED***
		if pm.executor != nil && p.IsEnabled() ***REMOVED***
			c.restart = false
			shutdownPlugin(p, c, pm.executor)
		***REMOVED***
	***REMOVED***
	mount.Unmount(pm.config.Root)
***REMOVED***

func (pm *Manager) upgradePlugin(p *v2.Plugin, configDigest digest.Digest, blobsums []digest.Digest, tmpRootFSDir string, privileges *types.PluginPrivileges) (err error) ***REMOVED***
	config, err := pm.setupNewPlugin(configDigest, blobsums, privileges)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	pdir := filepath.Join(pm.config.Root, p.PluginObj.ID)
	orig := filepath.Join(pdir, "rootfs")

	// Make sure nothing is mounted
	// This could happen if the plugin was disabled with `-f` with active mounts.
	// If there is anything in `orig` is still mounted, this should error out.
	if err := mount.RecursiveUnmount(orig); err != nil ***REMOVED***
		return errdefs.System(err)
	***REMOVED***

	backup := orig + "-old"
	if err := os.Rename(orig, backup); err != nil ***REMOVED***
		return errors.Wrap(errdefs.System(err), "error backing up plugin data before upgrade")
	***REMOVED***

	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if rmErr := os.RemoveAll(orig); rmErr != nil && !os.IsNotExist(rmErr) ***REMOVED***
				logrus.WithError(rmErr).WithField("dir", backup).Error("error cleaning up after failed upgrade")
				return
			***REMOVED***
			if mvErr := os.Rename(backup, orig); mvErr != nil ***REMOVED***
				err = errors.Wrap(mvErr, "error restoring old plugin root on upgrade failure")
			***REMOVED***
			if rmErr := os.RemoveAll(tmpRootFSDir); rmErr != nil && !os.IsNotExist(rmErr) ***REMOVED***
				logrus.WithError(rmErr).WithField("plugin", p.Name()).Errorf("error cleaning up plugin upgrade dir: %s", tmpRootFSDir)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if rmErr := os.RemoveAll(backup); rmErr != nil && !os.IsNotExist(rmErr) ***REMOVED***
				logrus.WithError(rmErr).WithField("dir", backup).Error("error cleaning up old plugin root after successful upgrade")
			***REMOVED***

			p.Config = configDigest
			p.Blobsums = blobsums
		***REMOVED***
	***REMOVED***()

	if err := os.Rename(tmpRootFSDir, orig); err != nil ***REMOVED***
		return errors.Wrap(errdefs.System(err), "error upgrading")
	***REMOVED***

	p.PluginObj.Config = config
	err = pm.save(p)
	return errors.Wrap(err, "error saving upgraded plugin config")
***REMOVED***

func (pm *Manager) setupNewPlugin(configDigest digest.Digest, blobsums []digest.Digest, privileges *types.PluginPrivileges) (types.PluginConfig, error) ***REMOVED***
	configRC, err := pm.blobStore.Get(configDigest)
	if err != nil ***REMOVED***
		return types.PluginConfig***REMOVED******REMOVED***, err
	***REMOVED***
	defer configRC.Close()

	var config types.PluginConfig
	dec := json.NewDecoder(configRC)
	if err := dec.Decode(&config); err != nil ***REMOVED***
		return types.PluginConfig***REMOVED******REMOVED***, errors.Wrapf(err, "failed to parse config")
	***REMOVED***
	if dec.More() ***REMOVED***
		return types.PluginConfig***REMOVED******REMOVED***, errors.New("invalid config json")
	***REMOVED***

	requiredPrivileges := computePrivileges(config)
	if err != nil ***REMOVED***
		return types.PluginConfig***REMOVED******REMOVED***, err
	***REMOVED***
	if privileges != nil ***REMOVED***
		if err := validatePrivileges(requiredPrivileges, *privileges); err != nil ***REMOVED***
			return types.PluginConfig***REMOVED******REMOVED***, err
		***REMOVED***
	***REMOVED***

	return config, nil
***REMOVED***

// createPlugin creates a new plugin. take lock before calling.
func (pm *Manager) createPlugin(name string, configDigest digest.Digest, blobsums []digest.Digest, rootFSDir string, privileges *types.PluginPrivileges, opts ...CreateOpt) (p *v2.Plugin, err error) ***REMOVED***
	if err := pm.config.Store.validateName(name); err != nil ***REMOVED*** // todo: this check is wrong. remove store
		return nil, errdefs.InvalidParameter(err)
	***REMOVED***

	config, err := pm.setupNewPlugin(configDigest, blobsums, privileges)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	p = &v2.Plugin***REMOVED***
		PluginObj: types.Plugin***REMOVED***
			Name:   name,
			ID:     stringid.GenerateRandomID(),
			Config: config,
		***REMOVED***,
		Config:   configDigest,
		Blobsums: blobsums,
	***REMOVED***
	p.InitEmptySettings()
	for _, o := range opts ***REMOVED***
		o(p)
	***REMOVED***

	pdir := filepath.Join(pm.config.Root, p.PluginObj.ID)
	if err := os.MkdirAll(pdir, 0700); err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "failed to mkdir %v", pdir)
	***REMOVED***

	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			os.RemoveAll(pdir)
		***REMOVED***
	***REMOVED***()

	if err := os.Rename(rootFSDir, filepath.Join(pdir, rootFSFileName)); err != nil ***REMOVED***
		return nil, errors.Wrap(err, "failed to rename rootfs")
	***REMOVED***

	if err := pm.save(p); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	pm.config.Store.Add(p) // todo: remove

	return p, nil
***REMOVED***
