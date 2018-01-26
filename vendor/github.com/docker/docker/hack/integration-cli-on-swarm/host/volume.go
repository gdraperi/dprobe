package main

import (
	"archive/tar"
	"bytes"
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

func createTar(data map[string][]byte) (io.Reader, error) ***REMOVED***
	var b bytes.Buffer
	tw := tar.NewWriter(&b)
	for path, datum := range data ***REMOVED***
		hdr := tar.Header***REMOVED***
			Name: path,
			Mode: 0644,
			Size: int64(len(datum)),
		***REMOVED***
		if err := tw.WriteHeader(&hdr); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		_, err := tw.Write(datum)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	if err := tw.Close(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &b, nil
***REMOVED***

// createVolumeWithData creates a volume with the given data (e.g. data["/foo"] = []byte("bar"))
// Internally, a container is created from the image so as to provision the data to the volume,
// which is attached to the container.
func createVolumeWithData(cli *client.Client, volumeName string, data map[string][]byte, image string) error ***REMOVED***
	_, err := cli.VolumeCreate(context.Background(),
		volume.VolumesCreateBody***REMOVED***
			Driver: "local",
			Name:   volumeName,
		***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	mnt := "/mnt"
	miniContainer, err := cli.ContainerCreate(context.Background(),
		&container.Config***REMOVED***
			Image: image,
		***REMOVED***,
		&container.HostConfig***REMOVED***
			Mounts: []mount.Mount***REMOVED***
				***REMOVED***
					Type:   mount.TypeVolume,
					Source: volumeName,
					Target: mnt,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***, nil, "")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	tr, err := createTar(data)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if cli.CopyToContainer(context.Background(),
		miniContainer.ID, mnt, tr, types.CopyToContainerOptions***REMOVED******REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***
	return cli.ContainerRemove(context.Background(),
		miniContainer.ID,
		types.ContainerRemoveOptions***REMOVED******REMOVED***)
***REMOVED***

func hasVolume(cli *client.Client, volumeName string) bool ***REMOVED***
	_, err := cli.VolumeInspect(context.Background(), volumeName)
	return err == nil
***REMOVED***

func removeVolume(cli *client.Client, volumeName string) error ***REMOVED***
	return cli.VolumeRemove(context.Background(), volumeName, true)
***REMOVED***
