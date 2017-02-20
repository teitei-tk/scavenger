package main

import (
	"fmt"

	"github.com/teitei-tk/scavenger"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	s := scavenger.New(3)
	u := "http://jp.techcrunch.com"
	printTitle := func(d *goquery.Document) {
		d.Find("h2.post-title").ChildrenFiltered("a").Each(func(_ int, sec *goquery.Selection) {
			fmt.Println(sec.Text())
		})
	}

	url := func(d *goquery.Document) string {
		url := u + d.Find("li.next").ChildrenFiltered("a").AttrOr("href", "")
		return url
	}

	p := func(d *goquery.Document) {
		r := url(d)
		if r != "" {
			s.Enqueue(scavenger.Job{URL: r, Parser: printTitle})
		}
	}

	s.Run([]string{"http://jp.techcrunch.com/"}, p)
}
