package dockerfile

import (
	"encoding/json"
	"io"
	"runtime"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/builder"
	containerpkg "github.com/docker/docker/container"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/containerfs"
	"golang.org/x/net/context"
)

// MockBackend implements the builder.Backend interface for unit testing
type MockBackend struct ***REMOVED***
	containerCreateFunc func(config types.ContainerCreateConfig) (container.ContainerCreateCreatedBody, error)
	commitFunc          func(string, *backend.ContainerCommitConfig) (string, error)
	getImageFunc        func(string) (builder.Image, builder.ReleaseableLayer, error)
	makeImageCacheFunc  func(cacheFrom []string) builder.ImageCache
***REMOVED***

func (m *MockBackend) ContainerAttachRaw(cID string, stdin io.ReadCloser, stdout, stderr io.Writer, stream bool, attached chan struct***REMOVED******REMOVED***) error ***REMOVED***
	return nil
***REMOVED***

func (m *MockBackend) ContainerCreate(config types.ContainerCreateConfig) (container.ContainerCreateCreatedBody, error) ***REMOVED***
	if m.containerCreateFunc != nil ***REMOVED***
		return m.containerCreateFunc(config)
	***REMOVED***
	return container.ContainerCreateCreatedBody***REMOVED******REMOVED***, nil
***REMOVED***

func (m *MockBackend) ContainerRm(name string, config *types.ContainerRmConfig) error ***REMOVED***
	return nil
***REMOVED***

func (m *MockBackend) Commit(cID string, cfg *backend.ContainerCommitConfig) (string, error) ***REMOVED***
	if m.commitFunc != nil ***REMOVED***
		return m.commitFunc(cID, cfg)
	***REMOVED***
	return "", nil
***REMOVED***

func (m *MockBackend) ContainerKill(containerID string, sig uint64) error ***REMOVED***
	return nil
***REMOVED***

func (m *MockBackend) ContainerStart(containerID string, hostConfig *container.HostConfig, checkpoint string, checkpointDir string) error ***REMOVED***
	return nil
***REMOVED***

func (m *MockBackend) ContainerWait(ctx context.Context, containerID string, condition containerpkg.WaitCondition) (<-chan containerpkg.StateStatus, error) ***REMOVED***
	return nil, nil
***REMOVED***

func (m *MockBackend) ContainerCreateWorkdir(containerID string) error ***REMOVED***
	return nil
***REMOVED***

func (m *MockBackend) CopyOnBuild(containerID string, destPath string, srcRoot string, srcPath string, decompress bool) error ***REMOVED***
	return nil
***REMOVED***

func (m *MockBackend) GetImageAndReleasableLayer(ctx context.Context, refOrID string, opts backend.GetImageAndLayerOptions) (builder.Image, builder.ReleaseableLayer, error) ***REMOVED***
	if m.getImageFunc != nil ***REMOVED***
		return m.getImageFunc(refOrID)
	***REMOVED***

	return &mockImage***REMOVED***id: "theid"***REMOVED***, &mockLayer***REMOVED******REMOVED***, nil
***REMOVED***

func (m *MockBackend) MakeImageCache(cacheFrom []string) builder.ImageCache ***REMOVED***
	if m.makeImageCacheFunc != nil ***REMOVED***
		return m.makeImageCacheFunc(cacheFrom)
	***REMOVED***
	return nil
***REMOVED***

func (m *MockBackend) CreateImage(config []byte, parent string) (builder.Image, error) ***REMOVED***
	return nil, nil
***REMOVED***

type mockImage struct ***REMOVED***
	id     string
	config *container.Config
***REMOVED***

func (i *mockImage) ImageID() string ***REMOVED***
	return i.id
***REMOVED***

func (i *mockImage) RunConfig() *container.Config ***REMOVED***
	return i.config
***REMOVED***

func (i *mockImage) OperatingSystem() string ***REMOVED***
	return runtime.GOOS
***REMOVED***

func (i *mockImage) MarshalJSON() ([]byte, error) ***REMOVED***
	type rawImage mockImage
	return json.Marshal(rawImage(*i))
***REMOVED***

type mockImageCache struct ***REMOVED***
	getCacheFunc func(parentID string, cfg *container.Config) (string, error)
***REMOVED***

func (mic *mockImageCache) GetCache(parentID string, cfg *container.Config) (string, error) ***REMOVED***
	if mic.getCacheFunc != nil ***REMOVED***
		return mic.getCacheFunc(parentID, cfg)
	***REMOVED***
	return "", nil
***REMOVED***

type mockLayer struct***REMOVED******REMOVED***

func (l *mockLayer) Release() error ***REMOVED***
	return nil
***REMOVED***

func (l *mockLayer) Mount() (containerfs.ContainerFS, error) ***REMOVED***
	return containerfs.NewLocalContainerFS("mountPath"), nil
***REMOVED***

func (l *mockLayer) Commit() (builder.ReleaseableLayer, error) ***REMOVED***
	return nil, nil
***REMOVED***

func (l *mockLayer) DiffID() layer.DiffID ***REMOVED***
	return layer.DiffID("abcdef")
***REMOVED***
