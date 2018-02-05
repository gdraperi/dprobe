package toml

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	burntsushi "github.com/BurntSushi/toml"
	yaml "gopkg.in/yaml.v2"
)

type benchmarkDoc struct ***REMOVED***
	Table struct ***REMOVED***
		Key      string
		Subtable struct ***REMOVED***
			Key string
		***REMOVED***
		Inline struct ***REMOVED***
			Name struct ***REMOVED***
				First string
				Last  string
			***REMOVED***
			Point struct ***REMOVED***
				X int64
				U int64
			***REMOVED***
		***REMOVED***
	***REMOVED***
	String struct ***REMOVED***
		Basic struct ***REMOVED***
			Basic string
		***REMOVED***
		Multiline struct ***REMOVED***
			Key1      string
			Key2      string
			Key3      string
			Continued struct ***REMOVED***
				Key1 string
				Key2 string
				Key3 string
			***REMOVED***
		***REMOVED***
		Literal struct ***REMOVED***
			Winpath   string
			Winpath2  string
			Quoted    string
			Regex     string
			Multiline struct ***REMOVED***
				Regex2 string
				Lines  string
			***REMOVED***
		***REMOVED***
	***REMOVED***
	Integer struct ***REMOVED***
		Key1        int64
		Key2        int64
		Key3        int64
		Key4        int64
		Underscores struct ***REMOVED***
			Key1 int64
			Key2 int64
			Key3 int64
		***REMOVED***
	***REMOVED***
	Float struct ***REMOVED***
		Fractional struct ***REMOVED***
			Key1 float64
			Key2 float64
			Key3 float64
		***REMOVED***
		Exponent struct ***REMOVED***
			Key1 float64
			Key2 float64
			Key3 float64
		***REMOVED***
		Both struct ***REMOVED***
			Key float64
		***REMOVED***
		Underscores struct ***REMOVED***
			Key1 float64
			Key2 float64
		***REMOVED***
	***REMOVED***
	Boolean struct ***REMOVED***
		True  bool
		False bool
	***REMOVED***
	Datetime struct ***REMOVED***
		Key1 time.Time
		Key2 time.Time
		Key3 time.Time
	***REMOVED***
	Array struct ***REMOVED***
		Key1 []int64
		Key2 []string
		Key3 [][]int64
		// TODO: Key4 not supported by go-toml's Unmarshal
		Key5 []int64
		Key6 []int64
	***REMOVED***
	Products []struct ***REMOVED***
		Name  string
		Sku   int64
		Color string
	***REMOVED***
	Fruit []struct ***REMOVED***
		Name     string
		Physical struct ***REMOVED***
			Color   string
			Shape   string
			Variety []struct ***REMOVED***
				Name string
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkParseToml(b *testing.B) ***REMOVED***
	fileBytes, err := ioutil.ReadFile("benchmark.toml")
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		_, err := LoadReader(bytes.NewReader(fileBytes))
		if err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkUnmarshalToml(b *testing.B) ***REMOVED***
	bytes, err := ioutil.ReadFile("benchmark.toml")
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		target := benchmarkDoc***REMOVED******REMOVED***
		err := Unmarshal(bytes, &target)
		if err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkUnmarshalBurntSushiToml(b *testing.B) ***REMOVED***
	bytes, err := ioutil.ReadFile("benchmark.toml")
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		target := benchmarkDoc***REMOVED******REMOVED***
		err := burntsushi.Unmarshal(bytes, &target)
		if err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkUnmarshalJson(b *testing.B) ***REMOVED***
	bytes, err := ioutil.ReadFile("benchmark.json")
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		target := benchmarkDoc***REMOVED******REMOVED***
		err := json.Unmarshal(bytes, &target)
		if err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkUnmarshalYaml(b *testing.B) ***REMOVED***
	bytes, err := ioutil.ReadFile("benchmark.yml")
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		target := benchmarkDoc***REMOVED******REMOVED***
		err := yaml.Unmarshal(bytes, &target)
		if err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***
