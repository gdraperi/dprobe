package zfs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/pborman/uuid"
)

type command struct ***REMOVED***
	Command string
	Stdin   io.Reader
	Stdout  io.Writer
***REMOVED***

func (c *command) Run(arg ...string) ([][]string, error) ***REMOVED***

	cmd := exec.Command(c.Command, arg...)

	var stdout, stderr bytes.Buffer

	if c.Stdout == nil ***REMOVED***
		cmd.Stdout = &stdout
	***REMOVED*** else ***REMOVED***
		cmd.Stdout = c.Stdout
	***REMOVED***

	if c.Stdin != nil ***REMOVED***
		cmd.Stdin = c.Stdin

	***REMOVED***
	cmd.Stderr = &stderr

	id := uuid.New()
	joinedArgs := strings.Join(cmd.Args, " ")

	logger.Log([]string***REMOVED***"ID:" + id, "START", joinedArgs***REMOVED***)
	err := cmd.Run()
	logger.Log([]string***REMOVED***"ID:" + id, "FINISH"***REMOVED***)

	if err != nil ***REMOVED***
		return nil, &Error***REMOVED***
			Err:    err,
			Debug:  strings.Join([]string***REMOVED***cmd.Path, joinedArgs***REMOVED***, " "),
			Stderr: stderr.String(),
		***REMOVED***
	***REMOVED***

	// assume if you passed in something for stdout, that you know what to do with it
	if c.Stdout != nil ***REMOVED***
		return nil, nil
	***REMOVED***

	lines := strings.Split(stdout.String(), "\n")

	//last line is always blank
	lines = lines[0 : len(lines)-1]
	output := make([][]string, len(lines))

	for i, l := range lines ***REMOVED***
		output[i] = strings.Fields(l)
	***REMOVED***

	return output, nil
***REMOVED***

func setString(field *string, value string) ***REMOVED***
	v := ""
	if value != "-" ***REMOVED***
		v = value
	***REMOVED***
	*field = v
***REMOVED***

func setUint(field *uint64, value string) error ***REMOVED***
	var v uint64
	if value != "-" ***REMOVED***
		var err error
		v, err = strconv.ParseUint(value, 10, 64)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	*field = v
	return nil
***REMOVED***

func (ds *Dataset) parseLine(line []string) error ***REMOVED***
	var err error

	if len(line) != len(dsPropList) ***REMOVED***
		return errors.New("Output does not match what is expected on this platform")
	***REMOVED***
	setString(&ds.Name, line[0])
	setString(&ds.Origin, line[1])

	if err = setUint(&ds.Used, line[2]); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err = setUint(&ds.Avail, line[3]); err != nil ***REMOVED***
		return err
	***REMOVED***

	setString(&ds.Mountpoint, line[4])
	setString(&ds.Compression, line[5])
	setString(&ds.Type, line[6])

	if err = setUint(&ds.Volsize, line[7]); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err = setUint(&ds.Quota, line[8]); err != nil ***REMOVED***
		return err
	***REMOVED***

	if runtime.GOOS == "solaris" ***REMOVED***
		return nil
	***REMOVED***

	if err = setUint(&ds.Written, line[9]); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err = setUint(&ds.Logicalused, line[10]); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err = setUint(&ds.Usedbydataset, line[11]); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

/*
 * from zfs diff`s escape function:
 *
 * Prints a file name out a character at a time.  If the character is
 * not in the range of what we consider "printable" ASCII, display it
 * as an escaped 3-digit octal value.  ASCII values less than a space
 * are all control characters and we declare the upper end as the
 * DELete character.  This also is the last 7-bit ASCII character.
 * We choose to treat all 8-bit ASCII as not printable for this
 * application.
 */
func unescapeFilepath(path string) (string, error) ***REMOVED***
	buf := make([]byte, 0, len(path))
	llen := len(path)
	for i := 0; i < llen; ***REMOVED***
		if path[i] == '\\' ***REMOVED***
			if llen < i+4 ***REMOVED***
				return "", fmt.Errorf("Invalid octal code: too short")
			***REMOVED***
			octalCode := path[(i + 1):(i + 4)]
			val, err := strconv.ParseUint(octalCode, 8, 8)
			if err != nil ***REMOVED***
				return "", fmt.Errorf("Invalid octal code: %v", err)
			***REMOVED***
			buf = append(buf, byte(val))
			i += 4
		***REMOVED*** else ***REMOVED***
			buf = append(buf, path[i])
			i++
		***REMOVED***
	***REMOVED***
	return string(buf), nil
***REMOVED***

var changeTypeMap = map[string]ChangeType***REMOVED***
	"-": Removed,
	"+": Created,
	"M": Modified,
	"R": Renamed,
***REMOVED***
var inodeTypeMap = map[string]InodeType***REMOVED***
	"B": BlockDevice,
	"C": CharacterDevice,
	"/": Directory,
	">": Door,
	"|": NamedPipe,
	"@": SymbolicLink,
	"P": EventPort,
	"=": Socket,
	"F": File,
***REMOVED***

// matches (+1) or (-1)
var referenceCountRegex = regexp.MustCompile("\\(([+-]\\d+?)\\)")

func parseReferenceCount(field string) (int, error) ***REMOVED***
	matches := referenceCountRegex.FindStringSubmatch(field)
	if matches == nil ***REMOVED***
		return 0, fmt.Errorf("Regexp does not match")
	***REMOVED***
	return strconv.Atoi(matches[1])
***REMOVED***

func parseInodeChange(line []string) (*InodeChange, error) ***REMOVED***
	llen := len(line)
	if llen < 1 ***REMOVED***
		return nil, fmt.Errorf("Empty line passed")
	***REMOVED***

	changeType := changeTypeMap[line[0]]
	if changeType == 0 ***REMOVED***
		return nil, fmt.Errorf("Unknown change type '%s'", line[0])
	***REMOVED***

	switch changeType ***REMOVED***
	case Renamed:
		if llen != 4 ***REMOVED***
			return nil, fmt.Errorf("Mismatching number of fields: expect 4, got: %d", llen)
		***REMOVED***
	case Modified:
		if llen != 4 && llen != 3 ***REMOVED***
			return nil, fmt.Errorf("Mismatching number of fields: expect 3..4, got: %d", llen)
		***REMOVED***
	default:
		if llen != 3 ***REMOVED***
			return nil, fmt.Errorf("Mismatching number of fields: expect 3, got: %d", llen)
		***REMOVED***
	***REMOVED***

	inodeType := inodeTypeMap[line[1]]
	if inodeType == 0 ***REMOVED***
		return nil, fmt.Errorf("Unknown inode type '%s'", line[1])
	***REMOVED***

	path, err := unescapeFilepath(line[2])
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("Failed to parse filename: %v", err)
	***REMOVED***

	var newPath string
	var referenceCount int
	switch changeType ***REMOVED***
	case Renamed:
		newPath, err = unescapeFilepath(line[3])
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("Failed to parse filename: %v", err)
		***REMOVED***
	case Modified:
		if llen == 4 ***REMOVED***
			referenceCount, err = parseReferenceCount(line[3])
			if err != nil ***REMOVED***
				return nil, fmt.Errorf("Failed to parse reference count: %v", err)
			***REMOVED***
		***REMOVED***
	default:
		newPath = ""
	***REMOVED***

	return &InodeChange***REMOVED***
		Change:               changeType,
		Type:                 inodeType,
		Path:                 path,
		NewPath:              newPath,
		ReferenceCountChange: referenceCount,
	***REMOVED***, nil
***REMOVED***

// example input
//M       /       /testpool/bar/
//+       F       /testpool/bar/hello.txt
//M       /       /testpool/bar/hello.txt (+1)
//M       /       /testpool/bar/hello-hardlink
func parseInodeChanges(lines [][]string) ([]*InodeChange, error) ***REMOVED***
	changes := make([]*InodeChange, len(lines))

	for i, line := range lines ***REMOVED***
		c, err := parseInodeChange(line)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("Failed to parse line %d of zfs diff: %v, got: '%s'", i, err, line)
		***REMOVED***
		changes[i] = c
	***REMOVED***
	return changes, nil
***REMOVED***

func listByType(t, filter string) ([]*Dataset, error) ***REMOVED***
	args := []string***REMOVED***"list", "-rHp", "-t", t, "-o", dsPropListOptions***REMOVED***

	if filter != "" ***REMOVED***
		args = append(args, filter)
	***REMOVED***
	out, err := zfs(args...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var datasets []*Dataset

	name := ""
	var ds *Dataset
	for _, line := range out ***REMOVED***
		if name != line[0] ***REMOVED***
			name = line[0]
			ds = &Dataset***REMOVED***Name: name***REMOVED***
			datasets = append(datasets, ds)
		***REMOVED***
		if err := ds.parseLine(line); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return datasets, nil
***REMOVED***

func propsSlice(properties map[string]string) []string ***REMOVED***
	args := make([]string, 0, len(properties)*3)
	for k, v := range properties ***REMOVED***
		args = append(args, "-o")
		args = append(args, fmt.Sprintf("%s=%s", k, v))
	***REMOVED***
	return args
***REMOVED***

func (z *Zpool) parseLine(line []string) error ***REMOVED***
	prop := line[1]
	val := line[2]

	var err error

	switch prop ***REMOVED***
	case "name":
		setString(&z.Name, val)
	case "health":
		setString(&z.Health, val)
	case "allocated":
		err = setUint(&z.Allocated, val)
	case "size":
		err = setUint(&z.Size, val)
	case "free":
		err = setUint(&z.Free, val)
	case "fragmentation":
		// Trim trailing "%" before parsing uint
		err = setUint(&z.Fragmentation, val[:len(val)-1])
	case "readonly":
		z.ReadOnly = val == "on"
	case "freeing":
		err = setUint(&z.Freeing, val)
	case "leaked":
		err = setUint(&z.Leaked, val)
	case "dedupratio":
		// Trim trailing "x" before parsing float64
		z.DedupRatio, err = strconv.ParseFloat(val[:len(val)-1], 64)
	***REMOVED***
	return err
***REMOVED***
