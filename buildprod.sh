#! /bin/bash

# Build web and other services

cd D:/goPath/src/video_server/api
env GOOS=linux GOARCH=amd64 go build -o ../bin/api

cd D:/goPath/src/video_server/scheduler
env GOOS=linux GOARCH=amd64 go build -o ../bin/scheduler

cd D:/goPath/src/video_server/streamserver
env GOOS=linux GOARCH=amd64 go build -o ../bin/streamserver

cd D:/goPath/src/video_server/web
env GOOS=linux GOARCH=amd64 go build -o ../bin/web