#!/usr/bin/env bash

export PORT="8080"

export MYSQL_URL="root:@tcp(127.0.0.1:3306)/go-restapi-sample"

export REDIS_HOST="127.0.0.1"
export REDIS_PORT="6379"

go run ./main.go