package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/docker/docker/pkg/namesgenerator"
)

func main() ***REMOVED***
	rand.Seed(time.Now().UnixNano())
	fmt.Println(namesgenerator.GetRandomName(0))
***REMOVED***
