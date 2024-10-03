@echo off
SET "PROJECTDIR=%~dp0"
pushd ..\..\..\..\..\functions-go\tools\hypbuild > nul
go run . "%PROJECTDIR%"
popd > nul
