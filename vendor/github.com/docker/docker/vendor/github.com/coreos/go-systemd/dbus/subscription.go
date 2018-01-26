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

package dbus

import (
	"errors"
	"time"

	"github.com/godbus/dbus"
)

const (
	cleanIgnoreInterval = int64(10 * time.Second)
	ignoreInterval      = int64(30 * time.Millisecond)
)

// Subscribe sets up this connection to subscribe to all systemd dbus events.
// This is required before calling SubscribeUnits. When the connection closes
// systemd will automatically stop sending signals so there is no need to
// explicitly call Unsubscribe().
func (c *Conn) Subscribe() error ***REMOVED***
	c.sigconn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
		"type='signal',interface='org.freedesktop.systemd1.Manager',member='UnitNew'")
	c.sigconn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
		"type='signal',interface='org.freedesktop.DBus.Properties',member='PropertiesChanged'")

	err := c.sigobj.Call("org.freedesktop.systemd1.Manager.Subscribe", 0).Store()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

// Unsubscribe this connection from systemd dbus events.
func (c *Conn) Unsubscribe() error ***REMOVED***
	err := c.sigobj.Call("org.freedesktop.systemd1.Manager.Unsubscribe", 0).Store()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (c *Conn) dispatch() ***REMOVED***
	ch := make(chan *dbus.Signal, signalBuffer)

	c.sigconn.Signal(ch)

	go func() ***REMOVED***
		for ***REMOVED***
			signal, ok := <-ch
			if !ok ***REMOVED***
				return
			***REMOVED***

			if signal.Name == "org.freedesktop.systemd1.Manager.JobRemoved" ***REMOVED***
				c.jobComplete(signal)
			***REMOVED***

			if c.subscriber.updateCh == nil ***REMOVED***
				continue
			***REMOVED***

			var unitPath dbus.ObjectPath
			switch signal.Name ***REMOVED***
			case "org.freedesktop.systemd1.Manager.JobRemoved":
				unitName := signal.Body[2].(string)
				c.sysobj.Call("org.freedesktop.systemd1.Manager.GetUnit", 0, unitName).Store(&unitPath)
			case "org.freedesktop.systemd1.Manager.UnitNew":
				unitPath = signal.Body[1].(dbus.ObjectPath)
			case "org.freedesktop.DBus.Properties.PropertiesChanged":
				if signal.Body[0].(string) == "org.freedesktop.systemd1.Unit" ***REMOVED***
					unitPath = signal.Path
				***REMOVED***
			***REMOVED***

			if unitPath == dbus.ObjectPath("") ***REMOVED***
				continue
			***REMOVED***

			c.sendSubStateUpdate(unitPath)
		***REMOVED***
	***REMOVED***()
***REMOVED***

// Returns two unbuffered channels which will receive all changed units every
// interval.  Deleted units are sent as nil.
func (c *Conn) SubscribeUnits(interval time.Duration) (<-chan map[string]*UnitStatus, <-chan error) ***REMOVED***
	return c.SubscribeUnitsCustom(interval, 0, func(u1, u2 *UnitStatus) bool ***REMOVED*** return *u1 != *u2 ***REMOVED***, nil)
***REMOVED***

// SubscribeUnitsCustom is like SubscribeUnits but lets you specify the buffer
// size of the channels, the comparison function for detecting changes and a filter
// function for cutting down on the noise that your channel receives.
func (c *Conn) SubscribeUnitsCustom(interval time.Duration, buffer int, isChanged func(*UnitStatus, *UnitStatus) bool, filterUnit func(string) bool) (<-chan map[string]*UnitStatus, <-chan error) ***REMOVED***
	old := make(map[string]*UnitStatus)
	statusChan := make(chan map[string]*UnitStatus, buffer)
	errChan := make(chan error, buffer)

	go func() ***REMOVED***
		for ***REMOVED***
			timerChan := time.After(interval)

			units, err := c.ListUnits()
			if err == nil ***REMOVED***
				cur := make(map[string]*UnitStatus)
				for i := range units ***REMOVED***
					if filterUnit != nil && filterUnit(units[i].Name) ***REMOVED***
						continue
					***REMOVED***
					cur[units[i].Name] = &units[i]
				***REMOVED***

				// add all new or changed units
				changed := make(map[string]*UnitStatus)
				for n, u := range cur ***REMOVED***
					if oldU, ok := old[n]; !ok || isChanged(oldU, u) ***REMOVED***
						changed[n] = u
					***REMOVED***
					delete(old, n)
				***REMOVED***

				// add all deleted units
				for oldN := range old ***REMOVED***
					changed[oldN] = nil
				***REMOVED***

				old = cur

				if len(changed) != 0 ***REMOVED***
					statusChan <- changed
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				errChan <- err
			***REMOVED***

			<-timerChan
		***REMOVED***
	***REMOVED***()

	return statusChan, errChan
***REMOVED***

type SubStateUpdate struct ***REMOVED***
	UnitName string
	SubState string
***REMOVED***

// SetSubStateSubscriber writes to updateCh when any unit's substate changes.
// Although this writes to updateCh on every state change, the reported state
// may be more recent than the change that generated it (due to an unavoidable
// race in the systemd dbus interface).  That is, this method provides a good
// way to keep a current view of all units' states, but is not guaranteed to
// show every state transition they go through.  Furthermore, state changes
// will only be written to the channel with non-blocking writes.  If updateCh
// is full, it attempts to write an error to errCh; if errCh is full, the error
// passes silently.
func (c *Conn) SetSubStateSubscriber(updateCh chan<- *SubStateUpdate, errCh chan<- error) ***REMOVED***
	c.subscriber.Lock()
	defer c.subscriber.Unlock()
	c.subscriber.updateCh = updateCh
	c.subscriber.errCh = errCh
***REMOVED***

func (c *Conn) sendSubStateUpdate(path dbus.ObjectPath) ***REMOVED***
	c.subscriber.Lock()
	defer c.subscriber.Unlock()

	if c.shouldIgnore(path) ***REMOVED***
		return
	***REMOVED***

	info, err := c.GetUnitProperties(string(path))
	if err != nil ***REMOVED***
		select ***REMOVED***
		case c.subscriber.errCh <- err:
		default:
		***REMOVED***
	***REMOVED***

	name := info["Id"].(string)
	substate := info["SubState"].(string)

	update := &SubStateUpdate***REMOVED***name, substate***REMOVED***
	select ***REMOVED***
	case c.subscriber.updateCh <- update:
	default:
		select ***REMOVED***
		case c.subscriber.errCh <- errors.New("update channel full!"):
		default:
		***REMOVED***
	***REMOVED***

	c.updateIgnore(path, info)
***REMOVED***

// The ignore functions work around a wart in the systemd dbus interface.
// Requesting the properties of an unloaded unit will cause systemd to send a
// pair of UnitNew/UnitRemoved signals.  Because we need to get a unit's
// properties on UnitNew (as that's the only indication of a new unit coming up
// for the first time), we would enter an infinite loop if we did not attempt
// to detect and ignore these spurious signals.  The signal themselves are
// indistinguishable from relevant ones, so we (somewhat hackishly) ignore an
// unloaded unit's signals for a short time after requesting its properties.
// This means that we will miss e.g. a transient unit being restarted
// *immediately* upon failure and also a transient unit being started
// immediately after requesting its status (with systemctl status, for example,
// because this causes a UnitNew signal to be sent which then causes us to fetch
// the properties).

func (c *Conn) shouldIgnore(path dbus.ObjectPath) bool ***REMOVED***
	t, ok := c.subscriber.ignore[path]
	return ok && t >= time.Now().UnixNano()
***REMOVED***

func (c *Conn) updateIgnore(path dbus.ObjectPath, info map[string]interface***REMOVED******REMOVED***) ***REMOVED***
	c.cleanIgnore()

	// unit is unloaded - it will trigger bad systemd dbus behavior
	if info["LoadState"].(string) == "not-found" ***REMOVED***
		c.subscriber.ignore[path] = time.Now().UnixNano() + ignoreInterval
	***REMOVED***
***REMOVED***

// without this, ignore would grow unboundedly over time
func (c *Conn) cleanIgnore() ***REMOVED***
	now := time.Now().UnixNano()
	if c.subscriber.cleanIgnore < now ***REMOVED***
		c.subscriber.cleanIgnore = now + cleanIgnoreInterval

		for p, t := range c.subscriber.ignore ***REMOVED***
			if t < now ***REMOVED***
				delete(c.subscriber.ignore, p)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
