#!/bin/bash

curl -v -s -X 'GET' localhost:8000/swagger/api.swagger.json | jq .
