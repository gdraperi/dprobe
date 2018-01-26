package metadata

import (
	"strings"

	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/filters"
	"github.com/containerd/containerd/images"
)

func adaptImage(o interface***REMOVED******REMOVED***) filters.Adaptor ***REMOVED***
	obj := o.(images.Image)
	return filters.AdapterFunc(func(fieldpath []string) (string, bool) ***REMOVED***
		if len(fieldpath) == 0 ***REMOVED***
			return "", false
		***REMOVED***

		switch fieldpath[0] ***REMOVED***
		case "name":
			return obj.Name, len(obj.Name) > 0
		case "target":
			if len(fieldpath) < 2 ***REMOVED***
				return "", false
			***REMOVED***

			switch fieldpath[1] ***REMOVED***
			case "digest":
				return obj.Target.Digest.String(), len(obj.Target.Digest) > 0
			case "mediatype":
				return obj.Target.MediaType, len(obj.Target.MediaType) > 0
			***REMOVED***
		case "labels":
			return checkMap(fieldpath[1:], obj.Labels)
			// TODO(stevvooe): Greater/Less than filters would be awesome for
			// size. Let's do it!
		***REMOVED***

		return "", false
	***REMOVED***)
***REMOVED***
func adaptContainer(o interface***REMOVED******REMOVED***) filters.Adaptor ***REMOVED***
	obj := o.(containers.Container)
	return filters.AdapterFunc(func(fieldpath []string) (string, bool) ***REMOVED***
		if len(fieldpath) == 0 ***REMOVED***
			return "", false
		***REMOVED***

		switch fieldpath[0] ***REMOVED***
		case "id":
			return obj.ID, len(obj.ID) > 0
		case "runtime":
			if len(fieldpath) <= 1 ***REMOVED***
				return "", false
			***REMOVED***

			switch fieldpath[1] ***REMOVED***
			case "name":
				return obj.Runtime.Name, len(obj.Runtime.Name) > 0
			default:
				return "", false
			***REMOVED***
		case "image":
			return obj.Image, len(obj.Image) > 0
		case "labels":
			return checkMap(fieldpath[1:], obj.Labels)
		***REMOVED***

		return "", false
	***REMOVED***)
***REMOVED***

func adaptContentInfo(info content.Info) filters.Adaptor ***REMOVED***
	return filters.AdapterFunc(func(fieldpath []string) (string, bool) ***REMOVED***
		if len(fieldpath) == 0 ***REMOVED***
			return "", false
		***REMOVED***

		switch fieldpath[0] ***REMOVED***
		case "digest":
			return info.Digest.String(), true
		case "size":
			// TODO: support size based filtering
		case "labels":
			return checkMap(fieldpath[1:], info.Labels)
		***REMOVED***

		return "", false
	***REMOVED***)
***REMOVED***

func adaptContentStatus(status content.Status) filters.Adaptor ***REMOVED***
	return filters.AdapterFunc(func(fieldpath []string) (string, bool) ***REMOVED***
		if len(fieldpath) == 0 ***REMOVED***
			return "", false
		***REMOVED***
		switch fieldpath[0] ***REMOVED***
		case "ref":
			return status.Ref, true
		***REMOVED***

		return "", false
	***REMOVED***)
***REMOVED***

func checkMap(fieldpath []string, m map[string]string) (string, bool) ***REMOVED***
	if len(m) == 0 ***REMOVED***
		return "", false
	***REMOVED***

	value, ok := m[strings.Join(fieldpath, ".")]
	return value, ok
***REMOVED***
