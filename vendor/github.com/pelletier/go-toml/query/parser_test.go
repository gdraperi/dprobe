package query

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"sort"
	"strings"
	"testing"
	"time"
)

type queryTestNode struct ***REMOVED***
	value    interface***REMOVED******REMOVED***
	position toml.Position
***REMOVED***

func valueString(root interface***REMOVED******REMOVED***) string ***REMOVED***
	result := "" //fmt.Sprintf("%T:", root)
	switch node := root.(type) ***REMOVED***
	case *Result:
		items := []string***REMOVED******REMOVED***
		for i, v := range node.Values() ***REMOVED***
			items = append(items, fmt.Sprintf("%s:%s",
				node.Positions()[i].String(), valueString(v)))
		***REMOVED***
		sort.Strings(items)
		result = "[" + strings.Join(items, ", ") + "]"
	case queryTestNode:
		result = fmt.Sprintf("%s:%s",
			node.position.String(), valueString(node.value))
	case []interface***REMOVED******REMOVED***:
		items := []string***REMOVED******REMOVED***
		for _, v := range node ***REMOVED***
			items = append(items, valueString(v))
		***REMOVED***
		sort.Strings(items)
		result = "[" + strings.Join(items, ", ") + "]"
	case *toml.Tree:
		// workaround for unreliable map key ordering
		items := []string***REMOVED******REMOVED***
		for _, k := range node.Keys() ***REMOVED***
			v := node.GetPath([]string***REMOVED***k***REMOVED***)
			items = append(items, k+":"+valueString(v))
		***REMOVED***
		sort.Strings(items)
		result = "***REMOVED***" + strings.Join(items, ", ") + "***REMOVED***"
	case map[string]interface***REMOVED******REMOVED***:
		// workaround for unreliable map key ordering
		items := []string***REMOVED******REMOVED***
		for k, v := range node ***REMOVED***
			items = append(items, k+":"+valueString(v))
		***REMOVED***
		sort.Strings(items)
		result = "***REMOVED***" + strings.Join(items, ", ") + "***REMOVED***"
	case int64:
		result += fmt.Sprintf("%d", node)
	case string:
		result += "'" + node + "'"
	case float64:
		result += fmt.Sprintf("%f", node)
	case bool:
		result += fmt.Sprintf("%t", node)
	case time.Time:
		result += fmt.Sprintf("'%v'", node)
	***REMOVED***
	return result
***REMOVED***

func assertValue(t *testing.T, result, ref interface***REMOVED******REMOVED***) ***REMOVED***
	pathStr := valueString(result)
	refStr := valueString(ref)
	if pathStr != refStr ***REMOVED***
		t.Errorf("values do not match")
		t.Log("test:", pathStr)
		t.Log("ref: ", refStr)
	***REMOVED***
***REMOVED***

func assertQueryPositions(t *testing.T, tomlDoc string, query string, ref []interface***REMOVED******REMOVED***) ***REMOVED***
	tree, err := toml.Load(tomlDoc)
	if err != nil ***REMOVED***
		t.Errorf("Non-nil toml parse error: %v", err)
		return
	***REMOVED***
	q, err := Compile(query)
	if err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***
	results := q.Execute(tree)
	assertValue(t, results, ref)
***REMOVED***

func TestQueryRoot(t *testing.T) ***REMOVED***
	assertQueryPositions(t,
		"a = 42",
		"$",
		[]interface***REMOVED******REMOVED******REMOVED***
			queryTestNode***REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"a": int64(42),
				***REMOVED***, toml.Position***REMOVED***1, 1***REMOVED***,
			***REMOVED***,
		***REMOVED***)
***REMOVED***

func TestQueryKey(t *testing.T) ***REMOVED***
	assertQueryPositions(t,
		"[foo]\na = 42",
		"$.foo.a",
		[]interface***REMOVED******REMOVED******REMOVED***
			queryTestNode***REMOVED***
				int64(42), toml.Position***REMOVED***2, 1***REMOVED***,
			***REMOVED***,
		***REMOVED***)
***REMOVED***

func TestQueryKeyString(t *testing.T) ***REMOVED***
	assertQueryPositions(t,
		"[foo]\na = 42",
		"$.foo['a']",
		[]interface***REMOVED******REMOVED******REMOVED***
			queryTestNode***REMOVED***
				int64(42), toml.Position***REMOVED***2, 1***REMOVED***,
			***REMOVED***,
		***REMOVED***)
***REMOVED***

func TestQueryIndex(t *testing.T) ***REMOVED***
	assertQueryPositions(t,
		"[foo]\na = [1,2,3,4,5,6,7,8,9,0]",
		"$.foo.a[5]",
		[]interface***REMOVED******REMOVED******REMOVED***
			queryTestNode***REMOVED***
				int64(6), toml.Position***REMOVED***2, 1***REMOVED***,
			***REMOVED***,
		***REMOVED***)
***REMOVED***

func TestQuerySliceRange(t *testing.T) ***REMOVED***
	assertQueryPositions(t,
		"[foo]\na = [1,2,3,4,5,6,7,8,9,0]",
		"$.foo.a[0:5]",
		[]interface***REMOVED******REMOVED******REMOVED***
			queryTestNode***REMOVED***
				int64(1), toml.Position***REMOVED***2, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				int64(2), toml.Position***REMOVED***2, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				int64(3), toml.Position***REMOVED***2, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				int64(4), toml.Position***REMOVED***2, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				int64(5), toml.Position***REMOVED***2, 1***REMOVED***,
			***REMOVED***,
		***REMOVED***)
***REMOVED***

func TestQuerySliceStep(t *testing.T) ***REMOVED***
	assertQueryPositions(t,
		"[foo]\na = [1,2,3,4,5,6,7,8,9,0]",
		"$.foo.a[0:5:2]",
		[]interface***REMOVED******REMOVED******REMOVED***
			queryTestNode***REMOVED***
				int64(1), toml.Position***REMOVED***2, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				int64(3), toml.Position***REMOVED***2, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				int64(5), toml.Position***REMOVED***2, 1***REMOVED***,
			***REMOVED***,
		***REMOVED***)
***REMOVED***

func TestQueryAny(t *testing.T) ***REMOVED***
	assertQueryPositions(t,
		"[foo.bar]\na=1\nb=2\n[foo.baz]\na=3\nb=4",
		"$.foo.*",
		[]interface***REMOVED******REMOVED******REMOVED***
			queryTestNode***REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"a": int64(1),
					"b": int64(2),
				***REMOVED***, toml.Position***REMOVED***1, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"a": int64(3),
					"b": int64(4),
				***REMOVED***, toml.Position***REMOVED***4, 1***REMOVED***,
			***REMOVED***,
		***REMOVED***)
***REMOVED***
func TestQueryUnionSimple(t *testing.T) ***REMOVED***
	assertQueryPositions(t,
		"[foo.bar]\na=1\nb=2\n[baz.foo]\na=3\nb=4\n[gorf.foo]\na=5\nb=6",
		"$.*[bar,foo]",
		[]interface***REMOVED******REMOVED******REMOVED***
			queryTestNode***REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"a": int64(1),
					"b": int64(2),
				***REMOVED***, toml.Position***REMOVED***1, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"a": int64(3),
					"b": int64(4),
				***REMOVED***, toml.Position***REMOVED***4, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"a": int64(5),
					"b": int64(6),
				***REMOVED***, toml.Position***REMOVED***7, 1***REMOVED***,
			***REMOVED***,
		***REMOVED***)
***REMOVED***

func TestQueryRecursionAll(t *testing.T) ***REMOVED***
	assertQueryPositions(t,
		"[foo.bar]\na=1\nb=2\n[baz.foo]\na=3\nb=4\n[gorf.foo]\na=5\nb=6",
		"$..*",
		[]interface***REMOVED******REMOVED******REMOVED***
			queryTestNode***REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"foo": map[string]interface***REMOVED******REMOVED******REMOVED***
						"bar": map[string]interface***REMOVED******REMOVED******REMOVED***
							"a": int64(1),
							"b": int64(2),
						***REMOVED***,
					***REMOVED***,
					"baz": map[string]interface***REMOVED******REMOVED******REMOVED***
						"foo": map[string]interface***REMOVED******REMOVED******REMOVED***
							"a": int64(3),
							"b": int64(4),
						***REMOVED***,
					***REMOVED***,
					"gorf": map[string]interface***REMOVED******REMOVED******REMOVED***
						"foo": map[string]interface***REMOVED******REMOVED******REMOVED***
							"a": int64(5),
							"b": int64(6),
						***REMOVED***,
					***REMOVED***,
				***REMOVED***, toml.Position***REMOVED***1, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"bar": map[string]interface***REMOVED******REMOVED******REMOVED***
						"a": int64(1),
						"b": int64(2),
					***REMOVED***,
				***REMOVED***, toml.Position***REMOVED***1, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"a": int64(1),
					"b": int64(2),
				***REMOVED***, toml.Position***REMOVED***1, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				int64(1), toml.Position***REMOVED***2, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				int64(2), toml.Position***REMOVED***3, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"foo": map[string]interface***REMOVED******REMOVED******REMOVED***
						"a": int64(3),
						"b": int64(4),
					***REMOVED***,
				***REMOVED***, toml.Position***REMOVED***4, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"a": int64(3),
					"b": int64(4),
				***REMOVED***, toml.Position***REMOVED***4, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				int64(3), toml.Position***REMOVED***5, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				int64(4), toml.Position***REMOVED***6, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"foo": map[string]interface***REMOVED******REMOVED******REMOVED***
						"a": int64(5),
						"b": int64(6),
					***REMOVED***,
				***REMOVED***, toml.Position***REMOVED***7, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"a": int64(5),
					"b": int64(6),
				***REMOVED***, toml.Position***REMOVED***7, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				int64(5), toml.Position***REMOVED***8, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				int64(6), toml.Position***REMOVED***9, 1***REMOVED***,
			***REMOVED***,
		***REMOVED***)
***REMOVED***

func TestQueryRecursionUnionSimple(t *testing.T) ***REMOVED***
	assertQueryPositions(t,
		"[foo.bar]\na=1\nb=2\n[baz.foo]\na=3\nb=4\n[gorf.foo]\na=5\nb=6",
		"$..['foo','bar']",
		[]interface***REMOVED******REMOVED******REMOVED***
			queryTestNode***REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"bar": map[string]interface***REMOVED******REMOVED******REMOVED***
						"a": int64(1),
						"b": int64(2),
					***REMOVED***,
				***REMOVED***, toml.Position***REMOVED***1, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"a": int64(3),
					"b": int64(4),
				***REMOVED***, toml.Position***REMOVED***4, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"a": int64(1),
					"b": int64(2),
				***REMOVED***, toml.Position***REMOVED***1, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"a": int64(5),
					"b": int64(6),
				***REMOVED***, toml.Position***REMOVED***7, 1***REMOVED***,
			***REMOVED***,
		***REMOVED***)
***REMOVED***

func TestQueryFilterFn(t *testing.T) ***REMOVED***
	buff, err := ioutil.ReadFile("../example.toml")
	if err != nil ***REMOVED***
		t.Error(err)
		return
	***REMOVED***

	assertQueryPositions(t, string(buff),
		"$..[?(int)]",
		[]interface***REMOVED******REMOVED******REMOVED***
			queryTestNode***REMOVED***
				int64(8001), toml.Position***REMOVED***13, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				int64(8001), toml.Position***REMOVED***13, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				int64(8002), toml.Position***REMOVED***13, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				int64(5000), toml.Position***REMOVED***14, 1***REMOVED***,
			***REMOVED***,
		***REMOVED***)

	assertQueryPositions(t, string(buff),
		"$..[?(string)]",
		[]interface***REMOVED******REMOVED******REMOVED***
			queryTestNode***REMOVED***
				"TOML Example", toml.Position***REMOVED***3, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				"Tom Preston-Werner", toml.Position***REMOVED***6, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				"GitHub", toml.Position***REMOVED***7, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				"GitHub Cofounder & CEO\nLikes tater tots and beer.",
				toml.Position***REMOVED***8, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				"192.168.1.1", toml.Position***REMOVED***12, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				"10.0.0.1", toml.Position***REMOVED***21, 3***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				"eqdc10", toml.Position***REMOVED***22, 3***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				"10.0.0.2", toml.Position***REMOVED***25, 3***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				"eqdc10", toml.Position***REMOVED***26, 3***REMOVED***,
			***REMOVED***,
		***REMOVED***)

	assertQueryPositions(t, string(buff),
		"$..[?(float)]",
		[]interface***REMOVED******REMOVED******REMOVED***
		// no float values in document
		***REMOVED***)

	tv, _ := time.Parse(time.RFC3339, "1979-05-27T07:32:00Z")
	assertQueryPositions(t, string(buff),
		"$..[?(tree)]",
		[]interface***REMOVED******REMOVED******REMOVED***
			queryTestNode***REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"name":         "Tom Preston-Werner",
					"organization": "GitHub",
					"bio":          "GitHub Cofounder & CEO\nLikes tater tots and beer.",
					"dob":          tv,
				***REMOVED***, toml.Position***REMOVED***5, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"server":         "192.168.1.1",
					"ports":          []interface***REMOVED******REMOVED******REMOVED***int64(8001), int64(8001), int64(8002)***REMOVED***,
					"connection_max": int64(5000),
					"enabled":        true,
				***REMOVED***, toml.Position***REMOVED***11, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"alpha": map[string]interface***REMOVED******REMOVED******REMOVED***
						"ip": "10.0.0.1",
						"dc": "eqdc10",
					***REMOVED***,
					"beta": map[string]interface***REMOVED******REMOVED******REMOVED***
						"ip": "10.0.0.2",
						"dc": "eqdc10",
					***REMOVED***,
				***REMOVED***, toml.Position***REMOVED***17, 1***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"ip": "10.0.0.1",
					"dc": "eqdc10",
				***REMOVED***, toml.Position***REMOVED***20, 3***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"ip": "10.0.0.2",
					"dc": "eqdc10",
				***REMOVED***, toml.Position***REMOVED***24, 3***REMOVED***,
			***REMOVED***,
			queryTestNode***REMOVED***
				map[string]interface***REMOVED******REMOVED******REMOVED***
					"data": []interface***REMOVED******REMOVED******REMOVED***
						[]interface***REMOVED******REMOVED******REMOVED***"gamma", "delta"***REMOVED***,
						[]interface***REMOVED******REMOVED******REMOVED***int64(1), int64(2)***REMOVED***,
					***REMOVED***,
				***REMOVED***, toml.Position***REMOVED***28, 1***REMOVED***,
			***REMOVED***,
		***REMOVED***)

	assertQueryPositions(t, string(buff),
		"$..[?(time)]",
		[]interface***REMOVED******REMOVED******REMOVED***
			queryTestNode***REMOVED***
				tv, toml.Position***REMOVED***9, 1***REMOVED***,
			***REMOVED***,
		***REMOVED***)

	assertQueryPositions(t, string(buff),
		"$..[?(bool)]",
		[]interface***REMOVED******REMOVED******REMOVED***
			queryTestNode***REMOVED***
				true, toml.Position***REMOVED***15, 1***REMOVED***,
			***REMOVED***,
		***REMOVED***)
***REMOVED***
