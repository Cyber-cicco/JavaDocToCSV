package scrapper

import (
	"errors"

	sitter "github.com/smacker/go-tree-sitter"
)

type DOMElement struct {
	Node    *sitter.Node
	content []byte
}

func ToDOM(n *sitter.Node, content []byte) (*DOMElement, error) {
	if n.Type() == "element" {
		return &DOMElement{Node: n, content: content}, nil
	}
	return nil, errors.New("node is not an HTML element")
}

func (s *DOMElement) QuerySelector(selector string) *sitter.Node {
	return nil
}

func (s *DOMElement) QuerySelectorAll(selector string) []*sitter.Node {
	return nil
}

func (s *DOMElement) InnerText() []byte {
	innertText := []byte{}
	return innertText
}

func (s *DOMElement) InnerHTML() *sitter.Node {
	innerHtml := &sitter.Node{}
	return innerHtml
}

func (s *DOMElement) TagName() string {
	return s.Node.Child(0).Child(1).Content(s.content)
}
