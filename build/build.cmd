@echo off

set OUT_DIR=%~dp0..\bin
set CMD_DIR=%~dp0..\cmd

echo Build into %OUT_DIR%
echo Sources from %CMD_DIR%

go build -o %OUT_DIR%\add.exe %CMD_DIR%\add\main.go
go build -o %OUT_DIR%\metadata-db-create.exe %CMD_DIR%\metadata-db-create\main.go
go build -o %OUT_DIR%\rest-server.exe %CMD_DIR%\rest-server\main.go
go build -o %OUT_DIR%\spa-server.exe %CMD_DIR%\spa-server\main.go
go build -o %OUT_DIR%\ssh-server.exe %CMD_DIR%\ssh-server\main.go