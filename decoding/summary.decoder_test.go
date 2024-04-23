package decoding

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	sitter "github.com/smacker/go-tree-sitter"
)

func TestGetTabNodes(t *testing.T) {

	nodes, _ := initTest(t)

	if len(nodes) != 6 {
		t.Fatalf("expected len of 6, got %d", len(nodes))
	}
}

func initTest(t *testing.T) ([]*sitter.Node, []byte) {

	path, err := filepath.Abs("../test-env/index.html")

	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	file, err := os.ReadFile(path)

	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	return getTabNodes(file), file
}

func TestFindTitleNode(t *testing.T) {

	nodes, content := initTest(t)

	titleNode := findTitleNode(content, nodes[0])
	title := titleNode.Child(1).Content(content)
	if title != CT_INTERFACE {
		t.Fatalf("Should have been Interface, got %s\n", title)
	}

	titleNode = findTitleNode(content, nodes[1])
	title = titleNode.Child(1).Content(content)
	if title != CT_CLASS {
		t.Fatalf("Should have been Class, got %s\n", title)
	}
}

func TestGetTabNode(t *testing.T) {

	nodes, content := initTest(t)
    links := findHrefNodesFromTable(content, nodes[1])

    for _, link := range links {
        fmt.Printf("node.Content(content): %v\n", link)
    }

}
