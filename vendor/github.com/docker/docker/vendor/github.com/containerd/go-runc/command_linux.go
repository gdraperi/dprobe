package runc

import (
	"context"
	"os/exec"
	"syscall"
)

func (r *Runc) command(context context.Context, args ...string) *exec.Cmd ***REMOVED***
	command := r.Command
	if command == "" ***REMOVED***
		command = DefaultCommand
	***REMOVED***
	cmd := exec.CommandContext(context, command, append(r.args(), args...)...)
	cmd.SysProcAttr = &syscall.SysProcAttr***REMOVED***
		Setpgid: r.Setpgid,
	***REMOVED***
	if r.PdeathSignal != 0 ***REMOVED***
		cmd.SysProcAttr.Pdeathsig = r.PdeathSignal
	***REMOVED***

	return cmd
***REMOVED***
