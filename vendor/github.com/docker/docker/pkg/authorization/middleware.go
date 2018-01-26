package authorization

import (
	"net/http"
	"sync"

	"github.com/docker/docker/pkg/plugingetter"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// Middleware uses a list of plugins to
// handle authorization in the API requests.
type Middleware struct ***REMOVED***
	mu      sync.Mutex
	plugins []Plugin
***REMOVED***

// NewMiddleware creates a new Middleware
// with a slice of plugins names.
func NewMiddleware(names []string, pg plugingetter.PluginGetter) *Middleware ***REMOVED***
	SetPluginGetter(pg)
	return &Middleware***REMOVED***
		plugins: newPlugins(names),
	***REMOVED***
***REMOVED***

func (m *Middleware) getAuthzPlugins() []Plugin ***REMOVED***
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.plugins
***REMOVED***

// SetPlugins sets the plugin used for authorization
func (m *Middleware) SetPlugins(names []string) ***REMOVED***
	m.mu.Lock()
	m.plugins = newPlugins(names)
	m.mu.Unlock()
***REMOVED***

// RemovePlugin removes a single plugin from this authz middleware chain
func (m *Middleware) RemovePlugin(name string) ***REMOVED***
	m.mu.Lock()
	defer m.mu.Unlock()
	plugins := m.plugins[:0]
	for _, authPlugin := range m.plugins ***REMOVED***
		if authPlugin.Name() != name ***REMOVED***
			plugins = append(plugins, authPlugin)
		***REMOVED***
	***REMOVED***
	m.plugins = plugins
***REMOVED***

// WrapHandler returns a new handler function wrapping the previous one in the request chain.
func (m *Middleware) WrapHandler(handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error) func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
		plugins := m.getAuthzPlugins()
		if len(plugins) == 0 ***REMOVED***
			return handler(ctx, w, r, vars)
		***REMOVED***

		user := ""
		userAuthNMethod := ""

		// Default authorization using existing TLS connection credentials
		// FIXME: Non trivial authorization mechanisms (such as advanced certificate validations, kerberos support
		// and ldap) will be extracted using AuthN feature, which is tracked under:
		// https://github.com/docker/docker/pull/20883
		if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 ***REMOVED***
			user = r.TLS.PeerCertificates[0].Subject.CommonName
			userAuthNMethod = "TLS"
		***REMOVED***

		authCtx := NewCtx(plugins, user, userAuthNMethod, r.Method, r.RequestURI)

		if err := authCtx.AuthZRequest(w, r); err != nil ***REMOVED***
			logrus.Errorf("AuthZRequest for %s %s returned error: %s", r.Method, r.RequestURI, err)
			return err
		***REMOVED***

		rw := NewResponseModifier(w)

		var errD error

		if errD = handler(ctx, rw, r, vars); errD != nil ***REMOVED***
			logrus.Errorf("Handler for %s %s returned error: %s", r.Method, r.RequestURI, errD)
		***REMOVED***

		// There's a chance that the authCtx.plugins was updated. One of the reasons
		// this can happen is when an authzplugin is disabled.
		plugins = m.getAuthzPlugins()
		if len(plugins) == 0 ***REMOVED***
			logrus.Debug("There are no authz plugins in the chain")
			return nil
		***REMOVED***

		authCtx.plugins = plugins

		if err := authCtx.AuthZResponse(rw, r); errD == nil && err != nil ***REMOVED***
			logrus.Errorf("AuthZResponse for %s %s returned error: %s", r.Method, r.RequestURI, err)
			return err
		***REMOVED***

		if errD != nil ***REMOVED***
			return errD
		***REMOVED***

		return nil
	***REMOVED***
***REMOVED***
