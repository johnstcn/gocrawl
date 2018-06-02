# gocrawl

Gocrawl is a simple command-line utility to crawl websites.
It uses `gopkg.in/xmlpath.v2` to perform HTML parsing and XPath evaluation.
Example usage and output:

```
./gocrawl -input example/job.json 
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
``` 