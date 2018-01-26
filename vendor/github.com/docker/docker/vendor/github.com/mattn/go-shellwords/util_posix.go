// +build !windows

package shellwords

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

func shellRun(line string) (string, error) ***REMOVED***
	shell := os.Getenv("SHELL")
	b, err := exec.Command(shell, "-c", line).Output()
	if err != nil ***REMOVED***
		return "", errors.New(err.Error() + ":" + string(b))
	***REMOVED***
	return strings.TrimSpace(string(b)), nil
***REMOVED***
