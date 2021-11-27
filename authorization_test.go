package crusch

import (
	"crypto/rsa"
	"fmt"
	"testing"
)

func TestRSAPrivateKeyFromPEMFile(t *testing.T) {

	_, err := RSAPrivateKeyFromPEMFile("random_key.pem")
	if err != nil {
		t.Errorf("valid file: returned error: %v", err)
	}

	_, err = RSAPrivateKeyFromPEMFile("random_key_non_existant.pem")
	if err == nil {
		t.Errorf("invalid file: returned no error")
	}

	_, err = RSAPrivateKeyFromPEMFile("readme.md")
	if err == nil {
		t.Errorf("valid file invalid key: returned no error")
	}
}

func TestApplicationAuthorizer(t *testing.T) {
	key := getKey()
	auth, err := NewApplicationAuth(6000, key)
	if err != nil {
		t.Errorf("application auth: unexpected %v", err)
	}
	defer auth.Dispose()

	if auth.ApplicationID != 6000 {
		t.Errorf(
			"application auth: application id %v want %v",
			auth.ApplicationID, 6000)
	}

	if auth.Key != key {
		t.Errorf(
			"application auth: application id %v want %v",
			auth.ApplicationID, 6000)
	}

	_, err = auth.GetHeader()
	if err != nil {
		t.Errorf("application auth: unexpected %v", err)
	}

	// test jwt claims etc
}

func TestOAuthAuthorizer(t *testing.T) {
	token := "123456789abcd"
	auth, err := NewOAuth(token)
	if err != nil {
		t.Errorf("oauth auth: unexpected %v", err)
	}
	defer auth.Dispose()

	if auth.Token != token {
		t.Errorf(
			"oauth auth: token %v want %v",
			auth.Token, token)
	}

	h, err := auth.GetHeader()
	if err != nil {
		t.Errorf("oauth auth: unexpected %v", err)
	}

	if h != fmt.Sprintf("bearer %s", token) {
		t.Errorf(
			"oauth auth: header %v want %v",
			auth.Token, fmt.Sprintf("bearer %s", token))
	}
}

func TestInstallationAuthorizer(t *testing.T) {
	var appID int64 = 123456
	var installationID int64 = 678903

	type tokenResponse struct {
		Token string `json:"token"`
	}

	client := setupClient(tokenResponse{Token: "testtokenstring"})
	key := getKey()
	auth, err := NewInstallationAuth(appID, installationID, key)
	auth.Client = client
	if err != nil {
		t.Errorf("installation auth: unexpected %v", err)
	}
	defer auth.Dispose()

	if auth.ApplicationID != appID {
		t.Errorf(
			"installation auth: applicationID %v want %v",
			auth.ApplicationID, appID)
	}

	if auth.InstallationID != installationID {
		t.Errorf(
			"installation auth: applicationID %v want %v",
			auth.ApplicationID, appID)
	}

	if auth.Key != key {
		t.Errorf(
			"installation auth: key %v want %v",
			auth.Key, key)
	}

	h, err := auth.GetHeader()
	if err != nil {
		t.Errorf("installation auth: unexpected %v", err)
	}

	if h != "token testtokenstring" {
		t.Errorf("installation auth: returned %v want %s", h, "token testtokenstring")
	}

	// this time should go thorugh lastUsed
	_, err = auth.GetHeader()
	if err != nil {
		t.Errorf("installation auth: unexpected %v", err)
	}
}

func getKey() *rsa.PrivateKey {
	key, err := RSAPrivateKeyFromPEMFile("random_key.pem")
	if err != nil {
		panic(err)
	}
	return key
}
