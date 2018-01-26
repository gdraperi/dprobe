# libnetwork - networking for containers

[![Circle CI](https://circleci.com/gh/docker/libnetwork/tree/master.svg?style=svg)](https://circleci.com/gh/docker/libnetwork/tree/master) [![Coverage Status](https://coveralls.io/repos/docker/libnetwork/badge.svg)](https://coveralls.io/r/docker/libnetwork) [![GoDoc](https://godoc.org/github.com/docker/libnetwork?status.svg)](https://godoc.org/github.com/docker/libnetwork)

Libnetwork provides a native Go implementation for connecting containers

The goal of libnetwork is to deliver a robust Container Network Model that provides a consistent programming interface and the required network abstractions for applications.

#### Design
Please refer to the [design](docs/design.md) for more information.

#### Using libnetwork

There are many networking solutions available to suit a broad range of use-cases. libnetwork uses a driver / plugin model to support all of these solutions while abstracting the complexity of the driver implementations by exposing a simple and consistent Network Model to users.


```go
func main() ***REMOVED***
	if reexec.Init() ***REMOVED***
		return
	***REMOVED***

	// Select and configure the network driver
	networkType := "bridge"

	// Create a new controller instance
	driverOptions := options.Generic***REMOVED******REMOVED***
	genericOption := make(map[string]interface***REMOVED******REMOVED***)
	genericOption[netlabel.GenericData] = driverOptions
	controller, err := libnetwork.New(config.OptionDriverConfig(networkType, genericOption))
	if err != nil ***REMOVED***
		log.Fatalf("libnetwork.New: %s", err)
	***REMOVED***

	// Create a network for containers to join.
	// NewNetwork accepts Variadic optional arguments that libnetwork and Drivers can use.
	network, err := controller.NewNetwork(networkType, "network1", "")
	if err != nil ***REMOVED***
		log.Fatalf("controller.NewNetwork: %s", err)
	***REMOVED***

	// For each new container: allocate IP and interfaces. The returned network
	// settings will be used for container infos (inspect and such), as well as
	// iptables rules for port publishing. This info is contained or accessible
	// from the returned endpoint.
	ep, err := network.CreateEndpoint("Endpoint1")
	if err != nil ***REMOVED***
		log.Fatalf("network.CreateEndpoint: %s", err)
	***REMOVED***

	// Create the sandbox for the container.
	// NewSandbox accepts Variadic optional arguments which libnetwork can use.
	sbx, err := controller.NewSandbox("container1",
		libnetwork.OptionHostname("test"),
		libnetwork.OptionDomainname("docker.io"))
	if err != nil ***REMOVED***
		log.Fatalf("controller.NewSandbox: %s", err)
	***REMOVED***

	// A sandbox can join the endpoint via the join api.
	err = ep.Join(sbx)
	if err != nil ***REMOVED***
		log.Fatalf("ep.Join: %s", err)
	***REMOVED***

	// libnetwork client can check the endpoint's operational data via the Info() API
	epInfo, err := ep.DriverInfo()
	if err != nil ***REMOVED***
		log.Fatalf("ep.DriverInfo: %s", err)
	***REMOVED***

	macAddress, ok := epInfo[netlabel.MacAddress]
	if !ok ***REMOVED***
		log.Fatalf("failed to get mac address from endpoint info")
	***REMOVED***

	fmt.Printf("Joined endpoint %s (%s) to sandbox %s (%s)\n", ep.Name(), macAddress, sbx.ContainerID(), sbx.Key())
***REMOVED***
```

## Future
Please refer to [roadmap](ROADMAP.md) for more information.

## Contributing

Want to hack on libnetwork? [Docker's contributions guidelines](https://github.com/docker/docker/blob/master/CONTRIBUTING.md) apply.

## Copyright and license
Code and documentation copyright 2015 Docker, inc. Code released under the Apache 2.0 license. Docs released under Creative commons.
