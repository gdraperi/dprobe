package etchosts

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"
)

// Record Structure for a single host record
type Record struct ***REMOVED***
	Hosts string
	IP    string
***REMOVED***

// WriteTo writes record to file and returns bytes written or error
func (r Record) WriteTo(w io.Writer) (int64, error) ***REMOVED***
	n, err := fmt.Fprintf(w, "%s\t%s\n", r.IP, r.Hosts)
	return int64(n), err
***REMOVED***

var (
	// Default hosts config records slice
	defaultContent = []Record***REMOVED***
		***REMOVED***Hosts: "localhost", IP: "127.0.0.1"***REMOVED***,
		***REMOVED***Hosts: "localhost ip6-localhost ip6-loopback", IP: "::1"***REMOVED***,
		***REMOVED***Hosts: "ip6-localnet", IP: "fe00::0"***REMOVED***,
		***REMOVED***Hosts: "ip6-mcastprefix", IP: "ff00::0"***REMOVED***,
		***REMOVED***Hosts: "ip6-allnodes", IP: "ff02::1"***REMOVED***,
		***REMOVED***Hosts: "ip6-allrouters", IP: "ff02::2"***REMOVED***,
	***REMOVED***

	// A cache of path level locks for synchronizing /etc/hosts
	// updates on a file level
	pathMap = make(map[string]*sync.Mutex)

	// A package level mutex to synchronize the cache itself
	pathMutex sync.Mutex
)

func pathLock(path string) func() ***REMOVED***
	pathMutex.Lock()
	defer pathMutex.Unlock()

	pl, ok := pathMap[path]
	if !ok ***REMOVED***
		pl = &sync.Mutex***REMOVED******REMOVED***
		pathMap[path] = pl
	***REMOVED***

	pl.Lock()
	return func() ***REMOVED***
		pl.Unlock()
	***REMOVED***
***REMOVED***

// Drop drops the path string from the path cache
func Drop(path string) ***REMOVED***
	pathMutex.Lock()
	defer pathMutex.Unlock()

	delete(pathMap, path)
***REMOVED***

// Build function
// path is path to host file string required
// IP, hostname, and domainname set main record leave empty for no master record
// extraContent is an array of extra host records.
func Build(path, IP, hostname, domainname string, extraContent []Record) error ***REMOVED***
	defer pathLock(path)()

	content := bytes.NewBuffer(nil)
	if IP != "" ***REMOVED***
		//set main record
		var mainRec Record
		mainRec.IP = IP
		// User might have provided a FQDN in hostname or split it across hostname
		// and domainname.  We want the FQDN and the bare hostname.
		fqdn := hostname
		if domainname != "" ***REMOVED***
			fqdn = fmt.Sprintf("%s.%s", fqdn, domainname)
		***REMOVED***
		parts := strings.SplitN(fqdn, ".", 2)
		if len(parts) == 2 ***REMOVED***
			mainRec.Hosts = fmt.Sprintf("%s %s", fqdn, parts[0])
		***REMOVED*** else ***REMOVED***
			mainRec.Hosts = fqdn
		***REMOVED***
		if _, err := mainRec.WriteTo(content); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	// Write defaultContent slice to buffer
	for _, r := range defaultContent ***REMOVED***
		if _, err := r.WriteTo(content); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	// Write extra content from function arguments
	for _, r := range extraContent ***REMOVED***
		if _, err := r.WriteTo(content); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return ioutil.WriteFile(path, content.Bytes(), 0644)
***REMOVED***

// Add adds an arbitrary number of Records to an already existing /etc/hosts file
func Add(path string, recs []Record) error ***REMOVED***
	defer pathLock(path)()

	if len(recs) == 0 ***REMOVED***
		return nil
	***REMOVED***

	b, err := mergeRecords(path, recs)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return ioutil.WriteFile(path, b, 0644)
***REMOVED***

func mergeRecords(path string, recs []Record) ([]byte, error) ***REMOVED***
	f, err := os.Open(path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer f.Close()

	content := bytes.NewBuffer(nil)

	if _, err := content.ReadFrom(f); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	for _, r := range recs ***REMOVED***
		if _, err := r.WriteTo(content); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return content.Bytes(), nil
***REMOVED***

// Delete deletes an arbitrary number of Records already existing in /etc/hosts file
func Delete(path string, recs []Record) error ***REMOVED***
	defer pathLock(path)()

	if len(recs) == 0 ***REMOVED***
		return nil
	***REMOVED***
	old, err := os.Open(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var buf bytes.Buffer

	s := bufio.NewScanner(old)
	eol := []byte***REMOVED***'\n'***REMOVED***
loop:
	for s.Scan() ***REMOVED***
		b := s.Bytes()
		if len(b) == 0 ***REMOVED***
			continue
		***REMOVED***

		if b[0] == '#' ***REMOVED***
			buf.Write(b)
			buf.Write(eol)
			continue
		***REMOVED***
		for _, r := range recs ***REMOVED***
			if bytes.HasSuffix(b, []byte("\t"+r.Hosts)) ***REMOVED***
				continue loop
			***REMOVED***
		***REMOVED***
		buf.Write(b)
		buf.Write(eol)
	***REMOVED***
	old.Close()
	if err := s.Err(); err != nil ***REMOVED***
		return err
	***REMOVED***
	return ioutil.WriteFile(path, buf.Bytes(), 0644)
***REMOVED***

// Update all IP addresses where hostname matches.
// path is path to host file
// IP is new IP address
// hostname is hostname to search for to replace IP
func Update(path, IP, hostname string) error ***REMOVED***
	defer pathLock(path)()

	old, err := ioutil.ReadFile(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	var re = regexp.MustCompile(fmt.Sprintf("(\\S*)(\\t%s)(\\s|\\.)", regexp.QuoteMeta(hostname)))
	return ioutil.WriteFile(path, re.ReplaceAll(old, []byte(IP+"$2"+"$3")), 0644)
***REMOVED***
