package xmlutil

import (
	"encoding/xml"
	"fmt"
	"io"
	"sort"
)

// A XMLNode contains the values to be encoded or decoded.
type XMLNode struct ***REMOVED***
	Name     xml.Name              `json:",omitempty"`
	Children map[string][]*XMLNode `json:",omitempty"`
	Text     string                `json:",omitempty"`
	Attr     []xml.Attr            `json:",omitempty"`

	namespaces map[string]string
	parent     *XMLNode
***REMOVED***

// NewXMLElement returns a pointer to a new XMLNode initialized to default values.
func NewXMLElement(name xml.Name) *XMLNode ***REMOVED***
	return &XMLNode***REMOVED***
		Name:     name,
		Children: map[string][]*XMLNode***REMOVED******REMOVED***,
		Attr:     []xml.Attr***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

// AddChild adds child to the XMLNode.
func (n *XMLNode) AddChild(child *XMLNode) ***REMOVED***
	if _, ok := n.Children[child.Name.Local]; !ok ***REMOVED***
		n.Children[child.Name.Local] = []*XMLNode***REMOVED******REMOVED***
	***REMOVED***
	n.Children[child.Name.Local] = append(n.Children[child.Name.Local], child)
***REMOVED***

// XMLToStruct converts a xml.Decoder stream to XMLNode with nested values.
func XMLToStruct(d *xml.Decoder, s *xml.StartElement) (*XMLNode, error) ***REMOVED***
	out := &XMLNode***REMOVED******REMOVED***
	for ***REMOVED***
		tok, err := d.Token()
		if err != nil ***REMOVED***
			if err == io.EOF ***REMOVED***
				break
			***REMOVED*** else ***REMOVED***
				return out, err
			***REMOVED***
		***REMOVED***

		if tok == nil ***REMOVED***
			break
		***REMOVED***

		switch typed := tok.(type) ***REMOVED***
		case xml.CharData:
			out.Text = string(typed.Copy())
		case xml.StartElement:
			el := typed.Copy()
			out.Attr = el.Attr
			if out.Children == nil ***REMOVED***
				out.Children = map[string][]*XMLNode***REMOVED******REMOVED***
			***REMOVED***

			name := typed.Name.Local
			slice := out.Children[name]
			if slice == nil ***REMOVED***
				slice = []*XMLNode***REMOVED******REMOVED***
			***REMOVED***
			node, e := XMLToStruct(d, &el)
			out.findNamespaces()
			if e != nil ***REMOVED***
				return out, e
			***REMOVED***
			node.Name = typed.Name
			node.findNamespaces()
			tempOut := *out
			// Save into a temp variable, simply because out gets squashed during
			// loop iterations
			node.parent = &tempOut
			slice = append(slice, node)
			out.Children[name] = slice
		case xml.EndElement:
			if s != nil && s.Name.Local == typed.Name.Local ***REMOVED*** // matching end token
				return out, nil
			***REMOVED***
			out = &XMLNode***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
	return out, nil
***REMOVED***

func (n *XMLNode) findNamespaces() ***REMOVED***
	ns := map[string]string***REMOVED******REMOVED***
	for _, a := range n.Attr ***REMOVED***
		if a.Name.Space == "xmlns" ***REMOVED***
			ns[a.Value] = a.Name.Local
		***REMOVED***
	***REMOVED***

	n.namespaces = ns
***REMOVED***

func (n *XMLNode) findElem(name string) (string, bool) ***REMOVED***
	for node := n; node != nil; node = node.parent ***REMOVED***
		for _, a := range node.Attr ***REMOVED***
			namespace := a.Name.Space
			if v, ok := node.namespaces[namespace]; ok ***REMOVED***
				namespace = v
			***REMOVED***
			if name == fmt.Sprintf("%s:%s", namespace, a.Name.Local) ***REMOVED***
				return a.Value, true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return "", false
***REMOVED***

// StructToXML writes an XMLNode to a xml.Encoder as tokens.
func StructToXML(e *xml.Encoder, node *XMLNode, sorted bool) error ***REMOVED***
	e.EncodeToken(xml.StartElement***REMOVED***Name: node.Name, Attr: node.Attr***REMOVED***)

	if node.Text != "" ***REMOVED***
		e.EncodeToken(xml.CharData([]byte(node.Text)))
	***REMOVED*** else if sorted ***REMOVED***
		sortedNames := []string***REMOVED******REMOVED***
		for k := range node.Children ***REMOVED***
			sortedNames = append(sortedNames, k)
		***REMOVED***
		sort.Strings(sortedNames)

		for _, k := range sortedNames ***REMOVED***
			for _, v := range node.Children[k] ***REMOVED***
				StructToXML(e, v, sorted)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for _, c := range node.Children ***REMOVED***
			for _, v := range c ***REMOVED***
				StructToXML(e, v, sorted)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	e.EncodeToken(xml.EndElement***REMOVED***Name: node.Name***REMOVED***)
	return e.Flush()
***REMOVED***
