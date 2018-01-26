# Go client for the Docker Engine API

The `docker` command uses this package to communicate with the daemon. It can also be used by your own Go applications to do anything the command-line interface does – running containers, pulling images, managing swarms, etc.

For example, to list running containers (the equivalent of `docker ps`):

```go
package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() ***REMOVED***
	cli, err := client.NewEnvClient()
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	for _, container := range containers ***REMOVED***
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	***REMOVED***
***REMOVED***
```

[Full documentation is available on GoDoc.](https://godoc.org/github.com/docker/docker/client)
