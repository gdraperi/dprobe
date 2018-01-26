package cio

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
)

// Config holds the IO configurations.
type Config struct ***REMOVED***
	// Terminal is true if one has been allocated
	Terminal bool
	// Stdin path
	Stdin string
	// Stdout path
	Stdout string
	// Stderr path
	Stderr string
***REMOVED***

// IO holds the io information for a task or process
type IO interface ***REMOVED***
	// Config returns the IO configuration.
	Config() Config
	// Cancel aborts all current io operations.
	Cancel()
	// Wait blocks until all io copy operations have completed.
	Wait()
	// Close cleans up all open io resources. Cancel() is always called before
	// Close()
	Close() error
***REMOVED***

// Creator creates new IO sets for a task
type Creator func(id string) (IO, error)

// Attach allows callers to reattach to running tasks
//
// There should only be one reader for a task's IO set
// because fifo's can only be read from one reader or the output
// will be sent only to the first reads
type Attach func(*FIFOSet) (IO, error)

// FIFOSet is a set of file paths to FIFOs for a task's standard IO streams
type FIFOSet struct ***REMOVED***
	Config
	close func() error
***REMOVED***

// Close the FIFOSet
func (f *FIFOSet) Close() error ***REMOVED***
	if f.close != nil ***REMOVED***
		return f.close()
	***REMOVED***
	return nil
***REMOVED***

// NewFIFOSet returns a new FIFOSet from a Config and a close function
func NewFIFOSet(config Config, close func() error) *FIFOSet ***REMOVED***
	return &FIFOSet***REMOVED***Config: config, close: close***REMOVED***
***REMOVED***

// Streams used to configure a Creator or Attach
type Streams struct ***REMOVED***
	Stdin    io.Reader
	Stdout   io.Writer
	Stderr   io.Writer
	Terminal bool
***REMOVED***

// Opt customize options for creating a Creator or Attach
type Opt func(*Streams)

// WithStdio sets stream options to the standard input/output streams
func WithStdio(opt *Streams) ***REMOVED***
	WithStreams(os.Stdin, os.Stdout, os.Stderr)(opt)
***REMOVED***

// WithTerminal sets the terminal option
func WithTerminal(opt *Streams) ***REMOVED***
	opt.Terminal = true
***REMOVED***

// WithStreams sets the stream options to the specified Reader and Writers
func WithStreams(stdin io.Reader, stdout, stderr io.Writer) Opt ***REMOVED***
	return func(opt *Streams) ***REMOVED***
		opt.Stdin = stdin
		opt.Stdout = stdout
		opt.Stderr = stderr
	***REMOVED***
***REMOVED***

// NewCreator returns an IO creator from the options
func NewCreator(opts ...Opt) Creator ***REMOVED***
	streams := &Streams***REMOVED******REMOVED***
	for _, opt := range opts ***REMOVED***
		opt(streams)
	***REMOVED***
	return func(id string) (IO, error) ***REMOVED***
		// TODO: accept root as a param
		root := "/run/containerd/fifo"
		fifos, err := NewFIFOSetInDir(root, id, streams.Terminal)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return copyIO(fifos, streams)
	***REMOVED***
***REMOVED***

// NewAttach attaches the existing io for a task to the provided io.Reader/Writers
func NewAttach(opts ...Opt) Attach ***REMOVED***
	streams := &Streams***REMOVED******REMOVED***
	for _, opt := range opts ***REMOVED***
		opt(streams)
	***REMOVED***
	return func(fifos *FIFOSet) (IO, error) ***REMOVED***
		if fifos == nil ***REMOVED***
			return nil, fmt.Errorf("cannot attach, missing fifos")
		***REMOVED***
		return copyIO(fifos, streams)
	***REMOVED***
***REMOVED***

// NullIO redirects the container's IO into /dev/null
func NullIO(_ string) (IO, error) ***REMOVED***
	return &cio***REMOVED******REMOVED***, nil
***REMOVED***

// cio is a basic container IO implementation.
type cio struct ***REMOVED***
	config  Config
	wg      *sync.WaitGroup
	closers []io.Closer
	cancel  context.CancelFunc
***REMOVED***

func (c *cio) Config() Config ***REMOVED***
	return c.config
***REMOVED***

func (c *cio) Wait() ***REMOVED***
	if c.wg != nil ***REMOVED***
		c.wg.Wait()
	***REMOVED***
***REMOVED***

func (c *cio) Close() error ***REMOVED***
	var lastErr error
	for _, closer := range c.closers ***REMOVED***
		if closer == nil ***REMOVED***
			continue
		***REMOVED***
		if err := closer.Close(); err != nil ***REMOVED***
			lastErr = err
		***REMOVED***
	***REMOVED***
	return lastErr
***REMOVED***

func (c *cio) Cancel() ***REMOVED***
	if c.cancel != nil ***REMOVED***
		c.cancel()
	***REMOVED***
***REMOVED***

type pipes struct ***REMOVED***
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser
***REMOVED***

// DirectIO allows task IO to be handled externally by the caller
type DirectIO struct ***REMOVED***
	pipes
	cio
***REMOVED***

var _ IO = &DirectIO***REMOVED******REMOVED***
