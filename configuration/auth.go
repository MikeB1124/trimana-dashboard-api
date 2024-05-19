package configuration

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type PoyntJWTResponse struct {
	ExpiresIn    int    `json:"expiresIn"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Scope        string `json:"scope"`
	TokenType    string `json:"tokenType"`
}

func GetPoyntJWTAccessToken() (string, error) {
	config := GetConfig()

	block, _ := pem.Decode([]byte(config.Poynt.PrivateKey))
	if block == nil {
		log.Fatalf("failed to decode PEM block containing private key")
		return "", fmt.Errorf("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("failed to parse private key: %v", err)
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	// Prepare JWT payload
	jwtPayload := jwt.MapClaims{
		"exp": time.Now().Add(time.Second * 300).Unix(),
		"iat": time.Now().Unix(),
		"iss": config.Poynt.ApplicationID,
		"sub": config.Poynt.ApplicationID,
		"aud": config.Poynt.URL,
		"jti": uuid.New().String(),
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtPayload)

	// Sign the token with the private key
	jwtToken, err := token.SignedString(privateKey)
	if err != nil {
		log.Fatalf("failed to sign the token: %v", err)
		return "", fmt.Errorf("failed to sign the token: %v", err)
	}

	// Prepare the form payload
	payload := url.Values{}
	payload.Set("grantType", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	payload.Set("assertion", jwtToken)

	// Prepare the HTTP request
	tokenEndpoint := config.Poynt.URL + "/token"
	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(payload.Encode()))
	if err != nil {
		log.Fatalf("failed to create request: %v", err)
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("api-version", "1.2")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed to send request: %v", err)
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read and process the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read response body: %v", err)
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var poyntJWTResponse PoyntJWTResponse
	err = json.Unmarshal(body, &poyntJWTResponse)
	if err != nil {
		log.Fatalf("failed to unmarshal poynt jwt response: %v", err)
		return "", fmt.Errorf("failed to unmarshal poynt jwt response: %v", err)
	}

	return poyntJWTResponse.AccessToken, nil
}
