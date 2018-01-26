// Copyright 2014 Unknwon
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
	"errors"
	"fmt"
	"strings"
)

// Section represents a config section.
type Section struct ***REMOVED***
	f        *File
	Comment  string
	name     string
	keys     map[string]*Key
	keyList  []string
	keysHash map[string]string

	isRawSection bool
	rawBody      string
***REMOVED***

func newSection(f *File, name string) *Section ***REMOVED***
	return &Section***REMOVED***
		f:        f,
		name:     name,
		keys:     make(map[string]*Key),
		keyList:  make([]string, 0, 10),
		keysHash: make(map[string]string),
	***REMOVED***
***REMOVED***

// Name returns name of Section.
func (s *Section) Name() string ***REMOVED***
	return s.name
***REMOVED***

// Body returns rawBody of Section if the section was marked as unparseable.
// It still follows the other rules of the INI format surrounding leading/trailing whitespace.
func (s *Section) Body() string ***REMOVED***
	return strings.TrimSpace(s.rawBody)
***REMOVED***

// NewKey creates a new key to given section.
func (s *Section) NewKey(name, val string) (*Key, error) ***REMOVED***
	if len(name) == 0 ***REMOVED***
		return nil, errors.New("error creating new key: empty key name")
	***REMOVED*** else if s.f.options.Insensitive ***REMOVED***
		name = strings.ToLower(name)
	***REMOVED***

	if s.f.BlockMode ***REMOVED***
		s.f.lock.Lock()
		defer s.f.lock.Unlock()
	***REMOVED***

	if inSlice(name, s.keyList) ***REMOVED***
		if s.f.options.AllowShadows ***REMOVED***
			if err := s.keys[name].addShadow(val); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			s.keys[name].value = val
		***REMOVED***
		return s.keys[name], nil
	***REMOVED***

	s.keyList = append(s.keyList, name)
	s.keys[name] = newKey(s, name, val)
	s.keysHash[name] = val
	return s.keys[name], nil
***REMOVED***

// NewBooleanKey creates a new boolean type key to given section.
func (s *Section) NewBooleanKey(name string) (*Key, error) ***REMOVED***
	key, err := s.NewKey(name, "true")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	key.isBooleanType = true
	return key, nil
***REMOVED***

// GetKey returns key in section by given name.
func (s *Section) GetKey(name string) (*Key, error) ***REMOVED***
	// FIXME: change to section level lock?
	if s.f.BlockMode ***REMOVED***
		s.f.lock.RLock()
	***REMOVED***
	if s.f.options.Insensitive ***REMOVED***
		name = strings.ToLower(name)
	***REMOVED***
	key := s.keys[name]
	if s.f.BlockMode ***REMOVED***
		s.f.lock.RUnlock()
	***REMOVED***

	if key == nil ***REMOVED***
		// Check if it is a child-section.
		sname := s.name
		for ***REMOVED***
			if i := strings.LastIndex(sname, "."); i > -1 ***REMOVED***
				sname = sname[:i]
				sec, err := s.f.GetSection(sname)
				if err != nil ***REMOVED***
					continue
				***REMOVED***
				return sec.GetKey(name)
			***REMOVED*** else ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		return nil, fmt.Errorf("error when getting key of section '%s': key '%s' not exists", s.name, name)
	***REMOVED***
	return key, nil
***REMOVED***

// HasKey returns true if section contains a key with given name.
func (s *Section) HasKey(name string) bool ***REMOVED***
	key, _ := s.GetKey(name)
	return key != nil
***REMOVED***

// Haskey is a backwards-compatible name for HasKey.
func (s *Section) Haskey(name string) bool ***REMOVED***
	return s.HasKey(name)
***REMOVED***

// HasValue returns true if section contains given raw value.
func (s *Section) HasValue(value string) bool ***REMOVED***
	if s.f.BlockMode ***REMOVED***
		s.f.lock.RLock()
		defer s.f.lock.RUnlock()
	***REMOVED***

	for _, k := range s.keys ***REMOVED***
		if value == k.value ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// Key assumes named Key exists in section and returns a zero-value when not.
func (s *Section) Key(name string) *Key ***REMOVED***
	key, err := s.GetKey(name)
	if err != nil ***REMOVED***
		// It's OK here because the only possible error is empty key name,
		// but if it's empty, this piece of code won't be executed.
		key, _ = s.NewKey(name, "")
		return key
	***REMOVED***
	return key
***REMOVED***

// Keys returns list of keys of section.
func (s *Section) Keys() []*Key ***REMOVED***
	keys := make([]*Key, len(s.keyList))
	for i := range s.keyList ***REMOVED***
		keys[i] = s.Key(s.keyList[i])
	***REMOVED***
	return keys
***REMOVED***

// ParentKeys returns list of keys of parent section.
func (s *Section) ParentKeys() []*Key ***REMOVED***
	var parentKeys []*Key
	sname := s.name
	for ***REMOVED***
		if i := strings.LastIndex(sname, "."); i > -1 ***REMOVED***
			sname = sname[:i]
			sec, err := s.f.GetSection(sname)
			if err != nil ***REMOVED***
				continue
			***REMOVED***
			parentKeys = append(parentKeys, sec.Keys()...)
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***

	***REMOVED***
	return parentKeys
***REMOVED***

// KeyStrings returns list of key names of section.
func (s *Section) KeyStrings() []string ***REMOVED***
	list := make([]string, len(s.keyList))
	copy(list, s.keyList)
	return list
***REMOVED***

// KeysHash returns keys hash consisting of names and values.
func (s *Section) KeysHash() map[string]string ***REMOVED***
	if s.f.BlockMode ***REMOVED***
		s.f.lock.RLock()
		defer s.f.lock.RUnlock()
	***REMOVED***

	hash := map[string]string***REMOVED******REMOVED***
	for key, value := range s.keysHash ***REMOVED***
		hash[key] = value
	***REMOVED***
	return hash
***REMOVED***

// DeleteKey deletes a key from section.
func (s *Section) DeleteKey(name string) ***REMOVED***
	if s.f.BlockMode ***REMOVED***
		s.f.lock.Lock()
		defer s.f.lock.Unlock()
	***REMOVED***

	for i, k := range s.keyList ***REMOVED***
		if k == name ***REMOVED***
			s.keyList = append(s.keyList[:i], s.keyList[i+1:]...)
			delete(s.keys, name)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***
