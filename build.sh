#!/bin/sh

GOOS=windows GOARCH=amd64 go build -o bin/planner-amd64-windows.exe main.go
GOOS=darwin GOARCH=amd64 go build -o bin/planner-amd64-mac main.go
GOOS=linux GOARCH=amd64 go build -o bin/planner-amd64-linux main.go
