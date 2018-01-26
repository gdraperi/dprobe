package container

import (
	"testing"
	"time"
)

func TestNewMemoryStore(t *testing.T) ***REMOVED***
	s := NewMemoryStore()
	m, ok := s.(*memoryStore)
	if !ok ***REMOVED***
		t.Fatalf("store is not a memory store %v", s)
	***REMOVED***
	if m.s == nil ***REMOVED***
		t.Fatal("expected store map to not be nil")
	***REMOVED***
***REMOVED***

func TestAddContainers(t *testing.T) ***REMOVED***
	s := NewMemoryStore()
	s.Add("id", NewBaseContainer("id", "root"))
	if s.Size() != 1 ***REMOVED***
		t.Fatalf("expected store size 1, got %v", s.Size())
	***REMOVED***
***REMOVED***

func TestGetContainer(t *testing.T) ***REMOVED***
	s := NewMemoryStore()
	s.Add("id", NewBaseContainer("id", "root"))
	c := s.Get("id")
	if c == nil ***REMOVED***
		t.Fatal("expected container to not be nil")
	***REMOVED***
***REMOVED***

func TestDeleteContainer(t *testing.T) ***REMOVED***
	s := NewMemoryStore()
	s.Add("id", NewBaseContainer("id", "root"))
	s.Delete("id")
	if c := s.Get("id"); c != nil ***REMOVED***
		t.Fatalf("expected container to be nil after removal, got %v", c)
	***REMOVED***

	if s.Size() != 0 ***REMOVED***
		t.Fatalf("expected store size to be 0, got %v", s.Size())
	***REMOVED***
***REMOVED***

func TestListContainers(t *testing.T) ***REMOVED***
	s := NewMemoryStore()

	cont := NewBaseContainer("id", "root")
	cont.Created = time.Now()
	cont2 := NewBaseContainer("id2", "root")
	cont2.Created = time.Now().Add(24 * time.Hour)

	s.Add("id", cont)
	s.Add("id2", cont2)

	list := s.List()
	if len(list) != 2 ***REMOVED***
		t.Fatalf("expected list size 2, got %v", len(list))
	***REMOVED***
	if list[0].ID != "id2" ***REMOVED***
		t.Fatalf("expected id2, got %v", list[0].ID)
	***REMOVED***
***REMOVED***

func TestFirstContainer(t *testing.T) ***REMOVED***
	s := NewMemoryStore()

	s.Add("id", NewBaseContainer("id", "root"))
	s.Add("id2", NewBaseContainer("id2", "root"))

	first := s.First(func(cont *Container) bool ***REMOVED***
		return cont.ID == "id2"
	***REMOVED***)

	if first == nil ***REMOVED***
		t.Fatal("expected container to not be nil")
	***REMOVED***
	if first.ID != "id2" ***REMOVED***
		t.Fatalf("expected id2, got %v", first)
	***REMOVED***
***REMOVED***

func TestApplyAllContainer(t *testing.T) ***REMOVED***
	s := NewMemoryStore()

	s.Add("id", NewBaseContainer("id", "root"))
	s.Add("id2", NewBaseContainer("id2", "root"))

	s.ApplyAll(func(cont *Container) ***REMOVED***
		if cont.ID == "id2" ***REMOVED***
			cont.ID = "newID"
		***REMOVED***
	***REMOVED***)

	cont := s.Get("id2")
	if cont == nil ***REMOVED***
		t.Fatal("expected container to not be nil")
	***REMOVED***
	if cont.ID != "newID" ***REMOVED***
		t.Fatalf("expected newID, got %v", cont.ID)
	***REMOVED***
***REMOVED***
