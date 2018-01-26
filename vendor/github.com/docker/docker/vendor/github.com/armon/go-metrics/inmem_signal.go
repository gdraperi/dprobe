package metrics

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// InmemSignal is used to listen for a given signal, and when received,
// to dump the current metrics from the InmemSink to an io.Writer
type InmemSignal struct ***REMOVED***
	signal syscall.Signal
	inm    *InmemSink
	w      io.Writer
	sigCh  chan os.Signal

	stop     bool
	stopCh   chan struct***REMOVED******REMOVED***
	stopLock sync.Mutex
***REMOVED***

// NewInmemSignal creates a new InmemSignal which listens for a given signal,
// and dumps the current metrics out to a writer
func NewInmemSignal(inmem *InmemSink, sig syscall.Signal, w io.Writer) *InmemSignal ***REMOVED***
	i := &InmemSignal***REMOVED***
		signal: sig,
		inm:    inmem,
		w:      w,
		sigCh:  make(chan os.Signal, 1),
		stopCh: make(chan struct***REMOVED******REMOVED***),
	***REMOVED***
	signal.Notify(i.sigCh, sig)
	go i.run()
	return i
***REMOVED***

// DefaultInmemSignal returns a new InmemSignal that responds to SIGUSR1
// and writes output to stderr. Windows uses SIGBREAK
func DefaultInmemSignal(inmem *InmemSink) *InmemSignal ***REMOVED***
	return NewInmemSignal(inmem, DefaultSignal, os.Stderr)
***REMOVED***

// Stop is used to stop the InmemSignal from listening
func (i *InmemSignal) Stop() ***REMOVED***
	i.stopLock.Lock()
	defer i.stopLock.Unlock()

	if i.stop ***REMOVED***
		return
	***REMOVED***
	i.stop = true
	close(i.stopCh)
	signal.Stop(i.sigCh)
***REMOVED***

// run is a long running routine that handles signals
func (i *InmemSignal) run() ***REMOVED***
	for ***REMOVED***
		select ***REMOVED***
		case <-i.sigCh:
			i.dumpStats()
		case <-i.stopCh:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// dumpStats is used to dump the data to output writer
func (i *InmemSignal) dumpStats() ***REMOVED***
	buf := bytes.NewBuffer(nil)

	data := i.inm.Data()
	// Skip the last period which is still being aggregated
	for i := 0; i < len(data)-1; i++ ***REMOVED***
		intv := data[i]
		intv.RLock()
		for name, val := range intv.Gauges ***REMOVED***
			fmt.Fprintf(buf, "[%v][G] '%s': %0.3f\n", intv.Interval, name, val)
		***REMOVED***
		for name, vals := range intv.Points ***REMOVED***
			for _, val := range vals ***REMOVED***
				fmt.Fprintf(buf, "[%v][P] '%s': %0.3f\n", intv.Interval, name, val)
			***REMOVED***
		***REMOVED***
		for name, agg := range intv.Counters ***REMOVED***
			fmt.Fprintf(buf, "[%v][C] '%s': %s\n", intv.Interval, name, agg)
		***REMOVED***
		for name, agg := range intv.Samples ***REMOVED***
			fmt.Fprintf(buf, "[%v][S] '%s': %s\n", intv.Interval, name, agg)
		***REMOVED***
		intv.RUnlock()
	***REMOVED***

	// Write out the bytes
	i.w.Write(buf.Bytes())
***REMOVED***
