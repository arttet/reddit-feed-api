#!/bin/bash

curl -v -s -X 'POST' \
  -d '{"player": 1}' \
  -H "Content-Type: application/json" \
  http://localhost:8080/v1/posts | jq .
