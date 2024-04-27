package decoding

import (
	"context"
	"os"
	"testing"

	"github.com/Cyber-cicco/HTMLtoDB/config"
)

func TestParseFields(t *testing.T) {
    content := initTest(t)
    tree, err := Parser.ParseCtx(context.Background(), nil, content)

    if err != nil {
        t.Fatalf("Shouldn't have had error, got %s", err)
    }
    
    node := getFieldNode(tree, content)

    if node == nil {
        t.Fatalf("Expected a value, got nil")
    }
}

func initTest(t *testing.T) []byte {
    path := config.URL_RESOURCES + "Class/Float.html"
    content, err := os.ReadFile(path)

    if err != nil {
        t.Fatalf("Shouldn't have had error, got %s", err)
    }

    return content
}
