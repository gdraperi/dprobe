package distribution

import (
	"fmt"
	"mime"

	"github.com/docker/distribution/context"
	"github.com/opencontainers/go-digest"
)

// Manifest represents a registry object specifying a set of
// references and an optional target
type Manifest interface ***REMOVED***
	// References returns a list of objects which make up this manifest.
	// A reference is anything which can be represented by a
	// distribution.Descriptor. These can consist of layers, resources or other
	// manifests.
	//
	// While no particular order is required, implementations should return
	// them from highest to lowest priority. For example, one might want to
	// return the base layer before the top layer.
	References() []Descriptor

	// Payload provides the serialized format of the manifest, in addition to
	// the media type.
	Payload() (mediaType string, payload []byte, err error)
***REMOVED***

// ManifestBuilder creates a manifest allowing one to include dependencies.
// Instances can be obtained from a version-specific manifest package.  Manifest
// specific data is passed into the function which creates the builder.
type ManifestBuilder interface ***REMOVED***
	// Build creates the manifest from his builder.
	Build(ctx context.Context) (Manifest, error)

	// References returns a list of objects which have been added to this
	// builder. The dependencies are returned in the order they were added,
	// which should be from base to head.
	References() []Descriptor

	// AppendReference includes the given object in the manifest after any
	// existing dependencies. If the add fails, such as when adding an
	// unsupported dependency, an error may be returned.
	//
	// The destination of the reference is dependent on the manifest type and
	// the dependency type.
	AppendReference(dependency Describable) error
***REMOVED***

// ManifestService describes operations on image manifests.
type ManifestService interface ***REMOVED***
	// Exists returns true if the manifest exists.
	Exists(ctx context.Context, dgst digest.Digest) (bool, error)

	// Get retrieves the manifest specified by the given digest
	Get(ctx context.Context, dgst digest.Digest, options ...ManifestServiceOption) (Manifest, error)

	// Put creates or updates the given manifest returning the manifest digest
	Put(ctx context.Context, manifest Manifest, options ...ManifestServiceOption) (digest.Digest, error)

	// Delete removes the manifest specified by the given digest. Deleting
	// a manifest that doesn't exist will return ErrManifestNotFound
	Delete(ctx context.Context, dgst digest.Digest) error
***REMOVED***

// ManifestEnumerator enables iterating over manifests
type ManifestEnumerator interface ***REMOVED***
	// Enumerate calls ingester for each manifest.
	Enumerate(ctx context.Context, ingester func(digest.Digest) error) error
***REMOVED***

// Describable is an interface for descriptors
type Describable interface ***REMOVED***
	Descriptor() Descriptor
***REMOVED***

// ManifestMediaTypes returns the supported media types for manifests.
func ManifestMediaTypes() (mediaTypes []string) ***REMOVED***
	for t := range mappings ***REMOVED***
		if t != "" ***REMOVED***
			mediaTypes = append(mediaTypes, t)
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// UnmarshalFunc implements manifest unmarshalling a given MediaType
type UnmarshalFunc func([]byte) (Manifest, Descriptor, error)

var mappings = make(map[string]UnmarshalFunc, 0)

// UnmarshalManifest looks up manifest unmarshal functions based on
// MediaType
func UnmarshalManifest(ctHeader string, p []byte) (Manifest, Descriptor, error) ***REMOVED***
	// Need to look up by the actual media type, not the raw contents of
	// the header. Strip semicolons and anything following them.
	var mediaType string
	if ctHeader != "" ***REMOVED***
		var err error
		mediaType, _, err = mime.ParseMediaType(ctHeader)
		if err != nil ***REMOVED***
			return nil, Descriptor***REMOVED******REMOVED***, err
		***REMOVED***
	***REMOVED***

	unmarshalFunc, ok := mappings[mediaType]
	if !ok ***REMOVED***
		unmarshalFunc, ok = mappings[""]
		if !ok ***REMOVED***
			return nil, Descriptor***REMOVED******REMOVED***, fmt.Errorf("unsupported manifest media type and no default available: %s", mediaType)
		***REMOVED***
	***REMOVED***

	return unmarshalFunc(p)
***REMOVED***

// RegisterManifestSchema registers an UnmarshalFunc for a given schema type.  This
// should be called from specific
func RegisterManifestSchema(mediaType string, u UnmarshalFunc) error ***REMOVED***
	if _, ok := mappings[mediaType]; ok ***REMOVED***
		return fmt.Errorf("manifest media type registration would overwrite existing: %s", mediaType)
	***REMOVED***
	mappings[mediaType] = u
	return nil
***REMOVED***
