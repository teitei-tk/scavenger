package scavenger

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type Request struct {
	URL    string
	Parser Parser
	Method string
}

func (s *Scavenger) fetch(r *Request) {
	res, err := s.newRequest(r)
	if err != nil {
		fmt.Println(err)
		return
	}

	r.Parser(res)
}

func (s *Scavenger) newRequest(r *Request) (*goquery.Document, error) {
	u, err := url.Parse(r.URL)
	if err != nil {
		return nil, err
	}

	if r.Method == "" {
		r.Method = http.MethodGet
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := s.Crawler.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	d, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}

	return d, nil
}
