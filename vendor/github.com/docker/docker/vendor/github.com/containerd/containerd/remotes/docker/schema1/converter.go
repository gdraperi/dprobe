package schema1

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/remotes"
	digest "github.com/opencontainers/go-digest"
	specs "github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

const manifestSizeLimit = 8e6 // 8MB

type blobState struct ***REMOVED***
	diffID digest.Digest
	empty  bool
***REMOVED***

// Converter converts schema1 manifests to schema2 on fetch
type Converter struct ***REMOVED***
	contentStore content.Store
	fetcher      remotes.Fetcher

	pulledManifest *manifest

	mu         sync.Mutex
	blobMap    map[digest.Digest]blobState
	layerBlobs map[digest.Digest]ocispec.Descriptor
***REMOVED***

// NewConverter returns a new converter
func NewConverter(contentStore content.Store, fetcher remotes.Fetcher) *Converter ***REMOVED***
	return &Converter***REMOVED***
		contentStore: contentStore,
		fetcher:      fetcher,
		blobMap:      map[digest.Digest]blobState***REMOVED******REMOVED***,
		layerBlobs:   map[digest.Digest]ocispec.Descriptor***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

// Handle fetching descriptors for a docker media type
func (c *Converter) Handle(ctx context.Context, desc ocispec.Descriptor) ([]ocispec.Descriptor, error) ***REMOVED***
	switch desc.MediaType ***REMOVED***
	case images.MediaTypeDockerSchema1Manifest:
		if err := c.fetchManifest(ctx, desc); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		m := c.pulledManifest
		if len(m.FSLayers) != len(m.History) ***REMOVED***
			return nil, errors.New("invalid schema 1 manifest, history and layer mismatch")
		***REMOVED***
		descs := make([]ocispec.Descriptor, 0, len(c.pulledManifest.FSLayers))

		for i := range m.FSLayers ***REMOVED***
			if _, ok := c.blobMap[c.pulledManifest.FSLayers[i].BlobSum]; !ok ***REMOVED***
				empty, err := isEmptyLayer([]byte(m.History[i].V1Compatibility))
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***

				// Do no attempt to download a known empty blob
				if !empty ***REMOVED***
					descs = append([]ocispec.Descriptor***REMOVED***
						***REMOVED***
							MediaType: images.MediaTypeDockerSchema2LayerGzip,
							Digest:    c.pulledManifest.FSLayers[i].BlobSum,
							Size:      -1,
						***REMOVED***,
					***REMOVED***, descs...)
				***REMOVED***
				c.blobMap[c.pulledManifest.FSLayers[i].BlobSum] = blobState***REMOVED***
					empty: empty,
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return descs, nil
	case images.MediaTypeDockerSchema2LayerGzip:
		if c.pulledManifest == nil ***REMOVED***
			return nil, errors.New("manifest required for schema 1 blob pull")
		***REMOVED***
		return nil, c.fetchBlob(ctx, desc)
	default:
		return nil, fmt.Errorf("%v not support for schema 1 manifests", desc.MediaType)
	***REMOVED***
***REMOVED***

// Convert a docker manifest to an OCI descriptor
func (c *Converter) Convert(ctx context.Context) (ocispec.Descriptor, error) ***REMOVED***
	history, diffIDs, err := c.schema1ManifestHistory()
	if err != nil ***REMOVED***
		return ocispec.Descriptor***REMOVED******REMOVED***, errors.Wrap(err, "schema 1 conversion failed")
	***REMOVED***

	var img ocispec.Image
	if err := json.Unmarshal([]byte(c.pulledManifest.History[0].V1Compatibility), &img); err != nil ***REMOVED***
		return ocispec.Descriptor***REMOVED******REMOVED***, errors.Wrap(err, "failed to unmarshal image from schema 1 history")
	***REMOVED***

	img.History = history
	img.RootFS = ocispec.RootFS***REMOVED***
		Type:    "layers",
		DiffIDs: diffIDs,
	***REMOVED***

	b, err := json.Marshal(img)
	if err != nil ***REMOVED***
		return ocispec.Descriptor***REMOVED******REMOVED***, errors.Wrap(err, "failed to marshal image")
	***REMOVED***

	config := ocispec.Descriptor***REMOVED***
		MediaType: ocispec.MediaTypeImageConfig,
		Digest:    digest.Canonical.FromBytes(b),
		Size:      int64(len(b)),
	***REMOVED***

	layers := make([]ocispec.Descriptor, len(diffIDs))
	for i, diffID := range diffIDs ***REMOVED***
		layers[i] = c.layerBlobs[diffID]
	***REMOVED***

	manifest := ocispec.Manifest***REMOVED***
		Versioned: specs.Versioned***REMOVED***
			SchemaVersion: 2,
		***REMOVED***,
		Config: config,
		Layers: layers,
	***REMOVED***

	mb, err := json.Marshal(manifest)
	if err != nil ***REMOVED***
		return ocispec.Descriptor***REMOVED******REMOVED***, errors.Wrap(err, "failed to marshal image")
	***REMOVED***

	desc := ocispec.Descriptor***REMOVED***
		MediaType: ocispec.MediaTypeImageManifest,
		Digest:    digest.Canonical.FromBytes(mb),
		Size:      int64(len(mb)),
	***REMOVED***

	labels := map[string]string***REMOVED******REMOVED***
	labels["containerd.io/gc.ref.content.0"] = manifest.Config.Digest.String()
	for i, ch := range manifest.Layers ***REMOVED***
		labels[fmt.Sprintf("containerd.io/gc.ref.content.%d", i+1)] = ch.Digest.String()
	***REMOVED***

	ref := remotes.MakeRefKey(ctx, desc)
	if err := content.WriteBlob(ctx, c.contentStore, ref, bytes.NewReader(mb), desc.Size, desc.Digest, content.WithLabels(labels)); err != nil ***REMOVED***
		return ocispec.Descriptor***REMOVED******REMOVED***, errors.Wrap(err, "failed to write config")
	***REMOVED***

	ref = remotes.MakeRefKey(ctx, config)
	if err := content.WriteBlob(ctx, c.contentStore, ref, bytes.NewReader(b), config.Size, config.Digest); err != nil ***REMOVED***
		return ocispec.Descriptor***REMOVED******REMOVED***, errors.Wrap(err, "failed to write config")
	***REMOVED***

	return desc, nil
***REMOVED***

func (c *Converter) fetchManifest(ctx context.Context, desc ocispec.Descriptor) error ***REMOVED***
	log.G(ctx).Debug("fetch schema 1")

	rc, err := c.fetcher.Fetch(ctx, desc)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	b, err := ioutil.ReadAll(io.LimitReader(rc, manifestSizeLimit)) // limit to 8MB
	rc.Close()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	b, err = stripSignature(b)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var m manifest
	if err := json.Unmarshal(b, &m); err != nil ***REMOVED***
		return err
	***REMOVED***
	c.pulledManifest = &m

	return nil
***REMOVED***

func (c *Converter) fetchBlob(ctx context.Context, desc ocispec.Descriptor) error ***REMOVED***
	log.G(ctx).Debug("fetch blob")

	var (
		ref   = remotes.MakeRefKey(ctx, desc)
		calc  = newBlobStateCalculator()
		retry = 16
		size  = desc.Size
	)

	// size may be unknown, set to zero for content ingest
	if size == -1 ***REMOVED***
		size = 0
	***REMOVED***

tryit:
	cw, err := c.contentStore.Writer(ctx, ref, size, desc.Digest)
	if err != nil ***REMOVED***
		if errdefs.IsUnavailable(err) ***REMOVED***
			select ***REMOVED***
			case <-time.After(time.Millisecond * time.Duration(rand.Intn(retry))):
				if retry < 2048 ***REMOVED***
					retry = retry << 1
				***REMOVED***
				goto tryit
			case <-ctx.Done():
				return err
			***REMOVED***
		***REMOVED*** else if !errdefs.IsAlreadyExists(err) ***REMOVED***
			return err
		***REMOVED***

		// TODO: Check if blob -> diff id mapping already exists
		// TODO: Check if blob empty label exists

		ra, err := c.contentStore.ReaderAt(ctx, desc.Digest)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		defer ra.Close()

		gr, err := gzip.NewReader(content.NewReader(ra))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		defer gr.Close()

		_, err = io.Copy(calc, gr)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		defer cw.Close()

		rc, err := c.fetcher.Fetch(ctx, desc)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		defer rc.Close()

		eg, _ := errgroup.WithContext(ctx)
		pr, pw := io.Pipe()

		eg.Go(func() error ***REMOVED***
			gr, err := gzip.NewReader(pr)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			defer gr.Close()

			_, err = io.Copy(calc, gr)
			pr.CloseWithError(err)
			return err
		***REMOVED***)

		eg.Go(func() error ***REMOVED***
			defer pw.Close()

			return content.Copy(ctx, cw, io.TeeReader(rc, pw), size, desc.Digest)
		***REMOVED***)

		if err := eg.Wait(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if desc.Size == -1 ***REMOVED***
		info, err := c.contentStore.Info(ctx, desc.Digest)
		if err != nil ***REMOVED***
			return errors.Wrap(err, "failed to get blob info")
		***REMOVED***
		desc.Size = info.Size
	***REMOVED***

	state := calc.State()

	c.mu.Lock()
	c.blobMap[desc.Digest] = state
	c.layerBlobs[state.diffID] = desc
	c.mu.Unlock()

	return nil
***REMOVED***
func (c *Converter) schema1ManifestHistory() ([]ocispec.History, []digest.Digest, error) ***REMOVED***
	if c.pulledManifest == nil ***REMOVED***
		return nil, nil, errors.New("missing schema 1 manifest for conversion")
	***REMOVED***
	m := *c.pulledManifest

	if len(m.History) == 0 ***REMOVED***
		return nil, nil, errors.New("no history")
	***REMOVED***

	history := make([]ocispec.History, len(m.History))
	diffIDs := []digest.Digest***REMOVED******REMOVED***
	for i := range m.History ***REMOVED***
		var h v1History
		if err := json.Unmarshal([]byte(m.History[i].V1Compatibility), &h); err != nil ***REMOVED***
			return nil, nil, errors.Wrap(err, "failed to unmarshal history")
		***REMOVED***

		blobSum := m.FSLayers[i].BlobSum

		state := c.blobMap[blobSum]

		history[len(history)-i-1] = ocispec.History***REMOVED***
			Author:     h.Author,
			Comment:    h.Comment,
			Created:    &h.Created,
			CreatedBy:  strings.Join(h.ContainerConfig.Cmd, " "),
			EmptyLayer: state.empty,
		***REMOVED***

		if !state.empty ***REMOVED***
			diffIDs = append([]digest.Digest***REMOVED***state.diffID***REMOVED***, diffIDs...)

		***REMOVED***
	***REMOVED***

	return history, diffIDs, nil
***REMOVED***

type fsLayer struct ***REMOVED***
	BlobSum digest.Digest `json:"blobSum"`
***REMOVED***

type history struct ***REMOVED***
	V1Compatibility string `json:"v1Compatibility"`
***REMOVED***

type manifest struct ***REMOVED***
	FSLayers []fsLayer `json:"fsLayers"`
	History  []history `json:"history"`
***REMOVED***

type v1History struct ***REMOVED***
	Author          string    `json:"author,omitempty"`
	Created         time.Time `json:"created"`
	Comment         string    `json:"comment,omitempty"`
	ThrowAway       *bool     `json:"throwaway,omitempty"`
	Size            *int      `json:"Size,omitempty"` // used before ThrowAway field
	ContainerConfig struct ***REMOVED***
		Cmd []string `json:"Cmd,omitempty"`
	***REMOVED*** `json:"container_config,omitempty"`
***REMOVED***

// isEmptyLayer returns whether the v1 compatibility history describes an
// empty layer. A return value of true indicates the layer is empty,
// however false does not indicate non-empty.
func isEmptyLayer(compatHistory []byte) (bool, error) ***REMOVED***
	var h v1History
	if err := json.Unmarshal(compatHistory, &h); err != nil ***REMOVED***
		return false, err
	***REMOVED***

	if h.ThrowAway != nil ***REMOVED***
		return *h.ThrowAway, nil
	***REMOVED***
	if h.Size != nil ***REMOVED***
		return *h.Size == 0, nil
	***REMOVED***

	// If no `Size` or `throwaway` field is given, then
	// it cannot be determined whether the layer is empty
	// from the history, return false
	return false, nil
***REMOVED***

type signature struct ***REMOVED***
	Signatures []jsParsedSignature `json:"signatures"`
***REMOVED***

type jsParsedSignature struct ***REMOVED***
	Protected string `json:"protected"`
***REMOVED***

type protectedBlock struct ***REMOVED***
	Length int    `json:"formatLength"`
	Tail   string `json:"formatTail"`
***REMOVED***

// joseBase64UrlDecode decodes the given string using the standard base64 url
// decoder but first adds the appropriate number of trailing '=' characters in
// accordance with the jose specification.
// http://tools.ietf.org/html/draft-ietf-jose-json-web-signature-31#section-2
func joseBase64UrlDecode(s string) ([]byte, error) ***REMOVED***
	switch len(s) % 4 ***REMOVED***
	case 0:
	case 2:
		s += "=="
	case 3:
		s += "="
	default:
		return nil, errors.New("illegal base64url string")
	***REMOVED***
	return base64.URLEncoding.DecodeString(s)
***REMOVED***

func stripSignature(b []byte) ([]byte, error) ***REMOVED***
	var sig signature
	if err := json.Unmarshal(b, &sig); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(sig.Signatures) == 0 ***REMOVED***
		return nil, errors.New("no signatures")
	***REMOVED***
	pb, err := joseBase64UrlDecode(sig.Signatures[0].Protected)
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "could not decode %s", sig.Signatures[0].Protected)
	***REMOVED***

	var protected protectedBlock
	if err := json.Unmarshal(pb, &protected); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if protected.Length > len(b) ***REMOVED***
		return nil, errors.New("invalid protected length block")
	***REMOVED***

	tail, err := joseBase64UrlDecode(protected.Tail)
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "invalid tail base 64 value")
	***REMOVED***

	return append(b[:protected.Length], tail...), nil
***REMOVED***

type blobStateCalculator struct ***REMOVED***
	empty    bool
	digester digest.Digester
***REMOVED***

func newBlobStateCalculator() *blobStateCalculator ***REMOVED***
	return &blobStateCalculator***REMOVED***
		empty:    true,
		digester: digest.Canonical.Digester(),
	***REMOVED***
***REMOVED***

func (c *blobStateCalculator) Write(p []byte) (int, error) ***REMOVED***
	if c.empty ***REMOVED***
		for _, b := range p ***REMOVED***
			if b != 0x00 ***REMOVED***
				c.empty = false
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return c.digester.Hash().Write(p)
***REMOVED***

func (c *blobStateCalculator) State() blobState ***REMOVED***
	return blobState***REMOVED***
		empty:  c.empty,
		diffID: c.digester.Digest(),
	***REMOVED***
***REMOVED***
