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

func TestContainerCommitError(t *testing.T) ***REMOVED***
	client := &Client***REMOVED***
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	***REMOVED***
	_, err := client.ContainerCommit(context.Background(), "nothing", types.ContainerCommitOptions***REMOVED******REMOVED***)
	if err == nil || err.Error() != "Error response from daemon: Server error" ***REMOVED***
		t.Fatalf("expected a Server Error, got %v", err)
	***REMOVED***
***REMOVED***

func TestContainerCommit(t *testing.T) ***REMOVED***
	expectedURL := "/commit"
	expectedContainerID := "container_id"
	specifiedReference := "repository_name:tag"
	expectedRepositoryName := "repository_name"
	expectedTag := "tag"
	expectedComment := "comment"
	expectedAuthor := "author"
	expectedChanges := []string***REMOVED***"change1", "change2"***REMOVED***

	client := &Client***REMOVED***
		client: newMockClient(func(req *http.Request) (*http.Response, error) ***REMOVED***
			if !strings.HasPrefix(req.URL.Path, expectedURL) ***REMOVED***
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			***REMOVED***
			query := req.URL.Query()
			containerID := query.Get("container")
			if containerID != expectedContainerID ***REMOVED***
				return nil, fmt.Errorf("container id not set in URL query properly. Expected '%s', got %s", expectedContainerID, containerID)
			***REMOVED***
			repo := query.Get("repo")
			if repo != expectedRepositoryName ***REMOVED***
				return nil, fmt.Errorf("container repo not set in URL query properly. Expected '%s', got %s", expectedRepositoryName, repo)
			***REMOVED***
			tag := query.Get("tag")
			if tag != expectedTag ***REMOVED***
				return nil, fmt.Errorf("container tag not set in URL query properly. Expected '%s', got %s'", expectedTag, tag)
			***REMOVED***
			comment := query.Get("comment")
			if comment != expectedComment ***REMOVED***
				return nil, fmt.Errorf("container comment not set in URL query properly. Expected '%s', got %s'", expectedComment, comment)
			***REMOVED***
			author := query.Get("author")
			if author != expectedAuthor ***REMOVED***
				return nil, fmt.Errorf("container author not set in URL query properly. Expected '%s', got %s'", expectedAuthor, author)
			***REMOVED***
			pause := query.Get("pause")
			if pause != "0" ***REMOVED***
				return nil, fmt.Errorf("container pause not set in URL query properly. Expected 'true', got %v'", pause)
			***REMOVED***
			changes := query["changes"]
			if len(changes) != len(expectedChanges) ***REMOVED***
				return nil, fmt.Errorf("expected container changes size to be '%d', got %d", len(expectedChanges), len(changes))
			***REMOVED***
			b, err := json.Marshal(types.IDResponse***REMOVED***
				ID: "new_container_id",
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

	r, err := client.ContainerCommit(context.Background(), expectedContainerID, types.ContainerCommitOptions***REMOVED***
		Reference: specifiedReference,
		Comment:   expectedComment,
		Author:    expectedAuthor,
		Changes:   expectedChanges,
		Pause:     false,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if r.ID != "new_container_id" ***REMOVED***
		t.Fatalf("expected `new_container_id`, got %s", r.ID)
	***REMOVED***
***REMOVED***
