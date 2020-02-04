## Crusch

Crusch provides tools for communicating Githubs V3 API with Github applications and installations, hopefully without too much unnecessary hassle.

If you are looking for something more complete, [`google/go-github`](https://github.com/google/go-github) is probably for you. But if you are looking for something quick and lightweight this might be worth a shot. 

You just create a client with the necassary details i.e. ApplicationID, InstallationID and a keyfile and you can make requests to githubs API. When using the `Client.`[`GET`, `POST`, `PUT`, `PATCH`, `DELETE`] methods, the basic headers github expects (Accept, Authorization etc) are all attached to the requests by default. Other requests might need something abit different (i.e. reactions), here you can use the `Client.DO` method which takes a pre-built `*http.Request` and attaches only the required authorization headers. 

If provided with a struct, the Client requests will also attempt to bind the responses body to it.  

### Usage

```go
import "github.com/aixr/crusch"
```

basic example
```go
client := crusch.NewDefault()
client.NewInstallationAuthFile(<ApplicationID>, <InstallationID>, <PEM keyfile location>)

respose, err := client.GET(
    fmt.Sprintf("/repos/%s/%s/issues", <user>, <repo>), 
    <query model>, 
    <binding model>)
```

Clients created through `crusch.New(<name>, <baseURL>)` or `crusch.NewDefault()` are added to a client pool, accessed via `crusch.Pool` variable, allowing you to get premade clients by name or authentication details.
```go
client1 := crusch.New("client_1", "api.github.com")
client1.NewInstallationAuth(<ApplicationID>, <InstallationID>, <private key>)

...

client = crusch.Pool.Get("client_1")
client.GET(fmt.Sprintf("/repos/%s/%s/issues", <owner>, <repo>), <query model>, <binding model>)
```

If more control over the request is required, alongside the `GET`, `POST`, `PUT`, `PATCH` and `DELETE` methods there is also a `DO` method which will take a `*http.Request` and apply the authentication headers and make the request.
```go
req := http.Request{}

req.Header.Add("Accept", "application/vnd.github.squirrel-girl-preview+json")
req.URL = &url.URL{
    Scheme: "https", 
    Host: "api.github.com", 
    Path: "/repos/<owner>/<repo>/comments/<comment_id>/reactions",
}
req.Method = http.MethodPost

client = crusch.Pool.Get("client_1")
res, err := client.DO(req)
```

### Todo
- [ ] ~~models/types~~ using go-githubs for now
- [ ] use context
- [ ] docs
- [ ] tests
- [ ] github actions
- [ ] check response success before parsing body
- [ ] check input information on new auths
- [ ] Oauth/token type, github sends down the type in the response
- [ ] remove or allow setting of https/http
- [ ] rename module to github.com/aixr/crusch
- [ ] Redo readme after large changes

### Notes
Im extemely new to golang (this is the second project, made alongside the 'first') so there are certainly things that can be done better, I will happily take feedback/contributions.
