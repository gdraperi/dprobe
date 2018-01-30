// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dnsmessage provides a mostly RFC 1035 compliant implementation of
// DNS message packing and unpacking.
//
// This implementation is designed to minimize heap allocations and avoid
// unnecessary packing and unpacking as much as possible.
package dnsmessage

import (
	"errors"
)

// Message formats

// A Type is a type of DNS request and response.
type Type uint16

// A Class is a type of network.
type Class uint16

// An OpCode is a DNS operation code.
type OpCode uint16

// An RCode is a DNS response status code.
type RCode uint16

// Wire constants.
const (
	// ResourceHeader.Type and Question.Type
	TypeA     Type = 1
	TypeNS    Type = 2
	TypeCNAME Type = 5
	TypeSOA   Type = 6
	TypePTR   Type = 12
	TypeMX    Type = 15
	TypeTXT   Type = 16
	TypeAAAA  Type = 28
	TypeSRV   Type = 33

	// Question.Type
	TypeWKS   Type = 11
	TypeHINFO Type = 13
	TypeMINFO Type = 14
	TypeAXFR  Type = 252
	TypeALL   Type = 255

	// ResourceHeader.Class and Question.Class
	ClassINET   Class = 1
	ClassCSNET  Class = 2
	ClassCHAOS  Class = 3
	ClassHESIOD Class = 4

	// Question.Class
	ClassANY Class = 255

	// Message.Rcode
	RCodeSuccess        RCode = 0
	RCodeFormatError    RCode = 1
	RCodeServerFailure  RCode = 2
	RCodeNameError      RCode = 3
	RCodeNotImplemented RCode = 4
	RCodeRefused        RCode = 5
)

var (
	// ErrNotStarted indicates that the prerequisite information isn't
	// available yet because the previous records haven't been appropriately
	// parsed, skipped or finished.
	ErrNotStarted = errors.New("parsing/packing of this type isn't available yet")

	// ErrSectionDone indicated that all records in the section have been
	// parsed or finished.
	ErrSectionDone = errors.New("parsing/packing of this section has completed")

	errBaseLen            = errors.New("insufficient data for base length type")
	errCalcLen            = errors.New("insufficient data for calculated length type")
	errReserved           = errors.New("segment prefix is reserved")
	errTooManyPtr         = errors.New("too many pointers (>10)")
	errInvalidPtr         = errors.New("invalid pointer")
	errNilResouceBody     = errors.New("nil resource body")
	errResourceLen        = errors.New("insufficient data for resource body length")
	errSegTooLong         = errors.New("segment length too long")
	errZeroSegLen         = errors.New("zero length segment")
	errResTooLong         = errors.New("resource length too long")
	errTooManyQuestions   = errors.New("too many Questions to pack (>65535)")
	errTooManyAnswers     = errors.New("too many Answers to pack (>65535)")
	errTooManyAuthorities = errors.New("too many Authorities to pack (>65535)")
	errTooManyAdditionals = errors.New("too many Additionals to pack (>65535)")
	errNonCanonicalName   = errors.New("name is not in canonical format (it must end with a .)")
)

// Internal constants.
const (
	// packStartingCap is the default initial buffer size allocated during
	// packing.
	//
	// The starting capacity doesn't matter too much, but most DNS responses
	// Will be <= 512 bytes as it is the limit for DNS over UDP.
	packStartingCap = 512

	// uint16Len is the length (in bytes) of a uint16.
	uint16Len = 2

	// uint32Len is the length (in bytes) of a uint32.
	uint32Len = 4

	// headerLen is the length (in bytes) of a DNS header.
	//
	// A header is comprised of 6 uint16s and no padding.
	headerLen = 6 * uint16Len
)

type nestedError struct ***REMOVED***
	// s is the current level's error message.
	s string

	// err is the nested error.
	err error
***REMOVED***

// nestedError implements error.Error.
func (e *nestedError) Error() string ***REMOVED***
	return e.s + ": " + e.err.Error()
***REMOVED***

// Header is a representation of a DNS message header.
type Header struct ***REMOVED***
	ID                 uint16
	Response           bool
	OpCode             OpCode
	Authoritative      bool
	Truncated          bool
	RecursionDesired   bool
	RecursionAvailable bool
	RCode              RCode
***REMOVED***

func (m *Header) pack() (id uint16, bits uint16) ***REMOVED***
	id = m.ID
	bits = uint16(m.OpCode)<<11 | uint16(m.RCode)
	if m.RecursionAvailable ***REMOVED***
		bits |= headerBitRA
	***REMOVED***
	if m.RecursionDesired ***REMOVED***
		bits |= headerBitRD
	***REMOVED***
	if m.Truncated ***REMOVED***
		bits |= headerBitTC
	***REMOVED***
	if m.Authoritative ***REMOVED***
		bits |= headerBitAA
	***REMOVED***
	if m.Response ***REMOVED***
		bits |= headerBitQR
	***REMOVED***
	return
***REMOVED***

// Message is a representation of a DNS message.
type Message struct ***REMOVED***
	Header
	Questions   []Question
	Answers     []Resource
	Authorities []Resource
	Additionals []Resource
***REMOVED***

type section uint8

const (
	sectionNotStarted section = iota
	sectionHeader
	sectionQuestions
	sectionAnswers
	sectionAuthorities
	sectionAdditionals
	sectionDone

	headerBitQR = 1 << 15 // query/response (response=1)
	headerBitAA = 1 << 10 // authoritative
	headerBitTC = 1 << 9  // truncated
	headerBitRD = 1 << 8  // recursion desired
	headerBitRA = 1 << 7  // recursion available
)

var sectionNames = map[section]string***REMOVED***
	sectionHeader:      "header",
	sectionQuestions:   "Question",
	sectionAnswers:     "Answer",
	sectionAuthorities: "Authority",
	sectionAdditionals: "Additional",
***REMOVED***

// header is the wire format for a DNS message header.
type header struct ***REMOVED***
	id          uint16
	bits        uint16
	questions   uint16
	answers     uint16
	authorities uint16
	additionals uint16
***REMOVED***

func (h *header) count(sec section) uint16 ***REMOVED***
	switch sec ***REMOVED***
	case sectionQuestions:
		return h.questions
	case sectionAnswers:
		return h.answers
	case sectionAuthorities:
		return h.authorities
	case sectionAdditionals:
		return h.additionals
	***REMOVED***
	return 0
***REMOVED***

func (h *header) pack(msg []byte) []byte ***REMOVED***
	msg = packUint16(msg, h.id)
	msg = packUint16(msg, h.bits)
	msg = packUint16(msg, h.questions)
	msg = packUint16(msg, h.answers)
	msg = packUint16(msg, h.authorities)
	return packUint16(msg, h.additionals)
***REMOVED***

func (h *header) unpack(msg []byte, off int) (int, error) ***REMOVED***
	newOff := off
	var err error
	if h.id, newOff, err = unpackUint16(msg, newOff); err != nil ***REMOVED***
		return off, &nestedError***REMOVED***"id", err***REMOVED***
	***REMOVED***
	if h.bits, newOff, err = unpackUint16(msg, newOff); err != nil ***REMOVED***
		return off, &nestedError***REMOVED***"bits", err***REMOVED***
	***REMOVED***
	if h.questions, newOff, err = unpackUint16(msg, newOff); err != nil ***REMOVED***
		return off, &nestedError***REMOVED***"questions", err***REMOVED***
	***REMOVED***
	if h.answers, newOff, err = unpackUint16(msg, newOff); err != nil ***REMOVED***
		return off, &nestedError***REMOVED***"answers", err***REMOVED***
	***REMOVED***
	if h.authorities, newOff, err = unpackUint16(msg, newOff); err != nil ***REMOVED***
		return off, &nestedError***REMOVED***"authorities", err***REMOVED***
	***REMOVED***
	if h.additionals, newOff, err = unpackUint16(msg, newOff); err != nil ***REMOVED***
		return off, &nestedError***REMOVED***"additionals", err***REMOVED***
	***REMOVED***
	return newOff, nil
***REMOVED***

func (h *header) header() Header ***REMOVED***
	return Header***REMOVED***
		ID:                 h.id,
		Response:           (h.bits & headerBitQR) != 0,
		OpCode:             OpCode(h.bits>>11) & 0xF,
		Authoritative:      (h.bits & headerBitAA) != 0,
		Truncated:          (h.bits & headerBitTC) != 0,
		RecursionDesired:   (h.bits & headerBitRD) != 0,
		RecursionAvailable: (h.bits & headerBitRA) != 0,
		RCode:              RCode(h.bits & 0xF),
	***REMOVED***
***REMOVED***

// A Resource is a DNS resource record.
type Resource struct ***REMOVED***
	Header ResourceHeader
	Body   ResourceBody
***REMOVED***

// A ResourceBody is a DNS resource record minus the header.
type ResourceBody interface ***REMOVED***
	// pack packs a Resource except for its header.
	pack(msg []byte, compression map[string]int) ([]byte, error)

	// realType returns the actual type of the Resource. This is used to
	// fill in the header Type field.
	realType() Type
***REMOVED***

func (r *Resource) pack(msg []byte, compression map[string]int) ([]byte, error) ***REMOVED***
	if r.Body == nil ***REMOVED***
		return msg, errNilResouceBody
	***REMOVED***
	oldMsg := msg
	r.Header.Type = r.Body.realType()
	msg, length, err := r.Header.pack(msg, compression)
	if err != nil ***REMOVED***
		return msg, &nestedError***REMOVED***"ResourceHeader", err***REMOVED***
	***REMOVED***
	preLen := len(msg)
	msg, err = r.Body.pack(msg, compression)
	if err != nil ***REMOVED***
		return msg, &nestedError***REMOVED***"content", err***REMOVED***
	***REMOVED***
	if err := r.Header.fixLen(msg, length, preLen); err != nil ***REMOVED***
		return oldMsg, err
	***REMOVED***
	return msg, nil
***REMOVED***

// A Parser allows incrementally parsing a DNS message.
//
// When parsing is started, the Header is parsed. Next, each Question can be
// either parsed or skipped. Alternatively, all Questions can be skipped at
// once. When all Questions have been parsed, attempting to parse Questions
// will return (nil, nil) and attempting to skip Questions will return
// (true, nil). After all Questions have been either parsed or skipped, all
// Answers, Authorities and Additionals can be either parsed or skipped in the
// same way, and each type of Resource must be fully parsed or skipped before
// proceeding to the next type of Resource.
//
// Note that there is no requirement to fully skip or parse the message.
type Parser struct ***REMOVED***
	msg    []byte
	header header

	section        section
	off            int
	index          int
	resHeaderValid bool
	resHeader      ResourceHeader
***REMOVED***

// Start parses the header and enables the parsing of Questions.
func (p *Parser) Start(msg []byte) (Header, error) ***REMOVED***
	if p.msg != nil ***REMOVED***
		*p = Parser***REMOVED******REMOVED***
	***REMOVED***
	p.msg = msg
	var err error
	if p.off, err = p.header.unpack(msg, 0); err != nil ***REMOVED***
		return Header***REMOVED******REMOVED***, &nestedError***REMOVED***"unpacking header", err***REMOVED***
	***REMOVED***
	p.section = sectionQuestions
	return p.header.header(), nil
***REMOVED***

func (p *Parser) checkAdvance(sec section) error ***REMOVED***
	if p.section < sec ***REMOVED***
		return ErrNotStarted
	***REMOVED***
	if p.section > sec ***REMOVED***
		return ErrSectionDone
	***REMOVED***
	p.resHeaderValid = false
	if p.index == int(p.header.count(sec)) ***REMOVED***
		p.index = 0
		p.section++
		return ErrSectionDone
	***REMOVED***
	return nil
***REMOVED***

func (p *Parser) resource(sec section) (Resource, error) ***REMOVED***
	var r Resource
	var err error
	r.Header, err = p.resourceHeader(sec)
	if err != nil ***REMOVED***
		return r, err
	***REMOVED***
	p.resHeaderValid = false
	r.Body, p.off, err = unpackResourceBody(p.msg, p.off, r.Header)
	if err != nil ***REMOVED***
		return Resource***REMOVED******REMOVED***, &nestedError***REMOVED***"unpacking " + sectionNames[sec], err***REMOVED***
	***REMOVED***
	p.index++
	return r, nil
***REMOVED***

func (p *Parser) resourceHeader(sec section) (ResourceHeader, error) ***REMOVED***
	if p.resHeaderValid ***REMOVED***
		return p.resHeader, nil
	***REMOVED***
	if err := p.checkAdvance(sec); err != nil ***REMOVED***
		return ResourceHeader***REMOVED******REMOVED***, err
	***REMOVED***
	var hdr ResourceHeader
	off, err := hdr.unpack(p.msg, p.off)
	if err != nil ***REMOVED***
		return ResourceHeader***REMOVED******REMOVED***, err
	***REMOVED***
	p.resHeaderValid = true
	p.resHeader = hdr
	p.off = off
	return hdr, nil
***REMOVED***

func (p *Parser) skipResource(sec section) error ***REMOVED***
	if p.resHeaderValid ***REMOVED***
		newOff := p.off + int(p.resHeader.Length)
		if newOff > len(p.msg) ***REMOVED***
			return errResourceLen
		***REMOVED***
		p.off = newOff
		p.resHeaderValid = false
		p.index++
		return nil
	***REMOVED***
	if err := p.checkAdvance(sec); err != nil ***REMOVED***
		return err
	***REMOVED***
	var err error
	p.off, err = skipResource(p.msg, p.off)
	if err != nil ***REMOVED***
		return &nestedError***REMOVED***"skipping: " + sectionNames[sec], err***REMOVED***
	***REMOVED***
	p.index++
	return nil
***REMOVED***

// Question parses a single Question.
func (p *Parser) Question() (Question, error) ***REMOVED***
	if err := p.checkAdvance(sectionQuestions); err != nil ***REMOVED***
		return Question***REMOVED******REMOVED***, err
	***REMOVED***
	var name Name
	off, err := name.unpack(p.msg, p.off)
	if err != nil ***REMOVED***
		return Question***REMOVED******REMOVED***, &nestedError***REMOVED***"unpacking Question.Name", err***REMOVED***
	***REMOVED***
	typ, off, err := unpackType(p.msg, off)
	if err != nil ***REMOVED***
		return Question***REMOVED******REMOVED***, &nestedError***REMOVED***"unpacking Question.Type", err***REMOVED***
	***REMOVED***
	class, off, err := unpackClass(p.msg, off)
	if err != nil ***REMOVED***
		return Question***REMOVED******REMOVED***, &nestedError***REMOVED***"unpacking Question.Class", err***REMOVED***
	***REMOVED***
	p.off = off
	p.index++
	return Question***REMOVED***name, typ, class***REMOVED***, nil
***REMOVED***

// AllQuestions parses all Questions.
func (p *Parser) AllQuestions() ([]Question, error) ***REMOVED***
	qs := make([]Question, 0, p.header.questions)
	for ***REMOVED***
		q, err := p.Question()
		if err == ErrSectionDone ***REMOVED***
			return qs, nil
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		qs = append(qs, q)
	***REMOVED***
***REMOVED***

// SkipQuestion skips a single Question.
func (p *Parser) SkipQuestion() error ***REMOVED***
	if err := p.checkAdvance(sectionQuestions); err != nil ***REMOVED***
		return err
	***REMOVED***
	off, err := skipName(p.msg, p.off)
	if err != nil ***REMOVED***
		return &nestedError***REMOVED***"skipping Question Name", err***REMOVED***
	***REMOVED***
	if off, err = skipType(p.msg, off); err != nil ***REMOVED***
		return &nestedError***REMOVED***"skipping Question Type", err***REMOVED***
	***REMOVED***
	if off, err = skipClass(p.msg, off); err != nil ***REMOVED***
		return &nestedError***REMOVED***"skipping Question Class", err***REMOVED***
	***REMOVED***
	p.off = off
	p.index++
	return nil
***REMOVED***

// SkipAllQuestions skips all Questions.
func (p *Parser) SkipAllQuestions() error ***REMOVED***
	for ***REMOVED***
		if err := p.SkipQuestion(); err == ErrSectionDone ***REMOVED***
			return nil
		***REMOVED*** else if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
***REMOVED***

// AnswerHeader parses a single Answer ResourceHeader.
func (p *Parser) AnswerHeader() (ResourceHeader, error) ***REMOVED***
	return p.resourceHeader(sectionAnswers)
***REMOVED***

// Answer parses a single Answer Resource.
func (p *Parser) Answer() (Resource, error) ***REMOVED***
	return p.resource(sectionAnswers)
***REMOVED***

// AllAnswers parses all Answer Resources.
func (p *Parser) AllAnswers() ([]Resource, error) ***REMOVED***
	as := make([]Resource, 0, p.header.answers)
	for ***REMOVED***
		a, err := p.Answer()
		if err == ErrSectionDone ***REMOVED***
			return as, nil
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		as = append(as, a)
	***REMOVED***
***REMOVED***

// SkipAnswer skips a single Answer Resource.
func (p *Parser) SkipAnswer() error ***REMOVED***
	return p.skipResource(sectionAnswers)
***REMOVED***

// SkipAllAnswers skips all Answer Resources.
func (p *Parser) SkipAllAnswers() error ***REMOVED***
	for ***REMOVED***
		if err := p.SkipAnswer(); err == ErrSectionDone ***REMOVED***
			return nil
		***REMOVED*** else if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
***REMOVED***

// AuthorityHeader parses a single Authority ResourceHeader.
func (p *Parser) AuthorityHeader() (ResourceHeader, error) ***REMOVED***
	return p.resourceHeader(sectionAuthorities)
***REMOVED***

// Authority parses a single Authority Resource.
func (p *Parser) Authority() (Resource, error) ***REMOVED***
	return p.resource(sectionAuthorities)
***REMOVED***

// AllAuthorities parses all Authority Resources.
func (p *Parser) AllAuthorities() ([]Resource, error) ***REMOVED***
	as := make([]Resource, 0, p.header.authorities)
	for ***REMOVED***
		a, err := p.Authority()
		if err == ErrSectionDone ***REMOVED***
			return as, nil
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		as = append(as, a)
	***REMOVED***
***REMOVED***

// SkipAuthority skips a single Authority Resource.
func (p *Parser) SkipAuthority() error ***REMOVED***
	return p.skipResource(sectionAuthorities)
***REMOVED***

// SkipAllAuthorities skips all Authority Resources.
func (p *Parser) SkipAllAuthorities() error ***REMOVED***
	for ***REMOVED***
		if err := p.SkipAuthority(); err == ErrSectionDone ***REMOVED***
			return nil
		***REMOVED*** else if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
***REMOVED***

// AdditionalHeader parses a single Additional ResourceHeader.
func (p *Parser) AdditionalHeader() (ResourceHeader, error) ***REMOVED***
	return p.resourceHeader(sectionAdditionals)
***REMOVED***

// Additional parses a single Additional Resource.
func (p *Parser) Additional() (Resource, error) ***REMOVED***
	return p.resource(sectionAdditionals)
***REMOVED***

// AllAdditionals parses all Additional Resources.
func (p *Parser) AllAdditionals() ([]Resource, error) ***REMOVED***
	as := make([]Resource, 0, p.header.additionals)
	for ***REMOVED***
		a, err := p.Additional()
		if err == ErrSectionDone ***REMOVED***
			return as, nil
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		as = append(as, a)
	***REMOVED***
***REMOVED***

// SkipAdditional skips a single Additional Resource.
func (p *Parser) SkipAdditional() error ***REMOVED***
	return p.skipResource(sectionAdditionals)
***REMOVED***

// SkipAllAdditionals skips all Additional Resources.
func (p *Parser) SkipAllAdditionals() error ***REMOVED***
	for ***REMOVED***
		if err := p.SkipAdditional(); err == ErrSectionDone ***REMOVED***
			return nil
		***REMOVED*** else if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
***REMOVED***

// CNAMEResource parses a single CNAMEResource.
//
// One of the XXXHeader methods must have been called before calling this
// method.
func (p *Parser) CNAMEResource() (CNAMEResource, error) ***REMOVED***
	if !p.resHeaderValid || p.resHeader.Type != TypeCNAME ***REMOVED***
		return CNAMEResource***REMOVED******REMOVED***, ErrNotStarted
	***REMOVED***
	r, err := unpackCNAMEResource(p.msg, p.off)
	if err != nil ***REMOVED***
		return CNAMEResource***REMOVED******REMOVED***, err
	***REMOVED***
	p.off += int(p.resHeader.Length)
	p.resHeaderValid = false
	p.index++
	return r, nil
***REMOVED***

// MXResource parses a single MXResource.
//
// One of the XXXHeader methods must have been called before calling this
// method.
func (p *Parser) MXResource() (MXResource, error) ***REMOVED***
	if !p.resHeaderValid || p.resHeader.Type != TypeMX ***REMOVED***
		return MXResource***REMOVED******REMOVED***, ErrNotStarted
	***REMOVED***
	r, err := unpackMXResource(p.msg, p.off)
	if err != nil ***REMOVED***
		return MXResource***REMOVED******REMOVED***, err
	***REMOVED***
	p.off += int(p.resHeader.Length)
	p.resHeaderValid = false
	p.index++
	return r, nil
***REMOVED***

// NSResource parses a single NSResource.
//
// One of the XXXHeader methods must have been called before calling this
// method.
func (p *Parser) NSResource() (NSResource, error) ***REMOVED***
	if !p.resHeaderValid || p.resHeader.Type != TypeNS ***REMOVED***
		return NSResource***REMOVED******REMOVED***, ErrNotStarted
	***REMOVED***
	r, err := unpackNSResource(p.msg, p.off)
	if err != nil ***REMOVED***
		return NSResource***REMOVED******REMOVED***, err
	***REMOVED***
	p.off += int(p.resHeader.Length)
	p.resHeaderValid = false
	p.index++
	return r, nil
***REMOVED***

// PTRResource parses a single PTRResource.
//
// One of the XXXHeader methods must have been called before calling this
// method.
func (p *Parser) PTRResource() (PTRResource, error) ***REMOVED***
	if !p.resHeaderValid || p.resHeader.Type != TypePTR ***REMOVED***
		return PTRResource***REMOVED******REMOVED***, ErrNotStarted
	***REMOVED***
	r, err := unpackPTRResource(p.msg, p.off)
	if err != nil ***REMOVED***
		return PTRResource***REMOVED******REMOVED***, err
	***REMOVED***
	p.off += int(p.resHeader.Length)
	p.resHeaderValid = false
	p.index++
	return r, nil
***REMOVED***

// SOAResource parses a single SOAResource.
//
// One of the XXXHeader methods must have been called before calling this
// method.
func (p *Parser) SOAResource() (SOAResource, error) ***REMOVED***
	if !p.resHeaderValid || p.resHeader.Type != TypeSOA ***REMOVED***
		return SOAResource***REMOVED******REMOVED***, ErrNotStarted
	***REMOVED***
	r, err := unpackSOAResource(p.msg, p.off)
	if err != nil ***REMOVED***
		return SOAResource***REMOVED******REMOVED***, err
	***REMOVED***
	p.off += int(p.resHeader.Length)
	p.resHeaderValid = false
	p.index++
	return r, nil
***REMOVED***

// TXTResource parses a single TXTResource.
//
// One of the XXXHeader methods must have been called before calling this
// method.
func (p *Parser) TXTResource() (TXTResource, error) ***REMOVED***
	if !p.resHeaderValid || p.resHeader.Type != TypeTXT ***REMOVED***
		return TXTResource***REMOVED******REMOVED***, ErrNotStarted
	***REMOVED***
	r, err := unpackTXTResource(p.msg, p.off, p.resHeader.Length)
	if err != nil ***REMOVED***
		return TXTResource***REMOVED******REMOVED***, err
	***REMOVED***
	p.off += int(p.resHeader.Length)
	p.resHeaderValid = false
	p.index++
	return r, nil
***REMOVED***

// SRVResource parses a single SRVResource.
//
// One of the XXXHeader methods must have been called before calling this
// method.
func (p *Parser) SRVResource() (SRVResource, error) ***REMOVED***
	if !p.resHeaderValid || p.resHeader.Type != TypeSRV ***REMOVED***
		return SRVResource***REMOVED******REMOVED***, ErrNotStarted
	***REMOVED***
	r, err := unpackSRVResource(p.msg, p.off)
	if err != nil ***REMOVED***
		return SRVResource***REMOVED******REMOVED***, err
	***REMOVED***
	p.off += int(p.resHeader.Length)
	p.resHeaderValid = false
	p.index++
	return r, nil
***REMOVED***

// AResource parses a single AResource.
//
// One of the XXXHeader methods must have been called before calling this
// method.
func (p *Parser) AResource() (AResource, error) ***REMOVED***
	if !p.resHeaderValid || p.resHeader.Type != TypeA ***REMOVED***
		return AResource***REMOVED******REMOVED***, ErrNotStarted
	***REMOVED***
	r, err := unpackAResource(p.msg, p.off)
	if err != nil ***REMOVED***
		return AResource***REMOVED******REMOVED***, err
	***REMOVED***
	p.off += int(p.resHeader.Length)
	p.resHeaderValid = false
	p.index++
	return r, nil
***REMOVED***

// AAAAResource parses a single AAAAResource.
//
// One of the XXXHeader methods must have been called before calling this
// method.
func (p *Parser) AAAAResource() (AAAAResource, error) ***REMOVED***
	if !p.resHeaderValid || p.resHeader.Type != TypeAAAA ***REMOVED***
		return AAAAResource***REMOVED******REMOVED***, ErrNotStarted
	***REMOVED***
	r, err := unpackAAAAResource(p.msg, p.off)
	if err != nil ***REMOVED***
		return AAAAResource***REMOVED******REMOVED***, err
	***REMOVED***
	p.off += int(p.resHeader.Length)
	p.resHeaderValid = false
	p.index++
	return r, nil
***REMOVED***

// Unpack parses a full Message.
func (m *Message) Unpack(msg []byte) error ***REMOVED***
	var p Parser
	var err error
	if m.Header, err = p.Start(msg); err != nil ***REMOVED***
		return err
	***REMOVED***
	if m.Questions, err = p.AllQuestions(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if m.Answers, err = p.AllAnswers(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if m.Authorities, err = p.AllAuthorities(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if m.Additionals, err = p.AllAdditionals(); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// Pack packs a full Message.
func (m *Message) Pack() ([]byte, error) ***REMOVED***
	return m.AppendPack(make([]byte, 0, packStartingCap))
***REMOVED***

// AppendPack is like Pack but appends the full Message to b and returns the
// extended buffer.
func (m *Message) AppendPack(b []byte) ([]byte, error) ***REMOVED***
	// Validate the lengths. It is very unlikely that anyone will try to
	// pack more than 65535 of any particular type, but it is possible and
	// we should fail gracefully.
	if len(m.Questions) > int(^uint16(0)) ***REMOVED***
		return nil, errTooManyQuestions
	***REMOVED***
	if len(m.Answers) > int(^uint16(0)) ***REMOVED***
		return nil, errTooManyAnswers
	***REMOVED***
	if len(m.Authorities) > int(^uint16(0)) ***REMOVED***
		return nil, errTooManyAuthorities
	***REMOVED***
	if len(m.Additionals) > int(^uint16(0)) ***REMOVED***
		return nil, errTooManyAdditionals
	***REMOVED***

	var h header
	h.id, h.bits = m.Header.pack()

	h.questions = uint16(len(m.Questions))
	h.answers = uint16(len(m.Answers))
	h.authorities = uint16(len(m.Authorities))
	h.additionals = uint16(len(m.Additionals))

	msg := h.pack(b)

	// RFC 1035 allows (but does not require) compression for packing. RFC
	// 1035 requires unpacking implementations to support compression, so
	// unconditionally enabling it is fine.
	//
	// DNS lookups are typically done over UDP, and RFC 1035 states that UDP
	// DNS messages can be a maximum of 512 bytes long. Without compression,
	// many DNS response messages are over this limit, so enabling
	// compression will help ensure compliance.
	compression := map[string]int***REMOVED******REMOVED***

	for i := range m.Questions ***REMOVED***
		var err error
		if msg, err = m.Questions[i].pack(msg, compression); err != nil ***REMOVED***
			return nil, &nestedError***REMOVED***"packing Question", err***REMOVED***
		***REMOVED***
	***REMOVED***
	for i := range m.Answers ***REMOVED***
		var err error
		if msg, err = m.Answers[i].pack(msg, compression); err != nil ***REMOVED***
			return nil, &nestedError***REMOVED***"packing Answer", err***REMOVED***
		***REMOVED***
	***REMOVED***
	for i := range m.Authorities ***REMOVED***
		var err error
		if msg, err = m.Authorities[i].pack(msg, compression); err != nil ***REMOVED***
			return nil, &nestedError***REMOVED***"packing Authority", err***REMOVED***
		***REMOVED***
	***REMOVED***
	for i := range m.Additionals ***REMOVED***
		var err error
		if msg, err = m.Additionals[i].pack(msg, compression); err != nil ***REMOVED***
			return nil, &nestedError***REMOVED***"packing Additional", err***REMOVED***
		***REMOVED***
	***REMOVED***

	return msg, nil
***REMOVED***

// A Builder allows incrementally packing a DNS message.
type Builder struct ***REMOVED***
	msg         []byte
	header      header
	section     section
	compression map[string]int
***REMOVED***

// Start initializes the builder.
//
// buf is optional (nil is fine), but if provided, Start takes ownership of buf.
func (b *Builder) Start(buf []byte, h Header) ***REMOVED***
	b.StartWithoutCompression(buf, h)
	b.compression = map[string]int***REMOVED******REMOVED***
***REMOVED***

// StartWithoutCompression initializes the builder with compression disabled.
//
// This avoids compression related allocations, but can result in larger message
// sizes. Be careful with this mode as it can cause messages to exceed the UDP
// size limit.
//
// buf is optional (nil is fine), but if provided, Start takes ownership of buf.
func (b *Builder) StartWithoutCompression(buf []byte, h Header) ***REMOVED***
	*b = Builder***REMOVED***msg: buf***REMOVED***
	b.header.id, b.header.bits = h.pack()
	if cap(b.msg) < headerLen ***REMOVED***
		b.msg = make([]byte, 0, packStartingCap)
	***REMOVED***
	b.msg = b.msg[:headerLen]
	b.section = sectionHeader
***REMOVED***

func (b *Builder) startCheck(s section) error ***REMOVED***
	if b.section <= sectionNotStarted ***REMOVED***
		return ErrNotStarted
	***REMOVED***
	if b.section > s ***REMOVED***
		return ErrSectionDone
	***REMOVED***
	return nil
***REMOVED***

// StartQuestions prepares the builder for packing Questions.
func (b *Builder) StartQuestions() error ***REMOVED***
	if err := b.startCheck(sectionQuestions); err != nil ***REMOVED***
		return err
	***REMOVED***
	b.section = sectionQuestions
	return nil
***REMOVED***

// StartAnswers prepares the builder for packing Answers.
func (b *Builder) StartAnswers() error ***REMOVED***
	if err := b.startCheck(sectionAnswers); err != nil ***REMOVED***
		return err
	***REMOVED***
	b.section = sectionAnswers
	return nil
***REMOVED***

// StartAuthorities prepares the builder for packing Authorities.
func (b *Builder) StartAuthorities() error ***REMOVED***
	if err := b.startCheck(sectionAuthorities); err != nil ***REMOVED***
		return err
	***REMOVED***
	b.section = sectionAuthorities
	return nil
***REMOVED***

// StartAdditionals prepares the builder for packing Additionals.
func (b *Builder) StartAdditionals() error ***REMOVED***
	if err := b.startCheck(sectionAdditionals); err != nil ***REMOVED***
		return err
	***REMOVED***
	b.section = sectionAdditionals
	return nil
***REMOVED***

func (b *Builder) incrementSectionCount() error ***REMOVED***
	var count *uint16
	var err error
	switch b.section ***REMOVED***
	case sectionQuestions:
		count = &b.header.questions
		err = errTooManyQuestions
	case sectionAnswers:
		count = &b.header.answers
		err = errTooManyAnswers
	case sectionAuthorities:
		count = &b.header.authorities
		err = errTooManyAuthorities
	case sectionAdditionals:
		count = &b.header.additionals
		err = errTooManyAdditionals
	***REMOVED***
	if *count == ^uint16(0) ***REMOVED***
		return err
	***REMOVED***
	*count++
	return nil
***REMOVED***

// Question adds a single Question.
func (b *Builder) Question(q Question) error ***REMOVED***
	if b.section < sectionQuestions ***REMOVED***
		return ErrNotStarted
	***REMOVED***
	if b.section > sectionQuestions ***REMOVED***
		return ErrSectionDone
	***REMOVED***
	msg, err := q.pack(b.msg, b.compression)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := b.incrementSectionCount(); err != nil ***REMOVED***
		return err
	***REMOVED***
	b.msg = msg
	return nil
***REMOVED***

func (b *Builder) checkResourceSection() error ***REMOVED***
	if b.section < sectionAnswers ***REMOVED***
		return ErrNotStarted
	***REMOVED***
	if b.section > sectionAdditionals ***REMOVED***
		return ErrSectionDone
	***REMOVED***
	return nil
***REMOVED***

// CNAMEResource adds a single CNAMEResource.
func (b *Builder) CNAMEResource(h ResourceHeader, r CNAMEResource) error ***REMOVED***
	if err := b.checkResourceSection(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.Type = r.realType()
	msg, length, err := h.pack(b.msg, b.compression)
	if err != nil ***REMOVED***
		return &nestedError***REMOVED***"ResourceHeader", err***REMOVED***
	***REMOVED***
	preLen := len(msg)
	if msg, err = r.pack(msg, b.compression); err != nil ***REMOVED***
		return &nestedError***REMOVED***"CNAMEResource body", err***REMOVED***
	***REMOVED***
	if err := h.fixLen(msg, length, preLen); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := b.incrementSectionCount(); err != nil ***REMOVED***
		return err
	***REMOVED***
	b.msg = msg
	return nil
***REMOVED***

// MXResource adds a single MXResource.
func (b *Builder) MXResource(h ResourceHeader, r MXResource) error ***REMOVED***
	if err := b.checkResourceSection(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.Type = r.realType()
	msg, length, err := h.pack(b.msg, b.compression)
	if err != nil ***REMOVED***
		return &nestedError***REMOVED***"ResourceHeader", err***REMOVED***
	***REMOVED***
	preLen := len(msg)
	if msg, err = r.pack(msg, b.compression); err != nil ***REMOVED***
		return &nestedError***REMOVED***"MXResource body", err***REMOVED***
	***REMOVED***
	if err := h.fixLen(msg, length, preLen); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := b.incrementSectionCount(); err != nil ***REMOVED***
		return err
	***REMOVED***
	b.msg = msg
	return nil
***REMOVED***

// NSResource adds a single NSResource.
func (b *Builder) NSResource(h ResourceHeader, r NSResource) error ***REMOVED***
	if err := b.checkResourceSection(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.Type = r.realType()
	msg, length, err := h.pack(b.msg, b.compression)
	if err != nil ***REMOVED***
		return &nestedError***REMOVED***"ResourceHeader", err***REMOVED***
	***REMOVED***
	preLen := len(msg)
	if msg, err = r.pack(msg, b.compression); err != nil ***REMOVED***
		return &nestedError***REMOVED***"NSResource body", err***REMOVED***
	***REMOVED***
	if err := h.fixLen(msg, length, preLen); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := b.incrementSectionCount(); err != nil ***REMOVED***
		return err
	***REMOVED***
	b.msg = msg
	return nil
***REMOVED***

// PTRResource adds a single PTRResource.
func (b *Builder) PTRResource(h ResourceHeader, r PTRResource) error ***REMOVED***
	if err := b.checkResourceSection(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.Type = r.realType()
	msg, length, err := h.pack(b.msg, b.compression)
	if err != nil ***REMOVED***
		return &nestedError***REMOVED***"ResourceHeader", err***REMOVED***
	***REMOVED***
	preLen := len(msg)
	if msg, err = r.pack(msg, b.compression); err != nil ***REMOVED***
		return &nestedError***REMOVED***"PTRResource body", err***REMOVED***
	***REMOVED***
	if err := h.fixLen(msg, length, preLen); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := b.incrementSectionCount(); err != nil ***REMOVED***
		return err
	***REMOVED***
	b.msg = msg
	return nil
***REMOVED***

// SOAResource adds a single SOAResource.
func (b *Builder) SOAResource(h ResourceHeader, r SOAResource) error ***REMOVED***
	if err := b.checkResourceSection(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.Type = r.realType()
	msg, length, err := h.pack(b.msg, b.compression)
	if err != nil ***REMOVED***
		return &nestedError***REMOVED***"ResourceHeader", err***REMOVED***
	***REMOVED***
	preLen := len(msg)
	if msg, err = r.pack(msg, b.compression); err != nil ***REMOVED***
		return &nestedError***REMOVED***"SOAResource body", err***REMOVED***
	***REMOVED***
	if err := h.fixLen(msg, length, preLen); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := b.incrementSectionCount(); err != nil ***REMOVED***
		return err
	***REMOVED***
	b.msg = msg
	return nil
***REMOVED***

// TXTResource adds a single TXTResource.
func (b *Builder) TXTResource(h ResourceHeader, r TXTResource) error ***REMOVED***
	if err := b.checkResourceSection(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.Type = r.realType()
	msg, length, err := h.pack(b.msg, b.compression)
	if err != nil ***REMOVED***
		return &nestedError***REMOVED***"ResourceHeader", err***REMOVED***
	***REMOVED***
	preLen := len(msg)
	if msg, err = r.pack(msg, b.compression); err != nil ***REMOVED***
		return &nestedError***REMOVED***"TXTResource body", err***REMOVED***
	***REMOVED***
	if err := h.fixLen(msg, length, preLen); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := b.incrementSectionCount(); err != nil ***REMOVED***
		return err
	***REMOVED***
	b.msg = msg
	return nil
***REMOVED***

// SRVResource adds a single SRVResource.
func (b *Builder) SRVResource(h ResourceHeader, r SRVResource) error ***REMOVED***
	if err := b.checkResourceSection(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.Type = r.realType()
	msg, length, err := h.pack(b.msg, b.compression)
	if err != nil ***REMOVED***
		return &nestedError***REMOVED***"ResourceHeader", err***REMOVED***
	***REMOVED***
	preLen := len(msg)
	if msg, err = r.pack(msg, b.compression); err != nil ***REMOVED***
		return &nestedError***REMOVED***"SRVResource body", err***REMOVED***
	***REMOVED***
	if err := h.fixLen(msg, length, preLen); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := b.incrementSectionCount(); err != nil ***REMOVED***
		return err
	***REMOVED***
	b.msg = msg
	return nil
***REMOVED***

// AResource adds a single AResource.
func (b *Builder) AResource(h ResourceHeader, r AResource) error ***REMOVED***
	if err := b.checkResourceSection(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.Type = r.realType()
	msg, length, err := h.pack(b.msg, b.compression)
	if err != nil ***REMOVED***
		return &nestedError***REMOVED***"ResourceHeader", err***REMOVED***
	***REMOVED***
	preLen := len(msg)
	if msg, err = r.pack(msg, b.compression); err != nil ***REMOVED***
		return &nestedError***REMOVED***"AResource body", err***REMOVED***
	***REMOVED***
	if err := h.fixLen(msg, length, preLen); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := b.incrementSectionCount(); err != nil ***REMOVED***
		return err
	***REMOVED***
	b.msg = msg
	return nil
***REMOVED***

// AAAAResource adds a single AAAAResource.
func (b *Builder) AAAAResource(h ResourceHeader, r AAAAResource) error ***REMOVED***
	if err := b.checkResourceSection(); err != nil ***REMOVED***
		return err
	***REMOVED***
	h.Type = r.realType()
	msg, length, err := h.pack(b.msg, b.compression)
	if err != nil ***REMOVED***
		return &nestedError***REMOVED***"ResourceHeader", err***REMOVED***
	***REMOVED***
	preLen := len(msg)
	if msg, err = r.pack(msg, b.compression); err != nil ***REMOVED***
		return &nestedError***REMOVED***"AAAAResource body", err***REMOVED***
	***REMOVED***
	if err := h.fixLen(msg, length, preLen); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := b.incrementSectionCount(); err != nil ***REMOVED***
		return err
	***REMOVED***
	b.msg = msg
	return nil
***REMOVED***

// Finish ends message building and generates a binary message.
func (b *Builder) Finish() ([]byte, error) ***REMOVED***
	if b.section < sectionHeader ***REMOVED***
		return nil, ErrNotStarted
	***REMOVED***
	b.section = sectionDone
	b.header.pack(b.msg[:0])
	return b.msg, nil
***REMOVED***

// A ResourceHeader is the header of a DNS resource record. There are
// many types of DNS resource records, but they all share the same header.
type ResourceHeader struct ***REMOVED***
	// Name is the domain name for which this resource record pertains.
	Name Name

	// Type is the type of DNS resource record.
	//
	// This field will be set automatically during packing.
	Type Type

	// Class is the class of network to which this DNS resource record
	// pertains.
	Class Class

	// TTL is the length of time (measured in seconds) which this resource
	// record is valid for (time to live). All Resources in a set should
	// have the same TTL (RFC 2181 Section 5.2).
	TTL uint32

	// Length is the length of data in the resource record after the header.
	//
	// This field will be set automatically during packing.
	Length uint16
***REMOVED***

// pack packs all of the fields in a ResourceHeader except for the length. The
// length bytes are returned as a slice so they can be filled in after the rest
// of the Resource has been packed.
func (h *ResourceHeader) pack(oldMsg []byte, compression map[string]int) (msg []byte, length []byte, err error) ***REMOVED***
	msg = oldMsg
	if msg, err = h.Name.pack(msg, compression); err != nil ***REMOVED***
		return oldMsg, nil, &nestedError***REMOVED***"Name", err***REMOVED***
	***REMOVED***
	msg = packType(msg, h.Type)
	msg = packClass(msg, h.Class)
	msg = packUint32(msg, h.TTL)
	lenBegin := len(msg)
	msg = packUint16(msg, h.Length)
	return msg, msg[lenBegin : lenBegin+uint16Len], nil
***REMOVED***

func (h *ResourceHeader) unpack(msg []byte, off int) (int, error) ***REMOVED***
	newOff := off
	var err error
	if newOff, err = h.Name.unpack(msg, newOff); err != nil ***REMOVED***
		return off, &nestedError***REMOVED***"Name", err***REMOVED***
	***REMOVED***
	if h.Type, newOff, err = unpackType(msg, newOff); err != nil ***REMOVED***
		return off, &nestedError***REMOVED***"Type", err***REMOVED***
	***REMOVED***
	if h.Class, newOff, err = unpackClass(msg, newOff); err != nil ***REMOVED***
		return off, &nestedError***REMOVED***"Class", err***REMOVED***
	***REMOVED***
	if h.TTL, newOff, err = unpackUint32(msg, newOff); err != nil ***REMOVED***
		return off, &nestedError***REMOVED***"TTL", err***REMOVED***
	***REMOVED***
	if h.Length, newOff, err = unpackUint16(msg, newOff); err != nil ***REMOVED***
		return off, &nestedError***REMOVED***"Length", err***REMOVED***
	***REMOVED***
	return newOff, nil
***REMOVED***

func (h *ResourceHeader) fixLen(msg []byte, length []byte, preLen int) error ***REMOVED***
	conLen := len(msg) - preLen
	if conLen > int(^uint16(0)) ***REMOVED***
		return errResTooLong
	***REMOVED***

	// Fill in the length now that we know how long the content is.
	packUint16(length[:0], uint16(conLen))
	h.Length = uint16(conLen)

	return nil
***REMOVED***

func skipResource(msg []byte, off int) (int, error) ***REMOVED***
	newOff, err := skipName(msg, off)
	if err != nil ***REMOVED***
		return off, &nestedError***REMOVED***"Name", err***REMOVED***
	***REMOVED***
	if newOff, err = skipType(msg, newOff); err != nil ***REMOVED***
		return off, &nestedError***REMOVED***"Type", err***REMOVED***
	***REMOVED***
	if newOff, err = skipClass(msg, newOff); err != nil ***REMOVED***
		return off, &nestedError***REMOVED***"Class", err***REMOVED***
	***REMOVED***
	if newOff, err = skipUint32(msg, newOff); err != nil ***REMOVED***
		return off, &nestedError***REMOVED***"TTL", err***REMOVED***
	***REMOVED***
	length, newOff, err := unpackUint16(msg, newOff)
	if err != nil ***REMOVED***
		return off, &nestedError***REMOVED***"Length", err***REMOVED***
	***REMOVED***
	if newOff += int(length); newOff > len(msg) ***REMOVED***
		return off, errResourceLen
	***REMOVED***
	return newOff, nil
***REMOVED***

func packUint16(msg []byte, field uint16) []byte ***REMOVED***
	return append(msg, byte(field>>8), byte(field))
***REMOVED***

func unpackUint16(msg []byte, off int) (uint16, int, error) ***REMOVED***
	if off+uint16Len > len(msg) ***REMOVED***
		return 0, off, errBaseLen
	***REMOVED***
	return uint16(msg[off])<<8 | uint16(msg[off+1]), off + uint16Len, nil
***REMOVED***

func skipUint16(msg []byte, off int) (int, error) ***REMOVED***
	if off+uint16Len > len(msg) ***REMOVED***
		return off, errBaseLen
	***REMOVED***
	return off + uint16Len, nil
***REMOVED***

func packType(msg []byte, field Type) []byte ***REMOVED***
	return packUint16(msg, uint16(field))
***REMOVED***

func unpackType(msg []byte, off int) (Type, int, error) ***REMOVED***
	t, o, err := unpackUint16(msg, off)
	return Type(t), o, err
***REMOVED***

func skipType(msg []byte, off int) (int, error) ***REMOVED***
	return skipUint16(msg, off)
***REMOVED***

func packClass(msg []byte, field Class) []byte ***REMOVED***
	return packUint16(msg, uint16(field))
***REMOVED***

func unpackClass(msg []byte, off int) (Class, int, error) ***REMOVED***
	c, o, err := unpackUint16(msg, off)
	return Class(c), o, err
***REMOVED***

func skipClass(msg []byte, off int) (int, error) ***REMOVED***
	return skipUint16(msg, off)
***REMOVED***

func packUint32(msg []byte, field uint32) []byte ***REMOVED***
	return append(
		msg,
		byte(field>>24),
		byte(field>>16),
		byte(field>>8),
		byte(field),
	)
***REMOVED***

func unpackUint32(msg []byte, off int) (uint32, int, error) ***REMOVED***
	if off+uint32Len > len(msg) ***REMOVED***
		return 0, off, errBaseLen
	***REMOVED***
	v := uint32(msg[off])<<24 | uint32(msg[off+1])<<16 | uint32(msg[off+2])<<8 | uint32(msg[off+3])
	return v, off + uint32Len, nil
***REMOVED***

func skipUint32(msg []byte, off int) (int, error) ***REMOVED***
	if off+uint32Len > len(msg) ***REMOVED***
		return off, errBaseLen
	***REMOVED***
	return off + uint32Len, nil
***REMOVED***

func packText(msg []byte, field string) []byte ***REMOVED***
	for len(field) > 0 ***REMOVED***
		l := len(field)
		if l > 255 ***REMOVED***
			l = 255
		***REMOVED***
		msg = append(msg, byte(l))
		msg = append(msg, field[:l]...)
		field = field[l:]
	***REMOVED***
	return msg
***REMOVED***

func unpackText(msg []byte, off int) (string, int, error) ***REMOVED***
	if off >= len(msg) ***REMOVED***
		return "", off, errBaseLen
	***REMOVED***
	beginOff := off + 1
	endOff := beginOff + int(msg[off])
	if endOff > len(msg) ***REMOVED***
		return "", off, errCalcLen
	***REMOVED***
	return string(msg[beginOff:endOff]), endOff, nil
***REMOVED***

func skipText(msg []byte, off int) (int, error) ***REMOVED***
	if off >= len(msg) ***REMOVED***
		return off, errBaseLen
	***REMOVED***
	endOff := off + 1 + int(msg[off])
	if endOff > len(msg) ***REMOVED***
		return off, errCalcLen
	***REMOVED***
	return endOff, nil
***REMOVED***

func packBytes(msg []byte, field []byte) []byte ***REMOVED***
	return append(msg, field...)
***REMOVED***

func unpackBytes(msg []byte, off int, field []byte) (int, error) ***REMOVED***
	newOff := off + len(field)
	if newOff > len(msg) ***REMOVED***
		return off, errBaseLen
	***REMOVED***
	copy(field, msg[off:newOff])
	return newOff, nil
***REMOVED***

func skipBytes(msg []byte, off int, field []byte) (int, error) ***REMOVED***
	newOff := off + len(field)
	if newOff > len(msg) ***REMOVED***
		return off, errBaseLen
	***REMOVED***
	return newOff, nil
***REMOVED***

const nameLen = 255

// A Name is a non-encoded domain name. It is used instead of strings to avoid
// allocations.
type Name struct ***REMOVED***
	Data   [nameLen]byte
	Length uint8
***REMOVED***

// NewName creates a new Name from a string.
func NewName(name string) (Name, error) ***REMOVED***
	if len([]byte(name)) > nameLen ***REMOVED***
		return Name***REMOVED******REMOVED***, errCalcLen
	***REMOVED***
	n := Name***REMOVED***Length: uint8(len(name))***REMOVED***
	copy(n.Data[:], []byte(name))
	return n, nil
***REMOVED***

func (n Name) String() string ***REMOVED***
	return string(n.Data[:n.Length])
***REMOVED***

// pack packs a domain name.
//
// Domain names are a sequence of counted strings split at the dots. They end
// with a zero-length string. Compression can be used to reuse domain suffixes.
//
// The compression map will be updated with new domain suffixes. If compression
// is nil, compression will not be used.
func (n *Name) pack(msg []byte, compression map[string]int) ([]byte, error) ***REMOVED***
	oldMsg := msg

	// Add a trailing dot to canonicalize name.
	if n.Length == 0 || n.Data[n.Length-1] != '.' ***REMOVED***
		return oldMsg, errNonCanonicalName
	***REMOVED***

	// Allow root domain.
	if n.Data[0] == '.' && n.Length == 1 ***REMOVED***
		return append(msg, 0), nil
	***REMOVED***

	// Emit sequence of counted strings, chopping at dots.
	for i, begin := 0, 0; i < int(n.Length); i++ ***REMOVED***
		// Check for the end of the segment.
		if n.Data[i] == '.' ***REMOVED***
			// The two most significant bits have special meaning.
			// It isn't allowed for segments to be long enough to
			// need them.
			if i-begin >= 1<<6 ***REMOVED***
				return oldMsg, errSegTooLong
			***REMOVED***

			// Segments must have a non-zero length.
			if i-begin == 0 ***REMOVED***
				return oldMsg, errZeroSegLen
			***REMOVED***

			msg = append(msg, byte(i-begin))

			for j := begin; j < i; j++ ***REMOVED***
				msg = append(msg, n.Data[j])
			***REMOVED***

			begin = i + 1
			continue
		***REMOVED***

		// We can only compress domain suffixes starting with a new
		// segment. A pointer is two bytes with the two most significant
		// bits set to 1 to indicate that it is a pointer.
		if (i == 0 || n.Data[i-1] == '.') && compression != nil ***REMOVED***
			if ptr, ok := compression[string(n.Data[i:])]; ok ***REMOVED***
				// Hit. Emit a pointer instead of the rest of
				// the domain.
				return append(msg, byte(ptr>>8|0xC0), byte(ptr)), nil
			***REMOVED***

			// Miss. Add the suffix to the compression table if the
			// offset can be stored in the available 14 bytes.
			if len(msg) <= int(^uint16(0)>>2) ***REMOVED***
				compression[string(n.Data[i:])] = len(msg)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return append(msg, 0), nil
***REMOVED***

// unpack unpacks a domain name.
func (n *Name) unpack(msg []byte, off int) (int, error) ***REMOVED***
	// currOff is the current working offset.
	currOff := off

	// newOff is the offset where the next record will start. Pointers lead
	// to data that belongs to other names and thus doesn't count towards to
	// the usage of this name.
	newOff := off

	// ptr is the number of pointers followed.
	var ptr int

	// Name is a slice representation of the name data.
	name := n.Data[:0]

Loop:
	for ***REMOVED***
		if currOff >= len(msg) ***REMOVED***
			return off, errBaseLen
		***REMOVED***
		c := int(msg[currOff])
		currOff++
		switch c & 0xC0 ***REMOVED***
		case 0x00: // String segment
			if c == 0x00 ***REMOVED***
				// A zero length signals the end of the name.
				break Loop
			***REMOVED***
			endOff := currOff + c
			if endOff > len(msg) ***REMOVED***
				return off, errCalcLen
			***REMOVED***
			name = append(name, msg[currOff:endOff]...)
			name = append(name, '.')
			currOff = endOff
		case 0xC0: // Pointer
			if currOff >= len(msg) ***REMOVED***
				return off, errInvalidPtr
			***REMOVED***
			c1 := msg[currOff]
			currOff++
			if ptr == 0 ***REMOVED***
				newOff = currOff
			***REMOVED***
			// Don't follow too many pointers, maybe there's a loop.
			if ptr++; ptr > 10 ***REMOVED***
				return off, errTooManyPtr
			***REMOVED***
			currOff = (c^0xC0)<<8 | int(c1)
		default:
			// Prefixes 0x80 and 0x40 are reserved.
			return off, errReserved
		***REMOVED***
	***REMOVED***
	if len(name) == 0 ***REMOVED***
		name = append(name, '.')
	***REMOVED***
	if len(name) > len(n.Data) ***REMOVED***
		return off, errCalcLen
	***REMOVED***
	n.Length = uint8(len(name))
	if ptr == 0 ***REMOVED***
		newOff = currOff
	***REMOVED***
	return newOff, nil
***REMOVED***

func skipName(msg []byte, off int) (int, error) ***REMOVED***
	// newOff is the offset where the next record will start. Pointers lead
	// to data that belongs to other names and thus doesn't count towards to
	// the usage of this name.
	newOff := off

Loop:
	for ***REMOVED***
		if newOff >= len(msg) ***REMOVED***
			return off, errBaseLen
		***REMOVED***
		c := int(msg[newOff])
		newOff++
		switch c & 0xC0 ***REMOVED***
		case 0x00:
			if c == 0x00 ***REMOVED***
				// A zero length signals the end of the name.
				break Loop
			***REMOVED***
			// literal string
			newOff += c
			if newOff > len(msg) ***REMOVED***
				return off, errCalcLen
			***REMOVED***
		case 0xC0:
			// Pointer to somewhere else in msg.

			// Pointers are two bytes.
			newOff++

			// Don't follow the pointer as the data here has ended.
			break Loop
		default:
			// Prefixes 0x80 and 0x40 are reserved.
			return off, errReserved
		***REMOVED***
	***REMOVED***

	return newOff, nil
***REMOVED***

// A Question is a DNS query.
type Question struct ***REMOVED***
	Name  Name
	Type  Type
	Class Class
***REMOVED***

func (q *Question) pack(msg []byte, compression map[string]int) ([]byte, error) ***REMOVED***
	msg, err := q.Name.pack(msg, compression)
	if err != nil ***REMOVED***
		return msg, &nestedError***REMOVED***"Name", err***REMOVED***
	***REMOVED***
	msg = packType(msg, q.Type)
	return packClass(msg, q.Class), nil
***REMOVED***

func unpackResourceBody(msg []byte, off int, hdr ResourceHeader) (ResourceBody, int, error) ***REMOVED***
	var (
		r    ResourceBody
		err  error
		name string
	)
	switch hdr.Type ***REMOVED***
	case TypeA:
		var rb AResource
		rb, err = unpackAResource(msg, off)
		r = &rb
		name = "A"
	case TypeNS:
		var rb NSResource
		rb, err = unpackNSResource(msg, off)
		r = &rb
		name = "NS"
	case TypeCNAME:
		var rb CNAMEResource
		rb, err = unpackCNAMEResource(msg, off)
		r = &rb
		name = "CNAME"
	case TypeSOA:
		var rb SOAResource
		rb, err = unpackSOAResource(msg, off)
		r = &rb
		name = "SOA"
	case TypePTR:
		var rb PTRResource
		rb, err = unpackPTRResource(msg, off)
		r = &rb
		name = "PTR"
	case TypeMX:
		var rb MXResource
		rb, err = unpackMXResource(msg, off)
		r = &rb
		name = "MX"
	case TypeTXT:
		var rb TXTResource
		rb, err = unpackTXTResource(msg, off, hdr.Length)
		r = &rb
		name = "TXT"
	case TypeAAAA:
		var rb AAAAResource
		rb, err = unpackAAAAResource(msg, off)
		r = &rb
		name = "AAAA"
	case TypeSRV:
		var rb SRVResource
		rb, err = unpackSRVResource(msg, off)
		r = &rb
		name = "SRV"
	***REMOVED***
	if err != nil ***REMOVED***
		return nil, off, &nestedError***REMOVED***name + " record", err***REMOVED***
	***REMOVED***
	if r == nil ***REMOVED***
		return nil, off, errors.New("invalid resource type: " + string(hdr.Type+'0'))
	***REMOVED***
	return r, off + int(hdr.Length), nil
***REMOVED***

// A CNAMEResource is a CNAME Resource record.
type CNAMEResource struct ***REMOVED***
	CNAME Name
***REMOVED***

func (r *CNAMEResource) realType() Type ***REMOVED***
	return TypeCNAME
***REMOVED***

func (r *CNAMEResource) pack(msg []byte, compression map[string]int) ([]byte, error) ***REMOVED***
	return r.CNAME.pack(msg, compression)
***REMOVED***

func unpackCNAMEResource(msg []byte, off int) (CNAMEResource, error) ***REMOVED***
	var cname Name
	if _, err := cname.unpack(msg, off); err != nil ***REMOVED***
		return CNAMEResource***REMOVED******REMOVED***, err
	***REMOVED***
	return CNAMEResource***REMOVED***cname***REMOVED***, nil
***REMOVED***

// An MXResource is an MX Resource record.
type MXResource struct ***REMOVED***
	Pref uint16
	MX   Name
***REMOVED***

func (r *MXResource) realType() Type ***REMOVED***
	return TypeMX
***REMOVED***

func (r *MXResource) pack(msg []byte, compression map[string]int) ([]byte, error) ***REMOVED***
	oldMsg := msg
	msg = packUint16(msg, r.Pref)
	msg, err := r.MX.pack(msg, compression)
	if err != nil ***REMOVED***
		return oldMsg, &nestedError***REMOVED***"MXResource.MX", err***REMOVED***
	***REMOVED***
	return msg, nil
***REMOVED***

func unpackMXResource(msg []byte, off int) (MXResource, error) ***REMOVED***
	pref, off, err := unpackUint16(msg, off)
	if err != nil ***REMOVED***
		return MXResource***REMOVED******REMOVED***, &nestedError***REMOVED***"Pref", err***REMOVED***
	***REMOVED***
	var mx Name
	if _, err := mx.unpack(msg, off); err != nil ***REMOVED***
		return MXResource***REMOVED******REMOVED***, &nestedError***REMOVED***"MX", err***REMOVED***
	***REMOVED***
	return MXResource***REMOVED***pref, mx***REMOVED***, nil
***REMOVED***

// An NSResource is an NS Resource record.
type NSResource struct ***REMOVED***
	NS Name
***REMOVED***

func (r *NSResource) realType() Type ***REMOVED***
	return TypeNS
***REMOVED***

func (r *NSResource) pack(msg []byte, compression map[string]int) ([]byte, error) ***REMOVED***
	return r.NS.pack(msg, compression)
***REMOVED***

func unpackNSResource(msg []byte, off int) (NSResource, error) ***REMOVED***
	var ns Name
	if _, err := ns.unpack(msg, off); err != nil ***REMOVED***
		return NSResource***REMOVED******REMOVED***, err
	***REMOVED***
	return NSResource***REMOVED***ns***REMOVED***, nil
***REMOVED***

// A PTRResource is a PTR Resource record.
type PTRResource struct ***REMOVED***
	PTR Name
***REMOVED***

func (r *PTRResource) realType() Type ***REMOVED***
	return TypePTR
***REMOVED***

func (r *PTRResource) pack(msg []byte, compression map[string]int) ([]byte, error) ***REMOVED***
	return r.PTR.pack(msg, compression)
***REMOVED***

func unpackPTRResource(msg []byte, off int) (PTRResource, error) ***REMOVED***
	var ptr Name
	if _, err := ptr.unpack(msg, off); err != nil ***REMOVED***
		return PTRResource***REMOVED******REMOVED***, err
	***REMOVED***
	return PTRResource***REMOVED***ptr***REMOVED***, nil
***REMOVED***

// An SOAResource is an SOA Resource record.
type SOAResource struct ***REMOVED***
	NS      Name
	MBox    Name
	Serial  uint32
	Refresh uint32
	Retry   uint32
	Expire  uint32

	// MinTTL the is the default TTL of Resources records which did not
	// contain a TTL value and the TTL of negative responses. (RFC 2308
	// Section 4)
	MinTTL uint32
***REMOVED***

func (r *SOAResource) realType() Type ***REMOVED***
	return TypeSOA
***REMOVED***

func (r *SOAResource) pack(msg []byte, compression map[string]int) ([]byte, error) ***REMOVED***
	oldMsg := msg
	msg, err := r.NS.pack(msg, compression)
	if err != nil ***REMOVED***
		return oldMsg, &nestedError***REMOVED***"SOAResource.NS", err***REMOVED***
	***REMOVED***
	msg, err = r.MBox.pack(msg, compression)
	if err != nil ***REMOVED***
		return oldMsg, &nestedError***REMOVED***"SOAResource.MBox", err***REMOVED***
	***REMOVED***
	msg = packUint32(msg, r.Serial)
	msg = packUint32(msg, r.Refresh)
	msg = packUint32(msg, r.Retry)
	msg = packUint32(msg, r.Expire)
	return packUint32(msg, r.MinTTL), nil
***REMOVED***

func unpackSOAResource(msg []byte, off int) (SOAResource, error) ***REMOVED***
	var ns Name
	off, err := ns.unpack(msg, off)
	if err != nil ***REMOVED***
		return SOAResource***REMOVED******REMOVED***, &nestedError***REMOVED***"NS", err***REMOVED***
	***REMOVED***
	var mbox Name
	if off, err = mbox.unpack(msg, off); err != nil ***REMOVED***
		return SOAResource***REMOVED******REMOVED***, &nestedError***REMOVED***"MBox", err***REMOVED***
	***REMOVED***
	serial, off, err := unpackUint32(msg, off)
	if err != nil ***REMOVED***
		return SOAResource***REMOVED******REMOVED***, &nestedError***REMOVED***"Serial", err***REMOVED***
	***REMOVED***
	refresh, off, err := unpackUint32(msg, off)
	if err != nil ***REMOVED***
		return SOAResource***REMOVED******REMOVED***, &nestedError***REMOVED***"Refresh", err***REMOVED***
	***REMOVED***
	retry, off, err := unpackUint32(msg, off)
	if err != nil ***REMOVED***
		return SOAResource***REMOVED******REMOVED***, &nestedError***REMOVED***"Retry", err***REMOVED***
	***REMOVED***
	expire, off, err := unpackUint32(msg, off)
	if err != nil ***REMOVED***
		return SOAResource***REMOVED******REMOVED***, &nestedError***REMOVED***"Expire", err***REMOVED***
	***REMOVED***
	minTTL, _, err := unpackUint32(msg, off)
	if err != nil ***REMOVED***
		return SOAResource***REMOVED******REMOVED***, &nestedError***REMOVED***"MinTTL", err***REMOVED***
	***REMOVED***
	return SOAResource***REMOVED***ns, mbox, serial, refresh, retry, expire, minTTL***REMOVED***, nil
***REMOVED***

// A TXTResource is a TXT Resource record.
type TXTResource struct ***REMOVED***
	Txt string // Not a domain name.
***REMOVED***

func (r *TXTResource) realType() Type ***REMOVED***
	return TypeTXT
***REMOVED***

func (r *TXTResource) pack(msg []byte, compression map[string]int) ([]byte, error) ***REMOVED***
	return packText(msg, r.Txt), nil
***REMOVED***

func unpackTXTResource(msg []byte, off int, length uint16) (TXTResource, error) ***REMOVED***
	var txt string
	for n := uint16(0); n < length; ***REMOVED***
		var t string
		var err error
		if t, off, err = unpackText(msg, off); err != nil ***REMOVED***
			return TXTResource***REMOVED******REMOVED***, &nestedError***REMOVED***"text", err***REMOVED***
		***REMOVED***
		// Check if we got too many bytes.
		if length-n < uint16(len(t))+1 ***REMOVED***
			return TXTResource***REMOVED******REMOVED***, errCalcLen
		***REMOVED***
		n += uint16(len(t)) + 1
		txt += t
	***REMOVED***
	return TXTResource***REMOVED***txt***REMOVED***, nil
***REMOVED***

// An SRVResource is an SRV Resource record.
type SRVResource struct ***REMOVED***
	Priority uint16
	Weight   uint16
	Port     uint16
	Target   Name // Not compressed as per RFC 2782.
***REMOVED***

func (r *SRVResource) realType() Type ***REMOVED***
	return TypeSRV
***REMOVED***

func (r *SRVResource) pack(msg []byte, compression map[string]int) ([]byte, error) ***REMOVED***
	oldMsg := msg
	msg = packUint16(msg, r.Priority)
	msg = packUint16(msg, r.Weight)
	msg = packUint16(msg, r.Port)
	msg, err := r.Target.pack(msg, nil)
	if err != nil ***REMOVED***
		return oldMsg, &nestedError***REMOVED***"SRVResource.Target", err***REMOVED***
	***REMOVED***
	return msg, nil
***REMOVED***

func unpackSRVResource(msg []byte, off int) (SRVResource, error) ***REMOVED***
	priority, off, err := unpackUint16(msg, off)
	if err != nil ***REMOVED***
		return SRVResource***REMOVED******REMOVED***, &nestedError***REMOVED***"Priority", err***REMOVED***
	***REMOVED***
	weight, off, err := unpackUint16(msg, off)
	if err != nil ***REMOVED***
		return SRVResource***REMOVED******REMOVED***, &nestedError***REMOVED***"Weight", err***REMOVED***
	***REMOVED***
	port, off, err := unpackUint16(msg, off)
	if err != nil ***REMOVED***
		return SRVResource***REMOVED******REMOVED***, &nestedError***REMOVED***"Port", err***REMOVED***
	***REMOVED***
	var target Name
	if _, err := target.unpack(msg, off); err != nil ***REMOVED***
		return SRVResource***REMOVED******REMOVED***, &nestedError***REMOVED***"Target", err***REMOVED***
	***REMOVED***
	return SRVResource***REMOVED***priority, weight, port, target***REMOVED***, nil
***REMOVED***

// An AResource is an A Resource record.
type AResource struct ***REMOVED***
	A [4]byte
***REMOVED***

func (r *AResource) realType() Type ***REMOVED***
	return TypeA
***REMOVED***

func (r *AResource) pack(msg []byte, compression map[string]int) ([]byte, error) ***REMOVED***
	return packBytes(msg, r.A[:]), nil
***REMOVED***

func unpackAResource(msg []byte, off int) (AResource, error) ***REMOVED***
	var a [4]byte
	if _, err := unpackBytes(msg, off, a[:]); err != nil ***REMOVED***
		return AResource***REMOVED******REMOVED***, err
	***REMOVED***
	return AResource***REMOVED***a***REMOVED***, nil
***REMOVED***

// An AAAAResource is an AAAA Resource record.
type AAAAResource struct ***REMOVED***
	AAAA [16]byte
***REMOVED***

func (r *AAAAResource) realType() Type ***REMOVED***
	return TypeAAAA
***REMOVED***

func (r *AAAAResource) pack(msg []byte, compression map[string]int) ([]byte, error) ***REMOVED***
	return packBytes(msg, r.AAAA[:]), nil
***REMOVED***

func unpackAAAAResource(msg []byte, off int) (AAAAResource, error) ***REMOVED***
	var aaaa [16]byte
	if _, err := unpackBytes(msg, off, aaaa[:]); err != nil ***REMOVED***
		return AAAAResource***REMOVED******REMOVED***, err
	***REMOVED***
	return AAAAResource***REMOVED***aaaa***REMOVED***, nil
***REMOVED***
