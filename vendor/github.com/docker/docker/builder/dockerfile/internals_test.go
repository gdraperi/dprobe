package dockerfile

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/builder"
	"github.com/docker/docker/builder/remotecontext"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyDockerfile(t *testing.T) ***REMOVED***
	contextDir, cleanup := createTestTempDir(t, "", "builder-dockerfile-test")
	defer cleanup()

	createTestTempFile(t, contextDir, builder.DefaultDockerfileName, "", 0777)

	readAndCheckDockerfile(t, "emptyDockerfile", contextDir, "", "the Dockerfile (Dockerfile) cannot be empty")
***REMOVED***

func TestSymlinkDockerfile(t *testing.T) ***REMOVED***
	contextDir, cleanup := createTestTempDir(t, "", "builder-dockerfile-test")
	defer cleanup()

	createTestSymlink(t, contextDir, builder.DefaultDockerfileName, "/etc/passwd")

	// The reason the error is "Cannot locate specified Dockerfile" is because
	// in the builder, the symlink is resolved within the context, therefore
	// Dockerfile -> /etc/passwd becomes etc/passwd from the context which is
	// a nonexistent file.
	expectedError := fmt.Sprintf("Cannot locate specified Dockerfile: %s", builder.DefaultDockerfileName)

	readAndCheckDockerfile(t, "symlinkDockerfile", contextDir, builder.DefaultDockerfileName, expectedError)
***REMOVED***

func TestDockerfileOutsideTheBuildContext(t *testing.T) ***REMOVED***
	contextDir, cleanup := createTestTempDir(t, "", "builder-dockerfile-test")
	defer cleanup()

	expectedError := "Forbidden path outside the build context: ../../Dockerfile ()"

	readAndCheckDockerfile(t, "DockerfileOutsideTheBuildContext", contextDir, "../../Dockerfile", expectedError)
***REMOVED***

func TestNonExistingDockerfile(t *testing.T) ***REMOVED***
	contextDir, cleanup := createTestTempDir(t, "", "builder-dockerfile-test")
	defer cleanup()

	expectedError := "Cannot locate specified Dockerfile: Dockerfile"

	readAndCheckDockerfile(t, "NonExistingDockerfile", contextDir, "Dockerfile", expectedError)
***REMOVED***

func readAndCheckDockerfile(t *testing.T, testName, contextDir, dockerfilePath, expectedError string) ***REMOVED***
	tarStream, err := archive.Tar(contextDir, archive.Uncompressed)
	require.NoError(t, err)

	defer func() ***REMOVED***
		if err = tarStream.Close(); err != nil ***REMOVED***
			t.Fatalf("Error when closing tar stream: %s", err)
		***REMOVED***
	***REMOVED***()

	if dockerfilePath == "" ***REMOVED*** // handled in BuildWithContext
		dockerfilePath = builder.DefaultDockerfileName
	***REMOVED***

	config := backend.BuildConfig***REMOVED***
		Options: &types.ImageBuildOptions***REMOVED***Dockerfile: dockerfilePath***REMOVED***,
		Source:  tarStream,
	***REMOVED***
	_, _, err = remotecontext.Detect(config)
	assert.EqualError(t, err, expectedError)
***REMOVED***

func TestCopyRunConfig(t *testing.T) ***REMOVED***
	defaultEnv := []string***REMOVED***"foo=1"***REMOVED***
	defaultCmd := []string***REMOVED***"old"***REMOVED***

	var testcases = []struct ***REMOVED***
		doc       string
		modifiers []runConfigModifier
		expected  *container.Config
	***REMOVED******REMOVED***
		***REMOVED***
			doc:       "Set the command",
			modifiers: []runConfigModifier***REMOVED***withCmd([]string***REMOVED***"new"***REMOVED***)***REMOVED***,
			expected: &container.Config***REMOVED***
				Cmd: []string***REMOVED***"new"***REMOVED***,
				Env: defaultEnv,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			doc:       "Set the command to a comment",
			modifiers: []runConfigModifier***REMOVED***withCmdComment("comment", runtime.GOOS)***REMOVED***,
			expected: &container.Config***REMOVED***
				Cmd: append(defaultShellForOS(runtime.GOOS), "#(nop) ", "comment"),
				Env: defaultEnv,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			doc: "Set the command and env",
			modifiers: []runConfigModifier***REMOVED***
				withCmd([]string***REMOVED***"new"***REMOVED***),
				withEnv([]string***REMOVED***"one", "two"***REMOVED***),
			***REMOVED***,
			expected: &container.Config***REMOVED***
				Cmd: []string***REMOVED***"new"***REMOVED***,
				Env: []string***REMOVED***"one", "two"***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for _, testcase := range testcases ***REMOVED***
		runConfig := &container.Config***REMOVED***
			Cmd: defaultCmd,
			Env: defaultEnv,
		***REMOVED***
		runConfigCopy := copyRunConfig(runConfig, testcase.modifiers...)
		assert.Equal(t, testcase.expected, runConfigCopy, testcase.doc)
		// Assert the original was not modified
		assert.NotEqual(t, runConfig, runConfigCopy, testcase.doc)
	***REMOVED***

***REMOVED***

func fullMutableRunConfig() *container.Config ***REMOVED***
	return &container.Config***REMOVED***
		Cmd: []string***REMOVED***"command", "arg1"***REMOVED***,
		Env: []string***REMOVED***"env1=foo", "env2=bar"***REMOVED***,
		ExposedPorts: nat.PortSet***REMOVED***
			"1000/tcp": ***REMOVED******REMOVED***,
			"1001/tcp": ***REMOVED******REMOVED***,
		***REMOVED***,
		Volumes: map[string]struct***REMOVED******REMOVED******REMOVED***
			"one": ***REMOVED******REMOVED***,
			"two": ***REMOVED******REMOVED***,
		***REMOVED***,
		Entrypoint: []string***REMOVED***"entry", "arg1"***REMOVED***,
		OnBuild:    []string***REMOVED***"first", "next"***REMOVED***,
		Labels: map[string]string***REMOVED***
			"label1": "value1",
			"label2": "value2",
		***REMOVED***,
		Shell: []string***REMOVED***"shell", "-c"***REMOVED***,
	***REMOVED***
***REMOVED***

func TestDeepCopyRunConfig(t *testing.T) ***REMOVED***
	runConfig := fullMutableRunConfig()
	copy := copyRunConfig(runConfig)
	assert.Equal(t, fullMutableRunConfig(), copy)

	copy.Cmd[1] = "arg2"
	copy.Env[1] = "env2=new"
	copy.ExposedPorts["10002"] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	copy.Volumes["three"] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	copy.Entrypoint[1] = "arg2"
	copy.OnBuild[0] = "start"
	copy.Labels["label3"] = "value3"
	copy.Shell[0] = "sh"
	assert.Equal(t, fullMutableRunConfig(), runConfig)
***REMOVED***
