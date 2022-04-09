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
	StateHandler    func(w http.ResponseWriter, r *http.Request) string
}

type AccessToken struct {
	AccessToken           string
	ExpiresIn             int64
	RefreshToken          string
	RefreshTokenExpiresIn int64
	Scope                 string
	TokenType             string
}
