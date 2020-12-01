package modo2auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type tokenHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}
type tokenPayload struct {
	Iat           int64  `json:"iat"`
	APIIdentifier string `json:"api_identifier"`
	APIUri        string `json:"api_uri"`
	BodyHash      string `json:"body_hash"`
}

// Config holds API credentials for communicating with Modo servers, as well as an optional `debug` flag for testing
type Config struct {
	APIIdentifier string
	APISecret     string
	Debug         bool // enable static values for testing
}

// Sign receives an http.Request and signs an Authorization Header
func (modo Config) Sign(req *http.Request) (*http.Request, error) {
	// get body depending on method of request
	var body []byte
	var err error
	if req.Method == "GET" {
		body = []byte("")
	} else {
		// get the body data
		body, err = ioutil.ReadAll(req.Body)
		// restore the io.ReadCloser to its original state
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		if err != nil {
			fmt.Println("err", err)
			return req, err
		}
	}

	token, err := GetToken(req.URL.Path, body, modo)
	if err != nil {
		return req, err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// GetToken returns an Authorization header string
func GetToken(apiURI string, body []byte, modo Config) (string, error) {
	// get components
	header := makeHeader()
	payload, err := makePayload(apiURI, modo.APIIdentifier, body, modo.Debug)
	if err != nil {
		return "", err
	}
	signature := makeSignature(header, payload, modo.APISecret)

	// concat final string
	token := "MODO2 " + string(header) + "." + string(payload) + "." + string(signature)

	return token, nil
}

func makeHeader() []byte {
	data := tokenHeader{Alg: "HS256", Typ: "JWT"}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	return base64URLEncode(jsonData)
}

func makePayload(apiURI string, apiIdentifier string, body []byte, debug bool) ([]byte, error) {
	iat := time.Now().Unix() // in seconds
	if debug {
		// static time for testing
		iat = int64(1590072685)
	}
	payload := tokenPayload{
		Iat:           iat,
		APIIdentifier: apiIdentifier,
		APIUri:        apiURI,
		BodyHash:      bodyHash(body),
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return base64URLEncode(jsonData), nil
}

func makeSignature(header []byte, payload []byte, secret string) []byte {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write(header)
	hash.Write([]byte{'.'})
	hash.Write(payload)
	signature := hash.Sum(nil)

	return base64URLEncode(signature)
}

func bodyHash(body []byte) string {
	hasher := sha256.New()
	hasher.Write(body)
	hash := hasher.Sum(nil)

	hashHex := make([]byte, hex.EncodedLen(len(hash)))
	hex.Encode(hashHex, hash)
	return string(hashHex)
}

func base64URLEncode(data []byte) []byte {
	enc := base64.URLEncoding.WithPadding(base64.NoPadding)
	encoded := make([]byte, enc.EncodedLen(len(data)))
	enc.Encode(encoded, data)
	return encoded
}
