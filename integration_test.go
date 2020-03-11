package crusch

import (
	"os"
	"strconv"
	"strings"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
)

// These methods do basic authorization testing with githubs API
// It will use private_key, application_id and installation_id env variables
// to get the required values for the test, if these aren't provided the tests will be skipped

func TestApplication(t *testing.T) {
	PrivateKey := os.Getenv("private_key")
	AppID := os.Getenv("application_id")

	if len(strings.TrimSpace(AppID)) <= 0 ||
		len(strings.TrimSpace(PrivateKey)) <= 0 {
		t.Skip("application_id and private_key environment variables must be provided for this test")
	}

	bytekey := []byte(PrivateKey)
	key, err := jwt.ParseRSAPrivateKeyFromPEM(bytekey)
	if err != nil {
		t.Errorf("%v", err)
	}

	id, err := strconv.ParseInt(AppID, 10, 64)
	if err != nil {
		t.Errorf("%v", err)
	}

	a, err := NewApplicationAuth(id, key)

	var v map[string]interface{}
	res, err := GithubClient.Get(a, "app", nil, &v)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Request error: %d", res.StatusCode)
	}
}

func TestInstallation(t *testing.T) {
	PrivateKey := os.Getenv("private_key")
	AppID := os.Getenv("application_id")
	InstallationID := os.Getenv("installation_id")

	if len(strings.TrimSpace(AppID)) <= 0 ||
		len(strings.TrimSpace(PrivateKey)) <= 0 ||
		len(strings.TrimSpace(InstallationID)) <= 0 {
		t.Skip("application_id, installation_id and private_key environment variables must be provided for this test")
	}

	bytekey := []byte(PrivateKey)
	key, err := jwt.ParseRSAPrivateKeyFromPEM(bytekey)
	if err != nil {
		t.Errorf("%v", err)
	}

	id, err := strconv.ParseInt(AppID, 10, 64)
	if err != nil {
		t.Errorf("%v", err)
	}

	iid, err := strconv.ParseInt(InstallationID, 10, 64)
	if err != nil {
		t.Errorf("%v", err)
	}

	a, err := NewInstallationAuth(id, iid, key)

	var v map[string]interface{}
	res, err := GithubClient.Get(a, "installation/repositories", nil, &v)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Request error: %d", res.StatusCode)
		t.Log(res)
	}
}
