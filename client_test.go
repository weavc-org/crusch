package crusch

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

var (
	generalResponse map[string]string = map[string]string{
		"weavc": "crusch",
		"one":   "1",
	}
)

type m struct {
	Weavc string `json:"weavc" url:"weavc"`
	One   string `json:"one" url:"one"`
}

func TestNewGithubClient(t *testing.T) {
	c := NewGithubClient("test.url", "http")
	if c.URL != "test.url" || c.Protocol != "http" {
		t.Errorf("valid client: returned %v", c)
	}
}

func TestGet(t *testing.T) {
	client := setupClient(generalResponse)

	var v m
	_, err := client.Get(setupAuth(), "test/uri", nil, &v)
	if err != nil {
		t.Errorf("valid get: unexpected %v", err)
	}

	want := m{
		Weavc: "crusch",
		One:   "1",
	}

	if !reflect.DeepEqual(v, want) {
		t.Errorf("valid get: returned %v want %v", v, want)
	}

	_, err = client.Get(setupAuth(), "test/uri", nil, nil)
	if err != nil {
		t.Errorf("valid get, no binding: unexpected %v", err)
	}

	_, err = client.Get(setupAuth(), "test/uri", "a=b&c=1", nil)
	if err != nil {
		t.Errorf("valid get, querystring: unexpected %v", err)
	}

	_, err = client.Get(setupAuth(), "test/uri", 123, nil)
	if err == nil {
		t.Errorf("invalid get, bad querystring: unexpected nil error %v", err)
	}

	type badBinding struct {
		w string
		h string
	}
	b := badBinding{w: "", h: ""}
	_, err = client.Get(setupAuth(), "test/uri", nil, b)
	if err == nil {
		t.Errorf("invalid get, bad binding: unexpected nil error %v", err)
	}

	client.Protocol = "http+_"
	_, err = client.Get(setupAuth(), "test/uri", nil, &v)
	if err == nil {
		t.Errorf("invalid get: unexpected %v", err)
	}
}

func TestDelete(t *testing.T) {
	client := setupClient(generalResponse)

	_, err := client.Delete(setupAuth(), "test/uri")
	if err != nil {
		t.Errorf("valid delete: unexpected %v", err)
	}

	client.Protocol = "http+_"
	_, err = client.Delete(setupAuth(), "test/uri")
	if err == nil {
		t.Errorf("invalid delete: unexpected %v", err)
	}
}

func TestPut(t *testing.T) {
	client := setupClient(generalResponse)

	var v m
	_, err := client.Put(setupAuth(), "test/uri", &m{Weavc: "crusch", One: "1"}, &v)
	if err != nil {
		t.Errorf("valid put: unexpected %v", err)
	}

	want := m{
		Weavc: "crusch",
		One:   "1",
	}

	if !reflect.DeepEqual(v, want) {
		t.Errorf("valid put: returned %v want %v", v, want)
	}

	_, err = client.Put(setupAuth(), "test/uri", nil, &v)
	if err != nil {
		t.Errorf("valid put: unexpected %v", err)
	}

	client.Protocol = "http+_"
	_, err = client.Put(setupAuth(), "test/uri", nil, &v)
	if err == nil {
		t.Errorf("invalid put: unexpected %v", err)
	}
}

func TestPatch(t *testing.T) {
	client := setupClient(generalResponse)

	var v m
	_, err := client.Patch(setupAuth(), "test/uri", &m{Weavc: "crusch", One: "1"}, &v)
	if err != nil {
		t.Errorf("valid patch: unexpected %v", err)
	}

	want := m{
		Weavc: "crusch",
		One:   "1",
	}

	if !reflect.DeepEqual(v, want) {
		t.Errorf("valid patch: returned %v want %v", v, want)
	}

	_, err = client.Patch(setupAuth(), "test/uri", nil, &v)
	if err != nil {
		t.Errorf("valid patch: unexpected %v", err)
	}

	client.Protocol = "http+_"
	_, err = client.Patch(setupAuth(), "test/uri", nil, &v)
	if err == nil {
		t.Errorf("invalid patch: unexpected %v", err)
	}
}

func TestPost(t *testing.T) {
	client := setupClient(generalResponse)

	var v m
	_, err := client.Post(setupAuth(), "test/uri", &m{Weavc: "crusch", One: "1"}, &v)
	if err != nil {
		t.Errorf("valid post: unexpected %v", err)
	}

	want := m{
		Weavc: "crusch",
		One:   "1",
	}

	if !reflect.DeepEqual(v, want) {
		t.Errorf("valid post: returned %v want %v", v, want)
	}

	_, err = client.Post(setupAuth(), "test/uri", nil, &v)
	if err != nil {
		t.Errorf("valid post: unexpected %v", err)
	}

	client.Protocol = "http+_"
	_, err = client.Post(setupAuth(), "test/uri", nil, &v)
	if err == nil {
		t.Errorf("invalid post: unexpected %v", err)
	}
}

func TestParseQuery(t *testing.T) {
	var s string = "one=1&weavc=crusch"

	query, err := parseQuery(s)
	if err != nil {
		t.Errorf("valid string query: unexpected %v", err)
	}
	if query != s {
		t.Errorf("valid string query: returned %s want %s", query, s)
	}

	query, err = parseQuery(&m{Weavc: "crusch", One: "1"})
	if err != nil {
		t.Errorf("valid struct query: unexpected %v", err)
	}
	if query != s {
		t.Errorf("valid struct query: returned %s want %s", query, s)
	}

	query, err = parseQuery(nil)
	if err != nil {
		t.Errorf("valid nil query: unexpected %v", err)
	}
	if query != "" {
		t.Errorf("valid nil query: returned %s want %s", query, "")
	}

	_, err = parseQuery(123)
	if err == nil {
		t.Errorf("invalid number query: unexpected nil err")
	}

	var v map[interface{}]interface{} = map[interface{}]interface{}{
		"weavc": map[interface{}]interface{}{"crusch": "1"},
	}
	_, err = parseQuery(v)
	if err == nil {
		t.Errorf("invalid struct query: unexpected nil err")
	}
}

// the following methods are used for getting and setting different objects for testing purposes

func setupClient(body interface{}) *Client {
	client := NewGithubClient("doesnt.matter", "http")

	var t http.RoundTripper = &testTransport{body: body}
	httpClient := &http.Client{Transport: t}
	client.SetHTTPClient(httpClient)

	return client
}

func setupAuth() Authorizer {
	return AuthorizerFunc(func() (string, error) {
		return "bearer randombearertokenexample", nil
	})
}

type testTransport struct {
	body interface{}
}

func (t *testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	body, err := jsonifyBody(t.body)
	if err != nil {
		return nil, err
	}
	return &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       ioutil.NopCloser(body),
		Request:    req,
	}, nil
}
