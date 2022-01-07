package modo2auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type JSONData []byte

func (d *JSONData) UnmarshalText(data []byte) error {
	*d = data
	return nil
}

func (d JSONData) MarshalText() ([]byte, error) {
	return d, nil
}

type authenticationHeader struct {
	Iat           int64    `json:"iat"`
	ApiIdentifier JSONData `json:"api_identifier"`
	ApiUri        string   `json:"api_uri"`
	BodyHash      JSONData `json:"body_hash"`
}

type tokenHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// ID holds API credentials for communicating with Modo servers, as well as an optional `debug` flag for testing
type ID struct {
	Key    string
	Secret string
}

// Sign adds an Authorization header to an http.request
func (id ID) Sign(req *http.Request) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	signature, err := Sign(req.URL.Path, time.Now(), body, []byte(id.Key), []byte(id.Secret))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", signature)
	req.Body = ioutil.NopCloser(bytes.NewReader(body))

	return nil
}

// Sign returns the signature to use for a request of [data] to the endpoint [api]
func Sign(api string, iat time.Time, data []byte, key, secret []byte) (signature string, err error) {
	token, authData := generateToken(api, data)
	authData.ApiIdentifier = key
	authData.Iat = iat.Unix()
	var authBytes []byte
	if authBytes, err = json.Marshal(authData); err != nil {
		return "", err
	}
	authData64 := base64URLEncode(authBytes)
	sig := generateSignature(token[len("MODO2 "):], authData64, secret)
	return string(token) + "." + string(authData64) + "." + string(sig), nil
}

func generateSignature(token, authData []byte, secret []byte) []byte {
	hash := hmac.New(sha256.New, secret)
	hash.Write(token)
	hash.Write([]byte{'.'})
	hash.Write(authData)

	signature := hash.Sum(nil)
	base64Signature := base64URLEncode(signature)

	return base64Signature
}

func generateToken(endpoint string, body []byte) (token []byte, auth authenticationHeader) {
	jsonTokenHeader, _ := json.Marshal(tokenHeader{Alg: "HS256", Typ: "JWT"})
	tokenData := base64URLEncode(jsonTokenHeader)
	token = make([]byte, len(tokenData)+len("MODO2 "))
	copy(token[copy(token, "MODO2 "):], tokenData)
	auth = authenticationHeader{
		BodyHash: bodyHash(body),
		ApiUri:   endpoint,
	}
	return
}

func bodyHash(body []byte) []byte {
	hasher := sha256.New()
	hasher.Write(body)
	hash := hasher.Sum(nil)

	hashHex := make([]byte, hex.EncodedLen(len(hash)))
	hex.Encode(hashHex, hash)
	return hashHex
}

func base64URLEncode(data []byte) []byte {
	enc := base64.URLEncoding.WithPadding(base64.NoPadding)
	encoded := make([]byte, enc.EncodedLen(len(data)))
	enc.Encode(encoded, data)
	return encoded
}
