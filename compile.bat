@echo off
rem Compile script for go-manga-parser
rem Usage: double-click this file or run from PowerShell/CMD

setlocal
cd /d "%~dp0"

echo Building manga-parser...
go build -o manga-parser.tmp.exe ./cmd/parser
if %ERRORLEVEL% NEQ 0 (
  echo Build failed with error %ERRORLEVEL%.
  pause
  exit /b %ERRORLEVEL%
)

rem Move the temp binary into place atomically (overwrite existing)
move /Y manga-parser.tmp.exe manga-parser.exe >nul
if %ERRORLEVEL% NEQ 0 (
  echo Failed to move compiled binary into place.
  pause
  exit /b %ERRORLEVEL%
)

echo Build succeeded. Updated manga-parser.exe
pause
endlocal
