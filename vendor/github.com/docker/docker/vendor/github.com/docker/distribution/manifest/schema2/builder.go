package schema2

import (
	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
	"github.com/opencontainers/go-digest"
)

// builder is a type for constructing manifests.
type builder struct ***REMOVED***
	// bs is a BlobService used to publish the configuration blob.
	bs distribution.BlobService

	// configMediaType is media type used to describe configuration
	configMediaType string

	// configJSON references
	configJSON []byte

	// dependencies is a list of descriptors that gets built by successive
	// calls to AppendReference. In case of image configuration these are layers.
	dependencies []distribution.Descriptor
***REMOVED***

// NewManifestBuilder is used to build new manifests for the current schema
// version. It takes a BlobService so it can publish the configuration blob
// as part of the Build process.
func NewManifestBuilder(bs distribution.BlobService, configMediaType string, configJSON []byte) distribution.ManifestBuilder ***REMOVED***
	mb := &builder***REMOVED***
		bs:              bs,
		configMediaType: configMediaType,
		configJSON:      make([]byte, len(configJSON)),
	***REMOVED***
	copy(mb.configJSON, configJSON)

	return mb
***REMOVED***

// Build produces a final manifest from the given references.
func (mb *builder) Build(ctx context.Context) (distribution.Manifest, error) ***REMOVED***
	m := Manifest***REMOVED***
		Versioned: SchemaVersion,
		Layers:    make([]distribution.Descriptor, len(mb.dependencies)),
	***REMOVED***
	copy(m.Layers, mb.dependencies)

	configDigest := digest.FromBytes(mb.configJSON)

	var err error
	m.Config, err = mb.bs.Stat(ctx, configDigest)
	switch err ***REMOVED***
	case nil:
		// Override MediaType, since Put always replaces the specified media
		// type with application/octet-stream in the descriptor it returns.
		m.Config.MediaType = mb.configMediaType
		return FromStruct(m)
	case distribution.ErrBlobUnknown:
		// nop
	default:
		return nil, err
	***REMOVED***

	// Add config to the blob store
	m.Config, err = mb.bs.Put(ctx, mb.configMediaType, mb.configJSON)
	// Override MediaType, since Put always replaces the specified media
	// type with application/octet-stream in the descriptor it returns.
	m.Config.MediaType = mb.configMediaType
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return FromStruct(m)
***REMOVED***

// AppendReference adds a reference to the current ManifestBuilder.
func (mb *builder) AppendReference(d distribution.Describable) error ***REMOVED***
	mb.dependencies = append(mb.dependencies, d.Descriptor())
	return nil
***REMOVED***

// References returns the current references added to this builder.
func (mb *builder) References() []distribution.Descriptor ***REMOVED***
	return mb.dependencies
***REMOVED***
