@echo off

node --loader ts-node/esm --no-warnings "%~dp0\dev" %*
