package containerd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	containersapi "github.com/containerd/containerd/api/services/containers/v1"
	contentapi "github.com/containerd/containerd/api/services/content/v1"
	diffapi "github.com/containerd/containerd/api/services/diff/v1"
	eventsapi "github.com/containerd/containerd/api/services/events/v1"
	imagesapi "github.com/containerd/containerd/api/services/images/v1"
	introspectionapi "github.com/containerd/containerd/api/services/introspection/v1"
	namespacesapi "github.com/containerd/containerd/api/services/namespaces/v1"
	snapshotsapi "github.com/containerd/containerd/api/services/snapshots/v1"
	"github.com/containerd/containerd/api/services/tasks/v1"
	versionservice "github.com/containerd/containerd/api/services/version/v1"
	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/dialer"
	"github.com/containerd/containerd/diff"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/platforms"
	"github.com/containerd/containerd/plugin"
	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
	"github.com/containerd/containerd/remotes/docker/schema1"
	"github.com/containerd/containerd/snapshots"
	"github.com/containerd/typeurl"
	ptypes "github.com/gogo/protobuf/types"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func init() ***REMOVED***
	const prefix = "types.containerd.io"
	// register TypeUrls for commonly marshaled external types
	major := strconv.Itoa(specs.VersionMajor)
	typeurl.Register(&specs.Spec***REMOVED******REMOVED***, prefix, "opencontainers/runtime-spec", major, "Spec")
	typeurl.Register(&specs.Process***REMOVED******REMOVED***, prefix, "opencontainers/runtime-spec", major, "Process")
	typeurl.Register(&specs.LinuxResources***REMOVED******REMOVED***, prefix, "opencontainers/runtime-spec", major, "LinuxResources")
	typeurl.Register(&specs.WindowsResources***REMOVED******REMOVED***, prefix, "opencontainers/runtime-spec", major, "WindowsResources")
***REMOVED***

// New returns a new containerd client that is connected to the containerd
// instance provided by address
func New(address string, opts ...ClientOpt) (*Client, error) ***REMOVED***
	var copts clientOpts
	for _, o := range opts ***REMOVED***
		if err := o(&copts); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	gopts := []grpc.DialOption***REMOVED***
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithTimeout(60 * time.Second),
		grpc.FailOnNonTempDialError(true),
		grpc.WithBackoffMaxDelay(3 * time.Second),
		grpc.WithDialer(dialer.Dialer),
	***REMOVED***
	if len(copts.dialOptions) > 0 ***REMOVED***
		gopts = copts.dialOptions
	***REMOVED***
	if copts.defaultns != "" ***REMOVED***
		unary, stream := newNSInterceptors(copts.defaultns)
		gopts = append(gopts,
			grpc.WithUnaryInterceptor(unary),
			grpc.WithStreamInterceptor(stream),
		)
	***REMOVED***
	conn, err := grpc.Dial(dialer.DialAddress(address), gopts...)
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "failed to dial %q", address)
	***REMOVED***
	return NewWithConn(conn, opts...)
***REMOVED***

// NewWithConn returns a new containerd client that is connected to the containerd
// instance provided by the connection
func NewWithConn(conn *grpc.ClientConn, opts ...ClientOpt) (*Client, error) ***REMOVED***
	return &Client***REMOVED***
		conn:    conn,
		runtime: fmt.Sprintf("%s.%s", plugin.RuntimePlugin, runtime.GOOS),
	***REMOVED***, nil
***REMOVED***

// Client is the client to interact with containerd and its various services
// using a uniform interface
type Client struct ***REMOVED***
	conn    *grpc.ClientConn
	runtime string
***REMOVED***

// IsServing returns true if the client can successfully connect to the
// containerd daemon and the healthcheck service returns the SERVING
// response.
// This call will block if a transient error is encountered during
// connection. A timeout can be set in the context to ensure it returns
// early.
func (c *Client) IsServing(ctx context.Context) (bool, error) ***REMOVED***
	r, err := c.HealthService().Check(ctx, &grpc_health_v1.HealthCheckRequest***REMOVED******REMOVED***, grpc.FailFast(false))
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	return r.Status == grpc_health_v1.HealthCheckResponse_SERVING, nil
***REMOVED***

// Containers returns all containers created in containerd
func (c *Client) Containers(ctx context.Context, filters ...string) ([]Container, error) ***REMOVED***
	r, err := c.ContainerService().List(ctx, filters...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var out []Container
	for _, container := range r ***REMOVED***
		out = append(out, containerFromRecord(c, container))
	***REMOVED***
	return out, nil
***REMOVED***

// NewContainer will create a new container in container with the provided id
// the id must be unique within the namespace
func (c *Client) NewContainer(ctx context.Context, id string, opts ...NewContainerOpts) (Container, error) ***REMOVED***
	ctx, done, err := c.WithLease(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer done()

	container := containers.Container***REMOVED***
		ID: id,
		Runtime: containers.RuntimeInfo***REMOVED***
			Name: c.runtime,
		***REMOVED***,
	***REMOVED***
	for _, o := range opts ***REMOVED***
		if err := o(ctx, c, &container); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	r, err := c.ContainerService().Create(ctx, container)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return containerFromRecord(c, r), nil
***REMOVED***

// LoadContainer loads an existing container from metadata
func (c *Client) LoadContainer(ctx context.Context, id string) (Container, error) ***REMOVED***
	r, err := c.ContainerService().Get(ctx, id)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return containerFromRecord(c, r), nil
***REMOVED***

// RemoteContext is used to configure object resolutions and transfers with
// remote content stores and image providers.
type RemoteContext struct ***REMOVED***
	// Resolver is used to resolve names to objects, fetchers, and pushers.
	// If no resolver is provided, defaults to Docker registry resolver.
	Resolver remotes.Resolver

	// Unpack is done after an image is pulled to extract into a snapshotter.
	// If an image is not unpacked on pull, it can be unpacked any time
	// afterwards. Unpacking is required to run an image.
	Unpack bool

	// Snapshotter used for unpacking
	Snapshotter string

	// Labels to be applied to the created image
	Labels map[string]string

	// BaseHandlers are a set of handlers which get are called on dispatch.
	// These handlers always get called before any operation specific
	// handlers.
	BaseHandlers []images.Handler

	// ConvertSchema1 is whether to convert Docker registry schema 1
	// manifests. If this option is false then any image which resolves
	// to schema 1 will return an error since schema 1 is not supported.
	ConvertSchema1 bool
***REMOVED***

func defaultRemoteContext() *RemoteContext ***REMOVED***
	return &RemoteContext***REMOVED***
		Resolver: docker.NewResolver(docker.ResolverOptions***REMOVED***
			Client: http.DefaultClient,
		***REMOVED***),
		Snapshotter: DefaultSnapshotter,
	***REMOVED***
***REMOVED***

// Pull downloads the provided content into containerd's content store
func (c *Client) Pull(ctx context.Context, ref string, opts ...RemoteOpt) (Image, error) ***REMOVED***
	pullCtx := defaultRemoteContext()
	for _, o := range opts ***REMOVED***
		if err := o(c, pullCtx); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	store := c.ContentStore()

	ctx, done, err := c.WithLease(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer done()

	name, desc, err := pullCtx.Resolver.Resolve(ctx, ref)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	fetcher, err := pullCtx.Resolver.Fetcher(ctx, name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var (
		schema1Converter *schema1.Converter
		handler          images.Handler
	)
	if desc.MediaType == images.MediaTypeDockerSchema1Manifest && pullCtx.ConvertSchema1 ***REMOVED***
		schema1Converter = schema1.NewConverter(store, fetcher)
		handler = images.Handlers(append(pullCtx.BaseHandlers, schema1Converter)...)
	***REMOVED*** else ***REMOVED***
		handler = images.Handlers(append(pullCtx.BaseHandlers,
			remotes.FetchHandler(store, fetcher),
			images.ChildrenHandler(store, platforms.Default()))...,
		)
	***REMOVED***

	if err := images.Dispatch(ctx, handler, desc); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if schema1Converter != nil ***REMOVED***
		desc, err = schema1Converter.Convert(ctx)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	imgrec := images.Image***REMOVED***
		Name:   name,
		Target: desc,
		Labels: pullCtx.Labels,
	***REMOVED***

	is := c.ImageService()
	if created, err := is.Create(ctx, imgrec); err != nil ***REMOVED***
		if !errdefs.IsAlreadyExists(err) ***REMOVED***
			return nil, err
		***REMOVED***

		updated, err := is.Update(ctx, imgrec)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		imgrec = updated
	***REMOVED*** else ***REMOVED***
		imgrec = created
	***REMOVED***

	img := &image***REMOVED***
		client: c,
		i:      imgrec,
	***REMOVED***
	if pullCtx.Unpack ***REMOVED***
		if err := img.Unpack(ctx, pullCtx.Snapshotter); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return img, nil
***REMOVED***

// Push uploads the provided content to a remote resource
func (c *Client) Push(ctx context.Context, ref string, desc ocispec.Descriptor, opts ...RemoteOpt) error ***REMOVED***
	pushCtx := defaultRemoteContext()
	for _, o := range opts ***REMOVED***
		if err := o(c, pushCtx); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	pusher, err := pushCtx.Resolver.Pusher(ctx, ref)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var m sync.Mutex
	manifestStack := []ocispec.Descriptor***REMOVED******REMOVED***

	filterHandler := images.HandlerFunc(func(ctx context.Context, desc ocispec.Descriptor) ([]ocispec.Descriptor, error) ***REMOVED***
		switch desc.MediaType ***REMOVED***
		case images.MediaTypeDockerSchema2Manifest, ocispec.MediaTypeImageManifest,
			images.MediaTypeDockerSchema2ManifestList, ocispec.MediaTypeImageIndex:
			m.Lock()
			manifestStack = append(manifestStack, desc)
			m.Unlock()
			return nil, images.ErrStopHandler
		default:
			return nil, nil
		***REMOVED***
	***REMOVED***)

	cs := c.ContentStore()
	pushHandler := remotes.PushHandler(cs, pusher)

	handlers := append(pushCtx.BaseHandlers,
		images.ChildrenHandler(cs, platforms.Default()),
		filterHandler,
		pushHandler,
	)

	if err := images.Dispatch(ctx, images.Handlers(handlers...), desc); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Iterate in reverse order as seen, parent always uploaded after child
	for i := len(manifestStack) - 1; i >= 0; i-- ***REMOVED***
		_, err := pushHandler(ctx, manifestStack[i])
		if err != nil ***REMOVED***
			// TODO(estesp): until we have a more complete method for index push, we need to report
			// missing dependencies in an index/manifest list by sensing the "400 Bad Request"
			// as a marker for this problem
			if (manifestStack[i].MediaType == ocispec.MediaTypeImageIndex ||
				manifestStack[i].MediaType == images.MediaTypeDockerSchema2ManifestList) &&
				errors.Cause(err) != nil && strings.Contains(errors.Cause(err).Error(), "400 Bad Request") ***REMOVED***
				return errors.Wrap(err, "manifest list/index references to blobs and/or manifests are missing in your target registry")
			***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// GetImage returns an existing image
func (c *Client) GetImage(ctx context.Context, ref string) (Image, error) ***REMOVED***
	i, err := c.ImageService().Get(ctx, ref)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &image***REMOVED***
		client: c,
		i:      i,
	***REMOVED***, nil
***REMOVED***

// ListImages returns all existing images
func (c *Client) ListImages(ctx context.Context, filters ...string) ([]Image, error) ***REMOVED***
	imgs, err := c.ImageService().List(ctx, filters...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	images := make([]Image, len(imgs))
	for i, img := range imgs ***REMOVED***
		images[i] = &image***REMOVED***
			client: c,
			i:      img,
		***REMOVED***
	***REMOVED***
	return images, nil
***REMOVED***

// Subscribe to events that match one or more of the provided filters.
//
// Callers should listen on both the envelope channel and errs channel. If the
// errs channel returns nil or an error, the subscriber should terminate.
//
// To cancel shutdown reciept of events, cancel the provided context. The errs
// channel will be closed and return a nil error.
func (c *Client) Subscribe(ctx context.Context, filters ...string) (ch <-chan *eventsapi.Envelope, errs <-chan error) ***REMOVED***
	var (
		evq  = make(chan *eventsapi.Envelope)
		errq = make(chan error, 1)
	)

	errs = errq
	ch = evq

	session, err := c.EventService().Subscribe(ctx, &eventsapi.SubscribeRequest***REMOVED***
		Filters: filters,
	***REMOVED***)
	if err != nil ***REMOVED***
		errq <- err
		close(errq)
		return
	***REMOVED***

	go func() ***REMOVED***
		defer close(errq)

		for ***REMOVED***
			ev, err := session.Recv()
			if err != nil ***REMOVED***
				errq <- err
				return
			***REMOVED***

			select ***REMOVED***
			case evq <- ev:
			case <-ctx.Done():
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return ch, errs
***REMOVED***

// Close closes the clients connection to containerd
func (c *Client) Close() error ***REMOVED***
	return c.conn.Close()
***REMOVED***

// NamespaceService returns the underlying Namespaces Store
func (c *Client) NamespaceService() namespaces.Store ***REMOVED***
	return NewNamespaceStoreFromClient(namespacesapi.NewNamespacesClient(c.conn))
***REMOVED***

// ContainerService returns the underlying container Store
func (c *Client) ContainerService() containers.Store ***REMOVED***
	return NewRemoteContainerStore(containersapi.NewContainersClient(c.conn))
***REMOVED***

// ContentStore returns the underlying content Store
func (c *Client) ContentStore() content.Store ***REMOVED***
	return NewContentStoreFromClient(contentapi.NewContentClient(c.conn))
***REMOVED***

// SnapshotService returns the underlying snapshotter for the provided snapshotter name
func (c *Client) SnapshotService(snapshotterName string) snapshots.Snapshotter ***REMOVED***
	return NewSnapshotterFromClient(snapshotsapi.NewSnapshotsClient(c.conn), snapshotterName)
***REMOVED***

// TaskService returns the underlying TasksClient
func (c *Client) TaskService() tasks.TasksClient ***REMOVED***
	return tasks.NewTasksClient(c.conn)
***REMOVED***

// ImageService returns the underlying image Store
func (c *Client) ImageService() images.Store ***REMOVED***
	return NewImageStoreFromClient(imagesapi.NewImagesClient(c.conn))
***REMOVED***

// DiffService returns the underlying Differ
func (c *Client) DiffService() diff.Differ ***REMOVED***
	return NewDiffServiceFromClient(diffapi.NewDiffClient(c.conn))
***REMOVED***

// IntrospectionService returns the underlying Introspection Client
func (c *Client) IntrospectionService() introspectionapi.IntrospectionClient ***REMOVED***
	return introspectionapi.NewIntrospectionClient(c.conn)
***REMOVED***

// HealthService returns the underlying GRPC HealthClient
func (c *Client) HealthService() grpc_health_v1.HealthClient ***REMOVED***
	return grpc_health_v1.NewHealthClient(c.conn)
***REMOVED***

// EventService returns the underlying EventsClient
func (c *Client) EventService() eventsapi.EventsClient ***REMOVED***
	return eventsapi.NewEventsClient(c.conn)
***REMOVED***

// VersionService returns the underlying VersionClient
func (c *Client) VersionService() versionservice.VersionClient ***REMOVED***
	return versionservice.NewVersionClient(c.conn)
***REMOVED***

// Version of containerd
type Version struct ***REMOVED***
	// Version number
	Version string
	// Revision from git that was built
	Revision string
***REMOVED***

// Version returns the version of containerd that the client is connected to
func (c *Client) Version(ctx context.Context) (Version, error) ***REMOVED***
	response, err := c.VersionService().Version(ctx, &ptypes.Empty***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		return Version***REMOVED******REMOVED***, err
	***REMOVED***
	return Version***REMOVED***
		Version:  response.Version,
		Revision: response.Revision,
	***REMOVED***, nil
***REMOVED***

type importOpts struct ***REMOVED***
***REMOVED***

// ImportOpt allows the caller to specify import specific options
type ImportOpt func(c *importOpts) error

func resolveImportOpt(opts ...ImportOpt) (importOpts, error) ***REMOVED***
	var iopts importOpts
	for _, o := range opts ***REMOVED***
		if err := o(&iopts); err != nil ***REMOVED***
			return iopts, err
		***REMOVED***
	***REMOVED***
	return iopts, nil
***REMOVED***

// Import imports an image from a Tar stream using reader.
// Caller needs to specify importer. Future version may use oci.v1 as the default.
// Note that unreferrenced blobs may be imported to the content store as well.
func (c *Client) Import(ctx context.Context, importer images.Importer, reader io.Reader, opts ...ImportOpt) ([]Image, error) ***REMOVED***
	_, err := resolveImportOpt(opts...) // unused now
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ctx, done, err := c.WithLease(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer done()

	imgrecs, err := importer.Import(ctx, c.ContentStore(), reader)
	if err != nil ***REMOVED***
		// is.Update() is not called on error
		return nil, err
	***REMOVED***

	is := c.ImageService()
	var images []Image
	for _, imgrec := range imgrecs ***REMOVED***
		if updated, err := is.Update(ctx, imgrec, "target"); err != nil ***REMOVED***
			if !errdefs.IsNotFound(err) ***REMOVED***
				return nil, err
			***REMOVED***

			created, err := is.Create(ctx, imgrec)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			imgrec = created
		***REMOVED*** else ***REMOVED***
			imgrec = updated
		***REMOVED***

		images = append(images, &image***REMOVED***
			client: c,
			i:      imgrec,
		***REMOVED***)
	***REMOVED***
	return images, nil
***REMOVED***

type exportOpts struct ***REMOVED***
***REMOVED***

// ExportOpt allows the caller to specify export-specific options
type ExportOpt func(c *exportOpts) error

func resolveExportOpt(opts ...ExportOpt) (exportOpts, error) ***REMOVED***
	var eopts exportOpts
	for _, o := range opts ***REMOVED***
		if err := o(&eopts); err != nil ***REMOVED***
			return eopts, err
		***REMOVED***
	***REMOVED***
	return eopts, nil
***REMOVED***

// Export exports an image to a Tar stream.
// OCI format is used by default.
// It is up to caller to put "org.opencontainers.image.ref.name" annotation to desc.
// TODO(AkihiroSuda): support exporting multiple descriptors at once to a single archive stream.
func (c *Client) Export(ctx context.Context, exporter images.Exporter, desc ocispec.Descriptor, opts ...ExportOpt) (io.ReadCloser, error) ***REMOVED***
	_, err := resolveExportOpt(opts...) // unused now
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	pr, pw := io.Pipe()
	go func() ***REMOVED***
		pw.CloseWithError(exporter.Export(ctx, c.ContentStore(), desc, pw))
	***REMOVED***()
	return pr, nil
***REMOVED***
