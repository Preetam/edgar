package edgar

import (
	"bufio"
	"io"
	"path"
	"strconv"
	"strings"
)

type FormsIndex struct {
	year    int
	quarter int

	entries []IndexEntry
}

type IndexEntry struct {
	FormType    string
	CompanyName string
	CIK         int
	DateFiled   string
	FileName    string
}

func (idx *FormsIndex) parseEntries(r io.Reader) error {
	br := bufio.NewReader(r)
	line, err := br.ReadString('\n')

	passedHeader := false
	for ; err == nil; line, err = br.ReadString('\n') {
		// Skip until we pass "---"
		if !passedHeader {
			if strings.HasPrefix(line, "---") {
				passedHeader = true
				continue
			}
			continue
		}

		line = strings.TrimSpace(line)
		//  0 = type
		// 12 = company name
		// 74 = CIK
		// 86 = date filed
		// 98 = file name
		parts := []string{
			line[:12],
			line[12:74],
			line[74:86],
			line[86:98],
			line[98:],
		}
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}

		entry := IndexEntry{}
		entry.FormType = parts[0]
		entry.CompanyName = parts[1]
		i, err := strconv.ParseInt(parts[2], 10, 32)
		if err != nil {
			return err
		}
		entry.CIK = int(i)
		entry.DateFiled = parts[3]
		entry.FileName = parts[4]

		idx.entries = append(idx.entries, entry)
	}
	return nil
}

func filenameToFormURL(filename string, file string) string {
	// edgar/data/1472468/0001062993-18-001250.txt
	parts := strings.Split(filename, "/")

	dir := strings.Replace(strings.TrimSuffix(parts[3], ".txt"), "-", "", -1)
	// 0001062993-18-001250.txt

	return "https://www.sec.gov/Archives/" + path.Join(append(parts[:3], dir, file)...)
}
