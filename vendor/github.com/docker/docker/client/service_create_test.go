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
	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/api/types/swarm"
	"github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestServiceCreateError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, err := client.ServiceCreate(context.Background(), swarm.ServiceSpec***REMOVED******REMOVED***, types.ServiceCreateOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestServiceCreate(t *testing.T) ***REMOVED***
	expectedURL := "/services/create"
	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			if req.Method != "POST" ***REMOVED***
				return nil, fmt.Errorf("expected POST method, got %s", req.Method)
			***REMOVED***
			b, err := json.Marshal(types.ServiceCreateResponse***REMOVED***
				ID: "service_id",
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

	r, err := client.ServiceCreate(context.Background(), swarm.ServiceSpec***REMOVED******REMOVED***, types.ServiceCreateOptions***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if r.ID != "service_id" ***REMOVED***
		t.Fatalf("expected `service_id`, got %s", r.ID)
	***REMOVED***
***REMOVED***

func TestServiceCreateCompatiblePlatforms(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		version: "1.30",
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if strings.HasPrefix(req.URL.Path, "/v1.30/services/create") ***REMOVED***
				var serviceSpec swarm.ServiceSpec

				// check if the /distribution endpoint returned correct output
				err := json.NewDecoder(req.Body).Decode(&serviceSpec)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***

				assert.Equal(t, "foobar:1.0@sha256:c0537ff6a5218ef531ece93d4984efc99bbf3f7497c0a7726c88e2bb7584dc96", serviceSpec.TaskTemplate.ContainerSpec.Image)
				assert.Len(t, serviceSpec.TaskTemplate.Placement.Platforms, 1)

				p := serviceSpec.TaskTemplate.Placement.Platforms[0]
				b, err := json.Marshal(types.ServiceCreateResponse***REMOVED***
					ID: "service_" + p.OS + "_" + p.Architecture,
				***REMOVED***)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				return &http.Response***REMOVED***
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(b)),
				***REMOVED***, nil
			***REMOVED*** else if strings.HasPrefix(req.URL.Path, "/v1.30/distribution/") ***REMOVED***
				b, err := json.Marshal(registrytypes.DistributionInspect***REMOVED***
					Descriptor: v1.Descriptor***REMOVED***
						Digest: "sha256:c0537ff6a5218ef531ece93d4984efc99bbf3f7497c0a7726c88e2bb7584dc96",
					***REMOVED***,
					Platforms: []v1.Platform***REMOVED***
						***REMOVED***
							Architecture: "amd64",
							OS:           "linux",
						***REMOVED***,
					***REMOVED***,
				***REMOVED***)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				return &http.Response***REMOVED***
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(b)),
				***REMOVED***, nil
			***REMOVED*** else ***REMOVED***
				return nil, fmt.Errorf("unexpected URL '%s'", req.URL.Path)
			***REMOVED***
		***REMOVED***),
	***REMOVED***

	spec := swarm.ServiceSpec***REMOVED***TaskTemplate: swarm.TaskSpec***REMOVED***ContainerSpec: &swarm.ContainerSpec***REMOVED***Image: "foobar:1.0"***REMOVED******REMOVED******REMOVED***

	r, err := client.ServiceCreate(context.Background(), spec, types.ServiceCreateOptions***REMOVED***QueryRegistry: true***REMOVED***)
	assert.NoError(t, err)
	assert.Equal(t, "service_linux_amd64", r.ID)
***REMOVED***

func TestServiceCreateDigestPinning(t *testing.T) ***REMOVED***
	dgst := "sha256:c0537ff6a5218ef531ece93d4984efc99bbf3f7497c0a7726c88e2bb7584dc96"
	dgstAlt := "sha256:37ffbf3f7497c07584dc9637ffbf3f7497c0758c0537ffbf3f7497c0c88e2bb7"
	serviceCreateImage := ""
	pinByDigestTests := []struct ***REMOVED***
		img      string // input image provided by the user
		expected string // expected image after digest pinning
	***REMOVED******REMOVED***
		// default registry returns familiar string
		***REMOVED***"docker.io/library/alpine", "alpine:latest@" + dgst***REMOVED***,
		// provided tag is preserved and digest added
		***REMOVED***"alpine:edge", "alpine:edge@" + dgst***REMOVED***,
		// image with provided alternative digest remains unchanged
		***REMOVED***"alpine@" + dgstAlt, "alpine@" + dgstAlt***REMOVED***,
		// image with provided tag and alternative digest remains unchanged
		***REMOVED***"alpine:edge@" + dgstAlt, "alpine:edge@" + dgstAlt***REMOVED***,
		// image on alternative registry does not result in familiar string
		***REMOVED***"alternate.registry/library/alpine", "alternate.registry/library/alpine:latest@" + dgst***REMOVED***,
		// unresolvable image does not get a digest
		***REMOVED***"cannotresolve", "cannotresolve:latest"***REMOVED***,
	***REMOVED***

	client := &Client***REMOVED***
		version: "1.30",
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if strings.HasPrefix(req.URL.Path, "/v1.30/services/create") ***REMOVED***
				// reset and set image received by the service create endpoint
				serviceCreateImage = ""
				var service swarm.ServiceSpec
				if err := json.NewDecoder(req.Body).Decode(&service); err != nil ***REMOVED***
					return nil, fmt.Errorf("could not parse service create request")
				***REMOVED***
				serviceCreateImage = service.TaskTemplate.ContainerSpec.Image

				b, err := json.Marshal(types.ServiceCreateResponse***REMOVED***
					ID: "service_id",
				***REMOVED***)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				return &http.Response***REMOVED***
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(b)),
				***REMOVED***, nil
			***REMOVED*** else if strings.HasPrefix(req.URL.Path, "/v1.30/distribution/cannotresolve") ***REMOVED***
				// unresolvable image
				return nil, fmt.Errorf("cannot resolve image")
			***REMOVED*** else if strings.HasPrefix(req.URL.Path, "/v1.30/distribution/") ***REMOVED***
				// resolvable images
				b, err := json.Marshal(registrytypes.DistributionInspect***REMOVED***
					Descriptor: v1.Descriptor***REMOVED***
						Digest: digest.Digest(dgst),
					***REMOVED***,
				***REMOVED***)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				return &http.Response***REMOVED***
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(b)),
				***REMOVED***, nil
			***REMOVED***
			return nil, fmt.Errorf("unexpected URL '%s'", req.URL.Path)
		***REMOVED***),
	***REMOVED***

	// run pin by digest tests
	for _, p := range pinByDigestTests ***REMOVED***
		r, err := client.ServiceCreate(context.Background(), swarm.ServiceSpec***REMOVED***
			TaskTemplate: swarm.TaskSpec***REMOVED***
				ContainerSpec: &swarm.ContainerSpec***REMOVED***
					Image: p.img,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***, types.ServiceCreateOptions***REMOVED***QueryRegistry: true***REMOVED***)

		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		if r.ID != "service_id" ***REMOVED***
			t.Fatalf("expected `service_id`, got %s", r.ID)
		***REMOVED***

		if p.expected != serviceCreateImage ***REMOVED***
			t.Fatalf("expected image %s, got %s", p.expected, serviceCreateImage)
		***REMOVED***
	***REMOVED***
***REMOVED***
