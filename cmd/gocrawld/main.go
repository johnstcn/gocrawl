package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/johnstcn/gocrawl/pkg/crawl"
)

type CrawlerDaemon struct {
	crawler crawl.Crawler
}

func (d *CrawlerDaemon) handleWork(w http.ResponseWriter, r *http.Request) {
	var jobSpec crawl.Job
	if err := json.NewDecoder(r.Body).Decode(&jobSpec); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid job spec"))
		return
	}

	res, err := d.crawler.Crawl(jobSpec)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("error executing job: %s", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	resp, _ := json.Marshal(res)
	w.Write(resp)
}

func main() {
	var host string
	var port int
	flag.StringVar(&host, "host", "127.0.0.1", "host to bind to")
	flag.IntVar(&port, "port", 12345, "port to bind to")
	flag.Parse()

	hostport := fmt.Sprintf("%s:%d", host, port)
	cd := &CrawlerDaemon{
		crawler: crawl.New(crawl.DefaultCrawlerOpts),
	}
	http.HandleFunc("/", cd.handleWork)
	log.Fatal(http.ListenAndServe(hostport, http.DefaultServeMux))
}
