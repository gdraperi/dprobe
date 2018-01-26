package fakegit

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/docker/docker/integration-cli/cli/build/fakecontext"
	"github.com/docker/docker/integration-cli/cli/build/fakestorage"
	"github.com/stretchr/testify/require"
)

type testingT interface ***REMOVED***
	require.TestingT
	logT
	Fatal(args ...interface***REMOVED******REMOVED***)
	Fatalf(string, ...interface***REMOVED******REMOVED***)
***REMOVED***

type logT interface ***REMOVED***
	Logf(string, ...interface***REMOVED******REMOVED***)
***REMOVED***

type gitServer interface ***REMOVED***
	URL() string
	Close() error
***REMOVED***

type localGitServer struct ***REMOVED***
	*httptest.Server
***REMOVED***

func (r *localGitServer) Close() error ***REMOVED***
	r.Server.Close()
	return nil
***REMOVED***

func (r *localGitServer) URL() string ***REMOVED***
	return r.Server.URL
***REMOVED***

// FakeGit is a fake git server
type FakeGit struct ***REMOVED***
	root    string
	server  gitServer
	RepoURL string
***REMOVED***

// Close closes the server, implements Closer interface
func (g *FakeGit) Close() ***REMOVED***
	g.server.Close()
	os.RemoveAll(g.root)
***REMOVED***

// New create a fake git server that can be used for git related tests
func New(c testingT, name string, files map[string]string, enforceLocalServer bool) *FakeGit ***REMOVED***
	ctx := fakecontext.New(c, "", fakecontext.WithFiles(files))
	defer ctx.Close()
	curdir, err := os.Getwd()
	if err != nil ***REMOVED***
		c.Fatal(err)
	***REMOVED***
	defer os.Chdir(curdir)

	if output, err := exec.Command("git", "init", ctx.Dir).CombinedOutput(); err != nil ***REMOVED***
		c.Fatalf("error trying to init repo: %s (%s)", err, output)
	***REMOVED***
	err = os.Chdir(ctx.Dir)
	if err != nil ***REMOVED***
		c.Fatal(err)
	***REMOVED***
	if output, err := exec.Command("git", "config", "user.name", "Fake User").CombinedOutput(); err != nil ***REMOVED***
		c.Fatalf("error trying to set 'user.name': %s (%s)", err, output)
	***REMOVED***
	if output, err := exec.Command("git", "config", "user.email", "fake.user@example.com").CombinedOutput(); err != nil ***REMOVED***
		c.Fatalf("error trying to set 'user.email': %s (%s)", err, output)
	***REMOVED***
	if output, err := exec.Command("git", "add", "*").CombinedOutput(); err != nil ***REMOVED***
		c.Fatalf("error trying to add files to repo: %s (%s)", err, output)
	***REMOVED***
	if output, err := exec.Command("git", "commit", "-a", "-m", "Initial commit").CombinedOutput(); err != nil ***REMOVED***
		c.Fatalf("error trying to commit to repo: %s (%s)", err, output)
	***REMOVED***

	root, err := ioutil.TempDir("", "docker-test-git-repo")
	if err != nil ***REMOVED***
		c.Fatal(err)
	***REMOVED***
	repoPath := filepath.Join(root, name+".git")
	if output, err := exec.Command("git", "clone", "--bare", ctx.Dir, repoPath).CombinedOutput(); err != nil ***REMOVED***
		os.RemoveAll(root)
		c.Fatalf("error trying to clone --bare: %s (%s)", err, output)
	***REMOVED***
	err = os.Chdir(repoPath)
	if err != nil ***REMOVED***
		os.RemoveAll(root)
		c.Fatal(err)
	***REMOVED***
	if output, err := exec.Command("git", "update-server-info").CombinedOutput(); err != nil ***REMOVED***
		os.RemoveAll(root)
		c.Fatalf("error trying to git update-server-info: %s (%s)", err, output)
	***REMOVED***
	err = os.Chdir(curdir)
	if err != nil ***REMOVED***
		os.RemoveAll(root)
		c.Fatal(err)
	***REMOVED***

	var server gitServer
	if !enforceLocalServer ***REMOVED***
		// use fakeStorage server, which might be local or remote (at test daemon)
		server = fakestorage.New(c, root)
	***REMOVED*** else ***REMOVED***
		// always start a local http server on CLI test machine
		httpServer := httptest.NewServer(http.FileServer(http.Dir(root)))
		server = &localGitServer***REMOVED***httpServer***REMOVED***
	***REMOVED***
	return &FakeGit***REMOVED***
		root:    root,
		server:  server,
		RepoURL: fmt.Sprintf("%s/%s.git", server.URL(), name),
	***REMOVED***
***REMOVED***
