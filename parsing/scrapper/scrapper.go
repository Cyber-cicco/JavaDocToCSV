package scrapper

import (
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
	QuerySelector(selector string) *DOMElement
	QuerySelectorAll(selector string) []*DOMElement
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
		return nil, true
	}

	selector, err := parseSelector(query)

	if err != nil {
		return nil, true
	}

	switch selector.sType {

	case ST_BASE:
		element = querier.GetFirstMatch(s.RootNode, func(n *sitter.Node) bool {
            isEl := n.Type() == "element"
            if !isEl {
                return false
            }
			return n.Child(0).Child(1).Content(s.content) == selector.matched
		})
	}

	return &DOMElement{
		Node:     element,
		document: s,
	}, element != nil
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
	innertText := []byte{}
	return innertText
}

func (s *DOMElement) ToString() string {
	return s.Node.Content(s.document.content)
}

func (s *DOMElement) InnerHTML() *sitter.Node {
	innerHtml := &sitter.Node{}
	return innerHtml
}

func (s *DOMElement) TagName() string {
	return s.Node.Child(0).Child(1).Content(s.document.content)
}

func (s *DOMElement) GetClass() string {
	return ""
}
