package main

//-------------------------------------------------------------------------
// scope
//-------------------------------------------------------------------------

type scope struct ***REMOVED***
	// the package name that this scope resides in
	pkgname  string
	parent   *scope // nil for universe scope
	entities map[string]*decl
***REMOVED***

func new_named_scope(outer *scope, name string) *scope ***REMOVED***
	s := new_scope(outer)
	s.pkgname = name
	return s
***REMOVED***

func new_scope(outer *scope) *scope ***REMOVED***
	s := new(scope)
	if outer != nil ***REMOVED***
		s.pkgname = outer.pkgname
	***REMOVED***
	s.parent = outer
	s.entities = make(map[string]*decl)
	return s
***REMOVED***

// returns: new, prev
func advance_scope(s *scope) (*scope, *scope) ***REMOVED***
	if len(s.entities) == 0 ***REMOVED***
		return s, s.parent
	***REMOVED***
	return new_scope(s), s
***REMOVED***

// adds declaration or returns an existing one
func (s *scope) add_named_decl(d *decl) *decl ***REMOVED***
	return s.add_decl(d.name, d)
***REMOVED***

func (s *scope) add_decl(name string, d *decl) *decl ***REMOVED***
	decl, ok := s.entities[name]
	if !ok ***REMOVED***
		s.entities[name] = d
		return d
	***REMOVED***
	return decl
***REMOVED***

func (s *scope) replace_decl(name string, d *decl) ***REMOVED***
	s.entities[name] = d
***REMOVED***

func (s *scope) merge_decl(d *decl) ***REMOVED***
	decl, ok := s.entities[d.name]
	if !ok ***REMOVED***
		s.entities[d.name] = d
	***REMOVED*** else ***REMOVED***
		decl := decl.deep_copy()
		decl.expand_or_replace(d)
		s.entities[d.name] = decl
	***REMOVED***
***REMOVED***

func (s *scope) lookup(name string) *decl ***REMOVED***
	decl, ok := s.entities[name]
	if !ok ***REMOVED***
		if s.parent != nil ***REMOVED***
			return s.parent.lookup(name)
		***REMOVED*** else ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	return decl
***REMOVED***
