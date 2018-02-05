package printer

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/token"
)

const (
	blank    = byte(' ')
	newline  = byte('\n')
	tab      = byte('\t')
	infinity = 1 << 30 // offset or line
)

var (
	unindent = []byte("\uE123") // in the private use space
)

type printer struct ***REMOVED***
	cfg  Config
	prev token.Pos

	comments           []*ast.CommentGroup // may be nil, contains all comments
	standaloneComments []*ast.CommentGroup // contains all standalone comments (not assigned to any node)

	enableTrace bool
	indentTrace int
***REMOVED***

type ByPosition []*ast.CommentGroup

func (b ByPosition) Len() int           ***REMOVED*** return len(b) ***REMOVED***
func (b ByPosition) Swap(i, j int)      ***REMOVED*** b[i], b[j] = b[j], b[i] ***REMOVED***
func (b ByPosition) Less(i, j int) bool ***REMOVED*** return b[i].Pos().Before(b[j].Pos()) ***REMOVED***

// collectComments comments all standalone comments which are not lead or line
// comment
func (p *printer) collectComments(node ast.Node) ***REMOVED***
	// first collect all comments. This is already stored in
	// ast.File.(comments)
	ast.Walk(node, func(nn ast.Node) (ast.Node, bool) ***REMOVED***
		switch t := nn.(type) ***REMOVED***
		case *ast.File:
			p.comments = t.Comments
			return nn, false
		***REMOVED***
		return nn, true
	***REMOVED***)

	standaloneComments := make(map[token.Pos]*ast.CommentGroup, 0)
	for _, c := range p.comments ***REMOVED***
		standaloneComments[c.Pos()] = c
	***REMOVED***

	// next remove all lead and line comments from the overall comment map.
	// This will give us comments which are standalone, comments which are not
	// assigned to any kind of node.
	ast.Walk(node, func(nn ast.Node) (ast.Node, bool) ***REMOVED***
		switch t := nn.(type) ***REMOVED***
		case *ast.LiteralType:
			if t.LeadComment != nil ***REMOVED***
				for _, comment := range t.LeadComment.List ***REMOVED***
					if _, ok := standaloneComments[comment.Pos()]; ok ***REMOVED***
						delete(standaloneComments, comment.Pos())
					***REMOVED***
				***REMOVED***
			***REMOVED***

			if t.LineComment != nil ***REMOVED***
				for _, comment := range t.LineComment.List ***REMOVED***
					if _, ok := standaloneComments[comment.Pos()]; ok ***REMOVED***
						delete(standaloneComments, comment.Pos())
					***REMOVED***
				***REMOVED***
			***REMOVED***
		case *ast.ObjectItem:
			if t.LeadComment != nil ***REMOVED***
				for _, comment := range t.LeadComment.List ***REMOVED***
					if _, ok := standaloneComments[comment.Pos()]; ok ***REMOVED***
						delete(standaloneComments, comment.Pos())
					***REMOVED***
				***REMOVED***
			***REMOVED***

			if t.LineComment != nil ***REMOVED***
				for _, comment := range t.LineComment.List ***REMOVED***
					if _, ok := standaloneComments[comment.Pos()]; ok ***REMOVED***
						delete(standaloneComments, comment.Pos())
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***

		return nn, true
	***REMOVED***)

	for _, c := range standaloneComments ***REMOVED***
		p.standaloneComments = append(p.standaloneComments, c)
	***REMOVED***

	sort.Sort(ByPosition(p.standaloneComments))
***REMOVED***

// output prints creates b printable HCL output and returns it.
func (p *printer) output(n interface***REMOVED******REMOVED***) []byte ***REMOVED***
	var buf bytes.Buffer

	switch t := n.(type) ***REMOVED***
	case *ast.File:
		// File doesn't trace so we add the tracing here
		defer un(trace(p, "File"))
		return p.output(t.Node)
	case *ast.ObjectList:
		defer un(trace(p, "ObjectList"))

		var index int
		for ***REMOVED***
			// Determine the location of the next actual non-comment
			// item. If we're at the end, the next item is at "infinity"
			var nextItem token.Pos
			if index != len(t.Items) ***REMOVED***
				nextItem = t.Items[index].Pos()
			***REMOVED*** else ***REMOVED***
				nextItem = token.Pos***REMOVED***Offset: infinity, Line: infinity***REMOVED***
			***REMOVED***

			// Go through the standalone comments in the file and print out
			// the comments that we should be for this object item.
			for _, c := range p.standaloneComments ***REMOVED***
				// Go through all the comments in the group. The group
				// should be printed together, not separated by double newlines.
				printed := false
				newlinePrinted := false
				for _, comment := range c.List ***REMOVED***
					// We only care about comments after the previous item
					// we've printed so that comments are printed in the
					// correct locations (between two objects for example).
					// And before the next item.
					if comment.Pos().After(p.prev) && comment.Pos().Before(nextItem) ***REMOVED***
						// if we hit the end add newlines so we can print the comment
						// we don't do this if prev is invalid which means the
						// beginning of the file since the first comment should
						// be at the first line.
						if !newlinePrinted && p.prev.IsValid() && index == len(t.Items) ***REMOVED***
							buf.Write([]byte***REMOVED***newline, newline***REMOVED***)
							newlinePrinted = true
						***REMOVED***

						// Write the actual comment.
						buf.WriteString(comment.Text)
						buf.WriteByte(newline)

						// Set printed to true to note that we printed something
						printed = true
					***REMOVED***
				***REMOVED***

				// If we're not at the last item, write a new line so
				// that there is a newline separating this comment from
				// the next object.
				if printed && index != len(t.Items) ***REMOVED***
					buf.WriteByte(newline)
				***REMOVED***
			***REMOVED***

			if index == len(t.Items) ***REMOVED***
				break
			***REMOVED***

			buf.Write(p.output(t.Items[index]))
			if index != len(t.Items)-1 ***REMOVED***
				// Always write a newline to separate us from the next item
				buf.WriteByte(newline)

				// Need to determine if we're going to separate the next item
				// with a blank line. The logic here is simple, though there
				// are a few conditions:
				//
				//   1. The next object is more than one line away anyways,
				//      so we need an empty line.
				//
				//   2. The next object is not a "single line" object, so
				//      we need an empty line.
				//
				//   3. This current object is not a single line object,
				//      so we need an empty line.
				current := t.Items[index]
				next := t.Items[index+1]
				if next.Pos().Line != t.Items[index].Pos().Line+1 ||
					!p.isSingleLineObject(next) ||
					!p.isSingleLineObject(current) ***REMOVED***
					buf.WriteByte(newline)
				***REMOVED***
			***REMOVED***
			index++
		***REMOVED***
	case *ast.ObjectKey:
		buf.WriteString(t.Token.Text)
	case *ast.ObjectItem:
		p.prev = t.Pos()
		buf.Write(p.objectItem(t))
	case *ast.LiteralType:
		buf.Write(p.literalType(t))
	case *ast.ListType:
		buf.Write(p.list(t))
	case *ast.ObjectType:
		buf.Write(p.objectType(t))
	default:
		fmt.Printf(" unknown type: %T\n", n)
	***REMOVED***

	return buf.Bytes()
***REMOVED***

func (p *printer) literalType(lit *ast.LiteralType) []byte ***REMOVED***
	result := []byte(lit.Token.Text)
	switch lit.Token.Type ***REMOVED***
	case token.HEREDOC:
		// Clear the trailing newline from heredocs
		if result[len(result)-1] == '\n' ***REMOVED***
			result = result[:len(result)-1]
		***REMOVED***

		// Poison lines 2+ so that we don't indent them
		result = p.heredocIndent(result)
	case token.STRING:
		// If this is a multiline string, poison lines 2+ so we don't
		// indent them.
		if bytes.IndexRune(result, '\n') >= 0 ***REMOVED***
			result = p.heredocIndent(result)
		***REMOVED***
	***REMOVED***

	return result
***REMOVED***

// objectItem returns the printable HCL form of an object item. An object type
// starts with one/multiple keys and has a value. The value might be of any
// type.
func (p *printer) objectItem(o *ast.ObjectItem) []byte ***REMOVED***
	defer un(trace(p, fmt.Sprintf("ObjectItem: %s", o.Keys[0].Token.Text)))
	var buf bytes.Buffer

	if o.LeadComment != nil ***REMOVED***
		for _, comment := range o.LeadComment.List ***REMOVED***
			buf.WriteString(comment.Text)
			buf.WriteByte(newline)
		***REMOVED***
	***REMOVED***

	for i, k := range o.Keys ***REMOVED***
		buf.WriteString(k.Token.Text)
		buf.WriteByte(blank)

		// reach end of key
		if o.Assign.IsValid() && i == len(o.Keys)-1 && len(o.Keys) == 1 ***REMOVED***
			buf.WriteString("=")
			buf.WriteByte(blank)
		***REMOVED***
	***REMOVED***

	buf.Write(p.output(o.Val))

	if o.Val.Pos().Line == o.Keys[0].Pos().Line && o.LineComment != nil ***REMOVED***
		buf.WriteByte(blank)
		for _, comment := range o.LineComment.List ***REMOVED***
			buf.WriteString(comment.Text)
		***REMOVED***
	***REMOVED***

	return buf.Bytes()
***REMOVED***

// objectType returns the printable HCL form of an object type. An object type
// begins with a brace and ends with a brace.
func (p *printer) objectType(o *ast.ObjectType) []byte ***REMOVED***
	defer un(trace(p, "ObjectType"))
	var buf bytes.Buffer
	buf.WriteString("***REMOVED***")

	var index int
	var nextItem token.Pos
	var commented, newlinePrinted bool
	for ***REMOVED***
		// Determine the location of the next actual non-comment
		// item. If we're at the end, the next item is the closing brace
		if index != len(o.List.Items) ***REMOVED***
			nextItem = o.List.Items[index].Pos()
		***REMOVED*** else ***REMOVED***
			nextItem = o.Rbrace
		***REMOVED***

		// Go through the standalone comments in the file and print out
		// the comments that we should be for this object item.
		for _, c := range p.standaloneComments ***REMOVED***
			printed := false
			var lastCommentPos token.Pos
			for _, comment := range c.List ***REMOVED***
				// We only care about comments after the previous item
				// we've printed so that comments are printed in the
				// correct locations (between two objects for example).
				// And before the next item.
				if comment.Pos().After(p.prev) && comment.Pos().Before(nextItem) ***REMOVED***
					// If there are standalone comments and the initial newline has not
					// been printed yet, do it now.
					if !newlinePrinted ***REMOVED***
						newlinePrinted = true
						buf.WriteByte(newline)
					***REMOVED***

					// add newline if it's between other printed nodes
					if index > 0 ***REMOVED***
						commented = true
						buf.WriteByte(newline)
					***REMOVED***

					// Store this position
					lastCommentPos = comment.Pos()

					// output the comment itself
					buf.Write(p.indent(p.heredocIndent([]byte(comment.Text))))

					// Set printed to true to note that we printed something
					printed = true

					/*
						if index != len(o.List.Items) ***REMOVED***
							buf.WriteByte(newline) // do not print on the end
						***REMOVED***
					*/
				***REMOVED***
			***REMOVED***

			// Stuff to do if we had comments
			if printed ***REMOVED***
				// Always write a newline
				buf.WriteByte(newline)

				// If there is another item in the object and our comment
				// didn't hug it directly, then make sure there is a blank
				// line separating them.
				if nextItem != o.Rbrace && nextItem.Line != lastCommentPos.Line+1 ***REMOVED***
					buf.WriteByte(newline)
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if index == len(o.List.Items) ***REMOVED***
			p.prev = o.Rbrace
			break
		***REMOVED***

		// At this point we are sure that it's not a totally empty block: print
		// the initial newline if it hasn't been printed yet by the previous
		// block about standalone comments.
		if !newlinePrinted ***REMOVED***
			buf.WriteByte(newline)
			newlinePrinted = true
		***REMOVED***

		// check if we have adjacent one liner items. If yes we'll going to align
		// the comments.
		var aligned []*ast.ObjectItem
		for _, item := range o.List.Items[index:] ***REMOVED***
			// we don't group one line lists
			if len(o.List.Items) == 1 ***REMOVED***
				break
			***REMOVED***

			// one means a oneliner with out any lead comment
			// two means a oneliner with lead comment
			// anything else might be something else
			cur := lines(string(p.objectItem(item)))
			if cur > 2 ***REMOVED***
				break
			***REMOVED***

			curPos := item.Pos()

			nextPos := token.Pos***REMOVED******REMOVED***
			if index != len(o.List.Items)-1 ***REMOVED***
				nextPos = o.List.Items[index+1].Pos()
			***REMOVED***

			prevPos := token.Pos***REMOVED******REMOVED***
			if index != 0 ***REMOVED***
				prevPos = o.List.Items[index-1].Pos()
			***REMOVED***

			// fmt.Println("DEBUG ----------------")
			// fmt.Printf("prev = %+v prevPos: %s\n", prev, prevPos)
			// fmt.Printf("cur = %+v curPos: %s\n", cur, curPos)
			// fmt.Printf("next = %+v nextPos: %s\n", next, nextPos)

			if curPos.Line+1 == nextPos.Line ***REMOVED***
				aligned = append(aligned, item)
				index++
				continue
			***REMOVED***

			if curPos.Line-1 == prevPos.Line ***REMOVED***
				aligned = append(aligned, item)
				index++

				// finish if we have a new line or comment next. This happens
				// if the next item is not adjacent
				if curPos.Line+1 != nextPos.Line ***REMOVED***
					break
				***REMOVED***
				continue
			***REMOVED***

			break
		***REMOVED***

		// put newlines if the items are between other non aligned items.
		// newlines are also added if there is a standalone comment already, so
		// check it too
		if !commented && index != len(aligned) ***REMOVED***
			buf.WriteByte(newline)
		***REMOVED***

		if len(aligned) >= 1 ***REMOVED***
			p.prev = aligned[len(aligned)-1].Pos()

			items := p.alignedItems(aligned)
			buf.Write(p.indent(items))
		***REMOVED*** else ***REMOVED***
			p.prev = o.List.Items[index].Pos()

			buf.Write(p.indent(p.objectItem(o.List.Items[index])))
			index++
		***REMOVED***

		buf.WriteByte(newline)
	***REMOVED***

	buf.WriteString("***REMOVED***")
	return buf.Bytes()
***REMOVED***

func (p *printer) alignedItems(items []*ast.ObjectItem) []byte ***REMOVED***
	var buf bytes.Buffer

	// find the longest key and value length, needed for alignment
	var longestKeyLen int // longest key length
	var longestValLen int // longest value length
	for _, item := range items ***REMOVED***
		key := len(item.Keys[0].Token.Text)
		val := len(p.output(item.Val))

		if key > longestKeyLen ***REMOVED***
			longestKeyLen = key
		***REMOVED***

		if val > longestValLen ***REMOVED***
			longestValLen = val
		***REMOVED***
	***REMOVED***

	for i, item := range items ***REMOVED***
		if item.LeadComment != nil ***REMOVED***
			for _, comment := range item.LeadComment.List ***REMOVED***
				buf.WriteString(comment.Text)
				buf.WriteByte(newline)
			***REMOVED***
		***REMOVED***

		for i, k := range item.Keys ***REMOVED***
			keyLen := len(k.Token.Text)
			buf.WriteString(k.Token.Text)
			for i := 0; i < longestKeyLen-keyLen+1; i++ ***REMOVED***
				buf.WriteByte(blank)
			***REMOVED***

			// reach end of key
			if i == len(item.Keys)-1 && len(item.Keys) == 1 ***REMOVED***
				buf.WriteString("=")
				buf.WriteByte(blank)
			***REMOVED***
		***REMOVED***

		val := p.output(item.Val)
		valLen := len(val)
		buf.Write(val)

		if item.Val.Pos().Line == item.Keys[0].Pos().Line && item.LineComment != nil ***REMOVED***
			for i := 0; i < longestValLen-valLen+1; i++ ***REMOVED***
				buf.WriteByte(blank)
			***REMOVED***

			for _, comment := range item.LineComment.List ***REMOVED***
				buf.WriteString(comment.Text)
			***REMOVED***
		***REMOVED***

		// do not print for the last item
		if i != len(items)-1 ***REMOVED***
			buf.WriteByte(newline)
		***REMOVED***
	***REMOVED***

	return buf.Bytes()
***REMOVED***

// list returns the printable HCL form of an list type.
func (p *printer) list(l *ast.ListType) []byte ***REMOVED***
	var buf bytes.Buffer
	buf.WriteString("[")

	var longestLine int
	for _, item := range l.List ***REMOVED***
		// for now we assume that the list only contains literal types
		if lit, ok := item.(*ast.LiteralType); ok ***REMOVED***
			lineLen := len(lit.Token.Text)
			if lineLen > longestLine ***REMOVED***
				longestLine = lineLen
			***REMOVED***
		***REMOVED***
	***REMOVED***

	insertSpaceBeforeItem := false
	lastHadLeadComment := false
	for i, item := range l.List ***REMOVED***
		// Keep track of whether this item is a heredoc since that has
		// unique behavior.
		heredoc := false
		if lit, ok := item.(*ast.LiteralType); ok && lit.Token.Type == token.HEREDOC ***REMOVED***
			heredoc = true
		***REMOVED***

		if item.Pos().Line != l.Lbrack.Line ***REMOVED***
			// multiline list, add newline before we add each item
			buf.WriteByte(newline)
			insertSpaceBeforeItem = false

			// If we have a lead comment, then we want to write that first
			leadComment := false
			if lit, ok := item.(*ast.LiteralType); ok && lit.LeadComment != nil ***REMOVED***
				leadComment = true

				// If this isn't the first item and the previous element
				// didn't have a lead comment, then we need to add an extra
				// newline to properly space things out. If it did have a
				// lead comment previously then this would be done
				// automatically.
				if i > 0 && !lastHadLeadComment ***REMOVED***
					buf.WriteByte(newline)
				***REMOVED***

				for _, comment := range lit.LeadComment.List ***REMOVED***
					buf.Write(p.indent([]byte(comment.Text)))
					buf.WriteByte(newline)
				***REMOVED***
			***REMOVED***

			// also indent each line
			val := p.output(item)
			curLen := len(val)
			buf.Write(p.indent(val))

			// if this item is a heredoc, then we output the comma on
			// the next line. This is the only case this happens.
			comma := []byte***REMOVED***','***REMOVED***
			if heredoc ***REMOVED***
				buf.WriteByte(newline)
				comma = p.indent(comma)
			***REMOVED***

			buf.Write(comma)

			if lit, ok := item.(*ast.LiteralType); ok && lit.LineComment != nil ***REMOVED***
				// if the next item doesn't have any comments, do not align
				buf.WriteByte(blank) // align one space
				for i := 0; i < longestLine-curLen; i++ ***REMOVED***
					buf.WriteByte(blank)
				***REMOVED***

				for _, comment := range lit.LineComment.List ***REMOVED***
					buf.WriteString(comment.Text)
				***REMOVED***
			***REMOVED***

			lastItem := i == len(l.List)-1
			if lastItem ***REMOVED***
				buf.WriteByte(newline)
			***REMOVED***

			if leadComment && !lastItem ***REMOVED***
				buf.WriteByte(newline)
			***REMOVED***

			lastHadLeadComment = leadComment
		***REMOVED*** else ***REMOVED***
			if insertSpaceBeforeItem ***REMOVED***
				buf.WriteByte(blank)
				insertSpaceBeforeItem = false
			***REMOVED***

			// Output the item itself
			// also indent each line
			val := p.output(item)
			curLen := len(val)
			buf.Write(val)

			// If this is a heredoc item we always have to output a newline
			// so that it parses properly.
			if heredoc ***REMOVED***
				buf.WriteByte(newline)
			***REMOVED***

			// If this isn't the last element, write a comma.
			if i != len(l.List)-1 ***REMOVED***
				buf.WriteString(",")
				insertSpaceBeforeItem = true
			***REMOVED***

			if lit, ok := item.(*ast.LiteralType); ok && lit.LineComment != nil ***REMOVED***
				// if the next item doesn't have any comments, do not align
				buf.WriteByte(blank) // align one space
				for i := 0; i < longestLine-curLen; i++ ***REMOVED***
					buf.WriteByte(blank)
				***REMOVED***

				for _, comment := range lit.LineComment.List ***REMOVED***
					buf.WriteString(comment.Text)
				***REMOVED***
			***REMOVED***
		***REMOVED***

	***REMOVED***

	buf.WriteString("]")
	return buf.Bytes()
***REMOVED***

// indent indents the lines of the given buffer for each non-empty line
func (p *printer) indent(buf []byte) []byte ***REMOVED***
	var prefix []byte
	if p.cfg.SpacesWidth != 0 ***REMOVED***
		for i := 0; i < p.cfg.SpacesWidth; i++ ***REMOVED***
			prefix = append(prefix, blank)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		prefix = []byte***REMOVED***tab***REMOVED***
	***REMOVED***

	var res []byte
	bol := true
	for _, c := range buf ***REMOVED***
		if bol && c != '\n' ***REMOVED***
			res = append(res, prefix...)
		***REMOVED***

		res = append(res, c)
		bol = c == '\n'
	***REMOVED***
	return res
***REMOVED***

// unindent removes all the indentation from the tombstoned lines
func (p *printer) unindent(buf []byte) []byte ***REMOVED***
	var res []byte
	for i := 0; i < len(buf); i++ ***REMOVED***
		skip := len(buf)-i <= len(unindent)
		if !skip ***REMOVED***
			skip = !bytes.Equal(unindent, buf[i:i+len(unindent)])
		***REMOVED***
		if skip ***REMOVED***
			res = append(res, buf[i])
			continue
		***REMOVED***

		// We have a marker. we have to backtrace here and clean out
		// any whitespace ahead of our tombstone up to a \n
		for j := len(res) - 1; j >= 0; j-- ***REMOVED***
			if res[j] == '\n' ***REMOVED***
				break
			***REMOVED***

			res = res[:j]
		***REMOVED***

		// Skip the entire unindent marker
		i += len(unindent) - 1
	***REMOVED***

	return res
***REMOVED***

// heredocIndent marks all the 2nd and further lines as unindentable
func (p *printer) heredocIndent(buf []byte) []byte ***REMOVED***
	var res []byte
	bol := false
	for _, c := range buf ***REMOVED***
		if bol && c != '\n' ***REMOVED***
			res = append(res, unindent...)
		***REMOVED***
		res = append(res, c)
		bol = c == '\n'
	***REMOVED***
	return res
***REMOVED***

// isSingleLineObject tells whether the given object item is a single
// line object such as "obj ***REMOVED******REMOVED***".
//
// A single line object:
//
//   * has no lead comments (hence multi-line)
//   * has no assignment
//   * has no values in the stanza (within ***REMOVED******REMOVED***)
//
func (p *printer) isSingleLineObject(val *ast.ObjectItem) bool ***REMOVED***
	// If there is a lead comment, can't be one line
	if val.LeadComment != nil ***REMOVED***
		return false
	***REMOVED***

	// If there is assignment, we always break by line
	if val.Assign.IsValid() ***REMOVED***
		return false
	***REMOVED***

	// If it isn't an object type, then its not a single line object
	ot, ok := val.Val.(*ast.ObjectType)
	if !ok ***REMOVED***
		return false
	***REMOVED***

	// If the object has no items, it is single line!
	return len(ot.List.Items) == 0
***REMOVED***

func lines(txt string) int ***REMOVED***
	endline := 1
	for i := 0; i < len(txt); i++ ***REMOVED***
		if txt[i] == '\n' ***REMOVED***
			endline++
		***REMOVED***
	***REMOVED***
	return endline
***REMOVED***

// ----------------------------------------------------------------------------
// Tracing support

func (p *printer) printTrace(a ...interface***REMOVED******REMOVED***) ***REMOVED***
	if !p.enableTrace ***REMOVED***
		return
	***REMOVED***

	const dots = ". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . "
	const n = len(dots)
	i := 2 * p.indentTrace
	for i > n ***REMOVED***
		fmt.Print(dots)
		i -= n
	***REMOVED***
	// i <= n
	fmt.Print(dots[0:i])
	fmt.Println(a...)
***REMOVED***

func trace(p *printer, msg string) *printer ***REMOVED***
	p.printTrace(msg, "(")
	p.indentTrace++
	return p
***REMOVED***

// Usage pattern: defer un(trace(p, "..."))
func un(p *printer) ***REMOVED***
	p.indentTrace--
	p.printTrace(")")
***REMOVED***
