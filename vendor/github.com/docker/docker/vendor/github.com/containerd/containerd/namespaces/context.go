package namespaces

import (
	"os"

	"github.com/containerd/containerd/errdefs"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const (
	// NamespaceEnvVar is the environment variable key name
	NamespaceEnvVar = "CONTAINERD_NAMESPACE"
	// Default is the name of the default namespace
	Default = "default"
)

type namespaceKey struct***REMOVED******REMOVED***

// WithNamespace sets a given namespace on the context
func WithNamespace(ctx context.Context, namespace string) context.Context ***REMOVED***
	ctx = context.WithValue(ctx, namespaceKey***REMOVED******REMOVED***, namespace) // set our key for namespace

	// also store on the grpc headers so it gets picked up by any clients that
	// are using this.
	return withGRPCNamespaceHeader(ctx, namespace)
***REMOVED***

// NamespaceFromEnv uses the namespace defined in CONTAINERD_NAMESPACE or
// default
func NamespaceFromEnv(ctx context.Context) context.Context ***REMOVED***
	namespace := os.Getenv(NamespaceEnvVar)
	if namespace == "" ***REMOVED***
		namespace = Default
	***REMOVED***
	return WithNamespace(ctx, namespace)
***REMOVED***

// Namespace returns the namespace from the context.
//
// The namespace is not guaranteed to be valid.
func Namespace(ctx context.Context) (string, bool) ***REMOVED***
	namespace, ok := ctx.Value(namespaceKey***REMOVED******REMOVED***).(string)
	if !ok ***REMOVED***
		return fromGRPCHeader(ctx)
	***REMOVED***

	return namespace, ok
***REMOVED***

// NamespaceRequired returns the valid namepace from the context or an error.
func NamespaceRequired(ctx context.Context) (string, error) ***REMOVED***
	namespace, ok := Namespace(ctx)
	if !ok || namespace == "" ***REMOVED***
		return "", errors.Wrapf(errdefs.ErrFailedPrecondition, "namespace is required")
	***REMOVED***

	if err := Validate(namespace); err != nil ***REMOVED***
		return "", errors.Wrap(err, "namespace validation")
	***REMOVED***

	return namespace, nil
***REMOVED***
