#!/bin/bash

curl -v -s -X 'GET' localhost:8000/log/level | jq .
