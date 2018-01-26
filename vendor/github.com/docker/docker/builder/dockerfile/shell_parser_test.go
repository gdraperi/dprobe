package dockerfile

import (
	"bufio"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShellParser4EnvVars(t *testing.T) ***REMOVED***
	fn := "envVarTest"
	lineCount := 0

	file, err := os.Open(fn)
	assert.NoError(t, err)
	defer file.Close()

	shlex := NewShellLex('\\')
	scanner := bufio.NewScanner(file)
	envs := []string***REMOVED***"PWD=/home", "SHELL=bash", "KOREAN=한국어"***REMOVED***
	for scanner.Scan() ***REMOVED***
		line := scanner.Text()
		lineCount++

		// Trim comments and blank lines
		i := strings.Index(line, "#")
		if i >= 0 ***REMOVED***
			line = line[:i]
		***REMOVED***
		line = strings.TrimSpace(line)

		if line == "" ***REMOVED***
			continue
		***REMOVED***

		words := strings.Split(line, "|")
		assert.Len(t, words, 3)

		platform := strings.TrimSpace(words[0])
		source := strings.TrimSpace(words[1])
		expected := strings.TrimSpace(words[2])

		// Key W=Windows; A=All; U=Unix
		if platform != "W" && platform != "A" && platform != "U" ***REMOVED***
			t.Fatalf("Invalid tag %s at line %d of %s. Must be W, A or U", platform, lineCount, fn)
		***REMOVED***

		if ((platform == "W" || platform == "A") && runtime.GOOS == "windows") ||
			((platform == "U" || platform == "A") && runtime.GOOS != "windows") ***REMOVED***
			newWord, err := shlex.ProcessWord(source, envs)
			if expected == "error" ***REMOVED***
				assert.Error(t, err)
			***REMOVED*** else ***REMOVED***
				assert.NoError(t, err)
				assert.Equal(t, newWord, expected)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestShellParser4Words(t *testing.T) ***REMOVED***
	fn := "wordsTest"

	file, err := os.Open(fn)
	if err != nil ***REMOVED***
		t.Fatalf("Can't open '%s': %s", err, fn)
	***REMOVED***
	defer file.Close()

	shlex := NewShellLex('\\')
	envs := []string***REMOVED******REMOVED***
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() ***REMOVED***
		line := scanner.Text()
		lineNum = lineNum + 1

		if strings.HasPrefix(line, "#") ***REMOVED***
			continue
		***REMOVED***

		if strings.HasPrefix(line, "ENV ") ***REMOVED***
			line = strings.TrimLeft(line[3:], " ")
			envs = append(envs, line)
			continue
		***REMOVED***

		words := strings.Split(line, "|")
		if len(words) != 2 ***REMOVED***
			t.Fatalf("Error in '%s'(line %d) - should be exactly one | in: %q", fn, lineNum, line)
		***REMOVED***
		test := strings.TrimSpace(words[0])
		expected := strings.Split(strings.TrimLeft(words[1], " "), ",")

		result, err := shlex.ProcessWords(test, envs)

		if err != nil ***REMOVED***
			result = []string***REMOVED***"error"***REMOVED***
		***REMOVED***

		if len(result) != len(expected) ***REMOVED***
			t.Fatalf("Error on line %d. %q was suppose to result in %q, but got %q instead", lineNum, test, expected, result)
		***REMOVED***
		for i, w := range expected ***REMOVED***
			if w != result[i] ***REMOVED***
				t.Fatalf("Error on line %d. %q was suppose to result in %q, but got %q instead", lineNum, test, expected, result)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestGetEnv(t *testing.T) ***REMOVED***
	sw := &shellWord***REMOVED***envs: nil***REMOVED***

	sw.envs = []string***REMOVED******REMOVED***
	if sw.getEnv("foo") != "" ***REMOVED***
		t.Fatal("2 - 'foo' should map to ''")
	***REMOVED***

	sw.envs = []string***REMOVED***"foo"***REMOVED***
	if sw.getEnv("foo") != "" ***REMOVED***
		t.Fatal("3 - 'foo' should map to ''")
	***REMOVED***

	sw.envs = []string***REMOVED***"foo="***REMOVED***
	if sw.getEnv("foo") != "" ***REMOVED***
		t.Fatal("4 - 'foo' should map to ''")
	***REMOVED***

	sw.envs = []string***REMOVED***"foo=bar"***REMOVED***
	if sw.getEnv("foo") != "bar" ***REMOVED***
		t.Fatal("5 - 'foo' should map to 'bar'")
	***REMOVED***

	sw.envs = []string***REMOVED***"foo=bar", "car=hat"***REMOVED***
	if sw.getEnv("foo") != "bar" ***REMOVED***
		t.Fatal("6 - 'foo' should map to 'bar'")
	***REMOVED***
	if sw.getEnv("car") != "hat" ***REMOVED***
		t.Fatal("7 - 'car' should map to 'hat'")
	***REMOVED***

	// Make sure we grab the first 'car' in the list
	sw.envs = []string***REMOVED***"foo=bar", "car=hat", "car=bike"***REMOVED***
	if sw.getEnv("car") != "hat" ***REMOVED***
		t.Fatal("8 - 'car' should map to 'hat'")
	***REMOVED***
***REMOVED***
