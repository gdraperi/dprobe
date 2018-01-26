package multireader

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestMultiReadSeekerReadAll(t *testing.T) ***REMOVED***
	str := "hello world"
	s1 := strings.NewReader(str + " 1")
	s2 := strings.NewReader(str + " 2")
	s3 := strings.NewReader(str + " 3")
	mr := MultiReadSeeker(s1, s2, s3)

	expectedSize := int64(s1.Len() + s2.Len() + s3.Len())

	b, err := ioutil.ReadAll(mr)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	expected := "hello world 1hello world 2hello world 3"
	if string(b) != expected ***REMOVED***
		t.Fatalf("ReadAll failed, got: %q, expected %q", string(b), expected)
	***REMOVED***

	size, err := mr.Seek(0, os.SEEK_END)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if size != expectedSize ***REMOVED***
		t.Fatalf("reader size does not match, got %d, expected %d", size, expectedSize)
	***REMOVED***

	// Reset the position and read again
	pos, err := mr.Seek(0, os.SEEK_SET)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if pos != 0 ***REMOVED***
		t.Fatalf("expected position to be set to 0, got %d", pos)
	***REMOVED***

	b, err = ioutil.ReadAll(mr)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if string(b) != expected ***REMOVED***
		t.Fatalf("ReadAll failed, got: %q, expected %q", string(b), expected)
	***REMOVED***

	// The positions of some readers are not 0
	s1.Seek(0, os.SEEK_SET)
	s2.Seek(0, os.SEEK_END)
	s3.Seek(0, os.SEEK_SET)
	mr = MultiReadSeeker(s1, s2, s3)
	b, err = ioutil.ReadAll(mr)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if string(b) != expected ***REMOVED***
		t.Fatalf("ReadAll failed, got: %q, expected %q", string(b), expected)
	***REMOVED***
***REMOVED***

func TestMultiReadSeekerReadEach(t *testing.T) ***REMOVED***
	str := "hello world"
	s1 := strings.NewReader(str + " 1")
	s2 := strings.NewReader(str + " 2")
	s3 := strings.NewReader(str + " 3")
	mr := MultiReadSeeker(s1, s2, s3)

	var totalBytes int64
	for i, s := range []*strings.Reader***REMOVED***s1, s2, s3***REMOVED*** ***REMOVED***
		sLen := int64(s.Len())
		buf := make([]byte, s.Len())
		expected := []byte(fmt.Sprintf("%s %d", str, i+1))

		if _, err := mr.Read(buf); err != nil && err != io.EOF ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		if !bytes.Equal(buf, expected) ***REMOVED***
			t.Fatalf("expected %q to be %q", string(buf), string(expected))
		***REMOVED***

		pos, err := mr.Seek(0, os.SEEK_CUR)
		if err != nil ***REMOVED***
			t.Fatalf("iteration: %d, error: %v", i+1, err)
		***REMOVED***

		// check that the total bytes read is the current position of the seeker
		totalBytes += sLen
		if pos != totalBytes ***REMOVED***
			t.Fatalf("expected current position to be: %d, got: %d, iteration: %d", totalBytes, pos, i+1)
		***REMOVED***

		// This tests not only that SEEK_SET and SEEK_CUR give the same values, but that the next iteration is in the expected position as well
		newPos, err := mr.Seek(pos, os.SEEK_SET)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if newPos != pos ***REMOVED***
			t.Fatalf("expected to get same position when calling SEEK_SET with value from SEEK_CUR, cur: %d, set: %d", pos, newPos)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestMultiReadSeekerReadSpanningChunks(t *testing.T) ***REMOVED***
	str := "hello world"
	s1 := strings.NewReader(str + " 1")
	s2 := strings.NewReader(str + " 2")
	s3 := strings.NewReader(str + " 3")
	mr := MultiReadSeeker(s1, s2, s3)

	buf := make([]byte, s1.Len()+3)
	_, err := mr.Read(buf)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// expected is the contents of s1 + 3 bytes from s2, ie, the `hel` at the end of this string
	expected := "hello world 1hel"
	if string(buf) != expected ***REMOVED***
		t.Fatalf("expected %s to be %s", string(buf), expected)
	***REMOVED***
***REMOVED***

func TestMultiReadSeekerNegativeSeek(t *testing.T) ***REMOVED***
	str := "hello world"
	s1 := strings.NewReader(str + " 1")
	s2 := strings.NewReader(str + " 2")
	s3 := strings.NewReader(str + " 3")
	mr := MultiReadSeeker(s1, s2, s3)

	s1Len := s1.Len()
	s2Len := s2.Len()
	s3Len := s3.Len()

	s, err := mr.Seek(int64(-1*s3.Len()), os.SEEK_END)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if s != int64(s1Len+s2Len) ***REMOVED***
		t.Fatalf("expected %d to be %d", s, s1.Len()+s2.Len())
	***REMOVED***

	buf := make([]byte, s3Len)
	if _, err := mr.Read(buf); err != nil && err != io.EOF ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	expected := fmt.Sprintf("%s %d", str, 3)
	if string(buf) != fmt.Sprintf("%s %d", str, 3) ***REMOVED***
		t.Fatalf("expected %q to be %q", string(buf), expected)
	***REMOVED***
***REMOVED***

func TestMultiReadSeekerCurAfterSet(t *testing.T) ***REMOVED***
	str := "hello world"
	s1 := strings.NewReader(str + " 1")
	s2 := strings.NewReader(str + " 2")
	s3 := strings.NewReader(str + " 3")
	mr := MultiReadSeeker(s1, s2, s3)

	mid := int64(s1.Len() + s2.Len()/2)

	size, err := mr.Seek(mid, os.SEEK_SET)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if size != mid ***REMOVED***
		t.Fatalf("reader size does not match, got %d, expected %d", size, mid)
	***REMOVED***

	size, err = mr.Seek(3, os.SEEK_CUR)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if size != mid+3 ***REMOVED***
		t.Fatalf("reader size does not match, got %d, expected %d", size, mid+3)
	***REMOVED***
	size, err = mr.Seek(5, os.SEEK_CUR)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if size != mid+8 ***REMOVED***
		t.Fatalf("reader size does not match, got %d, expected %d", size, mid+8)
	***REMOVED***

	size, err = mr.Seek(10, os.SEEK_CUR)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if size != mid+18 ***REMOVED***
		t.Fatalf("reader size does not match, got %d, expected %d", size, mid+18)
	***REMOVED***
***REMOVED***

func TestMultiReadSeekerSmallReads(t *testing.T) ***REMOVED***
	readers := []io.ReadSeeker***REMOVED******REMOVED***
	for i := 0; i < 10; i++ ***REMOVED***
		integer := make([]byte, 4)
		binary.BigEndian.PutUint32(integer, uint32(i))
		readers = append(readers, bytes.NewReader(integer))
	***REMOVED***

	reader := MultiReadSeeker(readers...)
	for i := 0; i < 10; i++ ***REMOVED***
		var integer uint32
		if err := binary.Read(reader, binary.BigEndian, &integer); err != nil ***REMOVED***
			t.Fatalf("Read from NewMultiReadSeeker failed: %v", err)
		***REMOVED***
		if uint32(i) != integer ***REMOVED***
			t.Fatalf("Read wrong value from NewMultiReadSeeker: %d != %d", i, integer)
		***REMOVED***
	***REMOVED***
***REMOVED***
