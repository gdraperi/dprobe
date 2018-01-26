package logger

import (
	"encoding/binary"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/docker/docker/api/types/plugins/logdriver"
	protoio "github.com/gogo/protobuf/io"
	"github.com/stretchr/testify/assert"
)

// mockLoggingPlugin implements the loggingPlugin interface for testing purposes
// it only supports a single log stream
type mockLoggingPlugin struct ***REMOVED***
	inStream io.ReadCloser
	f        *os.File
	closed   chan struct***REMOVED******REMOVED***
	t        *testing.T
***REMOVED***

func (l *mockLoggingPlugin) StartLogging(file string, info Info) error ***REMOVED***
	go func() ***REMOVED***
		io.Copy(l.f, l.inStream)
		close(l.closed)
	***REMOVED***()
	return nil
***REMOVED***

func (l *mockLoggingPlugin) StopLogging(file string) error ***REMOVED***
	l.inStream.Close()
	l.f.Close()
	os.Remove(l.f.Name())
	return nil
***REMOVED***

func (l *mockLoggingPlugin) Capabilities() (cap Capability, err error) ***REMOVED***
	return Capability***REMOVED***ReadLogs: true***REMOVED***, nil
***REMOVED***

func (l *mockLoggingPlugin) ReadLogs(info Info, config ReadConfig) (io.ReadCloser, error) ***REMOVED***
	r, w := io.Pipe()
	f, err := os.Open(l.f.Name())
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	go func() ***REMOVED***
		defer f.Close()
		dec := protoio.NewUint32DelimitedReader(f, binary.BigEndian, 1e6)
		enc := logdriver.NewLogEntryEncoder(w)

		for ***REMOVED***
			select ***REMOVED***
			case <-l.closed:
				w.Close()
				return
			default:
			***REMOVED***

			var msg logdriver.LogEntry
			if err := dec.ReadMsg(&msg); err != nil ***REMOVED***
				if err == io.EOF ***REMOVED***
					if !config.Follow ***REMOVED***
						w.Close()
						return
					***REMOVED***
					dec = protoio.NewUint32DelimitedReader(f, binary.BigEndian, 1e6)
					continue
				***REMOVED***

				l.t.Fatal(err)
				continue
			***REMOVED***

			if err := enc.Encode(&msg); err != nil ***REMOVED***
				w.CloseWithError(err)
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return r, nil
***REMOVED***

func newMockPluginAdapter(t *testing.T) Logger ***REMOVED***
	r, w := io.Pipe()
	f, err := ioutil.TempFile("", "mock-plugin-adapter")
	assert.NoError(t, err)

	enc := logdriver.NewLogEntryEncoder(w)
	a := &pluginAdapterWithRead***REMOVED***
		&pluginAdapter***REMOVED***
			plugin: &mockLoggingPlugin***REMOVED***
				inStream: r,
				f:        f,
				closed:   make(chan struct***REMOVED******REMOVED***),
				t:        t,
			***REMOVED***,
			stream: w,
			enc:    enc,
		***REMOVED***,
	***REMOVED***
	a.plugin.StartLogging("", Info***REMOVED******REMOVED***)
	return a
***REMOVED***

func TestAdapterReadLogs(t *testing.T) ***REMOVED***
	l := newMockPluginAdapter(t)

	testMsg := []Message***REMOVED***
		***REMOVED***Line: []byte("Are you the keymaker?"), Timestamp: time.Now()***REMOVED***,
		***REMOVED***Line: []byte("Follow the white rabbit"), Timestamp: time.Now()***REMOVED***,
	***REMOVED***
	for _, msg := range testMsg ***REMOVED***
		m := msg.copy()
		assert.NoError(t, l.Log(m))
	***REMOVED***

	lr, ok := l.(LogReader)
	assert.NotNil(t, ok)

	lw := lr.ReadLogs(ReadConfig***REMOVED******REMOVED***)

	for _, x := range testMsg ***REMOVED***
		select ***REMOVED***
		case msg := <-lw.Msg:
			testMessageEqual(t, &x, msg)
		case <-time.After(10 * time.Second):
			t.Fatal("timeout reading logs")
		***REMOVED***
	***REMOVED***

	select ***REMOVED***
	case _, ok := <-lw.Msg:
		assert.False(t, ok, "expected message channel to be closed")
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for message channel to close")

	***REMOVED***
	lw.Close()

	lw = lr.ReadLogs(ReadConfig***REMOVED***Follow: true***REMOVED***)
	for _, x := range testMsg ***REMOVED***
		select ***REMOVED***
		case msg := <-lw.Msg:
			testMessageEqual(t, &x, msg)
		case <-time.After(10 * time.Second):
			t.Fatal("timeout reading logs")
		***REMOVED***
	***REMOVED***

	x := Message***REMOVED***Line: []byte("Too infinity and beyond!"), Timestamp: time.Now()***REMOVED***
	assert.NoError(t, l.Log(x.copy()))

	select ***REMOVED***
	case msg, ok := <-lw.Msg:
		assert.NotNil(t, ok, "message channel unexpectedly closed")
		testMessageEqual(t, &x, msg)
	case <-time.After(10 * time.Second):
		t.Fatal("timeout reading logs")
	***REMOVED***

	l.Close()
	select ***REMOVED***
	case msg, ok := <-lw.Msg:
		assert.False(t, ok, "expected message channel to be closed")
		assert.Nil(t, msg)
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for logger to close")
	***REMOVED***
***REMOVED***

func testMessageEqual(t *testing.T, a, b *Message) ***REMOVED***
	assert.Equal(t, a.Line, b.Line)
	assert.Equal(t, a.Timestamp.UnixNano(), b.Timestamp.UnixNano())
	assert.Equal(t, a.Source, b.Source)
***REMOVED***
