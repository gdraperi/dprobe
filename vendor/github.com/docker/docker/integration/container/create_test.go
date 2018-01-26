package container

import (
	"context"
	"strconv"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/integration/util/request"
	"github.com/docker/docker/internal/testutil"
)

func TestCreateFailsWhenIdentifierDoesNotExist(t *testing.T) ***REMOVED***
	defer setupTest(t)()
	client := request.NewAPIClient(t)

	testCases := []struct ***REMOVED***
		doc           string
		image         string
		expectedError string
	***REMOVED******REMOVED***
		***REMOVED***
			doc:           "image and tag",
			image:         "test456:v1",
			expectedError: "No such image: test456:v1",
		***REMOVED***,
		***REMOVED***
			doc:           "image no tag",
			image:         "test456",
			expectedError: "No such image: test456",
		***REMOVED***,
		***REMOVED***
			doc:           "digest",
			image:         "sha256:0cb40641836c461bc97c793971d84d758371ed682042457523e4ae701efeaaaa",
			expectedError: "No such image: sha256:0cb40641836c461bc97c793971d84d758371ed682042457523e4ae701efeaaaa",
		***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(tc.doc, func(t *testing.T) ***REMOVED***
			t.Parallel()
			_, err := client.ContainerCreate(context.Background(),
				&container.Config***REMOVED***Image: tc.image***REMOVED***,
				&container.HostConfig***REMOVED******REMOVED***,
				&network.NetworkingConfig***REMOVED******REMOVED***,
				"foo",
			)
			testutil.ErrorContains(t, err, tc.expectedError)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestCreateWithInvalidEnv(t *testing.T) ***REMOVED***
	defer setupTest(t)()
	client := request.NewAPIClient(t)

	testCases := []struct ***REMOVED***
		env           string
		expectedError string
	***REMOVED******REMOVED***
		***REMOVED***
			env:           "",
			expectedError: "invalid environment variable:",
		***REMOVED***,
		***REMOVED***
			env:           "=",
			expectedError: "invalid environment variable: =",
		***REMOVED***,
		***REMOVED***
			env:           "=foo",
			expectedError: "invalid environment variable: =foo",
		***REMOVED***,
	***REMOVED***

	for index, tc := range testCases ***REMOVED***
		tc := tc
		t.Run(strconv.Itoa(index), func(t *testing.T) ***REMOVED***
			t.Parallel()
			_, err := client.ContainerCreate(context.Background(),
				&container.Config***REMOVED***
					Image: "busybox",
					Env:   []string***REMOVED***tc.env***REMOVED***,
				***REMOVED***,
				&container.HostConfig***REMOVED******REMOVED***,
				&network.NetworkingConfig***REMOVED******REMOVED***,
				"foo",
			)
			testutil.ErrorContains(t, err, tc.expectedError)
		***REMOVED***)
	***REMOVED***
***REMOVED***
