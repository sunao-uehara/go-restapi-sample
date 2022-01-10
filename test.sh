#!/usr/bin/env bash

go test -v -coverpkg=./... -coverprofile=profile.cov ./...
go tool cover -html=profile.cov
