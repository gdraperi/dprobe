package containerd

import (
	"context"

	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/oci"
	"github.com/containerd/containerd/platforms"
	"github.com/containerd/typeurl"
	"github.com/gogo/protobuf/types"
	"github.com/opencontainers/image-spec/identity"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
)

// DeleteOpts allows the caller to set options for the deletion of a container
type DeleteOpts func(ctx context.Context, client *Client, c containers.Container) error

// NewContainerOpts allows the caller to set additional options when creating a container
type NewContainerOpts func(ctx context.Context, client *Client, c *containers.Container) error

// UpdateContainerOpts allows the caller to set additional options when updating a container
type UpdateContainerOpts func(ctx context.Context, client *Client, c *containers.Container) error

// WithRuntime allows a user to specify the runtime name and additional options that should
// be used to create tasks for the container
func WithRuntime(name string, options interface***REMOVED******REMOVED***) NewContainerOpts ***REMOVED***
	return func(ctx context.Context, client *Client, c *containers.Container) error ***REMOVED***
		var (
			any *types.Any
			err error
		)
		if options != nil ***REMOVED***
			any, err = typeurl.MarshalAny(options)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		c.Runtime = containers.RuntimeInfo***REMOVED***
			Name:    name,
			Options: any,
		***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// WithImage sets the provided image as the base for the container
func WithImage(i Image) NewContainerOpts ***REMOVED***
	return func(ctx context.Context, client *Client, c *containers.Container) error ***REMOVED***
		c.Image = i.Name()
		return nil
	***REMOVED***
***REMOVED***

// WithContainerLabels adds the provided labels to the container
func WithContainerLabels(labels map[string]string) NewContainerOpts ***REMOVED***
	return func(_ context.Context, _ *Client, c *containers.Container) error ***REMOVED***
		c.Labels = labels
		return nil
	***REMOVED***
***REMOVED***

// WithSnapshotter sets the provided snapshotter for use by the container
//
// This option must appear before other snapshotter options to have an effect.
func WithSnapshotter(name string) NewContainerOpts ***REMOVED***
	return func(ctx context.Context, client *Client, c *containers.Container) error ***REMOVED***
		c.Snapshotter = name
		return nil
	***REMOVED***
***REMOVED***

// WithSnapshot uses an existing root filesystem for the container
func WithSnapshot(id string) NewContainerOpts ***REMOVED***
	return func(ctx context.Context, client *Client, c *containers.Container) error ***REMOVED***
		setSnapshotterIfEmpty(c)
		// check that the snapshot exists, if not, fail on creation
		if _, err := client.SnapshotService(c.Snapshotter).Mounts(ctx, id); err != nil ***REMOVED***
			return err
		***REMOVED***
		c.SnapshotKey = id
		return nil
	***REMOVED***
***REMOVED***

// WithNewSnapshot allocates a new snapshot to be used by the container as the
// root filesystem in read-write mode
func WithNewSnapshot(id string, i Image) NewContainerOpts ***REMOVED***
	return func(ctx context.Context, client *Client, c *containers.Container) error ***REMOVED***
		diffIDs, err := i.(*image).i.RootFS(ctx, client.ContentStore(), platforms.Default())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		setSnapshotterIfEmpty(c)
		parent := identity.ChainID(diffIDs).String()
		if _, err := client.SnapshotService(c.Snapshotter).Prepare(ctx, id, parent); err != nil ***REMOVED***
			return err
		***REMOVED***
		c.SnapshotKey = id
		c.Image = i.Name()
		return nil
	***REMOVED***
***REMOVED***

// WithSnapshotCleanup deletes the rootfs snapshot allocated for the container
func WithSnapshotCleanup(ctx context.Context, client *Client, c containers.Container) error ***REMOVED***
	if c.SnapshotKey != "" ***REMOVED***
		if c.Snapshotter == "" ***REMOVED***
			return errors.Wrapf(errdefs.ErrInvalidArgument, "container.Snapshotter must be set to cleanup rootfs snapshot")
		***REMOVED***
		return client.SnapshotService(c.Snapshotter).Remove(ctx, c.SnapshotKey)
	***REMOVED***
	return nil
***REMOVED***

// WithNewSnapshotView allocates a new snapshot to be used by the container as the
// root filesystem in read-only mode
func WithNewSnapshotView(id string, i Image) NewContainerOpts ***REMOVED***
	return func(ctx context.Context, client *Client, c *containers.Container) error ***REMOVED***
		diffIDs, err := i.(*image).i.RootFS(ctx, client.ContentStore(), platforms.Default())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		setSnapshotterIfEmpty(c)
		parent := identity.ChainID(diffIDs).String()
		if _, err := client.SnapshotService(c.Snapshotter).View(ctx, id, parent); err != nil ***REMOVED***
			return err
		***REMOVED***
		c.SnapshotKey = id
		c.Image = i.Name()
		return nil
	***REMOVED***
***REMOVED***

func setSnapshotterIfEmpty(c *containers.Container) ***REMOVED***
	if c.Snapshotter == "" ***REMOVED***
		c.Snapshotter = DefaultSnapshotter
	***REMOVED***
***REMOVED***

// WithContainerExtension appends extension data to the container object.
// Use this to decorate the container object with additional data for the client
// integration.
//
// Make sure to register the type of `extension` in the typeurl package via
// `typeurl.Register` or container creation may fail.
func WithContainerExtension(name string, extension interface***REMOVED******REMOVED***) NewContainerOpts ***REMOVED***
	return func(ctx context.Context, client *Client, c *containers.Container) error ***REMOVED***
		if name == "" ***REMOVED***
			return errors.Wrapf(errdefs.ErrInvalidArgument, "extension key must not be zero-length")
		***REMOVED***

		any, err := typeurl.MarshalAny(extension)
		if err != nil ***REMOVED***
			if errors.Cause(err) == typeurl.ErrNotFound ***REMOVED***
				return errors.Wrapf(err, "extension %q is not registered with the typeurl package, see `typeurl.Register`", name)
			***REMOVED***
			return errors.Wrap(err, "error marshalling extension")
		***REMOVED***

		if c.Extensions == nil ***REMOVED***
			c.Extensions = make(map[string]types.Any)
		***REMOVED***
		c.Extensions[name] = *any
		return nil
	***REMOVED***
***REMOVED***

// WithNewSpec generates a new spec for a new container
func WithNewSpec(opts ...oci.SpecOpts) NewContainerOpts ***REMOVED***
	return func(ctx context.Context, client *Client, c *containers.Container) error ***REMOVED***
		s, err := oci.GenerateSpec(ctx, client, c, opts...)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		c.Spec, err = typeurl.MarshalAny(s)
		return err
	***REMOVED***
***REMOVED***

// WithSpec sets the provided spec on the container
func WithSpec(s *specs.Spec, opts ...oci.SpecOpts) NewContainerOpts ***REMOVED***
	return func(ctx context.Context, client *Client, c *containers.Container) error ***REMOVED***
		for _, o := range opts ***REMOVED***
			if err := o(ctx, client, c, s); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		var err error
		c.Spec, err = typeurl.MarshalAny(s)
		return err
	***REMOVED***
***REMOVED***
