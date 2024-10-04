@echo off
SET "PROJECTDIR=%~dp0"
pushd ..\..\..\..\sdk\go\tools\modus-go-build > nul
go run . "%PROJECTDIR%"
popd > nul
