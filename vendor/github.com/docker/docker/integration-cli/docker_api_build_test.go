package main

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/integration-cli/checker"
	"github.com/docker/docker/integration-cli/cli/build/fakecontext"
	"github.com/docker/docker/integration-cli/cli/build/fakegit"
	"github.com/docker/docker/integration-cli/cli/build/fakestorage"
	"github.com/docker/docker/integration-cli/request"
	"github.com/go-check/check"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/filesync"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

func (s *DockerSuite) TestBuildAPIDockerFileRemote(c *check.C) ***REMOVED***
	testRequires(c, NotUserNamespace)

	var testD string
	if testEnv.OSType == "windows" ***REMOVED***
		testD = `FROM busybox
RUN find / -name ba*
RUN find /tmp/`
	***REMOVED*** else ***REMOVED***
		// -xdev is required because sysfs can cause EPERM
		testD = `FROM busybox
RUN find / -xdev -name ba*
RUN find /tmp/`
	***REMOVED***
	server := fakestorage.New(c, "", fakecontext.WithFiles(map[string]string***REMOVED***"testD": testD***REMOVED***))
	defer server.Close()

	res, body, err := request.Post("/build?dockerfile=baz&remote="+server.URL()+"/testD", request.JSON)
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusOK)

	buf, err := request.ReadBody(body)
	c.Assert(err, checker.IsNil)

	// Make sure Dockerfile exists.
	// Make sure 'baz' doesn't exist ANYWHERE despite being mentioned in the URL
	out := string(buf)
	c.Assert(out, checker.Contains, "RUN find /tmp")
	c.Assert(out, checker.Not(checker.Contains), "baz")
***REMOVED***

func (s *DockerSuite) TestBuildAPIRemoteTarballContext(c *check.C) ***REMOVED***
	buffer := new(bytes.Buffer)
	tw := tar.NewWriter(buffer)
	defer tw.Close()

	dockerfile := []byte("FROM busybox")
	err := tw.WriteHeader(&tar.Header***REMOVED***
		Name: "Dockerfile",
		Size: int64(len(dockerfile)),
	***REMOVED***)
	// failed to write tar file header
	c.Assert(err, checker.IsNil)

	_, err = tw.Write(dockerfile)
	// failed to write tar file content
	c.Assert(err, checker.IsNil)

	// failed to close tar archive
	c.Assert(tw.Close(), checker.IsNil)

	server := fakestorage.New(c, "", fakecontext.WithBinaryFiles(map[string]*bytes.Buffer***REMOVED***
		"testT.tar": buffer,
	***REMOVED***))
	defer server.Close()

	res, b, err := request.Post("/build?remote="+server.URL()+"/testT.tar", request.ContentType("application/tar"))
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusOK)
	b.Close()
***REMOVED***

func (s *DockerSuite) TestBuildAPIRemoteTarballContextWithCustomDockerfile(c *check.C) ***REMOVED***
	buffer := new(bytes.Buffer)
	tw := tar.NewWriter(buffer)
	defer tw.Close()

	dockerfile := []byte(`FROM busybox
RUN echo 'wrong'`)
	err := tw.WriteHeader(&tar.Header***REMOVED***
		Name: "Dockerfile",
		Size: int64(len(dockerfile)),
	***REMOVED***)
	// failed to write tar file header
	c.Assert(err, checker.IsNil)

	_, err = tw.Write(dockerfile)
	// failed to write tar file content
	c.Assert(err, checker.IsNil)

	custom := []byte(`FROM busybox
RUN echo 'right'
`)
	err = tw.WriteHeader(&tar.Header***REMOVED***
		Name: "custom",
		Size: int64(len(custom)),
	***REMOVED***)

	// failed to write tar file header
	c.Assert(err, checker.IsNil)

	_, err = tw.Write(custom)
	// failed to write tar file content
	c.Assert(err, checker.IsNil)

	// failed to close tar archive
	c.Assert(tw.Close(), checker.IsNil)

	server := fakestorage.New(c, "", fakecontext.WithBinaryFiles(map[string]*bytes.Buffer***REMOVED***
		"testT.tar": buffer,
	***REMOVED***))
	defer server.Close()

	url := "/build?dockerfile=custom&remote=" + server.URL() + "/testT.tar"
	res, body, err := request.Post(url, request.ContentType("application/tar"))
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusOK)

	defer body.Close()
	content, err := request.ReadBody(body)
	c.Assert(err, checker.IsNil)

	// Build used the wrong dockerfile.
	c.Assert(string(content), checker.Not(checker.Contains), "wrong")
***REMOVED***

func (s *DockerSuite) TestBuildAPILowerDockerfile(c *check.C) ***REMOVED***
	git := fakegit.New(c, "repo", map[string]string***REMOVED***
		"dockerfile": `FROM busybox
RUN echo from dockerfile`,
	***REMOVED***, false)
	defer git.Close()

	res, body, err := request.Post("/build?remote="+git.RepoURL, request.JSON)
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusOK)

	buf, err := request.ReadBody(body)
	c.Assert(err, checker.IsNil)

	out := string(buf)
	c.Assert(out, checker.Contains, "from dockerfile")
***REMOVED***

func (s *DockerSuite) TestBuildAPIBuildGitWithF(c *check.C) ***REMOVED***
	git := fakegit.New(c, "repo", map[string]string***REMOVED***
		"baz": `FROM busybox
RUN echo from baz`,
		"Dockerfile": `FROM busybox
RUN echo from Dockerfile`,
	***REMOVED***, false)
	defer git.Close()

	// Make sure it tries to 'dockerfile' query param value
	res, body, err := request.Post("/build?dockerfile=baz&remote="+git.RepoURL, request.JSON)
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusOK)

	buf, err := request.ReadBody(body)
	c.Assert(err, checker.IsNil)

	out := string(buf)
	c.Assert(out, checker.Contains, "from baz")
***REMOVED***

func (s *DockerSuite) TestBuildAPIDoubleDockerfile(c *check.C) ***REMOVED***
	testRequires(c, UnixCli) // dockerfile overwrites Dockerfile on Windows
	git := fakegit.New(c, "repo", map[string]string***REMOVED***
		"Dockerfile": `FROM busybox
RUN echo from Dockerfile`,
		"dockerfile": `FROM busybox
RUN echo from dockerfile`,
	***REMOVED***, false)
	defer git.Close()

	// Make sure it tries to 'dockerfile' query param value
	res, body, err := request.Post("/build?remote="+git.RepoURL, request.JSON)
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusOK)

	buf, err := request.ReadBody(body)
	c.Assert(err, checker.IsNil)

	out := string(buf)
	c.Assert(out, checker.Contains, "from Dockerfile")
***REMOVED***

func (s *DockerSuite) TestBuildAPIUnnormalizedTarPaths(c *check.C) ***REMOVED***
	// Make sure that build context tars with entries of the form
	// x/./y don't cause caching false positives.

	buildFromTarContext := func(fileContents []byte) string ***REMOVED***
		buffer := new(bytes.Buffer)
		tw := tar.NewWriter(buffer)
		defer tw.Close()

		dockerfile := []byte(`FROM busybox
	COPY dir /dir/`)
		err := tw.WriteHeader(&tar.Header***REMOVED***
			Name: "Dockerfile",
			Size: int64(len(dockerfile)),
		***REMOVED***)
		//failed to write tar file header
		c.Assert(err, checker.IsNil)

		_, err = tw.Write(dockerfile)
		// failed to write Dockerfile in tar file content
		c.Assert(err, checker.IsNil)

		err = tw.WriteHeader(&tar.Header***REMOVED***
			Name: "dir/./file",
			Size: int64(len(fileContents)),
		***REMOVED***)
		//failed to write tar file header
		c.Assert(err, checker.IsNil)

		_, err = tw.Write(fileContents)
		// failed to write file contents in tar file content
		c.Assert(err, checker.IsNil)

		// failed to close tar archive
		c.Assert(tw.Close(), checker.IsNil)

		res, body, err := request.Post("/build", request.RawContent(ioutil.NopCloser(buffer)), request.ContentType("application/x-tar"))
		c.Assert(err, checker.IsNil)
		c.Assert(res.StatusCode, checker.Equals, http.StatusOK)

		out, err := request.ReadBody(body)
		c.Assert(err, checker.IsNil)
		lines := strings.Split(string(out), "\n")
		c.Assert(len(lines), checker.GreaterThan, 1)
		c.Assert(lines[len(lines)-2], checker.Matches, ".*Successfully built [0-9a-f]***REMOVED***12***REMOVED***.*")

		re := regexp.MustCompile("Successfully built ([0-9a-f]***REMOVED***12***REMOVED***)")
		matches := re.FindStringSubmatch(lines[len(lines)-2])
		return matches[1]
	***REMOVED***

	imageA := buildFromTarContext([]byte("abc"))
	imageB := buildFromTarContext([]byte("def"))

	c.Assert(imageA, checker.Not(checker.Equals), imageB)
***REMOVED***

func (s *DockerSuite) TestBuildOnBuildWithCopy(c *check.C) ***REMOVED***
	dockerfile := `
		FROM ` + minimalBaseImage() + ` as onbuildbase
		ONBUILD COPY file /file

		FROM onbuildbase
	`
	ctx := fakecontext.New(c, "",
		fakecontext.WithDockerfile(dockerfile),
		fakecontext.WithFile("file", "some content"),
	)
	defer ctx.Close()

	res, body, err := request.Post(
		"/build",
		request.RawContent(ctx.AsTarReader(c)),
		request.ContentType("application/x-tar"))
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusOK)

	out, err := request.ReadBody(body)
	c.Assert(err, checker.IsNil)
	c.Assert(string(out), checker.Contains, "Successfully built")
***REMOVED***

func (s *DockerSuite) TestBuildOnBuildCache(c *check.C) ***REMOVED***
	build := func(dockerfile string) []byte ***REMOVED***
		ctx := fakecontext.New(c, "",
			fakecontext.WithDockerfile(dockerfile),
		)
		defer ctx.Close()

		res, body, err := request.Post(
			"/build",
			request.RawContent(ctx.AsTarReader(c)),
			request.ContentType("application/x-tar"))
		require.NoError(c, err)
		assert.Equal(c, http.StatusOK, res.StatusCode)

		out, err := request.ReadBody(body)
		require.NoError(c, err)
		assert.Contains(c, string(out), "Successfully built")
		return out
	***REMOVED***

	dockerfile := `
		FROM ` + minimalBaseImage() + ` as onbuildbase
		ENV something=bar
		ONBUILD ENV foo=bar
	`
	build(dockerfile)

	dockerfile += "FROM onbuildbase"
	out := build(dockerfile)

	imageIDs := getImageIDsFromBuild(c, out)
	assert.Len(c, imageIDs, 2)
	parentID, childID := imageIDs[0], imageIDs[1]

	client, err := request.NewClient()
	require.NoError(c, err)

	// check parentID is correct
	image, _, err := client.ImageInspectWithRaw(context.Background(), childID)
	require.NoError(c, err)
	assert.Equal(c, parentID, image.Parent)
***REMOVED***

func (s *DockerRegistrySuite) TestBuildCopyFromForcePull(c *check.C) ***REMOVED***
	client, err := request.NewClient()
	require.NoError(c, err)

	repoName := fmt.Sprintf("%v/dockercli/busybox", privateRegistryURL)
	// tag the image to upload it to the private registry
	err = client.ImageTag(context.TODO(), "busybox", repoName)
	assert.Nil(c, err)
	// push the image to the registry
	rc, err := client.ImagePush(context.TODO(), repoName, types.ImagePushOptions***REMOVED***RegistryAuth: "***REMOVED******REMOVED***"***REMOVED***)
	assert.Nil(c, err)
	_, err = io.Copy(ioutil.Discard, rc)
	assert.Nil(c, err)

	dockerfile := fmt.Sprintf(`
		FROM %s AS foo
		RUN touch abc
		FROM %s
		COPY --from=foo /abc /
		`, repoName, repoName)

	ctx := fakecontext.New(c, "",
		fakecontext.WithDockerfile(dockerfile),
	)
	defer ctx.Close()

	res, body, err := request.Post(
		"/build?pull=1",
		request.RawContent(ctx.AsTarReader(c)),
		request.ContentType("application/x-tar"))
	require.NoError(c, err)
	assert.Equal(c, http.StatusOK, res.StatusCode)

	out, err := request.ReadBody(body)
	require.NoError(c, err)
	assert.Contains(c, string(out), "Successfully built")
***REMOVED***

func (s *DockerSuite) TestBuildAddRemoteNoDecompress(c *check.C) ***REMOVED***
	buffer := new(bytes.Buffer)
	tw := tar.NewWriter(buffer)
	dt := []byte("contents")
	err := tw.WriteHeader(&tar.Header***REMOVED***
		Name:     "foo",
		Size:     int64(len(dt)),
		Mode:     0600,
		Typeflag: tar.TypeReg,
	***REMOVED***)
	require.NoError(c, err)
	_, err = tw.Write(dt)
	require.NoError(c, err)
	err = tw.Close()
	require.NoError(c, err)

	server := fakestorage.New(c, "", fakecontext.WithBinaryFiles(map[string]*bytes.Buffer***REMOVED***
		"test.tar": buffer,
	***REMOVED***))
	defer server.Close()

	dockerfile := fmt.Sprintf(`
		FROM busybox
		ADD %s/test.tar /
		RUN [ -f test.tar ]
		`, server.URL())

	ctx := fakecontext.New(c, "",
		fakecontext.WithDockerfile(dockerfile),
	)
	defer ctx.Close()

	res, body, err := request.Post(
		"/build",
		request.RawContent(ctx.AsTarReader(c)),
		request.ContentType("application/x-tar"))
	require.NoError(c, err)
	assert.Equal(c, http.StatusOK, res.StatusCode)

	out, err := request.ReadBody(body)
	require.NoError(c, err)
	assert.Contains(c, string(out), "Successfully built")
***REMOVED***

func (s *DockerSuite) TestBuildChownOnCopy(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	dockerfile := `FROM busybox
		RUN echo 'test1:x:1001:1001::/bin:/bin/false' >> /etc/passwd
		RUN echo 'test1:x:1001:' >> /etc/group
		RUN echo 'test2:x:1002:' >> /etc/group
		COPY --chown=test1:1002 . /new_dir
		RUN ls -l /
		RUN [ $(ls -l / | grep new_dir | awk '***REMOVED***print $3":"$4***REMOVED***') = 'test1:test2' ]
		RUN [ $(ls -nl / | grep new_dir | awk '***REMOVED***print $3":"$4***REMOVED***') = '1001:1002' ]
	`
	ctx := fakecontext.New(c, "",
		fakecontext.WithDockerfile(dockerfile),
		fakecontext.WithFile("test_file1", "some test content"),
	)
	defer ctx.Close()

	res, body, err := request.Post(
		"/build",
		request.RawContent(ctx.AsTarReader(c)),
		request.ContentType("application/x-tar"))
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusOK)

	out, err := request.ReadBody(body)
	require.NoError(c, err)
	assert.Contains(c, string(out), "Successfully built")
***REMOVED***

func (s *DockerSuite) TestBuildCopyCacheOnFileChange(c *check.C) ***REMOVED***

	dockerfile := `FROM busybox
COPY file /file`

	ctx1 := fakecontext.New(c, "",
		fakecontext.WithDockerfile(dockerfile),
		fakecontext.WithFile("file", "foo"))
	ctx2 := fakecontext.New(c, "",
		fakecontext.WithDockerfile(dockerfile),
		fakecontext.WithFile("file", "bar"))

	var build = func(ctx *fakecontext.Fake) string ***REMOVED***
		res, body, err := request.Post("/build",
			request.RawContent(ctx.AsTarReader(c)),
			request.ContentType("application/x-tar"))

		require.NoError(c, err)
		assert.Equal(c, http.StatusOK, res.StatusCode)

		out, err := request.ReadBody(body)
		require.NoError(c, err)

		ids := getImageIDsFromBuild(c, out)
		return ids[len(ids)-1]
	***REMOVED***

	id1 := build(ctx1)
	id2 := build(ctx1)
	id3 := build(ctx2)

	if id1 != id2 ***REMOVED***
		c.Fatal("didn't use the cache")
	***REMOVED***
	if id1 == id3 ***REMOVED***
		c.Fatal("COPY With different source file should not share same cache")
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestBuildAddCacheOnFileChange(c *check.C) ***REMOVED***

	dockerfile := `FROM busybox
ADD file /file`

	ctx1 := fakecontext.New(c, "",
		fakecontext.WithDockerfile(dockerfile),
		fakecontext.WithFile("file", "foo"))
	ctx2 := fakecontext.New(c, "",
		fakecontext.WithDockerfile(dockerfile),
		fakecontext.WithFile("file", "bar"))

	var build = func(ctx *fakecontext.Fake) string ***REMOVED***
		res, body, err := request.Post("/build",
			request.RawContent(ctx.AsTarReader(c)),
			request.ContentType("application/x-tar"))

		require.NoError(c, err)
		assert.Equal(c, http.StatusOK, res.StatusCode)

		out, err := request.ReadBody(body)
		require.NoError(c, err)

		ids := getImageIDsFromBuild(c, out)
		return ids[len(ids)-1]
	***REMOVED***

	id1 := build(ctx1)
	id2 := build(ctx1)
	id3 := build(ctx2)

	if id1 != id2 ***REMOVED***
		c.Fatal("didn't use the cache")
	***REMOVED***
	if id1 == id3 ***REMOVED***
		c.Fatal("COPY With different source file should not share same cache")
	***REMOVED***
***REMOVED***

func (s *DockerSuite) TestBuildWithSession(c *check.C) ***REMOVED***
	testRequires(c, ExperimentalDaemon)

	dockerfile := `
		FROM busybox
		COPY file /
		RUN cat /file
	`

	fctx := fakecontext.New(c, "",
		fakecontext.WithFile("file", "some content"),
	)
	defer fctx.Close()

	out := testBuildWithSession(c, fctx.Dir, dockerfile)
	assert.Contains(c, out, "some content")

	fctx.Add("second", "contentcontent")

	dockerfile += `
	COPY second /
	RUN cat /second
	`

	out = testBuildWithSession(c, fctx.Dir, dockerfile)
	assert.Equal(c, strings.Count(out, "Using cache"), 2)
	assert.Contains(c, out, "contentcontent")

	client, err := request.NewClient()
	require.NoError(c, err)

	du, err := client.DiskUsage(context.TODO())
	assert.Nil(c, err)
	assert.True(c, du.BuilderSize > 10)

	out = testBuildWithSession(c, fctx.Dir, dockerfile)
	assert.Equal(c, strings.Count(out, "Using cache"), 4)

	du2, err := client.DiskUsage(context.TODO())
	assert.Nil(c, err)
	assert.Equal(c, du.BuilderSize, du2.BuilderSize)

	// rebuild with regular tar, confirm cache still applies
	fctx.Add("Dockerfile", dockerfile)
	res, body, err := request.Post(
		"/build",
		request.RawContent(fctx.AsTarReader(c)),
		request.ContentType("application/x-tar"))
	require.NoError(c, err)
	assert.Equal(c, http.StatusOK, res.StatusCode)

	outBytes, err := request.ReadBody(body)
	require.NoError(c, err)
	assert.Contains(c, string(outBytes), "Successfully built")
	assert.Equal(c, strings.Count(string(outBytes), "Using cache"), 4)

	_, err = client.BuildCachePrune(context.TODO())
	assert.Nil(c, err)

	du, err = client.DiskUsage(context.TODO())
	assert.Nil(c, err)
	assert.Equal(c, du.BuilderSize, int64(0))
***REMOVED***

func testBuildWithSession(c *check.C, dir, dockerfile string) (outStr string) ***REMOVED***
	client, err := request.NewClient()
	require.NoError(c, err)

	sess, err := session.NewSession("foo1", "foo")
	assert.Nil(c, err)

	fsProvider := filesync.NewFSSyncProvider([]filesync.SyncedDir***REMOVED***
		***REMOVED***Dir: dir***REMOVED***,
	***REMOVED***)
	sess.Allow(fsProvider)

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error ***REMOVED***
		return sess.Run(ctx, client.DialSession)
	***REMOVED***)

	g.Go(func() error ***REMOVED***
		res, body, err := request.Post("/build?remote=client-session&session="+sess.ID(), func(req *http.Request) error ***REMOVED***
			req.Body = ioutil.NopCloser(strings.NewReader(dockerfile))
			return nil
		***REMOVED***)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		assert.Equal(c, res.StatusCode, http.StatusOK)
		out, err := request.ReadBody(body)
		require.NoError(c, err)
		assert.Contains(c, string(out), "Successfully built")
		sess.Close()
		outStr = string(out)
		return nil
	***REMOVED***)

	err = g.Wait()
	assert.Nil(c, err)
	return
***REMOVED***

func (s *DockerSuite) TestBuildScratchCopy(c *check.C) ***REMOVED***
	testRequires(c, DaemonIsLinux)
	dockerfile := `FROM scratch
ADD Dockerfile /
ENV foo bar`
	ctx := fakecontext.New(c, "",
		fakecontext.WithDockerfile(dockerfile),
	)
	defer ctx.Close()

	res, body, err := request.Post(
		"/build",
		request.RawContent(ctx.AsTarReader(c)),
		request.ContentType("application/x-tar"))
	c.Assert(err, checker.IsNil)
	c.Assert(res.StatusCode, checker.Equals, http.StatusOK)

	out, err := request.ReadBody(body)
	require.NoError(c, err)
	assert.Contains(c, string(out), "Successfully built")
***REMOVED***

type buildLine struct ***REMOVED***
	Stream string
	Aux    struct ***REMOVED***
		ID string
	***REMOVED***
***REMOVED***

func getImageIDsFromBuild(c *check.C, output []byte) []string ***REMOVED***
	ids := []string***REMOVED******REMOVED***
	for _, line := range bytes.Split(output, []byte("\n")) ***REMOVED***
		if len(line) == 0 ***REMOVED***
			continue
		***REMOVED***
		entry := buildLine***REMOVED******REMOVED***
		require.NoError(c, json.Unmarshal(line, &entry))
		if entry.Aux.ID != "" ***REMOVED***
			ids = append(ids, entry.Aux.ID)
		***REMOVED***
	***REMOVED***
	return ids
***REMOVED***
