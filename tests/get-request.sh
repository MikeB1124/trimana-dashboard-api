#!/bin/bash


curl "http://localhost:9000/2015-03-31/functions/function/invocations" \
  --data '{"httpMethod": "GET", "path": "/poynt/sales", "queryStringParameters": {"period": "DAILY"}}'