package scavenger

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

var (
	JobQueue chan Job
	Quit     chan bool
	Result   chan interface{}
)

type (
	Parser func(doc *goquery.Document) interface{}

	Job struct {
		URL    string
		Parser Parser
	}

	Scavenger struct {
		WorkerPool chan chan Job
		maxWorkers int
	}

	RequestResult struct {
		Items []interface{}
	}
)

func init() {
	JobQueue = make(chan Job)
	Quit = make(chan bool)
	Result = make(chan interface{})
}

func New(startURLs []string, p Parser, maxWorkers int) *Scavenger {
	for _, url := range startURLs {
		Enqueue(Job{URL: url, Parser: p})
	}

	return &Scavenger{
		WorkerPool: make(chan chan Job, maxWorkers),
		maxWorkers: maxWorkers,
	}
}

func (s *Scavenger) Run() {
	for i := 0; i < s.maxWorkers; i++ {
		NewWorker(s.WorkerPool).Start()
	}

	s.dispatch()
}

func Request(url string, p Parser) {
	Enqueue(Job{URL: url, Parser: p})
}

func Enqueue(job Job) {
	go func() { JobQueue <- job }()
}

func (s *Scavenger) Terminate() {
	go func() { Quit <- true }()
}

func (s *Scavenger) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			go func(job Job) {
				jobChannel := <-s.WorkerPool
				jobChannel <- job
			}(job)

		case <-Quit:
			return
		}
	}
}

func (s *Scavenger) Completion() {
	fmt.Println("finished")
}
