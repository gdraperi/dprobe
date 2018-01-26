package distribution

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/docker/distribution"
	"github.com/docker/distribution/context"
	"github.com/docker/distribution/manifest/manifestlist"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/registry/client/transport"
	"github.com/docker/docker/pkg/system"
	"github.com/sirupsen/logrus"
)

var _ distribution.Describable = &v2LayerDescriptor***REMOVED******REMOVED***

func (ld *v2LayerDescriptor) Descriptor() distribution.Descriptor ***REMOVED***
	if ld.src.MediaType == schema2.MediaTypeForeignLayer && len(ld.src.URLs) > 0 ***REMOVED***
		return ld.src
	***REMOVED***
	return distribution.Descriptor***REMOVED******REMOVED***
***REMOVED***

func (ld *v2LayerDescriptor) open(ctx context.Context) (distribution.ReadSeekCloser, error) ***REMOVED***
	blobs := ld.repo.Blobs(ctx)
	rsc, err := blobs.Open(ctx, ld.digest)

	if len(ld.src.URLs) == 0 ***REMOVED***
		return rsc, err
	***REMOVED***

	// We're done if the registry has this blob.
	if err == nil ***REMOVED***
		// Seek does an HTTP GET.  If it succeeds, the blob really is accessible.
		if _, err = rsc.Seek(0, os.SEEK_SET); err == nil ***REMOVED***
			return rsc, nil
		***REMOVED***
		rsc.Close()
	***REMOVED***

	// Find the first URL that results in a 200 result code.
	for _, url := range ld.src.URLs ***REMOVED***
		logrus.Debugf("Pulling %v from foreign URL %v", ld.digest, url)
		rsc = transport.NewHTTPReadSeeker(http.DefaultClient, url, nil)

		// Seek does an HTTP GET.  If it succeeds, the blob really is accessible.
		_, err = rsc.Seek(0, os.SEEK_SET)
		if err == nil ***REMOVED***
			break
		***REMOVED***
		logrus.Debugf("Download for %v failed: %v", ld.digest, err)
		rsc.Close()
		rsc = nil
	***REMOVED***
	return rsc, err
***REMOVED***

func filterManifests(manifests []manifestlist.ManifestDescriptor, os string) []manifestlist.ManifestDescriptor ***REMOVED***
	osVersion := ""
	if os == "windows" ***REMOVED***
		// TODO: Add UBR (Update Build Release) component after build
		version := system.GetOSVersion()
		osVersion = fmt.Sprintf("%d.%d.%d", version.MajorVersion, version.MinorVersion, version.Build)
		logrus.Debugf("will prefer entries with version %s", osVersion)
	***REMOVED***

	var matches []manifestlist.ManifestDescriptor
	for _, manifestDescriptor := range manifests ***REMOVED***
		// TODO: Consider filtering out greater versions, including only greater UBR
		if manifestDescriptor.Platform.Architecture == runtime.GOARCH && manifestDescriptor.Platform.OS == os ***REMOVED***
			matches = append(matches, manifestDescriptor)
			logrus.Debugf("found match for %s/%s with media type %s, digest %s", os, runtime.GOARCH, manifestDescriptor.MediaType, manifestDescriptor.Digest.String())
		***REMOVED***
	***REMOVED***
	if os == "windows" ***REMOVED***
		sort.Stable(manifestsByVersion***REMOVED***osVersion, matches***REMOVED***)
	***REMOVED***
	return matches
***REMOVED***

func versionMatch(actual, expected string) bool ***REMOVED***
	// Check whether the version matches up to the build, ignoring UBR
	return strings.HasPrefix(actual, expected+".")
***REMOVED***

type manifestsByVersion struct ***REMOVED***
	version string
	list    []manifestlist.ManifestDescriptor
***REMOVED***

func (mbv manifestsByVersion) Less(i, j int) bool ***REMOVED***
	// TODO: Split version by parts and compare
	// TODO: Prefer versions which have a greater version number
	// Move compatible versions to the top, with no other ordering changes
	return versionMatch(mbv.list[i].Platform.OSVersion, mbv.version) && !versionMatch(mbv.list[j].Platform.OSVersion, mbv.version)
***REMOVED***

func (mbv manifestsByVersion) Len() int ***REMOVED***
	return len(mbv.list)
***REMOVED***

func (mbv manifestsByVersion) Swap(i, j int) ***REMOVED***
	mbv.list[i], mbv.list[j] = mbv.list[j], mbv.list[i]
***REMOVED***
