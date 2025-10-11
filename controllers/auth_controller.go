package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/MikeB1124/trimana-dashboard-api/payroll"
	"github.com/aws/aws-lambda-go/events"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func Login(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Login request received")

	// Unmarshal the request
	var request LoginRequest
	err := json.Unmarshal([]byte(event.Body), &request)
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 400, Result: "Invalid request body"})
	}

	// Authenticate user
	user, err := payroll.AuthenticateUser(request.Username, request.Password)
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		return createResponse(payroll.Response{StatusCode: 401, Result: "Invalid credentials"})
	}

	// Generate a simple token (in production, use JWT)
	token, err := generateToken()
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 500, Result: "Error generating token"})
	}

	response := LoginResponse{
		Username: user.Username,
		Token:    token,
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 500, Result: "Error creating response"})
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
