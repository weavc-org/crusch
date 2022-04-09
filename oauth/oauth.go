package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/weavc/crusch/internal"
)

type OauthService struct {
	clientConfig *ClientConfig
	githubConfig *GithubConfig
}

func (s *OauthService) AccessCode(code string, state string) (AccessToken, error) {
	v := AccessToken{}

	m := map[string]string{
		"code":          code,
		"client_secret": s.githubConfig.Secret,
		"state":         state,
		"client_id":     s.githubConfig.ClientId,
		"redirect_uri":  s.githubConfig.RedirectUri,
	}
	body, err := internal.JsonifyBody(m)
	if err != nil {
		return v, err
	}

	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodPost, "https://github.com/login/oauth/access_token", body)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return v, err
	}

	if res.StatusCode != 200 {
		return v, fmt.Errorf("request failed with status code %b", res.StatusCode)
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&v)
	if err != nil {
		return v, err
	}

	return v, nil
}

func (s *OauthService) Redirect(state string) string {
	q := map[string]string{
		"state":        state,
		"redirect_uri": s.githubConfig.RedirectUri,
		"client_id":    s.githubConfig.ClientId,
	}

	qs, err := internal.ParseQuery(q)
	if err != nil {
		panic("Unable to format querystring")
	}

	return fmt.Sprintf("https://github.com/login/oauth/authorize?%s", qs)
}

func (s *OauthService) OauthRedirectHandler(w http.ResponseWriter, r *http.Request) {
	state := s.clientConfig.StateHandler(w, r)
	redirect := s.Redirect(state)
	http.Redirect(w, r, redirect, http.StatusFound)
}

func (s *OauthService) OauthResponseHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	code := params.Get("code")
	state := params.Get("state")

	accessCode, err := s.AccessCode(code, state)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	c := &http.Cookie{
		Name:     "gh_token",
		Value:    accessCode.AccessToken,
		Expires:  internal.ParseUnix(accessCode.ExpiresIn),
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
	}

	cr := &http.Cookie{
		Name:     "gh_refresh_token",
		Value:    accessCode.RefreshToken,
		Expires:  internal.ParseUnix(accessCode.RefreshTokenExpiresIn),
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
	}

	http.SetCookie(w, c)
	http.SetCookie(w, cr)

	s.clientConfig.ResponseHandler(w, r)
}
