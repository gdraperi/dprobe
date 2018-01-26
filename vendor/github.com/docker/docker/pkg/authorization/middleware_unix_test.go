// +build !windows

package authorization

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/docker/docker/pkg/plugingetter"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestMiddlewareWrapHandler(t *testing.T) ***REMOVED***
	server := authZPluginTestServer***REMOVED***t: t***REMOVED***
	server.start()
	defer server.stop()

	authZPlugin := createTestPlugin(t)
	pluginNames := []string***REMOVED***authZPlugin.name***REMOVED***

	var pluginGetter plugingetter.PluginGetter
	middleWare := NewMiddleware(pluginNames, pluginGetter)
	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
		return nil
	***REMOVED***

	authList := []Plugin***REMOVED***authZPlugin***REMOVED***
	middleWare.SetPlugins([]string***REMOVED***"My Test Plugin"***REMOVED***)
	setAuthzPlugins(middleWare, authList)
	mdHandler := middleWare.WrapHandler(handler)
	require.NotNil(t, mdHandler)

	addr := "www.example.com/auth"
	req, _ := http.NewRequest("GET", addr, nil)
	req.RequestURI = addr
	req.Header.Add("header", "value")

	resp := httptest.NewRecorder()
	ctx := context.Background()

	t.Run("Error Test Case :", func(t *testing.T) ***REMOVED***
		server.replayResponse = Response***REMOVED***
			Allow: false,
			Msg:   "Server Auth Not Allowed",
		***REMOVED***
		if err := mdHandler(ctx, resp, req, map[string]string***REMOVED******REMOVED***); err == nil ***REMOVED***
			require.Error(t, err)
		***REMOVED***

	***REMOVED***)

	t.Run("Positive Test Case :", func(t *testing.T) ***REMOVED***
		server.replayResponse = Response***REMOVED***
			Allow: true,
			Msg:   "Server Auth Allowed",
		***REMOVED***
		if err := mdHandler(ctx, resp, req, map[string]string***REMOVED******REMOVED***); err != nil ***REMOVED***
			require.NoError(t, err)
		***REMOVED***

	***REMOVED***)

***REMOVED***
