package edgar

type FilingSummary struct {
	Reports []Report `xml:"MyReports>Report"`
}

type Report struct {
	ShortName    string `xml:"ShortName"`
	HTMLFileName string `xml:"HtmlFileName"`
}
