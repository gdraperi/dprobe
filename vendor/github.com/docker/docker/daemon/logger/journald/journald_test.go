// +build linux

package journald

import (
	"testing"
)

func TestSanitizeKeyMod(t *testing.T) ***REMOVED***
	entries := map[string]string***REMOVED***
		"io.kubernetes.pod.name":      "IO_KUBERNETES_POD_NAME",
		"io?.kubernetes.pod.name":     "IO__KUBERNETES_POD_NAME",
		"?io.kubernetes.pod.name":     "IO_KUBERNETES_POD_NAME",
		"io123.kubernetes.pod.name":   "IO123_KUBERNETES_POD_NAME",
		"_io123.kubernetes.pod.name":  "IO123_KUBERNETES_POD_NAME",
		"__io123_kubernetes.pod.name": "IO123_KUBERNETES_POD_NAME",
	***REMOVED***
	for k, v := range entries ***REMOVED***
		if sanitizeKeyMod(k) != v ***REMOVED***
			t.Fatalf("Failed to sanitize %s, got %s, expected %s", k, sanitizeKeyMod(k), v)
		***REMOVED***
	***REMOVED***
***REMOVED***
