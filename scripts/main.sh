#!/bin/bash

./down.sh

go mod tidy
go run ../cmd/main.go