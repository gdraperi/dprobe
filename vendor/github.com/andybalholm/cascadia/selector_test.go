package cascadia

import (
	"bytes"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

type selectorTest struct ***REMOVED***
	HTML, selector string
	results        []string
***REMOVED***

func nodeString(n *html.Node) string ***REMOVED***
	buf := bytes.NewBufferString("")
	html.Render(buf, n)
	return buf.String()
***REMOVED***

var selectorTests = []selectorTest***REMOVED***
	***REMOVED***
		`<body><address>This address...</address></body>`,
		"address",
		[]string***REMOVED***
			"<address>This address...</address>",
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<!-- comment --><html><head></head><body>text</body></html>`,
		"*",
		[]string***REMOVED***
			"<html><head></head><body>text</body></html>",
			"<head></head>",
			"<body>text</body>",
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<html><head></head><body></body></html>`,
		"*",
		[]string***REMOVED***
			"<html><head></head><body></body></html>",
			"<head></head>",
			"<body></body>",
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p id="foo"><p id="bar">`,
		"#foo",
		[]string***REMOVED***
			`<p id="foo"></p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<ul><li id="t1"><p id="t1">`,
		"li#t1",
		[]string***REMOVED***
			`<li id="t1"><p id="t1"></p></li>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<ol><li id="t4"><li id="t44">`,
		"*#t4",
		[]string***REMOVED***
			`<li id="t4"></li>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<ul><li class="t1"><li class="t2">`,
		".t1",
		[]string***REMOVED***
			`<li class="t1"></li>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p class="t1 t2">`,
		"p.t1",
		[]string***REMOVED***
			`<p class="t1 t2"></p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<div class="test">`,
		"div.teST",
		[]string***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p class="t1 t2">`,
		".t1.fail",
		[]string***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p class="t1 t2">`,
		"p.t1.t2",
		[]string***REMOVED***
			`<p class="t1 t2"></p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p><p title="title">`,
		"p[title]",
		[]string***REMOVED***
			`<p title="title"></p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<address><address title="foo"><address title="bar">`,
		`address[title="foo"]`,
		[]string***REMOVED***
			`<address title="foo"><address title="bar"></address></address>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<address><address title="foo"><address title="bar">`,
		`address[title!="foo"]`,
		[]string***REMOVED***
			`<address><address title="foo"><address title="bar"></address></address></address>`,
			`<address title="bar"></address>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p title="tot foo bar">`,
		`[    	title        ~=       foo    ]`,
		[]string***REMOVED***
			`<p title="tot foo bar"></p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p title="hello world">`,
		`[title~="hello world"]`,
		[]string***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p lang="en"><p lang="en-gb"><p lang="enough"><p lang="fr-en">`,
		`[lang|="en"]`,
		[]string***REMOVED***
			`<p lang="en"></p>`,
			`<p lang="en-gb"></p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p title="foobar"><p title="barfoo">`,
		`[title^="foo"]`,
		[]string***REMOVED***
			`<p title="foobar"></p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p title="foobar"><p title="barfoo">`,
		`[title$="bar"]`,
		[]string***REMOVED***
			`<p title="foobar"></p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p title="foobarufoo">`,
		`[title*="bar"]`,
		[]string***REMOVED***
			`<p title="foobarufoo"></p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p class=" ">This text should be green.</p><p>This text should be green.</p>`,
		`p[class$=" "]`,
		[]string***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p class="">This text should be green.</p><p>This text should be green.</p>`,
		`p[class$=""]`,
		[]string***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p class=" ">This text should be green.</p><p>This text should be green.</p>`,
		`p[class^=" "]`,
		[]string***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p class="">This text should be green.</p><p>This text should be green.</p>`,
		`p[class^=""]`,
		[]string***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p class=" ">This text should be green.</p><p>This text should be green.</p>`,
		`p[class*=" "]`,
		[]string***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p class="">This text should be green.</p><p>This text should be green.</p>`,
		`p[class*=""]`,
		[]string***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<input type="radio" name="Sex" value="F"/>`,
		`input[name=Sex][value=F]`,
		[]string***REMOVED***
			`<input type="radio" name="Sex" value="F"/>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<table border="0" cellpadding="0" cellspacing="0" style="table-layout: fixed; width: 100%; border: 0 dashed; border-color: #FFFFFF"><tr style="height:64px">aaa</tr></table>`,
		`table[border="0"][cellpadding="0"][cellspacing="0"]`,
		[]string***REMOVED***
			`<table border="0" cellpadding="0" cellspacing="0" style="table-layout: fixed; width: 100%; border: 0 dashed; border-color: #FFFFFF"><tbody><tr style="height:64px"></tr></tbody></table>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p class="t1 t2">`,
		".t1:not(.t2)",
		[]string***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<div class="t3">`,
		`div:not(.t1)`,
		[]string***REMOVED***
			`<div class="t3"></div>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<div><div class="t2"><div class="t3">`,
		`div:not([class="t2"])`,
		[]string***REMOVED***
			`<div><div class="t2"><div class="t3"></div></div></div>`,
			`<div class="t3"></div>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<ol><li id=1><li id=2><li id=3></ol>`,
		`li:nth-child(odd)`,
		[]string***REMOVED***
			`<li id="1"></li>`,
			`<li id="3"></li>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<ol><li id=1><li id=2><li id=3></ol>`,
		`li:nth-child(even)`,
		[]string***REMOVED***
			`<li id="2"></li>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<ol><li id=1><li id=2><li id=3></ol>`,
		`li:nth-child(-n+2)`,
		[]string***REMOVED***
			`<li id="1"></li>`,
			`<li id="2"></li>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<ol><li id=1><li id=2><li id=3></ol>`,
		`li:nth-child(3n+1)`,
		[]string***REMOVED***
			`<li id="1"></li>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<ol><li id=1><li id=2><li id=3><li id=4></ol>`,
		`li:nth-last-child(odd)`,
		[]string***REMOVED***
			`<li id="2"></li>`,
			`<li id="4"></li>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<ol><li id=1><li id=2><li id=3><li id=4></ol>`,
		`li:nth-last-child(even)`,
		[]string***REMOVED***
			`<li id="1"></li>`,
			`<li id="3"></li>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<ol><li id=1><li id=2><li id=3><li id=4></ol>`,
		`li:nth-last-child(-n+2)`,
		[]string***REMOVED***
			`<li id="3"></li>`,
			`<li id="4"></li>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<ol><li id=1><li id=2><li id=3><li id=4></ol>`,
		`li:nth-last-child(3n+1)`,
		[]string***REMOVED***
			`<li id="1"></li>`,
			`<li id="4"></li>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p>some text <span id="1">and a span</span><span id="2"> and another</span></p>`,
		`span:first-child`,
		[]string***REMOVED***
			`<span id="1">and a span</span>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<span>a span</span> and some text`,
		`span:last-child`,
		[]string***REMOVED***
			`<span>a span</span>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<address></address><p id=1><p id=2>`,
		`p:nth-of-type(2)`,
		[]string***REMOVED***
			`<p id="2"></p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<address></address><p id=1><p id=2></p><a>`,
		`p:nth-last-of-type(2)`,
		[]string***REMOVED***
			`<p id="1"></p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<address></address><p id=1><p id=2></p><a>`,
		`p:last-of-type`,
		[]string***REMOVED***
			`<p id="2"></p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<address></address><p id=1><p id=2></p><a>`,
		`p:first-of-type`,
		[]string***REMOVED***
			`<p id="1"></p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<div><p id="1"></p><a></a></div><div><p id="2"></p></div>`,
		`p:only-child`,
		[]string***REMOVED***
			`<p id="2"></p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<div><p id="1"></p><a></a></div><div><p id="2"></p><p id="3"></p></div>`,
		`p:only-of-type`,
		[]string***REMOVED***
			`<p id="1"></p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p id="1"><!-- --><p id="2">Hello<p id="3"><span>`,
		`:empty`,
		[]string***REMOVED***
			`<head></head>`,
			`<p id="1"><!-- --></p>`,
			`<span></span>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<div><p id="1"><table><tr><td><p id="2"></table></div><p id="3">`,
		`div p`,
		[]string***REMOVED***
			`<p id="1"><table><tbody><tr><td><p id="2"></p></td></tr></tbody></table></p>`,
			`<p id="2"></p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<div><p id="1"><table><tr><td><p id="2"></table></div><p id="3">`,
		`div table p`,
		[]string***REMOVED***
			`<p id="2"></p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<div><p id="1"><div><p id="2"></div><table><tr><td><p id="3"></table></div>`,
		`div > p`,
		[]string***REMOVED***
			`<p id="1"></p>`,
			`<p id="2"></p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p id="1"><p id="2"></p><address></address><p id="3">`,
		`p ~ p`,
		[]string***REMOVED***
			`<p id="2"></p>`,
			`<p id="3"></p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p id="1"></p>
		 <!--comment-->
		 <p id="2"></p><address></address><p id="3">`,
		`p + p`,
		[]string***REMOVED***
			`<p id="2"></p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<ul><li></li><li></li></ul><p>`,
		`li, p`,
		[]string***REMOVED***
			"<li></li>",
			"<li></li>",
			"<p></p>",
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p id="1"><p id="2"></p><address></address><p id="3">`,
		`p +/*This is a comment*/ p`,
		[]string***REMOVED***
			`<p id="2"></p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p>Text block that <span>wraps inner text</span> and continues</p>`,
		`p:contains("that wraps")`,
		[]string***REMOVED***
			`<p>Text block that <span>wraps inner text</span> and continues</p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p>Text block that <span>wraps inner text</span> and continues</p>`,
		`p:containsOwn("that wraps")`,
		[]string***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p>Text block that <span>wraps inner text</span> and continues</p>`,
		`:containsOwn("inner")`,
		[]string***REMOVED***
			`<span>wraps inner text</span>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p>Text block that <span>wraps inner text</span> and continues</p>`,
		`p:containsOwn("block")`,
		[]string***REMOVED***
			`<p>Text block that <span>wraps inner text</span> and continues</p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<div id="d1"><p id="p1"><span>text content</span></p></div><div id="d2"/>`,
		`div:has(#p1)`,
		[]string***REMOVED***
			`<div id="d1"><p id="p1"><span>text content</span></p></div>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<div id="d1"><p id="p1"><span>contents 1</span></p></div>
		<div id="d2"><p>contents <em>2</em></p></div>`,
		`div:has(:containsOwn("2"))`,
		[]string***REMOVED***
			`<div id="d2"><p>contents <em>2</em></p></div>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<body><div id="d1"><p id="p1"><span>contents 1</span></p></div>
		<div id="d2"><p id="p2">contents <em>2</em></p></div></body>`,
		`body :has(:containsOwn("2"))`,
		[]string***REMOVED***
			`<div id="d2"><p id="p2">contents <em>2</em></p></div>`,
			`<p id="p2">contents <em>2</em></p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<body><div id="d1"><p id="p1"><span>contents 1</span></p></div>
		<div id="d2"><p id="p2">contents <em>2</em></p></div></body>`,
		`body :haschild(:containsOwn("2"))`,
		[]string***REMOVED***
			`<p id="p2">contents <em>2</em></p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p id="p1">0123456789</p><p id="p2">abcdef</p><p id="p3">0123ABCD</p>`,
		`p:matches([\d])`,
		[]string***REMOVED***
			`<p id="p1">0123456789</p>`,
			`<p id="p3">0123ABCD</p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p id="p1">0123456789</p><p id="p2">abcdef</p><p id="p3">0123ABCD</p>`,
		`p:matches([a-z])`,
		[]string***REMOVED***
			`<p id="p2">abcdef</p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p id="p1">0123456789</p><p id="p2">abcdef</p><p id="p3">0123ABCD</p>`,
		`p:matches([a-zA-Z])`,
		[]string***REMOVED***
			`<p id="p2">abcdef</p>`,
			`<p id="p3">0123ABCD</p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p id="p1">0123456789</p><p id="p2">abcdef</p><p id="p3">0123ABCD</p>`,
		`p:matches([^\d])`,
		[]string***REMOVED***
			`<p id="p2">abcdef</p>`,
			`<p id="p3">0123ABCD</p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p id="p1">0123456789</p><p id="p2">abcdef</p><p id="p3">0123ABCD</p>`,
		`p:matches(^(0|a))`,
		[]string***REMOVED***
			`<p id="p1">0123456789</p>`,
			`<p id="p2">abcdef</p>`,
			`<p id="p3">0123ABCD</p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p id="p1">0123456789</p><p id="p2">abcdef</p><p id="p3">0123ABCD</p>`,
		`p:matches(^\d+$)`,
		[]string***REMOVED***
			`<p id="p1">0123456789</p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<p id="p1">0123456789</p><p id="p2">abcdef</p><p id="p3">0123ABCD</p>`,
		`p:not(:matches(^\d+$))`,
		[]string***REMOVED***
			`<p id="p2">abcdef</p>`,
			`<p id="p3">0123ABCD</p>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<div><p id="p1">01234<em>567</em>89</p><div>`,
		`div :matchesOwn(^\d+$)`,
		[]string***REMOVED***
			`<p id="p1">01234<em>567</em>89</p>`,
			`<em>567</em>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<ul>
			<li><a id="a1" href="http://www.google.com/finance"></a>
			<li><a id="a2" href="http://finance.yahoo.com/"></a>
			<li><a id="a2" href="http://finance.untrusted.com/"/>
			<li><a id="a3" href="https://www.google.com/news"/>
			<li><a id="a4" href="http://news.yahoo.com"/>
		</ul>`,
		`[href#=(fina)]:not([href#=(\/\/[^\/]+untrusted)])`,
		[]string***REMOVED***
			`<a id="a1" href="http://www.google.com/finance"></a>`,
			`<a id="a2" href="http://finance.yahoo.com/"></a>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<ul>
			<li><a id="a1" href="http://www.google.com/finance"/>
			<li><a id="a2" href="http://finance.yahoo.com/"/>
			<li><a id="a3" href="https://www.google.com/news"></a>
			<li><a id="a4" href="http://news.yahoo.com"/>
		</ul>`,
		`[href#=(^https:\/\/[^\/]*\/?news)]`,
		[]string***REMOVED***
			`<a id="a3" href="https://www.google.com/news"></a>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<form>
			<label>Username <input type="text" name="username" /></label>
			<label>Password <input type="password" name="password" /></label>
			<label>Country
				<select name="country">
					<option value="ca">Canada</option>
					<option value="us">United States</option>
				</select>
			</label>
			<label>Bio <textarea name="bio"></textarea></label>
			<button>Sign up</button>
		</form>`,
		`:input`,
		[]string***REMOVED***
			`<input type="text" name="username"/>`,
			`<input type="password" name="password"/>`,
			`<select name="country">
					<option value="ca">Canada</option>
					<option value="us">United States</option>
				</select>`,
			`<textarea name="bio"></textarea>`,
			`<button>Sign up</button>`,
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<html><head></head><body></body></html>`,
		":root",
		[]string***REMOVED***
			"<html><head></head><body></body></html>",
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<html><head></head><body></body></html>`,
		"*:root",
		[]string***REMOVED***
			"<html><head></head><body></body></html>",
		***REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<html><head></head><body></body></html>`,
		"*:root:first-child",
		[]string***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<html><head></head><body></body></html>`,
		"*:root:nth-child(1)",
		[]string***REMOVED******REMOVED***,
	***REMOVED***,
	***REMOVED***
		`<html><head></head><body><a href="http://www.foo.com"></a></body></html>`,
		"a:not(:root)",
		[]string***REMOVED***
			`<a href="http://www.foo.com"></a>`,
		***REMOVED***,
	***REMOVED***,
***REMOVED***

func TestSelectors(t *testing.T) ***REMOVED***
	for _, test := range selectorTests ***REMOVED***
		s, err := Compile(test.selector)
		if err != nil ***REMOVED***
			t.Errorf("error compiling %q: %s", test.selector, err)
			continue
		***REMOVED***

		doc, err := html.Parse(strings.NewReader(test.HTML))
		if err != nil ***REMOVED***
			t.Errorf("error parsing %q: %s", test.HTML, err)
			continue
		***REMOVED***

		matches := s.MatchAll(doc)
		if len(matches) != len(test.results) ***REMOVED***
			t.Errorf("selector %s wanted %d elements, got %d instead", test.selector, len(test.results), len(matches))
			continue
		***REMOVED***

		for i, m := range matches ***REMOVED***
			got := nodeString(m)
			if got != test.results[i] ***REMOVED***
				t.Errorf("selector %s wanted %s, got %s instead", test.selector, test.results[i], got)
			***REMOVED***
		***REMOVED***

		firstMatch := s.MatchFirst(doc)
		if len(test.results) == 0 ***REMOVED***
			if firstMatch != nil ***REMOVED***
				t.Errorf("MatchFirst: selector %s want nil, got %s", test.selector, nodeString(firstMatch))
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			got := nodeString(firstMatch)
			if got != test.results[0] ***REMOVED***
				t.Errorf("MatchFirst: selector %s want %s, got %s", test.selector, test.results[0], got)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
