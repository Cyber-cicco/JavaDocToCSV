package scrapper

import (
	"context"
	"os"
	"testing"

	"github.com/Cyber-cicco/HTMLtoDB/config"
	"github.com/Cyber-cicco/HTMLtoDB/decoding"
)

const DIV_1 = "<div>JavaScript is disabled on your browser.</div>"

const CLASS_1 = `<th class="colOne" scope="col">Constructor and Description</th>`

const ID_1 = 
        `<ul class="navList" id="allclasses_navbar_top">
            <li><a href="../../allclasses-noframe.html">All&nbsp;Classes</a></li>
        </ul>`

func TestDOMStructure(t *testing.T) {

    content := initTest(t)

	tree, err := decoding.Parser.ParseCtx(context.Background(), nil, content)

	if err != nil {
		t.Fatalf("got error %s", err)
	}

    _, err = ToDOM(tree.RootNode(), content)

	if err != nil {
		t.Fatalf("got error %s", err)
	}

    _, err = ToDOM(tree.RootNode().Child(0), content)

	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestQuerySelector(t *testing.T) {
    
    content := initTest(t)

	tree, err := decoding.Parser.ParseCtx(context.Background(), nil, content)

	if err != nil {
		t.Fatalf("got error %s", err)
	}

    document, err := ToDOM(tree.RootNode(), content)


    div1, ok := document.QuerySelector("div")

    if !ok {
		t.Fatalf("QuerySelector returned nil")
    }

    if div1.ToString() != DIV_1 {
        t.Fatalf("Expected %s, got %s", DIV_1, div1.ToString())
    }

    class1, ok := document.QuerySelector(".colOne")

    if !ok {
		t.Fatalf("QuerySelector returned nil")
    }

    if class1.ToString() != CLASS_1 {
        t.Fatalf("Expected %s, got %s", CLASS_1, class1.ToString())
    }

    id1, ok := document.QuerySelector("#allclasses_navbar_top")

    if id1.ToString() != ID_1 {
        t.Fatalf("Expected\n%s, got \n%s", ID_1, id1.ToString())
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
