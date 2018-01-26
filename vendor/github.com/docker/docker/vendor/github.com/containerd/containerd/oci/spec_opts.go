package oci

import (
	"context"

	"github.com/containerd/containerd/containers"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// SpecOpts sets spec specific information to a newly generated OCI spec
type SpecOpts func(context.Context, Client, *containers.Container, *specs.Spec) error

// WithProcessArgs replaces the args on the generated spec
func WithProcessArgs(args ...string) SpecOpts ***REMOVED***
	return func(_ context.Context, _ Client, _ *containers.Container, s *specs.Spec) error ***REMOVED***
		s.Process.Args = args
		return nil
	***REMOVED***
***REMOVED***

// WithProcessCwd replaces the current working directory on the generated spec
func WithProcessCwd(cwd string) SpecOpts ***REMOVED***
	return func(_ context.Context, _ Client, _ *containers.Container, s *specs.Spec) error ***REMOVED***
		s.Process.Cwd = cwd
		return nil
	***REMOVED***
***REMOVED***

// WithHostname sets the container's hostname
func WithHostname(name string) SpecOpts ***REMOVED***
	return func(_ context.Context, _ Client, _ *containers.Container, s *specs.Spec) error ***REMOVED***
		s.Hostname = name
		return nil
	***REMOVED***
***REMOVED***
