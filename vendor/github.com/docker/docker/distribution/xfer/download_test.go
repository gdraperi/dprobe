package xfer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/docker/distribution"
	"github.com/docker/docker/image"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/progress"
	"github.com/opencontainers/go-digest"
	"golang.org/x/net/context"
)

const maxDownloadConcurrency = 3

type mockLayer struct ***REMOVED***
	layerData bytes.Buffer
	diffID    layer.DiffID
	chainID   layer.ChainID
	parent    layer.Layer
	os        string
***REMOVED***

func (ml *mockLayer) TarStream() (io.ReadCloser, error) ***REMOVED***
	return ioutil.NopCloser(bytes.NewBuffer(ml.layerData.Bytes())), nil
***REMOVED***

func (ml *mockLayer) TarStreamFrom(layer.ChainID) (io.ReadCloser, error) ***REMOVED***
	return nil, fmt.Errorf("not implemented")
***REMOVED***

func (ml *mockLayer) ChainID() layer.ChainID ***REMOVED***
	return ml.chainID
***REMOVED***

func (ml *mockLayer) DiffID() layer.DiffID ***REMOVED***
	return ml.diffID
***REMOVED***

func (ml *mockLayer) Parent() layer.Layer ***REMOVED***
	return ml.parent
***REMOVED***

func (ml *mockLayer) Size() (size int64, err error) ***REMOVED***
	return 0, nil
***REMOVED***

func (ml *mockLayer) DiffSize() (size int64, err error) ***REMOVED***
	return 0, nil
***REMOVED***

func (ml *mockLayer) Metadata() (map[string]string, error) ***REMOVED***
	return make(map[string]string), nil
***REMOVED***

type mockLayerStore struct ***REMOVED***
	layers map[layer.ChainID]*mockLayer
***REMOVED***

func createChainIDFromParent(parent layer.ChainID, dgsts ...layer.DiffID) layer.ChainID ***REMOVED***
	if len(dgsts) == 0 ***REMOVED***
		return parent
	***REMOVED***
	if parent == "" ***REMOVED***
		return createChainIDFromParent(layer.ChainID(dgsts[0]), dgsts[1:]...)
	***REMOVED***
	// H = "H(n-1) SHA256(n)"
	dgst := digest.FromBytes([]byte(string(parent) + " " + string(dgsts[0])))
	return createChainIDFromParent(layer.ChainID(dgst), dgsts[1:]...)
***REMOVED***

func (ls *mockLayerStore) Map() map[layer.ChainID]layer.Layer ***REMOVED***
	layers := map[layer.ChainID]layer.Layer***REMOVED******REMOVED***

	for k, v := range ls.layers ***REMOVED***
		layers[k] = v
	***REMOVED***

	return layers
***REMOVED***

func (ls *mockLayerStore) Register(reader io.Reader, parentID layer.ChainID) (layer.Layer, error) ***REMOVED***
	return ls.RegisterWithDescriptor(reader, parentID, distribution.Descriptor***REMOVED******REMOVED***)
***REMOVED***

func (ls *mockLayerStore) RegisterWithDescriptor(reader io.Reader, parentID layer.ChainID, _ distribution.Descriptor) (layer.Layer, error) ***REMOVED***
	var (
		parent layer.Layer
		err    error
	)

	if parentID != "" ***REMOVED***
		parent, err = ls.Get(parentID)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	l := &mockLayer***REMOVED***parent: parent***REMOVED***
	_, err = l.layerData.ReadFrom(reader)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	l.diffID = layer.DiffID(digest.FromBytes(l.layerData.Bytes()))
	l.chainID = createChainIDFromParent(parentID, l.diffID)

	ls.layers[l.chainID] = l
	return l, nil
***REMOVED***

func (ls *mockLayerStore) Get(chainID layer.ChainID) (layer.Layer, error) ***REMOVED***
	l, ok := ls.layers[chainID]
	if !ok ***REMOVED***
		return nil, layer.ErrLayerDoesNotExist
	***REMOVED***
	return l, nil
***REMOVED***

func (ls *mockLayerStore) Release(l layer.Layer) ([]layer.Metadata, error) ***REMOVED***
	return []layer.Metadata***REMOVED******REMOVED***, nil
***REMOVED***
func (ls *mockLayerStore) CreateRWLayer(string, layer.ChainID, *layer.CreateRWLayerOpts) (layer.RWLayer, error) ***REMOVED***
	return nil, errors.New("not implemented")
***REMOVED***

func (ls *mockLayerStore) GetRWLayer(string) (layer.RWLayer, error) ***REMOVED***
	return nil, errors.New("not implemented")
***REMOVED***

func (ls *mockLayerStore) ReleaseRWLayer(layer.RWLayer) ([]layer.Metadata, error) ***REMOVED***
	return nil, errors.New("not implemented")
***REMOVED***
func (ls *mockLayerStore) GetMountID(string) (string, error) ***REMOVED***
	return "", errors.New("not implemented")
***REMOVED***

func (ls *mockLayerStore) Cleanup() error ***REMOVED***
	return nil
***REMOVED***

func (ls *mockLayerStore) DriverStatus() [][2]string ***REMOVED***
	return [][2]string***REMOVED******REMOVED***
***REMOVED***

func (ls *mockLayerStore) DriverName() string ***REMOVED***
	return "mock"
***REMOVED***

type mockDownloadDescriptor struct ***REMOVED***
	currentDownloads *int32
	id               string
	diffID           layer.DiffID
	registeredDiffID layer.DiffID
	expectedDiffID   layer.DiffID
	simulateRetries  int
***REMOVED***

// Key returns the key used to deduplicate downloads.
func (d *mockDownloadDescriptor) Key() string ***REMOVED***
	return d.id
***REMOVED***

// ID returns the ID for display purposes.
func (d *mockDownloadDescriptor) ID() string ***REMOVED***
	return d.id
***REMOVED***

// DiffID should return the DiffID for this layer, or an error
// if it is unknown (for example, if it has not been downloaded
// before).
func (d *mockDownloadDescriptor) DiffID() (layer.DiffID, error) ***REMOVED***
	if d.diffID != "" ***REMOVED***
		return d.diffID, nil
	***REMOVED***
	return "", errors.New("no diffID available")
***REMOVED***

func (d *mockDownloadDescriptor) Registered(diffID layer.DiffID) ***REMOVED***
	d.registeredDiffID = diffID
***REMOVED***

func (d *mockDownloadDescriptor) mockTarStream() io.ReadCloser ***REMOVED***
	// The mock implementation returns the ID repeated 5 times as a tar
	// stream instead of actual tar data. The data is ignored except for
	// computing IDs.
	return ioutil.NopCloser(bytes.NewBuffer([]byte(d.id + d.id + d.id + d.id + d.id)))
***REMOVED***

// Download is called to perform the download.
func (d *mockDownloadDescriptor) Download(ctx context.Context, progressOutput progress.Output) (io.ReadCloser, int64, error) ***REMOVED***
	if d.currentDownloads != nil ***REMOVED***
		defer atomic.AddInt32(d.currentDownloads, -1)

		if atomic.AddInt32(d.currentDownloads, 1) > maxDownloadConcurrency ***REMOVED***
			return nil, 0, errors.New("concurrency limit exceeded")
		***REMOVED***
	***REMOVED***

	// Sleep a bit to simulate a time-consuming download.
	for i := int64(0); i <= 10; i++ ***REMOVED***
		select ***REMOVED***
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		case <-time.After(10 * time.Millisecond):
			progressOutput.WriteProgress(progress.Progress***REMOVED***ID: d.ID(), Action: "Downloading", Current: i, Total: 10***REMOVED***)
		***REMOVED***
	***REMOVED***

	if d.simulateRetries != 0 ***REMOVED***
		d.simulateRetries--
		return nil, 0, errors.New("simulating retry")
	***REMOVED***

	return d.mockTarStream(), 0, nil
***REMOVED***

func (d *mockDownloadDescriptor) Close() ***REMOVED***
***REMOVED***

func downloadDescriptors(currentDownloads *int32) []DownloadDescriptor ***REMOVED***
	return []DownloadDescriptor***REMOVED***
		&mockDownloadDescriptor***REMOVED***
			currentDownloads: currentDownloads,
			id:               "id1",
			expectedDiffID:   layer.DiffID("sha256:68e2c75dc5c78ea9240689c60d7599766c213ae210434c53af18470ae8c53ec1"),
		***REMOVED***,
		&mockDownloadDescriptor***REMOVED***
			currentDownloads: currentDownloads,
			id:               "id2",
			expectedDiffID:   layer.DiffID("sha256:64a636223116aa837973a5d9c2bdd17d9b204e4f95ac423e20e65dfbb3655473"),
		***REMOVED***,
		&mockDownloadDescriptor***REMOVED***
			currentDownloads: currentDownloads,
			id:               "id3",
			expectedDiffID:   layer.DiffID("sha256:58745a8bbd669c25213e9de578c4da5c8ee1c836b3581432c2b50e38a6753300"),
		***REMOVED***,
		&mockDownloadDescriptor***REMOVED***
			currentDownloads: currentDownloads,
			id:               "id2",
			expectedDiffID:   layer.DiffID("sha256:64a636223116aa837973a5d9c2bdd17d9b204e4f95ac423e20e65dfbb3655473"),
		***REMOVED***,
		&mockDownloadDescriptor***REMOVED***
			currentDownloads: currentDownloads,
			id:               "id4",
			expectedDiffID:   layer.DiffID("sha256:0dfb5b9577716cc173e95af7c10289322c29a6453a1718addc00c0c5b1330936"),
			simulateRetries:  1,
		***REMOVED***,
		&mockDownloadDescriptor***REMOVED***
			currentDownloads: currentDownloads,
			id:               "id5",
			expectedDiffID:   layer.DiffID("sha256:0a5f25fa1acbc647f6112a6276735d0fa01e4ee2aa7ec33015e337350e1ea23d"),
		***REMOVED***,
	***REMOVED***
***REMOVED***

func TestSuccessfulDownload(t *testing.T) ***REMOVED***
	// TODO Windows: Fix this unit text
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Needs fixing on Windows")
	***REMOVED***

	layerStore := &mockLayerStore***REMOVED***make(map[layer.ChainID]*mockLayer)***REMOVED***
	lsMap := make(map[string]layer.Store)
	lsMap[runtime.GOOS] = layerStore
	ldm := NewLayerDownloadManager(lsMap, maxDownloadConcurrency, func(m *LayerDownloadManager) ***REMOVED*** m.waitDuration = time.Millisecond ***REMOVED***)

	progressChan := make(chan progress.Progress)
	progressDone := make(chan struct***REMOVED******REMOVED***)
	receivedProgress := make(map[string]progress.Progress)

	go func() ***REMOVED***
		for p := range progressChan ***REMOVED***
			receivedProgress[p.ID] = p
		***REMOVED***
		close(progressDone)
	***REMOVED***()

	var currentDownloads int32
	descriptors := downloadDescriptors(&currentDownloads)

	firstDescriptor := descriptors[0].(*mockDownloadDescriptor)

	// Pre-register the first layer to simulate an already-existing layer
	l, err := layerStore.Register(firstDescriptor.mockTarStream(), "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	firstDescriptor.diffID = l.DiffID()

	rootFS, releaseFunc, err := ldm.Download(context.Background(), *image.NewRootFS(), runtime.GOOS, descriptors, progress.ChanOutput(progressChan))
	if err != nil ***REMOVED***
		t.Fatalf("download error: %v", err)
	***REMOVED***

	releaseFunc()

	close(progressChan)
	<-progressDone

	if len(rootFS.DiffIDs) != len(descriptors) ***REMOVED***
		t.Fatal("got wrong number of diffIDs in rootfs")
	***REMOVED***

	for i, d := range descriptors ***REMOVED***
		descriptor := d.(*mockDownloadDescriptor)

		if descriptor.diffID != "" ***REMOVED***
			if receivedProgress[d.ID()].Action != "Already exists" ***REMOVED***
				t.Fatalf("did not get 'Already exists' message for %v", d.ID())
			***REMOVED***
		***REMOVED*** else if receivedProgress[d.ID()].Action != "Pull complete" ***REMOVED***
			t.Fatalf("did not get 'Pull complete' message for %v", d.ID())
		***REMOVED***

		if rootFS.DiffIDs[i] != descriptor.expectedDiffID ***REMOVED***
			t.Fatalf("rootFS item %d has the wrong diffID (expected: %v got: %v)", i, descriptor.expectedDiffID, rootFS.DiffIDs[i])
		***REMOVED***

		if descriptor.diffID == "" && descriptor.registeredDiffID != rootFS.DiffIDs[i] ***REMOVED***
			t.Fatal("diffID mismatch between rootFS and Registered callback")
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestCancelledDownload(t *testing.T) ***REMOVED***
	layerStore := &mockLayerStore***REMOVED***make(map[layer.ChainID]*mockLayer)***REMOVED***
	lsMap := make(map[string]layer.Store)
	lsMap[runtime.GOOS] = layerStore
	ldm := NewLayerDownloadManager(lsMap, maxDownloadConcurrency, func(m *LayerDownloadManager) ***REMOVED*** m.waitDuration = time.Millisecond ***REMOVED***)
	progressChan := make(chan progress.Progress)
	progressDone := make(chan struct***REMOVED******REMOVED***)

	go func() ***REMOVED***
		for range progressChan ***REMOVED***
		***REMOVED***
		close(progressDone)
	***REMOVED***()

	ctx, cancel := context.WithCancel(context.Background())

	go func() ***REMOVED***
		<-time.After(time.Millisecond)
		cancel()
	***REMOVED***()

	descriptors := downloadDescriptors(nil)
	_, _, err := ldm.Download(ctx, *image.NewRootFS(), runtime.GOOS, descriptors, progress.ChanOutput(progressChan))
	if err != context.Canceled ***REMOVED***
		t.Fatal("expected download to be cancelled")
	***REMOVED***

	close(progressChan)
	<-progressDone
***REMOVED***
