package layer

import (
	"io"

	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/containerfs"
)

type mountedLayer struct ***REMOVED***
	name       string
	mountID    string
	initID     string
	parent     *roLayer
	path       string
	layerStore *layerStore

	references map[RWLayer]*referencedRWLayer
***REMOVED***

func (ml *mountedLayer) cacheParent() string ***REMOVED***
	if ml.initID != "" ***REMOVED***
		return ml.initID
	***REMOVED***
	if ml.parent != nil ***REMOVED***
		return ml.parent.cacheID
	***REMOVED***
	return ""
***REMOVED***

func (ml *mountedLayer) TarStream() (io.ReadCloser, error) ***REMOVED***
	return ml.layerStore.driver.Diff(ml.mountID, ml.cacheParent())
***REMOVED***

func (ml *mountedLayer) Name() string ***REMOVED***
	return ml.name
***REMOVED***

func (ml *mountedLayer) Parent() Layer ***REMOVED***
	if ml.parent != nil ***REMOVED***
		return ml.parent
	***REMOVED***

	// Return a nil interface instead of an interface wrapping a nil
	// pointer.
	return nil
***REMOVED***

func (ml *mountedLayer) Size() (int64, error) ***REMOVED***
	return ml.layerStore.driver.DiffSize(ml.mountID, ml.cacheParent())
***REMOVED***

func (ml *mountedLayer) Changes() ([]archive.Change, error) ***REMOVED***
	return ml.layerStore.driver.Changes(ml.mountID, ml.cacheParent())
***REMOVED***

func (ml *mountedLayer) Metadata() (map[string]string, error) ***REMOVED***
	return ml.layerStore.driver.GetMetadata(ml.mountID)
***REMOVED***

func (ml *mountedLayer) getReference() RWLayer ***REMOVED***
	ref := &referencedRWLayer***REMOVED***
		mountedLayer: ml,
	***REMOVED***
	ml.references[ref] = ref

	return ref
***REMOVED***

func (ml *mountedLayer) hasReferences() bool ***REMOVED***
	return len(ml.references) > 0
***REMOVED***

func (ml *mountedLayer) deleteReference(ref RWLayer) error ***REMOVED***
	if _, ok := ml.references[ref]; !ok ***REMOVED***
		return ErrLayerNotRetained
	***REMOVED***
	delete(ml.references, ref)
	return nil
***REMOVED***

func (ml *mountedLayer) retakeReference(r RWLayer) ***REMOVED***
	if ref, ok := r.(*referencedRWLayer); ok ***REMOVED***
		ml.references[ref] = ref
	***REMOVED***
***REMOVED***

type referencedRWLayer struct ***REMOVED***
	*mountedLayer
***REMOVED***

func (rl *referencedRWLayer) Mount(mountLabel string) (containerfs.ContainerFS, error) ***REMOVED***
	return rl.layerStore.driver.Get(rl.mountedLayer.mountID, mountLabel)
***REMOVED***

// Unmount decrements the activity count and unmounts the underlying layer
// Callers should only call `Unmount` once per call to `Mount`, even on error.
func (rl *referencedRWLayer) Unmount() error ***REMOVED***
	return rl.layerStore.driver.Put(rl.mountedLayer.mountID)
***REMOVED***
