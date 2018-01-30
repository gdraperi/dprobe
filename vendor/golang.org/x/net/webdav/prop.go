// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webdav

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/net/context"
)

// Proppatch describes a property update instruction as defined in RFC 4918.
// See http://www.webdav.org/specs/rfc4918.html#METHOD_PROPPATCH
type Proppatch struct ***REMOVED***
	// Remove specifies whether this patch removes properties. If it does not
	// remove them, it sets them.
	Remove bool
	// Props contains the properties to be set or removed.
	Props []Property
***REMOVED***

// Propstat describes a XML propstat element as defined in RFC 4918.
// See http://www.webdav.org/specs/rfc4918.html#ELEMENT_propstat
type Propstat struct ***REMOVED***
	// Props contains the properties for which Status applies.
	Props []Property

	// Status defines the HTTP status code of the properties in Prop.
	// Allowed values include, but are not limited to the WebDAV status
	// code extensions for HTTP/1.1.
	// http://www.webdav.org/specs/rfc4918.html#status.code.extensions.to.http11
	Status int

	// XMLError contains the XML representation of the optional error element.
	// XML content within this field must not rely on any predefined
	// namespace declarations or prefixes. If empty, the XML error element
	// is omitted.
	XMLError string

	// ResponseDescription contains the contents of the optional
	// responsedescription field. If empty, the XML element is omitted.
	ResponseDescription string
***REMOVED***

// makePropstats returns a slice containing those of x and y whose Props slice
// is non-empty. If both are empty, it returns a slice containing an otherwise
// zero Propstat whose HTTP status code is 200 OK.
func makePropstats(x, y Propstat) []Propstat ***REMOVED***
	pstats := make([]Propstat, 0, 2)
	if len(x.Props) != 0 ***REMOVED***
		pstats = append(pstats, x)
	***REMOVED***
	if len(y.Props) != 0 ***REMOVED***
		pstats = append(pstats, y)
	***REMOVED***
	if len(pstats) == 0 ***REMOVED***
		pstats = append(pstats, Propstat***REMOVED***
			Status: http.StatusOK,
		***REMOVED***)
	***REMOVED***
	return pstats
***REMOVED***

// DeadPropsHolder holds the dead properties of a resource.
//
// Dead properties are those properties that are explicitly defined. In
// comparison, live properties, such as DAV:getcontentlength, are implicitly
// defined by the underlying resource, and cannot be explicitly overridden or
// removed. See the Terminology section of
// http://www.webdav.org/specs/rfc4918.html#rfc.section.3
//
// There is a whitelist of the names of live properties. This package handles
// all live properties, and will only pass non-whitelisted names to the Patch
// method of DeadPropsHolder implementations.
type DeadPropsHolder interface ***REMOVED***
	// DeadProps returns a copy of the dead properties held.
	DeadProps() (map[xml.Name]Property, error)

	// Patch patches the dead properties held.
	//
	// Patching is atomic; either all or no patches succeed. It returns (nil,
	// non-nil) if an internal server error occurred, otherwise the Propstats
	// collectively contain one Property for each proposed patch Property. If
	// all patches succeed, Patch returns a slice of length one and a Propstat
	// element with a 200 OK HTTP status code. If none succeed, for reasons
	// other than an internal server error, no Propstat has status 200 OK.
	//
	// For more details on when various HTTP status codes apply, see
	// http://www.webdav.org/specs/rfc4918.html#PROPPATCH-status
	Patch([]Proppatch) ([]Propstat, error)
***REMOVED***

// liveProps contains all supported, protected DAV: properties.
var liveProps = map[xml.Name]struct ***REMOVED***
	// findFn implements the propfind function of this property. If nil,
	// it indicates a hidden property.
	findFn func(context.Context, FileSystem, LockSystem, string, os.FileInfo) (string, error)
	// dir is true if the property applies to directories.
	dir bool
***REMOVED******REMOVED***
	***REMOVED***Space: "DAV:", Local: "resourcetype"***REMOVED***: ***REMOVED***
		findFn: findResourceType,
		dir:    true,
	***REMOVED***,
	***REMOVED***Space: "DAV:", Local: "displayname"***REMOVED***: ***REMOVED***
		findFn: findDisplayName,
		dir:    true,
	***REMOVED***,
	***REMOVED***Space: "DAV:", Local: "getcontentlength"***REMOVED***: ***REMOVED***
		findFn: findContentLength,
		dir:    false,
	***REMOVED***,
	***REMOVED***Space: "DAV:", Local: "getlastmodified"***REMOVED***: ***REMOVED***
		findFn: findLastModified,
		// http://webdav.org/specs/rfc4918.html#PROPERTY_getlastmodified
		// suggests that getlastmodified should only apply to GETable
		// resources, and this package does not support GET on directories.
		//
		// Nonetheless, some WebDAV clients expect child directories to be
		// sortable by getlastmodified date, so this value is true, not false.
		// See golang.org/issue/15334.
		dir: true,
	***REMOVED***,
	***REMOVED***Space: "DAV:", Local: "creationdate"***REMOVED***: ***REMOVED***
		findFn: nil,
		dir:    false,
	***REMOVED***,
	***REMOVED***Space: "DAV:", Local: "getcontentlanguage"***REMOVED***: ***REMOVED***
		findFn: nil,
		dir:    false,
	***REMOVED***,
	***REMOVED***Space: "DAV:", Local: "getcontenttype"***REMOVED***: ***REMOVED***
		findFn: findContentType,
		dir:    false,
	***REMOVED***,
	***REMOVED***Space: "DAV:", Local: "getetag"***REMOVED***: ***REMOVED***
		findFn: findETag,
		// findETag implements ETag as the concatenated hex values of a file's
		// modification time and size. This is not a reliable synchronization
		// mechanism for directories, so we do not advertise getetag for DAV
		// collections.
		dir: false,
	***REMOVED***,

	// TODO: The lockdiscovery property requires LockSystem to list the
	// active locks on a resource.
	***REMOVED***Space: "DAV:", Local: "lockdiscovery"***REMOVED***: ***REMOVED******REMOVED***,
	***REMOVED***Space: "DAV:", Local: "supportedlock"***REMOVED***: ***REMOVED***
		findFn: findSupportedLock,
		dir:    true,
	***REMOVED***,
***REMOVED***

// TODO(nigeltao) merge props and allprop?

// Props returns the status of the properties named pnames for resource name.
//
// Each Propstat has a unique status and each property name will only be part
// of one Propstat element.
func props(ctx context.Context, fs FileSystem, ls LockSystem, name string, pnames []xml.Name) ([]Propstat, error) ***REMOVED***
	f, err := fs.OpenFile(ctx, name, os.O_RDONLY, 0)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer f.Close()
	fi, err := f.Stat()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	isDir := fi.IsDir()

	var deadProps map[xml.Name]Property
	if dph, ok := f.(DeadPropsHolder); ok ***REMOVED***
		deadProps, err = dph.DeadProps()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	pstatOK := Propstat***REMOVED***Status: http.StatusOK***REMOVED***
	pstatNotFound := Propstat***REMOVED***Status: http.StatusNotFound***REMOVED***
	for _, pn := range pnames ***REMOVED***
		// If this file has dead properties, check if they contain pn.
		if dp, ok := deadProps[pn]; ok ***REMOVED***
			pstatOK.Props = append(pstatOK.Props, dp)
			continue
		***REMOVED***
		// Otherwise, it must either be a live property or we don't know it.
		if prop := liveProps[pn]; prop.findFn != nil && (prop.dir || !isDir) ***REMOVED***
			innerXML, err := prop.findFn(ctx, fs, ls, name, fi)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			pstatOK.Props = append(pstatOK.Props, Property***REMOVED***
				XMLName:  pn,
				InnerXML: []byte(innerXML),
			***REMOVED***)
		***REMOVED*** else ***REMOVED***
			pstatNotFound.Props = append(pstatNotFound.Props, Property***REMOVED***
				XMLName: pn,
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	return makePropstats(pstatOK, pstatNotFound), nil
***REMOVED***

// Propnames returns the property names defined for resource name.
func propnames(ctx context.Context, fs FileSystem, ls LockSystem, name string) ([]xml.Name, error) ***REMOVED***
	f, err := fs.OpenFile(ctx, name, os.O_RDONLY, 0)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer f.Close()
	fi, err := f.Stat()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	isDir := fi.IsDir()

	var deadProps map[xml.Name]Property
	if dph, ok := f.(DeadPropsHolder); ok ***REMOVED***
		deadProps, err = dph.DeadProps()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	pnames := make([]xml.Name, 0, len(liveProps)+len(deadProps))
	for pn, prop := range liveProps ***REMOVED***
		if prop.findFn != nil && (prop.dir || !isDir) ***REMOVED***
			pnames = append(pnames, pn)
		***REMOVED***
	***REMOVED***
	for pn := range deadProps ***REMOVED***
		pnames = append(pnames, pn)
	***REMOVED***
	return pnames, nil
***REMOVED***

// Allprop returns the properties defined for resource name and the properties
// named in include.
//
// Note that RFC 4918 defines 'allprop' to return the DAV: properties defined
// within the RFC plus dead properties. Other live properties should only be
// returned if they are named in 'include'.
//
// See http://www.webdav.org/specs/rfc4918.html#METHOD_PROPFIND
func allprop(ctx context.Context, fs FileSystem, ls LockSystem, name string, include []xml.Name) ([]Propstat, error) ***REMOVED***
	pnames, err := propnames(ctx, fs, ls, name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// Add names from include if they are not already covered in pnames.
	nameset := make(map[xml.Name]bool)
	for _, pn := range pnames ***REMOVED***
		nameset[pn] = true
	***REMOVED***
	for _, pn := range include ***REMOVED***
		if !nameset[pn] ***REMOVED***
			pnames = append(pnames, pn)
		***REMOVED***
	***REMOVED***
	return props(ctx, fs, ls, name, pnames)
***REMOVED***

// Patch patches the properties of resource name. The return values are
// constrained in the same manner as DeadPropsHolder.Patch.
func patch(ctx context.Context, fs FileSystem, ls LockSystem, name string, patches []Proppatch) ([]Propstat, error) ***REMOVED***
	conflict := false
loop:
	for _, patch := range patches ***REMOVED***
		for _, p := range patch.Props ***REMOVED***
			if _, ok := liveProps[p.XMLName]; ok ***REMOVED***
				conflict = true
				break loop
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if conflict ***REMOVED***
		pstatForbidden := Propstat***REMOVED***
			Status:   http.StatusForbidden,
			XMLError: `<D:cannot-modify-protected-property xmlns:D="DAV:"/>`,
		***REMOVED***
		pstatFailedDep := Propstat***REMOVED***
			Status: StatusFailedDependency,
		***REMOVED***
		for _, patch := range patches ***REMOVED***
			for _, p := range patch.Props ***REMOVED***
				if _, ok := liveProps[p.XMLName]; ok ***REMOVED***
					pstatForbidden.Props = append(pstatForbidden.Props, Property***REMOVED***XMLName: p.XMLName***REMOVED***)
				***REMOVED*** else ***REMOVED***
					pstatFailedDep.Props = append(pstatFailedDep.Props, Property***REMOVED***XMLName: p.XMLName***REMOVED***)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return makePropstats(pstatForbidden, pstatFailedDep), nil
	***REMOVED***

	f, err := fs.OpenFile(ctx, name, os.O_RDWR, 0)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer f.Close()
	if dph, ok := f.(DeadPropsHolder); ok ***REMOVED***
		ret, err := dph.Patch(patches)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		// http://www.webdav.org/specs/rfc4918.html#ELEMENT_propstat says that
		// "The contents of the prop XML element must only list the names of
		// properties to which the result in the status element applies."
		for _, pstat := range ret ***REMOVED***
			for i, p := range pstat.Props ***REMOVED***
				pstat.Props[i] = Property***REMOVED***XMLName: p.XMLName***REMOVED***
			***REMOVED***
		***REMOVED***
		return ret, nil
	***REMOVED***
	// The file doesn't implement the optional DeadPropsHolder interface, so
	// all patches are forbidden.
	pstat := Propstat***REMOVED***Status: http.StatusForbidden***REMOVED***
	for _, patch := range patches ***REMOVED***
		for _, p := range patch.Props ***REMOVED***
			pstat.Props = append(pstat.Props, Property***REMOVED***XMLName: p.XMLName***REMOVED***)
		***REMOVED***
	***REMOVED***
	return []Propstat***REMOVED***pstat***REMOVED***, nil
***REMOVED***

func escapeXML(s string) string ***REMOVED***
	for i := 0; i < len(s); i++ ***REMOVED***
		// As an optimization, if s contains only ASCII letters, digits or a
		// few special characters, the escaped value is s itself and we don't
		// need to allocate a buffer and convert between string and []byte.
		switch c := s[i]; ***REMOVED***
		case c == ' ' || c == '_' ||
			('+' <= c && c <= '9') || // Digits as well as + , - . and /
			('A' <= c && c <= 'Z') ||
			('a' <= c && c <= 'z'):
			continue
		***REMOVED***
		// Otherwise, go through the full escaping process.
		var buf bytes.Buffer
		xml.EscapeText(&buf, []byte(s))
		return buf.String()
	***REMOVED***
	return s
***REMOVED***

func findResourceType(ctx context.Context, fs FileSystem, ls LockSystem, name string, fi os.FileInfo) (string, error) ***REMOVED***
	if fi.IsDir() ***REMOVED***
		return `<D:collection xmlns:D="DAV:"/>`, nil
	***REMOVED***
	return "", nil
***REMOVED***

func findDisplayName(ctx context.Context, fs FileSystem, ls LockSystem, name string, fi os.FileInfo) (string, error) ***REMOVED***
	if slashClean(name) == "/" ***REMOVED***
		// Hide the real name of a possibly prefixed root directory.
		return "", nil
	***REMOVED***
	return escapeXML(fi.Name()), nil
***REMOVED***

func findContentLength(ctx context.Context, fs FileSystem, ls LockSystem, name string, fi os.FileInfo) (string, error) ***REMOVED***
	return strconv.FormatInt(fi.Size(), 10), nil
***REMOVED***

func findLastModified(ctx context.Context, fs FileSystem, ls LockSystem, name string, fi os.FileInfo) (string, error) ***REMOVED***
	return fi.ModTime().Format(http.TimeFormat), nil
***REMOVED***

func findContentType(ctx context.Context, fs FileSystem, ls LockSystem, name string, fi os.FileInfo) (string, error) ***REMOVED***
	f, err := fs.OpenFile(ctx, name, os.O_RDONLY, 0)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	defer f.Close()
	// This implementation is based on serveContent's code in the standard net/http package.
	ctype := mime.TypeByExtension(filepath.Ext(name))
	if ctype != "" ***REMOVED***
		return ctype, nil
	***REMOVED***
	// Read a chunk to decide between utf-8 text and binary.
	var buf [512]byte
	n, err := io.ReadFull(f, buf[:])
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF ***REMOVED***
		return "", err
	***REMOVED***
	ctype = http.DetectContentType(buf[:n])
	// Rewind file.
	_, err = f.Seek(0, os.SEEK_SET)
	return ctype, err
***REMOVED***

func findETag(ctx context.Context, fs FileSystem, ls LockSystem, name string, fi os.FileInfo) (string, error) ***REMOVED***
	// The Apache http 2.4 web server by default concatenates the
	// modification time and size of a file. We replicate the heuristic
	// with nanosecond granularity.
	return fmt.Sprintf(`"%x%x"`, fi.ModTime().UnixNano(), fi.Size()), nil
***REMOVED***

func findSupportedLock(ctx context.Context, fs FileSystem, ls LockSystem, name string, fi os.FileInfo) (string, error) ***REMOVED***
	return `` +
		`<D:lockentry xmlns:D="DAV:">` +
		`<D:lockscope><D:exclusive/></D:lockscope>` +
		`<D:locktype><D:write/></D:locktype>` +
		`</D:lockentry>`, nil
***REMOVED***
