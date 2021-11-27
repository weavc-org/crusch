package crusch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/google/go-querystring/query"
)

// Client is used to process requests to and from Githubs v3 api
// URL and protocols can be changed
type Client struct {
	URL      string
	Protocol string
	Headers  []header
	client   *http.Client
}

type header struct {
	Name  string
	Value string
}

var (
	// GithubClient is the default GithubClient, using standard API url and https
	GithubClient = NewGithubClient("api.github.com", "https")
)

// NewGithubClient creates and returns a new GithubClient structure with given values
func NewGithubClient(url string, protocol string) *Client {
	c := &Client{
		URL:      url,
		Protocol: protocol,
		client:   http.DefaultClient,
	}
	c.AddHeader("Accept", "application/vnd.github.machine-man-preview+json")
	return c
}

// SetHTTPClient allows the http.Client on GithubClient to be changed
// http.DefaultClient is used by default
func (c *Client) SetHTTPClient(client *http.Client) {
	c.client = client
}

// AddHeader adds headers to the array of headers used in the request
func (c *Client) AddHeader(name string, value string) {
	c.RemoveHeader(name)
	h := header{Name: name, Value: value}
	c.Headers = append(c.Headers, h)
}

// RemoveHeader headers to the array of headers used in the request
func (c *Client) RemoveHeader(name string) {
	for i, h := range c.Headers {
		if h.Name == name {
			c.Headers = append(c.Headers[:i], c.Headers[i+1:]...)
		}
	}
}

// Get makes GET requests using the providers information
// Additional parameters/querystring can be passed through either as a string, struct or left as nil
// The response body will be bound to v
func (c *Client) Get(authorizer Authorizer, uri string, params interface{}, v interface{}) (*http.Response, error) {
	var req *http.Request

	query, err := parseQuery(params)
	if err != nil {
		return nil, err
	}

	uri = strings.TrimLeft(uri, "/")

	req, err = http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s://%s/%s?%s", c.Protocol, c.URL, uri, query),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	return c.Do(authorizer, req, v)
}

// Delete makes DELETE using the providers information
func (c *Client) Delete(authorizer Authorizer, uri string) (*http.Response, error) {
	var req *http.Request

	uri = strings.TrimLeft(uri, "/")

	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s://%s/%s", c.Protocol, c.URL, uri),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	return c.Do(authorizer, req, nil)
}

// Put makes PUT requests using the providers information
// A request body can be passed through and attempt to be converted to JSON, this can also be left as nil
// The response body will be bound to v
func (c *Client) Put(authorizer Authorizer, uri string, body interface{}, v interface{}) (*http.Response, error) {
	var req *http.Request

	uri = strings.TrimLeft(uri, "/")

	b, err := jsonifyBody(body)
	if err != nil {
		return nil, err
	}

	req, err = http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s://%s/%s", c.Protocol, c.URL, uri),
		b,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	return c.Do(authorizer, req, v)
}

// Patch makes PATCH requests using the providers information
// A request body can be passed through and attempt to be converted to JSON, this can also be left as nil
// The response body will be bound to v
func (c *Client) Patch(authorizer Authorizer, uri string, body interface{}, v interface{}) (*http.Response, error) {
	var req *http.Request

	b, err := jsonifyBody(body)
	if err != nil {
		return nil, err
	}

	uri = strings.TrimLeft(uri, "/")

	req, err = http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("%s://%s/%s", c.Protocol, c.URL, uri),
		b,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	return c.Do(authorizer, req, v)
}

// Post makes POST requests using the providers information
// A request body can be passed through and attempt to be converted to JSON, this can also be left as nil
// The response body will be bound to v
func (c *Client) Post(authorizer Authorizer, uri string, body interface{}, v interface{}) (*http.Response, error) {
	var req *http.Request

	b, err := jsonifyBody(body)
	if err != nil {
		return nil, err
	}

	uri = strings.TrimLeft(uri, "/")

	req, err = http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s://%s/%s", c.Protocol, c.URL, uri),
		b,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	return c.Do(authorizer, req, v)
}

// Do performs the given request using the providers details
// This will also bind the JSON response to v
func (c *Client) Do(authorizer Authorizer, req *http.Request, v interface{}) (*http.Response, error) {
	if req.Header == nil {
		req.Header = http.Header{}
	}

	auth, err := authorizer.GetHeader()
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", auth)

	for _, h := range c.Headers {
		req.Header.Add(h.Name, h.Value)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return res, err
	}

	if v != nil && (res.StatusCode >= 200 && res.StatusCode < 300) {
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(v)
		if err != nil {
			return res, err
		}
	} else if res.StatusCode >= 400 && res.Body != nil {
		// return the response body as error string if request failed/errored
		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		err = fmt.Errorf("%v", buf.String())
	}

	return res, err
}

func parseQuery(params interface{}) (string, error) {
	switch v := params.(type) {
	case string:
		return v, nil
	case nil:
		return "", nil
	default:
		val := reflect.ValueOf(params)
		if val.Kind() == reflect.Struct {
			return "", fmt.Errorf("unknown type of params, must be string, struct or nil")
		}
		q := reflect.ValueOf(params)
		if q.Kind() == reflect.Ptr && q.IsNil() {
			return "", nil
		}

		qs, err := query.Values(params)
		if err != nil {
			return "", err
		}

		return qs.Encode(), nil
	}

}

func jsonifyBody(body interface{}) (*bytes.Buffer, error) {

	if body == nil {
		return bytes.NewBufferString(""), nil
	}

	var buf = &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(body)
	if err != nil {
		return nil, fmt.Errorf("failed to encode body: %v", err)
	}

	return buf, nil
}
