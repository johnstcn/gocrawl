# Gocrawl

Gocrawl is a small Go library for scraping data from websites. Some command-line utilities are provided as well.

## gocrawl

Gocrawl is a simple command-line utility to crawl websites.
It uses `gopkg.in/xmlpath.v2` to perform HTML parsing and XPath evaluation.

See the `examples` directory for an example job specification.
Example usage and output:

    $ gocrawl -input example/job.json
    {
        "day": {
            "error": "",
            "values": [
                "Saturday"
            ]
        },
        "invalid_regexp": {
            "error": "error parsing regexp: missing closing ]: `[a-z+`",
            "values": [
                "Saturday, 2 June 2018"
            ]
        },
        "invalid_xpath": {
            "error": "compiling xml path \"//*[@id=\\\\\\\"ctdat\\\"]\":8: expected a literal string",
            "values": []
        },
        "month": {
            "error": "",
            "values": [
                "June"
            ]
        },
        "no_matching_xpath": {
            "error": "no match for xpath //*[@id=\"cttdat\"]",
            "values": []
        },
        "time": {
            "error": "",
            "values": [
                "13:18:59"
            ]
        }
    }


## gocrawld

Gocrawld is a daemon version of gocrawl. It accepts a POST request containing a job specification identical to that of `gocrawl` and returns the result of executing the crawl job, encoded as JSON.

Example usage:

    $ gocrawld -host localhost -port 12345 &
    <pid>
    $ curl -XPOST localhost:12345 --data @example/job.json
    <job output will be the same as above except less pretty>

Example docker usage:

```docker run --rm --net=host --detach johnstcn/gocrawld```
