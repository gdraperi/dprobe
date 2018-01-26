// +build !linux,!freebsd

package logger

import (
	"errors"
	"io"
)

func openPluginStream(a *pluginAdapter) (io.WriteCloser, error) ***REMOVED***
	return nil, errors.New("log plugin not supported")
***REMOVED***
