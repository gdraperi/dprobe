// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wal

import (
	"encoding/binary"
	"hash"
	"io"
	"os"
	"sync"

	"github.com/coreos/etcd/pkg/crc"
	"github.com/coreos/etcd/pkg/ioutil"
	"github.com/coreos/etcd/wal/walpb"
)

// walPageBytes is the alignment for flushing records to the backing Writer.
// It should be a multiple of the minimum sector size so that WAL can safely
// distinguish between torn writes and ordinary data corruption.
const walPageBytes = 8 * minSectorSize

type encoder struct ***REMOVED***
	mu sync.Mutex
	bw *ioutil.PageWriter

	crc       hash.Hash32
	buf       []byte
	uint64buf []byte
***REMOVED***

func newEncoder(w io.Writer, prevCrc uint32, pageOffset int) *encoder ***REMOVED***
	return &encoder***REMOVED***
		bw:  ioutil.NewPageWriter(w, walPageBytes, pageOffset),
		crc: crc.New(prevCrc, crcTable),
		// 1MB buffer
		buf:       make([]byte, 1024*1024),
		uint64buf: make([]byte, 8),
	***REMOVED***
***REMOVED***

// newFileEncoder creates a new encoder with current file offset for the page writer.
func newFileEncoder(f *os.File, prevCrc uint32) (*encoder, error) ***REMOVED***
	offset, err := f.Seek(0, io.SeekCurrent)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return newEncoder(f, prevCrc, int(offset)), nil
***REMOVED***

func (e *encoder) encode(rec *walpb.Record) error ***REMOVED***
	e.mu.Lock()
	defer e.mu.Unlock()

	e.crc.Write(rec.Data)
	rec.Crc = e.crc.Sum32()
	var (
		data []byte
		err  error
		n    int
	)

	if rec.Size() > len(e.buf) ***REMOVED***
		data, err = rec.Marshal()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		n, err = rec.MarshalTo(e.buf)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		data = e.buf[:n]
	***REMOVED***

	lenField, padBytes := encodeFrameSize(len(data))
	if err = writeUint64(e.bw, lenField, e.uint64buf); err != nil ***REMOVED***
		return err
	***REMOVED***

	if padBytes != 0 ***REMOVED***
		data = append(data, make([]byte, padBytes)...)
	***REMOVED***
	_, err = e.bw.Write(data)
	return err
***REMOVED***

func encodeFrameSize(dataBytes int) (lenField uint64, padBytes int) ***REMOVED***
	lenField = uint64(dataBytes)
	// force 8 byte alignment so length never gets a torn write
	padBytes = (8 - (dataBytes % 8)) % 8
	if padBytes != 0 ***REMOVED***
		lenField |= uint64(0x80|padBytes) << 56
	***REMOVED***
	return
***REMOVED***

func (e *encoder) flush() error ***REMOVED***
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.bw.Flush()
***REMOVED***

func writeUint64(w io.Writer, n uint64, buf []byte) error ***REMOVED***
	// http://golang.org/src/encoding/binary/binary.go
	binary.LittleEndian.PutUint64(buf, n)
	_, err := w.Write(buf)
	return err
***REMOVED***
