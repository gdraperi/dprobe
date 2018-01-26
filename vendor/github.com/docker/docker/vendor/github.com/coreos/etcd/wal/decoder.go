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
	"bufio"
	"encoding/binary"
	"hash"
	"io"
	"sync"

	"github.com/coreos/etcd/pkg/crc"
	"github.com/coreos/etcd/pkg/pbutil"
	"github.com/coreos/etcd/raft/raftpb"
	"github.com/coreos/etcd/wal/walpb"
)

const minSectorSize = 512

type decoder struct ***REMOVED***
	mu  sync.Mutex
	brs []*bufio.Reader

	// lastValidOff file offset following the last valid decoded record
	lastValidOff int64
	crc          hash.Hash32
***REMOVED***

func newDecoder(r ...io.Reader) *decoder ***REMOVED***
	readers := make([]*bufio.Reader, len(r))
	for i := range r ***REMOVED***
		readers[i] = bufio.NewReader(r[i])
	***REMOVED***
	return &decoder***REMOVED***
		brs: readers,
		crc: crc.New(0, crcTable),
	***REMOVED***
***REMOVED***

func (d *decoder) decode(rec *walpb.Record) error ***REMOVED***
	rec.Reset()
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.decodeRecord(rec)
***REMOVED***

func (d *decoder) decodeRecord(rec *walpb.Record) error ***REMOVED***
	if len(d.brs) == 0 ***REMOVED***
		return io.EOF
	***REMOVED***

	l, err := readInt64(d.brs[0])
	if err == io.EOF || (err == nil && l == 0) ***REMOVED***
		// hit end of file or preallocated space
		d.brs = d.brs[1:]
		if len(d.brs) == 0 ***REMOVED***
			return io.EOF
		***REMOVED***
		d.lastValidOff = 0
		return d.decodeRecord(rec)
	***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	recBytes, padBytes := decodeFrameSize(l)

	data := make([]byte, recBytes+padBytes)
	if _, err = io.ReadFull(d.brs[0], data); err != nil ***REMOVED***
		// ReadFull returns io.EOF only if no bytes were read
		// the decoder should treat this as an ErrUnexpectedEOF instead.
		if err == io.EOF ***REMOVED***
			err = io.ErrUnexpectedEOF
		***REMOVED***
		return err
	***REMOVED***
	if err := rec.Unmarshal(data[:recBytes]); err != nil ***REMOVED***
		if d.isTornEntry(data) ***REMOVED***
			return io.ErrUnexpectedEOF
		***REMOVED***
		return err
	***REMOVED***

	// skip crc checking if the record type is crcType
	if rec.Type != crcType ***REMOVED***
		d.crc.Write(rec.Data)
		if err := rec.Validate(d.crc.Sum32()); err != nil ***REMOVED***
			if d.isTornEntry(data) ***REMOVED***
				return io.ErrUnexpectedEOF
			***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	// record decoded as valid; point last valid offset to end of record
	d.lastValidOff += recBytes + padBytes + 8
	return nil
***REMOVED***

func decodeFrameSize(lenField int64) (recBytes int64, padBytes int64) ***REMOVED***
	// the record size is stored in the lower 56 bits of the 64-bit length
	recBytes = int64(uint64(lenField) & ^(uint64(0xff) << 56))
	// non-zero padding is indicated by set MSb / a negative length
	if lenField < 0 ***REMOVED***
		// padding is stored in lower 3 bits of length MSB
		padBytes = int64((uint64(lenField) >> 56) & 0x7)
	***REMOVED***
	return
***REMOVED***

// isTornEntry determines whether the last entry of the WAL was partially written
// and corrupted because of a torn write.
func (d *decoder) isTornEntry(data []byte) bool ***REMOVED***
	if len(d.brs) != 1 ***REMOVED***
		return false
	***REMOVED***

	fileOff := d.lastValidOff + 8
	curOff := 0
	chunks := [][]byte***REMOVED******REMOVED***
	// split data on sector boundaries
	for curOff < len(data) ***REMOVED***
		chunkLen := int(minSectorSize - (fileOff % minSectorSize))
		if chunkLen > len(data)-curOff ***REMOVED***
			chunkLen = len(data) - curOff
		***REMOVED***
		chunks = append(chunks, data[curOff:curOff+chunkLen])
		fileOff += int64(chunkLen)
		curOff += chunkLen
	***REMOVED***

	// if any data for a sector chunk is all 0, it's a torn write
	for _, sect := range chunks ***REMOVED***
		isZero := true
		for _, v := range sect ***REMOVED***
			if v != 0 ***REMOVED***
				isZero = false
				break
			***REMOVED***
		***REMOVED***
		if isZero ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (d *decoder) updateCRC(prevCrc uint32) ***REMOVED***
	d.crc = crc.New(prevCrc, crcTable)
***REMOVED***

func (d *decoder) lastCRC() uint32 ***REMOVED***
	return d.crc.Sum32()
***REMOVED***

func (d *decoder) lastOffset() int64 ***REMOVED*** return d.lastValidOff ***REMOVED***

func mustUnmarshalEntry(d []byte) raftpb.Entry ***REMOVED***
	var e raftpb.Entry
	pbutil.MustUnmarshal(&e, d)
	return e
***REMOVED***

func mustUnmarshalState(d []byte) raftpb.HardState ***REMOVED***
	var s raftpb.HardState
	pbutil.MustUnmarshal(&s, d)
	return s
***REMOVED***

func readInt64(r io.Reader) (int64, error) ***REMOVED***
	var n int64
	err := binary.Read(r, binary.LittleEndian, &n)
	return n, err
***REMOVED***
