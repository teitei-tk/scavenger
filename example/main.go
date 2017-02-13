package main

import (
	"fmt"

	"github.com/teitei-tk/scavenger"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	s, err := scavenger.New("http://teitei-tk.hatenablog.com/", func(d *goquery.Document) {
		d.Find("a.entry-title-link").Each(func(_ int, s *goquery.Selection) {
			fmt.Println(s.Text())
			fmt.Println(s.AttrOr("href", ""))
		})
	})

	if err != nil {
		panic(err)
	}

	s.Run()
}
