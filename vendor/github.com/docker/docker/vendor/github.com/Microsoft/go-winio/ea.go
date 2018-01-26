package winio

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type fileFullEaInformation struct ***REMOVED***
	NextEntryOffset uint32
	Flags           uint8
	NameLength      uint8
	ValueLength     uint16
***REMOVED***

var (
	fileFullEaInformationSize = binary.Size(&fileFullEaInformation***REMOVED******REMOVED***)

	errInvalidEaBuffer = errors.New("invalid extended attribute buffer")
	errEaNameTooLarge  = errors.New("extended attribute name too large")
	errEaValueTooLarge = errors.New("extended attribute value too large")
)

// ExtendedAttribute represents a single Windows EA.
type ExtendedAttribute struct ***REMOVED***
	Name  string
	Value []byte
	Flags uint8
***REMOVED***

func parseEa(b []byte) (ea ExtendedAttribute, nb []byte, err error) ***REMOVED***
	var info fileFullEaInformation
	err = binary.Read(bytes.NewReader(b), binary.LittleEndian, &info)
	if err != nil ***REMOVED***
		err = errInvalidEaBuffer
		return
	***REMOVED***

	nameOffset := fileFullEaInformationSize
	nameLen := int(info.NameLength)
	valueOffset := nameOffset + int(info.NameLength) + 1
	valueLen := int(info.ValueLength)
	nextOffset := int(info.NextEntryOffset)
	if valueLen+valueOffset > len(b) || nextOffset < 0 || nextOffset > len(b) ***REMOVED***
		err = errInvalidEaBuffer
		return
	***REMOVED***

	ea.Name = string(b[nameOffset : nameOffset+nameLen])
	ea.Value = b[valueOffset : valueOffset+valueLen]
	ea.Flags = info.Flags
	if info.NextEntryOffset != 0 ***REMOVED***
		nb = b[info.NextEntryOffset:]
	***REMOVED***
	return
***REMOVED***

// DecodeExtendedAttributes decodes a list of EAs from a FILE_FULL_EA_INFORMATION
// buffer retrieved from BackupRead, ZwQueryEaFile, etc.
func DecodeExtendedAttributes(b []byte) (eas []ExtendedAttribute, err error) ***REMOVED***
	for len(b) != 0 ***REMOVED***
		ea, nb, err := parseEa(b)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		eas = append(eas, ea)
		b = nb
	***REMOVED***
	return
***REMOVED***

func writeEa(buf *bytes.Buffer, ea *ExtendedAttribute, last bool) error ***REMOVED***
	if int(uint8(len(ea.Name))) != len(ea.Name) ***REMOVED***
		return errEaNameTooLarge
	***REMOVED***
	if int(uint16(len(ea.Value))) != len(ea.Value) ***REMOVED***
		return errEaValueTooLarge
	***REMOVED***
	entrySize := uint32(fileFullEaInformationSize + len(ea.Name) + 1 + len(ea.Value))
	withPadding := (entrySize + 3) &^ 3
	nextOffset := uint32(0)
	if !last ***REMOVED***
		nextOffset = withPadding
	***REMOVED***
	info := fileFullEaInformation***REMOVED***
		NextEntryOffset: nextOffset,
		Flags:           ea.Flags,
		NameLength:      uint8(len(ea.Name)),
		ValueLength:     uint16(len(ea.Value)),
	***REMOVED***

	err := binary.Write(buf, binary.LittleEndian, &info)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	_, err = buf.Write([]byte(ea.Name))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = buf.WriteByte(0)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	_, err = buf.Write(ea.Value)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	_, err = buf.Write([]byte***REMOVED***0, 0, 0***REMOVED***[0 : withPadding-entrySize])
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

// EncodeExtendedAttributes encodes a list of EAs into a FILE_FULL_EA_INFORMATION
// buffer for use with BackupWrite, ZwSetEaFile, etc.
func EncodeExtendedAttributes(eas []ExtendedAttribute) ([]byte, error) ***REMOVED***
	var buf bytes.Buffer
	for i := range eas ***REMOVED***
		last := false
		if i == len(eas)-1 ***REMOVED***
			last = true
		***REMOVED***

		err := writeEa(&buf, &eas[i], last)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return buf.Bytes(), nil
***REMOVED***
