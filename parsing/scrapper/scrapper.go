package scrapper

import (
	"bytes"
	"errors"

	"github.com/Cyber-cicco/tree-sitter-query-builder/querier"
	sitter "github.com/smacker/go-tree-sitter"
)

type SelectorType string

const (
	ST_ID    = "ST_ID"
	ST_CLASS = "ST_CLASS"
	ST_BASE  = "ST_BASE"
)

type DOMElement struct {
	Node     *sitter.Node
	document *DOMStructure
}

type DOMStructure struct {
	RootNode *sitter.Node
	content  []byte
}

type Selectable interface {
	QuerySelector(selector string) (*DOMElement, bool)
	QuerySelectorAll(selector string) ([]*DOMElement, bool)
}

func ToDOM(n *sitter.Node, content []byte) (*DOMStructure, error) {
	if n.Type() == "document" {
		return &DOMStructure{
			RootNode: n,
			content:  content,
		}, nil
	}
	return nil, errors.New("node is not an HTML element")
}

func (s *DOMElement) QuerySelector(query string) *DOMElement {
	return nil
}

func (s *DOMStructure) QuerySelector(query string) (*DOMElement, bool) {

	var element *sitter.Node

	if len(query) == 0 {
		return nil, false
	}

	selector, err := parseSelector(query)

	if err != nil {
		return nil, false
	}

	switch selector.sType {

	case ST_BASE:
		element = querier.GetFirstMatch(s.RootNode, func(n *sitter.Node) bool {
			isEl := n.Type() == "element"

			if !isEl {
				return false
			}

			return getTagName(n, s.content) == selector.matched
		})

	case ST_ID:
		element = querier.GetFirstMatch(s.RootNode, func(n *sitter.Node) bool {
			return elementWithAttributeEquals(n, "id", selector.matched, s.content)
		}).Parent()

	case ST_CLASS:
		element = querier.GetFirstMatch(s.RootNode, func(n *sitter.Node) bool {
			return elementWithAttributeEquals(n, "class", selector.matched, s.content)
		}).Parent()
	}

	return &DOMElement{
		Node:     element,
		document: s,
	}, element != nil
}

func elementWithAttributeEquals(n *sitter.Node, attributeName, matched string, content []byte) bool {

	isEl := n.Type() == "start_tag"

	if !isEl {
		return false
	}

	el := querier.GetFirstMatch(n, func(n *sitter.Node) bool {
		return attributeEquals(n, attributeName, matched, content)
	})

	return el != nil
}

func attributeEquals(n *sitter.Node, attributeName, matched string, content []byte) bool {

	isSeachedAttribute := n.Type() == "attribute" && n.Child(0) != nil && n.Child(0).Content(content) == attributeName

	if !isSeachedAttribute {
		return false
	}

	if n.Child(2) == nil || n.Child(2).Child(1) == nil {
		return false
	}

	return n.Child(2).Child(1).Content(content) == matched
}

type selector struct {
	matched string
	sType   SelectorType
}

func parseSelector(query string) (selector, error) {
	switch query[0] {
	case '.':

		if len(query) < 2 {
			return selector{}, errors.New("Erreur dans la syntaxe du sélecteur")
		}

		return selector{
			matched: query[1:],
			sType:   ST_CLASS,
		}, nil

	case '#':

		if len(query) < 2 {
			return selector{}, errors.New("Erreur dans la syntaxe du sélecteur")
		}

		return selector{
			matched: query[1:],
			sType:   ST_ID,
		}, nil

	default:

		return selector{
			matched: query,
			sType:   ST_BASE,
		}, nil
	}
}

func (s *DOMStructure) QuerySelectorAll(selector string) []*sitter.Node {
	return nil
}

func (s *DOMElement) QuerySelectorAll(selector string) []*sitter.Node {
	return nil
}

func (s *DOMElement) InnerText() []byte {

    var buffer bytes.Buffer
    nodes := []*sitter.Node{}
    nodes = querier.GetChildrenMatching(s.Node, func(n *sitter.Node) bool {
        return n.Type() == "text" || n.Type() == "entity"
    }, nodes)

    for _, match := range nodes {

        if match.Type() == "text" {
            buffer.Write([]byte(match.Content(s.document.content)))
        }

        if match.Type() == "entity" {
            char, ok := specialChars[match.Content(s.document.content)]
            if !ok {
                continue
            }
            buffer.Write([]byte{char})
        }
    }

	return buffer.Bytes()
}

func (s *DOMElement) ToString() string {
	return s.Node.Content(s.document.content)
}

func (s *DOMElement) TagName() string {
	return getTagName(s.Node, s.document.content)
}

func getTagName(n *sitter.Node, content []byte) string {
	return n.Child(0).Child(1).Content(content)
}

func (s *DOMElement) GetClass() string {
	return ""
}
