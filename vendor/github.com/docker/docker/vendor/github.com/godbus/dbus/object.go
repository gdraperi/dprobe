package dbus

import (
	"errors"
	"strings"
)

// BusObject is the interface of a remote object on which methods can be
// invoked.
type BusObject interface ***REMOVED***
	Call(method string, flags Flags, args ...interface***REMOVED******REMOVED***) *Call
	Go(method string, flags Flags, ch chan *Call, args ...interface***REMOVED******REMOVED***) *Call
	GetProperty(p string) (Variant, error)
	Destination() string
	Path() ObjectPath
***REMOVED***

// Object represents a remote object on which methods can be invoked.
type Object struct ***REMOVED***
	conn *Conn
	dest string
	path ObjectPath
***REMOVED***

// Call calls a method with (*Object).Go and waits for its reply.
func (o *Object) Call(method string, flags Flags, args ...interface***REMOVED******REMOVED***) *Call ***REMOVED***
	return <-o.Go(method, flags, make(chan *Call, 1), args...).Done
***REMOVED***

// AddMatchSignal subscribes BusObject to signals from specified interface and
// method (member).
func (o *Object) AddMatchSignal(iface, member string) *Call ***REMOVED***
	return o.Call(
		"org.freedesktop.DBus.AddMatch",
		0,
		"type='signal',interface='"+iface+"',member='"+member+"'",
	)
***REMOVED***

// Go calls a method with the given arguments asynchronously. It returns a
// Call structure representing this method call. The passed channel will
// return the same value once the call is done. If ch is nil, a new channel
// will be allocated. Otherwise, ch has to be buffered or Go will panic.
//
// If the flags include FlagNoReplyExpected, ch is ignored and a Call structure
// is returned of which only the Err member is valid.
//
// If the method parameter contains a dot ('.'), the part before the last dot
// specifies the interface on which the method is called.
func (o *Object) Go(method string, flags Flags, ch chan *Call, args ...interface***REMOVED******REMOVED***) *Call ***REMOVED***
	iface := ""
	i := strings.LastIndex(method, ".")
	if i != -1 ***REMOVED***
		iface = method[:i]
	***REMOVED***
	method = method[i+1:]
	msg := new(Message)
	msg.Type = TypeMethodCall
	msg.serial = o.conn.getSerial()
	msg.Flags = flags & (FlagNoAutoStart | FlagNoReplyExpected)
	msg.Headers = make(map[HeaderField]Variant)
	msg.Headers[FieldPath] = MakeVariant(o.path)
	msg.Headers[FieldDestination] = MakeVariant(o.dest)
	msg.Headers[FieldMember] = MakeVariant(method)
	if iface != "" ***REMOVED***
		msg.Headers[FieldInterface] = MakeVariant(iface)
	***REMOVED***
	msg.Body = args
	if len(args) > 0 ***REMOVED***
		msg.Headers[FieldSignature] = MakeVariant(SignatureOf(args...))
	***REMOVED***
	if msg.Flags&FlagNoReplyExpected == 0 ***REMOVED***
		if ch == nil ***REMOVED***
			ch = make(chan *Call, 10)
		***REMOVED*** else if cap(ch) == 0 ***REMOVED***
			panic("dbus: unbuffered channel passed to (*Object).Go")
		***REMOVED***
		call := &Call***REMOVED***
			Destination: o.dest,
			Path:        o.path,
			Method:      method,
			Args:        args,
			Done:        ch,
		***REMOVED***
		o.conn.callsLck.Lock()
		o.conn.calls[msg.serial] = call
		o.conn.callsLck.Unlock()
		o.conn.outLck.RLock()
		if o.conn.closed ***REMOVED***
			call.Err = ErrClosed
			call.Done <- call
		***REMOVED*** else ***REMOVED***
			o.conn.out <- msg
		***REMOVED***
		o.conn.outLck.RUnlock()
		return call
	***REMOVED***
	o.conn.outLck.RLock()
	defer o.conn.outLck.RUnlock()
	if o.conn.closed ***REMOVED***
		return &Call***REMOVED***Err: ErrClosed***REMOVED***
	***REMOVED***
	o.conn.out <- msg
	return &Call***REMOVED***Err: nil***REMOVED***
***REMOVED***

// GetProperty calls org.freedesktop.DBus.Properties.GetProperty on the given
// object. The property name must be given in interface.member notation.
func (o *Object) GetProperty(p string) (Variant, error) ***REMOVED***
	idx := strings.LastIndex(p, ".")
	if idx == -1 || idx+1 == len(p) ***REMOVED***
		return Variant***REMOVED******REMOVED***, errors.New("dbus: invalid property " + p)
	***REMOVED***

	iface := p[:idx]
	prop := p[idx+1:]

	result := Variant***REMOVED******REMOVED***
	err := o.Call("org.freedesktop.DBus.Properties.Get", 0, iface, prop).Store(&result)

	if err != nil ***REMOVED***
		return Variant***REMOVED******REMOVED***, err
	***REMOVED***

	return result, nil
***REMOVED***

// Destination returns the destination that calls on o are sent to.
func (o *Object) Destination() string ***REMOVED***
	return o.dest
***REMOVED***

// Path returns the path that calls on o are sent to.
func (o *Object) Path() ObjectPath ***REMOVED***
	return o.path
***REMOVED***
