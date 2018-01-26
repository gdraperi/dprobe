package remotecontext

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"testing"

	"github.com/docker/docker/builder"
	"github.com/docker/docker/pkg/containerfs"
)

const (
	dockerfileContents   = "FROM busybox"
	dockerignoreFilename = ".dockerignore"
	testfileContents     = "test"
)

const shouldStayFilename = "should_stay"

func extractFilenames(files []os.FileInfo) []string ***REMOVED***
	filenames := make([]string, len(files))

	for i, file := range files ***REMOVED***
		filenames[i] = file.Name()
	***REMOVED***

	return filenames
***REMOVED***

func checkDirectory(t *testing.T, dir string, expectedFiles []string) ***REMOVED***
	files, err := ioutil.ReadDir(dir)

	if err != nil ***REMOVED***
		t.Fatalf("Could not read directory: %s", err)
	***REMOVED***

	if len(files) != len(expectedFiles) ***REMOVED***
		log.Fatalf("Directory should contain exactly %d file(s), got %d", len(expectedFiles), len(files))
	***REMOVED***

	filenames := extractFilenames(files)
	sort.Strings(filenames)
	sort.Strings(expectedFiles)

	for i, filename := range filenames ***REMOVED***
		if filename != expectedFiles[i] ***REMOVED***
			t.Fatalf("File %s should be in the directory, got: %s", expectedFiles[i], filename)
		***REMOVED***
	***REMOVED***
***REMOVED***

func executeProcess(t *testing.T, contextDir string) ***REMOVED***
	modifiableCtx := &stubRemote***REMOVED***root: containerfs.NewLocalContainerFS(contextDir)***REMOVED***

	err := removeDockerfile(modifiableCtx, builder.DefaultDockerfileName)

	if err != nil ***REMOVED***
		t.Fatalf("Error when executing Process: %s", err)
	***REMOVED***
***REMOVED***

func TestProcessShouldRemoveDockerfileDockerignore(t *testing.T) ***REMOVED***
	contextDir, cleanup := createTestTempDir(t, "", "builder-dockerignore-process-test")
	defer cleanup()

	createTestTempFile(t, contextDir, shouldStayFilename, testfileContents, 0777)
	createTestTempFile(t, contextDir, dockerignoreFilename, "Dockerfile\n.dockerignore", 0777)
	createTestTempFile(t, contextDir, builder.DefaultDockerfileName, dockerfileContents, 0777)

	executeProcess(t, contextDir)

	checkDirectory(t, contextDir, []string***REMOVED***shouldStayFilename***REMOVED***)

***REMOVED***

func TestProcessNoDockerignore(t *testing.T) ***REMOVED***
	contextDir, cleanup := createTestTempDir(t, "", "builder-dockerignore-process-test")
	defer cleanup()

	createTestTempFile(t, contextDir, shouldStayFilename, testfileContents, 0777)
	createTestTempFile(t, contextDir, builder.DefaultDockerfileName, dockerfileContents, 0777)

	executeProcess(t, contextDir)

	checkDirectory(t, contextDir, []string***REMOVED***shouldStayFilename, builder.DefaultDockerfileName***REMOVED***)

***REMOVED***

func TestProcessShouldLeaveAllFiles(t *testing.T) ***REMOVED***
	contextDir, cleanup := createTestTempDir(t, "", "builder-dockerignore-process-test")
	defer cleanup()

	createTestTempFile(t, contextDir, shouldStayFilename, testfileContents, 0777)
	createTestTempFile(t, contextDir, builder.DefaultDockerfileName, dockerfileContents, 0777)
	createTestTempFile(t, contextDir, dockerignoreFilename, "input1\ninput2", 0777)

	executeProcess(t, contextDir)

	checkDirectory(t, contextDir, []string***REMOVED***shouldStayFilename, builder.DefaultDockerfileName, dockerignoreFilename***REMOVED***)

***REMOVED***

// TODO: remove after moving to a separate pkg
type stubRemote struct ***REMOVED***
	root containerfs.ContainerFS
***REMOVED***

func (r *stubRemote) Hash(path string) (string, error) ***REMOVED***
	return "", errors.New("not implemented")
***REMOVED***

func (r *stubRemote) Root() containerfs.ContainerFS ***REMOVED***
	return r.root
***REMOVED***
func (r *stubRemote) Close() error ***REMOVED***
	return errors.New("not implemented")
***REMOVED***
func (r *stubRemote) Remove(p string) error ***REMOVED***
	return r.root.Remove(r.root.Join(r.root.Path(), p))
***REMOVED***
