package crusch

import (
	"fmt"
	"strings"
)

// Client for making http requests to Githubs v3 json api
type Client struct {
	Name     string
	BaseURL  string
	Protocol string
	Auth     *Authorization
}

// NewDefault creates a new Client with default values and no auth
func NewDefault() *Client {
	return New("default", "api.github.com")
}

// New creates a new Client struct using given arguments
func New(name string, baseURL string) *Client {
	s := Client{
		Name:     name,
		BaseURL:  baseURL,
		Protocol: "https",
	}
	Pool.Pool = append(Pool.Pool, &s)
	return &s
}

// GetURL looks at values in client and formats the full url
// this method isn't exactly fool-proof but it should kind of work...
func (s *Client) GetURL() string {
	// sets defaults
	protocol := "https"
	url := "api.github.com"

	// if there is a protocol set
	if len(strings.TrimSpace(s.Protocol)) != 0 {
		// if the set protocol is 'http' or 'https' use it
		if s.Protocol == "http" || s.Protocol == "https" {
			protocol = s.Protocol
		}
	}

	// check if baseurl has value, use it
	if len(strings.TrimSpace(s.BaseURL)) != 0 {
		url = s.BaseURL
	}

	// if there isn't a prefix already set on baseurl
	if !strings.HasPrefix(s.BaseURL, "http:") &&
		!strings.HasPrefix(s.BaseURL, "https:") {
		url = fmt.Sprintf("%s://%s", protocol, url)
	}

	return url
}

// Dispose of client and references to it
func (s *Client) Dispose() *Client {
	i, err := Pool.getClientIndex(s)
	if err != nil {
		Pool.Pool = append(Pool.Pool[:i], Pool.Pool[i+1:]...)
	}

	s.Auth = nil
	s.BaseURL = ""
	s.Name = ""
	s = nil

	return s
}

// ClientPool stores a slice of Clients created in this session
// Clients are automatically added when created using New or NewDefault
type ClientPool struct {
	Pool []*Client
}

var (
	// Pool is an array of Clients used within this session
	Pool ClientPool = ClientPool{}
)

// GetByInstallationAuth tries to find an existing installation client that matches the auth details
func (cp *ClientPool) GetByInstallationAuth(applicationID int64, installationID int64) *Client {
	for _, Client := range cp.Pool {
		if Client.Auth != nil && Client.Auth.AuthType == Installation &&
			Client.Auth.ApplicationID == applicationID &&
			Client.Auth.InstallationID == installationID {
			return Client
		}
	}
	return nil
}

// GetByInstallationAuth tries to find an existing installation client that matches the auth details
func (cp *ClientPool) getClientIndex(c *Client) (int, error) {
	for i, client := range cp.Pool {
		if c == client {
			return i, nil
		}
	}
	return 0, fmt.Errorf("No client found")
}

// GetByApplicationAuth tries to find an existing application client that matches the auth details
func (cp *ClientPool) GetByApplicationAuth(applicationID int64) *Client {
	for _, Client := range cp.Pool {
		if Client.Auth != nil && Client.Auth.AuthType == Application &&
			Client.Auth.ApplicationID == applicationID &&
			Client.Auth.InstallationID == 0 {
			return Client
		}
	}
	return nil
}

// GetByOauthToken tries to find an existing application client that matches the auth details
func (cp *ClientPool) GetByOauthToken(token string) *Client {
	for _, Client := range cp.Pool {
		if Client.Auth != nil &&
			Client.Auth.AuthType == OAuth &&
			Client.Auth.OAuthAccessToken == token {
			return Client
		}
	}
	return nil
}

// Get will find a client by the given name,
// if multiple are in the array, only the first will be returned
func (cp *ClientPool) Get(name string) *Client {
	for _, Client := range cp.Pool {
		if Client.Name == name {
			return Client
		}
	}
	return nil
}
