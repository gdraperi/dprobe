// Copyright 2015 Unknwon
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package ini

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type tokenType int

const (
	_TOKEN_INVALID tokenType = iota
	_TOKEN_COMMENT
	_TOKEN_SECTION
	_TOKEN_KEY
)

type parser struct ***REMOVED***
	buf     *bufio.Reader
	isEOF   bool
	count   int
	comment *bytes.Buffer
***REMOVED***

func newParser(r io.Reader) *parser ***REMOVED***
	return &parser***REMOVED***
		buf:     bufio.NewReader(r),
		count:   1,
		comment: &bytes.Buffer***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

// BOM handles header of UTF-8, UTF-16 LE and UTF-16 BE's BOM format.
// http://en.wikipedia.org/wiki/Byte_order_mark#Representations_of_byte_order_marks_by_encoding
func (p *parser) BOM() error ***REMOVED***
	mask, err := p.buf.Peek(2)
	if err != nil && err != io.EOF ***REMOVED***
		return err
	***REMOVED*** else if len(mask) < 2 ***REMOVED***
		return nil
	***REMOVED***

	switch ***REMOVED***
	case mask[0] == 254 && mask[1] == 255:
		fallthrough
	case mask[0] == 255 && mask[1] == 254:
		p.buf.Read(mask)
	case mask[0] == 239 && mask[1] == 187:
		mask, err := p.buf.Peek(3)
		if err != nil && err != io.EOF ***REMOVED***
			return err
		***REMOVED*** else if len(mask) < 3 ***REMOVED***
			return nil
		***REMOVED***
		if mask[2] == 191 ***REMOVED***
			p.buf.Read(mask)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (p *parser) readUntil(delim byte) ([]byte, error) ***REMOVED***
	data, err := p.buf.ReadBytes(delim)
	if err != nil ***REMOVED***
		if err == io.EOF ***REMOVED***
			p.isEOF = true
		***REMOVED*** else ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return data, nil
***REMOVED***

func cleanComment(in []byte) ([]byte, bool) ***REMOVED***
	i := bytes.IndexAny(in, "#;")
	if i == -1 ***REMOVED***
		return nil, false
	***REMOVED***
	return in[i:], true
***REMOVED***

func readKeyName(in []byte) (string, int, error) ***REMOVED***
	line := string(in)

	// Check if key name surrounded by quotes.
	var keyQuote string
	if line[0] == '"' ***REMOVED***
		if len(line) > 6 && string(line[0:3]) == `"""` ***REMOVED***
			keyQuote = `"""`
		***REMOVED*** else ***REMOVED***
			keyQuote = `"`
		***REMOVED***
	***REMOVED*** else if line[0] == '`' ***REMOVED***
		keyQuote = "`"
	***REMOVED***

	// Get out key name
	endIdx := -1
	if len(keyQuote) > 0 ***REMOVED***
		startIdx := len(keyQuote)
		// FIXME: fail case -> """"""name"""=value
		pos := strings.Index(line[startIdx:], keyQuote)
		if pos == -1 ***REMOVED***
			return "", -1, fmt.Errorf("missing closing key quote: %s", line)
		***REMOVED***
		pos += startIdx

		// Find key-value delimiter
		i := strings.IndexAny(line[pos+startIdx:], "=:")
		if i < 0 ***REMOVED***
			return "", -1, ErrDelimiterNotFound***REMOVED***line***REMOVED***
		***REMOVED***
		endIdx = pos + i
		return strings.TrimSpace(line[startIdx:pos]), endIdx + startIdx + 1, nil
	***REMOVED***

	endIdx = strings.IndexAny(line, "=:")
	if endIdx < 0 ***REMOVED***
		return "", -1, ErrDelimiterNotFound***REMOVED***line***REMOVED***
	***REMOVED***
	return strings.TrimSpace(line[0:endIdx]), endIdx + 1, nil
***REMOVED***

func (p *parser) readMultilines(line, val, valQuote string) (string, error) ***REMOVED***
	for ***REMOVED***
		data, err := p.readUntil('\n')
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		next := string(data)

		pos := strings.LastIndex(next, valQuote)
		if pos > -1 ***REMOVED***
			val += next[:pos]

			comment, has := cleanComment([]byte(next[pos:]))
			if has ***REMOVED***
				p.comment.Write(bytes.TrimSpace(comment))
			***REMOVED***
			break
		***REMOVED***
		val += next
		if p.isEOF ***REMOVED***
			return "", fmt.Errorf("missing closing key quote from '%s' to '%s'", line, next)
		***REMOVED***
	***REMOVED***
	return val, nil
***REMOVED***

func (p *parser) readContinuationLines(val string) (string, error) ***REMOVED***
	for ***REMOVED***
		data, err := p.readUntil('\n')
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		next := strings.TrimSpace(string(data))

		if len(next) == 0 ***REMOVED***
			break
		***REMOVED***
		val += next
		if val[len(val)-1] != '\\' ***REMOVED***
			break
		***REMOVED***
		val = val[:len(val)-1]
	***REMOVED***
	return val, nil
***REMOVED***

// hasSurroundedQuote check if and only if the first and last characters
// are quotes \" or \'.
// It returns false if any other parts also contain same kind of quotes.
func hasSurroundedQuote(in string, quote byte) bool ***REMOVED***
	return len(in) > 2 && in[0] == quote && in[len(in)-1] == quote &&
		strings.IndexByte(in[1:], quote) == len(in)-2
***REMOVED***

func (p *parser) readValue(in []byte, ignoreContinuation bool) (string, error) ***REMOVED***
	line := strings.TrimLeftFunc(string(in), unicode.IsSpace)
	if len(line) == 0 ***REMOVED***
		return "", nil
	***REMOVED***

	var valQuote string
	if len(line) > 3 && string(line[0:3]) == `"""` ***REMOVED***
		valQuote = `"""`
	***REMOVED*** else if line[0] == '`' ***REMOVED***
		valQuote = "`"
	***REMOVED***

	if len(valQuote) > 0 ***REMOVED***
		startIdx := len(valQuote)
		pos := strings.LastIndex(line[startIdx:], valQuote)
		// Check for multi-line value
		if pos == -1 ***REMOVED***
			return p.readMultilines(line, line[startIdx:], valQuote)
		***REMOVED***

		return line[startIdx : pos+startIdx], nil
	***REMOVED***

	// Won't be able to reach here if value only contains whitespace.
	line = strings.TrimSpace(line)

	// Check continuation lines when desired.
	if !ignoreContinuation && line[len(line)-1] == '\\' ***REMOVED***
		return p.readContinuationLines(line[:len(line)-1])
	***REMOVED***

	i := strings.IndexAny(line, "#;")
	if i > -1 ***REMOVED***
		p.comment.WriteString(line[i:])
		line = strings.TrimSpace(line[:i])
	***REMOVED***

	// Trim single quotes
	if hasSurroundedQuote(line, '\'') ||
		hasSurroundedQuote(line, '"') ***REMOVED***
		line = line[1 : len(line)-1]
	***REMOVED***
	return line, nil
***REMOVED***

// parse parses data through an io.Reader.
func (f *File) parse(reader io.Reader) (err error) ***REMOVED***
	p := newParser(reader)
	if err = p.BOM(); err != nil ***REMOVED***
		return fmt.Errorf("BOM: %v", err)
	***REMOVED***

	// Ignore error because default section name is never empty string.
	section, _ := f.NewSection(DEFAULT_SECTION)

	var line []byte
	var inUnparseableSection bool
	for !p.isEOF ***REMOVED***
		line, err = p.readUntil('\n')
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		line = bytes.TrimLeftFunc(line, unicode.IsSpace)
		if len(line) == 0 ***REMOVED***
			continue
		***REMOVED***

		// Comments
		if line[0] == '#' || line[0] == ';' ***REMOVED***
			// Note: we do not care ending line break,
			// it is needed for adding second line,
			// so just clean it once at the end when set to value.
			p.comment.Write(line)
			continue
		***REMOVED***

		// Section
		if line[0] == '[' ***REMOVED***
			// Read to the next ']' (TODO: support quoted strings)
			// TODO(unknwon): use LastIndexByte when stop supporting Go1.4
			closeIdx := bytes.LastIndex(line, []byte("]"))
			if closeIdx == -1 ***REMOVED***
				return fmt.Errorf("unclosed section: %s", line)
			***REMOVED***

			name := string(line[1:closeIdx])
			section, err = f.NewSection(name)
			if err != nil ***REMOVED***
				return err
			***REMOVED***

			comment, has := cleanComment(line[closeIdx+1:])
			if has ***REMOVED***
				p.comment.Write(comment)
			***REMOVED***

			section.Comment = strings.TrimSpace(p.comment.String())

			// Reset aotu-counter and comments
			p.comment.Reset()
			p.count = 1

			inUnparseableSection = false
			for i := range f.options.UnparseableSections ***REMOVED***
				if f.options.UnparseableSections[i] == name ||
					(f.options.Insensitive && strings.ToLower(f.options.UnparseableSections[i]) == strings.ToLower(name)) ***REMOVED***
					inUnparseableSection = true
					continue
				***REMOVED***
			***REMOVED***
			continue
		***REMOVED***

		if inUnparseableSection ***REMOVED***
			section.isRawSection = true
			section.rawBody += string(line)
			continue
		***REMOVED***

		kname, offset, err := readKeyName(line)
		if err != nil ***REMOVED***
			// Treat as boolean key when desired, and whole line is key name.
			if IsErrDelimiterNotFound(err) && f.options.AllowBooleanKeys ***REMOVED***
				kname, err := p.readValue(line, f.options.IgnoreContinuation)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				key, err := section.NewBooleanKey(kname)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				key.Comment = strings.TrimSpace(p.comment.String())
				p.comment.Reset()
				continue
			***REMOVED***
			return err
		***REMOVED***

		// Auto increment.
		isAutoIncr := false
		if kname == "-" ***REMOVED***
			isAutoIncr = true
			kname = "#" + strconv.Itoa(p.count)
			p.count++
		***REMOVED***

		value, err := p.readValue(line[offset:], f.options.IgnoreContinuation)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		key, err := section.NewKey(kname, value)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		key.isAutoIncrement = isAutoIncr
		key.Comment = strings.TrimSpace(p.comment.String())
		p.comment.Reset()
	***REMOVED***
	return nil
***REMOVED***
