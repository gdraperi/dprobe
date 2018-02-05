package yaml_test

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
)

// An example showing how to unmarshal embedded
// structs from YAML.

type StructA struct ***REMOVED***
	A string `yaml:"a"`
***REMOVED***

type StructB struct ***REMOVED***
	// Embedded structs are not treated as embedded in YAML by default. To do that,
	// add the ",inline" annotation below
	StructA `yaml:",inline"`
	B       string `yaml:"b"`
***REMOVED***

var data = `
a: a string from struct A
b: a string from struct B
`

func ExampleUnmarshal_embedded() ***REMOVED***
	var b StructB

	err := yaml.Unmarshal([]byte(data), &b)
	if err != nil ***REMOVED***
		log.Fatalf("cannot unmarshal data: %v", err)
	***REMOVED***
	fmt.Println(b.A)
	fmt.Println(b.B)
	// Output:
	// a string from struct A
	// a string from struct B
***REMOVED***
