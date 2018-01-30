// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
)

var (
	proj     = flag.String("project", "symbolic-datum-552", "name of Project")
	zone     = flag.String("zone", "us-central1-a", "GCE zone")
	mach     = flag.String("machinetype", "n1-standard-1", "Machine type")
	instName = flag.String("instance_name", "http2-demo", "Name of VM instance.")
	sshPub   = flag.String("ssh_public_key", "", "ssh public key file to authorize. Can modify later in Google's web UI anyway.")
	staticIP = flag.String("static_ip", "130.211.116.44", "Static IP to use. If empty, automatic.")

	writeObject  = flag.String("write_object", "", "If non-empty, a VM isn't created and the flag value is Google Cloud Storage bucket/object to write. The contents from stdin.")
	publicObject = flag.Bool("write_object_is_public", false, "Whether the object created by --write_object should be public.")
)

func readFile(v string) string ***REMOVED***
	slurp, err := ioutil.ReadFile(v)
	if err != nil ***REMOVED***
		log.Fatalf("Error reading %s: %v", v, err)
	***REMOVED***
	return strings.TrimSpace(string(slurp))
***REMOVED***

var config = &oauth2.Config***REMOVED***
	// The client-id and secret should be for an "Installed Application" when using
	// the CLI. Later we'll use a web application with a callback.
	ClientID:     readFile("client-id.dat"),
	ClientSecret: readFile("client-secret.dat"),
	Endpoint:     google.Endpoint,
	Scopes: []string***REMOVED***
		compute.DevstorageFullControlScope,
		compute.ComputeScope,
		"https://www.googleapis.com/auth/sqlservice",
		"https://www.googleapis.com/auth/sqlservice.admin",
	***REMOVED***,
	RedirectURL: "urn:ietf:wg:oauth:2.0:oob",
***REMOVED***

const baseConfig = `#cloud-config
coreos:
  units:
    - name: h2demo.service
      command: start
      content: |
        [Unit]
        Description=HTTP2 Demo
        
        [Service]
        ExecStartPre=/bin/bash -c 'mkdir -p /opt/bin && curl -s -o /opt/bin/h2demo http://storage.googleapis.com/http2-demo-server-tls/h2demo && chmod +x /opt/bin/h2demo'
        ExecStart=/opt/bin/h2demo --prod
        RestartSec=5s
        Restart=always
        Type=simple
        
        [Install]
        WantedBy=multi-user.target
`

func main() ***REMOVED***
	flag.Parse()
	if *proj == "" ***REMOVED***
		log.Fatalf("Missing --project flag")
	***REMOVED***
	prefix := "https://www.googleapis.com/compute/v1/projects/" + *proj
	machType := prefix + "/zones/" + *zone + "/machineTypes/" + *mach

	const tokenFileName = "token.dat"
	tokenFile := tokenCacheFile(tokenFileName)
	tokenSource := oauth2.ReuseTokenSource(nil, tokenFile)
	token, err := tokenSource.Token()
	if err != nil ***REMOVED***
		if *writeObject != "" ***REMOVED***
			log.Fatalf("Can't use --write_object without a valid token.dat file already cached.")
		***REMOVED***
		log.Printf("Error getting token from %s: %v", tokenFileName, err)
		log.Printf("Get auth code from %v", config.AuthCodeURL("my-state"))
		fmt.Print("\nEnter auth code: ")
		sc := bufio.NewScanner(os.Stdin)
		sc.Scan()
		authCode := strings.TrimSpace(sc.Text())
		token, err = config.Exchange(oauth2.NoContext, authCode)
		if err != nil ***REMOVED***
			log.Fatalf("Error exchanging auth code for a token: %v", err)
		***REMOVED***
		if err := tokenFile.WriteToken(token); err != nil ***REMOVED***
			log.Fatalf("Error writing to %s: %v", tokenFileName, err)
		***REMOVED***
		tokenSource = oauth2.ReuseTokenSource(token, nil)
	***REMOVED***

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)

	if *writeObject != "" ***REMOVED***
		writeCloudStorageObject(oauthClient)
		return
	***REMOVED***

	computeService, _ := compute.New(oauthClient)

	natIP := *staticIP
	if natIP == "" ***REMOVED***
		// Try to find it by name.
		aggAddrList, err := computeService.Addresses.AggregatedList(*proj).Do()
		if err != nil ***REMOVED***
			log.Fatal(err)
		***REMOVED***
		// http://godoc.org/code.google.com/p/google-api-go-client/compute/v1#AddressAggregatedList
	IPLoop:
		for _, asl := range aggAddrList.Items ***REMOVED***
			for _, addr := range asl.Addresses ***REMOVED***
				if addr.Name == *instName+"-ip" && addr.Status == "RESERVED" ***REMOVED***
					natIP = addr.Address
					break IPLoop
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	cloudConfig := baseConfig
	if *sshPub != "" ***REMOVED***
		key := strings.TrimSpace(readFile(*sshPub))
		cloudConfig += fmt.Sprintf("\nssh_authorized_keys:\n    - %s\n", key)
	***REMOVED***
	if os.Getenv("USER") == "bradfitz" ***REMOVED***
		cloudConfig += fmt.Sprintf("\nssh_authorized_keys:\n    - %s\n", "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAIEAwks9dwWKlRC+73gRbvYtVg0vdCwDSuIlyt4z6xa/YU/jTDynM4R4W10hm2tPjy8iR1k8XhDv4/qdxe6m07NjG/By1tkmGpm1mGwho4Pr5kbAAy/Qg+NLCSdAYnnE00FQEcFOC15GFVMOW2AzDGKisReohwH9eIzHPzdYQNPRWXE= bradfitz@papag.bradfitz.com")
	***REMOVED***
	const maxCloudConfig = 32 << 10 // per compute API docs
	if len(cloudConfig) > maxCloudConfig ***REMOVED***
		log.Fatalf("cloud config length of %d bytes is over %d byte limit", len(cloudConfig), maxCloudConfig)
	***REMOVED***

	instance := &compute.Instance***REMOVED***
		Name:        *instName,
		Description: "Go Builder",
		MachineType: machType,
		Disks:       []*compute.AttachedDisk***REMOVED***instanceDisk(computeService)***REMOVED***,
		Tags: &compute.Tags***REMOVED***
			Items: []string***REMOVED***"http-server", "https-server"***REMOVED***,
		***REMOVED***,
		Metadata: &compute.Metadata***REMOVED***
			Items: []*compute.MetadataItems***REMOVED***
				***REMOVED***
					Key:   "user-data",
					Value: &cloudConfig,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		NetworkInterfaces: []*compute.NetworkInterface***REMOVED***
			***REMOVED***
				AccessConfigs: []*compute.AccessConfig***REMOVED***
					***REMOVED***
						Type:  "ONE_TO_ONE_NAT",
						Name:  "External NAT",
						NatIP: natIP,
					***REMOVED***,
				***REMOVED***,
				Network: prefix + "/global/networks/default",
			***REMOVED***,
		***REMOVED***,
		ServiceAccounts: []*compute.ServiceAccount***REMOVED***
			***REMOVED***
				Email: "default",
				Scopes: []string***REMOVED***
					compute.DevstorageFullControlScope,
					compute.ComputeScope,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	log.Printf("Creating instance...")
	op, err := computeService.Instances.Insert(*proj, *zone, instance).Do()
	if err != nil ***REMOVED***
		log.Fatalf("Failed to create instance: %v", err)
	***REMOVED***
	opName := op.Name
	log.Printf("Created. Waiting on operation %v", opName)
OpLoop:
	for ***REMOVED***
		time.Sleep(2 * time.Second)
		op, err := computeService.ZoneOperations.Get(*proj, *zone, opName).Do()
		if err != nil ***REMOVED***
			log.Fatalf("Failed to get op %s: %v", opName, err)
		***REMOVED***
		switch op.Status ***REMOVED***
		case "PENDING", "RUNNING":
			log.Printf("Waiting on operation %v", opName)
			continue
		case "DONE":
			if op.Error != nil ***REMOVED***
				for _, operr := range op.Error.Errors ***REMOVED***
					log.Printf("Error: %+v", operr)
				***REMOVED***
				log.Fatalf("Failed to start.")
			***REMOVED***
			log.Printf("Success. %+v", op)
			break OpLoop
		default:
			log.Fatalf("Unknown status %q: %+v", op.Status, op)
		***REMOVED***
	***REMOVED***

	inst, err := computeService.Instances.Get(*proj, *zone, *instName).Do()
	if err != nil ***REMOVED***
		log.Fatalf("Error getting instance after creation: %v", err)
	***REMOVED***
	ij, _ := json.MarshalIndent(inst, "", "    ")
	log.Printf("Instance: %s", ij)
***REMOVED***

func instanceDisk(svc *compute.Service) *compute.AttachedDisk ***REMOVED***
	const imageURL = "https://www.googleapis.com/compute/v1/projects/coreos-cloud/global/images/coreos-stable-444-5-0-v20141016"
	diskName := *instName + "-disk"

	return &compute.AttachedDisk***REMOVED***
		AutoDelete: true,
		Boot:       true,
		Type:       "PERSISTENT",
		InitializeParams: &compute.AttachedDiskInitializeParams***REMOVED***
			DiskName:    diskName,
			SourceImage: imageURL,
			DiskSizeGb:  50,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func writeCloudStorageObject(httpClient *http.Client) ***REMOVED***
	content := os.Stdin
	const maxSlurp = 1 << 20
	var buf bytes.Buffer
	n, err := io.CopyN(&buf, content, maxSlurp)
	if err != nil && err != io.EOF ***REMOVED***
		log.Fatalf("Error reading from stdin: %v, %v", n, err)
	***REMOVED***
	contentType := http.DetectContentType(buf.Bytes())

	req, err := http.NewRequest("PUT", "https://storage.googleapis.com/"+*writeObject, io.MultiReader(&buf, content))
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	req.Header.Set("x-goog-api-version", "2")
	if *publicObject ***REMOVED***
		req.Header.Set("x-goog-acl", "public-read")
	***REMOVED***
	req.Header.Set("Content-Type", contentType)
	res, err := httpClient.Do(req)
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***
	if res.StatusCode != 200 ***REMOVED***
		res.Write(os.Stderr)
		log.Fatalf("Failed.")
	***REMOVED***
	log.Printf("Success.")
	os.Exit(0)
***REMOVED***

type tokenCacheFile string

func (f tokenCacheFile) Token() (*oauth2.Token, error) ***REMOVED***
	slurp, err := ioutil.ReadFile(string(f))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	t := new(oauth2.Token)
	if err := json.Unmarshal(slurp, t); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return t, nil
***REMOVED***

func (f tokenCacheFile) WriteToken(t *oauth2.Token) error ***REMOVED***
	jt, err := json.Marshal(t)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return ioutil.WriteFile(string(f), jt, 0600)
***REMOVED***
