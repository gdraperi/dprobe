package rootfs

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/containerd/containerd/diff"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/snapshots"
	"github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/identity"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// Layer represents the descriptors for a layer diff. These descriptions
// include the descriptor for the uncompressed tar diff as well as a blob
// used to transport that tar. The blob descriptor may or may not describe
// a compressed object.
type Layer struct ***REMOVED***
	Diff ocispec.Descriptor
	Blob ocispec.Descriptor
***REMOVED***

// ApplyLayers applies all the layers using the given snapshotter and applier.
// The returned result is a chain id digest representing all the applied layers.
// Layers are applied in order they are given, making the first layer the
// bottom-most layer in the layer chain.
func ApplyLayers(ctx context.Context, layers []Layer, sn snapshots.Snapshotter, a diff.Differ) (digest.Digest, error) ***REMOVED***
	var chain []digest.Digest
	for _, layer := range layers ***REMOVED***
		if _, err := ApplyLayer(ctx, layer, chain, sn, a); err != nil ***REMOVED***
			// TODO: possibly wait and retry if extraction of same chain id was in progress
			return "", err
		***REMOVED***

		chain = append(chain, layer.Diff.Digest)
	***REMOVED***
	return identity.ChainID(chain), nil
***REMOVED***

// ApplyLayer applies a single layer on top of the given provided layer chain,
// using the provided snapshotter and applier. If the layer was unpacked true
// is returned, if the layer already exists false is returned.
func ApplyLayer(ctx context.Context, layer Layer, chain []digest.Digest, sn snapshots.Snapshotter, a diff.Differ, opts ...snapshots.Opt) (bool, error) ***REMOVED***
	var (
		parent  = identity.ChainID(chain)
		chainID = identity.ChainID(append(chain, layer.Diff.Digest))
		diff    ocispec.Descriptor
	)

	_, err := sn.Stat(ctx, chainID.String())
	if err == nil ***REMOVED***
		log.G(ctx).Debugf("Extraction not needed, layer snapshot %s exists", chainID)
		return false, nil
	***REMOVED*** else if !errdefs.IsNotFound(err) ***REMOVED***
		return false, errors.Wrapf(err, "failed to stat snapshot %s", chainID)
	***REMOVED***

	key := fmt.Sprintf("extract-%s %s", uniquePart(), chainID)

	// Prepare snapshot with from parent, label as root
	mounts, err := sn.Prepare(ctx, key, parent.String(), opts...)
	if err != nil ***REMOVED***
		//TODO: If is snapshot exists error, retry
		return false, errors.Wrapf(err, "failed to prepare extraction snapshot %q", key)
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			log.G(ctx).WithError(err).WithField("key", key).Infof("Apply failure, attempting cleanup")
			if rerr := sn.Remove(ctx, key); rerr != nil ***REMOVED***
				log.G(ctx).WithError(rerr).Warnf("Extraction snapshot %q removal failed", key)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	diff, err = a.Apply(ctx, layer.Blob, mounts)
	if err != nil ***REMOVED***
		return false, errors.Wrapf(err, "failed to extract layer %s", layer.Diff.Digest)
	***REMOVED***
	if diff.Digest != layer.Diff.Digest ***REMOVED***
		err = errors.Errorf("wrong diff id calculated on extraction %q", diff.Digest)
		return false, err
	***REMOVED***

	if err = sn.Commit(ctx, chainID.String(), key, opts...); err != nil ***REMOVED***
		if !errdefs.IsAlreadyExists(err) ***REMOVED***
			return false, errors.Wrapf(err, "failed to commit snapshot %s", key)
		***REMOVED***

		// Destination already exists, cleanup key and return without error
		err = nil
		if err := sn.Remove(ctx, key); err != nil ***REMOVED***
			return false, errors.Wrapf(err, "failed to cleanup aborted apply %s", key)
		***REMOVED***
		return false, nil
	***REMOVED***

	return true, nil
***REMOVED***

func uniquePart() string ***REMOVED***
	t := time.Now()
	var b [3]byte
	// Ignore read failures, just decreases uniqueness
	rand.Read(b[:])
	return fmt.Sprintf("%d-%s", t.Nanosecond(), base64.URLEncoding.EncodeToString(b[:]))
***REMOVED***
