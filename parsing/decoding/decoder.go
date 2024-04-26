package decoding

import (
	"context"
	"log"
	"strings"

	"github.com/Cyber-cicco/tree-sitter-query-builder/querier"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/html"
)

var lang *sitter.Language
var parser *sitter.Parser

type Method struct {
	Modifier      bool
	ReturnType    string
	Signature     string
	Documentation string
}

type Field struct {
	Modifier      bool
	Type          string
	Identifier    string
	Documentation string
}

func init() {
	lang = html.GetLanguage()
	parser = sitter.NewParser()
	parser.SetLanguage(lang)
}

func ParseSingleFile(content []byte, filePath string) {

	tree, err := parser.ParseCtx(context.Background(), nil, content)

	if err != nil {
		log.Fatalf("got error %s", err)
	}

	className := strings.TrimSuffix(filePath, ".html")

	parseFields(tree, content, className)

}

func parseFields(tree *sitter.Tree, content []byte, className string) {

	table := getFieldNode(tree, content)

    rows := getRows(table, content)
}

func getRows(table *sitter.Node, content []byte) []*sitter.Node {

    nodes := []*sitter.Node{}

    querier.GetChildrenMatching(table, func(n *sitter.Node) bool {
        isEl := n.Type() == "element"
        if !isEl {
            return false
        }
        tr := 
    }, nodes)

    return nodes

}
func getFieldNode(tree *sitter.Tree, content []byte) *sitter.Node {

	qb := querier.NewPQ(PQ_TABLE)
	qb.AddValue("text", "Field Summary")
	query, err := qb.GetQuery()

	if err != nil {
		log.Fatalf("got error %s", err)
	}

	query.Lang = lang
	query.Content = content
	query.Tree = tree
	var node *sitter.Node
	err = query.ExecuteQuery(func(c *sitter.QueryCapture) error {
		if c.Node.Type() == "element" {
			node = c.Node
		}
		return nil
	})

	return node
}

func parseMethods(tree *sitter.Tree, content []byte, className string) {

}

func parseEnumConstants(tree *sitter.Tree, content []byte, className string) {

}
