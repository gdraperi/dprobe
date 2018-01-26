package build

import (
	"io"
	"strings"

	"github.com/docker/docker/integration-cli/cli/build/fakecontext"
	"github.com/gotestyourself/gotestyourself/icmd"
)

type testingT interface ***REMOVED***
	Fatal(args ...interface***REMOVED******REMOVED***)
	Fatalf(string, ...interface***REMOVED******REMOVED***)
***REMOVED***

// WithStdinContext sets the build context from the standard input with the specified reader
func WithStdinContext(closer io.ReadCloser) func(*icmd.Cmd) func() ***REMOVED***
	return func(cmd *icmd.Cmd) func() ***REMOVED***
		cmd.Command = append(cmd.Command, "-")
		cmd.Stdin = closer
		return func() ***REMOVED***
			// FIXME(vdemeester) we should not ignore the error here…
			closer.Close()
		***REMOVED***
	***REMOVED***
***REMOVED***

// WithDockerfile creates / returns a CmdOperator to set the Dockerfile for a build operation
func WithDockerfile(dockerfile string) func(*icmd.Cmd) func() ***REMOVED***
	return func(cmd *icmd.Cmd) func() ***REMOVED***
		cmd.Command = append(cmd.Command, "-")
		cmd.Stdin = strings.NewReader(dockerfile)
		return nil
	***REMOVED***
***REMOVED***

// WithoutCache makes the build ignore cache
func WithoutCache(cmd *icmd.Cmd) func() ***REMOVED***
	cmd.Command = append(cmd.Command, "--no-cache")
	return nil
***REMOVED***

// WithContextPath sets the build context path
func WithContextPath(path string) func(*icmd.Cmd) func() ***REMOVED***
	return func(cmd *icmd.Cmd) func() ***REMOVED***
		cmd.Command = append(cmd.Command, path)
		return nil
	***REMOVED***
***REMOVED***

// WithExternalBuildContext use the specified context as build context
func WithExternalBuildContext(ctx *fakecontext.Fake) func(*icmd.Cmd) func() ***REMOVED***
	return func(cmd *icmd.Cmd) func() ***REMOVED***
		cmd.Dir = ctx.Dir
		cmd.Command = append(cmd.Command, ".")
		return nil
	***REMOVED***
***REMOVED***

// WithBuildContext sets up the build context
func WithBuildContext(t testingT, contextOperators ...func(*fakecontext.Fake) error) func(*icmd.Cmd) func() ***REMOVED***
	// FIXME(vdemeester) de-duplicate that
	ctx := fakecontext.New(t, "", contextOperators...)
	return func(cmd *icmd.Cmd) func() ***REMOVED***
		cmd.Dir = ctx.Dir
		cmd.Command = append(cmd.Command, ".")
		return closeBuildContext(t, ctx)
	***REMOVED***
***REMOVED***

// WithFile adds the specified file (with content) in the build context
func WithFile(name, content string) func(*fakecontext.Fake) error ***REMOVED***
	return fakecontext.WithFile(name, content)
***REMOVED***

func closeBuildContext(t testingT, ctx *fakecontext.Fake) func() ***REMOVED***
	return func() ***REMOVED***
		if err := ctx.Close(); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***
