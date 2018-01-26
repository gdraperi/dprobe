package main

import (
	"strings"

	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/internal/testutil"
	"github.com/go-check/check"
)

// tagging a named image in a new unprefixed repo should work
func (s *DockerSuite) TestTagUnprefixedRepoByName(c *check.C) ***REMOVED***
	dockerCmd(c, "tag", "busybox:latest", "testfoobarbaz")
***REMOVED***

// tagging an image by ID in a new unprefixed repo should work
func (s *DockerSuite) TestTagUnprefixedRepoByID(c *check.C) ***REMOVED***
	imageID := inspectField(c, "busybox", "Id")
	dockerCmd(c, "tag", imageID, "testfoobarbaz")
***REMOVED***

// ensure we don't allow the use of invalid repository names; these tag operations should fail
func (s *DockerSuite) TestTagInvalidUnprefixedRepo(c *check.C) ***REMOVED***
	invalidRepos := []string***REMOVED***"fo$z$", "Foo@3cc", "Foo$3", "Foo*3", "Fo^3", "Foo!3", "F)xcz(", "fo%asd", "FOO/bar"***REMOVED***

	for _, repo := range invalidRepos ***REMOVED***
		out, _, err := dockerCmdWithError("tag", "busybox", repo)
		c.Assert(err, checker.NotNil, check.Commentf("tag busybox %v should have failed : %v", repo, out))
	***REMOVED***
***REMOVED***

// ensure we don't allow the use of invalid tags; these tag operations should fail
func (s *DockerSuite) TestTagInvalidPrefixedRepo(c *check.C) ***REMOVED***
	longTag := testutil.GenerateRandomAlphaOnlyString(121)

	invalidTags := []string***REMOVED***"repo:fo$z$", "repo:Foo@3cc", "repo:Foo$3", "repo:Foo*3", "repo:Fo^3", "repo:Foo!3", "repo:%goodbye", "repo:#hashtagit", "repo:F)xcz(", "repo:-foo", "repo:..", longTag***REMOVED***

	for _, repotag := range invalidTags ***REMOVED***
		out, _, err := dockerCmdWithError("tag", "busybox", repotag)
		c.Assert(err, checker.NotNil, check.Commentf("tag busybox %v should have failed : %v", repotag, out))
	***REMOVED***
***REMOVED***

// ensure we allow the use of valid tags
func (s *DockerSuite) TestTagValidPrefixedRepo(c *check.C) ***REMOVED***
	validRepos := []string***REMOVED***"fooo/bar", "fooaa/test", "foooo:t", "HOSTNAME.DOMAIN.COM:443/foo/bar"***REMOVED***

	for _, repo := range validRepos ***REMOVED***
		_, _, err := dockerCmdWithError("tag", "busybox:latest", repo)
		if err != nil ***REMOVED***
			c.Errorf("tag busybox %v should have worked: %s", repo, err)
			continue
		***REMOVED***
		deleteImages(repo)
	***REMOVED***
***REMOVED***

// tag an image with an existed tag name without -f option should work
func (s *DockerSuite) TestTagExistedNameWithoutForce(c *check.C) ***REMOVED***
	dockerCmd(c, "tag", "busybox:latest", "busybox:test")
***REMOVED***

func (s *DockerSuite) TestTagWithPrefixHyphen(c *check.C) ***REMOVED***
	// test repository name begin with '-'
	out, _, err := dockerCmdWithError("tag", "busybox:latest", "-busybox:test")
	c.Assert(err, checker.NotNil, check.Commentf(out))
	c.Assert(out, checker.Contains, "Error parsing reference", check.Commentf("tag a name begin with '-' should failed"))

	// test namespace name begin with '-'
	out, _, err = dockerCmdWithError("tag", "busybox:latest", "-test/busybox:test")
	c.Assert(err, checker.NotNil, check.Commentf(out))
	c.Assert(out, checker.Contains, "Error parsing reference", check.Commentf("tag a name begin with '-' should failed"))

	// test index name begin with '-'
	out, _, err = dockerCmdWithError("tag", "busybox:latest", "-index:5000/busybox:test")
	c.Assert(err, checker.NotNil, check.Commentf(out))
	c.Assert(out, checker.Contains, "Error parsing reference", check.Commentf("tag a name begin with '-' should failed"))
***REMOVED***

// ensure tagging using official names works
// ensure all tags result in the same name
func (s *DockerSuite) TestTagOfficialNames(c *check.C) ***REMOVED***
	names := []string***REMOVED***
		"docker.io/busybox",
		"index.docker.io/busybox",
		"library/busybox",
		"docker.io/library/busybox",
		"index.docker.io/library/busybox",
	***REMOVED***

	for _, name := range names ***REMOVED***
		out, exitCode, err := dockerCmdWithError("tag", "busybox:latest", name+":latest")
		if err != nil || exitCode != 0 ***REMOVED***
			c.Errorf("tag busybox %v should have worked: %s, %s", name, err, out)
			continue
		***REMOVED***

		// ensure we don't have multiple tag names.
		out, _, err = dockerCmdWithError("images")
		if err != nil ***REMOVED***
			c.Errorf("listing images failed with errors: %v, %s", err, out)
		***REMOVED*** else if strings.Contains(out, name) ***REMOVED***
			c.Errorf("images should not have listed '%s'", name)
			deleteImages(name + ":latest")
		***REMOVED***
	***REMOVED***

	for _, name := range names ***REMOVED***
		_, exitCode, err := dockerCmdWithError("tag", name+":latest", "fooo/bar:latest")
		if err != nil || exitCode != 0 ***REMOVED***
			c.Errorf("tag %v fooo/bar should have worked: %s", name, err)
			continue
		***REMOVED***
		deleteImages("fooo/bar:latest")
	***REMOVED***
***REMOVED***

// ensure tags can not match digests
func (s *DockerSuite) TestTagMatchesDigest(c *check.C) ***REMOVED***
	digest := "busybox@sha256:abcdef76720241213f5303bda7704ec4c2ef75613173910a56fb1b6e20251507"
	// test setting tag fails
	_, _, err := dockerCmdWithError("tag", "busybox:latest", digest)
	if err == nil ***REMOVED***
		c.Fatal("digest tag a name should have failed")
	***REMOVED***
	// check that no new image matches the digest
	_, _, err = dockerCmdWithError("inspect", digest)
	if err == nil ***REMOVED***
		c.Fatal("inspecting by digest should have failed")
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestTagInvalidRepoName(c *check.C) ***REMOVED***
	// test setting tag fails
	_, _, err := dockerCmdWithError("tag", "busybox:latest", "sha256:sometag")
	if err == nil ***REMOVED***
		c.Fatal("tagging with image named \"sha256\" should have failed")
	***REMOVED***
***REMOVED***
