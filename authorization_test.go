package crusch

import (
	"os"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

// TestApplication tests application authorization specific methods
// but does not do full integration tests with Githubs api
func TestApplication(t *testing.T) {
	var appid int64 = 123456
	key, err := RSAPrivateKeyFromPEMFile("random_key.pem")
	if err != nil {
		t.Log("Failed to decode pem file: ", err)
		t.FailNow()
	}

	client := NewDefault()
	client.NewApplicationAuth(appid, key)

	auth, err := client.Authorization()
	if err != nil {
		t.Log("Failed to generate authorization header: ", err)
		t.FailNow()
	}

	authToken := strings.TrimLeft(auth, "Bearer")
	authToken = strings.TrimSpace(authToken)

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(authToken, claims, func(t *jwt.Token) (interface{}, error) {
		return &key.PublicKey, nil
	})
	if err != nil {
		t.Log("Failed jwt parse: ", err)
		t.FailNow()
	}

	_ = token

}

// TestApplicationIntegration will do full integration testing with Githubs API
// requires private_key and application_id environment variables to be set and valid
func TestApplicationIntegration(t *testing.T) {
	privateKey := os.Getenv("private_key")
	applicationID := os.Getenv("application_id")

	if len(strings.TrimSpace(privateKey)) == 0 ||
		len(strings.TrimSpace(applicationID)) == 0 {
		t.Skip("private_key and application_id environment variables expected")
		return
	}
}

// TestInstallation tests installation authorization specific methods
// but does not do full integration tests with Githubs api
// since this requires Githubs api to generate a valid installation token it doesn't do much
func TestInstallation(t *testing.T) {
	var appid int64 = 123456
	var instid int64 = 45678910
	key, err := RSAPrivateKeyFromPEMFile("random_key.pem")
	if err != nil {
		t.Log("Failed to decode pem file: ", err)
		t.FailNow()
	}

	client := NewDefault()
	client.NewInstallationAuth(appid, instid, key)
}

// TestInstallationIntegration will do full integration testing with Githubs API
// requires private_key, application_id and installation_id environment variables to be set and valid
func TestInstallationIntegration(t *testing.T) {
	privateKey := os.Getenv("private_key")
	applicationID := os.Getenv("application_id")
	installationID := os.Getenv("installation_id")

	if len(strings.TrimSpace(privateKey)) == 0 ||
		len(strings.TrimSpace(applicationID)) == 0 ||
		len(strings.TrimSpace(installationID)) == 0 {
		t.Skip("private_key, application_id and installation_id environment variables expected")
		return
	}
}

// TestOauth tests oauth authorization specific methods
// but does not do full integration tests with Githubs api
// since this requires Githubs api to validate a token, it doesn't do much
func TestOauth(t *testing.T) {

	client := NewDefault()
	client.NewOAuth("qwertyuiop", "")

	header, err := client.Authorization()
	if err != nil {
		t.Log("Failed to generate authorization header: ", err)
		t.FailNow()
	}

	if header != "bearer qwertyuiop" {
		t.Log("Incorrect header created: ", header)
		t.FailNow()
	}

	client.NewOAuth("asdfghjkl", "token")
	header, err = client.Authorization()
	if err != nil {
		t.Log("Failed to generate authorization header: ", err)
		t.FailNow()
	}

	if header != "token asdfghjkl" {
		t.Log("Incorrect header created: ", header)
		t.FailNow()
	}

}

// TestOauthIntegration will do full integration testing with Githubs API
// requires oauth_token environment variable to be set and valid
func TestOauthIntegration(t *testing.T) {
	oauthToken := os.Getenv("oauth_token")

	if len(strings.TrimSpace(oauthToken)) == 0 {
		t.Skip("oauth_token environment variables expected")
		return
	}

}
