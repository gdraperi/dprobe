package volumedrivers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/docker/docker/pkg/plugins"
	"github.com/docker/go-connections/tlsconfig"
)

func TestVolumeRequestError(t *testing.T) ***REMOVED***
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/VolumeDriver.Create", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
		fmt.Fprintln(w, `***REMOVED***"Err": "Cannot create volume"***REMOVED***`)
	***REMOVED***)

	mux.HandleFunc("/VolumeDriver.Remove", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
		fmt.Fprintln(w, `***REMOVED***"Err": "Cannot remove volume"***REMOVED***`)
	***REMOVED***)

	mux.HandleFunc("/VolumeDriver.Mount", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
		fmt.Fprintln(w, `***REMOVED***"Err": "Cannot mount volume"***REMOVED***`)
	***REMOVED***)

	mux.HandleFunc("/VolumeDriver.Unmount", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
		fmt.Fprintln(w, `***REMOVED***"Err": "Cannot unmount volume"***REMOVED***`)
	***REMOVED***)

	mux.HandleFunc("/VolumeDriver.Path", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
		fmt.Fprintln(w, `***REMOVED***"Err": "Unknown volume"***REMOVED***`)
	***REMOVED***)

	mux.HandleFunc("/VolumeDriver.List", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
		fmt.Fprintln(w, `***REMOVED***"Err": "Cannot list volumes"***REMOVED***`)
	***REMOVED***)

	mux.HandleFunc("/VolumeDriver.Get", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
		fmt.Fprintln(w, `***REMOVED***"Err": "Cannot get volume"***REMOVED***`)
	***REMOVED***)

	mux.HandleFunc("/VolumeDriver.Capabilities", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
		http.Error(w, "error", 500)
	***REMOVED***)

	u, _ := url.Parse(server.URL)
	client, err := plugins.NewClient("tcp://"+u.Host, &tlsconfig.Options***REMOVED***InsecureSkipVerify: true***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	driver := volumeDriverProxy***REMOVED***client***REMOVED***

	if err = driver.Create("volume", nil); err == nil ***REMOVED***
		t.Fatal("Expected error, was nil")
	***REMOVED***

	if !strings.Contains(err.Error(), "Cannot create volume") ***REMOVED***
		t.Fatalf("Unexpected error: %v\n", err)
	***REMOVED***

	_, err = driver.Mount("volume", "123")
	if err == nil ***REMOVED***
		t.Fatal("Expected error, was nil")
	***REMOVED***

	if !strings.Contains(err.Error(), "Cannot mount volume") ***REMOVED***
		t.Fatalf("Unexpected error: %v\n", err)
	***REMOVED***

	err = driver.Unmount("volume", "123")
	if err == nil ***REMOVED***
		t.Fatal("Expected error, was nil")
	***REMOVED***

	if !strings.Contains(err.Error(), "Cannot unmount volume") ***REMOVED***
		t.Fatalf("Unexpected error: %v\n", err)
	***REMOVED***

	err = driver.Remove("volume")
	if err == nil ***REMOVED***
		t.Fatal("Expected error, was nil")
	***REMOVED***

	if !strings.Contains(err.Error(), "Cannot remove volume") ***REMOVED***
		t.Fatalf("Unexpected error: %v\n", err)
	***REMOVED***

	_, err = driver.Path("volume")
	if err == nil ***REMOVED***
		t.Fatal("Expected error, was nil")
	***REMOVED***

	if !strings.Contains(err.Error(), "Unknown volume") ***REMOVED***
		t.Fatalf("Unexpected error: %v\n", err)
	***REMOVED***

	_, err = driver.List()
	if err == nil ***REMOVED***
		t.Fatal("Expected error, was nil")
	***REMOVED***
	if !strings.Contains(err.Error(), "Cannot list volumes") ***REMOVED***
		t.Fatalf("Unexpected error: %v\n", err)
	***REMOVED***

	_, err = driver.Get("volume")
	if err == nil ***REMOVED***
		t.Fatal("Expected error, was nil")
	***REMOVED***
	if !strings.Contains(err.Error(), "Cannot get volume") ***REMOVED***
		t.Fatalf("Unexpected error: %v\n", err)
	***REMOVED***

	_, err = driver.Capabilities()
	if err == nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
