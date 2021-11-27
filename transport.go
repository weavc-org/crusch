package crusch

import (
	"fmt"
	"net/http"
)

type transport struct {
	authorizer Authorizer
	rt         http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	h, err := t.authorizer.GetHeader()
	if err != nil {
		return nil, fmt.Errorf("failed to get authorization header: %v", err)
	}

	req.Header.Add("Authorization", h)

	return t.rt.RoundTrip(req)
}

// AttachAuthorizer attaches a new http.Transport layer that adds authorization headers to the request
// this new layer wraps any existing transport layers
// this can be used in conjuction with go-github to provide authorization headers to requests
func AttachAuthorizer(authorizer Authorizer, httpClient *http.Client) error {
	ct := transport{authorizer, httpClient.Transport}

	httpClient.Transport = &ct
	return nil
}
