@echo off
setlocal EnableDelayedExpansion

:: Properties
set DIR=%~dp0
set OS=win

:: Config
set NODE_WANTED=22.9.0
set NODE_MIN=22
set NODE_INSTALL_PATH=%DIR%node-bin
set ARCH=%PROCESSOR_ARCHITECTURE%

set "BASE_DIR=%DIR:~0,-4%"
set "PATH_FILE=%DIR%scripts\modus-path"
set "LOCAL_FILE=%DIR%scripts\modus-local"

:: Determine architecture
if /I "%ARCH%"=="AMD64" (
    set ARCH=x64
) else if /I "%ARCH%"=="ARM" (
    set ARCH=arm
) else if /I "%ARCH%"=="ARM64" (
    set ARCH=arm64
) else (
    echo Unsupported architecture: %ARCH%
    exit /b 1
)


:: Entry point
call :check_node
if errorlevel 1 (
    ::if not exist "%LOCAL_FILE%" (
        mkdir "%DIR%scripts\" >nul 2>&1
        (
            echo import { execute } from "@oclif/core";
            echo await execute({ dir: import.meta.url }^);
        ) > "%LOCAL_FILE%.mjs"
        (
            echo @echo off
            echo %NODE_INSTALL_PATH%\bin\node.exe "%%~dp0modus-local.mjs" %*
        ) > "%LOCAL_FILE%.cmd"
    ::)
    call :install_deps
    cls
    %~dp0scripts\modus-local.cmd %*
    exit
) else (
    if not exist "%PATH_FILE%" (
        mkdir "%DIR%scripts\" >nul 2>&1
        (
            echo import { execute } from "@oclif/core";
            echo await execute({ dir: import.meta.url }^);
        ) > "%PATH_FILE%.mjs"
        (
            echo @echo off
            echo node.exe "%%~dp0modus-path.mjs" %%*
        ) > "%PATH_FILE%.cmd"
    )
    call :install_deps
    cls
    %~dp0scripts\modus-path.cmd %*
    exit
)

endlocal

:: Install Node
:install_node
    set DOWNLOAD_FILE=node-v%NODE_WANTED%-%OS%-%ARCH%.zip
    set DOWNLOAD_URL=https://nodejs.org/dist/v%NODE_WANTED%/%DOWNLOAD_FILE%
    echo Installing Node from %DOWNLOAD_URL%

    set TEMP_DIR=%BASE_DIR%.modus-temp
    if exist "%TEMP_DIR%" (
        echo Removing existing temporary directory...
        rmdir /s /q "%TEMP_DIR%"
    )
    mkdir "%TEMP_DIR%" || (echo Failed to create temporary directory & exit /b 1)

    :: Download Node.js
    curl --progress-bar --show-error --location --fail "%DOWNLOAD_URL%" --output "%TEMP_DIR%\%DOWNLOAD_FILE%"
    if errorlevel 1 (
        echo Failed to download Node.js from %DOWNLOAD_URL%
        exit /b 1
    )

    tar -xf "%TEMP_DIR%\%DOWNLOAD_FILE%" -C "%TEMP_DIR%"
    if errorlevel 1 (
        echo Failed to unpack Node.js archive
        exit /b 1
    )

    rmdir /s /q "%NODE_INSTALL_PATH%\bin"
    mkdir "%NODE_INSTALL_PATH%\bin"
    if errorlevel 1 (
        echo Failed to create Node.js installation directory
        exit /b 1
    )
    del "%TEMP_DIR%\%DOWNLOAD_FILE%"
    xcopy "%TEMP_DIR%\node-v%NODE_WANTED%-%OS%-%ARCH%\*" "%NODE_INSTALL_PATH%\bin" /E /I >nul 2>&1
    rmdir /s /q "%TEMP_DIR%\node-v%NODE_WANTED%-%OS%-%ARCH%"
    if errorlevel 1 (
        echo Failed to move Node.js files to installation directory
        exit /b 1
    )
    rmdir /s /q "%TEMP_DIR%"
exit /b

:: Check Node.js version
:check_node
where node >nul 2>&1

if errorlevel 1 (
    if exist "%NODE_INSTALL_PATH%\bin\node.exe" (
        call :install_node
        exit /b 1
    )
) else (
    if not exist "%NODE_INSTALL_PATH%\bin\node.exe" (
        for /f "delims=v" %%v in ('node -v 2^>nul') do (
            set "NODE_VERSION=%%v"
            for /f "tokens=1,2 delims=." %%a in ("!NODE_VERSION!") do (
                set "MAJOR_VERSION=%%a"
                set "MINOR_VERSION=%%b"
                
                if not !MAJOR_VERSION! geq 22 (
                    call :install_node
                    exit /b 1
                )
            )
        )
    )
)
exit /b

:: Install dependencies
:install_deps
    if not exist "%BASE_DIR%node_modules" (
        echo Installing dependencies with npm
        where npm >nul 2>&1
        if errorlevel 1 (
            cd %BASE_DIR%
            "%NODE_INSTALL_PATH%\bin\npm" install >nul 2>&1
        ) else (
            cd %BASE_DIR%
            npm install >nul 2>&1
        )
    )
exit /b