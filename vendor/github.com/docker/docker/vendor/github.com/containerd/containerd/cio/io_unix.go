// +build !windows

package cio

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/containerd/fifo"
	"github.com/pkg/errors"
)

// NewFIFOSetInDir returns a new FIFOSet with paths in a temporary directory under root
func NewFIFOSetInDir(root, id string, terminal bool) (*FIFOSet, error) ***REMOVED***
	if root != "" ***REMOVED***
		if err := os.MkdirAll(root, 0700); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	dir, err := ioutil.TempDir(root, "")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	closer := func() error ***REMOVED***
		return os.RemoveAll(dir)
	***REMOVED***
	return NewFIFOSet(Config***REMOVED***
		Stdin:    filepath.Join(dir, id+"-stdin"),
		Stdout:   filepath.Join(dir, id+"-stdout"),
		Stderr:   filepath.Join(dir, id+"-stderr"),
		Terminal: terminal,
	***REMOVED***, closer), nil
***REMOVED***

func copyIO(fifos *FIFOSet, ioset *Streams) (*cio, error) ***REMOVED***
	var ctx, cancel = context.WithCancel(context.Background())
	pipes, err := openFifos(ctx, fifos)
	if err != nil ***REMOVED***
		cancel()
		return nil, err
	***REMOVED***

	if fifos.Stdin != "" ***REMOVED***
		go func() ***REMOVED***
			io.Copy(pipes.Stdin, ioset.Stdin)
			pipes.Stdin.Close()
		***REMOVED***()
	***REMOVED***

	var wg = &sync.WaitGroup***REMOVED******REMOVED***
	wg.Add(1)
	go func() ***REMOVED***
		io.Copy(ioset.Stdout, pipes.Stdout)
		pipes.Stdout.Close()
		wg.Done()
	***REMOVED***()

	if !fifos.Terminal ***REMOVED***
		wg.Add(1)
		go func() ***REMOVED***
			io.Copy(ioset.Stderr, pipes.Stderr)
			pipes.Stderr.Close()
			wg.Done()
		***REMOVED***()
	***REMOVED***
	return &cio***REMOVED***
		config:  fifos.Config,
		wg:      wg,
		closers: append(pipes.closers(), fifos),
		cancel:  cancel,
	***REMOVED***, nil
***REMOVED***

func openFifos(ctx context.Context, fifos *FIFOSet) (pipes, error) ***REMOVED***
	var err error
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			fifos.Close()
		***REMOVED***
	***REMOVED***()

	var f pipes
	if fifos.Stdin != "" ***REMOVED***
		if f.Stdin, err = fifo.OpenFifo(ctx, fifos.Stdin, syscall.O_WRONLY|syscall.O_CREAT|syscall.O_NONBLOCK, 0700); err != nil ***REMOVED***
			return f, errors.Wrapf(err, "failed to open stdin fifo")
		***REMOVED***
	***REMOVED***
	if fifos.Stdout != "" ***REMOVED***
		if f.Stdout, err = fifo.OpenFifo(ctx, fifos.Stdout, syscall.O_RDONLY|syscall.O_CREAT|syscall.O_NONBLOCK, 0700); err != nil ***REMOVED***
			f.Stdin.Close()
			return f, errors.Wrapf(err, "failed to open stdout fifo")
		***REMOVED***
	***REMOVED***
	if fifos.Stderr != "" ***REMOVED***
		if f.Stderr, err = fifo.OpenFifo(ctx, fifos.Stderr, syscall.O_RDONLY|syscall.O_CREAT|syscall.O_NONBLOCK, 0700); err != nil ***REMOVED***
			f.Stdin.Close()
			f.Stdout.Close()
			return f, errors.Wrapf(err, "failed to open stderr fifo")
		***REMOVED***
	***REMOVED***
	return f, nil
***REMOVED***

// NewDirectIO returns an IO implementation that exposes the IO streams as io.ReadCloser
// and io.WriteCloser.
func NewDirectIO(ctx context.Context, fifos *FIFOSet) (*DirectIO, error) ***REMOVED***
	ctx, cancel := context.WithCancel(ctx)
	pipes, err := openFifos(ctx, fifos)
	return &DirectIO***REMOVED***
		pipes: pipes,
		cio: cio***REMOVED***
			config:  fifos.Config,
			closers: append(pipes.closers(), fifos),
			cancel:  cancel,
		***REMOVED***,
	***REMOVED***, err
***REMOVED***

func (p *pipes) closers() []io.Closer ***REMOVED***
	return []io.Closer***REMOVED***p.Stdin, p.Stdout, p.Stderr***REMOVED***
***REMOVED***
