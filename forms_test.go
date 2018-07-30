package edgar

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
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
	for _, e := range idx.Entries {
		fmt.Println(e)
		if e.FormType != "10-Q" {
			continue
		}
		bs, err := e.Get10QBalanceSheet()
		if err != nil {
			t.Log(e)
			continue
		}
		fmt.Printf("%#v\n", bs)
		//t.Log(e, filenameToFormURL(e.FileName, "R2.htm"))
		//t.Log(e, filenameToFormURL(e.FileName, "FilingSummary.xml"))
		time.Sleep(time.Second)
	}
}

func Test2(t *testing.T) {
	f, err := os.Open("./FilingSummary.xml")
	if err != nil {
		log.Fatal(err)
	}
	filingSummary := &FilingSummary{}
	d := xml.NewDecoder(f)
	err = d.Decode(filingSummary)
	if err != nil {
		log.Fatal(err)
	}
	t.Log(filingSummary)
}

func TestFilenameToFormURL(t *testing.T) {
	t.Log(filenameToFormURL("edgar/data/1472468/0001062993-18-001250.txt", "R2.htm"))
}

func TestExtractRows(t *testing.T) {
	b, err := ioutil.ReadFile("./balance_sheet.htm")
	if err != nil {
		t.Fatal(err)
	}
	bs, err := ParseBalanceSheetHTML(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	log.Println(bs)
}
