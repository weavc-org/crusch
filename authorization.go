package crusch

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Authorizer implements a GetHeader() method which
// returns the value used inside the authorization header i.e. bearer token
// It is used inside client request methods and provides the request with authorization headers
type Authorizer interface {
	GetHeader() (string, error)
}

// AuthorizerFunc is a wrapper for the Authorizer interface
// it follows the a similar design to http.Handler
type AuthorizerFunc func() (string, error)

// GetHeader wraps AuthorizerFunc, implementing the Authorizer interface
func (a AuthorizerFunc) GetHeader() (string, error) {
	return a()
}

// NewApplicationAuth generates and returns a new ApplicationAuth struct using given values
func NewApplicationAuth(applicationID int64, key *rsa.PrivateKey) (*ApplicationAuth, error) {
	a := &ApplicationAuth{
		ApplicationID: applicationID,
		Key:           key,
	}

	return a, nil
}

// NewInstallationAuth generates a new ApplicationAuth structure using given values
func NewInstallationAuth(applicationID int64, installationID int64, key *rsa.PrivateKey) (*InstallationAuth, error) {
	a := &InstallationAuth{
		ApplicationID:  applicationID,
		InstallationID: installationID,
		Key:            key,
		Client:         GithubClient,
	}

	return a, nil
}

// NewOAuth generates and returns an OAuth authorizer that uses given values
func NewOAuth(token string) (*OAuth, error) {
	a := &OAuth{Token: token}
	return a, nil
}

// ApplicationAuth creates authorization headers based on the given ApplicationID and private key
// Both are provided by Github, see: https://developer.github.com/v3/apps/#get-the-authenticated-github-app
type ApplicationAuth struct {
	ApplicationID int64
	Key           *rsa.PrivateKey
}

// Dispose of values in ApplicationAuth struct
func (a *ApplicationAuth) Dispose() {
	a.ApplicationID = 0
	a.Key = nil
}

// GetHeader to implement Authorizer
// GetHeader generates a new JWT token using the ApplicationID and PEM from Github
// This header is used for authenticating a Github application against Githubs api
func (a *ApplicationAuth) GetHeader() (string, error) {
	claims := &jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute).Unix(),
		Issuer:    strconv.FormatInt(a.ApplicationID, 10),
	}

	bearer := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signed, err := bearer.SignedString(a.Key)
	if err != nil {
		return "", fmt.Errorf("could not sign jwt: %s", err)
	}

	return fmt.Sprintf("bearer %s", signed), nil
}

// InstallationAuth allows applications to authenticate as an installation against Githubs API
// Uses the ApplcationID and Key to generate an ApplicationAuth which is inturn used to get
// installation token from githubs api. This will also store the Last used header and time instead of
// getting a new token for each request.
// https://developer.github.com/v3/apps/#create-a-new-installation-token
type InstallationAuth struct {
	ApplicationID  int64
	InstallationID int64
	Key            *rsa.PrivateKey
	Client         *Client
	LastUsed
}

// Dispose of values in InstallationAuth
func (a *InstallationAuth) Dispose() {
	a.ApplicationID = 0
	a.InstallationID = 0
	a.Key = nil
	a.Client = nil
	a.header = ""
	a.validUntil = 0
	a.time = 0
}

// GetHeader to implement Authorizer
// This produces the required auth headers for the installation
// It will make a request to the Clients API (by default Github) to get an access token
// for the application and installation IDs provided by InstallationAuth
// If InstallationAuth has already generated an auth token and it is still valid, this will be used instead
// https://developer.github.com/v3/apps/#create-a-new-installation-token
func (a *InstallationAuth) GetHeader() (string, error) {

	if time.Now().Unix() <= a.validUntil && a.header != "" {
		return a.header, nil
	}

	var client *Client

	// use client attached to auth, or the default GithubClient
	if a.Client != nil {
		client = a.Client
	} else {
		client = GithubClient
	}

	auth, err := NewApplicationAuth(a.ApplicationID, a.Key)
	if err != nil {
		return "", fmt.Errorf("unable to create application authorizer: %v", err)
	}

	var v map[string]interface{}
	res, err := client.Post(
		auth,
		fmt.Sprintf("app/installations/%d/access_tokens", a.InstallationID),
		nil,
		&v,
	)

	if err != nil {
		return "", err
	}

	if res.StatusCode != 201 && res.StatusCode != 200 {
		return "", fmt.Errorf("%d error when trying to create access token", res.StatusCode)
	}

	t, ok := v["token"].(string)
	if !ok {
		return "", fmt.Errorf("error mapping token value")
	}

	a.header = fmt.Sprintf("token %s", t)
	a.validUntil = time.Now().Add((time.Hour - time.Minute)).Unix()
	a.time = time.Now().Unix()

	return a.header, nil
}

// OAuth authorizor for Githubs v3 API
// https://developer.github.com/v3/#authentication
type OAuth struct {
	Token string
}

// Dispose of values stored inside of OAuth
func (a *OAuth) Dispose() {
	a.Token = ""
}

// GetHeader to implement Authorizer
// returns the Authorization header required for oauth authentication
func (a *OAuth) GetHeader() (string, error) {
	return fmt.Sprintf("bearer %s", a.Token), nil
}

// RSAPrivateKeyFromPEMFile produces a *rsa.PrivateKey from a .pem file
// The .pem file is provided by Github to authorize your application against their API
func RSAPrivateKeyFromPEMFile(keyfile string) (*rsa.PrivateKey, error) {
	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}
	rsakey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		return nil, err
	}
	return rsakey, nil
}

// LastUsed details when a token/header was last used
// and when it needs to be reused
type LastUsed struct {
	header     string
	validUntil int64
	time       int64
}
