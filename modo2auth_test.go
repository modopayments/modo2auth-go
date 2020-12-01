package modo2auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

// mocks
func mockBody() []byte {
	body := map[string]string{
		"someKey": "someValue",
	}
	jsonData, _ := json.Marshal(body)
	return jsonData
}

func mockModo() Config {
	return Config{
		APIIdentifier: "7g0UApipMpuJ1VOOOHWJNIZH7VZINb08",
		APISecret:     "20I1s7GH7-pgn9041cgWlBKU8pcA1I4CCNpGuvu_xL4K-GnRSy3Q6IBtA5LYlIjy",
		Debug:         true,
	}
}

// TODO: Resolve how/if/when to run simple integration test against real server
// func TestIntegrationGet(test *testing.T) {
// 	modo := mockModo()
// 	modo.Debug = false
// 	apiHost := "http://localhost:82"
// 	apiURI := "/v2/vault/public_key"
// 	req, _ := http.NewRequest("GET", apiHost+apiURI, nil) // Test path, Test body absence is okay?
// 	signedReq, _ := modo.Sign(req)
// 	resp, _ := http.DefaultClient.Do(signedReq)

// 	if resp.StatusCode != 200 {
// 		test.Errorf("Expected 200 but got: %v", resp.StatusCode)
// 	}
// }
// func TestIntegrationPost(test *testing.T) {
// 	// format body data
// 	data := map[string]string{
// 		"start_date": "2020-05-01T00:00:00Z",
// 		"end_date":   "2020-05-26T00:00:00Z",
// 	}
// 	jsonData, _ := json.Marshal(data)
// 	body := bytes.NewBuffer(jsonData)

// 	modo := mockModo()
// 	modo.Debug = false
// 	apiHost := "http://localhost:82"
// 	apiURI := "/v2/reports"

// 	req, _ := http.NewRequest("POST", apiHost+apiURI, body)
// 	signedReq, _ := modo.Sign(req)
// 	resp, err := http.DefaultClient.Do(signedReq)
// 	if err != nil {
// 		test.Error("err", err)
// 	}

// 	if resp.StatusCode != 200 {
// 		test.Errorf("Expected 200 but got: %v", resp.StatusCode)
// 	}
// }

func TestSign(test *testing.T) {
	apiURI := "/test"
	body := mockBody()
	req, _ := http.NewRequest("POST", apiURI, bytes.NewBuffer(body))
	modo := mockModo()

	signedReq, err := modo.Sign(req)
	if err != nil {
		test.Error(err)
	}

	authHeader := signedReq.Header.Get("Authorization")
	expected := "MODO2 eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTAwNzI2ODUsImFwaV9pZGVudGlmaWVyIjoiN2cwVUFwaXBNcHVKMVZPT09IV0pOSVpIN1ZaSU5iMDgiLCJhcGlfdXJpIjoiL3Rlc3QiLCJib2R5X2hhc2giOiI0NmI4ZGFkMWM2OWNiZDUwMGU5ZDFmZTJmMmVjZTM1N2M4NGM2ZTM2Y2U3YTg2MGJmMTQ2NzJiNGI3NDBhZjE5In0.sWjjz_MpnSv8Z31gNpEm1cmDhN7MK7Z3ix61RbDRL7g"

	if authHeader != expected {
		test.Errorf("Expected token: "+expected+" \nBut got: %v", authHeader)
	}

}

func TestGetToken(test *testing.T) {
	modo := mockModo()
	body := mockBody()
	expected :=
		"MODO2 eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTAwNzI2ODUsImFwaV9pZGVudGlmaWVyIjoiN2cwVUFwaXBNcHVKMVZPT09IV0pOSVpIN1ZaSU5iMDgiLCJhcGlfdXJpIjoiL3Rlc3QiLCJib2R5X2hhc2giOiI0NmI4ZGFkMWM2OWNiZDUwMGU5ZDFmZTJmMmVjZTM1N2M4NGM2ZTM2Y2U3YTg2MGJmMTQ2NzJiNGI3NDBhZjE5In0.sWjjz_MpnSv8Z31gNpEm1cmDhN7MK7Z3ix61RbDRL7g"

	token, err := GetToken("/test", body, modo)

	if err != nil {
		test.Error(err)
	}
	if token != expected {
		test.Errorf("Expected token: "+expected+" \nBut got: %v", token)
	}
}

func TestMakeHeader(test *testing.T) {
	header := makeHeader()
	expected := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	data := string(header)

	if data != expected {
		test.Errorf("Expected: "+expected+"\nGot:%v", data)
	}
}

func TestBodyHash(test *testing.T) {
	jsonData := mockBody()
	hashedBody := bodyHash(jsonData)
	data := string(hashedBody)
	expected := "46b8dad1c69cbd500e9d1fe2f2ece357c84c6e36ce7a860bf14672b4b740af19"

	if data != expected {
		test.Errorf("Expected: "+expected+"\nGot:%v", data)
	}
}

func TestMakePayload(test *testing.T) {
	body := mockBody()
	payload, err := makePayload("/payouts/minimum?action=complete", "7g0UApipMpuJ1VOOOHWJNIZH7VZINb08", body, true)
	data := string(payload)
	expected := "eyJpYXQiOjE1OTAwNzI2ODUsImFwaV9pZGVudGlmaWVyIjoiN2cwVUFwaXBNcHVKMVZPT09IV0pOSVpIN1ZaSU5iMDgiLCJhcGlfdXJpIjoiL3BheW91dHMvbWluaW11bT9hY3Rpb249Y29tcGxldGUiLCJib2R5X2hhc2giOiI0NmI4ZGFkMWM2OWNiZDUwMGU5ZDFmZTJmMmVjZTM1N2M4NGM2ZTM2Y2U3YTg2MGJmMTQ2NzJiNGI3NDBhZjE5In0"

	if err != nil {
		test.Error(err)
	}
	if data != expected {
		test.Errorf("Expected: "+expected+"\nGot:%v", data)
	}
}

func TestMakeSignature(test *testing.T) {
	header := makeHeader()
	payload := []byte("eyJpYXQiOjE1OTAwNzI2ODUsImFwaV9pZGVudGlmaWVyIjoiN2cwVUFwaXBNcHVKMVZPT09IV0pOSVpIN1ZaSU5iMDgiLCJhcGlfdXJpIjoiL3BheW91dHMvbWluaW11bT9hY3Rpb249Y29tcGxldGUiLCJib2R5X2hhc2giOiI0NmI4ZGFkMWM2OWNiZDUwMGU5ZDFmZTJmMmVjZTM1N2M4NGM2ZTM2Y2U3YTg2MGJmMTQ2NzJiNGI3NDBhZjE5In0")
	signature := makeSignature(header, payload, "20I1s7GH7-pgn9041cgWlBKU8pcA1I4CCNpGuvu_xL4K-GnRSy3Q6IBtA5LYlIjy")
	data := string(signature)
	expected := "lE865md6iVe42QyAGMpcm4bJntMACcDISfCxMrKzOuo"

	if data != expected {
		test.Errorf("Expected: "+expected+"\nGot:%v", data)
	}
}
