package crusch

import (
	"net/http"
	"testing"
)

func TestAttachAuthorizer(t *testing.T) {
	var rt http.RoundTripper = &testTransport{body: "test 123"}
	c := &http.Client{Transport: rt}

	a, _ := NewOAuth("testingoauthroundtripperwrapper")

	err := AttachAuthorizer(a, c)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	req, _ := http.NewRequest("GET", "test.com", nil)
	res, err := c.Do(req)

	s := res.Request.Header.Get("Authorization")
	if s != "bearer testingoauthroundtripperwrapper" {
		t.Errorf("returned %s want bearer testingoauthroundtripperwrapper", s)
	}
}
