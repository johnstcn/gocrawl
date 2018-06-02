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
	Crawl(j *Job) (Result, error)
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
		httpClient: httpClient,
		parse:      xmlpathParse,
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

// ParseFunc is a function called on body to run rules and return a Result
type ParseFunc func(body []byte, rules []Rule) (Result, error)

var _ Crawler = (*crawler)(nil)

type crawler struct {
	httpClient Doer
	parse      ParseFunc
}

// Doer wraps http.Client
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

var _ Doer = (*http.Client)(nil)

func (c *crawler) Crawl(j *Job) (Result, error) {

	req, err := c.prepareRequest(j)
	if err != nil {
		return Result{}, errors.Wrap(err, "invalid request")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Result{}, errors.Wrap(err, "doing request")
	}

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

func (c *crawler) prepareRequest(j *Job) (*http.Request, error) {
	req, err := http.NewRequest(j.Request.Method, j.Request.URL, strings.NewReader(j.Request.Body))
	if err != nil {
		return &http.Request{}, err
	}

	for k, v := range j.Request.Headers {
		req.Header.Add(k, v)
	}

	return req, nil
}
