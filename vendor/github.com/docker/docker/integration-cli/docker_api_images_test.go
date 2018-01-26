package main

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/cli"
	"github.com/docker/docker/integration-cli/cli/build"
	"github.com/docker/docker/integration-cli/request"
	"github.com/go-check/check"
	"golang.org/x/net/context"
)

func (s *DockerSuite) TestAPIImagesFilter(c *check.C) ***REMOVED***
	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	name := "utest:tag1"
	name2 := "utest/docker:tag2"
	name3 := "utest:5000/docker:tag3"
	for _, n := range []string***REMOVED***name, name2, name3***REMOVED*** ***REMOVED***
		dockerCmd(c, "tag", "busybox", n)
	***REMOVED***
	getImages := func(filter string) []types.ImageSummary ***REMOVED***
		filters := filters.NewArgs()
		filters.Add("reference", filter)
		options := types.ImageListOptions***REMOVED***
			All:     false,
			Filters: filters,
		***REMOVED***
		images, err := cli.ImageList(context.Background(), options)
		c.Assert(err, checker.IsNil)

		return images
	***REMOVED***

	//incorrect number of matches returned
	images := getImages("utest*/*")
	c.Assert(images[0].RepoTags, checker.HasLen, 2)

	images = getImages("utest")
	c.Assert(images[0].RepoTags, checker.HasLen, 1)

	images = getImages("utest*")
	c.Assert(images[0].RepoTags, checker.HasLen, 1)

	images = getImages("*5000*/*")
	c.Assert(images[0].RepoTags, checker.HasLen, 1)
***REMOVED***

func (s *DockerSuite) TestAPIImagesSaveAndLoad(c *check.C) ***REMOVED***
	testRequires(c, Network)
	buildImageSuccessfully(c, "saveandload", build.WithDockerfile("FROM busybox\nENV FOO bar"))
	id := getIDByName(c, "saveandload")

	res, body, err := request.Get("/images/" + id + "/get")
	c.Assert(err, checker.IsNil)
	defer body.Close()
	c.Assert(res.StatusCode, checker.Equals, http.StatusOK)

	dockerCmd(c, "rmi", id)

	res, loadBody, err := request.Post("/images/load", request.RawContent(body), request.ContentType("application/x-tar"))
	c.Assert(err, checker.IsNil)
	defer loadBody.Close()
	c.Assert(res.StatusCode, checker.Equals, http.StatusOK)

	inspectOut := cli.InspectCmd(c, id, cli.Format(".Id")).Combined()
	c.Assert(strings.TrimSpace(string(inspectOut)), checker.Equals, id, check.Commentf("load did not work properly"))
***REMOVED***

func (s *DockerSuite) TestAPIImagesDelete(c *check.C) ***REMOVED***
	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	if testEnv.OSType != "windows" ***REMOVED***
		testRequires(c, Network)
	***REMOVED***
	name := "test-api-images-delete"
	buildImageSuccessfully(c, name, build.WithDockerfile("FROM busybox\nENV FOO bar"))
	id := getIDByName(c, name)

	dockerCmd(c, "tag", name, "test:tag1")

	_, err = cli.ImageRemove(context.Background(), id, types.ImageRemoveOptions***REMOVED******REMOVED***)
	c.Assert(err.Error(), checker.Contains, "unable to delete")

	_, err = cli.ImageRemove(context.Background(), "test:noexist", types.ImageRemoveOptions***REMOVED******REMOVED***)
	c.Assert(err.Error(), checker.Contains, "No such image")

	_, err = cli.ImageRemove(context.Background(), "test:tag1", types.ImageRemoveOptions***REMOVED******REMOVED***)
	c.Assert(err, checker.IsNil)
***REMOVED***

func (s *DockerSuite) TestAPIImagesHistory(c *check.C) ***REMOVED***
	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	if testEnv.OSType != "windows" ***REMOVED***
		testRequires(c, Network)
	***REMOVED***
	name := "test-api-images-history"
	buildImageSuccessfully(c, name, build.WithDockerfile("FROM busybox\nENV FOO bar"))
	id := getIDByName(c, name)

	historydata, err := cli.ImageHistory(context.Background(), id)
	c.Assert(err, checker.IsNil)

	c.Assert(historydata, checker.Not(checker.HasLen), 0)
	c.Assert(historydata[0].Tags[0], checker.Equals, "test-api-images-history:latest")
***REMOVED***

func (s *DockerSuite) TestAPIImagesImportBadSrc(c *check.C) ***REMOVED***
	testRequires(c, Network, SameHostDaemon)

	server := httptest.NewServer(http.NewServeMux())
	defer server.Close()

	tt := []struct ***REMOVED***
		statusExp int
		fromSrc   string
	***REMOVED******REMOVED***
		***REMOVED***http.StatusNotFound, server.URL + "/nofile.tar"***REMOVED***,
		***REMOVED***http.StatusNotFound, strings.TrimPrefix(server.URL, "http://") + "/nofile.tar"***REMOVED***,
		***REMOVED***http.StatusNotFound, strings.TrimPrefix(server.URL, "http://") + "%2Fdata%2Ffile.tar"***REMOVED***,
		***REMOVED***http.StatusInternalServerError, "%2Fdata%2Ffile.tar"***REMOVED***,
	***REMOVED***

	for _, te := range tt ***REMOVED***
		res, _, err := request.Post(strings.Join([]string***REMOVED***"/images/create?fromSrc=", te.fromSrc***REMOVED***, ""), request.JSON)
		c.Assert(err, check.IsNil)
		c.Assert(res.StatusCode, checker.Equals, te.statusExp)
		c.Assert(res.Header.Get("Content-Type"), checker.Equals, "application/json")
	***REMOVED***

***REMOVED***

// #14846
func (s *DockerSuite) TestAPIImagesSearchJSONContentType(c *check.C) ***REMOVED***
	testRequires(c, Network)

	res, b, err := request.Get("/images/search?term=test", request.JSON)
	c.Assert(err, check.IsNil)
	b.Close()
	c.Assert(res.StatusCode, checker.Equals, http.StatusOK)
	c.Assert(res.Header.Get("Content-Type"), checker.Equals, "application/json")
***REMOVED***

// Test case for 30027: image size reported as -1 in v1.12 client against v1.13 daemon.
// This test checks to make sure both v1.12 and v1.13 client against v1.13 daemon get correct `Size` after the fix.
func (s *DockerSuite) TestAPIImagesSizeCompatibility(c *check.C) ***REMOVED***
	cli, err := client.NewEnvClient()
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	images, err := cli.ImageList(context.Background(), types.ImageListOptions***REMOVED******REMOVED***)
	c.Assert(err, checker.IsNil)
	c.Assert(len(images), checker.Not(checker.Equals), 0)
	for _, image := range images ***REMOVED***
		c.Assert(image.Size, checker.Not(checker.Equals), int64(-1))
	***REMOVED***

	type v124Image struct ***REMOVED***
		ID          string `json:"Id"`
		ParentID    string `json:"ParentId"`
		RepoTags    []string
		RepoDigests []string
		Created     int64
		Size        int64
		VirtualSize int64
		Labels      map[string]string
	***REMOVED***

	cli, err = request.NewEnvClientWithVersion("v1.24")
	c.Assert(err, checker.IsNil)
	defer cli.Close()

	v124Images, err := cli.ImageList(context.Background(), types.ImageListOptions***REMOVED******REMOVED***)
	c.Assert(err, checker.IsNil)
	c.Assert(len(v124Images), checker.Not(checker.Equals), 0)
	for _, image := range v124Images ***REMOVED***
		c.Assert(image.Size, checker.Not(checker.Equals), int64(-1))
	***REMOVED***
***REMOVED***
