package layer

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/docker/distribution"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/opencontainers/go-digest"
	"github.com/sirupsen/logrus"
)

var (
	stringIDRegexp      = regexp.MustCompile(`^[a-f0-9]***REMOVED***64***REMOVED***(-init)?$`)
	supportedAlgorithms = []digest.Algorithm***REMOVED***
		digest.SHA256,
		// digest.SHA384, // Currently not used
		// digest.SHA512, // Currently not used
	***REMOVED***
)

type fileMetadataStore struct ***REMOVED***
	root string
***REMOVED***

type fileMetadataTransaction struct ***REMOVED***
	store *fileMetadataStore
	ws    *ioutils.AtomicWriteSet
***REMOVED***

// NewFSMetadataStore returns an instance of a metadata store
// which is backed by files on disk using the provided root
// as the root of metadata files.
func NewFSMetadataStore(root string) (MetadataStore, error) ***REMOVED***
	if err := os.MkdirAll(root, 0700); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &fileMetadataStore***REMOVED***
		root: root,
	***REMOVED***, nil
***REMOVED***

func (fms *fileMetadataStore) getLayerDirectory(layer ChainID) string ***REMOVED***
	dgst := digest.Digest(layer)
	return filepath.Join(fms.root, string(dgst.Algorithm()), dgst.Hex())
***REMOVED***

func (fms *fileMetadataStore) getLayerFilename(layer ChainID, filename string) string ***REMOVED***
	return filepath.Join(fms.getLayerDirectory(layer), filename)
***REMOVED***

func (fms *fileMetadataStore) getMountDirectory(mount string) string ***REMOVED***
	return filepath.Join(fms.root, "mounts", mount)
***REMOVED***

func (fms *fileMetadataStore) getMountFilename(mount, filename string) string ***REMOVED***
	return filepath.Join(fms.getMountDirectory(mount), filename)
***REMOVED***

func (fms *fileMetadataStore) StartTransaction() (MetadataTransaction, error) ***REMOVED***
	tmpDir := filepath.Join(fms.root, "tmp")
	if err := os.MkdirAll(tmpDir, 0755); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ws, err := ioutils.NewAtomicWriteSet(tmpDir)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &fileMetadataTransaction***REMOVED***
		store: fms,
		ws:    ws,
	***REMOVED***, nil
***REMOVED***

func (fm *fileMetadataTransaction) SetSize(size int64) error ***REMOVED***
	content := fmt.Sprintf("%d", size)
	return fm.ws.WriteFile("size", []byte(content), 0644)
***REMOVED***

func (fm *fileMetadataTransaction) SetParent(parent ChainID) error ***REMOVED***
	return fm.ws.WriteFile("parent", []byte(digest.Digest(parent).String()), 0644)
***REMOVED***

func (fm *fileMetadataTransaction) SetDiffID(diff DiffID) error ***REMOVED***
	return fm.ws.WriteFile("diff", []byte(digest.Digest(diff).String()), 0644)
***REMOVED***

func (fm *fileMetadataTransaction) SetCacheID(cacheID string) error ***REMOVED***
	return fm.ws.WriteFile("cache-id", []byte(cacheID), 0644)
***REMOVED***

func (fm *fileMetadataTransaction) SetDescriptor(ref distribution.Descriptor) error ***REMOVED***
	jsonRef, err := json.Marshal(ref)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return fm.ws.WriteFile("descriptor.json", jsonRef, 0644)
***REMOVED***

func (fm *fileMetadataTransaction) TarSplitWriter(compressInput bool) (io.WriteCloser, error) ***REMOVED***
	f, err := fm.ws.FileWriter("tar-split.json.gz", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var wc io.WriteCloser
	if compressInput ***REMOVED***
		wc = gzip.NewWriter(f)
	***REMOVED*** else ***REMOVED***
		wc = f
	***REMOVED***

	return ioutils.NewWriteCloserWrapper(wc, func() error ***REMOVED***
		wc.Close()
		return f.Close()
	***REMOVED***), nil
***REMOVED***

func (fm *fileMetadataTransaction) Commit(layer ChainID) error ***REMOVED***
	finalDir := fm.store.getLayerDirectory(layer)
	if err := os.MkdirAll(filepath.Dir(finalDir), 0755); err != nil ***REMOVED***
		return err
	***REMOVED***

	return fm.ws.Commit(finalDir)
***REMOVED***

func (fm *fileMetadataTransaction) Cancel() error ***REMOVED***
	return fm.ws.Cancel()
***REMOVED***

func (fm *fileMetadataTransaction) String() string ***REMOVED***
	return fm.ws.String()
***REMOVED***

func (fms *fileMetadataStore) GetSize(layer ChainID) (int64, error) ***REMOVED***
	content, err := ioutil.ReadFile(fms.getLayerFilename(layer, "size"))
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	size, err := strconv.ParseInt(string(content), 10, 64)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	return size, nil
***REMOVED***

func (fms *fileMetadataStore) GetParent(layer ChainID) (ChainID, error) ***REMOVED***
	content, err := ioutil.ReadFile(fms.getLayerFilename(layer, "parent"))
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return "", nil
		***REMOVED***
		return "", err
	***REMOVED***

	dgst, err := digest.Parse(strings.TrimSpace(string(content)))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return ChainID(dgst), nil
***REMOVED***

func (fms *fileMetadataStore) GetDiffID(layer ChainID) (DiffID, error) ***REMOVED***
	content, err := ioutil.ReadFile(fms.getLayerFilename(layer, "diff"))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	dgst, err := digest.Parse(strings.TrimSpace(string(content)))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return DiffID(dgst), nil
***REMOVED***

func (fms *fileMetadataStore) GetCacheID(layer ChainID) (string, error) ***REMOVED***
	contentBytes, err := ioutil.ReadFile(fms.getLayerFilename(layer, "cache-id"))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	content := strings.TrimSpace(string(contentBytes))

	if !stringIDRegexp.MatchString(content) ***REMOVED***
		return "", errors.New("invalid cache id value")
	***REMOVED***

	return content, nil
***REMOVED***

func (fms *fileMetadataStore) GetDescriptor(layer ChainID) (distribution.Descriptor, error) ***REMOVED***
	content, err := ioutil.ReadFile(fms.getLayerFilename(layer, "descriptor.json"))
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			// only return empty descriptor to represent what is stored
			return distribution.Descriptor***REMOVED******REMOVED***, nil
		***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***

	var ref distribution.Descriptor
	err = json.Unmarshal(content, &ref)
	if err != nil ***REMOVED***
		return distribution.Descriptor***REMOVED******REMOVED***, err
	***REMOVED***
	return ref, err
***REMOVED***

func (fms *fileMetadataStore) TarSplitReader(layer ChainID) (io.ReadCloser, error) ***REMOVED***
	fz, err := os.Open(fms.getLayerFilename(layer, "tar-split.json.gz"))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	f, err := gzip.NewReader(fz)
	if err != nil ***REMOVED***
		fz.Close()
		return nil, err
	***REMOVED***

	return ioutils.NewReadCloserWrapper(f, func() error ***REMOVED***
		f.Close()
		return fz.Close()
	***REMOVED***), nil
***REMOVED***

func (fms *fileMetadataStore) SetMountID(mount string, mountID string) error ***REMOVED***
	if err := os.MkdirAll(fms.getMountDirectory(mount), 0755); err != nil ***REMOVED***
		return err
	***REMOVED***
	return ioutil.WriteFile(fms.getMountFilename(mount, "mount-id"), []byte(mountID), 0644)
***REMOVED***

func (fms *fileMetadataStore) SetInitID(mount string, init string) error ***REMOVED***
	if err := os.MkdirAll(fms.getMountDirectory(mount), 0755); err != nil ***REMOVED***
		return err
	***REMOVED***
	return ioutil.WriteFile(fms.getMountFilename(mount, "init-id"), []byte(init), 0644)
***REMOVED***

func (fms *fileMetadataStore) SetMountParent(mount string, parent ChainID) error ***REMOVED***
	if err := os.MkdirAll(fms.getMountDirectory(mount), 0755); err != nil ***REMOVED***
		return err
	***REMOVED***
	return ioutil.WriteFile(fms.getMountFilename(mount, "parent"), []byte(digest.Digest(parent).String()), 0644)
***REMOVED***

func (fms *fileMetadataStore) GetMountID(mount string) (string, error) ***REMOVED***
	contentBytes, err := ioutil.ReadFile(fms.getMountFilename(mount, "mount-id"))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	content := strings.TrimSpace(string(contentBytes))

	if !stringIDRegexp.MatchString(content) ***REMOVED***
		return "", errors.New("invalid mount id value")
	***REMOVED***

	return content, nil
***REMOVED***

func (fms *fileMetadataStore) GetInitID(mount string) (string, error) ***REMOVED***
	contentBytes, err := ioutil.ReadFile(fms.getMountFilename(mount, "init-id"))
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return "", nil
		***REMOVED***
		return "", err
	***REMOVED***
	content := strings.TrimSpace(string(contentBytes))

	if !stringIDRegexp.MatchString(content) ***REMOVED***
		return "", errors.New("invalid init id value")
	***REMOVED***

	return content, nil
***REMOVED***

func (fms *fileMetadataStore) GetMountParent(mount string) (ChainID, error) ***REMOVED***
	content, err := ioutil.ReadFile(fms.getMountFilename(mount, "parent"))
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return "", nil
		***REMOVED***
		return "", err
	***REMOVED***

	dgst, err := digest.Parse(strings.TrimSpace(string(content)))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return ChainID(dgst), nil
***REMOVED***

func (fms *fileMetadataStore) List() ([]ChainID, []string, error) ***REMOVED***
	var ids []ChainID
	for _, algorithm := range supportedAlgorithms ***REMOVED***
		fileInfos, err := ioutil.ReadDir(filepath.Join(fms.root, string(algorithm)))
		if err != nil ***REMOVED***
			if os.IsNotExist(err) ***REMOVED***
				continue
			***REMOVED***
			return nil, nil, err
		***REMOVED***

		for _, fi := range fileInfos ***REMOVED***
			if fi.IsDir() && fi.Name() != "mounts" ***REMOVED***
				dgst := digest.NewDigestFromHex(string(algorithm), fi.Name())
				if err := dgst.Validate(); err != nil ***REMOVED***
					logrus.Debugf("Ignoring invalid digest %s:%s", algorithm, fi.Name())
				***REMOVED*** else ***REMOVED***
					ids = append(ids, ChainID(dgst))
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	fileInfos, err := ioutil.ReadDir(filepath.Join(fms.root, "mounts"))
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return ids, []string***REMOVED******REMOVED***, nil
		***REMOVED***
		return nil, nil, err
	***REMOVED***

	var mounts []string
	for _, fi := range fileInfos ***REMOVED***
		if fi.IsDir() ***REMOVED***
			mounts = append(mounts, fi.Name())
		***REMOVED***
	***REMOVED***

	return ids, mounts, nil
***REMOVED***

func (fms *fileMetadataStore) Remove(layer ChainID) error ***REMOVED***
	return os.RemoveAll(fms.getLayerDirectory(layer))
***REMOVED***

func (fms *fileMetadataStore) RemoveMount(mount string) error ***REMOVED***
	return os.RemoveAll(fms.getMountDirectory(mount))
***REMOVED***
