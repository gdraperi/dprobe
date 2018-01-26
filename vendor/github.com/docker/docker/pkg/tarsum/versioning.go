package tarsum

import (
	"archive/tar"
	"errors"
	"io"
	"sort"
	"strconv"
	"strings"
)

// Version is used for versioning of the TarSum algorithm
// based on the prefix of the hash used
// i.e. "tarsum+sha256:e58fcf7418d4390dec8e8fb69d88c06ec07039d651fedd3aa72af9972e7d046b"
type Version int

// Prefix of "tarsum"
const (
	Version0 Version = iota
	Version1
	// VersionDev this constant will be either the latest or an unsettled next-version of the TarSum calculation
	VersionDev
)

// WriteV1Header writes a tar header to a writer in V1 tarsum format.
func WriteV1Header(h *tar.Header, w io.Writer) ***REMOVED***
	for _, elem := range v1TarHeaderSelect(h) ***REMOVED***
		w.Write([]byte(elem[0] + elem[1]))
	***REMOVED***
***REMOVED***

// VersionLabelForChecksum returns the label for the given tarsum
// checksum, i.e., everything before the first `+` character in
// the string or an empty string if no label separator is found.
func VersionLabelForChecksum(checksum string) string ***REMOVED***
	// Checksums are in the form: ***REMOVED***versionLabel***REMOVED***+***REMOVED***hashID***REMOVED***:***REMOVED***hex***REMOVED***
	sepIndex := strings.Index(checksum, "+")
	if sepIndex < 0 ***REMOVED***
		return ""
	***REMOVED***
	return checksum[:sepIndex]
***REMOVED***

// GetVersions gets a list of all known tarsum versions.
func GetVersions() []Version ***REMOVED***
	v := []Version***REMOVED******REMOVED***
	for k := range tarSumVersions ***REMOVED***
		v = append(v, k)
	***REMOVED***
	return v
***REMOVED***

var (
	tarSumVersions = map[Version]string***REMOVED***
		Version0:   "tarsum",
		Version1:   "tarsum.v1",
		VersionDev: "tarsum.dev",
	***REMOVED***
	tarSumVersionsByName = map[string]Version***REMOVED***
		"tarsum":     Version0,
		"tarsum.v1":  Version1,
		"tarsum.dev": VersionDev,
	***REMOVED***
)

func (tsv Version) String() string ***REMOVED***
	return tarSumVersions[tsv]
***REMOVED***

// GetVersionFromTarsum returns the Version from the provided string.
func GetVersionFromTarsum(tarsum string) (Version, error) ***REMOVED***
	tsv := tarsum
	if strings.Contains(tarsum, "+") ***REMOVED***
		tsv = strings.SplitN(tarsum, "+", 2)[0]
	***REMOVED***
	for v, s := range tarSumVersions ***REMOVED***
		if s == tsv ***REMOVED***
			return v, nil
		***REMOVED***
	***REMOVED***
	return -1, ErrNotVersion
***REMOVED***

// Errors that may be returned by functions in this package
var (
	ErrNotVersion            = errors.New("string does not include a TarSum Version")
	ErrVersionNotImplemented = errors.New("TarSum Version is not yet implemented")
)

// tarHeaderSelector is the interface which different versions
// of tarsum should use for selecting and ordering tar headers
// for each item in the archive.
type tarHeaderSelector interface ***REMOVED***
	selectHeaders(h *tar.Header) (orderedHeaders [][2]string)
***REMOVED***

type tarHeaderSelectFunc func(h *tar.Header) (orderedHeaders [][2]string)

func (f tarHeaderSelectFunc) selectHeaders(h *tar.Header) (orderedHeaders [][2]string) ***REMOVED***
	return f(h)
***REMOVED***

func v0TarHeaderSelect(h *tar.Header) (orderedHeaders [][2]string) ***REMOVED***
	return [][2]string***REMOVED***
		***REMOVED***"name", h.Name***REMOVED***,
		***REMOVED***"mode", strconv.FormatInt(h.Mode, 10)***REMOVED***,
		***REMOVED***"uid", strconv.Itoa(h.Uid)***REMOVED***,
		***REMOVED***"gid", strconv.Itoa(h.Gid)***REMOVED***,
		***REMOVED***"size", strconv.FormatInt(h.Size, 10)***REMOVED***,
		***REMOVED***"mtime", strconv.FormatInt(h.ModTime.UTC().Unix(), 10)***REMOVED***,
		***REMOVED***"typeflag", string([]byte***REMOVED***h.Typeflag***REMOVED***)***REMOVED***,
		***REMOVED***"linkname", h.Linkname***REMOVED***,
		***REMOVED***"uname", h.Uname***REMOVED***,
		***REMOVED***"gname", h.Gname***REMOVED***,
		***REMOVED***"devmajor", strconv.FormatInt(h.Devmajor, 10)***REMOVED***,
		***REMOVED***"devminor", strconv.FormatInt(h.Devminor, 10)***REMOVED***,
	***REMOVED***
***REMOVED***

func v1TarHeaderSelect(h *tar.Header) (orderedHeaders [][2]string) ***REMOVED***
	// Get extended attributes.
	xAttrKeys := make([]string, len(h.Xattrs))
	for k := range h.Xattrs ***REMOVED***
		xAttrKeys = append(xAttrKeys, k)
	***REMOVED***
	sort.Strings(xAttrKeys)

	// Make the slice with enough capacity to hold the 11 basic headers
	// we want from the v0 selector plus however many xattrs we have.
	orderedHeaders = make([][2]string, 0, 11+len(xAttrKeys))

	// Copy all headers from v0 excluding the 'mtime' header (the 5th element).
	v0headers := v0TarHeaderSelect(h)
	orderedHeaders = append(orderedHeaders, v0headers[0:5]...)
	orderedHeaders = append(orderedHeaders, v0headers[6:]...)

	// Finally, append the sorted xattrs.
	for _, k := range xAttrKeys ***REMOVED***
		orderedHeaders = append(orderedHeaders, [2]string***REMOVED***k, h.Xattrs[k]***REMOVED***)
	***REMOVED***

	return
***REMOVED***

var registeredHeaderSelectors = map[Version]tarHeaderSelectFunc***REMOVED***
	Version0:   v0TarHeaderSelect,
	Version1:   v1TarHeaderSelect,
	VersionDev: v1TarHeaderSelect,
***REMOVED***

func getTarHeaderSelector(v Version) (tarHeaderSelector, error) ***REMOVED***
	headerSelector, ok := registeredHeaderSelectors[v]
	if !ok ***REMOVED***
		return nil, ErrVersionNotImplemented
	***REMOVED***

	return headerSelector, nil
***REMOVED***
