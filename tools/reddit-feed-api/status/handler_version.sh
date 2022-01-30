#!/bin/bash

curl -v -s -k -X 'GET' localhost:8000/version | jq .
