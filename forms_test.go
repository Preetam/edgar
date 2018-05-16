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
	err = idx.parseEntries(f)
	if err != nil {
		t.Fatal(err)
	}

	// for _, e := range idx.entries {
	// 	fmt.Println(e, filenameToFormURL(e.FileName, "R2.htm"))
	// }
}

func TestFilenameToFormURL(t *testing.T) {
	t.Log(filenameToFormURL("edgar/data/1472468/0001062993-18-001250.txt", "R2.htm"))
}

func TestExtractRows(t *testing.T) {
	b, err := ioutil.ReadFile("./balance_sheet.htm")
	if err != nil {
		t.Fatal(err)
	}

	n, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	var f func(*html.Node)

	row := 0
	currentRowHeader := ""
	quarterIdx := -1
	bs := &BalanceSheet{}
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			row++
			fmt.Println("============= row", row, "=============")
			quarterIdx = -1
		}
		if n.Type == html.TextNode && strings.TrimSpace(n.Data) != "" {
			fmt.Println(n.Data)
			if quarterIdx == -1 {
				currentRowHeader = n.Data
			} else {
				if row == 1 {
					bs.AddQuarter(n.Data)
				}

				bs.SetQuarterInfo(quarterIdx, currentRowHeader, n.Data)
			}
			quarterIdx++
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)

	fmt.Printf("%#v\n", bs)
}
