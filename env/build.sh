#!/bin/bash

rm -f app
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ../cmd/web/*
docker build -t registry.gitlab.com/alvox-env/registry/flcrd-api:latest .
docker push registry.gitlab.com/alvox-env/registry/flcrd-api:latest
rm -f app