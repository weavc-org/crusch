package crusch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

// GET makes get requests to Githubs API, using the Client for authorization and other details
func (s *Client) GET(uri string, params interface{}, v interface{}) (*http.Response, error) {
	res, err := get(s, uri, params)
	if err != nil {
		return res, err
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(v)
	if err != nil {
		log.Print("decode error: ", err)
		return res, err
	}

	return res, err
}

func get(s *Client, uri string, params interface{}) (*http.Response, error) {
	req, err := baseRequest(s, uri)
	if err != nil {
		return nil, err
	}

	req.Method = http.MethodGet

	if params != nil {
		q, err := parseQuery(params)
		if err != nil {
			return nil, err
		}
		req.URL.RawQuery = q
	}

	res, err := http.DefaultClient.Do(req)

	return res, err
}

// DELETE makes delete requests to Githubs API, using the Client for authorization and other details
func (s *Client) DELETE(uri string, params interface{}, v interface{}) (*http.Response, error) {
	res, err := delete(s, uri, params)
	if err != nil {
		return res, err
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(v)
	if err != nil {
		return res, err
	}

	return res, err
}

func delete(s *Client, uri string, params interface{}) (*http.Response, error) {
	req, err := baseRequest(s, uri)
	if err != nil {
		return nil, err
	}

	req.Method = http.MethodDelete

	if params != nil {
		q, err := parseQuery(params)
		if err != nil {
			return nil, err
		}
		req.URL.RawQuery = q
	}

	res, err := http.DefaultClient.Do(req)

	return res, err
}

// POST makes post requests to Githubs API, using Client for authorization and other details
func (s *Client) POST(uri string, body interface{}, v interface{}) (*http.Response, error) {
	res, err := post(s, uri, body)
	if err != nil {
		return res, err
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(v)
	if err != nil {
		return res, err
	}

	return res, err
}

func post(s *Client, uri string, body interface{}) (*http.Response, error) {
	req, err := baseRequest(s, uri)
	if err != nil {
		return nil, err
	}

	req.Method = http.MethodPost

	if body != nil {
		err = setBody(req, body)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := http.DefaultClient.Do(req)

	return res, err
}

// PATCH makes patch requests to Githubs API, using Client for authorization and other details
func (s *Client) PATCH(uri string, body interface{}, v interface{}) (*http.Response, error) {
	res, err := patch(s, uri, body)
	if err != nil {
		return res, err
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(v)
	if err != nil {
		return res, err
	}

	return res, err
}

func patch(s *Client, uri string, body interface{}) (*http.Response, error) {
	req, err := baseRequest(s, uri)
	if err != nil {
		return nil, err
	}

	req.Method = http.MethodPatch

	if body != nil {
		err = setBody(req, body)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := http.DefaultClient.Do(req)

	return res, err
}

// PUT makes put requests to Githubs API, using Client for authorization and other details
func (s *Client) PUT(uri string, body interface{}, v interface{}) (*http.Response, error) {
	res, err := put(s, uri, body)
	if err != nil {
		return res, err
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(v)
	if err != nil {
		return res, err
	}

	return res, err
}

func put(s *Client, uri string, body interface{}) (*http.Response, error) {
	req, err := baseRequest(s, uri)
	if err != nil {
		return nil, err
	}

	req.Method = http.MethodPut

	if body != nil {
		err = setBody(req, body)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := http.DefaultClient.Do(req)

	return res, err
}

// DO will add authentication headers and execute the given http request
// useful for when you need more control over the request than the other methods offer
func (s *Client) DO(req *http.Request) (*http.Response, error) {
	if req.Header == nil {
		req.Header = http.Header{}
	}

	auth, err := s.Authorization()
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", auth)

	res, err := http.DefaultClient.Do(req)

	return res, err
}

func setBody(req *http.Request, body interface{}) error {
	if body == nil {
		return fmt.Errorf("Body cannot be empty")
	}
	var buf io.ReadWriter
	buf = &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(body)
	if err != nil {
		return err
	}

	r, ok := buf.(io.Reader)
	if ok {
		rc, ok := r.(io.ReadCloser)
		if !ok && buf != nil {
			rc = ioutil.NopCloser(buf)
		}
		req.Body = rc
	}

	buffer, ok := r.(*bytes.Buffer)
	if !ok {
		return fmt.Errorf("Cannot create buffer")
	}
	req.ContentLength = int64(buffer.Len())
	b := buffer.Bytes()
	req.GetBody = func() (io.ReadCloser, error) {
		r := bytes.NewReader(b)
		return ioutil.NopCloser(r), nil
	}

	return nil
}

func parseQuery(opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return "", nil
	}

	qs, err := query.Values(opt)
	if err != nil {
		return "", err
	}

	return qs.Encode(), nil
}

func baseRequest(s *Client, uri string) (*http.Request, error) {
	req := http.Request{}

	auth, err := s.Authorization()
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{}

	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/vnd.github.machine-man-preview+json")
	req.URL = &url.URL{Scheme: "https", Host: s.BaseURL, Path: uri}

	return &req, nil
}
