package crusch

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	// Application authentication type
	// This type works using the ApplicationID and
	// private kd by Github
	Application AuthType = 0
	// Installation authentication type
	// Similar to application but authorizes as a specific
	// application installion using the ApplicationID, InstallationID
	// and privatovided by Github
	Installation AuthType = 1
	// OAuth authentication type
	// Just uses an OAuth token to authorize application as a specific user
	OAuth AuthType = 2
)

// AuthType allows us to identify the authorization type to use when
// generating authorization headers
type AuthType int

// Authorization for github requests
type Authorization struct {
	AuthType         AuthType
	ApplicationID    int64
	InstallationID   int64
	OAuthAccessToken string
	OAuthType        string
	Key              *rsa.PrivateKey
	LastUsed         *LastUsed
}

// LastUsed authentication header & valid until date for this Client
type LastUsed struct {
	AuthHeader string
	ValidUntil int64
	Time       int64
}

// RSAPrivateKeyFromPEMFile returns the *rsa.PrivateKey decoded from given PEM file,
// can be used elsewhere when generating Clients without having to read the file over and over
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

// NewApplicationAuthFile creates new application authentication Client from keyfile
func (s *Client) NewApplicationAuthFile(applicationID int64, keyfile string) {
	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		panic(fmt.Errorf("could not read private key: %s", err))
	}
	s.NewApplicationAuthBytes(applicationID, key)
}

// NewApplicationAuthBytes creates new application authentication Client from byte array
func (s *Client) NewApplicationAuthBytes(applicationID int64, key []byte) {
	rsakey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		panic(err)
	}
	s.NewApplicationAuth(applicationID, rsakey)
}

// NewApplicationAuth creates new application authentication Client
func (s *Client) NewApplicationAuth(applicationID int64, key *rsa.PrivateKey) {
	s.Auth = &Authorization{
		AuthType:      Application,
		ApplicationID: applicationID,
		Key:           key,
	}
}

// NewInstallationAuthFile creates new installation authentication Client from keyfile
func (s *Client) NewInstallationAuthFile(applicationID int64, installationID int64, keyfile string) {
	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		panic(fmt.Errorf("could not read private key: %s", err))
	}
	s.NewInstallationAuthBytes(applicationID, installationID, key)
}

// NewInstallationAuthBytes creates new installation authentication Client from byte array
func (s *Client) NewInstallationAuthBytes(applicationID int64, installationID int64, key []byte) {
	rsakey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		panic(err)
	}
	s.NewInstallationAuth(applicationID, installationID, rsakey)
}

// NewInstallationAuth creates new installation authentication Client
func (s *Client) NewInstallationAuth(applicationID int64, installationID int64, key *rsa.PrivateKey) {
	s.Auth = &Authorization{
		AuthType:       Installation,
		ApplicationID:  applicationID,
		InstallationID: installationID,
		Key:            key,
	}
}

// NewOAuth creates a new OAuth authentication client
func (s *Client) NewOAuth(token string, oauthType string) {
	if len(strings.TrimSpace(oauthType)) == 0 {
		oauthType = "bearer"
	}
	s.Auth = &Authorization{
		AuthType:         OAuth,
		OAuthAccessToken: token,
		OAuthType:        oauthType,
	}
}

// Authorization creates and returns the Clients authorization header
func (s *Client) Authorization() (string, error) {
	if s.Auth.AuthType == Application {
		return applicationAuthorization(s)
	} else if s.Auth.AuthType == Installation {
		return installationAuthorization(s)
	} else if s.Auth.AuthType == OAuth {
		return oauthAuthorization(s)
	}

	return "", fmt.Errorf("auth type not found")
}

func applicationAuthorization(s *Client) (string, error) {
	claims := &jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute).Unix(),
		Issuer:    strconv.FormatInt(s.Auth.ApplicationID, 10),
	}

	bearer := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signed, err := bearer.SignedString(s.Auth.Key)
	if err != nil {
		return "", fmt.Errorf("could not sign jwt: %s", err)
	}

	s.Auth.LastUsed = &LastUsed{
		AuthHeader: "Bearer " + signed,
		ValidUntil: time.Now().Add(time.Minute).Unix(),
		Time:       time.Now().Unix(),
	}

	return s.Auth.LastUsed.AuthHeader, nil
}

func installationAuthorization(s *Client) (string, error) {

	if s.Auth.LastUsed != nil {
		if time.Now().Unix() <= s.Auth.LastUsed.ValidUntil {
			s.Auth.LastUsed.Time = time.Now().Unix()
			return s.Auth.LastUsed.AuthHeader, nil
		}
	}

	var client = Pool.Get(fmt.Sprintf("crusch_application_%d", s.Auth.ApplicationID))
	if client == nil {
		client = New(fmt.Sprintf("crusch_application_%d", s.Auth.ApplicationID), s.BaseURL)
		client.NewApplicationAuth(s.Auth.ApplicationID, s.Auth.Key)
	}

	var v = make(map[string]interface{})
	var uri = fmt.Sprintf("/app/installations/%d/access_tokens", s.Auth.InstallationID)
	res, err := client.PostJson(uri, nil, &v)
	if err != nil {
		return "", err
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		t, ok := v["token"].(string)
		if !ok {
			return "", fmt.Errorf("error mapping result")
		}

		s.Auth.LastUsed = &LastUsed{
			AuthHeader: "token " + t,
			ValidUntil: time.Now().Add((time.Hour - time.Minute)).Unix(),
			Time:       time.Now().Unix(),
		}
		return s.Auth.LastUsed.AuthHeader, nil
	}

	return "", fmt.Errorf("Unable to create authentication token")

}

func oauthAuthorization(s *Client) (string, error) {

	s.Auth.LastUsed = &LastUsed{
		AuthHeader: fmt.Sprintf("%s %s", s.Auth.OAuthType, s.Auth.OAuthAccessToken),
		ValidUntil: 0,
		Time:       time.Now().Unix(),
	}

	return s.Auth.LastUsed.AuthHeader, nil
}
