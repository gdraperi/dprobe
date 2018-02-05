package cmd

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

// TestGoldenAddCmd initializes the project "github.com/spf13/testproject"
// in GOPATH, adds "test" command
// and compares the content of all files in cmd directory of testproject
// with appropriate golden files.
// Use -update to update existing golden files.
func TestGoldenAddCmd(t *testing.T) ***REMOVED***
	projectName := "github.com/spf13/testproject"
	project := NewProject(projectName)
	defer os.RemoveAll(project.AbsPath())

	viper.Set("author", "NAME HERE <EMAIL ADDRESS>")
	viper.Set("license", "apache")
	viper.Set("year", 2017)
	defer viper.Set("author", nil)
	defer viper.Set("license", nil)
	defer viper.Set("year", nil)

	// Initialize the project first.
	initializeProject(project)

	// Then add the "test" command.
	cmdName := "test"
	cmdPath := filepath.Join(project.CmdPath(), cmdName+".go")
	createCmdFile(project.License(), cmdPath, cmdName)

	expectedFiles := []string***REMOVED***".", "root.go", "test.go"***REMOVED***
	gotFiles := []string***REMOVED******REMOVED***

	// Check project file hierarchy and compare the content of every single file
	// with appropriate golden file.
	err := filepath.Walk(project.CmdPath(), func(path string, info os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Make path relative to project.CmdPath().
		// E.g. path = "/home/user/go/src/github.com/spf13/testproject/cmd/root.go"
		// then it returns just "root.go".
		relPath, err := filepath.Rel(project.CmdPath(), path)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		relPath = filepath.ToSlash(relPath)
		gotFiles = append(gotFiles, relPath)
		goldenPath := filepath.Join("testdata", filepath.Base(path)+".golden")

		switch relPath ***REMOVED***
		// Known directories.
		case ".":
			return nil
		// Known files.
		case "root.go", "test.go":
			if *update ***REMOVED***
				got, err := ioutil.ReadFile(path)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				ioutil.WriteFile(goldenPath, got, 0644)
			***REMOVED***
			return compareFiles(path, goldenPath)
		***REMOVED***
		// Unknown file.
		return errors.New("unknown file: " + path)
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// Check if some files lack.
	if err := checkLackFiles(expectedFiles, gotFiles); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestValidateCmdName(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		input    string
		expected string
	***REMOVED******REMOVED***
		***REMOVED***"cmdName", "cmdName"***REMOVED***,
		***REMOVED***"cmd_name", "cmdName"***REMOVED***,
		***REMOVED***"cmd-name", "cmdName"***REMOVED***,
		***REMOVED***"cmd______Name", "cmdName"***REMOVED***,
		***REMOVED***"cmd------Name", "cmdName"***REMOVED***,
		***REMOVED***"cmd______name", "cmdName"***REMOVED***,
		***REMOVED***"cmd------name", "cmdName"***REMOVED***,
		***REMOVED***"cmdName-----", "cmdName"***REMOVED***,
		***REMOVED***"cmdname-", "cmdname"***REMOVED***,
	***REMOVED***

	for _, testCase := range testCases ***REMOVED***
		got := validateCmdName(testCase.input)
		if testCase.expected != got ***REMOVED***
			t.Errorf("Expected %q, got %q", testCase.expected, got)
		***REMOVED***
	***REMOVED***
***REMOVED***
