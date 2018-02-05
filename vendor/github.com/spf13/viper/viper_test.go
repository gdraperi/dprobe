// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package viper

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/cast"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

var yamlExample = []byte(`Hacker: true
name: steve
hobbies:
- skateboarding
- snowboarding
- go
clothing:
  jacket: leather
  trousers: denim
  pants:
    size: large
age: 35
eyes : brown
beard: true
`)

var yamlExampleWithExtras = []byte(`Existing: true
Bogus: true
`)

type testUnmarshalExtra struct ***REMOVED***
	Existing bool
***REMOVED***

var tomlExample = []byte(`
title = "TOML Example"

[owner]
organization = "MongoDB"
Bio = "MongoDB Chief Developer Advocate & Hacker at Large"
dob = 1979-05-27T07:32:00Z # First class dates? Why not?`)

var jsonExample = []byte(`***REMOVED***
"id": "0001",
"type": "donut",
"name": "Cake",
"ppu": 0.55,
"batters": ***REMOVED***
        "batter": [
                ***REMOVED*** "type": "Regular" ***REMOVED***,
                ***REMOVED*** "type": "Chocolate" ***REMOVED***,
                ***REMOVED*** "type": "Blueberry" ***REMOVED***,
                ***REMOVED*** "type": "Devil's Food" ***REMOVED***
            ]
***REMOVED***
***REMOVED***`)

var hclExample = []byte(`
id = "0001"
type = "donut"
name = "Cake"
ppu = 0.55
foos ***REMOVED***
	foo ***REMOVED***
		key = 1
	***REMOVED***
	foo ***REMOVED***
		key = 2
	***REMOVED***
	foo ***REMOVED***
		key = 3
	***REMOVED***
	foo ***REMOVED***
		key = 4
	***REMOVED***
***REMOVED***`)

var propertiesExample = []byte(`
p_id: 0001
p_type: donut
p_name: Cake
p_ppu: 0.55
p_batters.batter.type: Regular
`)

var remoteExample = []byte(`***REMOVED***
"id":"0002",
"type":"cronut",
"newkey":"remote"
***REMOVED***`)

func initConfigs() ***REMOVED***
	Reset()
	var r io.Reader
	SetConfigType("yaml")
	r = bytes.NewReader(yamlExample)
	unmarshalReader(r, v.config)

	SetConfigType("json")
	r = bytes.NewReader(jsonExample)
	unmarshalReader(r, v.config)

	SetConfigType("hcl")
	r = bytes.NewReader(hclExample)
	unmarshalReader(r, v.config)

	SetConfigType("properties")
	r = bytes.NewReader(propertiesExample)
	unmarshalReader(r, v.config)

	SetConfigType("toml")
	r = bytes.NewReader(tomlExample)
	unmarshalReader(r, v.config)

	SetConfigType("json")
	remote := bytes.NewReader(remoteExample)
	unmarshalReader(remote, v.kvstore)
***REMOVED***

func initConfig(typ, config string) ***REMOVED***
	Reset()
	SetConfigType(typ)
	r := strings.NewReader(config)

	if err := unmarshalReader(r, v.config); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func initYAML() ***REMOVED***
	initConfig("yaml", string(yamlExample))
***REMOVED***

func initJSON() ***REMOVED***
	Reset()
	SetConfigType("json")
	r := bytes.NewReader(jsonExample)

	unmarshalReader(r, v.config)
***REMOVED***

func initProperties() ***REMOVED***
	Reset()
	SetConfigType("properties")
	r := bytes.NewReader(propertiesExample)

	unmarshalReader(r, v.config)
***REMOVED***

func initTOML() ***REMOVED***
	Reset()
	SetConfigType("toml")
	r := bytes.NewReader(tomlExample)

	unmarshalReader(r, v.config)
***REMOVED***

func initHcl() ***REMOVED***
	Reset()
	SetConfigType("hcl")
	r := bytes.NewReader(hclExample)

	unmarshalReader(r, v.config)
***REMOVED***

// make directories for testing
func initDirs(t *testing.T) (string, string, func()) ***REMOVED***

	var (
		testDirs = []string***REMOVED***`a a`, `b`, `c\c`, `D_`***REMOVED***
		config   = `improbable`
	)

	root, err := ioutil.TempDir("", "")

	cleanup := true
	defer func() ***REMOVED***
		if cleanup ***REMOVED***
			os.Chdir("..")
			os.RemoveAll(root)
		***REMOVED***
	***REMOVED***()

	assert.Nil(t, err)

	err = os.Chdir(root)
	assert.Nil(t, err)

	for _, dir := range testDirs ***REMOVED***
		err = os.Mkdir(dir, 0750)
		assert.Nil(t, err)

		err = ioutil.WriteFile(
			path.Join(dir, config+".toml"),
			[]byte("key = \"value is "+dir+"\"\n"),
			0640)
		assert.Nil(t, err)
	***REMOVED***

	cleanup = false
	return root, config, func() ***REMOVED***
		os.Chdir("..")
		os.RemoveAll(root)
	***REMOVED***
***REMOVED***

//stubs for PFlag Values
type stringValue string

func newStringValue(val string, p *string) *stringValue ***REMOVED***
	*p = val
	return (*stringValue)(p)
***REMOVED***

func (s *stringValue) Set(val string) error ***REMOVED***
	*s = stringValue(val)
	return nil
***REMOVED***

func (s *stringValue) Type() string ***REMOVED***
	return "string"
***REMOVED***

func (s *stringValue) String() string ***REMOVED***
	return fmt.Sprintf("%s", *s)
***REMOVED***

func TestBasics(t *testing.T) ***REMOVED***
	SetConfigFile("/tmp/config.yaml")
	filename, err := v.getConfigFile()
	assert.Equal(t, "/tmp/config.yaml", filename)
	assert.NoError(t, err)
***REMOVED***

func TestDefault(t *testing.T) ***REMOVED***
	SetDefault("age", 45)
	assert.Equal(t, 45, Get("age"))

	SetDefault("clothing.jacket", "slacks")
	assert.Equal(t, "slacks", Get("clothing.jacket"))

	SetConfigType("yaml")
	err := ReadConfig(bytes.NewBuffer(yamlExample))

	assert.NoError(t, err)
	assert.Equal(t, "leather", Get("clothing.jacket"))
***REMOVED***

func TestUnmarshaling(t *testing.T) ***REMOVED***
	SetConfigType("yaml")
	r := bytes.NewReader(yamlExample)

	unmarshalReader(r, v.config)
	assert.True(t, InConfig("name"))
	assert.False(t, InConfig("state"))
	assert.Equal(t, "steve", Get("name"))
	assert.Equal(t, []interface***REMOVED******REMOVED******REMOVED***"skateboarding", "snowboarding", "go"***REMOVED***, Get("hobbies"))
	assert.Equal(t, map[string]interface***REMOVED******REMOVED******REMOVED***"jacket": "leather", "trousers": "denim", "pants": map[string]interface***REMOVED******REMOVED******REMOVED***"size": "large"***REMOVED******REMOVED***, Get("clothing"))
	assert.Equal(t, 35, Get("age"))
***REMOVED***

func TestUnmarshalExact(t *testing.T) ***REMOVED***
	vip := New()
	target := &testUnmarshalExtra***REMOVED******REMOVED***
	vip.SetConfigType("yaml")
	r := bytes.NewReader(yamlExampleWithExtras)
	vip.ReadConfig(r)
	err := vip.UnmarshalExact(target)
	if err == nil ***REMOVED***
		t.Fatal("UnmarshalExact should error when populating a struct from a conf that contains unused fields")
	***REMOVED***
***REMOVED***

func TestOverrides(t *testing.T) ***REMOVED***
	Set("age", 40)
	assert.Equal(t, 40, Get("age"))
***REMOVED***

func TestDefaultPost(t *testing.T) ***REMOVED***
	assert.NotEqual(t, "NYC", Get("state"))
	SetDefault("state", "NYC")
	assert.Equal(t, "NYC", Get("state"))
***REMOVED***

func TestAliases(t *testing.T) ***REMOVED***
	RegisterAlias("years", "age")
	assert.Equal(t, 40, Get("years"))
	Set("years", 45)
	assert.Equal(t, 45, Get("age"))
***REMOVED***

func TestAliasInConfigFile(t *testing.T) ***REMOVED***
	// the config file specifies "beard".  If we make this an alias for
	// "hasbeard", we still want the old config file to work with beard.
	RegisterAlias("beard", "hasbeard")
	assert.Equal(t, true, Get("hasbeard"))
	Set("hasbeard", false)
	assert.Equal(t, false, Get("beard"))
***REMOVED***

func TestYML(t *testing.T) ***REMOVED***
	initYAML()
	assert.Equal(t, "steve", Get("name"))
***REMOVED***

func TestJSON(t *testing.T) ***REMOVED***
	initJSON()
	assert.Equal(t, "0001", Get("id"))
***REMOVED***

func TestProperties(t *testing.T) ***REMOVED***
	initProperties()
	assert.Equal(t, "0001", Get("p_id"))
***REMOVED***

func TestTOML(t *testing.T) ***REMOVED***
	initTOML()
	assert.Equal(t, "TOML Example", Get("title"))
***REMOVED***

func TestHCL(t *testing.T) ***REMOVED***
	initHcl()
	assert.Equal(t, "0001", Get("id"))
	assert.Equal(t, 0.55, Get("ppu"))
	assert.Equal(t, "donut", Get("type"))
	assert.Equal(t, "Cake", Get("name"))
	Set("id", "0002")
	assert.Equal(t, "0002", Get("id"))
	assert.NotEqual(t, "cronut", Get("type"))
***REMOVED***

func TestRemotePrecedence(t *testing.T) ***REMOVED***
	initJSON()

	remote := bytes.NewReader(remoteExample)
	assert.Equal(t, "0001", Get("id"))
	unmarshalReader(remote, v.kvstore)
	assert.Equal(t, "0001", Get("id"))
	assert.NotEqual(t, "cronut", Get("type"))
	assert.Equal(t, "remote", Get("newkey"))
	Set("newkey", "newvalue")
	assert.NotEqual(t, "remote", Get("newkey"))
	assert.Equal(t, "newvalue", Get("newkey"))
	Set("newkey", "remote")
***REMOVED***

func TestEnv(t *testing.T) ***REMOVED***
	initJSON()

	BindEnv("id")
	BindEnv("f", "FOOD")

	os.Setenv("ID", "13")
	os.Setenv("FOOD", "apple")
	os.Setenv("NAME", "crunk")

	assert.Equal(t, "13", Get("id"))
	assert.Equal(t, "apple", Get("f"))
	assert.Equal(t, "Cake", Get("name"))

	AutomaticEnv()

	assert.Equal(t, "crunk", Get("name"))

***REMOVED***

func TestEnvPrefix(t *testing.T) ***REMOVED***
	initJSON()

	SetEnvPrefix("foo") // will be uppercased automatically
	BindEnv("id")
	BindEnv("f", "FOOD") // not using prefix

	os.Setenv("FOO_ID", "13")
	os.Setenv("FOOD", "apple")
	os.Setenv("FOO_NAME", "crunk")

	assert.Equal(t, "13", Get("id"))
	assert.Equal(t, "apple", Get("f"))
	assert.Equal(t, "Cake", Get("name"))

	AutomaticEnv()

	assert.Equal(t, "crunk", Get("name"))
***REMOVED***

func TestAutoEnv(t *testing.T) ***REMOVED***
	Reset()

	AutomaticEnv()
	os.Setenv("FOO_BAR", "13")
	assert.Equal(t, "13", Get("foo_bar"))
***REMOVED***

func TestAutoEnvWithPrefix(t *testing.T) ***REMOVED***
	Reset()

	AutomaticEnv()
	SetEnvPrefix("Baz")
	os.Setenv("BAZ_BAR", "13")
	assert.Equal(t, "13", Get("bar"))
***REMOVED***

func TestSetEnvKeyReplacer(t *testing.T) ***REMOVED***
	Reset()

	AutomaticEnv()
	os.Setenv("REFRESH_INTERVAL", "30s")

	replacer := strings.NewReplacer("-", "_")
	SetEnvKeyReplacer(replacer)

	assert.Equal(t, "30s", Get("refresh-interval"))
***REMOVED***

func TestAllKeys(t *testing.T) ***REMOVED***
	initConfigs()

	ks := sort.StringSlice***REMOVED***"title", "newkey", "owner.organization", "owner.dob", "owner.bio", "name", "beard", "ppu", "batters.batter", "hobbies", "clothing.jacket", "clothing.trousers", "clothing.pants.size", "age", "hacker", "id", "type", "eyes", "p_id", "p_ppu", "p_batters.batter.type", "p_type", "p_name", "foos"***REMOVED***
	dob, _ := time.Parse(time.RFC3339, "1979-05-27T07:32:00Z")
	all := map[string]interface***REMOVED******REMOVED******REMOVED***"owner": map[string]interface***REMOVED******REMOVED******REMOVED***"organization": "MongoDB", "bio": "MongoDB Chief Developer Advocate & Hacker at Large", "dob": dob***REMOVED***, "title": "TOML Example", "ppu": 0.55, "eyes": "brown", "clothing": map[string]interface***REMOVED******REMOVED******REMOVED***"trousers": "denim", "jacket": "leather", "pants": map[string]interface***REMOVED******REMOVED******REMOVED***"size": "large"***REMOVED******REMOVED***, "id": "0001", "batters": map[string]interface***REMOVED******REMOVED******REMOVED***"batter": []interface***REMOVED******REMOVED******REMOVED***map[string]interface***REMOVED******REMOVED******REMOVED***"type": "Regular"***REMOVED***, map[string]interface***REMOVED******REMOVED******REMOVED***"type": "Chocolate"***REMOVED***, map[string]interface***REMOVED******REMOVED******REMOVED***"type": "Blueberry"***REMOVED***, map[string]interface***REMOVED******REMOVED******REMOVED***"type": "Devil's Food"***REMOVED******REMOVED******REMOVED***, "hacker": true, "beard": true, "hobbies": []interface***REMOVED******REMOVED******REMOVED***"skateboarding", "snowboarding", "go"***REMOVED***, "age": 35, "type": "donut", "newkey": "remote", "name": "Cake", "p_id": "0001", "p_ppu": "0.55", "p_name": "Cake", "p_batters": map[string]interface***REMOVED******REMOVED******REMOVED***"batter": map[string]interface***REMOVED******REMOVED******REMOVED***"type": "Regular"***REMOVED******REMOVED***, "p_type": "donut", "foos": []map[string]interface***REMOVED******REMOVED******REMOVED***map[string]interface***REMOVED******REMOVED******REMOVED***"foo": []map[string]interface***REMOVED******REMOVED******REMOVED***map[string]interface***REMOVED******REMOVED******REMOVED***"key": 1***REMOVED***, map[string]interface***REMOVED******REMOVED******REMOVED***"key": 2***REMOVED***, map[string]interface***REMOVED******REMOVED******REMOVED***"key": 3***REMOVED***, map[string]interface***REMOVED******REMOVED******REMOVED***"key": 4***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***

	var allkeys sort.StringSlice
	allkeys = AllKeys()
	allkeys.Sort()
	ks.Sort()

	assert.Equal(t, ks, allkeys)
	assert.Equal(t, all, AllSettings())
***REMOVED***

func TestAllKeysWithEnv(t *testing.T) ***REMOVED***
	v := New()

	// bind and define environment variables (including a nested one)
	v.BindEnv("id")
	v.BindEnv("foo.bar")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	os.Setenv("ID", "13")
	os.Setenv("FOO_BAR", "baz")

	expectedKeys := sort.StringSlice***REMOVED***"id", "foo.bar"***REMOVED***
	expectedKeys.Sort()
	keys := sort.StringSlice(v.AllKeys())
	keys.Sort()
	assert.Equal(t, expectedKeys, keys)
***REMOVED***

func TestAliasesOfAliases(t *testing.T) ***REMOVED***
	Set("Title", "Checking Case")
	RegisterAlias("Foo", "Bar")
	RegisterAlias("Bar", "Title")
	assert.Equal(t, "Checking Case", Get("FOO"))
***REMOVED***

func TestRecursiveAliases(t *testing.T) ***REMOVED***
	RegisterAlias("Baz", "Roo")
	RegisterAlias("Roo", "baz")
***REMOVED***

func TestUnmarshal(t *testing.T) ***REMOVED***
	SetDefault("port", 1313)
	Set("name", "Steve")
	Set("duration", "1s1ms")

	type config struct ***REMOVED***
		Port     int
		Name     string
		Duration time.Duration
	***REMOVED***

	var C config

	err := Unmarshal(&C)
	if err != nil ***REMOVED***
		t.Fatalf("unable to decode into struct, %v", err)
	***REMOVED***

	assert.Equal(t, &config***REMOVED***Name: "Steve", Port: 1313, Duration: time.Second + time.Millisecond***REMOVED***, &C)

	Set("port", 1234)
	err = Unmarshal(&C)
	if err != nil ***REMOVED***
		t.Fatalf("unable to decode into struct, %v", err)
	***REMOVED***
	assert.Equal(t, &config***REMOVED***Name: "Steve", Port: 1234, Duration: time.Second + time.Millisecond***REMOVED***, &C)
***REMOVED***

func TestBindPFlags(t *testing.T) ***REMOVED***
	v := New() // create independent Viper object
	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)

	var testValues = map[string]*string***REMOVED***
		"host":     nil,
		"port":     nil,
		"endpoint": nil,
	***REMOVED***

	var mutatedTestValues = map[string]string***REMOVED***
		"host":     "localhost",
		"port":     "6060",
		"endpoint": "/public",
	***REMOVED***

	for name := range testValues ***REMOVED***
		testValues[name] = flagSet.String(name, "", "test")
	***REMOVED***

	err := v.BindPFlags(flagSet)
	if err != nil ***REMOVED***
		t.Fatalf("error binding flag set, %v", err)
	***REMOVED***

	flagSet.VisitAll(func(flag *pflag.Flag) ***REMOVED***
		flag.Value.Set(mutatedTestValues[flag.Name])
		flag.Changed = true
	***REMOVED***)

	for name, expected := range mutatedTestValues ***REMOVED***
		assert.Equal(t, expected, v.Get(name))
	***REMOVED***

***REMOVED***

func TestBindPFlagsStringSlice(t *testing.T) ***REMOVED***
	for _, testValue := range []struct ***REMOVED***
		Expected []string
		Value    string
	***REMOVED******REMOVED***
		***REMOVED***[]string***REMOVED******REMOVED***, ""***REMOVED***,
		***REMOVED***[]string***REMOVED***"jeden"***REMOVED***, "jeden"***REMOVED***,
		***REMOVED***[]string***REMOVED***"dwa", "trzy"***REMOVED***, "dwa,trzy"***REMOVED***,
		***REMOVED***[]string***REMOVED***"cztery", "piec , szesc"***REMOVED***, "cztery,\"piec , szesc\""***REMOVED******REMOVED*** ***REMOVED***

		for _, changed := range []bool***REMOVED***true, false***REMOVED*** ***REMOVED***
			v := New() // create independent Viper object
			flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
			flagSet.StringSlice("stringslice", testValue.Expected, "test")
			flagSet.Visit(func(f *pflag.Flag) ***REMOVED***
				if len(testValue.Value) > 0 ***REMOVED***
					f.Value.Set(testValue.Value)
					f.Changed = changed
				***REMOVED***
			***REMOVED***)

			err := v.BindPFlags(flagSet)
			if err != nil ***REMOVED***
				t.Fatalf("error binding flag set, %v", err)
			***REMOVED***

			type TestStr struct ***REMOVED***
				StringSlice []string
			***REMOVED***
			val := &TestStr***REMOVED******REMOVED***
			if err := v.Unmarshal(val); err != nil ***REMOVED***
				t.Fatalf("%+#v cannot unmarshal: %s", testValue.Value, err)
			***REMOVED***
			assert.Equal(t, testValue.Expected, val.StringSlice)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestBindPFlag(t *testing.T) ***REMOVED***
	var testString = "testing"
	var testValue = newStringValue(testString, &testString)

	flag := &pflag.Flag***REMOVED***
		Name:    "testflag",
		Value:   testValue,
		Changed: false,
	***REMOVED***

	BindPFlag("testvalue", flag)

	assert.Equal(t, testString, Get("testvalue"))

	flag.Value.Set("testing_mutate")
	flag.Changed = true //hack for pflag usage

	assert.Equal(t, "testing_mutate", Get("testvalue"))

***REMOVED***

func TestBoundCaseSensitivity(t *testing.T) ***REMOVED***
	assert.Equal(t, "brown", Get("eyes"))

	BindEnv("eYEs", "TURTLE_EYES")
	os.Setenv("TURTLE_EYES", "blue")

	assert.Equal(t, "blue", Get("eyes"))

	var testString = "green"
	var testValue = newStringValue(testString, &testString)

	flag := &pflag.Flag***REMOVED***
		Name:    "eyeballs",
		Value:   testValue,
		Changed: true,
	***REMOVED***

	BindPFlag("eYEs", flag)
	assert.Equal(t, "green", Get("eyes"))

***REMOVED***

func TestSizeInBytes(t *testing.T) ***REMOVED***
	input := map[string]uint***REMOVED***
		"":               0,
		"b":              0,
		"12 bytes":       0,
		"200000000000gb": 0,
		"12 b":           12,
		"43 MB":          43 * (1 << 20),
		"10mb":           10 * (1 << 20),
		"1gb":            1 << 30,
	***REMOVED***

	for str, expected := range input ***REMOVED***
		assert.Equal(t, expected, parseSizeInBytes(str), str)
	***REMOVED***
***REMOVED***

func TestFindsNestedKeys(t *testing.T) ***REMOVED***
	initConfigs()
	dob, _ := time.Parse(time.RFC3339, "1979-05-27T07:32:00Z")

	Set("super", map[string]interface***REMOVED******REMOVED******REMOVED***
		"deep": map[string]interface***REMOVED******REMOVED******REMOVED***
			"nested": "value",
		***REMOVED***,
	***REMOVED***)

	expected := map[string]interface***REMOVED******REMOVED******REMOVED***
		"super": map[string]interface***REMOVED******REMOVED******REMOVED***
			"deep": map[string]interface***REMOVED******REMOVED******REMOVED***
				"nested": "value",
			***REMOVED***,
		***REMOVED***,
		"super.deep": map[string]interface***REMOVED******REMOVED******REMOVED***
			"nested": "value",
		***REMOVED***,
		"super.deep.nested":  "value",
		"owner.organization": "MongoDB",
		"batters.batter": []interface***REMOVED******REMOVED******REMOVED***
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"type": "Regular",
			***REMOVED***,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"type": "Chocolate",
			***REMOVED***,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"type": "Blueberry",
			***REMOVED***,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"type": "Devil's Food",
			***REMOVED***,
		***REMOVED***,
		"hobbies": []interface***REMOVED******REMOVED******REMOVED***
			"skateboarding", "snowboarding", "go",
		***REMOVED***,
		"title":  "TOML Example",
		"newkey": "remote",
		"batters": map[string]interface***REMOVED******REMOVED******REMOVED***
			"batter": []interface***REMOVED******REMOVED******REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"type": "Regular",
				***REMOVED***,
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"type": "Chocolate",
				***REMOVED***, map[string]interface***REMOVED******REMOVED******REMOVED***
					"type": "Blueberry",
				***REMOVED***, map[string]interface***REMOVED******REMOVED******REMOVED***
					"type": "Devil's Food",
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"eyes": "brown",
		"age":  35,
		"owner": map[string]interface***REMOVED******REMOVED******REMOVED***
			"organization": "MongoDB",
			"bio":          "MongoDB Chief Developer Advocate & Hacker at Large",
			"dob":          dob,
		***REMOVED***,
		"owner.bio": "MongoDB Chief Developer Advocate & Hacker at Large",
		"type":      "donut",
		"id":        "0001",
		"name":      "Cake",
		"hacker":    true,
		"ppu":       0.55,
		"clothing": map[string]interface***REMOVED******REMOVED******REMOVED***
			"jacket":   "leather",
			"trousers": "denim",
			"pants": map[string]interface***REMOVED******REMOVED******REMOVED***
				"size": "large",
			***REMOVED***,
		***REMOVED***,
		"clothing.jacket":     "leather",
		"clothing.pants.size": "large",
		"clothing.trousers":   "denim",
		"owner.dob":           dob,
		"beard":               true,
		"foos": []map[string]interface***REMOVED******REMOVED******REMOVED***
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"foo": []map[string]interface***REMOVED******REMOVED******REMOVED***
					map[string]interface***REMOVED******REMOVED******REMOVED***
						"key": 1,
					***REMOVED***,
					map[string]interface***REMOVED******REMOVED******REMOVED***
						"key": 2,
					***REMOVED***,
					map[string]interface***REMOVED******REMOVED******REMOVED***
						"key": 3,
					***REMOVED***,
					map[string]interface***REMOVED******REMOVED******REMOVED***
						"key": 4,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for key, expectedValue := range expected ***REMOVED***

		assert.Equal(t, expectedValue, v.Get(key))
	***REMOVED***

***REMOVED***

func TestReadBufConfig(t *testing.T) ***REMOVED***
	v := New()
	v.SetConfigType("yaml")
	v.ReadConfig(bytes.NewBuffer(yamlExample))
	t.Log(v.AllKeys())

	assert.True(t, v.InConfig("name"))
	assert.False(t, v.InConfig("state"))
	assert.Equal(t, "steve", v.Get("name"))
	assert.Equal(t, []interface***REMOVED******REMOVED******REMOVED***"skateboarding", "snowboarding", "go"***REMOVED***, v.Get("hobbies"))
	assert.Equal(t, map[string]interface***REMOVED******REMOVED******REMOVED***"jacket": "leather", "trousers": "denim", "pants": map[string]interface***REMOVED******REMOVED******REMOVED***"size": "large"***REMOVED******REMOVED***, v.Get("clothing"))
	assert.Equal(t, 35, v.Get("age"))
***REMOVED***

func TestIsSet(t *testing.T) ***REMOVED***
	v := New()
	v.SetConfigType("yaml")
	v.ReadConfig(bytes.NewBuffer(yamlExample))
	assert.True(t, v.IsSet("clothing.jacket"))
	assert.False(t, v.IsSet("clothing.jackets"))
	assert.False(t, v.IsSet("helloworld"))
	v.Set("helloworld", "fubar")
	assert.True(t, v.IsSet("helloworld"))
***REMOVED***

func TestDirsSearch(t *testing.T) ***REMOVED***

	root, config, cleanup := initDirs(t)
	defer cleanup()

	v := New()
	v.SetConfigName(config)
	v.SetDefault(`key`, `default`)

	entries, err := ioutil.ReadDir(root)
	for _, e := range entries ***REMOVED***
		if e.IsDir() ***REMOVED***
			v.AddConfigPath(e.Name())
		***REMOVED***
	***REMOVED***

	err = v.ReadInConfig()
	assert.Nil(t, err)

	assert.Equal(t, `value is `+path.Base(v.configPaths[0]), v.GetString(`key`))
***REMOVED***

func TestWrongDirsSearchNotFound(t *testing.T) ***REMOVED***

	_, config, cleanup := initDirs(t)
	defer cleanup()

	v := New()
	v.SetConfigName(config)
	v.SetDefault(`key`, `default`)

	v.AddConfigPath(`whattayoutalkingbout`)
	v.AddConfigPath(`thispathaintthere`)

	err := v.ReadInConfig()
	assert.Equal(t, reflect.TypeOf(ConfigFileNotFoundError***REMOVED***"", ""***REMOVED***), reflect.TypeOf(err))

	// Even though config did not load and the error might have
	// been ignored by the client, the default still loads
	assert.Equal(t, `default`, v.GetString(`key`))
***REMOVED***

func TestWrongDirsSearchNotFoundForMerge(t *testing.T) ***REMOVED***

	_, config, cleanup := initDirs(t)
	defer cleanup()

	v := New()
	v.SetConfigName(config)
	v.SetDefault(`key`, `default`)

	v.AddConfigPath(`whattayoutalkingbout`)
	v.AddConfigPath(`thispathaintthere`)

	err := v.MergeInConfig()
	assert.Equal(t, reflect.TypeOf(ConfigFileNotFoundError***REMOVED***"", ""***REMOVED***), reflect.TypeOf(err))

	// Even though config did not load and the error might have
	// been ignored by the client, the default still loads
	assert.Equal(t, `default`, v.GetString(`key`))
***REMOVED***

func TestSub(t *testing.T) ***REMOVED***
	v := New()
	v.SetConfigType("yaml")
	v.ReadConfig(bytes.NewBuffer(yamlExample))

	subv := v.Sub("clothing")
	assert.Equal(t, v.Get("clothing.pants.size"), subv.Get("pants.size"))

	subv = v.Sub("clothing.pants")
	assert.Equal(t, v.Get("clothing.pants.size"), subv.Get("size"))

	subv = v.Sub("clothing.pants.size")
	assert.Equal(t, (*Viper)(nil), subv)

	subv = v.Sub("missing.key")
	assert.Equal(t, (*Viper)(nil), subv)
***REMOVED***

var hclWriteExpected = []byte(`"foos" = ***REMOVED***
  "foo" = ***REMOVED***
    "key" = 1
  ***REMOVED***

  "foo" = ***REMOVED***
    "key" = 2
  ***REMOVED***

  "foo" = ***REMOVED***
    "key" = 3
  ***REMOVED***

  "foo" = ***REMOVED***
    "key" = 4
  ***REMOVED***
***REMOVED***

"id" = "0001"

"name" = "Cake"

"ppu" = 0.55

"type" = "donut"`)

func TestWriteConfigHCL(t *testing.T) ***REMOVED***
	v := New()
	fs := afero.NewMemMapFs()
	v.SetFs(fs)
	v.SetConfigName("c")
	v.SetConfigType("hcl")
	err := v.ReadConfig(bytes.NewBuffer(hclExample))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := v.WriteConfigAs("c.hcl"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	read, err := afero.ReadFile(fs, "c.hcl")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	assert.Equal(t, hclWriteExpected, read)
***REMOVED***

var jsonWriteExpected = []byte(`***REMOVED***
  "batters": ***REMOVED***
    "batter": [
      ***REMOVED***
        "type": "Regular"
  ***REMOVED***,
      ***REMOVED***
        "type": "Chocolate"
  ***REMOVED***,
      ***REMOVED***
        "type": "Blueberry"
  ***REMOVED***,
      ***REMOVED***
        "type": "Devil's Food"
  ***REMOVED***
    ]
  ***REMOVED***,
  "id": "0001",
  "name": "Cake",
  "ppu": 0.55,
  "type": "donut"
***REMOVED***`)

func TestWriteConfigJson(t *testing.T) ***REMOVED***
	v := New()
	fs := afero.NewMemMapFs()
	v.SetFs(fs)
	v.SetConfigName("c")
	v.SetConfigType("json")
	err := v.ReadConfig(bytes.NewBuffer(jsonExample))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := v.WriteConfigAs("c.json"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	read, err := afero.ReadFile(fs, "c.json")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	assert.Equal(t, jsonWriteExpected, read)
***REMOVED***

var propertiesWriteExpected = []byte(`p_id = 0001
p_type = donut
p_name = Cake
p_ppu = 0.55
p_batters.batter.type = Regular
`)

func TestWriteConfigProperties(t *testing.T) ***REMOVED***
	v := New()
	fs := afero.NewMemMapFs()
	v.SetFs(fs)
	v.SetConfigName("c")
	v.SetConfigType("properties")
	err := v.ReadConfig(bytes.NewBuffer(propertiesExample))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := v.WriteConfigAs("c.properties"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	read, err := afero.ReadFile(fs, "c.properties")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	assert.Equal(t, propertiesWriteExpected, read)
***REMOVED***

func TestWriteConfigTOML(t *testing.T) ***REMOVED***
	fs := afero.NewMemMapFs()
	v := New()
	v.SetFs(fs)
	v.SetConfigName("c")
	v.SetConfigType("toml")
	err := v.ReadConfig(bytes.NewBuffer(tomlExample))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := v.WriteConfigAs("c.toml"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// The TOML String method does not order the contents.
	// Therefore, we must read the generated file and compare the data.
	v2 := New()
	v2.SetFs(fs)
	v2.SetConfigName("c")
	v2.SetConfigType("toml")
	v2.SetConfigFile("c.toml")
	err = v2.ReadInConfig()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	assert.Equal(t, v.GetString("title"), v2.GetString("title"))
	assert.Equal(t, v.GetString("owner.bio"), v2.GetString("owner.bio"))
	assert.Equal(t, v.GetString("owner.dob"), v2.GetString("owner.dob"))
	assert.Equal(t, v.GetString("owner.organization"), v2.GetString("owner.organization"))
***REMOVED***

var yamlWriteExpected = []byte(`age: 35
beard: true
clothing:
  jacket: leather
  pants:
    size: large
  trousers: denim
eyes: brown
hacker: true
hobbies:
- skateboarding
- snowboarding
- go
name: steve
`)

func TestWriteConfigYAML(t *testing.T) ***REMOVED***
	v := New()
	fs := afero.NewMemMapFs()
	v.SetFs(fs)
	v.SetConfigName("c")
	v.SetConfigType("yaml")
	err := v.ReadConfig(bytes.NewBuffer(yamlExample))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := v.WriteConfigAs("c.yaml"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	read, err := afero.ReadFile(fs, "c.yaml")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	assert.Equal(t, yamlWriteExpected, read)
***REMOVED***

var yamlMergeExampleTgt = []byte(`
hello:
    pop: 37890
    lagrenum: 765432101234567
    world:
    - us
    - uk
    - fr
    - de
`)

var yamlMergeExampleSrc = []byte(`
hello:
    pop: 45000
    lagrenum: 7654321001234567
    universe:
    - mw
    - ad
fu: bar
`)

func TestMergeConfig(t *testing.T) ***REMOVED***
	v := New()
	v.SetConfigType("yml")
	if err := v.ReadConfig(bytes.NewBuffer(yamlMergeExampleTgt)); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if pop := v.GetInt("hello.pop"); pop != 37890 ***REMOVED***
		t.Fatalf("pop != 37890, = %d", pop)
	***REMOVED***

	if pop := v.GetInt("hello.lagrenum"); pop != 765432101234567 ***REMOVED***
		t.Fatalf("lagrenum != 765432101234567, = %d", pop)
	***REMOVED***

	if pop := v.GetInt64("hello.lagrenum"); pop != int64(765432101234567) ***REMOVED***
		t.Fatalf("int64 lagrenum != 765432101234567, = %d", pop)
	***REMOVED***

	if world := v.GetStringSlice("hello.world"); len(world) != 4 ***REMOVED***
		t.Fatalf("len(world) != 4, = %d", len(world))
	***REMOVED***

	if fu := v.GetString("fu"); fu != "" ***REMOVED***
		t.Fatalf("fu != \"\", = %s", fu)
	***REMOVED***

	if err := v.MergeConfig(bytes.NewBuffer(yamlMergeExampleSrc)); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if pop := v.GetInt("hello.pop"); pop != 45000 ***REMOVED***
		t.Fatalf("pop != 45000, = %d", pop)
	***REMOVED***

	if pop := v.GetInt("hello.lagrenum"); pop != 7654321001234567 ***REMOVED***
		t.Fatalf("lagrenum != 7654321001234567, = %d", pop)
	***REMOVED***

	if pop := v.GetInt64("hello.lagrenum"); pop != int64(7654321001234567) ***REMOVED***
		t.Fatalf("int64 lagrenum != 7654321001234567, = %d", pop)
	***REMOVED***

	if world := v.GetStringSlice("hello.world"); len(world) != 4 ***REMOVED***
		t.Fatalf("len(world) != 4, = %d", len(world))
	***REMOVED***

	if universe := v.GetStringSlice("hello.universe"); len(universe) != 2 ***REMOVED***
		t.Fatalf("len(universe) != 2, = %d", len(universe))
	***REMOVED***

	if fu := v.GetString("fu"); fu != "bar" ***REMOVED***
		t.Fatalf("fu != \"bar\", = %s", fu)
	***REMOVED***
***REMOVED***

func TestMergeConfigNoMerge(t *testing.T) ***REMOVED***
	v := New()
	v.SetConfigType("yml")
	if err := v.ReadConfig(bytes.NewBuffer(yamlMergeExampleTgt)); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if pop := v.GetInt("hello.pop"); pop != 37890 ***REMOVED***
		t.Fatalf("pop != 37890, = %d", pop)
	***REMOVED***

	if world := v.GetStringSlice("hello.world"); len(world) != 4 ***REMOVED***
		t.Fatalf("len(world) != 4, = %d", len(world))
	***REMOVED***

	if fu := v.GetString("fu"); fu != "" ***REMOVED***
		t.Fatalf("fu != \"\", = %s", fu)
	***REMOVED***

	if err := v.ReadConfig(bytes.NewBuffer(yamlMergeExampleSrc)); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if pop := v.GetInt("hello.pop"); pop != 45000 ***REMOVED***
		t.Fatalf("pop != 45000, = %d", pop)
	***REMOVED***

	if world := v.GetStringSlice("hello.world"); len(world) != 0 ***REMOVED***
		t.Fatalf("len(world) != 0, = %d", len(world))
	***REMOVED***

	if universe := v.GetStringSlice("hello.universe"); len(universe) != 2 ***REMOVED***
		t.Fatalf("len(universe) != 2, = %d", len(universe))
	***REMOVED***

	if fu := v.GetString("fu"); fu != "bar" ***REMOVED***
		t.Fatalf("fu != \"bar\", = %s", fu)
	***REMOVED***
***REMOVED***

func TestUnmarshalingWithAliases(t *testing.T) ***REMOVED***
	v := New()
	v.SetDefault("ID", 1)
	v.Set("name", "Steve")
	v.Set("lastname", "Owen")

	v.RegisterAlias("UserID", "ID")
	v.RegisterAlias("Firstname", "name")
	v.RegisterAlias("Surname", "lastname")

	type config struct ***REMOVED***
		ID        int
		FirstName string
		Surname   string
	***REMOVED***

	var C config
	err := v.Unmarshal(&C)
	if err != nil ***REMOVED***
		t.Fatalf("unable to decode into struct, %v", err)
	***REMOVED***

	assert.Equal(t, &config***REMOVED***ID: 1, FirstName: "Steve", Surname: "Owen"***REMOVED***, &C)
***REMOVED***

func TestSetConfigNameClearsFileCache(t *testing.T) ***REMOVED***
	SetConfigFile("/tmp/config.yaml")
	SetConfigName("default")
	f, err := v.getConfigFile()
	if err == nil ***REMOVED***
		t.Fatalf("config file cache should have been cleared")
	***REMOVED***
	assert.Empty(t, f)
***REMOVED***

func TestShadowedNestedValue(t *testing.T) ***REMOVED***

	config := `name: steve
clothing:
  jacket: leather
  trousers: denim
  pants:
    size: large
`
	initConfig("yaml", config)

	assert.Equal(t, "steve", GetString("name"))

	polyester := "polyester"
	SetDefault("clothing.shirt", polyester)
	SetDefault("clothing.jacket.price", 100)

	assert.Equal(t, "leather", GetString("clothing.jacket"))
	assert.Nil(t, Get("clothing.jacket.price"))
	assert.Equal(t, polyester, GetString("clothing.shirt"))

	clothingSettings := AllSettings()["clothing"].(map[string]interface***REMOVED******REMOVED***)
	assert.Equal(t, "leather", clothingSettings["jacket"])
	assert.Equal(t, polyester, clothingSettings["shirt"])
***REMOVED***

func TestDotParameter(t *testing.T) ***REMOVED***
	initJSON()
	// shoud take precedence over batters defined in jsonExample
	r := bytes.NewReader([]byte(`***REMOVED*** "batters.batter": [ ***REMOVED*** "type": "Small" ***REMOVED*** ] ***REMOVED***`))
	unmarshalReader(r, v.config)

	actual := Get("batters.batter")
	expected := []interface***REMOVED******REMOVED******REMOVED***map[string]interface***REMOVED******REMOVED******REMOVED***"type": "Small"***REMOVED******REMOVED***
	assert.Equal(t, expected, actual)
***REMOVED***

func TestCaseInsensitive(t *testing.T) ***REMOVED***
	for _, config := range []struct ***REMOVED***
		typ     string
		content string
	***REMOVED******REMOVED***
		***REMOVED***"yaml", `
aBcD: 1
eF:
  gH: 2
  iJk: 3
  Lm:
    nO: 4
    P:
      Q: 5
      R: 6
`***REMOVED***,
		***REMOVED***"json", `***REMOVED***
  "aBcD": 1,
  "eF": ***REMOVED***
    "iJk": 3,
    "Lm": ***REMOVED***
      "P": ***REMOVED***
        "Q": 5,
        "R": 6
  ***REMOVED***,
      "nO": 4
***REMOVED***,
    "gH": 2
  ***REMOVED***
***REMOVED***`***REMOVED***,
		***REMOVED***"toml", `aBcD = 1
[eF]
gH = 2
iJk = 3
[eF.Lm]
nO = 4
[eF.Lm.P]
Q = 5
R = 6
`***REMOVED***,
	***REMOVED*** ***REMOVED***
		doTestCaseInsensitive(t, config.typ, config.content)
	***REMOVED***
***REMOVED***

func TestCaseInsensitiveSet(t *testing.T) ***REMOVED***
	Reset()
	m1 := map[string]interface***REMOVED******REMOVED******REMOVED***
		"Foo": 32,
		"Bar": map[interface***REMOVED******REMOVED***]interface ***REMOVED***
		***REMOVED******REMOVED***
			"ABc": "A",
			"cDE": "B"***REMOVED***,
	***REMOVED***

	m2 := map[string]interface***REMOVED******REMOVED******REMOVED***
		"Foo": 52,
		"Bar": map[interface***REMOVED******REMOVED***]interface ***REMOVED***
		***REMOVED******REMOVED***
			"bCd": "A",
			"eFG": "B"***REMOVED***,
	***REMOVED***

	Set("Given1", m1)
	Set("Number1", 42)

	SetDefault("Given2", m2)
	SetDefault("Number2", 52)

	// Verify SetDefault
	if v := Get("number2"); v != 52 ***REMOVED***
		t.Fatalf("Expected 52 got %q", v)
	***REMOVED***

	if v := Get("given2.foo"); v != 52 ***REMOVED***
		t.Fatalf("Expected 52 got %q", v)
	***REMOVED***

	if v := Get("given2.bar.bcd"); v != "A" ***REMOVED***
		t.Fatalf("Expected A got %q", v)
	***REMOVED***

	if _, ok := m2["Foo"]; !ok ***REMOVED***
		t.Fatal("Input map changed")
	***REMOVED***

	// Verify Set
	if v := Get("number1"); v != 42 ***REMOVED***
		t.Fatalf("Expected 42 got %q", v)
	***REMOVED***

	if v := Get("given1.foo"); v != 32 ***REMOVED***
		t.Fatalf("Expected 32 got %q", v)
	***REMOVED***

	if v := Get("given1.bar.abc"); v != "A" ***REMOVED***
		t.Fatalf("Expected A got %q", v)
	***REMOVED***

	if _, ok := m1["Foo"]; !ok ***REMOVED***
		t.Fatal("Input map changed")
	***REMOVED***
***REMOVED***

func TestParseNested(t *testing.T) ***REMOVED***
	type duration struct ***REMOVED***
		Delay time.Duration
	***REMOVED***

	type item struct ***REMOVED***
		Name   string
		Delay  time.Duration
		Nested duration
	***REMOVED***

	config := `[[parent]]
	delay="100ms"
	[parent.nested]
	delay="200ms"
`
	initConfig("toml", config)

	var items []item
	err := v.UnmarshalKey("parent", &items)
	if err != nil ***REMOVED***
		t.Fatalf("unable to decode into struct, %v", err)
	***REMOVED***

	assert.Equal(t, 1, len(items))
	assert.Equal(t, 100*time.Millisecond, items[0].Delay)
	assert.Equal(t, 200*time.Millisecond, items[0].Nested.Delay)
***REMOVED***

func doTestCaseInsensitive(t *testing.T, typ, config string) ***REMOVED***
	initConfig(typ, config)
	Set("RfD", true)
	assert.Equal(t, true, Get("rfd"))
	assert.Equal(t, true, Get("rFD"))
	assert.Equal(t, 1, cast.ToInt(Get("abcd")))
	assert.Equal(t, 1, cast.ToInt(Get("Abcd")))
	assert.Equal(t, 2, cast.ToInt(Get("ef.gh")))
	assert.Equal(t, 3, cast.ToInt(Get("ef.ijk")))
	assert.Equal(t, 4, cast.ToInt(Get("ef.lm.no")))
	assert.Equal(t, 5, cast.ToInt(Get("ef.lm.p.q")))

***REMOVED***

func BenchmarkGetBool(b *testing.B) ***REMOVED***
	key := "BenchmarkGetBool"
	v = New()
	v.Set(key, true)

	for i := 0; i < b.N; i++ ***REMOVED***
		if !v.GetBool(key) ***REMOVED***
			b.Fatal("GetBool returned false")
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkGet(b *testing.B) ***REMOVED***
	key := "BenchmarkGet"
	v = New()
	v.Set(key, true)

	for i := 0; i < b.N; i++ ***REMOVED***
		if !v.Get(key).(bool) ***REMOVED***
			b.Fatal("Get returned false")
		***REMOVED***
	***REMOVED***
***REMOVED***

// This is the "perfect result" for the above.
func BenchmarkGetBoolFromMap(b *testing.B) ***REMOVED***
	m := make(map[string]bool)
	key := "BenchmarkGetBool"
	m[key] = true

	for i := 0; i < b.N; i++ ***REMOVED***
		if !m[key] ***REMOVED***
			b.Fatal("Map value was false")
		***REMOVED***
	***REMOVED***
***REMOVED***
