package shellwords

import (
	"errors"
	"os"
	"regexp"
)

var (
	ParseEnv      bool = false
	ParseBacktick bool = false
)

var envRe = regexp.MustCompile(`\$(***REMOVED***[a-zA-Z0-9_]+***REMOVED***|[a-zA-Z0-9_]+)`)

func isSpace(r rune) bool ***REMOVED***
	switch r ***REMOVED***
	case ' ', '\t', '\r', '\n':
		return true
	***REMOVED***
	return false
***REMOVED***

func replaceEnv(s string) string ***REMOVED***
	return envRe.ReplaceAllStringFunc(s, func(s string) string ***REMOVED***
		s = s[1:]
		if s[0] == '***REMOVED***' ***REMOVED***
			s = s[1 : len(s)-1]
		***REMOVED***
		return os.Getenv(s)
	***REMOVED***)
***REMOVED***

type Parser struct ***REMOVED***
	ParseEnv      bool
	ParseBacktick bool
	Position      int
***REMOVED***

func NewParser() *Parser ***REMOVED***
	return &Parser***REMOVED***ParseEnv, ParseBacktick, 0***REMOVED***
***REMOVED***

func (p *Parser) Parse(line string) ([]string, error) ***REMOVED***
	args := []string***REMOVED******REMOVED***
	buf := ""
	var escaped, doubleQuoted, singleQuoted, backQuote bool
	backtick := ""

	pos := -1
	got := false

loop:
	for i, r := range line ***REMOVED***
		if escaped ***REMOVED***
			buf += string(r)
			escaped = false
			continue
		***REMOVED***

		if r == '\\' ***REMOVED***
			if singleQuoted ***REMOVED***
				buf += string(r)
			***REMOVED*** else ***REMOVED***
				escaped = true
			***REMOVED***
			continue
		***REMOVED***

		if isSpace(r) ***REMOVED***
			if singleQuoted || doubleQuoted || backQuote ***REMOVED***
				buf += string(r)
				backtick += string(r)
			***REMOVED*** else if got ***REMOVED***
				if p.ParseEnv ***REMOVED***
					buf = replaceEnv(buf)
				***REMOVED***
				args = append(args, buf)
				buf = ""
				got = false
			***REMOVED***
			continue
		***REMOVED***

		switch r ***REMOVED***
		case '`':
			if !singleQuoted && !doubleQuoted ***REMOVED***
				if p.ParseBacktick ***REMOVED***
					if backQuote ***REMOVED***
						out, err := shellRun(backtick)
						if err != nil ***REMOVED***
							return nil, err
						***REMOVED***
						buf = out
					***REMOVED***
					backtick = ""
					backQuote = !backQuote
					continue
				***REMOVED***
				backtick = ""
				backQuote = !backQuote
			***REMOVED***
		case '"':
			if !singleQuoted ***REMOVED***
				doubleQuoted = !doubleQuoted
				continue
			***REMOVED***
		case '\'':
			if !doubleQuoted ***REMOVED***
				singleQuoted = !singleQuoted
				continue
			***REMOVED***
		case ';', '&', '|', '<', '>':
			if !(escaped || singleQuoted || doubleQuoted || backQuote) ***REMOVED***
				pos = i
				break loop
			***REMOVED***
		***REMOVED***

		got = true
		buf += string(r)
		if backQuote ***REMOVED***
			backtick += string(r)
		***REMOVED***
	***REMOVED***

	if got ***REMOVED***
		if p.ParseEnv ***REMOVED***
			buf = replaceEnv(buf)
		***REMOVED***
		args = append(args, buf)
	***REMOVED***

	if escaped || singleQuoted || doubleQuoted || backQuote ***REMOVED***
		return nil, errors.New("invalid command line string")
	***REMOVED***

	p.Position = pos

	return args, nil
***REMOVED***

func Parse(line string) ([]string, error) ***REMOVED***
	return NewParser().Parse(line)
***REMOVED***
