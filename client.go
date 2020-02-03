package crusch

// Client for making http requests to Githubs v3 json api
type Client struct {
	Name    string
	BaseURL string
	Auth    *Authorization
}

// NewDefault creates a new Client with default values and no auth
func NewDefault() *Client {
	return New("default", "api.github.com")
}

// New creates a new Client struct using given arguments
func New(name string, baseURL string) *Client {
	s := Client{Name: name, BaseURL: baseURL}
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
	for _, Client := range cp.Pool {
		if Client.Auth.LastUsed != nil &&
			Client.Auth.AuthType == Type &&
			Client.Auth.ApplicationID == applicationID &&
			Client.Auth.InstallationID == installationID {
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
