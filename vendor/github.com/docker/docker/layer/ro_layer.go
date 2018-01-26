package layer

import (
	"fmt"
	"io"

	"github.com/docker/distribution"
	"github.com/opencontainers/go-digest"
)

type roLayer struct ***REMOVED***
	chainID    ChainID
	diffID     DiffID
	parent     *roLayer
	cacheID    string
	size       int64
	layerStore *layerStore
	descriptor distribution.Descriptor

	referenceCount int
	references     map[Layer]struct***REMOVED******REMOVED***
***REMOVED***

// TarStream for roLayer guarantees that the data that is produced is the exact
// data that the layer was registered with.
func (rl *roLayer) TarStream() (io.ReadCloser, error) ***REMOVED***
	rc, err := rl.layerStore.getTarStream(rl)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	vrc, err := newVerifiedReadCloser(rc, digest.Digest(rl.diffID))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return vrc, nil
***REMOVED***

// TarStreamFrom does not make any guarantees to the correctness of the produced
// data. As such it should not be used when the layer content must be verified
// to be an exact match to the registered layer.
func (rl *roLayer) TarStreamFrom(parent ChainID) (io.ReadCloser, error) ***REMOVED***
	var parentCacheID string
	for pl := rl.parent; pl != nil; pl = pl.parent ***REMOVED***
		if pl.chainID == parent ***REMOVED***
			parentCacheID = pl.cacheID
			break
		***REMOVED***
	***REMOVED***

	if parent != ChainID("") && parentCacheID == "" ***REMOVED***
		return nil, fmt.Errorf("layer ID '%s' is not a parent of the specified layer: cannot provide diff to non-parent", parent)
	***REMOVED***
	return rl.layerStore.driver.Diff(rl.cacheID, parentCacheID)
***REMOVED***

func (rl *roLayer) ChainID() ChainID ***REMOVED***
	return rl.chainID
***REMOVED***

func (rl *roLayer) DiffID() DiffID ***REMOVED***
	return rl.diffID
***REMOVED***

func (rl *roLayer) Parent() Layer ***REMOVED***
	if rl.parent == nil ***REMOVED***
		return nil
	***REMOVED***
	return rl.parent
***REMOVED***

func (rl *roLayer) Size() (size int64, err error) ***REMOVED***
	if rl.parent != nil ***REMOVED***
		size, err = rl.parent.Size()
		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	return size + rl.size, nil
***REMOVED***

func (rl *roLayer) DiffSize() (size int64, err error) ***REMOVED***
	return rl.size, nil
***REMOVED***

func (rl *roLayer) Metadata() (map[string]string, error) ***REMOVED***
	return rl.layerStore.driver.GetMetadata(rl.cacheID)
***REMOVED***

type referencedCacheLayer struct ***REMOVED***
	*roLayer
***REMOVED***

func (rl *roLayer) getReference() Layer ***REMOVED***
	ref := &referencedCacheLayer***REMOVED***
		roLayer: rl,
	***REMOVED***
	rl.references[ref] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

	return ref
***REMOVED***

func (rl *roLayer) hasReference(ref Layer) bool ***REMOVED***
	_, ok := rl.references[ref]
	return ok
***REMOVED***

func (rl *roLayer) hasReferences() bool ***REMOVED***
	return len(rl.references) > 0
***REMOVED***

func (rl *roLayer) deleteReference(ref Layer) ***REMOVED***
	delete(rl.references, ref)
***REMOVED***

func (rl *roLayer) depth() int ***REMOVED***
	if rl.parent == nil ***REMOVED***
		return 1
	***REMOVED***
	return rl.parent.depth() + 1
***REMOVED***

func storeLayer(tx MetadataTransaction, layer *roLayer) error ***REMOVED***
	if err := tx.SetDiffID(layer.diffID); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := tx.SetSize(layer.size); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := tx.SetCacheID(layer.cacheID); err != nil ***REMOVED***
		return err
	***REMOVED***
	// Do not store empty descriptors
	if layer.descriptor.Digest != "" ***REMOVED***
		if err := tx.SetDescriptor(layer.descriptor); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if layer.parent != nil ***REMOVED***
		if err := tx.SetParent(layer.parent.chainID); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return tx.setOS(layer.layerStore.os)
***REMOVED***

func newVerifiedReadCloser(rc io.ReadCloser, dgst digest.Digest) (io.ReadCloser, error) ***REMOVED***
	return &verifiedReadCloser***REMOVED***
		rc:       rc,
		dgst:     dgst,
		verifier: dgst.Verifier(),
	***REMOVED***, nil
***REMOVED***

type verifiedReadCloser struct ***REMOVED***
	rc       io.ReadCloser
	dgst     digest.Digest
	verifier digest.Verifier
***REMOVED***

func (vrc *verifiedReadCloser) Read(p []byte) (n int, err error) ***REMOVED***
	n, err = vrc.rc.Read(p)
	if n > 0 ***REMOVED***
		if n, err := vrc.verifier.Write(p[:n]); err != nil ***REMOVED***
			return n, err
		***REMOVED***
	***REMOVED***
	if err == io.EOF ***REMOVED***
		if !vrc.verifier.Verified() ***REMOVED***
			err = fmt.Errorf("could not verify layer data for: %s. This may be because internal files in the layer store were modified. Re-pulling or rebuilding this image may resolve the issue", vrc.dgst)
		***REMOVED***
	***REMOVED***
	return
***REMOVED***
func (vrc *verifiedReadCloser) Close() error ***REMOVED***
	return vrc.rc.Close()
***REMOVED***
