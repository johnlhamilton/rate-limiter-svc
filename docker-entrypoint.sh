#!/usr/bin/env bash
go mod tidy
go mod vendor
go run main.go config.go api.go