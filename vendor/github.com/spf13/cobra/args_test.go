package cobra

import (
	"strings"
	"testing"
)

func TestNoArgs(t *testing.T) ***REMOVED***
	c := &Command***REMOVED***Use: "c", Args: NoArgs, Run: emptyRun***REMOVED***

	output, err := executeCommand(c)
	if output != "" ***REMOVED***
		t.Errorf("Unexpected string: %v", output)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Fatalf("Unexpected error: %v", err)
	***REMOVED***
***REMOVED***

func TestNoArgsWithArgs(t *testing.T) ***REMOVED***
	c := &Command***REMOVED***Use: "c", Args: NoArgs, Run: emptyRun***REMOVED***

	_, err := executeCommand(c, "illegal")
	if err == nil ***REMOVED***
		t.Fatal("Expected an error")
	***REMOVED***

	got := err.Error()
	expected := `unknown command "illegal" for "c"`
	if got != expected ***REMOVED***
		t.Errorf("Expected: %q, got: %q", expected, got)
	***REMOVED***
***REMOVED***

func TestOnlyValidArgs(t *testing.T) ***REMOVED***
	c := &Command***REMOVED***
		Use:       "c",
		Args:      OnlyValidArgs,
		ValidArgs: []string***REMOVED***"one", "two"***REMOVED***,
		Run:       emptyRun,
	***REMOVED***

	output, err := executeCommand(c, "one", "two")
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", output)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Fatalf("Unexpected error: %v", err)
	***REMOVED***
***REMOVED***

func TestOnlyValidArgsWithInvalidArgs(t *testing.T) ***REMOVED***
	c := &Command***REMOVED***
		Use:       "c",
		Args:      OnlyValidArgs,
		ValidArgs: []string***REMOVED***"one", "two"***REMOVED***,
		Run:       emptyRun,
	***REMOVED***

	_, err := executeCommand(c, "three")
	if err == nil ***REMOVED***
		t.Fatal("Expected an error")
	***REMOVED***

	got := err.Error()
	expected := `invalid argument "three" for "c"`
	if got != expected ***REMOVED***
		t.Errorf("Expected: %q, got: %q", expected, got)
	***REMOVED***
***REMOVED***

func TestArbitraryArgs(t *testing.T) ***REMOVED***
	c := &Command***REMOVED***Use: "c", Args: ArbitraryArgs, Run: emptyRun***REMOVED***
	output, err := executeCommand(c, "a", "b")
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", output)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***
***REMOVED***

func TestMinimumNArgs(t *testing.T) ***REMOVED***
	c := &Command***REMOVED***Use: "c", Args: MinimumNArgs(2), Run: emptyRun***REMOVED***
	output, err := executeCommand(c, "a", "b", "c")
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", output)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***
***REMOVED***

func TestMinimumNArgsWithLessArgs(t *testing.T) ***REMOVED***
	c := &Command***REMOVED***Use: "c", Args: MinimumNArgs(2), Run: emptyRun***REMOVED***
	_, err := executeCommand(c, "a")

	if err == nil ***REMOVED***
		t.Fatal("Expected an error")
	***REMOVED***

	got := err.Error()
	expected := "requires at least 2 arg(s), only received 1"
	if got != expected ***REMOVED***
		t.Fatalf("Expected %q, got %q", expected, got)
	***REMOVED***
***REMOVED***

func TestMaximumNArgs(t *testing.T) ***REMOVED***
	c := &Command***REMOVED***Use: "c", Args: MaximumNArgs(3), Run: emptyRun***REMOVED***
	output, err := executeCommand(c, "a", "b")
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", output)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***
***REMOVED***

func TestMaximumNArgsWithMoreArgs(t *testing.T) ***REMOVED***
	c := &Command***REMOVED***Use: "c", Args: MaximumNArgs(2), Run: emptyRun***REMOVED***
	_, err := executeCommand(c, "a", "b", "c")

	if err == nil ***REMOVED***
		t.Fatal("Expected an error")
	***REMOVED***

	got := err.Error()
	expected := "accepts at most 2 arg(s), received 3"
	if got != expected ***REMOVED***
		t.Fatalf("Expected %q, got %q", expected, got)
	***REMOVED***
***REMOVED***

func TestExactArgs(t *testing.T) ***REMOVED***
	c := &Command***REMOVED***Use: "c", Args: ExactArgs(3), Run: emptyRun***REMOVED***
	output, err := executeCommand(c, "a", "b", "c")
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", output)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***
***REMOVED***

func TestExactArgsWithInvalidCount(t *testing.T) ***REMOVED***
	c := &Command***REMOVED***Use: "c", Args: ExactArgs(2), Run: emptyRun***REMOVED***
	_, err := executeCommand(c, "a", "b", "c")

	if err == nil ***REMOVED***
		t.Fatal("Expected an error")
	***REMOVED***

	got := err.Error()
	expected := "accepts 2 arg(s), received 3"
	if got != expected ***REMOVED***
		t.Fatalf("Expected %q, got %q", expected, got)
	***REMOVED***
***REMOVED***

func TestRangeArgs(t *testing.T) ***REMOVED***
	c := &Command***REMOVED***Use: "c", Args: RangeArgs(2, 4), Run: emptyRun***REMOVED***
	output, err := executeCommand(c, "a", "b", "c")
	if output != "" ***REMOVED***
		t.Errorf("Unexpected output: %v", output)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***
***REMOVED***

func TestRangeArgsWithInvalidCount(t *testing.T) ***REMOVED***
	c := &Command***REMOVED***Use: "c", Args: RangeArgs(2, 4), Run: emptyRun***REMOVED***
	_, err := executeCommand(c, "a")

	if err == nil ***REMOVED***
		t.Fatal("Expected an error")
	***REMOVED***

	got := err.Error()
	expected := "accepts between 2 and 4 arg(s), received 1"
	if got != expected ***REMOVED***
		t.Fatalf("Expected %q, got %q", expected, got)
	***REMOVED***
***REMOVED***

func TestRootTakesNoArgs(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Run: emptyRun***REMOVED***
	childCmd := &Command***REMOVED***Use: "child", Run: emptyRun***REMOVED***
	rootCmd.AddCommand(childCmd)

	_, err := executeCommand(rootCmd, "illegal", "args")
	if err == nil ***REMOVED***
		t.Fatal("Expected an error")
	***REMOVED***

	got := err.Error()
	expected := `unknown command "illegal" for "root"`
	if !strings.Contains(got, expected) ***REMOVED***
		t.Errorf("expected %q, got %q", expected, got)
	***REMOVED***
***REMOVED***

func TestRootTakesArgs(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Args: ArbitraryArgs, Run: emptyRun***REMOVED***
	childCmd := &Command***REMOVED***Use: "child", Run: emptyRun***REMOVED***
	rootCmd.AddCommand(childCmd)

	_, err := executeCommand(rootCmd, "legal", "args")
	if err != nil ***REMOVED***
		t.Errorf("Unexpected error: %v", err)
	***REMOVED***
***REMOVED***

func TestChildTakesNoArgs(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Run: emptyRun***REMOVED***
	childCmd := &Command***REMOVED***Use: "child", Args: NoArgs, Run: emptyRun***REMOVED***
	rootCmd.AddCommand(childCmd)

	_, err := executeCommand(rootCmd, "child", "illegal", "args")
	if err == nil ***REMOVED***
		t.Fatal("Expected an error")
	***REMOVED***

	got := err.Error()
	expected := `unknown command "illegal" for "root child"`
	if !strings.Contains(got, expected) ***REMOVED***
		t.Errorf("expected %q, got %q", expected, got)
	***REMOVED***
***REMOVED***

func TestChildTakesArgs(t *testing.T) ***REMOVED***
	rootCmd := &Command***REMOVED***Use: "root", Run: emptyRun***REMOVED***
	childCmd := &Command***REMOVED***Use: "child", Args: ArbitraryArgs, Run: emptyRun***REMOVED***
	rootCmd.AddCommand(childCmd)

	_, err := executeCommand(rootCmd, "child", "legal", "args")
	if err != nil ***REMOVED***
		t.Fatalf("Unexpected error: %v", err)
	***REMOVED***
***REMOVED***
