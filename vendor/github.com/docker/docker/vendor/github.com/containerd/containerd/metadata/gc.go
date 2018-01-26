package metadata

import (
	"context"
	"fmt"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/containerd/containerd/gc"
	"github.com/containerd/containerd/log"
	"github.com/pkg/errors"
)

const (
	// ResourceUnknown specifies an unknown resource
	ResourceUnknown gc.ResourceType = iota
	// ResourceContent specifies a content resource
	ResourceContent
	// ResourceSnapshot specifies a snapshot resource
	ResourceSnapshot
	// ResourceContainer specifies a container resource
	ResourceContainer
	// ResourceTask specifies a task resource
	ResourceTask
)

var (
	labelGCRoot       = []byte("containerd.io/gc.root")
	labelGCSnapRef    = []byte("containerd.io/gc.ref.snapshot.")
	labelGCContentRef = []byte("containerd.io/gc.ref.content")
)

func scanRoots(ctx context.Context, tx *bolt.Tx, nc chan<- gc.Node) error ***REMOVED***
	v1bkt := tx.Bucket(bucketKeyVersion)
	if v1bkt == nil ***REMOVED***
		return nil
	***REMOVED***

	// iterate through each namespace
	v1c := v1bkt.Cursor()

	for k, v := v1c.First(); k != nil; k, v = v1c.Next() ***REMOVED***
		if v != nil ***REMOVED***
			continue
		***REMOVED***
		nbkt := v1bkt.Bucket(k)
		ns := string(k)

		lbkt := nbkt.Bucket(bucketKeyObjectLeases)
		if lbkt != nil ***REMOVED***
			if err := lbkt.ForEach(func(k, v []byte) error ***REMOVED***
				if v != nil ***REMOVED***
					return nil
				***REMOVED***
				libkt := lbkt.Bucket(k)

				cbkt := libkt.Bucket(bucketKeyObjectContent)
				if cbkt != nil ***REMOVED***
					if err := cbkt.ForEach(func(k, v []byte) error ***REMOVED***
						select ***REMOVED***
						case nc <- gcnode(ResourceContent, ns, string(k)):
						case <-ctx.Done():
							return ctx.Err()
						***REMOVED***
						return nil
					***REMOVED***); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***

				sbkt := libkt.Bucket(bucketKeyObjectSnapshots)
				if sbkt != nil ***REMOVED***
					if err := sbkt.ForEach(func(sk, sv []byte) error ***REMOVED***
						if sv != nil ***REMOVED***
							return nil
						***REMOVED***
						snbkt := sbkt.Bucket(sk)

						return snbkt.ForEach(func(k, v []byte) error ***REMOVED***
							select ***REMOVED***
							case nc <- gcnode(ResourceSnapshot, ns, fmt.Sprintf("%s/%s", sk, k)):
							case <-ctx.Done():
								return ctx.Err()
							***REMOVED***
							return nil
						***REMOVED***)
					***REMOVED***); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***

				return nil
			***REMOVED***); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		ibkt := nbkt.Bucket(bucketKeyObjectImages)
		if ibkt != nil ***REMOVED***
			if err := ibkt.ForEach(func(k, v []byte) error ***REMOVED***
				if v != nil ***REMOVED***
					return nil
				***REMOVED***

				target := ibkt.Bucket(k).Bucket(bucketKeyTarget)
				if target != nil ***REMOVED***
					contentKey := string(target.Get(bucketKeyDigest))
					select ***REMOVED***
					case nc <- gcnode(ResourceContent, ns, contentKey):
					case <-ctx.Done():
						return ctx.Err()
					***REMOVED***
				***REMOVED***
				return sendSnapshotRefs(ns, ibkt.Bucket(k), func(n gc.Node) ***REMOVED***
					select ***REMOVED***
					case nc <- n:
					case <-ctx.Done():
					***REMOVED***
				***REMOVED***)
			***REMOVED***); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		cbkt := nbkt.Bucket(bucketKeyObjectContent)
		if cbkt != nil ***REMOVED***
			cbkt = cbkt.Bucket(bucketKeyObjectBlob)
		***REMOVED***
		if cbkt != nil ***REMOVED***
			if err := cbkt.ForEach(func(k, v []byte) error ***REMOVED***
				if v != nil ***REMOVED***
					return nil
				***REMOVED***
				return sendRootRef(ctx, nc, gcnode(ResourceContent, ns, string(k)), cbkt.Bucket(k))
			***REMOVED***); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		cbkt = nbkt.Bucket(bucketKeyObjectContainers)
		if cbkt != nil ***REMOVED***
			if err := cbkt.ForEach(func(k, v []byte) error ***REMOVED***
				if v != nil ***REMOVED***
					return nil
				***REMOVED***
				snapshotter := string(cbkt.Bucket(k).Get(bucketKeySnapshotter))
				if snapshotter != "" ***REMOVED***
					ss := string(cbkt.Bucket(k).Get(bucketKeySnapshotKey))
					select ***REMOVED***
					case nc <- gcnode(ResourceSnapshot, ns, fmt.Sprintf("%s/%s", snapshotter, ss)):
					case <-ctx.Done():
						return ctx.Err()
					***REMOVED***
				***REMOVED***

				// TODO: Send additional snapshot refs through labels
				return sendSnapshotRefs(ns, cbkt.Bucket(k), func(n gc.Node) ***REMOVED***
					select ***REMOVED***
					case nc <- n:
					case <-ctx.Done():
					***REMOVED***
				***REMOVED***)
			***REMOVED***); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		sbkt := nbkt.Bucket(bucketKeyObjectSnapshots)
		if sbkt != nil ***REMOVED***
			if err := sbkt.ForEach(func(sk, sv []byte) error ***REMOVED***
				if sv != nil ***REMOVED***
					return nil
				***REMOVED***
				snbkt := sbkt.Bucket(sk)

				return snbkt.ForEach(func(k, v []byte) error ***REMOVED***
					if v != nil ***REMOVED***
						return nil
					***REMOVED***

					return sendRootRef(ctx, nc, gcnode(ResourceSnapshot, ns, fmt.Sprintf("%s/%s", sk, k)), snbkt.Bucket(k))
				***REMOVED***)
			***REMOVED***); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func references(ctx context.Context, tx *bolt.Tx, node gc.Node, fn func(gc.Node)) error ***REMOVED***
	if node.Type == ResourceContent ***REMOVED***
		bkt := getBucket(tx, bucketKeyVersion, []byte(node.Namespace), bucketKeyObjectContent, bucketKeyObjectBlob, []byte(node.Key))
		if bkt == nil ***REMOVED***
			// Node may be created from dead edge
			return nil
		***REMOVED***

		if err := sendSnapshotRefs(node.Namespace, bkt, fn); err != nil ***REMOVED***
			return err
		***REMOVED***
		return sendContentRefs(node.Namespace, bkt, fn)
	***REMOVED*** else if node.Type == ResourceSnapshot ***REMOVED***
		parts := strings.SplitN(node.Key, "/", 2)
		if len(parts) != 2 ***REMOVED***
			return errors.Errorf("invalid snapshot gc key %s", node.Key)
		***REMOVED***
		ss := parts[0]
		name := parts[1]

		bkt := getBucket(tx, bucketKeyVersion, []byte(node.Namespace), bucketKeyObjectSnapshots, []byte(ss), []byte(name))
		if bkt == nil ***REMOVED***
			getBucket(tx, bucketKeyVersion, []byte(node.Namespace), bucketKeyObjectSnapshots).ForEach(func(k, v []byte) error ***REMOVED***
				return nil
			***REMOVED***)

			// Node may be created from dead edge
			return nil
		***REMOVED***

		if pv := bkt.Get(bucketKeyParent); len(pv) > 0 ***REMOVED***
			fn(gcnode(ResourceSnapshot, node.Namespace, fmt.Sprintf("%s/%s", ss, pv)))
		***REMOVED***

		return sendSnapshotRefs(node.Namespace, bkt, fn)
	***REMOVED***

	return nil
***REMOVED***

func scanAll(ctx context.Context, tx *bolt.Tx, fn func(ctx context.Context, n gc.Node) error) error ***REMOVED***
	v1bkt := tx.Bucket(bucketKeyVersion)
	if v1bkt == nil ***REMOVED***
		return nil
	***REMOVED***

	// iterate through each namespace
	v1c := v1bkt.Cursor()

	for k, v := v1c.First(); k != nil; k, v = v1c.Next() ***REMOVED***
		if v != nil ***REMOVED***
			continue
		***REMOVED***
		nbkt := v1bkt.Bucket(k)
		ns := string(k)

		sbkt := nbkt.Bucket(bucketKeyObjectSnapshots)
		if sbkt != nil ***REMOVED***
			if err := sbkt.ForEach(func(sk, sv []byte) error ***REMOVED***
				if sv != nil ***REMOVED***
					return nil
				***REMOVED***
				snbkt := sbkt.Bucket(sk)
				return snbkt.ForEach(func(k, v []byte) error ***REMOVED***
					if v != nil ***REMOVED***
						return nil
					***REMOVED***
					node := gcnode(ResourceSnapshot, ns, fmt.Sprintf("%s/%s", sk, k))
					return fn(ctx, node)
				***REMOVED***)
			***REMOVED***); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		cbkt := nbkt.Bucket(bucketKeyObjectContent)
		if cbkt != nil ***REMOVED***
			cbkt = cbkt.Bucket(bucketKeyObjectBlob)
		***REMOVED***
		if cbkt != nil ***REMOVED***
			if err := cbkt.ForEach(func(k, v []byte) error ***REMOVED***
				if v != nil ***REMOVED***
					return nil
				***REMOVED***
				node := gcnode(ResourceContent, ns, string(k))
				return fn(ctx, node)
			***REMOVED***); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func remove(ctx context.Context, tx *bolt.Tx, node gc.Node) error ***REMOVED***
	v1bkt := tx.Bucket(bucketKeyVersion)
	if v1bkt == nil ***REMOVED***
		return nil
	***REMOVED***

	nsbkt := v1bkt.Bucket([]byte(node.Namespace))
	if nsbkt == nil ***REMOVED***
		return nil
	***REMOVED***

	switch node.Type ***REMOVED***
	case ResourceContent:
		cbkt := nsbkt.Bucket(bucketKeyObjectContent)
		if cbkt != nil ***REMOVED***
			cbkt = cbkt.Bucket(bucketKeyObjectBlob)
		***REMOVED***
		if cbkt != nil ***REMOVED***
			log.G(ctx).WithField("key", node.Key).Debug("remove content")
			return cbkt.DeleteBucket([]byte(node.Key))
		***REMOVED***
	case ResourceSnapshot:
		sbkt := nsbkt.Bucket(bucketKeyObjectSnapshots)
		if sbkt != nil ***REMOVED***
			parts := strings.SplitN(node.Key, "/", 2)
			if len(parts) != 2 ***REMOVED***
				return errors.Errorf("invalid snapshot gc key %s", node.Key)
			***REMOVED***
			ssbkt := sbkt.Bucket([]byte(parts[0]))
			if ssbkt != nil ***REMOVED***
				log.G(ctx).WithField("key", parts[1]).WithField("snapshotter", parts[0]).Debug("remove snapshot")
				return ssbkt.DeleteBucket([]byte(parts[1]))
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// sendSnapshotRefs sends all snapshot references referred to by the labels in the bkt
func sendSnapshotRefs(ns string, bkt *bolt.Bucket, fn func(gc.Node)) error ***REMOVED***
	lbkt := bkt.Bucket(bucketKeyObjectLabels)
	if lbkt != nil ***REMOVED***
		lc := lbkt.Cursor()

		for k, v := lc.Seek(labelGCSnapRef); k != nil && strings.HasPrefix(string(k), string(labelGCSnapRef)); k, v = lc.Next() ***REMOVED***
			snapshotter := string(k[len(labelGCSnapRef):])
			fn(gcnode(ResourceSnapshot, ns, fmt.Sprintf("%s/%s", snapshotter, v)))
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// sendContentRefs sends all content references referred to by the labels in the bkt
func sendContentRefs(ns string, bkt *bolt.Bucket, fn func(gc.Node)) error ***REMOVED***
	lbkt := bkt.Bucket(bucketKeyObjectLabels)
	if lbkt != nil ***REMOVED***
		lc := lbkt.Cursor()

		labelRef := string(labelGCContentRef)
		for k, v := lc.Seek(labelGCContentRef); k != nil && strings.HasPrefix(string(k), labelRef); k, v = lc.Next() ***REMOVED***
			if ks := string(k); ks != labelRef ***REMOVED***
				// Allow reference naming, ignore names
				if ks[len(labelRef)] != '.' ***REMOVED***
					continue
				***REMOVED***
			***REMOVED***

			fn(gcnode(ResourceContent, ns, string(v)))
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func isRootRef(bkt *bolt.Bucket) bool ***REMOVED***
	lbkt := bkt.Bucket(bucketKeyObjectLabels)
	if lbkt != nil ***REMOVED***
		rv := lbkt.Get(labelGCRoot)
		if rv != nil ***REMOVED***
			// TODO: interpret rv as a timestamp and skip if expired
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func sendRootRef(ctx context.Context, nc chan<- gc.Node, n gc.Node, bkt *bolt.Bucket) error ***REMOVED***
	if isRootRef(bkt) ***REMOVED***
		select ***REMOVED***
		case nc <- n:
		case <-ctx.Done():
			return ctx.Err()
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func gcnode(t gc.ResourceType, ns, key string) gc.Node ***REMOVED***
	return gc.Node***REMOVED***
		Type:      t,
		Namespace: ns,
		Key:       key,
	***REMOVED***
***REMOVED***
