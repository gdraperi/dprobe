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
	"github.com/godbus/dbus"
)

// From the systemd docs:
//
// The properties array of StartTransientUnit() may take many of the settings
// that may also be configured in unit files. Not all parameters are currently
// accepted though, but we plan to cover more properties with future release.
// Currently you may set the Description, Slice and all dependency types of
// units, as well as RemainAfterExit, ExecStart for service units,
// TimeoutStopUSec and PIDs for scope units, and CPUAccounting, CPUShares,
// BlockIOAccounting, BlockIOWeight, BlockIOReadBandwidth,
// BlockIOWriteBandwidth, BlockIODeviceWeight, MemoryAccounting, MemoryLimit,
// DevicePolicy, DeviceAllow for services/scopes/slices. These fields map
// directly to their counterparts in unit files and as normal D-Bus object
// properties. The exception here is the PIDs field of scope units which is
// used for construction of the scope only and specifies the initial PIDs to
// add to the scope object.

type Property struct ***REMOVED***
	Name  string
	Value dbus.Variant
***REMOVED***

type PropertyCollection struct ***REMOVED***
	Name       string
	Properties []Property
***REMOVED***

type execStart struct ***REMOVED***
	Path             string   // the binary path to execute
	Args             []string // an array with all arguments to pass to the executed command, starting with argument 0
	UncleanIsFailure bool     // a boolean whether it should be considered a failure if the process exits uncleanly
***REMOVED***

// PropExecStart sets the ExecStart service property.  The first argument is a
// slice with the binary path to execute followed by the arguments to pass to
// the executed command. See
// http://www.freedesktop.org/software/systemd/man/systemd.service.html#ExecStart=
func PropExecStart(command []string, uncleanIsFailure bool) Property ***REMOVED***
	execStarts := []execStart***REMOVED***
		execStart***REMOVED***
			Path:             command[0],
			Args:             command,
			UncleanIsFailure: uncleanIsFailure,
		***REMOVED***,
	***REMOVED***

	return Property***REMOVED***
		Name:  "ExecStart",
		Value: dbus.MakeVariant(execStarts),
	***REMOVED***
***REMOVED***

// PropRemainAfterExit sets the RemainAfterExit service property. See
// http://www.freedesktop.org/software/systemd/man/systemd.service.html#RemainAfterExit=
func PropRemainAfterExit(b bool) Property ***REMOVED***
	return Property***REMOVED***
		Name:  "RemainAfterExit",
		Value: dbus.MakeVariant(b),
	***REMOVED***
***REMOVED***

// PropType sets the Type service property. See
// http://www.freedesktop.org/software/systemd/man/systemd.service.html#Type=
func PropType(t string) Property ***REMOVED***
	return Property***REMOVED***
		Name:  "Type",
		Value: dbus.MakeVariant(t),
	***REMOVED***
***REMOVED***

// PropDescription sets the Description unit property. See
// http://www.freedesktop.org/software/systemd/man/systemd.unit#Description=
func PropDescription(desc string) Property ***REMOVED***
	return Property***REMOVED***
		Name:  "Description",
		Value: dbus.MakeVariant(desc),
	***REMOVED***
***REMOVED***

func propDependency(name string, units []string) Property ***REMOVED***
	return Property***REMOVED***
		Name:  name,
		Value: dbus.MakeVariant(units),
	***REMOVED***
***REMOVED***

// PropRequires sets the Requires unit property.  See
// http://www.freedesktop.org/software/systemd/man/systemd.unit.html#Requires=
func PropRequires(units ...string) Property ***REMOVED***
	return propDependency("Requires", units)
***REMOVED***

// PropRequiresOverridable sets the RequiresOverridable unit property.  See
// http://www.freedesktop.org/software/systemd/man/systemd.unit.html#RequiresOverridable=
func PropRequiresOverridable(units ...string) Property ***REMOVED***
	return propDependency("RequiresOverridable", units)
***REMOVED***

// PropRequisite sets the Requisite unit property.  See
// http://www.freedesktop.org/software/systemd/man/systemd.unit.html#Requisite=
func PropRequisite(units ...string) Property ***REMOVED***
	return propDependency("Requisite", units)
***REMOVED***

// PropRequisiteOverridable sets the RequisiteOverridable unit property.  See
// http://www.freedesktop.org/software/systemd/man/systemd.unit.html#RequisiteOverridable=
func PropRequisiteOverridable(units ...string) Property ***REMOVED***
	return propDependency("RequisiteOverridable", units)
***REMOVED***

// PropWants sets the Wants unit property.  See
// http://www.freedesktop.org/software/systemd/man/systemd.unit.html#Wants=
func PropWants(units ...string) Property ***REMOVED***
	return propDependency("Wants", units)
***REMOVED***

// PropBindsTo sets the BindsTo unit property.  See
// http://www.freedesktop.org/software/systemd/man/systemd.unit.html#BindsTo=
func PropBindsTo(units ...string) Property ***REMOVED***
	return propDependency("BindsTo", units)
***REMOVED***

// PropRequiredBy sets the RequiredBy unit property.  See
// http://www.freedesktop.org/software/systemd/man/systemd.unit.html#RequiredBy=
func PropRequiredBy(units ...string) Property ***REMOVED***
	return propDependency("RequiredBy", units)
***REMOVED***

// PropRequiredByOverridable sets the RequiredByOverridable unit property.  See
// http://www.freedesktop.org/software/systemd/man/systemd.unit.html#RequiredByOverridable=
func PropRequiredByOverridable(units ...string) Property ***REMOVED***
	return propDependency("RequiredByOverridable", units)
***REMOVED***

// PropWantedBy sets the WantedBy unit property.  See
// http://www.freedesktop.org/software/systemd/man/systemd.unit.html#WantedBy=
func PropWantedBy(units ...string) Property ***REMOVED***
	return propDependency("WantedBy", units)
***REMOVED***

// PropBoundBy sets the BoundBy unit property.  See
// http://www.freedesktop.org/software/systemd/main/systemd.unit.html#BoundBy=
func PropBoundBy(units ...string) Property ***REMOVED***
	return propDependency("BoundBy", units)
***REMOVED***

// PropConflicts sets the Conflicts unit property.  See
// http://www.freedesktop.org/software/systemd/man/systemd.unit.html#Conflicts=
func PropConflicts(units ...string) Property ***REMOVED***
	return propDependency("Conflicts", units)
***REMOVED***

// PropConflictedBy sets the ConflictedBy unit property.  See
// http://www.freedesktop.org/software/systemd/man/systemd.unit.html#ConflictedBy=
func PropConflictedBy(units ...string) Property ***REMOVED***
	return propDependency("ConflictedBy", units)
***REMOVED***

// PropBefore sets the Before unit property.  See
// http://www.freedesktop.org/software/systemd/man/systemd.unit.html#Before=
func PropBefore(units ...string) Property ***REMOVED***
	return propDependency("Before", units)
***REMOVED***

// PropAfter sets the After unit property.  See
// http://www.freedesktop.org/software/systemd/man/systemd.unit.html#After=
func PropAfter(units ...string) Property ***REMOVED***
	return propDependency("After", units)
***REMOVED***

// PropOnFailure sets the OnFailure unit property.  See
// http://www.freedesktop.org/software/systemd/man/systemd.unit.html#OnFailure=
func PropOnFailure(units ...string) Property ***REMOVED***
	return propDependency("OnFailure", units)
***REMOVED***

// PropTriggers sets the Triggers unit property.  See
// http://www.freedesktop.org/software/systemd/man/systemd.unit.html#Triggers=
func PropTriggers(units ...string) Property ***REMOVED***
	return propDependency("Triggers", units)
***REMOVED***

// PropTriggeredBy sets the TriggeredBy unit property.  See
// http://www.freedesktop.org/software/systemd/man/systemd.unit.html#TriggeredBy=
func PropTriggeredBy(units ...string) Property ***REMOVED***
	return propDependency("TriggeredBy", units)
***REMOVED***

// PropPropagatesReloadTo sets the PropagatesReloadTo unit property.  See
// http://www.freedesktop.org/software/systemd/man/systemd.unit.html#PropagatesReloadTo=
func PropPropagatesReloadTo(units ...string) Property ***REMOVED***
	return propDependency("PropagatesReloadTo", units)
***REMOVED***

// PropRequiresMountsFor sets the RequiresMountsFor unit property.  See
// http://www.freedesktop.org/software/systemd/man/systemd.unit.html#RequiresMountsFor=
func PropRequiresMountsFor(units ...string) Property ***REMOVED***
	return propDependency("RequiresMountsFor", units)
***REMOVED***

// PropSlice sets the Slice unit property.  See
// http://www.freedesktop.org/software/systemd/man/systemd.resource-control.html#Slice=
func PropSlice(slice string) Property ***REMOVED***
	return Property***REMOVED***
		Name:  "Slice",
		Value: dbus.MakeVariant(slice),
	***REMOVED***
***REMOVED***

// PropPids sets the PIDs field of scope units used in the initial construction
// of the scope only and specifies the initial PIDs to add to the scope object.
// See https://www.freedesktop.org/wiki/Software/systemd/ControlGroupInterface/#properties
func PropPids(pids ...uint32) Property ***REMOVED***
	return Property***REMOVED***
		Name:  "PIDs",
		Value: dbus.MakeVariant(pids),
	***REMOVED***
***REMOVED***
