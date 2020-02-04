package crusch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/google/go-querystring/query"
)

// GetJson executes a GET http request to given uri and params and binds the response to v
// params must be a Struct, string or nil
func (s *Client) GetJson(uri string, params interface{}, v interface{}) (*http.Response, error) {

	res, err := s.Get(uri, params)
	if err != nil {
		return res, err
	}

	if v != nil {
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(v)
		if err != nil {
			return res, err
		}
	}

	return res, err
}

// Get executes a GET http request to given uri and params
// params must be a Struct, string or nil
func (s *Client) Get(uri string, params interface{}) (*http.Response, error) {

	var p string
	switch v := params.(type) {
	case string:
		p = v
	case nil:
		p = ""
	default:
		val := reflect.ValueOf(params)
		if val.Kind() != reflect.Struct {
			parsed, err := ParseQuery(params)
			if err != nil {
				return nil, fmt.Errorf("Error parsing params: %s", err)
			}
			p = parsed
		} else {
			return nil, fmt.Errorf("Unknown type of params, only takes string, Struct or nil")
		}
	}

	url := fmt.Sprintf("https://%s%s?%s", s.BaseURL, uri, p)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %s", err)
	}

	auth, err := s.Authorization()
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/vnd.github.machine-man-preview+json")

	res, err := http.DefaultClient.Do(req)

	return res, err
}

// DeleteJson executes a Delete http request to given uri and binds the response to v
func (s *Client) DeleteJson(uri string, v interface{}) (*http.Response, error) {
	res, err := s.Delete(uri)
	if err != nil {
		return res, err
	}

	if v != nil {
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(v)
		if err != nil {
			return res, fmt.Errorf("Error formatting json to v: %s", err)
		}
	}

	return res, err
}

// Delete executes a Delete http request to given uri
func (s *Client) Delete(uri string) (*http.Response, error) {
	url := fmt.Sprintf("https://%s%s", s.BaseURL, uri)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %s", err)
	}

	auth, err := s.Authorization()
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/vnd.github.machine-man-preview+json")

	res, err := http.DefaultClient.Do(req)

	return res, err
}

// PostJson executes a POST http request to given uri and binds the response to v
// body will be encoded to json, use nil for no body
func (s *Client) PostJson(uri string, body interface{}, v interface{}) (*http.Response, error) {

	b, err := JsonBody(body)
	if err != nil {
		return nil, fmt.Errorf("Error parsing body to json: %s", err)
	}

	res, err := s.Post(uri, b)
	if err != nil {
		return res, err
	}

	if v != nil {
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(v)
		if err != nil {
			return res, fmt.Errorf("Error formatting json to v: %s", err)
		}
	}

	return res, err
}

// Post executes a POST http request to given uri and binds the response to v
// body will be encoded to json, use nil for no body
func (s *Client) Post(uri string, body *bytes.Buffer) (*http.Response, error) {

	url := fmt.Sprintf("https://%s%s", s.BaseURL, uri)

	var req *http.Request
	if body != nil {
		r, err := http.NewRequest("POST", url, body)
		if err != nil {
			return nil, fmt.Errorf("Error creating request: %s", err)
		}
		req = r
	} else {
		r, err := http.NewRequest("POST", url, nil)
		if err != nil {
			return nil, fmt.Errorf("Error creating request: %s", err)
		}
		req = r
	}

	auth, err := s.Authorization()
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/vnd.github.machine-man-preview+json")

	res, err := http.DefaultClient.Do(req)

	return res, err
}

// PatchJson executes a PATCH http request to given uri and binds the response to v
// body will be encoded to json, use nil for no body
func (s *Client) PatchJson(uri string, body interface{}, v interface{}) (*http.Response, error) {
	b, err := JsonBody(body)
	if err != nil {
		return nil, fmt.Errorf("Error parsing body to json: %s", err)
	}

	res, err := s.Patch(uri, b)
	if err != nil {
		return res, err
	}

	if v != nil {
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(v)
		if err != nil {
			return res, fmt.Errorf("Error formatting json to v: %s", err)
		}
	}

	return res, err
}

// Patch executes a PATCH http request to given uri and binds the response to v
// body will be encoded to json, use nil for no body
func (s *Client) Patch(uri string, body *bytes.Buffer) (*http.Response, error) {

	url := fmt.Sprintf("https://%s%s", s.BaseURL, uri)

	var req *http.Request
	if body != nil {
		r, err := http.NewRequest("POST", url, body)
		if err != nil {
			return nil, fmt.Errorf("Error creating request: %s", err)
		}
		req = r
	} else {
		r, err := http.NewRequest("POST", url, nil)
		if err != nil {
			return nil, fmt.Errorf("Error creating request: %s", err)
		}
		req = r
	}

	auth, err := s.Authorization()
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/vnd.github.machine-man-preview+json")

	res, err := http.DefaultClient.Do(req)

	return res, err
}

// PutJson executes a PUT http request to given uri and binds the response to v
// body will be encoded to json, use nil for no body
func (s *Client) PutJson(uri string, body interface{}, v interface{}) (*http.Response, error) {
	b, err := JsonBody(body)
	if err != nil {
		return nil, fmt.Errorf("Error parsing body to json: %s", err)
	}

	res, err := s.Put(uri, b)
	if err != nil {
		return res, err
	}

	if v != nil {
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(v)
		if err != nil {
			return res, fmt.Errorf("Error formatting json to v: %s", err)
		}
	}

	return res, err
}

// Put executes a PUT http request to given uri and binds the response to v
// body will be encoded to json, use nil for no body
func (s *Client) Put(uri string, body *bytes.Buffer) (*http.Response, error) {

	url := fmt.Sprintf("https://%s%s", s.BaseURL, uri)

	var req *http.Request
	if body != nil {
		r, err := http.NewRequest("POST", url, body)
		if err != nil {
			return nil, fmt.Errorf("Error creating request: %s", err)
		}
		req = r
	} else {
		r, err := http.NewRequest("POST", url, nil)
		if err != nil {
			return nil, fmt.Errorf("Error creating request: %s", err)
		}
		req = r
	}

	auth, err := s.Authorization()
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/vnd.github.machine-man-preview+json")

	res, err := http.DefaultClient.Do(req)

	return res, err
}

// Do takes a *http.Request adds authorization headers and performs the request
// useful for requests that don't follow Githubs usual headers and body etc i.e. reactions
func (s *Client) Do(req *http.Request) (*http.Response, error) {
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

// JsonBody takes a model/struct and attempts to make it into a *bytes.Buffer
func JsonBody(body interface{}) (*bytes.Buffer, error) {

	if body == nil {
		return nil, nil
	}

	var buf io.ReadWriter
	buf = &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(body)
	if err != nil {
		return nil, err
	}

	buffer, ok := buf.(*bytes.Buffer)
	if !ok {
		return nil, fmt.Errorf("Failed to create buffer")
	}

	return buffer, nil
}

// ParseQuery turns a struct into a url encoded query string
func ParseQuery(opt interface{}) (string, error) {
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
