# pty

Pty is a Go package for using unix pseudo-terminals.

## Install

    go get github.com/kr/pty

## Example

```go
package main

import (
	"github.com/kr/pty"
	"io"
	"os"
	"os/exec"
)

func main() ***REMOVED***
	c := exec.Command("grep", "--color=auto", "bar")
	f, err := pty.Start(c)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	go func() ***REMOVED***
		f.Write([]byte("foo\n"))
		f.Write([]byte("bar\n"))
		f.Write([]byte("baz\n"))
		f.Write([]byte***REMOVED***4***REMOVED***) // EOT
	***REMOVED***()
	io.Copy(os.Stdout, f)
***REMOVED***
```
