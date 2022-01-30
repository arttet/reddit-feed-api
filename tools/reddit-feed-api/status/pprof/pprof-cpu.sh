#!/bin/bash

go tool pprof http://localhost:8000/debug/pprof/profile?seconds=30
