package dockerfile

import (
	"bytes"
	"context"
	"runtime"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/builder"
	"github.com/docker/docker/builder/dockerfile/instructions"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newBuilderWithMockBackend() *Builder ***REMOVED***
	mockBackend := &MockBackend***REMOVED******REMOVED***
	ctx := context.Background()
	b := &Builder***REMOVED***
		options:       &types.ImageBuildOptions***REMOVED***Platform: runtime.GOOS***REMOVED***,
		docker:        mockBackend,
		Stdout:        new(bytes.Buffer),
		clientCtx:     ctx,
		disableCommit: true,
		imageSources: newImageSources(ctx, builderOptions***REMOVED***
			Options: &types.ImageBuildOptions***REMOVED***Platform: runtime.GOOS***REMOVED***,
			Backend: mockBackend,
		***REMOVED***),
		imageProber:      newImageProber(mockBackend, nil, false),
		containerManager: newContainerManager(mockBackend),
	***REMOVED***
	return b
***REMOVED***

func TestEnv2Variables(t *testing.T) ***REMOVED***
	b := newBuilderWithMockBackend()
	sb := newDispatchRequest(b, '\\', nil, newBuildArgs(make(map[string]*string)), newStagesBuildResults())
	envCommand := &instructions.EnvCommand***REMOVED***
		Env: instructions.KeyValuePairs***REMOVED***
			instructions.KeyValuePair***REMOVED***Key: "var1", Value: "val1"***REMOVED***,
			instructions.KeyValuePair***REMOVED***Key: "var2", Value: "val2"***REMOVED***,
		***REMOVED***,
	***REMOVED***
	err := dispatch(sb, envCommand)
	require.NoError(t, err)

	expected := []string***REMOVED***
		"var1=val1",
		"var2=val2",
	***REMOVED***
	assert.Equal(t, expected, sb.state.runConfig.Env)
***REMOVED***

func TestEnvValueWithExistingRunConfigEnv(t *testing.T) ***REMOVED***
	b := newBuilderWithMockBackend()
	sb := newDispatchRequest(b, '\\', nil, newBuildArgs(make(map[string]*string)), newStagesBuildResults())
	sb.state.runConfig.Env = []string***REMOVED***"var1=old", "var2=fromenv"***REMOVED***
	envCommand := &instructions.EnvCommand***REMOVED***
		Env: instructions.KeyValuePairs***REMOVED***
			instructions.KeyValuePair***REMOVED***Key: "var1", Value: "val1"***REMOVED***,
		***REMOVED***,
	***REMOVED***
	err := dispatch(sb, envCommand)
	require.NoError(t, err)
	expected := []string***REMOVED***
		"var1=val1",
		"var2=fromenv",
	***REMOVED***
	assert.Equal(t, expected, sb.state.runConfig.Env)
***REMOVED***

func TestMaintainer(t *testing.T) ***REMOVED***
	maintainerEntry := "Some Maintainer <maintainer@example.com>"
	b := newBuilderWithMockBackend()
	sb := newDispatchRequest(b, '\\', nil, newBuildArgs(make(map[string]*string)), newStagesBuildResults())
	cmd := &instructions.MaintainerCommand***REMOVED***Maintainer: maintainerEntry***REMOVED***
	err := dispatch(sb, cmd)
	require.NoError(t, err)
	assert.Equal(t, maintainerEntry, sb.state.maintainer)
***REMOVED***

func TestLabel(t *testing.T) ***REMOVED***
	labelName := "label"
	labelValue := "value"

	b := newBuilderWithMockBackend()
	sb := newDispatchRequest(b, '\\', nil, newBuildArgs(make(map[string]*string)), newStagesBuildResults())
	cmd := &instructions.LabelCommand***REMOVED***
		Labels: instructions.KeyValuePairs***REMOVED***
			instructions.KeyValuePair***REMOVED***Key: labelName, Value: labelValue***REMOVED***,
		***REMOVED***,
	***REMOVED***
	err := dispatch(sb, cmd)
	require.NoError(t, err)

	require.Contains(t, sb.state.runConfig.Labels, labelName)
	assert.Equal(t, sb.state.runConfig.Labels[labelName], labelValue)
***REMOVED***

func TestFromScratch(t *testing.T) ***REMOVED***
	b := newBuilderWithMockBackend()
	sb := newDispatchRequest(b, '\\', nil, newBuildArgs(make(map[string]*string)), newStagesBuildResults())
	cmd := &instructions.Stage***REMOVED***
		BaseName: "scratch",
	***REMOVED***
	err := initializeStage(sb, cmd)

	if runtime.GOOS == "windows" && !system.LCOWSupported() ***REMOVED***
		assert.EqualError(t, err, "Windows does not support FROM scratch")
		return
	***REMOVED***

	require.NoError(t, err)
	assert.True(t, sb.state.hasFromImage())
	assert.Equal(t, "", sb.state.imageID)
	expected := "PATH=" + system.DefaultPathEnv(runtime.GOOS)
	assert.Equal(t, []string***REMOVED***expected***REMOVED***, sb.state.runConfig.Env)
***REMOVED***

func TestFromWithArg(t *testing.T) ***REMOVED***
	tag, expected := ":sometag", "expectedthisid"

	getImage := func(name string) (builder.Image, builder.ReleaseableLayer, error) ***REMOVED***
		assert.Equal(t, "alpine"+tag, name)
		return &mockImage***REMOVED***id: "expectedthisid"***REMOVED***, nil, nil
	***REMOVED***
	b := newBuilderWithMockBackend()
	b.docker.(*MockBackend).getImageFunc = getImage
	args := newBuildArgs(make(map[string]*string))

	val := "sometag"
	metaArg := instructions.ArgCommand***REMOVED***
		Key:   "THETAG",
		Value: &val,
	***REMOVED***
	cmd := &instructions.Stage***REMOVED***
		BaseName: "alpine:$***REMOVED***THETAG***REMOVED***",
	***REMOVED***
	err := processMetaArg(metaArg, NewShellLex('\\'), args)

	sb := newDispatchRequest(b, '\\', nil, args, newStagesBuildResults())
	require.NoError(t, err)
	err = initializeStage(sb, cmd)
	require.NoError(t, err)

	assert.Equal(t, expected, sb.state.imageID)
	assert.Equal(t, expected, sb.state.baseImage.ImageID())
	assert.Len(t, sb.state.buildArgs.GetAllAllowed(), 0)
	assert.Len(t, sb.state.buildArgs.GetAllMeta(), 1)
***REMOVED***

func TestFromWithUndefinedArg(t *testing.T) ***REMOVED***
	tag, expected := "sometag", "expectedthisid"

	getImage := func(name string) (builder.Image, builder.ReleaseableLayer, error) ***REMOVED***
		assert.Equal(t, "alpine", name)
		return &mockImage***REMOVED***id: "expectedthisid"***REMOVED***, nil, nil
	***REMOVED***
	b := newBuilderWithMockBackend()
	b.docker.(*MockBackend).getImageFunc = getImage
	sb := newDispatchRequest(b, '\\', nil, newBuildArgs(make(map[string]*string)), newStagesBuildResults())

	b.options.BuildArgs = map[string]*string***REMOVED***"THETAG": &tag***REMOVED***

	cmd := &instructions.Stage***REMOVED***
		BaseName: "alpine$***REMOVED***THETAG***REMOVED***",
	***REMOVED***
	err := initializeStage(sb, cmd)
	require.NoError(t, err)
	assert.Equal(t, expected, sb.state.imageID)
***REMOVED***

func TestFromMultiStageWithNamedStage(t *testing.T) ***REMOVED***
	b := newBuilderWithMockBackend()
	firstFrom := &instructions.Stage***REMOVED***BaseName: "someimg", Name: "base"***REMOVED***
	secondFrom := &instructions.Stage***REMOVED***BaseName: "base"***REMOVED***
	previousResults := newStagesBuildResults()
	firstSB := newDispatchRequest(b, '\\', nil, newBuildArgs(make(map[string]*string)), previousResults)
	secondSB := newDispatchRequest(b, '\\', nil, newBuildArgs(make(map[string]*string)), previousResults)
	err := initializeStage(firstSB, firstFrom)
	require.NoError(t, err)
	assert.True(t, firstSB.state.hasFromImage())
	previousResults.indexed["base"] = firstSB.state.runConfig
	previousResults.flat = append(previousResults.flat, firstSB.state.runConfig)
	err = initializeStage(secondSB, secondFrom)
	require.NoError(t, err)
	assert.True(t, secondSB.state.hasFromImage())
***REMOVED***

func TestOnbuild(t *testing.T) ***REMOVED***
	b := newBuilderWithMockBackend()
	sb := newDispatchRequest(b, '\\', nil, newBuildArgs(make(map[string]*string)), newStagesBuildResults())
	cmd := &instructions.OnbuildCommand***REMOVED***
		Expression: "ADD . /app/src",
	***REMOVED***
	err := dispatch(sb, cmd)
	require.NoError(t, err)
	assert.Equal(t, "ADD . /app/src", sb.state.runConfig.OnBuild[0])
***REMOVED***

func TestWorkdir(t *testing.T) ***REMOVED***
	b := newBuilderWithMockBackend()
	sb := newDispatchRequest(b, '`', nil, newBuildArgs(make(map[string]*string)), newStagesBuildResults())
	workingDir := "/app"
	if runtime.GOOS == "windows" ***REMOVED***
		workingDir = "C:\\app"
	***REMOVED***
	cmd := &instructions.WorkdirCommand***REMOVED***
		Path: workingDir,
	***REMOVED***

	err := dispatch(sb, cmd)
	require.NoError(t, err)
	assert.Equal(t, workingDir, sb.state.runConfig.WorkingDir)
***REMOVED***

func TestCmd(t *testing.T) ***REMOVED***
	b := newBuilderWithMockBackend()
	sb := newDispatchRequest(b, '`', nil, newBuildArgs(make(map[string]*string)), newStagesBuildResults())
	command := "./executable"

	cmd := &instructions.CmdCommand***REMOVED***
		ShellDependantCmdLine: instructions.ShellDependantCmdLine***REMOVED***
			CmdLine:      strslice.StrSlice***REMOVED***command***REMOVED***,
			PrependShell: true,
		***REMOVED***,
	***REMOVED***
	err := dispatch(sb, cmd)
	require.NoError(t, err)

	var expectedCommand strslice.StrSlice
	if runtime.GOOS == "windows" ***REMOVED***
		expectedCommand = strslice.StrSlice(append([]string***REMOVED***"cmd"***REMOVED***, "/S", "/C", command))
	***REMOVED*** else ***REMOVED***
		expectedCommand = strslice.StrSlice(append([]string***REMOVED***"/bin/sh"***REMOVED***, "-c", command))
	***REMOVED***

	assert.Equal(t, expectedCommand, sb.state.runConfig.Cmd)
	assert.True(t, sb.state.cmdSet)
***REMOVED***

func TestHealthcheckNone(t *testing.T) ***REMOVED***
	b := newBuilderWithMockBackend()
	sb := newDispatchRequest(b, '`', nil, newBuildArgs(make(map[string]*string)), newStagesBuildResults())
	cmd := &instructions.HealthCheckCommand***REMOVED***
		Health: &container.HealthConfig***REMOVED***
			Test: []string***REMOVED***"NONE"***REMOVED***,
		***REMOVED***,
	***REMOVED***
	err := dispatch(sb, cmd)
	require.NoError(t, err)

	require.NotNil(t, sb.state.runConfig.Healthcheck)
	assert.Equal(t, []string***REMOVED***"NONE"***REMOVED***, sb.state.runConfig.Healthcheck.Test)
***REMOVED***

func TestHealthcheckCmd(t *testing.T) ***REMOVED***

	b := newBuilderWithMockBackend()
	sb := newDispatchRequest(b, '`', nil, newBuildArgs(make(map[string]*string)), newStagesBuildResults())
	expectedTest := []string***REMOVED***"CMD-SHELL", "curl -f http://localhost/ || exit 1"***REMOVED***
	cmd := &instructions.HealthCheckCommand***REMOVED***
		Health: &container.HealthConfig***REMOVED***
			Test: expectedTest,
		***REMOVED***,
	***REMOVED***
	err := dispatch(sb, cmd)
	require.NoError(t, err)

	require.NotNil(t, sb.state.runConfig.Healthcheck)
	assert.Equal(t, expectedTest, sb.state.runConfig.Healthcheck.Test)
***REMOVED***

func TestEntrypoint(t *testing.T) ***REMOVED***
	b := newBuilderWithMockBackend()
	sb := newDispatchRequest(b, '`', nil, newBuildArgs(make(map[string]*string)), newStagesBuildResults())
	entrypointCmd := "/usr/sbin/nginx"

	cmd := &instructions.EntrypointCommand***REMOVED***
		ShellDependantCmdLine: instructions.ShellDependantCmdLine***REMOVED***
			CmdLine:      strslice.StrSlice***REMOVED***entrypointCmd***REMOVED***,
			PrependShell: true,
		***REMOVED***,
	***REMOVED***
	err := dispatch(sb, cmd)
	require.NoError(t, err)
	require.NotNil(t, sb.state.runConfig.Entrypoint)

	var expectedEntrypoint strslice.StrSlice
	if runtime.GOOS == "windows" ***REMOVED***
		expectedEntrypoint = strslice.StrSlice(append([]string***REMOVED***"cmd"***REMOVED***, "/S", "/C", entrypointCmd))
	***REMOVED*** else ***REMOVED***
		expectedEntrypoint = strslice.StrSlice(append([]string***REMOVED***"/bin/sh"***REMOVED***, "-c", entrypointCmd))
	***REMOVED***
	assert.Equal(t, expectedEntrypoint, sb.state.runConfig.Entrypoint)
***REMOVED***

func TestExpose(t *testing.T) ***REMOVED***
	b := newBuilderWithMockBackend()
	sb := newDispatchRequest(b, '`', nil, newBuildArgs(make(map[string]*string)), newStagesBuildResults())

	exposedPort := "80"
	cmd := &instructions.ExposeCommand***REMOVED***
		Ports: []string***REMOVED***exposedPort***REMOVED***,
	***REMOVED***
	err := dispatch(sb, cmd)
	require.NoError(t, err)

	require.NotNil(t, sb.state.runConfig.ExposedPorts)
	require.Len(t, sb.state.runConfig.ExposedPorts, 1)

	portsMapping, err := nat.ParsePortSpec(exposedPort)
	require.NoError(t, err)
	assert.Contains(t, sb.state.runConfig.ExposedPorts, portsMapping[0].Port)
***REMOVED***

func TestUser(t *testing.T) ***REMOVED***
	b := newBuilderWithMockBackend()
	sb := newDispatchRequest(b, '`', nil, newBuildArgs(make(map[string]*string)), newStagesBuildResults())

	cmd := &instructions.UserCommand***REMOVED***
		User: "test",
	***REMOVED***
	err := dispatch(sb, cmd)
	require.NoError(t, err)
	assert.Equal(t, "test", sb.state.runConfig.User)
***REMOVED***

func TestVolume(t *testing.T) ***REMOVED***
	b := newBuilderWithMockBackend()
	sb := newDispatchRequest(b, '`', nil, newBuildArgs(make(map[string]*string)), newStagesBuildResults())

	exposedVolume := "/foo"

	cmd := &instructions.VolumeCommand***REMOVED***
		Volumes: []string***REMOVED***exposedVolume***REMOVED***,
	***REMOVED***
	err := dispatch(sb, cmd)
	require.NoError(t, err)
	require.NotNil(t, sb.state.runConfig.Volumes)
	assert.Len(t, sb.state.runConfig.Volumes, 1)
	assert.Contains(t, sb.state.runConfig.Volumes, exposedVolume)
***REMOVED***

func TestStopSignal(t *testing.T) ***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Windows does not support stopsignal")
		return
	***REMOVED***
	b := newBuilderWithMockBackend()
	sb := newDispatchRequest(b, '`', nil, newBuildArgs(make(map[string]*string)), newStagesBuildResults())
	signal := "SIGKILL"

	cmd := &instructions.StopSignalCommand***REMOVED***
		Signal: signal,
	***REMOVED***
	err := dispatch(sb, cmd)
	require.NoError(t, err)
	assert.Equal(t, signal, sb.state.runConfig.StopSignal)
***REMOVED***

func TestArg(t *testing.T) ***REMOVED***
	b := newBuilderWithMockBackend()
	sb := newDispatchRequest(b, '`', nil, newBuildArgs(make(map[string]*string)), newStagesBuildResults())

	argName := "foo"
	argVal := "bar"
	cmd := &instructions.ArgCommand***REMOVED***Key: argName, Value: &argVal***REMOVED***
	err := dispatch(sb, cmd)
	require.NoError(t, err)

	expected := map[string]string***REMOVED***argName: argVal***REMOVED***
	assert.Equal(t, expected, sb.state.buildArgs.GetAllAllowed())
***REMOVED***

func TestShell(t *testing.T) ***REMOVED***
	b := newBuilderWithMockBackend()
	sb := newDispatchRequest(b, '`', nil, newBuildArgs(make(map[string]*string)), newStagesBuildResults())

	shellCmd := "powershell"
	cmd := &instructions.ShellCommand***REMOVED***Shell: strslice.StrSlice***REMOVED***shellCmd***REMOVED******REMOVED***

	err := dispatch(sb, cmd)
	require.NoError(t, err)

	expectedShell := strslice.StrSlice([]string***REMOVED***shellCmd***REMOVED***)
	assert.Equal(t, expectedShell, sb.state.runConfig.Shell)
***REMOVED***

func TestPrependEnvOnCmd(t *testing.T) ***REMOVED***
	buildArgs := newBuildArgs(nil)
	buildArgs.AddArg("NO_PROXY", nil)

	args := []string***REMOVED***"sorted=nope", "args=not", "http_proxy=foo", "NO_PROXY=YA"***REMOVED***
	cmd := []string***REMOVED***"foo", "bar"***REMOVED***
	cmdWithEnv := prependEnvOnCmd(buildArgs, args, cmd)
	expected := strslice.StrSlice([]string***REMOVED***
		"|3", "NO_PROXY=YA", "args=not", "sorted=nope", "foo", "bar"***REMOVED***)
	assert.Equal(t, expected, cmdWithEnv)
***REMOVED***

func TestRunWithBuildArgs(t *testing.T) ***REMOVED***
	b := newBuilderWithMockBackend()
	args := newBuildArgs(make(map[string]*string))
	args.argsFromOptions["HTTP_PROXY"] = strPtr("FOO")
	b.disableCommit = false
	sb := newDispatchRequest(b, '`', nil, args, newStagesBuildResults())

	runConfig := &container.Config***REMOVED******REMOVED***
	origCmd := strslice.StrSlice([]string***REMOVED***"cmd", "in", "from", "image"***REMOVED***)
	cmdWithShell := strslice.StrSlice(append(getShell(runConfig, runtime.GOOS), "echo foo"))
	envVars := []string***REMOVED***"|1", "one=two"***REMOVED***
	cachedCmd := strslice.StrSlice(append(envVars, cmdWithShell...))

	imageCache := &mockImageCache***REMOVED***
		getCacheFunc: func(parentID string, cfg *container.Config) (string, error) ***REMOVED***
			// Check the runConfig.Cmd sent to probeCache()
			assert.Equal(t, cachedCmd, cfg.Cmd)
			assert.Equal(t, strslice.StrSlice(nil), cfg.Entrypoint)
			return "", nil
		***REMOVED***,
	***REMOVED***

	mockBackend := b.docker.(*MockBackend)
	mockBackend.makeImageCacheFunc = func(_ []string) builder.ImageCache ***REMOVED***
		return imageCache
	***REMOVED***
	b.imageProber = newImageProber(mockBackend, nil, false)
	mockBackend.getImageFunc = func(_ string) (builder.Image, builder.ReleaseableLayer, error) ***REMOVED***
		return &mockImage***REMOVED***
			id:     "abcdef",
			config: &container.Config***REMOVED***Cmd: origCmd***REMOVED***,
		***REMOVED***, nil, nil
	***REMOVED***
	mockBackend.containerCreateFunc = func(config types.ContainerCreateConfig) (container.ContainerCreateCreatedBody, error) ***REMOVED***
		// Check the runConfig.Cmd sent to create()
		assert.Equal(t, cmdWithShell, config.Config.Cmd)
		assert.Contains(t, config.Config.Env, "one=two")
		assert.Equal(t, strslice.StrSlice***REMOVED***""***REMOVED***, config.Config.Entrypoint)
		return container.ContainerCreateCreatedBody***REMOVED***ID: "12345"***REMOVED***, nil
	***REMOVED***
	mockBackend.commitFunc = func(cID string, cfg *backend.ContainerCommitConfig) (string, error) ***REMOVED***
		// Check the runConfig.Cmd sent to commit()
		assert.Equal(t, origCmd, cfg.Config.Cmd)
		assert.Equal(t, cachedCmd, cfg.ContainerConfig.Cmd)
		assert.Equal(t, strslice.StrSlice(nil), cfg.Config.Entrypoint)
		return "", nil
	***REMOVED***
	from := &instructions.Stage***REMOVED***BaseName: "abcdef"***REMOVED***
	err := initializeStage(sb, from)
	require.NoError(t, err)
	sb.state.buildArgs.AddArg("one", strPtr("two"))
	run := &instructions.RunCommand***REMOVED***
		ShellDependantCmdLine: instructions.ShellDependantCmdLine***REMOVED***
			CmdLine:      strslice.StrSlice***REMOVED***"echo foo"***REMOVED***,
			PrependShell: true,
		***REMOVED***,
	***REMOVED***
	require.NoError(t, dispatch(sb, run))

	// Check that runConfig.Cmd has not been modified by run
	assert.Equal(t, origCmd, sb.state.runConfig.Cmd)
***REMOVED***
