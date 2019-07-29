package main

import (
	"fmt"
	"os"
	"net/url"
	"strings"
	"github.com/PuerkitoBio/goquery"
)

/*
func makediffhtml() {
	os.Exec("/bin/sh", )
}
*/
func main() {
	f, err := os.Open("diff.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		panic(err)
	}
	base_url := "https://rlacollege.edu.in/"
	// use the goquery document...
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		notice_url, _ := s.Attr("href")

		parts := strings.Split(notice_url, "/")
		// ) is wrongfully parsed as markdown syntax, resulting in broken links in msg
		parts[2] = url.QueryEscape(parts[2])
		notice_url = strings.Join(parts[:], "/")

		title := s.Text()
		fmt.Printf("%d: [%s](%s)\n", i, title, base_url + notice_url)
	})

}