package main

import (
	"math/rand"
)

// chunkStrings chunks the string slice
func chunkStrings(x []string, numChunks int) [][]string ***REMOVED***
	var result [][]string
	chunkSize := (len(x) + numChunks - 1) / numChunks
	for i := 0; i < len(x); i += chunkSize ***REMOVED***
		ub := i + chunkSize
		if ub > len(x) ***REMOVED***
			ub = len(x)
		***REMOVED***
		result = append(result, x[i:ub])
	***REMOVED***
	return result
***REMOVED***

// shuffleStrings shuffles strings
func shuffleStrings(x []string, seed int64) ***REMOVED***
	r := rand.New(rand.NewSource(seed))
	for i := range x ***REMOVED***
		j := r.Intn(i + 1)
		x[i], x[j] = x[j], x[i]
	***REMOVED***
***REMOVED***
