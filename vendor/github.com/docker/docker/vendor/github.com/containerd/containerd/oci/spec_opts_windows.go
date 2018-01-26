// +build windows

package oci

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/images"
	"github.com/opencontainers/image-spec/specs-go/v1"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// WithImageConfig configures the spec to from the configuration of an Image
func WithImageConfig(image Image) SpecOpts ***REMOVED***
	return func(ctx context.Context, client Client, _ *containers.Container, s *specs.Spec) error ***REMOVED***
		ic, err := image.Config(ctx)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		var (
			ociimage v1.Image
			config   v1.ImageConfig
		)
		switch ic.MediaType ***REMOVED***
		case v1.MediaTypeImageConfig, images.MediaTypeDockerSchema2Config:
			p, err := content.ReadBlob(ctx, image.ContentStore(), ic.Digest)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if err := json.Unmarshal(p, &ociimage); err != nil ***REMOVED***
				return err
			***REMOVED***
			config = ociimage.Config
		default:
			return fmt.Errorf("unknown image config media type %s", ic.MediaType)
		***REMOVED***
		s.Process.Env = config.Env
		s.Process.Args = append(config.Entrypoint, config.Cmd...)
		s.Process.User = specs.User***REMOVED***
			Username: config.User,
		***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// WithTTY sets the information on the spec as well as the environment variables for
// using a TTY
func WithTTY(width, height int) SpecOpts ***REMOVED***
	return func(_ context.Context, _ Client, _ *containers.Container, s *specs.Spec) error ***REMOVED***
		s.Process.Terminal = true
		if s.Process.ConsoleSize == nil ***REMOVED***
			s.Process.ConsoleSize = &specs.Box***REMOVED******REMOVED***
		***REMOVED***
		s.Process.ConsoleSize.Width = uint(width)
		s.Process.ConsoleSize.Height = uint(height)
		return nil
	***REMOVED***
***REMOVED***

// WithUsername sets the username on the process
func WithUsername(username string) SpecOpts ***REMOVED***
	return func(ctx context.Context, client Client, c *containers.Container, s *specs.Spec) error ***REMOVED***
		s.Process.User.Username = username
		return nil
	***REMOVED***
***REMOVED***
