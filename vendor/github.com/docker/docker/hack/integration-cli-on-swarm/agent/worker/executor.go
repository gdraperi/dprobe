package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// testChunkExecutor executes integration-cli binary.
// image needs to be the worker image itself. testFlags are OR-set of regexp for filtering tests.
type testChunkExecutor func(image string, tests []string) (int64, string, error)

func dryTestChunkExecutor() testChunkExecutor ***REMOVED***
	return func(image string, tests []string) (int64, string, error) ***REMOVED***
		return 0, fmt.Sprintf("DRY RUN (image=%q, tests=%v)", image, tests), nil
	***REMOVED***
***REMOVED***

// privilegedTestChunkExecutor invokes a privileged container from the worker
// service via bind-mounted API socket so as to execute the test chunk
func privilegedTestChunkExecutor(autoRemove bool) testChunkExecutor ***REMOVED***
	return func(image string, tests []string) (int64, string, error) ***REMOVED***
		cli, err := client.NewEnvClient()
		if err != nil ***REMOVED***
			return 0, "", err
		***REMOVED***
		// propagate variables from the host (needs to be defined in the compose file)
		experimental := os.Getenv("DOCKER_EXPERIMENTAL")
		graphdriver := os.Getenv("DOCKER_GRAPHDRIVER")
		if graphdriver == "" ***REMOVED***
			info, err := cli.Info(context.Background())
			if err != nil ***REMOVED***
				return 0, "", err
			***REMOVED***
			graphdriver = info.Driver
		***REMOVED***
		// `daemon_dest` is similar to `$DEST` (e.g. `bundles/VERSION/test-integration`)
		// but it exists outside of `bundles` so as to make `$DOCKER_GRAPHDRIVER` work.
		//
		// Without this hack, `$DOCKER_GRAPHDRIVER` fails because of (e.g.) `overlay2 is not supported over overlayfs`
		//
		// see integration-cli/daemon/daemon.go
		daemonDest := "/daemon_dest"
		config := container.Config***REMOVED***
			Image: image,
			Env: []string***REMOVED***
				"TESTFLAGS=-check.f " + strings.Join(tests, "|"),
				"KEEPBUNDLE=1",
				"DOCKER_INTEGRATION_TESTS_VERIFIED=1", // for avoiding rebuilding integration-cli
				"DOCKER_EXPERIMENTAL=" + experimental,
				"DOCKER_GRAPHDRIVER=" + graphdriver,
				"DOCKER_INTEGRATION_DAEMON_DEST=" + daemonDest,
			***REMOVED***,
			Labels: map[string]string***REMOVED***
				"org.dockerproject.integration-cli-on-swarm":         "",
				"org.dockerproject.integration-cli-on-swarm.comment": "this non-service container is created for running privileged programs on Swarm. you can remove this container manually if the corresponding service is already stopped.",
			***REMOVED***,
			Entrypoint: []string***REMOVED***"hack/dind"***REMOVED***,
			Cmd:        []string***REMOVED***"hack/make.sh", "test-integration"***REMOVED***,
		***REMOVED***
		hostConfig := container.HostConfig***REMOVED***
			AutoRemove: autoRemove,
			Privileged: true,
			Mounts: []mount.Mount***REMOVED***
				***REMOVED***
					Type:   mount.TypeVolume,
					Target: daemonDest,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***
		id, stream, err := runContainer(context.Background(), cli, config, hostConfig)
		if err != nil ***REMOVED***
			return 0, "", err
		***REMOVED***
		var b bytes.Buffer
		teeContainerStream(&b, os.Stdout, os.Stderr, stream)
		resultC, errC := cli.ContainerWait(context.Background(), id, "")
		select ***REMOVED***
		case err := <-errC:
			return 0, "", err
		case result := <-resultC:
			return result.StatusCode, b.String(), nil
		***REMOVED***
	***REMOVED***
***REMOVED***

func runContainer(ctx context.Context, cli *client.Client, config container.Config, hostConfig container.HostConfig) (string, io.ReadCloser, error) ***REMOVED***
	created, err := cli.ContainerCreate(context.Background(),
		&config, &hostConfig, nil, "")
	if err != nil ***REMOVED***
		return "", nil, err
	***REMOVED***
	if err = cli.ContainerStart(ctx, created.ID, types.ContainerStartOptions***REMOVED******REMOVED***); err != nil ***REMOVED***
		return "", nil, err
	***REMOVED***
	stream, err := cli.ContainerLogs(ctx,
		created.ID,
		types.ContainerLogsOptions***REMOVED***
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
		***REMOVED***)
	return created.ID, stream, err
***REMOVED***

func teeContainerStream(w, stdout, stderr io.Writer, stream io.ReadCloser) ***REMOVED***
	stdcopy.StdCopy(io.MultiWriter(w, stdout), io.MultiWriter(w, stderr), stream)
	stream.Close()
***REMOVED***
