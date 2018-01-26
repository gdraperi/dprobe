package rootfs

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/mount"
	"github.com/containerd/containerd/snapshots"
	digest "github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

var (
	initializers = map[string]initializerFunc***REMOVED******REMOVED***
)

type initializerFunc func(string) error

// Mounter handles mount and unmount
type Mounter interface ***REMOVED***
	Mount(target string, mounts ...mount.Mount) error
	Unmount(target string) error
***REMOVED***

// InitRootFS initializes the snapshot for use as a rootfs
func InitRootFS(ctx context.Context, name string, parent digest.Digest, readonly bool, snapshotter snapshots.Snapshotter, mounter Mounter) ([]mount.Mount, error) ***REMOVED***
	_, err := snapshotter.Stat(ctx, name)
	if err == nil ***REMOVED***
		return nil, errors.Errorf("rootfs already exists")
	***REMOVED***
	// TODO: ensure not exist error once added to snapshot package

	parentS := parent.String()

	initName := defaultInitializer
	initFn := initializers[initName]
	if initFn != nil ***REMOVED***
		parentS, err = createInitLayer(ctx, parentS, initName, initFn, snapshotter, mounter)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	if readonly ***REMOVED***
		return snapshotter.View(ctx, name, parentS)
	***REMOVED***

	return snapshotter.Prepare(ctx, name, parentS)
***REMOVED***

func createInitLayer(ctx context.Context, parent, initName string, initFn func(string) error, snapshotter snapshots.Snapshotter, mounter Mounter) (string, error) ***REMOVED***
	initS := fmt.Sprintf("%s %s", parent, initName)
	if _, err := snapshotter.Stat(ctx, initS); err == nil ***REMOVED***
		return initS, nil
	***REMOVED***
	// TODO: ensure not exist error once added to snapshot package

	// Create tempdir
	td, err := ioutil.TempDir("", "create-init-")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	defer os.RemoveAll(td)

	mounts, err := snapshotter.Prepare(ctx, td, parent)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			if rerr := snapshotter.Remove(ctx, td); rerr != nil ***REMOVED***
				log.G(ctx).Errorf("Failed to remove snapshot %s: %v", td, rerr)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	if err = mounter.Mount(td, mounts...); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	if err = initFn(td); err != nil ***REMOVED***
		if merr := mounter.Unmount(td); merr != nil ***REMOVED***
			log.G(ctx).Errorf("Failed to unmount %s: %v", td, merr)
		***REMOVED***
		return "", err
	***REMOVED***

	if err = mounter.Unmount(td); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	if err := snapshotter.Commit(ctx, initS, td); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return initS, nil
***REMOVED***
