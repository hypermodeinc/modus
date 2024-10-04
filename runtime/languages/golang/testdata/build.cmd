@echo off
SET "PROJECTDIR=%~dp0"
pushd ..\..\..\..\sdk\go\tools\modus-go-build > nul
go run . "%PROJECTDIR%"
set "exit_code=%ERRORLEVEL%"
popd > nul
exit /b %exit_code%
