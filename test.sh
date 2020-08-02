#! /bin/bash
set -euo pipefail

go test -v ./cmd/web
go test -v ./pkg/models/pg
#go test -v -short ./pkg/models/pg