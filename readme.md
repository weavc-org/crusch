![Go](https://github.com/aixr/crusch/workflows/Go/badge.svg?branch=master)

Crusch is a lightweight libary which provides tools for Github Apps communicating with Githubs V3 API, without too much unnecessary hassle.

This libary provides a simple client structure to make requests to their API, the clients aid with adding and creating the required authorization headers and keeping track of when they might need to be renewed etc.
 
If you are looking for something more complete then [`go-github`](https://github.com/google/go-github) is probably for you. [`go-github`](https://github.com/google/go-github) is a more complete libary with types, seperate methods and bindings for every request github offers, which for what I was working on, was too complicated and quite annoying to work with when all I wanted was something simple, hence crusch.

### Usage

```go
import "github.com/aixr/crusch"
```

basic example
```go
client := crusch.NewDefault()
client.NewInstallationAuthFile(<ApplicationID>, <InstallationID>, <PEM keyfile location>)

v := make(map[string]interface{})
respose, err := client.GetJson(
    fmt.Sprintf("/repos/%s/%s/issues", <user>, <repo>), 
    "assignee=aixr&state=open", &v)
```

Clients created through `crusch.New(<name>, <baseURL>)` or `crusch.NewDefault()` are added to a client pool, accessed via `crusch.Pool` variable, allowing you to get premade clients by name or authentication details.
```go
client1 := crusch.New("client_1", "api.github.com")
client1.NewInstallationAuth(<ApplicationID>, <InstallationID>, <private key>)

...

client = crusch.Pool.Get("client_1")

issue := []byte(`
    "title": "Issue title",
    "body": "Issue body"
`)
buffer := bytes.NewBuffer(issue)
client.Post(fmt.Sprintf("/repos/%s/%s/issues", <owner>, <repo>), buffer)
```

If more control over the request is required, there is also a `Do(*http.Request)` method which will take a `*http.Request`, apply the authentication headers and send the request.
```go
req := http.Request{}
req.Header = http.Header{}

req.Method = http.MethodPost
req.Header.Add("Accept", "application/vnd.github.squirrel-girl-preview+json")
req.URL = &url.URL{
    Scheme: "https", 
    Host: "api.github.com", 
    Path: "/repos/<owner>/<repo>/comments/<comment_id>/reactions",
}

client = crusch.NewDefault()
client.NewOAuth("11a6c2809da4bc163487481cf02fb210", "bearer")
res, err := client.Do(req)
```

### Notes
Im extemely new to golang (this is the second project, which was done alongside the 'first') so there are certainly things that can be done better, I will happily take feedback/contributions.

### Todo
1. use contexts
1. finish tests
1. docs
1. optional types/structs? maybe - go-github's types work well alongside this already 
