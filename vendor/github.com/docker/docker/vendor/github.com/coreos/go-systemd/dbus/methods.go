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
	"path"
	"strconv"

	"github.com/godbus/dbus"
)

func (c *Conn) jobComplete(signal *dbus.Signal) ***REMOVED***
	var id uint32
	var job dbus.ObjectPath
	var unit string
	var result string
	dbus.Store(signal.Body, &id, &job, &unit, &result)
	c.jobListener.Lock()
	out, ok := c.jobListener.jobs[job]
	if ok ***REMOVED***
		out <- result
		delete(c.jobListener.jobs, job)
	***REMOVED***
	c.jobListener.Unlock()
***REMOVED***

func (c *Conn) startJob(ch chan<- string, job string, args ...interface***REMOVED******REMOVED***) (int, error) ***REMOVED***
	if ch != nil ***REMOVED***
		c.jobListener.Lock()
		defer c.jobListener.Unlock()
	***REMOVED***

	var p dbus.ObjectPath
	err := c.sysobj.Call(job, 0, args...).Store(&p)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	if ch != nil ***REMOVED***
		c.jobListener.jobs[p] = ch
	***REMOVED***

	// ignore error since 0 is fine if conversion fails
	jobID, _ := strconv.Atoi(path.Base(string(p)))

	return jobID, nil
***REMOVED***

// StartUnit enqueues a start job and depending jobs, if any (unless otherwise
// specified by the mode string).
//
// Takes the unit to activate, plus a mode string. The mode needs to be one of
// replace, fail, isolate, ignore-dependencies, ignore-requirements. If
// "replace" the call will start the unit and its dependencies, possibly
// replacing already queued jobs that conflict with this. If "fail" the call
// will start the unit and its dependencies, but will fail if this would change
// an already queued job. If "isolate" the call will start the unit in question
// and terminate all units that aren't dependencies of it. If
// "ignore-dependencies" it will start a unit but ignore all its dependencies.
// If "ignore-requirements" it will start a unit but only ignore the
// requirement dependencies. It is not recommended to make use of the latter
// two options.
//
// If the provided channel is non-nil, a result string will be sent to it upon
// job completion: one of done, canceled, timeout, failed, dependency, skipped.
// done indicates successful execution of a job. canceled indicates that a job
// has been canceled  before it finished execution. timeout indicates that the
// job timeout was reached. failed indicates that the job failed. dependency
// indicates that a job this job has been depending on failed and the job hence
// has been removed too. skipped indicates that a job was skipped because it
// didn't apply to the units current state.
//
// If no error occurs, the ID of the underlying systemd job will be returned. There
// does exist the possibility for no error to be returned, but for the returned job
// ID to be 0. In this case, the actual underlying ID is not 0 and this datapoint
// should not be considered authoritative.
//
// If an error does occur, it will be returned to the user alongside a job ID of 0.
func (c *Conn) StartUnit(name string, mode string, ch chan<- string) (int, error) ***REMOVED***
	return c.startJob(ch, "org.freedesktop.systemd1.Manager.StartUnit", name, mode)
***REMOVED***

// StopUnit is similar to StartUnit but stops the specified unit rather
// than starting it.
func (c *Conn) StopUnit(name string, mode string, ch chan<- string) (int, error) ***REMOVED***
	return c.startJob(ch, "org.freedesktop.systemd1.Manager.StopUnit", name, mode)
***REMOVED***

// ReloadUnit reloads a unit.  Reloading is done only if the unit is already running and fails otherwise.
func (c *Conn) ReloadUnit(name string, mode string, ch chan<- string) (int, error) ***REMOVED***
	return c.startJob(ch, "org.freedesktop.systemd1.Manager.ReloadUnit", name, mode)
***REMOVED***

// RestartUnit restarts a service.  If a service is restarted that isn't
// running it will be started.
func (c *Conn) RestartUnit(name string, mode string, ch chan<- string) (int, error) ***REMOVED***
	return c.startJob(ch, "org.freedesktop.systemd1.Manager.RestartUnit", name, mode)
***REMOVED***

// TryRestartUnit is like RestartUnit, except that a service that isn't running
// is not affected by the restart.
func (c *Conn) TryRestartUnit(name string, mode string, ch chan<- string) (int, error) ***REMOVED***
	return c.startJob(ch, "org.freedesktop.systemd1.Manager.TryRestartUnit", name, mode)
***REMOVED***

// ReloadOrRestart attempts a reload if the unit supports it and use a restart
// otherwise.
func (c *Conn) ReloadOrRestartUnit(name string, mode string, ch chan<- string) (int, error) ***REMOVED***
	return c.startJob(ch, "org.freedesktop.systemd1.Manager.ReloadOrRestartUnit", name, mode)
***REMOVED***

// ReloadOrTryRestart attempts a reload if the unit supports it and use a "Try"
// flavored restart otherwise.
func (c *Conn) ReloadOrTryRestartUnit(name string, mode string, ch chan<- string) (int, error) ***REMOVED***
	return c.startJob(ch, "org.freedesktop.systemd1.Manager.ReloadOrTryRestartUnit", name, mode)
***REMOVED***

// StartTransientUnit() may be used to create and start a transient unit, which
// will be released as soon as it is not running or referenced anymore or the
// system is rebooted. name is the unit name including suffix, and must be
// unique. mode is the same as in StartUnit(), properties contains properties
// of the unit.
func (c *Conn) StartTransientUnit(name string, mode string, properties []Property, ch chan<- string) (int, error) ***REMOVED***
	return c.startJob(ch, "org.freedesktop.systemd1.Manager.StartTransientUnit", name, mode, properties, make([]PropertyCollection, 0))
***REMOVED***

// KillUnit takes the unit name and a UNIX signal number to send.  All of the unit's
// processes are killed.
func (c *Conn) KillUnit(name string, signal int32) ***REMOVED***
	c.sysobj.Call("org.freedesktop.systemd1.Manager.KillUnit", 0, name, "all", signal).Store()
***REMOVED***

// ResetFailedUnit resets the "failed" state of a specific unit.
func (c *Conn) ResetFailedUnit(name string) error ***REMOVED***
	return c.sysobj.Call("org.freedesktop.systemd1.Manager.ResetFailedUnit", 0, name).Store()
***REMOVED***

// getProperties takes the unit name and returns all of its dbus object properties, for the given dbus interface
func (c *Conn) getProperties(unit string, dbusInterface string) (map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	var err error
	var props map[string]dbus.Variant

	path := unitPath(unit)
	if !path.IsValid() ***REMOVED***
		return nil, errors.New("invalid unit name: " + unit)
	***REMOVED***

	obj := c.sysconn.Object("org.freedesktop.systemd1", path)
	err = obj.Call("org.freedesktop.DBus.Properties.GetAll", 0, dbusInterface).Store(&props)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	out := make(map[string]interface***REMOVED******REMOVED***, len(props))
	for k, v := range props ***REMOVED***
		out[k] = v.Value()
	***REMOVED***

	return out, nil
***REMOVED***

// GetUnitProperties takes the unit name and returns all of its dbus object properties.
func (c *Conn) GetUnitProperties(unit string) (map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	return c.getProperties(unit, "org.freedesktop.systemd1.Unit")
***REMOVED***

func (c *Conn) getProperty(unit string, dbusInterface string, propertyName string) (*Property, error) ***REMOVED***
	var err error
	var prop dbus.Variant

	path := unitPath(unit)
	if !path.IsValid() ***REMOVED***
		return nil, errors.New("invalid unit name: " + unit)
	***REMOVED***

	obj := c.sysconn.Object("org.freedesktop.systemd1", path)
	err = obj.Call("org.freedesktop.DBus.Properties.Get", 0, dbusInterface, propertyName).Store(&prop)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &Property***REMOVED***Name: propertyName, Value: prop***REMOVED***, nil
***REMOVED***

func (c *Conn) GetUnitProperty(unit string, propertyName string) (*Property, error) ***REMOVED***
	return c.getProperty(unit, "org.freedesktop.systemd1.Unit", propertyName)
***REMOVED***

// GetServiceProperty returns property for given service name and property name
func (c *Conn) GetServiceProperty(service string, propertyName string) (*Property, error) ***REMOVED***
	return c.getProperty(service, "org.freedesktop.systemd1.Service", propertyName)
***REMOVED***

// GetUnitTypeProperties returns the extra properties for a unit, specific to the unit type.
// Valid values for unitType: Service, Socket, Target, Device, Mount, Automount, Snapshot, Timer, Swap, Path, Slice, Scope
// return "dbus.Error: Unknown interface" if the unitType is not the correct type of the unit
func (c *Conn) GetUnitTypeProperties(unit string, unitType string) (map[string]interface***REMOVED******REMOVED***, error) ***REMOVED***
	return c.getProperties(unit, "org.freedesktop.systemd1."+unitType)
***REMOVED***

// SetUnitProperties() may be used to modify certain unit properties at runtime.
// Not all properties may be changed at runtime, but many resource management
// settings (primarily those in systemd.cgroup(5)) may. The changes are applied
// instantly, and stored on disk for future boots, unless runtime is true, in which
// case the settings only apply until the next reboot. name is the name of the unit
// to modify. properties are the settings to set, encoded as an array of property
// name and value pairs.
func (c *Conn) SetUnitProperties(name string, runtime bool, properties ...Property) error ***REMOVED***
	return c.sysobj.Call("org.freedesktop.systemd1.Manager.SetUnitProperties", 0, name, runtime, properties).Store()
***REMOVED***

func (c *Conn) GetUnitTypeProperty(unit string, unitType string, propertyName string) (*Property, error) ***REMOVED***
	return c.getProperty(unit, "org.freedesktop.systemd1."+unitType, propertyName)
***REMOVED***

type UnitStatus struct ***REMOVED***
	Name        string          // The primary unit name as string
	Description string          // The human readable description string
	LoadState   string          // The load state (i.e. whether the unit file has been loaded successfully)
	ActiveState string          // The active state (i.e. whether the unit is currently started or not)
	SubState    string          // The sub state (a more fine-grained version of the active state that is specific to the unit type, which the active state is not)
	Followed    string          // A unit that is being followed in its state by this unit, if there is any, otherwise the empty string.
	Path        dbus.ObjectPath // The unit object path
	JobId       uint32          // If there is a job queued for the job unit the numeric job id, 0 otherwise
	JobType     string          // The job type as string
	JobPath     dbus.ObjectPath // The job object path
***REMOVED***

type storeFunc func(retvalues ...interface***REMOVED******REMOVED***) error

func (c *Conn) listUnitsInternal(f storeFunc) ([]UnitStatus, error) ***REMOVED***
	result := make([][]interface***REMOVED******REMOVED***, 0)
	err := f(&result)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	resultInterface := make([]interface***REMOVED******REMOVED***, len(result))
	for i := range result ***REMOVED***
		resultInterface[i] = result[i]
	***REMOVED***

	status := make([]UnitStatus, len(result))
	statusInterface := make([]interface***REMOVED******REMOVED***, len(status))
	for i := range status ***REMOVED***
		statusInterface[i] = &status[i]
	***REMOVED***

	err = dbus.Store(resultInterface, statusInterface...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return status, nil
***REMOVED***

// ListUnits returns an array with all currently loaded units. Note that
// units may be known by multiple names at the same time, and hence there might
// be more unit names loaded than actual units behind them.
func (c *Conn) ListUnits() ([]UnitStatus, error) ***REMOVED***
	return c.listUnitsInternal(c.sysobj.Call("org.freedesktop.systemd1.Manager.ListUnits", 0).Store)
***REMOVED***

// ListUnitsFiltered returns an array with units filtered by state.
// It takes a list of units' statuses to filter.
func (c *Conn) ListUnitsFiltered(states []string) ([]UnitStatus, error) ***REMOVED***
	return c.listUnitsInternal(c.sysobj.Call("org.freedesktop.systemd1.Manager.ListUnitsFiltered", 0, states).Store)
***REMOVED***

// ListUnitsByPatterns returns an array with units.
// It takes a list of units' statuses and names to filter.
// Note that units may be known by multiple names at the same time,
// and hence there might be more unit names loaded than actual units behind them.
func (c *Conn) ListUnitsByPatterns(states []string, patterns []string) ([]UnitStatus, error) ***REMOVED***
	return c.listUnitsInternal(c.sysobj.Call("org.freedesktop.systemd1.Manager.ListUnitsByPatterns", 0, states, patterns).Store)
***REMOVED***

// ListUnitsByNames returns an array with units. It takes a list of units'
// names and returns an UnitStatus array. Comparing to ListUnitsByPatterns
// method, this method returns statuses even for inactive or non-existing
// units. Input array should contain exact unit names, but not patterns.
func (c *Conn) ListUnitsByNames(units []string) ([]UnitStatus, error) ***REMOVED***
	return c.listUnitsInternal(c.sysobj.Call("org.freedesktop.systemd1.Manager.ListUnitsByNames", 0, units).Store)
***REMOVED***

type UnitFile struct ***REMOVED***
	Path string
	Type string
***REMOVED***

func (c *Conn) listUnitFilesInternal(f storeFunc) ([]UnitFile, error) ***REMOVED***
	result := make([][]interface***REMOVED******REMOVED***, 0)
	err := f(&result)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	resultInterface := make([]interface***REMOVED******REMOVED***, len(result))
	for i := range result ***REMOVED***
		resultInterface[i] = result[i]
	***REMOVED***

	files := make([]UnitFile, len(result))
	fileInterface := make([]interface***REMOVED******REMOVED***, len(files))
	for i := range files ***REMOVED***
		fileInterface[i] = &files[i]
	***REMOVED***

	err = dbus.Store(resultInterface, fileInterface...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return files, nil
***REMOVED***

// ListUnitFiles returns an array of all available units on disk.
func (c *Conn) ListUnitFiles() ([]UnitFile, error) ***REMOVED***
	return c.listUnitFilesInternal(c.sysobj.Call("org.freedesktop.systemd1.Manager.ListUnitFiles", 0).Store)
***REMOVED***

// ListUnitFilesByPatterns returns an array of all available units on disk matched the patterns.
func (c *Conn) ListUnitFilesByPatterns(states []string, patterns []string) ([]UnitFile, error) ***REMOVED***
	return c.listUnitFilesInternal(c.sysobj.Call("org.freedesktop.systemd1.Manager.ListUnitFilesByPatterns", 0, states, patterns).Store)
***REMOVED***

type LinkUnitFileChange EnableUnitFileChange

// LinkUnitFiles() links unit files (that are located outside of the
// usual unit search paths) into the unit search path.
//
// It takes a list of absolute paths to unit files to link and two
// booleans. The first boolean controls whether the unit shall be
// enabled for runtime only (true, /run), or persistently (false,
// /etc).
// The second controls whether symlinks pointing to other units shall
// be replaced if necessary.
//
// This call returns a list of the changes made. The list consists of
// structures with three strings: the type of the change (one of symlink
// or unlink), the file name of the symlink and the destination of the
// symlink.
func (c *Conn) LinkUnitFiles(files []string, runtime bool, force bool) ([]LinkUnitFileChange, error) ***REMOVED***
	result := make([][]interface***REMOVED******REMOVED***, 0)
	err := c.sysobj.Call("org.freedesktop.systemd1.Manager.LinkUnitFiles", 0, files, runtime, force).Store(&result)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	resultInterface := make([]interface***REMOVED******REMOVED***, len(result))
	for i := range result ***REMOVED***
		resultInterface[i] = result[i]
	***REMOVED***

	changes := make([]LinkUnitFileChange, len(result))
	changesInterface := make([]interface***REMOVED******REMOVED***, len(changes))
	for i := range changes ***REMOVED***
		changesInterface[i] = &changes[i]
	***REMOVED***

	err = dbus.Store(resultInterface, changesInterface...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return changes, nil
***REMOVED***

// EnableUnitFiles() may be used to enable one or more units in the system (by
// creating symlinks to them in /etc or /run).
//
// It takes a list of unit files to enable (either just file names or full
// absolute paths if the unit files are residing outside the usual unit
// search paths), and two booleans: the first controls whether the unit shall
// be enabled for runtime only (true, /run), or persistently (false, /etc).
// The second one controls whether symlinks pointing to other units shall
// be replaced if necessary.
//
// This call returns one boolean and an array with the changes made. The
// boolean signals whether the unit files contained any enablement
// information (i.e. an [Install]) section. The changes list consists of
// structures with three strings: the type of the change (one of symlink
// or unlink), the file name of the symlink and the destination of the
// symlink.
func (c *Conn) EnableUnitFiles(files []string, runtime bool, force bool) (bool, []EnableUnitFileChange, error) ***REMOVED***
	var carries_install_info bool

	result := make([][]interface***REMOVED******REMOVED***, 0)
	err := c.sysobj.Call("org.freedesktop.systemd1.Manager.EnableUnitFiles", 0, files, runtime, force).Store(&carries_install_info, &result)
	if err != nil ***REMOVED***
		return false, nil, err
	***REMOVED***

	resultInterface := make([]interface***REMOVED******REMOVED***, len(result))
	for i := range result ***REMOVED***
		resultInterface[i] = result[i]
	***REMOVED***

	changes := make([]EnableUnitFileChange, len(result))
	changesInterface := make([]interface***REMOVED******REMOVED***, len(changes))
	for i := range changes ***REMOVED***
		changesInterface[i] = &changes[i]
	***REMOVED***

	err = dbus.Store(resultInterface, changesInterface...)
	if err != nil ***REMOVED***
		return false, nil, err
	***REMOVED***

	return carries_install_info, changes, nil
***REMOVED***

type EnableUnitFileChange struct ***REMOVED***
	Type        string // Type of the change (one of symlink or unlink)
	Filename    string // File name of the symlink
	Destination string // Destination of the symlink
***REMOVED***

// DisableUnitFiles() may be used to disable one or more units in the system (by
// removing symlinks to them from /etc or /run).
//
// It takes a list of unit files to disable (either just file names or full
// absolute paths if the unit files are residing outside the usual unit
// search paths), and one boolean: whether the unit was enabled for runtime
// only (true, /run), or persistently (false, /etc).
//
// This call returns an array with the changes made. The changes list
// consists of structures with three strings: the type of the change (one of
// symlink or unlink), the file name of the symlink and the destination of the
// symlink.
func (c *Conn) DisableUnitFiles(files []string, runtime bool) ([]DisableUnitFileChange, error) ***REMOVED***
	result := make([][]interface***REMOVED******REMOVED***, 0)
	err := c.sysobj.Call("org.freedesktop.systemd1.Manager.DisableUnitFiles", 0, files, runtime).Store(&result)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	resultInterface := make([]interface***REMOVED******REMOVED***, len(result))
	for i := range result ***REMOVED***
		resultInterface[i] = result[i]
	***REMOVED***

	changes := make([]DisableUnitFileChange, len(result))
	changesInterface := make([]interface***REMOVED******REMOVED***, len(changes))
	for i := range changes ***REMOVED***
		changesInterface[i] = &changes[i]
	***REMOVED***

	err = dbus.Store(resultInterface, changesInterface...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return changes, nil
***REMOVED***

type DisableUnitFileChange struct ***REMOVED***
	Type        string // Type of the change (one of symlink or unlink)
	Filename    string // File name of the symlink
	Destination string // Destination of the symlink
***REMOVED***

// MaskUnitFiles masks one or more units in the system
//
// It takes three arguments:
//   * list of units to mask (either just file names or full
//     absolute paths if the unit files are residing outside
//     the usual unit search paths)
//   * runtime to specify whether the unit was enabled for runtime
//     only (true, /run/systemd/..), or persistently (false, /etc/systemd/..)
//   * force flag
func (c *Conn) MaskUnitFiles(files []string, runtime bool, force bool) ([]MaskUnitFileChange, error) ***REMOVED***
	result := make([][]interface***REMOVED******REMOVED***, 0)
	err := c.sysobj.Call("org.freedesktop.systemd1.Manager.MaskUnitFiles", 0, files, runtime, force).Store(&result)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	resultInterface := make([]interface***REMOVED******REMOVED***, len(result))
	for i := range result ***REMOVED***
		resultInterface[i] = result[i]
	***REMOVED***

	changes := make([]MaskUnitFileChange, len(result))
	changesInterface := make([]interface***REMOVED******REMOVED***, len(changes))
	for i := range changes ***REMOVED***
		changesInterface[i] = &changes[i]
	***REMOVED***

	err = dbus.Store(resultInterface, changesInterface...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return changes, nil
***REMOVED***

type MaskUnitFileChange struct ***REMOVED***
	Type        string // Type of the change (one of symlink or unlink)
	Filename    string // File name of the symlink
	Destination string // Destination of the symlink
***REMOVED***

// UnmaskUnitFiles unmasks one or more units in the system
//
// It takes two arguments:
//   * list of unit files to mask (either just file names or full
//     absolute paths if the unit files are residing outside
//     the usual unit search paths)
//   * runtime to specify whether the unit was enabled for runtime
//     only (true, /run/systemd/..), or persistently (false, /etc/systemd/..)
func (c *Conn) UnmaskUnitFiles(files []string, runtime bool) ([]UnmaskUnitFileChange, error) ***REMOVED***
	result := make([][]interface***REMOVED******REMOVED***, 0)
	err := c.sysobj.Call("org.freedesktop.systemd1.Manager.UnmaskUnitFiles", 0, files, runtime).Store(&result)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	resultInterface := make([]interface***REMOVED******REMOVED***, len(result))
	for i := range result ***REMOVED***
		resultInterface[i] = result[i]
	***REMOVED***

	changes := make([]UnmaskUnitFileChange, len(result))
	changesInterface := make([]interface***REMOVED******REMOVED***, len(changes))
	for i := range changes ***REMOVED***
		changesInterface[i] = &changes[i]
	***REMOVED***

	err = dbus.Store(resultInterface, changesInterface...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return changes, nil
***REMOVED***

type UnmaskUnitFileChange struct ***REMOVED***
	Type        string // Type of the change (one of symlink or unlink)
	Filename    string // File name of the symlink
	Destination string // Destination of the symlink
***REMOVED***

// Reload instructs systemd to scan for and reload unit files. This is
// equivalent to a 'systemctl daemon-reload'.
func (c *Conn) Reload() error ***REMOVED***
	return c.sysobj.Call("org.freedesktop.systemd1.Manager.Reload", 0).Store()
***REMOVED***

func unitPath(name string) dbus.ObjectPath ***REMOVED***
	return dbus.ObjectPath("/org/freedesktop/systemd1/unit/" + PathBusEscape(name))
***REMOVED***
