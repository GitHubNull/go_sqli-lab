@echo off
rem Go SQLi Lab - Windows double-click launcher
rem First run will compile and init the database, then start the server and open the browser.
rem Set SQLI_LAB_NO_BROWSER=1 to skip auto-opening the browser.
setlocal
cd /d "%~dp0.."
title Go SQLi Lab - Quick Start

echo ==============================================
echo   Go SQLi Lab Quick Start
echo ==============================================

rem ---- check Go ----
where go >nul 2>nul
if errorlevel 1 (
    echo [ERROR] Go not found.
    echo Please install Go 1.21+ and make sure go is in PATH.
    echo Download: https://go.dev/dl/
    pause
    exit /b 1
)

rem ---- build ----
if not exist bin mkdir bin
echo [1/3] Building...
go build -o bin\go-sqli-lab src\main.go
if errorlevel 1 (
    echo [ERROR] Build failed. Check the Go environment or code.
    pause
    exit /b 1
)

rem ---- init database (idempotent) ----
echo [2/3] Initializing database...
bin\go-sqli-lab --setup-db
if errorlevel 1 (
    echo [ERROR] Database init failed. See the log above.
    pause
    exit /b 1
)

set "PORT=8080"

echo [3/3] Starting server: http://localhost:%PORT%
echo The browser will open automatically. Press Ctrl+C to stop.

if not defined SQLI_LAB_NO_BROWSER (
    start "" "http://localhost:%PORT%"
)

bin\go-sqli-lab

echo.
echo Server stopped.
pause
endlocal
