package xfer

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"time"

	"github.com/docker/distribution"
	"github.com/docker/docker/image"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/system"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const maxDownloadAttempts = 5

// LayerDownloadManager figures out which layers need to be downloaded, then
// registers and downloads those, taking into account dependencies between
// layers.
type LayerDownloadManager struct ***REMOVED***
	layerStores  map[string]layer.Store
	tm           TransferManager
	waitDuration time.Duration
***REMOVED***

// SetConcurrency sets the max concurrent downloads for each pull
func (ldm *LayerDownloadManager) SetConcurrency(concurrency int) ***REMOVED***
	ldm.tm.SetConcurrency(concurrency)
***REMOVED***

// NewLayerDownloadManager returns a new LayerDownloadManager.
func NewLayerDownloadManager(layerStores map[string]layer.Store, concurrencyLimit int, options ...func(*LayerDownloadManager)) *LayerDownloadManager ***REMOVED***
	manager := LayerDownloadManager***REMOVED***
		layerStores:  layerStores,
		tm:           NewTransferManager(concurrencyLimit),
		waitDuration: time.Second,
	***REMOVED***
	for _, option := range options ***REMOVED***
		option(&manager)
	***REMOVED***
	return &manager
***REMOVED***

type downloadTransfer struct ***REMOVED***
	Transfer

	layerStore layer.Store
	layer      layer.Layer
	err        error
***REMOVED***

// result returns the layer resulting from the download, if the download
// and registration were successful.
func (d *downloadTransfer) result() (layer.Layer, error) ***REMOVED***
	return d.layer, d.err
***REMOVED***

// A DownloadDescriptor references a layer that may need to be downloaded.
type DownloadDescriptor interface ***REMOVED***
	// Key returns the key used to deduplicate downloads.
	Key() string
	// ID returns the ID for display purposes.
	ID() string
	// DiffID should return the DiffID for this layer, or an error
	// if it is unknown (for example, if it has not been downloaded
	// before).
	DiffID() (layer.DiffID, error)
	// Download is called to perform the download.
	Download(ctx context.Context, progressOutput progress.Output) (io.ReadCloser, int64, error)
	// Close is called when the download manager is finished with this
	// descriptor and will not call Download again or read from the reader
	// that Download returned.
	Close()
***REMOVED***

// DownloadDescriptorWithRegistered is a DownloadDescriptor that has an
// additional Registered method which gets called after a downloaded layer is
// registered. This allows the user of the download manager to know the DiffID
// of each registered layer. This method is called if a cast to
// DownloadDescriptorWithRegistered is successful.
type DownloadDescriptorWithRegistered interface ***REMOVED***
	DownloadDescriptor
	Registered(diffID layer.DiffID)
***REMOVED***

// Download is a blocking function which ensures the requested layers are
// present in the layer store. It uses the string returned by the Key method to
// deduplicate downloads. If a given layer is not already known to present in
// the layer store, and the key is not used by an in-progress download, the
// Download method is called to get the layer tar data. Layers are then
// registered in the appropriate order.  The caller must call the returned
// release function once it is done with the returned RootFS object.
func (ldm *LayerDownloadManager) Download(ctx context.Context, initialRootFS image.RootFS, os string, layers []DownloadDescriptor, progressOutput progress.Output) (image.RootFS, func(), error) ***REMOVED***
	var (
		topLayer       layer.Layer
		topDownload    *downloadTransfer
		watcher        *Watcher
		missingLayer   bool
		transferKey    = ""
		downloadsByKey = make(map[string]*downloadTransfer)
	)

	// Assume that the operating system is the host OS if blank, and validate it
	// to ensure we don't cause a panic by an invalid index into the layerstores.
	if os == "" ***REMOVED***
		os = runtime.GOOS
	***REMOVED***
	if !system.IsOSSupported(os) ***REMOVED***
		return image.RootFS***REMOVED******REMOVED***, nil, system.ErrNotSupportedOperatingSystem
	***REMOVED***

	rootFS := initialRootFS
	for _, descriptor := range layers ***REMOVED***
		key := descriptor.Key()
		transferKey += key

		if !missingLayer ***REMOVED***
			missingLayer = true
			diffID, err := descriptor.DiffID()
			if err == nil ***REMOVED***
				getRootFS := rootFS
				getRootFS.Append(diffID)
				l, err := ldm.layerStores[os].Get(getRootFS.ChainID())
				if err == nil ***REMOVED***
					// Layer already exists.
					logrus.Debugf("Layer already exists: %s", descriptor.ID())
					progress.Update(progressOutput, descriptor.ID(), "Already exists")
					if topLayer != nil ***REMOVED***
						layer.ReleaseAndLog(ldm.layerStores[os], topLayer)
					***REMOVED***
					topLayer = l
					missingLayer = false
					rootFS.Append(diffID)
					// Register this repository as a source of this layer.
					withRegistered, hasRegistered := descriptor.(DownloadDescriptorWithRegistered)
					if hasRegistered ***REMOVED*** // As layerstore may set the driver
						withRegistered.Registered(diffID)
					***REMOVED***
					continue
				***REMOVED***
			***REMOVED***
		***REMOVED***

		// Does this layer have the same data as a previous layer in
		// the stack? If so, avoid downloading it more than once.
		var topDownloadUncasted Transfer
		if existingDownload, ok := downloadsByKey[key]; ok ***REMOVED***
			xferFunc := ldm.makeDownloadFuncFromDownload(descriptor, existingDownload, topDownload, os)
			defer topDownload.Transfer.Release(watcher)
			topDownloadUncasted, watcher = ldm.tm.Transfer(transferKey, xferFunc, progressOutput)
			topDownload = topDownloadUncasted.(*downloadTransfer)
			continue
		***REMOVED***

		// Layer is not known to exist - download and register it.
		progress.Update(progressOutput, descriptor.ID(), "Pulling fs layer")

		var xferFunc DoFunc
		if topDownload != nil ***REMOVED***
			xferFunc = ldm.makeDownloadFunc(descriptor, "", topDownload, os)
			defer topDownload.Transfer.Release(watcher)
		***REMOVED*** else ***REMOVED***
			xferFunc = ldm.makeDownloadFunc(descriptor, rootFS.ChainID(), nil, os)
		***REMOVED***
		topDownloadUncasted, watcher = ldm.tm.Transfer(transferKey, xferFunc, progressOutput)
		topDownload = topDownloadUncasted.(*downloadTransfer)
		downloadsByKey[key] = topDownload
	***REMOVED***

	if topDownload == nil ***REMOVED***
		return rootFS, func() ***REMOVED***
			if topLayer != nil ***REMOVED***
				layer.ReleaseAndLog(ldm.layerStores[os], topLayer)
			***REMOVED***
		***REMOVED***, nil
	***REMOVED***

	// Won't be using the list built up so far - will generate it
	// from downloaded layers instead.
	rootFS.DiffIDs = []layer.DiffID***REMOVED******REMOVED***

	defer func() ***REMOVED***
		if topLayer != nil ***REMOVED***
			layer.ReleaseAndLog(ldm.layerStores[os], topLayer)
		***REMOVED***
	***REMOVED***()

	select ***REMOVED***
	case <-ctx.Done():
		topDownload.Transfer.Release(watcher)
		return rootFS, func() ***REMOVED******REMOVED***, ctx.Err()
	case <-topDownload.Done():
		break
	***REMOVED***

	l, err := topDownload.result()
	if err != nil ***REMOVED***
		topDownload.Transfer.Release(watcher)
		return rootFS, func() ***REMOVED******REMOVED***, err
	***REMOVED***

	// Must do this exactly len(layers) times, so we don't include the
	// base layer on Windows.
	for range layers ***REMOVED***
		if l == nil ***REMOVED***
			topDownload.Transfer.Release(watcher)
			return rootFS, func() ***REMOVED******REMOVED***, errors.New("internal error: too few parent layers")
		***REMOVED***
		rootFS.DiffIDs = append([]layer.DiffID***REMOVED***l.DiffID()***REMOVED***, rootFS.DiffIDs...)
		l = l.Parent()
	***REMOVED***
	return rootFS, func() ***REMOVED*** topDownload.Transfer.Release(watcher) ***REMOVED***, err
***REMOVED***

// makeDownloadFunc returns a function that performs the layer download and
// registration. If parentDownload is non-nil, it waits for that download to
// complete before the registration step, and registers the downloaded data
// on top of parentDownload's resulting layer. Otherwise, it registers the
// layer on top of the ChainID given by parentLayer.
func (ldm *LayerDownloadManager) makeDownloadFunc(descriptor DownloadDescriptor, parentLayer layer.ChainID, parentDownload *downloadTransfer, os string) DoFunc ***REMOVED***
	return func(progressChan chan<- progress.Progress, start <-chan struct***REMOVED******REMOVED***, inactive chan<- struct***REMOVED******REMOVED***) Transfer ***REMOVED***
		d := &downloadTransfer***REMOVED***
			Transfer:   NewTransfer(),
			layerStore: ldm.layerStores[os],
		***REMOVED***

		go func() ***REMOVED***
			defer func() ***REMOVED***
				close(progressChan)
			***REMOVED***()

			progressOutput := progress.ChanOutput(progressChan)

			select ***REMOVED***
			case <-start:
			default:
				progress.Update(progressOutput, descriptor.ID(), "Waiting")
				<-start
			***REMOVED***

			if parentDownload != nil ***REMOVED***
				// Did the parent download already fail or get
				// cancelled?
				select ***REMOVED***
				case <-parentDownload.Done():
					_, err := parentDownload.result()
					if err != nil ***REMOVED***
						d.err = err
						return
					***REMOVED***
				default:
				***REMOVED***
			***REMOVED***

			var (
				downloadReader io.ReadCloser
				size           int64
				err            error
				retries        int
			)

			defer descriptor.Close()

			for ***REMOVED***
				downloadReader, size, err = descriptor.Download(d.Transfer.Context(), progressOutput)
				if err == nil ***REMOVED***
					break
				***REMOVED***

				// If an error was returned because the context
				// was cancelled, we shouldn't retry.
				select ***REMOVED***
				case <-d.Transfer.Context().Done():
					d.err = err
					return
				default:
				***REMOVED***

				retries++
				if _, isDNR := err.(DoNotRetry); isDNR || retries == maxDownloadAttempts ***REMOVED***
					logrus.Errorf("Download failed: %v", err)
					d.err = err
					return
				***REMOVED***

				logrus.Errorf("Download failed, retrying: %v", err)
				delay := retries * 5
				ticker := time.NewTicker(ldm.waitDuration)

			selectLoop:
				for ***REMOVED***
					progress.Updatef(progressOutput, descriptor.ID(), "Retrying in %d second%s", delay, (map[bool]string***REMOVED***true: "s"***REMOVED***)[delay != 1])
					select ***REMOVED***
					case <-ticker.C:
						delay--
						if delay == 0 ***REMOVED***
							ticker.Stop()
							break selectLoop
						***REMOVED***
					case <-d.Transfer.Context().Done():
						ticker.Stop()
						d.err = errors.New("download cancelled during retry delay")
						return
					***REMOVED***

				***REMOVED***
			***REMOVED***

			close(inactive)

			if parentDownload != nil ***REMOVED***
				select ***REMOVED***
				case <-d.Transfer.Context().Done():
					d.err = errors.New("layer registration cancelled")
					downloadReader.Close()
					return
				case <-parentDownload.Done():
				***REMOVED***

				l, err := parentDownload.result()
				if err != nil ***REMOVED***
					d.err = err
					downloadReader.Close()
					return
				***REMOVED***
				parentLayer = l.ChainID()
			***REMOVED***

			reader := progress.NewProgressReader(ioutils.NewCancelReadCloser(d.Transfer.Context(), downloadReader), progressOutput, size, descriptor.ID(), "Extracting")
			defer reader.Close()

			inflatedLayerData, err := archive.DecompressStream(reader)
			if err != nil ***REMOVED***
				d.err = fmt.Errorf("could not get decompression stream: %v", err)
				return
			***REMOVED***

			var src distribution.Descriptor
			if fs, ok := descriptor.(distribution.Describable); ok ***REMOVED***
				src = fs.Descriptor()
			***REMOVED***
			if ds, ok := d.layerStore.(layer.DescribableStore); ok ***REMOVED***
				d.layer, err = ds.RegisterWithDescriptor(inflatedLayerData, parentLayer, src)
			***REMOVED*** else ***REMOVED***
				d.layer, err = d.layerStore.Register(inflatedLayerData, parentLayer)
			***REMOVED***
			if err != nil ***REMOVED***
				select ***REMOVED***
				case <-d.Transfer.Context().Done():
					d.err = errors.New("layer registration cancelled")
				default:
					d.err = fmt.Errorf("failed to register layer: %v", err)
				***REMOVED***
				return
			***REMOVED***

			progress.Update(progressOutput, descriptor.ID(), "Pull complete")
			withRegistered, hasRegistered := descriptor.(DownloadDescriptorWithRegistered)
			if hasRegistered ***REMOVED***
				withRegistered.Registered(d.layer.DiffID())
			***REMOVED***

			// Doesn't actually need to be its own goroutine, but
			// done like this so we can defer close(c).
			go func() ***REMOVED***
				<-d.Transfer.Released()
				if d.layer != nil ***REMOVED***
					layer.ReleaseAndLog(d.layerStore, d.layer)
				***REMOVED***
			***REMOVED***()
		***REMOVED***()

		return d
	***REMOVED***
***REMOVED***

// makeDownloadFuncFromDownload returns a function that performs the layer
// registration when the layer data is coming from an existing download. It
// waits for sourceDownload and parentDownload to complete, and then
// reregisters the data from sourceDownload's top layer on top of
// parentDownload. This function does not log progress output because it would
// interfere with the progress reporting for sourceDownload, which has the same
// Key.
func (ldm *LayerDownloadManager) makeDownloadFuncFromDownload(descriptor DownloadDescriptor, sourceDownload *downloadTransfer, parentDownload *downloadTransfer, os string) DoFunc ***REMOVED***
	return func(progressChan chan<- progress.Progress, start <-chan struct***REMOVED******REMOVED***, inactive chan<- struct***REMOVED******REMOVED***) Transfer ***REMOVED***
		d := &downloadTransfer***REMOVED***
			Transfer:   NewTransfer(),
			layerStore: ldm.layerStores[os],
		***REMOVED***

		go func() ***REMOVED***
			defer func() ***REMOVED***
				close(progressChan)
			***REMOVED***()

			<-start

			close(inactive)

			select ***REMOVED***
			case <-d.Transfer.Context().Done():
				d.err = errors.New("layer registration cancelled")
				return
			case <-parentDownload.Done():
			***REMOVED***

			l, err := parentDownload.result()
			if err != nil ***REMOVED***
				d.err = err
				return
			***REMOVED***
			parentLayer := l.ChainID()

			// sourceDownload should have already finished if
			// parentDownload finished, but wait for it explicitly
			// to be sure.
			select ***REMOVED***
			case <-d.Transfer.Context().Done():
				d.err = errors.New("layer registration cancelled")
				return
			case <-sourceDownload.Done():
			***REMOVED***

			l, err = sourceDownload.result()
			if err != nil ***REMOVED***
				d.err = err
				return
			***REMOVED***

			layerReader, err := l.TarStream()
			if err != nil ***REMOVED***
				d.err = err
				return
			***REMOVED***
			defer layerReader.Close()

			var src distribution.Descriptor
			if fs, ok := l.(distribution.Describable); ok ***REMOVED***
				src = fs.Descriptor()
			***REMOVED***
			if ds, ok := d.layerStore.(layer.DescribableStore); ok ***REMOVED***
				d.layer, err = ds.RegisterWithDescriptor(layerReader, parentLayer, src)
			***REMOVED*** else ***REMOVED***
				d.layer, err = d.layerStore.Register(layerReader, parentLayer)
			***REMOVED***
			if err != nil ***REMOVED***
				d.err = fmt.Errorf("failed to register layer: %v", err)
				return
			***REMOVED***

			withRegistered, hasRegistered := descriptor.(DownloadDescriptorWithRegistered)
			if hasRegistered ***REMOVED***
				withRegistered.Registered(d.layer.DiffID())
			***REMOVED***

			// Doesn't actually need to be its own goroutine, but
			// done like this so we can defer close(c).
			go func() ***REMOVED***
				<-d.Transfer.Released()
				if d.layer != nil ***REMOVED***
					layer.ReleaseAndLog(d.layerStore, d.layer)
				***REMOVED***
			***REMOVED***()
		***REMOVED***()

		return d
	***REMOVED***
***REMOVED***
