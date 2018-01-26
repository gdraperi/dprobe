package dockerfile

import (
	"testing"

	"github.com/docker/docker/builder/dockerfile/instructions"
	"github.com/docker/docker/builder/remotecontext"
	"github.com/docker/docker/internal/testutil"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/reexec"
)

type dispatchTestCase struct ***REMOVED***
	name, expectedError string
	cmd                 instructions.Command
	files               map[string]string
***REMOVED***

func init() ***REMOVED***
	reexec.Init()
***REMOVED***

func initDispatchTestCases() []dispatchTestCase ***REMOVED***
	dispatchTestCases := []dispatchTestCase***REMOVED***
		***REMOVED***
			name: "ADD multiple files to file",
			cmd: &instructions.AddCommand***REMOVED***SourcesAndDest: instructions.SourcesAndDest***REMOVED***
				"file1.txt",
				"file2.txt",
				"test",
			***REMOVED******REMOVED***,
			expectedError: "When using ADD with more than one source file, the destination must be a directory and end with a /",
			files:         map[string]string***REMOVED***"file1.txt": "test1", "file2.txt": "test2"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "Wildcard ADD multiple files to file",
			cmd: &instructions.AddCommand***REMOVED***SourcesAndDest: instructions.SourcesAndDest***REMOVED***
				"file*.txt",
				"test",
			***REMOVED******REMOVED***,
			expectedError: "When using ADD with more than one source file, the destination must be a directory and end with a /",
			files:         map[string]string***REMOVED***"file1.txt": "test1", "file2.txt": "test2"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "COPY multiple files to file",
			cmd: &instructions.CopyCommand***REMOVED***SourcesAndDest: instructions.SourcesAndDest***REMOVED***
				"file1.txt",
				"file2.txt",
				"test",
			***REMOVED******REMOVED***,
			expectedError: "When using COPY with more than one source file, the destination must be a directory and end with a /",
			files:         map[string]string***REMOVED***"file1.txt": "test1", "file2.txt": "test2"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "ADD multiple files to file with whitespace",
			cmd: &instructions.AddCommand***REMOVED***SourcesAndDest: instructions.SourcesAndDest***REMOVED***
				"test file1.txt",
				"test file2.txt",
				"test",
			***REMOVED******REMOVED***,
			expectedError: "When using ADD with more than one source file, the destination must be a directory and end with a /",
			files:         map[string]string***REMOVED***"test file1.txt": "test1", "test file2.txt": "test2"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "COPY multiple files to file with whitespace",
			cmd: &instructions.CopyCommand***REMOVED***SourcesAndDest: instructions.SourcesAndDest***REMOVED***
				"test file1.txt",
				"test file2.txt",
				"test",
			***REMOVED******REMOVED***,
			expectedError: "When using COPY with more than one source file, the destination must be a directory and end with a /",
			files:         map[string]string***REMOVED***"test file1.txt": "test1", "test file2.txt": "test2"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			name: "COPY wildcard no files",
			cmd: &instructions.CopyCommand***REMOVED***SourcesAndDest: instructions.SourcesAndDest***REMOVED***
				"file*.txt",
				"/tmp/",
			***REMOVED******REMOVED***,
			expectedError: "COPY failed: no source files were specified",
			files:         nil,
		***REMOVED***,
		***REMOVED***
			name: "COPY url",
			cmd: &instructions.CopyCommand***REMOVED***SourcesAndDest: instructions.SourcesAndDest***REMOVED***
				"https://index.docker.io/robots.txt",
				"/",
			***REMOVED******REMOVED***,
			expectedError: "source can't be a URL for COPY",
			files:         nil,
		***REMOVED******REMOVED***

	return dispatchTestCases
***REMOVED***

func TestDispatch(t *testing.T) ***REMOVED***
	testCases := initDispatchTestCases()

	for _, testCase := range testCases ***REMOVED***
		executeTestCase(t, testCase)
	***REMOVED***
***REMOVED***

func executeTestCase(t *testing.T, testCase dispatchTestCase) ***REMOVED***
	contextDir, cleanup := createTestTempDir(t, "", "builder-dockerfile-test")
	defer cleanup()

	for filename, content := range testCase.files ***REMOVED***
		createTestTempFile(t, contextDir, filename, content, 0777)
	***REMOVED***

	tarStream, err := archive.Tar(contextDir, archive.Uncompressed)

	if err != nil ***REMOVED***
		t.Fatalf("Error when creating tar stream: %s", err)
	***REMOVED***

	defer func() ***REMOVED***
		if err = tarStream.Close(); err != nil ***REMOVED***
			t.Fatalf("Error when closing tar stream: %s", err)
		***REMOVED***
	***REMOVED***()

	context, err := remotecontext.FromArchive(tarStream)

	if err != nil ***REMOVED***
		t.Fatalf("Error when creating tar context: %s", err)
	***REMOVED***

	defer func() ***REMOVED***
		if err = context.Close(); err != nil ***REMOVED***
			t.Fatalf("Error when closing tar context: %s", err)
		***REMOVED***
	***REMOVED***()

	b := newBuilderWithMockBackend()
	sb := newDispatchRequest(b, '`', context, newBuildArgs(make(map[string]*string)), newStagesBuildResults())
	err = dispatch(sb, testCase.cmd)
	testutil.ErrorContains(t, err, testCase.expectedError)
***REMOVED***
