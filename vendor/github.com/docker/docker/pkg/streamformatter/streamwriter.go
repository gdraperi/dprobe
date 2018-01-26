package streamformatter

import (
	"encoding/json"
	"io"

	"github.com/docker/docker/pkg/jsonmessage"
)

type streamWriter struct ***REMOVED***
	io.Writer
	lineFormat func([]byte) string
***REMOVED***

func (sw *streamWriter) Write(buf []byte) (int, error) ***REMOVED***
	formattedBuf := sw.format(buf)
	n, err := sw.Writer.Write(formattedBuf)
	if n != len(formattedBuf) ***REMOVED***
		return n, io.ErrShortWrite
	***REMOVED***
	return len(buf), err
***REMOVED***

func (sw *streamWriter) format(buf []byte) []byte ***REMOVED***
	msg := &jsonmessage.JSONMessage***REMOVED***Stream: sw.lineFormat(buf)***REMOVED***
	b, err := json.Marshal(msg)
	if err != nil ***REMOVED***
		return FormatError(err)
	***REMOVED***
	return appendNewline(b)
***REMOVED***

// NewStdoutWriter returns a writer which formats the output as json message
// representing stdout lines
func NewStdoutWriter(out io.Writer) io.Writer ***REMOVED***
	return &streamWriter***REMOVED***Writer: out, lineFormat: func(buf []byte) string ***REMOVED***
		return string(buf)
	***REMOVED******REMOVED***
***REMOVED***

// NewStderrWriter returns a writer which formats the output as json message
// representing stderr lines
func NewStderrWriter(out io.Writer) io.Writer ***REMOVED***
	return &streamWriter***REMOVED***Writer: out, lineFormat: func(buf []byte) string ***REMOVED***
		return "\033[91m" + string(buf) + "\033[0m"
	***REMOVED******REMOVED***
***REMOVED***
