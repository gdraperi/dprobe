package stream

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/containerd/containerd/cio"
	"github.com/docker/docker/pkg/broadcaster"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/pools"
	"github.com/sirupsen/logrus"
)

// Config holds information about I/O streams managed together.
//
// config.StdinPipe returns a WriteCloser which can be used to feed data
// to the standard input of the streamConfig's active process.
// config.StdoutPipe and streamConfig.StderrPipe each return a ReadCloser
// which can be used to retrieve the standard output (and error) generated
// by the container's active process. The output (and error) are actually
// copied and delivered to all StdoutPipe and StderrPipe consumers, using
// a kind of "broadcaster".
type Config struct ***REMOVED***
	sync.WaitGroup
	stdout    *broadcaster.Unbuffered
	stderr    *broadcaster.Unbuffered
	stdin     io.ReadCloser
	stdinPipe io.WriteCloser
***REMOVED***

// NewConfig creates a stream config and initializes
// the standard err and standard out to new unbuffered broadcasters.
func NewConfig() *Config ***REMOVED***
	return &Config***REMOVED***
		stderr: new(broadcaster.Unbuffered),
		stdout: new(broadcaster.Unbuffered),
	***REMOVED***
***REMOVED***

// Stdout returns the standard output in the configuration.
func (c *Config) Stdout() *broadcaster.Unbuffered ***REMOVED***
	return c.stdout
***REMOVED***

// Stderr returns the standard error in the configuration.
func (c *Config) Stderr() *broadcaster.Unbuffered ***REMOVED***
	return c.stderr
***REMOVED***

// Stdin returns the standard input in the configuration.
func (c *Config) Stdin() io.ReadCloser ***REMOVED***
	return c.stdin
***REMOVED***

// StdinPipe returns an input writer pipe as an io.WriteCloser.
func (c *Config) StdinPipe() io.WriteCloser ***REMOVED***
	return c.stdinPipe
***REMOVED***

// StdoutPipe creates a new io.ReadCloser with an empty bytes pipe.
// It adds this new out pipe to the Stdout broadcaster.
// This will block stdout if unconsumed.
func (c *Config) StdoutPipe() io.ReadCloser ***REMOVED***
	bytesPipe := ioutils.NewBytesPipe()
	c.stdout.Add(bytesPipe)
	return bytesPipe
***REMOVED***

// StderrPipe creates a new io.ReadCloser with an empty bytes pipe.
// It adds this new err pipe to the Stderr broadcaster.
// This will block stderr if unconsumed.
func (c *Config) StderrPipe() io.ReadCloser ***REMOVED***
	bytesPipe := ioutils.NewBytesPipe()
	c.stderr.Add(bytesPipe)
	return bytesPipe
***REMOVED***

// NewInputPipes creates new pipes for both standard inputs, Stdin and StdinPipe.
func (c *Config) NewInputPipes() ***REMOVED***
	c.stdin, c.stdinPipe = io.Pipe()
***REMOVED***

// NewNopInputPipe creates a new input pipe that will silently drop all messages in the input.
func (c *Config) NewNopInputPipe() ***REMOVED***
	c.stdinPipe = ioutils.NopWriteCloser(ioutil.Discard)
***REMOVED***

// CloseStreams ensures that the configured streams are properly closed.
func (c *Config) CloseStreams() error ***REMOVED***
	var errors []string

	if c.stdin != nil ***REMOVED***
		if err := c.stdin.Close(); err != nil ***REMOVED***
			errors = append(errors, fmt.Sprintf("error close stdin: %s", err))
		***REMOVED***
	***REMOVED***

	if err := c.stdout.Clean(); err != nil ***REMOVED***
		errors = append(errors, fmt.Sprintf("error close stdout: %s", err))
	***REMOVED***

	if err := c.stderr.Clean(); err != nil ***REMOVED***
		errors = append(errors, fmt.Sprintf("error close stderr: %s", err))
	***REMOVED***

	if len(errors) > 0 ***REMOVED***
		return fmt.Errorf(strings.Join(errors, "\n"))
	***REMOVED***

	return nil
***REMOVED***

// CopyToPipe connects streamconfig with a libcontainerd.IOPipe
func (c *Config) CopyToPipe(iop *cio.DirectIO) ***REMOVED***
	copyFunc := func(w io.Writer, r io.ReadCloser) ***REMOVED***
		c.Add(1)
		go func() ***REMOVED***
			if _, err := pools.Copy(w, r); err != nil ***REMOVED***
				logrus.Errorf("stream copy error: %v", err)
			***REMOVED***
			r.Close()
			c.Done()
		***REMOVED***()
	***REMOVED***

	if iop.Stdout != nil ***REMOVED***
		copyFunc(c.Stdout(), iop.Stdout)
	***REMOVED***
	if iop.Stderr != nil ***REMOVED***
		copyFunc(c.Stderr(), iop.Stderr)
	***REMOVED***

	if stdin := c.Stdin(); stdin != nil ***REMOVED***
		if iop.Stdin != nil ***REMOVED***
			go func() ***REMOVED***
				pools.Copy(iop.Stdin, stdin)
				if err := iop.Stdin.Close(); err != nil ***REMOVED***
					logrus.Warnf("failed to close stdin: %v", err)
				***REMOVED***
			***REMOVED***()
		***REMOVED***
	***REMOVED***
***REMOVED***
