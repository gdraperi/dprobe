package dns

// Holds a bunch of helper functions for dealing with labels.

// SplitDomainName splits a name string into it's labels.
// www.miek.nl. returns []string***REMOVED***"www", "miek", "nl"***REMOVED***
// The root label (.) returns nil. Note that using
// strings.Split(s) will work in most cases, but does not handle
// escaped dots (\.) for instance.
func SplitDomainName(s string) (labels []string) ***REMOVED***
	if len(s) == 0 ***REMOVED***
		return nil
	***REMOVED***
	fqdnEnd := 0 // offset of the final '.' or the length of the name
	idx := Split(s)
	begin := 0
	if s[len(s)-1] == '.' ***REMOVED***
		fqdnEnd = len(s) - 1
	***REMOVED*** else ***REMOVED***
		fqdnEnd = len(s)
	***REMOVED***

	switch len(idx) ***REMOVED***
	case 0:
		return nil
	case 1:
		// no-op
	default:
		end := 0
		for i := 1; i < len(idx); i++ ***REMOVED***
			end = idx[i]
			labels = append(labels, s[begin:end-1])
			begin = end
		***REMOVED***
	***REMOVED***

	labels = append(labels, s[begin:fqdnEnd])
	return labels
***REMOVED***

// CompareDomainName compares the names s1 and s2 and
// returns how many labels they have in common starting from the *right*.
// The comparison stops at the first inequality. The names are not downcased
// before the comparison.
//
// www.miek.nl. and miek.nl. have two labels in common: miek and nl
// www.miek.nl. and www.bla.nl. have one label in common: nl
func CompareDomainName(s1, s2 string) (n int) ***REMOVED***
	s1 = Fqdn(s1)
	s2 = Fqdn(s2)
	l1 := Split(s1)
	l2 := Split(s2)

	// the first check: root label
	if l1 == nil || l2 == nil ***REMOVED***
		return
	***REMOVED***

	j1 := len(l1) - 1 // end
	i1 := len(l1) - 2 // start
	j2 := len(l2) - 1
	i2 := len(l2) - 2
	// the second check can be done here: last/only label
	// before we fall through into the for-loop below
	if s1[l1[j1]:] == s2[l2[j2]:] ***REMOVED***
		n++
	***REMOVED*** else ***REMOVED***
		return
	***REMOVED***
	for ***REMOVED***
		if i1 < 0 || i2 < 0 ***REMOVED***
			break
		***REMOVED***
		if s1[l1[i1]:l1[j1]] == s2[l2[i2]:l2[j2]] ***REMOVED***
			n++
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
		j1--
		i1--
		j2--
		i2--
	***REMOVED***
	return
***REMOVED***

// CountLabel counts the the number of labels in the string s.
func CountLabel(s string) (labels int) ***REMOVED***
	if s == "." ***REMOVED***
		return
	***REMOVED***
	off := 0
	end := false
	for ***REMOVED***
		off, end = NextLabel(s, off)
		labels++
		if end ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// Split splits a name s into its label indexes.
// www.miek.nl. returns []int***REMOVED***0, 4, 9***REMOVED***, www.miek.nl also returns []int***REMOVED***0, 4, 9***REMOVED***.
// The root name (.) returns nil. Also see SplitDomainName.
func Split(s string) []int ***REMOVED***
	if s == "." ***REMOVED***
		return nil
	***REMOVED***
	idx := make([]int, 1, 3)
	off := 0
	end := false

	for ***REMOVED***
		off, end = NextLabel(s, off)
		if end ***REMOVED***
			return idx
		***REMOVED***
		idx = append(idx, off)
	***REMOVED***
***REMOVED***

// NextLabel returns the index of the start of the next label in the
// string s starting at offset.
// The bool end is true when the end of the string has been reached.
// Also see PrevLabel.
func NextLabel(s string, offset int) (i int, end bool) ***REMOVED***
	quote := false
	for i = offset; i < len(s)-1; i++ ***REMOVED***
		switch s[i] ***REMOVED***
		case '\\':
			quote = !quote
		default:
			quote = false
		case '.':
			if quote ***REMOVED***
				quote = !quote
				continue
			***REMOVED***
			return i + 1, false
		***REMOVED***
	***REMOVED***
	return i + 1, true
***REMOVED***

// PrevLabel returns the index of the label when starting from the right and
// jumping n labels to the left.
// The bool start is true when the start of the string has been overshot.
// Also see NextLabel.
func PrevLabel(s string, n int) (i int, start bool) ***REMOVED***
	if n == 0 ***REMOVED***
		return len(s), false
	***REMOVED***
	lab := Split(s)
	if lab == nil ***REMOVED***
		return 0, true
	***REMOVED***
	if n > len(lab) ***REMOVED***
		return 0, true
	***REMOVED***
	return lab[len(lab)-n], false
***REMOVED***
