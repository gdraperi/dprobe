package dockerfile

// internals for handling commands. Covers many areas and a lot of
// non-contiguous functionality. Please read the comments.

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/image"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/chrootarchive"
	"github.com/docker/docker/pkg/containerfs"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
)

// Archiver defines an interface for copying files from one destination to
// another using Tar/Untar.
type Archiver interface ***REMOVED***
	TarUntar(src, dst string) error
	UntarPath(src, dst string) error
	CopyWithTar(src, dst string) error
	CopyFileWithTar(src, dst string) error
	IDMappings() *idtools.IDMappings
***REMOVED***

// The builder will use the following interfaces if the container fs implements
// these for optimized copies to and from the container.
type extractor interface ***REMOVED***
	ExtractArchive(src io.Reader, dst string, opts *archive.TarOptions) error
***REMOVED***

type archiver interface ***REMOVED***
	ArchivePath(src string, opts *archive.TarOptions) (io.ReadCloser, error)
***REMOVED***

// helper functions to get tar/untar func
func untarFunc(i interface***REMOVED******REMOVED***) containerfs.UntarFunc ***REMOVED***
	if ea, ok := i.(extractor); ok ***REMOVED***
		return ea.ExtractArchive
	***REMOVED***
	return chrootarchive.Untar
***REMOVED***

func tarFunc(i interface***REMOVED******REMOVED***) containerfs.TarFunc ***REMOVED***
	if ap, ok := i.(archiver); ok ***REMOVED***
		return ap.ArchivePath
	***REMOVED***
	return archive.TarWithOptions
***REMOVED***

func (b *Builder) getArchiver(src, dst containerfs.Driver) Archiver ***REMOVED***
	t, u := tarFunc(src), untarFunc(dst)
	return &containerfs.Archiver***REMOVED***
		SrcDriver:     src,
		DstDriver:     dst,
		Tar:           t,
		Untar:         u,
		IDMappingsVar: b.idMappings,
	***REMOVED***
***REMOVED***

func (b *Builder) commit(dispatchState *dispatchState, comment string) error ***REMOVED***
	if b.disableCommit ***REMOVED***
		return nil
	***REMOVED***
	if !dispatchState.hasFromImage() ***REMOVED***
		return errors.New("Please provide a source image with `from` prior to commit")
	***REMOVED***

	optionsPlatform := system.ParsePlatform(b.options.Platform)
	runConfigWithCommentCmd := copyRunConfig(dispatchState.runConfig, withCmdComment(comment, optionsPlatform.OS))
	hit, err := b.probeCache(dispatchState, runConfigWithCommentCmd)
	if err != nil || hit ***REMOVED***
		return err
	***REMOVED***
	id, err := b.create(runConfigWithCommentCmd)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return b.commitContainer(dispatchState, id, runConfigWithCommentCmd)
***REMOVED***

func (b *Builder) commitContainer(dispatchState *dispatchState, id string, containerConfig *container.Config) error ***REMOVED***
	if b.disableCommit ***REMOVED***
		return nil
	***REMOVED***

	commitCfg := &backend.ContainerCommitConfig***REMOVED***
		ContainerCommitConfig: types.ContainerCommitConfig***REMOVED***
			Author: dispatchState.maintainer,
			Pause:  true,
			// TODO: this should be done by Commit()
			Config: copyRunConfig(dispatchState.runConfig),
		***REMOVED***,
		ContainerConfig: containerConfig,
	***REMOVED***

	// Commit the container
	imageID, err := b.docker.Commit(id, commitCfg)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	dispatchState.imageID = imageID
	return nil
***REMOVED***

func (b *Builder) exportImage(state *dispatchState, imageMount *imageMount, runConfig *container.Config) error ***REMOVED***
	newLayer, err := imageMount.Layer().Commit()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// add an image mount without an image so the layer is properly unmounted
	// if there is an error before we can add the full mount with image
	b.imageSources.Add(newImageMount(nil, newLayer))

	parentImage, ok := imageMount.Image().(*image.Image)
	if !ok ***REMOVED***
		return errors.Errorf("unexpected image type")
	***REMOVED***

	newImage := image.NewChildImage(parentImage, image.ChildConfig***REMOVED***
		Author:          state.maintainer,
		ContainerConfig: runConfig,
		DiffID:          newLayer.DiffID(),
		Config:          copyRunConfig(state.runConfig),
	***REMOVED***, parentImage.OS)

	// TODO: it seems strange to marshal this here instead of just passing in the
	// image struct
	config, err := newImage.MarshalJSON()
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to encode image config")
	***REMOVED***

	exportedImage, err := b.docker.CreateImage(config, state.imageID)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to export image")
	***REMOVED***

	state.imageID = exportedImage.ImageID()
	b.imageSources.Add(newImageMount(exportedImage, newLayer))
	return nil
***REMOVED***

func (b *Builder) performCopy(state *dispatchState, inst copyInstruction) error ***REMOVED***
	srcHash := getSourceHashFromInfos(inst.infos)

	var chownComment string
	if inst.chownStr != "" ***REMOVED***
		chownComment = fmt.Sprintf("--chown=%s", inst.chownStr)
	***REMOVED***
	commentStr := fmt.Sprintf("%s %s%s in %s ", inst.cmdName, chownComment, srcHash, inst.dest)

	// TODO: should this have been using origPaths instead of srcHash in the comment?
	optionsPlatform := system.ParsePlatform(b.options.Platform)
	runConfigWithCommentCmd := copyRunConfig(
		state.runConfig,
		withCmdCommentString(commentStr, optionsPlatform.OS))
	hit, err := b.probeCache(state, runConfigWithCommentCmd)
	if err != nil || hit ***REMOVED***
		return err
	***REMOVED***

	imageMount, err := b.imageSources.Get(state.imageID, true)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to get destination image %q", state.imageID)
	***REMOVED***

	destInfo, err := createDestInfo(state.runConfig.WorkingDir, inst, imageMount, b.options.Platform)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	chownPair := b.idMappings.RootPair()
	// if a chown was requested, perform the steps to get the uid, gid
	// translated (if necessary because of user namespaces), and replace
	// the root pair with the chown pair for copy operations
	if inst.chownStr != "" ***REMOVED***
		chownPair, err = parseChownFlag(inst.chownStr, destInfo.root.Path(), b.idMappings)
		if err != nil ***REMOVED***
			return errors.Wrapf(err, "unable to convert uid/gid chown string to host mapping")
		***REMOVED***
	***REMOVED***

	for _, info := range inst.infos ***REMOVED***
		opts := copyFileOptions***REMOVED***
			decompress: inst.allowLocalDecompression,
			archiver:   b.getArchiver(info.root, destInfo.root),
			chownPair:  chownPair,
		***REMOVED***
		if err := performCopyForInfo(destInfo, info, opts); err != nil ***REMOVED***
			return errors.Wrapf(err, "failed to copy files")
		***REMOVED***
	***REMOVED***
	return b.exportImage(state, imageMount, runConfigWithCommentCmd)
***REMOVED***

func createDestInfo(workingDir string, inst copyInstruction, imageMount *imageMount, platform string) (copyInfo, error) ***REMOVED***
	// Twiddle the destination when it's a relative path - meaning, make it
	// relative to the WORKINGDIR
	dest, err := normalizeDest(workingDir, inst.dest, platform)
	if err != nil ***REMOVED***
		return copyInfo***REMOVED******REMOVED***, errors.Wrapf(err, "invalid %s", inst.cmdName)
	***REMOVED***

	destMount, err := imageMount.Source()
	if err != nil ***REMOVED***
		return copyInfo***REMOVED******REMOVED***, errors.Wrapf(err, "failed to mount copy source")
	***REMOVED***

	return newCopyInfoFromSource(destMount, dest, ""), nil
***REMOVED***

// normalizeDest normalises the destination of a COPY/ADD command in a
// platform semantically consistent way.
func normalizeDest(workingDir, requested string, platform string) (string, error) ***REMOVED***
	dest := fromSlash(requested, platform)
	endsInSlash := strings.HasSuffix(dest, string(separator(platform)))

	if platform != "windows" ***REMOVED***
		if !path.IsAbs(requested) ***REMOVED***
			dest = path.Join("/", filepath.ToSlash(workingDir), dest)
			// Make sure we preserve any trailing slash
			if endsInSlash ***REMOVED***
				dest += "/"
			***REMOVED***
		***REMOVED***
		return dest, nil
	***REMOVED***

	// We are guaranteed that the working directory is already consistent,
	// However, Windows also has, for now, the limitation that ADD/COPY can
	// only be done to the system drive, not any drives that might be present
	// as a result of a bind mount.
	//
	// So... if the path requested is Linux-style absolute (/foo or \\foo),
	// we assume it is the system drive. If it is a Windows-style absolute
	// (DRIVE:\\foo), error if DRIVE is not C. And finally, ensure we
	// strip any configured working directories drive letter so that it
	// can be subsequently legitimately converted to a Windows volume-style
	// pathname.

	// Not a typo - filepath.IsAbs, not system.IsAbs on this next check as
	// we only want to validate where the DriveColon part has been supplied.
	if filepath.IsAbs(dest) ***REMOVED***
		if strings.ToUpper(string(dest[0])) != "C" ***REMOVED***
			return "", fmt.Errorf("Windows does not support destinations not on the system drive (C:)")
		***REMOVED***
		dest = dest[2:] // Strip the drive letter
	***REMOVED***

	// Cannot handle relative where WorkingDir is not the system drive.
	if len(workingDir) > 0 ***REMOVED***
		if ((len(workingDir) > 1) && !system.IsAbs(workingDir[2:])) || (len(workingDir) == 1) ***REMOVED***
			return "", fmt.Errorf("Current WorkingDir %s is not platform consistent", workingDir)
		***REMOVED***
		if !system.IsAbs(dest) ***REMOVED***
			if string(workingDir[0]) != "C" ***REMOVED***
				return "", fmt.Errorf("Windows does not support relative paths when WORKDIR is not the system drive")
			***REMOVED***
			dest = filepath.Join(string(os.PathSeparator), workingDir[2:], dest)
			// Make sure we preserve any trailing slash
			if endsInSlash ***REMOVED***
				dest += string(os.PathSeparator)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return dest, nil
***REMOVED***

// For backwards compat, if there's just one info then use it as the
// cache look-up string, otherwise hash 'em all into one
func getSourceHashFromInfos(infos []copyInfo) string ***REMOVED***
	if len(infos) == 1 ***REMOVED***
		return infos[0].hash
	***REMOVED***
	var hashs []string
	for _, info := range infos ***REMOVED***
		hashs = append(hashs, info.hash)
	***REMOVED***
	return hashStringSlice("multi", hashs)
***REMOVED***

func hashStringSlice(prefix string, slice []string) string ***REMOVED***
	hasher := sha256.New()
	hasher.Write([]byte(strings.Join(slice, ",")))
	return prefix + ":" + hex.EncodeToString(hasher.Sum(nil))
***REMOVED***

type runConfigModifier func(*container.Config)

func withCmd(cmd []string) runConfigModifier ***REMOVED***
	return func(runConfig *container.Config) ***REMOVED***
		runConfig.Cmd = cmd
	***REMOVED***
***REMOVED***

// withCmdComment sets Cmd to a nop comment string. See withCmdCommentString for
// why there are two almost identical versions of this.
func withCmdComment(comment string, platform string) runConfigModifier ***REMOVED***
	return func(runConfig *container.Config) ***REMOVED***
		runConfig.Cmd = append(getShell(runConfig, platform), "#(nop) ", comment)
	***REMOVED***
***REMOVED***

// withCmdCommentString exists to maintain compatibility with older versions.
// A few instructions (workdir, copy, add) used a nop comment that is a single arg
// where as all the other instructions used a two arg comment string. This
// function implements the single arg version.
func withCmdCommentString(comment string, platform string) runConfigModifier ***REMOVED***
	return func(runConfig *container.Config) ***REMOVED***
		runConfig.Cmd = append(getShell(runConfig, platform), "#(nop) "+comment)
	***REMOVED***
***REMOVED***

func withEnv(env []string) runConfigModifier ***REMOVED***
	return func(runConfig *container.Config) ***REMOVED***
		runConfig.Env = env
	***REMOVED***
***REMOVED***

// withEntrypointOverride sets an entrypoint on runConfig if the command is
// not empty. The entrypoint is left unmodified if command is empty.
//
// The dockerfile RUN instruction expect to run without an entrypoint
// so the runConfig entrypoint needs to be modified accordingly. ContainerCreate
// will change a []string***REMOVED***""***REMOVED*** entrypoint to nil, so we probe the cache with the
// nil entrypoint.
func withEntrypointOverride(cmd []string, entrypoint []string) runConfigModifier ***REMOVED***
	return func(runConfig *container.Config) ***REMOVED***
		if len(cmd) > 0 ***REMOVED***
			runConfig.Entrypoint = entrypoint
		***REMOVED***
	***REMOVED***
***REMOVED***

func copyRunConfig(runConfig *container.Config, modifiers ...runConfigModifier) *container.Config ***REMOVED***
	copy := *runConfig
	copy.Cmd = copyStringSlice(runConfig.Cmd)
	copy.Env = copyStringSlice(runConfig.Env)
	copy.Entrypoint = copyStringSlice(runConfig.Entrypoint)
	copy.OnBuild = copyStringSlice(runConfig.OnBuild)
	copy.Shell = copyStringSlice(runConfig.Shell)

	if copy.Volumes != nil ***REMOVED***
		copy.Volumes = make(map[string]struct***REMOVED******REMOVED***, len(runConfig.Volumes))
		for k, v := range runConfig.Volumes ***REMOVED***
			copy.Volumes[k] = v
		***REMOVED***
	***REMOVED***

	if copy.ExposedPorts != nil ***REMOVED***
		copy.ExposedPorts = make(nat.PortSet, len(runConfig.ExposedPorts))
		for k, v := range runConfig.ExposedPorts ***REMOVED***
			copy.ExposedPorts[k] = v
		***REMOVED***
	***REMOVED***

	if copy.Labels != nil ***REMOVED***
		copy.Labels = make(map[string]string, len(runConfig.Labels))
		for k, v := range runConfig.Labels ***REMOVED***
			copy.Labels[k] = v
		***REMOVED***
	***REMOVED***

	for _, modifier := range modifiers ***REMOVED***
		modifier(&copy)
	***REMOVED***
	return &copy
***REMOVED***

func copyStringSlice(orig []string) []string ***REMOVED***
	if orig == nil ***REMOVED***
		return nil
	***REMOVED***
	return append([]string***REMOVED******REMOVED***, orig...)
***REMOVED***

// getShell is a helper function which gets the right shell for prefixing the
// shell-form of RUN, ENTRYPOINT and CMD instructions
func getShell(c *container.Config, os string) []string ***REMOVED***
	if 0 == len(c.Shell) ***REMOVED***
		return append([]string***REMOVED******REMOVED***, defaultShellForOS(os)[:]...)
	***REMOVED***
	return append([]string***REMOVED******REMOVED***, c.Shell[:]...)
***REMOVED***

func (b *Builder) probeCache(dispatchState *dispatchState, runConfig *container.Config) (bool, error) ***REMOVED***
	cachedID, err := b.imageProber.Probe(dispatchState.imageID, runConfig)
	if cachedID == "" || err != nil ***REMOVED***
		return false, err
	***REMOVED***
	fmt.Fprint(b.Stdout, " ---> Using cache\n")

	dispatchState.imageID = cachedID
	return true, nil
***REMOVED***

var defaultLogConfig = container.LogConfig***REMOVED***Type: "none"***REMOVED***

func (b *Builder) probeAndCreate(dispatchState *dispatchState, runConfig *container.Config) (string, error) ***REMOVED***
	if hit, err := b.probeCache(dispatchState, runConfig); err != nil || hit ***REMOVED***
		return "", err
	***REMOVED***
	// Set a log config to override any default value set on the daemon
	hostConfig := &container.HostConfig***REMOVED***LogConfig: defaultLogConfig***REMOVED***
	container, err := b.containerManager.Create(runConfig, hostConfig)
	return container.ID, err
***REMOVED***

func (b *Builder) create(runConfig *container.Config) (string, error) ***REMOVED***
	hostConfig := hostConfigFromOptions(b.options)
	container, err := b.containerManager.Create(runConfig, hostConfig)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	// TODO: could this be moved into containerManager.Create() ?
	for _, warning := range container.Warnings ***REMOVED***
		fmt.Fprintf(b.Stdout, " ---> [Warning] %s\n", warning)
	***REMOVED***
	fmt.Fprintf(b.Stdout, " ---> Running in %s\n", stringid.TruncateID(container.ID))
	return container.ID, nil
***REMOVED***

func hostConfigFromOptions(options *types.ImageBuildOptions) *container.HostConfig ***REMOVED***
	resources := container.Resources***REMOVED***
		CgroupParent: options.CgroupParent,
		CPUShares:    options.CPUShares,
		CPUPeriod:    options.CPUPeriod,
		CPUQuota:     options.CPUQuota,
		CpusetCpus:   options.CPUSetCPUs,
		CpusetMems:   options.CPUSetMems,
		Memory:       options.Memory,
		MemorySwap:   options.MemorySwap,
		Ulimits:      options.Ulimits,
	***REMOVED***

	hc := &container.HostConfig***REMOVED***
		SecurityOpt: options.SecurityOpt,
		Isolation:   options.Isolation,
		ShmSize:     options.ShmSize,
		Resources:   resources,
		NetworkMode: container.NetworkMode(options.NetworkMode),
		// Set a log config to override any default value set on the daemon
		LogConfig:  defaultLogConfig,
		ExtraHosts: options.ExtraHosts,
	***REMOVED***

	// For WCOW, the default of 20GB hard-coded in the platform
	// is too small for builder scenarios where many users are
	// using RUN statements to install large amounts of data.
	// Use 127GB as that's the default size of a VHD in Hyper-V.
	if runtime.GOOS == "windows" && options.Platform == "windows" ***REMOVED***
		hc.StorageOpt = make(map[string]string)
		hc.StorageOpt["size"] = "127GB"
	***REMOVED***

	return hc
***REMOVED***

// fromSlash works like filepath.FromSlash but with a given OS platform field
func fromSlash(path, platform string) string ***REMOVED***
	if platform == "windows" ***REMOVED***
		return strings.Replace(path, "/", "\\", -1)
	***REMOVED***
	return path
***REMOVED***

// separator returns a OS path separator for the given OS platform
func separator(platform string) byte ***REMOVED***
	if platform == "windows" ***REMOVED***
		return '\\'
	***REMOVED***
	return '/'
***REMOVED***
