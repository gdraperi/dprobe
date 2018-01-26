// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Integration with the systemd D-Bus API.  See http://www.freedesktop.org/wiki/Software/systemd/dbus/
package dbus

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/godbus/dbus"
)

const (
	alpha        = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
	num          = `0123456789`
	alphanum     = alpha + num
	signalBuffer = 100
)

// needsEscape checks whether a byte in a potential dbus ObjectPath needs to be escaped
func needsEscape(i int, b byte) bool ***REMOVED***
	// Escape everything that is not a-z-A-Z-0-9
	// Also escape 0-9 if it's the first character
	return strings.IndexByte(alphanum, b) == -1 ||
		(i == 0 && strings.IndexByte(num, b) != -1)
***REMOVED***

// PathBusEscape sanitizes a constituent string of a dbus ObjectPath using the
// rules that systemd uses for serializing special characters.
func PathBusEscape(path string) string ***REMOVED***
	// Special case the empty string
	if len(path) == 0 ***REMOVED***
		return "_"
	***REMOVED***
	n := []byte***REMOVED******REMOVED***
	for i := 0; i < len(path); i++ ***REMOVED***
		c := path[i]
		if needsEscape(i, c) ***REMOVED***
			e := fmt.Sprintf("_%x", c)
			n = append(n, []byte(e)...)
		***REMOVED*** else ***REMOVED***
			n = append(n, c)
		***REMOVED***
	***REMOVED***
	return string(n)
***REMOVED***

// Conn is a connection to systemd's dbus endpoint.
type Conn struct ***REMOVED***
	// sysconn/sysobj are only used to call dbus methods
	sysconn *dbus.Conn
	sysobj  dbus.BusObject

	// sigconn/sigobj are only used to receive dbus signals
	sigconn *dbus.Conn
	sigobj  dbus.BusObject

	jobListener struct ***REMOVED***
		jobs map[dbus.ObjectPath]chan<- string
		sync.Mutex
	***REMOVED***
	subscriber struct ***REMOVED***
		updateCh chan<- *SubStateUpdate
		errCh    chan<- error
		sync.Mutex
		ignore      map[dbus.ObjectPath]int64
		cleanIgnore int64
	***REMOVED***
***REMOVED***

// New establishes a connection to any available bus and authenticates.
// Callers should call Close() when done with the connection.
func New() (*Conn, error) ***REMOVED***
	conn, err := NewSystemConnection()
	if err != nil && os.Geteuid() == 0 ***REMOVED***
		return NewSystemdConnection()
	***REMOVED***
	return conn, err
***REMOVED***

// NewSystemConnection establishes a connection to the system bus and authenticates.
// Callers should call Close() when done with the connection
func NewSystemConnection() (*Conn, error) ***REMOVED***
	return NewConnection(func() (*dbus.Conn, error) ***REMOVED***
		return dbusAuthHelloConnection(dbus.SystemBusPrivate)
	***REMOVED***)
***REMOVED***

// NewUserConnection establishes a connection to the session bus and
// authenticates. This can be used to connect to systemd user instances.
// Callers should call Close() when done with the connection.
func NewUserConnection() (*Conn, error) ***REMOVED***
	return NewConnection(func() (*dbus.Conn, error) ***REMOVED***
		return dbusAuthHelloConnection(dbus.SessionBusPrivate)
	***REMOVED***)
***REMOVED***

// NewSystemdConnection establishes a private, direct connection to systemd.
// This can be used for communicating with systemd without a dbus daemon.
// Callers should call Close() when done with the connection.
func NewSystemdConnection() (*Conn, error) ***REMOVED***
	return NewConnection(func() (*dbus.Conn, error) ***REMOVED***
		// We skip Hello when talking directly to systemd.
		return dbusAuthConnection(func() (*dbus.Conn, error) ***REMOVED***
			return dbus.Dial("unix:path=/run/systemd/private")
		***REMOVED***)
	***REMOVED***)
***REMOVED***

// Close closes an established connection
func (c *Conn) Close() ***REMOVED***
	c.sysconn.Close()
	c.sigconn.Close()
***REMOVED***

// NewConnection establishes a connection to a bus using a caller-supplied function.
// This allows connecting to remote buses through a user-supplied mechanism.
// The supplied function may be called multiple times, and should return independent connections.
// The returned connection must be fully initialised: the org.freedesktop.DBus.Hello call must have succeeded,
// and any authentication should be handled by the function.
func NewConnection(dialBus func() (*dbus.Conn, error)) (*Conn, error) ***REMOVED***
	sysconn, err := dialBus()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	sigconn, err := dialBus()
	if err != nil ***REMOVED***
		sysconn.Close()
		return nil, err
	***REMOVED***

	c := &Conn***REMOVED***
		sysconn: sysconn,
		sysobj:  systemdObject(sysconn),
		sigconn: sigconn,
		sigobj:  systemdObject(sigconn),
	***REMOVED***

	c.subscriber.ignore = make(map[dbus.ObjectPath]int64)
	c.jobListener.jobs = make(map[dbus.ObjectPath]chan<- string)

	// Setup the listeners on jobs so that we can get completions
	c.sigconn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
		"type='signal', interface='org.freedesktop.systemd1.Manager', member='JobRemoved'")

	c.dispatch()
	return c, nil
***REMOVED***

// GetManagerProperty returns the value of a property on the org.freedesktop.systemd1.Manager
// interface. The value is returned in its string representation, as defined at
// https://developer.gnome.org/glib/unstable/gvariant-text.html
func (c *Conn) GetManagerProperty(prop string) (string, error) ***REMOVED***
	variant, err := c.sysobj.GetProperty("org.freedesktop.systemd1.Manager." + prop)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return variant.String(), nil
***REMOVED***

func dbusAuthConnection(createBus func() (*dbus.Conn, error)) (*dbus.Conn, error) ***REMOVED***
	conn, err := createBus()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Only use EXTERNAL method, and hardcode the uid (not username)
	// to avoid a username lookup (which requires a dynamically linked
	// libc)
	methods := []dbus.Auth***REMOVED***dbus.AuthExternal(strconv.Itoa(os.Getuid()))***REMOVED***

	err = conn.Auth(methods)
	if err != nil ***REMOVED***
		conn.Close()
		return nil, err
	***REMOVED***

	return conn, nil
***REMOVED***

func dbusAuthHelloConnection(createBus func() (*dbus.Conn, error)) (*dbus.Conn, error) ***REMOVED***
	conn, err := dbusAuthConnection(createBus)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err = conn.Hello(); err != nil ***REMOVED***
		conn.Close()
		return nil, err
	***REMOVED***

	return conn, nil
***REMOVED***

func systemdObject(conn *dbus.Conn) dbus.BusObject ***REMOVED***
	return conn.Object("org.freedesktop.systemd1", dbus.ObjectPath("/org/freedesktop/systemd1"))
***REMOVED***
