package scavenger

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"runtime"

	"github.com/PuerkitoBio/goquery"
)

type (
	Parser func(doc *goquery.Document)

	Schedule struct {
		URL    *url.URL
		Parser Parser
	}

	Scavenger struct {
		Crawler *http.Client
	}
)

var (
	Schedulers = make(chan Schedule)

	Quit = make(chan bool)

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

	AppendSchedule(Schedule{
		URL:    url,
		Parser: p,
	})

	return s, nil
}

func AppendSchedule(sch Schedule) {
	fmt.Println("ppend new schedule")
	go func() {
		Schedulers <- sch
	}()
}

func Terminate() {
	go func() {
		Quit <- true
	}()
}

func (s *Scavenger) Run() {
	for {
		fmt.Println("running...")
		select {
		case sch := <-Schedulers:
			doc, _ := goquery.NewDocument(sch.URL.String())
			sch.Parser(doc)

		case <-Quit:
			fmt.Println("exit!!!!!!!!!")
			os.Exit(0)
			break
		}
	}
}
