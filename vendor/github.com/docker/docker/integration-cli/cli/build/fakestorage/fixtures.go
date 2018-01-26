package fakestorage

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/docker/docker/integration-cli/cli"
)

var ensureHTTPServerOnce sync.Once

func ensureHTTPServerImage(t testingT) ***REMOVED***
	var doIt bool
	ensureHTTPServerOnce.Do(func() ***REMOVED***
		doIt = true
	***REMOVED***)

	if !doIt ***REMOVED***
		return
	***REMOVED***

	defer testEnv.ProtectImage(t, "httpserver:latest")

	tmp, err := ioutil.TempDir("", "docker-http-server-test")
	if err != nil ***REMOVED***
		t.Fatalf("could not build http server: %v", err)
	***REMOVED***
	defer os.RemoveAll(tmp)

	goos := testEnv.OSType
	if goos == "" ***REMOVED***
		goos = "linux"
	***REMOVED***
	goarch := os.Getenv("DOCKER_ENGINE_GOARCH")
	if goarch == "" ***REMOVED***
		goarch = "amd64"
	***REMOVED***

	cpCmd, lookErr := exec.LookPath("cp")
	if lookErr != nil ***REMOVED***
		t.Fatalf("could not build http server: %v", lookErr)
	***REMOVED***

	if _, err = os.Stat("../contrib/httpserver/httpserver"); os.IsNotExist(err) ***REMOVED***
		goCmd, lookErr := exec.LookPath("go")
		if lookErr != nil ***REMOVED***
			t.Fatalf("could not build http server: %v", lookErr)
		***REMOVED***

		cmd := exec.Command(goCmd, "build", "-o", filepath.Join(tmp, "httpserver"), "github.com/docker/docker/contrib/httpserver")
		cmd.Env = append(os.Environ(), []string***REMOVED***
			"CGO_ENABLED=0",
			"GOOS=" + goos,
			"GOARCH=" + goarch,
		***REMOVED***...)
		var out []byte
		if out, err = cmd.CombinedOutput(); err != nil ***REMOVED***
			t.Fatalf("could not build http server: %s", string(out))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if out, err := exec.Command(cpCmd, "../contrib/httpserver/httpserver", filepath.Join(tmp, "httpserver")).CombinedOutput(); err != nil ***REMOVED***
			t.Fatalf("could not copy http server: %v", string(out))
		***REMOVED***
	***REMOVED***

	if out, err := exec.Command(cpCmd, "../contrib/httpserver/Dockerfile", filepath.Join(tmp, "Dockerfile")).CombinedOutput(); err != nil ***REMOVED***
		t.Fatalf("could not build http server: %v", string(out))
	***REMOVED***

	cli.DockerCmd(t, "build", "-q", "-t", "httpserver", tmp)
***REMOVED***
