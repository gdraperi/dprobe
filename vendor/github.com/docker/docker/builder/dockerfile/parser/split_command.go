package parser

import (
	"strings"
	"unicode"
)

// splitCommand takes a single line of text and parses out the cmd and args,
// which are used for dispatching to more exact parsing functions.
func splitCommand(line string) (string, []string, string, error) ***REMOVED***
	var args string
	var flags []string

	// Make sure we get the same results irrespective of leading/trailing spaces
	cmdline := tokenWhitespace.Split(strings.TrimSpace(line), 2)
	cmd := strings.ToLower(cmdline[0])

	if len(cmdline) == 2 ***REMOVED***
		var err error
		args, flags, err = extractBuilderFlags(cmdline[1])
		if err != nil ***REMOVED***
			return "", nil, "", err
		***REMOVED***
	***REMOVED***

	return cmd, flags, strings.TrimSpace(args), nil
***REMOVED***

func extractBuilderFlags(line string) (string, []string, error) ***REMOVED***
	// Parses the BuilderFlags and returns the remaining part of the line

	const (
		inSpaces = iota // looking for start of a word
		inWord
		inQuote
	)

	words := []string***REMOVED******REMOVED***
	phase := inSpaces
	word := ""
	quote := '\000'
	blankOK := false
	var ch rune

	for pos := 0; pos <= len(line); pos++ ***REMOVED***
		if pos != len(line) ***REMOVED***
			ch = rune(line[pos])
		***REMOVED***

		if phase == inSpaces ***REMOVED*** // Looking for start of word
			if pos == len(line) ***REMOVED*** // end of input
				break
			***REMOVED***
			if unicode.IsSpace(ch) ***REMOVED*** // skip spaces
				continue
			***REMOVED***

			// Only keep going if the next word starts with --
			if ch != '-' || pos+1 == len(line) || rune(line[pos+1]) != '-' ***REMOVED***
				return line[pos:], words, nil
			***REMOVED***

			phase = inWord // found something with "--", fall through
		***REMOVED***
		if (phase == inWord || phase == inQuote) && (pos == len(line)) ***REMOVED***
			if word != "--" && (blankOK || len(word) > 0) ***REMOVED***
				words = append(words, word)
			***REMOVED***
			break
		***REMOVED***
		if phase == inWord ***REMOVED***
			if unicode.IsSpace(ch) ***REMOVED***
				phase = inSpaces
				if word == "--" ***REMOVED***
					return line[pos:], words, nil
				***REMOVED***
				if blankOK || len(word) > 0 ***REMOVED***
					words = append(words, word)
				***REMOVED***
				word = ""
				blankOK = false
				continue
			***REMOVED***
			if ch == '\'' || ch == '"' ***REMOVED***
				quote = ch
				blankOK = true
				phase = inQuote
				continue
			***REMOVED***
			if ch == '\\' ***REMOVED***
				if pos+1 == len(line) ***REMOVED***
					continue // just skip \ at end
				***REMOVED***
				pos++
				ch = rune(line[pos])
			***REMOVED***
			word += string(ch)
			continue
		***REMOVED***
		if phase == inQuote ***REMOVED***
			if ch == quote ***REMOVED***
				phase = inWord
				continue
			***REMOVED***
			if ch == '\\' ***REMOVED***
				if pos+1 == len(line) ***REMOVED***
					phase = inWord
					continue // just skip \ at end
				***REMOVED***
				pos++
				ch = rune(line[pos])
			***REMOVED***
			word += string(ch)
		***REMOVED***
	***REMOVED***

	return "", words, nil
***REMOVED***
