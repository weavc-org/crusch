## Crusch

A lightweight golang module for making requests and authenticating against Githubs V3 json API. If you are looking for something more complete, [`google/go-github`](https://github.com/google/go-github) is probably the module for you. I made this as I didn't really have a need for what `go-github` offered, I just wanted something quick to make requests and authenticate with Githubs API without the hassle of models, cross-referencing multiple pieces of documentation etc. 

Crusch pretty much just handles authentication, headers and can do model binding etc if wanted.

### Usage

```go
import "github.com/aixr/crusch"
```

basic example
```go
client := crusch.NewDefault()
client.NewInstallationAuthFile(<ApplicationID>, <InstallationID>, <PEM keyfile location>)

model, respose, err := client.GET(
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
- [ ] optional types/models
- [ ] other forms of authentication i.e. oauth
- [ ] docs
- [ ] tests

### Notes
Im extemely new to golang (this is the second project, made alongside the 'first') so there are certainly things that can be done better, I will happily take feedback/contributions.
