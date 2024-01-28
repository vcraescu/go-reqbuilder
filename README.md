# go-reqbuilder

![Test Status](https://github.com/vcraescu/go-reqbuilder/actions/workflows/go.yml/badge.svg)
[![codecov](https://codecov.io/github/vcraescu/go-reqbuilder/branch/master/graph/badge.svg)](https://codecov.io/github/vcraescu/go-reqbuilder)

The package provides a flexible and customizable way to build HTTP requests in Go, allowing users to easily construct
requests with various methods, headers, bodies, and parameters.


## Installation

To use reqbuilder in your Go project, you can install it using go get:

```
go get github.com/vcraescu/go-reqbuilder
```

## Example

```go
type User struct {
    ID        string `json:"id,omitempty`
    FirstName string `json:"firstName,omitempty"`
    LastName  string `json:"lastName,omitempty"`
}

func main() {
    req, err := reqbuilder.NewBuilder("https://api.example.com").
        WithMethod(http.MethodPost).
        WithPath("/users").
        WithBody(User{
            ID:        "1",
            FirstName: "John",
            LastName:  "Doe",
        }).
        WithHeaders(reqbuilder.JSONAcceptHeader, reqbuilder.JSONContentHeader).
        Build(context.Background())
    if err != nil {
        log.Fatalf(err)
    }

    // ...
}
```

## License

This package is distributed under the MIT License. See the [LICENSE](LICENSE) file for details.

