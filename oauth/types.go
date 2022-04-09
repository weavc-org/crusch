package oauth

import "net/http"

type GithubConfig struct {
	Secret      string
	ClientId    string
	RedirectUri string
}

type ClientConfig struct {
	ValidateState   bool
	ResponseHandler http.HandlerFunc
	StateHandler    func(r *http.Request) string
}

type AccessToken struct {
	AccessToken           string `json:"access_token"`
	ExpiresIn             int64  `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in"`
	Scope                 string `json:"scope"`
	TokenType             string `json:"token_type"`
	Error                 string `json:"error"`
	ErrorDescription      string `json:"error_description"`
}

type AccessTokenRequest struct {
	Code     string `url:"code"`
	Secret   string `url:"client_secret"`
	State    string `url:"state"`
	Id       string `url:"client_id"`
	Redirect string `url:"redirect_uri"`
}

type RedirectRequest struct {
	State        string `url:"state"`
	Redirect_uri string `url:"redirect_uri"`
	Client_id    string `url:"client_id"`
}
