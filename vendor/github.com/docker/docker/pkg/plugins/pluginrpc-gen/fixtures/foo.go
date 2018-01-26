package foo

import (
	aliasedio "io"

	"github.com/docker/docker/pkg/plugins/pluginrpc-gen/fixtures/otherfixture"
)

type wobble struct ***REMOVED***
	Some      string
	Val       string
	Inception *wobble
***REMOVED***

// Fooer is an empty interface used for tests.
type Fooer interface***REMOVED******REMOVED***

// Fooer2 is an interface used for tests.
type Fooer2 interface ***REMOVED***
	Foo()
***REMOVED***

// Fooer3 is an interface used for tests.
type Fooer3 interface ***REMOVED***
	Foo()
	Bar(a string)
	Baz(a string) (err error)
	Qux(a, b string) (val string, err error)
	Wobble() (w *wobble)
	Wiggle() (w wobble)
	WiggleWobble(a []*wobble, b []wobble, c map[string]*wobble, d map[*wobble]wobble, e map[string][]wobble, f []*otherfixture.Spaceship) (g map[*wobble]wobble, h [][]*wobble, i otherfixture.Spaceship, j *otherfixture.Spaceship, k map[*otherfixture.Spaceship]otherfixture.Spaceship, l []otherfixture.Spaceship)
***REMOVED***

// Fooer4 is an interface used for tests.
type Fooer4 interface ***REMOVED***
	Foo() error
***REMOVED***

// Bar is an interface used for tests.
type Bar interface ***REMOVED***
	Boo(a string, b string) (s string, err error)
***REMOVED***

// Fooer5 is an interface used for tests.
type Fooer5 interface ***REMOVED***
	Foo()
	Bar
***REMOVED***

// Fooer6 is an interface used for tests.
type Fooer6 interface ***REMOVED***
	Foo(a otherfixture.Spaceship)
***REMOVED***

// Fooer7 is an interface used for tests.
type Fooer7 interface ***REMOVED***
	Foo(a *otherfixture.Spaceship)
***REMOVED***

// Fooer8 is an interface used for tests.
type Fooer8 interface ***REMOVED***
	Foo(a map[string]otherfixture.Spaceship)
***REMOVED***

// Fooer9 is an interface used for tests.
type Fooer9 interface ***REMOVED***
	Foo(a map[string]*otherfixture.Spaceship)
***REMOVED***

// Fooer10 is an interface used for tests.
type Fooer10 interface ***REMOVED***
	Foo(a []otherfixture.Spaceship)
***REMOVED***

// Fooer11 is an interface used for tests.
type Fooer11 interface ***REMOVED***
	Foo(a []*otherfixture.Spaceship)
***REMOVED***

// Fooer12 is an interface used for tests.
type Fooer12 interface ***REMOVED***
	Foo(a aliasedio.Reader)
***REMOVED***
