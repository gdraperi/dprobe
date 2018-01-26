package registry

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"testing"

	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/client/transport"
	"github.com/docker/docker/api/types"
	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/stretchr/testify/assert"
)

var (
	token = []string***REMOVED***"fake-token"***REMOVED***
)

const (
	imageID = "42d718c941f5c532ac049bf0b0ab53f0062f09a03afd4aa4a02c098e46032b9d"
	REPO    = "foo42/bar"
)

func spawnTestRegistrySession(t *testing.T) *Session ***REMOVED***
	authConfig := &types.AuthConfig***REMOVED******REMOVED***
	endpoint, err := NewV1Endpoint(makeIndex("/v1/"), "", nil)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	userAgent := "docker test client"
	var tr http.RoundTripper = debugTransport***REMOVED***NewTransport(nil), t.Log***REMOVED***
	tr = transport.NewTransport(AuthTransport(tr, authConfig, false), Headers(userAgent, nil)...)
	client := HTTPClient(tr)
	r, err := NewSession(client, authConfig, endpoint)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	// In a normal scenario for the v1 registry, the client should send a `X-Docker-Token: true`
	// header while authenticating, in order to retrieve a token that can be later used to
	// perform authenticated actions.
	//
	// The mock v1 registry does not support that, (TODO(tiborvass): support it), instead,
	// it will consider authenticated any request with the header `X-Docker-Token: fake-token`.
	//
	// Because we know that the client's transport is an `*authTransport` we simply cast it,
	// in order to set the internal cached token to the fake token, and thus send that fake token
	// upon every subsequent requests.
	r.client.Transport.(*authTransport).token = token
	return r
***REMOVED***

func TestPingRegistryEndpoint(t *testing.T) ***REMOVED***
	testPing := func(index *registrytypes.IndexInfo, expectedStandalone bool, assertMessage string) ***REMOVED***
		ep, err := NewV1Endpoint(index, "", nil)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		regInfo, err := ep.Ping()
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		assertEqual(t, regInfo.Standalone, expectedStandalone, assertMessage)
	***REMOVED***

	testPing(makeIndex("/v1/"), true, "Expected standalone to be true (default)")
	testPing(makeHTTPSIndex("/v1/"), true, "Expected standalone to be true (default)")
	testPing(makePublicIndex(), false, "Expected standalone to be false for public index")
***REMOVED***

func TestEndpoint(t *testing.T) ***REMOVED***
	// Simple wrapper to fail test if err != nil
	expandEndpoint := func(index *registrytypes.IndexInfo) *V1Endpoint ***REMOVED***
		endpoint, err := NewV1Endpoint(index, "", nil)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		return endpoint
	***REMOVED***

	assertInsecureIndex := func(index *registrytypes.IndexInfo) ***REMOVED***
		index.Secure = true
		_, err := NewV1Endpoint(index, "", nil)
		assertNotEqual(t, err, nil, index.Name+": Expected error for insecure index")
		assertEqual(t, strings.Contains(err.Error(), "insecure-registry"), true, index.Name+": Expected insecure-registry  error for insecure index")
		index.Secure = false
	***REMOVED***

	assertSecureIndex := func(index *registrytypes.IndexInfo) ***REMOVED***
		index.Secure = true
		_, err := NewV1Endpoint(index, "", nil)
		assertNotEqual(t, err, nil, index.Name+": Expected cert error for secure index")
		assertEqual(t, strings.Contains(err.Error(), "certificate signed by unknown authority"), true, index.Name+": Expected cert error for secure index")
		index.Secure = false
	***REMOVED***

	index := &registrytypes.IndexInfo***REMOVED******REMOVED***
	index.Name = makeURL("/v1/")
	endpoint := expandEndpoint(index)
	assertEqual(t, endpoint.String(), index.Name, "Expected endpoint to be "+index.Name)
	assertInsecureIndex(index)

	index.Name = makeURL("")
	endpoint = expandEndpoint(index)
	assertEqual(t, endpoint.String(), index.Name+"/v1/", index.Name+": Expected endpoint to be "+index.Name+"/v1/")
	assertInsecureIndex(index)

	httpURL := makeURL("")
	index.Name = strings.SplitN(httpURL, "://", 2)[1]
	endpoint = expandEndpoint(index)
	assertEqual(t, endpoint.String(), httpURL+"/v1/", index.Name+": Expected endpoint to be "+httpURL+"/v1/")
	assertInsecureIndex(index)

	index.Name = makeHTTPSURL("/v1/")
	endpoint = expandEndpoint(index)
	assertEqual(t, endpoint.String(), index.Name, "Expected endpoint to be "+index.Name)
	assertSecureIndex(index)

	index.Name = makeHTTPSURL("")
	endpoint = expandEndpoint(index)
	assertEqual(t, endpoint.String(), index.Name+"/v1/", index.Name+": Expected endpoint to be "+index.Name+"/v1/")
	assertSecureIndex(index)

	httpsURL := makeHTTPSURL("")
	index.Name = strings.SplitN(httpsURL, "://", 2)[1]
	endpoint = expandEndpoint(index)
	assertEqual(t, endpoint.String(), httpsURL+"/v1/", index.Name+": Expected endpoint to be "+httpsURL+"/v1/")
	assertSecureIndex(index)

	badEndpoints := []string***REMOVED***
		"http://127.0.0.1/v1/",
		"https://127.0.0.1/v1/",
		"http://127.0.0.1",
		"https://127.0.0.1",
		"127.0.0.1",
	***REMOVED***
	for _, address := range badEndpoints ***REMOVED***
		index.Name = address
		_, err := NewV1Endpoint(index, "", nil)
		checkNotEqual(t, err, nil, "Expected error while expanding bad endpoint")
	***REMOVED***
***REMOVED***

func TestGetRemoteHistory(t *testing.T) ***REMOVED***
	r := spawnTestRegistrySession(t)
	hist, err := r.GetRemoteHistory(imageID, makeURL("/v1/"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	assertEqual(t, len(hist), 2, "Expected 2 images in history")
	assertEqual(t, hist[0], imageID, "Expected "+imageID+"as first ancestry")
	assertEqual(t, hist[1], "77dbf71da1d00e3fbddc480176eac8994025630c6590d11cfc8fe1209c2a1d20",
		"Unexpected second ancestry")
***REMOVED***

func TestLookupRemoteImage(t *testing.T) ***REMOVED***
	r := spawnTestRegistrySession(t)
	err := r.LookupRemoteImage(imageID, makeURL("/v1/"))
	assertEqual(t, err, nil, "Expected error of remote lookup to nil")
	if err := r.LookupRemoteImage("abcdef", makeURL("/v1/")); err == nil ***REMOVED***
		t.Fatal("Expected error of remote lookup to not nil")
	***REMOVED***
***REMOVED***

func TestGetRemoteImageJSON(t *testing.T) ***REMOVED***
	r := spawnTestRegistrySession(t)
	json, size, err := r.GetRemoteImageJSON(imageID, makeURL("/v1/"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	assertEqual(t, size, int64(154), "Expected size 154")
	if len(json) == 0 ***REMOVED***
		t.Fatal("Expected non-empty json")
	***REMOVED***

	_, _, err = r.GetRemoteImageJSON("abcdef", makeURL("/v1/"))
	if err == nil ***REMOVED***
		t.Fatal("Expected image not found error")
	***REMOVED***
***REMOVED***

func TestGetRemoteImageLayer(t *testing.T) ***REMOVED***
	r := spawnTestRegistrySession(t)
	data, err := r.GetRemoteImageLayer(imageID, makeURL("/v1/"), 0)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if data == nil ***REMOVED***
		t.Fatal("Expected non-nil data result")
	***REMOVED***

	_, err = r.GetRemoteImageLayer("abcdef", makeURL("/v1/"), 0)
	if err == nil ***REMOVED***
		t.Fatal("Expected image not found error")
	***REMOVED***
***REMOVED***

func TestGetRemoteTag(t *testing.T) ***REMOVED***
	r := spawnTestRegistrySession(t)
	repoRef, err := reference.ParseNormalizedNamed(REPO)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	tag, err := r.GetRemoteTag([]string***REMOVED***makeURL("/v1/")***REMOVED***, repoRef, "test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	assertEqual(t, tag, imageID, "Expected tag test to map to "+imageID)

	bazRef, err := reference.ParseNormalizedNamed("foo42/baz")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	_, err = r.GetRemoteTag([]string***REMOVED***makeURL("/v1/")***REMOVED***, bazRef, "foo")
	if err != ErrRepoNotFound ***REMOVED***
		t.Fatal("Expected ErrRepoNotFound error when fetching tag for bogus repo")
	***REMOVED***
***REMOVED***

func TestGetRemoteTags(t *testing.T) ***REMOVED***
	r := spawnTestRegistrySession(t)
	repoRef, err := reference.ParseNormalizedNamed(REPO)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	tags, err := r.GetRemoteTags([]string***REMOVED***makeURL("/v1/")***REMOVED***, repoRef)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	assertEqual(t, len(tags), 2, "Expected two tags")
	assertEqual(t, tags["latest"], imageID, "Expected tag latest to map to "+imageID)
	assertEqual(t, tags["test"], imageID, "Expected tag test to map to "+imageID)

	bazRef, err := reference.ParseNormalizedNamed("foo42/baz")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	_, err = r.GetRemoteTags([]string***REMOVED***makeURL("/v1/")***REMOVED***, bazRef)
	if err != ErrRepoNotFound ***REMOVED***
		t.Fatal("Expected ErrRepoNotFound error when fetching tags for bogus repo")
	***REMOVED***
***REMOVED***

func TestGetRepositoryData(t *testing.T) ***REMOVED***
	r := spawnTestRegistrySession(t)
	parsedURL, err := url.Parse(makeURL("/v1/"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	host := "http://" + parsedURL.Host + "/v1/"
	repoRef, err := reference.ParseNormalizedNamed(REPO)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	data, err := r.GetRepositoryData(repoRef)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	assertEqual(t, len(data.ImgList), 2, "Expected 2 images in ImgList")
	assertEqual(t, len(data.Endpoints), 2,
		fmt.Sprintf("Expected 2 endpoints in Endpoints, found %d instead", len(data.Endpoints)))
	assertEqual(t, data.Endpoints[0], host,
		fmt.Sprintf("Expected first endpoint to be %s but found %s instead", host, data.Endpoints[0]))
	assertEqual(t, data.Endpoints[1], "http://test.example.com/v1/",
		fmt.Sprintf("Expected first endpoint to be http://test.example.com/v1/ but found %s instead", data.Endpoints[1]))

***REMOVED***

func TestPushImageJSONRegistry(t *testing.T) ***REMOVED***
	r := spawnTestRegistrySession(t)
	imgData := &ImgData***REMOVED***
		ID:       "77dbf71da1d00e3fbddc480176eac8994025630c6590d11cfc8fe1209c2a1d20",
		Checksum: "sha256:1ac330d56e05eef6d438586545ceff7550d3bdcb6b19961f12c5ba714ee1bb37",
	***REMOVED***

	err := r.PushImageJSONRegistry(imgData, []byte***REMOVED***0x42, 0xdf, 0x0***REMOVED***, makeURL("/v1/"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestPushImageLayerRegistry(t *testing.T) ***REMOVED***
	r := spawnTestRegistrySession(t)
	layer := strings.NewReader("")
	_, _, err := r.PushImageLayerRegistry(imageID, layer, makeURL("/v1/"), []byte***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestParseRepositoryInfo(t *testing.T) ***REMOVED***
	type staticRepositoryInfo struct ***REMOVED***
		Index         *registrytypes.IndexInfo
		RemoteName    string
		CanonicalName string
		LocalName     string
		Official      bool
	***REMOVED***

	expectedRepoInfos := map[string]staticRepositoryInfo***REMOVED***
		"fooo/bar": ***REMOVED***
			Index: &registrytypes.IndexInfo***REMOVED***
				Name:     IndexName,
				Official: true,
			***REMOVED***,
			RemoteName:    "fooo/bar",
			LocalName:     "fooo/bar",
			CanonicalName: "docker.io/fooo/bar",
			Official:      false,
		***REMOVED***,
		"library/ubuntu": ***REMOVED***
			Index: &registrytypes.IndexInfo***REMOVED***
				Name:     IndexName,
				Official: true,
			***REMOVED***,
			RemoteName:    "library/ubuntu",
			LocalName:     "ubuntu",
			CanonicalName: "docker.io/library/ubuntu",
			Official:      true,
		***REMOVED***,
		"nonlibrary/ubuntu": ***REMOVED***
			Index: &registrytypes.IndexInfo***REMOVED***
				Name:     IndexName,
				Official: true,
			***REMOVED***,
			RemoteName:    "nonlibrary/ubuntu",
			LocalName:     "nonlibrary/ubuntu",
			CanonicalName: "docker.io/nonlibrary/ubuntu",
			Official:      false,
		***REMOVED***,
		"ubuntu": ***REMOVED***
			Index: &registrytypes.IndexInfo***REMOVED***
				Name:     IndexName,
				Official: true,
			***REMOVED***,
			RemoteName:    "library/ubuntu",
			LocalName:     "ubuntu",
			CanonicalName: "docker.io/library/ubuntu",
			Official:      true,
		***REMOVED***,
		"other/library": ***REMOVED***
			Index: &registrytypes.IndexInfo***REMOVED***
				Name:     IndexName,
				Official: true,
			***REMOVED***,
			RemoteName:    "other/library",
			LocalName:     "other/library",
			CanonicalName: "docker.io/other/library",
			Official:      false,
		***REMOVED***,
		"127.0.0.1:8000/private/moonbase": ***REMOVED***
			Index: &registrytypes.IndexInfo***REMOVED***
				Name:     "127.0.0.1:8000",
				Official: false,
			***REMOVED***,
			RemoteName:    "private/moonbase",
			LocalName:     "127.0.0.1:8000/private/moonbase",
			CanonicalName: "127.0.0.1:8000/private/moonbase",
			Official:      false,
		***REMOVED***,
		"127.0.0.1:8000/privatebase": ***REMOVED***
			Index: &registrytypes.IndexInfo***REMOVED***
				Name:     "127.0.0.1:8000",
				Official: false,
			***REMOVED***,
			RemoteName:    "privatebase",
			LocalName:     "127.0.0.1:8000/privatebase",
			CanonicalName: "127.0.0.1:8000/privatebase",
			Official:      false,
		***REMOVED***,
		"localhost:8000/private/moonbase": ***REMOVED***
			Index: &registrytypes.IndexInfo***REMOVED***
				Name:     "localhost:8000",
				Official: false,
			***REMOVED***,
			RemoteName:    "private/moonbase",
			LocalName:     "localhost:8000/private/moonbase",
			CanonicalName: "localhost:8000/private/moonbase",
			Official:      false,
		***REMOVED***,
		"localhost:8000/privatebase": ***REMOVED***
			Index: &registrytypes.IndexInfo***REMOVED***
				Name:     "localhost:8000",
				Official: false,
			***REMOVED***,
			RemoteName:    "privatebase",
			LocalName:     "localhost:8000/privatebase",
			CanonicalName: "localhost:8000/privatebase",
			Official:      false,
		***REMOVED***,
		"example.com/private/moonbase": ***REMOVED***
			Index: &registrytypes.IndexInfo***REMOVED***
				Name:     "example.com",
				Official: false,
			***REMOVED***,
			RemoteName:    "private/moonbase",
			LocalName:     "example.com/private/moonbase",
			CanonicalName: "example.com/private/moonbase",
			Official:      false,
		***REMOVED***,
		"example.com/privatebase": ***REMOVED***
			Index: &registrytypes.IndexInfo***REMOVED***
				Name:     "example.com",
				Official: false,
			***REMOVED***,
			RemoteName:    "privatebase",
			LocalName:     "example.com/privatebase",
			CanonicalName: "example.com/privatebase",
			Official:      false,
		***REMOVED***,
		"example.com:8000/private/moonbase": ***REMOVED***
			Index: &registrytypes.IndexInfo***REMOVED***
				Name:     "example.com:8000",
				Official: false,
			***REMOVED***,
			RemoteName:    "private/moonbase",
			LocalName:     "example.com:8000/private/moonbase",
			CanonicalName: "example.com:8000/private/moonbase",
			Official:      false,
		***REMOVED***,
		"example.com:8000/privatebase": ***REMOVED***
			Index: &registrytypes.IndexInfo***REMOVED***
				Name:     "example.com:8000",
				Official: false,
			***REMOVED***,
			RemoteName:    "privatebase",
			LocalName:     "example.com:8000/privatebase",
			CanonicalName: "example.com:8000/privatebase",
			Official:      false,
		***REMOVED***,
		"localhost/private/moonbase": ***REMOVED***
			Index: &registrytypes.IndexInfo***REMOVED***
				Name:     "localhost",
				Official: false,
			***REMOVED***,
			RemoteName:    "private/moonbase",
			LocalName:     "localhost/private/moonbase",
			CanonicalName: "localhost/private/moonbase",
			Official:      false,
		***REMOVED***,
		"localhost/privatebase": ***REMOVED***
			Index: &registrytypes.IndexInfo***REMOVED***
				Name:     "localhost",
				Official: false,
			***REMOVED***,
			RemoteName:    "privatebase",
			LocalName:     "localhost/privatebase",
			CanonicalName: "localhost/privatebase",
			Official:      false,
		***REMOVED***,
		IndexName + "/public/moonbase": ***REMOVED***
			Index: &registrytypes.IndexInfo***REMOVED***
				Name:     IndexName,
				Official: true,
			***REMOVED***,
			RemoteName:    "public/moonbase",
			LocalName:     "public/moonbase",
			CanonicalName: "docker.io/public/moonbase",
			Official:      false,
		***REMOVED***,
		"index." + IndexName + "/public/moonbase": ***REMOVED***
			Index: &registrytypes.IndexInfo***REMOVED***
				Name:     IndexName,
				Official: true,
			***REMOVED***,
			RemoteName:    "public/moonbase",
			LocalName:     "public/moonbase",
			CanonicalName: "docker.io/public/moonbase",
			Official:      false,
		***REMOVED***,
		"ubuntu-12.04-base": ***REMOVED***
			Index: &registrytypes.IndexInfo***REMOVED***
				Name:     IndexName,
				Official: true,
			***REMOVED***,
			RemoteName:    "library/ubuntu-12.04-base",
			LocalName:     "ubuntu-12.04-base",
			CanonicalName: "docker.io/library/ubuntu-12.04-base",
			Official:      true,
		***REMOVED***,
		IndexName + "/ubuntu-12.04-base": ***REMOVED***
			Index: &registrytypes.IndexInfo***REMOVED***
				Name:     IndexName,
				Official: true,
			***REMOVED***,
			RemoteName:    "library/ubuntu-12.04-base",
			LocalName:     "ubuntu-12.04-base",
			CanonicalName: "docker.io/library/ubuntu-12.04-base",
			Official:      true,
		***REMOVED***,
		"index." + IndexName + "/ubuntu-12.04-base": ***REMOVED***
			Index: &registrytypes.IndexInfo***REMOVED***
				Name:     IndexName,
				Official: true,
			***REMOVED***,
			RemoteName:    "library/ubuntu-12.04-base",
			LocalName:     "ubuntu-12.04-base",
			CanonicalName: "docker.io/library/ubuntu-12.04-base",
			Official:      true,
		***REMOVED***,
	***REMOVED***

	for reposName, expectedRepoInfo := range expectedRepoInfos ***REMOVED***
		named, err := reference.ParseNormalizedNamed(reposName)
		if err != nil ***REMOVED***
			t.Error(err)
		***REMOVED***

		repoInfo, err := ParseRepositoryInfo(named)
		if err != nil ***REMOVED***
			t.Error(err)
		***REMOVED*** else ***REMOVED***
			checkEqual(t, repoInfo.Index.Name, expectedRepoInfo.Index.Name, reposName)
			checkEqual(t, reference.Path(repoInfo.Name), expectedRepoInfo.RemoteName, reposName)
			checkEqual(t, reference.FamiliarName(repoInfo.Name), expectedRepoInfo.LocalName, reposName)
			checkEqual(t, repoInfo.Name.Name(), expectedRepoInfo.CanonicalName, reposName)
			checkEqual(t, repoInfo.Index.Official, expectedRepoInfo.Index.Official, reposName)
			checkEqual(t, repoInfo.Official, expectedRepoInfo.Official, reposName)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestNewIndexInfo(t *testing.T) ***REMOVED***
	testIndexInfo := func(config *serviceConfig, expectedIndexInfos map[string]*registrytypes.IndexInfo) ***REMOVED***
		for indexName, expectedIndexInfo := range expectedIndexInfos ***REMOVED***
			index, err := newIndexInfo(config, indexName)
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED*** else ***REMOVED***
				checkEqual(t, index.Name, expectedIndexInfo.Name, indexName+" name")
				checkEqual(t, index.Official, expectedIndexInfo.Official, indexName+" is official")
				checkEqual(t, index.Secure, expectedIndexInfo.Secure, indexName+" is secure")
				checkEqual(t, len(index.Mirrors), len(expectedIndexInfo.Mirrors), indexName+" mirrors")
			***REMOVED***
		***REMOVED***
	***REMOVED***

	config := emptyServiceConfig
	noMirrors := []string***REMOVED******REMOVED***
	expectedIndexInfos := map[string]*registrytypes.IndexInfo***REMOVED***
		IndexName: ***REMOVED***
			Name:     IndexName,
			Official: true,
			Secure:   true,
			Mirrors:  noMirrors,
		***REMOVED***,
		"index." + IndexName: ***REMOVED***
			Name:     IndexName,
			Official: true,
			Secure:   true,
			Mirrors:  noMirrors,
		***REMOVED***,
		"example.com": ***REMOVED***
			Name:     "example.com",
			Official: false,
			Secure:   true,
			Mirrors:  noMirrors,
		***REMOVED***,
		"127.0.0.1:5000": ***REMOVED***
			Name:     "127.0.0.1:5000",
			Official: false,
			Secure:   false,
			Mirrors:  noMirrors,
		***REMOVED***,
	***REMOVED***
	testIndexInfo(config, expectedIndexInfos)

	publicMirrors := []string***REMOVED***"http://mirror1.local", "http://mirror2.local"***REMOVED***
	var err error
	config, err = makeServiceConfig(publicMirrors, []string***REMOVED***"example.com"***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	expectedIndexInfos = map[string]*registrytypes.IndexInfo***REMOVED***
		IndexName: ***REMOVED***
			Name:     IndexName,
			Official: true,
			Secure:   true,
			Mirrors:  publicMirrors,
		***REMOVED***,
		"index." + IndexName: ***REMOVED***
			Name:     IndexName,
			Official: true,
			Secure:   true,
			Mirrors:  publicMirrors,
		***REMOVED***,
		"example.com": ***REMOVED***
			Name:     "example.com",
			Official: false,
			Secure:   false,
			Mirrors:  noMirrors,
		***REMOVED***,
		"example.com:5000": ***REMOVED***
			Name:     "example.com:5000",
			Official: false,
			Secure:   true,
			Mirrors:  noMirrors,
		***REMOVED***,
		"127.0.0.1": ***REMOVED***
			Name:     "127.0.0.1",
			Official: false,
			Secure:   false,
			Mirrors:  noMirrors,
		***REMOVED***,
		"127.0.0.1:5000": ***REMOVED***
			Name:     "127.0.0.1:5000",
			Official: false,
			Secure:   false,
			Mirrors:  noMirrors,
		***REMOVED***,
		"other.com": ***REMOVED***
			Name:     "other.com",
			Official: false,
			Secure:   true,
			Mirrors:  noMirrors,
		***REMOVED***,
	***REMOVED***
	testIndexInfo(config, expectedIndexInfos)

	config, err = makeServiceConfig(nil, []string***REMOVED***"42.42.0.0/16"***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	expectedIndexInfos = map[string]*registrytypes.IndexInfo***REMOVED***
		"example.com": ***REMOVED***
			Name:     "example.com",
			Official: false,
			Secure:   false,
			Mirrors:  noMirrors,
		***REMOVED***,
		"example.com:5000": ***REMOVED***
			Name:     "example.com:5000",
			Official: false,
			Secure:   false,
			Mirrors:  noMirrors,
		***REMOVED***,
		"127.0.0.1": ***REMOVED***
			Name:     "127.0.0.1",
			Official: false,
			Secure:   false,
			Mirrors:  noMirrors,
		***REMOVED***,
		"127.0.0.1:5000": ***REMOVED***
			Name:     "127.0.0.1:5000",
			Official: false,
			Secure:   false,
			Mirrors:  noMirrors,
		***REMOVED***,
		"other.com": ***REMOVED***
			Name:     "other.com",
			Official: false,
			Secure:   true,
			Mirrors:  noMirrors,
		***REMOVED***,
	***REMOVED***
	testIndexInfo(config, expectedIndexInfos)
***REMOVED***

func TestMirrorEndpointLookup(t *testing.T) ***REMOVED***
	containsMirror := func(endpoints []APIEndpoint) bool ***REMOVED***
		for _, pe := range endpoints ***REMOVED***
			if pe.URL.Host == "my.mirror" ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
		return false
	***REMOVED***
	cfg, err := makeServiceConfig([]string***REMOVED***"https://my.mirror"***REMOVED***, nil)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	s := DefaultService***REMOVED***config: cfg***REMOVED***

	imageName, err := reference.WithName(IndexName + "/test/image")
	if err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
	pushAPIEndpoints, err := s.LookupPushEndpoints(reference.Domain(imageName))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if containsMirror(pushAPIEndpoints) ***REMOVED***
		t.Fatal("Push endpoint should not contain mirror")
	***REMOVED***

	pullAPIEndpoints, err := s.LookupPullEndpoints(reference.Domain(imageName))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if !containsMirror(pullAPIEndpoints) ***REMOVED***
		t.Fatal("Pull endpoint should contain mirror")
	***REMOVED***
***REMOVED***

func TestPushRegistryTag(t *testing.T) ***REMOVED***
	r := spawnTestRegistrySession(t)
	repoRef, err := reference.ParseNormalizedNamed(REPO)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	err = r.PushRegistryTag(repoRef, imageID, "stable", makeURL("/v1/"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestPushImageJSONIndex(t *testing.T) ***REMOVED***
	r := spawnTestRegistrySession(t)
	imgData := []*ImgData***REMOVED***
		***REMOVED***
			ID:       "77dbf71da1d00e3fbddc480176eac8994025630c6590d11cfc8fe1209c2a1d20",
			Checksum: "sha256:1ac330d56e05eef6d438586545ceff7550d3bdcb6b19961f12c5ba714ee1bb37",
		***REMOVED***,
		***REMOVED***
			ID:       "42d718c941f5c532ac049bf0b0ab53f0062f09a03afd4aa4a02c098e46032b9d",
			Checksum: "sha256:bea7bf2e4bacd479344b737328db47b18880d09096e6674165533aa994f5e9f2",
		***REMOVED***,
	***REMOVED***
	repoRef, err := reference.ParseNormalizedNamed(REPO)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	repoData, err := r.PushImageJSONIndex(repoRef, imgData, false, nil)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if repoData == nil ***REMOVED***
		t.Fatal("Expected RepositoryData object")
	***REMOVED***
	repoData, err = r.PushImageJSONIndex(repoRef, imgData, true, []string***REMOVED***r.indexEndpoint.String()***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if repoData == nil ***REMOVED***
		t.Fatal("Expected RepositoryData object")
	***REMOVED***
***REMOVED***

func TestSearchRepositories(t *testing.T) ***REMOVED***
	r := spawnTestRegistrySession(t)
	results, err := r.SearchRepositories("fakequery", 25)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if results == nil ***REMOVED***
		t.Fatal("Expected non-nil SearchResults object")
	***REMOVED***
	assertEqual(t, results.NumResults, 1, "Expected 1 search results")
	assertEqual(t, results.Query, "fakequery", "Expected 'fakequery' as query")
	assertEqual(t, results.Results[0].StarCount, 42, "Expected 'fakeimage' to have 42 stars")
***REMOVED***

func TestTrustedLocation(t *testing.T) ***REMOVED***
	for _, url := range []string***REMOVED***"http://example.com", "https://example.com:7777", "http://docker.io", "http://test.docker.com", "https://fakedocker.com"***REMOVED*** ***REMOVED***
		req, _ := http.NewRequest("GET", url, nil)
		assert.False(t, trustedLocation(req))
	***REMOVED***

	for _, url := range []string***REMOVED***"https://docker.io", "https://test.docker.com:80"***REMOVED*** ***REMOVED***
		req, _ := http.NewRequest("GET", url, nil)
		assert.True(t, trustedLocation(req))
	***REMOVED***
***REMOVED***

func TestAddRequiredHeadersToRedirectedRequests(t *testing.T) ***REMOVED***
	for _, urls := range [][]string***REMOVED***
		***REMOVED***"http://docker.io", "https://docker.com"***REMOVED***,
		***REMOVED***"https://foo.docker.io:7777", "http://bar.docker.com"***REMOVED***,
		***REMOVED***"https://foo.docker.io", "https://example.com"***REMOVED***,
	***REMOVED*** ***REMOVED***
		reqFrom, _ := http.NewRequest("GET", urls[0], nil)
		reqFrom.Header.Add("Content-Type", "application/json")
		reqFrom.Header.Add("Authorization", "super_secret")
		reqTo, _ := http.NewRequest("GET", urls[1], nil)

		addRequiredHeadersToRedirectedRequests(reqTo, []*http.Request***REMOVED***reqFrom***REMOVED***)

		if len(reqTo.Header) != 1 ***REMOVED***
			t.Fatalf("Expected 1 headers, got %d", len(reqTo.Header))
		***REMOVED***

		if reqTo.Header.Get("Content-Type") != "application/json" ***REMOVED***
			t.Fatal("'Content-Type' should be 'application/json'")
		***REMOVED***

		if reqTo.Header.Get("Authorization") != "" ***REMOVED***
			t.Fatal("'Authorization' should be empty")
		***REMOVED***
	***REMOVED***

	for _, urls := range [][]string***REMOVED***
		***REMOVED***"https://docker.io", "https://docker.com"***REMOVED***,
		***REMOVED***"https://foo.docker.io:7777", "https://bar.docker.com"***REMOVED***,
	***REMOVED*** ***REMOVED***
		reqFrom, _ := http.NewRequest("GET", urls[0], nil)
		reqFrom.Header.Add("Content-Type", "application/json")
		reqFrom.Header.Add("Authorization", "super_secret")
		reqTo, _ := http.NewRequest("GET", urls[1], nil)

		addRequiredHeadersToRedirectedRequests(reqTo, []*http.Request***REMOVED***reqFrom***REMOVED***)

		if len(reqTo.Header) != 2 ***REMOVED***
			t.Fatalf("Expected 2 headers, got %d", len(reqTo.Header))
		***REMOVED***

		if reqTo.Header.Get("Content-Type") != "application/json" ***REMOVED***
			t.Fatal("'Content-Type' should be 'application/json'")
		***REMOVED***

		if reqTo.Header.Get("Authorization") != "super_secret" ***REMOVED***
			t.Fatal("'Authorization' should be 'super_secret'")
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestAllowNondistributableArtifacts(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		addr       string
		registries []string
		expected   bool
	***REMOVED******REMOVED***
		***REMOVED***IndexName, nil, false***REMOVED***,
		***REMOVED***"example.com", []string***REMOVED******REMOVED***, false***REMOVED***,
		***REMOVED***"example.com", []string***REMOVED***"example.com"***REMOVED***, true***REMOVED***,
		***REMOVED***"localhost", []string***REMOVED***"localhost:5000"***REMOVED***, false***REMOVED***,
		***REMOVED***"localhost:5000", []string***REMOVED***"localhost:5000"***REMOVED***, true***REMOVED***,
		***REMOVED***"localhost", []string***REMOVED***"example.com"***REMOVED***, false***REMOVED***,
		***REMOVED***"127.0.0.1:5000", []string***REMOVED***"127.0.0.1:5000"***REMOVED***, true***REMOVED***,
		***REMOVED***"localhost", nil, false***REMOVED***,
		***REMOVED***"localhost:5000", nil, false***REMOVED***,
		***REMOVED***"127.0.0.1", nil, false***REMOVED***,
		***REMOVED***"localhost", []string***REMOVED***"example.com"***REMOVED***, false***REMOVED***,
		***REMOVED***"127.0.0.1", []string***REMOVED***"example.com"***REMOVED***, false***REMOVED***,
		***REMOVED***"example.com", nil, false***REMOVED***,
		***REMOVED***"example.com", []string***REMOVED***"example.com"***REMOVED***, true***REMOVED***,
		***REMOVED***"127.0.0.1", []string***REMOVED***"example.com"***REMOVED***, false***REMOVED***,
		***REMOVED***"127.0.0.1:5000", []string***REMOVED***"example.com"***REMOVED***, false***REMOVED***,
		***REMOVED***"example.com:5000", []string***REMOVED***"42.42.0.0/16"***REMOVED***, true***REMOVED***,
		***REMOVED***"example.com", []string***REMOVED***"42.42.0.0/16"***REMOVED***, true***REMOVED***,
		***REMOVED***"example.com:5000", []string***REMOVED***"42.42.42.42/8"***REMOVED***, true***REMOVED***,
		***REMOVED***"127.0.0.1:5000", []string***REMOVED***"127.0.0.0/8"***REMOVED***, true***REMOVED***,
		***REMOVED***"42.42.42.42:5000", []string***REMOVED***"42.1.1.1/8"***REMOVED***, true***REMOVED***,
		***REMOVED***"invalid.domain.com", []string***REMOVED***"42.42.0.0/16"***REMOVED***, false***REMOVED***,
		***REMOVED***"invalid.domain.com", []string***REMOVED***"invalid.domain.com"***REMOVED***, true***REMOVED***,
		***REMOVED***"invalid.domain.com:5000", []string***REMOVED***"invalid.domain.com"***REMOVED***, false***REMOVED***,
		***REMOVED***"invalid.domain.com:5000", []string***REMOVED***"invalid.domain.com:5000"***REMOVED***, true***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		config, err := newServiceConfig(ServiceOptions***REMOVED***
			AllowNondistributableArtifacts: tt.registries,
		***REMOVED***)
		if err != nil ***REMOVED***
			t.Error(err)
		***REMOVED***
		if v := allowNondistributableArtifacts(config, tt.addr); v != tt.expected ***REMOVED***
			t.Errorf("allowNondistributableArtifacts failed for %q %v, expected %v got %v", tt.addr, tt.registries, tt.expected, v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestIsSecureIndex(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		addr               string
		insecureRegistries []string
		expected           bool
	***REMOVED******REMOVED***
		***REMOVED***IndexName, nil, true***REMOVED***,
		***REMOVED***"example.com", []string***REMOVED******REMOVED***, true***REMOVED***,
		***REMOVED***"example.com", []string***REMOVED***"example.com"***REMOVED***, false***REMOVED***,
		***REMOVED***"localhost", []string***REMOVED***"localhost:5000"***REMOVED***, false***REMOVED***,
		***REMOVED***"localhost:5000", []string***REMOVED***"localhost:5000"***REMOVED***, false***REMOVED***,
		***REMOVED***"localhost", []string***REMOVED***"example.com"***REMOVED***, false***REMOVED***,
		***REMOVED***"127.0.0.1:5000", []string***REMOVED***"127.0.0.1:5000"***REMOVED***, false***REMOVED***,
		***REMOVED***"localhost", nil, false***REMOVED***,
		***REMOVED***"localhost:5000", nil, false***REMOVED***,
		***REMOVED***"127.0.0.1", nil, false***REMOVED***,
		***REMOVED***"localhost", []string***REMOVED***"example.com"***REMOVED***, false***REMOVED***,
		***REMOVED***"127.0.0.1", []string***REMOVED***"example.com"***REMOVED***, false***REMOVED***,
		***REMOVED***"example.com", nil, true***REMOVED***,
		***REMOVED***"example.com", []string***REMOVED***"example.com"***REMOVED***, false***REMOVED***,
		***REMOVED***"127.0.0.1", []string***REMOVED***"example.com"***REMOVED***, false***REMOVED***,
		***REMOVED***"127.0.0.1:5000", []string***REMOVED***"example.com"***REMOVED***, false***REMOVED***,
		***REMOVED***"example.com:5000", []string***REMOVED***"42.42.0.0/16"***REMOVED***, false***REMOVED***,
		***REMOVED***"example.com", []string***REMOVED***"42.42.0.0/16"***REMOVED***, false***REMOVED***,
		***REMOVED***"example.com:5000", []string***REMOVED***"42.42.42.42/8"***REMOVED***, false***REMOVED***,
		***REMOVED***"127.0.0.1:5000", []string***REMOVED***"127.0.0.0/8"***REMOVED***, false***REMOVED***,
		***REMOVED***"42.42.42.42:5000", []string***REMOVED***"42.1.1.1/8"***REMOVED***, false***REMOVED***,
		***REMOVED***"invalid.domain.com", []string***REMOVED***"42.42.0.0/16"***REMOVED***, true***REMOVED***,
		***REMOVED***"invalid.domain.com", []string***REMOVED***"invalid.domain.com"***REMOVED***, false***REMOVED***,
		***REMOVED***"invalid.domain.com:5000", []string***REMOVED***"invalid.domain.com"***REMOVED***, true***REMOVED***,
		***REMOVED***"invalid.domain.com:5000", []string***REMOVED***"invalid.domain.com:5000"***REMOVED***, false***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		config, err := makeServiceConfig(nil, tt.insecureRegistries)
		if err != nil ***REMOVED***
			t.Error(err)
		***REMOVED***
		if sec := isSecureIndex(config, tt.addr); sec != tt.expected ***REMOVED***
			t.Errorf("isSecureIndex failed for %q %v, expected %v got %v", tt.addr, tt.insecureRegistries, tt.expected, sec)
		***REMOVED***
	***REMOVED***
***REMOVED***

type debugTransport struct ***REMOVED***
	http.RoundTripper
	log func(...interface***REMOVED******REMOVED***)
***REMOVED***

func (tr debugTransport) RoundTrip(req *http.Request) (*http.Response, error) ***REMOVED***
	dump, err := httputil.DumpRequestOut(req, false)
	if err != nil ***REMOVED***
		tr.log("could not dump request")
	***REMOVED***
	tr.log(string(dump))
	resp, err := tr.RoundTripper.RoundTrip(req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	dump, err = httputil.DumpResponse(resp, false)
	if err != nil ***REMOVED***
		tr.log("could not dump response")
	***REMOVED***
	tr.log(string(dump))
	return resp, err
***REMOVED***
