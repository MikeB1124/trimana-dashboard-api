#!/bin/bash


curl "http://localhost:9000/2015-03-31/functions/function/invocations" \
  --data '{"httpMethod": "POST", "path": "/payroll/report", "body": "{\"cardID\":\"0002164531\"}"}'