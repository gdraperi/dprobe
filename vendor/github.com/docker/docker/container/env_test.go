package container

import "testing"

func TestReplaceAndAppendEnvVars(t *testing.T) ***REMOVED***
	var (
		d = []string***REMOVED***"HOME=/", "FOO=foo_default"***REMOVED***
		// remove FOO from env
		// remove BAR from env (nop)
		o = []string***REMOVED***"HOME=/root", "TERM=xterm", "FOO", "BAR"***REMOVED***
	)

	env := ReplaceOrAppendEnvValues(d, o)
	t.Logf("default=%v, override=%v, result=%v", d, o, env)
	if len(env) != 2 ***REMOVED***
		t.Fatalf("expected len of 2 got %d", len(env))
	***REMOVED***
	if env[0] != "HOME=/root" ***REMOVED***
		t.Fatalf("expected HOME=/root got '%s'", env[0])
	***REMOVED***
	if env[1] != "TERM=xterm" ***REMOVED***
		t.Fatalf("expected TERM=xterm got '%s'", env[1])
	***REMOVED***
***REMOVED***
