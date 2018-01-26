package tarexport

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/docker/distribution"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/image"
	"github.com/docker/docker/image/v1"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/system"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

type imageDescriptor struct ***REMOVED***
	refs     []reference.NamedTagged
	layers   []string
	image    *image.Image
	layerRef layer.Layer
***REMOVED***

type saveSession struct ***REMOVED***
	*tarexporter
	outDir      string
	images      map[image.ID]*imageDescriptor
	savedLayers map[string]struct***REMOVED******REMOVED***
	diffIDPaths map[layer.DiffID]string // cache every diffID blob to avoid duplicates
***REMOVED***

func (l *tarexporter) Save(names []string, outStream io.Writer) error ***REMOVED***
	images, err := l.parseNames(names)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Release all the image top layer references
	defer l.releaseLayerReferences(images)
	return (&saveSession***REMOVED***tarexporter: l, images: images***REMOVED***).save(outStream)
***REMOVED***

// parseNames will parse the image names to a map which contains image.ID to *imageDescriptor.
// Each imageDescriptor holds an image top layer reference named 'layerRef'. It is taken here, should be released later.
func (l *tarexporter) parseNames(names []string) (desc map[image.ID]*imageDescriptor, rErr error) ***REMOVED***
	imgDescr := make(map[image.ID]*imageDescriptor)
	defer func() ***REMOVED***
		if rErr != nil ***REMOVED***
			l.releaseLayerReferences(imgDescr)
		***REMOVED***
	***REMOVED***()

	addAssoc := func(id image.ID, ref reference.Named) error ***REMOVED***
		if _, ok := imgDescr[id]; !ok ***REMOVED***
			descr := &imageDescriptor***REMOVED******REMOVED***
			if err := l.takeLayerReference(id, descr); err != nil ***REMOVED***
				return err
			***REMOVED***
			imgDescr[id] = descr
		***REMOVED***

		if ref != nil ***REMOVED***
			if _, ok := ref.(reference.Canonical); ok ***REMOVED***
				return nil
			***REMOVED***
			tagged, ok := reference.TagNameOnly(ref).(reference.NamedTagged)
			if !ok ***REMOVED***
				return nil
			***REMOVED***

			for _, t := range imgDescr[id].refs ***REMOVED***
				if tagged.String() == t.String() ***REMOVED***
					return nil
				***REMOVED***
			***REMOVED***
			imgDescr[id].refs = append(imgDescr[id].refs, tagged)
		***REMOVED***
		return nil
	***REMOVED***

	for _, name := range names ***REMOVED***
		ref, err := reference.ParseAnyReference(name)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		namedRef, ok := ref.(reference.Named)
		if !ok ***REMOVED***
			// Check if digest ID reference
			if digested, ok := ref.(reference.Digested); ok ***REMOVED***
				id := image.IDFromDigest(digested.Digest())
				if err := addAssoc(id, nil); err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				continue
			***REMOVED***
			return nil, errors.Errorf("invalid reference: %v", name)
		***REMOVED***

		if reference.FamiliarName(namedRef) == string(digest.Canonical) ***REMOVED***
			imgID, err := l.is.Search(name)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if err := addAssoc(imgID, nil); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			continue
		***REMOVED***
		if reference.IsNameOnly(namedRef) ***REMOVED***
			assocs := l.rs.ReferencesByName(namedRef)
			for _, assoc := range assocs ***REMOVED***
				if err := addAssoc(image.IDFromDigest(assoc.ID), assoc.Ref); err != nil ***REMOVED***
					return nil, err
				***REMOVED***
			***REMOVED***
			if len(assocs) == 0 ***REMOVED***
				imgID, err := l.is.Search(name)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				if err := addAssoc(imgID, nil); err != nil ***REMOVED***
					return nil, err
				***REMOVED***
			***REMOVED***
			continue
		***REMOVED***
		id, err := l.rs.Get(namedRef)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if err := addAssoc(image.IDFromDigest(id), namedRef); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

	***REMOVED***
	return imgDescr, nil
***REMOVED***

// takeLayerReference will take/Get the image top layer reference
func (l *tarexporter) takeLayerReference(id image.ID, imgDescr *imageDescriptor) error ***REMOVED***
	img, err := l.is.Get(id)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	imgDescr.image = img
	topLayerID := img.RootFS.ChainID()
	if topLayerID == "" ***REMOVED***
		return nil
	***REMOVED***
	os := img.OS
	if os == "" ***REMOVED***
		os = runtime.GOOS
	***REMOVED***
	layer, err := l.lss[os].Get(topLayerID)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	imgDescr.layerRef = layer
	return nil
***REMOVED***

// releaseLayerReferences will release all the image top layer references
func (l *tarexporter) releaseLayerReferences(imgDescr map[image.ID]*imageDescriptor) error ***REMOVED***
	for _, descr := range imgDescr ***REMOVED***
		if descr.layerRef != nil ***REMOVED***
			os := descr.image.OS
			if os == "" ***REMOVED***
				os = runtime.GOOS
			***REMOVED***
			l.lss[os].Release(descr.layerRef)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (s *saveSession) save(outStream io.Writer) error ***REMOVED***
	s.savedLayers = make(map[string]struct***REMOVED******REMOVED***)
	s.diffIDPaths = make(map[layer.DiffID]string)

	// get image json
	tempDir, err := ioutil.TempDir("", "docker-export-")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer os.RemoveAll(tempDir)

	s.outDir = tempDir
	reposLegacy := make(map[string]map[string]string)

	var manifest []manifestItem
	var parentLinks []parentLink

	for id, imageDescr := range s.images ***REMOVED***
		foreignSrcs, err := s.saveImage(id)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		var repoTags []string
		var layers []string

		for _, ref := range imageDescr.refs ***REMOVED***
			familiarName := reference.FamiliarName(ref)
			if _, ok := reposLegacy[familiarName]; !ok ***REMOVED***
				reposLegacy[familiarName] = make(map[string]string)
			***REMOVED***
			reposLegacy[familiarName][ref.Tag()] = imageDescr.layers[len(imageDescr.layers)-1]
			repoTags = append(repoTags, reference.FamiliarString(ref))
		***REMOVED***

		for _, l := range imageDescr.layers ***REMOVED***
			layers = append(layers, filepath.Join(l, legacyLayerFileName))
		***REMOVED***

		manifest = append(manifest, manifestItem***REMOVED***
			Config:       id.Digest().Hex() + ".json",
			RepoTags:     repoTags,
			Layers:       layers,
			LayerSources: foreignSrcs,
		***REMOVED***)

		parentID, _ := s.is.GetParent(id)
		parentLinks = append(parentLinks, parentLink***REMOVED***id, parentID***REMOVED***)
		s.tarexporter.loggerImgEvent.LogImageEvent(id.String(), id.String(), "save")
	***REMOVED***

	for i, p := range validatedParentLinks(parentLinks) ***REMOVED***
		if p.parentID != "" ***REMOVED***
			manifest[i].Parent = p.parentID
		***REMOVED***
	***REMOVED***

	if len(reposLegacy) > 0 ***REMOVED***
		reposFile := filepath.Join(tempDir, legacyRepositoriesFileName)
		rf, err := os.OpenFile(reposFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := json.NewEncoder(rf).Encode(reposLegacy); err != nil ***REMOVED***
			rf.Close()
			return err
		***REMOVED***

		rf.Close()

		if err := system.Chtimes(reposFile, time.Unix(0, 0), time.Unix(0, 0)); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	manifestFileName := filepath.Join(tempDir, manifestFileName)
	f, err := os.OpenFile(manifestFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := json.NewEncoder(f).Encode(manifest); err != nil ***REMOVED***
		f.Close()
		return err
	***REMOVED***

	f.Close()

	if err := system.Chtimes(manifestFileName, time.Unix(0, 0), time.Unix(0, 0)); err != nil ***REMOVED***
		return err
	***REMOVED***

	fs, err := archive.Tar(tempDir, archive.Uncompressed)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer fs.Close()

	_, err = io.Copy(outStream, fs)
	return err
***REMOVED***

func (s *saveSession) saveImage(id image.ID) (map[layer.DiffID]distribution.Descriptor, error) ***REMOVED***
	img := s.images[id].image
	if len(img.RootFS.DiffIDs) == 0 ***REMOVED***
		return nil, fmt.Errorf("empty export - not implemented")
	***REMOVED***

	var parent digest.Digest
	var layers []string
	var foreignSrcs map[layer.DiffID]distribution.Descriptor
	for i := range img.RootFS.DiffIDs ***REMOVED***
		v1Img := image.V1Image***REMOVED***
			// This is for backward compatibility used for
			// pre v1.9 docker.
			Created: time.Unix(0, 0),
		***REMOVED***
		if i == len(img.RootFS.DiffIDs)-1 ***REMOVED***
			v1Img = img.V1Image
		***REMOVED***
		rootFS := *img.RootFS
		rootFS.DiffIDs = rootFS.DiffIDs[:i+1]
		v1ID, err := v1.CreateID(v1Img, rootFS.ChainID(), parent)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		v1Img.ID = v1ID.Hex()
		if parent != "" ***REMOVED***
			v1Img.Parent = parent.Hex()
		***REMOVED***

		src, err := s.saveLayer(rootFS.ChainID(), v1Img, img.Created)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		layers = append(layers, v1Img.ID)
		parent = v1ID
		if src.Digest != "" ***REMOVED***
			if foreignSrcs == nil ***REMOVED***
				foreignSrcs = make(map[layer.DiffID]distribution.Descriptor)
			***REMOVED***
			foreignSrcs[img.RootFS.DiffIDs[i]] = src
		***REMOVED***
	***REMOVED***

	configFile := filepath.Join(s.outDir, id.Digest().Hex()+".json")
	if err := ioutil.WriteFile(configFile, img.RawJSON(), 0644); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := system.Chtimes(configFile, img.Created, img.Created); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	s.images[id].layers = layers
	return foreignSrcs, nil
***REMOVED***

func (s *saveSession) saveLayer(id layer.ChainID, legacyImg image.V1Image, createdTime time.Time) (distribution.Descriptor, error) ***REMOVED***
	if _, exists := s.savedLayers[legacyImg.ID]; exists ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, nil
	***REMOVED***

	outDir := filepath.Join(s.outDir, legacyImg.ID)
	if err := os.Mkdir(outDir, 0755); err != nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***

	// todo: why is this version file here?
	if err := ioutil.WriteFile(filepath.Join(outDir, legacyVersionFileName), []byte("1.0"), 0644); err != nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***

	imageConfig, err := json.Marshal(legacyImg)
	if err != nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***

	if err := ioutil.WriteFile(filepath.Join(outDir, legacyConfigFileName), imageConfig, 0644); err != nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***

	// serialize filesystem
	layerPath := filepath.Join(outDir, legacyLayerFileName)
	operatingSystem := legacyImg.OS
	if operatingSystem == "" ***REMOVED***
		operatingSystem = runtime.GOOS
	***REMOVED***
	l, err := s.lss[operatingSystem].Get(id)
	if err != nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***
	defer layer.ReleaseAndLog(s.lss[operatingSystem], l)

	if oldPath, exists := s.diffIDPaths[l.DiffID()]; exists ***REMOVED***
		relPath, err := filepath.Rel(outDir, oldPath)
		if err != nil ***REMOVED***
			return distribution.Descriptor***REMOVED******REMOVED***, err
		***REMOVED***
		if err := os.Symlink(relPath, layerPath); err != nil ***REMOVED***
			return distribution.Descriptor***REMOVED******REMOVED***, errors.Wrap(err, "error creating symlink while saving layer")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// Use system.CreateSequential rather than os.Create. This ensures sequential
		// file access on Windows to avoid eating into MM standby list.
		// On Linux, this equates to a regular os.Create.
		tarFile, err := system.CreateSequential(layerPath)
		if err != nil ***REMOVED***
			return distribution.Descriptor***REMOVED******REMOVED***, err
		***REMOVED***
		defer tarFile.Close()

		arch, err := l.TarStream()
		if err != nil ***REMOVED***
			return distribution.Descriptor***REMOVED******REMOVED***, err
		***REMOVED***
		defer arch.Close()

		if _, err := io.Copy(tarFile, arch); err != nil ***REMOVED***
			return distribution.Descriptor***REMOVED******REMOVED***, err
		***REMOVED***

		for _, fname := range []string***REMOVED***"", legacyVersionFileName, legacyConfigFileName, legacyLayerFileName***REMOVED*** ***REMOVED***
			// todo: maybe save layer created timestamp?
			if err := system.Chtimes(filepath.Join(outDir, fname), createdTime, createdTime); err != nil ***REMOVED***
				return distribution.Descriptor***REMOVED******REMOVED***, err
			***REMOVED***
		***REMOVED***

		s.diffIDPaths[l.DiffID()] = layerPath
	***REMOVED***
	s.savedLayers[legacyImg.ID] = struct***REMOVED******REMOVED******REMOVED******REMOVED***

	var src distribution.Descriptor
	if fs, ok := l.(distribution.Describable); ok ***REMOVED***
		src = fs.Descriptor()
	***REMOVED***
	return src, nil
***REMOVED***
