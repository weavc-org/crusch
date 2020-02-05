package crusch

import (
	"testing"
)

func TestClient(t *testing.T) {

	Pool = ClientPool{}

	client1 := New("test_1", "test_1.github.com")
	client2 := New("test_2", "test_2.github.com")
	client3 := New("test_2", "test_3.github.com")
	_ = client3
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

	clientDefault = clientDefault.Dispose()

	if clientDefault != nil {
		t.Error("clientDefault should be disposed: ", clientDefault)
	}
}

func TestPool(t *testing.T) {
	Pool = ClientPool{}

	new1 := New("test_1", "test_1.github.com")
	new2 := New("test_2", "test_2.github.com")
	_ = new1
	_ = new2

	client2 := Pool.Get("test_2")
	if client2.Name != "test_2" || client2.BaseURL != "test_2.github.com" {
		t.Error("Invalid client2: ", client2)
	}

	client2.Name = "test_2_updated"

	client2Get := Pool.Get("test_2")
	if client2Get != nil {
		t.Error("found test_2 when shouldn't have: ", client2Get)
	}

	oauthClient := NewDefault()
	oauthClient.NewOAuth("testing_123", "")

	applicationClient := NewDefault()
	applicationClient.NewApplicationAuthFile(1234567, "random_key.pem")

	installationClient := NewDefault()
	installationClient.NewInstallationAuthFile(123456, 12345678, "random_key.pem")

	applicationClientGet := Pool.GetByApplicationAuth(1234567)
	oauthClientGet := Pool.GetByOauthToken("testing_123")
	installationClientGet := Pool.GetByInstallationAuth(123456, 12345678)

	if applicationClientGet == nil ||
		applicationClientGet.Auth.AuthType != Application ||
		applicationClientGet.Name != "default" ||
		applicationClientGet.Auth.ApplicationID != 1234567 ||
		applicationClientGet.Auth.InstallationID != 0 {
		t.Error("application client invalid: ", applicationClientGet)
	}

	if oauthClientGet == nil ||
		oauthClientGet.Auth.AuthType != OAuth ||
		oauthClientGet.Auth.ApplicationID != 0 ||
		oauthClientGet.Auth.InstallationID != 0 ||
		oauthClientGet.Auth.OAuthAccessToken != "testing_123" {
		t.Error("oauth client invalid: ", oauthClientGet)
	}

	if installationClientGet == nil ||
		installationClientGet.Auth.AuthType != Installation ||
		installationClientGet.Auth.ApplicationID != 123456 ||
		installationClientGet.Auth.InstallationID != 12345678 {
		t.Error("installation client invalid: ", installationClientGet)
	}

	applicationClientGet = Pool.GetByApplicationAuth(234567)
	if applicationClientGet != nil {
		t.Error("application client should be nil: ", applicationClientGet)
	}

	oauthClientGet = Pool.GetByOauthToken("40040400400400400000dkkk")
	if oauthClientGet != nil {
		t.Error("oauth client should be nil: ", oauthClientGet)
	}

	installationClientGet = Pool.GetByInstallationAuth(234567, 6969)
	if installationClientGet != nil {
		t.Error("installation client should be nil: ", installationClientGet)
	}
}

func TestGetURL(t *testing.T) {

	Pool = ClientPool{}

	client1 := New("1", "test_1.github.com")
	client2 := New("2", "http://test_2.github.com")
	client3 := New("3", "test_3.github.com")
	clientDefault := NewDefault()

	url1 := client1.GetURL()
	if url1 != "https://test_1.github.com" {
		t.Error("client1 url incorrect: ", url1)
	}

	url2 := client2.GetURL()
	if url2 != "http://test_2.github.com" {
		t.Error("client2 url incorrect: ", url2)
	}

	client3.Protocol = "http"
	url3 := client3.GetURL()
	if url3 != "http://test_3.github.com" {
		t.Error("client3 url incorrect: ", url3)
	}

	client3.Protocol = "ftp"
	url4 := client3.GetURL()
	if url4 != "https://test_3.github.com" {
		t.Error("client3 url incorrect: ", url3)
	}

	urlDefault := clientDefault.GetURL()
	if urlDefault != "https://api.github.com" {
		t.Error("Default url incorrect: ", urlDefault)
	}
}
