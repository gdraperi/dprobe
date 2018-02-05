package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/pelletier/go-toml"
)

func main() ***REMOVED***
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil ***REMOVED***
		log.Fatalf("Error during TOML read: %s", err)
		os.Exit(2)
	***REMOVED***
	tree, err := toml.Load(string(bytes))
	if err != nil ***REMOVED***
		log.Fatalf("Error during TOML load: %s", err)
		os.Exit(1)
	***REMOVED***

	typedTree := translate(*tree)

	if err := json.NewEncoder(os.Stdout).Encode(typedTree); err != nil ***REMOVED***
		log.Fatalf("Error encoding JSON: %s", err)
		os.Exit(3)
	***REMOVED***

	os.Exit(0)
***REMOVED***

func translate(tomlData interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	switch orig := tomlData.(type) ***REMOVED***
	case map[string]interface***REMOVED******REMOVED***:
		typed := make(map[string]interface***REMOVED******REMOVED***, len(orig))
		for k, v := range orig ***REMOVED***
			typed[k] = translate(v)
		***REMOVED***
		return typed
	case *toml.Tree:
		return translate(*orig)
	case toml.Tree:
		keys := orig.Keys()
		typed := make(map[string]interface***REMOVED******REMOVED***, len(keys))
		for _, k := range keys ***REMOVED***
			typed[k] = translate(orig.GetPath([]string***REMOVED***k***REMOVED***))
		***REMOVED***
		return typed
	case []*toml.Tree:
		typed := make([]map[string]interface***REMOVED******REMOVED***, len(orig))
		for i, v := range orig ***REMOVED***
			typed[i] = translate(v).(map[string]interface***REMOVED******REMOVED***)
		***REMOVED***
		return typed
	case []map[string]interface***REMOVED******REMOVED***:
		typed := make([]map[string]interface***REMOVED******REMOVED***, len(orig))
		for i, v := range orig ***REMOVED***
			typed[i] = translate(v).(map[string]interface***REMOVED******REMOVED***)
		***REMOVED***
		return typed
	case []interface***REMOVED******REMOVED***:
		typed := make([]interface***REMOVED******REMOVED***, len(orig))
		for i, v := range orig ***REMOVED***
			typed[i] = translate(v)
		***REMOVED***
		return tag("array", typed)
	case time.Time:
		return tag("datetime", orig.Format("2006-01-02T15:04:05Z"))
	case bool:
		return tag("bool", fmt.Sprintf("%v", orig))
	case int64:
		return tag("integer", fmt.Sprintf("%d", orig))
	case float64:
		return tag("float", fmt.Sprintf("%v", orig))
	case string:
		return tag("string", orig)
	***REMOVED***

	panic(fmt.Sprintf("Unknown type: %T", tomlData))
***REMOVED***

func tag(typeName string, data interface***REMOVED******REMOVED***) map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	return map[string]interface***REMOVED******REMOVED******REMOVED***
		"type":  typeName,
		"value": data,
	***REMOVED***
***REMOVED***
