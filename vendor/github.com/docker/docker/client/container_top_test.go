package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/container"
	"golang.org/x/net/context"
)

func TestContainerTopError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, err := client.ContainerTop(context.Background(), "nothing", []string***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestContainerTop(t *testing.T) ***REMOVED***
	expectedURL := "/containers/container_id/top"
	expectedProcesses := [][]string***REMOVED***
		***REMOVED***"p1", "p2"***REMOVED***,
		***REMOVED***"p3"***REMOVED***,
	***REMOVED***
	expectedTitles := []string***REMOVED***"title1", "title2"***REMOVED***

	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			query := req.URL.Query()
			args := query.Get("ps_args")
			if args != "arg1 arg2" ***REMOVED***
				return nil, fmt.Errorf("args not set in URL query properly. Expected 'arg1 arg2', got %v", args)
			***REMOVED***

			b, err := json.Marshal(container.ContainerTopOKBody***REMOVED***
				Processes: [][]string***REMOVED***
					***REMOVED***"p1", "p2"***REMOVED***,
					***REMOVED***"p3"***REMOVED***,
				***REMOVED***,
				Titles: []string***REMOVED***"title1", "title2"***REMOVED***,
			***REMOVED***)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(b)),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	processList, err := client.ContainerTop(context.Background(), "container_id", []string***REMOVED***"arg1", "arg2"***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !reflect.DeepEqual(expectedProcesses, processList.Processes) ***REMOVED***
		t.Fatalf("Processes: expected %v, got %v", expectedProcesses, processList.Processes)
	***REMOVED***
	if !reflect.DeepEqual(expectedTitles, processList.Titles) ***REMOVED***
		t.Fatalf("Titles: expected %v, got %v", expectedTitles, processList.Titles)
	***REMOVED***
***REMOVED***
