// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webdav

// The XML encoding is covered by Section 14.
// http://www.webdav.org/specs/rfc4918.html#xml.element.definitions

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	// As of https://go-review.googlesource.com/#/c/12772/ which was submitted
	// in July 2015, this package uses an internal fork of the standard
	// library's encoding/xml package, due to changes in the way namespaces
	// were encoded. Such changes were introduced in the Go 1.5 cycle, but were
	// rolled back in response to https://github.com/golang/go/issues/11841
	//
	// However, this package's exported API, specifically the Property and
	// DeadPropsHolder types, need to refer to the standard library's version
	// of the xml.Name type, as code that imports this package cannot refer to
	// the internal version.
	//
	// This file therefore imports both the internal and external versions, as
	// ixml and xml, and converts between them.
	//
	// In the long term, this package should use the standard library's version
	// only, and the internal fork deleted, once
	// https://github.com/golang/go/issues/13400 is resolved.
	ixml "golang.org/x/net/webdav/internal/xml"
)

// http://www.webdav.org/specs/rfc4918.html#ELEMENT_lockinfo
type lockInfo struct ***REMOVED***
	XMLName   ixml.Name `xml:"lockinfo"`
	Exclusive *struct***REMOVED******REMOVED*** `xml:"lockscope>exclusive"`
	Shared    *struct***REMOVED******REMOVED*** `xml:"lockscope>shared"`
	Write     *struct***REMOVED******REMOVED*** `xml:"locktype>write"`
	Owner     owner     `xml:"owner"`
***REMOVED***

// http://www.webdav.org/specs/rfc4918.html#ELEMENT_owner
type owner struct ***REMOVED***
	InnerXML string `xml:",innerxml"`
***REMOVED***

func readLockInfo(r io.Reader) (li lockInfo, status int, err error) ***REMOVED***
	c := &countingReader***REMOVED***r: r***REMOVED***
	if err = ixml.NewDecoder(c).Decode(&li); err != nil ***REMOVED***
		if err == io.EOF ***REMOVED***
			if c.n == 0 ***REMOVED***
				// An empty body means to refresh the lock.
				// http://www.webdav.org/specs/rfc4918.html#refreshing-locks
				return lockInfo***REMOVED******REMOVED***, 0, nil
			***REMOVED***
			err = errInvalidLockInfo
		***REMOVED***
		return lockInfo***REMOVED******REMOVED***, http.StatusBadRequest, err
	***REMOVED***
	// We only support exclusive (non-shared) write locks. In practice, these are
	// the only types of locks that seem to matter.
	if li.Exclusive == nil || li.Shared != nil || li.Write == nil ***REMOVED***
		return lockInfo***REMOVED******REMOVED***, http.StatusNotImplemented, errUnsupportedLockInfo
	***REMOVED***
	return li, 0, nil
***REMOVED***

type countingReader struct ***REMOVED***
	n int
	r io.Reader
***REMOVED***

func (c *countingReader) Read(p []byte) (int, error) ***REMOVED***
	n, err := c.r.Read(p)
	c.n += n
	return n, err
***REMOVED***

func writeLockInfo(w io.Writer, token string, ld LockDetails) (int, error) ***REMOVED***
	depth := "infinity"
	if ld.ZeroDepth ***REMOVED***
		depth = "0"
	***REMOVED***
	timeout := ld.Duration / time.Second
	return fmt.Fprintf(w, "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n"+
		"<D:prop xmlns:D=\"DAV:\"><D:lockdiscovery><D:activelock>\n"+
		"	<D:locktype><D:write/></D:locktype>\n"+
		"	<D:lockscope><D:exclusive/></D:lockscope>\n"+
		"	<D:depth>%s</D:depth>\n"+
		"	<D:owner>%s</D:owner>\n"+
		"	<D:timeout>Second-%d</D:timeout>\n"+
		"	<D:locktoken><D:href>%s</D:href></D:locktoken>\n"+
		"	<D:lockroot><D:href>%s</D:href></D:lockroot>\n"+
		"</D:activelock></D:lockdiscovery></D:prop>",
		depth, ld.OwnerXML, timeout, escape(token), escape(ld.Root),
	)
***REMOVED***

func escape(s string) string ***REMOVED***
	for i := 0; i < len(s); i++ ***REMOVED***
		switch s[i] ***REMOVED***
		case '"', '&', '\'', '<', '>':
			b := bytes.NewBuffer(nil)
			ixml.EscapeText(b, []byte(s))
			return b.String()
		***REMOVED***
	***REMOVED***
	return s
***REMOVED***

// Next returns the next token, if any, in the XML stream of d.
// RFC 4918 requires to ignore comments, processing instructions
// and directives.
// http://www.webdav.org/specs/rfc4918.html#property_values
// http://www.webdav.org/specs/rfc4918.html#xml-extensibility
func next(d *ixml.Decoder) (ixml.Token, error) ***REMOVED***
	for ***REMOVED***
		t, err := d.Token()
		if err != nil ***REMOVED***
			return t, err
		***REMOVED***
		switch t.(type) ***REMOVED***
		case ixml.Comment, ixml.Directive, ixml.ProcInst:
			continue
		default:
			return t, nil
		***REMOVED***
	***REMOVED***
***REMOVED***

// http://www.webdav.org/specs/rfc4918.html#ELEMENT_prop (for propfind)
type propfindProps []xml.Name

// UnmarshalXML appends the property names enclosed within start to pn.
//
// It returns an error if start does not contain any properties or if
// properties contain values. Character data between properties is ignored.
func (pn *propfindProps) UnmarshalXML(d *ixml.Decoder, start ixml.StartElement) error ***REMOVED***
	for ***REMOVED***
		t, err := next(d)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		switch t.(type) ***REMOVED***
		case ixml.EndElement:
			if len(*pn) == 0 ***REMOVED***
				return fmt.Errorf("%s must not be empty", start.Name.Local)
			***REMOVED***
			return nil
		case ixml.StartElement:
			name := t.(ixml.StartElement).Name
			t, err = next(d)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			if _, ok := t.(ixml.EndElement); !ok ***REMOVED***
				return fmt.Errorf("unexpected token %T", t)
			***REMOVED***
			*pn = append(*pn, xml.Name(name))
		***REMOVED***
	***REMOVED***
***REMOVED***

// http://www.webdav.org/specs/rfc4918.html#ELEMENT_propfind
type propfind struct ***REMOVED***
	XMLName  ixml.Name     `xml:"DAV: propfind"`
	Allprop  *struct***REMOVED******REMOVED***     `xml:"DAV: allprop"`
	Propname *struct***REMOVED******REMOVED***     `xml:"DAV: propname"`
	Prop     propfindProps `xml:"DAV: prop"`
	Include  propfindProps `xml:"DAV: include"`
***REMOVED***

func readPropfind(r io.Reader) (pf propfind, status int, err error) ***REMOVED***
	c := countingReader***REMOVED***r: r***REMOVED***
	if err = ixml.NewDecoder(&c).Decode(&pf); err != nil ***REMOVED***
		if err == io.EOF ***REMOVED***
			if c.n == 0 ***REMOVED***
				// An empty body means to propfind allprop.
				// http://www.webdav.org/specs/rfc4918.html#METHOD_PROPFIND
				return propfind***REMOVED***Allprop: new(struct***REMOVED******REMOVED***)***REMOVED***, 0, nil
			***REMOVED***
			err = errInvalidPropfind
		***REMOVED***
		return propfind***REMOVED******REMOVED***, http.StatusBadRequest, err
	***REMOVED***

	if pf.Allprop == nil && pf.Include != nil ***REMOVED***
		return propfind***REMOVED******REMOVED***, http.StatusBadRequest, errInvalidPropfind
	***REMOVED***
	if pf.Allprop != nil && (pf.Prop != nil || pf.Propname != nil) ***REMOVED***
		return propfind***REMOVED******REMOVED***, http.StatusBadRequest, errInvalidPropfind
	***REMOVED***
	if pf.Prop != nil && pf.Propname != nil ***REMOVED***
		return propfind***REMOVED******REMOVED***, http.StatusBadRequest, errInvalidPropfind
	***REMOVED***
	if pf.Propname == nil && pf.Allprop == nil && pf.Prop == nil ***REMOVED***
		return propfind***REMOVED******REMOVED***, http.StatusBadRequest, errInvalidPropfind
	***REMOVED***
	return pf, 0, nil
***REMOVED***

// Property represents a single DAV resource property as defined in RFC 4918.
// See http://www.webdav.org/specs/rfc4918.html#data.model.for.resource.properties
type Property struct ***REMOVED***
	// XMLName is the fully qualified name that identifies this property.
	XMLName xml.Name

	// Lang is an optional xml:lang attribute.
	Lang string `xml:"xml:lang,attr,omitempty"`

	// InnerXML contains the XML representation of the property value.
	// See http://www.webdav.org/specs/rfc4918.html#property_values
	//
	// Property values of complex type or mixed-content must have fully
	// expanded XML namespaces or be self-contained with according
	// XML namespace declarations. They must not rely on any XML
	// namespace declarations within the scope of the XML document,
	// even including the DAV: namespace.
	InnerXML []byte `xml:",innerxml"`
***REMOVED***

// ixmlProperty is the same as the Property type except it holds an ixml.Name
// instead of an xml.Name.
type ixmlProperty struct ***REMOVED***
	XMLName  ixml.Name
	Lang     string `xml:"xml:lang,attr,omitempty"`
	InnerXML []byte `xml:",innerxml"`
***REMOVED***

// http://www.webdav.org/specs/rfc4918.html#ELEMENT_error
// See multistatusWriter for the "D:" namespace prefix.
type xmlError struct ***REMOVED***
	XMLName  ixml.Name `xml:"D:error"`
	InnerXML []byte    `xml:",innerxml"`
***REMOVED***

// http://www.webdav.org/specs/rfc4918.html#ELEMENT_propstat
// See multistatusWriter for the "D:" namespace prefix.
type propstat struct ***REMOVED***
	Prop                []Property `xml:"D:prop>_ignored_"`
	Status              string     `xml:"D:status"`
	Error               *xmlError  `xml:"D:error"`
	ResponseDescription string     `xml:"D:responsedescription,omitempty"`
***REMOVED***

// ixmlPropstat is the same as the propstat type except it holds an ixml.Name
// instead of an xml.Name.
type ixmlPropstat struct ***REMOVED***
	Prop                []ixmlProperty `xml:"D:prop>_ignored_"`
	Status              string         `xml:"D:status"`
	Error               *xmlError      `xml:"D:error"`
	ResponseDescription string         `xml:"D:responsedescription,omitempty"`
***REMOVED***

// MarshalXML prepends the "D:" namespace prefix on properties in the DAV: namespace
// before encoding. See multistatusWriter.
func (ps propstat) MarshalXML(e *ixml.Encoder, start ixml.StartElement) error ***REMOVED***
	// Convert from a propstat to an ixmlPropstat.
	ixmlPs := ixmlPropstat***REMOVED***
		Prop:                make([]ixmlProperty, len(ps.Prop)),
		Status:              ps.Status,
		Error:               ps.Error,
		ResponseDescription: ps.ResponseDescription,
	***REMOVED***
	for k, prop := range ps.Prop ***REMOVED***
		ixmlPs.Prop[k] = ixmlProperty***REMOVED***
			XMLName:  ixml.Name(prop.XMLName),
			Lang:     prop.Lang,
			InnerXML: prop.InnerXML,
		***REMOVED***
	***REMOVED***

	for k, prop := range ixmlPs.Prop ***REMOVED***
		if prop.XMLName.Space == "DAV:" ***REMOVED***
			prop.XMLName = ixml.Name***REMOVED***Space: "", Local: "D:" + prop.XMLName.Local***REMOVED***
			ixmlPs.Prop[k] = prop
		***REMOVED***
	***REMOVED***
	// Distinct type to avoid infinite recursion of MarshalXML.
	type newpropstat ixmlPropstat
	return e.EncodeElement(newpropstat(ixmlPs), start)
***REMOVED***

// http://www.webdav.org/specs/rfc4918.html#ELEMENT_response
// See multistatusWriter for the "D:" namespace prefix.
type response struct ***REMOVED***
	XMLName             ixml.Name  `xml:"D:response"`
	Href                []string   `xml:"D:href"`
	Propstat            []propstat `xml:"D:propstat"`
	Status              string     `xml:"D:status,omitempty"`
	Error               *xmlError  `xml:"D:error"`
	ResponseDescription string     `xml:"D:responsedescription,omitempty"`
***REMOVED***

// MultistatusWriter marshals one or more Responses into a XML
// multistatus response.
// See http://www.webdav.org/specs/rfc4918.html#ELEMENT_multistatus
// TODO(rsto, mpl): As a workaround, the "D:" namespace prefix, defined as
// "DAV:" on this element, is prepended on the nested response, as well as on all
// its nested elements. All property names in the DAV: namespace are prefixed as
// well. This is because some versions of Mini-Redirector (on windows 7) ignore
// elements with a default namespace (no prefixed namespace). A less intrusive fix
// should be possible after golang.org/cl/11074. See https://golang.org/issue/11177
type multistatusWriter struct ***REMOVED***
	// ResponseDescription contains the optional responsedescription
	// of the multistatus XML element. Only the latest content before
	// close will be emitted. Empty response descriptions are not
	// written.
	responseDescription string

	w   http.ResponseWriter
	enc *ixml.Encoder
***REMOVED***

// Write validates and emits a DAV response as part of a multistatus response
// element.
//
// It sets the HTTP status code of its underlying http.ResponseWriter to 207
// (Multi-Status) and populates the Content-Type header. If r is the
// first, valid response to be written, Write prepends the XML representation
// of r with a multistatus tag. Callers must call close after the last response
// has been written.
func (w *multistatusWriter) write(r *response) error ***REMOVED***
	switch len(r.Href) ***REMOVED***
	case 0:
		return errInvalidResponse
	case 1:
		if len(r.Propstat) > 0 != (r.Status == "") ***REMOVED***
			return errInvalidResponse
		***REMOVED***
	default:
		if len(r.Propstat) > 0 || r.Status == "" ***REMOVED***
			return errInvalidResponse
		***REMOVED***
	***REMOVED***
	err := w.writeHeader()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return w.enc.Encode(r)
***REMOVED***

// writeHeader writes a XML multistatus start element on w's underlying
// http.ResponseWriter and returns the result of the write operation.
// After the first write attempt, writeHeader becomes a no-op.
func (w *multistatusWriter) writeHeader() error ***REMOVED***
	if w.enc != nil ***REMOVED***
		return nil
	***REMOVED***
	w.w.Header().Add("Content-Type", "text/xml; charset=utf-8")
	w.w.WriteHeader(StatusMulti)
	_, err := fmt.Fprintf(w.w, `<?xml version="1.0" encoding="UTF-8"?>`)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	w.enc = ixml.NewEncoder(w.w)
	return w.enc.EncodeToken(ixml.StartElement***REMOVED***
		Name: ixml.Name***REMOVED***
			Space: "DAV:",
			Local: "multistatus",
		***REMOVED***,
		Attr: []ixml.Attr***REMOVED******REMOVED***
			Name:  ixml.Name***REMOVED***Space: "xmlns", Local: "D"***REMOVED***,
			Value: "DAV:",
		***REMOVED******REMOVED***,
	***REMOVED***)
***REMOVED***

// Close completes the marshalling of the multistatus response. It returns
// an error if the multistatus response could not be completed. If both the
// return value and field enc of w are nil, then no multistatus response has
// been written.
func (w *multistatusWriter) close() error ***REMOVED***
	if w.enc == nil ***REMOVED***
		return nil
	***REMOVED***
	var end []ixml.Token
	if w.responseDescription != "" ***REMOVED***
		name := ixml.Name***REMOVED***Space: "DAV:", Local: "responsedescription"***REMOVED***
		end = append(end,
			ixml.StartElement***REMOVED***Name: name***REMOVED***,
			ixml.CharData(w.responseDescription),
			ixml.EndElement***REMOVED***Name: name***REMOVED***,
		)
	***REMOVED***
	end = append(end, ixml.EndElement***REMOVED***
		Name: ixml.Name***REMOVED***Space: "DAV:", Local: "multistatus"***REMOVED***,
	***REMOVED***)
	for _, t := range end ***REMOVED***
		err := w.enc.EncodeToken(t)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return w.enc.Flush()
***REMOVED***

var xmlLangName = ixml.Name***REMOVED***Space: "http://www.w3.org/XML/1998/namespace", Local: "lang"***REMOVED***

func xmlLang(s ixml.StartElement, d string) string ***REMOVED***
	for _, attr := range s.Attr ***REMOVED***
		if attr.Name == xmlLangName ***REMOVED***
			return attr.Value
		***REMOVED***
	***REMOVED***
	return d
***REMOVED***

type xmlValue []byte

func (v *xmlValue) UnmarshalXML(d *ixml.Decoder, start ixml.StartElement) error ***REMOVED***
	// The XML value of a property can be arbitrary, mixed-content XML.
	// To make sure that the unmarshalled value contains all required
	// namespaces, we encode all the property value XML tokens into a
	// buffer. This forces the encoder to redeclare any used namespaces.
	var b bytes.Buffer
	e := ixml.NewEncoder(&b)
	for ***REMOVED***
		t, err := next(d)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if e, ok := t.(ixml.EndElement); ok && e.Name == start.Name ***REMOVED***
			break
		***REMOVED***
		if err = e.EncodeToken(t); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	err := e.Flush()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*v = b.Bytes()
	return nil
***REMOVED***

// http://www.webdav.org/specs/rfc4918.html#ELEMENT_prop (for proppatch)
type proppatchProps []Property

// UnmarshalXML appends the property names and values enclosed within start
// to ps.
//
// An xml:lang attribute that is defined either on the DAV:prop or property
// name XML element is propagated to the property's Lang field.
//
// UnmarshalXML returns an error if start does not contain any properties or if
// property values contain syntactically incorrect XML.
func (ps *proppatchProps) UnmarshalXML(d *ixml.Decoder, start ixml.StartElement) error ***REMOVED***
	lang := xmlLang(start, "")
	for ***REMOVED***
		t, err := next(d)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		switch elem := t.(type) ***REMOVED***
		case ixml.EndElement:
			if len(*ps) == 0 ***REMOVED***
				return fmt.Errorf("%s must not be empty", start.Name.Local)
			***REMOVED***
			return nil
		case ixml.StartElement:
			p := Property***REMOVED***
				XMLName: xml.Name(t.(ixml.StartElement).Name),
				Lang:    xmlLang(t.(ixml.StartElement), lang),
			***REMOVED***
			err = d.DecodeElement(((*xmlValue)(&p.InnerXML)), &elem)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			*ps = append(*ps, p)
		***REMOVED***
	***REMOVED***
***REMOVED***

// http://www.webdav.org/specs/rfc4918.html#ELEMENT_set
// http://www.webdav.org/specs/rfc4918.html#ELEMENT_remove
type setRemove struct ***REMOVED***
	XMLName ixml.Name
	Lang    string         `xml:"xml:lang,attr,omitempty"`
	Prop    proppatchProps `xml:"DAV: prop"`
***REMOVED***

// http://www.webdav.org/specs/rfc4918.html#ELEMENT_propertyupdate
type propertyupdate struct ***REMOVED***
	XMLName   ixml.Name   `xml:"DAV: propertyupdate"`
	Lang      string      `xml:"xml:lang,attr,omitempty"`
	SetRemove []setRemove `xml:",any"`
***REMOVED***

func readProppatch(r io.Reader) (patches []Proppatch, status int, err error) ***REMOVED***
	var pu propertyupdate
	if err = ixml.NewDecoder(r).Decode(&pu); err != nil ***REMOVED***
		return nil, http.StatusBadRequest, err
	***REMOVED***
	for _, op := range pu.SetRemove ***REMOVED***
		remove := false
		switch op.XMLName ***REMOVED***
		case ixml.Name***REMOVED***Space: "DAV:", Local: "set"***REMOVED***:
			// No-op.
		case ixml.Name***REMOVED***Space: "DAV:", Local: "remove"***REMOVED***:
			for _, p := range op.Prop ***REMOVED***
				if len(p.InnerXML) > 0 ***REMOVED***
					return nil, http.StatusBadRequest, errInvalidProppatch
				***REMOVED***
			***REMOVED***
			remove = true
		default:
			return nil, http.StatusBadRequest, errInvalidProppatch
		***REMOVED***
		patches = append(patches, Proppatch***REMOVED***Remove: remove, Props: op.Prop***REMOVED***)
	***REMOVED***
	return patches, 0, nil
***REMOVED***
