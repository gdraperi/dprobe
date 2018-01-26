package client

import (
	"encoding/json"
	"net/url"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	timetypes "github.com/docker/docker/api/types/time"
)

// Events returns a stream of events in the daemon. It's up to the caller to close the stream
// by cancelling the context. Once the stream has been completely read an io.EOF error will
// be sent over the error channel. If an error is sent all processing will be stopped. It's up
// to the caller to reopen the stream in the event of an error by reinvoking this method.
func (cli *Client) Events(ctx context.Context, options types.EventsOptions) (<-chan events.Message, <-chan error) ***REMOVED***

	messages := make(chan events.Message)
	errs := make(chan error, 1)

	started := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		defer close(errs)

		query, err := buildEventsQueryParams(cli.version, options)
		if err != nil ***REMOVED***
			close(started)
			errs <- err
			return
		***REMOVED***

		resp, err := cli.get(ctx, "/events", query, nil)
		if err != nil ***REMOVED***
			close(started)
			errs <- err
			return
		***REMOVED***
		defer resp.body.Close()

		decoder := json.NewDecoder(resp.body)

		close(started)
		for ***REMOVED***
			select ***REMOVED***
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			default:
				var event events.Message
				if err := decoder.Decode(&event); err != nil ***REMOVED***
					errs <- err
					return
				***REMOVED***

				select ***REMOVED***
				case messages <- event:
				case <-ctx.Done():
					errs <- ctx.Err()
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	<-started

	return messages, errs
***REMOVED***

func buildEventsQueryParams(cliVersion string, options types.EventsOptions) (url.Values, error) ***REMOVED***
	query := url.Values***REMOVED******REMOVED***
	ref := time.Now()

	if options.Since != "" ***REMOVED***
		ts, err := timetypes.GetTimestamp(options.Since, ref)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		query.Set("since", ts)
	***REMOVED***

	if options.Until != "" ***REMOVED***
		ts, err := timetypes.GetTimestamp(options.Until, ref)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		query.Set("until", ts)
	***REMOVED***

	if options.Filters.Len() > 0 ***REMOVED***
		filterJSON, err := filters.ToParamWithVersion(cliVersion, options.Filters)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		query.Set("filters", filterJSON)
	***REMOVED***

	return query, nil
***REMOVED***
