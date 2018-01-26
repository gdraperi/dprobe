fluent-logger-golang
====

[![Build Status](https://travis-ci.org/fluent/fluent-logger-golang.png?branch=master)](https://travis-ci.org/fluent/fluent-logger-golang)

## A structured event logger for Fluentd (Golang)

## How to install

```
go get github.com/fluent/fluent-logger-golang/fluent
```

## Usage

Install the package with `go get` and use `import` to include it in your project.

```
import "github.com/fluent/fluent-logger-golang/fluent"
```

GoDoc: http://godoc.org/github.com/fluent/fluent-logger-golang/fluent

## Example

```go
package main

import (
  "github.com/fluent/fluent-logger-golang/fluent"
  "fmt"
  "time"
)

func main() ***REMOVED***
  logger, err := fluent.New(fluent.Config***REMOVED******REMOVED***)
  if err != nil ***REMOVED***
    fmt.Println(err)
  ***REMOVED***
  defer logger.Close()
  tag := "myapp.access"
  var data = map[string]string***REMOVED***
    "foo":  "bar",
    "hoge": "hoge",
  ***REMOVED***
  error := logger.Post(tag, data)
  // error := logger.PostWithTime(tag, time.Now(), data)
  if error != nil ***REMOVED***
    panic(error)
  ***REMOVED***
***REMOVED***
```

`data` must be a value like `map[string]literal`, `map[string]interface***REMOVED******REMOVED***`, `struct` or [`msgp.Marshaler`](http://godoc.org/github.com/tinylib/msgp/msgp#Marshaler). Logger refers tags `msg` or `codec` of each fields of structs.

## Setting config values

```go
f := fluent.New(fluent.Config***REMOVED***FluentPort: 80, FluentHost: "example.com"***REMOVED***)
```

### WriteTimeout

Sets the timeout for Write call of logger.Post.
Since the default is zero value, Write will not time out.

## Tests
```
go test
```
