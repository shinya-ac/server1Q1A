package auth0

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type JWKS struct {
	Keys []JSONWebKeys `json:"keys"`
}

func FetchJWKS(auth0Domain string) (*JWKS, error) {
	url := fmt.Sprintf("https://%s/.well-known/jwks.json", auth0Domain)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	jwks := &JWKS{}
	err = json.NewDecoder(resp.Body).Decode(jwks)
	if err != nil {
		return nil, err
	}

	return jwks, err
}
