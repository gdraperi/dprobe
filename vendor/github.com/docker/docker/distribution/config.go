package distribution

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"

	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/distribution/metadata"
	"github.com/docker/docker/distribution/xfer"
	"github.com/docker/docker/image"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/system"
	refstore "github.com/docker/docker/reference"
	"github.com/docker/docker/registry"
	"github.com/docker/libtrust"
	"github.com/opencontainers/go-digest"
	"golang.org/x/net/context"
)

// Config stores configuration for communicating
// with a registry.
type Config struct ***REMOVED***
	// MetaHeaders stores HTTP headers with metadata about the image
	MetaHeaders map[string][]string
	// AuthConfig holds authentication credentials for authenticating with
	// the registry.
	AuthConfig *types.AuthConfig
	// ProgressOutput is the interface for showing the status of the pull
	// operation.
	ProgressOutput progress.Output
	// RegistryService is the registry service to use for TLS configuration
	// and endpoint lookup.
	RegistryService registry.Service
	// ImageEventLogger notifies events for a given image
	ImageEventLogger func(id, name, action string)
	// MetadataStore is the storage backend for distribution-specific
	// metadata.
	MetadataStore metadata.Store
	// ImageStore manages images.
	ImageStore ImageConfigStore
	// ReferenceStore manages tags. This value is optional, when excluded
	// content will not be tagged.
	ReferenceStore refstore.Store
	// RequireSchema2 ensures that only schema2 manifests are used.
	RequireSchema2 bool
***REMOVED***

// ImagePullConfig stores pull configuration.
type ImagePullConfig struct ***REMOVED***
	Config

	// DownloadManager manages concurrent pulls.
	DownloadManager RootFSDownloadManager
	// Schema2Types is the valid schema2 configuration types allowed
	// by the pull operation.
	Schema2Types []string
	// OS is the requested operating system of the image being pulled to ensure it can be validated
	// when the host OS supports multiple image operating systems.
	OS string
***REMOVED***

// ImagePushConfig stores push configuration.
type ImagePushConfig struct ***REMOVED***
	Config

	// ConfigMediaType is the configuration media type for
	// schema2 manifests.
	ConfigMediaType string
	// LayerStores (indexed by operating system) manages layers.
	LayerStores map[string]PushLayerProvider
	// TrustKey is the private key for legacy signatures. This is typically
	// an ephemeral key, since these signatures are no longer verified.
	TrustKey libtrust.PrivateKey
	// UploadManager dispatches uploads.
	UploadManager *xfer.LayerUploadManager
***REMOVED***

// ImageConfigStore handles storing and getting image configurations
// by digest. Allows getting an image configurations rootfs from the
// configuration.
type ImageConfigStore interface ***REMOVED***
	Put([]byte) (digest.Digest, error)
	Get(digest.Digest) ([]byte, error)
	RootFSAndOSFromConfig([]byte) (*image.RootFS, string, error)
***REMOVED***

// PushLayerProvider provides layers to be pushed by ChainID.
type PushLayerProvider interface ***REMOVED***
	Get(layer.ChainID) (PushLayer, error)
***REMOVED***

// PushLayer is a pushable layer with metadata about the layer
// and access to the content of the layer.
type PushLayer interface ***REMOVED***
	ChainID() layer.ChainID
	DiffID() layer.DiffID
	Parent() PushLayer
	Open() (io.ReadCloser, error)
	Size() (int64, error)
	MediaType() string
	Release()
***REMOVED***

// RootFSDownloadManager handles downloading of the rootfs
type RootFSDownloadManager interface ***REMOVED***
	// Download downloads the layers into the given initial rootfs and
	// returns the final rootfs.
	// Given progress output to track download progress
	// Returns function to release download resources
	Download(ctx context.Context, initialRootFS image.RootFS, os string, layers []xfer.DownloadDescriptor, progressOutput progress.Output) (image.RootFS, func(), error)
***REMOVED***

type imageConfigStore struct ***REMOVED***
	image.Store
***REMOVED***

// NewImageConfigStoreFromStore returns an ImageConfigStore backed
// by an image.Store for container images.
func NewImageConfigStoreFromStore(is image.Store) ImageConfigStore ***REMOVED***
	return &imageConfigStore***REMOVED***
		Store: is,
	***REMOVED***
***REMOVED***

func (s *imageConfigStore) Put(c []byte) (digest.Digest, error) ***REMOVED***
	id, err := s.Store.Create(c)
	return digest.Digest(id), err
***REMOVED***

func (s *imageConfigStore) Get(d digest.Digest) ([]byte, error) ***REMOVED***
	img, err := s.Store.Get(image.IDFromDigest(d))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return img.RawJSON(), nil
***REMOVED***

func (s *imageConfigStore) RootFSAndOSFromConfig(c []byte) (*image.RootFS, string, error) ***REMOVED***
	var unmarshalledConfig image.Image
	if err := json.Unmarshal(c, &unmarshalledConfig); err != nil ***REMOVED***
		return nil, "", err
	***REMOVED***

	// fail immediately on Windows when downloading a non-Windows image
	// and vice versa. Exception on Windows if Linux Containers are enabled.
	if runtime.GOOS == "windows" && unmarshalledConfig.OS == "linux" && !system.LCOWSupported() ***REMOVED***
		return nil, "", fmt.Errorf("image operating system %q cannot be used on this platform", unmarshalledConfig.OS)
	***REMOVED*** else if runtime.GOOS != "windows" && unmarshalledConfig.OS == "windows" ***REMOVED***
		return nil, "", fmt.Errorf("image operating system %q cannot be used on this platform", unmarshalledConfig.OS)
	***REMOVED***

	os := unmarshalledConfig.OS
	if os == "" ***REMOVED***
		os = runtime.GOOS
	***REMOVED***
	if !system.IsOSSupported(os) ***REMOVED***
		return nil, "", system.ErrNotSupportedOperatingSystem
	***REMOVED***
	return unmarshalledConfig.RootFS, os, nil
***REMOVED***

type storeLayerProvider struct ***REMOVED***
	ls layer.Store
***REMOVED***

// NewLayerProvidersFromStores returns layer providers backed by
// an instance of LayerStore. Only getting layers as gzipped
// tars is supported.
func NewLayerProvidersFromStores(lss map[string]layer.Store) map[string]PushLayerProvider ***REMOVED***
	plps := make(map[string]PushLayerProvider)
	for os, ls := range lss ***REMOVED***
		plps[os] = &storeLayerProvider***REMOVED***ls: ls***REMOVED***
	***REMOVED***
	return plps
***REMOVED***

func (p *storeLayerProvider) Get(lid layer.ChainID) (PushLayer, error) ***REMOVED***
	if lid == "" ***REMOVED***
		return &storeLayer***REMOVED***
			Layer: layer.EmptyLayer,
		***REMOVED***, nil
	***REMOVED***
	l, err := p.ls.Get(lid)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	sl := storeLayer***REMOVED***
		Layer: l,
		ls:    p.ls,
	***REMOVED***
	if d, ok := l.(distribution.Describable); ok ***REMOVED***
		return &describableStoreLayer***REMOVED***
			storeLayer:  sl,
			describable: d,
		***REMOVED***, nil
	***REMOVED***

	return &sl, nil
***REMOVED***

type storeLayer struct ***REMOVED***
	layer.Layer
	ls layer.Store
***REMOVED***

func (l *storeLayer) Parent() PushLayer ***REMOVED***
	p := l.Layer.Parent()
	if p == nil ***REMOVED***
		return nil
	***REMOVED***
	sl := storeLayer***REMOVED***
		Layer: p,
		ls:    l.ls,
	***REMOVED***
	if d, ok := p.(distribution.Describable); ok ***REMOVED***
		return &describableStoreLayer***REMOVED***
			storeLayer:  sl,
			describable: d,
		***REMOVED***
	***REMOVED***

	return &sl
***REMOVED***

func (l *storeLayer) Open() (io.ReadCloser, error) ***REMOVED***
	return l.Layer.TarStream()
***REMOVED***

func (l *storeLayer) Size() (int64, error) ***REMOVED***
	return l.Layer.DiffSize()
***REMOVED***

func (l *storeLayer) MediaType() string ***REMOVED***
	// layer store always returns uncompressed tars
	return schema2.MediaTypeUncompressedLayer
***REMOVED***

func (l *storeLayer) Release() ***REMOVED***
	if l.ls != nil ***REMOVED***
		layer.ReleaseAndLog(l.ls, l.Layer)
	***REMOVED***
***REMOVED***

type describableStoreLayer struct ***REMOVED***
	storeLayer
	describable distribution.Describable
***REMOVED***

func (l *describableStoreLayer) Descriptor() distribution.Descriptor ***REMOVED***
	return l.describable.Descriptor()
***REMOVED***
