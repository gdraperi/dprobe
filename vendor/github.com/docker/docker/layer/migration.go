package layer

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/opencontainers/go-digest"
	"github.com/sirupsen/logrus"
	"github.com/vbatts/tar-split/tar/asm"
	"github.com/vbatts/tar-split/tar/storage"
)

// CreateRWLayerByGraphID creates a RWLayer in the layer store using
// the provided name with the given graphID. To get the RWLayer
// after migration the layer may be retrieved by the given name.
func (ls *layerStore) CreateRWLayerByGraphID(name, graphID string, parent ChainID) (err error) ***REMOVED***
	ls.mountL.Lock()
	defer ls.mountL.Unlock()
	m, ok := ls.mounts[name]
	if ok ***REMOVED***
		if m.parent.chainID != parent ***REMOVED***
			return errors.New("name conflict, mismatched parent")
		***REMOVED***
		if m.mountID != graphID ***REMOVED***
			return errors.New("mount already exists")
		***REMOVED***

		return nil
	***REMOVED***

	if !ls.driver.Exists(graphID) ***REMOVED***
		return fmt.Errorf("graph ID does not exist: %q", graphID)
	***REMOVED***

	var p *roLayer
	if string(parent) != "" ***REMOVED***
		p = ls.get(parent)
		if p == nil ***REMOVED***
			return ErrLayerDoesNotExist
		***REMOVED***

		// Release parent chain if error
		defer func() ***REMOVED***
			if err != nil ***REMOVED***
				ls.layerL.Lock()
				ls.releaseLayer(p)
				ls.layerL.Unlock()
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	// TODO: Ensure graphID has correct parent

	m = &mountedLayer***REMOVED***
		name:       name,
		parent:     p,
		mountID:    graphID,
		layerStore: ls,
		references: map[RWLayer]*referencedRWLayer***REMOVED******REMOVED***,
	***REMOVED***

	// Check for existing init layer
	initID := fmt.Sprintf("%s-init", graphID)
	if ls.driver.Exists(initID) ***REMOVED***
		m.initID = initID
	***REMOVED***

	return ls.saveMount(m)
***REMOVED***

func (ls *layerStore) ChecksumForGraphID(id, parent, oldTarDataPath, newTarDataPath string) (diffID DiffID, size int64, err error) ***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			logrus.Debugf("could not get checksum for %q with tar-split: %q", id, err)
			diffID, size, err = ls.checksumForGraphIDNoTarsplit(id, parent, newTarDataPath)
		***REMOVED***
	***REMOVED***()

	if oldTarDataPath == "" ***REMOVED***
		err = errors.New("no tar-split file")
		return
	***REMOVED***

	tarDataFile, err := os.Open(oldTarDataPath)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	defer tarDataFile.Close()
	uncompressed, err := gzip.NewReader(tarDataFile)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	dgst := digest.Canonical.Digester()
	err = ls.assembleTarTo(id, uncompressed, &size, dgst.Hash())
	if err != nil ***REMOVED***
		return
	***REMOVED***

	diffID = DiffID(dgst.Digest())
	err = os.RemoveAll(newTarDataPath)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	err = os.Link(oldTarDataPath, newTarDataPath)

	return
***REMOVED***

func (ls *layerStore) checksumForGraphIDNoTarsplit(id, parent, newTarDataPath string) (diffID DiffID, size int64, err error) ***REMOVED***
	rawarchive, err := ls.driver.Diff(id, parent)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	defer rawarchive.Close()

	f, err := os.Create(newTarDataPath)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	defer f.Close()
	mfz := gzip.NewWriter(f)
	defer mfz.Close()
	metaPacker := storage.NewJSONPacker(mfz)

	packerCounter := &packSizeCounter***REMOVED***metaPacker, &size***REMOVED***

	archive, err := asm.NewInputTarStream(rawarchive, packerCounter, nil)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	dgst, err := digest.FromReader(archive)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	diffID = DiffID(dgst)
	return
***REMOVED***

func (ls *layerStore) RegisterByGraphID(graphID string, parent ChainID, diffID DiffID, tarDataFile string, size int64) (Layer, error) ***REMOVED***
	// err is used to hold the error which will always trigger
	// cleanup of creates sources but may not be an error returned
	// to the caller (already exists).
	var err error
	var p *roLayer
	if string(parent) != "" ***REMOVED***
		p = ls.get(parent)
		if p == nil ***REMOVED***
			return nil, ErrLayerDoesNotExist
		***REMOVED***

		// Release parent chain if error
		defer func() ***REMOVED***
			if err != nil ***REMOVED***
				ls.layerL.Lock()
				ls.releaseLayer(p)
				ls.layerL.Unlock()
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	// Create new roLayer
	layer := &roLayer***REMOVED***
		parent:         p,
		cacheID:        graphID,
		referenceCount: 1,
		layerStore:     ls,
		references:     map[Layer]struct***REMOVED******REMOVED******REMOVED******REMOVED***,
		diffID:         diffID,
		size:           size,
		chainID:        createChainIDFromParent(parent, diffID),
	***REMOVED***

	ls.layerL.Lock()
	defer ls.layerL.Unlock()

	if existingLayer := ls.getWithoutLock(layer.chainID); existingLayer != nil ***REMOVED***
		// Set error for cleanup, but do not return
		err = errors.New("layer already exists")
		return existingLayer.getReference(), nil
	***REMOVED***

	tx, err := ls.store.StartTransaction()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			logrus.Debugf("Cleaning up transaction after failed migration for %s: %v", graphID, err)
			if err := tx.Cancel(); err != nil ***REMOVED***
				logrus.Errorf("Error canceling metadata transaction %q: %s", tx.String(), err)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	tsw, err := tx.TarSplitWriter(false)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer tsw.Close()
	tdf, err := os.Open(tarDataFile)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer tdf.Close()
	_, err = io.Copy(tsw, tdf)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err = storeLayer(tx, layer); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err = tx.Commit(layer.chainID); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ls.layerMap[layer.chainID] = layer

	return layer.getReference(), nil
***REMOVED***

type unpackSizeCounter struct ***REMOVED***
	unpacker storage.Unpacker
	size     *int64
***REMOVED***

func (u *unpackSizeCounter) Next() (*storage.Entry, error) ***REMOVED***
	e, err := u.unpacker.Next()
	if err == nil && u.size != nil ***REMOVED***
		*u.size += e.Size
	***REMOVED***
	return e, err
***REMOVED***

type packSizeCounter struct ***REMOVED***
	packer storage.Packer
	size   *int64
***REMOVED***

func (p *packSizeCounter) AddEntry(e storage.Entry) (int, error) ***REMOVED***
	n, err := p.packer.AddEntry(e)
	if err == nil && p.size != nil ***REMOVED***
		*p.size += e.Size
	***REMOVED***
	return n, err
***REMOVED***
