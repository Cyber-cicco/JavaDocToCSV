package decoding

import (
	"context"
	"log"
	"strings"

	"github.com/Cyber-cicco/tree-sitter-query-builder/querier"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/html"
)

const (
	CT_CLASS     = "Class"
	CT_INTERFACE = "Interface"
	CT_ENUM      = "Enum"
	CT_EX        = "Exception"
	CT_ERR       = "Error"
	CT_ANNO      = "Annotation Type"
)

var lang *sitter.Language
var parser *sitter.Parser

func init() {
	lang = html.GetLanguage()
	parser = sitter.NewParser()
	parser.SetLanguage(lang)
}

func FindLinksInSummary(content []byte) map[string][]string {
	classTypeToLink := make(map[string][]string)
	nodes := getTabNodes(content)
	for _, node := range nodes {
		titleNode := findTitleNode(content, node)
		title := titleNode.Child(1).Content(content)
		links := findHrefNodesFromTable(content, node)
		classTypeToLink[title] = links
	}
	return classTypeToLink
}

func getTabNodes(content []byte) []*sitter.Node {
	nodes := []*sitter.Node{}
	tree, err := parser.ParseCtx(context.Background(), nil, content)

	if err != nil {
		log.Fatalf("got error %s", err)
	}

	qb := querier.NewPQ(Q_LINK)

	query, err := qb.GetQuery()

	if err != nil {
		log.Fatalf("got error %s", err)
	}

	query.Tree = tree
	query.Lang = lang
	query.Content = content

	err = query.ExecuteQuery(func(c *sitter.QueryCapture) error {
		if c.Node.Type() == "tag_name" {
			nodes = append(nodes, c.Node.Parent().Parent())
		}
		return nil
	})

	if err != nil {
		log.Fatalf("got error %s", err)
	}
	return nodes

}

func findHrefNodesFromTable(content []byte, node *sitter.Node) []string {
	linkNodes := []*sitter.Node{}
    links := []string{}
	linkNodes = querier.GetChildrenMatching(node, func(node *sitter.Node) bool {
		isStartTag := node.Type() == "start_tag"
		if !isStartTag {
			return false
		}
		attribute := querier.GetFirstMatch(node, func(node *sitter.Node) bool {
			isHref := node.Type() == "attribute" && node.Child(0) != nil && node.Child(0).Content(content) == "href"
			if !isHref {
				return false
			}
			if node.Child(2) == nil || node.Child(2).Child(1) == nil {
				return false
			}
            href := node.Child(2).Child(1).Content(content)
			if strings.HasSuffix(href, ".html") {
                links = append(links, href)
                return true
            }
            return false
		})
		return attribute != nil
	}, linkNodes)
    return links
}

func findTitleNode(content []byte, node *sitter.Node) *sitter.Node {
	return querier.GetFirstMatch(node, func(node *sitter.Node) bool {
		isElement := node.Type() == "element"
		if !isElement {
			return false
		}
		if node.Child(0) == nil {
			return false
		}
		if node.Child(0).Child(0) == nil {
			return false
		}
		isTh := node.Child(0).Child(1).Content(content) == "th"
		if !isTh {
			return false
		}
		if node.Child(1) == nil || node.Child(1).Type() != "text" {
			return false
		}
		classNode := querier.GetFirstMatch(node, func(node *sitter.Node) bool {
			isClass := node.Type() == "attribute"
			if !isClass {
				return false
			}
			attributeValue := querier.GetFirstMatch(node, func(node *sitter.Node) bool {
				return node.Type() == "attribute_value"
			})
			return attributeValue != nil && attributeValue.Content(content) == "colFirst"
		})
		return classNode != nil
	})
}
