#!/bin/bash

rm -f app
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app cmd/web/*
docker build -t docker.pkg.github.com/alvox/flcrd-web/flcrd-api:latest .
#docker push docker.pkg.github.com/alvox/flcrd-web/flcrd-api:latest
rm -f app
docker rmi $(docker images --filter dangling=true -q)