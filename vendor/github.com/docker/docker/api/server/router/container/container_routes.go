package container

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"syscall"

	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/versions"
	containerpkg "github.com/docker/docker/container"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/signal"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

func (s *containerRouter) getContainersJSON(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***
	filter, err := filters.FromJSON(r.Form.Get("filters"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	config := &types.ContainerListOptions***REMOVED***
		All:     httputils.BoolValue(r, "all"),
		Size:    httputils.BoolValue(r, "size"),
		Since:   r.Form.Get("since"),
		Before:  r.Form.Get("before"),
		Filters: filter,
	***REMOVED***

	if tmpLimit := r.Form.Get("limit"); tmpLimit != "" ***REMOVED***
		limit, err := strconv.Atoi(tmpLimit)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		config.Limit = limit
	***REMOVED***

	containers, err := s.backend.Containers(config)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, containers)
***REMOVED***

func (s *containerRouter) getContainersStats(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	stream := httputils.BoolValueOrDefault(r, "stream", true)
	if !stream ***REMOVED***
		w.Header().Set("Content-Type", "application/json")
	***REMOVED***

	config := &backend.ContainerStatsConfig***REMOVED***
		Stream:    stream,
		OutStream: w,
		Version:   httputils.VersionFromContext(ctx),
	***REMOVED***

	return s.backend.ContainerStats(ctx, vars["name"], config)
***REMOVED***

func (s *containerRouter) getContainersLogs(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Args are validated before the stream starts because when it starts we're
	// sending HTTP 200 by writing an empty chunk of data to tell the client that
	// daemon is going to stream. By sending this initial HTTP 200 we can't report
	// any error after the stream starts (i.e. container not found, wrong parameters)
	// with the appropriate status code.
	stdout, stderr := httputils.BoolValue(r, "stdout"), httputils.BoolValue(r, "stderr")
	if !(stdout || stderr) ***REMOVED***
		return errdefs.InvalidParameter(errors.New("Bad parameters: you must choose at least one stream"))
	***REMOVED***

	containerName := vars["name"]
	logsConfig := &types.ContainerLogsOptions***REMOVED***
		Follow:     httputils.BoolValue(r, "follow"),
		Timestamps: httputils.BoolValue(r, "timestamps"),
		Since:      r.Form.Get("since"),
		Until:      r.Form.Get("until"),
		Tail:       r.Form.Get("tail"),
		ShowStdout: stdout,
		ShowStderr: stderr,
		Details:    httputils.BoolValue(r, "details"),
	***REMOVED***

	msgs, tty, err := s.backend.ContainerLogs(ctx, containerName, logsConfig)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// if has a tty, we're not muxing streams. if it doesn't, we are. simple.
	// this is the point of no return for writing a response. once we call
	// WriteLogStream, the response has been started and errors will be
	// returned in band by WriteLogStream
	httputils.WriteLogStream(ctx, w, msgs, logsConfig, !tty)
	return nil
***REMOVED***

func (s *containerRouter) getContainersExport(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	return s.backend.ContainerExport(vars["name"], w)
***REMOVED***

type bodyOnStartError struct***REMOVED******REMOVED***

func (bodyOnStartError) Error() string ***REMOVED***
	return "starting container with non-empty request body was deprecated since API v1.22 and removed in v1.24"
***REMOVED***

func (bodyOnStartError) InvalidParameter() ***REMOVED******REMOVED***

func (s *containerRouter) postContainersStart(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	// If contentLength is -1, we can assumed chunked encoding
	// or more technically that the length is unknown
	// https://golang.org/src/pkg/net/http/request.go#L139
	// net/http otherwise seems to swallow any headers related to chunked encoding
	// including r.TransferEncoding
	// allow a nil body for backwards compatibility

	version := httputils.VersionFromContext(ctx)
	var hostConfig *container.HostConfig
	// A non-nil json object is at least 7 characters.
	if r.ContentLength > 7 || r.ContentLength == -1 ***REMOVED***
		if versions.GreaterThanOrEqualTo(version, "1.24") ***REMOVED***
			return bodyOnStartError***REMOVED******REMOVED***
		***REMOVED***

		if err := httputils.CheckForJSON(r); err != nil ***REMOVED***
			return err
		***REMOVED***

		c, err := s.decoder.DecodeHostConfig(r.Body)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		hostConfig = c
	***REMOVED***

	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	checkpoint := r.Form.Get("checkpoint")
	checkpointDir := r.Form.Get("checkpoint-dir")
	if err := s.backend.ContainerStart(vars["name"], hostConfig, checkpoint, checkpointDir); err != nil ***REMOVED***
		return err
	***REMOVED***

	w.WriteHeader(http.StatusNoContent)
	return nil
***REMOVED***

func (s *containerRouter) postContainersStop(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	var seconds *int
	if tmpSeconds := r.Form.Get("t"); tmpSeconds != "" ***REMOVED***
		valSeconds, err := strconv.Atoi(tmpSeconds)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		seconds = &valSeconds
	***REMOVED***

	if err := s.backend.ContainerStop(vars["name"], seconds); err != nil ***REMOVED***
		return err
	***REMOVED***
	w.WriteHeader(http.StatusNoContent)

	return nil
***REMOVED***

func (s *containerRouter) postContainersKill(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	var sig syscall.Signal
	name := vars["name"]

	// If we have a signal, look at it. Otherwise, do nothing
	if sigStr := r.Form.Get("signal"); sigStr != "" ***REMOVED***
		var err error
		if sig, err = signal.ParseSignal(sigStr); err != nil ***REMOVED***
			return errdefs.InvalidParameter(err)
		***REMOVED***
	***REMOVED***

	if err := s.backend.ContainerKill(name, uint64(sig)); err != nil ***REMOVED***
		var isStopped bool
		if errdefs.IsConflict(err) ***REMOVED***
			isStopped = true
		***REMOVED***

		// Return error that's not caused because the container is stopped.
		// Return error if the container is not running and the api is >= 1.20
		// to keep backwards compatibility.
		version := httputils.VersionFromContext(ctx)
		if versions.GreaterThanOrEqualTo(version, "1.20") || !isStopped ***REMOVED***
			return errors.Wrapf(err, "Cannot kill container: %s", name)
		***REMOVED***
	***REMOVED***

	w.WriteHeader(http.StatusNoContent)
	return nil
***REMOVED***

func (s *containerRouter) postContainersRestart(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	var seconds *int
	if tmpSeconds := r.Form.Get("t"); tmpSeconds != "" ***REMOVED***
		valSeconds, err := strconv.Atoi(tmpSeconds)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		seconds = &valSeconds
	***REMOVED***

	if err := s.backend.ContainerRestart(vars["name"], seconds); err != nil ***REMOVED***
		return err
	***REMOVED***

	w.WriteHeader(http.StatusNoContent)

	return nil
***REMOVED***

func (s *containerRouter) postContainersPause(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := s.backend.ContainerPause(vars["name"]); err != nil ***REMOVED***
		return err
	***REMOVED***

	w.WriteHeader(http.StatusNoContent)

	return nil
***REMOVED***

func (s *containerRouter) postContainersUnpause(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := s.backend.ContainerUnpause(vars["name"]); err != nil ***REMOVED***
		return err
	***REMOVED***

	w.WriteHeader(http.StatusNoContent)

	return nil
***REMOVED***

func (s *containerRouter) postContainersWait(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	// Behavior changed in version 1.30 to handle wait condition and to
	// return headers immediately.
	version := httputils.VersionFromContext(ctx)
	legacyBehaviorPre130 := versions.LessThan(version, "1.30")
	legacyRemovalWaitPre134 := false

	// The wait condition defaults to "not-running".
	waitCondition := containerpkg.WaitConditionNotRunning
	if !legacyBehaviorPre130 ***REMOVED***
		if err := httputils.ParseForm(r); err != nil ***REMOVED***
			return err
		***REMOVED***
		switch container.WaitCondition(r.Form.Get("condition")) ***REMOVED***
		case container.WaitConditionNextExit:
			waitCondition = containerpkg.WaitConditionNextExit
		case container.WaitConditionRemoved:
			waitCondition = containerpkg.WaitConditionRemoved
			legacyRemovalWaitPre134 = versions.LessThan(version, "1.34")
		***REMOVED***
	***REMOVED***

	// Note: the context should get canceled if the client closes the
	// connection since this handler has been wrapped by the
	// router.WithCancel() wrapper.
	waitC, err := s.backend.ContainerWait(ctx, vars["name"], waitCondition)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	w.Header().Set("Content-Type", "application/json")

	if !legacyBehaviorPre130 ***REMOVED***
		// Write response header immediately.
		w.WriteHeader(http.StatusOK)
		if flusher, ok := w.(http.Flusher); ok ***REMOVED***
			flusher.Flush()
		***REMOVED***
	***REMOVED***

	// Block on the result of the wait operation.
	status := <-waitC

	// With API < 1.34, wait on WaitConditionRemoved did not return
	// in case container removal failed. The only way to report an
	// error back to the client is to not write anything (i.e. send
	// an empty response which will be treated as an error).
	if legacyRemovalWaitPre134 && status.Err() != nil ***REMOVED***
		return nil
	***REMOVED***

	var waitError *container.ContainerWaitOKBodyError
	if status.Err() != nil ***REMOVED***
		waitError = &container.ContainerWaitOKBodyError***REMOVED***Message: status.Err().Error()***REMOVED***
	***REMOVED***

	return json.NewEncoder(w).Encode(&container.ContainerWaitOKBody***REMOVED***
		StatusCode: int64(status.ExitCode()),
		Error:      waitError,
	***REMOVED***)
***REMOVED***

func (s *containerRouter) getContainersChanges(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	changes, err := s.backend.ContainerChanges(vars["name"])
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, changes)
***REMOVED***

func (s *containerRouter) getContainersTop(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	procList, err := s.backend.ContainerTop(vars["name"], r.Form.Get("ps_args"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, procList)
***REMOVED***

func (s *containerRouter) postContainerRename(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	name := vars["name"]
	newName := r.Form.Get("name")
	if err := s.backend.ContainerRename(name, newName); err != nil ***REMOVED***
		return err
	***REMOVED***
	w.WriteHeader(http.StatusNoContent)
	return nil
***REMOVED***

func (s *containerRouter) postContainerUpdate(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := httputils.CheckForJSON(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	var updateConfig container.UpdateConfig

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updateConfig); err != nil ***REMOVED***
		return err
	***REMOVED***

	hostConfig := &container.HostConfig***REMOVED***
		Resources:     updateConfig.Resources,
		RestartPolicy: updateConfig.RestartPolicy,
	***REMOVED***

	name := vars["name"]
	resp, err := s.backend.ContainerUpdate(name, hostConfig)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, resp)
***REMOVED***

func (s *containerRouter) postContainersCreate(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := httputils.CheckForJSON(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	name := r.Form.Get("name")

	config, hostConfig, networkingConfig, err := s.decoder.DecodeConfig(r.Body)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	version := httputils.VersionFromContext(ctx)
	adjustCPUShares := versions.LessThan(version, "1.19")

	// When using API 1.24 and under, the client is responsible for removing the container
	if hostConfig != nil && versions.LessThan(version, "1.25") ***REMOVED***
		hostConfig.AutoRemove = false
	***REMOVED***

	ccr, err := s.backend.ContainerCreate(types.ContainerCreateConfig***REMOVED***
		Name:             name,
		Config:           config,
		HostConfig:       hostConfig,
		NetworkingConfig: networkingConfig,
		AdjustCPUShares:  adjustCPUShares,
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusCreated, ccr)
***REMOVED***

func (s *containerRouter) deleteContainers(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	name := vars["name"]
	config := &types.ContainerRmConfig***REMOVED***
		ForceRemove:  httputils.BoolValue(r, "force"),
		RemoveVolume: httputils.BoolValue(r, "v"),
		RemoveLink:   httputils.BoolValue(r, "link"),
	***REMOVED***

	if err := s.backend.ContainerRm(name, config); err != nil ***REMOVED***
		return err
	***REMOVED***

	w.WriteHeader(http.StatusNoContent)

	return nil
***REMOVED***

func (s *containerRouter) postContainersResize(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
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

	return s.backend.ContainerResize(vars["name"], height, width)
***REMOVED***

func (s *containerRouter) postContainersAttach(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	err := httputils.ParseForm(r)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	containerName := vars["name"]

	_, upgrade := r.Header["Upgrade"]
	detachKeys := r.FormValue("detachKeys")

	hijacker, ok := w.(http.Hijacker)
	if !ok ***REMOVED***
		return errdefs.InvalidParameter(errors.Errorf("error attaching to container %s, hijack connection missing", containerName))
	***REMOVED***

	setupStreams := func() (io.ReadCloser, io.Writer, io.Writer, error) ***REMOVED***
		conn, _, err := hijacker.Hijack()
		if err != nil ***REMOVED***
			return nil, nil, nil, err
		***REMOVED***

		// set raw mode
		conn.Write([]byte***REMOVED******REMOVED***)

		if upgrade ***REMOVED***
			fmt.Fprintf(conn, "HTTP/1.1 101 UPGRADED\r\nContent-Type: application/vnd.docker.raw-stream\r\nConnection: Upgrade\r\nUpgrade: tcp\r\n\r\n")
		***REMOVED*** else ***REMOVED***
			fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: application/vnd.docker.raw-stream\r\n\r\n")
		***REMOVED***

		closer := func() error ***REMOVED***
			httputils.CloseStreams(conn)
			return nil
		***REMOVED***
		return ioutils.NewReadCloserWrapper(conn, closer), conn, conn, nil
	***REMOVED***

	attachConfig := &backend.ContainerAttachConfig***REMOVED***
		GetStreams: setupStreams,
		UseStdin:   httputils.BoolValue(r, "stdin"),
		UseStdout:  httputils.BoolValue(r, "stdout"),
		UseStderr:  httputils.BoolValue(r, "stderr"),
		Logs:       httputils.BoolValue(r, "logs"),
		Stream:     httputils.BoolValue(r, "stream"),
		DetachKeys: detachKeys,
		MuxStreams: true,
	***REMOVED***

	if err = s.backend.ContainerAttach(containerName, attachConfig); err != nil ***REMOVED***
		logrus.Errorf("Handler for %s %s returned error: %v", r.Method, r.URL.Path, err)
		// Remember to close stream if error happens
		conn, _, errHijack := hijacker.Hijack()
		if errHijack == nil ***REMOVED***
			statusCode := httputils.GetHTTPErrorStatusCode(err)
			statusText := http.StatusText(statusCode)
			fmt.Fprintf(conn, "HTTP/1.1 %d %s\r\nContent-Type: application/vnd.docker.raw-stream\r\n\r\n%s\r\n", statusCode, statusText, err.Error())
			httputils.CloseStreams(conn)
		***REMOVED*** else ***REMOVED***
			logrus.Errorf("Error Hijacking: %v", err)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (s *containerRouter) wsContainersAttach(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***
	containerName := vars["name"]

	var err error
	detachKeys := r.FormValue("detachKeys")

	done := make(chan struct***REMOVED******REMOVED***)
	started := make(chan struct***REMOVED******REMOVED***)

	version := httputils.VersionFromContext(ctx)

	setupStreams := func() (io.ReadCloser, io.Writer, io.Writer, error) ***REMOVED***
		wsChan := make(chan *websocket.Conn)
		h := func(conn *websocket.Conn) ***REMOVED***
			wsChan <- conn
			<-done
		***REMOVED***

		srv := websocket.Server***REMOVED***Handler: h, Handshake: nil***REMOVED***
		go func() ***REMOVED***
			close(started)
			srv.ServeHTTP(w, r)
		***REMOVED***()

		conn := <-wsChan
		// In case version 1.28 and above, a binary frame will be sent.
		// See 28176 for details.
		if versions.GreaterThanOrEqualTo(version, "1.28") ***REMOVED***
			conn.PayloadType = websocket.BinaryFrame
		***REMOVED***
		return conn, conn, conn, nil
	***REMOVED***

	attachConfig := &backend.ContainerAttachConfig***REMOVED***
		GetStreams: setupStreams,
		Logs:       httputils.BoolValue(r, "logs"),
		Stream:     httputils.BoolValue(r, "stream"),
		DetachKeys: detachKeys,
		UseStdin:   true,
		UseStdout:  true,
		UseStderr:  true,
		MuxStreams: false, // TODO: this should be true since it's a single stream for both stdout and stderr
	***REMOVED***

	err = s.backend.ContainerAttach(containerName, attachConfig)
	close(done)
	select ***REMOVED***
	case <-started:
		if err != nil ***REMOVED***
			logrus.Errorf("Error attaching websocket: %s", err)
		***REMOVED*** else ***REMOVED***
			logrus.Debug("websocket connection was closed by client")
		***REMOVED***
		return nil
	default:
	***REMOVED***
	return err
***REMOVED***

func (s *containerRouter) postContainersPrune(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	pruneFilters, err := filters.FromJSON(r.Form.Get("filters"))
	if err != nil ***REMOVED***
		return errdefs.InvalidParameter(err)
	***REMOVED***

	pruneReport, err := s.backend.ContainersPrune(ctx, pruneFilters)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return httputils.WriteJSON(w, http.StatusOK, pruneReport)
***REMOVED***
