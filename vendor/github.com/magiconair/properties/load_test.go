// Copyright 2017 Frank Schroeder. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package properties

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestLoadFailsWithNotExistingFile(t *testing.T) ***REMOVED***
	_, err := LoadFile("doesnotexist.properties", ISO_8859_1)
	assert.Equal(t, err != nil, true, "")
	assert.Matches(t, err.Error(), "open.*no such file or directory")
***REMOVED***

func TestLoadFilesFailsOnNotExistingFile(t *testing.T) ***REMOVED***
	_, err := LoadFile("doesnotexist.properties", ISO_8859_1)
	assert.Equal(t, err != nil, true, "")
	assert.Matches(t, err.Error(), "open.*no such file or directory")
***REMOVED***

func TestLoadFilesDoesNotFailOnNotExistingFileAndIgnoreMissing(t *testing.T) ***REMOVED***
	p, err := LoadFiles([]string***REMOVED***"doesnotexist.properties"***REMOVED***, ISO_8859_1, true)
	assert.Equal(t, err, nil)
	assert.Equal(t, p.Len(), 0)
***REMOVED***

func TestLoadString(t *testing.T) ***REMOVED***
	x := "key=äüö"
	p1 := MustLoadString(x)
	p2 := must(Load([]byte(x), UTF8))
	assert.Equal(t, p1, p2)
***REMOVED***

func TestLoadMap(t *testing.T) ***REMOVED***
	// LoadMap does not guarantee the same import order
	// of keys every time since map access is randomized.
	// Therefore, we need to compare the generated maps.
	m := map[string]string***REMOVED***"key": "value", "abc": "def"***REMOVED***
	assert.Equal(t, LoadMap(m).Map(), m)
***REMOVED***

func TestLoadFile(t *testing.T) ***REMOVED***
	tf := make(tempFiles, 0)
	defer tf.removeAll()

	filename := tf.makeFile("key=value")
	p := MustLoadFile(filename, ISO_8859_1)

	assert.Equal(t, p.Len(), 1)
	assertKeyValues(t, "", p, "key", "value")
***REMOVED***

func TestLoadFiles(t *testing.T) ***REMOVED***
	tf := make(tempFiles, 0)
	defer tf.removeAll()

	filename := tf.makeFile("key=value")
	filename2 := tf.makeFile("key2=value2")
	p := MustLoadFiles([]string***REMOVED***filename, filename2***REMOVED***, ISO_8859_1, false)
	assertKeyValues(t, "", p, "key", "value", "key2", "value2")
***REMOVED***

func TestLoadExpandedFile(t *testing.T) ***REMOVED***
	tf := make(tempFiles, 0)
	defer tf.removeAll()

	if err := os.Setenv("_VARX", "some-value"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	filename := tf.makeFilePrefix(os.Getenv("_VARX"), "key=value")
	filename = strings.Replace(filename, os.Getenv("_VARX"), "$***REMOVED***_VARX***REMOVED***", -1)
	p := MustLoadFile(filename, ISO_8859_1)
	assertKeyValues(t, "", p, "key", "value")
***REMOVED***

func TestLoadFilesAndIgnoreMissing(t *testing.T) ***REMOVED***
	tf := make(tempFiles, 0)
	defer tf.removeAll()

	filename := tf.makeFile("key=value")
	filename2 := tf.makeFile("key2=value2")
	p := MustLoadFiles([]string***REMOVED***filename, filename + "foo", filename2, filename2 + "foo"***REMOVED***, ISO_8859_1, true)
	assertKeyValues(t, "", p, "key", "value", "key2", "value2")
***REMOVED***

func TestLoadURL(t *testing.T) ***REMOVED***
	srv := testServer()
	defer srv.Close()
	p := MustLoadURL(srv.URL + "/a")
	assertKeyValues(t, "", p, "key", "value")
***REMOVED***

func TestLoadURLs(t *testing.T) ***REMOVED***
	srv := testServer()
	defer srv.Close()
	p := MustLoadURLs([]string***REMOVED***srv.URL + "/a", srv.URL + "/b"***REMOVED***, false)
	assertKeyValues(t, "", p, "key", "value", "key2", "value2")
***REMOVED***

func TestLoadURLsAndFailMissing(t *testing.T) ***REMOVED***
	srv := testServer()
	defer srv.Close()
	p, err := LoadURLs([]string***REMOVED***srv.URL + "/a", srv.URL + "/c"***REMOVED***, false)
	assert.Equal(t, p, (*Properties)(nil))
	assert.Matches(t, err.Error(), ".*returned 404.*")
***REMOVED***

func TestLoadURLsAndIgnoreMissing(t *testing.T) ***REMOVED***
	srv := testServer()
	defer srv.Close()
	p := MustLoadURLs([]string***REMOVED***srv.URL + "/a", srv.URL + "/b", srv.URL + "/c"***REMOVED***, true)
	assertKeyValues(t, "", p, "key", "value", "key2", "value2")
***REMOVED***

func TestLoadURLEncoding(t *testing.T) ***REMOVED***
	srv := testServer()
	defer srv.Close()

	uris := []string***REMOVED***"/none", "/utf8", "/plain", "/latin1", "/iso88591"***REMOVED***
	for i, uri := range uris ***REMOVED***
		p := MustLoadURL(srv.URL + uri)
		assert.Equal(t, p.GetString("key", ""), "äöü", fmt.Sprintf("%d", i))
	***REMOVED***
***REMOVED***

func TestLoadURLFailInvalidEncoding(t *testing.T) ***REMOVED***
	srv := testServer()
	defer srv.Close()

	p, err := LoadURL(srv.URL + "/json")
	assert.Equal(t, p, (*Properties)(nil))
	assert.Matches(t, err.Error(), ".*invalid content type.*")
***REMOVED***

func TestLoadAll(t *testing.T) ***REMOVED***
	tf := make(tempFiles, 0)
	defer tf.removeAll()

	filename := tf.makeFile("key=value")
	filename2 := tf.makeFile("key2=value3")
	filename3 := tf.makeFile("key=value4")
	srv := testServer()
	defer srv.Close()
	p := MustLoadAll([]string***REMOVED***filename, filename2, srv.URL + "/a", srv.URL + "/b", filename3***REMOVED***, UTF8, false)
	assertKeyValues(t, "", p, "key", "value4", "key2", "value2")
***REMOVED***

type tempFiles []string

func (tf *tempFiles) removeAll() ***REMOVED***
	for _, path := range *tf ***REMOVED***
		err := os.Remove(path)
		if err != nil ***REMOVED***
			fmt.Printf("os.Remove: %v", err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (tf *tempFiles) makeFile(data string) string ***REMOVED***
	return tf.makeFilePrefix("properties", data)
***REMOVED***

func (tf *tempFiles) makeFilePrefix(prefix, data string) string ***REMOVED***
	f, err := ioutil.TempFile("", prefix)
	if err != nil ***REMOVED***
		panic("ioutil.TempFile: " + err.Error())
	***REMOVED***

	// remember the temp file so that we can remove it later
	*tf = append(*tf, f.Name())

	n, err := fmt.Fprint(f, data)
	if err != nil ***REMOVED***
		panic("fmt.Fprintln: " + err.Error())
	***REMOVED***
	if n != len(data) ***REMOVED***
		panic(fmt.Sprintf("Data size mismatch. expected=%d wrote=%d\n", len(data), n))
	***REMOVED***

	err = f.Close()
	if err != nil ***REMOVED***
		panic("f.Close: " + err.Error())
	***REMOVED***

	return f.Name()
***REMOVED***

func testServer() *httptest.Server ***REMOVED***
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		send := func(data []byte, contentType string) ***REMOVED***
			w.Header().Set("Content-Type", contentType)
			if _, err := w.Write(data); err != nil ***REMOVED***
				panic(err)
			***REMOVED***
		***REMOVED***

		utf8 := []byte("key=äöü")
		iso88591 := []byte***REMOVED***0x6b, 0x65, 0x79, 0x3d, 0xe4, 0xf6, 0xfc***REMOVED*** // key=äöü

		switch r.RequestURI ***REMOVED***
		case "/a":
			send([]byte("key=value"), "")
		case "/b":
			send([]byte("key2=value2"), "")
		case "/none":
			send(utf8, "")
		case "/utf8":
			send(utf8, "text/plain; charset=utf-8")
		case "/json":
			send(utf8, "application/json; charset=utf-8")
		case "/plain":
			send(iso88591, "text/plain")
		case "/latin1":
			send(iso88591, "text/plain; charset=latin1")
		case "/iso88591":
			send(iso88591, "text/plain; charset=iso-8859-1")
		default:
			w.WriteHeader(404)
		***REMOVED***
	***REMOVED***))
***REMOVED***
