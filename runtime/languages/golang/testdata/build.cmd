@echo off
SET "PROJECTDIR=%~dp0"
pushd ..\..\..\..\sdk\go\tools\hypbuild > nul
go run . "%PROJECTDIR%"
popd > nul
