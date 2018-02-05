// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package catalog

import (
	"sync"

	"golang.org/x/text/internal"
	"golang.org/x/text/internal/catmsg"
	"golang.org/x/text/language"
)

// TODO:
// Dictionary returns a Dictionary that returns the first Message, using the
// given language tag, that matches:
//   1. the last one registered by one of the Set methods
//   2. returned by one of the Loaders
//   3. repeat from 1. using the parent language
// This approach allows messages to be underspecified.
// func (c *Catalog) Dictionary(tag language.Tag) (Dictionary, error) ***REMOVED***
// 	// TODO: verify dictionary exists.
// 	return &dict***REMOVED***&c.index, tag***REMOVED***, nil
// ***REMOVED***

type dict struct ***REMOVED***
	s   *store
	tag language.Tag // TODO: make compact tag.
***REMOVED***

func (d *dict) Lookup(key string) (data string, ok bool) ***REMOVED***
	return d.s.lookup(d.tag, key)
***REMOVED***

func (b *Builder) lookup(tag language.Tag, key string) (data string, ok bool) ***REMOVED***
	return b.index.lookup(tag, key)
***REMOVED***

func (c *Builder) set(tag language.Tag, key string, s *store, msg ...Message) error ***REMOVED***
	data, err := catmsg.Compile(tag, &dict***REMOVED***&c.macros, tag***REMOVED***, firstInSequence(msg))

	s.mutex.Lock()
	defer s.mutex.Unlock()

	m := s.index[tag]
	if m == nil ***REMOVED***
		m = msgMap***REMOVED******REMOVED***
		if s.index == nil ***REMOVED***
			s.index = map[language.Tag]msgMap***REMOVED******REMOVED***
		***REMOVED***
		c.matcher = nil
		s.index[tag] = m
	***REMOVED***

	m[key] = data
	return err
***REMOVED***

func (c *Builder) Matcher() language.Matcher ***REMOVED***
	c.index.mutex.RLock()
	m := c.matcher
	c.index.mutex.RUnlock()
	if m != nil ***REMOVED***
		return m
	***REMOVED***

	c.index.mutex.Lock()
	if c.matcher == nil ***REMOVED***
		c.matcher = language.NewMatcher(c.unlockedLanguages())
	***REMOVED***
	m = c.matcher
	c.index.mutex.Unlock()
	return m
***REMOVED***

type store struct ***REMOVED***
	mutex sync.RWMutex
	index map[language.Tag]msgMap
***REMOVED***

type msgMap map[string]string

func (s *store) lookup(tag language.Tag, key string) (data string, ok bool) ***REMOVED***
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for ; ; tag = tag.Parent() ***REMOVED***
		if msgs, ok := s.index[tag]; ok ***REMOVED***
			if msg, ok := msgs[key]; ok ***REMOVED***
				return msg, true
			***REMOVED***
		***REMOVED***
		if tag == language.Und ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return "", false
***REMOVED***

// Languages returns all languages for which the Catalog contains variants.
func (b *Builder) Languages() []language.Tag ***REMOVED***
	s := &b.index
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return b.unlockedLanguages()
***REMOVED***

func (b *Builder) unlockedLanguages() []language.Tag ***REMOVED***
	s := &b.index
	if len(s.index) == 0 ***REMOVED***
		return nil
	***REMOVED***
	tags := make([]language.Tag, 0, len(s.index))
	_, hasFallback := s.index[b.options.fallback]
	offset := 0
	if hasFallback ***REMOVED***
		tags = append(tags, b.options.fallback)
		offset = 1
	***REMOVED***
	for t := range s.index ***REMOVED***
		if t != b.options.fallback ***REMOVED***
			tags = append(tags, t)
		***REMOVED***
	***REMOVED***
	internal.SortTags(tags[offset:])
	return tags
***REMOVED***
