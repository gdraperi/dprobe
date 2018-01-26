package discovery

import (
	"fmt"
	"regexp"
	"strconv"
)

// Generate takes care of IP generation
func Generate(pattern string) []string ***REMOVED***
	re, _ := regexp.Compile(`\[(.+):(.+)\]`)
	submatch := re.FindStringSubmatch(pattern)
	if submatch == nil ***REMOVED***
		return []string***REMOVED***pattern***REMOVED***
	***REMOVED***

	from, err := strconv.Atoi(submatch[1])
	if err != nil ***REMOVED***
		return []string***REMOVED***pattern***REMOVED***
	***REMOVED***
	to, err := strconv.Atoi(submatch[2])
	if err != nil ***REMOVED***
		return []string***REMOVED***pattern***REMOVED***
	***REMOVED***

	template := re.ReplaceAllString(pattern, "%d")

	var result []string
	for val := from; val <= to; val++ ***REMOVED***
		entry := fmt.Sprintf(template, val)
		result = append(result, entry)
	***REMOVED***

	return result
***REMOVED***
