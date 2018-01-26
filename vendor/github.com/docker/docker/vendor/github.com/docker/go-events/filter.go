package events

// Matcher matches events.
type Matcher interface ***REMOVED***
	Match(event Event) bool
***REMOVED***

// MatcherFunc implements matcher with just a function.
type MatcherFunc func(event Event) bool

// Match calls the wrapped function.
func (fn MatcherFunc) Match(event Event) bool ***REMOVED***
	return fn(event)
***REMOVED***

// Filter provides an event sink that sends only events that are accepted by a
// Matcher. No methods on filter are goroutine safe.
type Filter struct ***REMOVED***
	dst     Sink
	matcher Matcher
	closed  bool
***REMOVED***

// NewFilter returns a new filter that will send to events to dst that return
// true for Matcher.
func NewFilter(dst Sink, matcher Matcher) Sink ***REMOVED***
	return &Filter***REMOVED***dst: dst, matcher: matcher***REMOVED***
***REMOVED***

// Write an event to the filter.
func (f *Filter) Write(event Event) error ***REMOVED***
	if f.closed ***REMOVED***
		return ErrSinkClosed
	***REMOVED***

	if f.matcher.Match(event) ***REMOVED***
		return f.dst.Write(event)
	***REMOVED***

	return nil
***REMOVED***

// Close the filter and allow no more events to pass through.
func (f *Filter) Close() error ***REMOVED***
	// TODO(stevvooe): Not all sinks should have Close.
	if f.closed ***REMOVED***
		return nil
	***REMOVED***

	f.closed = true
	return f.dst.Close()
***REMOVED***
