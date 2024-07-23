#!/bin/bash

go build cmd/cli.go
./cli "$@"
rm cli