package decoding

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Cyber-cicco/HTMLtoDB/scrapper"
	"github.com/Cyber-cicco/tree-sitter-query-builder/querier"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/html"
)

var Lang *sitter.Language
var Parser *sitter.Parser

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
	Lang = html.GetLanguage()
	Parser = sitter.NewParser()
	Parser.SetLanguage(Lang)
}

func ParseSingleFile(content []byte, filePath string) {

	tree, err := Parser.ParseCtx(context.Background(), nil, content)

	if err != nil {
		log.Fatalf("got error %s", err)
	}

	className := strings.TrimSuffix(filePath, ".html")

	parseFields(tree, content, className)

}

func parseFields(tree *sitter.Tree, content []byte, className string) {

	var buffer bytes.Buffer
	table, ok := getFieldNode(tree, content)

    //return if there is no field summary
    if !ok {
        return
    }

    fmt.Printf("className: %v\n", className)
	rows, err := getRows(table, content)
	if err != nil {
		panic(err)
	}

	for _, row := range rows {
		//Type of the field
		colfirst, ok := row.QuerySelector(".colFirst")
		if !ok {
			log.Fatalf("got nil result from querySelector")
		}

		buffer.Write(colfirst.InnerText())
		buffer.Write([]byte{';'})

		//Name of the field
		colLast, ok := row.QuerySelector(".colLast")
		if !ok {
			log.Fatalf("got nil result from querySelector")
		}
		code, ok := colLast.QuerySelector("code")
		buffer.Write(code.InnerText())
		buffer.Write([]byte{';'})

        //Documentation of the field
		block, ok := colLast.QuerySelector(".block")
		buffer.Write(block.InnerText())
		buffer.Write([]byte{';', '\n'})

	}

    fmt.Printf("%v\n", string(buffer.Bytes()))
}

func getRows(table *sitter.Node, content []byte) ([]*scrapper.DOMElement, error) {

	root, err := scrapper.ToDOM(table, content)

	if err != nil {
		return nil, err
	}

	nodes1, ok := root.QuerySelectorAll(".altColor")

	if !ok {
		return nil, errors.New("Got no results while trying")
	}

	nodes2, ok := root.QuerySelectorAll(".rowColor")

	return append(nodes1, nodes2...), nil

}

func getFieldNode(tree *sitter.Tree, content []byte) (*sitter.Node, bool) {

	qb := querier.NewPQ(PQ_TABLE)
	qb.AddValue("text", "Field Summary")
	query, err := qb.GetQuery()

	if err != nil {
		log.Fatalf("got error %s", err)
	}

	query.Lang = Lang
	query.Content = content
	query.Tree = tree
	var node *sitter.Node
	err = query.ExecuteQuery(func(c *sitter.QueryCapture) error {
		if c.Node.Type() == "element" {
			node = c.Node
		}
		return nil
	})

	return node, node != nil
}

func parseMethods(tree *sitter.Tree, content []byte, className string) {

}

func parseEnumConstants(tree *sitter.Tree, content []byte, className string) {

}
