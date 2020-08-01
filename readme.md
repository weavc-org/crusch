![tests](https://github.com/weavc/crusch/workflows/Go/badge.svg?branch=master) 
[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/weavc/crusch)

Crusch is a lightweight libary which provides tools for Github Apps to communicate with Githubs V3 API, without too much unnecessary hassle.

This libary provides a simple client structure to make requests to Githubs API. Clients aid with adding and creating the required authorization headers, keeping track of when they might need to be renewed and other small helper methods.
 
If you are looking for something more complete then [`go-github`](https://github.com/google/go-github) is probably for you. [`go-github`](https://github.com/google/go-github) is a more complete libary with types, seperate methods and bindings for every request github offers, which for what I was working on, was too complicated and quite annoying to work with when all I wanted was something simple, hence crusch.

### Usage

```go
import "github.com/weavc/crusch"
```

basic installation example
```go
var v []map[string]interface{}
authorizer, err := crusch
    .NewInstallationAuth(<ApplicationID int64>, <InstallationID int64>, <rsaKey *rsa.PrivateKey>)

res, err := crusch.Client.Get(
    authorizer, 
    "/repos/weavc/crusch/issues", 
    "assignee=weavc&state=open", 
    &v)
```

new client
```go
httpClient := &http.Client{}

client := crusch.NewGithubClient("api.github.com", "https")
client.SetHTTPClient(httpClient)
```

