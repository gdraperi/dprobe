// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pipeline

import (
	"encoding/json"
	"errors"
	"strings"

	"golang.org/x/text/language"
)

// TODO: these definitions should be moved to a package so that the can be used
// by other tools.

// The file contains the structures used to define translations of a certain
// messages.
//
// A translation may have multiple translations strings, or messages, depending
// on the feature values of the various arguments. For instance, consider
// a hypothetical translation from English to English, where the source defines
// the format string "%d file(s) remaining".
// See the examples directory for examples of extracted messages.

// Messages is used to store translations for a single language.
type Messages struct ***REMOVED***
	Language language.Tag    `json:"language"`
	Messages []Message       `json:"messages"`
	Macros   map[string]Text `json:"macros,omitempty"`
***REMOVED***

// A Message describes a message to be translated.
type Message struct ***REMOVED***
	// ID contains a list of identifiers for the message.
	ID IDList `json:"id"`
	// Key is the string that is used to look up the message at runtime.
	Key         string `json:"key,omitempty"`
	Meaning     string `json:"meaning,omitempty"`
	Message     Text   `json:"message"`
	Translation Text   `json:"translation"`

	Comment           string `json:"comment,omitempty"`
	TranslatorComment string `json:"translatorComment,omitempty"`

	Placeholders []Placeholder `json:"placeholders,omitempty"`

	// Fuzzy indicates that the provide translation needs review by a
	// translator, for instance because it was derived from automated
	// translation.
	Fuzzy bool `json:"fuzzy,omitempty"`

	// TODO: default placeholder syntax is ***REMOVED***foo***REMOVED***. Allow alternative escaping
	// like `foo`.

	// Extraction information.
	Position string `json:"position,omitempty"` // filePosition:line
***REMOVED***

// Placeholder reports the placeholder for the given ID if it is defined or nil
// otherwise.
func (m *Message) Placeholder(id string) *Placeholder ***REMOVED***
	for _, p := range m.Placeholders ***REMOVED***
		if p.ID == id ***REMOVED***
			return &p
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Substitute replaces placeholders in msg with their original value.
func (m *Message) Substitute(msg string) (sub string, err error) ***REMOVED***
	last := 0
	for i := 0; i < len(msg); ***REMOVED***
		pLeft := strings.IndexByte(msg[i:], '***REMOVED***')
		if pLeft == -1 ***REMOVED***
			break
		***REMOVED***
		pLeft += i
		pRight := strings.IndexByte(msg[pLeft:], '***REMOVED***')
		if pRight == -1 ***REMOVED***
			return "", errorf("unmatched '***REMOVED***'")
		***REMOVED***
		pRight += pLeft
		id := strings.TrimSpace(msg[pLeft+1 : pRight])
		i = pRight + 1
		if id != "" && id[0] == '$' ***REMOVED***
			continue
		***REMOVED***
		sub += msg[last:pLeft]
		last = i
		ph := m.Placeholder(id)
		if ph == nil ***REMOVED***
			return "", errorf("unknown placeholder %q in message %q", id, msg)
		***REMOVED***
		sub += ph.String
	***REMOVED***
	sub += msg[last:]
	return sub, err
***REMOVED***

var errIncompatibleMessage = errors.New("messages incompatible")

func checkEquivalence(a, b *Message) error ***REMOVED***
	for _, v := range a.ID ***REMOVED***
		for _, w := range b.ID ***REMOVED***
			if v == w ***REMOVED***
				return nil
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// TODO: canonicalize placeholders and check for type equivalence.
	return errIncompatibleMessage
***REMOVED***

// A Placeholder is a part of the message that should not be changed by a
// translator. It can be used to hide or prettify format strings (e.g. %d or
// ***REMOVED******REMOVED***.Count***REMOVED******REMOVED***), hide HTML, or mark common names that should not be translated.
type Placeholder struct ***REMOVED***
	// ID is the placeholder identifier without the curly braces.
	ID string `json:"id"`

	// String is the string with which to replace the placeholder. This may be a
	// formatting string (for instance "%d" or "***REMOVED******REMOVED***.Count***REMOVED******REMOVED***") or a literal string
	// (<div>).
	String string `json:"string"`

	Type           string `json:"type"`
	UnderlyingType string `json:"underlyingType"`
	// ArgNum and Expr are set if the placeholder is a substitution of an
	// argument.
	ArgNum int    `json:"argNum,omitempty"`
	Expr   string `json:"expr,omitempty"`

	Comment string `json:"comment,omitempty"`
	Example string `json:"example,omitempty"`

	// Features contains the features that are available for the implementation
	// of this argument.
	Features []Feature `json:"features,omitempty"`
***REMOVED***

// An argument contains information about the arguments passed to a message.
type argument struct ***REMOVED***
	// ArgNum corresponds to the number that should be used for explicit argument indexes (e.g.
	// "%[1]d").
	ArgNum int `json:"argNum,omitempty"`

	used           bool   // Used by Placeholder
	Type           string `json:"type"`
	UnderlyingType string `json:"underlyingType"`
	Expr           string `json:"expr"`
	Value          string `json:"value,omitempty"`
	Comment        string `json:"comment,omitempty"`
	Position       string `json:"position,omitempty"`
***REMOVED***

// Feature holds information about a feature that can be implemented by
// an Argument.
type Feature struct ***REMOVED***
	Type string `json:"type"` // Right now this is only gender and plural.

	// TODO: possible values and examples for the language under consideration.

***REMOVED***

// Text defines a message to be displayed.
type Text struct ***REMOVED***
	// Msg and Select contains the message to be displayed. Msg may be used as
	// a fallback value if none of the select cases match.
	Msg    string  `json:"msg,omitempty"`
	Select *Select `json:"select,omitempty"`

	// Var defines a map of variables that may be substituted in the selected
	// message.
	Var map[string]Text `json:"var,omitempty"`

	// Example contains an example message formatted with default values.
	Example string `json:"example,omitempty"`
***REMOVED***

// IsEmpty reports whether this Text can generate anything.
func (t *Text) IsEmpty() bool ***REMOVED***
	return t.Msg == "" && t.Select == nil && t.Var == nil
***REMOVED***

// rawText erases the UnmarshalJSON method.
type rawText Text

// UnmarshalJSON implements json.Unmarshaler.
func (t *Text) UnmarshalJSON(b []byte) error ***REMOVED***
	if b[0] == '"' ***REMOVED***
		return json.Unmarshal(b, &t.Msg)
	***REMOVED***
	return json.Unmarshal(b, (*rawText)(t))
***REMOVED***

// MarshalJSON implements json.Marshaler.
func (t *Text) MarshalJSON() ([]byte, error) ***REMOVED***
	if t.Select == nil && t.Var == nil && t.Example == "" ***REMOVED***
		return json.Marshal(t.Msg)
	***REMOVED***
	return json.Marshal((*rawText)(t))
***REMOVED***

// IDList is a set identifiers that each may refer to possibly different
// versions of the same message. When looking up a messages, the first
// identifier in the list takes precedence.
type IDList []string

// UnmarshalJSON implements json.Unmarshaler.
func (id *IDList) UnmarshalJSON(b []byte) error ***REMOVED***
	if b[0] == '"' ***REMOVED***
		*id = []string***REMOVED***""***REMOVED***
		return json.Unmarshal(b, &((*id)[0]))
	***REMOVED***
	return json.Unmarshal(b, (*[]string)(id))
***REMOVED***

// MarshalJSON implements json.Marshaler.
func (id *IDList) MarshalJSON() ([]byte, error) ***REMOVED***
	if len(*id) == 1 ***REMOVED***
		return json.Marshal((*id)[0])
	***REMOVED***
	return json.Marshal((*[]string)(id))
***REMOVED***

// Select selects a Text based on the feature value associated with a feature of
// a certain argument.
type Select struct ***REMOVED***
	Feature string          `json:"feature"` // Name of Feature type (e.g plural)
	Arg     string          `json:"arg"`     // The placeholder ID
	Cases   map[string]Text `json:"cases"`
***REMOVED***

// TODO: order matters, but can we derive the ordering from the case keys?
// type Case struct ***REMOVED***
// 	Key   string `json:"key"`
// 	Value Text   `json:"value"`
// ***REMOVED***
