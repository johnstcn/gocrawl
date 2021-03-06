package crawl

//go:generate mockery -package=crawltest -interface=Crawler

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Crawler is something that executes Jobs.
type Crawler interface {
	// Crawl executes Job synchonously.
	Crawl(j Job) (Result, error)
	// CrawlAsync executes Job asynchronously in a separate goroutine.
	CrawlAsync(j Job) (<-chan Result, <-chan error)
}

// New returns a new Crawler.
func New(opts CrawlerOpts) Crawler {
	return newCrawler(opts)
}

func newCrawler(opts CrawlerOpts) *crawler {
	httpClient := &http.Client{
		Timeout: opts.Timeout,
	}
	return &crawler{
		req:   prepareRequest,
		do:    httpClient.Do,
		parse: xmlpathParse,
	}
}

// CrawlerOpts are run-time configuration options for a Crawler
type CrawlerOpts struct {
	Timeout time.Duration
}

// DefaultCrawlerOpts are default crawler options
var DefaultCrawlerOpts = CrawlerOpts{
	Timeout: time.Duration(10) * time.Second,
}

// ReqFunc is a function called to prepare a HTTP request
type ReqFunc func(r *Request) (*http.Request, error)

// DoFunc is a function called to perform a HTTP request
type DoFunc func(*http.Request) (*http.Response, error)

// ParseFunc is a function called on body to run rules and return a Result
type ParseFunc func(body []byte, rules []Rule) (Result, error)

var _ Crawler = (*crawler)(nil)

type crawler struct {
	req   ReqFunc
	do    DoFunc
	parse ParseFunc
}

// Doer wraps http.Client
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

var _ Doer = (*http.Client)(nil)

func (c *crawler) Crawl(j Job) (Result, error) {
	return c.crawl(&j)
}

func (c *crawler) CrawlAsync(j Job) (<-chan Result, <-chan error) {
	rch := make(chan Result)
	ech := make(chan error)
	go func() {
		res, err := c.crawl(&j)
		rch <- res
		ech <- err
	}()
	return rch, ech
}

func (c *crawler) crawl(j *Job) (Result, error) {
	req, err := c.req(&j.Request)
	if err != nil {
		return Result{}, errors.Wrap(err, "invalid request")
	}

	resp, err := c.do(req)
	if err != nil {
		return Result{}, errors.Wrap(err, "doing request")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Result{}, errors.Wrap(err, "draining response body")
	}

	output, err := c.parse(body, j.Rules)
	if err != nil {
		return Result{}, errors.Wrap(err, "parsing retrieved page")
	}

	return output, nil
}

func prepareRequest(r *Request) (*http.Request, error) {
	req, err := http.NewRequest(r.Method, r.URL, strings.NewReader(r.Body))
	if err != nil {
		return &http.Request{}, err
	}

	for k, v := range r.Headers {
		req.Header.Add(k, v)
	}

	return req, nil
}
