package archive

import (
	"archive/tar"
	"bytes"
	"io"
)

// Generate generates a new archive from the content provided
// as input.
//
// `files` is a sequence of path/content pairs. A new file is
// added to the archive for each pair.
// If the last pair is incomplete, the file is created with an
// empty content. For example:
//
// Generate("foo.txt", "hello world", "emptyfile")
//
// The above call will return an archive with 2 files:
//  * ./foo.txt with content "hello world"
//  * ./empty with empty content
//
// FIXME: stream content instead of buffering
// FIXME: specify permissions and other archive metadata
func Generate(input ...string) (io.Reader, error) ***REMOVED***
	files := parseStringPairs(input...)
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	for _, file := range files ***REMOVED***
		name, content := file[0], file[1]
		hdr := &tar.Header***REMOVED***
			Name: name,
			Size: int64(len(content)),
		***REMOVED***
		if err := tw.WriteHeader(hdr); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if _, err := tw.Write([]byte(content)); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	if err := tw.Close(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return buf, nil
***REMOVED***

func parseStringPairs(input ...string) (output [][2]string) ***REMOVED***
	output = make([][2]string, 0, len(input)/2+1)
	for i := 0; i < len(input); i += 2 ***REMOVED***
		var pair [2]string
		pair[0] = input[i]
		if i+1 < len(input) ***REMOVED***
			pair[1] = input[i+1]
		***REMOVED***
		output = append(output, pair)
	***REMOVED***
	return
***REMOVED***
