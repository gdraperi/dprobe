package stream

import (
	"io"
	"sync"

	"golang.org/x/net/context"

	"github.com/docker/docker/pkg/pools"
	"github.com/docker/docker/pkg/term"
	"github.com/sirupsen/logrus"
)

var defaultEscapeSequence = []byte***REMOVED***16, 17***REMOVED*** // ctrl-p, ctrl-q

// AttachConfig is the config struct used to attach a client to a stream's stdio
type AttachConfig struct ***REMOVED***
	// Tells the attach copier that the stream's stdin is a TTY and to look for
	// escape sequences in stdin to detach from the stream.
	// When true the escape sequence is not passed to the underlying stream
	TTY bool
	// Specifies the detach keys the client will be using
	// Only useful when `TTY` is true
	DetachKeys []byte

	// CloseStdin signals that once done, stdin for the attached stream should be closed
	// For example, this would close the attached container's stdin.
	CloseStdin bool

	// UseStd* indicate whether the client has requested to be connected to the
	// given stream or not.  These flags are used instead of checking Std* != nil
	// at points before the client streams Std* are wired up.
	UseStdin, UseStdout, UseStderr bool

	// CStd* are the streams directly connected to the container
	CStdin           io.WriteCloser
	CStdout, CStderr io.ReadCloser

	// Provide client streams to wire up to
	Stdin          io.ReadCloser
	Stdout, Stderr io.Writer
***REMOVED***

// AttachStreams attaches the container's streams to the AttachConfig
func (c *Config) AttachStreams(cfg *AttachConfig) ***REMOVED***
	if cfg.UseStdin ***REMOVED***
		cfg.CStdin = c.StdinPipe()
	***REMOVED***

	if cfg.UseStdout ***REMOVED***
		cfg.CStdout = c.StdoutPipe()
	***REMOVED***

	if cfg.UseStderr ***REMOVED***
		cfg.CStderr = c.StderrPipe()
	***REMOVED***
***REMOVED***

// CopyStreams starts goroutines to copy data in and out to/from the container
func (c *Config) CopyStreams(ctx context.Context, cfg *AttachConfig) <-chan error ***REMOVED***
	var (
		wg     sync.WaitGroup
		errors = make(chan error, 3)
	)

	if cfg.Stdin != nil ***REMOVED***
		wg.Add(1)
	***REMOVED***

	if cfg.Stdout != nil ***REMOVED***
		wg.Add(1)
	***REMOVED***

	if cfg.Stderr != nil ***REMOVED***
		wg.Add(1)
	***REMOVED***

	// Connect stdin of container to the attach stdin stream.
	go func() ***REMOVED***
		if cfg.Stdin == nil ***REMOVED***
			return
		***REMOVED***
		logrus.Debug("attach: stdin: begin")

		var err error
		if cfg.TTY ***REMOVED***
			_, err = copyEscapable(cfg.CStdin, cfg.Stdin, cfg.DetachKeys)
		***REMOVED*** else ***REMOVED***
			_, err = pools.Copy(cfg.CStdin, cfg.Stdin)
		***REMOVED***
		if err == io.ErrClosedPipe ***REMOVED***
			err = nil
		***REMOVED***
		if err != nil ***REMOVED***
			logrus.Errorf("attach: stdin: %s", err)
			errors <- err
		***REMOVED***
		if cfg.CloseStdin && !cfg.TTY ***REMOVED***
			cfg.CStdin.Close()
		***REMOVED*** else ***REMOVED***
			// No matter what, when stdin is closed (io.Copy unblock), close stdout and stderr
			if cfg.CStdout != nil ***REMOVED***
				cfg.CStdout.Close()
			***REMOVED***
			if cfg.CStderr != nil ***REMOVED***
				cfg.CStderr.Close()
			***REMOVED***
		***REMOVED***
		logrus.Debug("attach: stdin: end")
		wg.Done()
	***REMOVED***()

	attachStream := func(name string, stream io.Writer, streamPipe io.ReadCloser) ***REMOVED***
		if stream == nil ***REMOVED***
			return
		***REMOVED***

		logrus.Debugf("attach: %s: begin", name)
		_, err := pools.Copy(stream, streamPipe)
		if err == io.ErrClosedPipe ***REMOVED***
			err = nil
		***REMOVED***
		if err != nil ***REMOVED***
			logrus.Errorf("attach: %s: %v", name, err)
			errors <- err
		***REMOVED***
		// Make sure stdin gets closed
		if cfg.Stdin != nil ***REMOVED***
			cfg.Stdin.Close()
		***REMOVED***
		streamPipe.Close()
		logrus.Debugf("attach: %s: end", name)
		wg.Done()
	***REMOVED***

	go attachStream("stdout", cfg.Stdout, cfg.CStdout)
	go attachStream("stderr", cfg.Stderr, cfg.CStderr)

	errs := make(chan error, 1)

	go func() ***REMOVED***
		defer close(errs)
		errs <- func() error ***REMOVED***
			done := make(chan struct***REMOVED******REMOVED***)
			go func() ***REMOVED***
				wg.Wait()
				close(done)
			***REMOVED***()
			select ***REMOVED***
			case <-done:
			case <-ctx.Done():
				// close all pipes
				if cfg.CStdin != nil ***REMOVED***
					cfg.CStdin.Close()
				***REMOVED***
				if cfg.CStdout != nil ***REMOVED***
					cfg.CStdout.Close()
				***REMOVED***
				if cfg.CStderr != nil ***REMOVED***
					cfg.CStderr.Close()
				***REMOVED***
				<-done
			***REMOVED***
			close(errors)
			for err := range errors ***REMOVED***
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
			return nil
		***REMOVED***()
	***REMOVED***()

	return errs
***REMOVED***

func copyEscapable(dst io.Writer, src io.ReadCloser, keys []byte) (written int64, err error) ***REMOVED***
	if len(keys) == 0 ***REMOVED***
		keys = defaultEscapeSequence
	***REMOVED***
	pr := term.NewEscapeProxy(src, keys)
	defer src.Close()

	return pools.Copy(dst, pr)
***REMOVED***
