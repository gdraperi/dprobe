package cgroups

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	specs "github.com/opencontainers/runtime-spec/specs-go"
)

const (
	allowDeviceFile = "devices.allow"
	denyDeviceFile  = "devices.deny"
	wildcard        = -1
)

func NewDevices(root string) *devicesController ***REMOVED***
	return &devicesController***REMOVED***
		root: filepath.Join(root, string(Devices)),
	***REMOVED***
***REMOVED***

type devicesController struct ***REMOVED***
	root string
***REMOVED***

func (d *devicesController) Name() Name ***REMOVED***
	return Devices
***REMOVED***

func (d *devicesController) Path(path string) string ***REMOVED***
	return filepath.Join(d.root, path)
***REMOVED***

func (d *devicesController) Create(path string, resources *specs.LinuxResources) error ***REMOVED***
	if err := os.MkdirAll(d.Path(path), defaultDirPerm); err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, device := range resources.Devices ***REMOVED***
		file := denyDeviceFile
		if device.Allow ***REMOVED***
			file = allowDeviceFile
		***REMOVED***
		if err := ioutil.WriteFile(
			filepath.Join(d.Path(path), file),
			[]byte(deviceString(device)),
			defaultFilePerm,
		); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (d *devicesController) Update(path string, resources *specs.LinuxResources) error ***REMOVED***
	return d.Create(path, resources)
***REMOVED***

func deviceString(device specs.LinuxDeviceCgroup) string ***REMOVED***
	return fmt.Sprintf("%c %s:%s %s",
		&device.Type,
		deviceNumber(device.Major),
		deviceNumber(device.Minor),
		&device.Access,
	)
***REMOVED***

func deviceNumber(number *int64) string ***REMOVED***
	if number == nil || *number == wildcard ***REMOVED***
		return "*"
	***REMOVED***
	return fmt.Sprint(*number)
***REMOVED***
