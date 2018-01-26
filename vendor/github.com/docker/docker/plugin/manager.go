package plugin

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/image"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/authorization"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/pkg/pubsub"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/docker/plugin/v2"
	"github.com/docker/docker/registry"
	"github.com/opencontainers/go-digest"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const configFileName = "config.json"
const rootFSFileName = "rootfs"

var validFullID = regexp.MustCompile(`^([a-f0-9]***REMOVED***64***REMOVED***)$`)

// Executor is the interface that the plugin manager uses to interact with for starting/stopping plugins
type Executor interface ***REMOVED***
	Create(id string, spec specs.Spec, stdout, stderr io.WriteCloser) error
	Restore(id string, stdout, stderr io.WriteCloser) error
	IsRunning(id string) (bool, error)
	Signal(id string, signal int) error
***REMOVED***

func (pm *Manager) restorePlugin(p *v2.Plugin) error ***REMOVED***
	if p.IsEnabled() ***REMOVED***
		return pm.restore(p)
	***REMOVED***
	return nil
***REMOVED***

type eventLogger func(id, name, action string)

// ManagerConfig defines configuration needed to start new manager.
type ManagerConfig struct ***REMOVED***
	Store              *Store // remove
	RegistryService    registry.Service
	LiveRestoreEnabled bool // TODO: remove
	LogPluginEvent     eventLogger
	Root               string
	ExecRoot           string
	CreateExecutor     ExecutorCreator
	AuthzMiddleware    *authorization.Middleware
***REMOVED***

// ExecutorCreator is used in the manager config to pass in an `Executor`
type ExecutorCreator func(*Manager) (Executor, error)

// Manager controls the plugin subsystem.
type Manager struct ***REMOVED***
	config    ManagerConfig
	mu        sync.RWMutex // protects cMap
	muGC      sync.RWMutex // protects blobstore deletions
	cMap      map[*v2.Plugin]*controller
	blobStore *basicBlobStore
	publisher *pubsub.Publisher
	executor  Executor
***REMOVED***

// controller represents the manager's control on a plugin.
type controller struct ***REMOVED***
	restart       bool
	exitChan      chan bool
	timeoutInSecs int
***REMOVED***

// pluginRegistryService ensures that all resolved repositories
// are of the plugin class.
type pluginRegistryService struct ***REMOVED***
	registry.Service
***REMOVED***

func (s pluginRegistryService) ResolveRepository(name reference.Named) (repoInfo *registry.RepositoryInfo, err error) ***REMOVED***
	repoInfo, err = s.Service.ResolveRepository(name)
	if repoInfo != nil ***REMOVED***
		repoInfo.Class = "plugin"
	***REMOVED***
	return
***REMOVED***

// NewManager returns a new plugin manager.
func NewManager(config ManagerConfig) (*Manager, error) ***REMOVED***
	if config.RegistryService != nil ***REMOVED***
		config.RegistryService = pluginRegistryService***REMOVED***config.RegistryService***REMOVED***
	***REMOVED***
	manager := &Manager***REMOVED***
		config: config,
	***REMOVED***
	for _, dirName := range []string***REMOVED***manager.config.Root, manager.config.ExecRoot, manager.tmpDir()***REMOVED*** ***REMOVED***
		if err := os.MkdirAll(dirName, 0700); err != nil ***REMOVED***
			return nil, errors.Wrapf(err, "failed to mkdir %v", dirName)
		***REMOVED***
	***REMOVED***

	if err := setupRoot(manager.config.Root); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var err error
	manager.executor, err = config.CreateExecutor(manager)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	manager.blobStore, err = newBasicBlobStore(filepath.Join(manager.config.Root, "storage/blobs"))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	manager.cMap = make(map[*v2.Plugin]*controller)
	if err := manager.reload(); err != nil ***REMOVED***
		return nil, errors.Wrap(err, "failed to restore plugins")
	***REMOVED***

	manager.publisher = pubsub.NewPublisher(0, 0)
	return manager, nil
***REMOVED***

func (pm *Manager) tmpDir() string ***REMOVED***
	return filepath.Join(pm.config.Root, "tmp")
***REMOVED***

// HandleExitEvent is called when the executor receives the exit event
// In the future we may change this, but for now all we care about is the exit event.
func (pm *Manager) HandleExitEvent(id string) error ***REMOVED***
	p, err := pm.config.Store.GetV2Plugin(id)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	os.RemoveAll(filepath.Join(pm.config.ExecRoot, id))

	if p.PropagatedMount != "" ***REMOVED***
		if err := mount.Unmount(p.PropagatedMount); err != nil ***REMOVED***
			logrus.Warnf("Could not unmount %s: %v", p.PropagatedMount, err)
		***REMOVED***
		propRoot := filepath.Join(filepath.Dir(p.Rootfs), "propagated-mount")
		if err := mount.Unmount(propRoot); err != nil ***REMOVED***
			logrus.Warn("Could not unmount %s: %v", propRoot, err)
		***REMOVED***
	***REMOVED***

	pm.mu.RLock()
	c := pm.cMap[p]
	if c.exitChan != nil ***REMOVED***
		close(c.exitChan)
	***REMOVED***
	restart := c.restart
	pm.mu.RUnlock()

	if restart ***REMOVED***
		pm.enable(p, c, true)
	***REMOVED***
	return nil
***REMOVED***

func handleLoadError(err error, id string) ***REMOVED***
	if err == nil ***REMOVED***
		return
	***REMOVED***
	logger := logrus.WithError(err).WithField("id", id)
	if os.IsNotExist(errors.Cause(err)) ***REMOVED***
		// Likely some error while removing on an older version of docker
		logger.Warn("missing plugin config, skipping: this may be caused due to a failed remove and requires manual cleanup.")
		return
	***REMOVED***
	logger.Error("error loading plugin, skipping")
***REMOVED***

func (pm *Manager) reload() error ***REMOVED*** // todo: restore
	dir, err := ioutil.ReadDir(pm.config.Root)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to read %v", pm.config.Root)
	***REMOVED***
	plugins := make(map[string]*v2.Plugin)
	for _, v := range dir ***REMOVED***
		if validFullID.MatchString(v.Name()) ***REMOVED***
			p, err := pm.loadPlugin(v.Name())
			if err != nil ***REMOVED***
				handleLoadError(err, v.Name())
				continue
			***REMOVED***
			plugins[p.GetID()] = p
		***REMOVED*** else ***REMOVED***
			if validFullID.MatchString(strings.TrimSuffix(v.Name(), "-removing")) ***REMOVED***
				// There was likely some error while removing this plugin, let's try to remove again here
				if err := system.EnsureRemoveAll(v.Name()); err != nil ***REMOVED***
					logrus.WithError(err).WithField("id", v.Name()).Warn("error while attempting to clean up previously removed plugin")
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	pm.config.Store.SetAll(plugins)

	var wg sync.WaitGroup
	wg.Add(len(plugins))
	for _, p := range plugins ***REMOVED***
		c := &controller***REMOVED******REMOVED*** // todo: remove this
		pm.cMap[p] = c
		go func(p *v2.Plugin) ***REMOVED***
			defer wg.Done()
			if err := pm.restorePlugin(p); err != nil ***REMOVED***
				logrus.Errorf("failed to restore plugin '%s': %s", p.Name(), err)
				return
			***REMOVED***

			if p.Rootfs != "" ***REMOVED***
				p.Rootfs = filepath.Join(pm.config.Root, p.PluginObj.ID, "rootfs")
			***REMOVED***

			// We should only enable rootfs propagation for certain plugin types that need it.
			for _, typ := range p.PluginObj.Config.Interface.Types ***REMOVED***
				if (typ.Capability == "volumedriver" || typ.Capability == "graphdriver") && typ.Prefix == "docker" && strings.HasPrefix(typ.Version, "1.") ***REMOVED***
					if p.PluginObj.Config.PropagatedMount != "" ***REMOVED***
						propRoot := filepath.Join(filepath.Dir(p.Rootfs), "propagated-mount")

						// check if we need to migrate an older propagated mount from before
						// these mounts were stored outside the plugin rootfs
						if _, err := os.Stat(propRoot); os.IsNotExist(err) ***REMOVED***
							if _, err := os.Stat(p.PropagatedMount); err == nil ***REMOVED***
								// make sure nothing is mounted here
								// don't care about errors
								mount.Unmount(p.PropagatedMount)
								if err := os.Rename(p.PropagatedMount, propRoot); err != nil ***REMOVED***
									logrus.WithError(err).WithField("dir", propRoot).Error("error migrating propagated mount storage")
								***REMOVED***
								if err := os.MkdirAll(p.PropagatedMount, 0755); err != nil ***REMOVED***
									logrus.WithError(err).WithField("dir", p.PropagatedMount).Error("error migrating propagated mount storage")
								***REMOVED***
							***REMOVED***
						***REMOVED***

						if err := os.MkdirAll(propRoot, 0755); err != nil ***REMOVED***
							logrus.Errorf("failed to create PropagatedMount directory at %s: %v", propRoot, err)
						***REMOVED***
						// TODO: sanitize PropagatedMount and prevent breakout
						p.PropagatedMount = filepath.Join(p.Rootfs, p.PluginObj.Config.PropagatedMount)
						if err := os.MkdirAll(p.PropagatedMount, 0755); err != nil ***REMOVED***
							logrus.Errorf("failed to create PropagatedMount directory at %s: %v", p.PropagatedMount, err)
							return
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***

			pm.save(p)
			requiresManualRestore := !pm.config.LiveRestoreEnabled && p.IsEnabled()

			if requiresManualRestore ***REMOVED***
				// if liveRestore is not enabled, the plugin will be stopped now so we should enable it
				if err := pm.enable(p, c, true); err != nil ***REMOVED***
					logrus.Errorf("failed to enable plugin '%s': %s", p.Name(), err)
				***REMOVED***
			***REMOVED***
		***REMOVED***(p)
	***REMOVED***
	wg.Wait()
	return nil
***REMOVED***

// Get looks up the requested plugin in the store.
func (pm *Manager) Get(idOrName string) (*v2.Plugin, error) ***REMOVED***
	return pm.config.Store.GetV2Plugin(idOrName)
***REMOVED***

func (pm *Manager) loadPlugin(id string) (*v2.Plugin, error) ***REMOVED***
	p := filepath.Join(pm.config.Root, id, configFileName)
	dt, err := ioutil.ReadFile(p)
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "error reading %v", p)
	***REMOVED***
	var plugin v2.Plugin
	if err := json.Unmarshal(dt, &plugin); err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "error decoding %v", p)
	***REMOVED***
	return &plugin, nil
***REMOVED***

func (pm *Manager) save(p *v2.Plugin) error ***REMOVED***
	pluginJSON, err := json.Marshal(p)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to marshal plugin json")
	***REMOVED***
	if err := ioutils.AtomicWriteFile(filepath.Join(pm.config.Root, p.GetID(), configFileName), pluginJSON, 0600); err != nil ***REMOVED***
		return errors.Wrap(err, "failed to write atomically plugin json")
	***REMOVED***
	return nil
***REMOVED***

// GC cleans up unreferenced blobs. This is recommended to run in a goroutine
func (pm *Manager) GC() ***REMOVED***
	pm.muGC.Lock()
	defer pm.muGC.Unlock()

	whitelist := make(map[digest.Digest]struct***REMOVED******REMOVED***)
	for _, p := range pm.config.Store.GetAll() ***REMOVED***
		whitelist[p.Config] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		for _, b := range p.Blobsums ***REMOVED***
			whitelist[b] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***

	pm.blobStore.gc(whitelist)
***REMOVED***

type logHook struct***REMOVED*** id string ***REMOVED***

func (logHook) Levels() []logrus.Level ***REMOVED***
	return logrus.AllLevels
***REMOVED***

func (l logHook) Fire(entry *logrus.Entry) error ***REMOVED***
	entry.Data = logrus.Fields***REMOVED***"plugin": l.id***REMOVED***
	return nil
***REMOVED***

func makeLoggerStreams(id string) (stdout, stderr io.WriteCloser) ***REMOVED***
	logger := logrus.New()
	logger.Hooks.Add(logHook***REMOVED***id***REMOVED***)
	return logger.WriterLevel(logrus.InfoLevel), logger.WriterLevel(logrus.ErrorLevel)
***REMOVED***

func validatePrivileges(requiredPrivileges, privileges types.PluginPrivileges) error ***REMOVED***
	if !isEqual(requiredPrivileges, privileges, isEqualPrivilege) ***REMOVED***
		return errors.New("incorrect privileges")
	***REMOVED***

	return nil
***REMOVED***

func isEqual(arrOne, arrOther types.PluginPrivileges, compare func(x, y types.PluginPrivilege) bool) bool ***REMOVED***
	if len(arrOne) != len(arrOther) ***REMOVED***
		return false
	***REMOVED***

	sort.Sort(arrOne)
	sort.Sort(arrOther)

	for i := 1; i < arrOne.Len(); i++ ***REMOVED***
		if !compare(arrOne[i], arrOther[i]) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

func isEqualPrivilege(a, b types.PluginPrivilege) bool ***REMOVED***
	if a.Name != b.Name ***REMOVED***
		return false
	***REMOVED***

	return reflect.DeepEqual(a.Value, b.Value)
***REMOVED***

func configToRootFS(c []byte) (*image.RootFS, string, error) ***REMOVED***
	// TODO @jhowardmsft LCOW - Will need to revisit this.
	os := runtime.GOOS
	var pluginConfig types.PluginConfig
	if err := json.Unmarshal(c, &pluginConfig); err != nil ***REMOVED***
		return nil, "", err
	***REMOVED***
	// validation for empty rootfs is in distribution code
	if pluginConfig.Rootfs == nil ***REMOVED***
		return nil, os, nil
	***REMOVED***

	return rootFSFromPlugin(pluginConfig.Rootfs), os, nil
***REMOVED***

func rootFSFromPlugin(pluginfs *types.PluginConfigRootfs) *image.RootFS ***REMOVED***
	rootFS := image.RootFS***REMOVED***
		Type:    pluginfs.Type,
		DiffIDs: make([]layer.DiffID, len(pluginfs.DiffIds)),
	***REMOVED***
	for i := range pluginfs.DiffIds ***REMOVED***
		rootFS.DiffIDs[i] = layer.DiffID(pluginfs.DiffIds[i])
	***REMOVED***

	return &rootFS
***REMOVED***
