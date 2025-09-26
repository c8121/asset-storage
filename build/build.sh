#!/bin/bash

BASE_DIR=$(realpath "$(dirname "$0")")
OUT_DIR=$BASE_DIR/../bin
CMD_DIR=$BASE_DIR/../cmd

echo Build into $OUT_DIR
echo Sources from $CMD_DIR

mkdir -p $OUT_DIR

go build -o $OUT_DIR/add $CMD_DIR/add/main.go
go build -o $OUT_DIR/metadata-db-create $CMD_DIR/metadata-sqlite-create/main.go
go build -o $OUT_DIR/rest-server $CMD_DIR/rest-server/main.go
go build -o $OUT_DIR/spa-server $CMD_DIR/spa-server/main.go