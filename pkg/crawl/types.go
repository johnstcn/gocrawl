package crawl

// Job is a single crawl job, consisting of a defined HTTP request and a number of Rules.
type Job struct {
	Request Request `json:"request"`
	Rules   []Rule  `json:"rules"`
}

// Request defines a HTTP request to be made as part of a Job.
type Request struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// Rule defines an element to be targeted on a page, with a number of filters to be run on the result of the XPath evaluation.
// Name is the unique name of the rule.
type Rule struct {
	Name    string   `json:"name"`
	XPath   string   `json:"xpath"`
	Filters []Filter `json:"filters"`
}

// Filter is a regular expression used as part of a Rule.
type Filter struct {
	Find    string `json:"find"`
	Replace string `json:"replace"`
}

// Result is the result of a crawl job. Outputs of rules are indexed by Rule name.
type Result map[string]RuleOutput

// RuleOutput is the result of running a Rule on a page.
type RuleOutput struct {
	Error  string   `json:"error"`
	Values []string `json:"values"`
}
