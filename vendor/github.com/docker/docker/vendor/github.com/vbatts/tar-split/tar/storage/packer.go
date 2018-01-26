package storage

import (
	"encoding/json"
	"errors"
	"io"
	"path/filepath"
	"unicode/utf8"
)

// ErrDuplicatePath occurs when a tar archive has more than one entry for the
// same file path
var ErrDuplicatePath = errors.New("duplicates of file paths not supported")

// Packer describes the methods to pack Entries to a storage destination
type Packer interface ***REMOVED***
	// AddEntry packs the Entry and returns its position
	AddEntry(e Entry) (int, error)
***REMOVED***

// Unpacker describes the methods to read Entries from a source
type Unpacker interface ***REMOVED***
	// Next returns the next Entry being unpacked, or error, until io.EOF
	Next() (*Entry, error)
***REMOVED***

/* TODO(vbatts) figure out a good model for this
type PackUnpacker interface ***REMOVED***
	Packer
	Unpacker
***REMOVED***
*/

type jsonUnpacker struct ***REMOVED***
	seen seenNames
	dec  *json.Decoder
***REMOVED***

func (jup *jsonUnpacker) Next() (*Entry, error) ***REMOVED***
	var e Entry
	err := jup.dec.Decode(&e)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// check for dup name
	if e.Type == FileType ***REMOVED***
		cName := filepath.Clean(e.GetName())
		if _, ok := jup.seen[cName]; ok ***REMOVED***
			return nil, ErrDuplicatePath
		***REMOVED***
		jup.seen[cName] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	return &e, err
***REMOVED***

// NewJSONUnpacker provides an Unpacker that reads Entries (SegmentType and
// FileType) as a json document.
//
// Each Entry read are expected to be delimited by new line.
func NewJSONUnpacker(r io.Reader) Unpacker ***REMOVED***
	return &jsonUnpacker***REMOVED***
		dec:  json.NewDecoder(r),
		seen: seenNames***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

type jsonPacker struct ***REMOVED***
	w    io.Writer
	e    *json.Encoder
	pos  int
	seen seenNames
***REMOVED***

type seenNames map[string]struct***REMOVED******REMOVED***

func (jp *jsonPacker) AddEntry(e Entry) (int, error) ***REMOVED***
	// if Name is not valid utf8, switch it to raw first.
	if e.Name != "" ***REMOVED***
		if !utf8.ValidString(e.Name) ***REMOVED***
			e.NameRaw = []byte(e.Name)
			e.Name = ""
		***REMOVED***
	***REMOVED***

	// check early for dup name
	if e.Type == FileType ***REMOVED***
		cName := filepath.Clean(e.GetName())
		if _, ok := jp.seen[cName]; ok ***REMOVED***
			return -1, ErrDuplicatePath
		***REMOVED***
		jp.seen[cName] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	e.Position = jp.pos
	err := jp.e.Encode(e)
	if err != nil ***REMOVED***
		return -1, err
	***REMOVED***

	// made it this far, increment now
	jp.pos++
	return e.Position, nil
***REMOVED***

// NewJSONPacker provides a Packer that writes each Entry (SegmentType and
// FileType) as a json document.
//
// The Entries are delimited by new line.
func NewJSONPacker(w io.Writer) Packer ***REMOVED***
	return &jsonPacker***REMOVED***
		w:    w,
		e:    json.NewEncoder(w),
		seen: seenNames***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

/*
TODO(vbatts) perhaps have a more compact packer/unpacker, maybe using msgapck
(https://github.com/ugorji/go)


Even though, since our jsonUnpacker and jsonPacker just take
io.Reader/io.Writer, then we can get away with passing them a
gzip.Reader/gzip.Writer
*/
