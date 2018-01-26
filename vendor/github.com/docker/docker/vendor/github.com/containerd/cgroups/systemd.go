package cgroups

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	systemdDbus "github.com/coreos/go-systemd/dbus"
	"github.com/godbus/dbus"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

const (
	SystemdDbus  Name = "systemd"
	defaultSlice      = "system.slice"
)

func Systemd() ([]Subsystem, error) ***REMOVED***
	root, err := v1MountPoint()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defaultSubsystems, err := defaults(root)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	s, err := NewSystemd(root)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// make sure the systemd controller is added first
	return append([]Subsystem***REMOVED***s***REMOVED***, defaultSubsystems...), nil
***REMOVED***

func Slice(slice, name string) Path ***REMOVED***
	if slice == "" ***REMOVED***
		slice = defaultSlice
	***REMOVED***
	return func(subsystem Name) (string, error) ***REMOVED***
		return filepath.Join(slice, unitName(name)), nil
	***REMOVED***
***REMOVED***

func NewSystemd(root string) (*SystemdController, error) ***REMOVED***
	return &SystemdController***REMOVED***
		root: root,
	***REMOVED***, nil
***REMOVED***

type SystemdController struct ***REMOVED***
	mu   sync.Mutex
	root string
***REMOVED***

func (s *SystemdController) Name() Name ***REMOVED***
	return SystemdDbus
***REMOVED***

func (s *SystemdController) Create(path string, resources *specs.LinuxResources) error ***REMOVED***
	conn, err := systemdDbus.New()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer conn.Close()
	slice, name := splitName(path)
	properties := []systemdDbus.Property***REMOVED***
		systemdDbus.PropDescription(fmt.Sprintf("cgroup %s", name)),
		systemdDbus.PropWants(slice),
		newProperty("DefaultDependencies", false),
		newProperty("Delegate", true),
		newProperty("MemoryAccounting", true),
		newProperty("CPUAccounting", true),
		newProperty("BlockIOAccounting", true),
	***REMOVED***
	ch := make(chan string)
	_, err = conn.StartTransientUnit(name, "replace", properties, ch)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	<-ch
	return nil
***REMOVED***

func (s *SystemdController) Delete(path string) error ***REMOVED***
	conn, err := systemdDbus.New()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer conn.Close()
	_, name := splitName(path)
	ch := make(chan string)
	_, err = conn.StopUnit(name, "replace", ch)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	<-ch
	return nil
***REMOVED***

func newProperty(name string, units interface***REMOVED******REMOVED***) systemdDbus.Property ***REMOVED***
	return systemdDbus.Property***REMOVED***
		Name:  name,
		Value: dbus.MakeVariant(units),
	***REMOVED***
***REMOVED***

func unitName(name string) string ***REMOVED***
	return fmt.Sprintf("%s.slice", name)
***REMOVED***

func splitName(path string) (slice string, unit string) ***REMOVED***
	slice, unit = filepath.Split(path)
	return strings.TrimSuffix(slice, "/"), unit
***REMOVED***
