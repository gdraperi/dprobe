package runc

import (
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

type IO interface ***REMOVED***
	io.Closer
	Stdin() io.WriteCloser
	Stdout() io.ReadCloser
	Stderr() io.ReadCloser
	Set(*exec.Cmd)
***REMOVED***

type StartCloser interface ***REMOVED***
	CloseAfterStart() error
***REMOVED***

// NewPipeIO creates pipe pairs to be used with runc
func NewPipeIO(uid, gid int) (i IO, err error) ***REMOVED***
	var pipes []*pipe
	// cleanup in case of an error
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			for _, p := range pipes ***REMOVED***
				p.Close()
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	stdin, err := newPipe()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	pipes = append(pipes, stdin)
	if err = unix.Fchown(int(stdin.r.Fd()), uid, gid); err != nil ***REMOVED***
		return nil, errors.Wrap(err, "failed to chown stdin")
	***REMOVED***

	stdout, err := newPipe()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	pipes = append(pipes, stdout)
	if err = unix.Fchown(int(stdout.w.Fd()), uid, gid); err != nil ***REMOVED***
		return nil, errors.Wrap(err, "failed to chown stdout")
	***REMOVED***

	stderr, err := newPipe()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	pipes = append(pipes, stderr)
	if err = unix.Fchown(int(stderr.w.Fd()), uid, gid); err != nil ***REMOVED***
		return nil, errors.Wrap(err, "failed to chown stderr")
	***REMOVED***

	return &pipeIO***REMOVED***
		in:  stdin,
		out: stdout,
		err: stderr,
	***REMOVED***, nil
***REMOVED***

func newPipe() (*pipe, error) ***REMOVED***
	r, w, err := os.Pipe()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &pipe***REMOVED***
		r: r,
		w: w,
	***REMOVED***, nil
***REMOVED***

type pipe struct ***REMOVED***
	r *os.File
	w *os.File
***REMOVED***

func (p *pipe) Close() error ***REMOVED***
	err := p.r.Close()
	if werr := p.w.Close(); err == nil ***REMOVED***
		err = werr
	***REMOVED***
	return err
***REMOVED***

type pipeIO struct ***REMOVED***
	in  *pipe
	out *pipe
	err *pipe
***REMOVED***

func (i *pipeIO) Stdin() io.WriteCloser ***REMOVED***
	return i.in.w
***REMOVED***

func (i *pipeIO) Stdout() io.ReadCloser ***REMOVED***
	return i.out.r
***REMOVED***

func (i *pipeIO) Stderr() io.ReadCloser ***REMOVED***
	return i.err.r
***REMOVED***

func (i *pipeIO) Close() error ***REMOVED***
	var err error
	for _, v := range []*pipe***REMOVED***
		i.in,
		i.out,
		i.err,
	***REMOVED*** ***REMOVED***
		if cerr := v.Close(); err == nil ***REMOVED***
			err = cerr
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***

func (i *pipeIO) CloseAfterStart() error ***REMOVED***
	for _, f := range []*os.File***REMOVED***
		i.out.w,
		i.err.w,
	***REMOVED*** ***REMOVED***
		f.Close()
	***REMOVED***
	return nil
***REMOVED***

// Set sets the io to the exec.Cmd
func (i *pipeIO) Set(cmd *exec.Cmd) ***REMOVED***
	cmd.Stdin = i.in.r
	cmd.Stdout = i.out.w
	cmd.Stderr = i.err.w
***REMOVED***

func NewSTDIO() (IO, error) ***REMOVED***
	return &stdio***REMOVED******REMOVED***, nil
***REMOVED***

type stdio struct ***REMOVED***
***REMOVED***

func (s *stdio) Close() error ***REMOVED***
	return nil
***REMOVED***

func (s *stdio) Set(cmd *exec.Cmd) ***REMOVED***
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
***REMOVED***

func (s *stdio) Stdin() io.WriteCloser ***REMOVED***
	return os.Stdin
***REMOVED***

func (s *stdio) Stdout() io.ReadCloser ***REMOVED***
	return os.Stdout
***REMOVED***

func (s *stdio) Stderr() io.ReadCloser ***REMOVED***
	return os.Stderr
***REMOVED***

// NewNullIO returns IO setup for /dev/null use with runc
func NewNullIO() (IO, error) ***REMOVED***
	f, err := os.Open(os.DevNull)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &nullIO***REMOVED***
		devNull: f,
	***REMOVED***, nil
***REMOVED***

type nullIO struct ***REMOVED***
	devNull *os.File
***REMOVED***

func (n *nullIO) Close() error ***REMOVED***
	// this should be closed after start but if not
	// make sure we close the file but don't return the error
	n.devNull.Close()
	return nil
***REMOVED***

func (n *nullIO) Stdin() io.WriteCloser ***REMOVED***
	return nil
***REMOVED***

func (n *nullIO) Stdout() io.ReadCloser ***REMOVED***
	return nil
***REMOVED***

func (n *nullIO) Stderr() io.ReadCloser ***REMOVED***
	return nil
***REMOVED***

func (n *nullIO) Set(c *exec.Cmd) ***REMOVED***
	// don't set STDIN here
	c.Stdout = n.devNull
	c.Stderr = n.devNull
***REMOVED***

func (n *nullIO) CloseAfterStart() error ***REMOVED***
	return n.devNull.Close()
***REMOVED***
