// +build !windows

package daemon

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	containertypes "github.com/docker/docker/api/types/container"
	mounttypes "github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/container"
	"github.com/docker/docker/volume"
)

func TestBackportMountSpec(t *testing.T) ***REMOVED***
	d := Daemon***REMOVED***containers: container.NewMemoryStore()***REMOVED***

	c := &container.Container***REMOVED***
		State: &container.State***REMOVED******REMOVED***,
		MountPoints: map[string]*volume.MountPoint***REMOVED***
			"/apple":      ***REMOVED***Destination: "/apple", Source: "/var/lib/docker/volumes/12345678", Name: "12345678", RW: true, CopyData: true***REMOVED***, // anonymous volume
			"/banana":     ***REMOVED***Destination: "/banana", Source: "/var/lib/docker/volumes/data", Name: "data", RW: true, CopyData: true***REMOVED***,        // named volume
			"/cherry":     ***REMOVED***Destination: "/cherry", Source: "/var/lib/docker/volumes/data", Name: "data", CopyData: true***REMOVED***,                  // RO named volume
			"/dates":      ***REMOVED***Destination: "/dates", Source: "/var/lib/docker/volumes/data", Name: "data"***REMOVED***,                                   // named volume nocopy
			"/elderberry": ***REMOVED***Destination: "/elderberry", Source: "/var/lib/docker/volumes/data", Name: "data"***REMOVED***,                              // masks anon vol
			"/fig":        ***REMOVED***Destination: "/fig", Source: "/data", RW: true***REMOVED***,                                                                // RW bind
			"/guava":      ***REMOVED***Destination: "/guava", Source: "/data", RW: false, Propagation: "shared"***REMOVED***,                                      // RO bind + propagation
			"/kumquat":    ***REMOVED***Destination: "/kumquat", Name: "data", RW: false, CopyData: true***REMOVED***,                                              // volumes-from

			// partially configured mountpoint due to #32613
			// specifically, `mp.Spec.Source` is not set
			"/honeydew": ***REMOVED***
				Type:        mounttypes.TypeVolume,
				Destination: "/honeydew",
				Name:        "data",
				Source:      "/var/lib/docker/volumes/data",
				Spec:        mounttypes.Mount***REMOVED***Type: mounttypes.TypeVolume, Target: "/honeydew", VolumeOptions: &mounttypes.VolumeOptions***REMOVED***NoCopy: true***REMOVED******REMOVED***,
			***REMOVED***,

			// from hostconfig.Mounts
			"/jambolan": ***REMOVED***
				Type:        mounttypes.TypeVolume,
				Destination: "/jambolan",
				Source:      "/var/lib/docker/volumes/data",
				RW:          true,
				Name:        "data",
				Spec:        mounttypes.Mount***REMOVED***Type: mounttypes.TypeVolume, Target: "/jambolan", Source: "data"***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		HostConfig: &containertypes.HostConfig***REMOVED***
			Binds: []string***REMOVED***
				"data:/banana",
				"data:/cherry:ro",
				"data:/dates:ro,nocopy",
				"data:/elderberry:ro,nocopy",
				"/data:/fig",
				"/data:/guava:ro,shared",
				"data:/honeydew:nocopy",
			***REMOVED***,
			VolumesFrom: []string***REMOVED***"1:ro"***REMOVED***,
			Mounts: []mounttypes.Mount***REMOVED***
				***REMOVED***Type: mounttypes.TypeVolume, Target: "/jambolan"***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		Config: &containertypes.Config***REMOVED***Volumes: map[string]struct***REMOVED******REMOVED******REMOVED***
			"/apple":      ***REMOVED******REMOVED***,
			"/elderberry": ***REMOVED******REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***

	d.containers.Add("1", &container.Container***REMOVED***
		State: &container.State***REMOVED******REMOVED***,
		ID:    "1",
		MountPoints: map[string]*volume.MountPoint***REMOVED***
			"/kumquat": ***REMOVED***Destination: "/kumquat", Name: "data", RW: false, CopyData: true***REMOVED***,
		***REMOVED***,
		HostConfig: &containertypes.HostConfig***REMOVED***
			Binds: []string***REMOVED***
				"data:/kumquat:ro",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***)

	type expected struct ***REMOVED***
		mp      *volume.MountPoint
		comment string
	***REMOVED***

	pretty := func(mp *volume.MountPoint) string ***REMOVED***
		b, err := json.MarshalIndent(mp, "\t", "    ")
		if err != nil ***REMOVED***
			return fmt.Sprintf("%#v", mp)
		***REMOVED***
		return string(b)
	***REMOVED***

	for _, x := range []expected***REMOVED***
		***REMOVED***
			mp: &volume.MountPoint***REMOVED***
				Type:        mounttypes.TypeVolume,
				Destination: "/apple",
				RW:          true,
				Name:        "12345678",
				Source:      "/var/lib/docker/volumes/12345678",
				CopyData:    true,
				Spec: mounttypes.Mount***REMOVED***
					Type:   mounttypes.TypeVolume,
					Source: "",
					Target: "/apple",
				***REMOVED***,
			***REMOVED***,
			comment: "anonymous volume",
		***REMOVED***,
		***REMOVED***
			mp: &volume.MountPoint***REMOVED***
				Type:        mounttypes.TypeVolume,
				Destination: "/banana",
				RW:          true,
				Name:        "data",
				Source:      "/var/lib/docker/volumes/data",
				CopyData:    true,
				Spec: mounttypes.Mount***REMOVED***
					Type:   mounttypes.TypeVolume,
					Source: "data",
					Target: "/banana",
				***REMOVED***,
			***REMOVED***,
			comment: "named volume",
		***REMOVED***,
		***REMOVED***
			mp: &volume.MountPoint***REMOVED***
				Type:        mounttypes.TypeVolume,
				Destination: "/cherry",
				Name:        "data",
				Source:      "/var/lib/docker/volumes/data",
				CopyData:    true,
				Spec: mounttypes.Mount***REMOVED***
					Type:     mounttypes.TypeVolume,
					Source:   "data",
					Target:   "/cherry",
					ReadOnly: true,
				***REMOVED***,
			***REMOVED***,
			comment: "read-only named volume",
		***REMOVED***,
		***REMOVED***
			mp: &volume.MountPoint***REMOVED***
				Type:        mounttypes.TypeVolume,
				Destination: "/dates",
				Name:        "data",
				Source:      "/var/lib/docker/volumes/data",
				Spec: mounttypes.Mount***REMOVED***
					Type:          mounttypes.TypeVolume,
					Source:        "data",
					Target:        "/dates",
					ReadOnly:      true,
					VolumeOptions: &mounttypes.VolumeOptions***REMOVED***NoCopy: true***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			comment: "named volume with nocopy",
		***REMOVED***,
		***REMOVED***
			mp: &volume.MountPoint***REMOVED***
				Type:        mounttypes.TypeVolume,
				Destination: "/elderberry",
				Name:        "data",
				Source:      "/var/lib/docker/volumes/data",
				Spec: mounttypes.Mount***REMOVED***
					Type:          mounttypes.TypeVolume,
					Source:        "data",
					Target:        "/elderberry",
					ReadOnly:      true,
					VolumeOptions: &mounttypes.VolumeOptions***REMOVED***NoCopy: true***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			comment: "masks an anonymous volume",
		***REMOVED***,
		***REMOVED***
			mp: &volume.MountPoint***REMOVED***
				Type:        mounttypes.TypeBind,
				Destination: "/fig",
				Source:      "/data",
				RW:          true,
				Spec: mounttypes.Mount***REMOVED***
					Type:   mounttypes.TypeBind,
					Source: "/data",
					Target: "/fig",
				***REMOVED***,
			***REMOVED***,
			comment: "bind mount with read/write",
		***REMOVED***,
		***REMOVED***
			mp: &volume.MountPoint***REMOVED***
				Type:        mounttypes.TypeBind,
				Destination: "/guava",
				Source:      "/data",
				RW:          false,
				Propagation: "shared",
				Spec: mounttypes.Mount***REMOVED***
					Type:        mounttypes.TypeBind,
					Source:      "/data",
					Target:      "/guava",
					ReadOnly:    true,
					BindOptions: &mounttypes.BindOptions***REMOVED***Propagation: "shared"***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			comment: "bind mount with read/write + shared propagation",
		***REMOVED***,
		***REMOVED***
			mp: &volume.MountPoint***REMOVED***
				Type:        mounttypes.TypeVolume,
				Destination: "/honeydew",
				Source:      "/var/lib/docker/volumes/data",
				RW:          true,
				Propagation: "shared",
				Spec: mounttypes.Mount***REMOVED***
					Type:          mounttypes.TypeVolume,
					Source:        "data",
					Target:        "/honeydew",
					VolumeOptions: &mounttypes.VolumeOptions***REMOVED***NoCopy: true***REMOVED***,
				***REMOVED***,
			***REMOVED***,
			comment: "partially configured named volume caused by #32613",
		***REMOVED***,
		***REMOVED***
			mp:      &(*c.MountPoints["/jambolan"]), // copy the mountpoint, expect no changes
			comment: "volume defined in mounts API",
		***REMOVED***,
		***REMOVED***
			mp: &volume.MountPoint***REMOVED***
				Type:        mounttypes.TypeVolume,
				Destination: "/kumquat",
				Source:      "/var/lib/docker/volumes/data",
				RW:          false,
				Name:        "data",
				Spec: mounttypes.Mount***REMOVED***
					Type:     mounttypes.TypeVolume,
					Source:   "data",
					Target:   "/kumquat",
					ReadOnly: true,
				***REMOVED***,
			***REMOVED***,
			comment: "partially configured named volume caused by #32613",
		***REMOVED***,
	***REMOVED*** ***REMOVED***

		mp := c.MountPoints[x.mp.Destination]
		d.backportMountSpec(c)

		if !reflect.DeepEqual(mp.Spec, x.mp.Spec) ***REMOVED***
			t.Fatalf("%s\nexpected:\n\t%s\n\ngot:\n\t%s", x.comment, pretty(x.mp), pretty(mp))
		***REMOVED***
	***REMOVED***
***REMOVED***
