package fakestorage

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"

	"github.com/docker/docker/integration-cli/cli"
	"github.com/docker/docker/integration-cli/cli/build"
	"github.com/docker/docker/integration-cli/cli/build/fakecontext"
	"github.com/docker/docker/integration-cli/request"
	"github.com/docker/docker/internal/test/environment"
	"github.com/docker/docker/internal/testutil"
	"github.com/stretchr/testify/require"
)

var testEnv *environment.Execution

type testingT interface ***REMOVED***
	require.TestingT
	logT
	Fatal(args ...interface***REMOVED******REMOVED***)
	Fatalf(string, ...interface***REMOVED******REMOVED***)
***REMOVED***

type logT interface ***REMOVED***
	Logf(string, ...interface***REMOVED******REMOVED***)
***REMOVED***

// Fake is a static file server. It might be running locally or remotely
// on test host.
type Fake interface ***REMOVED***
	Close() error
	URL() string
	CtxDir() string
***REMOVED***

// SetTestEnvironment sets a static test environment
// TODO: decouple this package from environment
func SetTestEnvironment(env *environment.Execution) ***REMOVED***
	testEnv = env
***REMOVED***

// New returns a static file server that will be use as build context.
func New(t testingT, dir string, modifiers ...func(*fakecontext.Fake) error) Fake ***REMOVED***
	if testEnv == nil ***REMOVED***
		t.Fatal("fakstorage package requires SetTestEnvironment() to be called before use.")
	***REMOVED***
	ctx := fakecontext.New(t, dir, modifiers...)
	if testEnv.IsLocalDaemon() ***REMOVED***
		return newLocalFakeStorage(ctx)
	***REMOVED***
	return newRemoteFileServer(t, ctx)
***REMOVED***

// localFileStorage is a file storage on the running machine
type localFileStorage struct ***REMOVED***
	*fakecontext.Fake
	*httptest.Server
***REMOVED***

func (s *localFileStorage) URL() string ***REMOVED***
	return s.Server.URL
***REMOVED***

func (s *localFileStorage) CtxDir() string ***REMOVED***
	return s.Fake.Dir
***REMOVED***

func (s *localFileStorage) Close() error ***REMOVED***
	defer s.Server.Close()
	return s.Fake.Close()
***REMOVED***

func newLocalFakeStorage(ctx *fakecontext.Fake) *localFileStorage ***REMOVED***
	handler := http.FileServer(http.Dir(ctx.Dir))
	server := httptest.NewServer(handler)
	return &localFileStorage***REMOVED***
		Fake:   ctx,
		Server: server,
	***REMOVED***
***REMOVED***

// remoteFileServer is a containerized static file server started on the remote
// testing machine to be used in URL-accepting docker build functionality.
type remoteFileServer struct ***REMOVED***
	host      string // hostname/port web server is listening to on docker host e.g. 0.0.0.0:43712
	container string
	image     string
	ctx       *fakecontext.Fake
***REMOVED***

func (f *remoteFileServer) URL() string ***REMOVED***
	u := url.URL***REMOVED***
		Scheme: "http",
		Host:   f.host***REMOVED***
	return u.String()
***REMOVED***

func (f *remoteFileServer) CtxDir() string ***REMOVED***
	return f.ctx.Dir
***REMOVED***

func (f *remoteFileServer) Close() error ***REMOVED***
	defer func() ***REMOVED***
		if f.ctx != nil ***REMOVED***
			f.ctx.Close()
		***REMOVED***
		if f.image != "" ***REMOVED***
			if err := cli.Docker(cli.Args("rmi", "-f", f.image)).Error; err != nil ***REMOVED***
				fmt.Fprintf(os.Stderr, "Error closing remote file server : %v\n", err)
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	if f.container == "" ***REMOVED***
		return nil
	***REMOVED***
	return cli.Docker(cli.Args("rm", "-fv", f.container)).Error
***REMOVED***

func newRemoteFileServer(t testingT, ctx *fakecontext.Fake) *remoteFileServer ***REMOVED***
	var (
		image     = fmt.Sprintf("fileserver-img-%s", strings.ToLower(testutil.GenerateRandomAlphaOnlyString(10)))
		container = fmt.Sprintf("fileserver-cnt-%s", strings.ToLower(testutil.GenerateRandomAlphaOnlyString(10)))
	)

	ensureHTTPServerImage(t)

	// Build the image
	if err := ctx.Add("Dockerfile", `FROM httpserver
COPY . /static`); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	cli.BuildCmd(t, image, build.WithoutCache, build.WithExternalBuildContext(ctx))

	// Start the container
	cli.DockerCmd(t, "run", "-d", "-P", "--name", container, image)

	// Find out the system assigned port
	out := cli.DockerCmd(t, "port", container, "80/tcp").Combined()
	fileserverHostPort := strings.Trim(out, "\n")
	_, port, err := net.SplitHostPort(fileserverHostPort)
	if err != nil ***REMOVED***
		t.Fatalf("unable to parse file server host:port: %v", err)
	***REMOVED***

	dockerHostURL, err := url.Parse(request.DaemonHost())
	if err != nil ***REMOVED***
		t.Fatalf("unable to parse daemon host URL: %v", err)
	***REMOVED***

	host, _, err := net.SplitHostPort(dockerHostURL.Host)
	if err != nil ***REMOVED***
		t.Fatalf("unable to parse docker daemon host:port: %v", err)
	***REMOVED***

	return &remoteFileServer***REMOVED***
		container: container,
		image:     image,
		host:      fmt.Sprintf("%s:%s", host, port),
		ctx:       ctx***REMOVED***
***REMOVED***
