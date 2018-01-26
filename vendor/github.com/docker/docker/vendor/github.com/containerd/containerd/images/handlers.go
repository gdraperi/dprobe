package images

import (
	"context"
	"fmt"

	"github.com/containerd/containerd/content"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var (
	// ErrSkipDesc is used to skip processing of a descriptor and
	// its descendants.
	ErrSkipDesc = fmt.Errorf("skip descriptor")

	// ErrStopHandler is used to signify that the descriptor
	// has been handled and should not be handled further.
	// This applies only to a single descriptor in a handler
	// chain and does not apply to descendant descriptors.
	ErrStopHandler = fmt.Errorf("stop handler")
)

// Handler handles image manifests
type Handler interface ***REMOVED***
	Handle(ctx context.Context, desc ocispec.Descriptor) (subdescs []ocispec.Descriptor, err error)
***REMOVED***

// HandlerFunc function implementing the Handler interface
type HandlerFunc func(ctx context.Context, desc ocispec.Descriptor) (subdescs []ocispec.Descriptor, err error)

// Handle image manifests
func (fn HandlerFunc) Handle(ctx context.Context, desc ocispec.Descriptor) (subdescs []ocispec.Descriptor, err error) ***REMOVED***
	return fn(ctx, desc)
***REMOVED***

// Handlers returns a handler that will run the handlers in sequence.
//
// A handler may return `ErrStopHandler` to stop calling additional handlers
func Handlers(handlers ...Handler) HandlerFunc ***REMOVED***
	return func(ctx context.Context, desc ocispec.Descriptor) (subdescs []ocispec.Descriptor, err error) ***REMOVED***
		var children []ocispec.Descriptor
		for _, handler := range handlers ***REMOVED***
			ch, err := handler.Handle(ctx, desc)
			if err != nil ***REMOVED***
				if errors.Cause(err) == ErrStopHandler ***REMOVED***
					break
				***REMOVED***
				return nil, err
			***REMOVED***

			children = append(children, ch...)
		***REMOVED***

		return children, nil
	***REMOVED***
***REMOVED***

// Walk the resources of an image and call the handler for each. If the handler
// decodes the sub-resources for each image,
//
// This differs from dispatch in that each sibling resource is considered
// synchronously.
func Walk(ctx context.Context, handler Handler, descs ...ocispec.Descriptor) error ***REMOVED***
	for _, desc := range descs ***REMOVED***

		children, err := handler.Handle(ctx, desc)
		if err != nil ***REMOVED***
			if errors.Cause(err) == ErrSkipDesc ***REMOVED***
				continue // don't traverse the children.
			***REMOVED***
			return err
		***REMOVED***

		if len(children) > 0 ***REMOVED***
			if err := Walk(ctx, handler, children...); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// Dispatch runs the provided handler for content specified by the descriptors.
// If the handler decode subresources, they will be visited, as well.
//
// Handlers for siblings are run in parallel on the provided descriptors. A
// handler may return `ErrSkipDesc` to signal to the dispatcher to not traverse
// any children.
//
// Typically, this function will be used with `FetchHandler`, often composed
// with other handlers.
//
// If any handler returns an error, the dispatch session will be canceled.
func Dispatch(ctx context.Context, handler Handler, descs ...ocispec.Descriptor) error ***REMOVED***
	eg, ctx := errgroup.WithContext(ctx)
	for _, desc := range descs ***REMOVED***
		desc := desc

		eg.Go(func() error ***REMOVED***
			desc := desc

			children, err := handler.Handle(ctx, desc)
			if err != nil ***REMOVED***
				if errors.Cause(err) == ErrSkipDesc ***REMOVED***
					return nil // don't traverse the children.
				***REMOVED***
				return err
			***REMOVED***

			if len(children) > 0 ***REMOVED***
				return Dispatch(ctx, handler, children...)
			***REMOVED***

			return nil
		***REMOVED***)
	***REMOVED***

	return eg.Wait()
***REMOVED***

// ChildrenHandler decodes well-known manifest types and returns their children.
//
// This is useful for supporting recursive fetch and other use cases where you
// want to do a full walk of resources.
//
// One can also replace this with another implementation to allow descending of
// arbitrary types.
func ChildrenHandler(provider content.Provider, platform string) HandlerFunc ***REMOVED***
	return func(ctx context.Context, desc ocispec.Descriptor) ([]ocispec.Descriptor, error) ***REMOVED***
		return Children(ctx, provider, desc, platform)
	***REMOVED***
***REMOVED***
