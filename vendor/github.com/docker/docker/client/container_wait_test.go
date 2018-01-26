package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"

	"golang.org/x/net/context"
)

func TestContainerWaitError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	resultC, errC := client.ContainerWait(context.Background(), "nothing", "")
	select ***REMOVED***
	case result := <-resultC:
		t.Fatalf("expected to not get a wait result, got %d", result.StatusCode)
	case err := <-errC:
		if err.Error() != "Error response from daemon: Server error" ***REMOVED***
			t.Fatalf("expected a Server Error, got %v", err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestContainerWait(t *testing.T) ***REMOVED***
	expectedURL := "/containers/container_id/wait"
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			b, err := json.Marshal(container.ContainerWaitOKBody***REMOVED***
				StatusCode: 15,
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

	resultC, errC := client.ContainerWait(context.Background(), "container_id", "")
	select ***REMOVED***
	case err := <-errC:
		t.Fatal(err)
	case result := <-resultC:
		if result.StatusCode != 15 ***REMOVED***
			t.Fatalf("expected a status code equal to '15', got %d", result.StatusCode)
		***REMOVED***
	***REMOVED***
***REMOVED***

func ExampleClient_ContainerWait_withTimeout() ***REMOVED***
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, _ := NewEnvClient()
	_, errC := client.ContainerWait(ctx, "container_id", "")
	if err := <-errC; err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
***REMOVED***
