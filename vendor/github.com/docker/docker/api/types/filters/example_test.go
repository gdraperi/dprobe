package filters

func ExampleArgs_MatchKVList() ***REMOVED***
	args := NewArgs(
		Arg("label", "image=foo"),
		Arg("label", "state=running"))

	// returns true because there are no values for bogus
	args.MatchKVList("bogus", nil)

	// returns false because there are no sources
	args.MatchKVList("label", nil)

	// returns true because all sources are matched
	args.MatchKVList("label", map[string]string***REMOVED***
		"image": "foo",
		"state": "running",
	***REMOVED***)

	// returns false because the values do not match
	args.MatchKVList("label", map[string]string***REMOVED***
		"image": "other",
	***REMOVED***)
***REMOVED***
