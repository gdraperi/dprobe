package layer

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"sync"

	"github.com/docker/distribution"
	"github.com/docker/docker/daemon/graphdriver"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/plugingetter"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/pkg/system"
	"github.com/opencontainers/go-digest"
	"github.com/sirupsen/logrus"
	"github.com/vbatts/tar-split/tar/asm"
	"github.com/vbatts/tar-split/tar/storage"
)

// maxLayerDepth represents the maximum number of
// layers which can be chained together. 125 was
// chosen to account for the 127 max in some
// graphdrivers plus the 2 additional layers
// used to create a rwlayer.
const maxLayerDepth = 125

type layerStore struct ***REMOVED***
	store       MetadataStore
	driver      graphdriver.Driver
	useTarSplit bool

	layerMap map[ChainID]*roLayer
	layerL   sync.Mutex

	mounts map[string]*mountedLayer
	mountL sync.Mutex
	os     string
***REMOVED***

// StoreOptions are the options used to create a new Store instance
type StoreOptions struct ***REMOVED***
	Root                      string
	MetadataStorePathTemplate string
	GraphDriver               string
	GraphDriverOptions        []string
	IDMappings                *idtools.IDMappings
	PluginGetter              plugingetter.PluginGetter
	ExperimentalEnabled       bool
	OS                        string
***REMOVED***

// NewStoreFromOptions creates a new Store instance
func NewStoreFromOptions(options StoreOptions) (Store, error) ***REMOVED***
	driver, err := graphdriver.New(options.GraphDriver, options.PluginGetter, graphdriver.Options***REMOVED***
		Root:                options.Root,
		DriverOptions:       options.GraphDriverOptions,
		UIDMaps:             options.IDMappings.UIDs(),
		GIDMaps:             options.IDMappings.GIDs(),
		ExperimentalEnabled: options.ExperimentalEnabled,
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("error initializing graphdriver: %v", err)
	***REMOVED***
	logrus.Debugf("Initialized graph driver %s", driver)

	fms, err := NewFSMetadataStore(fmt.Sprintf(options.MetadataStorePathTemplate, driver))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return NewStoreFromGraphDriver(fms, driver, options.OS)
***REMOVED***

// NewStoreFromGraphDriver creates a new Store instance using the provided
// metadata store and graph driver. The metadata store will be used to restore
// the Store.
func NewStoreFromGraphDriver(store MetadataStore, driver graphdriver.Driver, os string) (Store, error) ***REMOVED***
	if !system.IsOSSupported(os) ***REMOVED***
		return nil, fmt.Errorf("failed to initialize layer store as operating system '%s' is not supported", os)
	***REMOVED***
	caps := graphdriver.Capabilities***REMOVED******REMOVED***
	if capDriver, ok := driver.(graphdriver.CapabilityDriver); ok ***REMOVED***
		caps = capDriver.Capabilities()
	***REMOVED***

	ls := &layerStore***REMOVED***
		store:       store,
		driver:      driver,
		layerMap:    map[ChainID]*roLayer***REMOVED******REMOVED***,
		mounts:      map[string]*mountedLayer***REMOVED******REMOVED***,
		useTarSplit: !caps.ReproducesExactDiffs,
		os:          os,
	***REMOVED***

	ids, mounts, err := store.List()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	for _, id := range ids ***REMOVED***
		l, err := ls.loadLayer(id)
		if err != nil ***REMOVED***
			logrus.Debugf("Failed to load layer %s: %s", id, err)
			continue
		***REMOVED***
		if l.parent != nil ***REMOVED***
			l.parent.referenceCount++
		***REMOVED***
	***REMOVED***

	for _, mount := range mounts ***REMOVED***
		if err := ls.loadMount(mount); err != nil ***REMOVED***
			logrus.Debugf("Failed to load mount %s: %s", mount, err)
		***REMOVED***
	***REMOVED***

	return ls, nil
***REMOVED***

func (ls *layerStore) loadLayer(layer ChainID) (*roLayer, error) ***REMOVED***
	cl, ok := ls.layerMap[layer]
	if ok ***REMOVED***
		return cl, nil
	***REMOVED***

	diff, err := ls.store.GetDiffID(layer)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to get diff id for %s: %s", layer, err)
	***REMOVED***

	size, err := ls.store.GetSize(layer)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to get size for %s: %s", layer, err)
	***REMOVED***

	cacheID, err := ls.store.GetCacheID(layer)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to get cache id for %s: %s", layer, err)
	***REMOVED***

	parent, err := ls.store.GetParent(layer)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to get parent for %s: %s", layer, err)
	***REMOVED***

	descriptor, err := ls.store.GetDescriptor(layer)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to get descriptor for %s: %s", layer, err)
	***REMOVED***

	os, err := ls.store.getOS(layer)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to get operating system for %s: %s", layer, err)
	***REMOVED***

	if os != ls.os ***REMOVED***
		return nil, fmt.Errorf("failed to load layer with os %s into layerstore for %s", os, ls.os)
	***REMOVED***

	cl = &roLayer***REMOVED***
		chainID:    layer,
		diffID:     diff,
		size:       size,
		cacheID:    cacheID,
		layerStore: ls,
		references: map[Layer]struct***REMOVED******REMOVED******REMOVED******REMOVED***,
		descriptor: descriptor,
	***REMOVED***

	if parent != "" ***REMOVED***
		p, err := ls.loadLayer(parent)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		cl.parent = p
	***REMOVED***

	ls.layerMap[cl.chainID] = cl

	return cl, nil
***REMOVED***

func (ls *layerStore) loadMount(mount string) error ***REMOVED***
	if _, ok := ls.mounts[mount]; ok ***REMOVED***
		return nil
	***REMOVED***

	mountID, err := ls.store.GetMountID(mount)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	initID, err := ls.store.GetInitID(mount)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	parent, err := ls.store.GetMountParent(mount)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	ml := &mountedLayer***REMOVED***
		name:       mount,
		mountID:    mountID,
		initID:     initID,
		layerStore: ls,
		references: map[RWLayer]*referencedRWLayer***REMOVED******REMOVED***,
	***REMOVED***

	if parent != "" ***REMOVED***
		p, err := ls.loadLayer(parent)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		ml.parent = p

		p.referenceCount++
	***REMOVED***

	ls.mounts[ml.name] = ml

	return nil
***REMOVED***

func (ls *layerStore) applyTar(tx MetadataTransaction, ts io.Reader, parent string, layer *roLayer) error ***REMOVED***
	digester := digest.Canonical.Digester()
	tr := io.TeeReader(ts, digester.Hash())

	rdr := tr
	if ls.useTarSplit ***REMOVED***
		tsw, err := tx.TarSplitWriter(true)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		metaPacker := storage.NewJSONPacker(tsw)
		defer tsw.Close()

		// we're passing nil here for the file putter, because the ApplyDiff will
		// handle the extraction of the archive
		rdr, err = asm.NewInputTarStream(tr, metaPacker, nil)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	applySize, err := ls.driver.ApplyDiff(layer.cacheID, parent, rdr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Discard trailing data but ensure metadata is picked up to reconstruct stream
	io.Copy(ioutil.Discard, rdr) // ignore error as reader may be closed

	layer.size = applySize
	layer.diffID = DiffID(digester.Digest())

	logrus.Debugf("Applied tar %s to %s, size: %d", layer.diffID, layer.cacheID, applySize)

	return nil
***REMOVED***

func (ls *layerStore) Register(ts io.Reader, parent ChainID) (Layer, error) ***REMOVED***
	return ls.registerWithDescriptor(ts, parent, distribution.Descriptor***REMOVED******REMOVED***)
***REMOVED***

func (ls *layerStore) registerWithDescriptor(ts io.Reader, parent ChainID, descriptor distribution.Descriptor) (Layer, error) ***REMOVED***
	// err is used to hold the error which will always trigger
	// cleanup of creates sources but may not be an error returned
	// to the caller (already exists).
	var err error
	var pid string
	var p *roLayer

	if string(parent) != "" ***REMOVED***
		p = ls.get(parent)
		if p == nil ***REMOVED***
			return nil, ErrLayerDoesNotExist
		***REMOVED***
		pid = p.cacheID
		// Release parent chain if error
		defer func() ***REMOVED***
			if err != nil ***REMOVED***
				ls.layerL.Lock()
				ls.releaseLayer(p)
				ls.layerL.Unlock()
			***REMOVED***
		***REMOVED***()
		if p.depth() >= maxLayerDepth ***REMOVED***
			err = ErrMaxDepthExceeded
			return nil, err
		***REMOVED***
	***REMOVED***

	// Create new roLayer
	layer := &roLayer***REMOVED***
		parent:         p,
		cacheID:        stringid.GenerateRandomID(),
		referenceCount: 1,
		layerStore:     ls,
		references:     map[Layer]struct***REMOVED******REMOVED******REMOVED******REMOVED***,
		descriptor:     descriptor,
	***REMOVED***

	if err = ls.driver.Create(layer.cacheID, pid, nil); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	tx, err := ls.store.StartTransaction()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			logrus.Debugf("Cleaning up layer %s: %v", layer.cacheID, err)
			if err := ls.driver.Remove(layer.cacheID); err != nil ***REMOVED***
				logrus.Errorf("Error cleaning up cache layer %s: %v", layer.cacheID, err)
			***REMOVED***
			if err := tx.Cancel(); err != nil ***REMOVED***
				logrus.Errorf("Error canceling metadata transaction %q: %s", tx.String(), err)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	if err = ls.applyTar(tx, ts, pid, layer); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if layer.parent == nil ***REMOVED***
		layer.chainID = ChainID(layer.diffID)
	***REMOVED*** else ***REMOVED***
		layer.chainID = createChainIDFromParent(layer.parent.chainID, layer.diffID)
	***REMOVED***

	if err = storeLayer(tx, layer); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ls.layerL.Lock()
	defer ls.layerL.Unlock()

	if existingLayer := ls.getWithoutLock(layer.chainID); existingLayer != nil ***REMOVED***
		// Set error for cleanup, but do not return the error
		err = errors.New("layer already exists")
		return existingLayer.getReference(), nil
	***REMOVED***

	if err = tx.Commit(layer.chainID); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ls.layerMap[layer.chainID] = layer

	return layer.getReference(), nil
***REMOVED***

func (ls *layerStore) getWithoutLock(layer ChainID) *roLayer ***REMOVED***
	l, ok := ls.layerMap[layer]
	if !ok ***REMOVED***
		return nil
	***REMOVED***

	l.referenceCount++

	return l
***REMOVED***

func (ls *layerStore) get(l ChainID) *roLayer ***REMOVED***
	ls.layerL.Lock()
	defer ls.layerL.Unlock()
	return ls.getWithoutLock(l)
***REMOVED***

func (ls *layerStore) Get(l ChainID) (Layer, error) ***REMOVED***
	ls.layerL.Lock()
	defer ls.layerL.Unlock()

	layer := ls.getWithoutLock(l)
	if layer == nil ***REMOVED***
		return nil, ErrLayerDoesNotExist
	***REMOVED***

	return layer.getReference(), nil
***REMOVED***

func (ls *layerStore) Map() map[ChainID]Layer ***REMOVED***
	ls.layerL.Lock()
	defer ls.layerL.Unlock()

	layers := map[ChainID]Layer***REMOVED******REMOVED***

	for k, v := range ls.layerMap ***REMOVED***
		layers[k] = v
	***REMOVED***

	return layers
***REMOVED***

func (ls *layerStore) deleteLayer(layer *roLayer, metadata *Metadata) error ***REMOVED***
	err := ls.driver.Remove(layer.cacheID)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = ls.store.Remove(layer.chainID)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	metadata.DiffID = layer.diffID
	metadata.ChainID = layer.chainID
	metadata.Size, err = layer.Size()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	metadata.DiffSize = layer.size

	return nil
***REMOVED***

func (ls *layerStore) releaseLayer(l *roLayer) ([]Metadata, error) ***REMOVED***
	depth := 0
	removed := []Metadata***REMOVED******REMOVED***
	for ***REMOVED***
		if l.referenceCount == 0 ***REMOVED***
			panic("layer not retained")
		***REMOVED***
		l.referenceCount--
		if l.referenceCount != 0 ***REMOVED***
			return removed, nil
		***REMOVED***

		if len(removed) == 0 && depth > 0 ***REMOVED***
			panic("cannot remove layer with child")
		***REMOVED***
		if l.hasReferences() ***REMOVED***
			panic("cannot delete referenced layer")
		***REMOVED***
		var metadata Metadata
		if err := ls.deleteLayer(l, &metadata); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		delete(ls.layerMap, l.chainID)
		removed = append(removed, metadata)

		if l.parent == nil ***REMOVED***
			return removed, nil
		***REMOVED***

		depth++
		l = l.parent
	***REMOVED***
***REMOVED***

func (ls *layerStore) Release(l Layer) ([]Metadata, error) ***REMOVED***
	ls.layerL.Lock()
	defer ls.layerL.Unlock()
	layer, ok := ls.layerMap[l.ChainID()]
	if !ok ***REMOVED***
		return []Metadata***REMOVED******REMOVED***, nil
	***REMOVED***
	if !layer.hasReference(l) ***REMOVED***
		return nil, ErrLayerNotRetained
	***REMOVED***

	layer.deleteReference(l)

	return ls.releaseLayer(layer)
***REMOVED***

func (ls *layerStore) CreateRWLayer(name string, parent ChainID, opts *CreateRWLayerOpts) (RWLayer, error) ***REMOVED***
	var (
		storageOpt map[string]string
		initFunc   MountInit
		mountLabel string
	)

	if opts != nil ***REMOVED***
		mountLabel = opts.MountLabel
		storageOpt = opts.StorageOpt
		initFunc = opts.InitFunc
	***REMOVED***

	ls.mountL.Lock()
	defer ls.mountL.Unlock()
	m, ok := ls.mounts[name]
	if ok ***REMOVED***
		return nil, ErrMountNameConflict
	***REMOVED***

	var err error
	var pid string
	var p *roLayer
	if string(parent) != "" ***REMOVED***
		p = ls.get(parent)
		if p == nil ***REMOVED***
			return nil, ErrLayerDoesNotExist
		***REMOVED***
		pid = p.cacheID

		// Release parent chain if error
		defer func() ***REMOVED***
			if err != nil ***REMOVED***
				ls.layerL.Lock()
				ls.releaseLayer(p)
				ls.layerL.Unlock()
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	m = &mountedLayer***REMOVED***
		name:       name,
		parent:     p,
		mountID:    ls.mountID(name),
		layerStore: ls,
		references: map[RWLayer]*referencedRWLayer***REMOVED******REMOVED***,
	***REMOVED***

	if initFunc != nil ***REMOVED***
		pid, err = ls.initMount(m.mountID, pid, mountLabel, initFunc, storageOpt)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		m.initID = pid
	***REMOVED***

	createOpts := &graphdriver.CreateOpts***REMOVED***
		StorageOpt: storageOpt,
	***REMOVED***

	if err = ls.driver.CreateReadWrite(m.mountID, pid, createOpts); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err = ls.saveMount(m); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return m.getReference(), nil
***REMOVED***

func (ls *layerStore) GetRWLayer(id string) (RWLayer, error) ***REMOVED***
	ls.mountL.Lock()
	defer ls.mountL.Unlock()
	mount, ok := ls.mounts[id]
	if !ok ***REMOVED***
		return nil, ErrMountDoesNotExist
	***REMOVED***

	return mount.getReference(), nil
***REMOVED***

func (ls *layerStore) GetMountID(id string) (string, error) ***REMOVED***
	ls.mountL.Lock()
	defer ls.mountL.Unlock()
	mount, ok := ls.mounts[id]
	if !ok ***REMOVED***
		return "", ErrMountDoesNotExist
	***REMOVED***
	logrus.Debugf("GetMountID id: %s -> mountID: %s", id, mount.mountID)

	return mount.mountID, nil
***REMOVED***

func (ls *layerStore) ReleaseRWLayer(l RWLayer) ([]Metadata, error) ***REMOVED***
	ls.mountL.Lock()
	defer ls.mountL.Unlock()
	m, ok := ls.mounts[l.Name()]
	if !ok ***REMOVED***
		return []Metadata***REMOVED******REMOVED***, nil
	***REMOVED***

	if err := m.deleteReference(l); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if m.hasReferences() ***REMOVED***
		return []Metadata***REMOVED******REMOVED***, nil
	***REMOVED***

	if err := ls.driver.Remove(m.mountID); err != nil ***REMOVED***
		logrus.Errorf("Error removing mounted layer %s: %s", m.name, err)
		m.retakeReference(l)
		return nil, err
	***REMOVED***

	if m.initID != "" ***REMOVED***
		if err := ls.driver.Remove(m.initID); err != nil ***REMOVED***
			logrus.Errorf("Error removing init layer %s: %s", m.name, err)
			m.retakeReference(l)
			return nil, err
		***REMOVED***
	***REMOVED***

	if err := ls.store.RemoveMount(m.name); err != nil ***REMOVED***
		logrus.Errorf("Error removing mount metadata: %s: %s", m.name, err)
		m.retakeReference(l)
		return nil, err
	***REMOVED***

	delete(ls.mounts, m.Name())

	ls.layerL.Lock()
	defer ls.layerL.Unlock()
	if m.parent != nil ***REMOVED***
		return ls.releaseLayer(m.parent)
	***REMOVED***

	return []Metadata***REMOVED******REMOVED***, nil
***REMOVED***

func (ls *layerStore) saveMount(mount *mountedLayer) error ***REMOVED***
	if err := ls.store.SetMountID(mount.name, mount.mountID); err != nil ***REMOVED***
		return err
	***REMOVED***

	if mount.initID != "" ***REMOVED***
		if err := ls.store.SetInitID(mount.name, mount.initID); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if mount.parent != nil ***REMOVED***
		if err := ls.store.SetMountParent(mount.name, mount.parent.chainID); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	ls.mounts[mount.name] = mount

	return nil
***REMOVED***

func (ls *layerStore) initMount(graphID, parent, mountLabel string, initFunc MountInit, storageOpt map[string]string) (string, error) ***REMOVED***
	// Use "<graph-id>-init" to maintain compatibility with graph drivers
	// which are expecting this layer with this special name. If all
	// graph drivers can be updated to not rely on knowing about this layer
	// then the initID should be randomly generated.
	initID := fmt.Sprintf("%s-init", graphID)

	createOpts := &graphdriver.CreateOpts***REMOVED***
		MountLabel: mountLabel,
		StorageOpt: storageOpt,
	***REMOVED***

	if err := ls.driver.CreateReadWrite(initID, parent, createOpts); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	p, err := ls.driver.Get(initID, "")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	if err := initFunc(p); err != nil ***REMOVED***
		ls.driver.Put(initID)
		return "", err
	***REMOVED***

	if err := ls.driver.Put(initID); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return initID, nil
***REMOVED***

func (ls *layerStore) getTarStream(rl *roLayer) (io.ReadCloser, error) ***REMOVED***
	if !ls.useTarSplit ***REMOVED***
		var parentCacheID string
		if rl.parent != nil ***REMOVED***
			parentCacheID = rl.parent.cacheID
		***REMOVED***

		return ls.driver.Diff(rl.cacheID, parentCacheID)
	***REMOVED***

	r, err := ls.store.TarSplitReader(rl.chainID)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	pr, pw := io.Pipe()
	go func() ***REMOVED***
		err := ls.assembleTarTo(rl.cacheID, r, nil, pw)
		if err != nil ***REMOVED***
			pw.CloseWithError(err)
		***REMOVED*** else ***REMOVED***
			pw.Close()
		***REMOVED***
	***REMOVED***()

	return pr, nil
***REMOVED***

func (ls *layerStore) assembleTarTo(graphID string, metadata io.ReadCloser, size *int64, w io.Writer) error ***REMOVED***
	diffDriver, ok := ls.driver.(graphdriver.DiffGetterDriver)
	if !ok ***REMOVED***
		diffDriver = &naiveDiffPathDriver***REMOVED***ls.driver***REMOVED***
	***REMOVED***

	defer metadata.Close()

	// get our relative path to the container
	fileGetCloser, err := diffDriver.DiffGetter(graphID)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer fileGetCloser.Close()

	metaUnpacker := storage.NewJSONUnpacker(metadata)
	upackerCounter := &unpackSizeCounter***REMOVED***metaUnpacker, size***REMOVED***
	logrus.Debugf("Assembling tar data for %s", graphID)
	return asm.WriteOutputTarStream(fileGetCloser, upackerCounter, w)
***REMOVED***

func (ls *layerStore) Cleanup() error ***REMOVED***
	return ls.driver.Cleanup()
***REMOVED***

func (ls *layerStore) DriverStatus() [][2]string ***REMOVED***
	return ls.driver.Status()
***REMOVED***

func (ls *layerStore) DriverName() string ***REMOVED***
	return ls.driver.String()
***REMOVED***

type naiveDiffPathDriver struct ***REMOVED***
	graphdriver.Driver
***REMOVED***

type fileGetPutter struct ***REMOVED***
	storage.FileGetter
	driver graphdriver.Driver
	id     string
***REMOVED***

func (w *fileGetPutter) Close() error ***REMOVED***
	return w.driver.Put(w.id)
***REMOVED***

func (n *naiveDiffPathDriver) DiffGetter(id string) (graphdriver.FileGetCloser, error) ***REMOVED***
	p, err := n.Driver.Get(id, "")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &fileGetPutter***REMOVED***storage.NewPathFileGetter(p.Path()), n.Driver, id***REMOVED***, nil
***REMOVED***
