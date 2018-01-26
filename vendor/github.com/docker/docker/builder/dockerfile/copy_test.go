package dockerfile

import (
	"net/http"
	"testing"

	"github.com/docker/docker/pkg/containerfs"
	"github.com/gotestyourself/gotestyourself/fs"
	"github.com/stretchr/testify/assert"
)

func TestIsExistingDirectory(t *testing.T) ***REMOVED***
	tmpfile := fs.NewFile(t, "file-exists-test", fs.WithContent("something"))
	defer tmpfile.Remove()
	tmpdir := fs.NewDir(t, "dir-exists-test")
	defer tmpdir.Remove()

	var testcases = []struct ***REMOVED***
		doc      string
		path     string
		expected bool
	***REMOVED******REMOVED***
		***REMOVED***
			doc:      "directory exists",
			path:     tmpdir.Path(),
			expected: true,
		***REMOVED***,
		***REMOVED***
			doc:      "path doesn't exist",
			path:     "/bogus/path/does/not/exist",
			expected: false,
		***REMOVED***,
		***REMOVED***
			doc:      "file exists",
			path:     tmpfile.Path(),
			expected: false,
		***REMOVED***,
	***REMOVED***

	for _, testcase := range testcases ***REMOVED***
		result, err := isExistingDirectory(&copyEndpoint***REMOVED***driver: containerfs.NewLocalDriver(), path: testcase.path***REMOVED***)
		if !assert.NoError(t, err) ***REMOVED***
			continue
		***REMOVED***
		assert.Equal(t, testcase.expected, result, testcase.doc)
	***REMOVED***
***REMOVED***

func TestGetFilenameForDownload(t *testing.T) ***REMOVED***
	var testcases = []struct ***REMOVED***
		path        string
		disposition string
		expected    string
	***REMOVED******REMOVED***
		***REMOVED***
			path:     "http://www.example.com/",
			expected: "",
		***REMOVED***,
		***REMOVED***
			path:     "http://www.example.com/xyz",
			expected: "xyz",
		***REMOVED***,
		***REMOVED***
			path:     "http://www.example.com/xyz.html",
			expected: "xyz.html",
		***REMOVED***,
		***REMOVED***
			path:     "http://www.example.com/xyz/",
			expected: "",
		***REMOVED***,
		***REMOVED***
			path:     "http://www.example.com/xyz/uvw",
			expected: "uvw",
		***REMOVED***,
		***REMOVED***
			path:     "http://www.example.com/xyz/uvw.html",
			expected: "uvw.html",
		***REMOVED***,
		***REMOVED***
			path:     "http://www.example.com/xyz/uvw/",
			expected: "",
		***REMOVED***,
		***REMOVED***
			path:     "/",
			expected: "",
		***REMOVED***,
		***REMOVED***
			path:     "/xyz",
			expected: "xyz",
		***REMOVED***,
		***REMOVED***
			path:     "/xyz.html",
			expected: "xyz.html",
		***REMOVED***,
		***REMOVED***
			path:     "/xyz/",
			expected: "",
		***REMOVED***,
		***REMOVED***
			path:        "/xyz/",
			disposition: "attachment; filename=xyz.html",
			expected:    "xyz.html",
		***REMOVED***,
		***REMOVED***
			disposition: "",
			expected:    "",
		***REMOVED***,
		***REMOVED***
			disposition: "attachment; filename=xyz",
			expected:    "xyz",
		***REMOVED***,
		***REMOVED***
			disposition: "attachment; filename=xyz.html",
			expected:    "xyz.html",
		***REMOVED***,
		***REMOVED***
			disposition: "attachment; filename=\"xyz\"",
			expected:    "xyz",
		***REMOVED***,
		***REMOVED***
			disposition: "attachment; filename=\"xyz.html\"",
			expected:    "xyz.html",
		***REMOVED***,
		***REMOVED***
			disposition: "attachment; filename=\"/xyz.html\"",
			expected:    "xyz.html",
		***REMOVED***,
		***REMOVED***
			disposition: "attachment; filename=\"/xyz/uvw\"",
			expected:    "uvw",
		***REMOVED***,
		***REMOVED***
			disposition: "attachment; filename=\"Naïve file.txt\"",
			expected:    "Naïve file.txt",
		***REMOVED***,
	***REMOVED***
	for _, testcase := range testcases ***REMOVED***
		resp := http.Response***REMOVED***
			Header: make(map[string][]string),
		***REMOVED***
		if testcase.disposition != "" ***REMOVED***
			resp.Header.Add("Content-Disposition", testcase.disposition)
		***REMOVED***
		filename := getFilenameForDownload(testcase.path, &resp)
		assert.Equal(t, testcase.expected, filename)
	***REMOVED***
***REMOVED***
