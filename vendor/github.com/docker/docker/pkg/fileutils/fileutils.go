package fileutils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/scanner"

	"github.com/sirupsen/logrus"
)

// PatternMatcher allows checking paths agaist a list of patterns
type PatternMatcher struct ***REMOVED***
	patterns   []*Pattern
	exclusions bool
***REMOVED***

// NewPatternMatcher creates a new matcher object for specific patterns that can
// be used later to match against patterns against paths
func NewPatternMatcher(patterns []string) (*PatternMatcher, error) ***REMOVED***
	pm := &PatternMatcher***REMOVED***
		patterns: make([]*Pattern, 0, len(patterns)),
	***REMOVED***
	for _, p := range patterns ***REMOVED***
		// Eliminate leading and trailing whitespace.
		p = strings.TrimSpace(p)
		if p == "" ***REMOVED***
			continue
		***REMOVED***
		p = filepath.Clean(p)
		newp := &Pattern***REMOVED******REMOVED***
		if p[0] == '!' ***REMOVED***
			if len(p) == 1 ***REMOVED***
				return nil, errors.New("illegal exclusion pattern: \"!\"")
			***REMOVED***
			newp.exclusion = true
			p = p[1:]
			pm.exclusions = true
		***REMOVED***
		// Do some syntax checking on the pattern.
		// filepath's Match() has some really weird rules that are inconsistent
		// so instead of trying to dup their logic, just call Match() for its
		// error state and if there is an error in the pattern return it.
		// If this becomes an issue we can remove this since its really only
		// needed in the error (syntax) case - which isn't really critical.
		if _, err := filepath.Match(p, "."); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		newp.cleanedPattern = p
		newp.dirs = strings.Split(p, string(os.PathSeparator))
		pm.patterns = append(pm.patterns, newp)
	***REMOVED***
	return pm, nil
***REMOVED***

// Matches matches path against all the patterns. Matches is not safe to be
// called concurrently
func (pm *PatternMatcher) Matches(file string) (bool, error) ***REMOVED***
	matched := false
	file = filepath.FromSlash(file)
	parentPath := filepath.Dir(file)
	parentPathDirs := strings.Split(parentPath, string(os.PathSeparator))

	for _, pattern := range pm.patterns ***REMOVED***
		negative := false

		if pattern.exclusion ***REMOVED***
			negative = true
		***REMOVED***

		match, err := pattern.match(file)
		if err != nil ***REMOVED***
			return false, err
		***REMOVED***

		if !match && parentPath != "." ***REMOVED***
			// Check to see if the pattern matches one of our parent dirs.
			if len(pattern.dirs) <= len(parentPathDirs) ***REMOVED***
				match, _ = pattern.match(strings.Join(parentPathDirs[:len(pattern.dirs)], string(os.PathSeparator)))
			***REMOVED***
		***REMOVED***

		if match ***REMOVED***
			matched = !negative
		***REMOVED***
	***REMOVED***

	if matched ***REMOVED***
		logrus.Debugf("Skipping excluded path: %s", file)
	***REMOVED***

	return matched, nil
***REMOVED***

// Exclusions returns true if any of the patterns define exclusions
func (pm *PatternMatcher) Exclusions() bool ***REMOVED***
	return pm.exclusions
***REMOVED***

// Patterns returns array of active patterns
func (pm *PatternMatcher) Patterns() []*Pattern ***REMOVED***
	return pm.patterns
***REMOVED***

// Pattern defines a single regexp used used to filter file paths.
type Pattern struct ***REMOVED***
	cleanedPattern string
	dirs           []string
	regexp         *regexp.Regexp
	exclusion      bool
***REMOVED***

func (p *Pattern) String() string ***REMOVED***
	return p.cleanedPattern
***REMOVED***

// Exclusion returns true if this pattern defines exclusion
func (p *Pattern) Exclusion() bool ***REMOVED***
	return p.exclusion
***REMOVED***

func (p *Pattern) match(path string) (bool, error) ***REMOVED***

	if p.regexp == nil ***REMOVED***
		if err := p.compile(); err != nil ***REMOVED***
			return false, filepath.ErrBadPattern
		***REMOVED***
	***REMOVED***

	b := p.regexp.MatchString(path)

	return b, nil
***REMOVED***

func (p *Pattern) compile() error ***REMOVED***
	regStr := "^"
	pattern := p.cleanedPattern
	// Go through the pattern and convert it to a regexp.
	// We use a scanner so we can support utf-8 chars.
	var scan scanner.Scanner
	scan.Init(strings.NewReader(pattern))

	sl := string(os.PathSeparator)
	escSL := sl
	if sl == `\` ***REMOVED***
		escSL += `\`
	***REMOVED***

	for scan.Peek() != scanner.EOF ***REMOVED***
		ch := scan.Next()

		if ch == '*' ***REMOVED***
			if scan.Peek() == '*' ***REMOVED***
				// is some flavor of "**"
				scan.Next()

				// Treat **/ as ** so eat the "/"
				if string(scan.Peek()) == sl ***REMOVED***
					scan.Next()
				***REMOVED***

				if scan.Peek() == scanner.EOF ***REMOVED***
					// is "**EOF" - to align with .gitignore just accept all
					regStr += ".*"
				***REMOVED*** else ***REMOVED***
					// is "**"
					// Note that this allows for any # of /'s (even 0) because
					// the .* will eat everything, even /'s
					regStr += "(.*" + escSL + ")?"
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				// is "*" so map it to anything but "/"
				regStr += "[^" + escSL + "]*"
			***REMOVED***
		***REMOVED*** else if ch == '?' ***REMOVED***
			// "?" is any char except "/"
			regStr += "[^" + escSL + "]"
		***REMOVED*** else if ch == '.' || ch == '$' ***REMOVED***
			// Escape some regexp special chars that have no meaning
			// in golang's filepath.Match
			regStr += `\` + string(ch)
		***REMOVED*** else if ch == '\\' ***REMOVED***
			// escape next char. Note that a trailing \ in the pattern
			// will be left alone (but need to escape it)
			if sl == `\` ***REMOVED***
				// On windows map "\" to "\\", meaning an escaped backslash,
				// and then just continue because filepath.Match on
				// Windows doesn't allow escaping at all
				regStr += escSL
				continue
			***REMOVED***
			if scan.Peek() != scanner.EOF ***REMOVED***
				regStr += `\` + string(scan.Next())
			***REMOVED*** else ***REMOVED***
				regStr += `\`
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			regStr += string(ch)
		***REMOVED***
	***REMOVED***

	regStr += "$"

	re, err := regexp.Compile(regStr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	p.regexp = re
	return nil
***REMOVED***

// Matches returns true if file matches any of the patterns
// and isn't excluded by any of the subsequent patterns.
func Matches(file string, patterns []string) (bool, error) ***REMOVED***
	pm, err := NewPatternMatcher(patterns)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	file = filepath.Clean(file)

	if file == "." ***REMOVED***
		// Don't let them exclude everything, kind of silly.
		return false, nil
	***REMOVED***

	return pm.Matches(file)
***REMOVED***

// CopyFile copies from src to dst until either EOF is reached
// on src or an error occurs. It verifies src exists and removes
// the dst if it exists.
func CopyFile(src, dst string) (int64, error) ***REMOVED***
	cleanSrc := filepath.Clean(src)
	cleanDst := filepath.Clean(dst)
	if cleanSrc == cleanDst ***REMOVED***
		return 0, nil
	***REMOVED***
	sf, err := os.Open(cleanSrc)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	defer sf.Close()
	if err := os.Remove(cleanDst); err != nil && !os.IsNotExist(err) ***REMOVED***
		return 0, err
	***REMOVED***
	df, err := os.Create(cleanDst)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	defer df.Close()
	return io.Copy(df, sf)
***REMOVED***

// ReadSymlinkedDirectory returns the target directory of a symlink.
// The target of the symbolic link may not be a file.
func ReadSymlinkedDirectory(path string) (string, error) ***REMOVED***
	var realPath string
	var err error
	if realPath, err = filepath.Abs(path); err != nil ***REMOVED***
		return "", fmt.Errorf("unable to get absolute path for %s: %s", path, err)
	***REMOVED***
	if realPath, err = filepath.EvalSymlinks(realPath); err != nil ***REMOVED***
		return "", fmt.Errorf("failed to canonicalise path for %s: %s", path, err)
	***REMOVED***
	realPathInfo, err := os.Stat(realPath)
	if err != nil ***REMOVED***
		return "", fmt.Errorf("failed to stat target '%s' of '%s': %s", realPath, path, err)
	***REMOVED***
	if !realPathInfo.Mode().IsDir() ***REMOVED***
		return "", fmt.Errorf("canonical path points to a file '%s'", realPath)
	***REMOVED***
	return realPath, nil
***REMOVED***

// CreateIfNotExists creates a file or a directory only if it does not already exist.
func CreateIfNotExists(path string, isDir bool) error ***REMOVED***
	if _, err := os.Stat(path); err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			if isDir ***REMOVED***
				return os.MkdirAll(path, 0755)
			***REMOVED***
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil ***REMOVED***
				return err
			***REMOVED***
			f, err := os.OpenFile(path, os.O_CREATE, 0755)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			f.Close()
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
