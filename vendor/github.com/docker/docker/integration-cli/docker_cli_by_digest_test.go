package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/cli"
	"github.com/docker/docker/integration-cli/cli/build"
	"github.com/go-check/check"
	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
)

var (
	remoteRepoName  = "dockercli/busybox-by-dgst"
	repoName        = fmt.Sprintf("%s/%s", privateRegistryURL, remoteRepoName)
	pushDigestRegex = regexp.MustCompile("[\\S]+: digest: ([\\S]+) size: [0-9]+")
	digestRegex     = regexp.MustCompile("Digest: ([\\S]+)")
)

func setupImage(c *check.C) (digest.Digest, error) ***REMOVED***
	return setupImageWithTag(c, "latest")
***REMOVED***

func setupImageWithTag(c *check.C, tag string) (digest.Digest, error) ***REMOVED***
	containerName := "busyboxbydigest"

	// new file is committed because this layer is used for detecting malicious
	// changes. if this was committed as empty layer it would be skipped on pull
	// and malicious changes would never be detected.
	cli.DockerCmd(c, "run", "-e", "digest=1", "--name", containerName, "busybox", "touch", "anewfile")

	// tag the image to upload it to the private registry
	repoAndTag := repoName + ":" + tag
	cli.DockerCmd(c, "commit", containerName, repoAndTag)

	// delete the container as we don't need it any more
	cli.DockerCmd(c, "rm", "-fv", containerName)

	// push the image
	out := cli.DockerCmd(c, "push", repoAndTag).Combined()

	// delete our local repo that we previously tagged
	cli.DockerCmd(c, "rmi", repoAndTag)

	matches := pushDigestRegex.FindStringSubmatch(out)
	c.Assert(matches, checker.HasLen, 2, check.Commentf("unable to parse digest from push output: %s", out))
	pushDigest := matches[1]

	return digest.Digest(pushDigest), nil
***REMOVED***

func testPullByTagDisplaysDigest(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	pushDigest, err := setupImage(c)
	c.Assert(err, checker.IsNil, check.Commentf("error setting up image"))

	// pull from the registry using the tag
	out, _ := dockerCmd(c, "pull", repoName)

	// the pull output includes "Digest: <digest>", so find that
	matches := digestRegex.FindStringSubmatch(out)
	c.Assert(matches, checker.HasLen, 2, check.Commentf("unable to parse digest from pull output: %s", out))
	pullDigest := matches[1]

	// make sure the pushed and pull digests match
	c.Assert(pushDigest.String(), checker.Equals, pullDigest)
***REMOVED***

func (s *DockerRegistrySuite) TestPullByTagDisplaysDigest(c *check.C) ***REMOVED***
	testPullByTagDisplaysDigest(c)
***REMOVED***

func (s *DockerSchema1RegistrySuite) TestPullByTagDisplaysDigest(c *check.C) ***REMOVED***
	testPullByTagDisplaysDigest(c)
***REMOVED***

func testPullByDigest(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	pushDigest, err := setupImage(c)
	c.Assert(err, checker.IsNil, check.Commentf("error setting up image"))

	// pull from the registry using the <name>@<digest> reference
	imageReference := fmt.Sprintf("%s@%s", repoName, pushDigest)
	out, _ := dockerCmd(c, "pull", imageReference)

	// the pull output includes "Digest: <digest>", so find that
	matches := digestRegex.FindStringSubmatch(out)
	c.Assert(matches, checker.HasLen, 2, check.Commentf("unable to parse digest from pull output: %s", out))
	pullDigest := matches[1]

	// make sure the pushed and pull digests match
	c.Assert(pushDigest.String(), checker.Equals, pullDigest)
***REMOVED***

func (s *DockerRegistrySuite) TestPullByDigest(c *check.C) ***REMOVED***
	testPullByDigest(c)
***REMOVED***

func (s *DockerSchema1RegistrySuite) TestPullByDigest(c *check.C) ***REMOVED***
	testPullByDigest(c)
***REMOVED***

func testPullByDigestNoFallback(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	// pull from the registry using the <name>@<digest> reference
	imageReference := fmt.Sprintf("%s@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", repoName)
	out, _, err := dockerCmdWithError("pull", imageReference)
	c.Assert(err, checker.NotNil, check.Commentf("expected non-zero exit status and correct error message when pulling non-existing image"))
	c.Assert(out, checker.Contains, fmt.Sprintf("manifest for %s not found", imageReference), check.Commentf("expected non-zero exit status and correct error message when pulling non-existing image"))
***REMOVED***

func (s *DockerRegistrySuite) TestPullByDigestNoFallback(c *check.C) ***REMOVED***
	testPullByDigestNoFallback(c)
***REMOVED***

func (s *DockerSchema1RegistrySuite) TestPullByDigestNoFallback(c *check.C) ***REMOVED***
	testPullByDigestNoFallback(c)
***REMOVED***

func (s *DockerRegistrySuite) TestCreateByDigest(c *check.C) ***REMOVED***
	pushDigest, err := setupImage(c)
	c.Assert(err, checker.IsNil, check.Commentf("error setting up image"))

	imageReference := fmt.Sprintf("%s@%s", repoName, pushDigest)

	containerName := "createByDigest"
	dockerCmd(c, "create", "--name", containerName, imageReference)

	res := inspectField(c, containerName, "Config.Image")
	c.Assert(res, checker.Equals, imageReference)
***REMOVED***

func (s *DockerRegistrySuite) TestRunByDigest(c *check.C) ***REMOVED***
	pushDigest, err := setupImage(c)
	c.Assert(err, checker.IsNil)

	imageReference := fmt.Sprintf("%s@%s", repoName, pushDigest)

	containerName := "runByDigest"
	out, _ := dockerCmd(c, "run", "--name", containerName, imageReference, "sh", "-c", "echo found=$digest")

	foundRegex := regexp.MustCompile("found=([^\n]+)")
	matches := foundRegex.FindStringSubmatch(out)
	c.Assert(matches, checker.HasLen, 2, check.Commentf("unable to parse digest from pull output: %s", out))
	c.Assert(matches[1], checker.Equals, "1", check.Commentf("Expected %q, got %q", "1", matches[1]))

	res := inspectField(c, containerName, "Config.Image")
	c.Assert(res, checker.Equals, imageReference)
***REMOVED***

func (s *DockerRegistrySuite) TestRemoveImageByDigest(c *check.C) ***REMOVED***
	digest, err := setupImage(c)
	c.Assert(err, checker.IsNil, check.Commentf("error setting up image"))

	imageReference := fmt.Sprintf("%s@%s", repoName, digest)

	// pull from the registry using the <name>@<digest> reference
	dockerCmd(c, "pull", imageReference)

	// make sure inspect runs ok
	inspectField(c, imageReference, "Id")

	// do the delete
	err = deleteImages(imageReference)
	c.Assert(err, checker.IsNil, check.Commentf("unexpected error deleting image"))

	// try to inspect again - it should error this time
	_, err = inspectFieldWithError(imageReference, "Id")
	//unexpected nil err trying to inspect what should be a non-existent image
	c.Assert(err, checker.NotNil)
	c.Assert(err.Error(), checker.Contains, "No such object")
***REMOVED***

func (s *DockerRegistrySuite) TestBuildByDigest(c *check.C) ***REMOVED***
	digest, err := setupImage(c)
	c.Assert(err, checker.IsNil, check.Commentf("error setting up image"))

	imageReference := fmt.Sprintf("%s@%s", repoName, digest)

	// pull from the registry using the <name>@<digest> reference
	dockerCmd(c, "pull", imageReference)

	// get the image id
	imageID := inspectField(c, imageReference, "Id")

	// do the build
	name := "buildbydigest"
	buildImageSuccessfully(c, name, build.WithDockerfile(fmt.Sprintf(
		`FROM %s
     CMD ["/bin/echo", "Hello World"]`, imageReference)))
	c.Assert(err, checker.IsNil)

	// get the build's image id
	res := inspectField(c, name, "Config.Image")
	// make sure they match
	c.Assert(res, checker.Equals, imageID)
***REMOVED***

func (s *DockerRegistrySuite) TestTagByDigest(c *check.C) ***REMOVED***
	digest, err := setupImage(c)
	c.Assert(err, checker.IsNil, check.Commentf("error setting up image"))

	imageReference := fmt.Sprintf("%s@%s", repoName, digest)

	// pull from the registry using the <name>@<digest> reference
	dockerCmd(c, "pull", imageReference)

	// tag it
	tag := "tagbydigest"
	dockerCmd(c, "tag", imageReference, tag)

	expectedID := inspectField(c, imageReference, "Id")

	tagID := inspectField(c, tag, "Id")
	c.Assert(tagID, checker.Equals, expectedID)
***REMOVED***

func (s *DockerRegistrySuite) TestListImagesWithoutDigests(c *check.C) ***REMOVED***
	digest, err := setupImage(c)
	c.Assert(err, checker.IsNil, check.Commentf("error setting up image"))

	imageReference := fmt.Sprintf("%s@%s", repoName, digest)

	// pull from the registry using the <name>@<digest> reference
	dockerCmd(c, "pull", imageReference)

	out, _ := dockerCmd(c, "images")
	c.Assert(out, checker.Not(checker.Contains), "DIGEST", check.Commentf("list output should not have contained DIGEST header"))
***REMOVED***

func (s *DockerRegistrySuite) TestListImagesWithDigests(c *check.C) ***REMOVED***

	// setup image1
	digest1, err := setupImageWithTag(c, "tag1")
	c.Assert(err, checker.IsNil, check.Commentf("error setting up image"))
	imageReference1 := fmt.Sprintf("%s@%s", repoName, digest1)
	c.Logf("imageReference1 = %s", imageReference1)

	// pull image1 by digest
	dockerCmd(c, "pull", imageReference1)

	// list images
	out, _ := dockerCmd(c, "images", "--digests")

	// make sure repo shown, tag=<none>, digest = $digest1
	re1 := regexp.MustCompile(`\s*` + repoName + `\s*<none>\s*` + digest1.String() + `\s`)
	c.Assert(re1.MatchString(out), checker.True, check.Commentf("expected %q: %s", re1.String(), out))
	// setup image2
	digest2, err := setupImageWithTag(c, "tag2")
	//error setting up image
	c.Assert(err, checker.IsNil)
	imageReference2 := fmt.Sprintf("%s@%s", repoName, digest2)
	c.Logf("imageReference2 = %s", imageReference2)

	// pull image1 by digest
	dockerCmd(c, "pull", imageReference1)

	// pull image2 by digest
	dockerCmd(c, "pull", imageReference2)

	// list images
	out, _ = dockerCmd(c, "images", "--digests")

	// make sure repo shown, tag=<none>, digest = $digest1
	c.Assert(re1.MatchString(out), checker.True, check.Commentf("expected %q: %s", re1.String(), out))

	// make sure repo shown, tag=<none>, digest = $digest2
	re2 := regexp.MustCompile(`\s*` + repoName + `\s*<none>\s*` + digest2.String() + `\s`)
	c.Assert(re2.MatchString(out), checker.True, check.Commentf("expected %q: %s", re2.String(), out))

	// pull tag1
	dockerCmd(c, "pull", repoName+":tag1")

	// list images
	out, _ = dockerCmd(c, "images", "--digests")

	// make sure image 1 has repo, tag, <none> AND repo, <none>, digest
	reWithDigest1 := regexp.MustCompile(`\s*` + repoName + `\s*tag1\s*` + digest1.String() + `\s`)
	c.Assert(reWithDigest1.MatchString(out), checker.True, check.Commentf("expected %q: %s", reWithDigest1.String(), out))
	// make sure image 2 has repo, <none>, digest
	c.Assert(re2.MatchString(out), checker.True, check.Commentf("expected %q: %s", re2.String(), out))

	// pull tag 2
	dockerCmd(c, "pull", repoName+":tag2")

	// list images
	out, _ = dockerCmd(c, "images", "--digests")

	// make sure image 1 has repo, tag, digest
	c.Assert(reWithDigest1.MatchString(out), checker.True, check.Commentf("expected %q: %s", reWithDigest1.String(), out))

	// make sure image 2 has repo, tag, digest
	reWithDigest2 := regexp.MustCompile(`\s*` + repoName + `\s*tag2\s*` + digest2.String() + `\s`)
	c.Assert(reWithDigest2.MatchString(out), checker.True, check.Commentf("expected %q: %s", reWithDigest2.String(), out))

	// list images
	out, _ = dockerCmd(c, "images", "--digests")

	// make sure image 1 has repo, tag, digest
	c.Assert(reWithDigest1.MatchString(out), checker.True, check.Commentf("expected %q: %s", reWithDigest1.String(), out))
	// make sure image 2 has repo, tag, digest
	c.Assert(reWithDigest2.MatchString(out), checker.True, check.Commentf("expected %q: %s", reWithDigest2.String(), out))
	// make sure busybox has tag, but not digest
	busyboxRe := regexp.MustCompile(`\s*busybox\s*latest\s*<none>\s`)
	c.Assert(busyboxRe.MatchString(out), checker.True, check.Commentf("expected %q: %s", busyboxRe.String(), out))
***REMOVED***

func (s *DockerRegistrySuite) TestListDanglingImagesWithDigests(c *check.C) ***REMOVED***
	// setup image1
	digest1, err := setupImageWithTag(c, "dangle1")
	c.Assert(err, checker.IsNil, check.Commentf("error setting up image"))
	imageReference1 := fmt.Sprintf("%s@%s", repoName, digest1)
	c.Logf("imageReference1 = %s", imageReference1)

	// pull image1 by digest
	dockerCmd(c, "pull", imageReference1)

	// list images
	out, _ := dockerCmd(c, "images", "--digests")

	// make sure repo shown, tag=<none>, digest = $digest1
	re1 := regexp.MustCompile(`\s*` + repoName + `\s*<none>\s*` + digest1.String() + `\s`)
	c.Assert(re1.MatchString(out), checker.True, check.Commentf("expected %q: %s", re1.String(), out))
	// setup image2
	digest2, err := setupImageWithTag(c, "dangle2")
	//error setting up image
	c.Assert(err, checker.IsNil)
	imageReference2 := fmt.Sprintf("%s@%s", repoName, digest2)
	c.Logf("imageReference2 = %s", imageReference2)

	// pull image1 by digest
	dockerCmd(c, "pull", imageReference1)

	// pull image2 by digest
	dockerCmd(c, "pull", imageReference2)

	// list images
	out, _ = dockerCmd(c, "images", "--digests", "--filter=dangling=true")

	// make sure repo shown, tag=<none>, digest = $digest1
	c.Assert(re1.MatchString(out), checker.True, check.Commentf("expected %q: %s", re1.String(), out))

	// make sure repo shown, tag=<none>, digest = $digest2
	re2 := regexp.MustCompile(`\s*` + repoName + `\s*<none>\s*` + digest2.String() + `\s`)
	c.Assert(re2.MatchString(out), checker.True, check.Commentf("expected %q: %s", re2.String(), out))

	// pull dangle1 tag
	dockerCmd(c, "pull", repoName+":dangle1")

	// list images
	out, _ = dockerCmd(c, "images", "--digests", "--filter=dangling=true")

	// make sure image 1 has repo, tag, <none> AND repo, <none>, digest
	reWithDigest1 := regexp.MustCompile(`\s*` + repoName + `\s*dangle1\s*` + digest1.String() + `\s`)
	c.Assert(reWithDigest1.MatchString(out), checker.False, check.Commentf("unexpected %q: %s", reWithDigest1.String(), out))
	// make sure image 2 has repo, <none>, digest
	c.Assert(re2.MatchString(out), checker.True, check.Commentf("expected %q: %s", re2.String(), out))

	// pull dangle2 tag
	dockerCmd(c, "pull", repoName+":dangle2")

	// list images, show tagged images
	out, _ = dockerCmd(c, "images", "--digests")

	// make sure image 1 has repo, tag, digest
	c.Assert(reWithDigest1.MatchString(out), checker.True, check.Commentf("expected %q: %s", reWithDigest1.String(), out))

	// make sure image 2 has repo, tag, digest
	reWithDigest2 := regexp.MustCompile(`\s*` + repoName + `\s*dangle2\s*` + digest2.String() + `\s`)
	c.Assert(reWithDigest2.MatchString(out), checker.True, check.Commentf("expected %q: %s", reWithDigest2.String(), out))

	// list images, no longer dangling, should not match
	out, _ = dockerCmd(c, "images", "--digests", "--filter=dangling=true")

	// make sure image 1 has repo, tag, digest
	c.Assert(reWithDigest1.MatchString(out), checker.False, check.Commentf("unexpected %q: %s", reWithDigest1.String(), out))
	// make sure image 2 has repo, tag, digest
	c.Assert(reWithDigest2.MatchString(out), checker.False, check.Commentf("unexpected %q: %s", reWithDigest2.String(), out))
***REMOVED***

func (s *DockerRegistrySuite) TestInspectImageWithDigests(c *check.C) ***REMOVED***
	digest, err := setupImage(c)
	c.Assert(err, check.IsNil, check.Commentf("error setting up image"))

	imageReference := fmt.Sprintf("%s@%s", repoName, digest)

	// pull from the registry using the <name>@<digest> reference
	dockerCmd(c, "pull", imageReference)

	out, _ := dockerCmd(c, "inspect", imageReference)

	var imageJSON []types.ImageInspect
	err = json.Unmarshal([]byte(out), &imageJSON)
	c.Assert(err, checker.IsNil)
	c.Assert(imageJSON, checker.HasLen, 1)
	c.Assert(imageJSON[0].RepoDigests, checker.HasLen, 1)
	assert.Contains(c, imageJSON[0].RepoDigests, imageReference)
***REMOVED***

func (s *DockerRegistrySuite) TestPsListContainersFilterAncestorImageByDigest(c *check.C) ***REMOVED***
	existingContainers := ExistingContainerIDs(c)

	digest, err := setupImage(c)
	c.Assert(err, checker.IsNil, check.Commentf("error setting up image"))

	imageReference := fmt.Sprintf("%s@%s", repoName, digest)

	// pull from the registry using the <name>@<digest> reference
	dockerCmd(c, "pull", imageReference)

	// build an image from it
	imageName1 := "images_ps_filter_test"
	buildImageSuccessfully(c, imageName1, build.WithDockerfile(fmt.Sprintf(
		`FROM %s
		 LABEL match me 1`, imageReference)))

	// run a container based on that
	dockerCmd(c, "run", "--name=test1", imageReference, "echo", "hello")
	expectedID := getIDByName(c, "test1")

	// run a container based on the a descendant of that too
	dockerCmd(c, "run", "--name=test2", imageName1, "echo", "hello")
	expectedID1 := getIDByName(c, "test2")

	expectedIDs := []string***REMOVED***expectedID, expectedID1***REMOVED***

	// Invalid imageReference
	out, _ := dockerCmd(c, "ps", "-a", "-q", "--no-trunc", fmt.Sprintf("--filter=ancestor=busybox@%s", digest))
	// Filter container for ancestor filter should be empty
	c.Assert(strings.TrimSpace(out), checker.Equals, "")

	// Valid imageReference
	out, _ = dockerCmd(c, "ps", "-a", "-q", "--no-trunc", "--filter=ancestor="+imageReference)
	checkPsAncestorFilterOutput(c, RemoveOutputForExistingElements(out, existingContainers), imageReference, expectedIDs)
***REMOVED***

func (s *DockerRegistrySuite) TestDeleteImageByIDOnlyPulledByDigest(c *check.C) ***REMOVED***
	pushDigest, err := setupImage(c)
	c.Assert(err, checker.IsNil, check.Commentf("error setting up image"))

	// pull from the registry using the <name>@<digest> reference
	imageReference := fmt.Sprintf("%s@%s", repoName, pushDigest)
	dockerCmd(c, "pull", imageReference)
	// just in case...

	dockerCmd(c, "tag", imageReference, repoName+":sometag")

	imageID := inspectField(c, imageReference, "Id")

	dockerCmd(c, "rmi", imageID)

	_, err = inspectFieldWithError(imageID, "Id")
	c.Assert(err, checker.NotNil, check.Commentf("image should have been deleted"))
***REMOVED***

func (s *DockerRegistrySuite) TestDeleteImageWithDigestAndTag(c *check.C) ***REMOVED***
	pushDigest, err := setupImage(c)
	c.Assert(err, checker.IsNil, check.Commentf("error setting up image"))

	// pull from the registry using the <name>@<digest> reference
	imageReference := fmt.Sprintf("%s@%s", repoName, pushDigest)
	dockerCmd(c, "pull", imageReference)

	imageID := inspectField(c, imageReference, "Id")

	repoTag := repoName + ":sometag"
	repoTag2 := repoName + ":othertag"
	dockerCmd(c, "tag", imageReference, repoTag)
	dockerCmd(c, "tag", imageReference, repoTag2)

	dockerCmd(c, "rmi", repoTag2)

	// rmi should have deleted only repoTag2, because there's another tag
	inspectField(c, repoTag, "Id")

	dockerCmd(c, "rmi", repoTag)

	// rmi should have deleted the tag, the digest reference, and the image itself
	_, err = inspectFieldWithError(imageID, "Id")
	c.Assert(err, checker.NotNil, check.Commentf("image should have been deleted"))
***REMOVED***

func (s *DockerRegistrySuite) TestDeleteImageWithDigestAndMultiRepoTag(c *check.C) ***REMOVED***
	pushDigest, err := setupImage(c)
	c.Assert(err, checker.IsNil, check.Commentf("error setting up image"))

	repo2 := fmt.Sprintf("%s/%s", repoName, "repo2")

	// pull from the registry using the <name>@<digest> reference
	imageReference := fmt.Sprintf("%s@%s", repoName, pushDigest)
	dockerCmd(c, "pull", imageReference)

	imageID := inspectField(c, imageReference, "Id")

	repoTag := repoName + ":sometag"
	repoTag2 := repo2 + ":othertag"
	dockerCmd(c, "tag", imageReference, repoTag)
	dockerCmd(c, "tag", imageReference, repoTag2)

	dockerCmd(c, "rmi", repoTag)

	// rmi should have deleted repoTag and image reference, but left repoTag2
	inspectField(c, repoTag2, "Id")
	_, err = inspectFieldWithError(imageReference, "Id")
	c.Assert(err, checker.NotNil, check.Commentf("image digest reference should have been removed"))

	_, err = inspectFieldWithError(repoTag, "Id")
	c.Assert(err, checker.NotNil, check.Commentf("image tag reference should have been removed"))

	dockerCmd(c, "rmi", repoTag2)

	// rmi should have deleted the tag, the digest reference, and the image itself
	_, err = inspectFieldWithError(imageID, "Id")
	c.Assert(err, checker.NotNil, check.Commentf("image should have been deleted"))
***REMOVED***

// TestPullFailsWithAlteredManifest tests that a `docker pull` fails when
// we have modified a manifest blob and its digest cannot be verified.
// This is the schema2 version of the test.
func (s *DockerRegistrySuite) TestPullFailsWithAlteredManifest(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	manifestDigest, err := setupImage(c)
	c.Assert(err, checker.IsNil, check.Commentf("error setting up image"))

	// Load the target manifest blob.
	manifestBlob := s.reg.ReadBlobContents(c, manifestDigest)

	var imgManifest schema2.Manifest
	err = json.Unmarshal(manifestBlob, &imgManifest)
	c.Assert(err, checker.IsNil, check.Commentf("unable to decode image manifest from blob"))

	// Change a layer in the manifest.
	imgManifest.Layers[0].Digest = digest.Digest("sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")

	// Move the existing data file aside, so that we can replace it with a
	// malicious blob of data. NOTE: we defer the returned undo func.
	undo := s.reg.TempMoveBlobData(c, manifestDigest)
	defer undo()

	alteredManifestBlob, err := json.MarshalIndent(imgManifest, "", "   ")
	c.Assert(err, checker.IsNil, check.Commentf("unable to encode altered image manifest to JSON"))

	s.reg.WriteBlobContents(c, manifestDigest, alteredManifestBlob)

	// Now try pulling that image by digest. We should get an error about
	// digest verification for the manifest digest.

	// Pull from the registry using the <name>@<digest> reference.
	imageReference := fmt.Sprintf("%s@%s", repoName, manifestDigest)
	out, exitStatus, _ := dockerCmdWithError("pull", imageReference)
	c.Assert(exitStatus, checker.Not(check.Equals), 0)

	expectedErrorMsg := fmt.Sprintf("manifest verification failed for digest %s", manifestDigest)
	c.Assert(out, checker.Contains, expectedErrorMsg)
***REMOVED***

// TestPullFailsWithAlteredManifest tests that a `docker pull` fails when
// we have modified a manifest blob and its digest cannot be verified.
// This is the schema1 version of the test.
func (s *DockerSchema1RegistrySuite) TestPullFailsWithAlteredManifest(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	manifestDigest, err := setupImage(c)
	c.Assert(err, checker.IsNil, check.Commentf("error setting up image"))

	// Load the target manifest blob.
	manifestBlob := s.reg.ReadBlobContents(c, manifestDigest)

	var imgManifest schema1.Manifest
	err = json.Unmarshal(manifestBlob, &imgManifest)
	c.Assert(err, checker.IsNil, check.Commentf("unable to decode image manifest from blob"))

	// Change a layer in the manifest.
	imgManifest.FSLayers[0] = schema1.FSLayer***REMOVED***
		BlobSum: digest.Digest("sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"),
	***REMOVED***

	// Move the existing data file aside, so that we can replace it with a
	// malicious blob of data. NOTE: we defer the returned undo func.
	undo := s.reg.TempMoveBlobData(c, manifestDigest)
	defer undo()

	alteredManifestBlob, err := json.MarshalIndent(imgManifest, "", "   ")
	c.Assert(err, checker.IsNil, check.Commentf("unable to encode altered image manifest to JSON"))

	s.reg.WriteBlobContents(c, manifestDigest, alteredManifestBlob)

	// Now try pulling that image by digest. We should get an error about
	// digest verification for the manifest digest.

	// Pull from the registry using the <name>@<digest> reference.
	imageReference := fmt.Sprintf("%s@%s", repoName, manifestDigest)
	out, exitStatus, _ := dockerCmdWithError("pull", imageReference)
	c.Assert(exitStatus, checker.Not(check.Equals), 0)

	expectedErrorMsg := fmt.Sprintf("image verification failed for digest %s", manifestDigest)
	c.Assert(out, checker.Contains, expectedErrorMsg)
***REMOVED***

// TestPullFailsWithAlteredLayer tests that a `docker pull` fails when
// we have modified a layer blob and its digest cannot be verified.
// This is the schema2 version of the test.
func (s *DockerRegistrySuite) TestPullFailsWithAlteredLayer(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	manifestDigest, err := setupImage(c)
	c.Assert(err, checker.IsNil)

	// Load the target manifest blob.
	manifestBlob := s.reg.ReadBlobContents(c, manifestDigest)

	var imgManifest schema2.Manifest
	err = json.Unmarshal(manifestBlob, &imgManifest)
	c.Assert(err, checker.IsNil)

	// Next, get the digest of one of the layers from the manifest.
	targetLayerDigest := imgManifest.Layers[0].Digest

	// Move the existing data file aside, so that we can replace it with a
	// malicious blob of data. NOTE: we defer the returned undo func.
	undo := s.reg.TempMoveBlobData(c, targetLayerDigest)
	defer undo()

	// Now make a fake data blob in this directory.
	s.reg.WriteBlobContents(c, targetLayerDigest, []byte("This is not the data you are looking for."))

	// Now try pulling that image by digest. We should get an error about
	// digest verification for the target layer digest.

	// Remove distribution cache to force a re-pull of the blobs
	if err := os.RemoveAll(filepath.Join(testEnv.DaemonInfo.DockerRootDir, "image", s.d.StorageDriver(), "distribution")); err != nil ***REMOVED***
		c.Fatalf("error clearing distribution cache: %v", err)
	***REMOVED***

	// Pull from the registry using the <name>@<digest> reference.
	imageReference := fmt.Sprintf("%s@%s", repoName, manifestDigest)
	out, exitStatus, _ := dockerCmdWithError("pull", imageReference)
	c.Assert(exitStatus, checker.Not(check.Equals), 0, check.Commentf("expected a non-zero exit status"))

	expectedErrorMsg := fmt.Sprintf("filesystem layer verification failed for digest %s", targetLayerDigest)
	c.Assert(out, checker.Contains, expectedErrorMsg, check.Commentf("expected error message in output: %s", out))
***REMOVED***

// TestPullFailsWithAlteredLayer tests that a `docker pull` fails when
// we have modified a layer blob and its digest cannot be verified.
// This is the schema1 version of the test.
func (s *DockerSchema1RegistrySuite) TestPullFailsWithAlteredLayer(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	manifestDigest, err := setupImage(c)
	c.Assert(err, checker.IsNil)

	// Load the target manifest blob.
	manifestBlob := s.reg.ReadBlobContents(c, manifestDigest)

	var imgManifest schema1.Manifest
	err = json.Unmarshal(manifestBlob, &imgManifest)
	c.Assert(err, checker.IsNil)

	// Next, get the digest of one of the layers from the manifest.
	targetLayerDigest := imgManifest.FSLayers[0].BlobSum

	// Move the existing data file aside, so that we can replace it with a
	// malicious blob of data. NOTE: we defer the returned undo func.
	undo := s.reg.TempMoveBlobData(c, targetLayerDigest)
	defer undo()

	// Now make a fake data blob in this directory.
	s.reg.WriteBlobContents(c, targetLayerDigest, []byte("This is not the data you are looking for."))

	// Now try pulling that image by digest. We should get an error about
	// digest verification for the target layer digest.

	// Remove distribution cache to force a re-pull of the blobs
	if err := os.RemoveAll(filepath.Join(testEnv.DaemonInfo.DockerRootDir, "image", s.d.StorageDriver(), "distribution")); err != nil ***REMOVED***
		c.Fatalf("error clearing distribution cache: %v", err)
	***REMOVED***

	// Pull from the registry using the <name>@<digest> reference.
	imageReference := fmt.Sprintf("%s@%s", repoName, manifestDigest)
	out, exitStatus, _ := dockerCmdWithError("pull", imageReference)
	c.Assert(exitStatus, checker.Not(check.Equals), 0, check.Commentf("expected a non-zero exit status"))

	expectedErrorMsg := fmt.Sprintf("filesystem layer verification failed for digest %s", targetLayerDigest)
	c.Assert(out, checker.Contains, expectedErrorMsg, check.Commentf("expected error message in output: %s", out))
***REMOVED***
