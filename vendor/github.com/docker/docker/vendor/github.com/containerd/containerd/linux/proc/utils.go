// +build !windows

package proc

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"time"

	"github.com/containerd/containerd/errdefs"
	runc "github.com/containerd/go-runc"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

// TODO(mlaventure): move to runc package?
func getLastRuntimeError(r *runc.Runc) (string, error) ***REMOVED***
	if r.Log == "" ***REMOVED***
		return "", nil
	***REMOVED***

	f, err := os.OpenFile(r.Log, os.O_RDONLY, 0400)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	var (
		errMsg string
		log    struct ***REMOVED***
			Level string
			Msg   string
			Time  time.Time
		***REMOVED***
	)

	dec := json.NewDecoder(f)
	for err = nil; err == nil; ***REMOVED***
		if err = dec.Decode(&log); err != nil && err != io.EOF ***REMOVED***
			return "", err
		***REMOVED***
		if log.Level == "error" ***REMOVED***
			errMsg = strings.TrimSpace(log.Msg)
		***REMOVED***
	***REMOVED***

	return errMsg, nil
***REMOVED***

// criuError returns only the first line of the error message from criu
// it tries to add an invalid dump log location when returning the message
func criuError(err error) string ***REMOVED***
	parts := strings.Split(err.Error(), "\n")
	return parts[0]
***REMOVED***

func copyFile(to, from string) error ***REMOVED***
	ff, err := os.Open(from)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer ff.Close()
	tt, err := os.Create(to)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer tt.Close()
	_, err = io.Copy(tt, ff)
	return err
***REMOVED***

func checkKillError(err error) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***
	if strings.Contains(err.Error(), "os: process already finished") || err == unix.ESRCH ***REMOVED***
		return errors.Wrapf(errdefs.ErrNotFound, "process already finished")
	***REMOVED***
	return errors.Wrapf(err, "unknown error after kill")
***REMOVED***

func hasNoIO(r *CreateConfig) bool ***REMOVED***
	return r.Stdin == "" && r.Stdout == "" && r.Stderr == ""
***REMOVED***
