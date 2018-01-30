package cascadia

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func MustParseHTML(doc string) *html.Node ***REMOVED***
	dom, err := html.Parse(strings.NewReader(doc))
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return dom
***REMOVED***

var selector = MustCompile(`div.matched`)
var doc = `<!DOCTYPE html>
<html>
<body>
<div class="matched">
  <div>
    <div class="matched"></div>
    <div class="matched"></div>
    <div class="matched"></div>
    <div class="matched"></div>
    <div class="matched"></div>
    <div class="matched"></div>
    <div class="matched"></div>
    <div class="matched"></div>
    <div class="matched"></div>
    <div class="matched"></div>
    <div class="matched"></div>
    <div class="matched"></div>
    <div class="matched"></div>
    <div class="matched"></div>
    <div class="matched"></div>
    <div class="matched"></div>
  </div>
</div>
</body>
</html>
`
var dom = MustParseHTML(doc)

func BenchmarkMatchAll(b *testing.B) ***REMOVED***
	var matches []*html.Node
	for i := 0; i < b.N; i++ ***REMOVED***
		matches = selector.MatchAll(dom)
	***REMOVED***
	_ = matches
***REMOVED***
