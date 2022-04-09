package oauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/weavc/crusch/internal"
)

type OauthService struct {
	clientConfig *ClientConfig
	githubConfig *GithubConfig
}

func (s *OauthService) AccessCode(code string, state string) (*AccessToken, error) {
	v := &AccessToken{}

	m := &AccessTokenRequest{
		Code:     code,
		Secret:   s.githubConfig.Secret,
		State:    state,
		Id:       s.githubConfig.ClientId,
		Redirect: s.githubConfig.RedirectUri,
	}

	body, err := internal.ParseQuery(m)
	if err != nil {
		return nil, err
	}

	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodPost, "https://github.com/login/oauth/access_token", bytes.NewBuffer([]byte(body)))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("request failed with status code %d", res.StatusCode)
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(v)
	if err != nil {
		return nil, err
	}

	if len(v.Error) != 0 {
		return nil, fmt.Errorf("an error occured retrieving access tokens from github [%s]: %s", v.Error, v.ErrorDescription)
	}

	return v, nil
}

func (s *OauthService) Redirect(state string) string {
	q := &RedirectRequest{
		State:        state,
		Redirect_uri: s.githubConfig.RedirectUri,
		Client_id:    s.githubConfig.ClientId,
	}

	qs, err := internal.ParseQuery(q)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("https://github.com/login/oauth/authorize?%s", qs)
}

func (s *OauthService) OauthRedirectHandler(w http.ResponseWriter, r *http.Request) {
	state := s.clientConfig.StateHandler(r)
	redirect := s.Redirect(state)
	http.Redirect(w, r, redirect, http.StatusFound)
}

func (s *OauthService) OauthResponseHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	code := params.Get("code")
	state := params.Get("state")
	requestError := params.Get("error")

	if len(requestError) > 0 {
		panic(fmt.Errorf("an error occured while redirecting to github [%s]: %s", requestError, params.Get("error_description")))
	}

	if s.clientConfig.ValidateState {
		if state != s.clientConfig.StateHandler(r) {
			w.WriteHeader(400)
			w.Write([]byte("Invalid state"))
			return
		}
	}

	ac, err := s.AccessCode(code, state)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Could not retrieve access code from github"))
		return
	}

	if len(ac.AccessToken) > 0 {
		c := createCookie("gh_token", ac.AccessToken, ac.ExpiresIn)
		http.SetCookie(w, c)
	}

	if len(ac.RefreshToken) > 0 {
		cr := createCookie("gh_refresh_token", ac.RefreshToken, ac.RefreshTokenExpiresIn)
		http.SetCookie(w, cr)
	}

	s.clientConfig.ResponseHandler(w, r)
}

func createCookie(name string, value string, expires int64) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  time.Now().UTC().Add(time.Duration(expires * 1000000000)),
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		HttpOnly: true,
	}
}
