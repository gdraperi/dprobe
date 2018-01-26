package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

func TestPluginInspectError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***

	_, _, err := client.PluginInspectWithRaw(context.Background(), "nothing")
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestPluginInspect(t *testing.T) ***REMOVED***
	expectedURL := "/plugins/plugin_name"
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			content, err := json.Marshal(types.Plugin***REMOVED***
				ID: "plugin_id",
			***REMOVED***)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return &http.Response***REMOVED***
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(content)),
			***REMOVED***, nil
		***REMOVED***),
	***REMOVED***

	pluginInspect, _, err := client.PluginInspectWithRaw(context.Background(), "plugin_name")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if pluginInspect.ID != "plugin_id" ***REMOVED***
		t.Fatalf("expected `plugin_id`, got %s", pluginInspect.ID)
	***REMOVED***
***REMOVED***
