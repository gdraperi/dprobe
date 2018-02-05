package slack

import "testing"

func TestNewMessageItem(t *testing.T) ***REMOVED***
	c := "C1"
	m := &Message***REMOVED******REMOVED***
	mi := NewMessageItem(c, m)
	if mi.Type != TYPE_MESSAGE ***REMOVED***
		t.Errorf("want Type %s, got %s", mi.Type, TYPE_MESSAGE)
	***REMOVED***
	if mi.Channel != c ***REMOVED***
		t.Errorf("got Channel %s, want %s", mi.Channel, c)
	***REMOVED***
	if mi.Message != m ***REMOVED***
		t.Errorf("got Message %v, want %v", mi.Message, m)
	***REMOVED***
***REMOVED***

func TestNewFileItem(t *testing.T) ***REMOVED***
	f := &File***REMOVED******REMOVED***
	fi := NewFileItem(f)
	if fi.Type != TYPE_FILE ***REMOVED***
		t.Errorf("got Type %s, want %s", fi.Type, TYPE_FILE)
	***REMOVED***
	if fi.File != f ***REMOVED***
		t.Errorf("got File %v, want %v", fi.File, f)
	***REMOVED***
***REMOVED***

func TestNewFileCommentItem(t *testing.T) ***REMOVED***
	f := &File***REMOVED******REMOVED***
	c := &Comment***REMOVED******REMOVED***
	fci := NewFileCommentItem(f, c)
	if fci.Type != TYPE_FILE_COMMENT ***REMOVED***
		t.Errorf("got Type %s, want %s", fci.Type, TYPE_FILE_COMMENT)
	***REMOVED***
	if fci.File != f ***REMOVED***
		t.Errorf("got File %v, want %v", fci.File, f)
	***REMOVED***
	if fci.Comment != c ***REMOVED***
		t.Errorf("got Comment %v, want %v", fci.Comment, c)
	***REMOVED***
***REMOVED***

func TestNewChannelItem(t *testing.T) ***REMOVED***
	c := "C1"
	ci := NewChannelItem(c)
	if ci.Type != TYPE_CHANNEL ***REMOVED***
		t.Errorf("got Type %s, want %s", ci.Type, TYPE_CHANNEL)
	***REMOVED***
	if ci.Channel != "C1" ***REMOVED***
		t.Errorf("got Channel %v, want %v", ci.Channel, "C1")
	***REMOVED***
***REMOVED***

func TestNewIMItem(t *testing.T) ***REMOVED***
	c := "D1"
	ci := NewIMItem(c)
	if ci.Type != TYPE_IM ***REMOVED***
		t.Errorf("got Type %s, want %s", ci.Type, TYPE_IM)
	***REMOVED***
	if ci.Channel != "D1" ***REMOVED***
		t.Errorf("got Channel %v, want %v", ci.Channel, "D1")
	***REMOVED***
***REMOVED***

func TestNewGroupItem(t *testing.T) ***REMOVED***
	c := "G1"
	ci := NewGroupItem(c)
	if ci.Type != TYPE_GROUP ***REMOVED***
		t.Errorf("got Type %s, want %s", ci.Type, TYPE_GROUP)
	***REMOVED***
	if ci.Channel != "G1" ***REMOVED***
		t.Errorf("got Channel %v, want %v", ci.Channel, "G1")
	***REMOVED***
***REMOVED***

func TestNewRefToMessage(t *testing.T) ***REMOVED***
	ref := NewRefToMessage("chan", "ts")
	if got, want := ref.Channel, "chan"; got != want ***REMOVED***
		t.Errorf("Channel got %s, want %s", got, want)
	***REMOVED***
	if got, want := ref.Timestamp, "ts"; got != want ***REMOVED***
		t.Errorf("Timestamp got %s, want %s", got, want)
	***REMOVED***
	if got, want := ref.File, ""; got != want ***REMOVED***
		t.Errorf("File got %s, want %s", got, want)
	***REMOVED***
	if got, want := ref.Comment, ""; got != want ***REMOVED***
		t.Errorf("Comment got %s, want %s", got, want)
	***REMOVED***
***REMOVED***

func TestNewRefToFile(t *testing.T) ***REMOVED***
	ref := NewRefToFile("file")
	if got, want := ref.Channel, ""; got != want ***REMOVED***
		t.Errorf("Channel got %s, want %s", got, want)
	***REMOVED***
	if got, want := ref.Timestamp, ""; got != want ***REMOVED***
		t.Errorf("Timestamp got %s, want %s", got, want)
	***REMOVED***
	if got, want := ref.File, "file"; got != want ***REMOVED***
		t.Errorf("File got %s, want %s", got, want)
	***REMOVED***
	if got, want := ref.Comment, ""; got != want ***REMOVED***
		t.Errorf("Comment got %s, want %s", got, want)
	***REMOVED***
***REMOVED***

func TestNewRefToComment(t *testing.T) ***REMOVED***
	ref := NewRefToComment("file_comment")
	if got, want := ref.Channel, ""; got != want ***REMOVED***
		t.Errorf("Channel got %s, want %s", got, want)
	***REMOVED***
	if got, want := ref.Timestamp, ""; got != want ***REMOVED***
		t.Errorf("Timestamp got %s, want %s", got, want)
	***REMOVED***
	if got, want := ref.File, ""; got != want ***REMOVED***
		t.Errorf("File got %s, want %s", got, want)
	***REMOVED***
	if got, want := ref.Comment, "file_comment"; got != want ***REMOVED***
		t.Errorf("Comment got %s, want %s", got, want)
	***REMOVED***
***REMOVED***
