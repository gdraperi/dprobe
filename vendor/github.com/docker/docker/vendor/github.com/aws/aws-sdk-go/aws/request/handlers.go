package request

import (
	"fmt"
	"strings"
)

// A Handlers provides a collection of request handlers for various
// stages of handling requests.
type Handlers struct ***REMOVED***
	Validate         HandlerList
	Build            HandlerList
	Sign             HandlerList
	Send             HandlerList
	ValidateResponse HandlerList
	Unmarshal        HandlerList
	UnmarshalMeta    HandlerList
	UnmarshalError   HandlerList
	Retry            HandlerList
	AfterRetry       HandlerList
	Complete         HandlerList
***REMOVED***

// Copy returns of this handler's lists.
func (h *Handlers) Copy() Handlers ***REMOVED***
	return Handlers***REMOVED***
		Validate:         h.Validate.copy(),
		Build:            h.Build.copy(),
		Sign:             h.Sign.copy(),
		Send:             h.Send.copy(),
		ValidateResponse: h.ValidateResponse.copy(),
		Unmarshal:        h.Unmarshal.copy(),
		UnmarshalError:   h.UnmarshalError.copy(),
		UnmarshalMeta:    h.UnmarshalMeta.copy(),
		Retry:            h.Retry.copy(),
		AfterRetry:       h.AfterRetry.copy(),
		Complete:         h.Complete.copy(),
	***REMOVED***
***REMOVED***

// Clear removes callback functions for all handlers
func (h *Handlers) Clear() ***REMOVED***
	h.Validate.Clear()
	h.Build.Clear()
	h.Send.Clear()
	h.Sign.Clear()
	h.Unmarshal.Clear()
	h.UnmarshalMeta.Clear()
	h.UnmarshalError.Clear()
	h.ValidateResponse.Clear()
	h.Retry.Clear()
	h.AfterRetry.Clear()
	h.Complete.Clear()
***REMOVED***

// A HandlerListRunItem represents an entry in the HandlerList which
// is being run.
type HandlerListRunItem struct ***REMOVED***
	Index   int
	Handler NamedHandler
	Request *Request
***REMOVED***

// A HandlerList manages zero or more handlers in a list.
type HandlerList struct ***REMOVED***
	list []NamedHandler

	// Called after each request handler in the list is called. If set
	// and the func returns true the HandlerList will continue to iterate
	// over the request handlers. If false is returned the HandlerList
	// will stop iterating.
	//
	// Should be used if extra logic to be performed between each handler
	// in the list. This can be used to terminate a list's iteration
	// based on a condition such as error like, HandlerListStopOnError.
	// Or for logging like HandlerListLogItem.
	AfterEachFn func(item HandlerListRunItem) bool
***REMOVED***

// A NamedHandler is a struct that contains a name and function callback.
type NamedHandler struct ***REMOVED***
	Name string
	Fn   func(*Request)
***REMOVED***

// copy creates a copy of the handler list.
func (l *HandlerList) copy() HandlerList ***REMOVED***
	n := HandlerList***REMOVED***
		AfterEachFn: l.AfterEachFn,
	***REMOVED***
	if len(l.list) == 0 ***REMOVED***
		return n
	***REMOVED***

	n.list = append(make([]NamedHandler, 0, len(l.list)), l.list...)
	return n
***REMOVED***

// Clear clears the handler list.
func (l *HandlerList) Clear() ***REMOVED***
	l.list = l.list[0:0]
***REMOVED***

// Len returns the number of handlers in the list.
func (l *HandlerList) Len() int ***REMOVED***
	return len(l.list)
***REMOVED***

// PushBack pushes handler f to the back of the handler list.
func (l *HandlerList) PushBack(f func(*Request)) ***REMOVED***
	l.PushBackNamed(NamedHandler***REMOVED***"__anonymous", f***REMOVED***)
***REMOVED***

// PushBackNamed pushes named handler f to the back of the handler list.
func (l *HandlerList) PushBackNamed(n NamedHandler) ***REMOVED***
	if cap(l.list) == 0 ***REMOVED***
		l.list = make([]NamedHandler, 0, 5)
	***REMOVED***
	l.list = append(l.list, n)
***REMOVED***

// PushFront pushes handler f to the front of the handler list.
func (l *HandlerList) PushFront(f func(*Request)) ***REMOVED***
	l.PushFrontNamed(NamedHandler***REMOVED***"__anonymous", f***REMOVED***)
***REMOVED***

// PushFrontNamed pushes named handler f to the front of the handler list.
func (l *HandlerList) PushFrontNamed(n NamedHandler) ***REMOVED***
	if cap(l.list) == len(l.list) ***REMOVED***
		// Allocating new list required
		l.list = append([]NamedHandler***REMOVED***n***REMOVED***, l.list...)
	***REMOVED*** else ***REMOVED***
		// Enough room to prepend into list.
		l.list = append(l.list, NamedHandler***REMOVED******REMOVED***)
		copy(l.list[1:], l.list)
		l.list[0] = n
	***REMOVED***
***REMOVED***

// Remove removes a NamedHandler n
func (l *HandlerList) Remove(n NamedHandler) ***REMOVED***
	l.RemoveByName(n.Name)
***REMOVED***

// RemoveByName removes a NamedHandler by name.
func (l *HandlerList) RemoveByName(name string) ***REMOVED***
	for i := 0; i < len(l.list); i++ ***REMOVED***
		m := l.list[i]
		if m.Name == name ***REMOVED***
			// Shift array preventing creating new arrays
			copy(l.list[i:], l.list[i+1:])
			l.list[len(l.list)-1] = NamedHandler***REMOVED******REMOVED***
			l.list = l.list[:len(l.list)-1]

			// decrement list so next check to length is correct
			i--
		***REMOVED***
	***REMOVED***
***REMOVED***

// SwapNamed will swap out any existing handlers with the same name as the
// passed in NamedHandler returning true if handlers were swapped. False is
// returned otherwise.
func (l *HandlerList) SwapNamed(n NamedHandler) (swapped bool) ***REMOVED***
	for i := 0; i < len(l.list); i++ ***REMOVED***
		if l.list[i].Name == n.Name ***REMOVED***
			l.list[i].Fn = n.Fn
			swapped = true
		***REMOVED***
	***REMOVED***

	return swapped
***REMOVED***

// SetBackNamed will replace the named handler if it exists in the handler list.
// If the handler does not exist the handler will be added to the end of the list.
func (l *HandlerList) SetBackNamed(n NamedHandler) ***REMOVED***
	if !l.SwapNamed(n) ***REMOVED***
		l.PushBackNamed(n)
	***REMOVED***
***REMOVED***

// SetFrontNamed will replace the named handler if it exists in the handler list.
// If the handler does not exist the handler will be added to the beginning of
// the list.
func (l *HandlerList) SetFrontNamed(n NamedHandler) ***REMOVED***
	if !l.SwapNamed(n) ***REMOVED***
		l.PushFrontNamed(n)
	***REMOVED***
***REMOVED***

// Run executes all handlers in the list with a given request object.
func (l *HandlerList) Run(r *Request) ***REMOVED***
	for i, h := range l.list ***REMOVED***
		h.Fn(r)
		item := HandlerListRunItem***REMOVED***
			Index: i, Handler: h, Request: r,
		***REMOVED***
		if l.AfterEachFn != nil && !l.AfterEachFn(item) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// HandlerListLogItem logs the request handler and the state of the
// request's Error value. Always returns true to continue iterating
// request handlers in a HandlerList.
func HandlerListLogItem(item HandlerListRunItem) bool ***REMOVED***
	if item.Request.Config.Logger == nil ***REMOVED***
		return true
	***REMOVED***
	item.Request.Config.Logger.Log("DEBUG: RequestHandler",
		item.Index, item.Handler.Name, item.Request.Error)

	return true
***REMOVED***

// HandlerListStopOnError returns false to stop the HandlerList iterating
// over request handlers if Request.Error is not nil. True otherwise
// to continue iterating.
func HandlerListStopOnError(item HandlerListRunItem) bool ***REMOVED***
	return item.Request.Error == nil
***REMOVED***

// WithAppendUserAgent will add a string to the user agent prefixed with a
// single white space.
func WithAppendUserAgent(s string) Option ***REMOVED***
	return func(r *Request) ***REMOVED***
		r.Handlers.Build.PushBack(func(r2 *Request) ***REMOVED***
			AddToUserAgent(r, s)
		***REMOVED***)
	***REMOVED***
***REMOVED***

// MakeAddToUserAgentHandler will add the name/version pair to the User-Agent request
// header. If the extra parameters are provided they will be added as metadata to the
// name/version pair resulting in the following format.
// "name/version (extra0; extra1; ...)"
// The user agent part will be concatenated with this current request's user agent string.
func MakeAddToUserAgentHandler(name, version string, extra ...string) func(*Request) ***REMOVED***
	ua := fmt.Sprintf("%s/%s", name, version)
	if len(extra) > 0 ***REMOVED***
		ua += fmt.Sprintf(" (%s)", strings.Join(extra, "; "))
	***REMOVED***
	return func(r *Request) ***REMOVED***
		AddToUserAgent(r, ua)
	***REMOVED***
***REMOVED***

// MakeAddToUserAgentFreeFormHandler adds the input to the User-Agent request header.
// The input string will be concatenated with the current request's user agent string.
func MakeAddToUserAgentFreeFormHandler(s string) func(*Request) ***REMOVED***
	return func(r *Request) ***REMOVED***
		AddToUserAgent(r, s)
	***REMOVED***
***REMOVED***
