package ioutils

import (
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
	"testing"
	"time"
)

func TestBytesPipeRead(t *testing.T) ***REMOVED***
	buf := NewBytesPipe()
	buf.Write([]byte("12"))
	buf.Write([]byte("34"))
	buf.Write([]byte("56"))
	buf.Write([]byte("78"))
	buf.Write([]byte("90"))
	rd := make([]byte, 4)
	n, err := buf.Read(rd)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if n != 4 ***REMOVED***
		t.Fatalf("Wrong number of bytes read: %d, should be %d", n, 4)
	***REMOVED***
	if string(rd) != "1234" ***REMOVED***
		t.Fatalf("Read %s, but must be %s", rd, "1234")
	***REMOVED***
	n, err = buf.Read(rd)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if n != 4 ***REMOVED***
		t.Fatalf("Wrong number of bytes read: %d, should be %d", n, 4)
	***REMOVED***
	if string(rd) != "5678" ***REMOVED***
		t.Fatalf("Read %s, but must be %s", rd, "5679")
	***REMOVED***
	n, err = buf.Read(rd)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if n != 2 ***REMOVED***
		t.Fatalf("Wrong number of bytes read: %d, should be %d", n, 2)
	***REMOVED***
	if string(rd[:n]) != "90" ***REMOVED***
		t.Fatalf("Read %s, but must be %s", rd, "90")
	***REMOVED***
***REMOVED***

func TestBytesPipeWrite(t *testing.T) ***REMOVED***
	buf := NewBytesPipe()
	buf.Write([]byte("12"))
	buf.Write([]byte("34"))
	buf.Write([]byte("56"))
	buf.Write([]byte("78"))
	buf.Write([]byte("90"))
	if buf.buf[0].String() != "1234567890" ***REMOVED***
		t.Fatalf("Buffer %q, must be %q", buf.buf[0].String(), "1234567890")
	***REMOVED***
***REMOVED***

// Write and read in different speeds/chunk sizes and check valid data is read.
func TestBytesPipeWriteRandomChunks(t *testing.T) ***REMOVED***
	cases := []struct***REMOVED*** iterations, writesPerLoop, readsPerLoop int ***REMOVED******REMOVED***
		***REMOVED***100, 10, 1***REMOVED***,
		***REMOVED***1000, 10, 5***REMOVED***,
		***REMOVED***1000, 100, 0***REMOVED***,
		***REMOVED***1000, 5, 6***REMOVED***,
		***REMOVED***10000, 50, 25***REMOVED***,
	***REMOVED***

	testMessage := []byte("this is a random string for testing")
	// random slice sizes to read and write
	writeChunks := []int***REMOVED***25, 35, 15, 20***REMOVED***
	readChunks := []int***REMOVED***5, 45, 20, 25***REMOVED***

	for _, c := range cases ***REMOVED***
		// first pass: write directly to hash
		hash := sha1.New()
		for i := 0; i < c.iterations*c.writesPerLoop; i++ ***REMOVED***
			if _, err := hash.Write(testMessage[:writeChunks[i%len(writeChunks)]]); err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
		***REMOVED***
		expected := hex.EncodeToString(hash.Sum(nil))

		// write/read through buffer
		buf := NewBytesPipe()
		hash.Reset()

		done := make(chan struct***REMOVED******REMOVED***)

		go func() ***REMOVED***
			// random delay before read starts
			<-time.After(time.Duration(rand.Intn(10)) * time.Millisecond)
			for i := 0; ; i++ ***REMOVED***
				p := make([]byte, readChunks[(c.iterations*c.readsPerLoop+i)%len(readChunks)])
				n, _ := buf.Read(p)
				if n == 0 ***REMOVED***
					break
				***REMOVED***
				hash.Write(p[:n])
			***REMOVED***

			close(done)
		***REMOVED***()

		for i := 0; i < c.iterations; i++ ***REMOVED***
			for w := 0; w < c.writesPerLoop; w++ ***REMOVED***
				buf.Write(testMessage[:writeChunks[(i*c.writesPerLoop+w)%len(writeChunks)]])
			***REMOVED***
		***REMOVED***
		buf.Close()
		<-done

		actual := hex.EncodeToString(hash.Sum(nil))

		if expected != actual ***REMOVED***
			t.Fatalf("BytesPipe returned invalid data. Expected checksum %v, got %v", expected, actual)
		***REMOVED***

	***REMOVED***
***REMOVED***

func BenchmarkBytesPipeWrite(b *testing.B) ***REMOVED***
	testData := []byte("pretty short line, because why not?")
	for i := 0; i < b.N; i++ ***REMOVED***
		readBuf := make([]byte, 1024)
		buf := NewBytesPipe()
		go func() ***REMOVED***
			var err error
			for err == nil ***REMOVED***
				_, err = buf.Read(readBuf)
			***REMOVED***
		***REMOVED***()
		for j := 0; j < 1000; j++ ***REMOVED***
			buf.Write(testData)
		***REMOVED***
		buf.Close()
	***REMOVED***
***REMOVED***

func BenchmarkBytesPipeRead(b *testing.B) ***REMOVED***
	rd := make([]byte, 512)
	for i := 0; i < b.N; i++ ***REMOVED***
		b.StopTimer()
		buf := NewBytesPipe()
		for j := 0; j < 500; j++ ***REMOVED***
			buf.Write(make([]byte, 1024))
		***REMOVED***
		b.StartTimer()
		for j := 0; j < 1000; j++ ***REMOVED***
			if n, _ := buf.Read(rd); n != 512 ***REMOVED***
				b.Fatalf("Wrong number of bytes: %d", n)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
