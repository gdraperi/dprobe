package check

import (
	"fmt"
	"io"
	"sync"
)

// -----------------------------------------------------------------------
// Output writer manages atomic output writing according to settings.

type outputWriter struct ***REMOVED***
	m                    sync.Mutex
	writer               io.Writer
	wroteCallProblemLast bool
	Stream               bool
	Verbose              bool
***REMOVED***

func newOutputWriter(writer io.Writer, stream, verbose bool) *outputWriter ***REMOVED***
	return &outputWriter***REMOVED***writer: writer, Stream: stream, Verbose: verbose***REMOVED***
***REMOVED***

func (ow *outputWriter) Write(content []byte) (n int, err error) ***REMOVED***
	ow.m.Lock()
	n, err = ow.writer.Write(content)
	ow.m.Unlock()
	return
***REMOVED***

func (ow *outputWriter) WriteCallStarted(label string, c *C) ***REMOVED***
	if ow.Stream ***REMOVED***
		header := renderCallHeader(label, c, "", "\n")
		ow.m.Lock()
		ow.writer.Write([]byte(header))
		ow.m.Unlock()
	***REMOVED***
***REMOVED***

func (ow *outputWriter) WriteCallProblem(label string, c *C) ***REMOVED***
	var prefix string
	if !ow.Stream ***REMOVED***
		prefix = "\n-----------------------------------" +
			"-----------------------------------\n"
	***REMOVED***
	header := renderCallHeader(label, c, prefix, "\n\n")
	ow.m.Lock()
	ow.wroteCallProblemLast = true
	ow.writer.Write([]byte(header))
	if !ow.Stream ***REMOVED***
		c.logb.WriteTo(ow.writer)
	***REMOVED***
	ow.m.Unlock()
***REMOVED***

func (ow *outputWriter) WriteCallSuccess(label string, c *C) ***REMOVED***
	if ow.Stream || (ow.Verbose && c.kind == testKd) ***REMOVED***
		// TODO Use a buffer here.
		var suffix string
		if c.reason != "" ***REMOVED***
			suffix = " (" + c.reason + ")"
		***REMOVED***
		if c.status() == succeededSt ***REMOVED***
			suffix += "\t" + c.timerString()
		***REMOVED***
		suffix += "\n"
		if ow.Stream ***REMOVED***
			suffix += "\n"
		***REMOVED***
		header := renderCallHeader(label, c, "", suffix)
		ow.m.Lock()
		// Resist temptation of using line as prefix above due to race.
		if !ow.Stream && ow.wroteCallProblemLast ***REMOVED***
			header = "\n-----------------------------------" +
				"-----------------------------------\n" +
				header
		***REMOVED***
		ow.wroteCallProblemLast = false
		ow.writer.Write([]byte(header))
		ow.m.Unlock()
	***REMOVED***
***REMOVED***

func renderCallHeader(label string, c *C, prefix, suffix string) string ***REMOVED***
	pc := c.method.PC()
	return fmt.Sprintf("%s%s: %s: %s%s", prefix, label, niceFuncPath(pc),
		niceFuncName(pc), suffix)
***REMOVED***
