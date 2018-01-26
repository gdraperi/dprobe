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
	"io"
	"os"
	"path/filepath"

	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/coreos/etcd/wal/walpb"
)

// Repair tries to repair ErrUnexpectedEOF in the
// last wal file by truncating.
func Repair(dirpath string) bool ***REMOVED***
	f, err := openLast(dirpath)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	defer f.Close()

	rec := &walpb.Record***REMOVED******REMOVED***
	decoder := newDecoder(f)
	for ***REMOVED***
		lastOffset := decoder.lastOffset()
		err := decoder.decode(rec)
		switch err ***REMOVED***
		case nil:
			// update crc of the decoder when necessary
			switch rec.Type ***REMOVED***
			case crcType:
				crc := decoder.crc.Sum32()
				// current crc of decoder must match the crc of the record.
				// do no need to match 0 crc, since the decoder is a new one at this case.
				if crc != 0 && rec.Validate(crc) != nil ***REMOVED***
					return false
				***REMOVED***
				decoder.updateCRC(rec.Crc)
			***REMOVED***
			continue
		case io.EOF:
			return true
		case io.ErrUnexpectedEOF:
			plog.Noticef("repairing %v", f.Name())
			bf, bferr := os.Create(f.Name() + ".broken")
			if bferr != nil ***REMOVED***
				plog.Errorf("could not repair %v, failed to create backup file", f.Name())
				return false
			***REMOVED***
			defer bf.Close()

			if _, err = f.Seek(0, io.SeekStart); err != nil ***REMOVED***
				plog.Errorf("could not repair %v, failed to read file", f.Name())
				return false
			***REMOVED***

			if _, err = io.Copy(bf, f); err != nil ***REMOVED***
				plog.Errorf("could not repair %v, failed to copy file", f.Name())
				return false
			***REMOVED***

			if err = f.Truncate(int64(lastOffset)); err != nil ***REMOVED***
				plog.Errorf("could not repair %v, failed to truncate file", f.Name())
				return false
			***REMOVED***
			if err = fileutil.Fsync(f.File); err != nil ***REMOVED***
				plog.Errorf("could not repair %v, failed to sync file", f.Name())
				return false
			***REMOVED***
			return true
		default:
			plog.Errorf("could not repair error (%v)", err)
			return false
		***REMOVED***
	***REMOVED***
***REMOVED***

// openLast opens the last wal file for read and write.
func openLast(dirpath string) (*fileutil.LockedFile, error) ***REMOVED***
	names, err := readWalNames(dirpath)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	last := filepath.Join(dirpath, names[len(names)-1])
	return fileutil.LockFile(last, os.O_RDWR, fileutil.PrivateFileMode)
***REMOVED***
