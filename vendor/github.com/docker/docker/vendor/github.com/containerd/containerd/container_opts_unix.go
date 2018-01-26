// +build !windows

package containerd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"

	"github.com/containerd/containerd/api/types"
	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/mount"
	"github.com/containerd/containerd/platforms"
	"github.com/gogo/protobuf/proto"
	protobuf "github.com/gogo/protobuf/types"
	digest "github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/identity"
	"github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

// WithCheckpoint allows a container to be created from the checkpointed information
// provided by the descriptor. The image, snapshot, and runtime specifications are
// restored on the container
func WithCheckpoint(im Image, snapshotKey string) NewContainerOpts ***REMOVED***
	// set image and rw, and spec
	return func(ctx context.Context, client *Client, c *containers.Container) error ***REMOVED***
		var (
			desc  = im.Target()
			id    = desc.Digest
			store = client.ContentStore()
		)
		index, err := decodeIndex(ctx, store, id)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		var rw *v1.Descriptor
		for _, m := range index.Manifests ***REMOVED***
			switch m.MediaType ***REMOVED***
			case v1.MediaTypeImageLayer:
				fk := m
				rw = &fk
			case images.MediaTypeDockerSchema2Manifest, images.MediaTypeDockerSchema2ManifestList:
				config, err := images.Config(ctx, store, m, platforms.Default())
				if err != nil ***REMOVED***
					return errors.Wrap(err, "unable to resolve image config")
				***REMOVED***
				diffIDs, err := images.RootFS(ctx, store, config)
				if err != nil ***REMOVED***
					return errors.Wrap(err, "unable to get rootfs")
				***REMOVED***
				setSnapshotterIfEmpty(c)
				if _, err := client.SnapshotService(c.Snapshotter).Prepare(ctx, snapshotKey, identity.ChainID(diffIDs).String()); err != nil ***REMOVED***
					if !errdefs.IsAlreadyExists(err) ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
				c.Image = index.Annotations["image.name"]
			case images.MediaTypeContainerd1CheckpointConfig:
				data, err := content.ReadBlob(ctx, store, m.Digest)
				if err != nil ***REMOVED***
					return errors.Wrap(err, "unable to read checkpoint config")
				***REMOVED***
				var any protobuf.Any
				if err := proto.Unmarshal(data, &any); err != nil ***REMOVED***
					return err
				***REMOVED***
				c.Spec = &any
			***REMOVED***
		***REMOVED***
		if rw != nil ***REMOVED***
			// apply the rw snapshot to the new rw layer
			mounts, err := client.SnapshotService(c.Snapshotter).Mounts(ctx, snapshotKey)
			if err != nil ***REMOVED***
				return errors.Wrapf(err, "unable to get mounts for %s", snapshotKey)
			***REMOVED***
			if _, err := client.DiffService().Apply(ctx, *rw, mounts); err != nil ***REMOVED***
				return errors.Wrap(err, "unable to apply rw diff")
			***REMOVED***
		***REMOVED***
		c.SnapshotKey = snapshotKey
		return nil
	***REMOVED***
***REMOVED***

// WithTaskCheckpoint allows a task to be created with live runtime and memory data from a
// previous checkpoint. Additional software such as CRIU may be required to
// restore a task from a checkpoint
func WithTaskCheckpoint(im Image) NewTaskOpts ***REMOVED***
	return func(ctx context.Context, c *Client, info *TaskInfo) error ***REMOVED***
		desc := im.Target()
		id := desc.Digest
		index, err := decodeIndex(ctx, c.ContentStore(), id)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		for _, m := range index.Manifests ***REMOVED***
			if m.MediaType == images.MediaTypeContainerd1Checkpoint ***REMOVED***
				info.Checkpoint = &types.Descriptor***REMOVED***
					MediaType: m.MediaType,
					Size_:     m.Size,
					Digest:    m.Digest,
				***REMOVED***
				return nil
			***REMOVED***
		***REMOVED***
		return fmt.Errorf("checkpoint not found in index %s", id)
	***REMOVED***
***REMOVED***

func decodeIndex(ctx context.Context, store content.Store, id digest.Digest) (*v1.Index, error) ***REMOVED***
	var index v1.Index
	p, err := content.ReadBlob(ctx, store, id)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := json.Unmarshal(p, &index); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &index, nil
***REMOVED***

// WithRemappedSnapshot creates a new snapshot and remaps the uid/gid for the
// filesystem to be used by a container with user namespaces
func WithRemappedSnapshot(id string, i Image, uid, gid uint32) NewContainerOpts ***REMOVED***
	return withRemappedSnapshotBase(id, i, uid, gid, false)
***REMOVED***

// WithRemappedSnapshotView is similar to WithRemappedSnapshot but rootfs is mounted as read-only.
func WithRemappedSnapshotView(id string, i Image, uid, gid uint32) NewContainerOpts ***REMOVED***
	return withRemappedSnapshotBase(id, i, uid, gid, true)
***REMOVED***

func withRemappedSnapshotBase(id string, i Image, uid, gid uint32, readonly bool) NewContainerOpts ***REMOVED***
	return func(ctx context.Context, client *Client, c *containers.Container) error ***REMOVED***
		diffIDs, err := i.(*image).i.RootFS(ctx, client.ContentStore(), platforms.Default())
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		setSnapshotterIfEmpty(c)

		var (
			snapshotter = client.SnapshotService(c.Snapshotter)
			parent      = identity.ChainID(diffIDs).String()
			usernsID    = fmt.Sprintf("%s-%d-%d", parent, uid, gid)
		)
		if _, err := snapshotter.Stat(ctx, usernsID); err == nil ***REMOVED***
			if _, err := snapshotter.Prepare(ctx, id, usernsID); err == nil ***REMOVED***
				c.SnapshotKey = id
				c.Image = i.Name()
				return nil
			***REMOVED*** else if !errdefs.IsNotFound(err) ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		mounts, err := snapshotter.Prepare(ctx, usernsID+"-remap", parent)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := remapRootFS(mounts, uid, gid); err != nil ***REMOVED***
			snapshotter.Remove(ctx, usernsID)
			return err
		***REMOVED***
		if err := snapshotter.Commit(ctx, usernsID, usernsID+"-remap"); err != nil ***REMOVED***
			return err
		***REMOVED***
		if readonly ***REMOVED***
			_, err = snapshotter.View(ctx, id, usernsID)
		***REMOVED*** else ***REMOVED***
			_, err = snapshotter.Prepare(ctx, id, usernsID)
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		c.SnapshotKey = id
		c.Image = i.Name()
		return nil
	***REMOVED***
***REMOVED***

func remapRootFS(mounts []mount.Mount, uid, gid uint32) error ***REMOVED***
	root, err := ioutil.TempDir("", "ctd-remap")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer os.Remove(root)
	for _, m := range mounts ***REMOVED***
		if err := m.Mount(root); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	err = filepath.Walk(root, incrementFS(root, uid, gid))
	if uerr := mount.Unmount(root, 0); err == nil ***REMOVED***
		err = uerr
	***REMOVED***
	return err
***REMOVED***

func incrementFS(root string, uidInc, gidInc uint32) filepath.WalkFunc ***REMOVED***
	return func(path string, info os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		var (
			stat = info.Sys().(*syscall.Stat_t)
			u, g = int(stat.Uid + uidInc), int(stat.Gid + gidInc)
		)
		// be sure the lchown the path as to not de-reference the symlink to a host file
		return os.Lchown(path, u, g)
	***REMOVED***
***REMOVED***
