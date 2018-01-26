//+build !windows

package daemon

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/container"
)

func validatePSArgs(psArgs string) error ***REMOVED***
	// NOTE: \\s does not detect unicode whitespaces.
	// So we use fieldsASCII instead of strings.Fields in parsePSOutput.
	// See https://github.com/docker/docker/pull/24358
	// nolint: gosimple
	re := regexp.MustCompile("\\s+([^\\s]*)=\\s*(PID[^\\s]*)")
	for _, group := range re.FindAllStringSubmatch(psArgs, -1) ***REMOVED***
		if len(group) >= 3 ***REMOVED***
			k := group[1]
			v := group[2]
			if k != "pid" ***REMOVED***
				return fmt.Errorf("specifying \"%s=%s\" is not allowed", k, v)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// fieldsASCII is similar to strings.Fields but only allows ASCII whitespaces
func fieldsASCII(s string) []string ***REMOVED***
	fn := func(r rune) bool ***REMOVED***
		switch r ***REMOVED***
		case '\t', '\n', '\f', '\r', ' ':
			return true
		***REMOVED***
		return false
	***REMOVED***
	return strings.FieldsFunc(s, fn)
***REMOVED***

func appendProcess2ProcList(procList *container.ContainerTopOKBody, fields []string) ***REMOVED***
	// Make sure number of fields equals number of header titles
	// merging "overhanging" fields
	process := fields[:len(procList.Titles)-1]
	process = append(process, strings.Join(fields[len(procList.Titles)-1:], " "))
	procList.Processes = append(procList.Processes, process)
***REMOVED***

func hasPid(procs []uint32, pid int) bool ***REMOVED***
	for _, p := range procs ***REMOVED***
		if int(p) == pid ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func parsePSOutput(output []byte, procs []uint32) (*container.ContainerTopOKBody, error) ***REMOVED***
	procList := &container.ContainerTopOKBody***REMOVED******REMOVED***

	lines := strings.Split(string(output), "\n")
	procList.Titles = fieldsASCII(lines[0])

	pidIndex := -1
	for i, name := range procList.Titles ***REMOVED***
		if name == "PID" ***REMOVED***
			pidIndex = i
		***REMOVED***
	***REMOVED***
	if pidIndex == -1 ***REMOVED***
		return nil, fmt.Errorf("Couldn't find PID field in ps output")
	***REMOVED***

	// loop through the output and extract the PID from each line
	// fixing #30580, be able to display thread line also when "m" option used
	// in "docker top" client command
	preContainedPidFlag := false
	for _, line := range lines[1:] ***REMOVED***
		if len(line) == 0 ***REMOVED***
			continue
		***REMOVED***
		fields := fieldsASCII(line)

		var (
			p   int
			err error
		)

		if fields[pidIndex] == "-" ***REMOVED***
			if preContainedPidFlag ***REMOVED***
				appendProcess2ProcList(procList, fields)
			***REMOVED***
			continue
		***REMOVED***
		p, err = strconv.Atoi(fields[pidIndex])
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("Unexpected pid '%s': %s", fields[pidIndex], err)
		***REMOVED***

		if hasPid(procs, p) ***REMOVED***
			preContainedPidFlag = true
			appendProcess2ProcList(procList, fields)
			continue
		***REMOVED***
		preContainedPidFlag = false
	***REMOVED***
	return procList, nil
***REMOVED***

// ContainerTop lists the processes running inside of the given
// container by calling ps with the given args, or with the flags
// "-ef" if no args are given.  An error is returned if the container
// is not found, or is not running, or if there are any problems
// running ps, or parsing the output.
func (daemon *Daemon) ContainerTop(name string, psArgs string) (*container.ContainerTopOKBody, error) ***REMOVED***
	if psArgs == "" ***REMOVED***
		psArgs = "-ef"
	***REMOVED***

	if err := validatePSArgs(psArgs); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	container, err := daemon.GetContainer(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if !container.IsRunning() ***REMOVED***
		return nil, errNotRunning(container.ID)
	***REMOVED***

	if container.IsRestarting() ***REMOVED***
		return nil, errContainerIsRestarting(container.ID)
	***REMOVED***

	procs, err := daemon.containerd.ListPids(context.Background(), container.ID)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	output, err := exec.Command("ps", strings.Split(psArgs, " ")...).Output()
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("Error running ps: %v", err)
	***REMOVED***
	procList, err := parsePSOutput(output, procs)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	daemon.LogContainerEvent(container, "top")
	return procList, nil
***REMOVED***
