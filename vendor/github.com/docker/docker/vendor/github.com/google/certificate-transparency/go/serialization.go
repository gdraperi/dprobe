package ct

import (
	"bytes"
	"container/list"
	"crypto"
	"encoding/asn1"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Variable size structure prefix-header byte lengths
const (
	CertificateLengthBytes      = 3
	PreCertificateLengthBytes   = 3
	ExtensionsLengthBytes       = 2
	CertificateChainLengthBytes = 3
	SignatureLengthBytes        = 2
	JSONLengthBytes             = 3
)

// Max lengths
const (
	MaxCertificateLength = (1 << 24) - 1
	MaxExtensionsLength  = (1 << 16) - 1
	MaxSCTInListLength   = (1 << 16) - 1
	MaxSCTListLength     = (1 << 16) - 1
)

func writeUint(w io.Writer, value uint64, numBytes int) error ***REMOVED***
	buf := make([]uint8, numBytes)
	for i := 0; i < numBytes; i++ ***REMOVED***
		buf[numBytes-i-1] = uint8(value & 0xff)
		value >>= 8
	***REMOVED***
	if value != 0 ***REMOVED***
		return errors.New("numBytes was insufficiently large to represent value")
	***REMOVED***
	if _, err := w.Write(buf); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func writeVarBytes(w io.Writer, value []byte, numLenBytes int) error ***REMOVED***
	if err := writeUint(w, uint64(len(value)), numLenBytes); err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, err := w.Write(value); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func readUint(r io.Reader, numBytes int) (uint64, error) ***REMOVED***
	var l uint64
	for i := 0; i < numBytes; i++ ***REMOVED***
		l <<= 8
		var t uint8
		if err := binary.Read(r, binary.BigEndian, &t); err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		l |= uint64(t)
	***REMOVED***
	return l, nil
***REMOVED***

// Reads a variable length array of bytes from |r|. |numLenBytes| specifies the
// number of (BigEndian) prefix-bytes which contain the length of the actual
// array data bytes that follow.
// Allocates an array to hold the contents and returns a slice view into it if
// the read was successful, or an error otherwise.
func readVarBytes(r io.Reader, numLenBytes int) ([]byte, error) ***REMOVED***
	switch ***REMOVED***
	case numLenBytes > 8:
		return nil, fmt.Errorf("numLenBytes too large (%d)", numLenBytes)
	case numLenBytes == 0:
		return nil, errors.New("numLenBytes should be > 0")
	***REMOVED***
	l, err := readUint(r, numLenBytes)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	data := make([]byte, l)
	if n, err := io.ReadFull(r, data); err != nil ***REMOVED***
		if err == io.EOF || err == io.ErrUnexpectedEOF ***REMOVED***
			return nil, fmt.Errorf("short read: expected %d but got %d", l, n)
		***REMOVED***
		return nil, err
	***REMOVED***
	return data, nil
***REMOVED***

// Reads a list of ASN1Cert types from |r|
func readASN1CertList(r io.Reader, totalLenBytes int, elementLenBytes int) ([]ASN1Cert, error) ***REMOVED***
	listBytes, err := readVarBytes(r, totalLenBytes)
	if err != nil ***REMOVED***
		return []ASN1Cert***REMOVED******REMOVED***, err
	***REMOVED***
	list := list.New()
	listReader := bytes.NewReader(listBytes)
	var entry []byte
	for err == nil ***REMOVED***
		entry, err = readVarBytes(listReader, elementLenBytes)
		if err != nil ***REMOVED***
			if err != io.EOF ***REMOVED***
				return []ASN1Cert***REMOVED******REMOVED***, err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			list.PushBack(entry)
		***REMOVED***
	***REMOVED***
	ret := make([]ASN1Cert, list.Len())
	i := 0
	for e := list.Front(); e != nil; e = e.Next() ***REMOVED***
		ret[i] = e.Value.([]byte)
		i++
	***REMOVED***
	return ret, nil
***REMOVED***

// ReadTimestampedEntryInto parses the byte-stream representation of a
// TimestampedEntry from |r| and populates the struct |t| with the data.  See
// RFC section 3.4 for details on the format.
// Returns a non-nil error if there was a problem.
func ReadTimestampedEntryInto(r io.Reader, t *TimestampedEntry) error ***REMOVED***
	var err error
	if err = binary.Read(r, binary.BigEndian, &t.Timestamp); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err = binary.Read(r, binary.BigEndian, &t.EntryType); err != nil ***REMOVED***
		return err
	***REMOVED***
	switch t.EntryType ***REMOVED***
	case X509LogEntryType:
		if t.X509Entry, err = readVarBytes(r, CertificateLengthBytes); err != nil ***REMOVED***
			return err
		***REMOVED***
	case PrecertLogEntryType:
		if err := binary.Read(r, binary.BigEndian, &t.PrecertEntry.IssuerKeyHash); err != nil ***REMOVED***
			return err
		***REMOVED***
		if t.PrecertEntry.TBSCertificate, err = readVarBytes(r, PreCertificateLengthBytes); err != nil ***REMOVED***
			return err
		***REMOVED***
	case XJSONLogEntryType:
		if t.JSONData, err = readVarBytes(r, JSONLengthBytes); err != nil ***REMOVED***
			return err
		***REMOVED***
	default:
		return fmt.Errorf("unknown EntryType: %d", t.EntryType)
	***REMOVED***
	t.Extensions, err = readVarBytes(r, ExtensionsLengthBytes)
	return nil
***REMOVED***

// SerializeTimestampedEntry writes timestamped entry to Writer.
// In case of error, w may contain garbage.
func SerializeTimestampedEntry(w io.Writer, t *TimestampedEntry) error ***REMOVED***
	if err := binary.Write(w, binary.BigEndian, t.Timestamp); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := binary.Write(w, binary.BigEndian, t.EntryType); err != nil ***REMOVED***
		return err
	***REMOVED***
	switch t.EntryType ***REMOVED***
	case X509LogEntryType:
		if err := writeVarBytes(w, t.X509Entry, CertificateLengthBytes); err != nil ***REMOVED***
			return err
		***REMOVED***
	case PrecertLogEntryType:
		if err := binary.Write(w, binary.BigEndian, t.PrecertEntry.IssuerKeyHash); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := writeVarBytes(w, t.PrecertEntry.TBSCertificate, PreCertificateLengthBytes); err != nil ***REMOVED***
			return err
		***REMOVED***
	case XJSONLogEntryType:
		// TODO: Pending google/certificate-transparency#1243, replace
		// with ObjectHash once supported by CT server.
		//jsonhash := objecthash.CommonJSONHash(string(t.JSONData))
		if err := writeVarBytes(w, []byte(t.JSONData), JSONLengthBytes); err != nil ***REMOVED***
			return err
		***REMOVED***
	default:
		return fmt.Errorf("unknown EntryType: %d", t.EntryType)
	***REMOVED***
	writeVarBytes(w, t.Extensions, ExtensionsLengthBytes)
	return nil
***REMOVED***

// ReadMerkleTreeLeaf parses the byte-stream representation of a MerkleTreeLeaf
// and returns a pointer to a new MerkleTreeLeaf structure containing the
// parsed data.
// See RFC section 3.4 for details on the format.
// Returns a pointer to a new MerkleTreeLeaf or non-nil error if there was a
// problem
func ReadMerkleTreeLeaf(r io.Reader) (*MerkleTreeLeaf, error) ***REMOVED***
	var m MerkleTreeLeaf
	if err := binary.Read(r, binary.BigEndian, &m.Version); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if m.Version != V1 ***REMOVED***
		return nil, fmt.Errorf("unknown Version %d", m.Version)
	***REMOVED***
	if err := binary.Read(r, binary.BigEndian, &m.LeafType); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if m.LeafType != TimestampedEntryLeafType ***REMOVED***
		return nil, fmt.Errorf("unknown LeafType %d", m.LeafType)
	***REMOVED***
	if err := ReadTimestampedEntryInto(r, &m.TimestampedEntry); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &m, nil
***REMOVED***

// UnmarshalX509ChainArray unmarshalls the contents of the "chain:" entry in a
// GetEntries response in the case where the entry refers to an X509 leaf.
func UnmarshalX509ChainArray(b []byte) ([]ASN1Cert, error) ***REMOVED***
	return readASN1CertList(bytes.NewReader(b), CertificateChainLengthBytes, CertificateLengthBytes)
***REMOVED***

// UnmarshalPrecertChainArray unmarshalls the contents of the "chain:" entry in
// a GetEntries response in the case where the entry refers to a Precertificate
// leaf.
func UnmarshalPrecertChainArray(b []byte) ([]ASN1Cert, error) ***REMOVED***
	var chain []ASN1Cert

	reader := bytes.NewReader(b)
	// read the pre-cert entry:
	precert, err := readVarBytes(reader, CertificateLengthBytes)
	if err != nil ***REMOVED***
		return chain, err
	***REMOVED***
	chain = append(chain, precert)
	// and then read and return the chain up to the root:
	remainingChain, err := readASN1CertList(reader, CertificateChainLengthBytes, CertificateLengthBytes)
	if err != nil ***REMOVED***
		return chain, err
	***REMOVED***
	chain = append(chain, remainingChain...)
	return chain, nil
***REMOVED***

// UnmarshalDigitallySigned reconstructs a DigitallySigned structure from a Reader
func UnmarshalDigitallySigned(r io.Reader) (*DigitallySigned, error) ***REMOVED***
	var h byte
	if err := binary.Read(r, binary.BigEndian, &h); err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to read HashAlgorithm: %v", err)
	***REMOVED***

	var s byte
	if err := binary.Read(r, binary.BigEndian, &s); err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to read SignatureAlgorithm: %v", err)
	***REMOVED***

	sig, err := readVarBytes(r, SignatureLengthBytes)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("failed to read Signature bytes: %v", err)
	***REMOVED***

	return &DigitallySigned***REMOVED***
		HashAlgorithm:      HashAlgorithm(h),
		SignatureAlgorithm: SignatureAlgorithm(s),
		Signature:          sig,
	***REMOVED***, nil
***REMOVED***

func marshalDigitallySignedHere(ds DigitallySigned, here []byte) ([]byte, error) ***REMOVED***
	sigLen := len(ds.Signature)
	dsOutLen := 2 + SignatureLengthBytes + sigLen
	if here == nil ***REMOVED***
		here = make([]byte, dsOutLen)
	***REMOVED***
	if len(here) < dsOutLen ***REMOVED***
		return nil, ErrNotEnoughBuffer
	***REMOVED***
	here = here[0:dsOutLen]

	here[0] = byte(ds.HashAlgorithm)
	here[1] = byte(ds.SignatureAlgorithm)
	binary.BigEndian.PutUint16(here[2:4], uint16(sigLen))
	copy(here[4:], ds.Signature)

	return here, nil
***REMOVED***

// MarshalDigitallySigned marshalls a DigitallySigned structure into a byte array
func MarshalDigitallySigned(ds DigitallySigned) ([]byte, error) ***REMOVED***
	return marshalDigitallySignedHere(ds, nil)
***REMOVED***

func checkCertificateFormat(cert ASN1Cert) error ***REMOVED***
	if len(cert) == 0 ***REMOVED***
		return errors.New("certificate is zero length")
	***REMOVED***
	if len(cert) > MaxCertificateLength ***REMOVED***
		return errors.New("certificate too large")
	***REMOVED***
	return nil
***REMOVED***

func checkExtensionsFormat(ext CTExtensions) error ***REMOVED***
	if len(ext) > MaxExtensionsLength ***REMOVED***
		return errors.New("extensions too large")
	***REMOVED***
	return nil
***REMOVED***

func serializeV1CertSCTSignatureInput(timestamp uint64, cert ASN1Cert, ext CTExtensions) ([]byte, error) ***REMOVED***
	if err := checkCertificateFormat(cert); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := checkExtensionsFormat(ext); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, V1); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := binary.Write(&buf, binary.BigEndian, CertificateTimestampSignatureType); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := binary.Write(&buf, binary.BigEndian, timestamp); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := binary.Write(&buf, binary.BigEndian, X509LogEntryType); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := writeVarBytes(&buf, cert, CertificateLengthBytes); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := writeVarBytes(&buf, ext, ExtensionsLengthBytes); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return buf.Bytes(), nil
***REMOVED***

func serializeV1JSONSCTSignatureInput(timestamp uint64, j []byte) ([]byte, error) ***REMOVED***
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, V1); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := binary.Write(&buf, binary.BigEndian, CertificateTimestampSignatureType); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := binary.Write(&buf, binary.BigEndian, timestamp); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := binary.Write(&buf, binary.BigEndian, XJSONLogEntryType); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := writeVarBytes(&buf, j, JSONLengthBytes); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := writeVarBytes(&buf, nil, ExtensionsLengthBytes); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return buf.Bytes(), nil
***REMOVED***

func serializeV1PrecertSCTSignatureInput(timestamp uint64, issuerKeyHash [issuerKeyHashLength]byte, tbs []byte, ext CTExtensions) ([]byte, error) ***REMOVED***
	if err := checkCertificateFormat(tbs); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := checkExtensionsFormat(ext); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, V1); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := binary.Write(&buf, binary.BigEndian, CertificateTimestampSignatureType); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := binary.Write(&buf, binary.BigEndian, timestamp); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := binary.Write(&buf, binary.BigEndian, PrecertLogEntryType); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if _, err := buf.Write(issuerKeyHash[:]); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := writeVarBytes(&buf, tbs, CertificateLengthBytes); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := writeVarBytes(&buf, ext, ExtensionsLengthBytes); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return buf.Bytes(), nil
***REMOVED***

func serializeV1SCTSignatureInput(sct SignedCertificateTimestamp, entry LogEntry) ([]byte, error) ***REMOVED***
	if sct.SCTVersion != V1 ***REMOVED***
		return nil, fmt.Errorf("unsupported SCT version, expected V1, but got %s", sct.SCTVersion)
	***REMOVED***
	if entry.Leaf.LeafType != TimestampedEntryLeafType ***REMOVED***
		return nil, fmt.Errorf("Unsupported leaf type %s", entry.Leaf.LeafType)
	***REMOVED***
	switch entry.Leaf.TimestampedEntry.EntryType ***REMOVED***
	case X509LogEntryType:
		return serializeV1CertSCTSignatureInput(sct.Timestamp, entry.Leaf.TimestampedEntry.X509Entry, entry.Leaf.TimestampedEntry.Extensions)
	case PrecertLogEntryType:
		return serializeV1PrecertSCTSignatureInput(sct.Timestamp, entry.Leaf.TimestampedEntry.PrecertEntry.IssuerKeyHash,
			entry.Leaf.TimestampedEntry.PrecertEntry.TBSCertificate,
			entry.Leaf.TimestampedEntry.Extensions)
	case XJSONLogEntryType:
		return serializeV1JSONSCTSignatureInput(sct.Timestamp, entry.Leaf.TimestampedEntry.JSONData)
	default:
		return nil, fmt.Errorf("unknown TimestampedEntryLeafType %s", entry.Leaf.TimestampedEntry.EntryType)
	***REMOVED***
***REMOVED***

// SerializeSCTSignatureInput serializes the passed in sct and log entry into
// the correct format for signing.
func SerializeSCTSignatureInput(sct SignedCertificateTimestamp, entry LogEntry) ([]byte, error) ***REMOVED***
	switch sct.SCTVersion ***REMOVED***
	case V1:
		return serializeV1SCTSignatureInput(sct, entry)
	default:
		return nil, fmt.Errorf("unknown SCT version %d", sct.SCTVersion)
	***REMOVED***
***REMOVED***

// SerializedLength will return the space (in bytes)
func (sct SignedCertificateTimestamp) SerializedLength() (int, error) ***REMOVED***
	switch sct.SCTVersion ***REMOVED***
	case V1:
		extLen := len(sct.Extensions)
		sigLen := len(sct.Signature.Signature)
		return 1 + 32 + 8 + 2 + extLen + 2 + 2 + sigLen, nil
	default:
		return 0, ErrInvalidVersion
	***REMOVED***
***REMOVED***

func serializeV1SCTHere(sct SignedCertificateTimestamp, here []byte) ([]byte, error) ***REMOVED***
	if sct.SCTVersion != V1 ***REMOVED***
		return nil, ErrInvalidVersion
	***REMOVED***
	sctLen, err := sct.SerializedLength()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if here == nil ***REMOVED***
		here = make([]byte, sctLen)
	***REMOVED***
	if len(here) < sctLen ***REMOVED***
		return nil, ErrNotEnoughBuffer
	***REMOVED***
	if err := checkExtensionsFormat(sct.Extensions); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	here = here[0:sctLen]

	// Write Version
	here[0] = byte(sct.SCTVersion)

	// Write LogID
	copy(here[1:33], sct.LogID[:])

	// Write Timestamp
	binary.BigEndian.PutUint64(here[33:41], sct.Timestamp)

	// Write Extensions
	extLen := len(sct.Extensions)
	binary.BigEndian.PutUint16(here[41:43], uint16(extLen))
	n := 43 + extLen
	copy(here[43:n], sct.Extensions)

	// Write Signature
	_, err = marshalDigitallySignedHere(sct.Signature, here[n:])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return here, nil
***REMOVED***

// SerializeSCTHere serializes the passed in sct into the format specified
// by RFC6962 section 3.2.
// If a bytes slice here is provided then it will attempt to serialize into the
// provided byte slice, ErrNotEnoughBuffer will be returned if the buffer is
// too small.
// If a nil byte slice is provided, a buffer for will be allocated for you
// The returned slice will be sliced to the correct length.
func SerializeSCTHere(sct SignedCertificateTimestamp, here []byte) ([]byte, error) ***REMOVED***
	switch sct.SCTVersion ***REMOVED***
	case V1:
		return serializeV1SCTHere(sct, here)
	default:
		return nil, fmt.Errorf("unknown SCT version %d", sct.SCTVersion)
	***REMOVED***
***REMOVED***

// SerializeSCT serializes the passed in sct into the format specified
// by RFC6962 section 3.2
// Equivalent to SerializeSCTHere(sct, nil)
func SerializeSCT(sct SignedCertificateTimestamp) ([]byte, error) ***REMOVED***
	return SerializeSCTHere(sct, nil)
***REMOVED***

func deserializeSCTV1(r io.Reader, sct *SignedCertificateTimestamp) error ***REMOVED***
	if err := binary.Read(r, binary.BigEndian, &sct.LogID); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := binary.Read(r, binary.BigEndian, &sct.Timestamp); err != nil ***REMOVED***
		return err
	***REMOVED***
	ext, err := readVarBytes(r, ExtensionsLengthBytes)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	sct.Extensions = ext
	ds, err := UnmarshalDigitallySigned(r)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	sct.Signature = *ds
	return nil
***REMOVED***

// DeserializeSCT reads an SCT from Reader.
func DeserializeSCT(r io.Reader) (*SignedCertificateTimestamp, error) ***REMOVED***
	var sct SignedCertificateTimestamp
	if err := binary.Read(r, binary.BigEndian, &sct.SCTVersion); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	switch sct.SCTVersion ***REMOVED***
	case V1:
		return &sct, deserializeSCTV1(r, &sct)
	default:
		return nil, fmt.Errorf("unknown SCT version %d", sct.SCTVersion)
	***REMOVED***
***REMOVED***

func serializeV1STHSignatureInput(sth SignedTreeHead) ([]byte, error) ***REMOVED***
	if sth.Version != V1 ***REMOVED***
		return nil, fmt.Errorf("invalid STH version %d", sth.Version)
	***REMOVED***
	if sth.TreeSize < 0 ***REMOVED***
		return nil, fmt.Errorf("invalid tree size %d", sth.TreeSize)
	***REMOVED***
	if len(sth.SHA256RootHash) != crypto.SHA256.Size() ***REMOVED***
		return nil, fmt.Errorf("invalid TreeHash length, got %d expected %d", len(sth.SHA256RootHash), crypto.SHA256.Size())
	***REMOVED***

	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, V1); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := binary.Write(&buf, binary.BigEndian, TreeHashSignatureType); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := binary.Write(&buf, binary.BigEndian, sth.Timestamp); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := binary.Write(&buf, binary.BigEndian, sth.TreeSize); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if err := binary.Write(&buf, binary.BigEndian, sth.SHA256RootHash); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return buf.Bytes(), nil
***REMOVED***

// SerializeSTHSignatureInput serializes the passed in sth into the correct
// format for signing.
func SerializeSTHSignatureInput(sth SignedTreeHead) ([]byte, error) ***REMOVED***
	switch sth.Version ***REMOVED***
	case V1:
		return serializeV1STHSignatureInput(sth)
	default:
		return nil, fmt.Errorf("unsupported STH version %d", sth.Version)
	***REMOVED***
***REMOVED***

// SCTListSerializedLength determines the length of the required buffer should a SCT List need to be serialized
func SCTListSerializedLength(scts []SignedCertificateTimestamp) (int, error) ***REMOVED***
	if len(scts) == 0 ***REMOVED***
		return 0, fmt.Errorf("SCT List empty")
	***REMOVED***

	sctListLen := 0
	for i, sct := range scts ***REMOVED***
		n, err := sct.SerializedLength()
		if err != nil ***REMOVED***
			return 0, fmt.Errorf("unable to determine length of SCT in position %d: %v", i, err)
		***REMOVED***
		if n > MaxSCTInListLength ***REMOVED***
			return 0, fmt.Errorf("SCT in position %d too large: %d", i, n)
		***REMOVED***
		sctListLen += 2 + n
	***REMOVED***

	return sctListLen, nil
***REMOVED***

// SerializeSCTList serializes the passed-in slice of SignedCertificateTimestamp into a
// byte slice as a SignedCertificateTimestampList (see RFC6962 Section 3.3)
func SerializeSCTList(scts []SignedCertificateTimestamp) ([]byte, error) ***REMOVED***
	size, err := SCTListSerializedLength(scts)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	fullSize := 2 + size // 2 bytes for length + size of SCT list
	if fullSize > MaxSCTListLength ***REMOVED***
		return nil, fmt.Errorf("SCT List too large to serialize: %d", fullSize)
	***REMOVED***
	buf := new(bytes.Buffer)
	buf.Grow(fullSize)
	if err = writeUint(buf, uint64(size), 2); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	for _, sct := range scts ***REMOVED***
		serialized, err := SerializeSCT(sct)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if err = writeVarBytes(buf, serialized, 2); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return asn1.Marshal(buf.Bytes()) // transform to Octet String
***REMOVED***

// SerializeMerkleTreeLeaf writes MerkleTreeLeaf to Writer.
// In case of error, w may contain garbage.
func SerializeMerkleTreeLeaf(w io.Writer, m *MerkleTreeLeaf) error ***REMOVED***
	if m.Version != V1 ***REMOVED***
		return fmt.Errorf("unknown Version %d", m.Version)
	***REMOVED***
	if err := binary.Write(w, binary.BigEndian, m.Version); err != nil ***REMOVED***
		return err
	***REMOVED***
	if m.LeafType != TimestampedEntryLeafType ***REMOVED***
		return fmt.Errorf("unknown LeafType %d", m.LeafType)
	***REMOVED***
	if err := binary.Write(w, binary.BigEndian, m.LeafType); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := SerializeTimestampedEntry(w, &m.TimestampedEntry); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// CreateX509MerkleTreeLeaf generates a MerkleTreeLeaf for an X509 cert
func CreateX509MerkleTreeLeaf(cert ASN1Cert, timestamp uint64) *MerkleTreeLeaf ***REMOVED***
	return &MerkleTreeLeaf***REMOVED***
		Version:  V1,
		LeafType: TimestampedEntryLeafType,
		TimestampedEntry: TimestampedEntry***REMOVED***
			Timestamp: timestamp,
			EntryType: X509LogEntryType,
			X509Entry: cert,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// CreateJSONMerkleTreeLeaf creates the merkle tree leaf for json data.
func CreateJSONMerkleTreeLeaf(data interface***REMOVED******REMOVED***, timestamp uint64) *MerkleTreeLeaf ***REMOVED***
	jsonData, err := json.Marshal(AddJSONRequest***REMOVED***Data: data***REMOVED***)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	// Match the JSON serialization implemented by json-c
	jsonStr := strings.Replace(string(jsonData), ":", ": ", -1)
	jsonStr = strings.Replace(jsonStr, ",", ", ", -1)
	jsonStr = strings.Replace(jsonStr, "***REMOVED***", "***REMOVED*** ", -1)
	jsonStr = strings.Replace(jsonStr, "***REMOVED***", " ***REMOVED***", -1)
	jsonStr = strings.Replace(jsonStr, "/", `\/`, -1)
	// TODO: Pending google/certificate-transparency#1243, replace with
	// ObjectHash once supported by CT server.

	return &MerkleTreeLeaf***REMOVED***
		Version:  V1,
		LeafType: TimestampedEntryLeafType,
		TimestampedEntry: TimestampedEntry***REMOVED***
			Timestamp: timestamp,
			EntryType: XJSONLogEntryType,
			JSONData:  []byte(jsonStr),
		***REMOVED***,
	***REMOVED***
***REMOVED***
