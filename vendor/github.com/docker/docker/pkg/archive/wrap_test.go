package archive

import (
	"archive/tar"
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateEmptyFile(t *testing.T) ***REMOVED***
	archive, err := Generate("emptyFile")
	require.NoError(t, err)
	if archive == nil ***REMOVED***
		t.Fatal("The generated archive should not be nil.")
	***REMOVED***

	expectedFiles := [][]string***REMOVED***
		***REMOVED***"emptyFile", ""***REMOVED***,
	***REMOVED***

	tr := tar.NewReader(archive)
	actualFiles := make([][]string, 0, 10)
	i := 0
	for ***REMOVED***
		hdr, err := tr.Next()
		if err == io.EOF ***REMOVED***
			break
		***REMOVED***
		require.NoError(t, err)
		buf := new(bytes.Buffer)
		buf.ReadFrom(tr)
		content := buf.String()
		actualFiles = append(actualFiles, []string***REMOVED***hdr.Name, content***REMOVED***)
		i++
	***REMOVED***
	if len(actualFiles) != len(expectedFiles) ***REMOVED***
		t.Fatalf("Number of expected file %d, got %d.", len(expectedFiles), len(actualFiles))
	***REMOVED***
	for i := 0; i < len(expectedFiles); i++ ***REMOVED***
		actual := actualFiles[i]
		expected := expectedFiles[i]
		if actual[0] != expected[0] ***REMOVED***
			t.Fatalf("Expected name '%s', Actual name '%s'", expected[0], actual[0])
		***REMOVED***
		if actual[1] != expected[1] ***REMOVED***
			t.Fatalf("Expected content '%s', Actual content '%s'", expected[1], actual[1])
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestGenerateWithContent(t *testing.T) ***REMOVED***
	archive, err := Generate("file", "content")
	require.NoError(t, err)
	if archive == nil ***REMOVED***
		t.Fatal("The generated archive should not be nil.")
	***REMOVED***

	expectedFiles := [][]string***REMOVED***
		***REMOVED***"file", "content"***REMOVED***,
	***REMOVED***

	tr := tar.NewReader(archive)
	actualFiles := make([][]string, 0, 10)
	i := 0
	for ***REMOVED***
		hdr, err := tr.Next()
		if err == io.EOF ***REMOVED***
			break
		***REMOVED***
		require.NoError(t, err)
		buf := new(bytes.Buffer)
		buf.ReadFrom(tr)
		content := buf.String()
		actualFiles = append(actualFiles, []string***REMOVED***hdr.Name, content***REMOVED***)
		i++
	***REMOVED***
	if len(actualFiles) != len(expectedFiles) ***REMOVED***
		t.Fatalf("Number of expected file %d, got %d.", len(expectedFiles), len(actualFiles))
	***REMOVED***
	for i := 0; i < len(expectedFiles); i++ ***REMOVED***
		actual := actualFiles[i]
		expected := expectedFiles[i]
		if actual[0] != expected[0] ***REMOVED***
			t.Fatalf("Expected name '%s', Actual name '%s'", expected[0], actual[0])
		***REMOVED***
		if actual[1] != expected[1] ***REMOVED***
			t.Fatalf("Expected content '%s', Actual content '%s'", expected[1], actual[1])
		***REMOVED***
	***REMOVED***
***REMOVED***
