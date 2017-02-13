package main

import (
	"fmt"
	"net/url"

	"github.com/teitei-tk/scavenger"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	s, err := scavenger.New("http://teitei-tk.hatenablog.com/", func(d *goquery.Document) {
		d.Find("a.entry-title-link").Each(func(_ int, sec *goquery.Selection) {
			url, err := url.Parse(sec.AttrOr("href", ""))
			if err != nil {
				return
			}

			scavenger.AppendSchedule(scavenger.Schedule{
				URL: url,
				Parser: func(d2 *goquery.Document) {
					fmt.Println(d2.Find(".entry-title-link").Text())

					scavenger.Terminate()
				},
			})
		})
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("run")

	s.Run()
}
