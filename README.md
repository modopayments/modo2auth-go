# modo2auth-go

> A Go module to generate authentication details to communicate with Modo servers

# Prerequesites

**Credentials** that are created and shared by Modo. These will be different for each environment (`int`, `prod`, `local` etc...).

- `APIIdentifier` - API key from Modo
- `APISecret` - API secret from Modo

These values will be used when intantiating the library.

# Install
Presently, the github repo is private, so you'll need to tell `go` that:
```bash
# add env var via `~/.zshrc` or `~/.bashrc` or the like
export GOPRIVATE="github.com/modopayments-ux,github.com/modopayments"
```
Restart your terminal


```bash
# in the root of the project
go get github.com/modopayments-ux/modo2auth-go
```

# Example Usage

Here's an example using `...` to make requests. You can use your preferred method or library.

## `GET` Example

```go
package main

import (
  "fmt" // for example
	"io/ioutil" // for example
	"net/http"

  // 1. IMPORT
	"github.com/modopayments-ux/modo2auth-go"
)

func main() {
  // 2. INSTANTIATE - get these from MODO
	modo := modo2auth.ID{
		APIIdentifier: "...",
		APISecret:     "...",
  }
  
  // 3. SIGN & SEND REQUEST
	apiHOST := "http://localhost:82"
  apiURI := "/v2/vault/public_key"  
	req, _ := http.NewRequest("GET", apiHOST+apiURI, nil)
	signedReq, _ := modo.Sign(req)
	resp, _ := http.DefaultClient.Do(signedReq)
	readBody, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("Response body is", string(readBody))
}

```

## `POST` Example

```go
package main

import (
  "fmt" // for example
	"io/ioutil" // for example
	"net/http"

  // 1. IMPORT
	"github.com/modopayments-ux/modo2auth-go"
)

func main() {
  // 2. INSTANTIATE - get these from MODO
  modo := modo2auth.ID{
		APIIdentifier: "...",
		APISecret:     "...",
  }

  // 3. SIGN & SEND REQUEST
  apiHOST := "http://localhost:82"
  apiURI := "/v2/reports"  
  
  // format body data
	data := map[string]string{
		"start_date": "2020-05-01T00:00:00Z",
		"end_date":   "2020-05-26T00:00:00Z",
	}
	jsonData, _ := json.Marshal(data)
  body := bytes.NewBuffer(jsonData)
  
  // request
	req, _ := http.NewRequest("POST", apiHOST+apiURI, body)
	signedReq, _ := modo.Sign(req)
	resp, err := http.DefaultClient.Do(signedReq)
  
  // response
  respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response body is",respBody)
}

```

# API

## `Config({})`

Returns a Modo2Api Config struct.

- `APIIdentifier` (string) - API key from Modo
- `APISecret` (string) - API secret from Modo
- `Debug` (bool) - Enables static values for testing
- 
## `Sign(http.Request)`

Returns an `*http.Request` with a Modo2Auth token added as an Authorization header.

`http.Request` (object) - HTTP request being made

## `GetToken(apiURI, body, modo)`

Returns a Modo2Auth token (string) to be added to an HTTP request as an Authorization header (`{Authorization: value}`).

- `apiURI` (`[]byte`) - Api Uri intending to call to (ex: `"/v2/vault/public_key"`)
- `body` (`[]byte`) - Body of the request
  - should be stringified JSON
  - for `GET` requests, leave as `nil`
- `modo` (`modo2auth.ID`) - Credentials for requests

# Development

Prerequisite: `go` installed globally

1. Install ...
2. ...
3. Unit test - `...`

# Contributing
1. Fork this repo via Github
2. Create your feature branch (`git checkout -b feature/my-new-feature`)
3. Ensure unit tests are passing (`go test`)
4. Commit your changes (`git commit -am 'Add some feature'`)
5. Push to the branch (`git push origin feature/my-new-feature`)
6. Create new Pull Request via Github

# Publishing
Prequisite: Need to have `go` installed on your system. At the root of this directory, do the following:

1. Tag with new version `git tag v1.1.0` (According to Semantec Versioning guidelines)
2. Push tags `git push --tags`
3. ...
