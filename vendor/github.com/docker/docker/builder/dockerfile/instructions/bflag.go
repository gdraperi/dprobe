package instructions

import (
	"fmt"
	"strings"
)

// FlagType is the type of the build flag
type FlagType int

const (
	boolType FlagType = iota
	stringType
)

// BFlags contains all flags information for the builder
type BFlags struct ***REMOVED***
	Args  []string // actual flags/args from cmd line
	flags map[string]*Flag
	used  map[string]*Flag
	Err   error
***REMOVED***

// Flag contains all information for a flag
type Flag struct ***REMOVED***
	bf       *BFlags
	name     string
	flagType FlagType
	Value    string
***REMOVED***

// NewBFlags returns the new BFlags struct
func NewBFlags() *BFlags ***REMOVED***
	return &BFlags***REMOVED***
		flags: make(map[string]*Flag),
		used:  make(map[string]*Flag),
	***REMOVED***
***REMOVED***

// NewBFlagsWithArgs returns the new BFlags struct with Args set to args
func NewBFlagsWithArgs(args []string) *BFlags ***REMOVED***
	flags := NewBFlags()
	flags.Args = args
	return flags
***REMOVED***

// AddBool adds a bool flag to BFlags
// Note, any error will be generated when Parse() is called (see Parse).
func (bf *BFlags) AddBool(name string, def bool) *Flag ***REMOVED***
	flag := bf.addFlag(name, boolType)
	if flag == nil ***REMOVED***
		return nil
	***REMOVED***
	if def ***REMOVED***
		flag.Value = "true"
	***REMOVED*** else ***REMOVED***
		flag.Value = "false"
	***REMOVED***
	return flag
***REMOVED***

// AddString adds a string flag to BFlags
// Note, any error will be generated when Parse() is called (see Parse).
func (bf *BFlags) AddString(name string, def string) *Flag ***REMOVED***
	flag := bf.addFlag(name, stringType)
	if flag == nil ***REMOVED***
		return nil
	***REMOVED***
	flag.Value = def
	return flag
***REMOVED***

// addFlag is a generic func used by the other AddXXX() func
// to add a new flag to the BFlags struct.
// Note, any error will be generated when Parse() is called (see Parse).
func (bf *BFlags) addFlag(name string, flagType FlagType) *Flag ***REMOVED***
	if _, ok := bf.flags[name]; ok ***REMOVED***
		bf.Err = fmt.Errorf("Duplicate flag defined: %s", name)
		return nil
	***REMOVED***

	newFlag := &Flag***REMOVED***
		bf:       bf,
		name:     name,
		flagType: flagType,
	***REMOVED***
	bf.flags[name] = newFlag

	return newFlag
***REMOVED***

// IsUsed checks if the flag is used
func (fl *Flag) IsUsed() bool ***REMOVED***
	if _, ok := fl.bf.used[fl.name]; ok ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

// IsTrue checks if a bool flag is true
func (fl *Flag) IsTrue() bool ***REMOVED***
	if fl.flagType != boolType ***REMOVED***
		// Should never get here
		panic(fmt.Errorf("Trying to use IsTrue on a non-boolean: %s", fl.name))
	***REMOVED***
	return fl.Value == "true"
***REMOVED***

// Parse parses and checks if the BFlags is valid.
// Any error noticed during the AddXXX() funcs will be generated/returned
// here.  We do this because an error during AddXXX() is more like a
// compile time error so it doesn't matter too much when we stop our
// processing as long as we do stop it, so this allows the code
// around AddXXX() to be just:
//     defFlag := AddString("description", "")
// w/o needing to add an if-statement around each one.
func (bf *BFlags) Parse() error ***REMOVED***
	// If there was an error while defining the possible flags
	// go ahead and bubble it back up here since we didn't do it
	// earlier in the processing
	if bf.Err != nil ***REMOVED***
		return fmt.Errorf("Error setting up flags: %s", bf.Err)
	***REMOVED***

	for _, arg := range bf.Args ***REMOVED***
		if !strings.HasPrefix(arg, "--") ***REMOVED***
			return fmt.Errorf("Arg should start with -- : %s", arg)
		***REMOVED***

		if arg == "--" ***REMOVED***
			return nil
		***REMOVED***

		arg = arg[2:]
		value := ""

		index := strings.Index(arg, "=")
		if index >= 0 ***REMOVED***
			value = arg[index+1:]
			arg = arg[:index]
		***REMOVED***

		flag, ok := bf.flags[arg]
		if !ok ***REMOVED***
			return fmt.Errorf("Unknown flag: %s", arg)
		***REMOVED***

		if _, ok = bf.used[arg]; ok ***REMOVED***
			return fmt.Errorf("Duplicate flag specified: %s", arg)
		***REMOVED***

		bf.used[arg] = flag

		switch flag.flagType ***REMOVED***
		case boolType:
			// value == "" is only ok if no "=" was specified
			if index >= 0 && value == "" ***REMOVED***
				return fmt.Errorf("Missing a value on flag: %s", arg)
			***REMOVED***

			lower := strings.ToLower(value)
			if lower == "" ***REMOVED***
				flag.Value = "true"
			***REMOVED*** else if lower == "true" || lower == "false" ***REMOVED***
				flag.Value = lower
			***REMOVED*** else ***REMOVED***
				return fmt.Errorf("Expecting boolean value for flag %s, not: %s", arg, value)
			***REMOVED***

		case stringType:
			if index < 0 ***REMOVED***
				return fmt.Errorf("Missing a value on flag: %s", arg)
			***REMOVED***
			flag.Value = value

		default:
			panic("No idea what kind of flag we have! Should never get here!")
		***REMOVED***

	***REMOVED***

	return nil
***REMOVED***
