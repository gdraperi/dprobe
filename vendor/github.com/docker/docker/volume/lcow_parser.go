package volume

import (
	"errors"
	"path"

	"github.com/docker/docker/api/types/mount"
)

var lcowSpecificValidators mountValidator = func(m *mount.Mount) error ***REMOVED***
	if path.Clean(m.Target) == "/" ***REMOVED***
		return ErrVolumeTargetIsRoot
	***REMOVED***
	if m.Type == mount.TypeNamedPipe ***REMOVED***
		return errors.New("Linux containers on Windows do not support named pipe mounts")
	***REMOVED***
	return nil
***REMOVED***

type lcowParser struct ***REMOVED***
	windowsParser
***REMOVED***

func (p *lcowParser) ValidateMountConfig(mnt *mount.Mount) error ***REMOVED***
	return p.validateMountConfigReg(mnt, rxLCOWDestination, lcowSpecificValidators)
***REMOVED***

func (p *lcowParser) ParseMountRaw(raw, volumeDriver string) (*MountPoint, error) ***REMOVED***
	return p.parseMountRaw(raw, volumeDriver, rxLCOWDestination, false, lcowSpecificValidators)
***REMOVED***

func (p *lcowParser) ParseMountSpec(cfg mount.Mount) (*MountPoint, error) ***REMOVED***
	return p.parseMountSpec(cfg, rxLCOWDestination, false, lcowSpecificValidators)
***REMOVED***
