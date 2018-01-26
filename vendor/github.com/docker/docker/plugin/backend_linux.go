package plugin

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/distribution"
	progressutils "github.com/docker/docker/distribution/utils"
	"github.com/docker/docker/distribution/xfer"
	"github.com/docker/docker/dockerversion"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/image"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/authorization"
	"github.com/docker/docker/pkg/chrootarchive"
	"github.com/docker/docker/pkg/mount"
	"github.com/docker/docker/pkg/pools"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/docker/plugin/v2"
	refstore "github.com/docker/docker/reference"
	digest "github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

var acceptedPluginFilterTags = map[string]bool***REMOVED***
	"enabled":    true,
	"capability": true,
***REMOVED***

// Disable deactivates a plugin. This means resources (volumes, networks) cant use them.
func (pm *Manager) Disable(refOrID string, config *types.PluginDisableConfig) error ***REMOVED***
	p, err := pm.config.Store.GetV2Plugin(refOrID)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	pm.mu.RLock()
	c := pm.cMap[p]
	pm.mu.RUnlock()

	if !config.ForceDisable && p.GetRefCount() > 0 ***REMOVED***
		return errors.WithStack(inUseError(p.Name()))
	***REMOVED***

	for _, typ := range p.GetTypes() ***REMOVED***
		if typ.Capability == authorization.AuthZApiImplements ***REMOVED***
			pm.config.AuthzMiddleware.RemovePlugin(p.Name())
		***REMOVED***
	***REMOVED***

	if err := pm.disable(p, c); err != nil ***REMOVED***
		return err
	***REMOVED***
	pm.publisher.Publish(EventDisable***REMOVED***Plugin: p.PluginObj***REMOVED***)
	pm.config.LogPluginEvent(p.GetID(), refOrID, "disable")
	return nil
***REMOVED***

// Enable activates a plugin, which implies that they are ready to be used by containers.
func (pm *Manager) Enable(refOrID string, config *types.PluginEnableConfig) error ***REMOVED***
	p, err := pm.config.Store.GetV2Plugin(refOrID)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	c := &controller***REMOVED***timeoutInSecs: config.Timeout***REMOVED***
	if err := pm.enable(p, c, false); err != nil ***REMOVED***
		return err
	***REMOVED***
	pm.publisher.Publish(EventEnable***REMOVED***Plugin: p.PluginObj***REMOVED***)
	pm.config.LogPluginEvent(p.GetID(), refOrID, "enable")
	return nil
***REMOVED***

// Inspect examines a plugin config
func (pm *Manager) Inspect(refOrID string) (tp *types.Plugin, err error) ***REMOVED***
	p, err := pm.config.Store.GetV2Plugin(refOrID)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &p.PluginObj, nil
***REMOVED***

func (pm *Manager) pull(ctx context.Context, ref reference.Named, config *distribution.ImagePullConfig, outStream io.Writer) error ***REMOVED***
	if outStream != nil ***REMOVED***
		// Include a buffer so that slow client connections don't affect
		// transfer performance.
		progressChan := make(chan progress.Progress, 100)

		writesDone := make(chan struct***REMOVED******REMOVED***)

		defer func() ***REMOVED***
			close(progressChan)
			<-writesDone
		***REMOVED***()

		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithCancel(ctx)

		go func() ***REMOVED***
			progressutils.WriteDistributionProgress(cancelFunc, outStream, progressChan)
			close(writesDone)
		***REMOVED***()

		config.ProgressOutput = progress.ChanOutput(progressChan)
	***REMOVED*** else ***REMOVED***
		config.ProgressOutput = progress.DiscardOutput()
	***REMOVED***
	return distribution.Pull(ctx, ref, config)
***REMOVED***

type tempConfigStore struct ***REMOVED***
	config       []byte
	configDigest digest.Digest
***REMOVED***

func (s *tempConfigStore) Put(c []byte) (digest.Digest, error) ***REMOVED***
	dgst := digest.FromBytes(c)

	s.config = c
	s.configDigest = dgst

	return dgst, nil
***REMOVED***

func (s *tempConfigStore) Get(d digest.Digest) ([]byte, error) ***REMOVED***
	if d != s.configDigest ***REMOVED***
		return nil, errNotFound("digest not found")
	***REMOVED***
	return s.config, nil
***REMOVED***

func (s *tempConfigStore) RootFSAndOSFromConfig(c []byte) (*image.RootFS, string, error) ***REMOVED***
	return configToRootFS(c)
***REMOVED***

func computePrivileges(c types.PluginConfig) types.PluginPrivileges ***REMOVED***
	var privileges types.PluginPrivileges
	if c.Network.Type != "null" && c.Network.Type != "bridge" && c.Network.Type != "" ***REMOVED***
		privileges = append(privileges, types.PluginPrivilege***REMOVED***
			Name:        "network",
			Description: "permissions to access a network",
			Value:       []string***REMOVED***c.Network.Type***REMOVED***,
		***REMOVED***)
	***REMOVED***
	if c.IpcHost ***REMOVED***
		privileges = append(privileges, types.PluginPrivilege***REMOVED***
			Name:        "host ipc namespace",
			Description: "allow access to host ipc namespace",
			Value:       []string***REMOVED***"true"***REMOVED***,
		***REMOVED***)
	***REMOVED***
	if c.PidHost ***REMOVED***
		privileges = append(privileges, types.PluginPrivilege***REMOVED***
			Name:        "host pid namespace",
			Description: "allow access to host pid namespace",
			Value:       []string***REMOVED***"true"***REMOVED***,
		***REMOVED***)
	***REMOVED***
	for _, mount := range c.Mounts ***REMOVED***
		if mount.Source != nil ***REMOVED***
			privileges = append(privileges, types.PluginPrivilege***REMOVED***
				Name:        "mount",
				Description: "host path to mount",
				Value:       []string***REMOVED****mount.Source***REMOVED***,
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	for _, device := range c.Linux.Devices ***REMOVED***
		if device.Path != nil ***REMOVED***
			privileges = append(privileges, types.PluginPrivilege***REMOVED***
				Name:        "device",
				Description: "host device to access",
				Value:       []string***REMOVED****device.Path***REMOVED***,
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	if c.Linux.AllowAllDevices ***REMOVED***
		privileges = append(privileges, types.PluginPrivilege***REMOVED***
			Name:        "allow-all-devices",
			Description: "allow 'rwm' access to all devices",
			Value:       []string***REMOVED***"true"***REMOVED***,
		***REMOVED***)
	***REMOVED***
	if len(c.Linux.Capabilities) > 0 ***REMOVED***
		privileges = append(privileges, types.PluginPrivilege***REMOVED***
			Name:        "capabilities",
			Description: "list of additional capabilities required",
			Value:       c.Linux.Capabilities,
		***REMOVED***)
	***REMOVED***

	return privileges
***REMOVED***

// Privileges pulls a plugin config and computes the privileges required to install it.
func (pm *Manager) Privileges(ctx context.Context, ref reference.Named, metaHeader http.Header, authConfig *types.AuthConfig) (types.PluginPrivileges, error) ***REMOVED***
	// create image store instance
	cs := &tempConfigStore***REMOVED******REMOVED***

	// DownloadManager not defined because only pulling configuration.
	pluginPullConfig := &distribution.ImagePullConfig***REMOVED***
		Config: distribution.Config***REMOVED***
			MetaHeaders:      metaHeader,
			AuthConfig:       authConfig,
			RegistryService:  pm.config.RegistryService,
			ImageEventLogger: func(string, string, string) ***REMOVED******REMOVED***,
			ImageStore:       cs,
		***REMOVED***,
		Schema2Types: distribution.PluginTypes,
	***REMOVED***

	if err := pm.pull(ctx, ref, pluginPullConfig, nil); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if cs.config == nil ***REMOVED***
		return nil, errors.New("no configuration pulled")
	***REMOVED***
	var config types.PluginConfig
	if err := json.Unmarshal(cs.config, &config); err != nil ***REMOVED***
		return nil, errdefs.System(err)
	***REMOVED***

	return computePrivileges(config), nil
***REMOVED***

// Upgrade upgrades a plugin
func (pm *Manager) Upgrade(ctx context.Context, ref reference.Named, name string, metaHeader http.Header, authConfig *types.AuthConfig, privileges types.PluginPrivileges, outStream io.Writer) (err error) ***REMOVED***
	p, err := pm.config.Store.GetV2Plugin(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if p.IsEnabled() ***REMOVED***
		return errors.Wrap(enabledError(p.Name()), "plugin must be disabled before upgrading")
	***REMOVED***

	pm.muGC.RLock()
	defer pm.muGC.RUnlock()

	// revalidate because Pull is public
	if _, err := reference.ParseNormalizedNamed(name); err != nil ***REMOVED***
		return errors.Wrapf(errdefs.InvalidParameter(err), "failed to parse %q", name)
	***REMOVED***

	tmpRootFSDir, err := ioutil.TempDir(pm.tmpDir(), ".rootfs")
	if err != nil ***REMOVED***
		return errors.Wrap(errdefs.System(err), "error preparing upgrade")
	***REMOVED***
	defer os.RemoveAll(tmpRootFSDir)

	dm := &downloadManager***REMOVED***
		tmpDir:    tmpRootFSDir,
		blobStore: pm.blobStore,
	***REMOVED***

	pluginPullConfig := &distribution.ImagePullConfig***REMOVED***
		Config: distribution.Config***REMOVED***
			MetaHeaders:      metaHeader,
			AuthConfig:       authConfig,
			RegistryService:  pm.config.RegistryService,
			ImageEventLogger: pm.config.LogPluginEvent,
			ImageStore:       dm,
		***REMOVED***,
		DownloadManager: dm, // todo: reevaluate if possible to substitute distribution/xfer dependencies instead
		Schema2Types:    distribution.PluginTypes,
	***REMOVED***

	err = pm.pull(ctx, ref, pluginPullConfig, outStream)
	if err != nil ***REMOVED***
		go pm.GC()
		return err
	***REMOVED***

	if err := pm.upgradePlugin(p, dm.configDigest, dm.blobs, tmpRootFSDir, &privileges); err != nil ***REMOVED***
		return err
	***REMOVED***
	p.PluginObj.PluginReference = ref.String()
	return nil
***REMOVED***

// Pull pulls a plugin, check if the correct privileges are provided and install the plugin.
func (pm *Manager) Pull(ctx context.Context, ref reference.Named, name string, metaHeader http.Header, authConfig *types.AuthConfig, privileges types.PluginPrivileges, outStream io.Writer, opts ...CreateOpt) (err error) ***REMOVED***
	pm.muGC.RLock()
	defer pm.muGC.RUnlock()

	// revalidate because Pull is public
	nameref, err := reference.ParseNormalizedNamed(name)
	if err != nil ***REMOVED***
		return errors.Wrapf(errdefs.InvalidParameter(err), "failed to parse %q", name)
	***REMOVED***
	name = reference.FamiliarString(reference.TagNameOnly(nameref))

	if err := pm.config.Store.validateName(name); err != nil ***REMOVED***
		return errdefs.InvalidParameter(err)
	***REMOVED***

	tmpRootFSDir, err := ioutil.TempDir(pm.tmpDir(), ".rootfs")
	if err != nil ***REMOVED***
		return errors.Wrap(errdefs.System(err), "error preparing pull")
	***REMOVED***
	defer os.RemoveAll(tmpRootFSDir)

	dm := &downloadManager***REMOVED***
		tmpDir:    tmpRootFSDir,
		blobStore: pm.blobStore,
	***REMOVED***

	pluginPullConfig := &distribution.ImagePullConfig***REMOVED***
		Config: distribution.Config***REMOVED***
			MetaHeaders:      metaHeader,
			AuthConfig:       authConfig,
			RegistryService:  pm.config.RegistryService,
			ImageEventLogger: pm.config.LogPluginEvent,
			ImageStore:       dm,
		***REMOVED***,
		DownloadManager: dm, // todo: reevaluate if possible to substitute distribution/xfer dependencies instead
		Schema2Types:    distribution.PluginTypes,
	***REMOVED***

	err = pm.pull(ctx, ref, pluginPullConfig, outStream)
	if err != nil ***REMOVED***
		go pm.GC()
		return err
	***REMOVED***

	refOpt := func(p *v2.Plugin) ***REMOVED***
		p.PluginObj.PluginReference = ref.String()
	***REMOVED***
	optsList := make([]CreateOpt, 0, len(opts)+1)
	optsList = append(optsList, opts...)
	optsList = append(optsList, refOpt)

	p, err := pm.createPlugin(name, dm.configDigest, dm.blobs, tmpRootFSDir, &privileges, optsList...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	pm.publisher.Publish(EventCreate***REMOVED***Plugin: p.PluginObj***REMOVED***)
	return nil
***REMOVED***

// List displays the list of plugins and associated metadata.
func (pm *Manager) List(pluginFilters filters.Args) ([]types.Plugin, error) ***REMOVED***
	if err := pluginFilters.Validate(acceptedPluginFilterTags); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	enabledOnly := false
	disabledOnly := false
	if pluginFilters.Contains("enabled") ***REMOVED***
		if pluginFilters.ExactMatch("enabled", "true") ***REMOVED***
			enabledOnly = true
		***REMOVED*** else if pluginFilters.ExactMatch("enabled", "false") ***REMOVED***
			disabledOnly = true
		***REMOVED*** else ***REMOVED***
			return nil, invalidFilter***REMOVED***"enabled", pluginFilters.Get("enabled")***REMOVED***
		***REMOVED***
	***REMOVED***

	plugins := pm.config.Store.GetAll()
	out := make([]types.Plugin, 0, len(plugins))

next:
	for _, p := range plugins ***REMOVED***
		if enabledOnly && !p.PluginObj.Enabled ***REMOVED***
			continue
		***REMOVED***
		if disabledOnly && p.PluginObj.Enabled ***REMOVED***
			continue
		***REMOVED***
		if pluginFilters.Contains("capability") ***REMOVED***
			for _, f := range p.GetTypes() ***REMOVED***
				if !pluginFilters.Match("capability", f.Capability) ***REMOVED***
					continue next
				***REMOVED***
			***REMOVED***
		***REMOVED***
		out = append(out, p.PluginObj)
	***REMOVED***
	return out, nil
***REMOVED***

// Push pushes a plugin to the store.
func (pm *Manager) Push(ctx context.Context, name string, metaHeader http.Header, authConfig *types.AuthConfig, outStream io.Writer) error ***REMOVED***
	p, err := pm.config.Store.GetV2Plugin(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	ref, err := reference.ParseNormalizedNamed(p.Name())
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "plugin has invalid name %v for push", p.Name())
	***REMOVED***

	var po progress.Output
	if outStream != nil ***REMOVED***
		// Include a buffer so that slow client connections don't affect
		// transfer performance.
		progressChan := make(chan progress.Progress, 100)

		writesDone := make(chan struct***REMOVED******REMOVED***)

		defer func() ***REMOVED***
			close(progressChan)
			<-writesDone
		***REMOVED***()

		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithCancel(ctx)

		go func() ***REMOVED***
			progressutils.WriteDistributionProgress(cancelFunc, outStream, progressChan)
			close(writesDone)
		***REMOVED***()

		po = progress.ChanOutput(progressChan)
	***REMOVED*** else ***REMOVED***
		po = progress.DiscardOutput()
	***REMOVED***

	// TODO: replace these with manager
	is := &pluginConfigStore***REMOVED***
		pm:     pm,
		plugin: p,
	***REMOVED***
	lss := make(map[string]distribution.PushLayerProvider)
	lss[runtime.GOOS] = &pluginLayerProvider***REMOVED***
		pm:     pm,
		plugin: p,
	***REMOVED***
	rs := &pluginReference***REMOVED***
		name:     ref,
		pluginID: p.Config,
	***REMOVED***

	uploadManager := xfer.NewLayerUploadManager(3)

	imagePushConfig := &distribution.ImagePushConfig***REMOVED***
		Config: distribution.Config***REMOVED***
			MetaHeaders:      metaHeader,
			AuthConfig:       authConfig,
			ProgressOutput:   po,
			RegistryService:  pm.config.RegistryService,
			ReferenceStore:   rs,
			ImageEventLogger: pm.config.LogPluginEvent,
			ImageStore:       is,
			RequireSchema2:   true,
		***REMOVED***,
		ConfigMediaType: schema2.MediaTypePluginConfig,
		LayerStores:     lss,
		UploadManager:   uploadManager,
	***REMOVED***

	return distribution.Push(ctx, ref, imagePushConfig)
***REMOVED***

type pluginReference struct ***REMOVED***
	name     reference.Named
	pluginID digest.Digest
***REMOVED***

func (r *pluginReference) References(id digest.Digest) []reference.Named ***REMOVED***
	if r.pluginID != id ***REMOVED***
		return nil
	***REMOVED***
	return []reference.Named***REMOVED***r.name***REMOVED***
***REMOVED***

func (r *pluginReference) ReferencesByName(ref reference.Named) []refstore.Association ***REMOVED***
	return []refstore.Association***REMOVED***
		***REMOVED***
			Ref: r.name,
			ID:  r.pluginID,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (r *pluginReference) Get(ref reference.Named) (digest.Digest, error) ***REMOVED***
	if r.name.String() != ref.String() ***REMOVED***
		return digest.Digest(""), refstore.ErrDoesNotExist
	***REMOVED***
	return r.pluginID, nil
***REMOVED***

func (r *pluginReference) AddTag(ref reference.Named, id digest.Digest, force bool) error ***REMOVED***
	// Read only, ignore
	return nil
***REMOVED***
func (r *pluginReference) AddDigest(ref reference.Canonical, id digest.Digest, force bool) error ***REMOVED***
	// Read only, ignore
	return nil
***REMOVED***
func (r *pluginReference) Delete(ref reference.Named) (bool, error) ***REMOVED***
	// Read only, ignore
	return false, nil
***REMOVED***

type pluginConfigStore struct ***REMOVED***
	pm     *Manager
	plugin *v2.Plugin
***REMOVED***

func (s *pluginConfigStore) Put([]byte) (digest.Digest, error) ***REMOVED***
	return digest.Digest(""), errors.New("cannot store config on push")
***REMOVED***

func (s *pluginConfigStore) Get(d digest.Digest) ([]byte, error) ***REMOVED***
	if s.plugin.Config != d ***REMOVED***
		return nil, errors.New("plugin not found")
	***REMOVED***
	rwc, err := s.pm.blobStore.Get(d)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer rwc.Close()
	return ioutil.ReadAll(rwc)
***REMOVED***

func (s *pluginConfigStore) RootFSAndOSFromConfig(c []byte) (*image.RootFS, string, error) ***REMOVED***
	return configToRootFS(c)
***REMOVED***

type pluginLayerProvider struct ***REMOVED***
	pm     *Manager
	plugin *v2.Plugin
***REMOVED***

func (p *pluginLayerProvider) Get(id layer.ChainID) (distribution.PushLayer, error) ***REMOVED***
	rootFS := rootFSFromPlugin(p.plugin.PluginObj.Config.Rootfs)
	var i int
	for i = 1; i <= len(rootFS.DiffIDs); i++ ***REMOVED***
		if layer.CreateChainID(rootFS.DiffIDs[:i]) == id ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	if i > len(rootFS.DiffIDs) ***REMOVED***
		return nil, errors.New("layer not found")
	***REMOVED***
	return &pluginLayer***REMOVED***
		pm:      p.pm,
		diffIDs: rootFS.DiffIDs[:i],
		blobs:   p.plugin.Blobsums[:i],
	***REMOVED***, nil
***REMOVED***

type pluginLayer struct ***REMOVED***
	pm      *Manager
	diffIDs []layer.DiffID
	blobs   []digest.Digest
***REMOVED***

func (l *pluginLayer) ChainID() layer.ChainID ***REMOVED***
	return layer.CreateChainID(l.diffIDs)
***REMOVED***

func (l *pluginLayer) DiffID() layer.DiffID ***REMOVED***
	return l.diffIDs[len(l.diffIDs)-1]
***REMOVED***

func (l *pluginLayer) Parent() distribution.PushLayer ***REMOVED***
	if len(l.diffIDs) == 1 ***REMOVED***
		return nil
	***REMOVED***
	return &pluginLayer***REMOVED***
		pm:      l.pm,
		diffIDs: l.diffIDs[:len(l.diffIDs)-1],
		blobs:   l.blobs[:len(l.diffIDs)-1],
	***REMOVED***
***REMOVED***

func (l *pluginLayer) Open() (io.ReadCloser, error) ***REMOVED***
	return l.pm.blobStore.Get(l.blobs[len(l.diffIDs)-1])
***REMOVED***

func (l *pluginLayer) Size() (int64, error) ***REMOVED***
	return l.pm.blobStore.Size(l.blobs[len(l.diffIDs)-1])
***REMOVED***

func (l *pluginLayer) MediaType() string ***REMOVED***
	return schema2.MediaTypeLayer
***REMOVED***

func (l *pluginLayer) Release() ***REMOVED***
	// Nothing needs to be release, no references held
***REMOVED***

// Remove deletes plugin's root directory.
func (pm *Manager) Remove(name string, config *types.PluginRmConfig) error ***REMOVED***
	p, err := pm.config.Store.GetV2Plugin(name)
	pm.mu.RLock()
	c := pm.cMap[p]
	pm.mu.RUnlock()

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if !config.ForceRemove ***REMOVED***
		if p.GetRefCount() > 0 ***REMOVED***
			return inUseError(p.Name())
		***REMOVED***
		if p.IsEnabled() ***REMOVED***
			return enabledError(p.Name())
		***REMOVED***
	***REMOVED***

	if p.IsEnabled() ***REMOVED***
		if err := pm.disable(p, c); err != nil ***REMOVED***
			logrus.Errorf("failed to disable plugin '%s': %s", p.Name(), err)
		***REMOVED***
	***REMOVED***

	defer func() ***REMOVED***
		go pm.GC()
	***REMOVED***()

	id := p.GetID()
	pluginDir := filepath.Join(pm.config.Root, id)

	if err := mount.RecursiveUnmount(pluginDir); err != nil ***REMOVED***
		return errors.Wrap(err, "error unmounting plugin data")
	***REMOVED***

	if err := atomicRemoveAll(pluginDir); err != nil ***REMOVED***
		return err
	***REMOVED***

	pm.config.Store.Remove(p)
	pm.config.LogPluginEvent(id, name, "remove")
	pm.publisher.Publish(EventRemove***REMOVED***Plugin: p.PluginObj***REMOVED***)
	return nil
***REMOVED***

// Set sets plugin args
func (pm *Manager) Set(name string, args []string) error ***REMOVED***
	p, err := pm.config.Store.GetV2Plugin(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := p.Set(args); err != nil ***REMOVED***
		return err
	***REMOVED***
	return pm.save(p)
***REMOVED***

// CreateFromContext creates a plugin from the given pluginDir which contains
// both the rootfs and the config.json and a repoName with optional tag.
func (pm *Manager) CreateFromContext(ctx context.Context, tarCtx io.ReadCloser, options *types.PluginCreateOptions) (err error) ***REMOVED***
	pm.muGC.RLock()
	defer pm.muGC.RUnlock()

	ref, err := reference.ParseNormalizedNamed(options.RepoName)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to parse reference %v", options.RepoName)
	***REMOVED***
	if _, ok := ref.(reference.Canonical); ok ***REMOVED***
		return errors.Errorf("canonical references are not permitted")
	***REMOVED***
	name := reference.FamiliarString(reference.TagNameOnly(ref))

	if err := pm.config.Store.validateName(name); err != nil ***REMOVED*** // fast check, real check is in createPlugin()
		return err
	***REMOVED***

	tmpRootFSDir, err := ioutil.TempDir(pm.tmpDir(), ".rootfs")
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to create temp directory")
	***REMOVED***
	defer os.RemoveAll(tmpRootFSDir)

	var configJSON []byte
	rootFS := splitConfigRootFSFromTar(tarCtx, &configJSON)

	rootFSBlob, err := pm.blobStore.New()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer rootFSBlob.Close()
	gzw := gzip.NewWriter(rootFSBlob)
	layerDigester := digest.Canonical.Digester()
	rootFSReader := io.TeeReader(rootFS, io.MultiWriter(gzw, layerDigester.Hash()))

	if err := chrootarchive.Untar(rootFSReader, tmpRootFSDir, nil); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := rootFS.Close(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if configJSON == nil ***REMOVED***
		return errors.New("config not found")
	***REMOVED***

	if err := gzw.Close(); err != nil ***REMOVED***
		return errors.Wrap(err, "error closing gzip writer")
	***REMOVED***

	var config types.PluginConfig
	if err := json.Unmarshal(configJSON, &config); err != nil ***REMOVED***
		return errors.Wrap(err, "failed to parse config")
	***REMOVED***

	if err := pm.validateConfig(config); err != nil ***REMOVED***
		return err
	***REMOVED***

	pm.mu.Lock()
	defer pm.mu.Unlock()

	rootFSBlobsum, err := rootFSBlob.Commit()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			go pm.GC()
		***REMOVED***
	***REMOVED***()

	config.Rootfs = &types.PluginConfigRootfs***REMOVED***
		Type:    "layers",
		DiffIds: []string***REMOVED***layerDigester.Digest().String()***REMOVED***,
	***REMOVED***

	config.DockerVersion = dockerversion.Version

	configBlob, err := pm.blobStore.New()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer configBlob.Close()
	if err := json.NewEncoder(configBlob).Encode(config); err != nil ***REMOVED***
		return errors.Wrap(err, "error encoding json config")
	***REMOVED***
	configBlobsum, err := configBlob.Commit()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	p, err := pm.createPlugin(name, configBlobsum, []digest.Digest***REMOVED***rootFSBlobsum***REMOVED***, tmpRootFSDir, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	p.PluginObj.PluginReference = name

	pm.publisher.Publish(EventCreate***REMOVED***Plugin: p.PluginObj***REMOVED***)
	pm.config.LogPluginEvent(p.PluginObj.ID, name, "create")

	return nil
***REMOVED***

func (pm *Manager) validateConfig(config types.PluginConfig) error ***REMOVED***
	return nil // TODO:
***REMOVED***

func splitConfigRootFSFromTar(in io.ReadCloser, config *[]byte) io.ReadCloser ***REMOVED***
	pr, pw := io.Pipe()
	go func() ***REMOVED***
		tarReader := tar.NewReader(in)
		tarWriter := tar.NewWriter(pw)
		defer in.Close()

		hasRootFS := false

		for ***REMOVED***
			hdr, err := tarReader.Next()
			if err == io.EOF ***REMOVED***
				if !hasRootFS ***REMOVED***
					pw.CloseWithError(errors.Wrap(err, "no rootfs found"))
					return
				***REMOVED***
				// Signals end of archive.
				tarWriter.Close()
				pw.Close()
				return
			***REMOVED***
			if err != nil ***REMOVED***
				pw.CloseWithError(errors.Wrap(err, "failed to read from tar"))
				return
			***REMOVED***

			content := io.Reader(tarReader)
			name := path.Clean(hdr.Name)
			if path.IsAbs(name) ***REMOVED***
				name = name[1:]
			***REMOVED***
			if name == configFileName ***REMOVED***
				dt, err := ioutil.ReadAll(content)
				if err != nil ***REMOVED***
					pw.CloseWithError(errors.Wrapf(err, "failed to read %s", configFileName))
					return
				***REMOVED***
				*config = dt
			***REMOVED***
			if parts := strings.Split(name, "/"); len(parts) != 0 && parts[0] == rootFSFileName ***REMOVED***
				hdr.Name = path.Clean(path.Join(parts[1:]...))
				if hdr.Typeflag == tar.TypeLink && strings.HasPrefix(strings.ToLower(hdr.Linkname), rootFSFileName+"/") ***REMOVED***
					hdr.Linkname = hdr.Linkname[len(rootFSFileName)+1:]
				***REMOVED***
				if err := tarWriter.WriteHeader(hdr); err != nil ***REMOVED***
					pw.CloseWithError(errors.Wrap(err, "error writing tar header"))
					return
				***REMOVED***
				if _, err := pools.Copy(tarWriter, content); err != nil ***REMOVED***
					pw.CloseWithError(errors.Wrap(err, "error copying tar data"))
					return
				***REMOVED***
				hasRootFS = true
			***REMOVED*** else ***REMOVED***
				io.Copy(ioutil.Discard, content)
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	return pr
***REMOVED***

func atomicRemoveAll(dir string) error ***REMOVED***
	renamed := dir + "-removing"

	err := os.Rename(dir, renamed)
	switch ***REMOVED***
	case os.IsNotExist(err), err == nil:
		// even if `dir` doesn't exist, we can still try and remove `renamed`
	case os.IsExist(err):
		// Some previous remove failed, check if the origin dir exists
		if e := system.EnsureRemoveAll(renamed); e != nil ***REMOVED***
			return errors.Wrap(err, "rename target already exists and could not be removed")
		***REMOVED***
		if _, err := os.Stat(dir); os.IsNotExist(err) ***REMOVED***
			// origin doesn't exist, nothing left to do
			return nil
		***REMOVED***

		// attempt to rename again
		if err := os.Rename(dir, renamed); err != nil ***REMOVED***
			return errors.Wrap(err, "failed to rename dir for atomic removal")
		***REMOVED***
	default:
		return errors.Wrap(err, "failed to rename dir for atomic removal")
	***REMOVED***

	if err := system.EnsureRemoveAll(renamed); err != nil ***REMOVED***
		os.Rename(renamed, dir)
		return err
	***REMOVED***
	return nil
***REMOVED***
