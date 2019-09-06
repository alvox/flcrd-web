#!/bin/bash

go build -o app ../cmd/web/*
./app -port=":5000" -dsn="postgres://flcrd:flcrd@localhost/flcrd?sslmode=disable" -appkey="myappkey"