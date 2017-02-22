package scavenger

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
}

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
	}
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				w.fetch(job)

			case <-Quit:
				return
			}
		}
	}()
}

func (w Worker) fetch(job Job) {
	res, err := w.newRequest(job)
	if err != nil {
		fmt.Println(err)
		return
	}

	r := job.Parser(res)
	fmt.Println(r)
}

func (w Worker) newRequest(job Job) (*goquery.Document, error) {
	u, err := url.Parse(job.URL)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
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
