#!/bin/bash


curl "http://localhost:9000/2015-03-31/functions/function/invocations" \
  --data '{"path": "/hello", "body": "{\"key1\":\"value1\",\"key2\":\"value2\"}"}'