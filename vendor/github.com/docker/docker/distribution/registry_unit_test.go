package distribution

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"runtime"
	"strings"
	"testing"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/registry"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const secretRegistryToken = "mysecrettoken"

type tokenPassThruHandler struct ***REMOVED***
	reached       bool
	gotToken      bool
	shouldSend401 func(url string) bool
***REMOVED***

func (h *tokenPassThruHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) ***REMOVED***
	h.reached = true
	if strings.Contains(r.Header.Get("Authorization"), secretRegistryToken) ***REMOVED***
		logrus.Debug("Detected registry token in auth header")
		h.gotToken = true
	***REMOVED***
	if h.shouldSend401 == nil || h.shouldSend401(r.RequestURI) ***REMOVED***
		w.Header().Set("WWW-Authenticate", `Bearer realm="foorealm"`)
		w.WriteHeader(401)
	***REMOVED***
***REMOVED***

func testTokenPassThru(t *testing.T, ts *httptest.Server) ***REMOVED***
	uri, err := url.Parse(ts.URL)
	if err != nil ***REMOVED***
		t.Fatalf("could not parse url from test server: %v", err)
	***REMOVED***

	endpoint := registry.APIEndpoint***REMOVED***
		Mirror:       false,
		URL:          uri,
		Version:      2,
		Official:     false,
		TrimHostname: false,
		TLSConfig:    nil,
	***REMOVED***
	n, _ := reference.ParseNormalizedNamed("testremotename")
	repoInfo := &registry.RepositoryInfo***REMOVED***
		Name: n,
		Index: &registrytypes.IndexInfo***REMOVED***
			Name:     "testrepo",
			Mirrors:  nil,
			Secure:   false,
			Official: false,
		***REMOVED***,
		Official: false,
	***REMOVED***
	imagePullConfig := &ImagePullConfig***REMOVED***
		Config: Config***REMOVED***
			MetaHeaders: http.Header***REMOVED******REMOVED***,
			AuthConfig: &types.AuthConfig***REMOVED***
				RegistryToken: secretRegistryToken,
			***REMOVED***,
		***REMOVED***,
		Schema2Types: ImageTypes,
	***REMOVED***
	puller, err := newPuller(endpoint, repoInfo, imagePullConfig)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	p := puller.(*v2Puller)
	ctx := context.Background()
	p.repo, _, err = NewV2Repository(ctx, p.repoInfo, p.endpoint, p.config.MetaHeaders, p.config.AuthConfig, "pull")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	logrus.Debug("About to pull")
	// We expect it to fail, since we haven't mock'd the full registry exchange in our handler above
	tag, _ := reference.WithTag(n, "tag_goes_here")
	_ = p.pullV2Repository(ctx, tag, runtime.GOOS)
***REMOVED***

func TestTokenPassThru(t *testing.T) ***REMOVED***
	handler := &tokenPassThruHandler***REMOVED***shouldSend401: func(url string) bool ***REMOVED*** return url == "/v2/" ***REMOVED******REMOVED***
	ts := httptest.NewServer(handler)
	defer ts.Close()

	testTokenPassThru(t, ts)

	if !handler.reached ***REMOVED***
		t.Fatal("Handler not reached")
	***REMOVED***
	if !handler.gotToken ***REMOVED***
		t.Fatal("Failed to receive registry token")
	***REMOVED***
***REMOVED***

func TestTokenPassThruDifferentHost(t *testing.T) ***REMOVED***
	handler := new(tokenPassThruHandler)
	ts := httptest.NewServer(handler)
	defer ts.Close()

	tsredirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		if r.RequestURI == "/v2/" ***REMOVED***
			w.Header().Set("WWW-Authenticate", `Bearer realm="foorealm"`)
			w.WriteHeader(401)
			return
		***REMOVED***
		http.Redirect(w, r, ts.URL+r.URL.Path, http.StatusMovedPermanently)
	***REMOVED***))
	defer tsredirect.Close()

	testTokenPassThru(t, tsredirect)

	if !handler.reached ***REMOVED***
		t.Fatal("Handler not reached")
	***REMOVED***
	if handler.gotToken ***REMOVED***
		t.Fatal("Redirect should not forward Authorization header to another host")
	***REMOVED***
***REMOVED***
