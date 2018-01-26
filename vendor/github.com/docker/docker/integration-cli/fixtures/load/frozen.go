package load

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"github.com/pkg/errors"
)

const frozenImgDir = "/docker-frozen-images"

// FrozenImagesLinux loads the frozen image set for the integration suite
// If the images are not available locally it will download them
// TODO: This loads whatever is in the frozen image dir, regardless of what
// images were passed in. If the images need to be downloaded, then it will respect
// the passed in images
func FrozenImagesLinux(client client.APIClient, images ...string) error ***REMOVED***
	var loadImages []struct***REMOVED*** srcName, destName string ***REMOVED***
	for _, img := range images ***REMOVED***
		if !imageExists(client, img) ***REMOVED***
			srcName := img
			// hello-world:latest gets re-tagged as hello-world:frozen
			// there are some tests that use hello-world:latest specifically so it pulls
			// the image and hello-world:frozen is used for when we just want a super
			// small image
			if img == "hello-world:frozen" ***REMOVED***
				srcName = "hello-world:latest"
			***REMOVED***
			loadImages = append(loadImages, struct***REMOVED*** srcName, destName string ***REMOVED******REMOVED***
				srcName:  srcName,
				destName: img,
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	if len(loadImages) == 0 ***REMOVED***
		// everything is loaded, we're done
		return nil
	***REMOVED***

	ctx := context.Background()
	fi, err := os.Stat(frozenImgDir)
	if err != nil || !fi.IsDir() ***REMOVED***
		srcImages := make([]string, 0, len(loadImages))
		for _, img := range loadImages ***REMOVED***
			srcImages = append(srcImages, img.srcName)
		***REMOVED***
		if err := pullImages(ctx, client, srcImages); err != nil ***REMOVED***
			return errors.Wrap(err, "error pulling image list")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if err := loadFrozenImages(ctx, client); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	for _, img := range loadImages ***REMOVED***
		if img.srcName != img.destName ***REMOVED***
			if err := client.ImageTag(ctx, img.srcName, img.destName); err != nil ***REMOVED***
				return errors.Wrapf(err, "failed to tag %s as %s", img.srcName, img.destName)
			***REMOVED***
			if _, err := client.ImageRemove(ctx, img.srcName, types.ImageRemoveOptions***REMOVED******REMOVED***); err != nil ***REMOVED***
				return errors.Wrapf(err, "failed to remove %s", img.srcName)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func imageExists(client client.APIClient, name string) bool ***REMOVED***
	_, _, err := client.ImageInspectWithRaw(context.Background(), name)
	return err == nil
***REMOVED***

func loadFrozenImages(ctx context.Context, client client.APIClient) error ***REMOVED***
	tar, err := exec.LookPath("tar")
	if err != nil ***REMOVED***
		return errors.Wrap(err, "could not find tar binary")
	***REMOVED***
	tarCmd := exec.Command(tar, "-cC", frozenImgDir, ".")
	out, err := tarCmd.StdoutPipe()
	if err != nil ***REMOVED***
		return errors.Wrap(err, "error getting stdout pipe for tar command")
	***REMOVED***

	errBuf := bytes.NewBuffer(nil)
	tarCmd.Stderr = errBuf
	tarCmd.Start()
	defer tarCmd.Wait()

	resp, err := client.ImageLoad(ctx, out, true)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "failed to load frozen images")
	***REMOVED***
	defer resp.Body.Close()
	fd, isTerminal := term.GetFdInfo(os.Stdout)
	return jsonmessage.DisplayJSONMessagesStream(resp.Body, os.Stdout, fd, isTerminal, nil)
***REMOVED***

func pullImages(ctx context.Context, client client.APIClient, images []string) error ***REMOVED***
	cwd, err := os.Getwd()
	if err != nil ***REMOVED***
		return errors.Wrap(err, "error getting path to dockerfile")
	***REMOVED***
	dockerfile := os.Getenv("DOCKERFILE")
	if dockerfile == "" ***REMOVED***
		dockerfile = "Dockerfile"
	***REMOVED***
	dockerfilePath := filepath.Join(filepath.Dir(filepath.Clean(cwd)), dockerfile)
	pullRefs, err := readFrozenImageList(dockerfilePath, images)
	if err != nil ***REMOVED***
		return errors.Wrap(err, "error reading frozen image list")
	***REMOVED***

	var wg sync.WaitGroup
	chErr := make(chan error, len(images))
	for tag, ref := range pullRefs ***REMOVED***
		wg.Add(1)
		go func(tag, ref string) ***REMOVED***
			defer wg.Done()
			if err := pullTagAndRemove(ctx, client, ref, tag); err != nil ***REMOVED***
				chErr <- err
				return
			***REMOVED***
		***REMOVED***(tag, ref)
	***REMOVED***
	wg.Wait()
	close(chErr)
	return <-chErr
***REMOVED***

func pullTagAndRemove(ctx context.Context, client client.APIClient, ref string, tag string) error ***REMOVED***
	resp, err := client.ImagePull(ctx, ref, types.ImagePullOptions***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to pull %s", ref)
	***REMOVED***
	defer resp.Close()
	fd, isTerminal := term.GetFdInfo(os.Stdout)
	if err := jsonmessage.DisplayJSONMessagesStream(resp, os.Stdout, fd, isTerminal, nil); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := client.ImageTag(ctx, ref, tag); err != nil ***REMOVED***
		return errors.Wrapf(err, "failed to tag %s as %s", ref, tag)
	***REMOVED***
	_, err = client.ImageRemove(ctx, ref, types.ImageRemoveOptions***REMOVED******REMOVED***)
	return errors.Wrapf(err, "failed to remove %s", ref)

***REMOVED***

func readFrozenImageList(dockerfilePath string, images []string) (map[string]string, error) ***REMOVED***
	f, err := os.Open(dockerfilePath)
	if err != nil ***REMOVED***
		return nil, errors.Wrap(err, "error reading dockerfile")
	***REMOVED***
	defer f.Close()
	ls := make(map[string]string)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() ***REMOVED***
		line := strings.Fields(scanner.Text())
		if len(line) < 3 ***REMOVED***
			continue
		***REMOVED***
		if !(line[0] == "RUN" && line[1] == "./contrib/download-frozen-image-v2.sh") ***REMOVED***
			continue
		***REMOVED***

		for scanner.Scan() ***REMOVED***
			img := strings.TrimSpace(scanner.Text())
			img = strings.TrimSuffix(img, "\\")
			img = strings.TrimSpace(img)
			split := strings.Split(img, "@")
			if len(split) < 2 ***REMOVED***
				break
			***REMOVED***

			for _, i := range images ***REMOVED***
				if split[0] == i ***REMOVED***
					ls[i] = img
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return ls, nil
***REMOVED***
