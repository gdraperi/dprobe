package opts

import (
	"fmt"
	"strings"
	"testing"
)

func TestValidateIPAddress(t *testing.T) ***REMOVED***
	if ret, err := ValidateIPAddress(`1.2.3.4`); err != nil || ret == "" ***REMOVED***
		t.Fatalf("ValidateIPAddress(`1.2.3.4`) got %s %s", ret, err)
	***REMOVED***

	if ret, err := ValidateIPAddress(`127.0.0.1`); err != nil || ret == "" ***REMOVED***
		t.Fatalf("ValidateIPAddress(`127.0.0.1`) got %s %s", ret, err)
	***REMOVED***

	if ret, err := ValidateIPAddress(`::1`); err != nil || ret == "" ***REMOVED***
		t.Fatalf("ValidateIPAddress(`::1`) got %s %s", ret, err)
	***REMOVED***

	if ret, err := ValidateIPAddress(`127`); err == nil || ret != "" ***REMOVED***
		t.Fatalf("ValidateIPAddress(`127`) got %s %s", ret, err)
	***REMOVED***

	if ret, err := ValidateIPAddress(`random invalid string`); err == nil || ret != "" ***REMOVED***
		t.Fatalf("ValidateIPAddress(`random invalid string`) got %s %s", ret, err)
	***REMOVED***

***REMOVED***

func TestMapOpts(t *testing.T) ***REMOVED***
	tmpMap := make(map[string]string)
	o := NewMapOpts(tmpMap, logOptsValidator)
	o.Set("max-size=1")
	if o.String() != "map[max-size:1]" ***REMOVED***
		t.Errorf("%s != [map[max-size:1]", o.String())
	***REMOVED***

	o.Set("max-file=2")
	if len(tmpMap) != 2 ***REMOVED***
		t.Errorf("map length %d != 2", len(tmpMap))
	***REMOVED***

	if tmpMap["max-file"] != "2" ***REMOVED***
		t.Errorf("max-file = %s != 2", tmpMap["max-file"])
	***REMOVED***

	if tmpMap["max-size"] != "1" ***REMOVED***
		t.Errorf("max-size = %s != 1", tmpMap["max-size"])
	***REMOVED***
	if o.Set("dummy-val=3") == nil ***REMOVED***
		t.Error("validator is not being called")
	***REMOVED***
***REMOVED***

func TestListOptsWithoutValidator(t *testing.T) ***REMOVED***
	o := NewListOpts(nil)
	o.Set("foo")
	if o.String() != "[foo]" ***REMOVED***
		t.Errorf("%s != [foo]", o.String())
	***REMOVED***
	o.Set("bar")
	if o.Len() != 2 ***REMOVED***
		t.Errorf("%d != 2", o.Len())
	***REMOVED***
	o.Set("bar")
	if o.Len() != 3 ***REMOVED***
		t.Errorf("%d != 3", o.Len())
	***REMOVED***
	if !o.Get("bar") ***REMOVED***
		t.Error("o.Get(\"bar\") == false")
	***REMOVED***
	if o.Get("baz") ***REMOVED***
		t.Error("o.Get(\"baz\") == true")
	***REMOVED***
	o.Delete("foo")
	if o.String() != "[bar bar]" ***REMOVED***
		t.Errorf("%s != [bar bar]", o.String())
	***REMOVED***
	listOpts := o.GetAll()
	if len(listOpts) != 2 || listOpts[0] != "bar" || listOpts[1] != "bar" ***REMOVED***
		t.Errorf("Expected [[bar bar]], got [%v]", listOpts)
	***REMOVED***
	mapListOpts := o.GetMap()
	if len(mapListOpts) != 1 ***REMOVED***
		t.Errorf("Expected [map[bar:***REMOVED******REMOVED***]], got [%v]", mapListOpts)
	***REMOVED***

***REMOVED***

func TestListOptsWithValidator(t *testing.T) ***REMOVED***
	// Re-using logOptsvalidator (used by MapOpts)
	o := NewListOpts(logOptsValidator)
	o.Set("foo")
	if o.String() != "" ***REMOVED***
		t.Errorf(`%s != ""`, o.String())
	***REMOVED***
	o.Set("foo=bar")
	if o.String() != "" ***REMOVED***
		t.Errorf(`%s != ""`, o.String())
	***REMOVED***
	o.Set("max-file=2")
	if o.Len() != 1 ***REMOVED***
		t.Errorf("%d != 1", o.Len())
	***REMOVED***
	if !o.Get("max-file=2") ***REMOVED***
		t.Error("o.Get(\"max-file=2\") == false")
	***REMOVED***
	if o.Get("baz") ***REMOVED***
		t.Error("o.Get(\"baz\") == true")
	***REMOVED***
	o.Delete("max-file=2")
	if o.String() != "" ***REMOVED***
		t.Errorf(`%s != ""`, o.String())
	***REMOVED***
***REMOVED***

func TestValidateDNSSearch(t *testing.T) ***REMOVED***
	valid := []string***REMOVED***
		`.`,
		`a`,
		`a.`,
		`1.foo`,
		`17.foo`,
		`foo.bar`,
		`foo.bar.baz`,
		`foo.bar.`,
		`foo.bar.baz`,
		`foo1.bar2`,
		`foo1.bar2.baz`,
		`1foo.2bar.`,
		`1foo.2bar.baz`,
		`foo-1.bar-2`,
		`foo-1.bar-2.baz`,
		`foo-1.bar-2.`,
		`foo-1.bar-2.baz`,
		`1-foo.2-bar`,
		`1-foo.2-bar.baz`,
		`1-foo.2-bar.`,
		`1-foo.2-bar.baz`,
	***REMOVED***

	invalid := []string***REMOVED***
		``,
		` `,
		`  `,
		`17`,
		`17.`,
		`.17`,
		`17-.`,
		`17-.foo`,
		`.foo`,
		`foo-.bar`,
		`-foo.bar`,
		`foo.bar-`,
		`foo.bar-.baz`,
		`foo.-bar`,
		`foo.-bar.baz`,
		`foo.bar.baz.this.should.fail.on.long.name.because.it.is.longer.thanitshouldbethis.should.fail.on.long.name.because.it.is.longer.thanitshouldbethis.should.fail.on.long.name.because.it.is.longer.thanitshouldbethis.should.fail.on.long.name.because.it.is.longer.thanitshouldbe`,
	***REMOVED***

	for _, domain := range valid ***REMOVED***
		if ret, err := ValidateDNSSearch(domain); err != nil || ret == "" ***REMOVED***
			t.Fatalf("ValidateDNSSearch(`"+domain+"`) got %s %s", ret, err)
		***REMOVED***
	***REMOVED***

	for _, domain := range invalid ***REMOVED***
		if ret, err := ValidateDNSSearch(domain); err == nil || ret != "" ***REMOVED***
			t.Fatalf("ValidateDNSSearch(`"+domain+"`) got %s %s", ret, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestValidateLabel(t *testing.T) ***REMOVED***
	if _, err := ValidateLabel("label"); err == nil || err.Error() != "bad attribute format: label" ***REMOVED***
		t.Fatalf("Expected an error [bad attribute format: label], go %v", err)
	***REMOVED***
	if actual, err := ValidateLabel("key1=value1"); err != nil || actual != "key1=value1" ***REMOVED***
		t.Fatalf("Expected [key1=value1], got [%v,%v]", actual, err)
	***REMOVED***
	// Validate it's working with more than one =
	if actual, err := ValidateLabel("key1=value1=value2"); err != nil ***REMOVED***
		t.Fatalf("Expected [key1=value1=value2], got [%v,%v]", actual, err)
	***REMOVED***
	// Validate it's working with one more
	if actual, err := ValidateLabel("key1=value1=value2=value3"); err != nil ***REMOVED***
		t.Fatalf("Expected [key1=value1=value2=value2], got [%v,%v]", actual, err)
	***REMOVED***
***REMOVED***

func logOptsValidator(val string) (string, error) ***REMOVED***
	allowedKeys := map[string]string***REMOVED***"max-size": "1", "max-file": "2"***REMOVED***
	vals := strings.Split(val, "=")
	if allowedKeys[vals[0]] != "" ***REMOVED***
		return val, nil
	***REMOVED***
	return "", fmt.Errorf("invalid key %s", vals[0])
***REMOVED***

func TestNamedListOpts(t *testing.T) ***REMOVED***
	var v []string
	o := NewNamedListOptsRef("foo-name", &v, nil)

	o.Set("foo")
	if o.String() != "[foo]" ***REMOVED***
		t.Errorf("%s != [foo]", o.String())
	***REMOVED***
	if o.Name() != "foo-name" ***REMOVED***
		t.Errorf("%s != foo-name", o.Name())
	***REMOVED***
	if len(v) != 1 ***REMOVED***
		t.Errorf("expected foo to be in the values, got %v", v)
	***REMOVED***
***REMOVED***

func TestNamedMapOpts(t *testing.T) ***REMOVED***
	tmpMap := make(map[string]string)
	o := NewNamedMapOpts("max-name", tmpMap, nil)

	o.Set("max-size=1")
	if o.String() != "map[max-size:1]" ***REMOVED***
		t.Errorf("%s != [map[max-size:1]", o.String())
	***REMOVED***
	if o.Name() != "max-name" ***REMOVED***
		t.Errorf("%s != max-name", o.Name())
	***REMOVED***
	if _, exist := tmpMap["max-size"]; !exist ***REMOVED***
		t.Errorf("expected map-size to be in the values, got %v", tmpMap)
	***REMOVED***
***REMOVED***

func TestParseLink(t *testing.T) ***REMOVED***
	name, alias, err := ParseLink("name:alias")
	if err != nil ***REMOVED***
		t.Fatalf("Expected not to error out on a valid name:alias format but got: %v", err)
	***REMOVED***
	if name != "name" ***REMOVED***
		t.Fatalf("Link name should have been name, got %s instead", name)
	***REMOVED***
	if alias != "alias" ***REMOVED***
		t.Fatalf("Link alias should have been alias, got %s instead", alias)
	***REMOVED***
	// short format definition
	name, alias, err = ParseLink("name")
	if err != nil ***REMOVED***
		t.Fatalf("Expected not to error out on a valid name only format but got: %v", err)
	***REMOVED***
	if name != "name" ***REMOVED***
		t.Fatalf("Link name should have been name, got %s instead", name)
	***REMOVED***
	if alias != "name" ***REMOVED***
		t.Fatalf("Link alias should have been name, got %s instead", alias)
	***REMOVED***
	// empty string link definition is not allowed
	if _, _, err := ParseLink(""); err == nil || !strings.Contains(err.Error(), "empty string specified for links") ***REMOVED***
		t.Fatalf("Expected error 'empty string specified for links' but got: %v", err)
	***REMOVED***
	// more than two colons are not allowed
	if _, _, err := ParseLink("link:alias:wrong"); err == nil || !strings.Contains(err.Error(), "bad format for links: link:alias:wrong") ***REMOVED***
		t.Fatalf("Expected error 'bad format for links: link:alias:wrong' but got: %v", err)
	***REMOVED***
***REMOVED***
