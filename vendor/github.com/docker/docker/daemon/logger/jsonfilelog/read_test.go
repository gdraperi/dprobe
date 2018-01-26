package jsonfilelog

import (
	"bytes"
	"testing"
	"time"

	"github.com/docker/docker/daemon/logger"
	"github.com/gotestyourself/gotestyourself/fs"
	"github.com/stretchr/testify/require"
)

func BenchmarkJSONFileLoggerReadLogs(b *testing.B) ***REMOVED***
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

	chError := make(chan error, b.N+1)
	go func() ***REMOVED***
		for i := 0; i < b.N; i++ ***REMOVED***
			chError <- jsonlogger.Log(msg)
		***REMOVED***
		chError <- jsonlogger.Close()
	***REMOVED***()

	lw := jsonlogger.(*JSONFileLogger).ReadLogs(logger.ReadConfig***REMOVED***Follow: true***REMOVED***)
	watchClose := lw.WatchClose()
	for ***REMOVED***
		select ***REMOVED***
		case <-lw.Msg:
		case <-watchClose:
			return
		case err := <-chError:
			if err != nil ***REMOVED***
				b.Fatal(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
