package scavenger

import "github.com/PuerkitoBio/goquery"

var (
	JobQueue chan Job
	Quit     chan bool
)

type (
	Parser func(doc *goquery.Document)

	Job struct {
		URL    string
		Parser Parser
	}

	Scavenger struct {
		WorkerPool chan chan Job
		maxWorkers int
	}
)

func init() {
	JobQueue = make(chan Job)
	Quit = make(chan bool)
}

func New(maxWorkers int) *Scavenger {
	pool := make(chan chan Job, maxWorkers)
	return &Scavenger{
		WorkerPool: pool,
		maxWorkers: maxWorkers,
	}
}

func (s *Scavenger) Run(startURLs []string, p Parser) {
	for _, url := range startURLs {
		s.Enqueue(Job{URL: url, Parser: p})
	}

	for i := 0; i < s.maxWorkers; i++ {
		NewWorker(s.WorkerPool).Start()
	}

	s.dispatch()
}

func (s *Scavenger) Enqueue(job Job) {
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
