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

	parseTable(tree, content, className, "Field Summary", "fields")
	parseTable(tree, content, className, "Method Summary", "methods")
	parseEnumConstants(tree, content, className)

}

func parseTable(tree *sitter.Tree, content []byte, className, text, suffix string) {

	var buffer bytes.Buffer
	table, ok := getTextNode(tree, content, text)

    //return if there is no field summary
    if !ok {
        return
    }

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

func getTextNode(tree *sitter.Tree, content []byte, text string) (*sitter.Node, bool) {

	qb := querier.NewPQ(PQ_TABLE)
	qb.AddValue("text",text)
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

func parseEnumConstants(tree *sitter.Tree, content []byte, className string) {

	var buffer bytes.Buffer
	table, ok := getTextNode(tree, content, "Enum Constant Summary")

    //return if there is no field summary
    if !ok {
        return
    }

    fmt.Printf("enum: %v\n", className)
	rows, err := getRows(table, content)

	if err != nil {
		panic(err)
	}

	for _, row := range rows {
		//Type of the field
		colfirst, ok := row.QuerySelector(".colOne")
		if !ok {
			log.Fatalf("got nil result from querySelector")
		}

		code, ok := colfirst.QuerySelector("code")
		buffer.Write(code.InnerText())
		buffer.Write([]byte{';'})

        //Documentation of the field
		block, ok := colfirst.QuerySelector(".block")
		buffer.Write(block.InnerText())
		buffer.Write([]byte{';', '\n'})

	}

    fmt.Printf("%v\n", string(buffer.Bytes()))
}
