package edgar

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type BalanceSheet struct {
	url      string
	Quarters []BalanceSheetQuarter
	factor   string
}

func (bs *BalanceSheet) AddQuarter(date string) {
	bs.Quarters = append(bs.Quarters, BalanceSheetQuarter{
		Date: date,
	})
}

func (bs *BalanceSheet) SetQuarterInfo(quarterIndex int, rowHeader string, rowData string) {
	switch strings.ToLower(rowHeader) {
	case "cash and cash equivalents":
		bs.Quarters[quarterIndex].CashAndCashEquivalents = rowData
	case "total current assets":
		bs.Quarters[quarterIndex].TotalCurrentAssets = rowData
	case "total assets":
		bs.Quarters[quarterIndex].TotalAssets = rowData
	case "total current liabilities":
		bs.Quarters[quarterIndex].TotalCurrentLiabilities = rowData
	}
}

type BalanceSheetQuarter struct {
	Date                    string
	CashAndCashEquivalents  string
	TotalCurrentAssets      string
	TotalAssets             string
	TotalCurrentLiabilities string
}

func ParseBalanceSheetHTML(r io.Reader) (*BalanceSheet, error) {
	n, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	var f func(*html.Node)

	row := 0
	currentRowHeader := ""
	quarterIdx := -1
	bs := &BalanceSheet{}
	f = func(n *html.Node) {
		if n.Data == "br" {
			return
		}
		if n.Type == html.ElementNode && n.Data == "tr" {
			row++
			quarterIdx = -1
		}
		if n.Type == html.TextNode && strings.TrimSpace(n.Data) != "" {
			if quarterIdx == -1 {
				currentRowHeader = n.Data
			} else {
				if row == 1 {
					if trimmedFactor := strings.TrimSpace(n.Data); strings.HasPrefix(trimmedFactor, "$") {
						bs.factor = trimmedFactor
					} else {
						bs.AddQuarter(n.Data)
					}
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

	return bs, nil
}
