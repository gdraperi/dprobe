package winio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"unicode/utf16"
	"unsafe"
)

const (
	reparseTagMountPoint = 0xA0000003
	reparseTagSymlink    = 0xA000000C
)

type reparseDataBuffer struct ***REMOVED***
	ReparseTag           uint32
	ReparseDataLength    uint16
	Reserved             uint16
	SubstituteNameOffset uint16
	SubstituteNameLength uint16
	PrintNameOffset      uint16
	PrintNameLength      uint16
***REMOVED***

// ReparsePoint describes a Win32 symlink or mount point.
type ReparsePoint struct ***REMOVED***
	Target       string
	IsMountPoint bool
***REMOVED***

// UnsupportedReparsePointError is returned when trying to decode a non-symlink or
// mount point reparse point.
type UnsupportedReparsePointError struct ***REMOVED***
	Tag uint32
***REMOVED***

func (e *UnsupportedReparsePointError) Error() string ***REMOVED***
	return fmt.Sprintf("unsupported reparse point %x", e.Tag)
***REMOVED***

// DecodeReparsePoint decodes a Win32 REPARSE_DATA_BUFFER structure containing either a symlink
// or a mount point.
func DecodeReparsePoint(b []byte) (*ReparsePoint, error) ***REMOVED***
	tag := binary.LittleEndian.Uint32(b[0:4])
	return DecodeReparsePointData(tag, b[8:])
***REMOVED***

func DecodeReparsePointData(tag uint32, b []byte) (*ReparsePoint, error) ***REMOVED***
	isMountPoint := false
	switch tag ***REMOVED***
	case reparseTagMountPoint:
		isMountPoint = true
	case reparseTagSymlink:
	default:
		return nil, &UnsupportedReparsePointError***REMOVED***tag***REMOVED***
	***REMOVED***
	nameOffset := 8 + binary.LittleEndian.Uint16(b[4:6])
	if !isMountPoint ***REMOVED***
		nameOffset += 4
	***REMOVED***
	nameLength := binary.LittleEndian.Uint16(b[6:8])
	name := make([]uint16, nameLength/2)
	err := binary.Read(bytes.NewReader(b[nameOffset:nameOffset+nameLength]), binary.LittleEndian, &name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &ReparsePoint***REMOVED***string(utf16.Decode(name)), isMountPoint***REMOVED***, nil
***REMOVED***

func isDriveLetter(c byte) bool ***REMOVED***
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
***REMOVED***

// EncodeReparsePoint encodes a Win32 REPARSE_DATA_BUFFER structure describing a symlink or
// mount point.
func EncodeReparsePoint(rp *ReparsePoint) []byte ***REMOVED***
	// Generate an NT path and determine if this is a relative path.
	var ntTarget string
	relative := false
	if strings.HasPrefix(rp.Target, `\\?\`) ***REMOVED***
		ntTarget = `\??\` + rp.Target[4:]
	***REMOVED*** else if strings.HasPrefix(rp.Target, `\\`) ***REMOVED***
		ntTarget = `\??\UNC\` + rp.Target[2:]
	***REMOVED*** else if len(rp.Target) >= 2 && isDriveLetter(rp.Target[0]) && rp.Target[1] == ':' ***REMOVED***
		ntTarget = `\??\` + rp.Target
	***REMOVED*** else ***REMOVED***
		ntTarget = rp.Target
		relative = true
	***REMOVED***

	// The paths must be NUL-terminated even though they are counted strings.
	target16 := utf16.Encode([]rune(rp.Target + "\x00"))
	ntTarget16 := utf16.Encode([]rune(ntTarget + "\x00"))

	size := int(unsafe.Sizeof(reparseDataBuffer***REMOVED******REMOVED***)) - 8
	size += len(ntTarget16)*2 + len(target16)*2

	tag := uint32(reparseTagMountPoint)
	if !rp.IsMountPoint ***REMOVED***
		tag = reparseTagSymlink
		size += 4 // Add room for symlink flags
	***REMOVED***

	data := reparseDataBuffer***REMOVED***
		ReparseTag:           tag,
		ReparseDataLength:    uint16(size),
		SubstituteNameOffset: 0,
		SubstituteNameLength: uint16((len(ntTarget16) - 1) * 2),
		PrintNameOffset:      uint16(len(ntTarget16) * 2),
		PrintNameLength:      uint16((len(target16) - 1) * 2),
	***REMOVED***

	var b bytes.Buffer
	binary.Write(&b, binary.LittleEndian, &data)
	if !rp.IsMountPoint ***REMOVED***
		flags := uint32(0)
		if relative ***REMOVED***
			flags |= 1
		***REMOVED***
		binary.Write(&b, binary.LittleEndian, flags)
	***REMOVED***

	binary.Write(&b, binary.LittleEndian, ntTarget16)
	binary.Write(&b, binary.LittleEndian, target16)
	return b.Bytes()
***REMOVED***
