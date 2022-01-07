package modo2auth_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/modopayments-ux/modo2auth/v2"
)

func mockBody(t testing.TB) []byte {
	t.Helper()

	body := map[string]string{
		"someKey": "someValue",
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		t.Error(err)
	}
	return jsonData
}

func mockID() modo2auth.ID {
	return modo2auth.ID{
		Key:    "7g0UApipMpuJ1VOOOHWJNIZH7VZINb08",
		Secret: "20I1s7GH7-pgn9041cgWlBKU8pcA1I4CCNpGuvu_xL4K-GnRSy3Q6IBtA5LYlIjy",
	}
}

func TestSign(t *testing.T) {
	t.Parallel()

	apiURI := "/test"
	body := mockBody(t)
	modo := mockID()
	expected := "MODO2 eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTAwNzI2ODUsImFwaV9pZGVudGlmaWVyIjoiN2cwVUFwaXBNcHVKMVZPT09IV0pOSVpIN1ZaSU5iMDgiLCJhcGlfdXJpIjoiL3Rlc3QiLCJib2R5X2hhc2giOiI0NmI4ZGFkMWM2OWNiZDUwMGU5ZDFmZTJmMmVjZTM1N2M4NGM2ZTM2Y2U3YTg2MGJmMTQ2NzJiNGI3NDBhZjE5In0.sWjjz_MpnSv8Z31gNpEm1cmDhN7MK7Z3ix61RbDRL7g"

	signature, err := modo2auth.Sign(apiURI, time.Unix(1590072685, 0), body, []byte(modo.Key), []byte(modo.Secret))
	if err != nil {
		t.Error(err)
	}

	if signature != expected {
		t.Errorf("Expected token: "+expected+" \nBut got: %v", signature)
	}
}

func TestID_Sign(t *testing.T) {
	t.Parallel()

	apiURI := "/test"
	body := mockBody(t)
	req, _ := http.NewRequest(http.MethodPost, apiURI, bytes.NewBuffer(body))
	modo := mockID()

	err := modo.Sign(req)
	if err != nil {
		t.Error(err)
	}

	if req.Header.Get("Authorization") == "" {
		t.Errorf("Expected an Authorization header")
	}

	rbody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Error(err)
	}
	if string(rbody) != string(body) {
		t.Errorf("request body changed.\nHave: %s\nWant: %s", string(rbody), string(body))
	}
}
