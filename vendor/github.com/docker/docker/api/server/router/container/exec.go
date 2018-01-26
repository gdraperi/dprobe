package container

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/versions"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func (s *containerRouter) getExecByID(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	eConfig, err := s.backend.ContainerExecInspect(vars["id"])
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, eConfig)
***REMOVED***

type execCommandError struct***REMOVED******REMOVED***

func (execCommandError) Error() string ***REMOVED***
	return "No exec command specified"
***REMOVED***

func (execCommandError) InvalidParameter() ***REMOVED******REMOVED***

func (s *containerRouter) postContainerExecCreate(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := httputils.CheckForJSON(r); err != nil ***REMOVED***
		return err
	***REMOVED***
	name := vars["name"]

	execConfig := &types.ExecConfig***REMOVED******REMOVED***
	if err := json.NewDecoder(r.Body).Decode(execConfig); err != nil ***REMOVED***
		return err
	***REMOVED***

	if len(execConfig.Cmd) == 0 ***REMOVED***
		return execCommandError***REMOVED******REMOVED***
	***REMOVED***

	// Register an instance of Exec in container.
	id, err := s.backend.ContainerExecCreate(name, execConfig)
	if err != nil ***REMOVED***
		logrus.Errorf("Error setting up exec command in container %s: %v", name, err)
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusCreated, &types.IDResponse***REMOVED***
		ID: id,
	***REMOVED***)
***REMOVED***

// TODO(vishh): Refactor the code to avoid having to specify stream config as part of both create and start.
func (s *containerRouter) postContainerExecStart(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	version := httputils.VersionFromContext(ctx)
	if versions.GreaterThan(version, "1.21") ***REMOVED***
		if err := httputils.CheckForJSON(r); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	var (
		execName                  = vars["name"]
		stdin, inStream           io.ReadCloser
		stdout, stderr, outStream io.Writer
	)

	execStartCheck := &types.ExecStartCheck***REMOVED******REMOVED***
	if err := json.NewDecoder(r.Body).Decode(execStartCheck); err != nil ***REMOVED***
		return err
	***REMOVED***

	if exists, err := s.backend.ExecExists(execName); !exists ***REMOVED***
		return err
	***REMOVED***

	if !execStartCheck.Detach ***REMOVED***
		var err error
		// Setting up the streaming http interface.
		inStream, outStream, err = httputils.HijackConnection(w)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		defer httputils.CloseStreams(inStream, outStream)

		if _, ok := r.Header["Upgrade"]; ok ***REMOVED***
			fmt.Fprint(outStream, "HTTP/1.1 101 UPGRADED\r\nContent-Type: application/vnd.docker.raw-stream\r\nConnection: Upgrade\r\nUpgrade: tcp\r\n")
		***REMOVED*** else ***REMOVED***
			fmt.Fprint(outStream, "HTTP/1.1 200 OK\r\nContent-Type: application/vnd.docker.raw-stream\r\n")
		***REMOVED***

		// copy headers that were removed as part of hijack
		if err := w.Header().WriteSubset(outStream, nil); err != nil ***REMOVED***
			return err
		***REMOVED***
		fmt.Fprint(outStream, "\r\n")

		stdin = inStream
		stdout = outStream
		if !execStartCheck.Tty ***REMOVED***
			stderr = stdcopy.NewStdWriter(outStream, stdcopy.Stderr)
			stdout = stdcopy.NewStdWriter(outStream, stdcopy.Stdout)
		***REMOVED***
	***REMOVED***

	// Now run the user process in container.
	// Maybe we should we pass ctx here if we're not detaching?
	if err := s.backend.ContainerExecStart(context.Background(), execName, stdin, stdout, stderr); err != nil ***REMOVED***
		if execStartCheck.Detach ***REMOVED***
			return err
		***REMOVED***
		stdout.Write([]byte(err.Error() + "\r\n"))
		logrus.Errorf("Error running exec %s in container: %v", execName, err)
	***REMOVED***
	return nil
***REMOVED***

func (s *containerRouter) postContainerExecResize(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***
	height, err := strconv.Atoi(r.Form.Get("h"))
	if err != nil ***REMOVED***
		return errdefs.InvalidParameter(err)
	***REMOVED***
	width, err := strconv.Atoi(r.Form.Get("w"))
	if err != nil ***REMOVED***
		return errdefs.InvalidParameter(err)
	***REMOVED***

	return s.backend.ContainerExecResize(vars["name"], height, width)
***REMOVED***
