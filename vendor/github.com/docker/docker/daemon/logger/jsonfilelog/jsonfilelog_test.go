package jsonfilelog

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/daemon/logger/jsonfilelog/jsonlog"
	"github.com/gotestyourself/gotestyourself/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONFileLogger(t *testing.T) ***REMOVED***
	cid := "a7317399f3f857173c6179d44823594f8294678dea9999662e5c625b5a1c7657"
	tmp, err := ioutil.TempDir("", "docker-logger-")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmp)
	filename := filepath.Join(tmp, "container.log")
	l, err := New(logger.Info***REMOVED***
		ContainerID: cid,
		LogPath:     filename,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer l.Close()

	if err := l.Log(&logger.Message***REMOVED***Line: []byte("line1"), Source: "src1"***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := l.Log(&logger.Message***REMOVED***Line: []byte("line2"), Source: "src2"***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := l.Log(&logger.Message***REMOVED***Line: []byte("line3"), Source: "src3"***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	res, err := ioutil.ReadFile(filename)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	expected := `***REMOVED***"log":"line1\n","stream":"src1","time":"0001-01-01T00:00:00Z"***REMOVED***
***REMOVED***"log":"line2\n","stream":"src2","time":"0001-01-01T00:00:00Z"***REMOVED***
***REMOVED***"log":"line3\n","stream":"src3","time":"0001-01-01T00:00:00Z"***REMOVED***
`

	if string(res) != expected ***REMOVED***
		t.Fatalf("Wrong log content: %q, expected %q", res, expected)
	***REMOVED***
***REMOVED***

func TestJSONFileLoggerWithTags(t *testing.T) ***REMOVED***
	cid := "a7317399f3f857173c6179d44823594f8294678dea9999662e5c625b5a1c7657"
	cname := "test-container"
	tmp, err := ioutil.TempDir("", "docker-logger-")

	require.NoError(t, err)

	defer os.RemoveAll(tmp)
	filename := filepath.Join(tmp, "container.log")
	l, err := New(logger.Info***REMOVED***
		Config: map[string]string***REMOVED***
			"tag": "***REMOVED******REMOVED***.ID***REMOVED******REMOVED***/***REMOVED******REMOVED***.Name***REMOVED******REMOVED***", // first 12 characters of ContainerID and full ContainerName
		***REMOVED***,
		ContainerID:   cid,
		ContainerName: cname,
		LogPath:       filename,
	***REMOVED***)

	require.NoError(t, err)
	defer l.Close()

	err = l.Log(&logger.Message***REMOVED***Line: []byte("line1"), Source: "src1"***REMOVED***)
	require.NoError(t, err)

	err = l.Log(&logger.Message***REMOVED***Line: []byte("line2"), Source: "src2"***REMOVED***)
	require.NoError(t, err)

	err = l.Log(&logger.Message***REMOVED***Line: []byte("line3"), Source: "src3"***REMOVED***)
	require.NoError(t, err)

	res, err := ioutil.ReadFile(filename)
	require.NoError(t, err)

	expected := `***REMOVED***"log":"line1\n","stream":"src1","attrs":***REMOVED***"tag":"a7317399f3f8/test-container"***REMOVED***,"time":"0001-01-01T00:00:00Z"***REMOVED***
***REMOVED***"log":"line2\n","stream":"src2","attrs":***REMOVED***"tag":"a7317399f3f8/test-container"***REMOVED***,"time":"0001-01-01T00:00:00Z"***REMOVED***
***REMOVED***"log":"line3\n","stream":"src3","attrs":***REMOVED***"tag":"a7317399f3f8/test-container"***REMOVED***,"time":"0001-01-01T00:00:00Z"***REMOVED***
`
	assert.Equal(t, expected, string(res))
***REMOVED***

func BenchmarkJSONFileLoggerLog(b *testing.B) ***REMOVED***
	tmp := fs.NewDir(b, "bench-jsonfilelog")
	defer tmp.Remove()

	jsonlogger, err := New(logger.Info***REMOVED***
		ContainerID: "a7317399f3f857173c6179d44823594f8294678dea9999662e5c625b5a1c7657",
		LogPath:     tmp.Join("container.log"),
		Config: map[string]string***REMOVED***
			"labels": "first,second",
		***REMOVED***,
		ContainerLabels: map[string]string***REMOVED***
			"first":  "label_value",
			"second": "label_foo",
		***REMOVED***,
	***REMOVED***)
	require.NoError(b, err)
	defer jsonlogger.Close()

	msg := &logger.Message***REMOVED***
		Line:      []byte("Line that thinks that it is log line from docker\n"),
		Source:    "stderr",
		Timestamp: time.Now().UTC(),
	***REMOVED***

	buf := bytes.NewBuffer(nil)
	require.NoError(b, marshalMessage(msg, nil, buf))
	b.SetBytes(int64(buf.Len()))

	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if err := jsonlogger.Log(msg); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestJSONFileLoggerWithOpts(t *testing.T) ***REMOVED***
	cid := "a7317399f3f857173c6179d44823594f8294678dea9999662e5c625b5a1c7657"
	tmp, err := ioutil.TempDir("", "docker-logger-")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmp)
	filename := filepath.Join(tmp, "container.log")
	config := map[string]string***REMOVED***"max-file": "2", "max-size": "1k"***REMOVED***
	l, err := New(logger.Info***REMOVED***
		ContainerID: cid,
		LogPath:     filename,
		Config:      config,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer l.Close()
	for i := 0; i < 20; i++ ***REMOVED***
		if err := l.Log(&logger.Message***REMOVED***Line: []byte("line" + strconv.Itoa(i)), Source: "src1"***REMOVED***); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***
	res, err := ioutil.ReadFile(filename)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	penUlt, err := ioutil.ReadFile(filename + ".1")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	expectedPenultimate := `***REMOVED***"log":"line0\n","stream":"src1","time":"0001-01-01T00:00:00Z"***REMOVED***
***REMOVED***"log":"line1\n","stream":"src1","time":"0001-01-01T00:00:00Z"***REMOVED***
***REMOVED***"log":"line2\n","stream":"src1","time":"0001-01-01T00:00:00Z"***REMOVED***
***REMOVED***"log":"line3\n","stream":"src1","time":"0001-01-01T00:00:00Z"***REMOVED***
***REMOVED***"log":"line4\n","stream":"src1","time":"0001-01-01T00:00:00Z"***REMOVED***
***REMOVED***"log":"line5\n","stream":"src1","time":"0001-01-01T00:00:00Z"***REMOVED***
***REMOVED***"log":"line6\n","stream":"src1","time":"0001-01-01T00:00:00Z"***REMOVED***
***REMOVED***"log":"line7\n","stream":"src1","time":"0001-01-01T00:00:00Z"***REMOVED***
***REMOVED***"log":"line8\n","stream":"src1","time":"0001-01-01T00:00:00Z"***REMOVED***
***REMOVED***"log":"line9\n","stream":"src1","time":"0001-01-01T00:00:00Z"***REMOVED***
***REMOVED***"log":"line10\n","stream":"src1","time":"0001-01-01T00:00:00Z"***REMOVED***
***REMOVED***"log":"line11\n","stream":"src1","time":"0001-01-01T00:00:00Z"***REMOVED***
***REMOVED***"log":"line12\n","stream":"src1","time":"0001-01-01T00:00:00Z"***REMOVED***
***REMOVED***"log":"line13\n","stream":"src1","time":"0001-01-01T00:00:00Z"***REMOVED***
***REMOVED***"log":"line14\n","stream":"src1","time":"0001-01-01T00:00:00Z"***REMOVED***
***REMOVED***"log":"line15\n","stream":"src1","time":"0001-01-01T00:00:00Z"***REMOVED***
`
	expected := `***REMOVED***"log":"line16\n","stream":"src1","time":"0001-01-01T00:00:00Z"***REMOVED***
***REMOVED***"log":"line17\n","stream":"src1","time":"0001-01-01T00:00:00Z"***REMOVED***
***REMOVED***"log":"line18\n","stream":"src1","time":"0001-01-01T00:00:00Z"***REMOVED***
***REMOVED***"log":"line19\n","stream":"src1","time":"0001-01-01T00:00:00Z"***REMOVED***
`

	if string(res) != expected ***REMOVED***
		t.Fatalf("Wrong log content: %q, expected %q", res, expected)
	***REMOVED***
	if string(penUlt) != expectedPenultimate ***REMOVED***
		t.Fatalf("Wrong log content: %q, expected %q", penUlt, expectedPenultimate)
	***REMOVED***

***REMOVED***

func TestJSONFileLoggerWithLabelsEnv(t *testing.T) ***REMOVED***
	cid := "a7317399f3f857173c6179d44823594f8294678dea9999662e5c625b5a1c7657"
	tmp, err := ioutil.TempDir("", "docker-logger-")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmp)
	filename := filepath.Join(tmp, "container.log")
	config := map[string]string***REMOVED***"labels": "rack,dc", "env": "environ,debug,ssl", "env-regex": "^dc"***REMOVED***
	l, err := New(logger.Info***REMOVED***
		ContainerID:     cid,
		LogPath:         filename,
		Config:          config,
		ContainerLabels: map[string]string***REMOVED***"rack": "101", "dc": "lhr"***REMOVED***,
		ContainerEnv:    []string***REMOVED***"environ=production", "debug=false", "port=10001", "ssl=true", "dc_region=west"***REMOVED***,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer l.Close()
	if err := l.Log(&logger.Message***REMOVED***Line: []byte("line"), Source: "src1"***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	res, err := ioutil.ReadFile(filename)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	var jsonLog jsonlog.JSONLogs
	if err := json.Unmarshal(res, &jsonLog); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	extra := make(map[string]string)
	if err := json.Unmarshal(jsonLog.RawAttrs, &extra); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	expected := map[string]string***REMOVED***
		"rack":      "101",
		"dc":        "lhr",
		"environ":   "production",
		"debug":     "false",
		"ssl":       "true",
		"dc_region": "west",
	***REMOVED***
	if !reflect.DeepEqual(extra, expected) ***REMOVED***
		t.Fatalf("Wrong log attrs: %q, expected %q", extra, expected)
	***REMOVED***
***REMOVED***
