package pty

import (
	"os"
	"os/exec"
	"syscall"
)

// Start assigns a pseudo-terminal tty os.File to c.Stdin, c.Stdout,
// and c.Stderr, calls c.Start, and returns the File of the tty's
// corresponding pty.
func Start(c *exec.Cmd) (pty *os.File, err error) ***REMOVED***
	pty, tty, err := Open()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer tty.Close()
	c.Stdout = tty
	c.Stdin = tty
	c.Stderr = tty
	c.SysProcAttr = &syscall.SysProcAttr***REMOVED***Setctty: true, Setsid: true***REMOVED***
	err = c.Start()
	if err != nil ***REMOVED***
		pty.Close()
		return nil, err
	***REMOVED***
	return pty, err
***REMOVED***
