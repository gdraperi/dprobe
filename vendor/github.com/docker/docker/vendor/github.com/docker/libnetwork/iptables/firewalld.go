package iptables

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus"
	"github.com/sirupsen/logrus"
)

// IPV defines the table string
type IPV string

const (
	// Iptables point ipv4 table
	Iptables IPV = "ipv4"
	// IP6Tables point to ipv6 table
	IP6Tables IPV = "ipv6"
	// Ebtables point to bridge table
	Ebtables IPV = "eb"
)
const (
	dbusInterface = "org.fedoraproject.FirewallD1"
	dbusPath      = "/org/fedoraproject/FirewallD1"
)

// Conn is a connection to firewalld dbus endpoint.
type Conn struct ***REMOVED***
	sysconn *dbus.Conn
	sysobj  dbus.BusObject
	signal  chan *dbus.Signal
***REMOVED***

var (
	connection       *Conn
	firewalldRunning bool      // is Firewalld service running
	onReloaded       []*func() // callbacks when Firewalld has been reloaded
)

// FirewalldInit initializes firewalld management code.
func FirewalldInit() error ***REMOVED***
	var err error

	if connection, err = newConnection(); err != nil ***REMOVED***
		return fmt.Errorf("Failed to connect to D-Bus system bus: %v", err)
	***REMOVED***
	firewalldRunning = checkRunning()
	if !firewalldRunning ***REMOVED***
		connection.sysconn.Close()
		connection = nil
	***REMOVED***
	if connection != nil ***REMOVED***
		go signalHandler()
	***REMOVED***

	return nil
***REMOVED***

// New() establishes a connection to the system bus.
func newConnection() (*Conn, error) ***REMOVED***
	c := new(Conn)
	if err := c.initConnection(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return c, nil
***REMOVED***

// Innitialize D-Bus connection.
func (c *Conn) initConnection() error ***REMOVED***
	var err error

	c.sysconn, err = dbus.SystemBus()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// This never fails, even if the service is not running atm.
	c.sysobj = c.sysconn.Object(dbusInterface, dbus.ObjectPath(dbusPath))

	rule := fmt.Sprintf("type='signal',path='%s',interface='%s',sender='%s',member='Reloaded'",
		dbusPath, dbusInterface, dbusInterface)
	c.sysconn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, rule)

	rule = fmt.Sprintf("type='signal',interface='org.freedesktop.DBus',member='NameOwnerChanged',path='/org/freedesktop/DBus',sender='org.freedesktop.DBus',arg0='%s'",
		dbusInterface)
	c.sysconn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, rule)

	c.signal = make(chan *dbus.Signal, 10)
	c.sysconn.Signal(c.signal)

	return nil
***REMOVED***

func signalHandler() ***REMOVED***
	for signal := range connection.signal ***REMOVED***
		if strings.Contains(signal.Name, "NameOwnerChanged") ***REMOVED***
			firewalldRunning = checkRunning()
			dbusConnectionChanged(signal.Body)
		***REMOVED*** else if strings.Contains(signal.Name, "Reloaded") ***REMOVED***
			reloaded()
		***REMOVED***
	***REMOVED***
***REMOVED***

func dbusConnectionChanged(args []interface***REMOVED******REMOVED***) ***REMOVED***
	name := args[0].(string)
	oldOwner := args[1].(string)
	newOwner := args[2].(string)

	if name != dbusInterface ***REMOVED***
		return
	***REMOVED***

	if len(newOwner) > 0 ***REMOVED***
		connectionEstablished()
	***REMOVED*** else if len(oldOwner) > 0 ***REMOVED***
		connectionLost()
	***REMOVED***
***REMOVED***

func connectionEstablished() ***REMOVED***
	reloaded()
***REMOVED***

func connectionLost() ***REMOVED***
	// Doesn't do anything for now. Libvirt also doesn't react to this.
***REMOVED***

// call all callbacks
func reloaded() ***REMOVED***
	for _, pf := range onReloaded ***REMOVED***
		(*pf)()
	***REMOVED***
***REMOVED***

// OnReloaded add callback
func OnReloaded(callback func()) ***REMOVED***
	for _, pf := range onReloaded ***REMOVED***
		if pf == &callback ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	onReloaded = append(onReloaded, &callback)
***REMOVED***

// Call some remote method to see whether the service is actually running.
func checkRunning() bool ***REMOVED***
	var zone string
	var err error

	if connection != nil ***REMOVED***
		err = connection.sysobj.Call(dbusInterface+".getDefaultZone", 0).Store(&zone)
		return err == nil
	***REMOVED***
	return false
***REMOVED***

// Passthrough method simply passes args through to iptables/ip6tables
func Passthrough(ipv IPV, args ...string) ([]byte, error) ***REMOVED***
	var output string
	logrus.Debugf("Firewalld passthrough: %s, %s", ipv, args)
	if err := connection.sysobj.Call(dbusInterface+".direct.passthrough", 0, ipv, args).Store(&output); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return []byte(output), nil
***REMOVED***
