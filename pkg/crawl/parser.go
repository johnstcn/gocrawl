package crawl

import (
	"bytes"
	"fmt"
	"regexp"

	"gopkg.in/xmlpath.v2"
)

func xmlpathParse(body []byte, rules []Rule) (Result, error) {
	page, err := xmlpath.ParseHTML(bytes.NewReader(body))
	if err != nil {
		return Result{}, err
	}

	output := make(map[string]RuleOutput)
	for _, rule := range rules {
		output[rule.Name] = runRule(page, rule)
	}

	return output, nil
}

func runRule(page *xmlpath.Node, r Rule) RuleOutput {
	out := RuleOutput{
		Values: make([]string, 0),
		Error:  "",
	}
	xp, err := xmlpath.Compile(r.XPath)
	if err != nil {
		out.Error = err.Error()
		return out
	}

	var matched bool
	iter := xp.Iter(page)
	for iter.Next() {
		matched = true
		rawVal := iter.Node().String()
		filtered, err := runFilters(rawVal, r.Filters)
		if err != nil {
			out.Error = err.Error()
		}
		out.Values = append(out.Values, filtered)
	}

	if !matched {
		out.Error = fmt.Sprintf("no match for xpath %s", r.XPath)
	}

	return out
}

func runFilters(raw string, filters []Filter) (string, error) {
	var err error
	var expr *regexp.Regexp
	res := raw

	for _, filter := range filters {
		expr, err = regexp.Compile(filter.Find)
		if err != nil {
			return res, err
		}
		res = expr.ReplaceAllString(res, filter.Replace)
	}
	return res, nil
}
