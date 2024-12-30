#!/usr/bin/env sh

set -e

cd /workspace/scripts/download-csv
go run .
cd /workspace

go run ./cmd/brotato-unlocker-init
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o brotato-unlocker.exe ./cmd/brotato-unlocker/main.go