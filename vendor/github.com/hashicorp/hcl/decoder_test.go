package hcl

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/hcl/hcl/ast"
)

func TestDecode_interface(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		File string
		Err  bool
		Out  interface***REMOVED******REMOVED***
	***REMOVED******REMOVED***
		***REMOVED***
			"basic.hcl",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"foo": "bar",
				"bar": "$***REMOVED***file(\"bing/bong.txt\")***REMOVED***",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"basic_squish.hcl",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"foo":     "bar",
				"bar":     "$***REMOVED***file(\"bing/bong.txt\")***REMOVED***",
				"foo-bar": "baz",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"empty.hcl",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"resource": []map[string]interface***REMOVED******REMOVED******REMOVED***
					map[string]interface***REMOVED******REMOVED******REMOVED***
						"foo": []map[string]interface***REMOVED******REMOVED******REMOVED***
							map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"tfvars.hcl",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"regularvar": "Should work",
				"map.key1":   "Value",
				"map.key2":   "Other value",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"escape.hcl",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"foo":          "bar\"baz\\n",
				"qux":          "back\\slash",
				"bar":          "new\nline",
				"qax":          `slash\:colon`,
				"nested":       `$***REMOVED***HH\\:mm\\:ss***REMOVED***`,
				"nestedquotes": `$***REMOVED***"\"stringwrappedinquotes\""***REMOVED***`,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"float.hcl",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"a": 1.02,
				"b": 2,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"multiline_bad.hcl",
			true,
			nil,
		***REMOVED***,
		***REMOVED***
			"multiline_literal.hcl",
			true,
			nil,
		***REMOVED***,
		***REMOVED***
			"multiline_literal_with_hil.hcl",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***"multiline_literal_with_hil": "$***REMOVED***hello\n  world***REMOVED***"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"multiline_no_marker.hcl",
			true,
			nil,
		***REMOVED***,
		***REMOVED***
			"multiline.hcl",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***"foo": "bar\nbaz\n"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"multiline_indented.hcl",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***"foo": "  bar\n  baz\n"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"multiline_no_hanging_indent.hcl",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***"foo": "  baz\n    bar\n      foo\n"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"multiline_no_eof.hcl",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***"foo": "bar\nbaz\n", "key": "value"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"multiline.json",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***"foo": "bar\nbaz"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"null_strings.json",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"module": []map[string]interface***REMOVED******REMOVED******REMOVED***
					map[string]interface***REMOVED******REMOVED******REMOVED***
						"app": []map[string]interface***REMOVED******REMOVED******REMOVED***
							map[string]interface***REMOVED******REMOVED******REMOVED***"foo": ""***REMOVED***,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"scientific.json",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"a": 1e-10,
				"b": 1e+10,
				"c": 1e10,
				"d": 1.2e-10,
				"e": 1.2e+10,
				"f": 1.2e10,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"scientific.hcl",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"a": 1e-10,
				"b": 1e+10,
				"c": 1e10,
				"d": 1.2e-10,
				"e": 1.2e+10,
				"f": 1.2e10,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"terraform_heroku.hcl",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"name": "terraform-test-app",
				"config_vars": []map[string]interface***REMOVED******REMOVED******REMOVED***
					map[string]interface***REMOVED******REMOVED******REMOVED***
						"FOO": "bar",
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"structure_multi.hcl",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"foo": []map[string]interface***REMOVED******REMOVED******REMOVED***
					map[string]interface***REMOVED******REMOVED******REMOVED***
						"baz": []map[string]interface***REMOVED******REMOVED******REMOVED***
							map[string]interface***REMOVED******REMOVED******REMOVED***"key": 7***REMOVED***,
						***REMOVED***,
					***REMOVED***,
					map[string]interface***REMOVED******REMOVED******REMOVED***
						"bar": []map[string]interface***REMOVED******REMOVED******REMOVED***
							map[string]interface***REMOVED******REMOVED******REMOVED***"key": 12***REMOVED***,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"structure_multi.json",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"foo": []map[string]interface***REMOVED******REMOVED******REMOVED***
					map[string]interface***REMOVED******REMOVED******REMOVED***
						"baz": []map[string]interface***REMOVED******REMOVED******REMOVED***
							map[string]interface***REMOVED******REMOVED******REMOVED***"key": 7***REMOVED***,
						***REMOVED***,
					***REMOVED***,
					map[string]interface***REMOVED******REMOVED******REMOVED***
						"bar": []map[string]interface***REMOVED******REMOVED******REMOVED***
							map[string]interface***REMOVED******REMOVED******REMOVED***"key": 12***REMOVED***,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"list_of_lists.hcl",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"foo": []interface***REMOVED******REMOVED******REMOVED***
					[]interface***REMOVED******REMOVED******REMOVED***"foo"***REMOVED***,
					[]interface***REMOVED******REMOVED******REMOVED***"bar"***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"list_of_maps.hcl",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"foo": []interface***REMOVED******REMOVED******REMOVED***
					map[string]interface***REMOVED******REMOVED******REMOVED***"somekey1": "someval1"***REMOVED***,
					map[string]interface***REMOVED******REMOVED******REMOVED***"somekey2": "someval2", "someextrakey": "someextraval"***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"assign_deep.hcl",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"resource": []interface***REMOVED******REMOVED******REMOVED***
					map[string]interface***REMOVED******REMOVED******REMOVED***
						"foo": []interface***REMOVED******REMOVED******REMOVED***
							map[string]interface***REMOVED******REMOVED******REMOVED***
								"bar": []map[string]interface***REMOVED******REMOVED******REMOVED***
									map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED******REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***,
		***REMOVED***,
		***REMOVED***
			"structure_list.hcl",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"foo": []map[string]interface***REMOVED******REMOVED******REMOVED***
					map[string]interface***REMOVED******REMOVED******REMOVED***
						"key": 7,
					***REMOVED***,
					map[string]interface***REMOVED******REMOVED******REMOVED***
						"key": 12,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"structure_list.json",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"foo": []map[string]interface***REMOVED******REMOVED******REMOVED***
					map[string]interface***REMOVED******REMOVED******REMOVED***
						"key": 7,
					***REMOVED***,
					map[string]interface***REMOVED******REMOVED******REMOVED***
						"key": 12,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			"structure_list_deep.json",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"bar": []map[string]interface***REMOVED******REMOVED******REMOVED***
					map[string]interface***REMOVED******REMOVED******REMOVED***
						"foo": []map[string]interface***REMOVED******REMOVED******REMOVED***
							map[string]interface***REMOVED******REMOVED******REMOVED***
								"name": "terraform_example",
								"ingress": []map[string]interface***REMOVED******REMOVED******REMOVED***
									map[string]interface***REMOVED******REMOVED******REMOVED***
										"from_port": 22,
									***REMOVED***,
									map[string]interface***REMOVED******REMOVED******REMOVED***
										"from_port": 80,
									***REMOVED***,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,

		***REMOVED***
			"structure_list_empty.json",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"foo": []interface***REMOVED******REMOVED******REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,

		***REMOVED***
			"nested_block_comment.hcl",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"bar": "value",
			***REMOVED***,
		***REMOVED***,

		***REMOVED***
			"unterminated_block_comment.hcl",
			true,
			nil,
		***REMOVED***,

		***REMOVED***
			"unterminated_brace.hcl",
			true,
			nil,
		***REMOVED***,

		***REMOVED***
			"nested_provider_bad.hcl",
			true,
			nil,
		***REMOVED***,

		***REMOVED***
			"object_list.json",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"resource": []map[string]interface***REMOVED******REMOVED******REMOVED***
					map[string]interface***REMOVED******REMOVED******REMOVED***
						"aws_instance": []map[string]interface***REMOVED******REMOVED******REMOVED***
							map[string]interface***REMOVED******REMOVED******REMOVED***
								"db": []map[string]interface***REMOVED******REMOVED******REMOVED***
									map[string]interface***REMOVED******REMOVED******REMOVED***
										"vpc": "foo",
										"provisioner": []map[string]interface***REMOVED******REMOVED******REMOVED***
											map[string]interface***REMOVED******REMOVED******REMOVED***
												"file": []map[string]interface***REMOVED******REMOVED******REMOVED***
													map[string]interface***REMOVED******REMOVED******REMOVED***
														"source":      "foo",
														"destination": "bar",
													***REMOVED***,
												***REMOVED***,
											***REMOVED***,
										***REMOVED***,
									***REMOVED***,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,

		// Terraform GH-8295 sanity test that basic decoding into
		// interface***REMOVED******REMOVED*** works.
		***REMOVED***
			"terraform_variable_invalid.json",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"variable": []map[string]interface***REMOVED******REMOVED******REMOVED***
					map[string]interface***REMOVED******REMOVED******REMOVED***
						"whatever": "abc123",
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,

		***REMOVED***
			"interpolate.json",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"default": `$***REMOVED***replace("europe-west", "-", " ")***REMOVED***`,
			***REMOVED***,
		***REMOVED***,

		***REMOVED***
			"block_assign.hcl",
			true,
			nil,
		***REMOVED***,

		***REMOVED***
			"escape_backslash.hcl",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"output": []map[string]interface***REMOVED******REMOVED******REMOVED***
					map[string]interface***REMOVED******REMOVED******REMOVED***
						"one":  `$***REMOVED***replace(var.sub_domain, ".", "\\.")***REMOVED***`,
						"two":  `$***REMOVED***replace(var.sub_domain, ".", "\\\\.")***REMOVED***`,
						"many": `$***REMOVED***replace(var.sub_domain, ".", "\\\\\\\\.")***REMOVED***`,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,

		***REMOVED***
			"git_crypt.hcl",
			true,
			nil,
		***REMOVED***,

		***REMOVED***
			"object_with_bool.hcl",
			false,
			map[string]interface***REMOVED******REMOVED******REMOVED***
				"path": []map[string]interface***REMOVED******REMOVED******REMOVED***
					map[string]interface***REMOVED******REMOVED******REMOVED***
						"policy": "write",
						"permissions": []map[string]interface***REMOVED******REMOVED******REMOVED***
							map[string]interface***REMOVED******REMOVED******REMOVED***
								"bool": []interface***REMOVED******REMOVED******REMOVED***false***REMOVED***,
							***REMOVED***,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for _, tc := range cases ***REMOVED***
		t.Run(tc.File, func(t *testing.T) ***REMOVED***
			d, err := ioutil.ReadFile(filepath.Join(fixtureDir, tc.File))
			if err != nil ***REMOVED***
				t.Fatalf("err: %s", err)
			***REMOVED***

			var out interface***REMOVED******REMOVED***
			err = Decode(&out, string(d))
			if (err != nil) != tc.Err ***REMOVED***
				t.Fatalf("Input: %s\n\nError: %s", tc.File, err)
			***REMOVED***

			if !reflect.DeepEqual(out, tc.Out) ***REMOVED***
				t.Fatalf("Input: %s. Actual, Expected.\n\n%#v\n\n%#v", tc.File, out, tc.Out)
			***REMOVED***

			var v interface***REMOVED******REMOVED***
			err = Unmarshal(d, &v)
			if (err != nil) != tc.Err ***REMOVED***
				t.Fatalf("Input: %s\n\nError: %s", tc.File, err)
			***REMOVED***

			if !reflect.DeepEqual(v, tc.Out) ***REMOVED***
				t.Fatalf("Input: %s. Actual, Expected.\n\n%#v\n\n%#v", tc.File, out, tc.Out)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestDecode_interfaceInline(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		Value string
		Err   bool
		Out   interface***REMOVED******REMOVED***
	***REMOVED******REMOVED***
		***REMOVED***"t t e***REMOVED******REMOVED******REMOVED******REMOVED***", true, nil***REMOVED***,
		***REMOVED***"t=0t d ***REMOVED******REMOVED***", true, map[string]interface***REMOVED******REMOVED******REMOVED***"t": 0***REMOVED******REMOVED***,
		***REMOVED***"v=0E0v d***REMOVED******REMOVED***", true, map[string]interface***REMOVED******REMOVED******REMOVED***"v": float64(0)***REMOVED******REMOVED***,
	***REMOVED***

	for _, tc := range cases ***REMOVED***
		t.Logf("Testing: %q", tc.Value)

		var out interface***REMOVED******REMOVED***
		err := Decode(&out, tc.Value)
		if (err != nil) != tc.Err ***REMOVED***
			t.Fatalf("Input: %q\n\nError: %s", tc.Value, err)
		***REMOVED***

		if !reflect.DeepEqual(out, tc.Out) ***REMOVED***
			t.Fatalf("Input: %q. Actual, Expected.\n\n%#v\n\n%#v", tc.Value, out, tc.Out)
		***REMOVED***

		var v interface***REMOVED******REMOVED***
		err = Unmarshal([]byte(tc.Value), &v)
		if (err != nil) != tc.Err ***REMOVED***
			t.Fatalf("Input: %q\n\nError: %s", tc.Value, err)
		***REMOVED***

		if !reflect.DeepEqual(v, tc.Out) ***REMOVED***
			t.Fatalf("Input: %q. Actual, Expected.\n\n%#v\n\n%#v", tc.Value, out, tc.Out)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestDecode_equal(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		One, Two string
	***REMOVED******REMOVED***
		***REMOVED***
			"basic.hcl",
			"basic.json",
		***REMOVED***,
		***REMOVED***
			"float.hcl",
			"float.json",
		***REMOVED***,
		/*
			***REMOVED***
				"structure.hcl",
				"structure.json",
			***REMOVED***,
		*/
		***REMOVED***
			"structure.hcl",
			"structure_flat.json",
		***REMOVED***,
		***REMOVED***
			"terraform_heroku.hcl",
			"terraform_heroku.json",
		***REMOVED***,
	***REMOVED***

	for _, tc := range cases ***REMOVED***
		p1 := filepath.Join(fixtureDir, tc.One)
		p2 := filepath.Join(fixtureDir, tc.Two)

		d1, err := ioutil.ReadFile(p1)
		if err != nil ***REMOVED***
			t.Fatalf("err: %s", err)
		***REMOVED***

		d2, err := ioutil.ReadFile(p2)
		if err != nil ***REMOVED***
			t.Fatalf("err: %s", err)
		***REMOVED***

		var i1, i2 interface***REMOVED******REMOVED***
		err = Decode(&i1, string(d1))
		if err != nil ***REMOVED***
			t.Fatalf("err: %s", err)
		***REMOVED***

		err = Decode(&i2, string(d2))
		if err != nil ***REMOVED***
			t.Fatalf("err: %s", err)
		***REMOVED***

		if !reflect.DeepEqual(i1, i2) ***REMOVED***
			t.Fatalf(
				"%s != %s\n\n%#v\n\n%#v",
				tc.One, tc.Two,
				i1, i2)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestDecode_flatMap(t *testing.T) ***REMOVED***
	var val map[string]map[string]string

	err := Decode(&val, testReadFile(t, "structure_flatmap.hcl"))
	if err != nil ***REMOVED***
		t.Fatalf("err: %s", err)
	***REMOVED***

	expected := map[string]map[string]string***REMOVED***
		"foo": map[string]string***REMOVED***
			"foo": "bar",
			"key": "7",
		***REMOVED***,
	***REMOVED***

	if !reflect.DeepEqual(val, expected) ***REMOVED***
		t.Fatalf("Actual: %#v\n\nExpected: %#v", val, expected)
	***REMOVED***
***REMOVED***

func TestDecode_structure(t *testing.T) ***REMOVED***
	type Embedded interface***REMOVED******REMOVED***

	type V struct ***REMOVED***
		Embedded `hcl:"-"`
		Key      int
		Foo      string
	***REMOVED***

	var actual V

	err := Decode(&actual, testReadFile(t, "flat.hcl"))
	if err != nil ***REMOVED***
		t.Fatalf("err: %s", err)
	***REMOVED***

	expected := V***REMOVED***
		Key: 7,
		Foo: "bar",
	***REMOVED***

	if !reflect.DeepEqual(actual, expected) ***REMOVED***
		t.Fatalf("Actual: %#v\n\nExpected: %#v", actual, expected)
	***REMOVED***
***REMOVED***

func TestDecode_structurePtr(t *testing.T) ***REMOVED***
	type V struct ***REMOVED***
		Key int
		Foo string
	***REMOVED***

	var actual *V

	err := Decode(&actual, testReadFile(t, "flat.hcl"))
	if err != nil ***REMOVED***
		t.Fatalf("err: %s", err)
	***REMOVED***

	expected := &V***REMOVED***
		Key: 7,
		Foo: "bar",
	***REMOVED***

	if !reflect.DeepEqual(actual, expected) ***REMOVED***
		t.Fatalf("Actual: %#v\n\nExpected: %#v", actual, expected)
	***REMOVED***
***REMOVED***

func TestDecode_structureArray(t *testing.T) ***REMOVED***
	// This test is extracted from a failure in Consul (consul.io),
	// hence the interesting structure naming.

	type KeyPolicyType string

	type KeyPolicy struct ***REMOVED***
		Prefix string `hcl:",key"`
		Policy KeyPolicyType
	***REMOVED***

	type Policy struct ***REMOVED***
		Keys []KeyPolicy `hcl:"key,expand"`
	***REMOVED***

	expected := Policy***REMOVED***
		Keys: []KeyPolicy***REMOVED***
			KeyPolicy***REMOVED***
				Prefix: "",
				Policy: "read",
			***REMOVED***,
			KeyPolicy***REMOVED***
				Prefix: "foo/",
				Policy: "write",
			***REMOVED***,
			KeyPolicy***REMOVED***
				Prefix: "foo/bar/",
				Policy: "read",
			***REMOVED***,
			KeyPolicy***REMOVED***
				Prefix: "foo/bar/baz",
				Policy: "deny",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	files := []string***REMOVED***
		"decode_policy.hcl",
		"decode_policy.json",
	***REMOVED***

	for _, f := range files ***REMOVED***
		var actual Policy

		err := Decode(&actual, testReadFile(t, f))
		if err != nil ***REMOVED***
			t.Fatalf("Input: %s\n\nerr: %s", f, err)
		***REMOVED***

		if !reflect.DeepEqual(actual, expected) ***REMOVED***
			t.Fatalf("Input: %s\n\nActual: %#v\n\nExpected: %#v", f, actual, expected)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestDecode_sliceExpand(t *testing.T) ***REMOVED***
	type testInner struct ***REMOVED***
		Name string `hcl:",key"`
		Key  string
	***REMOVED***

	type testStruct struct ***REMOVED***
		Services []testInner `hcl:"service,expand"`
	***REMOVED***

	expected := testStruct***REMOVED***
		Services: []testInner***REMOVED***
			testInner***REMOVED***
				Name: "my-service-0",
				Key:  "value",
			***REMOVED***,
			testInner***REMOVED***
				Name: "my-service-1",
				Key:  "value",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	files := []string***REMOVED***
		"slice_expand.hcl",
	***REMOVED***

	for _, f := range files ***REMOVED***
		t.Logf("Testing: %s", f)

		var actual testStruct
		err := Decode(&actual, testReadFile(t, f))
		if err != nil ***REMOVED***
			t.Fatalf("Input: %s\n\nerr: %s", f, err)
		***REMOVED***

		if !reflect.DeepEqual(actual, expected) ***REMOVED***
			t.Fatalf("Input: %s\n\nActual: %#v\n\nExpected: %#v", f, actual, expected)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestDecode_structureMap(t *testing.T) ***REMOVED***
	// This test is extracted from a failure in Terraform (terraform.io),
	// hence the interesting structure naming.

	type hclVariable struct ***REMOVED***
		Default     interface***REMOVED******REMOVED***
		Description string
		Fields      []string `hcl:",decodedFields"`
	***REMOVED***

	type rawConfig struct ***REMOVED***
		Variable map[string]hclVariable
	***REMOVED***

	expected := rawConfig***REMOVED***
		Variable: map[string]hclVariable***REMOVED***
			"foo": hclVariable***REMOVED***
				Default:     "bar",
				Description: "bar",
				Fields:      []string***REMOVED***"Default", "Description"***REMOVED***,
			***REMOVED***,

			"amis": hclVariable***REMOVED***
				Default: []map[string]interface***REMOVED******REMOVED******REMOVED***
					map[string]interface***REMOVED******REMOVED******REMOVED***
						"east": "foo",
					***REMOVED***,
				***REMOVED***,
				Fields: []string***REMOVED***"Default"***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	files := []string***REMOVED***
		"decode_tf_variable.hcl",
		"decode_tf_variable.json",
	***REMOVED***

	for _, f := range files ***REMOVED***
		t.Logf("Testing: %s", f)

		var actual rawConfig
		err := Decode(&actual, testReadFile(t, f))
		if err != nil ***REMOVED***
			t.Fatalf("Input: %s\n\nerr: %s", f, err)
		***REMOVED***

		if !reflect.DeepEqual(actual, expected) ***REMOVED***
			t.Fatalf("Input: %s\n\nActual: %#v\n\nExpected: %#v", f, actual, expected)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestDecode_structureMapInvalid(t *testing.T) ***REMOVED***
	// Terraform GH-8295

	type hclVariable struct ***REMOVED***
		Default     interface***REMOVED******REMOVED***
		Description string
		Fields      []string `hcl:",decodedFields"`
	***REMOVED***

	type rawConfig struct ***REMOVED***
		Variable map[string]*hclVariable
	***REMOVED***

	var actual rawConfig
	err := Decode(&actual, testReadFile(t, "terraform_variable_invalid.json"))
	if err == nil ***REMOVED***
		t.Fatal("expected error")
	***REMOVED***
***REMOVED***

func TestDecode_interfaceNonPointer(t *testing.T) ***REMOVED***
	var value interface***REMOVED******REMOVED***
	err := Decode(value, testReadFile(t, "basic_int_string.hcl"))
	if err == nil ***REMOVED***
		t.Fatal("should error")
	***REMOVED***
***REMOVED***

func TestDecode_intString(t *testing.T) ***REMOVED***
	var value struct ***REMOVED***
		Count int
	***REMOVED***

	err := Decode(&value, testReadFile(t, "basic_int_string.hcl"))
	if err != nil ***REMOVED***
		t.Fatalf("err: %s", err)
	***REMOVED***

	if value.Count != 3 ***REMOVED***
		t.Fatalf("bad: %#v", value.Count)
	***REMOVED***
***REMOVED***

func TestDecode_float32(t *testing.T) ***REMOVED***
	var value struct ***REMOVED***
		A float32 `hcl:"a"`
		B float32 `hcl:"b"`
	***REMOVED***

	err := Decode(&value, testReadFile(t, "float.hcl"))
	if err != nil ***REMOVED***
		t.Fatalf("err: %s", err)
	***REMOVED***

	if got, want := value.A, float32(1.02); got != want ***REMOVED***
		t.Fatalf("wrong result %#v; want %#v", got, want)
	***REMOVED***
	if got, want := value.B, float32(2); got != want ***REMOVED***
		t.Fatalf("wrong result %#v; want %#v", got, want)
	***REMOVED***
***REMOVED***

func TestDecode_float64(t *testing.T) ***REMOVED***
	var value struct ***REMOVED***
		A float64 `hcl:"a"`
		B float64 `hcl:"b"`
	***REMOVED***

	err := Decode(&value, testReadFile(t, "float.hcl"))
	if err != nil ***REMOVED***
		t.Fatalf("err: %s", err)
	***REMOVED***

	if got, want := value.A, float64(1.02); got != want ***REMOVED***
		t.Fatalf("wrong result %#v; want %#v", got, want)
	***REMOVED***
	if got, want := value.B, float64(2); got != want ***REMOVED***
		t.Fatalf("wrong result %#v; want %#v", got, want)
	***REMOVED***
***REMOVED***

func TestDecode_intStringAliased(t *testing.T) ***REMOVED***
	var value struct ***REMOVED***
		Count time.Duration
	***REMOVED***

	err := Decode(&value, testReadFile(t, "basic_int_string.hcl"))
	if err != nil ***REMOVED***
		t.Fatalf("err: %s", err)
	***REMOVED***

	if value.Count != time.Duration(3) ***REMOVED***
		t.Fatalf("bad: %#v", value.Count)
	***REMOVED***
***REMOVED***

func TestDecode_Node(t *testing.T) ***REMOVED***
	// given
	var value struct ***REMOVED***
		Content ast.Node
		Nested  struct ***REMOVED***
			Content ast.Node
		***REMOVED***
	***REMOVED***

	content := `
content ***REMOVED***
	hello = "world"
***REMOVED***
`

	// when
	err := Decode(&value, content)

	// then
	if err != nil ***REMOVED***
		t.Errorf("unable to decode content, %v", err)
		return
	***REMOVED***

	// verify ast.Node can be decoded later
	var v map[string]interface***REMOVED******REMOVED***
	err = DecodeObject(&v, value.Content)
	if err != nil ***REMOVED***
		t.Errorf("unable to decode content, %v", err)
		return
	***REMOVED***

	if v["hello"] != "world" ***REMOVED***
		t.Errorf("expected mapping to be returned")
	***REMOVED***
***REMOVED***

func TestDecode_NestedNode(t *testing.T) ***REMOVED***
	// given
	var value struct ***REMOVED***
		Nested struct ***REMOVED***
			Content ast.Node
		***REMOVED***
	***REMOVED***

	content := `
nested "content" ***REMOVED***
	hello = "world"
***REMOVED***
`

	// when
	err := Decode(&value, content)

	// then
	if err != nil ***REMOVED***
		t.Errorf("unable to decode content, %v", err)
		return
	***REMOVED***

	// verify ast.Node can be decoded later
	var v map[string]interface***REMOVED******REMOVED***
	err = DecodeObject(&v, value.Nested.Content)
	if err != nil ***REMOVED***
		t.Errorf("unable to decode content, %v", err)
		return
	***REMOVED***

	if v["hello"] != "world" ***REMOVED***
		t.Errorf("expected mapping to be returned")
	***REMOVED***
***REMOVED***

// https://github.com/hashicorp/hcl/issues/60
func TestDecode_topLevelKeys(t *testing.T) ***REMOVED***
	type Template struct ***REMOVED***
		Source string
	***REMOVED***

	templates := struct ***REMOVED***
		Templates []*Template `hcl:"template"`
	***REMOVED******REMOVED******REMOVED***

	err := Decode(&templates, `
	template ***REMOVED***
	    source = "blah"
	***REMOVED***

	template ***REMOVED***
	    source = "blahblah"
	***REMOVED***`)

	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if templates.Templates[0].Source != "blah" ***REMOVED***
		t.Errorf("bad source: %s", templates.Templates[0].Source)
	***REMOVED***

	if templates.Templates[1].Source != "blahblah" ***REMOVED***
		t.Errorf("bad source: %s", templates.Templates[1].Source)
	***REMOVED***
***REMOVED***

func TestDecode_flattenedJSON(t *testing.T) ***REMOVED***
	// make sure we can also correctly extract a Name key too
	type V struct ***REMOVED***
		Name        string `hcl:",key"`
		Description string
		Default     map[string]string
	***REMOVED***
	type Vars struct ***REMOVED***
		Variable []*V
	***REMOVED***

	cases := []struct ***REMOVED***
		JSON     string
		Out      interface***REMOVED******REMOVED***
		Expected interface***REMOVED******REMOVED***
	***REMOVED******REMOVED***
		***REMOVED*** // Nested object, no sibling keys
			JSON: `
***REMOVED***
  "var_name": ***REMOVED***
    "default": ***REMOVED***
      "key1": "a",
      "key2": "b"
***REMOVED***
  ***REMOVED***
***REMOVED***
			`,
			Out: &[]*V***REMOVED******REMOVED***,
			Expected: &[]*V***REMOVED***
				&V***REMOVED***
					Name:    "var_name",
					Default: map[string]string***REMOVED***"key1": "a", "key2": "b"***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,

		***REMOVED*** // Nested object with a sibling key (this worked previously)
			JSON: `
***REMOVED***
  "var_name": ***REMOVED***
    "description": "Described",
    "default": ***REMOVED***
      "key1": "a",
      "key2": "b"
***REMOVED***
  ***REMOVED***
***REMOVED***
			`,
			Out: &[]*V***REMOVED******REMOVED***,
			Expected: &[]*V***REMOVED***
				&V***REMOVED***
					Name:        "var_name",
					Description: "Described",
					Default:     map[string]string***REMOVED***"key1": "a", "key2": "b"***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,

		***REMOVED*** // Multiple nested objects, one with a sibling key
			JSON: `
***REMOVED***
  "variable": ***REMOVED***
    "var_1": ***REMOVED***
      "default": ***REMOVED***
        "key1": "a",
        "key2": "b"
  ***REMOVED***
***REMOVED***,
    "var_2": ***REMOVED***
      "description": "Described",
      "default": ***REMOVED***
        "key1": "a",
        "key2": "b"
  ***REMOVED***
***REMOVED***
  ***REMOVED***
***REMOVED***
			`,
			Out: &Vars***REMOVED******REMOVED***,
			Expected: &Vars***REMOVED***
				Variable: []*V***REMOVED***
					&V***REMOVED***
						Name:    "var_1",
						Default: map[string]string***REMOVED***"key1": "a", "key2": "b"***REMOVED***,
					***REMOVED***,
					&V***REMOVED***
						Name:        "var_2",
						Description: "Described",
						Default:     map[string]string***REMOVED***"key1": "a", "key2": "b"***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,

		***REMOVED*** // Nested object to maps
			JSON: `
***REMOVED***
  "variable": ***REMOVED***
    "var_name": ***REMOVED***
      "description": "Described",
      "default": ***REMOVED***
        "key1": "a",
        "key2": "b"
  ***REMOVED***
***REMOVED***
  ***REMOVED***
***REMOVED***
			`,
			Out: &[]map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***,
			Expected: &[]map[string]interface***REMOVED******REMOVED******REMOVED***
				***REMOVED***
					"variable": []map[string]interface***REMOVED******REMOVED******REMOVED***
						***REMOVED***
							"var_name": []map[string]interface***REMOVED******REMOVED******REMOVED***
								***REMOVED***
									"description": "Described",
									"default": []map[string]interface***REMOVED******REMOVED******REMOVED***
										***REMOVED***
											"key1": "a",
											"key2": "b",
										***REMOVED***,
									***REMOVED***,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,

		***REMOVED*** // Nested object to maps without a sibling key should decode the same as above
			JSON: `
***REMOVED***
  "variable": ***REMOVED***
    "var_name": ***REMOVED***
      "default": ***REMOVED***
        "key1": "a",
        "key2": "b"
  ***REMOVED***
***REMOVED***
  ***REMOVED***
***REMOVED***
			`,
			Out: &[]map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***,
			Expected: &[]map[string]interface***REMOVED******REMOVED******REMOVED***
				***REMOVED***
					"variable": []map[string]interface***REMOVED******REMOVED******REMOVED***
						***REMOVED***
							"var_name": []map[string]interface***REMOVED******REMOVED******REMOVED***
								***REMOVED***
									"default": []map[string]interface***REMOVED******REMOVED******REMOVED***
										***REMOVED***
											"key1": "a",
											"key2": "b",
										***REMOVED***,
									***REMOVED***,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,

		***REMOVED*** // Nested objects, one with a sibling key, and one without
			JSON: `
***REMOVED***
  "variable": ***REMOVED***
    "var_1": ***REMOVED***
      "default": ***REMOVED***
        "key1": "a",
        "key2": "b"
  ***REMOVED***
***REMOVED***,
    "var_2": ***REMOVED***
      "description": "Described",
      "default": ***REMOVED***
        "key1": "a",
        "key2": "b"
  ***REMOVED***
***REMOVED***
  ***REMOVED***
***REMOVED***
			`,
			Out: &[]map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***,
			Expected: &[]map[string]interface***REMOVED******REMOVED******REMOVED***
				***REMOVED***
					"variable": []map[string]interface***REMOVED******REMOVED******REMOVED***
						***REMOVED***
							"var_1": []map[string]interface***REMOVED******REMOVED******REMOVED***
								***REMOVED***
									"default": []map[string]interface***REMOVED******REMOVED******REMOVED***
										***REMOVED***
											"key1": "a",
											"key2": "b",
										***REMOVED***,
									***REMOVED***,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
				***REMOVED***
					"variable": []map[string]interface***REMOVED******REMOVED******REMOVED***
						***REMOVED***
							"var_2": []map[string]interface***REMOVED******REMOVED******REMOVED***
								***REMOVED***
									"description": "Described",
									"default": []map[string]interface***REMOVED******REMOVED******REMOVED***
										***REMOVED***
											"key1": "a",
											"key2": "b",
										***REMOVED***,
									***REMOVED***,
								***REMOVED***,
							***REMOVED***,
						***REMOVED***,
					***REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for i, tc := range cases ***REMOVED***
		err := Decode(tc.Out, tc.JSON)
		if err != nil ***REMOVED***
			t.Fatalf("[%d] err: %s", i, err)
		***REMOVED***

		if !reflect.DeepEqual(tc.Out, tc.Expected) ***REMOVED***
			t.Fatalf("[%d]\ngot: %s\nexpected: %s\n", i, spew.Sdump(tc.Out), spew.Sdump(tc.Expected))
		***REMOVED***
	***REMOVED***
***REMOVED***
