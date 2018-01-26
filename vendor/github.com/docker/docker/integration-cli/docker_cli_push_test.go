package main

import (
	"archive/tar"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/cli/config"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/cli"
	"github.com/docker/docker/integration-cli/cli/build"
	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
)

// Pushing an image to a private registry.
func testPushBusyboxImage(c *check.C) ***REMOVED***
	repoName := fmt.Sprintf("%v/dockercli/busybox", privateRegistryURL)
	// tag the image to upload it to the private registry
	dockerCmd(c, "tag", "busybox", repoName)
	// push the image to the registry
	dockerCmd(c, "push", repoName)
***REMOVED***

func (s *DockerRegistrySuite) TestPushBusyboxImage(c *check.C) ***REMOVED***
	testPushBusyboxImage(c)
***REMOVED***

func (s *DockerSchema1RegistrySuite) TestPushBusyboxImage(c *check.C) ***REMOVED***
	testPushBusyboxImage(c)
***REMOVED***

// pushing an image without a prefix should throw an error
func (s *DockerSuite) TestPushUnprefixedRepo(c *check.C) ***REMOVED***
	out, _, err := dockerCmdWithError("push", "busybox")
	c.Assert(err, check.NotNil, check.Commentf("pushing an unprefixed repo didn't result in a non-zero exit status: %s", out))
***REMOVED***

func testPushUntagged(c *check.C) ***REMOVED***
	repoName := fmt.Sprintf("%v/dockercli/busybox", privateRegistryURL)
	expected := "An image does not exist locally with the tag"

	out, _, err := dockerCmdWithError("push", repoName)
	c.Assert(err, check.NotNil, check.Commentf("pushing the image to the private registry should have failed: output %q", out))
	c.Assert(out, checker.Contains, expected, check.Commentf("pushing the image failed"))
***REMOVED***

func (s *DockerRegistrySuite) TestPushUntagged(c *check.C) ***REMOVED***
	testPushUntagged(c)
***REMOVED***

func (s *DockerSchema1RegistrySuite) TestPushUntagged(c *check.C) ***REMOVED***
	testPushUntagged(c)
***REMOVED***

func testPushBadTag(c *check.C) ***REMOVED***
	repoName := fmt.Sprintf("%v/dockercli/busybox:latest", privateRegistryURL)
	expected := "does not exist"

	out, _, err := dockerCmdWithError("push", repoName)
	c.Assert(err, check.NotNil, check.Commentf("pushing the image to the private registry should have failed: output %q", out))
	c.Assert(out, checker.Contains, expected, check.Commentf("pushing the image failed"))
***REMOVED***

func (s *DockerRegistrySuite) TestPushBadTag(c *check.C) ***REMOVED***
	testPushBadTag(c)
***REMOVED***

func (s *DockerSchema1RegistrySuite) TestPushBadTag(c *check.C) ***REMOVED***
	testPushBadTag(c)
***REMOVED***

func testPushMultipleTags(c *check.C) ***REMOVED***
	repoName := fmt.Sprintf("%v/dockercli/busybox", privateRegistryURL)
	repoTag1 := fmt.Sprintf("%v/dockercli/busybox:t1", privateRegistryURL)
	repoTag2 := fmt.Sprintf("%v/dockercli/busybox:t2", privateRegistryURL)
	// tag the image and upload it to the private registry
	dockerCmd(c, "tag", "busybox", repoTag1)

	dockerCmd(c, "tag", "busybox", repoTag2)

	dockerCmd(c, "push", repoName)

	// Ensure layer list is equivalent for repoTag1 and repoTag2
	out1, _ := dockerCmd(c, "pull", repoTag1)

	imageAlreadyExists := ": Image already exists"
	var out1Lines []string
	for _, outputLine := range strings.Split(out1, "\n") ***REMOVED***
		if strings.Contains(outputLine, imageAlreadyExists) ***REMOVED***
			out1Lines = append(out1Lines, outputLine)
		***REMOVED***
	***REMOVED***

	out2, _ := dockerCmd(c, "pull", repoTag2)

	var out2Lines []string
	for _, outputLine := range strings.Split(out2, "\n") ***REMOVED***
		if strings.Contains(outputLine, imageAlreadyExists) ***REMOVED***
			out1Lines = append(out1Lines, outputLine)
		***REMOVED***
	***REMOVED***
	c.Assert(out2Lines, checker.HasLen, len(out1Lines))

	for i := range out1Lines ***REMOVED***
		c.Assert(out1Lines[i], checker.Equals, out2Lines[i])
	***REMOVED***
***REMOVED***

func (s *DockerRegistrySuite) TestPushMultipleTags(c *check.C) ***REMOVED***
	testPushMultipleTags(c)
***REMOVED***

func (s *DockerSchema1RegistrySuite) TestPushMultipleTags(c *check.C) ***REMOVED***
	testPushMultipleTags(c)
***REMOVED***

func testPushEmptyLayer(c *check.C) ***REMOVED***
	repoName := fmt.Sprintf("%v/dockercli/emptylayer", privateRegistryURL)
	emptyTarball, err := ioutil.TempFile("", "empty_tarball")
	c.Assert(err, check.IsNil, check.Commentf("Unable to create test file"))

	tw := tar.NewWriter(emptyTarball)
	err = tw.Close()
	c.Assert(err, check.IsNil, check.Commentf("Error creating empty tarball"))

	freader, err := os.Open(emptyTarball.Name())
	c.Assert(err, check.IsNil, check.Commentf("Could not open test tarball"))
	defer freader.Close()

	icmd.RunCmd(icmd.Cmd***REMOVED***
		Command: []string***REMOVED***dockerBinary, "import", "-", repoName***REMOVED***,
		Stdin:   freader,
	***REMOVED***).Assert(c, icmd.Success)

	// Now verify we can push it
	out, _, err := dockerCmdWithError("push", repoName)
	c.Assert(err, check.IsNil, check.Commentf("pushing the image to the private registry has failed: %s", out))
***REMOVED***

func (s *DockerRegistrySuite) TestPushEmptyLayer(c *check.C) ***REMOVED***
	testPushEmptyLayer(c)
***REMOVED***

func (s *DockerSchema1RegistrySuite) TestPushEmptyLayer(c *check.C) ***REMOVED***
	testPushEmptyLayer(c)
***REMOVED***

// testConcurrentPush pushes multiple tags to the same repo
// concurrently.
func testConcurrentPush(c *check.C) ***REMOVED***
	repoName := fmt.Sprintf("%v/dockercli/busybox", privateRegistryURL)

	repos := []string***REMOVED******REMOVED***
	for _, tag := range []string***REMOVED***"push1", "push2", "push3"***REMOVED*** ***REMOVED***
		repo := fmt.Sprintf("%v:%v", repoName, tag)
		buildImageSuccessfully(c, repo, build.WithDockerfile(fmt.Sprintf(`
	FROM busybox
	ENTRYPOINT ["/bin/echo"]
	ENV FOO foo
	ENV BAR bar
	CMD echo %s
`, repo)))
		repos = append(repos, repo)
	***REMOVED***

	// Push tags, in parallel
	results := make(chan error)

	for _, repo := range repos ***REMOVED***
		go func(repo string) ***REMOVED***
			result := icmd.RunCommand(dockerBinary, "push", repo)
			results <- result.Error
		***REMOVED***(repo)
	***REMOVED***

	for range repos ***REMOVED***
		err := <-results
		c.Assert(err, checker.IsNil, check.Commentf("concurrent push failed with error: %v", err))
	***REMOVED***

	// Clear local images store.
	args := append([]string***REMOVED***"rmi"***REMOVED***, repos...)
	dockerCmd(c, args...)

	// Re-pull and run individual tags, to make sure pushes succeeded
	for _, repo := range repos ***REMOVED***
		dockerCmd(c, "pull", repo)
		dockerCmd(c, "inspect", repo)
		out, _ := dockerCmd(c, "run", "--rm", repo)
		c.Assert(strings.TrimSpace(out), checker.Equals, "/bin/sh -c echo "+repo)
	***REMOVED***
***REMOVED***

func (s *DockerRegistrySuite) TestConcurrentPush(c *check.C) ***REMOVED***
	testConcurrentPush(c)
***REMOVED***

func (s *DockerSchema1RegistrySuite) TestConcurrentPush(c *check.C) ***REMOVED***
	testConcurrentPush(c)
***REMOVED***

func (s *DockerRegistrySuite) TestCrossRepositoryLayerPush(c *check.C) ***REMOVED***
	sourceRepoName := fmt.Sprintf("%v/dockercli/busybox", privateRegistryURL)
	// tag the image to upload it to the private registry
	dockerCmd(c, "tag", "busybox", sourceRepoName)
	// push the image to the registry
	out1, _, err := dockerCmdWithError("push", sourceRepoName)
	c.Assert(err, check.IsNil, check.Commentf("pushing the image to the private registry has failed: %s", out1))
	// ensure that none of the layers were mounted from another repository during push
	c.Assert(strings.Contains(out1, "Mounted from"), check.Equals, false)

	digest1 := reference.DigestRegexp.FindString(out1)
	c.Assert(len(digest1), checker.GreaterThan, 0, check.Commentf("no digest found for pushed manifest"))

	destRepoName := fmt.Sprintf("%v/dockercli/crossrepopush", privateRegistryURL)
	// retag the image to upload the same layers to another repo in the same registry
	dockerCmd(c, "tag", "busybox", destRepoName)
	// push the image to the registry
	out2, _, err := dockerCmdWithError("push", destRepoName)
	c.Assert(err, check.IsNil, check.Commentf("pushing the image to the private registry has failed: %s", out2))
	// ensure that layers were mounted from the first repo during push
	c.Assert(strings.Contains(out2, "Mounted from dockercli/busybox"), check.Equals, true)

	digest2 := reference.DigestRegexp.FindString(out2)
	c.Assert(len(digest2), checker.GreaterThan, 0, check.Commentf("no digest found for pushed manifest"))
	c.Assert(digest1, check.Equals, digest2)

	// ensure that pushing again produces the same digest
	out3, _, err := dockerCmdWithError("push", destRepoName)
	c.Assert(err, check.IsNil, check.Commentf("pushing the image to the private registry has failed: %s", out2))

	digest3 := reference.DigestRegexp.FindString(out3)
	c.Assert(len(digest2), checker.GreaterThan, 0, check.Commentf("no digest found for pushed manifest"))
	c.Assert(digest3, check.Equals, digest2)

	// ensure that we can pull and run the cross-repo-pushed repository
	dockerCmd(c, "rmi", destRepoName)
	dockerCmd(c, "pull", destRepoName)
	out4, _ := dockerCmd(c, "run", destRepoName, "echo", "-n", "hello world")
	c.Assert(out4, check.Equals, "hello world")
***REMOVED***

func (s *DockerSchema1RegistrySuite) TestCrossRepositoryLayerPushNotSupported(c *check.C) ***REMOVED***
	sourceRepoName := fmt.Sprintf("%v/dockercli/busybox", privateRegistryURL)
	// tag the image to upload it to the private registry
	dockerCmd(c, "tag", "busybox", sourceRepoName)
	// push the image to the registry
	out1, _, err := dockerCmdWithError("push", sourceRepoName)
	c.Assert(err, check.IsNil, check.Commentf("pushing the image to the private registry has failed: %s", out1))
	// ensure that none of the layers were mounted from another repository during push
	c.Assert(strings.Contains(out1, "Mounted from"), check.Equals, false)

	digest1 := reference.DigestRegexp.FindString(out1)
	c.Assert(len(digest1), checker.GreaterThan, 0, check.Commentf("no digest found for pushed manifest"))

	destRepoName := fmt.Sprintf("%v/dockercli/crossrepopush", privateRegistryURL)
	// retag the image to upload the same layers to another repo in the same registry
	dockerCmd(c, "tag", "busybox", destRepoName)
	// push the image to the registry
	out2, _, err := dockerCmdWithError("push", destRepoName)
	c.Assert(err, check.IsNil, check.Commentf("pushing the image to the private registry has failed: %s", out2))
	// schema1 registry should not support cross-repo layer mounts, so ensure that this does not happen
	c.Assert(strings.Contains(out2, "Mounted from"), check.Equals, false)

	digest2 := reference.DigestRegexp.FindString(out2)
	c.Assert(len(digest2), checker.GreaterThan, 0, check.Commentf("no digest found for pushed manifest"))
	c.Assert(digest1, check.Not(check.Equals), digest2)

	// ensure that we can pull and run the second pushed repository
	dockerCmd(c, "rmi", destRepoName)
	dockerCmd(c, "pull", destRepoName)
	out3, _ := dockerCmd(c, "run", destRepoName, "echo", "-n", "hello world")
	c.Assert(out3, check.Equals, "hello world")
***REMOVED***

func (s *DockerTrustSuite) TestTrustedPush(c *check.C) ***REMOVED***
	repoName := fmt.Sprintf("%v/dockerclitrusted/pushtest:latest", privateRegistryURL)
	// tag the image and upload it to the private registry
	cli.DockerCmd(c, "tag", "busybox", repoName)

	cli.Docker(cli.Args("push", repoName), trustedCmd).Assert(c, SuccessSigningAndPushing)

	// Try pull after push
	cli.Docker(cli.Args("pull", repoName), trustedCmd).Assert(c, icmd.Expected***REMOVED***
		Out: "Status: Image is up to date",
	***REMOVED***)

	// Assert that we rotated the snapshot key to the server by checking our local keystore
	contents, err := ioutil.ReadDir(filepath.Join(config.Dir(), "trust/private/tuf_keys", privateRegistryURL, "dockerclitrusted/pushtest"))
	c.Assert(err, check.IsNil, check.Commentf("Unable to read local tuf key files"))
	// Check that we only have 1 key (targets key)
	c.Assert(contents, checker.HasLen, 1)
***REMOVED***

func (s *DockerTrustSuite) TestTrustedPushWithEnvPasswords(c *check.C) ***REMOVED***
	repoName := fmt.Sprintf("%v/dockerclienv/trusted:latest", privateRegistryURL)
	// tag the image and upload it to the private registry
	cli.DockerCmd(c, "tag", "busybox", repoName)

	cli.Docker(cli.Args("push", repoName), trustedCmdWithPassphrases("12345678", "12345678")).Assert(c, SuccessSigningAndPushing)

	// Try pull after push
	cli.Docker(cli.Args("pull", repoName), trustedCmd).Assert(c, icmd.Expected***REMOVED***
		Out: "Status: Image is up to date",
	***REMOVED***)
***REMOVED***

func (s *DockerTrustSuite) TestTrustedPushWithFailingServer(c *check.C) ***REMOVED***
	repoName := fmt.Sprintf("%v/dockerclitrusted/failingserver:latest", privateRegistryURL)
	// tag the image and upload it to the private registry
	cli.DockerCmd(c, "tag", "busybox", repoName)

	// Using a name that doesn't resolve to an address makes this test faster
	cli.Docker(cli.Args("push", repoName), trustedCmdWithServer("https://server.invalid:81/")).Assert(c, icmd.Expected***REMOVED***
		ExitCode: 1,
		Err:      "error contacting notary server",
	***REMOVED***)
***REMOVED***

func (s *DockerTrustSuite) TestTrustedPushWithoutServerAndUntrusted(c *check.C) ***REMOVED***
	repoName := fmt.Sprintf("%v/dockerclitrusted/trustedandnot:latest", privateRegistryURL)
	// tag the image and upload it to the private registry
	cli.DockerCmd(c, "tag", "busybox", repoName)

	result := cli.Docker(cli.Args("push", "--disable-content-trust", repoName), trustedCmdWithServer("https://server.invalid:81/"))
	result.Assert(c, icmd.Success)
	c.Assert(result.Combined(), check.Not(checker.Contains), "Error establishing connection to notary repository", check.Commentf("Missing expected output on trusted push with --disable-content-trust:"))
***REMOVED***

func (s *DockerTrustSuite) TestTrustedPushWithExistingTag(c *check.C) ***REMOVED***
	repoName := fmt.Sprintf("%v/dockerclitag/trusted:latest", privateRegistryURL)
	// tag the image and upload it to the private registry
	cli.DockerCmd(c, "tag", "busybox", repoName)
	cli.DockerCmd(c, "push", repoName)

	cli.Docker(cli.Args("push", repoName), trustedCmd).Assert(c, SuccessSigningAndPushing)

	// Try pull after push
	cli.Docker(cli.Args("pull", repoName), trustedCmd).Assert(c, icmd.Expected***REMOVED***
		Out: "Status: Image is up to date",
	***REMOVED***)
***REMOVED***

func (s *DockerTrustSuite) TestTrustedPushWithExistingSignedTag(c *check.C) ***REMOVED***
	repoName := fmt.Sprintf("%v/dockerclipushpush/trusted:latest", privateRegistryURL)
	// tag the image and upload it to the private registry
	cli.DockerCmd(c, "tag", "busybox", repoName)

	// Do a trusted push
	cli.Docker(cli.Args("push", repoName), trustedCmd).Assert(c, SuccessSigningAndPushing)

	// Do another trusted push
	cli.Docker(cli.Args("push", repoName), trustedCmd).Assert(c, SuccessSigningAndPushing)
	cli.DockerCmd(c, "rmi", repoName)

	// Try pull to ensure the double push did not break our ability to pull
	cli.Docker(cli.Args("pull", repoName), trustedCmd).Assert(c, SuccessDownloaded)
***REMOVED***

func (s *DockerTrustSuite) TestTrustedPushWithIncorrectPassphraseForNonRoot(c *check.C) ***REMOVED***
	repoName := fmt.Sprintf("%v/dockercliincorretpwd/trusted:latest", privateRegistryURL)
	// tag the image and upload it to the private registry
	cli.DockerCmd(c, "tag", "busybox", repoName)

	// Push with default passphrases
	cli.Docker(cli.Args("push", repoName), trustedCmd).Assert(c, SuccessSigningAndPushing)

	// Push with wrong passphrases
	cli.Docker(cli.Args("push", repoName), trustedCmdWithPassphrases("12345678", "87654321")).Assert(c, icmd.Expected***REMOVED***
		ExitCode: 1,
		Err:      "could not find necessary signing keys",
	***REMOVED***)
***REMOVED***

func (s *DockerTrustSuite) TestTrustedPushWithReleasesDelegationOnly(c *check.C) ***REMOVED***
	testRequires(c, NotaryHosting)
	repoName := fmt.Sprintf("%v/dockerclireleasedelegationinitfirst/trusted", privateRegistryURL)
	targetName := fmt.Sprintf("%s:latest", repoName)
	s.notaryInitRepo(c, repoName)
	s.notaryCreateDelegation(c, repoName, "targets/releases", s.not.keys[0].Public)
	s.notaryPublish(c, repoName)

	s.notaryImportKey(c, repoName, "targets/releases", s.not.keys[0].Private)

	// tag the image and upload it to the private registry
	cli.DockerCmd(c, "tag", "busybox", targetName)

	cli.Docker(cli.Args("push", targetName), trustedCmd).Assert(c, SuccessSigningAndPushing)
	// check to make sure that the target has been added to targets/releases and not targets
	s.assertTargetInRoles(c, repoName, "latest", "targets/releases")
	s.assertTargetNotInRoles(c, repoName, "latest", "targets")

	// Try pull after push
	os.RemoveAll(filepath.Join(config.Dir(), "trust"))

	cli.Docker(cli.Args("pull", targetName), trustedCmd).Assert(c, icmd.Expected***REMOVED***
		Out: "Status: Image is up to date",
	***REMOVED***)
***REMOVED***

func (s *DockerTrustSuite) TestTrustedPushSignsAllFirstLevelRolesWeHaveKeysFor(c *check.C) ***REMOVED***
	testRequires(c, NotaryHosting)
	repoName := fmt.Sprintf("%v/dockerclimanyroles/trusted", privateRegistryURL)
	targetName := fmt.Sprintf("%s:latest", repoName)
	s.notaryInitRepo(c, repoName)
	s.notaryCreateDelegation(c, repoName, "targets/role1", s.not.keys[0].Public)
	s.notaryCreateDelegation(c, repoName, "targets/role2", s.not.keys[1].Public)
	s.notaryCreateDelegation(c, repoName, "targets/role3", s.not.keys[2].Public)

	// import everything except the third key
	s.notaryImportKey(c, repoName, "targets/role1", s.not.keys[0].Private)
	s.notaryImportKey(c, repoName, "targets/role2", s.not.keys[1].Private)

	s.notaryCreateDelegation(c, repoName, "targets/role1/subrole", s.not.keys[3].Public)
	s.notaryImportKey(c, repoName, "targets/role1/subrole", s.not.keys[3].Private)

	s.notaryPublish(c, repoName)

	// tag the image and upload it to the private registry
	cli.DockerCmd(c, "tag", "busybox", targetName)

	cli.Docker(cli.Args("push", targetName), trustedCmd).Assert(c, SuccessSigningAndPushing)

	// check to make sure that the target has been added to targets/role1 and targets/role2, and
	// not targets (because there are delegations) or targets/role3 (due to missing key) or
	// targets/role1/subrole (due to it being a second level delegation)
	s.assertTargetInRoles(c, repoName, "latest", "targets/role1", "targets/role2")
	s.assertTargetNotInRoles(c, repoName, "latest", "targets")

	// Try pull after push
	os.RemoveAll(filepath.Join(config.Dir(), "trust"))

	// pull should fail because none of these are the releases role
	cli.Docker(cli.Args("pull", targetName), trustedCmd).Assert(c, icmd.Expected***REMOVED***
		ExitCode: 1,
	***REMOVED***)
***REMOVED***

func (s *DockerTrustSuite) TestTrustedPushSignsForRolesWithKeysAndValidPaths(c *check.C) ***REMOVED***
	repoName := fmt.Sprintf("%v/dockerclirolesbykeysandpaths/trusted", privateRegistryURL)
	targetName := fmt.Sprintf("%s:latest", repoName)
	s.notaryInitRepo(c, repoName)
	s.notaryCreateDelegation(c, repoName, "targets/role1", s.not.keys[0].Public, "l", "z")
	s.notaryCreateDelegation(c, repoName, "targets/role2", s.not.keys[1].Public, "x", "y")
	s.notaryCreateDelegation(c, repoName, "targets/role3", s.not.keys[2].Public, "latest")
	s.notaryCreateDelegation(c, repoName, "targets/role4", s.not.keys[3].Public, "latest")

	// import everything except the third key
	s.notaryImportKey(c, repoName, "targets/role1", s.not.keys[0].Private)
	s.notaryImportKey(c, repoName, "targets/role2", s.not.keys[1].Private)
	s.notaryImportKey(c, repoName, "targets/role4", s.not.keys[3].Private)

	s.notaryPublish(c, repoName)

	// tag the image and upload it to the private registry
	cli.DockerCmd(c, "tag", "busybox", targetName)

	cli.Docker(cli.Args("push", targetName), trustedCmd).Assert(c, SuccessSigningAndPushing)

	// check to make sure that the target has been added to targets/role1 and targets/role4, and
	// not targets (because there are delegations) or targets/role2 (due to path restrictions) or
	// targets/role3 (due to missing key)
	s.assertTargetInRoles(c, repoName, "latest", "targets/role1", "targets/role4")
	s.assertTargetNotInRoles(c, repoName, "latest", "targets")

	// Try pull after push
	os.RemoveAll(filepath.Join(config.Dir(), "trust"))

	// pull should fail because none of these are the releases role
	cli.Docker(cli.Args("pull", targetName), trustedCmd).Assert(c, icmd.Expected***REMOVED***
		ExitCode: 1,
	***REMOVED***)
***REMOVED***

func (s *DockerTrustSuite) TestTrustedPushDoesntSignTargetsIfDelegationsExist(c *check.C) ***REMOVED***
	testRequires(c, NotaryHosting)
	repoName := fmt.Sprintf("%v/dockerclireleasedelegationnotsignable/trusted", privateRegistryURL)
	targetName := fmt.Sprintf("%s:latest", repoName)
	s.notaryInitRepo(c, repoName)
	s.notaryCreateDelegation(c, repoName, "targets/role1", s.not.keys[0].Public)
	s.notaryPublish(c, repoName)

	// do not import any delegations key

	// tag the image and upload it to the private registry
	cli.DockerCmd(c, "tag", "busybox", targetName)

	cli.Docker(cli.Args("push", targetName), trustedCmd).Assert(c, icmd.Expected***REMOVED***
		ExitCode: 1,
		Err:      "no valid signing keys",
	***REMOVED***)
	s.assertTargetNotInRoles(c, repoName, "latest", "targets", "targets/role1")
***REMOVED***

func (s *DockerRegistryAuthHtpasswdSuite) TestPushNoCredentialsNoRetry(c *check.C) ***REMOVED***
	repoName := fmt.Sprintf("%s/busybox", privateRegistryURL)
	dockerCmd(c, "tag", "busybox", repoName)
	out, _, err := dockerCmdWithError("push", repoName)
	c.Assert(err, check.NotNil, check.Commentf(out))
	c.Assert(out, check.Not(checker.Contains), "Retrying")
	c.Assert(out, checker.Contains, "no basic auth credentials")
***REMOVED***

// This may be flaky but it's needed not to regress on unauthorized push, see #21054
func (s *DockerSuite) TestPushToCentralRegistryUnauthorized(c *check.C) ***REMOVED***
	testRequires(c, Network)
	repoName := "test/busybox"
	dockerCmd(c, "tag", "busybox", repoName)
	out, _, err := dockerCmdWithError("push", repoName)
	c.Assert(err, check.NotNil, check.Commentf(out))
	c.Assert(out, check.Not(checker.Contains), "Retrying")
***REMOVED***

func getTestTokenService(status int, body string, retries int) *httptest.Server ***REMOVED***
	var mu sync.Mutex
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		mu.Lock()
		if retries > 0 ***REMOVED***
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`***REMOVED***"errors":[***REMOVED***"code":"UNAVAILABLE","message":"cannot create token at this time"***REMOVED***]***REMOVED***`))
			retries--
		***REMOVED*** else ***REMOVED***
			w.WriteHeader(status)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(body))
		***REMOVED***
		mu.Unlock()
	***REMOVED***))
***REMOVED***

func (s *DockerRegistryAuthTokenSuite) TestPushTokenServiceUnauthResponse(c *check.C) ***REMOVED***
	ts := getTestTokenService(http.StatusUnauthorized, `***REMOVED***"errors": [***REMOVED***"Code":"UNAUTHORIZED", "message": "a message", "detail": null***REMOVED***]***REMOVED***`, 0)
	defer ts.Close()
	s.setupRegistryWithTokenService(c, ts.URL)
	repoName := fmt.Sprintf("%s/busybox", privateRegistryURL)
	dockerCmd(c, "tag", "busybox", repoName)
	out, _, err := dockerCmdWithError("push", repoName)
	c.Assert(err, check.NotNil, check.Commentf(out))
	c.Assert(out, checker.Not(checker.Contains), "Retrying")
	c.Assert(out, checker.Contains, "unauthorized: a message")
***REMOVED***

func (s *DockerRegistryAuthTokenSuite) TestPushMisconfiguredTokenServiceResponseUnauthorized(c *check.C) ***REMOVED***
	ts := getTestTokenService(http.StatusUnauthorized, `***REMOVED***"error": "unauthorized"***REMOVED***`, 0)
	defer ts.Close()
	s.setupRegistryWithTokenService(c, ts.URL)
	repoName := fmt.Sprintf("%s/busybox", privateRegistryURL)
	dockerCmd(c, "tag", "busybox", repoName)
	out, _, err := dockerCmdWithError("push", repoName)
	c.Assert(err, check.NotNil, check.Commentf(out))
	c.Assert(out, checker.Not(checker.Contains), "Retrying")
	split := strings.Split(out, "\n")
	c.Assert(split[len(split)-2], check.Equals, "unauthorized: authentication required")
***REMOVED***

func (s *DockerRegistryAuthTokenSuite) TestPushMisconfiguredTokenServiceResponseError(c *check.C) ***REMOVED***
	ts := getTestTokenService(http.StatusTooManyRequests, `***REMOVED***"errors": [***REMOVED***"code":"TOOMANYREQUESTS","message":"out of tokens"***REMOVED***]***REMOVED***`, 3)
	defer ts.Close()
	s.setupRegistryWithTokenService(c, ts.URL)
	repoName := fmt.Sprintf("%s/busybox", privateRegistryURL)
	dockerCmd(c, "tag", "busybox", repoName)
	out, _, err := dockerCmdWithError("push", repoName)
	c.Assert(err, check.NotNil, check.Commentf(out))
	// TODO: isolate test so that it can be guaranteed that the 503 will trigger xfer retries
	//c.Assert(out, checker.Contains, "Retrying")
	//c.Assert(out, checker.Not(checker.Contains), "Retrying in 15")
	split := strings.Split(out, "\n")
	c.Assert(split[len(split)-2], check.Equals, "toomanyrequests: out of tokens")
***REMOVED***

func (s *DockerRegistryAuthTokenSuite) TestPushMisconfiguredTokenServiceResponseUnparsable(c *check.C) ***REMOVED***
	ts := getTestTokenService(http.StatusForbidden, `no way`, 0)
	defer ts.Close()
	s.setupRegistryWithTokenService(c, ts.URL)
	repoName := fmt.Sprintf("%s/busybox", privateRegistryURL)
	dockerCmd(c, "tag", "busybox", repoName)
	out, _, err := dockerCmdWithError("push", repoName)
	c.Assert(err, check.NotNil, check.Commentf(out))
	c.Assert(out, checker.Not(checker.Contains), "Retrying")
	split := strings.Split(out, "\n")
	c.Assert(split[len(split)-2], checker.Contains, "error parsing HTTP 403 response body: ")
***REMOVED***

func (s *DockerRegistryAuthTokenSuite) TestPushMisconfiguredTokenServiceResponseNoToken(c *check.C) ***REMOVED***
	ts := getTestTokenService(http.StatusOK, `***REMOVED***"something": "wrong"***REMOVED***`, 0)
	defer ts.Close()
	s.setupRegistryWithTokenService(c, ts.URL)
	repoName := fmt.Sprintf("%s/busybox", privateRegistryURL)
	dockerCmd(c, "tag", "busybox", repoName)
	out, _, err := dockerCmdWithError("push", repoName)
	c.Assert(err, check.NotNil, check.Commentf(out))
	c.Assert(out, checker.Not(checker.Contains), "Retrying")
	split := strings.Split(out, "\n")
	c.Assert(split[len(split)-2], check.Equals, "authorization server did not include a token in the response")
***REMOVED***
