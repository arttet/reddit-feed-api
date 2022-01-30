#!/bin/bash

curl -v -s -X 'PUT' \
  -d '{"level": "info"}' \
  -H "Content-Type: application/json" \
  localhost:8000/log/level | jq .
