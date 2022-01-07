# modo2auth-go

> A Go module to generate authentication details to communicate with Modo servers

# Prerequesites

**Credentials** that are created and shared by Modo. These will be different for each environment (`int`, `prod`, `local` etc...).

- `api_identifier` - API key from Modo
- `api_secret` - API secret from Modo

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

## Example

```
package main

import (
	"net/http"

	"github.com/modopayments-ux/modo2auth_go"
)

const endpoint = "http://modopayments.com"

func main() {
    id := modo2auth_go.ID{
        Key:    "7g0UApipMpuJ1VOOOHWJNIZH7VZINb08",
        Secret: "20I1s7GH7-pgn9041cgWlBKU8pcA1I4CCNpGuvu_xL4K-GnRSy3Q6IBtA5LYlIjy",
    }
    req, _ := http.NewRequest(http.MethodPost, endpoint+"/v3/checkout/list", nil)
    _ = id.Sign(req)
}

```

# API

`ID` (struct) with the following properties:

- `Key` (string) - API key from Modo
- `Secret` (string) - API secret from Modo

## `ID.Sign(r *http.Request)`

Adds an "Authorization" header to the provided request, signed using the credentials in ID. The request must not be modified after being signed

## `Sign(api string, iat time.Time, data []byte, key, secret []byte) (signature string, err error)`

Returns a signature for use in an Authorization header. iat should normally be the current time. api must be the API *path* (ex. "/v3/checkout/list"), NOT the full URL.

# Development

Prerequisite: `go` installed globally

1. Install `go install .`
3. Unit test - `go test ./...`

# Contributing
1. Fork this repo via Github
2. Create your feature branch (`git checkout -b feature/my-new-feature`)
3. Ensure unit tests are passing (`go test ./...`)
4. Commit your changes (`git commit -am 'Add some feature'`)
5. Push to the branch (`git push origin feature/my-new-feature`)
6. Create new Pull Request via Github

# Publishing
Prequisite: Need to have `go` installed on your system. At the root of this directory, do the following:

1. Commit and push (`git commit -am 'Version bump'`)
2. Tag with new version `git tag v1.1.0` (example)
3. Push tags `git push --tags`
