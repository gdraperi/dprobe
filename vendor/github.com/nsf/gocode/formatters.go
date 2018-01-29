package main

import (
	"fmt"
	"strings"
)

//-------------------------------------------------------------------------
// formatter interfaces
//-------------------------------------------------------------------------

type formatter interface ***REMOVED***
	write_candidates(candidates []candidate, num int)
***REMOVED***

//-------------------------------------------------------------------------
// nice_formatter (just for testing, simple textual output)
//-------------------------------------------------------------------------

type nice_formatter struct***REMOVED******REMOVED***

func (*nice_formatter) write_candidates(candidates []candidate, num int) ***REMOVED***
	if candidates == nil ***REMOVED***
		fmt.Printf("Nothing to complete.\n")
		return
	***REMOVED***

	fmt.Printf("Found %d candidates:\n", len(candidates))
	for _, c := range candidates ***REMOVED***
		abbr := fmt.Sprintf("%s %s %s", c.Class, c.Name, c.Type)
		if c.Class == decl_func ***REMOVED***
			abbr = fmt.Sprintf("%s %s%s", c.Class, c.Name, c.Type[len("func"):])
		***REMOVED***
		fmt.Printf("  %s\n", abbr)
	***REMOVED***
***REMOVED***

//-------------------------------------------------------------------------
// vim_formatter
//-------------------------------------------------------------------------

type vim_formatter struct***REMOVED******REMOVED***

func (*vim_formatter) write_candidates(candidates []candidate, num int) ***REMOVED***
	if candidates == nil ***REMOVED***
		fmt.Print("[0, []]")
		return
	***REMOVED***

	fmt.Printf("[%d, [", num)
	for i, c := range candidates ***REMOVED***
		if i != 0 ***REMOVED***
			fmt.Printf(", ")
		***REMOVED***

		word := c.Name
		if c.Class == decl_func ***REMOVED***
			word += "("
			if strings.HasPrefix(c.Type, "func()") ***REMOVED***
				word += ")"
			***REMOVED***
		***REMOVED***

		abbr := fmt.Sprintf("%s %s %s", c.Class, c.Name, c.Type)
		if c.Class == decl_func ***REMOVED***
			abbr = fmt.Sprintf("%s %s%s", c.Class, c.Name, c.Type[len("func"):])
		***REMOVED***
		fmt.Printf("***REMOVED***'word': '%s', 'abbr': '%s', 'info': '%s'***REMOVED***", word, abbr, abbr)
	***REMOVED***
	fmt.Printf("]]")
***REMOVED***

//-------------------------------------------------------------------------
// godit_formatter
//-------------------------------------------------------------------------

type godit_formatter struct***REMOVED******REMOVED***

func (*godit_formatter) write_candidates(candidates []candidate, num int) ***REMOVED***
	fmt.Printf("%d,,%d\n", num, len(candidates))
	for _, c := range candidates ***REMOVED***
		contents := c.Name
		if c.Class == decl_func ***REMOVED***
			contents += "("
			if strings.HasPrefix(c.Type, "func()") ***REMOVED***
				contents += ")"
			***REMOVED***
		***REMOVED***

		display := fmt.Sprintf("%s %s %s", c.Class, c.Name, c.Type)
		if c.Class == decl_func ***REMOVED***
			display = fmt.Sprintf("%s %s%s", c.Class, c.Name, c.Type[len("func"):])
		***REMOVED***
		fmt.Printf("%s,,%s\n", display, contents)
	***REMOVED***
***REMOVED***

//-------------------------------------------------------------------------
// emacs_formatter
//-------------------------------------------------------------------------

type emacs_formatter struct***REMOVED******REMOVED***

func (*emacs_formatter) write_candidates(candidates []candidate, num int) ***REMOVED***
	for _, c := range candidates ***REMOVED***
		var hint string
		switch ***REMOVED***
		case c.Class == decl_func:
			hint = c.Type
		case c.Type == "":
			hint = c.Class.String()
		default:
			hint = c.Class.String() + " " + c.Type
		***REMOVED***
		fmt.Printf("%s,,%s\n", c.Name, hint)
	***REMOVED***
***REMOVED***

//-------------------------------------------------------------------------
// csv_formatter
//-------------------------------------------------------------------------

type csv_formatter struct***REMOVED******REMOVED***

func (*csv_formatter) write_candidates(candidates []candidate, num int) ***REMOVED***
	for _, c := range candidates ***REMOVED***
		fmt.Printf("%s,,%s,,%s\n", c.Class, c.Name, c.Type)
	***REMOVED***
***REMOVED***

//-------------------------------------------------------------------------
// csv_with_package_formatter
//-------------------------------------------------------------------------

type csv_with_package_formatter struct***REMOVED******REMOVED***

func (*csv_with_package_formatter) write_candidates(candidates []candidate, num int) ***REMOVED***
	for _, c := range candidates ***REMOVED***
		fmt.Printf("%s,,%s,,%s,,%s\n", c.Class, c.Name, c.Type, c.Package)
	***REMOVED***
***REMOVED***

//-------------------------------------------------------------------------
// json_formatter
//-------------------------------------------------------------------------

type json_formatter struct***REMOVED******REMOVED***

func (*json_formatter) write_candidates(candidates []candidate, num int) ***REMOVED***
	if candidates == nil ***REMOVED***
		fmt.Print("[]")
		return
	***REMOVED***

	fmt.Printf(`[%d, [`, num)
	for i, c := range candidates ***REMOVED***
		if i != 0 ***REMOVED***
			fmt.Printf(", ")
		***REMOVED***
		fmt.Printf(`***REMOVED***"class": "%s", "name": "%s", "type": "%s", "package": "%s"***REMOVED***`,
			c.Class, c.Name, c.Type, c.Package)
	***REMOVED***
	fmt.Print("]]")
***REMOVED***

//-------------------------------------------------------------------------

func get_formatter(name string) formatter ***REMOVED***
	switch name ***REMOVED***
	case "vim":
		return new(vim_formatter)
	case "emacs":
		return new(emacs_formatter)
	case "nice":
		return new(nice_formatter)
	case "csv":
		return new(csv_formatter)
	case "csv-with-package":
		return new(csv_with_package_formatter)
	case "json":
		return new(json_formatter)
	case "godit":
		return new(godit_formatter)
	***REMOVED***
	return new(nice_formatter)
***REMOVED***
