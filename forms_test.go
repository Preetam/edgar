package edgar

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func Test1(t *testing.T) {
	idx := &FormsIndex{}
	f, err := os.Open("/tmp/form.idx")
	if err != nil {
		log.Fatal(err)
	}
	idx.parseEntries(f)
}

func TestFilenameToFormURL(t *testing.T) {
	t.Log(filenameToFormURL("edgar/data/1472468/0001062993-18-001250.txt", "R2.htm"))
}

func TestExtractRows(t *testing.T) {
	b, err := ioutil.ReadFile("./income.htm")
	if err != nil {
		t.Fatal(err)
	}

	// re := regexp.MustCompile(`<tr[^<]*>(\s+|.+)+?</tr[^<]*>`)
	// tdRe := regexp.MustCompile(`<td[^<]*>(?:\s+|(.+))+?</td[^<]*>`)
	// matches := re.FindAllString(string(b), -1)
	// for i, match := range matches {
	// 	fmt.Println("match", i)
	// 	fmt.Println(match + "\n")
	//
	// 	tdMatches := tdRe.FindAllString(match, -1)
	// 	for j, tdMatch := range tdMatches {
	// 		fmt.Println("match", j)
	// 		fmt.Println(tdMatch)
	// 	}
	//
	// 	fmt.Println("===")
	// }
	//
	n, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	var f func(*html.Node)

	row := 0
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			row++
			fmt.Println("============= row", row, "=============")
		}
		if n.Type == html.TextNode && strings.TrimSpace(n.Data) != "" {
			fmt.Println(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
}
