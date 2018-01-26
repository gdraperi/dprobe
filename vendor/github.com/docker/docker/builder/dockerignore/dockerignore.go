package dockerignore

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

// ReadAll reads a .dockerignore file and returns the list of file patterns
// to ignore. Note this will trim whitespace from each line as well
// as use GO's "clean" func to get the shortest/cleanest path for each.
func ReadAll(reader io.Reader) ([]string, error) ***REMOVED***
	if reader == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	scanner := bufio.NewScanner(reader)
	var excludes []string
	currentLine := 0

	utf8bom := []byte***REMOVED***0xEF, 0xBB, 0xBF***REMOVED***
	for scanner.Scan() ***REMOVED***
		scannedBytes := scanner.Bytes()
		// We trim UTF8 BOM
		if currentLine == 0 ***REMOVED***
			scannedBytes = bytes.TrimPrefix(scannedBytes, utf8bom)
		***REMOVED***
		pattern := string(scannedBytes)
		currentLine++
		// Lines starting with # (comments) are ignored before processing
		if strings.HasPrefix(pattern, "#") ***REMOVED***
			continue
		***REMOVED***
		pattern = strings.TrimSpace(pattern)
		if pattern == "" ***REMOVED***
			continue
		***REMOVED***
		// normalize absolute paths to paths relative to the context
		// (taking care of '!' prefix)
		invert := pattern[0] == '!'
		if invert ***REMOVED***
			pattern = strings.TrimSpace(pattern[1:])
		***REMOVED***
		if len(pattern) > 0 ***REMOVED***
			pattern = filepath.Clean(pattern)
			pattern = filepath.ToSlash(pattern)
			if len(pattern) > 1 && pattern[0] == '/' ***REMOVED***
				pattern = pattern[1:]
			***REMOVED***
		***REMOVED***
		if invert ***REMOVED***
			pattern = "!" + pattern
		***REMOVED***

		excludes = append(excludes, pattern)
	***REMOVED***
	if err := scanner.Err(); err != nil ***REMOVED***
		return nil, fmt.Errorf("Error reading .dockerignore: %v", err)
	***REMOVED***
	return excludes, nil
***REMOVED***
