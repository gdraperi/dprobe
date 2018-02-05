package cmd

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

// TestGoldenInitCmd initializes the project "github.com/spf13/testproject"
// in GOPATH and compares the content of files in initialized project with
// appropriate golden files ("testdata/*.golden").
// Use -update to update existing golden files.
func TestGoldenInitCmd(t *testing.T) ***REMOVED***
	projectName := "github.com/spf13/testproject"
	project := NewProject(projectName)
	defer os.RemoveAll(project.AbsPath())

	viper.Set("author", "NAME HERE <EMAIL ADDRESS>")
	viper.Set("license", "apache")
	viper.Set("year", 2017)
	defer viper.Set("author", nil)
	defer viper.Set("license", nil)
	defer viper.Set("year", nil)

	os.Args = []string***REMOVED***"cobra", "init", projectName***REMOVED***
	if err := rootCmd.Execute(); err != nil ***REMOVED***
		t.Fatal("Error by execution:", err)
	***REMOVED***

	expectedFiles := []string***REMOVED***".", "cmd", "LICENSE", "main.go", "cmd/root.go"***REMOVED***
	gotFiles := []string***REMOVED******REMOVED***

	// Check project file hierarchy and compare the content of every single file
	// with appropriate golden file.
	err := filepath.Walk(project.AbsPath(), func(path string, info os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Make path relative to project.AbsPath().
		// E.g. path = "/home/user/go/src/github.com/spf13/testproject/cmd/root.go"
		// then it returns just "cmd/root.go".
		relPath, err := filepath.Rel(project.AbsPath(), path)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		relPath = filepath.ToSlash(relPath)
		gotFiles = append(gotFiles, relPath)
		goldenPath := filepath.Join("testdata", filepath.Base(path)+".golden")

		switch relPath ***REMOVED***
		// Known directories.
		case ".", "cmd":
			return nil
		// Known files.
		case "LICENSE", "main.go", "cmd/root.go":
			if *update ***REMOVED***
				got, err := ioutil.ReadFile(path)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				if err := ioutil.WriteFile(goldenPath, got, 0644); err != nil ***REMOVED***
					t.Fatal("Error while updating file:", err)
				***REMOVED***
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
