package httputils

import (
	"fmt"
	"io"
	"net/url"
	"sort"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/stdcopy"
)

// WriteLogStream writes an encoded byte stream of log messages from the
// messages channel, multiplexing them with a stdcopy.Writer if mux is true
func WriteLogStream(_ context.Context, w io.Writer, msgs <-chan *backend.LogMessage, config *types.ContainerLogsOptions, mux bool) ***REMOVED***
	wf := ioutils.NewWriteFlusher(w)
	defer wf.Close()

	wf.Flush()

	outStream := io.Writer(wf)
	errStream := outStream
	sysErrStream := errStream
	if mux ***REMOVED***
		sysErrStream = stdcopy.NewStdWriter(outStream, stdcopy.Systemerr)
		errStream = stdcopy.NewStdWriter(outStream, stdcopy.Stderr)
		outStream = stdcopy.NewStdWriter(outStream, stdcopy.Stdout)
	***REMOVED***

	for ***REMOVED***
		msg, ok := <-msgs
		if !ok ***REMOVED***
			return
		***REMOVED***
		// check if the message contains an error. if so, write that error
		// and exit
		if msg.Err != nil ***REMOVED***
			fmt.Fprintf(sysErrStream, "Error grabbing logs: %v\n", msg.Err)
			continue
		***REMOVED***
		logLine := msg.Line
		if config.Details ***REMOVED***
			logLine = append(attrsByteSlice(msg.Attrs), ' ')
			logLine = append(logLine, msg.Line...)
		***REMOVED***
		if config.Timestamps ***REMOVED***
			logLine = append([]byte(msg.Timestamp.Format(jsonmessage.RFC3339NanoFixed)+" "), logLine...)
		***REMOVED***
		if msg.Source == "stdout" && config.ShowStdout ***REMOVED***
			outStream.Write(logLine)
		***REMOVED***
		if msg.Source == "stderr" && config.ShowStderr ***REMOVED***
			errStream.Write(logLine)
		***REMOVED***
	***REMOVED***
***REMOVED***

type byKey []backend.LogAttr

func (b byKey) Len() int           ***REMOVED*** return len(b) ***REMOVED***
func (b byKey) Less(i, j int) bool ***REMOVED*** return b[i].Key < b[j].Key ***REMOVED***
func (b byKey) Swap(i, j int)      ***REMOVED*** b[i], b[j] = b[j], b[i] ***REMOVED***

func attrsByteSlice(a []backend.LogAttr) []byte ***REMOVED***
	// Note this sorts "a" in-place. That is fine here - nothing else is
	// going to use Attrs or care about the order.
	sort.Sort(byKey(a))

	var ret []byte
	for i, pair := range a ***REMOVED***
		k, v := url.QueryEscape(pair.Key), url.QueryEscape(pair.Value)
		ret = append(ret, []byte(k)...)
		ret = append(ret, '=')
		ret = append(ret, []byte(v)...)
		if i != len(a)-1 ***REMOVED***
			ret = append(ret, ',')
		***REMOVED***
	***REMOVED***
	return ret
***REMOVED***
