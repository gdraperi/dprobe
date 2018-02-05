// +build gofuzz

package toml

func Fuzz(data []byte) int ***REMOVED***
	tree, err := LoadBytes(data)
	if err != nil ***REMOVED***
		if tree != nil ***REMOVED***
			panic("tree must be nil if there is an error")
		***REMOVED***
		return 0
	***REMOVED***

	str, err := tree.ToTomlString()
	if err != nil ***REMOVED***
		if str != "" ***REMOVED***
			panic(`str must be "" if there is an error`)
		***REMOVED***
		panic(err)
	***REMOVED***

	tree, err = Load(str)
	if err != nil ***REMOVED***
		if tree != nil ***REMOVED***
			panic("tree must be nil if there is an error")
		***REMOVED***
		return 0
	***REMOVED***

	return 1
***REMOVED***
