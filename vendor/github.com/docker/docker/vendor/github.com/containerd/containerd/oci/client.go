package oci

import (
	"context"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/snapshots"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// Client interface used by SpecOpt
type Client interface ***REMOVED***
	SnapshotService(snapshotterName string) snapshots.Snapshotter
***REMOVED***

// Image interface used by some SpecOpt to query image configuration
type Image interface ***REMOVED***
	// Config descriptor for the image.
	Config(ctx context.Context) (ocispec.Descriptor, error)
	// ContentStore provides a content store which contains image blob data
	ContentStore() content.Store
***REMOVED***
