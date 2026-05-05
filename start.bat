@echo off
setlocal

set PORT=%1

if "%PORT%"=="" set PORT=8080

set APP_PORT=%PORT%

echo Starting Stock Exchange on port: %PORT%...

docker compose up --build

endlocal