package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strconv"
)

//-------------------------------------------------------------------------
// config
//
// Structure represents persistent config storage of the gocode daemon. Usually
// the config is located somewhere in ~/.config/gocode directory.
//-------------------------------------------------------------------------

type config struct ***REMOVED***
	ProposeBuiltins    bool   `json:"propose-builtins"`
	LibPath            string `json:"lib-path"`
	CustomPkgPrefix    string `json:"custom-pkg-prefix"`
	CustomVendorDir    string `json:"custom-vendor-dir"`
	Autobuild          bool   `json:"autobuild"`
	ForceDebugOutput   string `json:"force-debug-output"`
	PackageLookupMode  string `json:"package-lookup-mode"`
	CloseTimeout       int    `json:"close-timeout"`
	UnimportedPackages bool   `json:"unimported-packages"`
	Partials           bool   `json:"partials"`
	IgnoreCase         bool   `json:"ignore-case"`
	ClassFiltering     bool   `json:"class-filtering"`
***REMOVED***

var g_config_desc = map[string]string***REMOVED***
	"propose-builtins":    "If set to ***REMOVED***true***REMOVED***, gocode will add built-in types, functions and constants to autocompletion proposals.",
	"lib-path":            "A string option. Allows you to add search paths for packages. By default, gocode only searches ***REMOVED***$GOPATH/pkg/$GOOS_$GOARCH***REMOVED*** and ***REMOVED***$GOROOT/pkg/$GOOS_$GOARCH***REMOVED*** in terms of previously existed environment variables. Also you can specify multiple paths using ':' (colon) as a separator (on Windows use semicolon ';'). The paths specified by ***REMOVED***lib-path***REMOVED*** are prepended to the default ones.",
	"custom-pkg-prefix":   "",
	"custom-vendor-dir":   "",
	"autobuild":           "If set to ***REMOVED***true***REMOVED***, gocode will try to automatically build out-of-date packages when their source files are modified, in order to obtain the freshest autocomplete results for them. This feature is experimental.",
	"force-debug-output":  "If is not empty, gocode will forcefully redirect the logging into that file. Also forces enabling of the debug mode on the server side.",
	"package-lookup-mode": "If set to ***REMOVED***go***REMOVED***, use standard Go package lookup rules. If set to ***REMOVED***gb***REMOVED***, use gb-specific lookup rules. See ***REMOVED***https://github.com/constabulary/gb***REMOVED*** for details.",
	"close-timeout":       "If there have been no completion requests after this number of seconds, the gocode process will terminate. Default is 30 minutes.",
	"unimported-packages": "If set to ***REMOVED***true***REMOVED***, gocode will try to import certain known packages automatically for identifiers which cannot be resolved otherwise. Currently only a limited set of standard library packages is supported.",
	"partials":            "If set to ***REMOVED***false***REMOVED***, gocode will not filter autocompletion results based on entered prefix before the cursor. Instead it will return all available autocompletion results viable for a given context. Whether this option is set to ***REMOVED***true***REMOVED*** or ***REMOVED***false***REMOVED***, gocode will return a valid prefix length for output formats which support it. Setting this option to a non-default value may result in editor misbehaviour.",
	"ignore-case":         "If set to ***REMOVED***true***REMOVED***, gocode will perform case-insensitive matching when doing prefix-based filtering.",
	"class-filtering":     "Enables or disables gocode's feature where it performs class-based filtering if partial input matches corresponding class keyword: const, var, type, func, package.",
***REMOVED***

var g_default_config = config***REMOVED***
	ProposeBuiltins:    false,
	LibPath:            "",
	CustomPkgPrefix:    "",
	Autobuild:          false,
	ForceDebugOutput:   "",
	PackageLookupMode:  "go",
	CloseTimeout:       1800,
	UnimportedPackages: false,
	Partials:           true,
	IgnoreCase:         false,
	ClassFiltering:     true,
***REMOVED***
var g_config = g_default_config

var g_string_to_bool = map[string]bool***REMOVED***
	"t":     true,
	"true":  true,
	"y":     true,
	"yes":   true,
	"on":    true,
	"1":     true,
	"f":     false,
	"false": false,
	"n":     false,
	"no":    false,
	"off":   false,
	"0":     false,
***REMOVED***

func set_value(v reflect.Value, value string) ***REMOVED***
	switch t := v; t.Kind() ***REMOVED***
	case reflect.Bool:
		v, ok := g_string_to_bool[value]
		if ok ***REMOVED***
			t.SetBool(v)
		***REMOVED***
	case reflect.String:
		t.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(value, 10, 64)
		if err == nil ***REMOVED***
			t.SetInt(v)
		***REMOVED***
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(value, 64)
		if err == nil ***REMOVED***
			t.SetFloat(v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func list_value(v reflect.Value, name string, w io.Writer) ***REMOVED***
	switch t := v; t.Kind() ***REMOVED***
	case reflect.Bool:
		fmt.Fprintf(w, "%s %v\n", name, t.Bool())
	case reflect.String:
		fmt.Fprintf(w, "%s \"%v\"\n", name, t.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(w, "%s %v\n", name, t.Int())
	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(w, "%s %v\n", name, t.Float())
	***REMOVED***
***REMOVED***

func (this *config) list() string ***REMOVED***
	str, typ := this.value_and_type()
	buf := bytes.NewBuffer(make([]byte, 0, 256))
	for i := 0; i < str.NumField(); i++ ***REMOVED***
		v := str.Field(i)
		name := typ.Field(i).Tag.Get("json")
		list_value(v, name, buf)
	***REMOVED***
	return buf.String()
***REMOVED***

func (this *config) list_option(name string) string ***REMOVED***
	str, typ := this.value_and_type()
	buf := bytes.NewBuffer(make([]byte, 0, 256))
	for i := 0; i < str.NumField(); i++ ***REMOVED***
		v := str.Field(i)
		nm := typ.Field(i).Tag.Get("json")
		if nm == name ***REMOVED***
			list_value(v, name, buf)
		***REMOVED***
	***REMOVED***
	return buf.String()
***REMOVED***

func (this *config) set_option(name, value string) string ***REMOVED***
	str, typ := this.value_and_type()
	buf := bytes.NewBuffer(make([]byte, 0, 256))
	for i := 0; i < str.NumField(); i++ ***REMOVED***
		v := str.Field(i)
		nm := typ.Field(i).Tag.Get("json")
		if nm == name ***REMOVED***
			set_value(v, value)
			list_value(v, name, buf)
		***REMOVED***
	***REMOVED***
	this.write()
	return buf.String()

***REMOVED***

func (this *config) value_and_type() (reflect.Value, reflect.Type) ***REMOVED***
	v := reflect.ValueOf(this).Elem()
	return v, v.Type()
***REMOVED***

func (this *config) write() error ***REMOVED***
	data, err := json.Marshal(this)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// make sure config dir exists
	dir := config_dir()
	if !file_exists(dir) ***REMOVED***
		os.MkdirAll(dir, 0755)
	***REMOVED***

	f, err := os.Create(config_file())
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer f.Close()

	_, err = f.Write(data)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (this *config) read() error ***REMOVED***
	data, err := ioutil.ReadFile(config_file())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = json.Unmarshal(data, this)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func quoted(v interface***REMOVED******REMOVED***) string ***REMOVED***
	switch v.(type) ***REMOVED***
	case string:
		return fmt.Sprintf("%q", v)
	case int:
		return fmt.Sprint(v)
	case bool:
		return fmt.Sprint(v)
	default:
		panic("unreachable")
	***REMOVED***
***REMOVED***

var descRE = regexp.MustCompile(`***REMOVED***[^***REMOVED***]+***REMOVED***`)

func preprocess_desc(v string) string ***REMOVED***
	return descRE.ReplaceAllStringFunc(v, func(v string) string ***REMOVED***
		return color_cyan + v[1:len(v)-1] + color_none
	***REMOVED***)
***REMOVED***

func (this *config) options() string ***REMOVED***
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%sConfig file location%s: %s\n", color_white_bold, color_none, config_file())
	dv := reflect.ValueOf(g_default_config)
	v, t := this.value_and_type()
	for i, n := 0, t.NumField(); i < n; i++ ***REMOVED***
		f := t.Field(i)
		index := f.Index
		tag := f.Tag.Get("json")
		fmt.Fprintf(&buf, "\n%s%s%s\n", color_yellow_bold, tag, color_none)
		fmt.Fprintf(&buf, "%stype%s: %s\n", color_yellow, color_none, f.Type)
		fmt.Fprintf(&buf, "%svalue%s: %s\n", color_yellow, color_none, quoted(v.FieldByIndex(index).Interface()))
		fmt.Fprintf(&buf, "%sdefault%s: %s\n", color_yellow, color_none, quoted(dv.FieldByIndex(index).Interface()))
		fmt.Fprintf(&buf, "%sdescription%s: %s\n", color_yellow, color_none, preprocess_desc(g_config_desc[tag]))
	***REMOVED***

	return buf.String()
***REMOVED***
