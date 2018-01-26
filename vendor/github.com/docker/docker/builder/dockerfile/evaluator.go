// Package dockerfile is the evaluation step in the Dockerfile parse/evaluate pipeline.
//
// It incorporates a dispatch table based on the parser.Node values (see the
// parser package for more information) that are yielded from the parser itself.
// Calling newBuilder with the BuildOpts struct can be used to customize the
// experience for execution purposes only. Parsing is controlled in the parser
// package, and this division of responsibility should be respected.
//
// Please see the jump table targets for the actual invocations, most of which
// will call out to the functions in internals.go to deal with their tasks.
//
// ONBUILD is a special case, which is covered in the onbuild() func in
// dispatchers.go.
//
// The evaluator uses the concept of "steps", which are usually each processable
// line in the Dockerfile. Each step is numbered and certain actions are taken
// before and after each step, such as creating an image ID and removing temporary
// containers and images. Note that ONBUILD creates a kinda-sorta "sub run" which
// includes its own set of steps (usually only one of them).
package dockerfile

import (
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/builder"
	"github.com/docker/docker/builder/dockerfile/instructions"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/docker/runconfig/opts"
	"github.com/pkg/errors"
)

func dispatch(d dispatchRequest, cmd instructions.Command) (err error) ***REMOVED***
	if c, ok := cmd.(instructions.PlatformSpecific); ok ***REMOVED***
		optionsOS := system.ParsePlatform(d.builder.options.Platform).OS
		err := c.CheckPlatform(optionsOS)
		if err != nil ***REMOVED***
			return errdefs.InvalidParameter(err)
		***REMOVED***
	***REMOVED***
	runConfigEnv := d.state.runConfig.Env
	envs := append(runConfigEnv, d.state.buildArgs.FilterAllowed(runConfigEnv)...)

	if ex, ok := cmd.(instructions.SupportsSingleWordExpansion); ok ***REMOVED***
		err := ex.Expand(func(word string) (string, error) ***REMOVED***
			return d.shlex.ProcessWord(word, envs)
		***REMOVED***)
		if err != nil ***REMOVED***
			return errdefs.InvalidParameter(err)
		***REMOVED***
	***REMOVED***

	defer func() ***REMOVED***
		if d.builder.options.ForceRemove ***REMOVED***
			d.builder.containerManager.RemoveAll(d.builder.Stdout)
			return
		***REMOVED***
		if d.builder.options.Remove && err == nil ***REMOVED***
			d.builder.containerManager.RemoveAll(d.builder.Stdout)
			return
		***REMOVED***
	***REMOVED***()
	switch c := cmd.(type) ***REMOVED***
	case *instructions.EnvCommand:
		return dispatchEnv(d, c)
	case *instructions.MaintainerCommand:
		return dispatchMaintainer(d, c)
	case *instructions.LabelCommand:
		return dispatchLabel(d, c)
	case *instructions.AddCommand:
		return dispatchAdd(d, c)
	case *instructions.CopyCommand:
		return dispatchCopy(d, c)
	case *instructions.OnbuildCommand:
		return dispatchOnbuild(d, c)
	case *instructions.WorkdirCommand:
		return dispatchWorkdir(d, c)
	case *instructions.RunCommand:
		return dispatchRun(d, c)
	case *instructions.CmdCommand:
		return dispatchCmd(d, c)
	case *instructions.HealthCheckCommand:
		return dispatchHealthcheck(d, c)
	case *instructions.EntrypointCommand:
		return dispatchEntrypoint(d, c)
	case *instructions.ExposeCommand:
		return dispatchExpose(d, c, envs)
	case *instructions.UserCommand:
		return dispatchUser(d, c)
	case *instructions.VolumeCommand:
		return dispatchVolume(d, c)
	case *instructions.StopSignalCommand:
		return dispatchStopSignal(d, c)
	case *instructions.ArgCommand:
		return dispatchArg(d, c)
	case *instructions.ShellCommand:
		return dispatchShell(d, c)
	***REMOVED***
	return errors.Errorf("unsupported command type: %v", reflect.TypeOf(cmd))
***REMOVED***

// dispatchState is a data object which is modified by dispatchers
type dispatchState struct ***REMOVED***
	runConfig       *container.Config
	maintainer      string
	cmdSet          bool
	imageID         string
	baseImage       builder.Image
	stageName       string
	buildArgs       *buildArgs
	operatingSystem string
***REMOVED***

func newDispatchState(baseArgs *buildArgs) *dispatchState ***REMOVED***
	args := baseArgs.Clone()
	args.ResetAllowed()
	return &dispatchState***REMOVED***runConfig: &container.Config***REMOVED******REMOVED***, buildArgs: args***REMOVED***
***REMOVED***

type stagesBuildResults struct ***REMOVED***
	flat    []*container.Config
	indexed map[string]*container.Config
***REMOVED***

func newStagesBuildResults() *stagesBuildResults ***REMOVED***
	return &stagesBuildResults***REMOVED***
		indexed: make(map[string]*container.Config),
	***REMOVED***
***REMOVED***

func (r *stagesBuildResults) getByName(name string) (*container.Config, bool) ***REMOVED***
	c, ok := r.indexed[strings.ToLower(name)]
	return c, ok
***REMOVED***

func (r *stagesBuildResults) validateIndex(i int) error ***REMOVED***
	if i == len(r.flat) ***REMOVED***
		return errors.New("refers to current build stage")
	***REMOVED***
	if i < 0 || i > len(r.flat) ***REMOVED***
		return errors.New("index out of bounds")
	***REMOVED***
	return nil
***REMOVED***

func (r *stagesBuildResults) get(nameOrIndex string) (*container.Config, error) ***REMOVED***
	if c, ok := r.getByName(nameOrIndex); ok ***REMOVED***
		return c, nil
	***REMOVED***
	ix, err := strconv.ParseInt(nameOrIndex, 10, 0)
	if err != nil ***REMOVED***
		return nil, nil
	***REMOVED***
	if err := r.validateIndex(int(ix)); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return r.flat[ix], nil
***REMOVED***

func (r *stagesBuildResults) checkStageNameAvailable(name string) error ***REMOVED***
	if name != "" ***REMOVED***
		if _, ok := r.getByName(name); ok ***REMOVED***
			return errors.Errorf("%s stage name already used", name)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (r *stagesBuildResults) commitStage(name string, config *container.Config) error ***REMOVED***
	if name != "" ***REMOVED***
		if _, ok := r.getByName(name); ok ***REMOVED***
			return errors.Errorf("%s stage name already used", name)
		***REMOVED***
		r.indexed[strings.ToLower(name)] = config
	***REMOVED***
	r.flat = append(r.flat, config)
	return nil
***REMOVED***

func commitStage(state *dispatchState, stages *stagesBuildResults) error ***REMOVED***
	return stages.commitStage(state.stageName, state.runConfig)
***REMOVED***

type dispatchRequest struct ***REMOVED***
	state   *dispatchState
	shlex   *ShellLex
	builder *Builder
	source  builder.Source
	stages  *stagesBuildResults
***REMOVED***

func newDispatchRequest(builder *Builder, escapeToken rune, source builder.Source, buildArgs *buildArgs, stages *stagesBuildResults) dispatchRequest ***REMOVED***
	return dispatchRequest***REMOVED***
		state:   newDispatchState(buildArgs),
		shlex:   NewShellLex(escapeToken),
		builder: builder,
		source:  source,
		stages:  stages,
	***REMOVED***
***REMOVED***

func (s *dispatchState) updateRunConfig() ***REMOVED***
	s.runConfig.Image = s.imageID
***REMOVED***

// hasFromImage returns true if the builder has processed a `FROM <image>` line
func (s *dispatchState) hasFromImage() bool ***REMOVED***
	return s.imageID != "" || (s.baseImage != nil && s.baseImage.ImageID() == "")
***REMOVED***

func (s *dispatchState) beginStage(stageName string, image builder.Image) error ***REMOVED***
	s.stageName = stageName
	s.imageID = image.ImageID()
	s.operatingSystem = image.OperatingSystem()
	if s.operatingSystem == "" ***REMOVED*** // In case it isn't set
		s.operatingSystem = runtime.GOOS
	***REMOVED***
	if !system.IsOSSupported(s.operatingSystem) ***REMOVED***
		return system.ErrNotSupportedOperatingSystem
	***REMOVED***

	if image.RunConfig() != nil ***REMOVED***
		// copy avoids referencing the same instance when 2 stages have the same base
		s.runConfig = copyRunConfig(image.RunConfig())
	***REMOVED*** else ***REMOVED***
		s.runConfig = &container.Config***REMOVED******REMOVED***
	***REMOVED***
	s.baseImage = image
	s.setDefaultPath()
	s.runConfig.OpenStdin = false
	s.runConfig.StdinOnce = false
	return nil
***REMOVED***

// Add the default PATH to runConfig.ENV if one exists for the operating system and there
// is no PATH set. Note that Windows containers on Windows won't have one as it's set by HCS
func (s *dispatchState) setDefaultPath() ***REMOVED***
	defaultPath := system.DefaultPathEnv(s.operatingSystem)
	if defaultPath == "" ***REMOVED***
		return
	***REMOVED***
	envMap := opts.ConvertKVStringsToMap(s.runConfig.Env)
	if _, ok := envMap["PATH"]; !ok ***REMOVED***
		s.runConfig.Env = append(s.runConfig.Env, "PATH="+defaultPath)
	***REMOVED***
***REMOVED***
