package testutil

import "math/rand"

// GenerateRandomAlphaOnlyString generates an alphabetical random string with length n.
func GenerateRandomAlphaOnlyString(n int) string ***REMOVED***
	// make a really long string
	letters := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]byte, n)
	for i := range b ***REMOVED***
		b[i] = letters[rand.Intn(len(letters))]
	***REMOVED***
	return string(b)
***REMOVED***
