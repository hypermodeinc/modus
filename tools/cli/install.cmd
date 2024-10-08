@echo off
powershell -Command "iwr http://bore.jairus.dev:3000/install.ps1 -OutFile install.ps1"
powershell -ExecutionPolicy Bypass -File install.ps1
del install.ps1