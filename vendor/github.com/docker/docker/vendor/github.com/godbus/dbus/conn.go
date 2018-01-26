package dbus

import (
	"errors"
	"io"
	"os"
	"reflect"
	"strings"
	"sync"
)

const defaultSystemBusAddress = "unix:path=/var/run/dbus/system_bus_socket"

var (
	systemBus     *Conn
	systemBusLck  sync.Mutex
	sessionBus    *Conn
	sessionBusLck sync.Mutex
	sessionEnvLck sync.Mutex
)

// ErrClosed is the error returned by calls on a closed connection.
var ErrClosed = errors.New("dbus: connection closed by user")

// Conn represents a connection to a message bus (usually, the system or
// session bus).
//
// Connections are either shared or private. Shared connections
// are shared between calls to the functions that return them. As a result,
// the methods Close, Auth and Hello must not be called on them.
//
// Multiple goroutines may invoke methods on a connection simultaneously.
type Conn struct ***REMOVED***
	transport

	busObj BusObject
	unixFD bool
	uuid   string

	names    []string
	namesLck sync.RWMutex

	serialLck  sync.Mutex
	nextSerial uint32
	serialUsed map[uint32]bool

	calls    map[uint32]*Call
	callsLck sync.RWMutex

	handlers    map[ObjectPath]map[string]exportedObj
	handlersLck sync.RWMutex

	out    chan *Message
	closed bool
	outLck sync.RWMutex

	signals    []chan<- *Signal
	signalsLck sync.Mutex

	eavesdropped    chan<- *Message
	eavesdroppedLck sync.Mutex
***REMOVED***

// SessionBus returns a shared connection to the session bus, connecting to it
// if not already done.
func SessionBus() (conn *Conn, err error) ***REMOVED***
	sessionBusLck.Lock()
	defer sessionBusLck.Unlock()
	if sessionBus != nil ***REMOVED***
		return sessionBus, nil
	***REMOVED***
	defer func() ***REMOVED***
		if conn != nil ***REMOVED***
			sessionBus = conn
		***REMOVED***
	***REMOVED***()
	conn, err = SessionBusPrivate()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if err = conn.Auth(nil); err != nil ***REMOVED***
		conn.Close()
		conn = nil
		return
	***REMOVED***
	if err = conn.Hello(); err != nil ***REMOVED***
		conn.Close()
		conn = nil
	***REMOVED***
	return
***REMOVED***

// SessionBusPrivate returns a new private connection to the session bus.
func SessionBusPrivate() (*Conn, error) ***REMOVED***
	sessionEnvLck.Lock()
	defer sessionEnvLck.Unlock()
	address := os.Getenv("DBUS_SESSION_BUS_ADDRESS")
	if address != "" && address != "autolaunch:" ***REMOVED***
		return Dial(address)
	***REMOVED***

	return sessionBusPlatform()
***REMOVED***

// SystemBus returns a shared connection to the system bus, connecting to it if
// not already done.
func SystemBus() (conn *Conn, err error) ***REMOVED***
	systemBusLck.Lock()
	defer systemBusLck.Unlock()
	if systemBus != nil ***REMOVED***
		return systemBus, nil
	***REMOVED***
	defer func() ***REMOVED***
		if conn != nil ***REMOVED***
			systemBus = conn
		***REMOVED***
	***REMOVED***()
	conn, err = SystemBusPrivate()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if err = conn.Auth(nil); err != nil ***REMOVED***
		conn.Close()
		conn = nil
		return
	***REMOVED***
	if err = conn.Hello(); err != nil ***REMOVED***
		conn.Close()
		conn = nil
	***REMOVED***
	return
***REMOVED***

// SystemBusPrivate returns a new private connection to the system bus.
func SystemBusPrivate() (*Conn, error) ***REMOVED***
	address := os.Getenv("DBUS_SYSTEM_BUS_ADDRESS")
	if address != "" ***REMOVED***
		return Dial(address)
	***REMOVED***
	return Dial(defaultSystemBusAddress)
***REMOVED***

// Dial establishes a new private connection to the message bus specified by address.
func Dial(address string) (*Conn, error) ***REMOVED***
	tr, err := getTransport(address)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return newConn(tr)
***REMOVED***

// NewConn creates a new private *Conn from an already established connection.
func NewConn(conn io.ReadWriteCloser) (*Conn, error) ***REMOVED***
	return newConn(genericTransport***REMOVED***conn***REMOVED***)
***REMOVED***

// newConn creates a new *Conn from a transport.
func newConn(tr transport) (*Conn, error) ***REMOVED***
	conn := new(Conn)
	conn.transport = tr
	conn.calls = make(map[uint32]*Call)
	conn.out = make(chan *Message, 10)
	conn.handlers = make(map[ObjectPath]map[string]exportedObj)
	conn.nextSerial = 1
	conn.serialUsed = map[uint32]bool***REMOVED***0: true***REMOVED***
	conn.busObj = conn.Object("org.freedesktop.DBus", "/org/freedesktop/DBus")
	return conn, nil
***REMOVED***

// BusObject returns the object owned by the bus daemon which handles
// administrative requests.
func (conn *Conn) BusObject() BusObject ***REMOVED***
	return conn.busObj
***REMOVED***

// Close closes the connection. Any blocked operations will return with errors
// and the channels passed to Eavesdrop and Signal are closed. This method must
// not be called on shared connections.
func (conn *Conn) Close() error ***REMOVED***
	conn.outLck.Lock()
	if conn.closed ***REMOVED***
		// inWorker calls Close on read error, the read error may
		// be caused by another caller calling Close to shutdown the
		// dbus connection, a double-close scenario we prevent here.
		conn.outLck.Unlock()
		return nil
	***REMOVED***
	close(conn.out)
	conn.closed = true
	conn.outLck.Unlock()
	conn.signalsLck.Lock()
	for _, ch := range conn.signals ***REMOVED***
		close(ch)
	***REMOVED***
	conn.signalsLck.Unlock()
	conn.eavesdroppedLck.Lock()
	if conn.eavesdropped != nil ***REMOVED***
		close(conn.eavesdropped)
	***REMOVED***
	conn.eavesdroppedLck.Unlock()
	return conn.transport.Close()
***REMOVED***

// Eavesdrop causes conn to send all incoming messages to the given channel
// without further processing. Method replies, errors and signals will not be
// sent to the appropiate channels and method calls will not be handled. If nil
// is passed, the normal behaviour is restored.
//
// The caller has to make sure that ch is sufficiently buffered;
// if a message arrives when a write to ch is not possible, the message is
// discarded.
func (conn *Conn) Eavesdrop(ch chan<- *Message) ***REMOVED***
	conn.eavesdroppedLck.Lock()
	conn.eavesdropped = ch
	conn.eavesdroppedLck.Unlock()
***REMOVED***

// getSerial returns an unused serial.
func (conn *Conn) getSerial() uint32 ***REMOVED***
	conn.serialLck.Lock()
	defer conn.serialLck.Unlock()
	n := conn.nextSerial
	for conn.serialUsed[n] ***REMOVED***
		n++
	***REMOVED***
	conn.serialUsed[n] = true
	conn.nextSerial = n + 1
	return n
***REMOVED***

// Hello sends the initial org.freedesktop.DBus.Hello call. This method must be
// called after authentication, but before sending any other messages to the
// bus. Hello must not be called for shared connections.
func (conn *Conn) Hello() error ***REMOVED***
	var s string
	err := conn.busObj.Call("org.freedesktop.DBus.Hello", 0).Store(&s)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	conn.namesLck.Lock()
	conn.names = make([]string, 1)
	conn.names[0] = s
	conn.namesLck.Unlock()
	return nil
***REMOVED***

// inWorker runs in an own goroutine, reading incoming messages from the
// transport and dispatching them appropiately.
func (conn *Conn) inWorker() ***REMOVED***
	for ***REMOVED***
		msg, err := conn.ReadMessage()
		if err == nil ***REMOVED***
			conn.eavesdroppedLck.Lock()
			if conn.eavesdropped != nil ***REMOVED***
				select ***REMOVED***
				case conn.eavesdropped <- msg:
				default:
				***REMOVED***
				conn.eavesdroppedLck.Unlock()
				continue
			***REMOVED***
			conn.eavesdroppedLck.Unlock()
			dest, _ := msg.Headers[FieldDestination].value.(string)
			found := false
			if dest == "" ***REMOVED***
				found = true
			***REMOVED*** else ***REMOVED***
				conn.namesLck.RLock()
				if len(conn.names) == 0 ***REMOVED***
					found = true
				***REMOVED***
				for _, v := range conn.names ***REMOVED***
					if dest == v ***REMOVED***
						found = true
						break
					***REMOVED***
				***REMOVED***
				conn.namesLck.RUnlock()
			***REMOVED***
			if !found ***REMOVED***
				// Eavesdropped a message, but no channel for it is registered.
				// Ignore it.
				continue
			***REMOVED***
			switch msg.Type ***REMOVED***
			case TypeMethodReply, TypeError:
				serial := msg.Headers[FieldReplySerial].value.(uint32)
				conn.callsLck.Lock()
				if c, ok := conn.calls[serial]; ok ***REMOVED***
					if msg.Type == TypeError ***REMOVED***
						name, _ := msg.Headers[FieldErrorName].value.(string)
						c.Err = Error***REMOVED***name, msg.Body***REMOVED***
					***REMOVED*** else ***REMOVED***
						c.Body = msg.Body
					***REMOVED***
					c.Done <- c
					conn.serialLck.Lock()
					delete(conn.serialUsed, serial)
					conn.serialLck.Unlock()
					delete(conn.calls, serial)
				***REMOVED***
				conn.callsLck.Unlock()
			case TypeSignal:
				iface := msg.Headers[FieldInterface].value.(string)
				member := msg.Headers[FieldMember].value.(string)
				// as per http://dbus.freedesktop.org/doc/dbus-specification.html ,
				// sender is optional for signals.
				sender, _ := msg.Headers[FieldSender].value.(string)
				if iface == "org.freedesktop.DBus" && sender == "org.freedesktop.DBus" ***REMOVED***
					if member == "NameLost" ***REMOVED***
						// If we lost the name on the bus, remove it from our
						// tracking list.
						name, ok := msg.Body[0].(string)
						if !ok ***REMOVED***
							panic("Unable to read the lost name")
						***REMOVED***
						conn.namesLck.Lock()
						for i, v := range conn.names ***REMOVED***
							if v == name ***REMOVED***
								conn.names = append(conn.names[:i],
									conn.names[i+1:]...)
							***REMOVED***
						***REMOVED***
						conn.namesLck.Unlock()
					***REMOVED*** else if member == "NameAcquired" ***REMOVED***
						// If we acquired the name on the bus, add it to our
						// tracking list.
						name, ok := msg.Body[0].(string)
						if !ok ***REMOVED***
							panic("Unable to read the acquired name")
						***REMOVED***
						conn.namesLck.Lock()
						conn.names = append(conn.names, name)
						conn.namesLck.Unlock()
					***REMOVED***
				***REMOVED***
				signal := &Signal***REMOVED***
					Sender: sender,
					Path:   msg.Headers[FieldPath].value.(ObjectPath),
					Name:   iface + "." + member,
					Body:   msg.Body,
				***REMOVED***
				conn.signalsLck.Lock()
				for _, ch := range conn.signals ***REMOVED***
					ch <- signal
				***REMOVED***
				conn.signalsLck.Unlock()
			case TypeMethodCall:
				go conn.handleCall(msg)
			***REMOVED***
		***REMOVED*** else if _, ok := err.(InvalidMessageError); !ok ***REMOVED***
			// Some read error occured (usually EOF); we can't really do
			// anything but to shut down all stuff and returns errors to all
			// pending replies.
			conn.Close()
			conn.callsLck.RLock()
			for _, v := range conn.calls ***REMOVED***
				v.Err = err
				v.Done <- v
			***REMOVED***
			conn.callsLck.RUnlock()
			return
		***REMOVED***
		// invalid messages are ignored
	***REMOVED***
***REMOVED***

// Names returns the list of all names that are currently owned by this
// connection. The slice is always at least one element long, the first element
// being the unique name of the connection.
func (conn *Conn) Names() []string ***REMOVED***
	conn.namesLck.RLock()
	// copy the slice so it can't be modified
	s := make([]string, len(conn.names))
	copy(s, conn.names)
	conn.namesLck.RUnlock()
	return s
***REMOVED***

// Object returns the object identified by the given destination name and path.
func (conn *Conn) Object(dest string, path ObjectPath) BusObject ***REMOVED***
	return &Object***REMOVED***conn, dest, path***REMOVED***
***REMOVED***

// outWorker runs in an own goroutine, encoding and sending messages that are
// sent to conn.out.
func (conn *Conn) outWorker() ***REMOVED***
	for msg := range conn.out ***REMOVED***
		err := conn.SendMessage(msg)
		conn.callsLck.RLock()
		if err != nil ***REMOVED***
			if c := conn.calls[msg.serial]; c != nil ***REMOVED***
				c.Err = err
				c.Done <- c
			***REMOVED***
			conn.serialLck.Lock()
			delete(conn.serialUsed, msg.serial)
			conn.serialLck.Unlock()
		***REMOVED*** else if msg.Type != TypeMethodCall ***REMOVED***
			conn.serialLck.Lock()
			delete(conn.serialUsed, msg.serial)
			conn.serialLck.Unlock()
		***REMOVED***
		conn.callsLck.RUnlock()
	***REMOVED***
***REMOVED***

// Send sends the given message to the message bus. You usually don't need to
// use this; use the higher-level equivalents (Call / Go, Emit and Export)
// instead. If msg is a method call and NoReplyExpected is not set, a non-nil
// call is returned and the same value is sent to ch (which must be buffered)
// once the call is complete. Otherwise, ch is ignored and a Call structure is
// returned of which only the Err member is valid.
func (conn *Conn) Send(msg *Message, ch chan *Call) *Call ***REMOVED***
	var call *Call

	msg.serial = conn.getSerial()
	if msg.Type == TypeMethodCall && msg.Flags&FlagNoReplyExpected == 0 ***REMOVED***
		if ch == nil ***REMOVED***
			ch = make(chan *Call, 5)
		***REMOVED*** else if cap(ch) == 0 ***REMOVED***
			panic("dbus: unbuffered channel passed to (*Conn).Send")
		***REMOVED***
		call = new(Call)
		call.Destination, _ = msg.Headers[FieldDestination].value.(string)
		call.Path, _ = msg.Headers[FieldPath].value.(ObjectPath)
		iface, _ := msg.Headers[FieldInterface].value.(string)
		member, _ := msg.Headers[FieldMember].value.(string)
		call.Method = iface + "." + member
		call.Args = msg.Body
		call.Done = ch
		conn.callsLck.Lock()
		conn.calls[msg.serial] = call
		conn.callsLck.Unlock()
		conn.outLck.RLock()
		if conn.closed ***REMOVED***
			call.Err = ErrClosed
			call.Done <- call
		***REMOVED*** else ***REMOVED***
			conn.out <- msg
		***REMOVED***
		conn.outLck.RUnlock()
	***REMOVED*** else ***REMOVED***
		conn.outLck.RLock()
		if conn.closed ***REMOVED***
			call = &Call***REMOVED***Err: ErrClosed***REMOVED***
		***REMOVED*** else ***REMOVED***
			conn.out <- msg
			call = &Call***REMOVED***Err: nil***REMOVED***
		***REMOVED***
		conn.outLck.RUnlock()
	***REMOVED***
	return call
***REMOVED***

// sendError creates an error message corresponding to the parameters and sends
// it to conn.out.
func (conn *Conn) sendError(e Error, dest string, serial uint32) ***REMOVED***
	msg := new(Message)
	msg.Type = TypeError
	msg.serial = conn.getSerial()
	msg.Headers = make(map[HeaderField]Variant)
	if dest != "" ***REMOVED***
		msg.Headers[FieldDestination] = MakeVariant(dest)
	***REMOVED***
	msg.Headers[FieldErrorName] = MakeVariant(e.Name)
	msg.Headers[FieldReplySerial] = MakeVariant(serial)
	msg.Body = e.Body
	if len(e.Body) > 0 ***REMOVED***
		msg.Headers[FieldSignature] = MakeVariant(SignatureOf(e.Body...))
	***REMOVED***
	conn.outLck.RLock()
	if !conn.closed ***REMOVED***
		conn.out <- msg
	***REMOVED***
	conn.outLck.RUnlock()
***REMOVED***

// sendReply creates a method reply message corresponding to the parameters and
// sends it to conn.out.
func (conn *Conn) sendReply(dest string, serial uint32, values ...interface***REMOVED******REMOVED***) ***REMOVED***
	msg := new(Message)
	msg.Type = TypeMethodReply
	msg.serial = conn.getSerial()
	msg.Headers = make(map[HeaderField]Variant)
	if dest != "" ***REMOVED***
		msg.Headers[FieldDestination] = MakeVariant(dest)
	***REMOVED***
	msg.Headers[FieldReplySerial] = MakeVariant(serial)
	msg.Body = values
	if len(values) > 0 ***REMOVED***
		msg.Headers[FieldSignature] = MakeVariant(SignatureOf(values...))
	***REMOVED***
	conn.outLck.RLock()
	if !conn.closed ***REMOVED***
		conn.out <- msg
	***REMOVED***
	conn.outLck.RUnlock()
***REMOVED***

// Signal registers the given channel to be passed all received signal messages.
// The caller has to make sure that ch is sufficiently buffered; if a message
// arrives when a write to c is not possible, it is discarded.
//
// Multiple of these channels can be registered at the same time.
//
// These channels are "overwritten" by Eavesdrop; i.e., if there currently is a
// channel for eavesdropped messages, this channel receives all signals, and
// none of the channels passed to Signal will receive any signals.
func (conn *Conn) Signal(ch chan<- *Signal) ***REMOVED***
	conn.signalsLck.Lock()
	conn.signals = append(conn.signals, ch)
	conn.signalsLck.Unlock()
***REMOVED***

// RemoveSignal removes the given channel from the list of the registered channels.
func (conn *Conn) RemoveSignal(ch chan<- *Signal) ***REMOVED***
	conn.signalsLck.Lock()
	for i := len(conn.signals) - 1; i >= 0; i-- ***REMOVED***
		if ch == conn.signals[i] ***REMOVED***
			copy(conn.signals[i:], conn.signals[i+1:])
			conn.signals[len(conn.signals)-1] = nil
			conn.signals = conn.signals[:len(conn.signals)-1]
		***REMOVED***
	***REMOVED***
	conn.signalsLck.Unlock()
***REMOVED***

// SupportsUnixFDs returns whether the underlying transport supports passing of
// unix file descriptors. If this is false, method calls containing unix file
// descriptors will return an error and emitted signals containing them will
// not be sent.
func (conn *Conn) SupportsUnixFDs() bool ***REMOVED***
	return conn.unixFD
***REMOVED***

// Error represents a D-Bus message of type Error.
type Error struct ***REMOVED***
	Name string
	Body []interface***REMOVED******REMOVED***
***REMOVED***

func NewError(name string, body []interface***REMOVED******REMOVED***) *Error ***REMOVED***
	return &Error***REMOVED***name, body***REMOVED***
***REMOVED***

func (e Error) Error() string ***REMOVED***
	if len(e.Body) >= 1 ***REMOVED***
		s, ok := e.Body[0].(string)
		if ok ***REMOVED***
			return s
		***REMOVED***
	***REMOVED***
	return e.Name
***REMOVED***

// Signal represents a D-Bus message of type Signal. The name member is given in
// "interface.member" notation, e.g. org.freedesktop.D-Bus.NameLost.
type Signal struct ***REMOVED***
	Sender string
	Path   ObjectPath
	Name   string
	Body   []interface***REMOVED******REMOVED***
***REMOVED***

// transport is a D-Bus transport.
type transport interface ***REMOVED***
	// Read and Write raw data (for example, for the authentication protocol).
	io.ReadWriteCloser

	// Send the initial null byte used for the EXTERNAL mechanism.
	SendNullByte() error

	// Returns whether this transport supports passing Unix FDs.
	SupportsUnixFDs() bool

	// Signal the transport that Unix FD passing is enabled for this connection.
	EnableUnixFDs()

	// Read / send a message, handling things like Unix FDs.
	ReadMessage() (*Message, error)
	SendMessage(*Message) error
***REMOVED***

var (
	transports = make(map[string]func(string) (transport, error))
)

func getTransport(address string) (transport, error) ***REMOVED***
	var err error
	var t transport

	addresses := strings.Split(address, ";")
	for _, v := range addresses ***REMOVED***
		i := strings.IndexRune(v, ':')
		if i == -1 ***REMOVED***
			err = errors.New("dbus: invalid bus address (no transport)")
			continue
		***REMOVED***
		f := transports[v[:i]]
		if f == nil ***REMOVED***
			err = errors.New("dbus: invalid bus address (invalid or unsupported transport)")
			continue
		***REMOVED***
		t, err = f(v[i+1:])
		if err == nil ***REMOVED***
			return t, nil
		***REMOVED***
	***REMOVED***
	return nil, err
***REMOVED***

// dereferenceAll returns a slice that, assuming that vs is a slice of pointers
// of arbitrary types, containes the values that are obtained from dereferencing
// all elements in vs.
func dereferenceAll(vs []interface***REMOVED******REMOVED***) []interface***REMOVED******REMOVED*** ***REMOVED***
	for i := range vs ***REMOVED***
		v := reflect.ValueOf(vs[i])
		v = v.Elem()
		vs[i] = v.Interface()
	***REMOVED***
	return vs
***REMOVED***

// getKey gets a key from a the list of keys. Returns "" on error / not found...
func getKey(s, key string) string ***REMOVED***
	for _, keyEqualsValue := range strings.Split(s, ",") ***REMOVED***
		keyValue := strings.SplitN(keyEqualsValue, "=", 2)
		if len(keyValue) == 2 && keyValue[0] == key ***REMOVED***
			return keyValue[1]
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***
