#!/bin/bash
go test -coverprofile=coverage.tmp ./...

grep -v "internal/server/db/" coverage.tmp |
	grep -v "internal/server/mocks/" |
	grep -v "swagger/" |
	grep -v "cmd/" |
	grep -v "pkg/proto/" >coverage.out

go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html

rm coverage.tmp
