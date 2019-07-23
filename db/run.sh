#!/bin/bash

docker run --name flcrd-db \
  -e POSTGRES_DB=flcrd \
  -e POSTGRES_USER=flcrd \
  -e POSTGRES_PASSWORD=flcrd \
  -p 5432:5432 \
  -v /tmp/flcrd-db:/var/lib/postgresql/data \
  flcrd-db:latest