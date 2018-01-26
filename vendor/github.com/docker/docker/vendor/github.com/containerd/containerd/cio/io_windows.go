package cio

import (
	"fmt"
	"io"
	"net"
	"sync"

	winio "github.com/Microsoft/go-winio"
	"github.com/containerd/containerd/log"
	"github.com/pkg/errors"
)

const pipeRoot = `\\.\pipe`

// NewFIFOSetInDir returns a new set of fifos for the task
func NewFIFOSetInDir(_, id string, terminal bool) (*FIFOSet, error) ***REMOVED***
	return NewFIFOSet(Config***REMOVED***
		Terminal: terminal,
		Stdin:    fmt.Sprintf(`%s\ctr-%s-stdin`, pipeRoot, id),
		Stdout:   fmt.Sprintf(`%s\ctr-%s-stdout`, pipeRoot, id),
		Stderr:   fmt.Sprintf(`%s\ctr-%s-stderr`, pipeRoot, id),
	***REMOVED***, nil), nil
***REMOVED***

func copyIO(fifos *FIFOSet, ioset *Streams) (*cio, error) ***REMOVED***
	var (
		wg  sync.WaitGroup
		set []io.Closer
	)

	if fifos.Stdin != "" ***REMOVED***
		l, err := winio.ListenPipe(fifos.Stdin, nil)
		if err != nil ***REMOVED***
			return nil, errors.Wrapf(err, "failed to create stdin pipe %s", fifos.Stdin)
		***REMOVED***
		defer func(l net.Listener) ***REMOVED***
			if err != nil ***REMOVED***
				l.Close()
			***REMOVED***
		***REMOVED***(l)
		set = append(set, l)

		go func() ***REMOVED***
			c, err := l.Accept()
			if err != nil ***REMOVED***
				log.L.WithError(err).Errorf("failed to accept stdin connection on %s", fifos.Stdin)
				return
			***REMOVED***
			io.Copy(c, ioset.Stdin)
			c.Close()
			l.Close()
		***REMOVED***()
	***REMOVED***

	if fifos.Stdout != "" ***REMOVED***
		l, err := winio.ListenPipe(fifos.Stdout, nil)
		if err != nil ***REMOVED***
			return nil, errors.Wrapf(err, "failed to create stdin pipe %s", fifos.Stdout)
		***REMOVED***
		defer func(l net.Listener) ***REMOVED***
			if err != nil ***REMOVED***
				l.Close()
			***REMOVED***
		***REMOVED***(l)
		set = append(set, l)

		wg.Add(1)
		go func() ***REMOVED***
			defer wg.Done()
			c, err := l.Accept()
			if err != nil ***REMOVED***
				log.L.WithError(err).Errorf("failed to accept stdout connection on %s", fifos.Stdout)
				return
			***REMOVED***
			io.Copy(ioset.Stdout, c)
			c.Close()
			l.Close()
		***REMOVED***()
	***REMOVED***

	if !fifos.Terminal && fifos.Stderr != "" ***REMOVED***
		l, err := winio.ListenPipe(fifos.Stderr, nil)
		if err != nil ***REMOVED***
			return nil, errors.Wrapf(err, "failed to create stderr pipe %s", fifos.Stderr)
		***REMOVED***
		defer func(l net.Listener) ***REMOVED***
			if err != nil ***REMOVED***
				l.Close()
			***REMOVED***
		***REMOVED***(l)
		set = append(set, l)

		wg.Add(1)
		go func() ***REMOVED***
			defer wg.Done()
			c, err := l.Accept()
			if err != nil ***REMOVED***
				log.L.WithError(err).Errorf("failed to accept stderr connection on %s", fifos.Stderr)
				return
			***REMOVED***
			io.Copy(ioset.Stderr, c)
			c.Close()
			l.Close()
		***REMOVED***()
	***REMOVED***

	return &cio***REMOVED***config: fifos.Config, closers: set***REMOVED***, nil
***REMOVED***

// NewDirectIO returns an IO implementation that exposes the IO streams as io.ReadCloser
// and io.WriteCloser.
func NewDirectIO(stdin io.WriteCloser, stdout, stderr io.ReadCloser, terminal bool) *DirectIO ***REMOVED***
	return &DirectIO***REMOVED***
		pipes: pipes***REMOVED***
			Stdin:  stdin,
			Stdout: stdout,
			Stderr: stderr,
		***REMOVED***,
		cio: cio***REMOVED***
			config: Config***REMOVED***Terminal: terminal***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***
