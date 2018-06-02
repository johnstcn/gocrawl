package crawl

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// errReader is a Reader that return an error when calling Read()
type errReader struct{}

func (r *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("a test error")
}

func Test_Crawl_ErrBadRequest(t *testing.T) {
	c := &crawler{
		req: func(j *Job) (*http.Request, error) {
			return nil, errors.New("a test error")
		},
	}

	_, err := c.Crawl(&Job{})
	require.EqualError(t, err, "invalid request: a test error")
}

func Test_Crawl_ErrDoRequest(t *testing.T) {
	c := &crawler{
		req: func(j *Job) (*http.Request, error) {
			return httptest.NewRequest(http.MethodGet, "http://test.com", nil), nil
		},
		do: func(*http.Request) (*http.Response, error) {
			return nil, errors.New("a test error")
		},
	}

	_, err := c.Crawl(&Job{})
	require.EqualError(t, err, "doing request: a test error")
}

func Test_Crawl_ErrDrainBody(t *testing.T) {
	c := &crawler{
		req: func(j *Job) (*http.Request, error) {
			return httptest.NewRequest(http.MethodGet, "http://test.com", nil), nil
		},
		do: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(&errReader{}),
			}, nil
		},
	}

	_, err := c.Crawl(&Job{})
	require.EqualError(t, err, "draining response body: a test error")
}

func Test_Crawl_ErrParse(t *testing.T) {
	c := &crawler{
		req: func(j *Job) (*http.Request, error) {
			return httptest.NewRequest(http.MethodGet, "http://test.com", nil), nil
		},
		do: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(strings.NewReader("")),
			}, nil
		},
		parse: func(b []byte, rules []Rule) (Result, error) {
			return Result{}, errors.New("a test error")
		},
	}

	_, err := c.Crawl(&Job{})
	require.EqualError(t, err, "parsing retrieved page: a test error")
}

func Test_Crawl_OK(t *testing.T) {
	expected := Result{
		"test": RuleOutput{
			Error:  "",
			Values: []string{"test result"},
		},
	}

	c := &crawler{
		req: func(j *Job) (*http.Request, error) {
			return httptest.NewRequest(http.MethodGet, "http://test.com", nil), nil
		},
		do: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(strings.NewReader("")),
			}, nil
		},
		parse: func(b []byte, rules []Rule) (Result, error) {
			return expected, nil
		},
	}

	res, err := c.Crawl(&Job{})
	require.NoError(t, err)
	require.EqualValues(t, expected, res)
}
