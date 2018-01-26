package containerd

import (
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/remotes"
	"google.golang.org/grpc"
)

type clientOpts struct ***REMOVED***
	defaultns   string
	dialOptions []grpc.DialOption
***REMOVED***

// ClientOpt allows callers to set options on the containerd client
type ClientOpt func(c *clientOpts) error

// WithDefaultNamespace sets the default namespace on the client
//
// Any operation that does not have a namespace set on the context will
// be provided the default namespace
func WithDefaultNamespace(ns string) ClientOpt ***REMOVED***
	return func(c *clientOpts) error ***REMOVED***
		c.defaultns = ns
		return nil
	***REMOVED***
***REMOVED***

// WithDialOpts allows grpc.DialOptions to be set on the connection
func WithDialOpts(opts []grpc.DialOption) ClientOpt ***REMOVED***
	return func(c *clientOpts) error ***REMOVED***
		c.dialOptions = opts
		return nil
	***REMOVED***
***REMOVED***

// RemoteOpt allows the caller to set distribution options for a remote
type RemoteOpt func(*Client, *RemoteContext) error

// WithPullUnpack is used to unpack an image after pull. This
// uses the snapshotter, content store, and diff service
// configured for the client.
func WithPullUnpack(_ *Client, c *RemoteContext) error ***REMOVED***
	c.Unpack = true
	return nil
***REMOVED***

// WithPullSnapshotter specifies snapshotter name used for unpacking
func WithPullSnapshotter(snapshotterName string) RemoteOpt ***REMOVED***
	return func(_ *Client, c *RemoteContext) error ***REMOVED***
		c.Snapshotter = snapshotterName
		return nil
	***REMOVED***
***REMOVED***

// WithPullLabel sets a label to be associated with a pulled reference
func WithPullLabel(key, value string) RemoteOpt ***REMOVED***
	return func(_ *Client, rc *RemoteContext) error ***REMOVED***
		if rc.Labels == nil ***REMOVED***
			rc.Labels = make(map[string]string)
		***REMOVED***

		rc.Labels[key] = value
		return nil
	***REMOVED***
***REMOVED***

// WithPullLabels associates a set of labels to a pulled reference
func WithPullLabels(labels map[string]string) RemoteOpt ***REMOVED***
	return func(_ *Client, rc *RemoteContext) error ***REMOVED***
		if rc.Labels == nil ***REMOVED***
			rc.Labels = make(map[string]string)
		***REMOVED***

		for k, v := range labels ***REMOVED***
			rc.Labels[k] = v
		***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// WithSchema1Conversion is used to convert Docker registry schema 1
// manifests to oci manifests on pull. Without this option schema 1
// manifests will return a not supported error.
func WithSchema1Conversion(client *Client, c *RemoteContext) error ***REMOVED***
	c.ConvertSchema1 = true
	return nil
***REMOVED***

// WithResolver specifies the resolver to use.
func WithResolver(resolver remotes.Resolver) RemoteOpt ***REMOVED***
	return func(client *Client, c *RemoteContext) error ***REMOVED***
		c.Resolver = resolver
		return nil
	***REMOVED***
***REMOVED***

// WithImageHandler adds a base handler to be called on dispatch.
func WithImageHandler(h images.Handler) RemoteOpt ***REMOVED***
	return func(client *Client, c *RemoteContext) error ***REMOVED***
		c.BaseHandlers = append(c.BaseHandlers, h)
		return nil
	***REMOVED***
***REMOVED***
