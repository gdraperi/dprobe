package volume

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/pkg/stringid"
)

type linuxParser struct ***REMOVED***
***REMOVED***

func linuxSplitRawSpec(raw string) ([]string, error) ***REMOVED***
	if strings.Count(raw, ":") > 2 ***REMOVED***
		return nil, errInvalidSpec(raw)
	***REMOVED***

	arr := strings.SplitN(raw, ":", 3)
	if arr[0] == "" ***REMOVED***
		return nil, errInvalidSpec(raw)
	***REMOVED***
	return arr, nil
***REMOVED***

func linuxValidateNotRoot(p string) error ***REMOVED***
	p = path.Clean(strings.Replace(p, `\`, `/`, -1))
	if p == "/" ***REMOVED***
		return ErrVolumeTargetIsRoot
	***REMOVED***
	return nil
***REMOVED***
func linuxValidateAbsolute(p string) error ***REMOVED***
	p = strings.Replace(p, `\`, `/`, -1)
	if path.IsAbs(p) ***REMOVED***
		return nil
	***REMOVED***
	return fmt.Errorf("invalid mount path: '%s' mount path must be absolute", p)
***REMOVED***
func (p *linuxParser) ValidateMountConfig(mnt *mount.Mount) error ***REMOVED***
	// there was something looking like a bug in existing codebase:
	// - validateMountConfig on linux was called with options skipping bind source existence when calling ParseMountRaw
	// - but not when calling ParseMountSpec directly... nor when the unit test called it directly
	return p.validateMountConfigImpl(mnt, true)
***REMOVED***
func (p *linuxParser) validateMountConfigImpl(mnt *mount.Mount, validateBindSourceExists bool) error ***REMOVED***
	if len(mnt.Target) == 0 ***REMOVED***
		return &errMountConfig***REMOVED***mnt, errMissingField("Target")***REMOVED***
	***REMOVED***

	if err := linuxValidateNotRoot(mnt.Target); err != nil ***REMOVED***
		return &errMountConfig***REMOVED***mnt, err***REMOVED***
	***REMOVED***

	if err := linuxValidateAbsolute(mnt.Target); err != nil ***REMOVED***
		return &errMountConfig***REMOVED***mnt, err***REMOVED***
	***REMOVED***

	switch mnt.Type ***REMOVED***
	case mount.TypeBind:
		if len(mnt.Source) == 0 ***REMOVED***
			return &errMountConfig***REMOVED***mnt, errMissingField("Source")***REMOVED***
		***REMOVED***
		// Don't error out just because the propagation mode is not supported on the platform
		if opts := mnt.BindOptions; opts != nil ***REMOVED***
			if len(opts.Propagation) > 0 && len(linuxPropagationModes) > 0 ***REMOVED***
				if _, ok := linuxPropagationModes[opts.Propagation]; !ok ***REMOVED***
					return &errMountConfig***REMOVED***mnt, fmt.Errorf("invalid propagation mode: %s", opts.Propagation)***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if mnt.VolumeOptions != nil ***REMOVED***
			return &errMountConfig***REMOVED***mnt, errExtraField("VolumeOptions")***REMOVED***
		***REMOVED***

		if err := linuxValidateAbsolute(mnt.Source); err != nil ***REMOVED***
			return &errMountConfig***REMOVED***mnt, err***REMOVED***
		***REMOVED***

		if validateBindSourceExists ***REMOVED***
			exists, _, _ := currentFileInfoProvider.fileInfo(mnt.Source)
			if !exists ***REMOVED***
				return &errMountConfig***REMOVED***mnt, errBindNotExist***REMOVED***
			***REMOVED***
		***REMOVED***

	case mount.TypeVolume:
		if mnt.BindOptions != nil ***REMOVED***
			return &errMountConfig***REMOVED***mnt, errExtraField("BindOptions")***REMOVED***
		***REMOVED***

		if len(mnt.Source) == 0 && mnt.ReadOnly ***REMOVED***
			return &errMountConfig***REMOVED***mnt, fmt.Errorf("must not set ReadOnly mode when using anonymous volumes")***REMOVED***
		***REMOVED***
	case mount.TypeTmpfs:
		if len(mnt.Source) != 0 ***REMOVED***
			return &errMountConfig***REMOVED***mnt, errExtraField("Source")***REMOVED***
		***REMOVED***
		if _, err := p.ConvertTmpfsOptions(mnt.TmpfsOptions, mnt.ReadOnly); err != nil ***REMOVED***
			return &errMountConfig***REMOVED***mnt, err***REMOVED***
		***REMOVED***
	default:
		return &errMountConfig***REMOVED***mnt, errors.New("mount type unknown")***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// read-write modes
var rwModes = map[string]bool***REMOVED***
	"rw": true,
	"ro": true,
***REMOVED***

// label modes
var linuxLabelModes = map[string]bool***REMOVED***
	"Z": true,
	"z": true,
***REMOVED***

// consistency modes
var linuxConsistencyModes = map[mount.Consistency]bool***REMOVED***
	mount.ConsistencyFull:      true,
	mount.ConsistencyCached:    true,
	mount.ConsistencyDelegated: true,
***REMOVED***
var linuxPropagationModes = map[mount.Propagation]bool***REMOVED***
	mount.PropagationPrivate:  true,
	mount.PropagationRPrivate: true,
	mount.PropagationSlave:    true,
	mount.PropagationRSlave:   true,
	mount.PropagationShared:   true,
	mount.PropagationRShared:  true,
***REMOVED***

const linuxDefaultPropagationMode = mount.PropagationRPrivate

func linuxGetPropagation(mode string) mount.Propagation ***REMOVED***
	for _, o := range strings.Split(mode, ",") ***REMOVED***
		prop := mount.Propagation(o)
		if linuxPropagationModes[prop] ***REMOVED***
			return prop
		***REMOVED***
	***REMOVED***
	return linuxDefaultPropagationMode
***REMOVED***

func linuxHasPropagation(mode string) bool ***REMOVED***
	for _, o := range strings.Split(mode, ",") ***REMOVED***
		if linuxPropagationModes[mount.Propagation(o)] ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func linuxValidMountMode(mode string) bool ***REMOVED***
	if mode == "" ***REMOVED***
		return true
	***REMOVED***

	rwModeCount := 0
	labelModeCount := 0
	propagationModeCount := 0
	copyModeCount := 0
	consistencyModeCount := 0

	for _, o := range strings.Split(mode, ",") ***REMOVED***
		switch ***REMOVED***
		case rwModes[o]:
			rwModeCount++
		case linuxLabelModes[o]:
			labelModeCount++
		case linuxPropagationModes[mount.Propagation(o)]:
			propagationModeCount++
		case copyModeExists(o):
			copyModeCount++
		case linuxConsistencyModes[mount.Consistency(o)]:
			consistencyModeCount++
		default:
			return false
		***REMOVED***
	***REMOVED***

	// Only one string for each mode is allowed.
	if rwModeCount > 1 || labelModeCount > 1 || propagationModeCount > 1 || copyModeCount > 1 || consistencyModeCount > 1 ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

func (p *linuxParser) ReadWrite(mode string) bool ***REMOVED***
	if !linuxValidMountMode(mode) ***REMOVED***
		return false
	***REMOVED***

	for _, o := range strings.Split(mode, ",") ***REMOVED***
		if o == "ro" ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func (p *linuxParser) ParseMountRaw(raw, volumeDriver string) (*MountPoint, error) ***REMOVED***
	arr, err := linuxSplitRawSpec(raw)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var spec mount.Mount
	var mode string
	switch len(arr) ***REMOVED***
	case 1:
		// Just a destination path in the container
		spec.Target = arr[0]
	case 2:
		if linuxValidMountMode(arr[1]) ***REMOVED***
			// Destination + Mode is not a valid volume - volumes
			// cannot include a mode. e.g. /foo:rw
			return nil, errInvalidSpec(raw)
		***REMOVED***
		// Host Source Path or Name + Destination
		spec.Source = arr[0]
		spec.Target = arr[1]
	case 3:
		// HostSourcePath+DestinationPath+Mode
		spec.Source = arr[0]
		spec.Target = arr[1]
		mode = arr[2]
	default:
		return nil, errInvalidSpec(raw)
	***REMOVED***

	if !linuxValidMountMode(mode) ***REMOVED***
		return nil, errInvalidMode(mode)
	***REMOVED***

	if path.IsAbs(spec.Source) ***REMOVED***
		spec.Type = mount.TypeBind
	***REMOVED*** else ***REMOVED***
		spec.Type = mount.TypeVolume
	***REMOVED***

	spec.ReadOnly = !p.ReadWrite(mode)

	// cannot assume that if a volume driver is passed in that we should set it
	if volumeDriver != "" && spec.Type == mount.TypeVolume ***REMOVED***
		spec.VolumeOptions = &mount.VolumeOptions***REMOVED***
			DriverConfig: &mount.Driver***REMOVED***Name: volumeDriver***REMOVED***,
		***REMOVED***
	***REMOVED***

	if copyData, isSet := getCopyMode(mode, p.DefaultCopyMode()); isSet ***REMOVED***
		if spec.VolumeOptions == nil ***REMOVED***
			spec.VolumeOptions = &mount.VolumeOptions***REMOVED******REMOVED***
		***REMOVED***
		spec.VolumeOptions.NoCopy = !copyData
	***REMOVED***
	if linuxHasPropagation(mode) ***REMOVED***
		spec.BindOptions = &mount.BindOptions***REMOVED***
			Propagation: linuxGetPropagation(mode),
		***REMOVED***
	***REMOVED***

	mp, err := p.parseMountSpec(spec, false)
	if mp != nil ***REMOVED***
		mp.Mode = mode
	***REMOVED***
	if err != nil ***REMOVED***
		err = fmt.Errorf("%v: %v", errInvalidSpec(raw), err)
	***REMOVED***
	return mp, err
***REMOVED***
func (p *linuxParser) ParseMountSpec(cfg mount.Mount) (*MountPoint, error) ***REMOVED***
	return p.parseMountSpec(cfg, true)
***REMOVED***
func (p *linuxParser) parseMountSpec(cfg mount.Mount, validateBindSourceExists bool) (*MountPoint, error) ***REMOVED***
	if err := p.validateMountConfigImpl(&cfg, validateBindSourceExists); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	mp := &MountPoint***REMOVED***
		RW:          !cfg.ReadOnly,
		Destination: path.Clean(filepath.ToSlash(cfg.Target)),
		Type:        cfg.Type,
		Spec:        cfg,
	***REMOVED***

	switch cfg.Type ***REMOVED***
	case mount.TypeVolume:
		if cfg.Source == "" ***REMOVED***
			mp.Name = stringid.GenerateNonCryptoID()
		***REMOVED*** else ***REMOVED***
			mp.Name = cfg.Source
		***REMOVED***
		mp.CopyData = p.DefaultCopyMode()

		if cfg.VolumeOptions != nil ***REMOVED***
			if cfg.VolumeOptions.DriverConfig != nil ***REMOVED***
				mp.Driver = cfg.VolumeOptions.DriverConfig.Name
			***REMOVED***
			if cfg.VolumeOptions.NoCopy ***REMOVED***
				mp.CopyData = false
			***REMOVED***
		***REMOVED***
	case mount.TypeBind:
		mp.Source = path.Clean(filepath.ToSlash(cfg.Source))
		if cfg.BindOptions != nil && len(cfg.BindOptions.Propagation) > 0 ***REMOVED***
			mp.Propagation = cfg.BindOptions.Propagation
		***REMOVED*** else ***REMOVED***
			// If user did not specify a propagation mode, get
			// default propagation mode.
			mp.Propagation = linuxDefaultPropagationMode
		***REMOVED***
	case mount.TypeTmpfs:
		// NOP
	***REMOVED***
	return mp, nil
***REMOVED***

func (p *linuxParser) ParseVolumesFrom(spec string) (string, string, error) ***REMOVED***
	if len(spec) == 0 ***REMOVED***
		return "", "", fmt.Errorf("volumes-from specification cannot be an empty string")
	***REMOVED***

	specParts := strings.SplitN(spec, ":", 2)
	id := specParts[0]
	mode := "rw"

	if len(specParts) == 2 ***REMOVED***
		mode = specParts[1]
		if !linuxValidMountMode(mode) ***REMOVED***
			return "", "", errInvalidMode(mode)
		***REMOVED***
		// For now don't allow propagation properties while importing
		// volumes from data container. These volumes will inherit
		// the same propagation property as of the original volume
		// in data container. This probably can be relaxed in future.
		if linuxHasPropagation(mode) ***REMOVED***
			return "", "", errInvalidMode(mode)
		***REMOVED***
		// Do not allow copy modes on volumes-from
		if _, isSet := getCopyMode(mode, p.DefaultCopyMode()); isSet ***REMOVED***
			return "", "", errInvalidMode(mode)
		***REMOVED***
	***REMOVED***
	return id, mode, nil
***REMOVED***

func (p *linuxParser) DefaultPropagationMode() mount.Propagation ***REMOVED***
	return linuxDefaultPropagationMode
***REMOVED***

func (p *linuxParser) ConvertTmpfsOptions(opt *mount.TmpfsOptions, readOnly bool) (string, error) ***REMOVED***
	var rawOpts []string
	if readOnly ***REMOVED***
		rawOpts = append(rawOpts, "ro")
	***REMOVED***

	if opt != nil && opt.Mode != 0 ***REMOVED***
		rawOpts = append(rawOpts, fmt.Sprintf("mode=%o", opt.Mode))
	***REMOVED***

	if opt != nil && opt.SizeBytes != 0 ***REMOVED***
		// calculate suffix here, making this linux specific, but that is
		// okay, since API is that way anyways.

		// we do this by finding the suffix that divides evenly into the
		// value, returning the value itself, with no suffix, if it fails.
		//
		// For the most part, we don't enforce any semantic to this values.
		// The operating system will usually align this and enforce minimum
		// and maximums.
		var (
			size   = opt.SizeBytes
			suffix string
		)
		for _, r := range []struct ***REMOVED***
			suffix  string
			divisor int64
		***REMOVED******REMOVED***
			***REMOVED***"g", 1 << 30***REMOVED***,
			***REMOVED***"m", 1 << 20***REMOVED***,
			***REMOVED***"k", 1 << 10***REMOVED***,
		***REMOVED*** ***REMOVED***
			if size%r.divisor == 0 ***REMOVED***
				size = size / r.divisor
				suffix = r.suffix
				break
			***REMOVED***
		***REMOVED***

		rawOpts = append(rawOpts, fmt.Sprintf("size=%d%s", size, suffix))
	***REMOVED***
	return strings.Join(rawOpts, ","), nil
***REMOVED***

func (p *linuxParser) DefaultCopyMode() bool ***REMOVED***
	return true
***REMOVED***
func (p *linuxParser) ValidateVolumeName(name string) error ***REMOVED***
	return nil
***REMOVED***

func (p *linuxParser) IsBackwardCompatible(m *MountPoint) bool ***REMOVED***
	return len(m.Source) > 0 || m.Driver == DefaultDriverName
***REMOVED***

func (p *linuxParser) ValidateTmpfsMountDestination(dest string) error ***REMOVED***
	if err := linuxValidateNotRoot(dest); err != nil ***REMOVED***
		return err
	***REMOVED***
	return linuxValidateAbsolute(dest)
***REMOVED***
