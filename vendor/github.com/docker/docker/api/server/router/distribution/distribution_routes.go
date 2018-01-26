package distribution

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/docker/distribution/manifest/manifestlist"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/types"
	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func (s *distributionRouter) getDistributionInfo(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error ***REMOVED***
	if err := httputils.ParseForm(r); err != nil ***REMOVED***
		return err
	***REMOVED***

	w.Header().Set("Content-Type", "application/json")

	var (
		config              = &types.AuthConfig***REMOVED******REMOVED***
		authEncoded         = r.Header.Get("X-Registry-Auth")
		distributionInspect registrytypes.DistributionInspect
	)

	if authEncoded != "" ***REMOVED***
		authJSON := base64.NewDecoder(base64.URLEncoding, strings.NewReader(authEncoded))
		if err := json.NewDecoder(authJSON).Decode(&config); err != nil ***REMOVED***
			// for a search it is not an error if no auth was given
			// to increase compatibility with the existing api it is defaulting to be empty
			config = &types.AuthConfig***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***

	image := vars["name"]

	ref, err := reference.ParseAnyReference(image)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	namedRef, ok := ref.(reference.Named)
	if !ok ***REMOVED***
		if _, ok := ref.(reference.Digested); ok ***REMOVED***
			// full image ID
			return errors.Errorf("no manifest found for full image ID")
		***REMOVED***
		return errors.Errorf("unknown image reference format: %s", image)
	***REMOVED***

	distrepo, _, err := s.backend.GetRepository(ctx, namedRef, config)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	blobsrvc := distrepo.Blobs(ctx)

	if canonicalRef, ok := namedRef.(reference.Canonical); !ok ***REMOVED***
		namedRef = reference.TagNameOnly(namedRef)

		taggedRef, ok := namedRef.(reference.NamedTagged)
		if !ok ***REMOVED***
			return errors.Errorf("image reference not tagged: %s", image)
		***REMOVED***

		descriptor, err := distrepo.Tags(ctx).Get(ctx, taggedRef.Tag())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		distributionInspect.Descriptor = v1.Descriptor***REMOVED***
			MediaType: descriptor.MediaType,
			Digest:    descriptor.Digest,
			Size:      descriptor.Size,
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// TODO(nishanttotla): Once manifests can be looked up as a blob, the
		// descriptor should be set using blobsrvc.Stat(ctx, canonicalRef.Digest())
		// instead of having to manually fill in the fields
		distributionInspect.Descriptor.Digest = canonicalRef.Digest()
	***REMOVED***

	// we have a digest, so we can retrieve the manifest
	mnfstsrvc, err := distrepo.Manifests(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	mnfst, err := mnfstsrvc.Get(ctx, distributionInspect.Descriptor.Digest)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	mediaType, payload, err := mnfst.Payload()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// update MediaType because registry might return something incorrect
	distributionInspect.Descriptor.MediaType = mediaType
	if distributionInspect.Descriptor.Size == 0 ***REMOVED***
		distributionInspect.Descriptor.Size = int64(len(payload))
	***REMOVED***

	// retrieve platform information depending on the type of manifest
	switch mnfstObj := mnfst.(type) ***REMOVED***
	case *manifestlist.DeserializedManifestList:
		for _, m := range mnfstObj.Manifests ***REMOVED***
			distributionInspect.Platforms = append(distributionInspect.Platforms, v1.Platform***REMOVED***
				Architecture: m.Platform.Architecture,
				OS:           m.Platform.OS,
				OSVersion:    m.Platform.OSVersion,
				OSFeatures:   m.Platform.OSFeatures,
				Variant:      m.Platform.Variant,
			***REMOVED***)
		***REMOVED***
	case *schema2.DeserializedManifest:
		configJSON, err := blobsrvc.Get(ctx, mnfstObj.Config.Digest)
		var platform v1.Platform
		if err == nil ***REMOVED***
			err := json.Unmarshal(configJSON, &platform)
			if err == nil && (platform.OS != "" || platform.Architecture != "") ***REMOVED***
				distributionInspect.Platforms = append(distributionInspect.Platforms, platform)
			***REMOVED***
		***REMOVED***
	case *schema1.SignedManifest:
		platform := v1.Platform***REMOVED***
			Architecture: mnfstObj.Architecture,
			OS:           "linux",
		***REMOVED***
		distributionInspect.Platforms = append(distributionInspect.Platforms, platform)
	***REMOVED***

	return httputils.WriteJSON(w, http.StatusOK, distributionInspect)
***REMOVED***
