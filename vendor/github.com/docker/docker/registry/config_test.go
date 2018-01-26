package registry

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadAllowNondistributableArtifacts(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		registries []string
		cidrStrs   []string
		hostnames  []string
		err        string
	***REMOVED******REMOVED***
		***REMOVED***
			registries: []string***REMOVED***"1.2.3.0/24"***REMOVED***,
			cidrStrs:   []string***REMOVED***"1.2.3.0/24"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***"2001:db8::/120"***REMOVED***,
			cidrStrs:   []string***REMOVED***"2001:db8::/120"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***"127.0.0.1"***REMOVED***,
			hostnames:  []string***REMOVED***"127.0.0.1"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***"127.0.0.1:8080"***REMOVED***,
			hostnames:  []string***REMOVED***"127.0.0.1:8080"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***"2001:db8::1"***REMOVED***,
			hostnames:  []string***REMOVED***"2001:db8::1"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***"[2001:db8::1]:80"***REMOVED***,
			hostnames:  []string***REMOVED***"[2001:db8::1]:80"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***"[2001:db8::1]:80"***REMOVED***,
			hostnames:  []string***REMOVED***"[2001:db8::1]:80"***REMOVED***,
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***"1.2.3.0/24", "2001:db8::/120", "127.0.0.1", "127.0.0.1:8080"***REMOVED***,
			cidrStrs:   []string***REMOVED***"1.2.3.0/24", "2001:db8::/120"***REMOVED***,
			hostnames:  []string***REMOVED***"127.0.0.1", "127.0.0.1:8080"***REMOVED***,
		***REMOVED***,

		***REMOVED***
			registries: []string***REMOVED***"http://mytest.com"***REMOVED***,
			err:        "allow-nondistributable-artifacts registry http://mytest.com should not contain '://'",
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***"https://mytest.com"***REMOVED***,
			err:        "allow-nondistributable-artifacts registry https://mytest.com should not contain '://'",
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***"HTTP://mytest.com"***REMOVED***,
			err:        "allow-nondistributable-artifacts registry HTTP://mytest.com should not contain '://'",
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***"svn://mytest.com"***REMOVED***,
			err:        "allow-nondistributable-artifacts registry svn://mytest.com should not contain '://'",
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***"-invalid-registry"***REMOVED***,
			err:        "Cannot begin or end with a hyphen",
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***`mytest-.com`***REMOVED***,
			err:        `allow-nondistributable-artifacts registry mytest-.com is not valid: invalid host "mytest-.com"`,
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***`1200:0000:AB00:1234:0000:2552:7777:1313:8080`***REMOVED***,
			err:        `allow-nondistributable-artifacts registry 1200:0000:AB00:1234:0000:2552:7777:1313:8080 is not valid: invalid host "1200:0000:AB00:1234:0000:2552:7777:1313:8080"`,
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***`mytest.com:500000`***REMOVED***,
			err:        `allow-nondistributable-artifacts registry mytest.com:500000 is not valid: invalid port "500000"`,
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***`"mytest.com"`***REMOVED***,
			err:        `allow-nondistributable-artifacts registry "mytest.com" is not valid: invalid host "\"mytest.com\""`,
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***`"mytest.com:5000"`***REMOVED***,
			err:        `allow-nondistributable-artifacts registry "mytest.com:5000" is not valid: invalid host "\"mytest.com"`,
		***REMOVED***,
	***REMOVED***
	for _, testCase := range testCases ***REMOVED***
		config := emptyServiceConfig
		err := config.LoadAllowNondistributableArtifacts(testCase.registries)
		if testCase.err == "" ***REMOVED***
			if err != nil ***REMOVED***
				t.Fatalf("expect no error, got '%s'", err)
			***REMOVED***

			cidrStrs := []string***REMOVED******REMOVED***
			for _, c := range config.AllowNondistributableArtifactsCIDRs ***REMOVED***
				cidrStrs = append(cidrStrs, c.String())
			***REMOVED***

			sort.Strings(testCase.cidrStrs)
			sort.Strings(cidrStrs)
			if (len(testCase.cidrStrs) > 0 || len(cidrStrs) > 0) && !reflect.DeepEqual(testCase.cidrStrs, cidrStrs) ***REMOVED***
				t.Fatalf("expect AllowNondistributableArtifactsCIDRs to be '%+v', got '%+v'", testCase.cidrStrs, cidrStrs)
			***REMOVED***

			sort.Strings(testCase.hostnames)
			sort.Strings(config.AllowNondistributableArtifactsHostnames)
			if (len(testCase.hostnames) > 0 || len(config.AllowNondistributableArtifactsHostnames) > 0) && !reflect.DeepEqual(testCase.hostnames, config.AllowNondistributableArtifactsHostnames) ***REMOVED***
				t.Fatalf("expect AllowNondistributableArtifactsHostnames to be '%+v', got '%+v'", testCase.hostnames, config.AllowNondistributableArtifactsHostnames)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if err == nil ***REMOVED***
				t.Fatalf("expect error '%s', got no error", testCase.err)
			***REMOVED***
			if !strings.Contains(err.Error(), testCase.err) ***REMOVED***
				t.Fatalf("expect error '%s', got '%s'", testCase.err, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestValidateMirror(t *testing.T) ***REMOVED***
	valid := []string***REMOVED***
		"http://mirror-1.com",
		"http://mirror-1.com/",
		"https://mirror-1.com",
		"https://mirror-1.com/",
		"http://localhost",
		"https://localhost",
		"http://localhost:5000",
		"https://localhost:5000",
		"http://127.0.0.1",
		"https://127.0.0.1",
		"http://127.0.0.1:5000",
		"https://127.0.0.1:5000",
	***REMOVED***

	invalid := []string***REMOVED***
		"!invalid!://%as%",
		"ftp://mirror-1.com",
		"http://mirror-1.com/?q=foo",
		"http://mirror-1.com/v1/",
		"http://mirror-1.com/v1/?q=foo",
		"http://mirror-1.com/v1/?q=foo#frag",
		"http://mirror-1.com?q=foo",
		"https://mirror-1.com#frag",
		"https://mirror-1.com/#frag",
		"http://foo:bar@mirror-1.com/",
		"https://mirror-1.com/v1/",
		"https://mirror-1.com/v1/#",
		"https://mirror-1.com?q",
	***REMOVED***

	for _, address := range valid ***REMOVED***
		if ret, err := ValidateMirror(address); err != nil || ret == "" ***REMOVED***
			t.Errorf("ValidateMirror(`"+address+"`) got %s %s", ret, err)
		***REMOVED***
	***REMOVED***

	for _, address := range invalid ***REMOVED***
		if ret, err := ValidateMirror(address); err == nil || ret != "" ***REMOVED***
			t.Errorf("ValidateMirror(`"+address+"`) got %s %s", ret, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestLoadInsecureRegistries(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		registries []string
		index      string
		err        string
	***REMOVED******REMOVED***
		***REMOVED***
			registries: []string***REMOVED***"127.0.0.1"***REMOVED***,
			index:      "127.0.0.1",
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***"127.0.0.1:8080"***REMOVED***,
			index:      "127.0.0.1:8080",
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***"2001:db8::1"***REMOVED***,
			index:      "2001:db8::1",
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***"[2001:db8::1]:80"***REMOVED***,
			index:      "[2001:db8::1]:80",
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***"http://mytest.com"***REMOVED***,
			index:      "mytest.com",
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***"https://mytest.com"***REMOVED***,
			index:      "mytest.com",
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***"HTTP://mytest.com"***REMOVED***,
			index:      "mytest.com",
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***"svn://mytest.com"***REMOVED***,
			err:        "insecure registry svn://mytest.com should not contain '://'",
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***"-invalid-registry"***REMOVED***,
			err:        "Cannot begin or end with a hyphen",
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***`mytest-.com`***REMOVED***,
			err:        `insecure registry mytest-.com is not valid: invalid host "mytest-.com"`,
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***`1200:0000:AB00:1234:0000:2552:7777:1313:8080`***REMOVED***,
			err:        `insecure registry 1200:0000:AB00:1234:0000:2552:7777:1313:8080 is not valid: invalid host "1200:0000:AB00:1234:0000:2552:7777:1313:8080"`,
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***`mytest.com:500000`***REMOVED***,
			err:        `insecure registry mytest.com:500000 is not valid: invalid port "500000"`,
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***`"mytest.com"`***REMOVED***,
			err:        `insecure registry "mytest.com" is not valid: invalid host "\"mytest.com\""`,
		***REMOVED***,
		***REMOVED***
			registries: []string***REMOVED***`"mytest.com:5000"`***REMOVED***,
			err:        `insecure registry "mytest.com:5000" is not valid: invalid host "\"mytest.com"`,
		***REMOVED***,
	***REMOVED***
	for _, testCase := range testCases ***REMOVED***
		config := emptyServiceConfig
		err := config.LoadInsecureRegistries(testCase.registries)
		if testCase.err == "" ***REMOVED***
			if err != nil ***REMOVED***
				t.Fatalf("expect no error, got '%s'", err)
			***REMOVED***
			match := false
			for index := range config.IndexConfigs ***REMOVED***
				if index == testCase.index ***REMOVED***
					match = true
				***REMOVED***
			***REMOVED***
			if !match ***REMOVED***
				t.Fatalf("expect index configs to contain '%s', got %+v", testCase.index, config.IndexConfigs)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if err == nil ***REMOVED***
				t.Fatalf("expect error '%s', got no error", testCase.err)
			***REMOVED***
			if !strings.Contains(err.Error(), testCase.err) ***REMOVED***
				t.Fatalf("expect error '%s', got '%s'", testCase.err, err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestNewServiceConfig(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		opts   ServiceOptions
		errStr string
	***REMOVED******REMOVED***
		***REMOVED***
			ServiceOptions***REMOVED******REMOVED***,
			"",
		***REMOVED***,
		***REMOVED***
			ServiceOptions***REMOVED***
				Mirrors: []string***REMOVED***"example.com:5000"***REMOVED***,
			***REMOVED***,
			`invalid mirror: unsupported scheme "example.com" in "example.com:5000"`,
		***REMOVED***,
		***REMOVED***
			ServiceOptions***REMOVED***
				Mirrors: []string***REMOVED***"http://example.com:5000"***REMOVED***,
			***REMOVED***,
			"",
		***REMOVED***,
		***REMOVED***
			ServiceOptions***REMOVED***
				InsecureRegistries: []string***REMOVED***"[fe80::]/64"***REMOVED***,
			***REMOVED***,
			`insecure registry [fe80::]/64 is not valid: invalid host "[fe80::]/64"`,
		***REMOVED***,
		***REMOVED***
			ServiceOptions***REMOVED***
				InsecureRegistries: []string***REMOVED***"102.10.8.1/24"***REMOVED***,
			***REMOVED***,
			"",
		***REMOVED***,
		***REMOVED***
			ServiceOptions***REMOVED***
				AllowNondistributableArtifacts: []string***REMOVED***"[fe80::]/64"***REMOVED***,
			***REMOVED***,
			`allow-nondistributable-artifacts registry [fe80::]/64 is not valid: invalid host "[fe80::]/64"`,
		***REMOVED***,
		***REMOVED***
			ServiceOptions***REMOVED***
				AllowNondistributableArtifacts: []string***REMOVED***"102.10.8.1/24"***REMOVED***,
			***REMOVED***,
			"",
		***REMOVED***,
	***REMOVED***

	for _, testCase := range testCases ***REMOVED***
		_, err := newServiceConfig(testCase.opts)
		if testCase.errStr != "" ***REMOVED***
			assert.EqualError(t, err, testCase.errStr)
		***REMOVED*** else ***REMOVED***
			assert.Nil(t, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestValidateIndexName(t *testing.T) ***REMOVED***
	valid := []struct ***REMOVED***
		index  string
		expect string
	***REMOVED******REMOVED***
		***REMOVED***
			index:  "index.docker.io",
			expect: "docker.io",
		***REMOVED***,
		***REMOVED***
			index:  "example.com",
			expect: "example.com",
		***REMOVED***,
		***REMOVED***
			index:  "127.0.0.1:8080",
			expect: "127.0.0.1:8080",
		***REMOVED***,
		***REMOVED***
			index:  "mytest-1.com",
			expect: "mytest-1.com",
		***REMOVED***,
		***REMOVED***
			index:  "mirror-1.com/v1/?q=foo",
			expect: "mirror-1.com/v1/?q=foo",
		***REMOVED***,
	***REMOVED***

	for _, testCase := range valid ***REMOVED***
		result, err := ValidateIndexName(testCase.index)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, testCase.expect, result)
		***REMOVED***

	***REMOVED***

***REMOVED***

func TestValidateIndexNameWithError(t *testing.T) ***REMOVED***
	invalid := []struct ***REMOVED***
		index string
		err   string
	***REMOVED******REMOVED***
		***REMOVED***
			index: "docker.io-",
			err:   "invalid index name (docker.io-). Cannot begin or end with a hyphen",
		***REMOVED***,
		***REMOVED***
			index: "-example.com",
			err:   "invalid index name (-example.com). Cannot begin or end with a hyphen",
		***REMOVED***,
		***REMOVED***
			index: "mirror-1.com/v1/?q=foo-",
			err:   "invalid index name (mirror-1.com/v1/?q=foo-). Cannot begin or end with a hyphen",
		***REMOVED***,
	***REMOVED***
	for _, testCase := range invalid ***REMOVED***
		_, err := ValidateIndexName(testCase.index)
		assert.EqualError(t, err, testCase.err)
	***REMOVED***
***REMOVED***
