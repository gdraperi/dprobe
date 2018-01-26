package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/docker/docker/client"
)

func system(commands [][]string) error ***REMOVED***
	for _, c := range commands ***REMOVED***
		cmd := exec.Command(c[0], c[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()
		if err := cmd.Run(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func pushImage(unusedCli *client.Client, remote, local string) error ***REMOVED***
	// FIXME: eliminate os/exec (but it is hard to pass auth without os/exec ...)
	return system([][]string***REMOVED***
		***REMOVED***"docker", "image", "tag", local, remote***REMOVED***,
		***REMOVED***"docker", "image", "push", remote***REMOVED***,
	***REMOVED***)
***REMOVED***

func deployStack(unusedCli *client.Client, stackName, composeFilePath string) error ***REMOVED***
	// FIXME: eliminate os/exec (but stack is implemented in CLI ...)
	return system([][]string***REMOVED***
		***REMOVED***"docker", "stack", "deploy",
			"--compose-file", composeFilePath,
			"--with-registry-auth",
			stackName***REMOVED***,
	***REMOVED***)
***REMOVED***

func hasStack(unusedCli *client.Client, stackName string) bool ***REMOVED***
	// FIXME: eliminate os/exec (but stack is implemented in CLI ...)
	out, err := exec.Command("docker", "stack", "ls").CombinedOutput()
	if err != nil ***REMOVED***
		panic(fmt.Errorf("`docker stack ls` failed with: %s", string(out)))
	***REMOVED***
	// FIXME: not accurate
	return strings.Contains(string(out), stackName)
***REMOVED***

func removeStack(unusedCli *client.Client, stackName string) error ***REMOVED***
	// FIXME: eliminate os/exec (but stack is implemented in CLI ...)
	if err := system([][]string***REMOVED***
		***REMOVED***"docker", "stack", "rm", stackName***REMOVED***,
	***REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***
	// FIXME
	time.Sleep(10 * time.Second)
	return nil
***REMOVED***
