#!/bin/bash
env GOOS=linux go build main.go
docker build -t testapi --no-cache .
