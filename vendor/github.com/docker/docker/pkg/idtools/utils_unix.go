// +build !windows

package idtools

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func resolveBinary(binname string) (string, error) ***REMOVED***
	binaryPath, err := exec.LookPath(binname)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	resolvedPath, err := filepath.EvalSymlinks(binaryPath)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	//only return no error if the final resolved binary basename
	//matches what was searched for
	if filepath.Base(resolvedPath) == binname ***REMOVED***
		return resolvedPath, nil
	***REMOVED***
	return "", fmt.Errorf("Binary %q does not resolve to a binary of that name in $PATH (%q)", binname, resolvedPath)
***REMOVED***

func execCmd(cmd, args string) ([]byte, error) ***REMOVED***
	execCmd := exec.Command(cmd, strings.Split(args, " ")...)
	return execCmd.CombinedOutput()
***REMOVED***
