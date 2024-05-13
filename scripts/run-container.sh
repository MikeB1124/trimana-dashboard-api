#!/bin/bash

docker stop trimana-api-lambda
docker rm trimana-api-lambda
GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap main.go
docker build -t trimana-api-image .
docker run --name trimana-api-lambda -p 9000:8080 trimana-api-image