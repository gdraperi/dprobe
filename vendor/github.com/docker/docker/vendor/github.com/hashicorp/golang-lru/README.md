golang-lru
==========

This provides the `lru` package which implements a fixed-size
thread safe LRU cache. It is based on the cache in Groupcache.

Documentation
=============

Full docs are available on [Godoc](http://godoc.org/github.com/hashicorp/golang-lru)

Example
=======

Using the LRU is very simple:

```go
l, _ := New(128)
for i := 0; i < 256; i++ ***REMOVED***
    l.Add(i, nil)
***REMOVED***
if l.Len() != 128 ***REMOVED***
    panic(fmt.Sprintf("bad len: %v", l.Len()))
***REMOVED***
```
