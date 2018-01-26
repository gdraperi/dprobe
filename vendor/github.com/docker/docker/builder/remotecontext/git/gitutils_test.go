package git

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRemoteURL(t *testing.T) ***REMOVED***
	dir, err := parseRemoteURL("git://github.com/user/repo.git")
	require.NoError(t, err)
	assert.NotEmpty(t, dir)
	assert.Equal(t, gitRepo***REMOVED***"git://github.com/user/repo.git", "master", ""***REMOVED***, dir)

	dir, err = parseRemoteURL("git://github.com/user/repo.git#mybranch:mydir/mysubdir/")
	require.NoError(t, err)
	assert.NotEmpty(t, dir)
	assert.Equal(t, gitRepo***REMOVED***"git://github.com/user/repo.git", "mybranch", "mydir/mysubdir/"***REMOVED***, dir)

	dir, err = parseRemoteURL("https://github.com/user/repo.git")
	require.NoError(t, err)
	assert.NotEmpty(t, dir)
	assert.Equal(t, gitRepo***REMOVED***"https://github.com/user/repo.git", "master", ""***REMOVED***, dir)

	dir, err = parseRemoteURL("https://github.com/user/repo.git#mybranch:mydir/mysubdir/")
	require.NoError(t, err)
	assert.NotEmpty(t, dir)
	assert.Equal(t, gitRepo***REMOVED***"https://github.com/user/repo.git", "mybranch", "mydir/mysubdir/"***REMOVED***, dir)

	dir, err = parseRemoteURL("git@github.com:user/repo.git")
	require.NoError(t, err)
	assert.NotEmpty(t, dir)
	assert.Equal(t, gitRepo***REMOVED***"git@github.com:user/repo.git", "master", ""***REMOVED***, dir)

	dir, err = parseRemoteURL("git@github.com:user/repo.git#mybranch:mydir/mysubdir/")
	require.NoError(t, err)
	assert.NotEmpty(t, dir)
	assert.Equal(t, gitRepo***REMOVED***"git@github.com:user/repo.git", "mybranch", "mydir/mysubdir/"***REMOVED***, dir)
***REMOVED***

func TestCloneArgsSmartHttp(t *testing.T) ***REMOVED***
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	serverURL, _ := url.Parse(server.URL)

	serverURL.Path = "/repo.git"

	mux.HandleFunc("/repo.git/info/refs", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		q := r.URL.Query().Get("service")
		w.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-advertisement", q))
	***REMOVED***)

	args := fetchArgs(serverURL.String(), "master")
	exp := []string***REMOVED***"fetch", "--depth", "1", "origin", "master"***REMOVED***
	assert.Equal(t, exp, args)
***REMOVED***

func TestCloneArgsDumbHttp(t *testing.T) ***REMOVED***
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	serverURL, _ := url.Parse(server.URL)

	serverURL.Path = "/repo.git"

	mux.HandleFunc("/repo.git/info/refs", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Content-Type", "text/plain")
	***REMOVED***)

	args := fetchArgs(serverURL.String(), "master")
	exp := []string***REMOVED***"fetch", "origin", "master"***REMOVED***
	assert.Equal(t, exp, args)
***REMOVED***

func TestCloneArgsGit(t *testing.T) ***REMOVED***
	args := fetchArgs("git://github.com/docker/docker", "master")
	exp := []string***REMOVED***"fetch", "--depth", "1", "origin", "master"***REMOVED***
	assert.Equal(t, exp, args)
***REMOVED***

func gitGetConfig(name string) string ***REMOVED***
	b, err := git([]string***REMOVED***"config", "--get", name***REMOVED***...)
	if err != nil ***REMOVED***
		// since we are interested in empty or non empty string,
		// we can safely ignore the err here.
		return ""
	***REMOVED***
	return strings.TrimSpace(string(b))
***REMOVED***

func TestCheckoutGit(t *testing.T) ***REMOVED***
	root, err := ioutil.TempDir("", "docker-build-git-checkout")
	require.NoError(t, err)
	defer os.RemoveAll(root)

	autocrlf := gitGetConfig("core.autocrlf")
	if !(autocrlf == "true" || autocrlf == "false" ||
		autocrlf == "input" || autocrlf == "") ***REMOVED***
		t.Logf("unknown core.autocrlf value: \"%s\"", autocrlf)
	***REMOVED***
	eol := "\n"
	if autocrlf == "true" ***REMOVED***
		eol = "\r\n"
	***REMOVED***

	gitDir := filepath.Join(root, "repo")
	_, err = git("init", gitDir)
	require.NoError(t, err)

	_, err = gitWithinDir(gitDir, "config", "user.email", "test@docker.com")
	require.NoError(t, err)

	_, err = gitWithinDir(gitDir, "config", "user.name", "Docker test")
	require.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(gitDir, "Dockerfile"), []byte("FROM scratch"), 0644)
	require.NoError(t, err)

	subDir := filepath.Join(gitDir, "subdir")
	require.NoError(t, os.Mkdir(subDir, 0755))

	err = ioutil.WriteFile(filepath.Join(subDir, "Dockerfile"), []byte("FROM scratch\nEXPOSE 5000"), 0644)
	require.NoError(t, err)

	if runtime.GOOS != "windows" ***REMOVED***
		if err = os.Symlink("../subdir", filepath.Join(gitDir, "parentlink")); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		if err = os.Symlink("/subdir", filepath.Join(gitDir, "absolutelink")); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***

	_, err = gitWithinDir(gitDir, "add", "-A")
	require.NoError(t, err)

	_, err = gitWithinDir(gitDir, "commit", "-am", "First commit")
	require.NoError(t, err)

	_, err = gitWithinDir(gitDir, "checkout", "-b", "test")
	require.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(gitDir, "Dockerfile"), []byte("FROM scratch\nEXPOSE 3000"), 0644)
	require.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(subDir, "Dockerfile"), []byte("FROM busybox\nEXPOSE 5000"), 0644)
	require.NoError(t, err)

	_, err = gitWithinDir(gitDir, "add", "-A")
	require.NoError(t, err)

	_, err = gitWithinDir(gitDir, "commit", "-am", "Branch commit")
	require.NoError(t, err)

	_, err = gitWithinDir(gitDir, "checkout", "master")
	require.NoError(t, err)

	// set up submodule
	subrepoDir := filepath.Join(root, "subrepo")
	_, err = git("init", subrepoDir)
	require.NoError(t, err)

	_, err = gitWithinDir(subrepoDir, "config", "user.email", "test@docker.com")
	require.NoError(t, err)

	_, err = gitWithinDir(subrepoDir, "config", "user.name", "Docker test")
	require.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(subrepoDir, "subfile"), []byte("subcontents"), 0644)
	require.NoError(t, err)

	_, err = gitWithinDir(subrepoDir, "add", "-A")
	require.NoError(t, err)

	_, err = gitWithinDir(subrepoDir, "commit", "-am", "Subrepo initial")
	require.NoError(t, err)

	cmd := exec.Command("git", "submodule", "add", subrepoDir, "sub") // this command doesn't work with --work-tree
	cmd.Dir = gitDir
	require.NoError(t, cmd.Run())

	_, err = gitWithinDir(gitDir, "add", "-A")
	require.NoError(t, err)

	_, err = gitWithinDir(gitDir, "commit", "-am", "With submodule")
	require.NoError(t, err)

	type singleCase struct ***REMOVED***
		frag      string
		exp       string
		fail      bool
		submodule bool
	***REMOVED***

	cases := []singleCase***REMOVED***
		***REMOVED***"", "FROM scratch", false, true***REMOVED***,
		***REMOVED***"master", "FROM scratch", false, true***REMOVED***,
		***REMOVED***":subdir", "FROM scratch" + eol + "EXPOSE 5000", false, false***REMOVED***,
		***REMOVED***":nosubdir", "", true, false***REMOVED***,   // missing directory error
		***REMOVED***":Dockerfile", "", true, false***REMOVED***, // not a directory error
		***REMOVED***"master:nosubdir", "", true, false***REMOVED***,
		***REMOVED***"master:subdir", "FROM scratch" + eol + "EXPOSE 5000", false, false***REMOVED***,
		***REMOVED***"master:../subdir", "", true, false***REMOVED***,
		***REMOVED***"test", "FROM scratch" + eol + "EXPOSE 3000", false, false***REMOVED***,
		***REMOVED***"test:", "FROM scratch" + eol + "EXPOSE 3000", false, false***REMOVED***,
		***REMOVED***"test:subdir", "FROM busybox" + eol + "EXPOSE 5000", false, false***REMOVED***,
	***REMOVED***

	if runtime.GOOS != "windows" ***REMOVED***
		// Windows GIT (2.7.1 x64) does not support parentlink/absolutelink. Sample output below
		// 	git --work-tree .\repo --git-dir .\repo\.git add -A
		//	error: readlink("absolutelink"): Function not implemented
		// 	error: unable to index file absolutelink
		// 	fatal: adding files failed
		cases = append(cases, singleCase***REMOVED***frag: "master:absolutelink", exp: "FROM scratch" + eol + "EXPOSE 5000", fail: false***REMOVED***)
		cases = append(cases, singleCase***REMOVED***frag: "master:parentlink", exp: "FROM scratch" + eol + "EXPOSE 5000", fail: false***REMOVED***)
	***REMOVED***

	for _, c := range cases ***REMOVED***
		ref, subdir := getRefAndSubdir(c.frag)
		r, err := cloneGitRepo(gitRepo***REMOVED***remote: gitDir, ref: ref, subdir: subdir***REMOVED***)

		if c.fail ***REMOVED***
			assert.Error(t, err)
			continue
		***REMOVED***
		require.NoError(t, err)
		defer os.RemoveAll(r)
		if c.submodule ***REMOVED***
			b, err := ioutil.ReadFile(filepath.Join(r, "sub/subfile"))
			require.NoError(t, err)
			assert.Equal(t, "subcontents", string(b))
		***REMOVED*** else ***REMOVED***
			_, err := os.Stat(filepath.Join(r, "sub/subfile"))
			require.Error(t, err)
			require.True(t, os.IsNotExist(err))
		***REMOVED***

		b, err := ioutil.ReadFile(filepath.Join(r, "Dockerfile"))
		require.NoError(t, err)
		assert.Equal(t, c.exp, string(b))
	***REMOVED***
***REMOVED***

func TestValidGitTransport(t *testing.T) ***REMOVED***
	gitUrls := []string***REMOVED***
		"git://github.com/docker/docker",
		"git@github.com:docker/docker.git",
		"git@bitbucket.org:atlassianlabs/atlassian-docker.git",
		"https://github.com/docker/docker.git",
		"http://github.com/docker/docker.git",
		"http://github.com/docker/docker.git#branch",
		"http://github.com/docker/docker.git#:dir",
	***REMOVED***
	incompleteGitUrls := []string***REMOVED***
		"github.com/docker/docker",
	***REMOVED***

	for _, url := range gitUrls ***REMOVED***
		if !isGitTransport(url) ***REMOVED***
			t.Fatalf("%q should be detected as valid Git prefix", url)
		***REMOVED***
	***REMOVED***

	for _, url := range incompleteGitUrls ***REMOVED***
		if isGitTransport(url) ***REMOVED***
			t.Fatalf("%q should not be detected as valid Git prefix", url)
		***REMOVED***
	***REMOVED***
***REMOVED***
