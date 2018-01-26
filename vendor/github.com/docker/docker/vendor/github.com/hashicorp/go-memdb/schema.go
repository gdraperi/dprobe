package memdb

import "fmt"

// DBSchema contains the full database schema used for MemDB
type DBSchema struct ***REMOVED***
	Tables map[string]*TableSchema
***REMOVED***

// Validate is used to validate the database schema
func (s *DBSchema) Validate() error ***REMOVED***
	if s == nil ***REMOVED***
		return fmt.Errorf("missing schema")
	***REMOVED***
	if len(s.Tables) == 0 ***REMOVED***
		return fmt.Errorf("no tables defined")
	***REMOVED***
	for name, table := range s.Tables ***REMOVED***
		if name != table.Name ***REMOVED***
			return fmt.Errorf("table name mis-match for '%s'", name)
		***REMOVED***
		if err := table.Validate(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// TableSchema contains the schema for a single table
type TableSchema struct ***REMOVED***
	Name    string
	Indexes map[string]*IndexSchema
***REMOVED***

// Validate is used to validate the table schema
func (s *TableSchema) Validate() error ***REMOVED***
	if s.Name == "" ***REMOVED***
		return fmt.Errorf("missing table name")
	***REMOVED***
	if len(s.Indexes) == 0 ***REMOVED***
		return fmt.Errorf("missing table schemas for '%s'", s.Name)
	***REMOVED***
	if _, ok := s.Indexes["id"]; !ok ***REMOVED***
		return fmt.Errorf("must have id index")
	***REMOVED***
	if !s.Indexes["id"].Unique ***REMOVED***
		return fmt.Errorf("id index must be unique")
	***REMOVED***
	if _, ok := s.Indexes["id"].Indexer.(SingleIndexer); !ok ***REMOVED***
		return fmt.Errorf("id index must be a SingleIndexer")
	***REMOVED***
	for name, index := range s.Indexes ***REMOVED***
		if name != index.Name ***REMOVED***
			return fmt.Errorf("index name mis-match for '%s'", name)
		***REMOVED***
		if err := index.Validate(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// IndexSchema contains the schema for an index
type IndexSchema struct ***REMOVED***
	Name         string
	AllowMissing bool
	Unique       bool
	Indexer      Indexer
***REMOVED***

func (s *IndexSchema) Validate() error ***REMOVED***
	if s.Name == "" ***REMOVED***
		return fmt.Errorf("missing index name")
	***REMOVED***
	if s.Indexer == nil ***REMOVED***
		return fmt.Errorf("missing index function for '%s'", s.Name)
	***REMOVED***
	switch s.Indexer.(type) ***REMOVED***
	case SingleIndexer:
	case MultiIndexer:
	default:
		return fmt.Errorf("indexer for '%s' must be a SingleIndexer or MultiIndexer", s.Name)
	***REMOVED***
	return nil
***REMOVED***
