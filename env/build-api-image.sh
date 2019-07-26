#!/bin/bash

rm -f app
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ../cmd/web/*
docker build . -t flcrd-api:latest -f Dockerfile-api