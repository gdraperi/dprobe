package cluster

import (
	"testing"

	"github.com/docker/docker/api/types/filters"
)

func TestNewListSecretsFilters(t *testing.T) ***REMOVED***
	validNameFilter := filters.NewArgs()
	validNameFilter.Add("name", "test_name")

	validIDFilter := filters.NewArgs()
	validIDFilter.Add("id", "7c9009d6720f6de3b492f5")

	validLabelFilter := filters.NewArgs()
	validLabelFilter.Add("label", "type=test")
	validLabelFilter.Add("label", "storage=ssd")
	validLabelFilter.Add("label", "memory")

	validNamesFilter := filters.NewArgs()
	validNamesFilter.Add("names", "test_name")

	validAllFilter := filters.NewArgs()
	validAllFilter.Add("name", "nodeName")
	validAllFilter.Add("id", "7c9009d6720f6de3b492f5")
	validAllFilter.Add("label", "type=test")
	validAllFilter.Add("label", "memory")
	validAllFilter.Add("names", "test_name")

	validFilters := []filters.Args***REMOVED***
		validNameFilter,
		validIDFilter,
		validLabelFilter,
		validNamesFilter,
		validAllFilter,
	***REMOVED***

	invalidTypeFilter := filters.NewArgs()
	invalidTypeFilter.Add("nonexist", "aaaa")

	invalidFilters := []filters.Args***REMOVED***
		invalidTypeFilter,
	***REMOVED***

	for _, filter := range validFilters ***REMOVED***
		if _, err := newListSecretsFilters(filter); err != nil ***REMOVED***
			t.Fatalf("Should get no error, got %v", err)
		***REMOVED***
	***REMOVED***

	for _, filter := range invalidFilters ***REMOVED***
		if _, err := newListSecretsFilters(filter); err == nil ***REMOVED***
			t.Fatalf("Should get an error for filter %v, while got nil", filter)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestNewListConfigsFilters(t *testing.T) ***REMOVED***
	validNameFilter := filters.NewArgs()
	validNameFilter.Add("name", "test_name")

	validIDFilter := filters.NewArgs()
	validIDFilter.Add("id", "7c9009d6720f6de3b492f5")

	validLabelFilter := filters.NewArgs()
	validLabelFilter.Add("label", "type=test")
	validLabelFilter.Add("label", "storage=ssd")
	validLabelFilter.Add("label", "memory")

	validAllFilter := filters.NewArgs()
	validAllFilter.Add("name", "nodeName")
	validAllFilter.Add("id", "7c9009d6720f6de3b492f5")
	validAllFilter.Add("label", "type=test")
	validAllFilter.Add("label", "memory")

	validFilters := []filters.Args***REMOVED***
		validNameFilter,
		validIDFilter,
		validLabelFilter,
		validAllFilter,
	***REMOVED***

	invalidTypeFilter := filters.NewArgs()
	invalidTypeFilter.Add("nonexist", "aaaa")

	invalidFilters := []filters.Args***REMOVED***
		invalidTypeFilter,
	***REMOVED***

	for _, filter := range validFilters ***REMOVED***
		if _, err := newListConfigsFilters(filter); err != nil ***REMOVED***
			t.Fatalf("Should get no error, got %v", err)
		***REMOVED***
	***REMOVED***

	for _, filter := range invalidFilters ***REMOVED***
		if _, err := newListConfigsFilters(filter); err == nil ***REMOVED***
			t.Fatalf("Should get an error for filter %v, while got nil", filter)
		***REMOVED***
	***REMOVED***
***REMOVED***
