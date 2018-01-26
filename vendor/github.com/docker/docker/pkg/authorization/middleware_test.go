package authorization

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/docker/docker/pkg/plugingetter"
	"github.com/stretchr/testify/require"
)

func TestMiddleware(t *testing.T) ***REMOVED***
	pluginNames := []string***REMOVED***"testPlugin1", "testPlugin2"***REMOVED***
	var pluginGetter plugingetter.PluginGetter
	m := NewMiddleware(pluginNames, pluginGetter)
	authPlugins := m.getAuthzPlugins()
	require.Equal(t, 2, len(authPlugins))
	require.EqualValues(t, pluginNames[0], authPlugins[0].Name())
	require.EqualValues(t, pluginNames[1], authPlugins[1].Name())
***REMOVED***

func TestNewResponseModifier(t *testing.T) ***REMOVED***
	recorder := httptest.NewRecorder()
	modifier := NewResponseModifier(recorder)
	modifier.Header().Set("H1", "V1")
	modifier.Write([]byte("body"))
	require.False(t, modifier.Hijacked())
	modifier.WriteHeader(http.StatusInternalServerError)
	require.NotNil(t, modifier.RawBody())

	raw, err := modifier.RawHeaders()
	require.NotNil(t, raw)
	require.Nil(t, err)

	headerData := strings.Split(strings.TrimSpace(string(raw)), ":")
	require.EqualValues(t, "H1", strings.TrimSpace(headerData[0]))
	require.EqualValues(t, "V1", strings.TrimSpace(headerData[1]))

	modifier.Flush()
	modifier.FlushAll()

	if recorder.Header().Get("H1") != "V1" ***REMOVED***
		t.Fatalf("Header value must exists %s", recorder.Header().Get("H1"))
	***REMOVED***

***REMOVED***

func setAuthzPlugins(m *Middleware, plugins []Plugin) ***REMOVED***
	m.mu.Lock()
	m.plugins = plugins
	m.mu.Unlock()
***REMOVED***
