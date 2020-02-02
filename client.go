package crusch

import "time"

// Client for making http requests to Githubs v3 json api
type Client struct {
	BaseURL string
	Auth    *Authorization
}

// NewDefault creates a new Client with default values and no auth
func NewDefault() *Client {
	return New("api.github.com")
}

// New creates a new Client struct using given arguments
func New(baseURL string) *Client {
	s := Client{BaseURL: baseURL}
	Pool.Pool = append(Pool.Pool, &s)
	return &s
}

// ClientPool stores a slice of Clients created in this session
type ClientPool struct {
	Pool []*Client
}

var (
	// Pool is an array of Clients used within this session
	Pool ClientPool = ClientPool{}
)

// GetByAuth tries to find an existing Client that matches the auth details
func (cp *ClientPool) GetByAuth(Type AuthType, applicationID int64, installationID int64) *Client {
	for i, Client := range cp.Pool {
		if Client.Auth.LastUsed != nil {
			if Client.Auth.LastUsed.Time > time.Now().Add((time.Hour * 24)).Unix() {
				cp.Pool = append(cp.Pool[:i], cp.Pool[i+1:]...)
				continue
			}
		} else {
			cp.Pool = append(cp.Pool[:i], cp.Pool[i+1:]...)
			continue
		}
		if Client.Auth.LastUsed != nil &&
			Client.Auth.AuthType == Type &&
			Client.Auth.ApplicationID == applicationID &&
			Client.Auth.InstallationID == installationID {
			return Client
		}
	}
	return nil
}
