@echo off

set OUT_DIR=%~dp0..\bin
set CMD_DIR=%~dp0..\cmd

echo Build into %OUT_DIR%
echo Sources from %CMD_DIR%

go build -o %OUT_DIR%\add.exe %CMD_DIR%\add\main.go
go build -o %OUT_DIR%\metadata-db-create.exe %CMD_DIR%\metadata_db_create\main.go
go build -o %OUT_DIR%\rest-server.exe %CMD_DIR%\rest_server\main.go
go build -o %OUT_DIR%\spa-server.exe %CMD_DIR%\spa_server\main.go
go build -o %OUT_DIR%\ssh-server.exe %CMD_DIR%\ssh_server\main.go
go build -o %OUT_DIR%\user-edit.exe %CMD_DIR%\user_edit\main.go
go build -o %OUT_DIR%\faces.exe %CMD_DIR%\faces\main.go