package backend

import (
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/streamformatter"
)

// PullOption defines different modes for accessing images
type PullOption int

const (
	// PullOptionNoPull only returns local images
	PullOptionNoPull PullOption = iota
	// PullOptionForcePull always tries to pull a ref from the registry first
	PullOptionForcePull
	// PullOptionPreferLocal uses local image if it exists, otherwise pulls
	PullOptionPreferLocal
)

// ProgressWriter is a data object to transport progress streams to the client
type ProgressWriter struct ***REMOVED***
	Output             io.Writer
	StdoutFormatter    io.Writer
	StderrFormatter    io.Writer
	AuxFormatter       *streamformatter.AuxFormatter
	ProgressReaderFunc func(io.ReadCloser) io.ReadCloser
***REMOVED***

// BuildConfig is the configuration used by a BuildManager to start a build
type BuildConfig struct ***REMOVED***
	Source         io.ReadCloser
	ProgressWriter ProgressWriter
	Options        *types.ImageBuildOptions
***REMOVED***

// GetImageAndLayerOptions are the options supported by GetImageAndReleasableLayer
type GetImageAndLayerOptions struct ***REMOVED***
	PullOption PullOption
	AuthConfig map[string]types.AuthConfig
	Output     io.Writer
	OS         string
***REMOVED***
