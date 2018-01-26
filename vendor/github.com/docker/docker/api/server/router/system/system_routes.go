package system

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/registry"
	timetypes "github.com/docker/docker/api/types/time"
	"github.com/docker/docker/api/types/versions"
	"github.com/docker/docker/pkg/ioutils"
	pkgerrors "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func optionsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	w.WriteHeader(http.StatusOK)
	return nil
***REMOVED***

func pingHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	_, err := w.Write([]byte***REMOVED***'O', 'K'***REMOVED***)
	return err
***REMOVED***

func (s *systemRouter) getInfo(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	info, err := s.backend.SystemInfo()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if s.cluster != nil ***REMOVED***
		info.Swarm = s.cluster.Info()
	***REMOVED***

	if versions.LessThan(httputils.VersionFromContext(ctx), "1.25") ***REMOVED***
		// TODO: handle this conversion in engine-api
		type oldInfo struct ***REMOVED***
			*types.Info
			ExecutionDriver string
		***REMOVED***
		old := &oldInfo***REMOVED***
			Info:            info,
			ExecutionDriver: "<not supported>",
		***REMOVED***
		nameOnlySecurityOptions := []string***REMOVED******REMOVED***
		kvSecOpts, err := types.DecodeSecurityOptions(old.SecurityOptions)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		for _, s := range kvSecOpts ***REMOVED***
			nameOnlySecurityOptions = append(nameOnlySecurityOptions, s.Name)
		***REMOVED***
		old.SecurityOptions = nameOnlySecurityOptions
		return httputils.WriteJSON(w, http.StatusOK, old)
	***REMOVED***
	return httputils.WriteJSON(w, http.StatusOK, info)
***REMOVED***

func (s *systemRouter) getVersion(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	info := s.backend.SystemVersion()

	return httputils.WriteJSON(w, http.StatusOK, info)
***REMOVED***

func (s *systemRouter) getDiskUsage(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	du, err := s.backend.SystemDiskUsage(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	builderSize, err := s.builder.DiskUsage()
	if err != nil ***REMOVED***
		return pkgerrors.Wrap(err, "error getting build cache usage")
	***REMOVED***
	du.BuilderSize = builderSize

	return httputils.WriteJSON(w, http.StatusOK, du)
***REMOVED***

type invalidRequestError struct ***REMOVED***
	Err error
***REMOVED***

func (e invalidRequestError) Error() string ***REMOVED***
	return e.Err.Error()
***REMOVED***

func (e invalidRequestError) InvalidParameter() ***REMOVED******REMOVED***

func (s *systemRouter) getEvents(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	since, err := eventTime(r.Form.Get("since"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	until, err := eventTime(r.Form.Get("until"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var (
		timeout        <-chan time.Time
		onlyPastEvents bool
	)
	if !until.IsZero() ***REMOVED***
		if until.Before(since) ***REMOVED***
			return invalidRequestError***REMOVED***fmt.Errorf("`since` time (%s) cannot be after `until` time (%s)", r.Form.Get("since"), r.Form.Get("until"))***REMOVED***
		***REMOVED***

		now := time.Now()

		onlyPastEvents = until.Before(now)

		if !onlyPastEvents ***REMOVED***
			dur := until.Sub(now)
			timeout = time.After(dur)
		***REMOVED***
	***REMOVED***

	ef, err := filters.FromJSON(r.Form.Get("filters"))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	w.Header().Set("Content-Type", "application/json")
	output := ioutils.NewWriteFlusher(w)
	defer output.Close()
	output.Flush()

	enc := json.NewEncoder(output)

	buffered, l := s.backend.SubscribeToEvents(since, until, ef)
	defer s.backend.UnsubscribeFromEvents(l)

	for _, ev := range buffered ***REMOVED***
		if err := enc.Encode(ev); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if onlyPastEvents ***REMOVED***
		return nil
	***REMOVED***

	for ***REMOVED***
		select ***REMOVED***
		case ev := <-l:
			jev, ok := ev.(events.Message)
			if !ok ***REMOVED***
				logrus.Warnf("unexpected event message: %q", ev)
				continue
			***REMOVED***
			if err := enc.Encode(jev); err != nil ***REMOVED***
				return err
			***REMOVED***
		case <-timeout:
			return nil
		case <-ctx.Done():
			logrus.Debug("Client context cancelled, stop sending events")
			return nil
		***REMOVED***
	***REMOVED***
***REMOVED***

func (s *systemRouter) postAuth(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	var config *types.AuthConfig
	err := json.NewDecoder(r.Body).Decode(&config)
	r.Body.Close()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	status, token, err := s.backend.AuthenticateToRegistry(ctx, config)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return httputils.WriteJSON(w, http.StatusOK, &registry.AuthenticateOKBody***REMOVED***
		Status:        status,
		IdentityToken: token,
	***REMOVED***)
***REMOVED***

func eventTime(formTime string) (time.Time, error) ***REMOVED***
	t, tNano, err := timetypes.ParseTimestamps(formTime, -1)
	if err != nil ***REMOVED***
		return time.Time***REMOVED******REMOVED***, err
	***REMOVED***
	if t == -1 ***REMOVED***
		return time.Time***REMOVED******REMOVED***, nil
	***REMOVED***
	return time.Unix(t, tNano), nil
***REMOVED***
