// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package message

// TODO: some types in this file will need to be made public at some time.
// Documentation and method names will reflect this by using the exported name.

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

// MatchLanguage reports the matched tag obtained from language.MatchStrings for
// the Matcher of the DefaultCatalog.
func MatchLanguage(preferred ...string) language.Tag ***REMOVED***
	c := DefaultCatalog
	tag, _ := language.MatchStrings(c.Matcher(), preferred...)
	return tag
***REMOVED***

// DefaultCatalog is used by SetString.
var DefaultCatalog catalog.Catalog = defaultCatalog

var defaultCatalog = catalog.NewBuilder()

// SetString calls SetString on the initial default Catalog.
func SetString(tag language.Tag, key string, msg string) error ***REMOVED***
	return defaultCatalog.SetString(tag, key, msg)
***REMOVED***

// Set calls Set on the initial default Catalog.
func Set(tag language.Tag, key string, msg ...catalog.Message) error ***REMOVED***
	return defaultCatalog.Set(tag, key, msg...)
***REMOVED***
