package image

import (
	"runtime"

	"github.com/docker/docker/layer"
	"github.com/sirupsen/logrus"
)

// TypeLayers is used for RootFS.Type for filesystems organized into layers.
const TypeLayers = "layers"

// typeLayersWithBase is an older format used by Windows up to v1.12. We
// explicitly handle this as an error case to ensure that a daemon which still
// has an older image like this on disk can still start, even though the
// image itself is not usable. See https://github.com/docker/docker/pull/25806.
const typeLayersWithBase = "layers+base"

// RootFS describes images root filesystem
// This is currently a placeholder that only supports layers. In the future
// this can be made into an interface that supports different implementations.
type RootFS struct ***REMOVED***
	Type    string         `json:"type"`
	DiffIDs []layer.DiffID `json:"diff_ids,omitempty"`
***REMOVED***

// NewRootFS returns empty RootFS struct
func NewRootFS() *RootFS ***REMOVED***
	return &RootFS***REMOVED***Type: TypeLayers***REMOVED***
***REMOVED***

// Append appends a new diffID to rootfs
func (r *RootFS) Append(id layer.DiffID) ***REMOVED***
	r.DiffIDs = append(r.DiffIDs, id)
***REMOVED***

// Clone returns a copy of the RootFS
func (r *RootFS) Clone() *RootFS ***REMOVED***
	newRoot := NewRootFS()
	newRoot.Type = r.Type
	newRoot.DiffIDs = append(r.DiffIDs)
	return newRoot
***REMOVED***

// ChainID returns the ChainID for the top layer in RootFS.
func (r *RootFS) ChainID() layer.ChainID ***REMOVED***
	if runtime.GOOS == "windows" && r.Type == typeLayersWithBase ***REMOVED***
		logrus.Warnf("Layer type is unsupported on this platform. DiffIDs: '%v'", r.DiffIDs)
		return ""
	***REMOVED***
	return layer.CreateChainID(r.DiffIDs)
***REMOVED***
