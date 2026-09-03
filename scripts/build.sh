#!/bin/bash

BASE_DIR=$(realpath "$(dirname "$0")")
OUT_DIR=$BASE_DIR/../bin
CMD_DIR=$BASE_DIR/../cmd

echo Build into $OUT_DIR
echo Sources from $CMD_DIR

mkdir -p $OUT_DIR

go build -o $OUT_DIR/add $CMD_DIR/add/main.go
go build -o $OUT_DIR/remove $CMD_DIR/remove/main.go
go build -o $OUT_DIR/metadata-db-create $CMD_DIR/metadata_db_create/main.go
go build -o $OUT_DIR/rest-server $CMD_DIR/rest_server/main.go
go build -o $OUT_DIR/spa-server $CMD_DIR/spa_server/main.go
go build -o $OUT_DIR/ssh-server $CMD_DIR/ssh_server/main.go
go build -o $OUT_DIR/user-edit $CMD_DIR/user_edit/main.go
go build -o $OUT_DIR/faces $CMD_DIR/faces/main.go