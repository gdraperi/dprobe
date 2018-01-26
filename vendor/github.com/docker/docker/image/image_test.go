package image

import (
	"encoding/json"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/layer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const sampleImageJSON = `***REMOVED***
	"architecture": "amd64",
	"os": "linux",
	"config": ***REMOVED******REMOVED***,
	"rootfs": ***REMOVED***
		"type": "layers",
		"diff_ids": []
	***REMOVED***
***REMOVED***`

func TestNewFromJSON(t *testing.T) ***REMOVED***
	img, err := NewFromJSON([]byte(sampleImageJSON))
	require.NoError(t, err)
	assert.Equal(t, sampleImageJSON, string(img.RawJSON()))
***REMOVED***

func TestNewFromJSONWithInvalidJSON(t *testing.T) ***REMOVED***
	_, err := NewFromJSON([]byte("***REMOVED******REMOVED***"))
	assert.EqualError(t, err, "invalid image JSON, no RootFS key")
***REMOVED***

func TestMarshalKeyOrder(t *testing.T) ***REMOVED***
	b, err := json.Marshal(&Image***REMOVED***
		V1Image: V1Image***REMOVED***
			Comment:      "a",
			Author:       "b",
			Architecture: "c",
		***REMOVED***,
	***REMOVED***)
	assert.NoError(t, err)

	expectedOrder := []string***REMOVED***"architecture", "author", "comment"***REMOVED***
	var indexes []int
	for _, k := range expectedOrder ***REMOVED***
		indexes = append(indexes, strings.Index(string(b), k))
	***REMOVED***

	if !sort.IntsAreSorted(indexes) ***REMOVED***
		t.Fatal("invalid key order in JSON: ", string(b))
	***REMOVED***
***REMOVED***

func TestImage(t *testing.T) ***REMOVED***
	cid := "50a16564e727"
	config := &container.Config***REMOVED***
		Hostname:   "hostname",
		Domainname: "domain",
		User:       "root",
	***REMOVED***
	os := runtime.GOOS

	img := &Image***REMOVED***
		V1Image: V1Image***REMOVED***
			Config: config,
		***REMOVED***,
		computedID: ID(cid),
	***REMOVED***

	assert.Equal(t, cid, img.ImageID())
	assert.Equal(t, cid, img.ID().String())
	assert.Equal(t, os, img.OperatingSystem())
	assert.Equal(t, config, img.RunConfig())
***REMOVED***

func TestImageOSNotEmpty(t *testing.T) ***REMOVED***
	os := "os"
	img := &Image***REMOVED***
		V1Image: V1Image***REMOVED***
			OS: os,
		***REMOVED***,
		OSVersion: "osversion",
	***REMOVED***
	assert.Equal(t, os, img.OperatingSystem())
***REMOVED***

func TestNewChildImageFromImageWithRootFS(t *testing.T) ***REMOVED***
	rootFS := NewRootFS()
	rootFS.Append(layer.DiffID("ba5e"))
	parent := &Image***REMOVED***
		RootFS: rootFS,
		History: []History***REMOVED***
			NewHistory("a", "c", "r", false),
		***REMOVED***,
	***REMOVED***
	childConfig := ChildConfig***REMOVED***
		DiffID:  layer.DiffID("abcdef"),
		Author:  "author",
		Comment: "comment",
		ContainerConfig: &container.Config***REMOVED***
			Cmd: []string***REMOVED***"echo", "foo"***REMOVED***,
		***REMOVED***,
		Config: &container.Config***REMOVED******REMOVED***,
	***REMOVED***

	newImage := NewChildImage(parent, childConfig, "platform")
	expectedDiffIDs := []layer.DiffID***REMOVED***layer.DiffID("ba5e"), layer.DiffID("abcdef")***REMOVED***
	assert.Equal(t, expectedDiffIDs, newImage.RootFS.DiffIDs)
	assert.Equal(t, childConfig.Author, newImage.Author)
	assert.Equal(t, childConfig.Config, newImage.Config)
	assert.Equal(t, *childConfig.ContainerConfig, newImage.ContainerConfig)
	assert.Equal(t, "platform", newImage.OS)
	assert.Equal(t, childConfig.Config, newImage.Config)

	assert.Len(t, newImage.History, 2)
	assert.Equal(t, childConfig.Comment, newImage.History[1].Comment)

	// RootFS should be copied not mutated
	assert.NotEqual(t, parent.RootFS.DiffIDs, newImage.RootFS.DiffIDs)
***REMOVED***
