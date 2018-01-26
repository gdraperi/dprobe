package dockerfile

import (
	"bytes"
	"strings"
	"text/scanner"
	"unicode"

	"github.com/pkg/errors"
)

// ShellLex performs shell word splitting and variable expansion.
//
// ShellLex takes a string and an array of env variables and
// process all quotes (" and ') as well as $xxx and $***REMOVED***xxx***REMOVED*** env variable
// tokens.  Tries to mimic bash shell process.
// It doesn't support all flavors of $***REMOVED***xx:...***REMOVED*** formats but new ones can
// be added by adding code to the "special $***REMOVED******REMOVED*** format processing" section
type ShellLex struct ***REMOVED***
	escapeToken rune
***REMOVED***

// NewShellLex creates a new ShellLex which uses escapeToken to escape quotes.
func NewShellLex(escapeToken rune) *ShellLex ***REMOVED***
	return &ShellLex***REMOVED***escapeToken: escapeToken***REMOVED***
***REMOVED***

// ProcessWord will use the 'env' list of environment variables,
// and replace any env var references in 'word'.
func (s *ShellLex) ProcessWord(word string, env []string) (string, error) ***REMOVED***
	word, _, err := s.process(word, env)
	return word, err
***REMOVED***

// ProcessWords will use the 'env' list of environment variables,
// and replace any env var references in 'word' then it will also
// return a slice of strings which represents the 'word'
// split up based on spaces - taking into account quotes.  Note that
// this splitting is done **after** the env var substitutions are done.
// Note, each one is trimmed to remove leading and trailing spaces (unless
// they are quoted", but ProcessWord retains spaces between words.
func (s *ShellLex) ProcessWords(word string, env []string) ([]string, error) ***REMOVED***
	_, words, err := s.process(word, env)
	return words, err
***REMOVED***

func (s *ShellLex) process(word string, env []string) (string, []string, error) ***REMOVED***
	sw := &shellWord***REMOVED***
		envs:        env,
		escapeToken: s.escapeToken,
	***REMOVED***
	sw.scanner.Init(strings.NewReader(word))
	return sw.process(word)
***REMOVED***

type shellWord struct ***REMOVED***
	scanner     scanner.Scanner
	envs        []string
	escapeToken rune
***REMOVED***

func (sw *shellWord) process(source string) (string, []string, error) ***REMOVED***
	word, words, err := sw.processStopOn(scanner.EOF)
	if err != nil ***REMOVED***
		err = errors.Wrapf(err, "failed to process %q", source)
	***REMOVED***
	return word, words, err
***REMOVED***

type wordsStruct struct ***REMOVED***
	word   string
	words  []string
	inWord bool
***REMOVED***

func (w *wordsStruct) addChar(ch rune) ***REMOVED***
	if unicode.IsSpace(ch) && w.inWord ***REMOVED***
		if len(w.word) != 0 ***REMOVED***
			w.words = append(w.words, w.word)
			w.word = ""
			w.inWord = false
		***REMOVED***
	***REMOVED*** else if !unicode.IsSpace(ch) ***REMOVED***
		w.addRawChar(ch)
	***REMOVED***
***REMOVED***

func (w *wordsStruct) addRawChar(ch rune) ***REMOVED***
	w.word += string(ch)
	w.inWord = true
***REMOVED***

func (w *wordsStruct) addString(str string) ***REMOVED***
	var scan scanner.Scanner
	scan.Init(strings.NewReader(str))
	for scan.Peek() != scanner.EOF ***REMOVED***
		w.addChar(scan.Next())
	***REMOVED***
***REMOVED***

func (w *wordsStruct) addRawString(str string) ***REMOVED***
	w.word += str
	w.inWord = true
***REMOVED***

func (w *wordsStruct) getWords() []string ***REMOVED***
	if len(w.word) > 0 ***REMOVED***
		w.words = append(w.words, w.word)

		// Just in case we're called again by mistake
		w.word = ""
		w.inWord = false
	***REMOVED***
	return w.words
***REMOVED***

// Process the word, starting at 'pos', and stop when we get to the
// end of the word or the 'stopChar' character
func (sw *shellWord) processStopOn(stopChar rune) (string, []string, error) ***REMOVED***
	var result bytes.Buffer
	var words wordsStruct

	var charFuncMapping = map[rune]func() (string, error)***REMOVED***
		'\'': sw.processSingleQuote,
		'"':  sw.processDoubleQuote,
		'$':  sw.processDollar,
	***REMOVED***

	for sw.scanner.Peek() != scanner.EOF ***REMOVED***
		ch := sw.scanner.Peek()

		if stopChar != scanner.EOF && ch == stopChar ***REMOVED***
			sw.scanner.Next()
			break
		***REMOVED***
		if fn, ok := charFuncMapping[ch]; ok ***REMOVED***
			// Call special processing func for certain chars
			tmp, err := fn()
			if err != nil ***REMOVED***
				return "", []string***REMOVED******REMOVED***, err
			***REMOVED***
			result.WriteString(tmp)

			if ch == rune('$') ***REMOVED***
				words.addString(tmp)
			***REMOVED*** else ***REMOVED***
				words.addRawString(tmp)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// Not special, just add it to the result
			ch = sw.scanner.Next()

			if ch == sw.escapeToken ***REMOVED***
				// '\' (default escape token, but ` allowed) escapes, except end of line
				ch = sw.scanner.Next()

				if ch == scanner.EOF ***REMOVED***
					break
				***REMOVED***

				words.addRawChar(ch)
			***REMOVED*** else ***REMOVED***
				words.addChar(ch)
			***REMOVED***

			result.WriteRune(ch)
		***REMOVED***
	***REMOVED***

	return result.String(), words.getWords(), nil
***REMOVED***

func (sw *shellWord) processSingleQuote() (string, error) ***REMOVED***
	// All chars between single quotes are taken as-is
	// Note, you can't escape '
	//
	// From the "sh" man page:
	// Single Quotes
	//   Enclosing characters in single quotes preserves the literal meaning of
	//   all the characters (except single quotes, making it impossible to put
	//   single-quotes in a single-quoted string).

	var result bytes.Buffer

	sw.scanner.Next()

	for ***REMOVED***
		ch := sw.scanner.Next()
		switch ch ***REMOVED***
		case scanner.EOF:
			return "", errors.New("unexpected end of statement while looking for matching single-quote")
		case '\'':
			return result.String(), nil
		***REMOVED***
		result.WriteRune(ch)
	***REMOVED***
***REMOVED***

func (sw *shellWord) processDoubleQuote() (string, error) ***REMOVED***
	// All chars up to the next " are taken as-is, even ', except any $ chars
	// But you can escape " with a \ (or ` if escape token set accordingly)
	//
	// From the "sh" man page:
	// Double Quotes
	//  Enclosing characters within double quotes preserves the literal meaning
	//  of all characters except dollarsign ($), backquote (`), and backslash
	//  (\).  The backslash inside double quotes is historically weird, and
	//  serves to quote only the following characters:
	//    $ ` " \ <newline>.
	//  Otherwise it remains literal.

	var result bytes.Buffer

	sw.scanner.Next()

	for ***REMOVED***
		switch sw.scanner.Peek() ***REMOVED***
		case scanner.EOF:
			return "", errors.New("unexpected end of statement while looking for matching double-quote")
		case '"':
			sw.scanner.Next()
			return result.String(), nil
		case '$':
			value, err := sw.processDollar()
			if err != nil ***REMOVED***
				return "", err
			***REMOVED***
			result.WriteString(value)
		default:
			ch := sw.scanner.Next()
			if ch == sw.escapeToken ***REMOVED***
				switch sw.scanner.Peek() ***REMOVED***
				case scanner.EOF:
					// Ignore \ at end of word
					continue
				case '"', '$', sw.escapeToken:
					// These chars can be escaped, all other \'s are left as-is
					// Note: for now don't do anything special with ` chars.
					// Not sure what to do with them anyway since we're not going
					// to execute the text in there (not now anyway).
					ch = sw.scanner.Next()
				***REMOVED***
			***REMOVED***
			result.WriteRune(ch)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (sw *shellWord) processDollar() (string, error) ***REMOVED***
	sw.scanner.Next()

	// $xxx case
	if sw.scanner.Peek() != '***REMOVED***' ***REMOVED***
		name := sw.processName()
		if name == "" ***REMOVED***
			return "$", nil
		***REMOVED***
		return sw.getEnv(name), nil
	***REMOVED***

	sw.scanner.Next()
	name := sw.processName()
	ch := sw.scanner.Peek()
	if ch == '***REMOVED***' ***REMOVED***
		// Normal $***REMOVED***xx***REMOVED*** case
		sw.scanner.Next()
		return sw.getEnv(name), nil
	***REMOVED***
	if ch == ':' ***REMOVED***
		// Special $***REMOVED***xx:...***REMOVED*** format processing
		// Yes it allows for recursive $'s in the ... spot

		sw.scanner.Next() // skip over :
		modifier := sw.scanner.Next()

		word, _, err := sw.processStopOn('***REMOVED***')
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***

		// Grab the current value of the variable in question so we
		// can use to to determine what to do based on the modifier
		newValue := sw.getEnv(name)

		switch modifier ***REMOVED***
		case '+':
			if newValue != "" ***REMOVED***
				newValue = word
			***REMOVED***
			return newValue, nil

		case '-':
			if newValue == "" ***REMOVED***
				newValue = word
			***REMOVED***
			return newValue, nil

		default:
			return "", errors.Errorf("unsupported modifier (%c) in substitution", modifier)
		***REMOVED***
	***REMOVED***
	return "", errors.Errorf("missing ':' in substitution")
***REMOVED***

func (sw *shellWord) processName() string ***REMOVED***
	// Read in a name (alphanumeric or _)
	// If it starts with a numeric then just return $#
	var name bytes.Buffer

	for sw.scanner.Peek() != scanner.EOF ***REMOVED***
		ch := sw.scanner.Peek()
		if name.Len() == 0 && unicode.IsDigit(ch) ***REMOVED***
			ch = sw.scanner.Next()
			return string(ch)
		***REMOVED***
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' ***REMOVED***
			break
		***REMOVED***
		ch = sw.scanner.Next()
		name.WriteRune(ch)
	***REMOVED***

	return name.String()
***REMOVED***

func (sw *shellWord) getEnv(name string) string ***REMOVED***
	for _, env := range sw.envs ***REMOVED***
		i := strings.Index(env, "=")
		if i < 0 ***REMOVED***
			if equalEnvKeys(name, env) ***REMOVED***
				// Should probably never get here, but just in case treat
				// it like "var" and "var=" are the same
				return ""
			***REMOVED***
			continue
		***REMOVED***
		compareName := env[:i]
		if !equalEnvKeys(name, compareName) ***REMOVED***
			continue
		***REMOVED***
		return env[i+1:]
	***REMOVED***
	return ""
***REMOVED***
