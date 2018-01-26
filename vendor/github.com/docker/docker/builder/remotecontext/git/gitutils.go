package git

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/symlink"
	"github.com/docker/docker/pkg/urlutil"
	"github.com/pkg/errors"
)

type gitRepo struct ***REMOVED***
	remote string
	ref    string
	subdir string
***REMOVED***

// Clone clones a repository into a newly created directory which
// will be under "docker-build-git"
func Clone(remoteURL string) (string, error) ***REMOVED***
	repo, err := parseRemoteURL(remoteURL)

	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return cloneGitRepo(repo)
***REMOVED***

func cloneGitRepo(repo gitRepo) (checkoutDir string, err error) ***REMOVED***
	fetch := fetchArgs(repo.remote, repo.ref)

	root, err := ioutil.TempDir("", "docker-build-git")
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			os.RemoveAll(root)
		***REMOVED***
	***REMOVED***()

	if out, err := gitWithinDir(root, "init"); err != nil ***REMOVED***
		return "", errors.Wrapf(err, "failed to init repo at %s: %s", root, out)
	***REMOVED***

	// Add origin remote for compatibility with previous implementation that
	// used "git clone" and also to make sure local refs are created for branches
	if out, err := gitWithinDir(root, "remote", "add", "origin", repo.remote); err != nil ***REMOVED***
		return "", errors.Wrapf(err, "failed add origin repo at %s: %s", repo.remote, out)
	***REMOVED***

	if output, err := gitWithinDir(root, fetch...); err != nil ***REMOVED***
		return "", errors.Wrapf(err, "error fetching: %s", output)
	***REMOVED***

	checkoutDir, err = checkoutGit(root, repo.ref, repo.subdir)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	cmd := exec.Command("git", "submodule", "update", "--init", "--recursive", "--depth=1")
	cmd.Dir = root
	output, err := cmd.CombinedOutput()
	if err != nil ***REMOVED***
		return "", errors.Wrapf(err, "error initializing submodules: %s", output)
	***REMOVED***

	return checkoutDir, nil
***REMOVED***

func parseRemoteURL(remoteURL string) (gitRepo, error) ***REMOVED***
	repo := gitRepo***REMOVED******REMOVED***

	if !isGitTransport(remoteURL) ***REMOVED***
		remoteURL = "https://" + remoteURL
	***REMOVED***

	var fragment string
	if strings.HasPrefix(remoteURL, "git@") ***REMOVED***
		// git@.. is not an URL, so cannot be parsed as URL
		parts := strings.SplitN(remoteURL, "#", 2)

		repo.remote = parts[0]
		if len(parts) == 2 ***REMOVED***
			fragment = parts[1]
		***REMOVED***
		repo.ref, repo.subdir = getRefAndSubdir(fragment)
	***REMOVED*** else ***REMOVED***
		u, err := url.Parse(remoteURL)
		if err != nil ***REMOVED***
			return repo, err
		***REMOVED***

		repo.ref, repo.subdir = getRefAndSubdir(u.Fragment)
		u.Fragment = ""
		repo.remote = u.String()
	***REMOVED***
	return repo, nil
***REMOVED***

func getRefAndSubdir(fragment string) (ref string, subdir string) ***REMOVED***
	refAndDir := strings.SplitN(fragment, ":", 2)
	ref = "master"
	if len(refAndDir[0]) != 0 ***REMOVED***
		ref = refAndDir[0]
	***REMOVED***
	if len(refAndDir) > 1 && len(refAndDir[1]) != 0 ***REMOVED***
		subdir = refAndDir[1]
	***REMOVED***
	return
***REMOVED***

func fetchArgs(remoteURL string, ref string) []string ***REMOVED***
	args := []string***REMOVED***"fetch"***REMOVED***

	if supportsShallowClone(remoteURL) ***REMOVED***
		args = append(args, "--depth", "1")
	***REMOVED***

	return append(args, "origin", ref)
***REMOVED***

// Check if a given git URL supports a shallow git clone,
// i.e. it is a non-HTTP server or a smart HTTP server.
func supportsShallowClone(remoteURL string) bool ***REMOVED***
	if urlutil.IsURL(remoteURL) ***REMOVED***
		// Check if the HTTP server is smart

		// Smart servers must correctly respond to a query for the git-upload-pack service
		serviceURL := remoteURL + "/info/refs?service=git-upload-pack"

		// Try a HEAD request and fallback to a Get request on error
		res, err := http.Head(serviceURL)
		if err != nil || res.StatusCode != http.StatusOK ***REMOVED***
			res, err = http.Get(serviceURL)
			if err == nil ***REMOVED***
				res.Body.Close()
			***REMOVED***
			if err != nil || res.StatusCode != http.StatusOK ***REMOVED***
				// request failed
				return false
			***REMOVED***
		***REMOVED***

		if res.Header.Get("Content-Type") != "application/x-git-upload-pack-advertisement" ***REMOVED***
			// Fallback, not a smart server
			return false
		***REMOVED***
		return true
	***REMOVED***
	// Non-HTTP protocols always support shallow clones
	return true
***REMOVED***

func checkoutGit(root, ref, subdir string) (string, error) ***REMOVED***
	// Try checking out by ref name first. This will work on branches and sets
	// .git/HEAD to the current branch name
	if output, err := gitWithinDir(root, "checkout", ref); err != nil ***REMOVED***
		// If checking out by branch name fails check out the last fetched ref
		if _, err2 := gitWithinDir(root, "checkout", "FETCH_HEAD"); err2 != nil ***REMOVED***
			return "", errors.Wrapf(err, "error checking out %s: %s", ref, output)
		***REMOVED***
	***REMOVED***

	if subdir != "" ***REMOVED***
		newCtx, err := symlink.FollowSymlinkInScope(filepath.Join(root, subdir), root)
		if err != nil ***REMOVED***
			return "", errors.Wrapf(err, "error setting git context, %q not within git root", subdir)
		***REMOVED***

		fi, err := os.Stat(newCtx)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		if !fi.IsDir() ***REMOVED***
			return "", errors.Errorf("error setting git context, not a directory: %s", newCtx)
		***REMOVED***
		root = newCtx
	***REMOVED***

	return root, nil
***REMOVED***

func gitWithinDir(dir string, args ...string) ([]byte, error) ***REMOVED***
	a := []string***REMOVED***"--work-tree", dir, "--git-dir", filepath.Join(dir, ".git")***REMOVED***
	return git(append(a, args...)...)
***REMOVED***

func git(args ...string) ([]byte, error) ***REMOVED***
	return exec.Command("git", args...).CombinedOutput()
***REMOVED***

// isGitTransport returns true if the provided str is a git transport by inspecting
// the prefix of the string for known protocols used in git.
func isGitTransport(str string) bool ***REMOVED***
	return urlutil.IsURL(str) || strings.HasPrefix(str, "git://") || strings.HasPrefix(str, "git@")
***REMOVED***
