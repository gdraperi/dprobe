package libcontainerd

import "fmt"

// WithOOMScore defines the oom_score_adj to set for the containerd process.
func WithOOMScore(score int) RemoteOption ***REMOVED***
	return oomScore(score)
***REMOVED***

type oomScore int

func (o oomScore) Apply(r Remote) error ***REMOVED***
	if remote, ok := r.(*remote); ok ***REMOVED***
		remote.OOMScore = int(o)
		return nil
	***REMOVED***
	return fmt.Errorf("WithOOMScore option not supported for this remote")
***REMOVED***

// WithSubreaper sets whether containerd should register itself as a
// subreaper
func WithSubreaper(reap bool) RemoteOption ***REMOVED***
	return subreaper(reap)
***REMOVED***

type subreaper bool

func (s subreaper) Apply(r Remote) error ***REMOVED***
	if remote, ok := r.(*remote); ok ***REMOVED***
		remote.NoSubreaper = !bool(s)
		return nil
	***REMOVED***
	return fmt.Errorf("WithSubreaper option not supported for this remote")
***REMOVED***
