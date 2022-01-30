#!/bin/bash

go tool pprof -alloc_space http://localhost:8000/debug/pprof/heap
