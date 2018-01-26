// +build !windows

package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/docker/docker/daemon/graphdriver/devmapper"
	"github.com/docker/docker/pkg/devicemapper"
	"github.com/sirupsen/logrus"
)

func usage() ***REMOVED***
	fmt.Fprintf(os.Stderr, "Usage: %s <flags>  [status] | [list] | [device id]  | [resize new-pool-size] | [snap new-id base-id] | [remove id] | [mount id mountpoint]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
***REMOVED***

func byteSizeFromString(arg string) (int64, error) ***REMOVED***
	digits := ""
	rest := ""
	last := strings.LastIndexAny(arg, "0123456789")
	if last >= 0 ***REMOVED***
		digits = arg[:last+1]
		rest = arg[last+1:]
	***REMOVED***

	val, err := strconv.ParseInt(digits, 10, 64)
	if err != nil ***REMOVED***
		return val, err
	***REMOVED***

	rest = strings.ToLower(strings.TrimSpace(rest))

	var multiplier int64 = 1
	switch rest ***REMOVED***
	case "":
		multiplier = 1
	case "k", "kb":
		multiplier = 1024
	case "m", "mb":
		multiplier = 1024 * 1024
	case "g", "gb":
		multiplier = 1024 * 1024 * 1024
	case "t", "tb":
		multiplier = 1024 * 1024 * 1024 * 1024
	default:
		return 0, fmt.Errorf("Unknown size unit: %s", rest)
	***REMOVED***

	return val * multiplier, nil
***REMOVED***

func main() ***REMOVED***
	root := flag.String("r", "/var/lib/docker", "Docker root dir")
	flDebug := flag.Bool("D", false, "Debug mode")

	flag.Parse()

	if *flDebug ***REMOVED***
		os.Setenv("DEBUG", "1")
		logrus.SetLevel(logrus.DebugLevel)
	***REMOVED***

	if flag.NArg() < 1 ***REMOVED***
		usage()
	***REMOVED***

	args := flag.Args()

	home := path.Join(*root, "devicemapper")
	devices, err := devmapper.NewDeviceSet(home, false, nil, nil, nil)
	if err != nil ***REMOVED***
		fmt.Println("Can't initialize device mapper: ", err)
		os.Exit(1)
	***REMOVED***

	switch args[0] ***REMOVED***
	case "status":
		status := devices.Status()
		fmt.Printf("Pool name: %s\n", status.PoolName)
		fmt.Printf("Data Loopback file: %s\n", status.DataLoopback)
		fmt.Printf("Metadata Loopback file: %s\n", status.MetadataLoopback)
		fmt.Printf("Sector size: %d\n", status.SectorSize)
		fmt.Printf("Data use: %d of %d (%.1f %%)\n", status.Data.Used, status.Data.Total, 100.0*float64(status.Data.Used)/float64(status.Data.Total))
		fmt.Printf("Metadata use: %d of %d (%.1f %%)\n", status.Metadata.Used, status.Metadata.Total, 100.0*float64(status.Metadata.Used)/float64(status.Metadata.Total))
	case "list":
		ids := devices.List()
		sort.Strings(ids)
		for _, id := range ids ***REMOVED***
			fmt.Println(id)
		***REMOVED***
	case "device":
		if flag.NArg() < 2 ***REMOVED***
			usage()
		***REMOVED***
		status, err := devices.GetDeviceStatus(args[1])
		if err != nil ***REMOVED***
			fmt.Println("Can't get device info: ", err)
			os.Exit(1)
		***REMOVED***
		fmt.Printf("Id: %d\n", status.DeviceID)
		fmt.Printf("Size: %d\n", status.Size)
		fmt.Printf("Transaction Id: %d\n", status.TransactionID)
		fmt.Printf("Size in Sectors: %d\n", status.SizeInSectors)
		fmt.Printf("Mapped Sectors: %d\n", status.MappedSectors)
		fmt.Printf("Highest Mapped Sector: %d\n", status.HighestMappedSector)
	case "resize":
		if flag.NArg() < 2 ***REMOVED***
			usage()
		***REMOVED***

		size, err := byteSizeFromString(args[1])
		if err != nil ***REMOVED***
			fmt.Println("Invalid size: ", err)
			os.Exit(1)
		***REMOVED***

		err = devices.ResizePool(size)
		if err != nil ***REMOVED***
			fmt.Println("Error resizing pool: ", err)
			os.Exit(1)
		***REMOVED***

	case "snap":
		if flag.NArg() < 3 ***REMOVED***
			usage()
		***REMOVED***

		err := devices.AddDevice(args[1], args[2], nil)
		if err != nil ***REMOVED***
			fmt.Println("Can't create snap device: ", err)
			os.Exit(1)
		***REMOVED***
	case "remove":
		if flag.NArg() < 2 ***REMOVED***
			usage()
		***REMOVED***

		err := devicemapper.RemoveDevice(args[1])
		if err != nil ***REMOVED***
			fmt.Println("Can't remove device: ", err)
			os.Exit(1)
		***REMOVED***
	case "mount":
		if flag.NArg() < 3 ***REMOVED***
			usage()
		***REMOVED***

		err := devices.MountDevice(args[1], args[2], "")
		if err != nil ***REMOVED***
			fmt.Println("Can't mount device: ", err)
			os.Exit(1)
		***REMOVED***
	default:
		fmt.Printf("Unknown command %s\n", args[0])
		usage()

		os.Exit(1)
	***REMOVED***
***REMOVED***
