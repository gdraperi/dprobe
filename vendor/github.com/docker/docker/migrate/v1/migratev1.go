package v1

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"

	"encoding/json"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/distribution/metadata"
	"github.com/docker/docker/image"
	imagev1 "github.com/docker/docker/image/v1"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/ioutils"
	refstore "github.com/docker/docker/reference"
	"github.com/opencontainers/go-digest"
	"github.com/sirupsen/logrus"
)

type graphIDRegistrar interface ***REMOVED***
	RegisterByGraphID(string, layer.ChainID, layer.DiffID, string, int64) (layer.Layer, error)
	Release(layer.Layer) ([]layer.Metadata, error)
***REMOVED***

type graphIDMounter interface ***REMOVED***
	CreateRWLayerByGraphID(string, string, layer.ChainID) error
***REMOVED***

type checksumCalculator interface ***REMOVED***
	ChecksumForGraphID(id, parent, oldTarDataPath, newTarDataPath string) (diffID layer.DiffID, size int64, err error)
***REMOVED***

const (
	graphDirName                 = "graph"
	tarDataFileName              = "tar-data.json.gz"
	migrationFileName            = ".migration-v1-images.json"
	migrationTagsFileName        = ".migration-v1-tags"
	migrationDiffIDFileName      = ".migration-diffid"
	migrationSizeFileName        = ".migration-size"
	migrationTarDataFileName     = ".migration-tardata"
	containersDirName            = "containers"
	configFileNameLegacy         = "config.json"
	configFileName               = "config.v2.json"
	repositoriesFilePrefixLegacy = "repositories-"
)

var (
	errUnsupported = errors.New("migration is not supported")
)

// Migrate takes an old graph directory and transforms the metadata into the
// new format.
func Migrate(root, driverName string, ls layer.Store, is image.Store, rs refstore.Store, ms metadata.Store) error ***REMOVED***
	graphDir := filepath.Join(root, graphDirName)
	if _, err := os.Lstat(graphDir); os.IsNotExist(err) ***REMOVED***
		return nil
	***REMOVED***

	mappings, err := restoreMappings(root)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if cc, ok := ls.(checksumCalculator); ok ***REMOVED***
		CalculateLayerChecksums(root, cc, mappings)
	***REMOVED***

	if registrar, ok := ls.(graphIDRegistrar); !ok ***REMOVED***
		return errUnsupported
	***REMOVED*** else if err := migrateImages(root, registrar, is, ms, mappings); err != nil ***REMOVED***
		return err
	***REMOVED***

	err = saveMappings(root, mappings)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if mounter, ok := ls.(graphIDMounter); !ok ***REMOVED***
		return errUnsupported
	***REMOVED*** else if err := migrateContainers(root, mounter, is, mappings); err != nil ***REMOVED***
		return err
	***REMOVED***

	return migrateRefs(root, driverName, rs, mappings)
***REMOVED***

// CalculateLayerChecksums walks an old graph directory and calculates checksums
// for each layer. These checksums are later used for migration.
func CalculateLayerChecksums(root string, ls checksumCalculator, mappings map[string]image.ID) ***REMOVED***
	graphDir := filepath.Join(root, graphDirName)
	// spawn some extra workers also for maximum performance because the process is bounded by both cpu and io
	workers := runtime.NumCPU() * 3
	workQueue := make(chan string, workers)

	wg := sync.WaitGroup***REMOVED******REMOVED***

	for i := 0; i < workers; i++ ***REMOVED***
		wg.Add(1)
		go func() ***REMOVED***
			for id := range workQueue ***REMOVED***
				start := time.Now()
				if err := calculateLayerChecksum(graphDir, id, ls); err != nil ***REMOVED***
					logrus.Errorf("could not calculate checksum for %q, %q", id, err)
				***REMOVED***
				elapsed := time.Since(start)
				logrus.Debugf("layer %s took %.2f seconds", id, elapsed.Seconds())
			***REMOVED***
			wg.Done()
		***REMOVED***()
	***REMOVED***

	dir, err := ioutil.ReadDir(graphDir)
	if err != nil ***REMOVED***
		logrus.Errorf("could not read directory %q", graphDir)
		return
	***REMOVED***
	for _, v := range dir ***REMOVED***
		v1ID := v.Name()
		if err := imagev1.ValidateID(v1ID); err != nil ***REMOVED***
			continue
		***REMOVED***
		if _, ok := mappings[v1ID]; ok ***REMOVED*** // support old migrations without helper files
			continue
		***REMOVED***
		workQueue <- v1ID
	***REMOVED***
	close(workQueue)
	wg.Wait()
***REMOVED***

func calculateLayerChecksum(graphDir, id string, ls checksumCalculator) error ***REMOVED***
	diffIDFile := filepath.Join(graphDir, id, migrationDiffIDFileName)
	if _, err := os.Lstat(diffIDFile); err == nil ***REMOVED***
		return nil
	***REMOVED*** else if !os.IsNotExist(err) ***REMOVED***
		return err
	***REMOVED***

	parent, err := getParent(filepath.Join(graphDir, id))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	diffID, size, err := ls.ChecksumForGraphID(id, parent, filepath.Join(graphDir, id, tarDataFileName), filepath.Join(graphDir, id, migrationTarDataFileName))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := ioutil.WriteFile(filepath.Join(graphDir, id, migrationSizeFileName), []byte(strconv.Itoa(int(size))), 0600); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := ioutils.AtomicWriteFile(filepath.Join(graphDir, id, migrationDiffIDFileName), []byte(diffID), 0600); err != nil ***REMOVED***
		return err
	***REMOVED***

	logrus.Infof("calculated checksum for layer %s: %s", id, diffID)
	return nil
***REMOVED***

func restoreMappings(root string) (map[string]image.ID, error) ***REMOVED***
	mappings := make(map[string]image.ID)

	mfile := filepath.Join(root, migrationFileName)
	f, err := os.Open(mfile)
	if err != nil && !os.IsNotExist(err) ***REMOVED***
		return nil, err
	***REMOVED*** else if err == nil ***REMOVED***
		err := json.NewDecoder(f).Decode(&mappings)
		if err != nil ***REMOVED***
			f.Close()
			return nil, err
		***REMOVED***
		f.Close()
	***REMOVED***

	return mappings, nil
***REMOVED***

func saveMappings(root string, mappings map[string]image.ID) error ***REMOVED***
	mfile := filepath.Join(root, migrationFileName)
	f, err := os.OpenFile(mfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer f.Close()
	return json.NewEncoder(f).Encode(mappings)
***REMOVED***

func migrateImages(root string, ls graphIDRegistrar, is image.Store, ms metadata.Store, mappings map[string]image.ID) error ***REMOVED***
	graphDir := filepath.Join(root, graphDirName)

	dir, err := ioutil.ReadDir(graphDir)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, v := range dir ***REMOVED***
		v1ID := v.Name()
		if err := imagev1.ValidateID(v1ID); err != nil ***REMOVED***
			continue
		***REMOVED***
		if _, exists := mappings[v1ID]; exists ***REMOVED***
			continue
		***REMOVED***
		if err := migrateImage(v1ID, root, ls, is, ms, mappings); err != nil ***REMOVED***
			continue
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func migrateContainers(root string, ls graphIDMounter, is image.Store, imageMappings map[string]image.ID) error ***REMOVED***
	containersDir := filepath.Join(root, containersDirName)
	dir, err := ioutil.ReadDir(containersDir)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, v := range dir ***REMOVED***
		id := v.Name()

		if _, err := os.Stat(filepath.Join(containersDir, id, configFileName)); err == nil ***REMOVED***
			continue
		***REMOVED***

		containerJSON, err := ioutil.ReadFile(filepath.Join(containersDir, id, configFileNameLegacy))
		if err != nil ***REMOVED***
			logrus.Errorf("migrate container error: %v", err)
			continue
		***REMOVED***

		var c map[string]*json.RawMessage
		if err := json.Unmarshal(containerJSON, &c); err != nil ***REMOVED***
			logrus.Errorf("migrate container error: %v", err)
			continue
		***REMOVED***

		imageStrJSON, ok := c["Image"]
		if !ok ***REMOVED***
			return fmt.Errorf("invalid container configuration for %v", id)
		***REMOVED***

		var image string
		if err := json.Unmarshal([]byte(*imageStrJSON), &image); err != nil ***REMOVED***
			logrus.Errorf("migrate container error: %v", err)
			continue
		***REMOVED***

		imageID, ok := imageMappings[image]
		if !ok ***REMOVED***
			logrus.Errorf("image not migrated %v", imageID) // non-fatal error
			continue
		***REMOVED***

		c["Image"] = rawJSON(imageID)

		containerJSON, err = json.Marshal(c)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := ioutil.WriteFile(filepath.Join(containersDir, id, configFileName), containerJSON, 0600); err != nil ***REMOVED***
			return err
		***REMOVED***

		img, err := is.Get(imageID)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := ls.CreateRWLayerByGraphID(id, id, img.RootFS.ChainID()); err != nil ***REMOVED***
			logrus.Errorf("migrate container error: %v", err)
			continue
		***REMOVED***

		logrus.Infof("migrated container %s to point to %s", id, imageID)

	***REMOVED***
	return nil
***REMOVED***

type refAdder interface ***REMOVED***
	AddTag(ref reference.Named, id digest.Digest, force bool) error
	AddDigest(ref reference.Canonical, id digest.Digest, force bool) error
***REMOVED***

func migrateRefs(root, driverName string, rs refAdder, mappings map[string]image.ID) error ***REMOVED***
	migrationFile := filepath.Join(root, migrationTagsFileName)
	if _, err := os.Lstat(migrationFile); !os.IsNotExist(err) ***REMOVED***
		return err
	***REMOVED***

	type repositories struct ***REMOVED***
		Repositories map[string]map[string]string
	***REMOVED***

	var repos repositories

	f, err := os.Open(filepath.Join(root, repositoriesFilePrefixLegacy+driverName))
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return nil
		***REMOVED***
		return err
	***REMOVED***
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&repos); err != nil ***REMOVED***
		return err
	***REMOVED***

	for name, repo := range repos.Repositories ***REMOVED***
		for tag, id := range repo ***REMOVED***
			if strongID, exists := mappings[id]; exists ***REMOVED***
				ref, err := reference.ParseNormalizedNamed(name)
				if err != nil ***REMOVED***
					logrus.Errorf("migrate tags: invalid name %q, %q", name, err)
					continue
				***REMOVED***
				if !reference.IsNameOnly(ref) ***REMOVED***
					logrus.Errorf("migrate tags: invalid name %q, unexpected tag or digest", name)
					continue
				***REMOVED***
				if dgst, err := digest.Parse(tag); err == nil ***REMOVED***
					canonical, err := reference.WithDigest(reference.TrimNamed(ref), dgst)
					if err != nil ***REMOVED***
						logrus.Errorf("migrate tags: invalid digest %q, %q", dgst, err)
						continue
					***REMOVED***
					if err := rs.AddDigest(canonical, strongID.Digest(), false); err != nil ***REMOVED***
						logrus.Errorf("can't migrate digest %q for %q, err: %q", reference.FamiliarString(ref), strongID, err)
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					tagRef, err := reference.WithTag(ref, tag)
					if err != nil ***REMOVED***
						logrus.Errorf("migrate tags: invalid tag %q, %q", tag, err)
						continue
					***REMOVED***
					if err := rs.AddTag(tagRef, strongID.Digest(), false); err != nil ***REMOVED***
						logrus.Errorf("can't migrate tag %q for %q, err: %q", reference.FamiliarString(ref), strongID, err)
					***REMOVED***
				***REMOVED***
				logrus.Infof("migrated tag %s:%s to point to %s", name, tag, strongID)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	mf, err := os.Create(migrationFile)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	mf.Close()

	return nil
***REMOVED***

func getParent(confDir string) (string, error) ***REMOVED***
	jsonFile := filepath.Join(confDir, "json")
	imageJSON, err := ioutil.ReadFile(jsonFile)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	var parent struct ***REMOVED***
		Parent   string
		ParentID digest.Digest `json:"parent_id"`
	***REMOVED***
	if err := json.Unmarshal(imageJSON, &parent); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if parent.Parent == "" && parent.ParentID != "" ***REMOVED*** // v1.9
		parent.Parent = parent.ParentID.Hex()
	***REMOVED***
	// compatibilityID for parent
	parentCompatibilityID, err := ioutil.ReadFile(filepath.Join(confDir, "parent"))
	if err == nil && len(parentCompatibilityID) > 0 ***REMOVED***
		parent.Parent = string(parentCompatibilityID)
	***REMOVED***
	return parent.Parent, nil
***REMOVED***

func migrateImage(id, root string, ls graphIDRegistrar, is image.Store, ms metadata.Store, mappings map[string]image.ID) (err error) ***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			logrus.Errorf("migration failed for %v, err: %v", id, err)
		***REMOVED***
	***REMOVED***()

	parent, err := getParent(filepath.Join(root, graphDirName, id))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var parentID image.ID
	if parent != "" ***REMOVED***
		var exists bool
		if parentID, exists = mappings[parent]; !exists ***REMOVED***
			if err := migrateImage(parent, root, ls, is, ms, mappings); err != nil ***REMOVED***
				// todo: fail or allow broken chains?
				return err
			***REMOVED***
			parentID = mappings[parent]
		***REMOVED***
	***REMOVED***

	rootFS := image.NewRootFS()
	var history []image.History

	if parentID != "" ***REMOVED***
		parentImg, err := is.Get(parentID)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		rootFS = parentImg.RootFS
		history = parentImg.History
	***REMOVED***

	diffIDData, err := ioutil.ReadFile(filepath.Join(root, graphDirName, id, migrationDiffIDFileName))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	diffID, err := digest.Parse(string(diffIDData))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	sizeStr, err := ioutil.ReadFile(filepath.Join(root, graphDirName, id, migrationSizeFileName))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	size, err := strconv.ParseInt(string(sizeStr), 10, 64)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	layer, err := ls.RegisterByGraphID(id, rootFS.ChainID(), layer.DiffID(diffID), filepath.Join(root, graphDirName, id, migrationTarDataFileName), size)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	logrus.Infof("migrated layer %s to %s", id, layer.DiffID())

	jsonFile := filepath.Join(root, graphDirName, id, "json")
	imageJSON, err := ioutil.ReadFile(jsonFile)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	h, err := imagev1.HistoryFromConfig(imageJSON, false)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	history = append(history, h)

	rootFS.Append(layer.DiffID())

	config, err := imagev1.MakeConfigFromV1Config(imageJSON, rootFS, history)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	strongID, err := is.Create(config)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	logrus.Infof("migrated image %s to %s", id, strongID)

	if parentID != "" ***REMOVED***
		if err := is.SetParent(strongID, parentID); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	checksum, err := ioutil.ReadFile(filepath.Join(root, graphDirName, id, "checksum"))
	if err == nil ***REMOVED*** // best effort
		dgst, err := digest.Parse(string(checksum))
		if err == nil ***REMOVED***
			V2MetadataService := metadata.NewV2MetadataService(ms)
			V2MetadataService.Add(layer.DiffID(), metadata.V2Metadata***REMOVED***Digest: dgst***REMOVED***)
		***REMOVED***
	***REMOVED***
	_, err = ls.Release(layer)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	mappings[id] = strongID
	return
***REMOVED***

func rawJSON(value interface***REMOVED******REMOVED***) *json.RawMessage ***REMOVED***
	jsonval, err := json.Marshal(value)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	return (*json.RawMessage)(&jsonval)
***REMOVED***
