// Package tarsum provides algorithms to perform checksum calculation on
// filesystem layers.
//
// The transportation of filesystems, regarding Docker, is done with tar(1)
// archives. There are a variety of tar serialization formats [2], and a key
// concern here is ensuring a repeatable checksum given a set of inputs from a
// generic tar archive. Types of transportation include distribution to and from a
// registry endpoint, saving and loading through commands or Docker daemon APIs,
// transferring the build context from client to Docker daemon, and committing the
// filesystem of a container to become an image.
//
// As tar archives are used for transit, but not preserved in many situations, the
// focus of the algorithm is to ensure the integrity of the preserved filesystem,
// while maintaining a deterministic accountability. This includes neither
// constraining the ordering or manipulation of the files during the creation or
// unpacking of the archive, nor include additional metadata state about the file
// system attributes.
package tarsum

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"path"
	"strings"
)

const (
	buf8K  = 8 * 1024
	buf16K = 16 * 1024
	buf32K = 32 * 1024
)

// NewTarSum creates a new interface for calculating a fixed time checksum of a
// tar archive.
//
// This is used for calculating checksums of layers of an image, in some cases
// including the byte payload of the image's json metadata as well, and for
// calculating the checksums for buildcache.
func NewTarSum(r io.Reader, dc bool, v Version) (TarSum, error) ***REMOVED***
	return NewTarSumHash(r, dc, v, DefaultTHash)
***REMOVED***

// NewTarSumHash creates a new TarSum, providing a THash to use rather than
// the DefaultTHash.
func NewTarSumHash(r io.Reader, dc bool, v Version, tHash THash) (TarSum, error) ***REMOVED***
	headerSelector, err := getTarHeaderSelector(v)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ts := &tarSum***REMOVED***Reader: r, DisableCompression: dc, tarSumVersion: v, headerSelector: headerSelector, tHash: tHash***REMOVED***
	err = ts.initTarSum()
	return ts, err
***REMOVED***

// NewTarSumForLabel creates a new TarSum using the provided TarSum version+hash label.
func NewTarSumForLabel(r io.Reader, disableCompression bool, label string) (TarSum, error) ***REMOVED***
	parts := strings.SplitN(label, "+", 2)
	if len(parts) != 2 ***REMOVED***
		return nil, errors.New("tarsum label string should be of the form: ***REMOVED***tarsum_version***REMOVED***+***REMOVED***hash_name***REMOVED***")
	***REMOVED***

	versionName, hashName := parts[0], parts[1]

	version, ok := tarSumVersionsByName[versionName]
	if !ok ***REMOVED***
		return nil, fmt.Errorf("unknown TarSum version name: %q", versionName)
	***REMOVED***

	hashConfig, ok := standardHashConfigs[hashName]
	if !ok ***REMOVED***
		return nil, fmt.Errorf("unknown TarSum hash name: %q", hashName)
	***REMOVED***

	tHash := NewTHash(hashConfig.name, hashConfig.hash.New)

	return NewTarSumHash(r, disableCompression, version, tHash)
***REMOVED***

// TarSum is the generic interface for calculating fixed time
// checksums of a tar archive.
type TarSum interface ***REMOVED***
	io.Reader
	GetSums() FileInfoSums
	Sum([]byte) string
	Version() Version
	Hash() THash
***REMOVED***

// tarSum struct is the structure for a Version0 checksum calculation.
type tarSum struct ***REMOVED***
	io.Reader
	tarR               *tar.Reader
	tarW               *tar.Writer
	writer             writeCloseFlusher
	bufTar             *bytes.Buffer
	bufWriter          *bytes.Buffer
	bufData            []byte
	h                  hash.Hash
	tHash              THash
	sums               FileInfoSums
	fileCounter        int64
	currentFile        string
	finished           bool
	first              bool
	DisableCompression bool              // false by default. When false, the output gzip compressed.
	tarSumVersion      Version           // this field is not exported so it can not be mutated during use
	headerSelector     tarHeaderSelector // handles selecting and ordering headers for files in the archive
***REMOVED***

func (ts tarSum) Hash() THash ***REMOVED***
	return ts.tHash
***REMOVED***

func (ts tarSum) Version() Version ***REMOVED***
	return ts.tarSumVersion
***REMOVED***

// THash provides a hash.Hash type generator and its name.
type THash interface ***REMOVED***
	Hash() hash.Hash
	Name() string
***REMOVED***

// NewTHash is a convenience method for creating a THash.
func NewTHash(name string, h func() hash.Hash) THash ***REMOVED***
	return simpleTHash***REMOVED***n: name, h: h***REMOVED***
***REMOVED***

type tHashConfig struct ***REMOVED***
	name string
	hash crypto.Hash
***REMOVED***

var (
	// NOTE: DO NOT include MD5 or SHA1, which are considered insecure.
	standardHashConfigs = map[string]tHashConfig***REMOVED***
		"sha256": ***REMOVED***name: "sha256", hash: crypto.SHA256***REMOVED***,
		"sha512": ***REMOVED***name: "sha512", hash: crypto.SHA512***REMOVED***,
	***REMOVED***
)

// DefaultTHash is default TarSum hashing algorithm - "sha256".
var DefaultTHash = NewTHash("sha256", sha256.New)

type simpleTHash struct ***REMOVED***
	n string
	h func() hash.Hash
***REMOVED***

func (sth simpleTHash) Name() string    ***REMOVED*** return sth.n ***REMOVED***
func (sth simpleTHash) Hash() hash.Hash ***REMOVED*** return sth.h() ***REMOVED***

func (ts *tarSum) encodeHeader(h *tar.Header) error ***REMOVED***
	for _, elem := range ts.headerSelector.selectHeaders(h) ***REMOVED***
		if _, err := ts.h.Write([]byte(elem[0] + elem[1])); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (ts *tarSum) initTarSum() error ***REMOVED***
	ts.bufTar = bytes.NewBuffer([]byte***REMOVED******REMOVED***)
	ts.bufWriter = bytes.NewBuffer([]byte***REMOVED******REMOVED***)
	ts.tarR = tar.NewReader(ts.Reader)
	ts.tarW = tar.NewWriter(ts.bufTar)
	if !ts.DisableCompression ***REMOVED***
		ts.writer = gzip.NewWriter(ts.bufWriter)
	***REMOVED*** else ***REMOVED***
		ts.writer = &nopCloseFlusher***REMOVED***Writer: ts.bufWriter***REMOVED***
	***REMOVED***
	if ts.tHash == nil ***REMOVED***
		ts.tHash = DefaultTHash
	***REMOVED***
	ts.h = ts.tHash.Hash()
	ts.h.Reset()
	ts.first = true
	ts.sums = FileInfoSums***REMOVED******REMOVED***
	return nil
***REMOVED***

func (ts *tarSum) Read(buf []byte) (int, error) ***REMOVED***
	if ts.finished ***REMOVED***
		return ts.bufWriter.Read(buf)
	***REMOVED***
	if len(ts.bufData) < len(buf) ***REMOVED***
		switch ***REMOVED***
		case len(buf) <= buf8K:
			ts.bufData = make([]byte, buf8K)
		case len(buf) <= buf16K:
			ts.bufData = make([]byte, buf16K)
		case len(buf) <= buf32K:
			ts.bufData = make([]byte, buf32K)
		default:
			ts.bufData = make([]byte, len(buf))
		***REMOVED***
	***REMOVED***
	buf2 := ts.bufData[:len(buf)]

	n, err := ts.tarR.Read(buf2)
	if err != nil ***REMOVED***
		if err == io.EOF ***REMOVED***
			if _, err := ts.h.Write(buf2[:n]); err != nil ***REMOVED***
				return 0, err
			***REMOVED***
			if !ts.first ***REMOVED***
				ts.sums = append(ts.sums, fileInfoSum***REMOVED***name: ts.currentFile, sum: hex.EncodeToString(ts.h.Sum(nil)), pos: ts.fileCounter***REMOVED***)
				ts.fileCounter++
				ts.h.Reset()
			***REMOVED*** else ***REMOVED***
				ts.first = false
			***REMOVED***

			currentHeader, err := ts.tarR.Next()
			if err != nil ***REMOVED***
				if err == io.EOF ***REMOVED***
					if err := ts.tarW.Close(); err != nil ***REMOVED***
						return 0, err
					***REMOVED***
					if _, err := io.Copy(ts.writer, ts.bufTar); err != nil ***REMOVED***
						return 0, err
					***REMOVED***
					if err := ts.writer.Close(); err != nil ***REMOVED***
						return 0, err
					***REMOVED***
					ts.finished = true
					return n, nil
				***REMOVED***
				return n, err
			***REMOVED***
			ts.currentFile = path.Clean(currentHeader.Name)
			if err := ts.encodeHeader(currentHeader); err != nil ***REMOVED***
				return 0, err
			***REMOVED***
			if err := ts.tarW.WriteHeader(currentHeader); err != nil ***REMOVED***
				return 0, err
			***REMOVED***
			if _, err := ts.tarW.Write(buf2[:n]); err != nil ***REMOVED***
				return 0, err
			***REMOVED***
			ts.tarW.Flush()
			if _, err := io.Copy(ts.writer, ts.bufTar); err != nil ***REMOVED***
				return 0, err
			***REMOVED***
			ts.writer.Flush()

			return ts.bufWriter.Read(buf)
		***REMOVED***
		return n, err
	***REMOVED***

	// Filling the hash buffer
	if _, err = ts.h.Write(buf2[:n]); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	// Filling the tar writer
	if _, err = ts.tarW.Write(buf2[:n]); err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	ts.tarW.Flush()

	// Filling the output writer
	if _, err = io.Copy(ts.writer, ts.bufTar); err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	ts.writer.Flush()

	return ts.bufWriter.Read(buf)
***REMOVED***

func (ts *tarSum) Sum(extra []byte) string ***REMOVED***
	ts.sums.SortBySums()
	h := ts.tHash.Hash()
	if extra != nil ***REMOVED***
		h.Write(extra)
	***REMOVED***
	for _, fis := range ts.sums ***REMOVED***
		h.Write([]byte(fis.Sum()))
	***REMOVED***
	checksum := ts.Version().String() + "+" + ts.tHash.Name() + ":" + hex.EncodeToString(h.Sum(nil))
	return checksum
***REMOVED***

func (ts *tarSum) GetSums() FileInfoSums ***REMOVED***
	return ts.sums
***REMOVED***
