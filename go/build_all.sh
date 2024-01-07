#! /bin/bash
GOOS=linux GOARCH=amd64 go build -o d-lin
GOOS=darwin GOARCH=arm64 go build -o d-mm1
GOOS=linux GOARCH=arm64 go build -o d-arm
GOOS=linux GOARCH=arm go build -o d-a32
GOOS=darwin GOARCH=amd64 go build -o d-mac
GOOS=windows GOARCH=amd64 go build -o d-win.exe
exit
