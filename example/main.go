package main

import (
	"github.com/teitei-tk/scavenger"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	/*
		t := func() interface{} {
			type res struct {
				Name string
			}

			r := res{Name: "test"}
			return r
		}

		fmt.Println(t())
	*/

	startURLs := []string{"http://jp.techcrunch.com"}
	parser := func(document *goquery.Document) interface{} {
		var texts []string
		document.Find("h2.post-title").ChildrenFiltered("a").Each(func(_ int, sec *goquery.Selection) {
			texts = append(texts, sec.Text())
		})
		return texts
	}

	s := scavenger.New(startURLs, parser, 3)

	s.Run()
}
