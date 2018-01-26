package schema1

import (
	"fmt"

	"errors"
	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
	"github.com/docker/distribution/manifest"
	"github.com/docker/distribution/reference"
	"github.com/docker/libtrust"
	"github.com/opencontainers/go-digest"
)

// referenceManifestBuilder is a type for constructing manifests from schema1
// dependencies.
type referenceManifestBuilder struct ***REMOVED***
	Manifest
	pk libtrust.PrivateKey
***REMOVED***

// NewReferenceManifestBuilder is used to build new manifests for the current
// schema version using schema1 dependencies.
func NewReferenceManifestBuilder(pk libtrust.PrivateKey, ref reference.Named, architecture string) distribution.ManifestBuilder ***REMOVED***
	tag := ""
	if tagged, isTagged := ref.(reference.Tagged); isTagged ***REMOVED***
		tag = tagged.Tag()
	***REMOVED***

	return &referenceManifestBuilder***REMOVED***
		Manifest: Manifest***REMOVED***
			Versioned: manifest.Versioned***REMOVED***
				SchemaVersion: 1,
			***REMOVED***,
			Name:         ref.Name(),
			Tag:          tag,
			Architecture: architecture,
		***REMOVED***,
		pk: pk,
	***REMOVED***
***REMOVED***

func (mb *referenceManifestBuilder) Build(ctx context.Context) (distribution.Manifest, error) ***REMOVED***
	m := mb.Manifest
	if len(m.FSLayers) == 0 ***REMOVED***
		return nil, errors.New("cannot build manifest with zero layers or history")
	***REMOVED***

	m.FSLayers = make([]FSLayer, len(mb.Manifest.FSLayers))
	m.History = make([]History, len(mb.Manifest.History))
	copy(m.FSLayers, mb.Manifest.FSLayers)
	copy(m.History, mb.Manifest.History)

	return Sign(&m, mb.pk)
***REMOVED***

// AppendReference adds a reference to the current ManifestBuilder
func (mb *referenceManifestBuilder) AppendReference(d distribution.Describable) error ***REMOVED***
	r, ok := d.(Reference)
	if !ok ***REMOVED***
		return fmt.Errorf("Unable to add non-reference type to v1 builder")
	***REMOVED***

	// Entries need to be prepended
	mb.Manifest.FSLayers = append([]FSLayer***REMOVED******REMOVED***BlobSum: r.Digest***REMOVED******REMOVED***, mb.Manifest.FSLayers...)
	mb.Manifest.History = append([]History***REMOVED***r.History***REMOVED***, mb.Manifest.History...)
	return nil

***REMOVED***

// References returns the current references added to this builder
func (mb *referenceManifestBuilder) References() []distribution.Descriptor ***REMOVED***
	refs := make([]distribution.Descriptor, len(mb.Manifest.FSLayers))
	for i := range mb.Manifest.FSLayers ***REMOVED***
		layerDigest := mb.Manifest.FSLayers[i].BlobSum
		history := mb.Manifest.History[i]
		ref := Reference***REMOVED***layerDigest, 0, history***REMOVED***
		refs[i] = ref.Descriptor()
	***REMOVED***
	return refs
***REMOVED***

// Reference describes a manifest v2, schema version 1 dependency.
// An FSLayer associated with a history entry.
type Reference struct ***REMOVED***
	Digest  digest.Digest
	Size    int64 // if we know it, set it for the descriptor.
	History History
***REMOVED***

// Descriptor describes a reference
func (r Reference) Descriptor() distribution.Descriptor ***REMOVED***
	return distribution.Descriptor***REMOVED***
		MediaType: MediaTypeManifestLayer,
		Digest:    r.Digest,
		Size:      r.Size,
	***REMOVED***
***REMOVED***
