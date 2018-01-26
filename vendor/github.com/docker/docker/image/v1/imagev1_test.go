package v1

import (
	"encoding/json"
	"testing"

	"github.com/docker/docker/image"
)

func TestMakeV1ConfigFromConfig(t *testing.T) ***REMOVED***
	img := &image.Image***REMOVED***
		V1Image: image.V1Image***REMOVED***
			ID:     "v2id",
			Parent: "v2parent",
			OS:     "os",
		***REMOVED***,
		OSVersion: "osversion",
		RootFS: &image.RootFS***REMOVED***
			Type: "layers",
		***REMOVED***,
	***REMOVED***
	v2js, err := json.Marshal(img)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// Convert the image back in order to get RawJSON() support.
	img, err = image.NewFromJSON(v2js)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	js, err := MakeV1ConfigFromConfig(img, "v1id", "v1parent", false)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	newimg := &image.Image***REMOVED******REMOVED***
	err = json.Unmarshal(js, newimg)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if newimg.V1Image.ID != "v1id" || newimg.Parent != "v1parent" ***REMOVED***
		t.Error("ids should have changed", newimg.V1Image.ID, newimg.V1Image.Parent)
	***REMOVED***

	if newimg.RootFS != nil ***REMOVED***
		t.Error("rootfs should have been removed")
	***REMOVED***

	if newimg.V1Image.OS != "os" ***REMOVED***
		t.Error("os should have been preserved")
	***REMOVED***
***REMOVED***
