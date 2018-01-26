package volume

import (
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/mount"
)

type parseMountRawTestSet struct ***REMOVED***
	valid   []string
	invalid map[string]string
***REMOVED***

func TestConvertTmpfsOptions(t *testing.T) ***REMOVED***
	type testCase struct ***REMOVED***
		opt                  mount.TmpfsOptions
		readOnly             bool
		expectedSubstrings   []string
		unexpectedSubstrings []string
	***REMOVED***
	cases := []testCase***REMOVED***
		***REMOVED***
			opt:                  mount.TmpfsOptions***REMOVED***SizeBytes: 1024 * 1024, Mode: 0700***REMOVED***,
			readOnly:             false,
			expectedSubstrings:   []string***REMOVED***"size=1m", "mode=700"***REMOVED***,
			unexpectedSubstrings: []string***REMOVED***"ro"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			opt:                  mount.TmpfsOptions***REMOVED******REMOVED***,
			readOnly:             true,
			expectedSubstrings:   []string***REMOVED***"ro"***REMOVED***,
			unexpectedSubstrings: []string***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***
	p := &linuxParser***REMOVED******REMOVED***
	for _, c := range cases ***REMOVED***
		data, err := p.ConvertTmpfsOptions(&c.opt, c.readOnly)
		if err != nil ***REMOVED***
			t.Fatalf("could not convert %+v (readOnly: %v) to string: %v",
				c.opt, c.readOnly, err)
		***REMOVED***
		t.Logf("data=%q", data)
		for _, s := range c.expectedSubstrings ***REMOVED***
			if !strings.Contains(data, s) ***REMOVED***
				t.Fatalf("expected substring: %s, got %v (case=%+v)", s, data, c)
			***REMOVED***
		***REMOVED***
		for _, s := range c.unexpectedSubstrings ***REMOVED***
			if strings.Contains(data, s) ***REMOVED***
				t.Fatalf("unexpected substring: %s, got %v (case=%+v)", s, data, c)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

type mockFiProvider struct***REMOVED******REMOVED***

func (mockFiProvider) fileInfo(path string) (exists, isDir bool, err error) ***REMOVED***
	dirs := map[string]struct***REMOVED******REMOVED******REMOVED***
		`c:\`:                    ***REMOVED******REMOVED***,
		`c:\windows\`:            ***REMOVED******REMOVED***,
		`c:\windows`:             ***REMOVED******REMOVED***,
		`c:\program files`:       ***REMOVED******REMOVED***,
		`c:\Windows`:             ***REMOVED******REMOVED***,
		`c:\Program Files (x86)`: ***REMOVED******REMOVED***,
		`\\?\c:\windows\`:        ***REMOVED******REMOVED***,
	***REMOVED***
	files := map[string]struct***REMOVED******REMOVED******REMOVED***
		`c:\windows\system32\ntdll.dll`: ***REMOVED******REMOVED***,
	***REMOVED***
	if _, ok := dirs[path]; ok ***REMOVED***
		return true, true, nil
	***REMOVED***
	if _, ok := files[path]; ok ***REMOVED***
		return true, false, nil
	***REMOVED***
	return false, false, nil
***REMOVED***

func TestParseMountRaw(t *testing.T) ***REMOVED***

	previousProvider := currentFileInfoProvider
	defer func() ***REMOVED*** currentFileInfoProvider = previousProvider ***REMOVED***()
	currentFileInfoProvider = mockFiProvider***REMOVED******REMOVED***
	windowsSet := parseMountRawTestSet***REMOVED***
		valid: []string***REMOVED***
			`d:\`,
			`d:`,
			`d:\path`,
			`d:\path with space`,
			`c:\:d:\`,
			`c:\windows\:d:`,
			`c:\windows:d:\s p a c e`,
			`c:\windows:d:\s p a c e:RW`,
			`c:\program files:d:\s p a c e i n h o s t d i r`,
			`0123456789name:d:`,
			`MiXeDcAsEnAmE:d:`,
			`name:D:`,
			`name:D::rW`,
			`name:D::RW`,
			`name:D::RO`,
			`c:/:d:/forward/slashes/are/good/too`,
			`c:/:d:/including with/spaces:ro`,
			`c:\Windows`,                // With capital
			`c:\Program Files (x86)`,    // With capitals and brackets
			`\\?\c:\windows\:d:`,        // Long path handling (source)
			`c:\windows\:\\?\d:\`,       // Long path handling (target)
			`\\.\pipe\foo:\\.\pipe\foo`, // named pipe
			`//./pipe/foo://./pipe/foo`, // named pipe forward slashes
		***REMOVED***,
		invalid: map[string]string***REMOVED***
			``:                                 "invalid volume specification: ",
			`.`:                                "invalid volume specification: ",
			`..\`:                              "invalid volume specification: ",
			`c:\:..\`:                          "invalid volume specification: ",
			`c:\:d:\:xyzzy`:                    "invalid volume specification: ",
			`c:`:                               "cannot be `c:`",
			`c:\`:                              "cannot be `c:`",
			`c:\notexist:d:`:                   `source path does not exist`,
			`c:\windows\system32\ntdll.dll:d:`: `source path must be a directory`,
			`name<:d:`:                         `invalid volume specification`,
			`name>:d:`:                         `invalid volume specification`,
			`name::d:`:                         `invalid volume specification`,
			`name":d:`:                         `invalid volume specification`,
			`name\:d:`:                         `invalid volume specification`,
			`name*:d:`:                         `invalid volume specification`,
			`name|:d:`:                         `invalid volume specification`,
			`name?:d:`:                         `invalid volume specification`,
			`name/:d:`:                         `invalid volume specification`,
			`d:\pathandmode:rw`:                `invalid volume specification`,
			`d:\pathandmode:ro`:                `invalid volume specification`,
			`con:d:`:                           `cannot be a reserved word for Windows filenames`,
			`PRN:d:`:                           `cannot be a reserved word for Windows filenames`,
			`aUx:d:`:                           `cannot be a reserved word for Windows filenames`,
			`nul:d:`:                           `cannot be a reserved word for Windows filenames`,
			`com1:d:`:                          `cannot be a reserved word for Windows filenames`,
			`com2:d:`:                          `cannot be a reserved word for Windows filenames`,
			`com3:d:`:                          `cannot be a reserved word for Windows filenames`,
			`com4:d:`:                          `cannot be a reserved word for Windows filenames`,
			`com5:d:`:                          `cannot be a reserved word for Windows filenames`,
			`com6:d:`:                          `cannot be a reserved word for Windows filenames`,
			`com7:d:`:                          `cannot be a reserved word for Windows filenames`,
			`com8:d:`:                          `cannot be a reserved word for Windows filenames`,
			`com9:d:`:                          `cannot be a reserved word for Windows filenames`,
			`lpt1:d:`:                          `cannot be a reserved word for Windows filenames`,
			`lpt2:d:`:                          `cannot be a reserved word for Windows filenames`,
			`lpt3:d:`:                          `cannot be a reserved word for Windows filenames`,
			`lpt4:d:`:                          `cannot be a reserved word for Windows filenames`,
			`lpt5:d:`:                          `cannot be a reserved word for Windows filenames`,
			`lpt6:d:`:                          `cannot be a reserved word for Windows filenames`,
			`lpt7:d:`:                          `cannot be a reserved word for Windows filenames`,
			`lpt8:d:`:                          `cannot be a reserved word for Windows filenames`,
			`lpt9:d:`:                          `cannot be a reserved word for Windows filenames`,
			`c:\windows\system32\ntdll.dll`:    `Only directories can be mapped on this platform`,
			`\\.\pipe\foo:c:\pipe`:             `'c:\pipe' is not a valid pipe path`,
		***REMOVED***,
	***REMOVED***
	lcowSet := parseMountRawTestSet***REMOVED***
		valid: []string***REMOVED***
			`/foo`,
			`/foo/`,
			`/foo bar`,
			`c:\:/foo`,
			`c:\windows\:/foo`,
			`c:\windows:/s p a c e`,
			`c:\windows:/s p a c e:RW`,
			`c:\program files:/s p a c e i n h o s t d i r`,
			`0123456789name:/foo`,
			`MiXeDcAsEnAmE:/foo`,
			`name:/foo`,
			`name:/foo:rW`,
			`name:/foo:RW`,
			`name:/foo:RO`,
			`c:/:/forward/slashes/are/good/too`,
			`c:/:/including with/spaces:ro`,
			`/Program Files (x86)`, // With capitals and brackets
		***REMOVED***,
		invalid: map[string]string***REMOVED***
			``:                                   "invalid volume specification: ",
			`.`:                                  "invalid volume specification: ",
			`c:`:                                 "invalid volume specification: ",
			`c:\`:                                "invalid volume specification: ",
			`../`:                                "invalid volume specification: ",
			`c:\:../`:                            "invalid volume specification: ",
			`c:\:/foo:xyzzy`:                     "invalid volume specification: ",
			`/`:                                  "destination can't be '/'",
			`/..`:                                "destination can't be '/'",
			`c:\notexist:/foo`:                   `source path does not exist`,
			`c:\windows\system32\ntdll.dll:/foo`: `source path must be a directory`,
			`name<:/foo`:                         `invalid volume specification`,
			`name>:/foo`:                         `invalid volume specification`,
			`name::/foo`:                         `invalid volume specification`,
			`name":/foo`:                         `invalid volume specification`,
			`name\:/foo`:                         `invalid volume specification`,
			`name*:/foo`:                         `invalid volume specification`,
			`name|:/foo`:                         `invalid volume specification`,
			`name?:/foo`:                         `invalid volume specification`,
			`name/:/foo`:                         `invalid volume specification`,
			`/foo:rw`:                            `invalid volume specification`,
			`/foo:ro`:                            `invalid volume specification`,
			`con:/foo`:                           `cannot be a reserved word for Windows filenames`,
			`PRN:/foo`:                           `cannot be a reserved word for Windows filenames`,
			`aUx:/foo`:                           `cannot be a reserved word for Windows filenames`,
			`nul:/foo`:                           `cannot be a reserved word for Windows filenames`,
			`com1:/foo`:                          `cannot be a reserved word for Windows filenames`,
			`com2:/foo`:                          `cannot be a reserved word for Windows filenames`,
			`com3:/foo`:                          `cannot be a reserved word for Windows filenames`,
			`com4:/foo`:                          `cannot be a reserved word for Windows filenames`,
			`com5:/foo`:                          `cannot be a reserved word for Windows filenames`,
			`com6:/foo`:                          `cannot be a reserved word for Windows filenames`,
			`com7:/foo`:                          `cannot be a reserved word for Windows filenames`,
			`com8:/foo`:                          `cannot be a reserved word for Windows filenames`,
			`com9:/foo`:                          `cannot be a reserved word for Windows filenames`,
			`lpt1:/foo`:                          `cannot be a reserved word for Windows filenames`,
			`lpt2:/foo`:                          `cannot be a reserved word for Windows filenames`,
			`lpt3:/foo`:                          `cannot be a reserved word for Windows filenames`,
			`lpt4:/foo`:                          `cannot be a reserved word for Windows filenames`,
			`lpt5:/foo`:                          `cannot be a reserved word for Windows filenames`,
			`lpt6:/foo`:                          `cannot be a reserved word for Windows filenames`,
			`lpt7:/foo`:                          `cannot be a reserved word for Windows filenames`,
			`lpt8:/foo`:                          `cannot be a reserved word for Windows filenames`,
			`lpt9:/foo`:                          `cannot be a reserved word for Windows filenames`,
			`\\.\pipe\foo:/foo`:                  `Linux containers on Windows do not support named pipe mounts`,
		***REMOVED***,
	***REMOVED***
	linuxSet := parseMountRawTestSet***REMOVED***
		valid: []string***REMOVED***
			"/home",
			"/home:/home",
			"/home:/something/else",
			"/with space",
			"/home:/with space",
			"relative:/absolute-path",
			"hostPath:/containerPath:ro",
			"/hostPath:/containerPath:rw",
			"/rw:/ro",
			"/hostPath:/containerPath:shared",
			"/hostPath:/containerPath:rshared",
			"/hostPath:/containerPath:slave",
			"/hostPath:/containerPath:rslave",
			"/hostPath:/containerPath:private",
			"/hostPath:/containerPath:rprivate",
			"/hostPath:/containerPath:ro,shared",
			"/hostPath:/containerPath:ro,slave",
			"/hostPath:/containerPath:ro,private",
			"/hostPath:/containerPath:ro,z,shared",
			"/hostPath:/containerPath:ro,Z,slave",
			"/hostPath:/containerPath:Z,ro,slave",
			"/hostPath:/containerPath:slave,Z,ro",
			"/hostPath:/containerPath:Z,slave,ro",
			"/hostPath:/containerPath:slave,ro,Z",
			"/hostPath:/containerPath:rslave,ro,Z",
			"/hostPath:/containerPath:ro,rshared,Z",
			"/hostPath:/containerPath:ro,Z,rprivate",
		***REMOVED***,
		invalid: map[string]string***REMOVED***
			"":                                "invalid volume specification",
			"./":                              "mount path must be absolute",
			"../":                             "mount path must be absolute",
			"/:../":                           "mount path must be absolute",
			"/:path":                          "mount path must be absolute",
			":":                               "invalid volume specification",
			"/tmp:":                           "invalid volume specification",
			":test":                           "invalid volume specification",
			":/test":                          "invalid volume specification",
			"tmp:":                            "invalid volume specification",
			":test:":                          "invalid volume specification",
			"::":                              "invalid volume specification",
			":::":                             "invalid volume specification",
			"/tmp:::":                         "invalid volume specification",
			":/tmp::":                         "invalid volume specification",
			"/path:rw":                        "invalid volume specification",
			"/path:ro":                        "invalid volume specification",
			"/rw:rw":                          "invalid volume specification",
			"path:ro":                         "invalid volume specification",
			"/path:/path:sw":                  `invalid mode`,
			"/path:/path:rwz":                 `invalid mode`,
			"/path:/path:ro,rshared,rslave":   `invalid mode`,
			"/path:/path:ro,z,rshared,rslave": `invalid mode`,
			"/path:shared":                    "invalid volume specification",
			"/path:slave":                     "invalid volume specification",
			"/path:private":                   "invalid volume specification",
			"name:/absolute-path:shared":      "invalid volume specification",
			"name:/absolute-path:rshared":     "invalid volume specification",
			"name:/absolute-path:slave":       "invalid volume specification",
			"name:/absolute-path:rslave":      "invalid volume specification",
			"name:/absolute-path:private":     "invalid volume specification",
			"name:/absolute-path:rprivate":    "invalid volume specification",
		***REMOVED***,
	***REMOVED***

	linParser := &linuxParser***REMOVED******REMOVED***
	winParser := &windowsParser***REMOVED******REMOVED***
	lcowParser := &lcowParser***REMOVED******REMOVED***
	tester := func(parser Parser, set parseMountRawTestSet) ***REMOVED***

		for _, path := range set.valid ***REMOVED***

			if _, err := parser.ParseMountRaw(path, "local"); err != nil ***REMOVED***
				t.Errorf("ParseMountRaw(`%q`) should succeed: error %q", path, err)
			***REMOVED***
		***REMOVED***

		for path, expectedError := range set.invalid ***REMOVED***
			if mp, err := parser.ParseMountRaw(path, "local"); err == nil ***REMOVED***
				t.Errorf("ParseMountRaw(`%q`) should have failed validation. Err '%v' - MP: %v", path, err, mp)
			***REMOVED*** else ***REMOVED***
				if !strings.Contains(err.Error(), expectedError) ***REMOVED***
					t.Errorf("ParseMountRaw(`%q`) error should contain %q, got %v", path, expectedError, err.Error())
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	tester(linParser, linuxSet)
	tester(winParser, windowsSet)
	tester(lcowParser, lcowSet)

***REMOVED***

// testParseMountRaw is a structure used by TestParseMountRawSplit for
// specifying test cases for the ParseMountRaw() function.
type testParseMountRaw struct ***REMOVED***
	bind      string
	driver    string
	expType   mount.Type
	expDest   string
	expSource string
	expName   string
	expDriver string
	expRW     bool
	fail      bool
***REMOVED***

func TestParseMountRawSplit(t *testing.T) ***REMOVED***
	previousProvider := currentFileInfoProvider
	defer func() ***REMOVED*** currentFileInfoProvider = previousProvider ***REMOVED***()
	currentFileInfoProvider = mockFiProvider***REMOVED******REMOVED***
	windowsCases := []testParseMountRaw***REMOVED***
		***REMOVED***`c:\:d:`, "local", mount.TypeBind, `d:`, `c:\`, ``, "", true, false***REMOVED***,
		***REMOVED***`c:\:d:\`, "local", mount.TypeBind, `d:\`, `c:\`, ``, "", true, false***REMOVED***,
		***REMOVED***`c:\:d:\:ro`, "local", mount.TypeBind, `d:\`, `c:\`, ``, "", false, false***REMOVED***,
		***REMOVED***`c:\:d:\:rw`, "local", mount.TypeBind, `d:\`, `c:\`, ``, "", true, false***REMOVED***,
		***REMOVED***`c:\:d:\:foo`, "local", mount.TypeBind, `d:\`, `c:\`, ``, "", false, true***REMOVED***,
		***REMOVED***`name:d::rw`, "local", mount.TypeVolume, `d:`, ``, `name`, "local", true, false***REMOVED***,
		***REMOVED***`name:d:`, "local", mount.TypeVolume, `d:`, ``, `name`, "local", true, false***REMOVED***,
		***REMOVED***`name:d::ro`, "local", mount.TypeVolume, `d:`, ``, `name`, "local", false, false***REMOVED***,
		***REMOVED***`name:c:`, "", mount.TypeVolume, ``, ``, ``, "", true, true***REMOVED***,
		***REMOVED***`driver/name:c:`, "", mount.TypeVolume, ``, ``, ``, "", true, true***REMOVED***,
		***REMOVED***`\\.\pipe\foo:\\.\pipe\bar`, "local", mount.TypeNamedPipe, `\\.\pipe\bar`, `\\.\pipe\foo`, "", "", true, false***REMOVED***,
		***REMOVED***`\\.\pipe\foo:c:\foo\bar`, "local", mount.TypeNamedPipe, ``, ``, "", "", true, true***REMOVED***,
		***REMOVED***`c:\foo\bar:\\.\pipe\foo`, "local", mount.TypeNamedPipe, ``, ``, "", "", true, true***REMOVED***,
	***REMOVED***
	lcowCases := []testParseMountRaw***REMOVED***
		***REMOVED***`c:\:/foo`, "local", mount.TypeBind, `/foo`, `c:\`, ``, "", true, false***REMOVED***,
		***REMOVED***`c:\:/foo:ro`, "local", mount.TypeBind, `/foo`, `c:\`, ``, "", false, false***REMOVED***,
		***REMOVED***`c:\:/foo:rw`, "local", mount.TypeBind, `/foo`, `c:\`, ``, "", true, false***REMOVED***,
		***REMOVED***`c:\:/foo:foo`, "local", mount.TypeBind, `/foo`, `c:\`, ``, "", false, true***REMOVED***,
		***REMOVED***`name:/foo:rw`, "local", mount.TypeVolume, `/foo`, ``, `name`, "local", true, false***REMOVED***,
		***REMOVED***`name:/foo`, "local", mount.TypeVolume, `/foo`, ``, `name`, "local", true, false***REMOVED***,
		***REMOVED***`name:/foo:ro`, "local", mount.TypeVolume, `/foo`, ``, `name`, "local", false, false***REMOVED***,
		***REMOVED***`name:/`, "", mount.TypeVolume, ``, ``, ``, "", true, true***REMOVED***,
		***REMOVED***`driver/name:/`, "", mount.TypeVolume, ``, ``, ``, "", true, true***REMOVED***,
		***REMOVED***`\\.\pipe\foo:\\.\pipe\bar`, "local", mount.TypeNamedPipe, `\\.\pipe\bar`, `\\.\pipe\foo`, "", "", true, true***REMOVED***,
		***REMOVED***`\\.\pipe\foo:/data`, "local", mount.TypeNamedPipe, ``, ``, "", "", true, true***REMOVED***,
		***REMOVED***`c:\foo\bar:\\.\pipe\foo`, "local", mount.TypeNamedPipe, ``, ``, "", "", true, true***REMOVED***,
	***REMOVED***
	linuxCases := []testParseMountRaw***REMOVED***
		***REMOVED***"/tmp:/tmp1", "", mount.TypeBind, "/tmp1", "/tmp", "", "", true, false***REMOVED***,
		***REMOVED***"/tmp:/tmp2:ro", "", mount.TypeBind, "/tmp2", "/tmp", "", "", false, false***REMOVED***,
		***REMOVED***"/tmp:/tmp3:rw", "", mount.TypeBind, "/tmp3", "/tmp", "", "", true, false***REMOVED***,
		***REMOVED***"/tmp:/tmp4:foo", "", mount.TypeBind, "", "", "", "", false, true***REMOVED***,
		***REMOVED***"name:/named1", "", mount.TypeVolume, "/named1", "", "name", "", true, false***REMOVED***,
		***REMOVED***"name:/named2", "external", mount.TypeVolume, "/named2", "", "name", "external", true, false***REMOVED***,
		***REMOVED***"name:/named3:ro", "local", mount.TypeVolume, "/named3", "", "name", "local", false, false***REMOVED***,
		***REMOVED***"local/name:/tmp:rw", "", mount.TypeVolume, "/tmp", "", "local/name", "", true, false***REMOVED***,
		***REMOVED***"/tmp:tmp", "", mount.TypeBind, "", "", "", "", true, true***REMOVED***,
	***REMOVED***
	linParser := &linuxParser***REMOVED******REMOVED***
	winParser := &windowsParser***REMOVED******REMOVED***
	lcowParser := &lcowParser***REMOVED******REMOVED***
	tester := func(parser Parser, cases []testParseMountRaw) ***REMOVED***
		for i, c := range cases ***REMOVED***
			t.Logf("case %d", i)
			m, err := parser.ParseMountRaw(c.bind, c.driver)
			if c.fail ***REMOVED***
				if err == nil ***REMOVED***
					t.Errorf("Expected error, was nil, for spec %s\n", c.bind)
				***REMOVED***
				continue
			***REMOVED***

			if m == nil || err != nil ***REMOVED***
				t.Errorf("ParseMountRaw failed for spec '%s', driver '%s', error '%v'", c.bind, c.driver, err.Error())
				continue
			***REMOVED***

			if m.Destination != c.expDest ***REMOVED***
				t.Errorf("Expected destination '%s, was %s', for spec '%s'", c.expDest, m.Destination, c.bind)
			***REMOVED***

			if m.Source != c.expSource ***REMOVED***
				t.Errorf("Expected source '%s', was '%s', for spec '%s'", c.expSource, m.Source, c.bind)
			***REMOVED***

			if m.Name != c.expName ***REMOVED***
				t.Errorf("Expected name '%s', was '%s' for spec '%s'", c.expName, m.Name, c.bind)
			***REMOVED***

			if m.Driver != c.expDriver ***REMOVED***
				t.Errorf("Expected driver '%s', was '%s', for spec '%s'", c.expDriver, m.Driver, c.bind)
			***REMOVED***

			if m.RW != c.expRW ***REMOVED***
				t.Errorf("Expected RW '%v', was '%v' for spec '%s'", c.expRW, m.RW, c.bind)
			***REMOVED***
			if m.Type != c.expType ***REMOVED***
				t.Fatalf("Expected type '%s', was '%s', for spec '%s'", c.expType, m.Type, c.bind)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	tester(linParser, linuxCases)
	tester(winParser, windowsCases)
	tester(lcowParser, lcowCases)
***REMOVED***

func TestParseMountSpec(t *testing.T) ***REMOVED***
	type c struct ***REMOVED***
		input    mount.Mount
		expected MountPoint
	***REMOVED***
	testDir, err := ioutil.TempDir("", "test-mount-config")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(testDir)
	parser := NewParser(runtime.GOOS)
	cases := []c***REMOVED***
		***REMOVED***mount.Mount***REMOVED***Type: mount.TypeBind, Source: testDir, Target: testDestinationPath, ReadOnly: true***REMOVED***, MountPoint***REMOVED***Type: mount.TypeBind, Source: testDir, Destination: testDestinationPath, Propagation: parser.DefaultPropagationMode()***REMOVED******REMOVED***,
		***REMOVED***mount.Mount***REMOVED***Type: mount.TypeBind, Source: testDir, Target: testDestinationPath***REMOVED***, MountPoint***REMOVED***Type: mount.TypeBind, Source: testDir, Destination: testDestinationPath, RW: true, Propagation: parser.DefaultPropagationMode()***REMOVED******REMOVED***,
		***REMOVED***mount.Mount***REMOVED***Type: mount.TypeBind, Source: testDir + string(os.PathSeparator), Target: testDestinationPath, ReadOnly: true***REMOVED***, MountPoint***REMOVED***Type: mount.TypeBind, Source: testDir, Destination: testDestinationPath, Propagation: parser.DefaultPropagationMode()***REMOVED******REMOVED***,
		***REMOVED***mount.Mount***REMOVED***Type: mount.TypeBind, Source: testDir, Target: testDestinationPath + string(os.PathSeparator), ReadOnly: true***REMOVED***, MountPoint***REMOVED***Type: mount.TypeBind, Source: testDir, Destination: testDestinationPath, Propagation: parser.DefaultPropagationMode()***REMOVED******REMOVED***,
		***REMOVED***mount.Mount***REMOVED***Type: mount.TypeVolume, Target: testDestinationPath***REMOVED***, MountPoint***REMOVED***Type: mount.TypeVolume, Destination: testDestinationPath, RW: true, CopyData: parser.DefaultCopyMode()***REMOVED******REMOVED***,
		***REMOVED***mount.Mount***REMOVED***Type: mount.TypeVolume, Target: testDestinationPath + string(os.PathSeparator)***REMOVED***, MountPoint***REMOVED***Type: mount.TypeVolume, Destination: testDestinationPath, RW: true, CopyData: parser.DefaultCopyMode()***REMOVED******REMOVED***,
	***REMOVED***

	for i, c := range cases ***REMOVED***
		t.Logf("case %d", i)
		mp, err := parser.ParseMountSpec(c.input)
		if err != nil ***REMOVED***
			t.Error(err)
		***REMOVED***

		if c.expected.Type != mp.Type ***REMOVED***
			t.Errorf("Expected mount types to match. Expected: '%s', Actual: '%s'", c.expected.Type, mp.Type)
		***REMOVED***
		if c.expected.Destination != mp.Destination ***REMOVED***
			t.Errorf("Expected mount destination to match. Expected: '%s', Actual: '%s'", c.expected.Destination, mp.Destination)
		***REMOVED***
		if c.expected.Source != mp.Source ***REMOVED***
			t.Errorf("Expected mount source to match. Expected: '%s', Actual: '%s'", c.expected.Source, mp.Source)
		***REMOVED***
		if c.expected.RW != mp.RW ***REMOVED***
			t.Errorf("Expected mount writable to match. Expected: '%v', Actual: '%v'", c.expected.RW, mp.RW)
		***REMOVED***
		if c.expected.Propagation != mp.Propagation ***REMOVED***
			t.Errorf("Expected mount propagation to match. Expected: '%v', Actual: '%s'", c.expected.Propagation, mp.Propagation)
		***REMOVED***
		if c.expected.Driver != mp.Driver ***REMOVED***
			t.Errorf("Expected mount driver to match. Expected: '%v', Actual: '%s'", c.expected.Driver, mp.Driver)
		***REMOVED***
		if c.expected.CopyData != mp.CopyData ***REMOVED***
			t.Errorf("Expected mount copy data to match. Expected: '%v', Actual: '%v'", c.expected.CopyData, mp.CopyData)
		***REMOVED***
	***REMOVED***
***REMOVED***
