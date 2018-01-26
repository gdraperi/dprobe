package dockerfile

// This file contains the dispatchers for each command. Note that
// `nullDispatch` is not actually a command, but support for commands we parse
// but do nothing with.
//
// See evaluator.go for a higher level discussion of the whole evaluator
// package.

import (
	"bytes"
	"fmt"
	"runtime"
	"sort"
	"strings"

	"github.com/docker/docker/api"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/builder"
	"github.com/docker/docker/builder/dockerfile/instructions"
	"github.com/docker/docker/builder/dockerfile/parser"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/image"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/signal"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ENV foo bar
//
// Sets the environment variable foo to bar, also makes interpolation
// in the dockerfile available from the next statement on via $***REMOVED***foo***REMOVED***.
//
func dispatchEnv(d dispatchRequest, c *instructions.EnvCommand) error ***REMOVED***
	runConfig := d.state.runConfig
	commitMessage := bytes.NewBufferString("ENV")
	for _, e := range c.Env ***REMOVED***
		name := e.Key
		newVar := e.String()

		commitMessage.WriteString(" " + newVar)
		gotOne := false
		for i, envVar := range runConfig.Env ***REMOVED***
			envParts := strings.SplitN(envVar, "=", 2)
			compareFrom := envParts[0]
			if equalEnvKeys(compareFrom, name) ***REMOVED***
				runConfig.Env[i] = newVar
				gotOne = true
				break
			***REMOVED***
		***REMOVED***
		if !gotOne ***REMOVED***
			runConfig.Env = append(runConfig.Env, newVar)
		***REMOVED***
	***REMOVED***
	return d.builder.commit(d.state, commitMessage.String())
***REMOVED***

// MAINTAINER some text <maybe@an.email.address>
//
// Sets the maintainer metadata.
func dispatchMaintainer(d dispatchRequest, c *instructions.MaintainerCommand) error ***REMOVED***

	d.state.maintainer = c.Maintainer
	return d.builder.commit(d.state, "MAINTAINER "+c.Maintainer)
***REMOVED***

// LABEL some json data describing the image
//
// Sets the Label variable foo to bar,
//
func dispatchLabel(d dispatchRequest, c *instructions.LabelCommand) error ***REMOVED***
	if d.state.runConfig.Labels == nil ***REMOVED***
		d.state.runConfig.Labels = make(map[string]string)
	***REMOVED***
	commitStr := "LABEL"
	for _, v := range c.Labels ***REMOVED***
		d.state.runConfig.Labels[v.Key] = v.Value
		commitStr += " " + v.String()
	***REMOVED***
	return d.builder.commit(d.state, commitStr)
***REMOVED***

// ADD foo /path
//
// Add the file 'foo' to '/path'. Tarball and Remote URL (git, http) handling
// exist here. If you do not wish to have this automatic handling, use COPY.
//
func dispatchAdd(d dispatchRequest, c *instructions.AddCommand) error ***REMOVED***
	downloader := newRemoteSourceDownloader(d.builder.Output, d.builder.Stdout)
	copier := copierFromDispatchRequest(d, downloader, nil)
	defer copier.Cleanup()

	copyInstruction, err := copier.createCopyInstruction(c.SourcesAndDest, "ADD")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	copyInstruction.chownStr = c.Chown
	copyInstruction.allowLocalDecompression = true

	return d.builder.performCopy(d.state, copyInstruction)
***REMOVED***

// COPY foo /path
//
// Same as 'ADD' but without the tar and remote url handling.
//
func dispatchCopy(d dispatchRequest, c *instructions.CopyCommand) error ***REMOVED***
	var im *imageMount
	var err error
	if c.From != "" ***REMOVED***
		im, err = d.getImageMount(c.From)
		if err != nil ***REMOVED***
			return errors.Wrapf(err, "invalid from flag value %s", c.From)
		***REMOVED***
	***REMOVED***
	copier := copierFromDispatchRequest(d, errOnSourceDownload, im)
	defer copier.Cleanup()
	copyInstruction, err := copier.createCopyInstruction(c.SourcesAndDest, "COPY")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	copyInstruction.chownStr = c.Chown

	return d.builder.performCopy(d.state, copyInstruction)
***REMOVED***

func (d *dispatchRequest) getImageMount(imageRefOrID string) (*imageMount, error) ***REMOVED***
	if imageRefOrID == "" ***REMOVED***
		// TODO: this could return the source in the default case as well?
		return nil, nil
	***REMOVED***

	var localOnly bool
	stage, err := d.stages.get(imageRefOrID)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if stage != nil ***REMOVED***
		imageRefOrID = stage.Image
		localOnly = true
	***REMOVED***
	return d.builder.imageSources.Get(imageRefOrID, localOnly)
***REMOVED***

// FROM imagename[:tag | @digest] [AS build-stage-name]
//
func initializeStage(d dispatchRequest, cmd *instructions.Stage) error ***REMOVED***
	d.builder.imageProber.Reset()
	image, err := d.getFromImage(d.shlex, cmd.BaseName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	state := d.state
	if err := state.beginStage(cmd.Name, image); err != nil ***REMOVED***
		return err
	***REMOVED***
	if len(state.runConfig.OnBuild) > 0 ***REMOVED***
		triggers := state.runConfig.OnBuild
		state.runConfig.OnBuild = nil
		return dispatchTriggeredOnBuild(d, triggers)
	***REMOVED***
	return nil
***REMOVED***

func dispatchTriggeredOnBuild(d dispatchRequest, triggers []string) error ***REMOVED***
	fmt.Fprintf(d.builder.Stdout, "# Executing %d build trigger", len(triggers))
	if len(triggers) > 1 ***REMOVED***
		fmt.Fprint(d.builder.Stdout, "s")
	***REMOVED***
	fmt.Fprintln(d.builder.Stdout)
	for _, trigger := range triggers ***REMOVED***
		d.state.updateRunConfig()
		ast, err := parser.Parse(strings.NewReader(trigger))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if len(ast.AST.Children) != 1 ***REMOVED***
			return errors.New("onbuild trigger should be a single expression")
		***REMOVED***
		cmd, err := instructions.ParseCommand(ast.AST.Children[0])
		if err != nil ***REMOVED***
			if instructions.IsUnknownInstruction(err) ***REMOVED***
				buildsFailed.WithValues(metricsUnknownInstructionError).Inc()
			***REMOVED***
			return err
		***REMOVED***
		err = dispatch(d, cmd)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (d *dispatchRequest) getExpandedImageName(shlex *ShellLex, name string) (string, error) ***REMOVED***
	substitutionArgs := []string***REMOVED******REMOVED***
	for key, value := range d.state.buildArgs.GetAllMeta() ***REMOVED***
		substitutionArgs = append(substitutionArgs, key+"="+value)
	***REMOVED***

	name, err := shlex.ProcessWord(name, substitutionArgs)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return name, nil
***REMOVED***
func (d *dispatchRequest) getImageOrStage(name string) (builder.Image, error) ***REMOVED***
	var localOnly bool
	if im, ok := d.stages.getByName(name); ok ***REMOVED***
		name = im.Image
		localOnly = true
	***REMOVED***

	// Windows cannot support a container with no base image unless it is LCOW.
	if name == api.NoBaseImageSpecifier ***REMOVED***
		imageImage := &image.Image***REMOVED******REMOVED***
		imageImage.OS = runtime.GOOS
		if runtime.GOOS == "windows" ***REMOVED***
			optionsOS := system.ParsePlatform(d.builder.options.Platform).OS
			switch optionsOS ***REMOVED***
			case "windows", "":
				return nil, errors.New("Windows does not support FROM scratch")
			case "linux":
				if !system.LCOWSupported() ***REMOVED***
					return nil, errors.New("Linux containers are not supported on this system")
				***REMOVED***
				imageImage.OS = "linux"
			default:
				return nil, errors.Errorf("operating system %q is not supported", optionsOS)
			***REMOVED***
		***REMOVED***
		return builder.Image(imageImage), nil
	***REMOVED***
	imageMount, err := d.builder.imageSources.Get(name, localOnly)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return imageMount.Image(), nil
***REMOVED***
func (d *dispatchRequest) getFromImage(shlex *ShellLex, name string) (builder.Image, error) ***REMOVED***
	name, err := d.getExpandedImageName(shlex, name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return d.getImageOrStage(name)
***REMOVED***

func dispatchOnbuild(d dispatchRequest, c *instructions.OnbuildCommand) error ***REMOVED***

	d.state.runConfig.OnBuild = append(d.state.runConfig.OnBuild, c.Expression)
	return d.builder.commit(d.state, "ONBUILD "+c.Expression)
***REMOVED***

// WORKDIR /tmp
//
// Set the working directory for future RUN/CMD/etc statements.
//
func dispatchWorkdir(d dispatchRequest, c *instructions.WorkdirCommand) error ***REMOVED***
	runConfig := d.state.runConfig
	var err error
	baseImageOS := system.ParsePlatform(d.state.operatingSystem).OS
	runConfig.WorkingDir, err = normalizeWorkdir(baseImageOS, runConfig.WorkingDir, c.Path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// For performance reasons, we explicitly do a create/mkdir now
	// This avoids having an unnecessary expensive mount/unmount calls
	// (on Windows in particular) during each container create.
	// Prior to 1.13, the mkdir was deferred and not executed at this step.
	if d.builder.disableCommit ***REMOVED***
		// Don't call back into the daemon if we're going through docker commit --change "WORKDIR /foo".
		// We've already updated the runConfig and that's enough.
		return nil
	***REMOVED***

	comment := "WORKDIR " + runConfig.WorkingDir
	runConfigWithCommentCmd := copyRunConfig(runConfig, withCmdCommentString(comment, baseImageOS))
	containerID, err := d.builder.probeAndCreate(d.state, runConfigWithCommentCmd)
	if err != nil || containerID == "" ***REMOVED***
		return err
	***REMOVED***
	if err := d.builder.docker.ContainerCreateWorkdir(containerID); err != nil ***REMOVED***
		return err
	***REMOVED***

	return d.builder.commitContainer(d.state, containerID, runConfigWithCommentCmd)
***REMOVED***

func resolveCmdLine(cmd instructions.ShellDependantCmdLine, runConfig *container.Config, os string) []string ***REMOVED***
	result := cmd.CmdLine
	if cmd.PrependShell && result != nil ***REMOVED***
		result = append(getShell(runConfig, os), result...)
	***REMOVED***
	return result
***REMOVED***

// RUN some command yo
//
// run a command and commit the image. Args are automatically prepended with
// the current SHELL which defaults to 'sh -c' under linux or 'cmd /S /C' under
// Windows, in the event there is only one argument The difference in processing:
//
// RUN echo hi          # sh -c echo hi       (Linux and LCOW)
// RUN echo hi          # cmd /S /C echo hi   (Windows)
// RUN [ "echo", "hi" ] # echo hi
//
func dispatchRun(d dispatchRequest, c *instructions.RunCommand) error ***REMOVED***
	if !system.IsOSSupported(d.state.operatingSystem) ***REMOVED***
		return system.ErrNotSupportedOperatingSystem
	***REMOVED***
	stateRunConfig := d.state.runConfig
	cmdFromArgs := resolveCmdLine(c.ShellDependantCmdLine, stateRunConfig, d.state.operatingSystem)
	buildArgs := d.state.buildArgs.FilterAllowed(stateRunConfig.Env)

	saveCmd := cmdFromArgs
	if len(buildArgs) > 0 ***REMOVED***
		saveCmd = prependEnvOnCmd(d.state.buildArgs, buildArgs, cmdFromArgs)
	***REMOVED***

	runConfigForCacheProbe := copyRunConfig(stateRunConfig,
		withCmd(saveCmd),
		withEntrypointOverride(saveCmd, nil))
	hit, err := d.builder.probeCache(d.state, runConfigForCacheProbe)
	if err != nil || hit ***REMOVED***
		return err
	***REMOVED***

	runConfig := copyRunConfig(stateRunConfig,
		withCmd(cmdFromArgs),
		withEnv(append(stateRunConfig.Env, buildArgs...)),
		withEntrypointOverride(saveCmd, strslice.StrSlice***REMOVED***""***REMOVED***))

	// set config as already being escaped, this prevents double escaping on windows
	runConfig.ArgsEscaped = true

	logrus.Debugf("[BUILDER] Command to be executed: %v", runConfig.Cmd)
	cID, err := d.builder.create(runConfig)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := d.builder.containerManager.Run(d.builder.clientCtx, cID, d.builder.Stdout, d.builder.Stderr); err != nil ***REMOVED***
		if err, ok := err.(*statusCodeError); ok ***REMOVED***
			// TODO: change error type, because jsonmessage.JSONError assumes HTTP
			return &jsonmessage.JSONError***REMOVED***
				Message: fmt.Sprintf(
					"The command '%s' returned a non-zero code: %d",
					strings.Join(runConfig.Cmd, " "), err.StatusCode()),
				Code: err.StatusCode(),
			***REMOVED***
		***REMOVED***
		return err
	***REMOVED***

	return d.builder.commitContainer(d.state, cID, runConfigForCacheProbe)
***REMOVED***

// Derive the command to use for probeCache() and to commit in this container.
// Note that we only do this if there are any build-time env vars.  Also, we
// use the special argument "|#" at the start of the args array. This will
// avoid conflicts with any RUN command since commands can not
// start with | (vertical bar). The "#" (number of build envs) is there to
// help ensure proper cache matches. We don't want a RUN command
// that starts with "foo=abc" to be considered part of a build-time env var.
//
// remove any unreferenced built-in args from the environment variables.
// These args are transparent so resulting image should be the same regardless
// of the value.
func prependEnvOnCmd(buildArgs *buildArgs, buildArgVars []string, cmd strslice.StrSlice) strslice.StrSlice ***REMOVED***
	var tmpBuildEnv []string
	for _, env := range buildArgVars ***REMOVED***
		key := strings.SplitN(env, "=", 2)[0]
		if buildArgs.IsReferencedOrNotBuiltin(key) ***REMOVED***
			tmpBuildEnv = append(tmpBuildEnv, env)
		***REMOVED***
	***REMOVED***

	sort.Strings(tmpBuildEnv)
	tmpEnv := append([]string***REMOVED***fmt.Sprintf("|%d", len(tmpBuildEnv))***REMOVED***, tmpBuildEnv...)
	return strslice.StrSlice(append(tmpEnv, cmd...))
***REMOVED***

// CMD foo
//
// Set the default command to run in the container (which may be empty).
// Argument handling is the same as RUN.
//
func dispatchCmd(d dispatchRequest, c *instructions.CmdCommand) error ***REMOVED***
	runConfig := d.state.runConfig
	optionsOS := system.ParsePlatform(d.builder.options.Platform).OS
	cmd := resolveCmdLine(c.ShellDependantCmdLine, runConfig, optionsOS)
	runConfig.Cmd = cmd
	// set config as already being escaped, this prevents double escaping on windows
	runConfig.ArgsEscaped = true

	if err := d.builder.commit(d.state, fmt.Sprintf("CMD %q", cmd)); err != nil ***REMOVED***
		return err
	***REMOVED***

	if len(c.ShellDependantCmdLine.CmdLine) != 0 ***REMOVED***
		d.state.cmdSet = true
	***REMOVED***

	return nil
***REMOVED***

// HEALTHCHECK foo
//
// Set the default healthcheck command to run in the container (which may be empty).
// Argument handling is the same as RUN.
//
func dispatchHealthcheck(d dispatchRequest, c *instructions.HealthCheckCommand) error ***REMOVED***
	runConfig := d.state.runConfig
	if runConfig.Healthcheck != nil ***REMOVED***
		oldCmd := runConfig.Healthcheck.Test
		if len(oldCmd) > 0 && oldCmd[0] != "NONE" ***REMOVED***
			fmt.Fprintf(d.builder.Stdout, "Note: overriding previous HEALTHCHECK: %v\n", oldCmd)
		***REMOVED***
	***REMOVED***
	runConfig.Healthcheck = c.Health
	return d.builder.commit(d.state, fmt.Sprintf("HEALTHCHECK %q", runConfig.Healthcheck))
***REMOVED***

// ENTRYPOINT /usr/sbin/nginx
//
// Set the entrypoint to /usr/sbin/nginx. Will accept the CMD as the arguments
// to /usr/sbin/nginx. Uses the default shell if not in JSON format.
//
// Handles command processing similar to CMD and RUN, only req.runConfig.Entrypoint
// is initialized at newBuilder time instead of through argument parsing.
//
func dispatchEntrypoint(d dispatchRequest, c *instructions.EntrypointCommand) error ***REMOVED***
	runConfig := d.state.runConfig
	optionsOS := system.ParsePlatform(d.builder.options.Platform).OS
	cmd := resolveCmdLine(c.ShellDependantCmdLine, runConfig, optionsOS)
	runConfig.Entrypoint = cmd
	if !d.state.cmdSet ***REMOVED***
		runConfig.Cmd = nil
	***REMOVED***

	return d.builder.commit(d.state, fmt.Sprintf("ENTRYPOINT %q", runConfig.Entrypoint))
***REMOVED***

// EXPOSE 6667/tcp 7000/tcp
//
// Expose ports for links and port mappings. This all ends up in
// req.runConfig.ExposedPorts for runconfig.
//
func dispatchExpose(d dispatchRequest, c *instructions.ExposeCommand, envs []string) error ***REMOVED***
	// custom multi word expansion
	// expose $FOO with FOO="80 443" is expanded as EXPOSE [80,443]. This is the only command supporting word to words expansion
	// so the word processing has been de-generalized
	ports := []string***REMOVED******REMOVED***
	for _, p := range c.Ports ***REMOVED***
		ps, err := d.shlex.ProcessWords(p, envs)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		ports = append(ports, ps...)
	***REMOVED***
	c.Ports = ports

	ps, _, err := nat.ParsePortSpecs(ports)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if d.state.runConfig.ExposedPorts == nil ***REMOVED***
		d.state.runConfig.ExposedPorts = make(nat.PortSet)
	***REMOVED***
	for p := range ps ***REMOVED***
		d.state.runConfig.ExposedPorts[p] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	return d.builder.commit(d.state, "EXPOSE "+strings.Join(c.Ports, " "))
***REMOVED***

// USER foo
//
// Set the user to 'foo' for future commands and when running the
// ENTRYPOINT/CMD at container run time.
//
func dispatchUser(d dispatchRequest, c *instructions.UserCommand) error ***REMOVED***
	d.state.runConfig.User = c.User
	return d.builder.commit(d.state, fmt.Sprintf("USER %v", c.User))
***REMOVED***

// VOLUME /foo
//
// Expose the volume /foo for use. Will also accept the JSON array form.
//
func dispatchVolume(d dispatchRequest, c *instructions.VolumeCommand) error ***REMOVED***
	if d.state.runConfig.Volumes == nil ***REMOVED***
		d.state.runConfig.Volumes = map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
	for _, v := range c.Volumes ***REMOVED***
		if v == "" ***REMOVED***
			return errors.New("VOLUME specified can not be an empty string")
		***REMOVED***
		d.state.runConfig.Volumes[v] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***
	return d.builder.commit(d.state, fmt.Sprintf("VOLUME %v", c.Volumes))
***REMOVED***

// STOPSIGNAL signal
//
// Set the signal that will be used to kill the container.
func dispatchStopSignal(d dispatchRequest, c *instructions.StopSignalCommand) error ***REMOVED***

	_, err := signal.ParseSignal(c.Signal)
	if err != nil ***REMOVED***
		return errdefs.InvalidParameter(err)
	***REMOVED***
	d.state.runConfig.StopSignal = c.Signal
	return d.builder.commit(d.state, fmt.Sprintf("STOPSIGNAL %v", c.Signal))
***REMOVED***

// ARG name[=value]
//
// Adds the variable foo to the trusted list of variables that can be passed
// to builder using the --build-arg flag for expansion/substitution or passing to 'run'.
// Dockerfile author may optionally set a default value of this variable.
func dispatchArg(d dispatchRequest, c *instructions.ArgCommand) error ***REMOVED***

	commitStr := "ARG " + c.Key
	if c.Value != nil ***REMOVED***
		commitStr += "=" + *c.Value
	***REMOVED***

	d.state.buildArgs.AddArg(c.Key, c.Value)
	return d.builder.commit(d.state, commitStr)
***REMOVED***

// SHELL powershell -command
//
// Set the non-default shell to use.
func dispatchShell(d dispatchRequest, c *instructions.ShellCommand) error ***REMOVED***
	d.state.runConfig.Shell = c.Shell
	return d.builder.commit(d.state, fmt.Sprintf("SHELL %v", d.state.runConfig.Shell))
***REMOVED***
