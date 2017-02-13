package scavenger

import (
	"net/http"
	"net/url"
	"runtime"

	"github.com/PuerkitoBio/goquery"
)

type (
	Parser func(doc *goquery.Document)

	Scheduler struct {
		URL    *url.URL
		Parser Parser
	}

	Scavenger struct {
		Schedulers []*Scheduler
		Crawler    *http.Client
	}
)

var (
	defaulReqConcurrency = runtime.NumCPU()
)

func New(rawURL string, p Parser) (*Scavenger, error) {
	url, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	reqConcurrency := defaulReqConcurrency
	runtime.GOMAXPROCS(reqConcurrency)

	if p == nil {
		panic("parser is not implemented")
	}

	// @TODO: Tuning
	client := http.DefaultClient

	s := &Scavenger{
		Crawler: client,
	}

	s.Schedulers = append(s.Schedulers, &Scheduler{
		URL:    url,
		Parser: p,
	})

	return s, nil
}

func (s *Scavenger) Run() {
	for _, s := range s.Schedulers {
		doc, err := goquery.NewDocument(s.URL.String())
		if err != nil {
			panic(err)
		}
		s.Parser(doc)
	}
}
