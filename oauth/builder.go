package oauth

import "net/http"

type Builder interface {
	Build() *OauthService
	ConfigureGithub(secret string, clientId string, redirectUri string) Builder
	Configure(func(*ClientConfig)) Builder
	RegisterCustomResponseHandler(http.HandlerFunc) Builder
	RegisterCustomStateHandler(func(w http.ResponseWriter, r *http.Request) string) Builder
}

func NewBuilder() Builder {
	return &builder{
		oauthService: &OauthService{
			githubConfig: nil,
			clientConfig: defaultConfig()}}
}

type builder struct {
	Builder
	oauthService *OauthService
}

func (b *builder) Build() *OauthService {
	if b.oauthService.githubConfig == nil {
		panic("Please configure your github credentials by calling ConfigureGithub(privateKey string, clientId string) before building the service.")
	}

	return b.oauthService
}

func (b *builder) ConfigureGithub(secret string, clientId string, redirectUri string) Builder {
	b.oauthService.githubConfig = &GithubConfig{Secret: secret, ClientId: clientId, RedirectUri: redirectUri}
	return b
}

func (b *builder) Configure(configureFunc func(*ClientConfig)) Builder {
	configureFunc(b.oauthService.clientConfig)
	return b
}

func (b *builder) RegisterCustomResponseHandler(f http.HandlerFunc) Builder {
	b.oauthService.clientConfig.ResponseHandler = f
	return b
}

func (b *builder) RegisterCustomStateHandler(f func(w http.ResponseWriter, r *http.Request) string) Builder {
	b.oauthService.clientConfig.StateHandler = f
	return b
}

func defaultConfig() *ClientConfig {
	return &ClientConfig{
		ValidateState:   true,
		ResponseHandler: defaultResponseHandler,
		StateHandler:    defaultStateHandler}
}

func defaultResponseHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Ok"))
}

func defaultStateHandler(w http.ResponseWriter, r *http.Request) string {
	return "12345"
}
