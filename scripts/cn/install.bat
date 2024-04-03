@echo off
rem Determine architecture
FOR /F "tokens=2 delims=:," %%F IN ('systeminfo ^| findstr /C:"System Type" ^| findstr x64') DO (
    SET ARCH=amd64
)
if not defined ARCH SET ARCH=arm64

rem Determine the latest release tag
FOR /F "tokens=2 delims=:," %%F IN ('curl -s https://api.github.com/repos/nevalang/neva/releases/latest ^| findstr /C:"\"tag_name\""' ) DO (
    SET LATEST_TAG=%%~F
)
set LATEST_TAG=%LATEST_TAG:"=%
set "LATEST_TAG=%LATEST_TAG: =%"

rem Build the release url
set  BIN_NAME=neva
set "BIN_URL=https://github.com/nevalang/neva/releases/download/%LATEST_TAG%/%BIN_NAME%-windows-%ARCH%.exe"
ECHO Downloading...

rem Download the binary
curl -L %BIN_URL% -o %BIN_NAME%.exe > NUL 2>&1

rem Move the binary to a location in the user's path
set "destination_path=C:\Windows"
move %BIN_NAME%.exe %destination_path%
Echo Installed successfully