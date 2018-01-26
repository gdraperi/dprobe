package schema1

import (
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
	"github.com/docker/distribution/manifest"
	"github.com/docker/distribution/reference"
	"github.com/docker/libtrust"
	"github.com/opencontainers/go-digest"
)

type diffID digest.Digest

// gzippedEmptyTar is a gzip-compressed version of an empty tar file
// (1024 NULL bytes)
var gzippedEmptyTar = []byte***REMOVED***
	31, 139, 8, 0, 0, 9, 110, 136, 0, 255, 98, 24, 5, 163, 96, 20, 140, 88,
	0, 8, 0, 0, 255, 255, 46, 175, 181, 239, 0, 4, 0, 0,
***REMOVED***

// digestSHA256GzippedEmptyTar is the canonical sha256 digest of
// gzippedEmptyTar
const digestSHA256GzippedEmptyTar = digest.Digest("sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4")

// configManifestBuilder is a type for constructing manifests from an image
// configuration and generic descriptors.
type configManifestBuilder struct ***REMOVED***
	// bs is a BlobService used to create empty layer tars in the
	// blob store if necessary.
	bs distribution.BlobService
	// pk is the libtrust private key used to sign the final manifest.
	pk libtrust.PrivateKey
	// configJSON is configuration supplied when the ManifestBuilder was
	// created.
	configJSON []byte
	// ref contains the name and optional tag provided to NewConfigManifestBuilder.
	ref reference.Named
	// descriptors is the set of descriptors referencing the layers.
	descriptors []distribution.Descriptor
	// emptyTarDigest is set to a valid digest if an empty tar has been
	// put in the blob store; otherwise it is empty.
	emptyTarDigest digest.Digest
***REMOVED***

// NewConfigManifestBuilder is used to build new manifests for the current
// schema version from an image configuration and a set of descriptors.
// It takes a BlobService so that it can add an empty tar to the blob store
// if the resulting manifest needs empty layers.
func NewConfigManifestBuilder(bs distribution.BlobService, pk libtrust.PrivateKey, ref reference.Named, configJSON []byte) distribution.ManifestBuilder ***REMOVED***
	return &configManifestBuilder***REMOVED***
		bs:         bs,
		pk:         pk,
		configJSON: configJSON,
		ref:        ref,
	***REMOVED***
***REMOVED***

// Build produces a final manifest from the given references
func (mb *configManifestBuilder) Build(ctx context.Context) (m distribution.Manifest, err error) ***REMOVED***
	type imageRootFS struct ***REMOVED***
		Type      string   `json:"type"`
		DiffIDs   []diffID `json:"diff_ids,omitempty"`
		BaseLayer string   `json:"base_layer,omitempty"`
	***REMOVED***

	type imageHistory struct ***REMOVED***
		Created    time.Time `json:"created"`
		Author     string    `json:"author,omitempty"`
		CreatedBy  string    `json:"created_by,omitempty"`
		Comment    string    `json:"comment,omitempty"`
		EmptyLayer bool      `json:"empty_layer,omitempty"`
	***REMOVED***

	type imageConfig struct ***REMOVED***
		RootFS       *imageRootFS   `json:"rootfs,omitempty"`
		History      []imageHistory `json:"history,omitempty"`
		Architecture string         `json:"architecture,omitempty"`
	***REMOVED***

	var img imageConfig

	if err := json.Unmarshal(mb.configJSON, &img); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if len(img.History) == 0 ***REMOVED***
		return nil, errors.New("empty history when trying to create schema1 manifest")
	***REMOVED***

	if len(img.RootFS.DiffIDs) != len(mb.descriptors) ***REMOVED***
		return nil, fmt.Errorf("number of descriptors and number of layers in rootfs must match: len(%v) != len(%v)", img.RootFS.DiffIDs, mb.descriptors)
	***REMOVED***

	// Generate IDs for each layer
	// For non-top-level layers, create fake V1Compatibility strings that
	// fit the format and don't collide with anything else, but don't
	// result in runnable images on their own.
	type v1Compatibility struct ***REMOVED***
		ID              string    `json:"id"`
		Parent          string    `json:"parent,omitempty"`
		Comment         string    `json:"comment,omitempty"`
		Created         time.Time `json:"created"`
		ContainerConfig struct ***REMOVED***
			Cmd []string
		***REMOVED*** `json:"container_config,omitempty"`
		Author    string `json:"author,omitempty"`
		ThrowAway bool   `json:"throwaway,omitempty"`
	***REMOVED***

	fsLayerList := make([]FSLayer, len(img.History))
	history := make([]History, len(img.History))

	parent := ""
	layerCounter := 0
	for i, h := range img.History[:len(img.History)-1] ***REMOVED***
		var blobsum digest.Digest
		if h.EmptyLayer ***REMOVED***
			if blobsum, err = mb.emptyTar(ctx); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if len(img.RootFS.DiffIDs) <= layerCounter ***REMOVED***
				return nil, errors.New("too many non-empty layers in History section")
			***REMOVED***
			blobsum = mb.descriptors[layerCounter].Digest
			layerCounter++
		***REMOVED***

		v1ID := digest.FromBytes([]byte(blobsum.Hex() + " " + parent)).Hex()

		if i == 0 && img.RootFS.BaseLayer != "" ***REMOVED***
			// windows-only baselayer setup
			baseID := sha512.Sum384([]byte(img.RootFS.BaseLayer))
			parent = fmt.Sprintf("%x", baseID[:32])
		***REMOVED***

		v1Compatibility := v1Compatibility***REMOVED***
			ID:      v1ID,
			Parent:  parent,
			Comment: h.Comment,
			Created: h.Created,
			Author:  h.Author,
		***REMOVED***
		v1Compatibility.ContainerConfig.Cmd = []string***REMOVED***img.History[i].CreatedBy***REMOVED***
		if h.EmptyLayer ***REMOVED***
			v1Compatibility.ThrowAway = true
		***REMOVED***
		jsonBytes, err := json.Marshal(&v1Compatibility)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		reversedIndex := len(img.History) - i - 1
		history[reversedIndex].V1Compatibility = string(jsonBytes)
		fsLayerList[reversedIndex] = FSLayer***REMOVED***BlobSum: blobsum***REMOVED***

		parent = v1ID
	***REMOVED***

	latestHistory := img.History[len(img.History)-1]

	var blobsum digest.Digest
	if latestHistory.EmptyLayer ***REMOVED***
		if blobsum, err = mb.emptyTar(ctx); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if len(img.RootFS.DiffIDs) <= layerCounter ***REMOVED***
			return nil, errors.New("too many non-empty layers in History section")
		***REMOVED***
		blobsum = mb.descriptors[layerCounter].Digest
	***REMOVED***

	fsLayerList[0] = FSLayer***REMOVED***BlobSum: blobsum***REMOVED***
	dgst := digest.FromBytes([]byte(blobsum.Hex() + " " + parent + " " + string(mb.configJSON)))

	// Top-level v1compatibility string should be a modified version of the
	// image config.
	transformedConfig, err := MakeV1ConfigFromConfig(mb.configJSON, dgst.Hex(), parent, latestHistory.EmptyLayer)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	history[0].V1Compatibility = string(transformedConfig)

	tag := ""
	if tagged, isTagged := mb.ref.(reference.Tagged); isTagged ***REMOVED***
		tag = tagged.Tag()
	***REMOVED***

	mfst := Manifest***REMOVED***
		Versioned: manifest.Versioned***REMOVED***
			SchemaVersion: 1,
		***REMOVED***,
		Name:         mb.ref.Name(),
		Tag:          tag,
		Architecture: img.Architecture,
		FSLayers:     fsLayerList,
		History:      history,
	***REMOVED***

	return Sign(&mfst, mb.pk)
***REMOVED***

// emptyTar pushes a compressed empty tar to the blob store if one doesn't
// already exist, and returns its blobsum.
func (mb *configManifestBuilder) emptyTar(ctx context.Context) (digest.Digest, error) ***REMOVED***
	if mb.emptyTarDigest != "" ***REMOVED***
		// Already put an empty tar
		return mb.emptyTarDigest, nil
	***REMOVED***

	descriptor, err := mb.bs.Stat(ctx, digestSHA256GzippedEmptyTar)
	switch err ***REMOVED***
	case nil:
		mb.emptyTarDigest = descriptor.Digest
		return descriptor.Digest, nil
	case distribution.ErrBlobUnknown:
		// nop
	default:
		return "", err
	***REMOVED***

	// Add gzipped empty tar to the blob store
	descriptor, err = mb.bs.Put(ctx, "", gzippedEmptyTar)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	mb.emptyTarDigest = descriptor.Digest

	return descriptor.Digest, nil
***REMOVED***

// AppendReference adds a reference to the current ManifestBuilder
func (mb *configManifestBuilder) AppendReference(d distribution.Describable) error ***REMOVED***
	descriptor := d.Descriptor()

	if err := descriptor.Digest.Validate(); err != nil ***REMOVED***
		return err
	***REMOVED***

	mb.descriptors = append(mb.descriptors, descriptor)
	return nil
***REMOVED***

// References returns the current references added to this builder
func (mb *configManifestBuilder) References() []distribution.Descriptor ***REMOVED***
	return mb.descriptors
***REMOVED***

// MakeV1ConfigFromConfig creates an legacy V1 image config from image config JSON
func MakeV1ConfigFromConfig(configJSON []byte, v1ID, parentV1ID string, throwaway bool) ([]byte, error) ***REMOVED***
	// Top-level v1compatibility string should be a modified version of the
	// image config.
	var configAsMap map[string]*json.RawMessage
	if err := json.Unmarshal(configJSON, &configAsMap); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Delete fields that didn't exist in old manifest
	delete(configAsMap, "rootfs")
	delete(configAsMap, "history")
	configAsMap["id"] = rawJSON(v1ID)
	if parentV1ID != "" ***REMOVED***
		configAsMap["parent"] = rawJSON(parentV1ID)
	***REMOVED***
	if throwaway ***REMOVED***
		configAsMap["throwaway"] = rawJSON(true)
	***REMOVED***

	return json.Marshal(configAsMap)
***REMOVED***

func rawJSON(value interface***REMOVED******REMOVED***) *json.RawMessage ***REMOVED***
	jsonval, err := json.Marshal(value)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	return (*json.RawMessage)(&jsonval)
***REMOVED***
