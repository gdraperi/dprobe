package dockerfile

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"runtime"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/builder"
	"github.com/docker/docker/builder/dockerfile/instructions"
	"github.com/docker/docker/builder/dockerfile/parser"
	"github.com/docker/docker/builder/fscache"
	"github.com/docker/docker/builder/remotecontext"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/streamformatter"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/pkg/system"
	"github.com/moby/buildkit/session"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/sync/syncmap"
)

var validCommitCommands = map[string]bool***REMOVED***
	"cmd":         true,
	"entrypoint":  true,
	"healthcheck": true,
	"env":         true,
	"expose":      true,
	"label":       true,
	"onbuild":     true,
	"user":        true,
	"volume":      true,
	"workdir":     true,
***REMOVED***

const (
	stepFormat = "Step %d/%d : %v"
)

// SessionGetter is object used to get access to a session by uuid
type SessionGetter interface ***REMOVED***
	Get(ctx context.Context, uuid string) (session.Caller, error)
***REMOVED***

// BuildManager is shared across all Builder objects
type BuildManager struct ***REMOVED***
	idMappings *idtools.IDMappings
	backend    builder.Backend
	pathCache  pathCache // TODO: make this persistent
	sg         SessionGetter
	fsCache    *fscache.FSCache
***REMOVED***

// NewBuildManager creates a BuildManager
func NewBuildManager(b builder.Backend, sg SessionGetter, fsCache *fscache.FSCache, idMappings *idtools.IDMappings) (*BuildManager, error) ***REMOVED***
	bm := &BuildManager***REMOVED***
		backend:    b,
		pathCache:  &syncmap.Map***REMOVED******REMOVED***,
		sg:         sg,
		idMappings: idMappings,
		fsCache:    fsCache,
	***REMOVED***
	if err := fsCache.RegisterTransport(remotecontext.ClientSessionRemote, NewClientSessionTransport()); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return bm, nil
***REMOVED***

// Build starts a new build from a BuildConfig
func (bm *BuildManager) Build(ctx context.Context, config backend.BuildConfig) (*builder.Result, error) ***REMOVED***
	buildsTriggered.Inc()
	if config.Options.Dockerfile == "" ***REMOVED***
		config.Options.Dockerfile = builder.DefaultDockerfileName
	***REMOVED***

	source, dockerfile, err := remotecontext.Detect(config)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer func() ***REMOVED***
		if source != nil ***REMOVED***
			if err := source.Close(); err != nil ***REMOVED***
				logrus.Debugf("[BUILDER] failed to remove temporary context: %v", err)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if src, err := bm.initializeClientSession(ctx, cancel, config.Options); err != nil ***REMOVED***
		return nil, err
	***REMOVED*** else if src != nil ***REMOVED***
		source = src
	***REMOVED***

	os := runtime.GOOS
	optionsPlatform := system.ParsePlatform(config.Options.Platform)
	if dockerfile.OS != "" ***REMOVED***
		if optionsPlatform.OS != "" && optionsPlatform.OS != dockerfile.OS ***REMOVED***
			return nil, fmt.Errorf("invalid platform")
		***REMOVED***
		os = dockerfile.OS
	***REMOVED*** else if optionsPlatform.OS != "" ***REMOVED***
		os = optionsPlatform.OS
	***REMOVED***
	config.Options.Platform = os
	dockerfile.OS = os

	builderOptions := builderOptions***REMOVED***
		Options:        config.Options,
		ProgressWriter: config.ProgressWriter,
		Backend:        bm.backend,
		PathCache:      bm.pathCache,
		IDMappings:     bm.idMappings,
	***REMOVED***
	return newBuilder(ctx, builderOptions).build(source, dockerfile)
***REMOVED***

func (bm *BuildManager) initializeClientSession(ctx context.Context, cancel func(), options *types.ImageBuildOptions) (builder.Source, error) ***REMOVED***
	if options.SessionID == "" || bm.sg == nil ***REMOVED***
		return nil, nil
	***REMOVED***
	logrus.Debug("client is session enabled")

	connectCtx, cancelCtx := context.WithTimeout(ctx, sessionConnectTimeout)
	defer cancelCtx()

	c, err := bm.sg.Get(connectCtx, options.SessionID)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	go func() ***REMOVED***
		<-c.Context().Done()
		cancel()
	***REMOVED***()
	if options.RemoteContext == remotecontext.ClientSessionRemote ***REMOVED***
		st := time.Now()
		csi, err := NewClientSessionSourceIdentifier(ctx, bm.sg, options.SessionID)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		src, err := bm.fsCache.SyncFrom(ctx, csi)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		logrus.Debugf("sync-time: %v", time.Since(st))
		return src, nil
	***REMOVED***
	return nil, nil
***REMOVED***

// builderOptions are the dependencies required by the builder
type builderOptions struct ***REMOVED***
	Options        *types.ImageBuildOptions
	Backend        builder.Backend
	ProgressWriter backend.ProgressWriter
	PathCache      pathCache
	IDMappings     *idtools.IDMappings
***REMOVED***

// Builder is a Dockerfile builder
// It implements the builder.Backend interface.
type Builder struct ***REMOVED***
	options *types.ImageBuildOptions

	Stdout io.Writer
	Stderr io.Writer
	Aux    *streamformatter.AuxFormatter
	Output io.Writer

	docker    builder.Backend
	clientCtx context.Context

	idMappings       *idtools.IDMappings
	disableCommit    bool
	imageSources     *imageSources
	pathCache        pathCache
	containerManager *containerManager
	imageProber      ImageProber
***REMOVED***

// newBuilder creates a new Dockerfile builder from an optional dockerfile and a Options.
func newBuilder(clientCtx context.Context, options builderOptions) *Builder ***REMOVED***
	config := options.Options
	if config == nil ***REMOVED***
		config = new(types.ImageBuildOptions)
	***REMOVED***

	b := &Builder***REMOVED***
		clientCtx:        clientCtx,
		options:          config,
		Stdout:           options.ProgressWriter.StdoutFormatter,
		Stderr:           options.ProgressWriter.StderrFormatter,
		Aux:              options.ProgressWriter.AuxFormatter,
		Output:           options.ProgressWriter.Output,
		docker:           options.Backend,
		idMappings:       options.IDMappings,
		imageSources:     newImageSources(clientCtx, options),
		pathCache:        options.PathCache,
		imageProber:      newImageProber(options.Backend, config.CacheFrom, config.NoCache),
		containerManager: newContainerManager(options.Backend),
	***REMOVED***

	return b
***REMOVED***

// Build runs the Dockerfile builder by parsing the Dockerfile and executing
// the instructions from the file.
func (b *Builder) build(source builder.Source, dockerfile *parser.Result) (*builder.Result, error) ***REMOVED***
	defer b.imageSources.Unmount()

	addNodesForLabelOption(dockerfile.AST, b.options.Labels)

	stages, metaArgs, err := instructions.Parse(dockerfile.AST)
	if err != nil ***REMOVED***
		if instructions.IsUnknownInstruction(err) ***REMOVED***
			buildsFailed.WithValues(metricsUnknownInstructionError).Inc()
		***REMOVED***
		return nil, errdefs.InvalidParameter(err)
	***REMOVED***
	if b.options.Target != "" ***REMOVED***
		targetIx, found := instructions.HasStage(stages, b.options.Target)
		if !found ***REMOVED***
			buildsFailed.WithValues(metricsBuildTargetNotReachableError).Inc()
			return nil, errors.Errorf("failed to reach build target %s in Dockerfile", b.options.Target)
		***REMOVED***
		stages = stages[:targetIx+1]
	***REMOVED***

	dockerfile.PrintWarnings(b.Stderr)
	dispatchState, err := b.dispatchDockerfileWithCancellation(stages, metaArgs, dockerfile.EscapeToken, source)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if dispatchState.imageID == "" ***REMOVED***
		buildsFailed.WithValues(metricsDockerfileEmptyError).Inc()
		return nil, errors.New("No image was generated. Is your Dockerfile empty?")
	***REMOVED***
	return &builder.Result***REMOVED***ImageID: dispatchState.imageID, FromImage: dispatchState.baseImage***REMOVED***, nil
***REMOVED***

func emitImageID(aux *streamformatter.AuxFormatter, state *dispatchState) error ***REMOVED***
	if aux == nil || state.imageID == "" ***REMOVED***
		return nil
	***REMOVED***
	return aux.Emit(types.BuildResult***REMOVED***ID: state.imageID***REMOVED***)
***REMOVED***

func processMetaArg(meta instructions.ArgCommand, shlex *ShellLex, args *buildArgs) error ***REMOVED***
	// ShellLex currently only support the concatenated string format
	envs := convertMapToEnvList(args.GetAllAllowed())
	if err := meta.Expand(func(word string) (string, error) ***REMOVED***
		return shlex.ProcessWord(word, envs)
	***REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***
	args.AddArg(meta.Key, meta.Value)
	args.AddMetaArg(meta.Key, meta.Value)
	return nil
***REMOVED***

func printCommand(out io.Writer, currentCommandIndex int, totalCommands int, cmd interface***REMOVED******REMOVED***) int ***REMOVED***
	fmt.Fprintf(out, stepFormat, currentCommandIndex, totalCommands, cmd)
	fmt.Fprintln(out)
	return currentCommandIndex + 1
***REMOVED***

func (b *Builder) dispatchDockerfileWithCancellation(parseResult []instructions.Stage, metaArgs []instructions.ArgCommand, escapeToken rune, source builder.Source) (*dispatchState, error) ***REMOVED***
	dispatchRequest := dispatchRequest***REMOVED******REMOVED***
	buildArgs := newBuildArgs(b.options.BuildArgs)
	totalCommands := len(metaArgs) + len(parseResult)
	currentCommandIndex := 1
	for _, stage := range parseResult ***REMOVED***
		totalCommands += len(stage.Commands)
	***REMOVED***
	shlex := NewShellLex(escapeToken)
	for _, meta := range metaArgs ***REMOVED***
		currentCommandIndex = printCommand(b.Stdout, currentCommandIndex, totalCommands, &meta)

		err := processMetaArg(meta, shlex, buildArgs)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	stagesResults := newStagesBuildResults()

	for _, stage := range parseResult ***REMOVED***
		if err := stagesResults.checkStageNameAvailable(stage.Name); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		dispatchRequest = newDispatchRequest(b, escapeToken, source, buildArgs, stagesResults)

		currentCommandIndex = printCommand(b.Stdout, currentCommandIndex, totalCommands, stage.SourceCode)
		if err := initializeStage(dispatchRequest, &stage); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		dispatchRequest.state.updateRunConfig()
		fmt.Fprintf(b.Stdout, " ---> %s\n", stringid.TruncateID(dispatchRequest.state.imageID))
		for _, cmd := range stage.Commands ***REMOVED***
			select ***REMOVED***
			case <-b.clientCtx.Done():
				logrus.Debug("Builder: build cancelled!")
				fmt.Fprint(b.Stdout, "Build cancelled\n")
				buildsFailed.WithValues(metricsBuildCanceled).Inc()
				return nil, errors.New("Build cancelled")
			default:
				// Not cancelled yet, keep going...
			***REMOVED***

			currentCommandIndex = printCommand(b.Stdout, currentCommandIndex, totalCommands, cmd)

			if err := dispatch(dispatchRequest, cmd); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			dispatchRequest.state.updateRunConfig()
			fmt.Fprintf(b.Stdout, " ---> %s\n", stringid.TruncateID(dispatchRequest.state.imageID))

		***REMOVED***
		if err := emitImageID(b.Aux, dispatchRequest.state); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		buildArgs.MergeReferencedArgs(dispatchRequest.state.buildArgs)
		if err := commitStage(dispatchRequest.state, stagesResults); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	buildArgs.WarnOnUnusedBuildArgs(b.Stdout)
	return dispatchRequest.state, nil
***REMOVED***

func addNodesForLabelOption(dockerfile *parser.Node, labels map[string]string) ***REMOVED***
	if len(labels) == 0 ***REMOVED***
		return
	***REMOVED***

	node := parser.NodeFromLabels(labels)
	dockerfile.Children = append(dockerfile.Children, node)
***REMOVED***

// BuildFromConfig builds directly from `changes`, treating it as if it were the contents of a Dockerfile
// It will:
// - Call parse.Parse() to get an AST root for the concatenated Dockerfile entries.
// - Do build by calling builder.dispatch() to call all entries' handling routines
//
// BuildFromConfig is used by the /commit endpoint, with the changes
// coming from the query parameter of the same name.
//
// TODO: Remove?
func BuildFromConfig(config *container.Config, changes []string, os string) (*container.Config, error) ***REMOVED***
	if !system.IsOSSupported(os) ***REMOVED***
		return nil, errdefs.InvalidParameter(system.ErrNotSupportedOperatingSystem)
	***REMOVED***
	if len(changes) == 0 ***REMOVED***
		return config, nil
	***REMOVED***

	dockerfile, err := parser.Parse(bytes.NewBufferString(strings.Join(changes, "\n")))
	if err != nil ***REMOVED***
		return nil, errdefs.InvalidParameter(err)
	***REMOVED***

	b := newBuilder(context.Background(), builderOptions***REMOVED***
		Options: &types.ImageBuildOptions***REMOVED***NoCache: true***REMOVED***,
	***REMOVED***)

	// ensure that the commands are valid
	for _, n := range dockerfile.AST.Children ***REMOVED***
		if !validCommitCommands[n.Value] ***REMOVED***
			return nil, errdefs.InvalidParameter(errors.Errorf("%s is not a valid change command", n.Value))
		***REMOVED***
	***REMOVED***

	b.Stdout = ioutil.Discard
	b.Stderr = ioutil.Discard
	b.disableCommit = true

	commands := []instructions.Command***REMOVED******REMOVED***
	for _, n := range dockerfile.AST.Children ***REMOVED***
		cmd, err := instructions.ParseCommand(n)
		if err != nil ***REMOVED***
			return nil, errdefs.InvalidParameter(err)
		***REMOVED***
		commands = append(commands, cmd)
	***REMOVED***

	dispatchRequest := newDispatchRequest(b, dockerfile.EscapeToken, nil, newBuildArgs(b.options.BuildArgs), newStagesBuildResults())
	// We make mutations to the configuration, ensure we have a copy
	dispatchRequest.state.runConfig = copyRunConfig(config)
	dispatchRequest.state.imageID = config.Image
	dispatchRequest.state.operatingSystem = os
	for _, cmd := range commands ***REMOVED***
		err := dispatch(dispatchRequest, cmd)
		if err != nil ***REMOVED***
			return nil, errdefs.InvalidParameter(err)
		***REMOVED***
		dispatchRequest.state.updateRunConfig()
	***REMOVED***

	return dispatchRequest.state.runConfig, nil
***REMOVED***

func convertMapToEnvList(m map[string]string) []string ***REMOVED***
	result := []string***REMOVED******REMOVED***
	for k, v := range m ***REMOVED***
		result = append(result, k+"="+v)
	***REMOVED***
	return result
***REMOVED***
