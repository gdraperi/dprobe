// +build !windows

package authz

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/docker/docker/integration-cli/daemon"
	"github.com/docker/docker/internal/test/environment"
	"github.com/docker/docker/pkg/authorization"
	"github.com/docker/docker/pkg/plugins"
)

var (
	testEnv *environment.Execution
	d       *daemon.Daemon
	server  *httptest.Server
)

const dockerdBinary = "dockerd"

func TestMain(m *testing.M) ***REMOVED***
	var err error
	testEnv, err = environment.New()
	if err != nil ***REMOVED***
		fmt.Println(err)
		os.Exit(1)
	***REMOVED***
	err = environment.EnsureFrozenImagesLinux(testEnv)
	if err != nil ***REMOVED***
		fmt.Println(err)
		os.Exit(1)
	***REMOVED***

	testEnv.Print()
	setupSuite()
	exitCode := m.Run()
	teardownSuite()

	os.Exit(exitCode)
***REMOVED***

func setupTest(t *testing.T) func() ***REMOVED***
	environment.ProtectAll(t, testEnv)

	d = daemon.New(t, "", dockerdBinary, daemon.Config***REMOVED***
		Experimental: testEnv.DaemonInfo.ExperimentalBuild,
	***REMOVED***)

	return func() ***REMOVED***
		if d != nil ***REMOVED***
			d.Stop(t)
		***REMOVED***
		testEnv.Clean(t)
	***REMOVED***
***REMOVED***

func setupSuite() ***REMOVED***
	mux := http.NewServeMux()
	server = httptest.NewServer(mux)

	mux.HandleFunc("/Plugin.Activate", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		b, err := json.Marshal(plugins.Manifest***REMOVED***Implements: []string***REMOVED***authorization.AuthZApiImplements***REMOVED******REMOVED***)
		if err != nil ***REMOVED***
			panic("could not marshal json for /Plugin.Activate: " + err.Error())
		***REMOVED***
		w.Write(b)
	***REMOVED***)

	mux.HandleFunc("/AuthZPlugin.AuthZReq", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil ***REMOVED***
			panic("could not read body for /AuthZPlugin.AuthZReq: " + err.Error())
		***REMOVED***
		authReq := authorization.Request***REMOVED******REMOVED***
		err = json.Unmarshal(body, &authReq)
		if err != nil ***REMOVED***
			panic("could not unmarshal json for /AuthZPlugin.AuthZReq: " + err.Error())
		***REMOVED***

		assertBody(authReq.RequestURI, authReq.RequestHeaders, authReq.RequestBody)
		assertAuthHeaders(authReq.RequestHeaders)

		// Count only server version api
		if strings.HasSuffix(authReq.RequestURI, serverVersionAPI) ***REMOVED***
			ctrl.versionReqCount++
		***REMOVED***

		ctrl.requestsURIs = append(ctrl.requestsURIs, authReq.RequestURI)

		reqRes := ctrl.reqRes
		if isAllowed(authReq.RequestURI) ***REMOVED***
			reqRes = authorization.Response***REMOVED***Allow: true***REMOVED***
		***REMOVED***
		if reqRes.Err != "" ***REMOVED***
			w.WriteHeader(http.StatusInternalServerError)
		***REMOVED***
		b, err := json.Marshal(reqRes)
		if err != nil ***REMOVED***
			panic("could not marshal json for /AuthZPlugin.AuthZReq: " + err.Error())
		***REMOVED***

		ctrl.reqUser = authReq.User
		w.Write(b)
	***REMOVED***)

	mux.HandleFunc("/AuthZPlugin.AuthZRes", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil ***REMOVED***
			panic("could not read body for /AuthZPlugin.AuthZRes: " + err.Error())
		***REMOVED***
		authReq := authorization.Request***REMOVED******REMOVED***
		err = json.Unmarshal(body, &authReq)
		if err != nil ***REMOVED***
			panic("could not unmarshal json for /AuthZPlugin.AuthZRes: " + err.Error())
		***REMOVED***

		assertBody(authReq.RequestURI, authReq.ResponseHeaders, authReq.ResponseBody)
		assertAuthHeaders(authReq.ResponseHeaders)

		// Count only server version api
		if strings.HasSuffix(authReq.RequestURI, serverVersionAPI) ***REMOVED***
			ctrl.versionResCount++
		***REMOVED***
		resRes := ctrl.resRes
		if isAllowed(authReq.RequestURI) ***REMOVED***
			resRes = authorization.Response***REMOVED***Allow: true***REMOVED***
		***REMOVED***
		if resRes.Err != "" ***REMOVED***
			w.WriteHeader(http.StatusInternalServerError)
		***REMOVED***
		b, err := json.Marshal(resRes)
		if err != nil ***REMOVED***
			panic("could not marshal json for /AuthZPlugin.AuthZRes: " + err.Error())
		***REMOVED***
		ctrl.resUser = authReq.User
		w.Write(b)
	***REMOVED***)
***REMOVED***

func teardownSuite() ***REMOVED***
	if server == nil ***REMOVED***
		return
	***REMOVED***

	server.Close()
***REMOVED***

// assertAuthHeaders validates authentication headers are removed
func assertAuthHeaders(headers map[string]string) error ***REMOVED***
	for k := range headers ***REMOVED***
		if strings.Contains(strings.ToLower(k), "auth") || strings.Contains(strings.ToLower(k), "x-registry") ***REMOVED***
			panic(fmt.Sprintf("Found authentication headers in request '%v'", headers))
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// assertBody asserts that body is removed for non text/json requests
func assertBody(requestURI string, headers map[string]string, body []byte) ***REMOVED***
	if strings.Contains(strings.ToLower(requestURI), "auth") && len(body) > 0 ***REMOVED***
		panic("Body included for authentication endpoint " + string(body))
	***REMOVED***

	for k, v := range headers ***REMOVED***
		if strings.EqualFold(k, "Content-Type") && strings.HasPrefix(v, "text/") || v == "application/json" ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	if len(body) > 0 ***REMOVED***
		panic(fmt.Sprintf("Body included while it should not (Headers: '%v')", headers))
	***REMOVED***
***REMOVED***
