// Copyright 2015 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

import (
	"errors"
	"fmt"
)

// Scan implements sql.Scanner so UUIDs can be read from databases transparently
// Currently, database types that map to string and []byte are supported. Please
// consult database-specific driver documentation for matching types.
func (uuid *UUID) Scan(src interface***REMOVED******REMOVED***) error ***REMOVED***
	switch src.(type) ***REMOVED***
	case string:
		// if an empty UUID comes from a table, we return a null UUID
		if src.(string) == "" ***REMOVED***
			return nil
		***REMOVED***

		// see uuid.Parse for required string format
		parsed := Parse(src.(string))

		if parsed == nil ***REMOVED***
			return errors.New("Scan: invalid UUID format")
		***REMOVED***

		*uuid = parsed
	case []byte:
		b := src.([]byte)

		// if an empty UUID comes from a table, we return a null UUID
		if len(b) == 0 ***REMOVED***
			return nil
		***REMOVED***

		// assumes a simple slice of bytes if 16 bytes
		// otherwise attempts to parse
		if len(b) == 16 ***REMOVED***
			*uuid = UUID(b)
		***REMOVED*** else ***REMOVED***
			u := Parse(string(b))

			if u == nil ***REMOVED***
				return errors.New("Scan: invalid UUID format")
			***REMOVED***

			*uuid = u
		***REMOVED***

	default:
		return fmt.Errorf("Scan: unable to scan type %T into UUID", src)
	***REMOVED***

	return nil
***REMOVED***
