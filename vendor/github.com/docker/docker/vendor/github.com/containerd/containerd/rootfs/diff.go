package rootfs

import (
	"fmt"

	"github.com/containerd/containerd/diff"
	"github.com/containerd/containerd/mount"
	"github.com/containerd/containerd/snapshots"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"golang.org/x/net/context"
)

// Diff creates a layer diff for the given snapshot identifier from the parent
// of the snapshot. A content ref is provided to track the progress of the
// content creation and the provided snapshotter and mount differ are used
// for calculating the diff. The descriptor for the layer diff is returned.
func Diff(ctx context.Context, snapshotID string, sn snapshots.Snapshotter, d diff.Differ, opts ...diff.Opt) (ocispec.Descriptor, error) ***REMOVED***
	info, err := sn.Stat(ctx, snapshotID)
	if err != nil ***REMOVED***
		return ocispec.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***

	lowerKey := fmt.Sprintf("%s-parent-view", info.Parent)
	lower, err := sn.View(ctx, lowerKey, info.Parent)
	if err != nil ***REMOVED***
		return ocispec.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***
	defer sn.Remove(ctx, lowerKey)

	var upper []mount.Mount
	if info.Kind == snapshots.KindActive ***REMOVED***
		upper, err = sn.Mounts(ctx, snapshotID)
		if err != nil ***REMOVED***
			return ocispec.Descriptor***REMOVED******REMOVED***, err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		upperKey := fmt.Sprintf("%s-view", snapshotID)
		upper, err = sn.View(ctx, upperKey, snapshotID)
		if err != nil ***REMOVED***
			return ocispec.Descriptor***REMOVED******REMOVED***, err
		***REMOVED***
		defer sn.Remove(ctx, lowerKey)
	***REMOVED***

	return d.DiffMounts(ctx, lower, upper, opts...)
***REMOVED***
