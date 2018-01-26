package dbus

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	errmsgInvalidArg = Error***REMOVED***
		"org.freedesktop.DBus.Error.InvalidArgs",
		[]interface***REMOVED******REMOVED******REMOVED***"Invalid type / number of args"***REMOVED***,
	***REMOVED***
	errmsgNoObject = Error***REMOVED***
		"org.freedesktop.DBus.Error.NoSuchObject",
		[]interface***REMOVED******REMOVED******REMOVED***"No such object"***REMOVED***,
	***REMOVED***
	errmsgUnknownMethod = Error***REMOVED***
		"org.freedesktop.DBus.Error.UnknownMethod",
		[]interface***REMOVED******REMOVED******REMOVED***"Unknown / invalid method"***REMOVED***,
	***REMOVED***
)

// exportedObj represents an exported object. It stores a precomputed
// method table that represents the methods exported on the bus.
type exportedObj struct ***REMOVED***
	methods map[string]reflect.Value

	// Whether or not this export is for the entire subtree
	includeSubtree bool
***REMOVED***

func (obj exportedObj) Method(name string) (reflect.Value, bool) ***REMOVED***
	out, exists := obj.methods[name]
	return out, exists
***REMOVED***

// Sender is a type which can be used in exported methods to receive the message
// sender.
type Sender string

func computeMethodName(name string, mapping map[string]string) string ***REMOVED***
	newname, ok := mapping[name]
	if ok ***REMOVED***
		name = newname
	***REMOVED***
	return name
***REMOVED***

func getMethods(in interface***REMOVED******REMOVED***, mapping map[string]string) map[string]reflect.Value ***REMOVED***
	if in == nil ***REMOVED***
		return nil
	***REMOVED***
	methods := make(map[string]reflect.Value)
	val := reflect.ValueOf(in)
	typ := val.Type()
	for i := 0; i < typ.NumMethod(); i++ ***REMOVED***
		methtype := typ.Method(i)
		method := val.Method(i)
		t := method.Type()
		// only track valid methods must return *Error as last arg
		// and must be exported
		if t.NumOut() == 0 ||
			t.Out(t.NumOut()-1) != reflect.TypeOf(&errmsgInvalidArg) ||
			methtype.PkgPath != "" ***REMOVED***
			continue
		***REMOVED***
		// map names while building table
		methods[computeMethodName(methtype.Name, mapping)] = method
	***REMOVED***
	return methods
***REMOVED***

// searchHandlers will look through all registered handlers looking for one
// to handle the given path. If a verbatim one isn't found, it will check for
// a subtree registration for the path as well.
func (conn *Conn) searchHandlers(path ObjectPath) (map[string]exportedObj, bool) ***REMOVED***
	conn.handlersLck.RLock()
	defer conn.handlersLck.RUnlock()

	handlers, ok := conn.handlers[path]
	if ok ***REMOVED***
		return handlers, ok
	***REMOVED***

	// If handlers weren't found for this exact path, look for a matching subtree
	// registration
	handlers = make(map[string]exportedObj)
	path = path[:strings.LastIndex(string(path), "/")]
	for len(path) > 0 ***REMOVED***
		var subtreeHandlers map[string]exportedObj
		subtreeHandlers, ok = conn.handlers[path]
		if ok ***REMOVED***
			for iface, handler := range subtreeHandlers ***REMOVED***
				// Only include this handler if it registered for the subtree
				if handler.includeSubtree ***REMOVED***
					handlers[iface] = handler
				***REMOVED***
			***REMOVED***

			break
		***REMOVED***

		path = path[:strings.LastIndex(string(path), "/")]
	***REMOVED***

	return handlers, ok
***REMOVED***

// handleCall handles the given method call (i.e. looks if it's one of the
// pre-implemented ones and searches for a corresponding handler if not).
func (conn *Conn) handleCall(msg *Message) ***REMOVED***
	name := msg.Headers[FieldMember].value.(string)
	path := msg.Headers[FieldPath].value.(ObjectPath)
	ifaceName, hasIface := msg.Headers[FieldInterface].value.(string)
	sender, hasSender := msg.Headers[FieldSender].value.(string)
	serial := msg.serial
	if ifaceName == "org.freedesktop.DBus.Peer" ***REMOVED***
		switch name ***REMOVED***
		case "Ping":
			conn.sendReply(sender, serial)
		case "GetMachineId":
			conn.sendReply(sender, serial, conn.uuid)
		default:
			conn.sendError(errmsgUnknownMethod, sender, serial)
		***REMOVED***
		return
	***REMOVED*** else if ifaceName == "org.freedesktop.DBus.Introspectable" && name == "Introspect" ***REMOVED***
		if _, ok := conn.handlers[path]; !ok ***REMOVED***
			subpath := make(map[string]struct***REMOVED******REMOVED***)
			var xml bytes.Buffer
			xml.WriteString("<node>")
			for h, _ := range conn.handlers ***REMOVED***
				p := string(path)
				if p != "/" ***REMOVED***
					p += "/"
				***REMOVED***
				if strings.HasPrefix(string(h), p) ***REMOVED***
					node_name := strings.Split(string(h[len(p):]), "/")[0]
					subpath[node_name] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
				***REMOVED***
			***REMOVED***
			for s, _ := range subpath ***REMOVED***
				xml.WriteString("\n\t<node name=\"" + s + "\"/>")
			***REMOVED***
			xml.WriteString("\n</node>")
			conn.sendReply(sender, serial, xml.String())
			return
		***REMOVED***
	***REMOVED***
	if len(name) == 0 ***REMOVED***
		conn.sendError(errmsgUnknownMethod, sender, serial)
	***REMOVED***

	// Find the exported handler (if any) for this path
	handlers, ok := conn.searchHandlers(path)
	if !ok ***REMOVED***
		conn.sendError(errmsgNoObject, sender, serial)
		return
	***REMOVED***

	var m reflect.Value
	var exists bool
	if hasIface ***REMOVED***
		iface := handlers[ifaceName]
		m, exists = iface.Method(name)
	***REMOVED*** else ***REMOVED***
		for _, v := range handlers ***REMOVED***
			m, exists = v.Method(name)
			if exists ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if !exists ***REMOVED***
		conn.sendError(errmsgUnknownMethod, sender, serial)
		return
	***REMOVED***

	t := m.Type()
	vs := msg.Body
	pointers := make([]interface***REMOVED******REMOVED***, t.NumIn())
	decode := make([]interface***REMOVED******REMOVED***, 0, len(vs))
	for i := 0; i < t.NumIn(); i++ ***REMOVED***
		tp := t.In(i)
		val := reflect.New(tp)
		pointers[i] = val.Interface()
		if tp == reflect.TypeOf((*Sender)(nil)).Elem() ***REMOVED***
			val.Elem().SetString(sender)
		***REMOVED*** else if tp == reflect.TypeOf((*Message)(nil)).Elem() ***REMOVED***
			val.Elem().Set(reflect.ValueOf(*msg))
		***REMOVED*** else ***REMOVED***
			decode = append(decode, pointers[i])
		***REMOVED***
	***REMOVED***

	if len(decode) != len(vs) ***REMOVED***
		conn.sendError(errmsgInvalidArg, sender, serial)
		return
	***REMOVED***

	if err := Store(vs, decode...); err != nil ***REMOVED***
		conn.sendError(errmsgInvalidArg, sender, serial)
		return
	***REMOVED***

	// Extract parameters
	params := make([]reflect.Value, len(pointers))
	for i := 0; i < len(pointers); i++ ***REMOVED***
		params[i] = reflect.ValueOf(pointers[i]).Elem()
	***REMOVED***

	// Call method
	ret := m.Call(params)
	if em := ret[t.NumOut()-1].Interface().(*Error); em != nil ***REMOVED***
		conn.sendError(*em, sender, serial)
		return
	***REMOVED***

	if msg.Flags&FlagNoReplyExpected == 0 ***REMOVED***
		reply := new(Message)
		reply.Type = TypeMethodReply
		reply.serial = conn.getSerial()
		reply.Headers = make(map[HeaderField]Variant)
		if hasSender ***REMOVED***
			reply.Headers[FieldDestination] = msg.Headers[FieldSender]
		***REMOVED***
		reply.Headers[FieldReplySerial] = MakeVariant(msg.serial)
		reply.Body = make([]interface***REMOVED******REMOVED***, len(ret)-1)
		for i := 0; i < len(ret)-1; i++ ***REMOVED***
			reply.Body[i] = ret[i].Interface()
		***REMOVED***
		if len(ret) != 1 ***REMOVED***
			reply.Headers[FieldSignature] = MakeVariant(SignatureOf(reply.Body...))
		***REMOVED***
		conn.outLck.RLock()
		if !conn.closed ***REMOVED***
			conn.out <- reply
		***REMOVED***
		conn.outLck.RUnlock()
	***REMOVED***
***REMOVED***

// Emit emits the given signal on the message bus. The name parameter must be
// formatted as "interface.member", e.g., "org.freedesktop.DBus.NameLost".
func (conn *Conn) Emit(path ObjectPath, name string, values ...interface***REMOVED******REMOVED***) error ***REMOVED***
	if !path.IsValid() ***REMOVED***
		return errors.New("dbus: invalid object path")
	***REMOVED***
	i := strings.LastIndex(name, ".")
	if i == -1 ***REMOVED***
		return errors.New("dbus: invalid method name")
	***REMOVED***
	iface := name[:i]
	member := name[i+1:]
	if !isValidMember(member) ***REMOVED***
		return errors.New("dbus: invalid method name")
	***REMOVED***
	if !isValidInterface(iface) ***REMOVED***
		return errors.New("dbus: invalid interface name")
	***REMOVED***
	msg := new(Message)
	msg.Type = TypeSignal
	msg.serial = conn.getSerial()
	msg.Headers = make(map[HeaderField]Variant)
	msg.Headers[FieldInterface] = MakeVariant(iface)
	msg.Headers[FieldMember] = MakeVariant(member)
	msg.Headers[FieldPath] = MakeVariant(path)
	msg.Body = values
	if len(values) > 0 ***REMOVED***
		msg.Headers[FieldSignature] = MakeVariant(SignatureOf(values...))
	***REMOVED***
	conn.outLck.RLock()
	defer conn.outLck.RUnlock()
	if conn.closed ***REMOVED***
		return ErrClosed
	***REMOVED***
	conn.out <- msg
	return nil
***REMOVED***

// Export registers the given value to be exported as an object on the
// message bus.
//
// If a method call on the given path and interface is received, an exported
// method with the same name is called with v as the receiver if the
// parameters match and the last return value is of type *Error. If this
// *Error is not nil, it is sent back to the caller as an error.
// Otherwise, a method reply is sent with the other return values as its body.
//
// Any parameters with the special type Sender are set to the sender of the
// dbus message when the method is called. Parameters of this type do not
// contribute to the dbus signature of the method (i.e. the method is exposed
// as if the parameters of type Sender were not there).
//
// Similarly, any parameters with the type Message are set to the raw message
// received on the bus. Again, parameters of this type do not contribute to the
// dbus signature of the method.
//
// Every method call is executed in a new goroutine, so the method may be called
// in multiple goroutines at once.
//
// Method calls on the interface org.freedesktop.DBus.Peer will be automatically
// handled for every object.
//
// Passing nil as the first parameter will cause conn to cease handling calls on
// the given combination of path and interface.
//
// Export returns an error if path is not a valid path name.
func (conn *Conn) Export(v interface***REMOVED******REMOVED***, path ObjectPath, iface string) error ***REMOVED***
	return conn.ExportWithMap(v, nil, path, iface)
***REMOVED***

// ExportWithMap works exactly like Export but provides the ability to remap
// method names (e.g. export a lower-case method).
//
// The keys in the map are the real method names (exported on the struct), and
// the values are the method names to be exported on DBus.
func (conn *Conn) ExportWithMap(v interface***REMOVED******REMOVED***, mapping map[string]string, path ObjectPath, iface string) error ***REMOVED***
	return conn.export(getMethods(v, mapping), path, iface, false)
***REMOVED***

// ExportSubtree works exactly like Export but registers the given value for
// an entire subtree rather under the root path provided.
//
// In order to make this useful, one parameter in each of the value's exported
// methods should be a Message, in which case it will contain the raw message
// (allowing one to get access to the path that caused the method to be called).
//
// Note that more specific export paths take precedence over less specific. For
// example, a method call using the ObjectPath /foo/bar/baz will call a method
// exported on /foo/bar before a method exported on /foo.
func (conn *Conn) ExportSubtree(v interface***REMOVED******REMOVED***, path ObjectPath, iface string) error ***REMOVED***
	return conn.ExportSubtreeWithMap(v, nil, path, iface)
***REMOVED***

// ExportSubtreeWithMap works exactly like ExportSubtree but provides the
// ability to remap method names (e.g. export a lower-case method).
//
// The keys in the map are the real method names (exported on the struct), and
// the values are the method names to be exported on DBus.
func (conn *Conn) ExportSubtreeWithMap(v interface***REMOVED******REMOVED***, mapping map[string]string, path ObjectPath, iface string) error ***REMOVED***
	return conn.export(getMethods(v, mapping), path, iface, true)
***REMOVED***

// ExportMethodTable like Export registers the given methods as an object
// on the message bus. Unlike Export the it uses a method table to define
// the object instead of a native go object.
//
// The method table is a map from method name to function closure
// representing the method. This allows an object exported on the bus to not
// necessarily be a native go object. It can be useful for generating exposed
// methods on the fly.
//
// Any non-function objects in the method table are ignored.
func (conn *Conn) ExportMethodTable(methods map[string]interface***REMOVED******REMOVED***, path ObjectPath, iface string) error ***REMOVED***
	return conn.exportMethodTable(methods, path, iface, false)
***REMOVED***

// Like ExportSubtree, but with the same caveats as ExportMethodTable.
func (conn *Conn) ExportSubtreeMethodTable(methods map[string]interface***REMOVED******REMOVED***, path ObjectPath, iface string) error ***REMOVED***
	return conn.exportMethodTable(methods, path, iface, true)
***REMOVED***

func (conn *Conn) exportMethodTable(methods map[string]interface***REMOVED******REMOVED***, path ObjectPath, iface string, includeSubtree bool) error ***REMOVED***
	out := make(map[string]reflect.Value)
	for name, method := range methods ***REMOVED***
		rval := reflect.ValueOf(method)
		if rval.Kind() != reflect.Func ***REMOVED***
			continue
		***REMOVED***
		t := rval.Type()
		// only track valid methods must return *Error as last arg
		if t.NumOut() == 0 ||
			t.Out(t.NumOut()-1) != reflect.TypeOf(&errmsgInvalidArg) ***REMOVED***
			continue
		***REMOVED***
		out[name] = rval
	***REMOVED***
	return conn.export(out, path, iface, includeSubtree)
***REMOVED***

// exportWithMap is the worker function for all exports/registrations.
func (conn *Conn) export(methods map[string]reflect.Value, path ObjectPath, iface string, includeSubtree bool) error ***REMOVED***
	if !path.IsValid() ***REMOVED***
		return fmt.Errorf(`dbus: Invalid path name: "%s"`, path)
	***REMOVED***

	conn.handlersLck.Lock()
	defer conn.handlersLck.Unlock()

	// Remove a previous export if the interface is nil
	if methods == nil ***REMOVED***
		if _, ok := conn.handlers[path]; ok ***REMOVED***
			delete(conn.handlers[path], iface)
			if len(conn.handlers[path]) == 0 ***REMOVED***
				delete(conn.handlers, path)
			***REMOVED***
		***REMOVED***

		return nil
	***REMOVED***

	// If this is the first handler for this path, make a new map to hold all
	// handlers for this path.
	if _, ok := conn.handlers[path]; !ok ***REMOVED***
		conn.handlers[path] = make(map[string]exportedObj)
	***REMOVED***

	// Finally, save this handler
	conn.handlers[path][iface] = exportedObj***REMOVED***
		methods:        methods,
		includeSubtree: includeSubtree,
	***REMOVED***

	return nil
***REMOVED***

// ReleaseName calls org.freedesktop.DBus.ReleaseName and awaits a response.
func (conn *Conn) ReleaseName(name string) (ReleaseNameReply, error) ***REMOVED***
	var r uint32
	err := conn.busObj.Call("org.freedesktop.DBus.ReleaseName", 0, name).Store(&r)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return ReleaseNameReply(r), nil
***REMOVED***

// RequestName calls org.freedesktop.DBus.RequestName and awaits a response.
func (conn *Conn) RequestName(name string, flags RequestNameFlags) (RequestNameReply, error) ***REMOVED***
	var r uint32
	err := conn.busObj.Call("org.freedesktop.DBus.RequestName", 0, name, flags).Store(&r)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return RequestNameReply(r), nil
***REMOVED***

// ReleaseNameReply is the reply to a ReleaseName call.
type ReleaseNameReply uint32

const (
	ReleaseNameReplyReleased ReleaseNameReply = 1 + iota
	ReleaseNameReplyNonExistent
	ReleaseNameReplyNotOwner
)

// RequestNameFlags represents the possible flags for a RequestName call.
type RequestNameFlags uint32

const (
	NameFlagAllowReplacement RequestNameFlags = 1 << iota
	NameFlagReplaceExisting
	NameFlagDoNotQueue
)

// RequestNameReply is the reply to a RequestName call.
type RequestNameReply uint32

const (
	RequestNameReplyPrimaryOwner RequestNameReply = 1 + iota
	RequestNameReplyInQueue
	RequestNameReplyExists
	RequestNameReplyAlreadyOwner
)
