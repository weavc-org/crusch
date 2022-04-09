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

	m := struct {
		Code     string `url:"code"`
		Secret   string `url:"client_secret"`
		State    string `url:"state"`
		Id       string `url:"client_id"`
		Redirect string `url:"redirect_uri"`
	}{
		Code:     code,
		Secret:   s.githubConfig.Secret,
		State:    state,
		Id:       s.githubConfig.ClientId,
		Redirect: s.githubConfig.RedirectUri,
	}

	body, err := internal.ParseQuery(&m)
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
	q := struct {
		State        string `url:"state"`
		Redirect_uri string `url:"redirect_uri"`
		Client_id    string `url:"client_id"`
	}{
		State:        state,
		Redirect_uri: s.githubConfig.RedirectUri,
		Client_id:    s.githubConfig.ClientId,
	}

	qs, err := internal.ParseQuery(&q)
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

	if s.clientConfig.ValidateState {
		if state != s.clientConfig.StateHandler(r) {
			w.WriteHeader(400)
			w.Write([]byte("Invalid state"))
			return
		}
	}

	accessCode, err := s.AccessCode(code, state)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Could not retrieve access code from github"))
		return
	}

	c := &http.Cookie{
		Name:     "gh_token",
		Value:    accessCode.AccessToken,
		Expires:  time.Now().UTC().Add(time.Duration(accessCode.ExpiresIn * 1000000000)),
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		HttpOnly: true,
	}
	http.SetCookie(w, c)

	if len(accessCode.RefreshToken) > 0 {
		cr := &http.Cookie{
			Name:     "gh_refresh_token",
			Value:    accessCode.RefreshToken,
			Expires:  time.Now().UTC().Add(time.Duration(accessCode.RefreshTokenExpiresIn * 1000000000)),
			Secure:   true,
			SameSite: http.SameSiteDefaultMode,
			HttpOnly: true,
		}

		http.SetCookie(w, cr)
	}

	s.clientConfig.ResponseHandler(w, r)
}
