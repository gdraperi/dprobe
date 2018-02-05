package websocket

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"testing"
)

type nopCloser struct***REMOVED*** io.Writer ***REMOVED***

func (nopCloser) Close() error ***REMOVED*** return nil ***REMOVED***

func TestTruncWriter(t *testing.T) ***REMOVED***
	const data = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijlkmnopqrstuvwxyz987654321"
	for n := 1; n <= 10; n++ ***REMOVED***
		var b bytes.Buffer
		w := &truncWriter***REMOVED***w: nopCloser***REMOVED***&b***REMOVED******REMOVED***
		p := []byte(data)
		for len(p) > 0 ***REMOVED***
			m := len(p)
			if m > n ***REMOVED***
				m = n
			***REMOVED***
			w.Write(p[:m])
			p = p[m:]
		***REMOVED***
		if b.String() != data[:len(data)-len(w.p)] ***REMOVED***
			t.Errorf("%d: %q", n, b.String())
		***REMOVED***
	***REMOVED***
***REMOVED***

func textMessages(num int) [][]byte ***REMOVED***
	messages := make([][]byte, num)
	for i := 0; i < num; i++ ***REMOVED***
		msg := fmt.Sprintf("planet: %d, country: %d, city: %d, street: %d", i, i, i, i)
		messages[i] = []byte(msg)
	***REMOVED***
	return messages
***REMOVED***

func BenchmarkWriteNoCompression(b *testing.B) ***REMOVED***
	w := ioutil.Discard
	c := newConn(fakeNetConn***REMOVED***Reader: nil, Writer: w***REMOVED***, false, 1024, 1024)
	messages := textMessages(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		c.WriteMessage(TextMessage, messages[i%len(messages)])
	***REMOVED***
	b.ReportAllocs()
***REMOVED***

func BenchmarkWriteWithCompression(b *testing.B) ***REMOVED***
	w := ioutil.Discard
	c := newConn(fakeNetConn***REMOVED***Reader: nil, Writer: w***REMOVED***, false, 1024, 1024)
	messages := textMessages(100)
	c.enableWriteCompression = true
	c.newCompressionWriter = compressNoContextTakeover
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		c.WriteMessage(TextMessage, messages[i%len(messages)])
	***REMOVED***
	b.ReportAllocs()
***REMOVED***

func TestValidCompressionLevel(t *testing.T) ***REMOVED***
	c := newConn(fakeNetConn***REMOVED******REMOVED***, false, 1024, 1024)
	for _, level := range []int***REMOVED***minCompressionLevel - 1, maxCompressionLevel + 1***REMOVED*** ***REMOVED***
		if err := c.SetCompressionLevel(level); err == nil ***REMOVED***
			t.Errorf("no error for level %d", level)
		***REMOVED***
	***REMOVED***
	for _, level := range []int***REMOVED***minCompressionLevel, maxCompressionLevel***REMOVED*** ***REMOVED***
		if err := c.SetCompressionLevel(level); err != nil ***REMOVED***
			t.Errorf("error for level %d", level)
		***REMOVED***
	***REMOVED***
***REMOVED***
