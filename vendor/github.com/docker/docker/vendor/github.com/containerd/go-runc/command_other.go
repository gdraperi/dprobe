// +build !linux

package runc

import (
	"context"
	"os/exec"
)

func (r *Runc) command(context context.Context, args ...string) *exec.Cmd ***REMOVED***
	command := r.Command
	if command == "" ***REMOVED***
		command = DefaultCommand
	***REMOVED***
	return exec.CommandContext(context, command, append(r.args(), args...)...)
***REMOVED***
