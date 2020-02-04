package crusch

import (
	"os"
	"strings"
	"testing"
)

func TestClient(t *testing.T) {
	envKey := false
	key := os.Getenv("private_key")
	if len(strings.TrimSpace(key)) != 0 {
		envKey = true
	}

	client1 := New("test_1", "test_1.github.com")
	client2 := New("test_2", "test_2.github.com")
	client3 := New("test_2", "test_3.github.com")
	clientDefault := NewDefault()

	if client1.Name != "test_1" || client1.BaseURL != "test_1.github.com" {
		t.Error("Invalid client1: ", client1)
	}

	if client2.Name != "test_2" || client2.BaseURL != "test_2.github.com" {
		t.Error("Invalid client2: ", client2)
	}

	if clientDefault.Name != "default" || clientDefault.BaseURL != "api.github.com" {
		t.Error("Invalid default client: ", clientDefault)
	}

	client2Get := Pool.Get("test_2")
	if client2Get.Name != "test_2" || client2Get.BaseURL != "test_2.github.com" {
		t.Error("Invalid client2Get: ", client2Get)
	}

	client2.Name = "test_2_updated"

	if client2Get.Name != "test_2_updated" {
		t.Error("Didnt update client2Get: ", client2Get)
	}

	client3.NewOAuth("testing_123")
	clientDefault.NewApplicationAuthFile(1234567, "random_key.pem")
	if envKey {
		client2Get.NewInstallationAuthBytes(123456, 12345678, []byte(key))
	} else {
		client2Get.NewInstallationAuthFile(123456, 12345678, "random_key.pem")
	}

	applicationClient := Pool.GetByApplicationAuth(1234567)
	oauthClient := Pool.GetByOauthToken("testing_123")
	installationClient := Pool.GetByInstallationAuth(123456, 12345678)

	if applicationClient.Auth.AuthType != Application ||
		applicationClient.Name != "default" ||
		applicationClient.Auth.ApplicationID != 1234567 ||
		applicationClient.Auth.InstallationID != 0 {
		t.Error("application client invalid: ", applicationClient)
	}

	if oauthClient.Auth.AuthType != OAuth ||
		oauthClient.Name != "test_2" ||
		oauthClient.Auth.ApplicationID != 0 ||
		oauthClient.Auth.InstallationID != 0 ||
		oauthClient.Auth.OAuthAccessToken != "testing_123" {
		t.Error("oauth client invalid: ", oauthClient)
	}

	if installationClient.Auth.AuthType != Installation ||
		installationClient.Name != "test_2_updated" ||
		installationClient.Auth.ApplicationID != 123456 ||
		installationClient.Auth.InstallationID != 12345678 {
		t.Error("installation client invalid: ", installationClient)
	}

	applicationClient = Pool.GetByApplicationAuth(234567)
	if applicationClient != nil {
		t.Error("application client should be nil: ", applicationClient)
	}

	installationClient = Pool.GetByInstallationAuth(234567, 6969)
	if installationClient != nil {
		t.Error("installation client should be nil: ", installationClient)
	}

	oauthClient = Pool.GetByOauthToken("40040400400400400000dkkk")
	if oauthClient != nil {
		t.Error("oauth client should be nil: ", oauthClient)
	}
}
