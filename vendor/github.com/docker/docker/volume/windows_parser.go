package volume

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/pkg/stringid"
)

type windowsParser struct ***REMOVED***
***REMOVED***

const (
	// Spec should be in the format [source:]destination[:mode]
	//
	// Examples: c:\foo bar:d:rw
	//           c:\foo:d:\bar
	//           myname:d:
	//           d:\
	//
	// Explanation of this regex! Thanks @thaJeztah on IRC and gist for help. See
	// https://gist.github.com/thaJeztah/6185659e4978789fb2b2. A good place to
	// test is https://regex-golang.appspot.com/assets/html/index.html
	//
	// Useful link for referencing named capturing groups:
	// http://stackoverflow.com/questions/20750843/using-named-matches-from-go-regex
	//
	// There are three match groups: source, destination and mode.
	//

	// rxHostDir is the first option of a source
	rxHostDir = `(?:\\\\\?\\)?[a-z]:[\\/](?:[^\\/:*?"<>|\r\n]+[\\/]?)*`
	// rxName is the second option of a source
	rxName = `[^\\/:*?"<>|\r\n]+`

	// RXReservedNames are reserved names not possible on Windows
	rxReservedNames = `(con)|(prn)|(nul)|(aux)|(com[1-9])|(lpt[1-9])`

	// rxPipe is a named path pipe (starts with `\\.\pipe\`, possibly with / instead of \)
	rxPipe = `[/\\]***REMOVED***2***REMOVED***.[/\\]pipe[/\\][^:*?"<>|\r\n]+`
	// rxSource is the combined possibilities for a source
	rxSource = `((?P<source>((` + rxHostDir + `)|(` + rxName + `)|(` + rxPipe + `))):)?`

	// Source. Can be either a host directory, a name, or omitted:
	//  HostDir:
	//    -  Essentially using the folder solution from
	//       https://www.safaribooksonline.com/library/view/regular-expressions-cookbook/9781449327453/ch08s18.html
	//       but adding case insensitivity.
	//    -  Must be an absolute path such as c:\path
	//    -  Can include spaces such as `c:\program files`
	//    -  And then followed by a colon which is not in the capture group
	//    -  And can be optional
	//  Name:
	//    -  Must not contain invalid NTFS filename characters (https://msdn.microsoft.com/en-us/library/windows/desktop/aa365247(v=vs.85).aspx)
	//    -  And then followed by a colon which is not in the capture group
	//    -  And can be optional

	// rxDestination is the regex expression for the mount destination
	rxDestination = `(?P<destination>((?:\\\\\?\\)?([a-z]):((?:[\\/][^\\/:*?"<>\r\n]+)*[\\/]?))|(` + rxPipe + `))`

	rxLCOWDestination = `(?P<destination>/(?:[^\\/:*?"<>\r\n]+[/]?)*)`
	// Destination (aka container path):
	//    -  Variation on hostdir but can be a drive followed by colon as well
	//    -  If a path, must be absolute. Can include spaces
	//    -  Drive cannot be c: (explicitly checked in code, not RegEx)

	// rxMode is the regex expression for the mode of the mount
	// Mode (optional):
	//    -  Hopefully self explanatory in comparison to above regex's.
	//    -  Colon is not in the capture group
	rxMode = `(:(?P<mode>(?i)ro|rw))?`
)

type mountValidator func(mnt *mount.Mount) error

func windowsSplitRawSpec(raw, destRegex string) ([]string, error) ***REMOVED***
	specExp := regexp.MustCompile(`^` + rxSource + destRegex + rxMode + `$`)
	match := specExp.FindStringSubmatch(strings.ToLower(raw))

	// Must have something back
	if len(match) == 0 ***REMOVED***
		return nil, errInvalidSpec(raw)
	***REMOVED***

	var split []string
	matchgroups := make(map[string]string)
	// Pull out the sub expressions from the named capture groups
	for i, name := range specExp.SubexpNames() ***REMOVED***
		matchgroups[name] = strings.ToLower(match[i])
	***REMOVED***
	if source, exists := matchgroups["source"]; exists ***REMOVED***
		if source != "" ***REMOVED***
			split = append(split, source)
		***REMOVED***
	***REMOVED***
	if destination, exists := matchgroups["destination"]; exists ***REMOVED***
		if destination != "" ***REMOVED***
			split = append(split, destination)
		***REMOVED***
	***REMOVED***
	if mode, exists := matchgroups["mode"]; exists ***REMOVED***
		if mode != "" ***REMOVED***
			split = append(split, mode)
		***REMOVED***
	***REMOVED***
	// Fix #26329. If the destination appears to be a file, and the source is null,
	// it may be because we've fallen through the possible naming regex and hit a
	// situation where the user intention was to map a file into a container through
	// a local volume, but this is not supported by the platform.
	if matchgroups["source"] == "" && matchgroups["destination"] != "" ***REMOVED***
		volExp := regexp.MustCompile(`^` + rxName + `$`)
		reservedNameExp := regexp.MustCompile(`^` + rxReservedNames + `$`)

		if volExp.MatchString(matchgroups["destination"]) ***REMOVED***
			if reservedNameExp.MatchString(matchgroups["destination"]) ***REMOVED***
				return nil, fmt.Errorf("volume name %q cannot be a reserved word for Windows filenames", matchgroups["destination"])
			***REMOVED***
		***REMOVED*** else ***REMOVED***

			exists, isDir, _ := currentFileInfoProvider.fileInfo(matchgroups["destination"])
			if exists && !isDir ***REMOVED***
				return nil, fmt.Errorf("file '%s' cannot be mapped. Only directories can be mapped on this platform", matchgroups["destination"])

			***REMOVED***
		***REMOVED***
	***REMOVED***
	return split, nil
***REMOVED***

func windowsValidMountMode(mode string) bool ***REMOVED***
	if mode == "" ***REMOVED***
		return true
	***REMOVED***
	return rwModes[strings.ToLower(mode)]
***REMOVED***
func windowsValidateNotRoot(p string) error ***REMOVED***
	p = strings.ToLower(strings.Replace(p, `/`, `\`, -1))
	if p == "c:" || p == `c:\` ***REMOVED***
		return fmt.Errorf("destination path cannot be `c:` or `c:\\`: %v", p)
	***REMOVED***
	return nil
***REMOVED***

var windowsSpecificValidators mountValidator = func(mnt *mount.Mount) error ***REMOVED***
	return windowsValidateNotRoot(mnt.Target)
***REMOVED***

func windowsValidateRegex(p, r string) error ***REMOVED***
	if regexp.MustCompile(`^` + r + `$`).MatchString(strings.ToLower(p)) ***REMOVED***
		return nil
	***REMOVED***
	return fmt.Errorf("invalid mount path: '%s'", p)
***REMOVED***
func windowsValidateAbsolute(p string) error ***REMOVED***
	if err := windowsValidateRegex(p, rxDestination); err != nil ***REMOVED***
		return fmt.Errorf("invalid mount path: '%s' mount path must be absolute", p)
	***REMOVED***
	return nil
***REMOVED***

func windowsDetectMountType(p string) mount.Type ***REMOVED***
	if strings.HasPrefix(p, `\\.\pipe\`) ***REMOVED***
		return mount.TypeNamedPipe
	***REMOVED*** else if regexp.MustCompile(`^` + rxHostDir + `$`).MatchString(p) ***REMOVED***
		return mount.TypeBind
	***REMOVED*** else ***REMOVED***
		return mount.TypeVolume
	***REMOVED***
***REMOVED***

func (p *windowsParser) ReadWrite(mode string) bool ***REMOVED***
	return strings.ToLower(mode) != "ro"
***REMOVED***

// IsVolumeNameValid checks a volume name in a platform specific manner.
func (p *windowsParser) ValidateVolumeName(name string) error ***REMOVED***
	nameExp := regexp.MustCompile(`^` + rxName + `$`)
	if !nameExp.MatchString(name) ***REMOVED***
		return errors.New("invalid volume name")
	***REMOVED***
	nameExp = regexp.MustCompile(`^` + rxReservedNames + `$`)
	if nameExp.MatchString(name) ***REMOVED***
		return fmt.Errorf("volume name %q cannot be a reserved word for Windows filenames", name)
	***REMOVED***
	return nil
***REMOVED***
func (p *windowsParser) ValidateMountConfig(mnt *mount.Mount) error ***REMOVED***
	return p.validateMountConfigReg(mnt, rxDestination, windowsSpecificValidators)
***REMOVED***

type fileInfoProvider interface ***REMOVED***
	fileInfo(path string) (exist, isDir bool, err error)
***REMOVED***

type defaultFileInfoProvider struct ***REMOVED***
***REMOVED***

func (defaultFileInfoProvider) fileInfo(path string) (exist, isDir bool, err error) ***REMOVED***
	fi, err := os.Stat(path)
	if err != nil ***REMOVED***
		if !os.IsNotExist(err) ***REMOVED***
			return false, false, err
		***REMOVED***
		return false, false, nil
	***REMOVED***
	return true, fi.IsDir(), nil
***REMOVED***

var currentFileInfoProvider fileInfoProvider = defaultFileInfoProvider***REMOVED******REMOVED***

func (p *windowsParser) validateMountConfigReg(mnt *mount.Mount, destRegex string, additionalValidators ...mountValidator) error ***REMOVED***

	for _, v := range additionalValidators ***REMOVED***
		if err := v(mnt); err != nil ***REMOVED***
			return &errMountConfig***REMOVED***mnt, err***REMOVED***
		***REMOVED***
	***REMOVED***
	if len(mnt.Target) == 0 ***REMOVED***
		return &errMountConfig***REMOVED***mnt, errMissingField("Target")***REMOVED***
	***REMOVED***

	if err := windowsValidateRegex(mnt.Target, destRegex); err != nil ***REMOVED***
		return &errMountConfig***REMOVED***mnt, err***REMOVED***
	***REMOVED***

	switch mnt.Type ***REMOVED***
	case mount.TypeBind:
		if len(mnt.Source) == 0 ***REMOVED***
			return &errMountConfig***REMOVED***mnt, errMissingField("Source")***REMOVED***
		***REMOVED***
		// Don't error out just because the propagation mode is not supported on the platform
		if opts := mnt.BindOptions; opts != nil ***REMOVED***
			if len(opts.Propagation) > 0 ***REMOVED***
				return &errMountConfig***REMOVED***mnt, fmt.Errorf("invalid propagation mode: %s", opts.Propagation)***REMOVED***
			***REMOVED***
		***REMOVED***
		if mnt.VolumeOptions != nil ***REMOVED***
			return &errMountConfig***REMOVED***mnt, errExtraField("VolumeOptions")***REMOVED***
		***REMOVED***

		if err := windowsValidateAbsolute(mnt.Source); err != nil ***REMOVED***
			return &errMountConfig***REMOVED***mnt, err***REMOVED***
		***REMOVED***

		exists, isdir, err := currentFileInfoProvider.fileInfo(mnt.Source)
		if err != nil ***REMOVED***
			return &errMountConfig***REMOVED***mnt, err***REMOVED***
		***REMOVED***
		if !exists ***REMOVED***
			return &errMountConfig***REMOVED***mnt, errBindNotExist***REMOVED***
		***REMOVED***
		if !isdir ***REMOVED***
			return &errMountConfig***REMOVED***mnt, fmt.Errorf("source path must be a directory")***REMOVED***
		***REMOVED***

	case mount.TypeVolume:
		if mnt.BindOptions != nil ***REMOVED***
			return &errMountConfig***REMOVED***mnt, errExtraField("BindOptions")***REMOVED***
		***REMOVED***

		if len(mnt.Source) == 0 && mnt.ReadOnly ***REMOVED***
			return &errMountConfig***REMOVED***mnt, fmt.Errorf("must not set ReadOnly mode when using anonymous volumes")***REMOVED***
		***REMOVED***

		if len(mnt.Source) != 0 ***REMOVED***
			if err := p.ValidateVolumeName(mnt.Source); err != nil ***REMOVED***
				return &errMountConfig***REMOVED***mnt, err***REMOVED***
			***REMOVED***
		***REMOVED***
	case mount.TypeNamedPipe:
		if len(mnt.Source) == 0 ***REMOVED***
			return &errMountConfig***REMOVED***mnt, errMissingField("Source")***REMOVED***
		***REMOVED***

		if mnt.BindOptions != nil ***REMOVED***
			return &errMountConfig***REMOVED***mnt, errExtraField("BindOptions")***REMOVED***
		***REMOVED***

		if mnt.ReadOnly ***REMOVED***
			return &errMountConfig***REMOVED***mnt, errExtraField("ReadOnly")***REMOVED***
		***REMOVED***

		if windowsDetectMountType(mnt.Source) != mount.TypeNamedPipe ***REMOVED***
			return &errMountConfig***REMOVED***mnt, fmt.Errorf("'%s' is not a valid pipe path", mnt.Source)***REMOVED***
		***REMOVED***

		if windowsDetectMountType(mnt.Target) != mount.TypeNamedPipe ***REMOVED***
			return &errMountConfig***REMOVED***mnt, fmt.Errorf("'%s' is not a valid pipe path", mnt.Target)***REMOVED***
		***REMOVED***
	default:
		return &errMountConfig***REMOVED***mnt, errors.New("mount type unknown")***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
func (p *windowsParser) ParseMountRaw(raw, volumeDriver string) (*MountPoint, error) ***REMOVED***
	return p.parseMountRaw(raw, volumeDriver, rxDestination, true, windowsSpecificValidators)
***REMOVED***

func (p *windowsParser) parseMountRaw(raw, volumeDriver, destRegex string, convertTargetToBackslash bool, additionalValidators ...mountValidator) (*MountPoint, error) ***REMOVED***
	arr, err := windowsSplitRawSpec(raw, destRegex)
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
		if windowsValidMountMode(arr[1]) ***REMOVED***
			// Destination + Mode is not a valid volume - volumes
			// cannot include a mode. e.g. /foo:rw
			return nil, errInvalidSpec(raw)
		***REMOVED***
		// Host Source Path or Name + Destination
		spec.Source = strings.Replace(arr[0], `/`, `\`, -1)
		spec.Target = arr[1]
	case 3:
		// HostSourcePath+DestinationPath+Mode
		spec.Source = strings.Replace(arr[0], `/`, `\`, -1)
		spec.Target = arr[1]
		mode = arr[2]
	default:
		return nil, errInvalidSpec(raw)
	***REMOVED***
	if convertTargetToBackslash ***REMOVED***
		spec.Target = strings.Replace(spec.Target, `/`, `\`, -1)
	***REMOVED***

	if !windowsValidMountMode(mode) ***REMOVED***
		return nil, errInvalidMode(mode)
	***REMOVED***

	spec.Type = windowsDetectMountType(spec.Source)
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

	mp, err := p.parseMountSpec(spec, destRegex, convertTargetToBackslash, additionalValidators...)
	if mp != nil ***REMOVED***
		mp.Mode = mode
	***REMOVED***
	if err != nil ***REMOVED***
		err = fmt.Errorf("%v: %v", errInvalidSpec(raw), err)
	***REMOVED***
	return mp, err
***REMOVED***

func (p *windowsParser) ParseMountSpec(cfg mount.Mount) (*MountPoint, error) ***REMOVED***
	return p.parseMountSpec(cfg, rxDestination, true, windowsSpecificValidators)
***REMOVED***
func (p *windowsParser) parseMountSpec(cfg mount.Mount, destRegex string, convertTargetToBackslash bool, additionalValidators ...mountValidator) (*MountPoint, error) ***REMOVED***
	if err := p.validateMountConfigReg(&cfg, destRegex, additionalValidators...); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	mp := &MountPoint***REMOVED***
		RW:          !cfg.ReadOnly,
		Destination: cfg.Target,
		Type:        cfg.Type,
		Spec:        cfg,
	***REMOVED***
	if convertTargetToBackslash ***REMOVED***
		mp.Destination = strings.Replace(cfg.Target, `/`, `\`, -1)
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
		mp.Source = strings.Replace(cfg.Source, `/`, `\`, -1)
	case mount.TypeNamedPipe:
		mp.Source = strings.Replace(cfg.Source, `/`, `\`, -1)
	***REMOVED***
	// cleanup trailing `\` except for paths like `c:\`
	if len(mp.Source) > 3 && mp.Source[len(mp.Source)-1] == '\\' ***REMOVED***
		mp.Source = mp.Source[:len(mp.Source)-1]
	***REMOVED***
	if len(mp.Destination) > 3 && mp.Destination[len(mp.Destination)-1] == '\\' ***REMOVED***
		mp.Destination = mp.Destination[:len(mp.Destination)-1]
	***REMOVED***
	return mp, nil
***REMOVED***

func (p *windowsParser) ParseVolumesFrom(spec string) (string, string, error) ***REMOVED***
	if len(spec) == 0 ***REMOVED***
		return "", "", fmt.Errorf("volumes-from specification cannot be an empty string")
	***REMOVED***

	specParts := strings.SplitN(spec, ":", 2)
	id := specParts[0]
	mode := "rw"

	if len(specParts) == 2 ***REMOVED***
		mode = specParts[1]
		if !windowsValidMountMode(mode) ***REMOVED***
			return "", "", errInvalidMode(mode)
		***REMOVED***

		// Do not allow copy modes on volumes-from
		if _, isSet := getCopyMode(mode, p.DefaultCopyMode()); isSet ***REMOVED***
			return "", "", errInvalidMode(mode)
		***REMOVED***
	***REMOVED***
	return id, mode, nil
***REMOVED***

func (p *windowsParser) DefaultPropagationMode() mount.Propagation ***REMOVED***
	return mount.Propagation("")
***REMOVED***

func (p *windowsParser) ConvertTmpfsOptions(opt *mount.TmpfsOptions, readOnly bool) (string, error) ***REMOVED***
	return "", fmt.Errorf("%s does not support tmpfs", runtime.GOOS)
***REMOVED***
func (p *windowsParser) DefaultCopyMode() bool ***REMOVED***
	return false
***REMOVED***
func (p *windowsParser) IsBackwardCompatible(m *MountPoint) bool ***REMOVED***
	return false
***REMOVED***

func (p *windowsParser) ValidateTmpfsMountDestination(dest string) error ***REMOVED***
	return errors.New("Platform does not support tmpfs")
***REMOVED***
