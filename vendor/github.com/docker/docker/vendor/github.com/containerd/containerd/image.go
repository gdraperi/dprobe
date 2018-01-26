package containerd

import (
	"context"
	"fmt"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/platforms"
	"github.com/containerd/containerd/rootfs"
	"github.com/containerd/containerd/snapshots"
	digest "github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/identity"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

// Image describes an image used by containers
type Image interface ***REMOVED***
	// Name of the image
	Name() string
	// Target descriptor for the image content
	Target() ocispec.Descriptor
	// Unpack unpacks the image's content into a snapshot
	Unpack(context.Context, string) error
	// RootFS returns the unpacked diffids that make up images rootfs.
	RootFS(ctx context.Context) ([]digest.Digest, error)
	// Size returns the total size of the image's packed resources.
	Size(ctx context.Context) (int64, error)
	// Config descriptor for the image.
	Config(ctx context.Context) (ocispec.Descriptor, error)
	// IsUnpacked returns whether or not an image is unpacked.
	IsUnpacked(context.Context, string) (bool, error)
	// ContentStore provides a content store which contains image blob data
	ContentStore() content.Store
***REMOVED***

var _ = (Image)(&image***REMOVED******REMOVED***)

type image struct ***REMOVED***
	client *Client

	i images.Image
***REMOVED***

func (i *image) Name() string ***REMOVED***
	return i.i.Name
***REMOVED***

func (i *image) Target() ocispec.Descriptor ***REMOVED***
	return i.i.Target
***REMOVED***

func (i *image) RootFS(ctx context.Context) ([]digest.Digest, error) ***REMOVED***
	provider := i.client.ContentStore()
	return i.i.RootFS(ctx, provider, platforms.Default())
***REMOVED***

func (i *image) Size(ctx context.Context) (int64, error) ***REMOVED***
	provider := i.client.ContentStore()
	return i.i.Size(ctx, provider, platforms.Default())
***REMOVED***

func (i *image) Config(ctx context.Context) (ocispec.Descriptor, error) ***REMOVED***
	provider := i.client.ContentStore()
	return i.i.Config(ctx, provider, platforms.Default())
***REMOVED***

func (i *image) IsUnpacked(ctx context.Context, snapshotterName string) (bool, error) ***REMOVED***
	sn := i.client.SnapshotService(snapshotterName)
	cs := i.client.ContentStore()

	diffs, err := i.i.RootFS(ctx, cs, platforms.Default())
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	chainID := identity.ChainID(diffs)
	_, err = sn.Stat(ctx, chainID.String())
	if err == nil ***REMOVED***
		return true, nil
	***REMOVED*** else if !errdefs.IsNotFound(err) ***REMOVED***
		return false, err
	***REMOVED***

	return false, nil
***REMOVED***

func (i *image) Unpack(ctx context.Context, snapshotterName string) error ***REMOVED***
	ctx, done, err := i.client.WithLease(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer done()

	layers, err := i.getLayers(ctx, platforms.Default())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var (
		sn = i.client.SnapshotService(snapshotterName)
		a  = i.client.DiffService()
		cs = i.client.ContentStore()

		chain    []digest.Digest
		unpacked bool
	)
	for _, layer := range layers ***REMOVED***
		labels := map[string]string***REMOVED***
			"containerd.io/uncompressed": layer.Diff.Digest.String(),
		***REMOVED***

		unpacked, err = rootfs.ApplyLayer(ctx, layer, chain, sn, a, snapshots.WithLabels(labels))
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		chain = append(chain, layer.Diff.Digest)
	***REMOVED***

	if unpacked ***REMOVED***
		desc, err := i.i.Config(ctx, cs, platforms.Default())
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		rootfs := identity.ChainID(chain).String()

		cinfo := content.Info***REMOVED***
			Digest: desc.Digest,
			Labels: map[string]string***REMOVED***
				fmt.Sprintf("containerd.io/gc.ref.snapshot.%s", snapshotterName): rootfs,
			***REMOVED***,
		***REMOVED***
		if _, err := cs.Update(ctx, cinfo, fmt.Sprintf("labels.containerd.io/gc.ref.snapshot.%s", snapshotterName)); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (i *image) getLayers(ctx context.Context, platform string) ([]rootfs.Layer, error) ***REMOVED***
	cs := i.client.ContentStore()

	manifest, err := images.Manifest(ctx, cs, i.i.Target, platform)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	diffIDs, err := i.i.RootFS(ctx, cs, platform)
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "failed to resolve rootfs")
	***REMOVED***
	if len(diffIDs) != len(manifest.Layers) ***REMOVED***
		return nil, errors.Errorf("mismatched image rootfs and manifest layers")
	***REMOVED***
	layers := make([]rootfs.Layer, len(diffIDs))
	for i := range diffIDs ***REMOVED***
		layers[i].Diff = ocispec.Descriptor***REMOVED***
			// TODO: derive media type from compressed type
			MediaType: ocispec.MediaTypeImageLayer,
			Digest:    diffIDs[i],
		***REMOVED***
		layers[i].Blob = manifest.Layers[i]
	***REMOVED***
	return layers, nil
***REMOVED***

func (i *image) ContentStore() content.Store ***REMOVED***
	return i.client.ContentStore()
***REMOVED***
