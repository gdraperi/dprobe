// +build linux freebsd

package logger

import (
	"context"
	"io"

	"github.com/containerd/fifo"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

func openPluginStream(a *pluginAdapter) (io.WriteCloser, error) ***REMOVED***
	f, err := fifo.OpenFifo(context.Background(), a.fifoPath, unix.O_WRONLY|unix.O_CREAT|unix.O_NONBLOCK, 0700)
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "error creating i/o pipe for log plugin: %s", a.Name())
	***REMOVED***
	return f, nil
***REMOVED***
