package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/sirupsen/logrus"
)

const (
	defaultStackName       = "integration-cli-on-swarm"
	defaultVolumeName      = "integration-cli-on-swarm"
	defaultMasterImageName = "integration-cli-master"
	defaultWorkerImageName = "integration-cli-worker"
)

func main() ***REMOVED***
	rc, err := xmain()
	if err != nil ***REMOVED***
		logrus.Fatalf("fatal error: %v", err)
	***REMOVED***
	os.Exit(rc)
***REMOVED***

func xmain() (int, error) ***REMOVED***
	// Should we use cobra maybe?
	replicas := flag.Int("replicas", 1, "Number of worker service replica")
	chunks := flag.Int("chunks", 0, "Number of test chunks executed in batch (0 == replicas)")
	pushWorkerImage := flag.String("push-worker-image", "", "Push the worker image to the registry. Required for distributed execution. (empty == not to push)")
	shuffle := flag.Bool("shuffle", false, "Shuffle the input so as to mitigate makespan nonuniformity")
	// flags below are rarely used
	randSeed := flag.Int64("rand-seed", int64(0), "Random seed used for shuffling (0 == current time)")
	filtersFile := flag.String("filters-file", "", "Path to optional file composed of `-check.f` filter strings")
	dryRun := flag.Bool("dry-run", false, "Dry run")
	keepExecutor := flag.Bool("keep-executor", false, "Do not auto-remove executor containers, which is used for running privileged programs on Swarm")
	flag.Parse()
	if *chunks == 0 ***REMOVED***
		*chunks = *replicas
	***REMOVED***
	if *randSeed == int64(0) ***REMOVED***
		*randSeed = time.Now().UnixNano()
	***REMOVED***
	cli, err := client.NewEnvClient()
	if err != nil ***REMOVED***
		return 1, err
	***REMOVED***
	if hasStack(cli, defaultStackName) ***REMOVED***
		logrus.Infof("Removing stack %s", defaultStackName)
		removeStack(cli, defaultStackName)
	***REMOVED***
	if hasVolume(cli, defaultVolumeName) ***REMOVED***
		logrus.Infof("Removing volume %s", defaultVolumeName)
		removeVolume(cli, defaultVolumeName)
	***REMOVED***
	if err = ensureImages(cli, []string***REMOVED***defaultWorkerImageName, defaultMasterImageName***REMOVED***); err != nil ***REMOVED***
		return 1, err
	***REMOVED***
	workerImageForStack := defaultWorkerImageName
	if *pushWorkerImage != "" ***REMOVED***
		logrus.Infof("Pushing %s to %s", defaultWorkerImageName, *pushWorkerImage)
		if err = pushImage(cli, *pushWorkerImage, defaultWorkerImageName); err != nil ***REMOVED***
			return 1, err
		***REMOVED***
		workerImageForStack = *pushWorkerImage
	***REMOVED***
	compose, err := createCompose("", cli, composeOptions***REMOVED***
		Replicas:     *replicas,
		Chunks:       *chunks,
		MasterImage:  defaultMasterImageName,
		WorkerImage:  workerImageForStack,
		Volume:       defaultVolumeName,
		Shuffle:      *shuffle,
		RandSeed:     *randSeed,
		DryRun:       *dryRun,
		KeepExecutor: *keepExecutor,
	***REMOVED***)
	if err != nil ***REMOVED***
		return 1, err
	***REMOVED***
	filters, err := filtersBytes(*filtersFile)
	if err != nil ***REMOVED***
		return 1, err
	***REMOVED***
	logrus.Infof("Creating volume %s with input data", defaultVolumeName)
	if err = createVolumeWithData(cli,
		defaultVolumeName,
		map[string][]byte***REMOVED***"/input": filters***REMOVED***,
		defaultMasterImageName); err != nil ***REMOVED***
		return 1, err
	***REMOVED***
	logrus.Infof("Deploying stack %s from %s", defaultStackName, compose)
	defer func() ***REMOVED***
		logrus.Infof("NOTE: You may want to inspect or clean up following resources:")
		logrus.Infof(" - Stack: %s", defaultStackName)
		logrus.Infof(" - Volume: %s", defaultVolumeName)
		logrus.Infof(" - Compose file: %s", compose)
		logrus.Infof(" - Master image: %s", defaultMasterImageName)
		logrus.Infof(" - Worker image: %s", workerImageForStack)
	***REMOVED***()
	if err = deployStack(cli, defaultStackName, compose); err != nil ***REMOVED***
		return 1, err
	***REMOVED***
	logrus.Infof("The log will be displayed here after some duration."+
		"You can watch the live status via `docker service logs %s_worker`",
		defaultStackName)
	masterContainerID, err := waitForMasterUp(cli, defaultStackName)
	if err != nil ***REMOVED***
		return 1, err
	***REMOVED***
	rc, err := waitForContainerCompletion(cli, os.Stdout, os.Stderr, masterContainerID)
	if err != nil ***REMOVED***
		return 1, err
	***REMOVED***
	logrus.Infof("Exit status: %d", rc)
	return int(rc), nil
***REMOVED***

func ensureImages(cli *client.Client, images []string) error ***REMOVED***
	for _, image := range images ***REMOVED***
		_, _, err := cli.ImageInspectWithRaw(context.Background(), image)
		if err != nil ***REMOVED***
			return fmt.Errorf("could not find image %s, please run `make build-integration-cli-on-swarm`: %v",
				image, err)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func filtersBytes(optionalFiltersFile string) ([]byte, error) ***REMOVED***
	var b []byte
	if optionalFiltersFile == "" ***REMOVED***
		tests, err := enumerateTests(".")
		if err != nil ***REMOVED***
			return b, err
		***REMOVED***
		b = []byte(strings.Join(tests, "\n") + "\n")
	***REMOVED*** else ***REMOVED***
		var err error
		b, err = ioutil.ReadFile(optionalFiltersFile)
		if err != nil ***REMOVED***
			return b, err
		***REMOVED***
	***REMOVED***
	return b, nil
***REMOVED***

func waitForMasterUp(cli *client.Client, stackName string) (string, error) ***REMOVED***
	// FIXME(AkihiroSuda): it should retry until master is up, rather than pre-sleeping
	time.Sleep(10 * time.Second)

	fil := filters.NewArgs()
	fil.Add("label", "com.docker.stack.namespace="+stackName)
	// FIXME(AkihiroSuda): we should not rely on internal service naming convention
	fil.Add("label", "com.docker.swarm.service.name="+stackName+"_master")
	masters, err := cli.ContainerList(context.Background(), types.ContainerListOptions***REMOVED***
		All:     true,
		Filters: fil,
	***REMOVED***)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if len(masters) == 0 ***REMOVED***
		return "", fmt.Errorf("master not running in stack %s?", stackName)
	***REMOVED***
	return masters[0].ID, nil
***REMOVED***

func waitForContainerCompletion(cli *client.Client, stdout, stderr io.Writer, containerID string) (int64, error) ***REMOVED***
	stream, err := cli.ContainerLogs(context.Background(),
		containerID,
		types.ContainerLogsOptions***REMOVED***
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
		***REMOVED***)
	if err != nil ***REMOVED***
		return 1, err
	***REMOVED***
	stdcopy.StdCopy(stdout, stderr, stream)
	stream.Close()
	resultC, errC := cli.ContainerWait(context.Background(), containerID, "")
	select ***REMOVED***
	case err := <-errC:
		return 1, err
	case result := <-resultC:
		return result.StatusCode, nil
	***REMOVED***
***REMOVED***
