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

// JSONData ...
type JSONData []byte

// UnmarshalText ...
func (d *JSONData) UnmarshalText(data []byte) error {
	*d = data
	return nil
}

// MarshalText ...
func (d JSONData) MarshalText() ([]byte, error) {
	return d, nil
}

type tokenHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}
type tokenPayload struct {
	Iat           int64    `json:"iat"`
	APIIdentifier JSONData `json:"api_identifier"`
	APIUri        JSONData `json:"api_uri"`
	BodyHash      JSONData `json:"body_hash"`
}

// ID holds API credentials for communicating with Modo servers, as well as an optional `debug` flag for testing
type ID struct {
	APIIdentifier string
	APISecret     string
	Debug         bool // used only for testing purposes
}

// Sign receives an http.Request and signs an Authorization Header
func (modo ID) Sign(req *http.Request) (*http.Request, error) {
	apiURI := []byte(req.URL.Path)

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

	token, err := getToken(apiURI, body, modo)
	if err != nil {
		return req, err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// getToken ...
func getToken(apiURI []byte, body []byte, modo ID) (string, error) {
	// get credentials
	apiIdentifier := []byte(modo.APIIdentifier)
	apiSecret := []byte(modo.APISecret)

	// get components
	header := _makeHeader()
	payload, err := _makePayload(apiURI, apiIdentifier, body, modo.Debug)
	if err != nil {
		return "", err
	}
	signature := _makeSignature(header, payload, apiSecret)

	// concat final string
	token := "MODO2 " + string(header) + "." + string(payload) + "." + string(signature)

	return token, nil
}

func _makeHeader() []byte {
	data := tokenHeader{Alg: "HS256", Typ: "JWT"}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	return _base64URLEncode(jsonData)
}

func _makePayload(apiURI []byte, apiIdentifier []byte, body []byte, debug bool) ([]byte, error) {
	iat := time.Now().Unix() // in seconds
	if debug {
		// static time for testing
		iat = int64(1590072685)
	}

	bodyHash := _bodyHash([]byte(body)) // hex digest of sha256 of data

	payload := tokenPayload{
		Iat:           iat,
		APIIdentifier: JSONData(apiIdentifier),
		APIUri:        JSONData(apiURI),
		BodyHash:      JSONData(bodyHash),
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return _base64URLEncode(jsonData), nil
}

func _makeSignature(header []byte, payload []byte, secret []byte) []byte {
	hash := hmac.New(sha256.New, secret)
	hash.Write(header)
	hash.Write([]byte{'.'})
	hash.Write(payload)
	signature := hash.Sum(nil)

	return _base64URLEncode(signature)
}

func _bodyHash(body []byte) []byte {
	hasher := sha256.New()
	hasher.Write(body)
	hash := hasher.Sum(nil)

	hashHex := make([]byte, hex.EncodedLen(len(hash)))
	hex.Encode(hashHex, hash)
	return hashHex
}

func _base64URLEncode(data []byte) []byte {
	enc := base64.URLEncoding.WithPadding(base64.NoPadding)
	encoded := make([]byte, enc.EncodedLen(len(data)))
	enc.Encode(encoded, data)
	return encoded
}
