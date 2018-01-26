package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
)

func TestEventsErrorInOptions(t *testing.T) ***REMOVED***
	errorCases := []struct ***REMOVED***
		options       types.EventsOptions
		expectedError string
	***REMOVED******REMOVED***
		***REMOVED***
			options: types.EventsOptions***REMOVED***
				Since: "2006-01-02TZ",
			***REMOVED***,
			expectedError: `parsing time "2006-01-02TZ"`,
		***REMOVED***,
		***REMOVED***
			options: types.EventsOptions***REMOVED***
				Until: "2006-01-02TZ",
			***REMOVED***,
			expectedError: `parsing time "2006-01-02TZ"`,
		***REMOVED***,
	***REMOVED***
	for _, e := range errorCases ***REMOVED***
		client := &Client***REMOVED***
			client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
		***REMOVED***
		_, errs := client.Events(context.Background(), e.options)
		err := <-errs
		if err == nil || !strings.Contains(err.Error(), e.expectedError) ***REMOVED***
			t.Fatalf("expected an error %q, got %v", e.expectedError, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestEventsErrorFromServer(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, errs := client.Events(context.Background(), types.EventsOptions***REMOVED******REMOVED***)
	err := <-errs
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestEvents(t *testing.T) ***REMOVED***

	expectedURL := "/events"

	filters := filters.NewArgs()
	filters.Add("type", events.ContainerEventType)
	expectedFiltersJSON := fmt.Sprintf(`***REMOVED***"type":***REMOVED***"%s":true***REMOVED******REMOVED***`, events.ContainerEventType)

	eventsCases := []struct ***REMOVED***
		options             types.EventsOptions
		events              []events.Message
		expectedEvents      map[string]bool
		expectedQueryParams map[string]string
	***REMOVED******REMOVED***
		***REMOVED***
			options: types.EventsOptions***REMOVED***
				Filters: filters,
			***REMOVED***,
			expectedQueryParams: map[string]string***REMOVED***
				"filters": expectedFiltersJSON,
			***REMOVED***,
			events:         []events.Message***REMOVED******REMOVED***,
			expectedEvents: make(map[string]bool),
		***REMOVED***,
		***REMOVED***
			options: types.EventsOptions***REMOVED***
				Filters: filters,
			***REMOVED***,
			expectedQueryParams: map[string]string***REMOVED***
				"filters": expectedFiltersJSON,
			***REMOVED***,
			events: []events.Message***REMOVED***
				***REMOVED***
					Type:   "container",
					ID:     "1",
					Action: "create",
				***REMOVED***,
				***REMOVED***
					Type:   "container",
					ID:     "2",
					Action: "die",
				***REMOVED***,
				***REMOVED***
					Type:   "container",
					ID:     "3",
					Action: "create",
				***REMOVED***,
			***REMOVED***,
			expectedEvents: map[string]bool***REMOVED***
				"1": true,
				"2": true,
				"3": true,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for _, eventsCase := range eventsCases ***REMOVED***
		client := &Client***REMOVED***
			client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
				if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
					return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
				***REMOVED***
				query := req.URL.Query()

				for key, expected := range eventsCase.expectedQueryParams ***REMOVED***
					actual := query.Get(key)
					if actual != expected ***REMOVED***
						return nil, fmt.Errorf("%s not set in URL query properly. Expected '%s', got %s", key, expected, actual)
					***REMOVED***
				***REMOVED***

				buffer := new(bytes.Buffer)

				for _, e := range eventsCase.events ***REMOVED***
					b, _ := json.Marshal(e)
					buffer.Write(b)
				***REMOVED***

				return &http.Response***REMOVED***
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(buffer),
				***REMOVED***, nil
			***REMOVED***),
		***REMOVED***

		messages, errs := client.Events(context.Background(), eventsCase.options)

	loop:
		for ***REMOVED***
			select ***REMOVED***
			case err := <-errs:
				if err != nil && err != io.EOF ***REMOVED***
					t.Fatal(err)
				***REMOVED***

				break loop
			case e := <-messages:
				_, ok := eventsCase.expectedEvents[e.ID]
				if !ok ***REMOVED***
					t.Fatalf("event received not expected with action %s & id %s", e.Action, e.ID)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
