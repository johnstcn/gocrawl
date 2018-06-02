package crawl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_xmlpathParse_OK(t *testing.T) {
	body := []byte(`<html>
		<body>
		<p>testing 123</p>
		</body>
		</html>`)
	rule := Rule{
		Name:  "testrule",
		XPath: "/html/body/p",
		Filters: []Filter{
			Filter{
				Find:    "^[a-z]+\\s+(\\d+)$",
				Replace: "$1",
			},
		},
	}

	expected := Result{
		"testrule": RuleOutput{
			Error:  "",
			Values: []string{"123"},
		},
	}

	actual, err := xmlpathParse(body, []Rule{rule})
	require.NoError(t, err)
	require.EqualValues(t, expected, actual)
}

func Test_xmlpathParse_ErrNoMatch(t *testing.T) {
	body := []byte(`<html>
		<body>
		<b>test</b>
		</body>
		</html>`)
	rule := Rule{
		Name:  "testrule",
		XPath: "/html/body/p",
	}

	expected := Result{
		"testrule": RuleOutput{
			Error:  "no match for xpath /html/body/p",
			Values: []string{},
		},
	}

	actual, err := xmlpathParse(body, []Rule{rule})
	require.NoError(t, err)
	require.EqualValues(t, expected, actual)
}

func Test_xmlpathParse_InvalidExpr(t *testing.T) {
	body := []byte(`<html>
		<body>
		<p>testing 123</p>
		</body>
		</html>`)
	rule := Rule{
		Name:  "testrule",
		XPath: "/html/body/p",
		Filters: []Filter{
			Filter{
				Find:    "^[a-z+\\s+\\d+)$",
				Replace: "$1",
			},
		},
	}

	expected := Result{
		"testrule": RuleOutput{
			Error:  "error parsing regexp: missing closing ]: `[a-z+\\s+\\d+)$`",
			Values: []string{"testing 123"},
		},
	}

	actual, err := xmlpathParse(body, []Rule{rule})
	require.NoError(t, err)
	require.EqualValues(t, expected, actual)
}

func Test_xmlpathParse_InvalidXPath(t *testing.T) {
	body := []byte(`<html>
		<body>
		<p>testing 123</p>
		</body>
		</html>`)
	rule := Rule{
		Name:    "testrule",
		XPath:   "//foo[@bar",
		Filters: []Filter{},
	}

	expected := Result{
		"testrule": RuleOutput{
			Error:  "compiling xml path \"//foo[@bar\":10: expected ']'",
			Values: []string{},
		},
	}

	actual, err := xmlpathParse(body, []Rule{rule})
	require.NoError(t, err)
	require.EqualValues(t, expected, actual)
}
