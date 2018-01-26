package tarexport

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/docker/distribution"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/image"
	"github.com/docker/docker/image/v1"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/chrootarchive"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/streamformatter"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/pkg/symlink"
	"github.com/docker/docker/pkg/system"
	digest "github.com/opencontainers/go-digest"
	"github.com/sirupsen/logrus"
)

func (l *tarexporter) Load(inTar io.ReadCloser, outStream io.Writer, quiet bool) error ***REMOVED***
	var progressOutput progress.Output
	if !quiet ***REMOVED***
		progressOutput = streamformatter.NewJSONProgressOutput(outStream, false)
	***REMOVED***
	outStream = streamformatter.NewStdoutWriter(outStream)

	tmpDir, err := ioutil.TempDir("", "docker-import-")
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer os.RemoveAll(tmpDir)

	if err := chrootarchive.Untar(inTar, tmpDir, nil); err != nil ***REMOVED***
		return err
	***REMOVED***
	// read manifest, if no file then load in legacy mode
	manifestPath, err := safePath(tmpDir, manifestFileName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	manifestFile, err := os.Open(manifestPath)
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return l.legacyLoad(tmpDir, outStream, progressOutput)
		***REMOVED***
		return err
	***REMOVED***
	defer manifestFile.Close()

	var manifest []manifestItem
	if err := json.NewDecoder(manifestFile).Decode(&manifest); err != nil ***REMOVED***
		return err
	***REMOVED***

	var parentLinks []parentLink
	var imageIDsStr string
	var imageRefCount int

	for _, m := range manifest ***REMOVED***
		configPath, err := safePath(tmpDir, m.Config)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		config, err := ioutil.ReadFile(configPath)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		img, err := image.NewFromJSON(config)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := checkCompatibleOS(img.OS); err != nil ***REMOVED***
			return err
		***REMOVED***
		rootFS := *img.RootFS
		rootFS.DiffIDs = nil

		if expected, actual := len(m.Layers), len(img.RootFS.DiffIDs); expected != actual ***REMOVED***
			return fmt.Errorf("invalid manifest, layers length mismatch: expected %d, got %d", expected, actual)
		***REMOVED***

		// On Windows, validate the platform, defaulting to windows if not present.
		os := img.OS
		if os == "" ***REMOVED***
			os = runtime.GOOS
		***REMOVED***
		if runtime.GOOS == "windows" ***REMOVED***
			if (os != "windows") && (os != "linux") ***REMOVED***
				return fmt.Errorf("configuration for this image has an unsupported operating system: %s", os)
			***REMOVED***
		***REMOVED***

		for i, diffID := range img.RootFS.DiffIDs ***REMOVED***
			layerPath, err := safePath(tmpDir, m.Layers[i])
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			r := rootFS
			r.Append(diffID)
			newLayer, err := l.lss[os].Get(r.ChainID())
			if err != nil ***REMOVED***
				newLayer, err = l.loadLayer(layerPath, rootFS, diffID.String(), os, m.LayerSources[diffID], progressOutput)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
			defer layer.ReleaseAndLog(l.lss[os], newLayer)
			if expected, actual := diffID, newLayer.DiffID(); expected != actual ***REMOVED***
				return fmt.Errorf("invalid diffID for layer %d: expected %q, got %q", i, expected, actual)
			***REMOVED***
			rootFS.Append(diffID)
		***REMOVED***

		imgID, err := l.is.Create(config)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		imageIDsStr += fmt.Sprintf("Loaded image ID: %s\n", imgID)

		imageRefCount = 0
		for _, repoTag := range m.RepoTags ***REMOVED***
			named, err := reference.ParseNormalizedNamed(repoTag)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			ref, ok := named.(reference.NamedTagged)
			if !ok ***REMOVED***
				return fmt.Errorf("invalid tag %q", repoTag)
			***REMOVED***
			l.setLoadedTag(ref, imgID.Digest(), outStream)
			outStream.Write([]byte(fmt.Sprintf("Loaded image: %s\n", reference.FamiliarString(ref))))
			imageRefCount++
		***REMOVED***

		parentLinks = append(parentLinks, parentLink***REMOVED***imgID, m.Parent***REMOVED***)
		l.loggerImgEvent.LogImageEvent(imgID.String(), imgID.String(), "load")
	***REMOVED***

	for _, p := range validatedParentLinks(parentLinks) ***REMOVED***
		if p.parentID != "" ***REMOVED***
			if err := l.setParentID(p.id, p.parentID); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if imageRefCount == 0 ***REMOVED***
		outStream.Write([]byte(imageIDsStr))
	***REMOVED***

	return nil
***REMOVED***

func (l *tarexporter) setParentID(id, parentID image.ID) error ***REMOVED***
	img, err := l.is.Get(id)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	parent, err := l.is.Get(parentID)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if !checkValidParent(img, parent) ***REMOVED***
		return fmt.Errorf("image %v is not a valid parent for %v", parent.ID(), img.ID())
	***REMOVED***
	return l.is.SetParent(id, parentID)
***REMOVED***

func (l *tarexporter) loadLayer(filename string, rootFS image.RootFS, id string, os string, foreignSrc distribution.Descriptor, progressOutput progress.Output) (layer.Layer, error) ***REMOVED***
	// We use system.OpenSequential to use sequential file access on Windows, avoiding
	// depleting the standby list. On Linux, this equates to a regular os.Open.
	rawTar, err := system.OpenSequential(filename)
	if err != nil ***REMOVED***
		logrus.Debugf("Error reading embedded tar: %v", err)
		return nil, err
	***REMOVED***
	defer rawTar.Close()

	var r io.Reader
	if progressOutput != nil ***REMOVED***
		fileInfo, err := rawTar.Stat()
		if err != nil ***REMOVED***
			logrus.Debugf("Error statting file: %v", err)
			return nil, err
		***REMOVED***

		r = progress.NewProgressReader(rawTar, progressOutput, fileInfo.Size(), stringid.TruncateID(id), "Loading layer")
	***REMOVED*** else ***REMOVED***
		r = rawTar
	***REMOVED***

	inflatedLayerData, err := archive.DecompressStream(r)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer inflatedLayerData.Close()

	if ds, ok := l.lss[os].(layer.DescribableStore); ok ***REMOVED***
		return ds.RegisterWithDescriptor(inflatedLayerData, rootFS.ChainID(), foreignSrc)
	***REMOVED***
	return l.lss[os].Register(inflatedLayerData, rootFS.ChainID())
***REMOVED***

func (l *tarexporter) setLoadedTag(ref reference.Named, imgID digest.Digest, outStream io.Writer) error ***REMOVED***
	if prevID, err := l.rs.Get(ref); err == nil && prevID != imgID ***REMOVED***
		fmt.Fprintf(outStream, "The image %s already exists, renaming the old one with ID %s to empty string\n", reference.FamiliarString(ref), string(prevID)) // todo: this message is wrong in case of multiple tags
	***REMOVED***

	return l.rs.AddTag(ref, imgID, true)
***REMOVED***

func (l *tarexporter) legacyLoad(tmpDir string, outStream io.Writer, progressOutput progress.Output) error ***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED***
		return errors.New("Windows does not support legacy loading of images")
	***REMOVED***

	legacyLoadedMap := make(map[string]image.ID)

	dirs, err := ioutil.ReadDir(tmpDir)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// every dir represents an image
	for _, d := range dirs ***REMOVED***
		if d.IsDir() ***REMOVED***
			if err := l.legacyLoadImage(d.Name(), tmpDir, legacyLoadedMap, progressOutput); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// load tags from repositories file
	repositoriesPath, err := safePath(tmpDir, legacyRepositoriesFileName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	repositoriesFile, err := os.Open(repositoriesPath)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer repositoriesFile.Close()

	repositories := make(map[string]map[string]string)
	if err := json.NewDecoder(repositoriesFile).Decode(&repositories); err != nil ***REMOVED***
		return err
	***REMOVED***

	for name, tagMap := range repositories ***REMOVED***
		for tag, oldID := range tagMap ***REMOVED***
			imgID, ok := legacyLoadedMap[oldID]
			if !ok ***REMOVED***
				return fmt.Errorf("invalid target ID: %v", oldID)
			***REMOVED***
			named, err := reference.ParseNormalizedNamed(name)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			ref, err := reference.WithTag(named, tag)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			l.setLoadedTag(ref, imgID.Digest(), outStream)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (l *tarexporter) legacyLoadImage(oldID, sourceDir string, loadedMap map[string]image.ID, progressOutput progress.Output) error ***REMOVED***
	if _, loaded := loadedMap[oldID]; loaded ***REMOVED***
		return nil
	***REMOVED***
	configPath, err := safePath(sourceDir, filepath.Join(oldID, legacyConfigFileName))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	imageJSON, err := ioutil.ReadFile(configPath)
	if err != nil ***REMOVED***
		logrus.Debugf("Error reading json: %v", err)
		return err
	***REMOVED***

	var img struct ***REMOVED***
		OS     string
		Parent string
	***REMOVED***
	if err := json.Unmarshal(imageJSON, &img); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := checkCompatibleOS(img.OS); err != nil ***REMOVED***
		return err
	***REMOVED***
	if img.OS == "" ***REMOVED***
		img.OS = runtime.GOOS
	***REMOVED***

	var parentID image.ID
	if img.Parent != "" ***REMOVED***
		for ***REMOVED***
			var loaded bool
			if parentID, loaded = loadedMap[img.Parent]; !loaded ***REMOVED***
				if err := l.legacyLoadImage(img.Parent, sourceDir, loadedMap, progressOutput); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// todo: try to connect with migrate code
	rootFS := image.NewRootFS()
	var history []image.History

	if parentID != "" ***REMOVED***
		parentImg, err := l.is.Get(parentID)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		rootFS = parentImg.RootFS
		history = parentImg.History
	***REMOVED***

	layerPath, err := safePath(sourceDir, filepath.Join(oldID, legacyLayerFileName))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	newLayer, err := l.loadLayer(layerPath, *rootFS, oldID, img.OS, distribution.Descriptor***REMOVED******REMOVED***, progressOutput)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	rootFS.Append(newLayer.DiffID())

	h, err := v1.HistoryFromConfig(imageJSON, false)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	history = append(history, h)

	config, err := v1.MakeConfigFromV1Config(imageJSON, rootFS, history)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	imgID, err := l.is.Create(config)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	metadata, err := l.lss[img.OS].Release(newLayer)
	layer.LogReleaseMetadata(metadata)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if parentID != "" ***REMOVED***
		if err := l.is.SetParent(imgID, parentID); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	loadedMap[oldID] = imgID
	return nil
***REMOVED***

func safePath(base, path string) (string, error) ***REMOVED***
	return symlink.FollowSymlinkInScope(filepath.Join(base, path), base)
***REMOVED***

type parentLink struct ***REMOVED***
	id, parentID image.ID
***REMOVED***

func validatedParentLinks(pl []parentLink) (ret []parentLink) ***REMOVED***
mainloop:
	for i, p := range pl ***REMOVED***
		ret = append(ret, p)
		for _, p2 := range pl ***REMOVED***
			if p2.id == p.parentID && p2.id != p.id ***REMOVED***
				continue mainloop
			***REMOVED***
		***REMOVED***
		ret[i].parentID = ""
	***REMOVED***
	return
***REMOVED***

func checkValidParent(img, parent *image.Image) bool ***REMOVED***
	if len(img.History) == 0 && len(parent.History) == 0 ***REMOVED***
		return true // having history is not mandatory
	***REMOVED***
	if len(img.History)-len(parent.History) != 1 ***REMOVED***
		return false
	***REMOVED***
	for i, h := range parent.History ***REMOVED***
		if !reflect.DeepEqual(h, img.History[i]) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func checkCompatibleOS(imageOS string) error ***REMOVED***
	// always compatible if the images OS matches the host OS; also match an empty image OS
	if imageOS == runtime.GOOS || imageOS == "" ***REMOVED***
		return nil
	***REMOVED***
	// On non-Windows hosts, for compatibility, fail if the image is Windows.
	if runtime.GOOS != "windows" && imageOS == "windows" ***REMOVED***
		return fmt.Errorf("cannot load %s image on %s", imageOS, runtime.GOOS)
	***REMOVED***
	// Finally, check the image OS is supported for the platform.
	if err := system.ValidatePlatform(system.ParsePlatform(imageOS)); err != nil ***REMOVED***
		return fmt.Errorf("cannot load %s image on %s: %s", imageOS, runtime.GOOS, err)
	***REMOVED***
	return nil
***REMOVED***
